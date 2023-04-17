package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hibiken/asynq"
)

const (
	CleanUpExpiredVerificationRequestsTask = "clean_up_expired_verification_requests"
	CleanUpExpiredTokensTask               = "clean_up_expired_tokens"
)

type (
	// Worker is a task handler for email delivery.
	Worker struct {
		repo workerRepository
		log  logger
	}

	workerRepository interface {
		CleanUpExpiredUserVerifications(ctx context.Context) error
		DeleteExpiredTokens(ctx context.Context) error
	}

	logger interface {
		Errorf(format string, args ...interface{})
	}
)

// NewWorker creates a new email task handler.
func NewWorker(repo workerRepository, log logger) *Worker {
	return &Worker{repo: repo, log: log}
}

// Schedule schedules tasks for the worker.
func (w *Worker) Schedule(s *asynq.Scheduler) {
	s.Register("@every 1h", asynq.NewTask(CleanUpExpiredVerificationRequestsTask, nil),
		asynq.Queue("auth-exp-ver-reqs"),
		asynq.Unique(time.Hour),
		asynq.MaxRetry(0),
	)
	s.Register("@every 1h", asynq.NewTask(CleanUpExpiredTokensTask, nil),
		asynq.Queue("auth-exp-tokens"),
		asynq.Unique(time.Hour),
		asynq.MaxRetry(0),
	)
}

// Register registers task handlers for email delivery.
func (w *Worker) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(CleanUpExpiredVerificationRequestsTask, w.CleanUpExpiredVerificationRequests)
	mux.HandleFunc(CleanUpExpiredTokensTask, w.CleanUpExpiredTokens)
}

// CleanUpExpiredVerificationRequests cleans up expired verification requests.
func (w *Worker) CleanUpExpiredVerificationRequests(ctx context.Context, t *asynq.Task) error {
	if err := w.repo.CleanUpExpiredUserVerifications(ctx); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			w.log.Errorf("failed to clean up expired verification requests: %w", err)
		}
	}

	return nil
}

// CleanUpExpiredTokens cleans up expired tokens.
func (w *Worker) CleanUpExpiredTokens(ctx context.Context, t *asynq.Task) error {
	if err := w.repo.DeleteExpiredTokens(ctx); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			w.log.Errorf("failed to clean up expired tokens: %w", err)
		}
	}

	return nil
}
