package repository

import "context"

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
