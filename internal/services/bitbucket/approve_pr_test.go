package bitbucket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ApprovePR(t *testing.T) {
	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
				username, repoSlug, pullRequestID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"user": {
					"account_id": "123456",
					"display_name": "Test User",
					"nickname": "testuser",
					"username": "testuser",
					"uuid": "{58021780-82b6-4517-b153-0ae73ce3e4b4}",
					"type": "user"
				},
				"role": "REVIEWER",
				"approved": true,
				"state": "approved",
				"type": "participant"
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "Test User", result.User.DisplayName)
		assert.Equal(t, "testuser", result.User.Username)
		assert.Equal(t, "REVIEWER", result.Role)
		assert.True(t, result.Approved)
		assert.Equal(t, "approved", result.State)
		assert.Equal(t, "participant", result.Type)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"user": {
					"username": "testuser",
					"type": "user"
				},
				"approved": true,
				"type": "participant"
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "testuser", result.User.Username)
		assert.True(t, result.Approved)
		assert.Equal(t, "participant", result.Type)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

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
					"message": "Cannot approve your own pull request"
				}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "approve pull request failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New(faker.Sentence()),
		}

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		// Act
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		expectedError := fmt.Errorf("failed to get token: %w", mockTokenProvider.Err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
