package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type (
	// Enqueuer is a helper struct for enqueuing email tasks.
	Enqueuer struct {
		client       *asynq.Client
		queueName    string
		taskDeadline time.Duration
		maxRetry     int
		log          Logger
	}

	// EnqueuerOption is a function that configures an enqueuer.
	EnqueuerOption func(*Enqueuer)

	Logger interface {
		Warnf(format string, args ...interface{})
	}
)

// NewEnqueuer creates a new email enqueuer.
// This function accepts EnqueuerOption to configure the enqueuer.
// Default values are used if no option is provided.
// Default values are:
//   - queue name: "default"
//   - task deadline: 1 minute
//   - max retry: 3
func NewEnqueuer(client *asynq.Client, opt ...EnqueuerOption) *Enqueuer {
	if client == nil {
		panic("client is nil")
	}

	e := &Enqueuer{
		client:       client,
		queueName:    "default",
		taskDeadline: time.Minute,
		maxRetry:     3,
	}

	for _, o := range opt {
		o(e)
	}

	return e
}

// WithQueueName configures the queue name.
func WithQueueName(name string) EnqueuerOption {
	return func(e *Enqueuer) {
		e.queueName = name
	}
}

// WithTaskDeadline configures the task deadline.
func WithTaskDeadline(d time.Duration) EnqueuerOption {
	return func(e *Enqueuer) {
		e.taskDeadline = d
	}
}

// WithMaxRetry configures the max retry.
func WithMaxRetry(n int) EnqueuerOption {
	return func(e *Enqueuer) {
		e.maxRetry = n
	}
}

// WithLogger configures the logger.
func WithLogger(l Logger) EnqueuerOption {
	return func(e *Enqueuer) {
		e.log = l
	}
}

// enqueueTask enqueues a task to the queue.
func (e *Enqueuer) enqueueTask(ctx context.Context, task *asynq.Task) error {
	if e.client == nil {
		if e.log != nil {
			e.log.Warnf("client is nil, skipping enqueue task: %s", task.Type())
			return nil
		}
		return fmt.Errorf("client is nil, skipping enqueue task: %s", task.Type())
	}

	if _, err := e.client.Enqueue(
		task,
		asynq.Queue(e.queueName),
		asynq.Deadline(time.Now().Add(e.taskDeadline)),
		asynq.MaxRetry(e.maxRetry),
		asynq.Unique(e.taskDeadline),
	); err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// SendConfirmationEmail sends confirmation email to user.
// This method returns a task to be added to the queue.
func (e *Enqueuer) SendConfirmationEmail(ctx context.Context, uid uuid.UUID, email, otp string) error {
	payload, err := json.Marshal(ConfirmationEmailPayload{
		UserID: uid.String(),
		Email:  email,
		OTP:    otp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return e.enqueueTask(ctx, asynq.NewTask(SendConfirmationEmailTask, payload))
}

// SendPasswordRecoveryEmail sends password reset email to user.
// This method returns a task to be added to the queue.
func (e *Enqueuer) SendPasswordRecoveryEmail(ctx context.Context, uid uuid.UUID, email, otp string) error {
	payload, err := json.Marshal(ConfirmationEmailPayload{
		UserID: uid.String(),
		Email:  email,
		OTP:    otp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return e.enqueueTask(ctx, asynq.NewTask(SendPasswordRecoveryEmailTask, payload))
}

// SendDestroyProfileEmail sends email to confirm profile destruction.
// This method returns a task to be added to the queue.
func (e *Enqueuer) SendDestroyProfileEmail(ctx context.Context, uid uuid.UUID, email, otp string) error {
	payload, err := json.Marshal(ConfirmationEmailPayload{
		UserID: uid.String(),
		Email:  email,
		OTP:    otp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return e.enqueueTask(ctx, asynq.NewTask(SendDestroyProfileEmailTask, payload))
}
