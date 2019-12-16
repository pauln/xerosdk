package auth

import (
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
)

// Repository will keep the API information for the user sessions between
// quicka and xero platform
type Repository interface {
	CreateSession(userID uuid.UUID, t *oauth2.Token) error
	UpdateSession(userID uuid.UUID, t *oauth2.Token) error
	GetSession(userID uuid.UUID) (*oauth2.Token, error)
}

// TokenRefresher keep the information needed for our custom TokenSource
type TokenRefresher struct {
	repo     Repository
	token    *oauth2.Token
	provider *Provider
	userID   uuid.UUID
}

// NewTokenRefresher function will build a new TokenSource based on a given token
// and a given repository
func NewTokenRefresher(repo Repository, token *oauth2.Token, provider *Provider, userID uuid.UUID) oauth2.TokenSource {
	return &TokenRefresher{
		repo:     repo,
		token:    token,
		userID:   userID,
		provider: provider,
	}
}

// Token method is the custom implementation of the refresh token process using
// a session repo as a base
func (t *TokenRefresher) Token() (*oauth2.Token, error) {
	if !t.token.Valid() {
		token, err := t.provider.Refresh(t.token)
		if err != nil {
			return nil, err
		}
		if err = t.repo.UpdateSession(t.userID, token); err != nil {
			return nil, err
		}
		return token, nil
	}
	return t.token, nil
}
