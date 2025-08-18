package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetFileContent(t *testing.T) {
	t.Run("success returns file content", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedContent := "package main\n\nfunc main() {}\n"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/src/%s/%s", username, repoSlug, commit, filePath), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, expectedContent)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedContent, result.Content)
	})
	t.Run("handles optional Account parameter", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"
		account := "acc-" + faker.Word()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedContent := "package main\n\nfunc main() {}\n"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/src/%s/%s", username, repoSlug, commit, filePath), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))
			assert.Equal(t, account, r.URL.Query().Get("account"))

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, expectedContent)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
			Account:    &account,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedContent, result.Content)
	})

	t.Run("handles API error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error": {"message": "Not found"}}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "get file content failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New(faker.Sentence()),
		}

		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to get token")
	})
}
