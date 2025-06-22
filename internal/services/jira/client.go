package jira

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"go.uber.org/dig"
)

// TokenProvider provides authentication tokens for Jira API requests.
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// Client provides access to Jira Cloud API operations.
type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     *slog.Logger
}

// ClientDeps contains dependencies for the Jira client.
type ClientDeps struct {
	dig.In

	ClientFactory *httpservices.ClientFactory
	RootLogger    *slog.Logger
	BaseURL       string `name:"config.atlassian.jira.baseUrl"`
}

// NewClient creates a new Jira API client.
func NewClient(deps ClientDeps) *Client {
	return &Client{
		httpClient: deps.ClientFactory.CreateClient(),
		baseURL:    deps.BaseURL,
		logger:     deps.RootLogger.WithGroup("jira-client"),
	}
}

// GetBaseURL returns the base URL with the domain replaced.
// The baseURL contains a placeholder {domain} that needs to be replaced with the actual domain.
func (c *Client) GetBaseURL(domain string) string {
	return strings.ReplaceAll(c.baseURL, "{domain}", domain)
}
