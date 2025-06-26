package app

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for authentication providers.
func TestTokenProviders(t *testing.T) {
	t.Run("StaticTokenProvider returns correct token and type", func(t *testing.T) {
		// Arrange
		randomTokenValue := faker.UUIDHyphenated()
		// 'Bearer' is the required static type for tokens in our system
		const expectedTokenType = "Bearer"
		provider := newStaticTokenProvider(randomTokenValue)

		// Act
		token, err := provider.GetToken(t.Context())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedTokenType, token.Type)
		assert.Equal(t, randomTokenValue, token.Value)
	})

	t.Run("StaticTokenProvider returns empty value when initialized with empty string", func(t *testing.T) {
		// Arrange
		provider := newStaticTokenProvider("")

		// Act
		token, err := provider.GetToken(t.Context())

		// Assert
		require.NoError(t, err)
		// 'Bearer' is the required static type for tokens in our system
		assert.Equal(t, "Bearer", token.Type)
		assert.Empty(t, token.Value)
	})

	t.Run("TokenProviderFunc delegates to underlying provider", func(t *testing.T) {
		// Arrange
		randomTokenValue := faker.UUIDHyphenated()
		provider := newStaticTokenProvider(randomTokenValue)
		tokenFunc := tokenProviderFunc(provider.GetToken)

		// Act
		token, err := tokenFunc.GetToken(t.Context())

		// Assert
		require.NoError(t, err)
		// 'Bearer' is the required static type for tokens in our system
		assert.Equal(t, "Bearer", token.Type)
		assert.Equal(t, randomTokenValue, token.Value)
	})
}
