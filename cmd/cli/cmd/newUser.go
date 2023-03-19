/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically

	"github.com/dmitrymomot/go-env"
	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// newUserCmd represents the newUser command
var newUserCmd = &cobra.Command{
	Use:   "new-user",
	Short: "Create a new user",
	Long:  `Create a new user with the given email and password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		connStr := cmd.Flag("db").Value.String()
		if connStr == "" {
			connStr = env.GetString("DATABASE_URL", "")
			if connStr == "" {
				return fmt.Errorf("db connection string is required")
			}
		}

		email := cmd.Flag("email").Value.String()
		password := cmd.Flag("password").Value.String()

		if err := createNewUser(connStr, email, password); err != nil {
			return fmt.Errorf("failed to create new user: %w", err)
		}

		color.Green("\nNew user created successfully!")
		bold := color.New(color.Bold).SprintFunc()
		fmt.Println("---------------------------------------------------------------------------------")
		fmt.Println(bold("Email:      "), email)
		fmt.Println(bold("Password:   "), password)
		fmt.Println("---------------------------------------------------------------------------------")
		color.Yellow("Use the above credentials to login. You can change the email and password later.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newUserCmd)

	// DB flag
	newUserCmd.Flags().String("db", "", "Database connection string")

	// Email flag
	newUserCmd.Flags().StringP("email", "e", "", "Email of the user")
	newUserCmd.MarkFlagRequired("email")

	// Password flag
	newUserCmd.Flags().StringP("password", "p", "", "Password of the user")
	newUserCmd.MarkFlagRequired("password")
}

// create a new user
func createNewUser(dbConnString string, email, password string) error {
	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init repository
	repo, err := repository.Prepare(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to prepare repository: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash client secret: %w", err)
	}

	// Create user
	if _, err := repo.CreateUser(ctx, repository.CreateUserParams{
		Email:    email,
		Password: passwordHash,
	}); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
