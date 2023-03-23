package client

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Service is the client service interface.
	Service interface {
		// Create creates a new client.
		Create(ctx context.Context, uid string, domain string, isPublic bool) (*Client, error)
		// GetByID returns a client by its ID.
		GetByID(ctx context.Context, id string) (*Client, error)
		// GetByUserID returns a clients list by its user ID.
		GetByUserID(ctx context.Context, uid string) ([]*Client, error)
		// Delete deletes a client by its ID.
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo clientRepository
	}

	clientRepository interface {
		CreateClient(ctx context.Context, arg repository.CreateClientParams) (repository.Client, error)
		DeleteClient(ctx context.Context, id string) error
		GetClientByID(ctx context.Context, id string) (repository.Client, error)
		GetClientByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Client, error)
	}
)

// NewService returns a new instance of a service.
func NewService(repo clientRepository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new client.
func (s *service) Create(ctx context.Context, userID string, domain string, isPublic bool) (*Client, error) {
	clientID := fmt.Sprintf("id_%s", random.String(32))
	clientSecret := fmt.Sprintf("secret_%s", random.String(32))

	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash client secret: %w", err)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %w", err)
	}

	allowedGrants := []string{
		"authorization_code",
		"refresh_token",
		"__implicit",
	}
	if !isPublic {
		allowedGrants = append(allowedGrants, "client_credentials")
	}

	// Create client
	c, err := s.repo.CreateClient(ctx, repository.CreateClientParams{
		ID:            clientID,
		Secret:        clientSecretHash,
		Domain:        domain,
		IsPublic:      isPublic,
		UserID:        uid,
		AllowedGrants: allowedGrants,
		Scope:         "client:* user:*",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return NewClient(c, clientSecret), nil
}

// GetByID returns a client by its ID.
func (s *service) GetByID(ctx context.Context, id string) (*Client, error) {
	client, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get client by id: %w", err)
	}

	return NewClient(client, ""), nil
}

// GetByUserID returns a clients list by its user ID.
func (s *service) GetByUserID(ctx context.Context, uid string) ([]*Client, error) {
	userID, err := uuid.Parse(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %w", err)
	}

	clients, err := s.repo.GetClientByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients by user id: %w", err)
	}

	result := make([]*Client, 0, len(clients))
	for _, c := range clients {
		result = append(result, NewClient(c, ""))
	}

	return result, nil
}

// Delete deletes a client by its ID.
func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.DeleteClient(ctx, id); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}
