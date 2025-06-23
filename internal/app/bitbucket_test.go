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
			mockClient, ok := deps.Client.(*MockBitbucketClient)
			require.True(t, ok, "Client should be a MockBitbucketClient")
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			// Create test data
			account := NewRandomAtlassianAccount(
				WithAtlassianAccountName("default"),
				WithAtlassianAccountDefault(true),
				WithAtlassianAccountBitbucket(),
			)

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Mock the accounts repo to return our test account
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(&account, nil)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, account.Bitbucket.Workspace, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)

					// Verify the PR request matches expected values
					expectedRequest := &bitbucket.PullRequest{
						Title:             expectedPR.Title,
						Description:       expectedPR.Description,
						CloseSourceBranch: true,
						Source: bitbucket.PullRequestSource{
							Branch: bitbucket.PullRequestBranch{
								Name: expectedPR.Source.Branch.Name,
							},
						},
						Destination: &bitbucket.PullRequestDestination{
							Branch: bitbucket.PullRequestBranch{
								Name: expectedPR.Destination.Branch.Name,
							},
						},
					}
					assert.Equal(t, expectedRequest, params.Request)
					assert.Empty(t, params.Request.Reviewers) // No reviewers in this test

					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:          repoName,
				Title:             expectedPR.Title,
				Description:       expectedPR.Description,
				SourceBranch:      expectedPR.Source.Branch.Name,
				DestBranch:        expectedPR.Destination.Branch.Name,
				CloseSourceBranch: true,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully creates pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, ok := deps.Client.(*MockBitbucketClient)
			require.True(t, ok, "Client should be a MockBitbucketClient")
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			account := NewRandomAtlassianAccount(
				WithAtlassianAccountName(accountName),
				WithAtlassianAccountDefault(false),
				WithAtlassianAccountBitbucket(
					WithBitbucketAccountWorkspace("custom-workspace-"+faker.Username()),
				),
			)

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Mock the accounts repo to return our test account
			mockAccounts.EXPECT().
				GetAccountByName(mock.Anything, accountName).
				Return(&account, nil)

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
				RepoName:     repoName,
				Title:        expectedPR.Title,
				Description:  expectedPR.Description,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully creates pull request with reviewers", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockClient, ok := deps.Client.(*MockBitbucketClient)
			require.True(t, ok, "Client should be a MockBitbucketClient")
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			account := NewRandomAtlassianAccount(
				WithAtlassianAccountName("default"),
				WithAtlassianAccountDefault(true),
				WithAtlassianAccountBitbucket(),
			)

			// Generate random reviewers
			reviewers := []string{
				"user-" + faker.Username(),
				"user-" + faker.Username(),
			}

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Mock the accounts repo
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(&account, nil)

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
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        expectedPR.Title,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
				Reviewers:    reviewers,
			})

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
		})

		t.Run("fails when account resolution fails", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Account not found
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(nil, ErrNoDefaultAccount)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        expectedPR.Title,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrNoDefaultAccount)
		})

		t.Run("fails when named account not found", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			accountName := "non-existent-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Account not found
			mockAccounts.EXPECT().
				GetAccountByName(mock.Anything, accountName).
				Return(nil, ErrAccountNotFound)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				AccountName:  accountName,
				RepoName:     repoName,
				Title:        expectedPR.Title,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrAccountNotFound)
		})

		t.Run("fails when account has no Bitbucket config", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps()
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			// Account with no Bitbucket config, only Jira
			account := NewRandomAtlassianAccount(
				WithAtlassianAccountName("default"),
				WithAtlassianAccountDefault(true),
				WithAtlassianAccountJira(),
			)
			// Remove Bitbucket config
			account.Bitbucket = nil

			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(&account, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        expectedPR.Title,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
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
						Title:        "PR-" + faker.Sentence(),
						SourceBranch: "feature/" + faker.Word(),
						DestBranch:   "main",
					},
					errMsg: "repository name is required",
				},
				{
					name: "missing title",
					params: BitbucketCreatePRParams{
						RepoName:     "repo-" + faker.Username(),
						SourceBranch: "feature/" + faker.Word(),
						DestBranch:   "main",
					},
					errMsg: "title is required",
				},
				{
					name: "missing source branch",
					params: BitbucketCreatePRParams{
						RepoName:   "repo-" + faker.Username(),
						Title:      "PR-" + faker.Sentence(),
						DestBranch: "main",
					},
					errMsg: "source branch is required",
				},
				{
					name: "missing destination branch",
					params: BitbucketCreatePRParams{
						RepoName:     "repo-" + faker.Username(),
						Title:        "PR-" + faker.Sentence(),
						SourceBranch: "feature/" + faker.Word(),
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
			mockClient, ok := deps.Client.(*MockBitbucketClient)
			require.True(t, ok, "Client should be a MockBitbucketClient")
			mockAccounts, ok := deps.AccountsRepo.(*MockAtlassianAccountsRepository)
			require.True(t, ok, "AccountsRepo should be a MockAtlassianAccountsRepository")
			service := NewBitbucketService(deps)

			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()

			account := NewRandomAtlassianAccount(
				WithAtlassianAccountName("default"),
				WithAtlassianAccountDefault(true),
				WithAtlassianAccountBitbucket(),
			)

			clientErr := errors.New("client error: " + faker.Sentence())

			// Mock account repo
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(&account, nil)

			// Mock client error
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, clientErr)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        expectedPR.Title,
				SourceBranch: expectedPR.Source.Branch.Name,
				DestBranch:   expectedPR.Destination.Branch.Name,
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

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(faker.RandomUnixTime())

			// Act
			result, err := service.ReadPR(t.Context(), BitbucketReadPRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
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

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(faker.RandomUnixTime())

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
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

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(faker.RandomUnixTime())

			// Act
			result, err := service.ApprovePR(t.Context(), BitbucketApprovePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
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

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(faker.RandomUnixTime())

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, "not implemented", err.Error())
		})
	})
}
