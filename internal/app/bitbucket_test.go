package app

import (
	"errors"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		t.Run("successfully creates pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, _ := deps.Client.(*MockBitbucketClient)
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			// Create test data
			account := &AtlassianAccount{
				Name:    "default",
				Default: true,
				Bitbucket: &BitbucketAccount{
					Token:     "test-token",
					Workspace: "test-workspace",
				},
			}

			expectedPR := &bitbucket.PullRequest{
				ID:          123,
				Title:       "Test PR",
				Description: "Test Description",
				Source: bitbucket.PullRequestSource{
					Branch: bitbucket.PullRequestBranch{
						Name: "feature-branch",
					},
				},
				Destination: &bitbucket.PullRequestDestination{
					Branch: bitbucket.PullRequestBranch{
						Name: "main",
					},
				},
				CloseSourceBranch: true,
			}

			// Mock the accounts repo to return our test account
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(account, nil)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, account.Bitbucket.Workspace, params.Username)
					assert.Equal(t, "test-repo", params.RepoSlug)

					// Verify the PR request
					pr := params.Request
					assert.Equal(t, "Test PR", pr.Title)
					assert.Equal(t, "Test Description", pr.Description)
					assert.Equal(t, "feature-branch", pr.Source.Branch.Name)
					assert.Equal(t, "main", pr.Destination.Branch.Name)
					assert.True(t, pr.CloseSourceBranch)
					assert.Empty(t, pr.Reviewers) // No reviewers in this test

					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:          "test-repo",
				Title:             "Test PR",
				Description:       "Test Description",
				SourceBranch:      "feature-branch",
				DestBranch:        "main",
				CloseSourceBranch: true,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully creates pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, _ := deps.Client.(*MockBitbucketClient)
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			account := &AtlassianAccount{
				Name:    accountName,
				Default: false,
				Bitbucket: &BitbucketAccount{
					Token:     "custom-token",
					Workspace: "custom-workspace",
				},
			}

			expectedPR := &bitbucket.PullRequest{
				ID:          123,
				Title:       "Test PR",
				Description: "Test Description",
			}

			// Mock the accounts repo to return our test account
			mockAccounts.EXPECT().
				GetAccountByName(mock.Anything, accountName).
				Return(account, nil)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify the workspace from the custom account is used
					assert.Equal(t, account.Bitbucket.Workspace, params.Username)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				AccountName:  accountName,
				RepoName:     "test-repo",
				Title:        "Test PR",
				Description:  "Test Description",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully creates pull request with reviewers", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, _ := deps.Client.(*MockBitbucketClient)
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			account := &AtlassianAccount{
				Name:    "default",
				Default: true,
				Bitbucket: &BitbucketAccount{
					Token:     "test-token",
					Workspace: "test-workspace",
				},
			}

			reviewers := []string{"user1", "user2"}

			// Mock the accounts repo
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(account, nil)

			// Mock the client
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify reviewers are properly added
					assert.Len(t, params.Request.Reviewers, len(reviewers))
					for i, reviewer := range params.Request.Reviewers {
						assert.Equal(t, reviewers[i], reviewer.Username)
					}
					return true
				})).
				Return(&bitbucket.PullRequest{ID: 123}, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     "test-repo",
				Title:        "Test PR",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
				Reviewers:    reviewers,
			})

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
		})

		t.Run("fails when account resolution fails", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			// Account not found
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(nil, ErrNoDefaultAccount)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     "test-repo",
				Title:        "Test PR",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrNoDefaultAccount)
		})

		t.Run("fails when named account not found", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			accountName := "non-existent"

			// Account not found
			mockAccounts.EXPECT().
				GetAccountByName(mock.Anything, accountName).
				Return(nil, ErrAccountNotFound)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				AccountName:  accountName,
				RepoName:     "test-repo",
				Title:        "Test PR",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrAccountNotFound)
		})

		t.Run("fails when account has no Bitbucket config", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			// Account with no Bitbucket config
			account := &AtlassianAccount{
				Name:    "default",
				Default: true,
				Jira: &JiraAccount{
					Token:  "jira-token",
					Domain: "test-domain",
				},
				// No Bitbucket config
			}

			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(account, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     "test-repo",
				Title:        "Test PR",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "bitbucket configuration not found")
		})

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketCreatePRParams
				errMsg string
			}{
				{
					name: "missing repo owner and name",
					params: BitbucketCreatePRParams{
						Title:        "Test PR",
						SourceBranch: "feature",
						DestBranch:   "main",
					},
					errMsg: "repository name is required",
				},
				{
					name: "missing title",
					params: BitbucketCreatePRParams{
						RepoName:     "test-repo",
						SourceBranch: "feature",
						DestBranch:   "main",
					},
					errMsg: "title is required",
				},
				{
					name: "missing source branch",
					params: BitbucketCreatePRParams{
						RepoName:   "test-repo",
						Title:      "Test PR",
						DestBranch: "main",
					},
					errMsg: "source branch is required",
				},
				{
					name: "missing destination branch",
					params: BitbucketCreatePRParams{
						RepoName:     "test-repo",
						Title:        "Test PR",
						SourceBranch: "feature",
					},
					errMsg: "destination branch is required",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Act
					result, err := service.CreatePR(t.Context(), tc.params)

					// Assert
					assert.Nil(t, result)
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})

		t.Run("fails when client returns error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, _ := deps.Client.(*MockBitbucketClient)
			mockAccounts, _ := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			service := NewBitbucketService(deps)

			account := &AtlassianAccount{
				Name:    "default",
				Default: true,
				Bitbucket: &BitbucketAccount{
					Token:     "test-token",
					Workspace: "test-workspace",
				},
			}

			clientErr := errors.New("client error")

			// Mock account repo
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(account, nil)

			// Mock client error
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, clientErr)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     "test-repo",
				Title:        "Test PR",
				SourceBranch: "feature-branch",
				DestBranch:   "main",
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, clientErr)
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
