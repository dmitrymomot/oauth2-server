package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dmitrymomot/oauth2-server/internal/mdw"
	postmarkClient "github.com/dmitrymomot/oauth2-server/internal/postmark"
	"github.com/dmitrymomot/oauth2-server/lib/middleware"
	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/dmitrymomot/oauth2-server/svc/api/client"
	"github.com/dmitrymomot/oauth2-server/svc/api/user"
	"github.com/dmitrymomot/oauth2-server/svc/auth"
	"github.com/dmitrymomot/oauth2-server/svc/mailer"
	"github.com/dmitrymomot/oauth2-server/svc/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"github.com/keighl/postmark"
	_ "github.com/lib/pq" // init pg driver
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Init logger
	logger := logrus.WithFields(logrus.Fields{
		"app":       appName,
		"build_tag": buildTagRuntime,
	})
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

	// mail enqueuer
	var mailEnqueuer *mailer.Enqueuer
	if redisConnString != "" {
		// Redis connect options for asynq client
		redisConnOpt, err := asynq.ParseRedisURI(redisConnString)
		if err != nil {
			logger.WithError(err).Fatal("Failed to parse redis connection string")
		}
		redisConn := redisConnOpt.MakeRedisClient()
		redisClient, ok := redisConn.(*redis.Client)
		if !ok {
			logger.Fatal("Failed to cast redis connection to *redis.Client")
		}

		// init the session manager
		initSessionManager(redisClient)

		// Init asynq client
		asynqClient := asynq.NewClient(customRedisConnOpt{redis: redisClient})
		defer asynqClient.Close()

		// Init mail enqueuer
		mailEnqueuer = mailer.NewEnqueuer(
			asynqClient,
			mailer.WithLogger(logger.WithField("component", "mailer")),
			mailer.WithQueueName(queueName),
			mailer.WithMaxRetry(queueMaxRetry),
			mailer.WithTaskDeadline(queueTaskDeadline),
		)

		baseURL := strings.TrimSuffix(appBaseURL, "/")
		if baseURL == "" || baseURL == "/" || !strings.HasPrefix(baseURL, "http") {
			logger.Fatal("Application base URL is invalid")
		}

		// Mailer service
		pc := postmarkClient.New(
			postmark.NewClient(postmarkServerToken, postmarkProjectToken),
			postmarkClient.Config{
				ProductName:  productName,
				ProductURL:   productURL,
				ProductLogo:  productLogoURL,
				SupportEmail: supportEmail,
				CompanyName:  companyName,
				FromEmail:    mailFromEmail,
				FromName:     mailFromName,

				VerificationCodeURL: fmt.Sprintf("%s/%s", baseURL, "/auth/verification/verify"),
				PasswordResetURL:    fmt.Sprintf("%s/%s", baseURL, "/auth/password/reset"),
				DestroyUserCodeURL:  fmt.Sprintf("%s/%s", baseURL, "/auth/account/destroy/verify"),
			},
		)

		// Run asynq worker
		eg.Go(runQueueServer(
			redisConnOpt,
			logger.WithField("component", "queue-worker"),
			mailer.NewWorker(pc),
		))

		// Run asynq scheduler
		eg.Go(runScheduler(
			redisConnOpt,
			logger.WithField("component", "scheduler"),
			auth.NewWorker(repo, logger.WithField("component", "auth-worker")),
		))
	} else {
		logger.Warn("Redis connection string is empty, skipping asynq client")
	}

	// Init HTTP router
	r := initRouter(logger.WithField("component", "http-router"))

	// Mount oauth2 server
	{
		storage := oauth.NewStore(repo)
		srv, manager := oauth.NewOauth2Server(
			generates.NewJWTAccessGenerate("", []byte(oauthSigningKey), jwt.SigningMethodHS512),
			generates.NewAuthorizeGenerate(),
			storage, storage,
			oauth.NewHandlerLogger(
				oauth.NewHandler(
					repo,
					oauth.WithClientScope("user:read client:read"),
					oauth.WithPasswordScope("user:*"),
					oauth.WithCodeScope("user:* client:*"),
				),
				logger.WithField("component", "oauth2"),
			),
		)

		r.Mount("/oauth", oauth.MakeHTTPHandler(
			srv,
			manager,
			logger.WithField("component", "oauth2"),
			"/auth/login",
		))
	}

	// Mount auth service
	r.Mount("/auth", auth.MakeHTTPHandler(
		auth.NewService(repo, db, mailEnqueuer),
		"/oauth/authorize",
		mdw.NotAuthOnly(authorizedHomeURI),
	))

	// Mount api services
	r.Route("/api", func(api chi.Router) {
		api.Mount("/user", user.MakeHTTPHandler(
			user.MakeEndpoints(
				user.NewService(repo, mailEnqueuer, db),
				middleware.GokitAuthMiddleware(
					middleware.VerifyJWT(oauthSigningKey),
				),
			),
			logger.WithField("component", "api-user"),
		))

		api.Mount("/client", client.MakeHTTPHandler(
			client.MakeEndpoints(
				client.NewService(repo),
				middleware.GokitAuthMiddleware(
					middleware.VerifyJWT(oauthSigningKey),
				),
			),
			logger.WithField("component", "api-client"),
		))
	})

	// Run HTTP server
	eg.Go(runServer(ctx, httpPort, r, logger.WithField("component", "http-server")))

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

type customRedisConnOpt struct {
	redis *redis.Client
}

func (c customRedisConnOpt) MakeRedisClient() interface{} {
	return c.redis
}
