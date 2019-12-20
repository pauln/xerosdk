package auth

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
)

const (
	authURL  = "https://login.xero.com/identity/connect/authorize"
	tokenURL = "https://identity.xero.com/connect/token"

	tenantIDHeader = "xero-tenant-id"
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

// XeroTransport represents the information needed for custom Xero transport
type XeroTransport struct {
	T        http.RoundTripper
	TenantID uuid.UUID
}

// RoundTrip method will add on each request the custom header for inform the
// tenantID
func (xt *XeroTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(tenantIDHeader, xt.TenantID.String())
	return xt.T.RoundTrip(req)
}

// NewXeroTransport will build a new XeroTransport based on the given Tenant
func NewXeroTransport(tenantID uuid.UUID) *XeroTransport {
	return &XeroTransport{
		T:        http.DefaultTransport,
		TenantID: tenantID,
	}
}

// Session type represents the information for each session using for connect
// to Xero
type Session struct {
	Token    *oauth2.Token
	UserID   uuid.UUID
	TenantID uuid.UUID
	Repo     Repository
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

// Client will build a custom http.Client for Xero
func (c *Provider) Client(s *Session) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Base:   NewXeroTransport(s.TenantID),
			Source: oauth2.ReuseTokenSource(nil, NewTokenRefresher(s.Repo, s.Token, c, s.UserID)),
		},
	}
}

// NewClient method will return a new http.Client for use in our calls, using
// the TokenSource will even refresh the token if needed
func (c *Provider) NewClient(t *oauth2.Token) *http.Client {
	return c.conf.Client(c.ctx, t)
}
