package bitbucket

import (
	"context"
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

type errorRoundTripper struct {
	err error
}

func (e *errorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, e.err
}

type errorReadCloser struct{}

func (e *errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}
func (e *errorReadCloser) Close() error { return nil }

func TestClient_GetFileContent_EdgeCases(t *testing.T) {
	t.Run("http client returns error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		client := &Client{
			httpClient: &http.Client{Transport: &errorRoundTripper{err: errors.New("network error")}},
			baseURL:    "http://localhost",
			logger:     nil,
		}

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

	t.Run("malformed baseURL triggers url.Parse error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"
		account := "acc"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		// baseURL missing scheme
		client := &Client{
			httpClient: http.DefaultClient,
			baseURL:    "://bad-url",
			logger:     nil,
		}

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
			Account:    &account,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to parse request URL")
	})

	t.Run("context canceled triggers request creation error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		client := &Client{
			httpClient: http.DefaultClient,
			baseURL:    "http://localhost",
			logger:     nil,
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel() // immediately cancel

		result, err := client.GetFileContent(ctx, mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		require.ErrorContains(t, err, "get file content failed")
		assert.ErrorContains(t, err, "context canceled")
	})

	t.Run("io.ReadAll returns error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		client := &Client{
			httpClient: &http.Client{
				Transport: roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       &errorReadCloser{},
						Header:     make(http.Header),
					}, nil
				}),
			},
			baseURL: "http://localhost",
			logger:  nil,
		}

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to read file content response")
	})

	t.Run("Account param is non-nil but empty string (should not add ?account=)", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		commit := faker.UUIDHyphenated()
		filePath := "src/" + faker.Word() + ".go"
		account := ""

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "dummy content")
		}))
		defer server.Close()

		client := &Client{
			httpClient: http.DefaultClient,
			baseURL:    server.URL,
			logger:     nil,
		}

		result, err := client.GetFileContent(t.Context(), mockTokenProvider, GetFileContentParams{
			RepoOwner:  username,
			RepoName:   repoSlug,
			CommitHash: commit,
			FilePath:   filePath,
			Account:    &account,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "", gotQuery)
	})
}

// roundTripperFunc allows inline definition of http.RoundTripper for tests.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
