package main

import (
	"database/sql"

	_ "github.com/lib/pq" // init pg driver
	"github.com/sirupsen/logrus"

	"github.com/dmitrymomot/go-env"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	appName = env.GetString("APP_NAME", "db-migrate")

	// DB
	dbConnString    = env.MustString("DATABASE_URL")
	migrationsDir   = env.GetString("DATABASE_MIGRATIONS_DIR", "./repository/sql/migrations")
	migrationsTable = env.GetString("DATABASE_MIGRATIONS_TABLE", "migrations")

	// Build tag is set up while deployment
	buildTag        = "undefined"
	buildTagRuntime = env.GetString("COMMIT_HASH", buildTag)
)

func main() {
	// Init logger
	logrus.SetReportCaller(false)
	logger := logrus.WithFields(logrus.Fields{
		"app":       appName,
		"build_tag": buildTagRuntime,
	})
	logger.Logger.SetLevel(logrus.InfoLevel)

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		logger.WithError(err).Fatal("failed to init db connection")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("failed to ping db")
	}

	m := migrate.MigrationSet{
		TableName: migrationsTable,
	}
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	n, err := m.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.WithError(err).Fatal("failed to apply migrations")
	}

	logger.Infof("applied %d migrations!", n)
}
