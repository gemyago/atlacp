package bitbucket

import (
	"context"

	"github.com/gemyago/atlacp/internal/diag"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
)

// MockTokenProvider is a simple mock implementation for testing.
type MockTokenProvider struct {
	Token string
	Err   error
}

// GetToken implements the TokenProvider interface for testing.
func (m *MockTokenProvider) GetToken(_ context.Context) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Token, nil
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
