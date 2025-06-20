package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientFactory(t *testing.T) {
	makeMockDeps := func() ClientFactoryDeps {
		return ClientFactoryDeps{
			RootLogger: diag.RootTestLogger(),
		}
	}

	t.Run("should create HTTP client with all middleware enabled", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		factory := NewClientFactory(deps)
		token := faker.Word()
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer "+token {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte(`{"error": "unauthorized"}`))
				assert.NoError(t, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"success": true, "all_middleware": true}`))
			assert.NoError(t, err)
		}))
		defer testServer.Close()

		// Act - all middleware enabled by default
		client := factory.CreateClient()

		// Create request with token in context
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		require.NoError(t, err)
		ctx := middleware.WithAuthToken(req.Context(), token)
		req = req.WithContext(ctx)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "all_middleware")

		// Should use default timeout
		assert.Equal(t, 30*time.Second, client.Timeout)
	})

	t.Run("should create HTTP client with all middleware disabled", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		factory := NewClientFactory(deps)
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Should receive request without auth header and return success
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"success": true, "no_middleware": true}`))
			assert.NoError(t, err)
		}))
		defer testServer.Close()

		// Act - disable all middleware
		client := factory.CreateClient(
			WithAuth(false),
			WithLogging(false),
			WithErrorHandling(false),
			WithTimeout(45*time.Second),
		)

		resp, err := client.Get(testServer.URL)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "no_middleware")

		// Should use custom timeout
		assert.Equal(t, 45*time.Second, client.Timeout)
	})
}
