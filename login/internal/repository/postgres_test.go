package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestNewPostgresPool(t *testing.T) {
	conn := "postgres://user:password@localhost:5432/gopher_db?sslmode=disable"

	db, err := NewPostgresPool(context.Background(), conn)

	fmt.Print(db, err)
	if err != nil {
		t.Errorf("Database connection failed")
	}

	defer db.Close()

	t.Log("Successfully connected to Postgres!")
}
