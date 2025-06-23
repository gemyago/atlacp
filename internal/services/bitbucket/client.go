package bitbucket

import (
	"context"
	"log/slog"
	"net/http"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"go.uber.org/dig"
)

// TokenProvider provides authentication tokens for Bitbucket API requests.
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// Client provides access to Bitbucket Cloud API operations.
type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     *slog.Logger
}

// ClientDeps contains dependencies for the Bitbucket client.
type ClientDeps struct {
	dig.In

	ClientFactory *httpservices.ClientFactory
	RootLogger    *slog.Logger
	BaseURL       string `name:"config.atlassian.bitbucket.baseUrl"`
}

// NewClient creates a new Bitbucket API client.
func NewClient(deps ClientDeps) *Client {
	return &Client{
		httpClient: deps.ClientFactory.CreateClient(),
		baseURL:    deps.BaseURL,
		logger:     deps.RootLogger.WithGroup("bitbucket-client"),
	}
}
