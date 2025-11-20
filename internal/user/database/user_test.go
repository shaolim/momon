package database

import (
	"context"
	"github/shaolim/momon/internal/user/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		testDB, _ := testDatabaseInstance.NewDatabase(t)
		userDB := New(testDB)
		ctx := context.Background()

		user := &model.User{
			LineUserID:  "line123",
			DisplayName: "surti",
			Status:      model.UserStatusActive,
		}

		err := userDB.AddUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to add user: %v", err)
		}

		// Verify ID was assigned
		if user.ID == 0 {
			t.Error("expected user ID to be assigned")
		}
	})

	t.Run("duplicate_line_user_id", func(t *testing.T) {
		t.Parallel()

		testDB, _ := testDatabaseInstance.NewDatabase(t)
		userDB := New(testDB)
		ctx := context.Background()

		user1 := &model.User{
			LineUserID:  "line456",
			DisplayName: "User One",
			Status:      model.UserStatusActive,
		}

		// Add first user - should succeed
		err := userDB.AddUser(ctx, user1)
		if err != nil {
			t.Fatalf("failed to add first user: %v", err)
		}

		// Try to add second user with same LineUserID - should fail
		user2 := &model.User{
			LineUserID:  "line456",
			DisplayName: "User Two",
			Status:      model.UserStatusActive,
		}

		err = userDB.AddUser(ctx, user2)
		if err == nil {
			t.Fatal("expected error when adding duplicate line_user_id, got nil")
		}

		// Verify it's a unique constraint violation
		if !strings.Contains(err.Error(), "duplicate") && !strings.Contains(err.Error(), "unique") {
			t.Errorf("expected unique constraint error, got: %v", err)
		}
	})

	t.Run("empty_line_user_id", func(t *testing.T) {
		t.Parallel()

		testDB, _ := testDatabaseInstance.NewDatabase(t)
		userDB := New(testDB)
		ctx := context.Background()

		user := &model.User{
			LineUserID:  "",
			DisplayName: "Test User",
			Status:      model.UserStatusActive,
		}

		err := userDB.AddUser(ctx, user)
		assert.EqualError(t, err, "line user id must not be empty")
	})
}
