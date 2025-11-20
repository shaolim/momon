package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const (
	// databaseName is the name of the template database to clone.
	databaseName = "test-db-template"

	// databaseUser and databasePassword are the username and password for
	// connecting to the database. These values are only used for testing.
	databaseUser     = "test-user"
	databasePassword = "testing123"

	defaultPostgresTag = "18-alpine"
)

type TestInstance struct {
	pool      *dockertest.Pool
	container *dockertest.Resource
	url       *url.URL
	conn      *pgx.Conn
	connLock  sync.Mutex
}

func MustTestInstance() *TestInstance {
	testDatabaseInstance, err := NewTestInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return testDatabaseInstance
}

func NewTestInstance() (*TestInstance, error) {
	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to create database docker pool: %w", err)
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        defaultPostgresTag,
		Env: []string{
			"LANG=C",
			"POSTGRES_DB=" + databaseName,
			"POSTGRES_USER=" + databaseUser,
			"POSTGRES_PASSWORD=" + databasePassword,
		},
	}, func(c *docker.HostConfig) {
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start database container: %w", err)
	}

	if err := container.Expire(120); err != nil {
		return nil, fmt.Errorf("failed to expire database container: %w", err)
	}

	connectionURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(databaseUser, databasePassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     databaseName,
		RawQuery: "sslmode=disable",
	}

	var conn *pgx.Conn
	if err := pool.Retry(func() error {
		var err error
		conn, err = pgx.Connect(ctx, connectionURL.String())
		if err != nil {
			return err
		}
		return conn.Ping(ctx)
	}); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := dbMigrate(connectionURL.String()); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &TestInstance{
		pool:      pool,
		container: container,
		conn:      conn,
		url:       connectionURL,
	}, nil
}

// dbMigrate runs the migrations. u is the connection URL string (e.g.
// postgres://...).
func dbMigrate(u string) error {
	// Run the migrations
	migrationsDir := fmt.Sprintf("file://%s", dbMigrationsDir())
	m, err := migrate.New(migrationsDir, u)
	if err != nil {
		return fmt.Errorf("failed create migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database error: %w", dbErr)
	}
	return nil
}

// dbMigrationsDir returns the path on disk to the migrations. It uses
// runtime.Caller() to get the path to the caller, since this package is
// imported by multiple others at different levels.
func dbMigrationsDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), "../../migrations")
}

func (i *TestInstance) MustClose() error {
	if err := i.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return nil
}

func (i *TestInstance) Close() (retErr error) {
	defer func() {
		if err := i.pool.Purge(i.container); err != nil {
			retErr = fmt.Errorf("failed to purge database container: %w", err)
			return
		}
	}()

	ctx := context.Background()
	if err := i.conn.Close(ctx); err != nil {
		retErr = fmt.Errorf("failed to close connection: %w", err)
		return
	}

	return nil
}

func (i *TestInstance) NewDatabase(tb testing.TB) (*DB, *Config) {
	newDatabaseName, err := i.clone()
	if err != nil {
		tb.Fatal(err)
	}

	// Build the new connection URL for the new database name. Query params are
	// dropped with ResolveReference, so we have to re-add disabling SSL over
	// localhost.
	connectionURL := i.url.ResolveReference(&url.URL{Path: newDatabaseName})
	connectionURL.RawQuery = "sslmode=disable"

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, connectionURL.String())
	if err != nil {
		tb.Fatalf("failed to connect to database %q: %s", newDatabaseName, err)
	}

	db := &DB{Pool: dbpool}

	tb.Cleanup(func() {
		ctx := context.Background()

		// Close connection first. It is an error to drop a database with active
		// connections.
		db.Close()

		// Drop the database to keep the container from running out of resources.
		// If the instance connection is already closed, we can skip this since the
		// container will be purged anyway.
		q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, newDatabaseName)

		i.connLock.Lock()
		defer i.connLock.Unlock()

		if i.conn.IsClosed() {
			// Connection is already closed, container will be purged anyway
			return
		}

		if _, err := i.conn.Exec(ctx, q); err != nil {
			tb.Errorf("failed to drop database %q: %s", newDatabaseName, err)
		}
	})

	host, port, err := net.SplitHostPort(i.url.Host)
	if err != nil {
		tb.Fatalf("failed to split host/port %q: %s", i.url.Host, err)
	}

	return db, &Config{
		Name:     newDatabaseName,
		User:     databaseUser,
		Host:     host,
		Port:     port,
		Password: databasePassword,
	}
}

func (i *TestInstance) clone() (string, error) {
	name, err := randomDatabaseName()
	if err != nil {
		return "", fmt.Errorf("failed to generate random database name: %w", err)
	}

	ctx := context.Background()
	q := fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, name, databaseName)

	// Unfortunately postgres does not allow parallel database creation from the
	// same template, so this is guarded with a lock.
	i.connLock.Lock()
	defer i.connLock.Unlock()

	// Clone the template database as the new random database name.
	if _, err := i.conn.Exec(ctx, q); err != nil {
		return "", fmt.Errorf("failed to clone template database: %w", err)
	}
	return name, nil
}

// randomDatabaseName returns a random database name.
func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
