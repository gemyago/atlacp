package bitbucket

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPRComment(t *testing.T) {
	t.Run("success returns comment with resolution", func(t *testing.T) {
		workspace := "ws-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		prID := int64(100 + rand.IntN(9000))
		commentID := int64(200 + rand.IntN(9000))
		raw := faker.Sentence()
		displayName := faker.Name()
		created := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
		updated := created.Add(30 * time.Minute)

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments/%d",
				workspace, repoSlug, prID, commentID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"content": {"raw": %q},
				"user": {"display_name": %q, "type": "user"},
				"created_on": %q,
				"updated_on": %q,
				"resolution": {"resolved": true}
			}`, commentID, raw, displayName, created.Format(time.RFC3339), updated.Format(time.RFC3339))
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		got, err := client.GetPRComment(t.Context(), mockTokenProvider, GetPRCommentParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
			CommentID: commentID,
		})

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, commentID, got.ID)
		assert.Equal(t, raw, got.Content.Raw)
		require.NotNil(t, got.Author)
		assert.Equal(t, displayName, got.Author.DisplayName)
		resolved, known := ResolvedStateFromResolutionJSON(got.Resolution)
		assert.True(t, resolved)
		assert.True(t, known)
	})

	t.Run("http error", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		prID := int64(1 + rand.IntN(100))
		commentID := int64(1 + rand.IntN(100))
		errMsg := "not found"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error":{"message":%q}}`, errMsg)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		got, err := client.GetPRComment(t.Context(), mockTokenProvider, GetPRCommentParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
			CommentID: commentID,
		})

		require.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "get pull request comment failed")
	})

	t.Run("handles token error", func(t *testing.T) {
		tokenErr := errors.New("token error")
		mockTokenProvider := &MockTokenProvider{Err: tokenErr}

		deps := makeMockDepsWithTestName(t, "http://dummy-url")
		client := NewClient(deps)

		got, err := client.GetPRComment(t.Context(), mockTokenProvider, GetPRCommentParams{
			Workspace: faker.Username(),
			RepoSlug:  faker.Username(),
			PRID:      1,
			CommentID: 2,
		})

		require.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
