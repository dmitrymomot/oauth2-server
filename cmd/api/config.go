package main

import (
	"time"

	"github.com/dmitrymomot/go-env"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

var (
	// Application
	appName  = env.GetString("APP_NAME", "api")
	appDebug = env.GetBool("APP_DEBUG", false)

	// Product
	productName    = env.GetString("PRODUCT_NAME", "OAuth2 API") // To show on client side
	productIconURI = env.GetString("PRODUCT_ICON", "")           // absolute URI to product icon

	// HTTP Router
	httpPort                  = env.GetInt("HTTP_PORT", 8080)
	httpRequestTimeout        = env.GetDuration("HTTP_REQUEST_TIMEOUT", time.Second*10)
	httpServerShutdownTimeout = env.GetDuration("HTTP_SERVER_SHUTDOWN_TIMEOUT", time.Second*5)
	// httpLimitRequestBodySize  = env.GetInt[int64]("HTTP_LIMIT_REQUEST_BODY_SIZE", 1<<20) // 1 MB
	httpRateLimit         = env.GetInt("HTTP_RATE_LIMIT", 100)
	httpRateLimitDuration = env.GetDuration("HTTP_RATE_LIMIT_DURATION", time.Minute)

	// Cors
	corsAllowedOrigins     = env.GetStrings("CORS_ALLOWED_ORIGINS", ",", []string{"*"})
	corsAllowedMethods     = env.GetStrings("CORS_ALLOWED_METHODS", ",", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"})
	corsAllowedHeaders     = env.GetStrings("CORS_ALLOWED_HEADERS", ",", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID", "X-Request-Id", "Origin", "User-Agent", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Pragma", "Referer"})
	corsAllowedCredentials = env.GetBool("CORS_ALLOWED_CREDENTIALS", true)
	corsMaxAge             = env.GetInt("CORS_MAX_AGE", 300)

	// Build tag is set up while deployment
	buildTag        = "undefined"
	buildTagRuntime = env.GetString("COMMIT_HASH", buildTag)

	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 20)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 2)

	// Redis
	redisConnString = env.GetString("REDIS_DATABASE_URL", "")

	// Worker
	workerConcurrency = env.GetInt("WORKER_CONCURRENCY", 10)
	queueName         = env.GetString("QUEUE_NAME", "default")

	// // Auth
	// oauthSigningKey = env.MustString("OAUTH_SIGNING_KEY")
	// accessTokenTTL  = env.GetDuration("ACCESS_TOKEN_TTL", time.Minute*5)
	// refreshTokenTTL = env.GetDuration("REFRESH_TOKEN_TTL", time.Hour)
	// clientID        = env.MustString("CLIENT_ID")
	// clientSecret    = env.MustString("CLIENT_SECRET")
)
