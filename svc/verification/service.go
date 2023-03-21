package verification

import "context"

type (
	// User verification service
	Service interface {
		// Verify user email with otp code
		VerifyEmail(ctx context.Context, email, code string) error
		// Resend verification email
		ResendEmail(ctx context.Context, email string) error
	}

	service struct {
		repo     verificationRepository
		enqueuer emailEnqueuer
	}

	verificationRepository interface{}

	emailEnqueuer interface{}
)

// NewService creates a new user verification service
func NewService(repo verificationRepository, enqueuer emailEnqueuer) Service {
	return &service{repo: repo, enqueuer: enqueuer}
}

// Verify user email with otp code
func (s *service) VerifyEmail(ctx context.Context, email, code string) error {
	return nil
}

// Resend verification email
func (s *service) ResendEmail(ctx context.Context, email string) error {
	return nil
}
