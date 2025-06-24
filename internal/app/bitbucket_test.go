package app

import (
	"errors"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/gemyago/atlacp/internal/testing/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Tests for BitbucketService.
func TestBitbucketService(t *testing.T) {
	// Helper function to create mock dependencies
	makeMockDeps := func(t *testing.T) BitbucketServiceDeps {
		return BitbucketServiceDeps{
			Client:      NewMockBitbucketClient(t),
			AuthFactory: NewMockbitbucketAuthFactory(t),
			RootLogger:  diag.RootTestLogger(),
		}
	}

	t.Run("CreatePR", func(t *testing.T) {
		t.Run("successfully creates pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
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
				RepoOwner:         repoOwner,
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
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify the workspace from the custom account is used
					assert.Equal(t, repoOwner, params.Username)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				AccountName:  accountName,
				RepoOwner:    repoOwner,
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
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Generate random reviewers
			reviewers := []string{
				"user-" + faker.Username(),
				"user-" + faker.Username(),
			}

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

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
				RepoOwner:    repoOwner,
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

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketCreatePRParams
				errMsg string
			}{
				{
					name: "missing repo name",
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
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			clientErr := errors.New("client error: " + faker.Sentence())

			// Mock auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock client error
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, clientErr)

			// Act
			result, err := service.CreatePR(t.Context(), BitbucketCreatePRParams{
				RepoOwner:    repoOwner,
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
		t.Run("successfully retrieves pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				GetPR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.GetPRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.ReadPR(t.Context(), BitbucketReadPRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully retrieves pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			expectedPR := bitbucket.NewRandomPullRequest()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				GetPR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.GetPRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.ReadPR(t.Context(), BitbucketReadPRParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketReadPRParams
				errMsg string
			}{
				{
					name: "missing repo owner",
					params: BitbucketReadPRParams{
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repo name",
					params: BitbucketReadPRParams{
						RepoOwner:     "owner-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketReadPRParams{
						RepoOwner:     "owner-" + faker.Username(),
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 0,
					},
					errMsg: "pull request ID must be positive",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Act
					result, err := service.ReadPR(t.Context(), tc.params)

					// Assert
					assert.Nil(t, result)
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})

		t.Run("handles client error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)
			expectedErr := errors.New("API error: " + faker.Sentence())

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return an error
			mockClient.EXPECT().
				GetPR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.ReadPR(t.Context(), BitbucketReadPRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, expectedErr, errors.Unwrap(err))
		})
	})

	t.Run("UpdatePR", func(t *testing.T) {
		t.Run("successfully updates pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000

			newTitle := "Updated: " + faker.Sentence()
			newDescription := "Updated description: " + faker.Paragraph()

			expectedPR := bitbucket.NewRandomPullRequest()
			expectedPR.Title = newTitle
			expectedPR.Description = newDescription

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				UpdatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.UpdatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)

					// Verify the update request
					assert.Equal(t, newTitle, params.Request.Title)
					assert.Equal(t, newDescription, params.Request.Description)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				Title:         newTitle,
				Description:   newDescription,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully updates pull request with title only", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000

			newTitle := "Title only update: " + faker.Sentence()

			expectedPR := bitbucket.NewRandomPullRequest()
			expectedPR.Title = newTitle

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				UpdatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.UpdatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)

					// Verify the update request has title only
					assert.Equal(t, newTitle, params.Request.Title)
					assert.Empty(t, params.Request.Description)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				Title:         newTitle,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully updates pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000

			newDescription := "Description only update: " + faker.Paragraph()

			expectedPR := bitbucket.NewRandomPullRequest()
			expectedPR.Description = newDescription

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider for the named account
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				UpdatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.UpdatePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)

					// Verify update request has description only
					assert.Empty(t, params.Request.Title)
					assert.Equal(t, newDescription, params.Request.Description)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				Description:   newDescription,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketUpdatePRParams
				errMsg string
			}{
				{
					name: "missing repo owner",
					params: BitbucketUpdatePRParams{
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
						Title:         "New title",
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repo name",
					params: BitbucketUpdatePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						PullRequestID: 1,
						Title:         "New title",
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketUpdatePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 0,
						Title:         "New title",
					},
					errMsg: "pull request ID must be positive",
				},
				{
					name: "missing both title and description",
					params: BitbucketUpdatePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "either title or description must be provided",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Act
					result, err := service.UpdatePR(t.Context(), tc.params)

					// Assert
					assert.Nil(t, result)
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})

		t.Run("handles client error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockBitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			newTitle := "Title: " + faker.Sentence()
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)
			expectedErr := errors.New("API error: " + faker.Sentence())

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return an error
			mockClient.EXPECT().
				UpdatePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.UpdatePR(t.Context(), BitbucketUpdatePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				Title:         newTitle,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.Equal(t, expectedErr, errors.Unwrap(err))
		})
	})

	t.Run("ApprovePR", func(t *testing.T) {
		t.Run("returns not implemented error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
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
			deps := makeMockDeps(t)
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
