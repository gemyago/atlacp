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

func TestClient_ResolvePRComment(t *testing.T) {
	t.Run("success returns resolution", func(t *testing.T) {
		workspace := "ws-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		prID := int64(100 + rand.IntN(9000))
		commentID := int64(200 + rand.IntN(9000))
		displayName := faker.Name()
		createdOn := time.Now().UTC().Add(-time.Minute).Truncate(time.Second)

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments/%d/resolve",
				workspace, repoSlug, prID, commentID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"type": "resolution",
				"user": {"type": "user", "display_name": %q},
				"created_on": %q
			}`, displayName, createdOn.Format(time.RFC3339))
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		got, err := client.ResolvePRComment(t.Context(), mockTokenProvider, ResolvePRCommentParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
			CommentID: commentID,
		})

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "resolution", got.Type)
		require.NotNil(t, got.User)
		assert.Equal(t, displayName, got.User.DisplayName)
		assert.True(t, got.CreatedOn.Equal(createdOn))
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

		got, err := client.ResolvePRComment(t.Context(), mockTokenProvider, ResolvePRCommentParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
			CommentID: commentID,
		})

		require.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "resolve pull request comment failed")
	})

	t.Run("handles token error", func(t *testing.T) {
		tokenErr := errors.New("token error")
		mockTokenProvider := &MockTokenProvider{Err: tokenErr}

		deps := makeMockDepsWithTestName(t, "http://dummy-url")
		client := NewClient(deps)

		got, err := client.ResolvePRComment(t.Context(), mockTokenProvider, ResolvePRCommentParams{
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
