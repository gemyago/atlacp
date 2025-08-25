package bitbucket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListPRComments(t *testing.T) {
	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Arrange
		workspace := "test-workspace-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := int64(rand.Intn(1000) + 1)

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		// Create test data for comments
		commentID1 := int64(rand.Intn(10000) + 1)
		commentID2 := int64(rand.Intn(10000) + 1)
		commentContent1 := "This is a general comment: " + faker.Sentence()
		commentContent2 := "This is an inline comment: " + faker.Sentence()
		authorDisplayName1 := "Author One " + faker.FirstName()
		authorDisplayName2 := "Author Two " + faker.FirstName()
		authorAccountID1 := faker.UUIDHyphenated()
		authorAccountID2 := faker.UUIDHyphenated()
		filePath := "src/main.go"
		lineFrom := rand.Intn(100) + 1
		lineTo := lineFrom + rand.Intn(10) + 1

		createdOn1 := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		createdOn2 := time.Now().UTC().Add(-12 * time.Hour).Truncate(time.Second)
		updatedOn1 := createdOn1.Add(time.Hour).Truncate(time.Second)
		updatedOn2 := createdOn2.Add(30 * time.Minute).Truncate(time.Second)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments", workspace, repoSlug, prID), r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return successful response with both general and inline comments
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"size": 2,
				"page": 1,
				"pagelen": 10,
				"values": [
					{
						"id": %d,
						"content": {
							"raw": "%s",
							"markup": "markdown",
							"html": "<p>%s</p>"
						},
						"user": {
							"account_id": "%s",
							"display_name": "%s",
							"nickname": "author1",
							"type": "user"
						},
						"created_on": "%s",
						"updated_on": "%s",
						"deleted": false
					},
					{
						"id": %d,
						"content": {
							"raw": "%s",
							"markup": "markdown",
							"html": "<p>%s</p>"
						},
						"user": {
							"account_id": "%s",
							"display_name": "%s",
							"nickname": "author2",
							"type": "user"
						},
						"created_on": "%s",
						"updated_on": "%s",
						"deleted": false,
						"inline": {
							"path": "%s",
							"from": %d,
							"to": %d
						}
					}
				]
			}`, commentID1, commentContent1, commentContent1, authorAccountID1, authorDisplayName1,
				createdOn1.Format(time.RFC3339), updatedOn1.Format(time.RFC3339),
				commentID2, commentContent2, commentContent2, authorAccountID2, authorDisplayName2,
				createdOn2.Format(time.RFC3339), updatedOn2.Format(time.RFC3339),
				filePath, lineFrom, lineTo)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := ListPRCommentsParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
		}

		// Act
		result, err := client.ListPRComments(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Values, 2)

		// Verify first comment (general comment)
		comment1 := result.Values[0]
		assert.Equal(t, commentID1, comment1.ID)
		assert.Equal(t, commentContent1, comment1.Content.Raw)
		assert.Equal(t, createdOn1, comment1.CreatedOn)
		assert.Equal(t, updatedOn1, comment1.UpdatedOn)
		assert.NotNil(t, comment1.Author)
		assert.Equal(t, authorAccountID1, comment1.Author.AccountID)
		assert.Equal(t, authorDisplayName1, comment1.Author.DisplayName)
		assert.Nil(t, comment1.Inline) // General comment has no inline context
		assert.Nil(t, comment1.Parent) // No parent comment

		// Verify second comment (inline comment)
		comment2 := result.Values[1]
		assert.Equal(t, commentID2, comment2.ID)
		assert.Equal(t, commentContent2, comment2.Content.Raw)
		assert.Equal(t, createdOn2, comment2.CreatedOn)
		assert.Equal(t, updatedOn2, comment2.UpdatedOn)
		assert.NotNil(t, comment2.Author)
		assert.Equal(t, authorAccountID2, comment2.Author.AccountID)
		assert.Equal(t, authorDisplayName2, comment2.Author.DisplayName)
		assert.NotNil(t, comment2.Inline) // Inline comment has inline context
		assert.Equal(t, filePath, comment2.Inline.Path)
		assert.Equal(t, lineFrom, comment2.Inline.From)
		assert.Equal(t, lineTo, comment2.Inline.To)
		assert.Nil(t, comment2.Parent) // No parent comment
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange
		workspace := "test-workspace-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := int64(rand.Intn(1000) + 1)

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return empty comments response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"size": 0,
				"page": 1,
				"pagelen": 10,
				"values": []
			}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := ListPRCommentsParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
		}

		// Act
		result, err := client.ListPRComments(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Values)
	})

	t.Run("api error", func(t *testing.T) {
		// Arrange
		workspace := "test-workspace-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := int64(rand.Intn(1000) + 1)

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{
				"type": "error",
				"error": {
					"message": "Pull request not found"
				}
			}`)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := ListPRCommentsParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
		}

		// Act
		result, err := client.ListPRComments(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "list pull request comments failed")
	})

	t.Run("token error", func(t *testing.T) {
		// Arrange
		workspace := "test-workspace-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prID := int64(rand.Intn(1000) + 1)

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New("token retrieval failed"),
		}

		deps := makeMockDepsWithTestName(t, "https://api.bitbucket.org/2.0")
		client := NewClient(deps)

		params := ListPRCommentsParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PRID:      prID,
		}

		// Act
		result, err := client.ListPRComments(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get token")
		assert.Contains(t, err.Error(), "token retrieval failed")
	})
}
