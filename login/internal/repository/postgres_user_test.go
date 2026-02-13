package repository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestPostgresUserRepository(t *testing.T) {
	ctx := context.Background()
	// Get this from your .env or hardcode for local docker testing
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/gopher_db?sslmode=disable"
	}

	// 1. Setup connection
	pool, err := NewPostgresPool(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	// 2. Run migrations
	if err := RunMigrations(connStr); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	repo := NewPostgresUserRepository(pool)

	t.Run("Create and Find User", func(t *testing.T) {
		userID := uuid.New().String()
		email := "tddq@example.com"
		user := &User{
			ID:       userID,
			Username: "ritik",
			Email:    email,
			Password: "hashed_password",
		}

		// Test Create
		err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		// Test Find
		found, err := repo.FindByEmail(ctx, email)
		if err != nil {
			t.Fatalf("FindByEmail failed: %v", err)
		}

		if found == nil {
			t.Fatal("Expected to find user, got nil")
		}

		if found.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, found.Username)
		}
	})
}
