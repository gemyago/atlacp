package jira

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ManageLabels(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}

	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/issue/TEST-123", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			bodyStr := string(body)
			assert.Contains(t, bodyStr, `"update":`)
			assert.Contains(t, bodyStr, `"labels":`)
			assert.Contains(t, bodyStr, `"add":"bug"`)
			assert.Contains(t, bodyStr, `"add":"critical"`)
			assert.Contains(t, bodyStr, `"remove":"wontfix"`)

			// Return successful response (empty 204 No Content)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request
		err := client.ManageLabels(t.Context(), mockTokenProvider, ManageLabelsParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-123",
			AddLabels:    []string{"bug", "critical"},
			RemoveLabels: []string{"wontfix"},
		})

		// Verify the result
		require.NoError(t, err)
	})

	t.Run("success with required parameters only (add labels only)", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/issue/TEST-456", r.URL.Path)

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			bodyStr := string(body)
			assert.Contains(t, bodyStr, `"add":"enhancement"`)
			assert.NotContains(t, bodyStr, `"remove"`)

			// Return successful response (empty 204 No Content)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request with minimal parameters
		err := client.ManageLabels(t.Context(), mockTokenProvider, ManageLabelsParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-456",
			AddLabels:    []string{"enhancement"},
			RemoveLabels: nil, // No labels to remove
		})

		// Verify the result
		require.NoError(t, err)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"errorMessages": ["Issue TEST-999 does not exist or you do not have permission to edit it."],
				"errors": {}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request with non-existent ticket key
		err := client.ManageLabels(t.Context(), mockTokenProvider, ManageLabelsParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-999",
			AddLabels:    []string{"bug"},
			RemoveLabels: nil,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "manage labels failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Create client with mock dependencies
		deps := makeMockDeps("https://example.atlassian.net/rest/api/3")
		client := NewClient(deps)

		// Setup token provider to return an error
		mockTokenProvider.Err = errors.New("token error")

		// Execute the request
		err := client.ManageLabels(t.Context(), mockTokenProvider, ManageLabelsParams{
			Domain:       "example",
			TicketKey:    "TEST-123",
			AddLabels:    []string{"bug"},
			RemoveLabels: nil,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
