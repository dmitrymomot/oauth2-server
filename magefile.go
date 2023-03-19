//go:build mage
// +build mage

package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// mg contains helpful utility functions, like Deps

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Run runs the application with environment variables
func Run() error {
	// clean up go build cache
	color.Yellow("Cleaning up go build cache...")
	sh.RunV("go", "clean", "-cache")

	// run the application
	color.Cyan("Starting the application...")
	return sh.RunV("go", "run", "./cmd/api/")
}

// MigrateUp runs the database migrations up
func MigrateUp() error {
	color.Cyan("Running database migrations...")
	return sh.RunV("go", "run", "./cmd/migrate/main.go")
}

// Up runs the application and the database migrations up
func Up() error {
	color.Cyan("Starting the database...")
	if err := sh.RunV("docker-compose", "up", "-d"); err != nil {
		return err
	}

	color.Yellow("Waiting for the database to start...")
	time.Sleep(5 * time.Second)

	if err := MigrateUp(); err != nil {
		return err
	}

	return Run()
}

// Down stops the application and the database
func Down() error {
	color.Yellow("Stopping the database and removing all data...")
	return sh.RunV("docker-compose", "down", "--volumes", "--rmi=local")
}

// Client namespace
type Client mg.Namespace

// Create a new client
func (Client) Create(userID, domain, isPublic string) error {
	color.Cyan("Creating a new client...")

	return sh.RunV(
		"go", "run", "./cmd/cli/", "new-client",
		"--user_id="+userID,
		"--domain="+domain,
		"--public="+isPublic,
	)
}

// User namespace
type User mg.Namespace

// Create a new user
func (User) Create(email, password string) error {
	color.Cyan("Creating a new user...")

	return sh.RunV(
		"go", "run", "./cmd/cli/", "new-user",
		"--email="+email,
		"--password="+password,
	)
}
