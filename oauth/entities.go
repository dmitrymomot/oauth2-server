package oauth

import (
	"time"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

// Client represents an OAuth client implements the oauth2.ClientInfo interface.
type Client struct {
	ID         string    `json:"id"`
	Secret     string    `json:"secret,omitempty"`
	secretHash []byte    `json:"-"` // hash of the secret
	Domain     string    `json:"domain"`
	Public     bool      `json:"is_public"`
	UserID     uuid.UUID `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewClient creates a new client instance.
// The secret is hashed before being stored.
// Client implements the ClientInfo interface.
func NewClient(source repository.Client, secret string) *Client {
	return &Client{
		ID:         source.ID,
		Secret:     secret,
		secretHash: source.Secret,
		Domain:     source.Domain,
		Public:     source.IsPublic,
		UserID:     source.UserID,
		CreatedAt:  source.CreatedAt,
	}
}

// GetID returns the client ID.
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret returns the client secret.
func (c *Client) GetSecret() string {
	return string(c.secretHash)
}

// GetDomain returns the client domain.
func (c *Client) GetDomain() string {
	return c.Domain
}

// IsPublic returns true if the client is public.
func (c *Client) IsPublic() bool {
	return c.Public
}

// GetUserID returns the client user ID.
func (c *Client) GetUserID() string {
	return c.UserID.String()
}

// Token represents an OAuth token implements the oauth2.TokenInfo interface.
type Token struct {
	ID                  uuid.UUID  `json:"id"`
	ClientID            string     `json:"client_id"`
	UserID              uuid.UUID  `json:"user_id,omitempty"`
	RedirectURI         string     `json:"redirect_uri,omitempty"`
	Scope               string     `json:"scope,omitempty"`
	Code                string     `json:"code,omitempty"`
	CodeCreatedAt       *time.Time `json:"code_created_at,omitempty"`
	CodeExpiresIn       int64      `json:"code_expires_in,omitempty"`
	CodeChallenge       string     `json:"code_challenge,omitempty"`
	CodeChallengeMethod string     `json:"code_challenge_method,omitempty"`
	Access              string     `json:"access,omitempty"`
	AccessCreatedAt     *time.Time `json:"access_created_at,omitempty"`
	AccessExpiresIn     int64      `json:"access_expires_in,omitempty"`
	Refresh             string     `json:"refresh,omitempty"`
	RefreshCreatedAt    *time.Time `json:"refresh_created_at,omitempty"`
	RefreshExpiresIn    int64      `json:"refresh_expires_in,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

// NewToken creates a new token instance from a repository token.
func NewToken(source repository.Token) *Token {
	t := &Token{
		ID:                  source.ID,
		ClientID:            source.ClientID,
		UserID:              source.UserID,
		RedirectURI:         source.RedirectURI,
		Scope:               source.Scope,
		Code:                source.Code,
		CodeExpiresIn:       source.CodeExpiresIn,
		CodeChallenge:       source.CodeChallenge,
		CodeChallengeMethod: source.CodeChallengeMethod,
		Access:              source.Access,
		AccessExpiresIn:     source.AccessExpiresIn,
		Refresh:             source.Refresh,
		RefreshExpiresIn:    source.RefreshExpiresIn,
		CreatedAt:           source.CreatedAt,
	}

	if source.CodeCreatedAt.Valid {
		t.CodeCreatedAt = &source.CodeCreatedAt.Time
	}

	if source.AccessCreatedAt.Valid {
		t.AccessCreatedAt = &source.AccessCreatedAt.Time
	}

	if source.RefreshCreatedAt.Valid {
		t.RefreshCreatedAt = &source.RefreshCreatedAt.Time
	}

	return t
}

func (t *Token) New() oauth2.TokenInfo {
	return &Token{}
}

func (t *Token) GetClientID() string {
	return t.ClientID
}

func (t *Token) SetClientID(id string) {
	t.ClientID = id
}

func (t *Token) GetUserID() string {
	return t.UserID.String()
}

func (t *Token) SetUserID(id string) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return
	}
	t.UserID = uid
}

func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

func (t *Token) SetRedirectURI(uri string) {
	t.RedirectURI = uri
}

func (t *Token) GetScope() string {
	return t.Scope
}

func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

func (t *Token) GetCode() string {
	return t.Code
}

func (t *Token) SetCode(code string) {
	t.Code = code
}

func (t *Token) GetCodeCreateAt() time.Time {
	if t.CodeCreatedAt == nil {
		return time.Time{}
	}
	return *t.CodeCreatedAt
}

func (t *Token) SetCodeCreateAt(createdAt time.Time) {
	t.CodeCreatedAt = &createdAt
}

func (t *Token) GetCodeExpiresIn() time.Duration {
	return time.Duration(t.CodeExpiresIn) * time.Second
}

func (t *Token) SetCodeExpiresIn(expIn time.Duration) {
	t.CodeExpiresIn = int64(expIn.Seconds())
}

func (t *Token) GetCodeChallenge() string {
	return t.CodeChallenge
}

func (t *Token) SetCodeChallenge(challenge string) {
	t.CodeChallenge = challenge
}

func (t *Token) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod(t.CodeChallengeMethod)
}

func (t *Token) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
	t.CodeChallengeMethod = string(method)
}

func (t *Token) GetAccess() string {
	return t.Access
}

func (t *Token) SetAccess(access string) {
	t.Access = access
}

func (t *Token) GetAccessCreateAt() time.Time {
	if t.AccessCreatedAt == nil {
		return time.Time{}
	}
	return *t.AccessCreatedAt
}

func (t *Token) SetAccessCreateAt(createdAt time.Time) {
	t.AccessCreatedAt = &createdAt
}

func (t *Token) GetAccessExpiresIn() time.Duration {
	return time.Duration(t.AccessExpiresIn) * time.Second
}

func (t *Token) SetAccessExpiresIn(expIn time.Duration) {
	t.AccessExpiresIn = int64(expIn.Seconds())
}

func (t *Token) GetRefresh() string {
	return t.Refresh
}

func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

func (t *Token) GetRefreshCreateAt() time.Time {
	if t.RefreshCreatedAt == nil {
		return time.Time{}
	}
	return *t.RefreshCreatedAt
}

func (t *Token) SetRefreshCreateAt(createdAt time.Time) {
	t.RefreshCreatedAt = &createdAt
}

func (t *Token) GetRefreshExpiresIn() time.Duration {
	return time.Duration(t.RefreshExpiresIn) * time.Second
}

func (t *Token) SetRefreshExpiresIn(expIn time.Duration) {
	t.RefreshExpiresIn = int64(expIn.Seconds())
}
