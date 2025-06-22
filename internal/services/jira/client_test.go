package jira

import (
	"context"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/stretchr/testify/assert"
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

func TestClient_GetBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		domain      string
		expectedURL string
	}{
		{
			name:        "replace domain placeholder",
			baseURL:     "https://{domain}.atlassian.net/rest/api/3",
			domain:      "example",
			expectedURL: "https://example.atlassian.net/rest/api/3",
		},
		{
			name:        "replace domain placeholder with multiple occurrences",
			baseURL:     "https://{domain}.atlassian.net/{domain}/api/3",
			domain:      "test",
			expectedURL: "https://test.atlassian.net/test/api/3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := makeMockDeps(tt.baseURL)
			client := NewClient(deps)

			result := client.GetBaseURL(tt.domain)
			assert.Equal(t, tt.expectedURL, result)
		})
	}
}
