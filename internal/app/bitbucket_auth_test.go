package app

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gemyago/atlacp/internal/diag"
)

func TestBitbucketAuthFactory(t *testing.T) {
	makeMockDeps := func(t *testing.T) (BitbucketAuthFactoryDeps, *MockAtlassianAccountsRepository) {
		mockRepo := NewMockAtlassianAccountsRepository(t)
		return BitbucketAuthFactoryDeps{
			AccountsRepo: mockRepo,
			RootLogger:   diag.RootTestLogger(),
		}, mockRepo
	}

	t.Run("should create new auth factory", func(t *testing.T) {
		deps, _ := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		assert.NotNil(t, auth)
		assert.NotNil(t, auth.logger)
	})

	t.Run("should get token provider for default account", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		expectedToken := faker.UUIDHyphenated()
		expectedAccount := &AtlassianAccount{
			Name:    faker.Name(),
			Default: true,
			Bitbucket: &BitbucketAccount{
				Token:     expectedToken,
				Workspace: faker.Username(),
			},
		}

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(expectedAccount, nil)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), "")

		require.NoError(t, err)
		assert.NotNil(t, tokenProvider)

		// Test that the token provider returns the correct token
		token, err := tokenProvider.GetToken(t.Context())
		require.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("should get token provider for named account", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		accountName := faker.Username()
		expectedToken := faker.UUIDHyphenated()
		expectedAccount := &AtlassianAccount{
			Name:    accountName,
			Default: false,
			Bitbucket: &BitbucketAccount{
				Token:     expectedToken,
				Workspace: faker.Username(),
			},
		}

		mockRepo.EXPECT().GetAccountByName(t.Context(), accountName).Return(expectedAccount, nil)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), accountName)

		require.NoError(t, err)
		assert.NotNil(t, tokenProvider)

		// Test that the token provider returns the correct token
		token, err := tokenProvider.GetToken(t.Context())
		require.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("should return error when default account not found", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(nil, ErrNoDefaultAccount)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), "")

		require.Error(t, err)
		assert.Nil(t, tokenProvider)
		assert.ErrorIs(t, err, ErrNoDefaultAccount)
	})

	t.Run("should return ErrAccountNotFound when default account returns ErrAccountNotFound", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(nil, ErrAccountNotFound)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), "")

		require.Error(t, err)
		assert.Nil(t, tokenProvider)
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("should return error when named account not found", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		accountName := faker.Username()
		mockRepo.EXPECT().GetAccountByName(t.Context(), accountName).Return(nil, ErrAccountNotFound)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), accountName)

		require.Error(t, err)
		assert.Nil(t, tokenProvider)
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("should return error when account has no bitbucket config", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		expectedAccount := &AtlassianAccount{
			Name:      faker.Name(),
			Default:   true,
			Bitbucket: nil, // No Bitbucket config
		}

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(expectedAccount, nil)

		tokenProvider, err := auth.GetTokenProvider(t.Context(), "")

		require.Error(t, err)
		assert.Nil(t, tokenProvider)
		assert.Contains(t, err.Error(), "bitbucket configuration not found")
	})
}
