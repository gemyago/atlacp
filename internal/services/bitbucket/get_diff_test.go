package bitbucket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPRDiff(t *testing.T) {
	t.Run("success returns diff content", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		expectedDiff := "diff --git a/file1.go b/file1.go\n..."

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff", username, repoSlug, prID), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, expectedDiff)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedDiff, string(*result))
	})

	t.Run("handles error reading diff response body", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1
		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			// Hijack the connection to simulate a read error
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
		}))
		defer server.Close()
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)
		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to read diff response")
	})

	t.Run("handles paginated diff via X-Next-Page header", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1
		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}
		page1 := "diff --git a/file1.go b/file1.go\n..."
		page2 := "diff --git a/file2.go b/file2.go\n..."
		var server *httptest.Server
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			switch r.URL.Path {
			case fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff", username, repoSlug, prID):
				// First request
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))
				w.Header().Set("X-Next-Page", server.URL+"/next")
				fmt.Fprint(w, page1)
			case "/next":
				fmt.Fprint(w, page2)
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
		}
		server = httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)
		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, string(*result), page1)
		assert.Contains(t, string(*result), page2)
	})

	t.Run("handles optional parameters: FilePaths, Context, Account", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1
		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}
		filePaths := []string{"foo.go", "bar.go"}
		contextLines := 7
		account := "acc-123"
		expectedDiff := "diff --git a/foo.go b/foo.go\n..."

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff", username, repoSlug, prID), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))
			// Check query params
			q := r.URL.Query()
			assert.ElementsMatch(t, filePaths, q["path"])
			assert.Equal(t, strconv.Itoa(contextLines), q.Get("context"))
			// Check header
			assert.Equal(t, account, r.Header.Get("X-Atlassian-Account"))
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, expectedDiff)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
			FilePaths: filePaths,
			Context:   &contextLines,
			Account:   &account,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedDiff, string(*result))
	})

	t.Run("handles API error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1

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

		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "get diff failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New(faker.Sentence()),
		}

		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		result, err := client.GetPRDiff(t.Context(), mockTokenProvider, GetPRDiffParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "failed to get token")
	})
	t.Run("parameter validation errors", func(t *testing.T) {
		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}
		tests := []struct {
			name   string
			params GetPRDiffParams
			errMsg string
		}{
			{"missing RepoOwner", GetPRDiffParams{RepoName: "repo", PRID: 1}, "RepoOwner is required"},
			{"missing RepoName", GetPRDiffParams{RepoOwner: "owner", PRID: 1}, "RepoName is required"},
			{"missing PRID", GetPRDiffParams{RepoOwner: "owner", RepoName: "repo"}, "PRID is required and must be non-zero"},
		}
		client := NewClient(makeMockDepsWithTestName(t, "http://example.com"))
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				result, err := client.GetPRDiff(t.Context(), mockTokenProvider, tc.params)
				require.Error(t, err)
				assert.Nil(t, result)
				assert.ErrorContains(t, err, tc.errMsg)
			})
		}
	})
}
