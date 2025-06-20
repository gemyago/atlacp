package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRoundTripper is a mock implementation of http.RoundTripper for testing.
type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	res, _ := args.Get(0).(*http.Response)
	return res, args.Error(1)
}

func TestAuthenticationMiddleware(t *testing.T) {
	t.Run("should add Bearer token when token is in context", func(t *testing.T) {
		// Arrange
		token := faker.Word()
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport)

		ctx := WithAuthToken(t.Context(), token)
		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
		req = req.WithContext(ctx)

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
			return r.Header.Get("Authorization") == "Bearer "+token
		})).Return(expectedResponse, nil)

		// Act
		resp, err := authMiddleware.RoundTrip(req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should pass through unchanged when no token in context", func(t *testing.T) {
		// Arrange
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport)

		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
			return r.Header.Get("Authorization") == ""
		})).Return(expectedResponse, nil)

		// Act
		resp, err := authMiddleware.RoundTrip(req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should not modify original request", func(t *testing.T) {
		// Arrange
		token := faker.Word()
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport)

		ctx := WithAuthToken(t.Context(), token)
		originalReq := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
		originalReq = originalReq.WithContext(ctx)

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(expectedResponse, nil)

		// Act
		_, err := authMiddleware.RoundTrip(originalReq)

		// Assert
		require.NoError(t, err)
		// Original request should not have Authorization header
		assert.Empty(t, originalReq.Header.Get("Authorization"))
		mockTransport.AssertExpectations(t)
	})
}
