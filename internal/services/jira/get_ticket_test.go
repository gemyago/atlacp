package jira

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetTicket(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}

	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/issue/TEST-123", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "summary,description,status", r.URL.Query().Get("fields"))
			assert.Equal(t, "renderedFields,transitions", r.URL.Query().Get("expand"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"id": "10000",
				"key": "TEST-123",
				"self": "https://example.atlassian.net/rest/api/3/issue/10000",
				"fields": {
					"summary": "Test Issue",
					"description": "Test description",
					"status": {
						"id": "10000",
						"name": "To Do",
						"description": "Issue is open and ready for the assignee to start work on it.",
						"statusCategory": {
							"id": 2,
							"key": "new",
							"name": "To Do"
						},
						"self": "https://example.atlassian.net/rest/api/3/status/10000"
					},
					"priority": {
						"id": "3",
						"name": "Medium",
						"self": "https://example.atlassian.net/rest/api/3/priority/3",
						"iconUrl": "https://example.atlassian.net/images/icons/priorities/medium.svg"
					},
					"issuetype": {
						"id": "10001",
						"name": "Task",
						"description": "A task that needs to be done.",
						"iconUrl": "https://example.atlassian.net/secure/viewavatar?size=medium&avatarId=10318&avatarType=issuetype",
						"subtask": false,
						"self": "https://example.atlassian.net/rest/api/3/issuetype/10001"
					},
					"project": {
						"id": "10000",
						"key": "TEST",
						"name": "Test Project",
						"self": "https://example.atlassian.net/rest/api/3/project/10000"
					},
					"created": "2023-01-01T00:00:00.000Z",
					"updated": "2023-01-02T00:00:00.000Z",
					"labels": ["bug", "critical"]
				},
				"transitions": [
					{
						"id": "11",
						"name": "In Progress",
						"to": {
							"id": "10001",
							"name": "In Progress",
							"description": "This issue is being actively worked on."
						}
					},
					{
						"id": "21",
						"name": "Done",
						"to": {
							"id": "10002",
							"name": "Done",
							"description": "This issue is complete."
						}
					}
				]
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request
		result, err := client.GetTicket(t.Context(), mockTokenProvider, GetTicketParams{
			Domain:    "example", // This will be ignored since we're using the mock server URL
			TicketKey: "TEST-123",
			Fields:    []string{"summary", "description", "status"},
			Expand:    []string{"renderedFields", "transitions"},
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, "10000", result.ID)
		assert.Equal(t, "TEST-123", result.Key)
		assert.Equal(t, "Test Issue", result.Fields.Summary)
		assert.Equal(t, "Test description", result.Fields.Description)
		assert.Equal(t, "To Do", result.Fields.Status.Name)
		assert.Len(t, result.Transitions, 2)
		assert.Equal(t, "In Progress", result.Transitions[0].Name)
		assert.Equal(t, "Done", result.Transitions[1].Name)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"id": "10001",
				"key": "TEST-456",
				"self": "https://example.atlassian.net/rest/api/3/issue/10001",
				"fields": {
					"summary": "Minimal Issue"
				}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request with minimal parameters
		result, err := client.GetTicket(t.Context(), mockTokenProvider, GetTicketParams{
			Domain:    "example", // This will be ignored since we're using the mock server URL
			TicketKey: "TEST-456",
			// No optional parameters
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, "10001", result.ID)
		assert.Equal(t, "TEST-456", result.Key)
		assert.Equal(t, "Minimal Issue", result.Fields.Summary)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{
				"errorMessages": ["Issue does not exist or you do not have permission to see it."],
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
		_, err := client.GetTicket(t.Context(), mockTokenProvider, GetTicketParams{
			Domain:    "example", // This will be ignored since we're using the mock server URL
			TicketKey: "NONEXISTENT-999",
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get ticket failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Create client with mock dependencies
		deps := makeMockDeps("https://example.atlassian.net/rest/api/3")
		client := NewClient(deps)

		// Setup token provider to return an error
		mockTokenProvider.Err = errors.New("token error")

		// Execute the request
		_, err := client.GetTicket(t.Context(), mockTokenProvider, GetTicketParams{
			Domain:    "example",
			TicketKey: "TEST-123",
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
