package service

import (
	"context"
	"login/internal/repository"
	"testing"
)

type MockRepository struct {
	repository.UserRepository
	SavedUser *repository.User
}

// FindByEmail simulates looking up a user
func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	// If the email is "exists@example.com", we pretend the user is in the DB
	if email == "exists@example.com" {
		return &repository.User{Email: email}, nil
	}
	return nil, nil // Otherwise, user not found
}

func (m *MockRepository) CreateUser(ctx context.Context, user *repository.User) error {
	m.SavedUser = user // Capture the user for the test to check
	return nil
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	// 1. Setup our service with the mock database
	mockRepo := &MockRepository{}
	svc := &UserService{repo: mockRepo}

	// 2. Try to register an email that our mock says ALREADY exists
	_, err := svc.Register(context.Background(), "gopher", "exists@example.com", "password123")

	if err == nil {
		t.Errorf("expected error for duplicate email, but got nil")
	}
}

func TestRegister_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := &UserService{repo: mockRepo}
	rawPassword := "secret123"

	_, err := svc.Register(context.Background(), "gopher", "new@example.com", rawPassword)

	if err != nil {
		t.Fatalf("Expected successful registration, got error: %v", err)
	}

	if mockRepo.SavedUser.Password == rawPassword {
		t.Errorf("Raw password is not hashed")
	}
}
