package cmd

import (
	"context"
	"fmt"
	"login/internal/auth"
	"login/internal/repository"
	"login/internal/service"
)

func main() {
	// 1. Initialize Repository (Postgres)
	ctx := context.Background()
	connStr := "postgres://user:password@localhost:5432/gopher_db?sslmode=disable"

	pool, err := repository.NewPostgresPool(ctx, connStr)
	if err != nil {
		fmt.Errorf("Postgres setup fail: %v", err)
	}
	if err := repository.RunMigrations(connStr); err != nil {
		fmt.Errorf("Failed to run migrations: %v", err)
	}

	repo := repository.NewPostgresUserRepository(pool)

	// 2. Initialize Rate Limiter (Redis)
	// You must provide the address of your Redis container
	limiter := repository.NewRateLimiter("localhost:6379")

	// 3. Initialize Token Manager (JWT)
	// In production, load the secret from an Environment Variable!
	tokenManager := auth.NewTokenManager("your-super-secret-key", "gopher-service")

	// 4. Inject everything into the Service
	userService := service.NewUserService(repo, limiter, tokenManager)

	userService.Login(ctx, "m", "","e", "")
	// Now userService.Login() will work because s.limiter and s.tokenManager are not nil!
}
