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
	"github.com/dmitrymomot/random"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// newClientCmd represents the newClient command
var newClientCmd = &cobra.Command{
	Use:   "new-client",
	Short: "Create a new client",
	Long:  `Create a new client and return client id and secret to use in API client.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		connStr := cmd.Flag("db").Value.String()
		if connStr == "" {
			connStr = env.GetString("DATABASE_URL", "")
			if connStr == "" {
				return fmt.Errorf("db connection string is required")
			}
		}

		isPublic, _ := cmd.Flags().GetBool("public")

		clientID, clientSecret, err := createNewClient(
			connStr,
			isPublic,
			cmd.Flag("domain").Value.String(),
			cmd.Flag("user_id").Value.String(),
		)
		if err != nil {
			return fmt.Errorf("failed to create new client: %w", err)
		}

		color.Green("\nNew client generated")
		bold := color.New(color.Bold).SprintFunc()
		fmt.Println("---------------------------------------------------------------------------------")
		fmt.Println(bold("Client ID:          "), clientID)
		fmt.Println(bold("Client Secret:      "), clientSecret)
		fmt.Println("---------------------------------------------------------------------------------")
		color.Yellow("Please save the client ID and client secret somewhere safe.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newClientCmd)
	newClientCmd.Flags().BoolP("public", "t", false, "Is the client public?")
	newClientCmd.Flags().String("db", "", "Database connection string")
	newClientCmd.Flags().StringP("domain", "d", "", "Client domain")
	newClientCmd.Flags().StringP("user_id", "u", "", "User ID")
}

func createNewClient(dbConnString string, public bool, domain, userID string) (id, secret string, err error) {
	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		return "", "", fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return "", "", fmt.Errorf("failed to ping db: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init repository
	repo, err := repository.Prepare(ctx, db)
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare repository: %w", err)
	}

	clientID := fmt.Sprintf("id_%s", random.String(32))
	clientSecret := fmt.Sprintf("secret_%s", random.String(32))

	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash client secret: %w", err)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user id: %w", err)
	}

	// Create client
	if _, err := repo.CreateClient(ctx, repository.CreateClientParams{
		ID:       clientID,
		Secret:   clientSecretHash,
		Domain:   domain,
		IsPublic: public,
		UserID:   uid,
	}); err != nil {
		return "", "", fmt.Errorf("failed to create client: %w", err)
	}

	return clientID, clientSecret, nil
}
