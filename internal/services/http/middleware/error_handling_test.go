package middleware

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gemyago/atlacp/internal/diag"
)

func TestErrorHandlingMiddleware(t *testing.T) {
	makeMockDeps := func() ErrorHandlingMiddlewareDeps {
		return ErrorHandlingMiddlewareDeps{
			RootLogger: diag.RootTestLogger(),
		}
	}

	t.Run("should pass through successful responses unchanged", func(t *testing.T) {
		// Arrange
		expectedResp := &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
		}
		mockTransport := &MockRoundTripper{}
		mockTransport.On("RoundTrip", mock.Anything).Return(expectedResp, nil)

		middleware := NewErrorHandlingMiddleware(mockTransport, makeMockDeps())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/test"},
		}

		// Act
		resp, err := middleware.RoundTrip(req)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should wrap 4xx client errors in HTTPError", func(t *testing.T) {
		// Arrange
		clientErrorResp := &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
		}
		mockTransport := &MockRoundTripper{}
		mockTransport.On("RoundTrip", mock.Anything).Return(clientErrorResp, nil)

		middleware := NewErrorHandlingMiddleware(mockTransport, makeMockDeps())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/missing"},
		}

		// Act
		resp, err := middleware.RoundTrip(req)

		// Assert
		require.Error(t, err)
		assert.Nil(t, resp)

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 404, httpErr.StatusCode)
		assert.Equal(t, "GET", httpErr.Method)
		assert.Equal(t, "https://api.example.com/missing", httpErr.URL)
		assert.Contains(t, httpErr.Message, "client error")
		assert.Contains(t, httpErr.Message, "404")
		mockTransport.AssertExpectations(t)
	})

	t.Run("should wrap 5xx server errors in HTTPError", func(t *testing.T) {
		// Arrange
		serverErrorResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
		}
		mockTransport := &MockRoundTripper{}
		mockTransport.On("RoundTrip", mock.Anything).Return(serverErrorResp, nil)

		middleware := NewErrorHandlingMiddleware(mockTransport, makeMockDeps())
		req := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/action"},
		}

		// Act
		resp, err := middleware.RoundTrip(req)

		// Assert
		require.Error(t, err)
		assert.Nil(t, resp)

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 500, httpErr.StatusCode)
		assert.Equal(t, "POST", httpErr.Method)
		assert.Equal(t, "https://api.example.com/action", httpErr.URL)
		assert.Contains(t, httpErr.Message, "server error")
		assert.Contains(t, httpErr.Message, "500")
		mockTransport.AssertExpectations(t)
	})

	t.Run("should wrap transport errors and preserve error chain", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("network connection failed")
		mockTransport := &MockRoundTripper{}
		mockTransport.On("RoundTrip", mock.Anything).Return(nil, originalErr)

		middleware := NewErrorHandlingMiddleware(mockTransport, makeMockDeps())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/test"},
		}

		// Act
		resp, err := middleware.RoundTrip(req)

		// Assert
		require.Error(t, err)
		assert.Nil(t, resp)

		// Should preserve original error in chain
		require.ErrorIs(t, err, originalErr)

		// Should also be wrapped in HTTPError
		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 0, httpErr.StatusCode) // No status code for transport errors
		assert.Equal(t, "GET", httpErr.Method)
		assert.Equal(t, "https://api.example.com/test", httpErr.URL)
		assert.Contains(t, httpErr.Message, "transport error")
		assert.Equal(t, originalErr, httpErr.Err)
		mockTransport.AssertExpectations(t)
	})

	t.Run("should not modify original request", func(t *testing.T) {
		// Arrange
		mockTransport := &MockRoundTripper{}
		mockTransport.On("RoundTrip", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

		middleware := NewErrorHandlingMiddleware(mockTransport, makeMockDeps())
		originalReq := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "api.example.com", Path: "/test"},
			Header: http.Header{"X-Test": []string{"original"}},
		}

		// Act
		_, _ = middleware.RoundTrip(originalReq)

		// Assert
		assert.Equal(t, "GET", originalReq.Method)
		assert.Equal(t, "https://api.example.com/test", originalReq.URL.String())
		assert.Equal(t, "original", originalReq.Header.Get("X-Test"))
		mockTransport.AssertExpectations(t)
	})
}
