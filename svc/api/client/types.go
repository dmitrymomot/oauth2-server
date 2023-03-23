package client

import (
	"time"

	"github.com/dmitrymomot/oauth2-server/repository"
)

// Client represents an OAuth client.
type Client struct {
	ID        string `json:"id"`
	Secret    string `json:"secret,omitempty"`
	Domain    string `json:"domain"`
	Public    bool   `json:"is_public"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// NewClient creates a new client instance.
// The secret is hashed before being stored in the database.
// So it can be returned only once after creation.
func NewClient(source repository.Client, secret string) *Client {
	return &Client{
		ID:        source.ID,
		Secret:    secret,
		Domain:    source.Domain,
		Public:    source.IsPublic,
		UserID:    source.UserID.String(),
		CreatedAt: source.CreatedAt.Format(time.RFC3339),
	}
}
