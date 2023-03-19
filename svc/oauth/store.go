package oauth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmitrymomot/oauth2-server/repository"
	"github.com/go-oauth2/oauth2/v4"
	oauth2Errors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/google/uuid"
)

type (
	Store struct {
		repo oauthRepository
	}

	oauthRepository interface {
		GetClientByID(ctx context.Context, id string) (repository.Client, error)

		CreateToken(ctx context.Context, arg repository.CreateTokenParams) (repository.Token, error)
		DeleteByAccess(ctx context.Context, access string) error
		DeleteByCode(ctx context.Context, code string) error
		DeleteByRefresh(ctx context.Context, refresh string) error
		DeleteExpiredTokens(ctx context.Context) error
		GetTokenByAccess(ctx context.Context, access string) (repository.Token, error)
		GetTokenByCode(ctx context.Context, code string) (repository.Token, error)
		GetTokenByRefresh(ctx context.Context, refresh string) (repository.Token, error)
	}
)

// NewStore creates a new store instance.
// The store is used to manage the client and token information.
// Implements the interface of the oauth2.ClientStore and oauth2.TokenStore.
func NewStore(repo oauthRepository) *Store {
	return &Store{
		repo: repo,
	}
}

// according to the ID for the client information
func (s *Store) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	client, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get client by id: %w", err)
	}

	return NewClient(client, ""), nil
}

// create and store the new token information
func (s *Store) Create(ctx context.Context, info oauth2.TokenInfo) error {
	var uid uuid.UUID
	if info.GetUserID() != "" {
		id, err := uuid.Parse(info.GetUserID())
		if err != nil {
			return fmt.Errorf("failed to parse user id: %w", err)
		}
		uid = id
	}

	if _, err := s.repo.CreateToken(ctx, repository.CreateTokenParams{
		ClientID:    info.GetClientID(),
		UserID:      uuid.NullUUID{UUID: uid, Valid: uid != uuid.Nil},
		RedirectURI: info.GetRedirectURI(),
		Scope:       info.GetScope(),
		Code:        info.GetCode(),
		CodeCreatedAt: func() sql.NullTime {
			if info.GetCodeCreateAt().IsZero() {
				return sql.NullTime{}
			}
			return sql.NullTime{
				Time:  info.GetCodeCreateAt(),
				Valid: true,
			}
		}(),
		CodeExpiresIn:       int64(info.GetCodeExpiresIn().Seconds()),
		CodeChallenge:       info.GetCodeChallenge(),
		CodeChallengeMethod: string(info.GetCodeChallengeMethod()),
		Access:              info.GetAccess(),
		AccessCreatedAt: func() sql.NullTime {
			if info.GetAccessCreateAt().IsZero() {
				return sql.NullTime{}
			}
			return sql.NullTime{
				Time:  info.GetAccessCreateAt(),
				Valid: true,
			}
		}(),
		AccessExpiresIn: int64(info.GetAccessExpiresIn().Seconds()),
		Refresh:         info.GetRefresh(),
		RefreshCreatedAt: func() sql.NullTime {
			if info.GetRefreshCreateAt().IsZero() {
				return sql.NullTime{}
			}
			return sql.NullTime{
				Time:  info.GetRefreshCreateAt(),
				Valid: true,
			}
		}(),
		RefreshExpiresIn: int64(info.GetRefreshExpiresIn().Seconds()),
	}); err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

// delete the authorization code
func (s *Store) RemoveByCode(ctx context.Context, code string) error {
	if err := s.repo.DeleteByCode(ctx, code); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete by code: %w", err)
		}
	}
	return nil
}

// use the access token to delete the token information
func (s *Store) RemoveByAccess(ctx context.Context, access string) error {
	if err := s.repo.DeleteByAccess(ctx, access); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete by access: %w", err)
		}
	}
	return nil
}

// use the refresh token to delete the token information
func (s *Store) RemoveByRefresh(ctx context.Context, refresh string) error {
	if err := s.repo.DeleteByRefresh(ctx, refresh); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete by refresh: %w", err)
		}
	}
	return nil
}

// use the authorization code for token information data
func (s *Store) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	token, err := s.repo.GetTokenByCode(ctx, code)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get token by code: %w", err)
		}
		return nil, oauth2Errors.ErrInvalidAuthorizeCode
	}

	return NewToken(token), nil
}

// use the access token for token information data
func (s *Store) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	token, err := s.repo.GetTokenByAccess(ctx, access)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get token by code: %w", err)
		}
		return nil, oauth2Errors.ErrInvalidAccessToken
	}

	return NewToken(token), nil
}

// use the refresh token for token information data
func (s *Store) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	token, err := s.repo.GetTokenByRefresh(ctx, refresh)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get token by code: %w", err)
		}
		return nil, oauth2Errors.ErrInvalidRefreshToken
	}

	return NewToken(token), nil
}
