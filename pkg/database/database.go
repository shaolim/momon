package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*DB, error) {
	dbpool, err := pgxpool.New(ctx, cfg.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{
		Pool: dbpool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) InTx(ctx context.Context, isoLevel pgx.TxIsoLevel, f func(tx pgx.Tx) error) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := f(tx); err != nil {
		if err1 := tx.Rollback(ctx); err1 != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", err1, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
