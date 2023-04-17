package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		// Login authenticates a user and returns a user ID.
		Login(ctx context.Context, email, password string) (uuid.UUID, error)
		// Register creates a new user and returns a user ID.
		Register(ctx context.Context, email, password string) (uuid.UUID, error)
		// PasswordRecovery sends a password recovery email.
		PasswordRecovery(ctx context.Context, email string) error
		// PasswordReset resets a user password.
		PasswordReset(ctx context.Context, email, otp, password string) error
		// Verify user email with otp code
		VerifyEmail(ctx context.Context, email, code string) error
		// Resend verification email
		ResendVerificationEmail(ctx context.Context, email string) error
		// DestroyProfileRequest sends a destroy profile email.
		DestroyProfileRequest(ctx context.Context, email string) error
		// DestroyProfile destroys a user profile.
		DestroyProfile(ctx context.Context, email, otp string) error
	}

	service struct {
		repo authRepository
		db   *sql.DB
		mail mailer
	}

	authRepository interface {
		WithTx(tx *sql.Tx) *repository.Queries
		GetUserByID(ctx context.Context, id uuid.UUID) (repository.User, error)
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
		UpdateUserPassword(ctx context.Context, arg repository.UpdateUserPasswordParams) (repository.User, error)
		UpdateUserVerifiedAt(ctx context.Context, id uuid.UUID) error
		DeleteUser(ctx context.Context, id uuid.UUID) error

		CreateUserVerification(ctx context.Context, arg repository.CreateUserVerificationParams) error
		GetUserVerificationByEmail(ctx context.Context, arg repository.GetUserVerificationByEmailParams) (repository.UserVerification, error)
		DeleteUserVerificationsByEmail(ctx context.Context, arg repository.DeleteUserVerificationsByEmailParams) error
		DeleteUserVerificationsByUserID(ctx context.Context, arg repository.DeleteUserVerificationsByUserIDParams) error
	}

	mailer interface {
		SendConfirmationEmail(ctx context.Context, uid uuid.UUID, email, otp string) error
		SendPasswordRecoveryEmail(ctx context.Context, uid uuid.UUID, email, otp string) error
		SendDestroyProfileEmail(ctx context.Context, uid uuid.UUID, email, otp string) error
	}
)

// NewService creates a new auth service.
func NewService(repo authRepository, db *sql.DB, m mailer) Service {
	return &service{
		repo: repo,
		db:   db,
		mail: m,
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

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if !user.VerifiedAt.Valid {
		return uuid.Nil, ErrUserNotVerified
	}

	return user.ID, nil
}

// Register creates a new user and returns a user ID.
func (s *service) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	email = strings.TrimSpace(strings.ToLower(email))
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

	tx, err := s.db.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	user, err = repo.CreateUser(ctx, repository.CreateUserParams{
		Email:    email,
		Password: passwordHash,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	otp := random.String(6, random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate otp hash: %w", err)
	}

	if err := repo.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.UserVerificationRequestTypeEmailVerification,
		UserID:           user.ID,
		Email:            user.Email,
		VerificationCode: otpHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user verification: %w", err)
	}

	if err := s.mail.SendConfirmationEmail(ctx, user.ID, user.Email, otp); err != nil {
		return uuid.Nil, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user.ID, nil
}

// PasswordRecovery sends a password recovery email.
func (s *service) PasswordRecovery(ctx context.Context, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	otp := random.String(6, random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate otp hash: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	if err := repo.DeleteUserVerificationsByEmail(ctx, repository.DeleteUserVerificationsByEmailParams{
		RequestType: repository.UserVerificationRequestTypePasswordReset,
		Email:       user.Email,
	}); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete user verifications by email: %w", err)
		}
	}

	if err := repo.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.UserVerificationRequestTypePasswordReset,
		UserID:           user.ID,
		Email:            user.Email,
		VerificationCode: otpHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}); err != nil {
		return fmt.Errorf("failed to create user verification: %w", err)
	}

	if err := s.mail.SendPasswordRecoveryEmail(ctx, user.ID, user.Email, otp); err != nil {
		return fmt.Errorf("failed to send password recovery email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// PasswordReset resets a user password.
func (s *service) PasswordReset(ctx context.Context, email, otp, password string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	uv, err := s.repo.GetUserVerificationByEmail(ctx, repository.GetUserVerificationByEmailParams{
		RequestType: repository.UserVerificationRequestTypePasswordReset,
		Email:       email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidVerificationRequest
		}
		return fmt.Errorf("failed to get user verification by email: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(uv.VerificationCode, []byte(otp)); err != nil {
		return ErrInvalidVerificationCode
	}
	if time.Now().After(uv.ExpiresAt) {
		return ErrVerificationCodeExpired
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	user, err := repo.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       uv.UserID,
		Password: passwordHash,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user password: %w", err)
	}

	if err := repo.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.UserVerificationRequestTypePasswordReset,
		UserID:      user.ID,
	}); err != nil {
		return fmt.Errorf("failed to delete user verifications by user id: %w", err)
	}

	if !user.VerifiedAt.Valid {
		if err := repo.UpdateUserVerifiedAt(ctx, user.ID); err != nil {
			return fmt.Errorf("failed to update user verified at: %w", err)
		}

		if err := repo.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
			RequestType: repository.UserVerificationRequestTypeEmailVerification,
			UserID:      user.ID,
		}); err != nil {
			return fmt.Errorf("failed to delete user verifications by user id: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Verify user email with otp code
func (s *service) VerifyEmail(ctx context.Context, email, otp string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	uv, err := s.repo.GetUserVerificationByEmail(ctx, repository.GetUserVerificationByEmailParams{
		RequestType: repository.UserVerificationRequestTypeEmailVerification,
		Email:       email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidVerificationRequest
		}
		return fmt.Errorf("failed to get user verification by email: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(uv.VerificationCode, []byte(otp)); err != nil {
		return ErrInvalidVerificationCode
	}
	if time.Now().After(uv.ExpiresAt) {
		return ErrVerificationCodeExpired
	}

	user, err := s.repo.GetUserByID(ctx, uv.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by id: %w", err)
	}
	if user.VerifiedAt.Valid {
		return ErrUserAlreadyVerified
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.repo.WithTx(tx)

	if err := repo.UpdateUserVerifiedAt(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to update user verified at: %w", err)
	}

	if err := repo.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.UserVerificationRequestTypeEmailVerification,
		UserID:      user.ID,
	}); err != nil {
		return fmt.Errorf("failed to delete user verifications by user id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Resend verification email
func (s *service) ResendVerificationEmail(ctx context.Context, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	if user.VerifiedAt.Valid {
		return ErrUserAlreadyVerified
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
		RequestType:      repository.UserVerificationRequestTypeEmailVerification,
		UserID:           user.ID,
		Email:            user.Email,
		VerificationCode: otpHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}); err != nil {
		return fmt.Errorf("failed to create user verification: %w", err)
	}

	if err := s.mail.SendConfirmationEmail(ctx, user.ID, user.Email, otp); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DestroyProfileRequest sends a destroy profile email.
func (s *service) DestroyProfileRequest(ctx context.Context, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by email: %w", err)
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
		return fmt.Errorf("failed to send destroy profile email: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DestroyProfile destroys a user profile.
func (s *service) DestroyProfile(ctx context.Context, email, otp string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	uv, err := s.repo.GetUserVerificationByEmail(ctx, repository.GetUserVerificationByEmailParams{
		RequestType: repository.UserVerificationRequestTypeDeleteAccount,
		Email:       email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidVerificationRequest
		}
		return fmt.Errorf("failed to get user verification by email: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(uv.VerificationCode, []byte(otp)); err != nil {
		return ErrInvalidVerificationCode
	}
	if time.Now().After(uv.ExpiresAt) {
		return ErrVerificationCodeExpired
	}

	user, err := s.repo.GetUserByID(ctx, uv.UserID)
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

	if err := repo.DeleteUser(ctx, user.ID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete user by id: %w", err)
		}
	}

	if err := repo.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.UserVerificationRequestTypeDeleteAccount,
		UserID:      user.ID,
	}); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete user verifications by user id: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
