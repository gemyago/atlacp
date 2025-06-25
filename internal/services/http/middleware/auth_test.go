package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
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

func makeMockDeps() AuthenticationMiddlewareDeps {
	return AuthenticationMiddlewareDeps{
		RootLogger: diag.RootTestLogger(),
	}
}

func TestAuthenticationMiddleware(t *testing.T) {
	t.Run("should add Bearer token when token is in context", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		tokenValue := faker.Word()
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport, deps)

		ctx := WithAuthToken(t.Context(), tokenValue)
		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
		req = req.WithContext(ctx)

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
			return r.Header.Get("Authorization") == "Bearer "+tokenValue
		})).Return(expectedResponse, nil)

		// Act
		resp, err := authMiddleware.RoundTrip(req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should add custom token type when using WithAuthTokenV2", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		tokenValue := faker.Word()
		tokenType := "CustomType"
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport, deps)

		token := Token{Type: tokenType, Value: tokenValue}
		ctx := WithAuthTokenV2(t.Context(), token)
		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
		req = req.WithContext(ctx)

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
			return r.Header.Get("Authorization") == tokenType+" "+tokenValue
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
		deps := makeMockDeps()
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport, deps)

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
		deps := makeMockDeps()
		tokenValue := faker.Word()
		mockTransport := &MockRoundTripper{}
		authMiddleware := NewAuthenticationMiddleware(mockTransport, deps)

		ctx := WithAuthToken(t.Context(), tokenValue)
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

	t.Run("should extract token from context using AuthTokenFromContext", func(t *testing.T) {
		// Arrange
		token := Token{Type: "Bearer", Value: faker.Word()}
		ctx := WithAuthTokenV2(t.Context(), token)

		// Act
		extracted, ok := AuthTokenFromContext(ctx)

		// Assert
		require.True(t, ok)
		assert.Equal(t, token, extracted)
	})
}
