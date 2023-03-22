package main

import (
	"time"

	"github.com/dmitrymomot/go-env"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

var (
	// Application
	appName    = env.GetString("APP_NAME", "api")
	appDebug   = env.GetBool("APP_DEBUG", false)
	appBaseURL = env.GetString("APP_BASE_URL", "http://localhost:8080")
	appLogoURL = env.GetString("APP_LOGO_URL", "")

	// Product
	productName    = env.GetString("PRODUCT_NAME", "OAuth2 API") // To show on client side
	productURL     = env.MustString("PRODUCT_URL")
	supportEmail   = env.MustString("SUPPORT_EMAIL")
	companyName    = env.GetString("COMPANY_NAME", productName)
	productLogoURL = env.GetString("LOGO_URL", "") // absolute URI to product logo

	// HTTP Router
	httpPort                  = env.GetInt("HTTP_PORT", 8080)
	httpRequestTimeout        = env.GetDuration("HTTP_REQUEST_TIMEOUT", time.Second*10)
	httpServerShutdownTimeout = env.GetDuration("HTTP_SERVER_SHUTDOWN_TIMEOUT", time.Second*5)
	httpRateLimit             = env.GetInt("HTTP_RATE_LIMIT", 100)
	httpRateLimitDuration     = env.GetDuration("HTTP_RATE_LIMIT_DURATION", time.Minute)

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
	queueTaskDeadline = env.GetDuration("QUEUE_TASK_DEADLINE", time.Minute)
	queueMaxRetry     = env.GetInt("QUEUE_TASK_RETRY_LIMIT", 3)

	// Auth
	oauthSigningKey   = env.MustString("OAUTH_SIGNING_KEY")
	authorizedHomeURI = env.GetString("AUTHORIZED_HOME_URI", "http://localhost:3000")

	// Postmark
	postmarkServerToken  = env.MustString("POSTMARK_SERVER_TOKEN")
	postmarkProjectToken = env.MustString("POSTMARK_ACCOUNT_TOKEN")

	// Mailer
	mailFromEmail = env.MustString("MAILER_FROM_EMAIL")
	mailFromName  = env.GetString("MAILER_FROM_NAME", productName)

	// Session
	sessionSigningKey     = env.MustString("SESSION_SIGNING_KEY")
	sessionCookieName     = env.GetString("SESSION_COOKIE_NAME", "session")
	sessionCookieLifeTime = env.GetInt("SESSION_COOKIE_LIFE_TIME", 86400)
	sessionCookieDomain   = env.GetString("SESSION_COOKIE_DOMAIN", "")
	sessionCookieSecure   = env.GetBool("SESSION_COOKIE_SECURE", false)
	sessionCookieHttpOnly = env.GetBool("SESSION_COOKIE_HTTP_ONLY", true)
	sessionCookieSameSite = env.GetString("SESSION_COOKIE_SAME_SITE", "lax")
	sessionExpiresIn      = env.GetInt("SESSION_EXPIRES_IN", int64(86400))
)
