package bitbucket

import (
	"context"

	"github.com/gemyago/atlacp/internal/diag"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// MockTokenProvider is a simple mock implementation for testing.
type MockTokenProvider struct {
	TokenType  string
	TokenValue string
	Err        error
}

// GetToken implements the TokenProvider interface for testing.
func (m *MockTokenProvider) GetToken(_ context.Context) (middleware.Token, error) {
	if m.Err != nil {
		return middleware.Token{}, m.Err
	}
	return middleware.Token{Type: m.TokenType, Value: m.TokenValue}, nil
}

// makeMockDeps creates mock dependencies for testing.
func makeMockDeps(baseURL string) ClientDeps {
	return ClientDeps{
		ClientFactory: httpservices.NewClientFactory(httpservices.ClientFactoryDeps{
			RootLogger: diag.RootTestLogger(),
		}),
		RootLogger: diag.RootTestLogger(),
		BaseURL:    baseURL,
	}
}
