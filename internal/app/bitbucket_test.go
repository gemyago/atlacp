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
			Client:      NewMockbitbucketClient(t),
			AuthFactory: NewMockbitbucketAuthFactory(t),
			RootLogger:  diag.RootTestLogger(),
		}
	}

	t.Run("CreatePR", func(t *testing.T) {
		t.Run("successfully creates pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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

		t.Run("successfully creates draft pull request", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			expectedPR := bitbucket.NewRandomPullRequest()
			expectedPR.Draft = true // Set draft status on expected PR
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				CreatePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.CreatePRParams) bool {
					// Verify draft flag is properly set
					assert.True(t, params.Request.Draft)
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
				Draft:        true, // Set as draft PR
			})

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, result.Draft)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
		t.Run("successfully approves pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			expectedParticipant := bitbucket.NewRandomParticipant(true)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected participant
			mockClient.EXPECT().
				ApprovePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.ApprovePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)
					return true
				})).
				Return(expectedParticipant, nil)

			// Act
			result, err := service.ApprovePR(t.Context(), BitbucketApprovePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedParticipant, result)
		})

		t.Run("successfully approves pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			expectedParticipant := bitbucket.NewRandomParticipant(true)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider for the named account
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client to return expected participant
			mockClient.EXPECT().
				ApprovePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.ApprovePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)
					return true
				})).
				Return(expectedParticipant, nil)

			// Act
			result, err := service.ApprovePR(t.Context(), BitbucketApprovePRParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedParticipant, result)
		})

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketApprovePRParams
				errMsg string
			}{
				{
					name: "missing repo owner",
					params: BitbucketApprovePRParams{
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repo name",
					params: BitbucketApprovePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketApprovePRParams{
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
					result, err := service.ApprovePR(t.Context(), tc.params)

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
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
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
				ApprovePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.ApprovePR(t.Context(), BitbucketApprovePRParams{
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

	t.Run("MergePR", func(t *testing.T) {
		t.Run("successfully merges pull request with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			mergeStrategy := "merge_commit"
			message := "Merging PR: " + faker.Sentence()
			closeSourceBranch := true

			expectedPR := bitbucket.NewRandomPullRequest(
				bitbucket.WithPullRequestState("MERGED"),
				bitbucket.WithPullRequestCloseSourceBranch(closeSourceBranch),
			)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				MergePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.MergePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)

					// Verify merge parameters
					require.NotNil(t, params.MergeParameters)
					assert.Equal(t, message, params.MergeParameters.Message)
					assert.Equal(t, mergeStrategy, params.MergeParameters.MergeStrategy)
					assert.Equal(t, closeSourceBranch, params.MergeParameters.CloseSourceBranch)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				RepoOwner:         repoOwner,
				RepoName:          repoName,
				PullRequestID:     pullRequestID,
				Message:           message,
				CloseSourceBranch: closeSourceBranch,
				MergeStrategy:     mergeStrategy,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully merges pull request with named account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			mergeStrategy := "squash"
			expectedPR := bitbucket.NewRandomPullRequest(
				bitbucket.WithPullRequestState("MERGED"),
			)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider for the named account
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				MergePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.MergePRParams) bool {
					// Verify the parameters
					assert.Equal(t, repoOwner, params.Username)
					assert.Equal(t, repoName, params.RepoSlug)
					assert.Equal(t, pullRequestID, params.PullRequestID)

					// Verify merge strategy
					require.NotNil(t, params.MergeParameters)
					assert.Equal(t, mergeStrategy, params.MergeParameters.MergeStrategy)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				MergeStrategy: mergeStrategy,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("successfully uses fast-forward strategy", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			mergeStrategy := "fast_forward"
			expectedPR := bitbucket.NewRandomPullRequest(
				bitbucket.WithPullRequestState("MERGED"),
			)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected PR
			mockClient.EXPECT().
				MergePR(mock.Anything, mock.Anything, mock.MatchedBy(func(params bitbucket.MergePRParams) bool {
					// Verify merge strategy
					require.NotNil(t, params.MergeParameters)
					assert.Equal(t, mergeStrategy, params.MergeParameters.MergeStrategy)
					return true
				})).
				Return(expectedPR, nil)

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
				MergeStrategy: mergeStrategy,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedPR, result)
		})

		t.Run("fails when client returns error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			pullRequestID := int(faker.RandomUnixTime()) % 10000
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			clientErr := errors.New("client error: " + faker.Sentence())

			// Mock auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock client error
			mockClient.EXPECT().
				MergePR(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, clientErr)

			// Act
			result, err := service.MergePR(t.Context(), BitbucketMergePRParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: pullRequestID,
			})

			// Assert
			assert.Nil(t, result)
			require.Error(t, err)
			assert.ErrorIs(t, err, clientErr)
		})

		t.Run("fails when missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketMergePRParams
				errMsg string
			}{
				{
					name: "missing repo owner",
					params: BitbucketMergePRParams{
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repo name",
					params: BitbucketMergePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						PullRequestID: 1,
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketMergePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 0,
					},
					errMsg: "pull request ID must be positive",
				},
				{
					name: "invalid merge strategy",
					params: BitbucketMergePRParams{
						RepoOwner:     "owner-" + faker.Username(),
						RepoName:      "repo-" + faker.Username(),
						PullRequestID: 1,
						MergeStrategy: "invalid_strategy",
					},
					errMsg: "invalid merge strategy",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Act
					result, err := service.MergePR(t.Context(), tc.params)

					// Assert
					assert.Nil(t, result)
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})
	})

	t.Run("ListTasks", func(t *testing.T) {
		t.Run("successfully lists tasks", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			expectedTasks := createRandomPaginatedTasks(3)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected tasks
			mockClient.EXPECT().
				ListPullRequestTasks(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.ListPullRequestTasksParams) bool {
						// Verify the parameters
						assert.Equal(t, repoOwner, params.Workspace)
						assert.Equal(t, repoName, params.RepoSlug)
						assert.Equal(t, prID, params.PullReqID)

						// No query parameters in this test
						assert.Empty(t, params.Query)
						assert.Empty(t, params.Sort)
						assert.Zero(t, params.PageLen)

						return true
					})).
				Return(expectedTasks, nil)

			// Act
			result, err := service.ListTasks(t.Context(), BitbucketListTasksParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTasks, result)
		})

		t.Run("successfully lists tasks with query parameters", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			query := "state=\"UNRESOLVED\""
			sort := "created_on"
			pageLen := 10
			expectedTasks := createRandomPaginatedTasks(2)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				ListPullRequestTasks(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.ListPullRequestTasksParams) bool {
						// Verify query parameters
						assert.Equal(t, query, params.Query)
						assert.Equal(t, sort, params.Sort)
						assert.Equal(t, pageLen, params.PageLen)
						return true
					})).
				Return(expectedTasks, nil)

			// Act
			result, err := service.ListTasks(t.Context(), BitbucketListTasksParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				Query:         query,
				Sort:          sort,
				PageLen:       pageLen,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTasks, result)
		})

		t.Run("successfully lists tasks with custom account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			expectedTasks := createRandomPaginatedTasks(1)
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				ListPullRequestTasks(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(_ bitbucket.ListPullRequestTasksParams) bool {
						return true
					})).
				Return(expectedTasks, nil)

			// Act
			result, err := service.ListTasks(t.Context(), BitbucketListTasksParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTasks, result)
		})

		t.Run("fails when client returns error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			expectedErr := errors.New("client error")
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return error
			mockClient.EXPECT().
				ListPullRequestTasks(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.ListTasks(t.Context(), BitbucketListTasksParams{
				RepoOwner:     "owner",
				RepoName:      "repo",
				PullRequestID: 123,
			})

			// Assert
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, expectedErr)
		})

		t.Run("fails when required parameters are missing", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketListTasksParams
				errMsg string
			}{
				{
					name: "missing repository owner",
					params: BitbucketListTasksParams{
						RepoName:      "repo",
						PullRequestID: 123,
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repository name",
					params: BitbucketListTasksParams{
						RepoOwner:     "owner",
						PullRequestID: 123,
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketListTasksParams{
						RepoOwner:     "owner",
						RepoName:      "repo",
						PullRequestID: 0,
					},
					errMsg: "pull request ID must be positive",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Act
					result, err := service.ListTasks(t.Context(), tc.params)

					// Assert
					require.Error(t, err)
					assert.Nil(t, result)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})
	})

	t.Run("UpdateTask", func(t *testing.T) {
		t.Run("successfully updates task content", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			taskID := int(100 + faker.RandomUnixTime()%900)
			content := "Updated task content - " + faker.Sentence()

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      int64(taskID),
						State:   "UNRESOLVED",
						Content: &bitbucket.TaskContent{Raw: content},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected task
			mockClient.EXPECT().
				UpdateTask(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.UpdateTaskParams) bool {
						// Verify the parameters
						assert.Equal(t, repoOwner, params.Workspace)
						assert.Equal(t, repoName, params.RepoSlug)
						assert.Equal(t, prID, params.PullReqID)
						assert.Equal(t, taskID, params.TaskID)
						assert.Equal(t, content, params.Content)
						assert.Empty(t, params.State)

						return true
					})).
				Return(expectedTask, nil)

			// Act
			result, err := service.UpdateTask(t.Context(), BitbucketUpdateTaskParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				TaskID:        taskID,
				Content:       content,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("successfully updates task state", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			taskID := int(100 + faker.RandomUnixTime()%900)
			state := "RESOLVED"

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      int64(taskID),
						State:   state,
						Content: &bitbucket.TaskContent{Raw: "Original content"},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				UpdateTask(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.UpdateTaskParams) bool {
						// Verify state is updated
						assert.Equal(t, state, params.State)
						assert.Empty(t, params.Content)
						return true
					})).
				Return(expectedTask, nil)

			// Act
			result, err := service.UpdateTask(t.Context(), BitbucketUpdateTaskParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				TaskID:        taskID,
				State:         state,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("successfully updates task with custom account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			content := "Updated content"
			state := "RESOLVED"
			taskID := int(100 + faker.RandomUnixTime()%900)

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      int64(taskID),
						State:   state,
						Content: &bitbucket.TaskContent{Raw: content},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to use custom account
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				UpdateTask(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(_ bitbucket.UpdateTaskParams) bool {
						return true
					})).
				Return(expectedTask, nil)

			// Act
			result, err := service.UpdateTask(t.Context(), BitbucketUpdateTaskParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				TaskID:        taskID,
				Content:       content,
				State:         state,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("fails when client returns error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			expectedErr := errors.New("client error")
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return error
			mockClient.EXPECT().
				UpdateTask(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.UpdateTask(t.Context(), BitbucketUpdateTaskParams{
				RepoOwner:     "owner",
				RepoName:      "repo",
				PullRequestID: 123,
				TaskID:        456,
				Content:       "content",
			})

			// Assert
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, expectedErr)
		})

		t.Run("fails when required parameters are missing", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketUpdateTaskParams
				errMsg string
			}{
				{
					name: "missing repository owner",
					params: BitbucketUpdateTaskParams{
						RepoName:      "repo",
						PullRequestID: 123,
						TaskID:        456,
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repository name",
					params: BitbucketUpdateTaskParams{
						RepoOwner:     "owner",
						PullRequestID: 123,
						TaskID:        456,
					},
					errMsg: "repository name is required",
				},
				{
					name: "missing pull request ID",
					params: BitbucketUpdateTaskParams{
						RepoOwner: "owner",
						RepoName:  "repo",
						TaskID:    456,
					},
					errMsg: "pull request ID must be positive",
				},
				{
					name: "missing task ID",
					params: BitbucketUpdateTaskParams{
						RepoOwner:     "owner",
						RepoName:      "repo",
						PullRequestID: 123,
					},
					errMsg: "task ID must be positive",
				},
				{
					name: "no update content or state specified",
					params: BitbucketUpdateTaskParams{
						RepoOwner:     "owner",
						RepoName:      "repo",
						PullRequestID: 123,
						TaskID:        456,
					},
					errMsg: "either content or state must be provided",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					result, err := service.UpdateTask(t.Context(), tc.params)
					require.Error(t, err)
					assert.Nil(t, result)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})
	})

	t.Run("CreateTask", func(t *testing.T) {
		t.Run("successfully creates task with default account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			content := "Task content - " + faker.Sentence()

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      faker.RandomUnixTime() % 1000,
						State:   "UNRESOLVED",
						Content: &bitbucket.TaskContent{Raw: content},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory to return our token provider
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return expected task
			mockClient.EXPECT().
				CreatePullRequestTask(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.CreatePullRequestTaskParams) bool {
						// Verify the parameters
						assert.Equal(t, repoOwner, params.Workspace)
						assert.Equal(t, repoName, params.RepoSlug)
						assert.Equal(t, prID, params.PullReqID)
						assert.Equal(t, content, params.Content)
						assert.Zero(t, params.CommentID)
						assert.Nil(t, params.Pending)
						return true
					})).
				Return(expectedTask, nil)

			// Act
			result, err := service.CreateTask(t.Context(), BitbucketCreateTaskParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				Content:       content,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("successfully creates task with comment ID", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			content := "Task content - " + faker.Sentence()
			commentID := faker.RandomUnixTime() % 1000

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      faker.RandomUnixTime() % 1000,
						State:   "UNRESOLVED",
						Content: &bitbucket.TaskContent{Raw: content},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				CreatePullRequestTask(
					mock.Anything,
					mock.Anything,
					mock.MatchedBy(func(params bitbucket.CreatePullRequestTaskParams) bool {
						// Verify the parameters
						assert.Equal(t, repoOwner, params.Workspace)
						assert.Equal(t, repoName, params.RepoSlug)
						assert.Equal(t, prID, params.PullReqID)
						assert.Equal(t, content, params.Content)
						assert.Equal(t, commentID, params.CommentID)
						return true
					})).
				Return(expectedTask, nil)

			// Act
			result, err := service.CreateTask(t.Context(), BitbucketCreateTaskParams{
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				Content:       content,
				CommentID:     commentID,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("successfully creates task with custom account", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			// Create test data
			accountName := "custom-account-" + faker.Username()
			repoOwner := "owner-" + faker.Username()
			repoName := "repo-" + faker.Username()
			prID := int(100 + faker.RandomUnixTime()%900)
			content := "Task content - " + faker.Sentence()

			expectedTask := &bitbucket.PullRequestCommentTask{
				PullRequestTask: bitbucket.PullRequestTask{
					Task: bitbucket.Task{
						ID:      faker.RandomUnixTime() % 1000,
						State:   "UNRESOLVED",
						Content: &bitbucket.TaskContent{Raw: content},
					},
				},
			}

			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, accountName).
				Return(tokenProvider)

			// Mock the client
			mockClient.EXPECT().
				CreatePullRequestTask(mock.Anything, mock.Anything, mock.Anything).
				Return(expectedTask, nil)

			// Act
			result, err := service.CreateTask(t.Context(), BitbucketCreateTaskParams{
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				PullRequestID: prID,
				Content:       content,
			})

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedTask, result)
		})

		t.Run("fails when client returns error", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			mockClient := mocks.GetMock[*MockbitbucketClient](t, deps.Client)
			mockAuth := mocks.GetMock[*MockbitbucketAuthFactory](t, deps.AuthFactory)
			service := NewBitbucketService(deps)

			expectedErr := errors.New("client error")
			token := "token-" + faker.UUIDHyphenated()
			tokenProvider := newStaticTokenProvider(token)

			// Mock the auth factory
			mockAuth.EXPECT().
				getTokenProvider(mock.Anything, "").
				Return(tokenProvider)

			// Mock the client to return error
			mockClient.EXPECT().
				CreatePullRequestTask(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, expectedErr)

			// Act
			result, err := service.CreateTask(t.Context(), BitbucketCreateTaskParams{
				RepoOwner:     "owner",
				RepoName:      "repo",
				PullRequestID: 123,
				Content:       "Content",
			})

			// Assert
			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, expectedErr)
		})

		t.Run("fails when required parameters are missing", func(t *testing.T) {
			// Arrange
			deps := makeMockDeps(t)
			service := NewBitbucketService(deps)

			testCases := []struct {
				name   string
				params BitbucketCreateTaskParams
				errMsg string
			}{
				{
					name: "missing repository owner",
					params: BitbucketCreateTaskParams{
						RepoName:      "repo",
						PullRequestID: 123,
						Content:       "Content",
					},
					errMsg: "repository owner is required",
				},
				{
					name: "missing repository name",
					params: BitbucketCreateTaskParams{
						RepoOwner:     "owner",
						PullRequestID: 123,
						Content:       "Content",
					},
					errMsg: "repository name is required",
				},
				{
					name: "invalid pull request ID",
					params: BitbucketCreateTaskParams{
						RepoOwner:     "owner",
						RepoName:      "repo",
						PullRequestID: 0,
						Content:       "Content",
					},
					errMsg: "pull request ID must be positive",
				},
				{
					name: "missing content",
					params: BitbucketCreateTaskParams{
						RepoOwner:     "owner",
						RepoName:      "repo",
						PullRequestID: 123,
					},
					errMsg: "content is required",
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					result, err := service.CreateTask(t.Context(), tc.params)
					require.Error(t, err)
					assert.Nil(t, result)
					assert.Contains(t, err.Error(), tc.errMsg)
				})
			}
		})
	})
}

// Helper for creating random tasks.
func createRandomPaginatedTasks(count int) *bitbucket.PaginatedTasks {
	tasks := make([]bitbucket.PullRequestCommentTask, count)
	for i := range count {
		tasks[i] = bitbucket.PullRequestCommentTask{
			PullRequestTask: bitbucket.PullRequestTask{
				Task: bitbucket.Task{
					ID:      int64(i + 1),
					State:   "UNRESOLVED",
					Content: &bitbucket.TaskContent{Raw: "Task " + faker.Sentence()},
					Creator: &bitbucket.Account{DisplayName: faker.Name()},
				},
			},
		}
	}
	return &bitbucket.PaginatedTasks{
		Size:    count,
		Page:    1,
		PageLen: count,
		Values:  tasks,
	}
}
