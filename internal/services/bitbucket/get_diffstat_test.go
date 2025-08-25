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

func TestClient_GetPRDiffStat(t *testing.T) {
	t.Run("success returns paginated diffstat", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat", username, repoSlug, prID), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"pagelen": 2,
				"values": [
					{
						"type": "diffstat",
						"status": "modified",
						"lines_added": 10,
						"lines_removed": 2,
						"path": "file1.go"
					},
					{
						"type": "diffstat",
						"status": "added",
						"lines_added": 20,
						"lines_removed": 0,
						"path": "file2.go"
					}
				],
				"page": 1,
				"size": 2
			}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		result, err := client.GetPRDiffStat(t.Context(), mockTokenProvider, GetPRDiffStatParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Size)
		assert.Len(t, result.Values, 2)
		assert.Equal(t, "file1.go", result.Values[0].Path)
		assert.Equal(t, "file2.go", result.Values[1].Path)
		t.Run("handles optional parameters: FilePaths, Context, Account", func(t *testing.T) {
			username2 := "test-user-" + faker.Word()
			repoSlug2 := "test-repo-" + faker.Word()
			prID2 := rand.Intn(1000) + 1
			mockTokenProvider2 := &MockTokenProvider{
				TokenType:  "Bearer",
				TokenValue: faker.UUIDHyphenated(),
			}
			filePaths := []string{"foo.go", "bar.go"}
			contextLines := 5
			account := "acc-456"
			server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(
					t,
					fmt.Sprintf(
						"/repositories/%s/%s/pullrequests/%d/diffstat",
						username2, repoSlug2, prID2,
					),
					r.URL.Path,
				)
				assert.Equal(t, "Bearer "+mockTokenProvider2.TokenValue, r.Header.Get("Authorization"))
				q := r.URL.Query()
				assert.ElementsMatch(t, filePaths, q["path"])
				assert.Equal(t, strconv.Itoa(contextLines), q.Get("context"))
				assert.Equal(t, account, q.Get("account_id"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"pagelen":1,"values":[],"page":1,"size":0}`)
			}))
			defer server2.Close()
			deps2 := makeMockDepsWithTestName(t, server2.URL)
			client2 := NewClient(deps2)
			result2, err2 := client2.GetPRDiffStat(t.Context(), mockTokenProvider2, GetPRDiffStatParams{
				RepoOwner: username2,
				RepoName:  repoSlug2,
				PRID:      prID2,
				FilePaths: filePaths,
				Context:   &contextLines,
				Account:   &account,
			})
			require.NoError(t, err2)
			require.NotNil(t, result2)
			assert.Equal(t, 0, result2.Size)
			assert.Empty(t, result2.Values)
		})
	})

	t.Run("handles paginated diffstat via next field", func(t *testing.T) {
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := rand.Intn(1000) + 1
		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}
		var server *httptest.Server
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat", username, repoSlug, prID):
				nextURL := server.URL + "/next"
				page1 := fmt.Sprintf(
					`{"pagelen":1,"values":[{"type":"diffstat","status":"modified","lines_added":1,`+
						`"lines_removed":0,"path":"file1.go"}],`+
						`"page":1,"size":2,"next":"%s"}`,
					nextURL,
				)
				fmt.Fprint(w, page1)
			case "/next":
				page2 := `{"pagelen":1,"values":[{"type":"diffstat","status":"added","lines_added":2,` +
					`"lines_removed":0,"path":"file2.go"}],` +
					`"page":2,"size":2}`
				fmt.Fprint(w, page2)
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
		}
		server = httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)
		result, err := client.GetPRDiffStat(t.Context(), mockTokenProvider, GetPRDiffStatParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Size)
		assert.Len(t, result.Values, 2)
		assert.Equal(t, "file1.go", result.Values[0].Path)
		assert.Equal(t, "file2.go", result.Values[1].Path)
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

		result, err := client.GetPRDiffStat(t.Context(), mockTokenProvider, GetPRDiffStatParams{
			RepoOwner: username,
			RepoName:  repoSlug,
			PRID:      prID,
		})

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "get diffstat failed")
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

		result, err := client.GetPRDiffStat(t.Context(), mockTokenProvider, GetPRDiffStatParams{
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
			params GetPRDiffStatParams
			errMsg string
		}{
			{"missing RepoOwner", GetPRDiffStatParams{RepoName: "repo", PRID: 1}, "RepoOwner is required"},
			{"missing RepoName", GetPRDiffStatParams{RepoOwner: "owner", PRID: 1}, "RepoName is required"},
			{"missing PRID", GetPRDiffStatParams{RepoOwner: "owner", RepoName: "repo"}, "PRID is required and must be non-zero"},
		}
		client := NewClient(makeMockDepsWithTestName(t, "http://example.com"))
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				result, err := client.GetPRDiffStat(t.Context(), mockTokenProvider, tc.params)
				require.Error(t, err)
				assert.Nil(t, result)
				assert.ErrorContains(t, err, tc.errMsg)
			})
		}
	})
}
