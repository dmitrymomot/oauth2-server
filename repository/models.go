// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        string    `json:"id"`
	Secret    []byte    `json:"secret"`
	Domain    string    `json:"domain"`
	IsPublic  bool      `json:"is_public"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Token struct {
	ID                  uuid.UUID    `json:"id"`
	ClientID            string       `json:"client_id"`
	UserID              uuid.UUID    `json:"user_id"`
	RedirectURI         string       `json:"redirect_uri"`
	Scope               string       `json:"scope"`
	Code                string       `json:"code"`
	CodeCreatedAt       sql.NullTime `json:"code_created_at"`
	CodeExpiresIn       int64        `json:"code_expires_in"`
	CodeChallenge       string       `json:"code_challenge"`
	CodeChallengeMethod string       `json:"code_challenge_method"`
	Access              string       `json:"access"`
	AccessCreatedAt     sql.NullTime `json:"access_created_at"`
	AccessExpiresIn     int64        `json:"access_expires_in"`
	Refresh             string       `json:"refresh"`
	RefreshCreatedAt    sql.NullTime `json:"refresh_created_at"`
	RefreshExpiresIn    int64        `json:"refresh_expires_in"`
	CreatedAt           time.Time    `json:"created_at"`
}

type User struct {
	ID         uuid.UUID    `json:"id"`
	Email      string       `json:"email"`
	Password   []byte       `json:"password"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	VerifiedAt sql.NullTime `json:"verified_at"`
}
