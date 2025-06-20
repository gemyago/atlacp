package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware(t *testing.T) {
	makeMockDeps := func() LoggingMiddlewareDeps {
		return LoggingMiddlewareDeps{
			RootLogger: diag.RootTestLogger(),
		}
	}

	t.Run("should call next transport and return response", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		mockTransport := &MockRoundTripper{}
		loggingMiddleware := NewLoggingMiddleware(mockTransport, deps)

		req := httptest.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", req).Return(expectedResponse, nil)

		// Act
		resp, err := loggingMiddleware.RoundTrip(req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should propagate errors from next transport", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		mockTransport := &MockRoundTripper{}
		loggingMiddleware := NewLoggingMiddleware(mockTransport, deps)

		req := httptest.NewRequest(http.MethodPost, "https://api.example.com/test", nil)
		expectedError := assert.AnError

		mockTransport.On("RoundTrip", req).Return((*http.Response)(nil), expectedError)

		// Act
		resp, err := loggingMiddleware.RoundTrip(req)

		// Assert
		assert.Nil(t, resp)
		assert.Equal(t, expectedError, err)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should not modify original request", func(t *testing.T) {
		// Arrange
		deps := makeMockDeps()
		mockTransport := &MockRoundTripper{}
		loggingMiddleware := NewLoggingMiddleware(mockTransport, deps)

		originalReq := httptest.NewRequest(http.MethodPut, "https://api.example.com/test", nil)
		originalReq.Header.Set("X-Original", "value")

		expectedResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		}

		mockTransport.On("RoundTrip", originalReq).Return(expectedResponse, nil)

		// Act
		_, err := loggingMiddleware.RoundTrip(originalReq)

		// Assert
		require.NoError(t, err)
		// Original request should be unchanged
		assert.Equal(t, "value", originalReq.Header.Get("X-Original"))
		mockTransport.AssertExpectations(t)
	})
}
