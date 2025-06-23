package app

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for authentication providers.
func TestTokenProviders(t *testing.T) {
	t.Run("StaticTokenProvider", func(t *testing.T) {
		t.Run("should return the provided token", func(t *testing.T) {
			// Arrange
			expectedToken := "token-" + faker.UUIDHyphenated()
			provider := newStaticTokenProvider(expectedToken)

			// Act
			token, err := provider.GetToken(t.Context())

			// Assert
			require.NoError(t, err)
			assert.Equal(t, expectedToken, token)
		})

		t.Run("should accept empty token", func(t *testing.T) {
			// Arrange
			provider := newStaticTokenProvider("")

			// Act
			token, err := provider.GetToken(t.Context())

			// Assert
			require.NoError(t, err)
			assert.Empty(t, token)
		})

		t.Run("should handle context cancellation gracefully", func(t *testing.T) {
			// Arrange
			provider := newStaticTokenProvider("test-token")
			ctx, cancel := context.WithCancel(t.Context())
			cancel() // Cancel the context

			// Act
			token, err := provider.GetToken(ctx)

			// Assert - should still work because the implementation doesn't use the context
			require.NoError(t, err)
			assert.Equal(t, "test-token", token)
		})
	})
}
