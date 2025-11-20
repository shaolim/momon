package database

import (
	"context"
	"testing"
)

var testDatabaseInstance *TestInstance

func TestMain(m *testing.M) {
	testDatabaseInstance = MustTestInstance()
	defer testDatabaseInstance.MustClose()
	m.Run()
}

func TestNewDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	db, cfg := testDatabaseInstance.NewDatabase(t)
	if db == nil {
		t.Fatal("expected db to be non-nil")
	}
	if cfg == nil {
		t.Fatal("expected config to be non-nil")
	}

	// Verify the database connection works
	ctx := context.Background()
	if err := db.Pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	// Verify the users table exists (from migrations)
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'users'
	)`
	if err := db.Pool.QueryRow(ctx, query).Scan(&exists); err != nil {
		t.Fatalf("failed to query for users table: %v", err)
	}

	if !exists {
		t.Error("users table does not exist in cloned database")
	}
}
