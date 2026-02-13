package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) error {
	m, err := migrate.New("file:///Users/ritiksrivastava/Documents/Coding/Go-Learning/login/internal/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("migration failed to init: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	m.Up()
	return nil
}
