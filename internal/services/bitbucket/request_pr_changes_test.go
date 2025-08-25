//go:build !release

package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_RequestPRChanges(t *testing.T) {
	t.Run("success requests changes", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedParticipant := Participant{
			User:     *NewRandomPullRequestAuthor(),
			Role:     "REVIEWER",
			Approved: false,
			State:    "changes_requested",
			Type:     "participant",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/request-changes", workspace, repoSlug, pullReqID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(expectedParticipant)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := RequestPRChangesParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		status, requestedAt, err := client.RequestPRChanges(t.Context(), mockTokenProvider, params)
		require.NoError(t, err)
		assert.Equal(t, "changes_requested", status)
		assert.WithinDuration(t, time.Now().UTC(), requestedAt, time.Second)
	})

	t.Run("handles API error", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": {"message": "Invalid request"}}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := RequestPRChangesParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		status, requestedAt, err := client.RequestPRChanges(t.Context(), mockTokenProvider, params)
		require.Error(t, err)
		assert.Equal(t, "", status)
		assert.True(t, requestedAt.IsZero())
		assert.Contains(t, err.Error(), "request changes failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000
		tokenErr := errors.New(faker.Sentence())

		mockTokenProvider := &MockTokenProvider{
			Err: tokenErr,
		}

		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		params := RequestPRChangesParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		status, requestedAt, err := client.RequestPRChanges(t.Context(), mockTokenProvider, params)
		require.Error(t, err)
		assert.Equal(t, "", status)
		assert.True(t, requestedAt.IsZero())
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
