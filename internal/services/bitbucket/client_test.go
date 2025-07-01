package bitbucket

import (
	"context"
	"testing"

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

// makeMockDepsWithTestName creates mock dependencies for testing with test name in the logger.
func makeMockDepsWithTestName(t *testing.T, baseURL string) ClientDeps {
	rootLogger := diag.RootTestLogger().With("test", t.Name())
	return ClientDeps{
		ClientFactory: httpservices.NewClientFactory(httpservices.ClientFactoryDeps{
			RootLogger: rootLogger,
		}),
		RootLogger: rootLogger,
		BaseURL:    baseURL,
	}
}
