package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	authURL  = "https://login.xero.com/identity/connect/authorize"
	tokenURL = "https://identity.xero.com/connect/token"
)

// Config keeps the information needed for do an OAuth2 connection
// with Xero
type Config struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	RedirectURL  string
}

// Provider type will keep the minimum structure for make the connection
// between quicka and Xero
type Provider struct {
	conf *oauth2.Config
	ctx  context.Context
}

// NewProvider function will build a new Provider with the given criteria
func NewProvider(c Config) *Provider {
	return &Provider{
		conf: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Scopes:       c.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
			RedirectURL: c.RedirectURL,
		},
		ctx: context.Background(),
	}
}

// GetAuthURL method will return the url for redirect and start the OAuth2
// process
func (c *Provider) GetAuthURL(state string) string {
	return c.conf.AuthCodeURL(state)
}

// GetTokenFromCode method will find the token with the given code, this method
// should be called after a success callback received from auth process
func (c *Provider) GetTokenFromCode(code string) (*oauth2.Token, error) {
	return c.conf.Exchange(c.ctx, code)
}

// Refresh method will refresh the given token
func (c *Provider) Refresh(t *oauth2.Token) (*oauth2.Token, error) {
	return c.conf.TokenSource(c.ctx, t).Token()
}

// Client method will build a new http.Client with the custom TokenRefresher
func (c *Provider) Client(t *oauth2.Token, repo Repository) *http.Client {
	return oauth2.NewClient(c.ctx, NewTokenRefresher(repo, t, c))
}

// NewClient method will return a new http.Client for use in our calls, using
// the TokenSource will even refresh the token if needed
func (c *Provider) NewClient(t *oauth2.Token) *http.Client {
	return c.conf.Client(c.ctx, t)
}
