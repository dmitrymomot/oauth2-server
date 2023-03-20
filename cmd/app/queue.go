package main

import (
	"github.com/hibiken/asynq"
)

type (
	// taskHandler is an interface for task handlers.
	taskHandler interface {
		Register(*asynq.ServeMux)
	}
)

// setupQueue creates a new queue client and registers task handlers.
func runQueueServer(redisConnOpt asynq.RedisConnOpt, log asynq.Logger, handlers ...taskHandler) func() error {
	return func() error {
		// Setup asynq server
		srv := asynq.NewServer(
			redisConnOpt,
			asynq.Config{
				Concurrency: workerConcurrency,
				Logger:      log,
				Queues: map[string]int{
					queueName: workerConcurrency,
				},
			},
		)

		// Run server
		return srv.Run(registerQueueHandlers(handlers...))
	}
}

// registerQueueHandlers registers handlers for each task type.
func registerQueueHandlers(handlers ...taskHandler) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	// Register handlers
	for _, h := range handlers {
		h.Register(mux)
	}

	return mux
}
