package bitbucket

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_DeclinePR(t *testing.T) {
	t.Run("declines a pull request with no request body", func(t *testing.T) {
		username := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1
		tokenProvider := &MockTokenProvider{TokenType: "Bearer", TokenValue: faker.UUIDHyphenated()}
		requests := 0

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			assert.Equal(t, http.MethodPost, r.Method)
			expectedPath := fmt.Sprintf(
				"/repositories/%s/%s/pullrequests/%d/decline", username, repoSlug, pullRequestID,
			)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+tokenProvider.TokenValue, r.Header.Get("Authorization"))
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Empty(t, body)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(
				w,
				`{"id": %d, "title": "Declined %s", "state": "DECLINED", "type": "pullrequest"}`,
				pullRequestID,
				faker.Word(),
			)
		}))
		defer server.Close()

		client := NewClient(makeMockDepsWithTestName(t, server.URL))
		pullRequest, err := client.DeclinePR(t.Context(), tokenProvider, DeclinePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		require.NoError(t, err)
		assert.Equal(t, 1, requests)
		assert.Equal(t, pullRequestID, pullRequest.ID)
		assert.Equal(t, "DECLINED", pullRequest.State)
	})

	t.Run("wraps API errors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client := NewClient(makeMockDepsWithTestName(t, server.URL))
		tokenProvider := &MockTokenProvider{TokenValue: faker.UUIDHyphenated()}
		pullRequest, err := client.DeclinePR(t.Context(), tokenProvider, DeclinePRParams{
			Username:      faker.Word(),
			RepoSlug:      faker.Word(),
			PullRequestID: rand.Intn(1000) + 1,
		})

		require.Error(t, err)
		assert.Nil(t, pullRequest)
		assert.ErrorContains(t, err, "decline pull request failed")
	})

	t.Run("wraps token provider errors", func(t *testing.T) {
		tokenErr := errors.New(faker.Sentence())
		client := NewClient(makeMockDepsWithTestName(t, "http://example.com"))

		pullRequest, err := client.DeclinePR(t.Context(), &MockTokenProvider{Err: tokenErr}, DeclinePRParams{})

		require.Error(t, err)
		assert.Nil(t, pullRequest)
		require.ErrorIs(t, err, tokenErr)
		assert.ErrorContains(t, err, "failed to get token")
	})
}
