package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq" // init pg driver

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Init logger
	logrus.SetReportCaller(false)
	logger := logrus.WithFields(logrus.Fields{
		"app":       appName,
		"build_tag": buildTagRuntime,
	})
	if appDebug {
		logger.Logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.Logger.SetLevel(logrus.InfoLevel)
	}

	defer func() { logger.Info("Server successfully shutdown") }()

	// Errgroup with context
	eg, ctx := errgroup.WithContext(newCtx(logger))

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init db connection")
	}
	defer db.Close()

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping db")
	}

	// Init repository
	repo, err := repository.Prepare(ctx, db)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init repository")
	}
	_ = repo // TODO: use repo

	if redisConnString != "" {
		// Redis connect options for asynq client
		redisConnOpt, err := asynq.ParseRedisURI(redisConnString)
		if err != nil {
			logger.WithError(err).Fatal("Failed to parse redis connection string")
		}

		// Init asynq client
		asynqClient := asynq.NewClient(redisConnOpt)
		defer asynqClient.Close()

		// Run asynq worker
		eg.Go(runQueueServer(
			redisConnOpt,
			logger,
			// TODO: add all queues here
		))

		// Run asynq scheduler
		eg.Go(runScheduler(
			redisConnOpt,
			logger,
			// TODO: add all schedulers here
		))
	} else {
		logger.Warn("Redis connection string is empty, skipping asynq client")
	}

	// Init HTTP router
	r := initRouter(logger)

	// Run HTTP server
	eg.Go(runServer(ctx, httpPort, r, logger))

	// Run all goroutines
	if err := eg.Wait(); err != nil {
		logger.WithError(err).Fatal("Error occurred")
	}
}

// newCtx creates a new context that is cancelled when an interrupt signal is received.
func newCtx(log *logrus.Entry) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		sCh := make(chan os.Signal, 1)
		signal.Notify(sCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGPIPE)
		<-sCh

		log.Debug("Received interrupt signal, shutting down")
	}()
	return ctx
}
