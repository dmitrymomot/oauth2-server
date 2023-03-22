package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		// GetByID returns the user with the specified the user ID.
		GetByID(ctx context.Context, id string) (*User, error)
		// UpdateEmail updates the email address of the user with the specified ID.
		UpdateEmail(ctx context.Context, id, email string) (*User, error)
		// UpdatePassword updates the password of the user with the specified ID.
		UpdatePassword(ctx context.Context, id, oldPassword, newPassword string) error
		// Delete deletes the user with the specified ID.
		Delete(ctx context.Context, id string) error
	}

	User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Verified  bool   `json:"verified"`
		CreatedAt string `json:"created_at"`
	}

	service struct {
		repo userRepository
		mail mailer
		db   *sql.DB
	}

	userRepository interface {
		WithTx(tx *sql.Tx) *repository.Queries
		GetUserByID(ctx context.Context, id uuid.UUID) (repository.User, error)
		UpdateUserEmail(ctx context.Context, arg repository.UpdateUserEmailParams) (repository.User, error)
		UpdateUserPassword(ctx context.Context, arg repository.UpdateUserPasswordParams) (repository.User, error)
		DeleteUser(ctx context.Context, id uuid.UUID) error
		CreateUserVerification(ctx context.Context, arg repository.CreateUserVerificationParams) error
	}

	mailer interface {
		SendConfirmationEmail(ctx context.Context, uid uuid.UUID, email, otp string) error
		SendDestroyProfileEmail(ctx context.Context, uid uuid.UUID, email, otp string) error
	}
)

// NewUser casts a repository.User to a user.User.
func NewUser(u repository.User) *User {
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Verified:  u.VerifiedAt.Valid,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// NewService creates a new user service.
// It is the concrete implementation of the Service interface.
func NewService(repo userRepository, m mailer, db *sql.DB) Service {
	return &service{repo: repo, mail: m, db: db}
}

// GetByID returns the user with the specified the user ID.
func (s *service) GetByID(ctx context.Context, id string) (*User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	u, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return NewUser(u), nil
}

// UpdateEmail updates the email address of the user with the specified ID.
func (s *service) UpdateEmail(ctx context.Context, id, email string) (*User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	user, err := repo.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
		ID:    uid,
		Email: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to update user email: %w", err)
	}

	otp := random.String(6, random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate otp hash: %w", err)
	}

	if err := repo.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.UserVerificationRequestTypeEmailVerification,
		UserID:           user.ID,
		Email:            user.Email,
		VerificationCode: otpHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}); err != nil {
		return nil, fmt.Errorf("failed to create user verification: %w", err)
	}

	if err := s.mail.SendConfirmationEmail(ctx, user.ID, user.Email, otp); err != nil {
		return nil, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return NewUser(user), nil
}

// UpdatePassword updates the password of the user with the specified ID.
func (s *service) UpdatePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by id: %w", err)
	}
	if bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)) != nil {
		return ErrInvalidPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	user, err = s.repo.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       uid,
		Password: passwordHash,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user email: %w", err)
	}

	return nil
}

// Delete deletes the user with the specified ID.
func (s *service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	otp := random.String(6, random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate otp hash: %w", err)
	}

	if err := repo.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.UserVerificationRequestTypeDeleteAccount,
		UserID:           user.ID,
		Email:            user.Email,
		VerificationCode: otpHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}); err != nil {
		return fmt.Errorf("failed to create user verification: %w", err)
	}

	if err := s.mail.SendDestroyProfileEmail(ctx, user.ID, user.Email, otp); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
