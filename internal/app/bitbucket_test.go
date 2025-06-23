package app

import (
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for BitbucketService.
func TestBitbucketService(t *testing.T) {
	// Helper function to create mock dependencies
	makeMockDeps := func() BitbucketServiceDeps {
		return BitbucketServiceDeps{
			Client:       NewMockBitbucketClient(t),
			AccountsRepo: NewMockAtlassianAccountsRepository(t),
			RootLogger:   diag.RootTestLogger(),
		}
	}

	t.Run("CreatePR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoOwner:    "test-owner",
				RepoName:     "test-repo",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})

	t.Run("ReadPR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			// Act
			result, err := service.ReadPR(t.Context(), BitbucketReadPRParams{
				RepoOwner:     "test-owner",
				RepoName:      "test-repo",
				PullRequestID: 123,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})

	t.Run("UpdatePR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				RepoOwner:     "test-owner",
				RepoName:      "test-repo",
				PullRequestID: 123,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})

	t.Run("ApprovePR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			// Act
			result, err := service.ApprovePR(t.Context(), BitbucketApprovePRParams{
				RepoOwner:     "test-owner",
				RepoName:      "test-repo",
				PullRequestID: 123,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})

	t.Run("MergePR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				RepoOwner:     "test-owner",
				RepoName:      "test-repo",
				PullRequestID: 123,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})
}
