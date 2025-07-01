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

func TestClient_CreatePR(t *testing.T) {
	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prTitle := "Test PR " + faker.Sentence()
		prDescription := faker.Paragraph()
		sourceBranch := "feature-" + faker.Word()
		targetBranch := "main-" + faker.Word()
		commitHash := faker.UUIDHyphenated()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		createdOn := time.Now().UTC().Truncate(time.Second)
		createdOnStr := createdOn.Format(time.RFC3339)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests", username, repoSlug), r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{
				"id": %d,
				"title": "%s",
				"description": "%s",
				"state": "OPEN",
				"author": {
					"account_id": "123456",
					"display_name": "Test User",
					"nickname": "testuser",
					"username": "testuser",
					"uuid": "{58021780-82b6-4517-b153-0ae73ce3e4b4}",
					"type": "user"
				},
				"source": {
					"branch": {
						"name": "%s"
					},
					"commit": {
						"hash": "%s"
					},
					"repository": {
						"full_name": "%s/%s",
						"name": "%s",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"destination": {
					"branch": {
						"name": "%s"
					},
					"repository": {
						"full_name": "%s/%s",
						"name": "%s",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"close_source_branch": true,
				"comment_count": 0,
				"task_count": 0,
				"created_on": "%s",
				"updated_on": "%s",
				"draft": false
			}`, rand.Intn(1000)+1, prTitle, prDescription, sourceBranch, commitHash,
				username, repoSlug, repoSlug, targetBranch, username, repoSlug, repoSlug,
				createdOnStr, createdOnStr)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.CreatePR(t.Context(), mockTokenProvider, CreatePRParams{
			Username: username,
			RepoSlug: repoSlug,
			Request: &PullRequest{
				Title:       prTitle,
				Description: prDescription,
				Source: PullRequestSource{
					Branch: PullRequestBranch{
						Name: sourceBranch,
					},
				},
				Destination: &PullRequestDestination{
					Branch: PullRequestBranch{
						Name: targetBranch,
					},
				},
				CloseSourceBranch: true,
			},
		})

		// Assert
		require.NoError(t, err)
		assert.NotZero(t, result.ID)
		assert.Equal(t, prTitle, result.Title)
		assert.Equal(t, prDescription, result.Description)
		assert.Equal(t, "OPEN", result.State)
		assert.Equal(t, sourceBranch, result.Source.Branch.Name)
		assert.Equal(t, targetBranch, result.Destination.Branch.Name)
		assert.Equal(t, commitHash, result.Source.Commit.Hash)
		assert.True(t, result.CloseSourceBranch)
		assert.Equal(t, createdOn, result.CreatedOn.UTC())
		assert.Equal(t, createdOn, result.UpdatedOn.UTC())
		assert.False(t, *result.Draft)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		prTitle := "Minimal PR " + faker.Word()
		sourceBranch := "feature-" + faker.Word()
		prID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{
				"id": %d,
				"title": "%s",
				"state": "OPEN",
				"source": {
					"branch": {
						"name": "%s"
					}
				}
			}`, prID, prTitle, sourceBranch)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.CreatePR(t.Context(), mockTokenProvider, CreatePRParams{
			Username: username,
			RepoSlug: repoSlug,
			Request: &PullRequest{
				Title: prTitle,
				Source: PullRequestSource{
					Branch: PullRequestBranch{
						Name: sourceBranch,
					},
				},
			},
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, prID, result.ID)
		assert.Equal(t, prTitle, result.Title)
		assert.Equal(t, "OPEN", result.State)
		assert.Equal(t, sourceBranch, result.Source.Branch.Name)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		sourceBranch := "feature-" + faker.Word()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"error": {
					"message": "Invalid request"
				}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.CreatePR(t.Context(), mockTokenProvider, CreatePRParams{
			Username: username,
			RepoSlug: repoSlug,
			Request: &PullRequest{
				// Missing required title field
				Source: PullRequestSource{
					Branch: PullRequestBranch{
						Name: sourceBranch,
					},
				},
			},
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "create pull request failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		sourceBranch := "feature-" + faker.Word()

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New(faker.Sentence()),
		}

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		// Act
		result, err := client.CreatePR(t.Context(), mockTokenProvider, CreatePRParams{
			Username: username,
			RepoSlug: repoSlug,
			Request: &PullRequest{
				Title: faker.Sentence(),
				Source: PullRequestSource{
					Branch: PullRequestBranch{
						Name: sourceBranch,
					},
				},
			},
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		expectedError := fmt.Errorf("failed to get token: %w", mockTokenProvider.Err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
