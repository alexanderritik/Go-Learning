package service

import (
	"context"
	"errors"
	"time"

	"login/internal/auth"
	"login/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.UserRepository
	limiter      *repository.RateLimiter
	tokenManager auth.TokenGenerator
}

func NewUserService(
	repo repository.UserRepository,
	limiter *repository.RateLimiter,
	tm auth.TokenGenerator,
) *UserService {
	return &UserService{
		repo:         repo,
		limiter:      limiter,
		tokenManager: tm,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (string, error) {

	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if existingUser != nil {
		return "", errors.New("user already exists")
	}

	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user_id := uuid.NewString()
	newUser := &repository.User{
		ID:       user_id,
		Username: username,
		Email:    email,
		Password: string(hashedBytes),
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return "", err
	}
	return user_id, nil
}

func (s *UserService) Login(ctx context.Context, username, email, password, ip string) (string, error) {
	limitKey := "limit:login:" + ip

	allowed, err := s.limiter.Allow(ctx, limitKey, 5, time.Minute)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", errors.New("too many login attempts: blocked for 1 minute")
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := s.tokenManager.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
