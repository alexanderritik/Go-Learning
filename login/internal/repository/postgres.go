package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Production tip: Set pool limits (e.g., max connections) here
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Ping the connection to ensure it's actually alive
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
