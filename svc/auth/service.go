package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Login(ctx context.Context, email, password string) (uuid.UUID, error)
		Register(ctx context.Context, email, password string) (uuid.UUID, error)
		PasswordRecovery(ctx context.Context, email string) error
		PasswordReset(ctx context.Context, email, otp, password string) error
	}

	service struct {
		repo authRepository
	}

	authRepository interface {
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
	}
)

// NewService creates a new auth service.
func NewService(repo authRepository) Service {
	return &service{
		repo: repo,
	}
}

// Login authenticates a user and returns a user ID.
func (s *service) Login(ctx context.Context, email, password string) (uuid.UUID, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrInvalidCredentials
		}
		return uuid.Nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user.ID, nil
}

// Register creates a new user and returns a user ID.
func (s *service) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}
	if user.ID != uuid.Nil {
		return uuid.Nil, ErrEmailTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	user, err = s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:    strings.TrimSpace(strings.ToLower(email)),
		Password: passwordHash,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

// PasswordRecovery sends a password recovery email.
func (s *service) PasswordRecovery(ctx context.Context, email string) error {
	return nil
}

// PasswordReset resets a user password.
func (s *service) PasswordReset(ctx context.Context, email, otp, password string) error {
	return nil
}
