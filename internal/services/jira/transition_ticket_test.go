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

func TestClient_TransitionTicket(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}

	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/issue/TEST-123/transitions", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), `"transition":{"id":"21"}`)
			assert.Contains(t, string(body), `"fields":{"resolution":{"name":"Done"}}`)

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

		// Create fields map
		fields := map[string]interface{}{
			"resolution": map[string]string{
				"name": "Done",
			},
		}

		// Execute the request
		err := client.TransitionTicket(t.Context(), mockTokenProvider, TransitionTicketParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-123",
			TransitionID: "21", // Done transition
			Fields:       fields,
		})

		// Verify the result
		require.NoError(t, err)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/issue/TEST-456/transitions", r.URL.Path)

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), `"transition":{"id":"11"}`)
			assert.NotContains(t, string(body), `"fields"`)

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
		err := client.TransitionTicket(t.Context(), mockTokenProvider, TransitionTicketParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-456",
			TransitionID: "11", // In Progress transition
			// No fields or updates
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
				"errorMessages": ["Transition '999' is not valid for issue 'TEST-123'"],
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

		// Execute the request with invalid transition ID
		err := client.TransitionTicket(t.Context(), mockTokenProvider, TransitionTicketParams{
			Domain:       "example", // This will be ignored since we're using the mock server URL
			TicketKey:    "TEST-123",
			TransitionID: "999", // Invalid transition ID
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transition ticket failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Create client with mock dependencies
		deps := makeMockDeps("https://example.atlassian.net/rest/api/3")
		client := NewClient(deps)

		// Setup token provider to return an error
		mockTokenProvider.Err = errors.New("token error")

		// Execute the request
		err := client.TransitionTicket(t.Context(), mockTokenProvider, TransitionTicketParams{
			Domain:       "example",
			TicketKey:    "TEST-123",
			TransitionID: "21",
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
