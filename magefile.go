//go:build mage
// +build mage

package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically

	"github.com/fatih/color"
	"github.com/google/uuid"
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
func (Client) Create(projectID, redirectURI, isInternal, desc string) error {
	color.Cyan("Creating a new client...")

	if projectID == "" {
		projectID = uuid.New().String()
	}

	if redirectURI == "" {
		redirectURI = "http://localhost:9094/oauth2"
	}

	return sh.RunV(
		"go", "run", "./cmd/client/main.go",
		"-project-id="+projectID,
		"-redirect-uri="+redirectURI,
		"-internal="+isInternal,
		"-desc="+desc,
	)
}

// App namespace
type App mg.Namespace

// Create a new app
func (App) Create(id, name, isInternal string) error {
	color.Cyan("Creating a new app...")

	return sh.RunV("go", "run", "./cmd/apps/main.go", "-id="+id, "-name="+name, "-internal="+isInternal)
}

// Add scope to an app
func (App) AddScope(appID, scopeID, scopeName string) error {
	color.Cyan("Adding scope to an app...")

	return sh.RunV("go", "run", "./cmd/scopes/main.go", "-id", scopeID, "-name", scopeName, "-app-id", appID)
}

// Add test data
func InitTestData() error {
	color.Cyan("Adding test data...")

	// Create a new internal client
	if err := sh.RunV("mage", "client:create", "", "", "1", "Internal client to manage clients"); err != nil {
		return err
	}

	// Create a new general client
	if err := sh.RunV("mage", "client:create", "", "", "0", "Example of general customer client, which can be used to user auth, etc"); err != nil {
		return err
	}

	// Create new apps
	if err := sh.RunV("mage", "app:create", "profile", "Profile", "0"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:create", "wallet", "Wallet", "0"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:create", "client", "Client Manager", "1"); err != nil {
		return err
	}

	// Add scopes to an app
	if err := sh.RunV("mage", "app:addScope", "profile", "read", "Read profile"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:addScope", "profile", "update", "Update profile"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:addScope", "wallet", "read", "Get wallet address and balance"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:addScope", "wallet", "transaction", "Make transactions"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:addScope", "client", "read", "Get client data"); err != nil {
		return err
	}
	if err := sh.RunV("mage", "app:addScope", "client", "write", "Create new clients and update existed"); err != nil {
		return err
	}

	return nil
}
