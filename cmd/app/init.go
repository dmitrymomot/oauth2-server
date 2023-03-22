package main

import (
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	sessionRedis "github.com/go-session/redis/v3"
	gosession "github.com/go-session/session/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	if appDebug {
		// SetReportCaller sets whether the standard logger will include the calling
		// method as a field.
		// logrus.SetReportCaller(true)

		// Only log the debug severity or above.
		logrus.SetLevel(logrus.DebugLevel)

		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	} else {
		// Only log the info severity or above.
		logrus.SetLevel(logrus.InfoLevel)

		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Register the types for gob
	gob.Register(url.Values{})
	gob.Register(map[string]string{})
	gob.Register(map[string][]string{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[string][]interface{}{})
	gob.Register(map[interface{}]interface{}{})

	// init the template engine with the default template path
	initTemplateEngine()
}

// init session manager
func initSessionManager(log *logrus.Entry) {
	if redisConnString != "" {
		connURI, err := url.Parse(redisConnString)
		if err != nil {
			log.WithError(err).Fatal("failed to parse redis connection string")
		}

		redisPort, err := strconv.Atoi(connURI.Port())
		if err != nil {
			log.WithError(err).Fatal("failed to parse redis port")
		}

		redisHost := fmt.Sprintf("%s@%s", connURI.User.String(), connURI.Hostname())
		redisHost = strings.Trim(redisHost, "@")

		// Init the session manager
		gosession.InitManager(
			gosession.SetSign([]byte(sessionSigningKey)),
			gosession.SetCookieName(sessionCookieName),
			gosession.SetCookieLifeTime(sessionCookieLifeTime),
			gosession.SetDomain(sessionCookieDomain),
			gosession.SetSecure(sessionCookieSecure),
			gosession.SetExpired(sessionExpiresIn),
			gosession.SetStore(sessionRedis.NewRedisStore(
				&sessionRedis.Options{
					Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
					DB:   15,
				},
			)),
		)
	}
}
