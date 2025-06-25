package app

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// Tests for authentication providers.
func TestTokenProviders(t *testing.T) {
	t.Run("StaticTokenProvider", func(t *testing.T) {
		t.Run("should return the provided token", func(t *testing.T) {
			// Arrange
			expectedTokenValue := "token-" + faker.UUIDHyphenated()
			expectedToken := middleware.Token{Type: "Bearer", Value: expectedTokenValue}
			provider := newStaticTokenProvider(expectedTokenValue)

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
			assert.Equal(t, middleware.Token{Type: "Bearer", Value: ""}, token)
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
			assert.Equal(t, middleware.Token{Type: "Bearer", Value: "test-token"}, token)
		})
	})
}
