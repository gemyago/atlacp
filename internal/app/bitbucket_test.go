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
			prTitle := "PR-" + faker.Sentence()
			prDesc := faker.Paragraph()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"
			prID := int(faker.RandomUnixTime())

			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
				Title:       prTitle,
				Description: prDesc,
				Source: bitbucket.PullRequestSource{
					Branch: bitbucket.PullRequestBranch{
						Name: sourceBranch,
					},
				},
				Destination: &bitbucket.PullRequestDestination{
					Branch: bitbucket.PullRequestBranch{
						Name: destBranch,
					},
				},
				CloseSourceBranch: true,
			}

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

					// Verify the PR request
					pr := params.Request
					assert.Equal(t, prTitle, pr.Title)
					assert.Equal(t, prDesc, pr.Description)
					assert.Equal(t, sourceBranch, pr.Source.Branch.Name)
					assert.Equal(t, destBranch, pr.Destination.Branch.Name)
					assert.True(t, pr.CloseSourceBranch)
					assert.Empty(t, pr.Reviewers) // No reviewers in this test

					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:          repoName,
				Title:             prTitle,
				Description:       prDesc,
				SourceBranch:      sourceBranch,
				DestBranch:        destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			prDesc := faker.Paragraph()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"
			prID := int(faker.RandomUnixTime())

			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
				Title:       prTitle,
				Description: prDesc,
			}

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
				Title:        prTitle,
				Description:  prDesc,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"
			prID := int(faker.RandomUnixTime())

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
				Return(&bitbucket.PullRequest{ID: prID}, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        prTitle,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"

			// Account not found
			mockAccounts.EXPECT().
				GetDefaultAccount(mock.Anything).
				Return(nil, ErrNoDefaultAccount)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoName:     repoName,
				Title:        prTitle,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"

			// Account not found
			mockAccounts.EXPECT().
				GetAccountByName(mock.Anything, accountName).
				Return(nil, ErrAccountNotFound)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				AccountName:  accountName,
				RepoName:     repoName,
				Title:        prTitle,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"

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
				Title:        prTitle,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
			prTitle := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Word()
			destBranch := "main"

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
				Title:        prTitle,
				SourceBranch: sourceBranch,
				DestBranch:   destBranch,
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
