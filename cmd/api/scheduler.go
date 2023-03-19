package main

import "github.com/hibiken/asynq"

type (
	schedulerHandler interface {
		Schedule(*asynq.Scheduler)
	}
)

// runScheduler creates a new scheduler server and registers task handlers.
func runScheduler(redisConnOpt asynq.RedisConnOpt, log asynq.Logger, handlers ...schedulerHandler) func() error {
	return func() error {
		// Setup asynq scheduler
		scheduler := asynq.NewScheduler(
			redisConnOpt,
			&asynq.SchedulerOpts{
				Logger:   log,
				LogLevel: asynq.InfoLevel,
			},
		)

		// Register handlers
		for _, h := range handlers {
			h.Schedule(scheduler)
		}

		// Run scheduler
		return scheduler.Run()
	}
}
