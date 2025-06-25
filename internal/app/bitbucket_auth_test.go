package app

import (
	"errors"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

func TestBitbucketAuthFactory(t *testing.T) {
	makeMockDeps := func(t *testing.T) (BitbucketAuthFactoryDeps, *MockAtlassianAccountsRepository) {
		mockRepo := NewMockAtlassianAccountsRepository(t)
		return BitbucketAuthFactoryDeps{
			AccountsRepo: mockRepo,
			RootLogger:   diag.RootTestLogger(),
		}, mockRepo
	}

	t.Run("should get token provider for default account", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		expectedTokenValue := faker.UUIDHyphenated()
		expectedToken := middleware.Token{Type: "Bearer", Value: expectedTokenValue}
		expectedAccount := &AtlassianAccount{
			Name:    faker.Name(),
			Default: true,
			Bitbucket: &BitbucketAccount{
				Token:     expectedTokenValue,
				Workspace: faker.Username(),
			},
		}

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(expectedAccount, nil)

		tokenProvider := auth.getTokenProvider(t.Context(), "")
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
		expectedTokenValue := faker.UUIDHyphenated()
		expectedToken := middleware.Token{Type: "Bearer", Value: expectedTokenValue}
		expectedAccount := &AtlassianAccount{
			Name:    accountName,
			Default: false,
			Bitbucket: &BitbucketAccount{
				Token:     expectedTokenValue,
				Workspace: faker.Username(),
			},
		}

		mockRepo.EXPECT().GetAccountByName(t.Context(), accountName).Return(expectedAccount, nil)

		tokenProvider := auth.getTokenProvider(t.Context(), accountName)
		assert.NotNil(t, tokenProvider)

		// Test that the token provider returns the correct token
		token, err := tokenProvider.GetToken(t.Context())
		require.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("should return error when default account not found", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		wantErr := errors.New(faker.Sentence())
		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(nil, wantErr)

		tokenProvider := auth.getTokenProvider(t.Context(), "")

		token, err := tokenProvider.GetToken(t.Context())
		assert.Equal(t, middleware.Token{}, token)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return ErrAccountNotFound when default account returns ErrAccountNotFound", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		mockRepo.EXPECT().GetDefaultAccount(t.Context()).Return(nil, ErrAccountNotFound)

		tokenProvider := auth.getTokenProvider(t.Context(), "")

		token, err := tokenProvider.GetToken(t.Context())
		assert.Equal(t, middleware.Token{}, token)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("should return error when named account not found", func(t *testing.T) {
		deps, mockRepo := makeMockDeps(t)
		auth := newBitbucketAuthFactory(deps)

		accountName := faker.Username()
		mockRepo.EXPECT().GetAccountByName(t.Context(), accountName).Return(nil, ErrAccountNotFound)

		tokenProvider := auth.getTokenProvider(t.Context(), accountName)

		token, err := tokenProvider.GetToken(t.Context())
		assert.Equal(t, middleware.Token{}, token)
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

		tokenProvider := auth.getTokenProvider(t.Context(), "")
		token, err := tokenProvider.GetToken(t.Context())

		assert.Equal(t, middleware.Token{}, token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bitbucket configuration not found")
	})
}
