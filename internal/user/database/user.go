package database

import (
	"context"
	"fmt"
	"github/shaolim/momon/internal/user/model"
	"github/shaolim/momon/pkg/database"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserDB interface {
	AddUser(ctx context.Context, user *model.User) error
}

type userDB struct {
	db *database.DB
}

func New(db *database.DB) UserDB {
	return &userDB{
		db: db,
	}
}

func (db *userDB) AddUser(ctx context.Context, user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	if err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, `
			INSERT INTO users (line_user_id, display_name, status, created_at, updated_at)
			VALUES($1, $2, $3, $4, $5)
			RETURNING id
		`, user.LineUserID, user.DisplayName, user.Status, user.CreatedAt, user.UpdatedAt)

		if err := row.Scan(&user.ID); err != nil {
			return fmt.Errorf("insert users: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
