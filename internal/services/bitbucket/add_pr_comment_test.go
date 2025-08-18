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

func TestClient_AddPRComment(t *testing.T) {
	t.Run("success with all parameters (inline comment)", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000
		commentText := faker.Sentence()
		filePath := "src/main.go"
		lineFrom := 10
		lineTo := 12

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedCommentID := faker.RandomUnixTime()
		expectedStatus := "success"
		expectedContent := &TaskContent{Raw: commentText}
		expectedComment := Comment{
			ID:        expectedCommentID,
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Content:   expectedContent,
			User:      &Account{DisplayName: faker.Name()},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments", workspace, repoSlug, pullReqID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, commentText, payload["content"].(map[string]interface{})["raw"])
			inline := payload["inline"].(map[string]interface{})
			assert.Equal(t, filePath, inline["path"])
			assert.Equal(t, float64(lineFrom), inline["from"])
			assert.Equal(t, float64(lineTo), inline["to"])

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(expectedComment)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := AddPRCommentParams{
			Workspace:   workspace,
			RepoSlug:    repoSlug,
			PullReqID:   pullReqID,
			CommentText: commentText,
			FilePath:    filePath,
			LineFrom:    lineFrom,
			LineTo:      lineTo,
		}

		commentID, status, err := client.AddPRComment(t.Context(), mockTokenProvider, params)
		require.NoError(t, err)
		assert.Equal(t, expectedCommentID, commentID)
		assert.Equal(t, expectedStatus, status)
	})

	t.Run("success with required parameters only (general comment)", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000
		commentText := faker.Sentence()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedCommentID := faker.RandomUnixTime()
		expectedStatus := "success"
		expectedContent := &TaskContent{Raw: commentText}
		expectedComment := Comment{
			ID:        expectedCommentID,
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			Content:   expectedContent,
			User:      &Account{DisplayName: faker.Name()},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments", workspace, repoSlug, pullReqID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			var payload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, commentText, payload["content"].(map[string]interface{})["raw"])
			assert.Nil(t, payload["inline"])

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(expectedComment)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := AddPRCommentParams{
			Workspace:   workspace,
			RepoSlug:    repoSlug,
			PullReqID:   pullReqID,
			CommentText: commentText,
		}

		commentID, status, err := client.AddPRComment(t.Context(), mockTokenProvider, params)
		require.NoError(t, err)
		assert.Equal(t, expectedCommentID, commentID)
		assert.Equal(t, expectedStatus, status)
	})

	t.Run("handles API error", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000
		commentText := faker.Sentence()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": {"message": "Invalid request"}}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := AddPRCommentParams{
			Workspace:   workspace,
			RepoSlug:    repoSlug,
			PullReqID:   pullReqID,
			CommentText: commentText,
		}

		commentID, status, err := client.AddPRComment(t.Context(), mockTokenProvider, params)
		require.Error(t, err)
		assert.Equal(t, int64(0), commentID)
		assert.Equal(t, "", status)
		assert.Contains(t, err.Error(), "add pull request comment failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000
		commentText := faker.Sentence()
		tokenErr := errors.New(faker.Sentence())

		mockTokenProvider := &MockTokenProvider{
			Err: tokenErr,
		}

		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		params := AddPRCommentParams{
			Workspace:   workspace,
			RepoSlug:    repoSlug,
			PullReqID:   pullReqID,
			CommentText: commentText,
		}

		commentID, status, err := client.AddPRComment(t.Context(), mockTokenProvider, params)
		require.Error(t, err)
		assert.Equal(t, int64(0), commentID)
		assert.Equal(t, "", status)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
