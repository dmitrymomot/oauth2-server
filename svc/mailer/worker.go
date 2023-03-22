package mailer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

type (
	// Worker is a task handler for email delivery.
	Worker struct {
		mail mailer
	}

	// WorkerOption is a function that configures Worker.
	WorkerOption func(*Worker)

	mailer interface {
		SendVerificationCode(ctx context.Context, uid, email, otp string) error
		SendResetPasswordCode(ctx context.Context, uid, email, otp string) error
		SendDestroyProfileCode(ctx context.Context, uid, email, otp string) error
	}
)

// NewWorker creates a new email task handler.
func NewWorker(mail mailer) *Worker {
	return &Worker{
		mail: mail,
	}
}

// Register registers task handlers for email delivery.
func (w *Worker) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(SendConfirmationEmailTask, w.TaskSendConfirmationEmail)
	mux.HandleFunc(SendPasswordRecoveryEmailTask, w.TaskSendPasswordResetEmail)
	mux.HandleFunc(SendDestroyProfileEmailTask, w.TaskSendDestroyProfileEmail)
}

// TaskSendConfirmationEmail sends confirmation email to user
// after successful registration or email change.
func (w *Worker) TaskSendConfirmationEmail(ctx context.Context, t *asynq.Task) error {
	var p ConfirmationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := w.mail.SendVerificationCode(ctx, p.UserID, p.Email, p.OTP); err != nil {
		return errors.Wrap(err, "failed to send email with email address verification code")
	}

	return nil
}

// TaskSendPasswordResetEmail sends password reset email to user.
func (w *Worker) TaskSendPasswordResetEmail(ctx context.Context, t *asynq.Task) error {
	var p ConfirmationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := w.mail.SendResetPasswordCode(ctx, p.UserID, p.Email, p.OTP); err != nil {
		return errors.Wrap(err, "failed to send email with password reset code")
	}

	return nil
}

// TaskSendDestroyProfileEmail sends email to confirm profile destruction.
func (w *Worker) TaskSendDestroyProfileEmail(ctx context.Context, t *asynq.Task) error {
	var p ConfirmationEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := w.mail.SendDestroyProfileCode(ctx, p.UserID, p.Email, p.OTP); err != nil {
		return errors.Wrap(err, "failed to send email with destroy account confirmation code")
	}

	return nil
}
