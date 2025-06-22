package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	// HTTP status code boundaries for error classification.
	httpStatusClientErrorMin = 400
	httpStatusServerErrorMin = 500
)

// HTTPError represents an HTTP-related error with context.
type HTTPError struct {
	StatusCode int
	Method     string
	URL        string
	Message    string
	Err        error
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements error unwrapping for error chain support.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// ErrorHandlingMiddlewareDeps contains dependencies for the error handling middleware.
type ErrorHandlingMiddlewareDeps struct {
	RootLogger *slog.Logger
}

// ErrorHandlingMiddleware wraps an http.RoundTripper to add generic HTTP error handling.
type ErrorHandlingMiddleware struct {
	transport http.RoundTripper
	logger    *slog.Logger
}

// NewErrorHandlingMiddleware creates a new error handling middleware.
func NewErrorHandlingMiddleware(transport http.RoundTripper, deps ErrorHandlingMiddlewareDeps) http.RoundTripper {
	return &ErrorHandlingMiddleware{
		transport: transport,
		logger:    deps.RootLogger.WithGroup("http-error-middleware"),
	}
}

// RoundTrip implements http.RoundTripper interface.
// Handles HTTP errors by wrapping non-2xx responses and transport errors in HTTPError.
func (e *ErrorHandlingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// Call next transport
	resp, err := e.transport.RoundTrip(req)

	// Handle transport errors (network issues, etc.)
	if err != nil {
		httpErr := &HTTPError{
			StatusCode: 0, // No status code for transport errors
			Method:     req.Method,
			URL:        req.URL.String(),
			Message:    "HTTP transport error",
			Err:        err,
		}
		e.logger.ErrorContext(req.Context(), "HTTP transport error",
			"method", req.Method,
			"url", req.URL.String(),
			"error", err,
		)
		return nil, httpErr
	}

	// Handle HTTP error status codes
	if resp.StatusCode >= httpStatusClientErrorMin {
		var message string
		if resp.StatusCode >= httpStatusServerErrorMin {
			message = fmt.Sprintf("HTTP server error (%d %s)", resp.StatusCode, resp.Status)
		} else {
			message = fmt.Sprintf("HTTP client error (%d %s)", resp.StatusCode, resp.Status)
		}

		httpErr := &HTTPError{
			StatusCode: resp.StatusCode,
			Method:     req.Method,
			URL:        req.URL.String(),
			Message:    message,
			Err:        nil, // No underlying error for HTTP status errors
		}

		// Read response body for logging, but preserve it for downstream code
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = io.ReadAll(resp.Body)
			// Close the original body to prevent resource leaks
			resp.Body.Close()
			// Replace with a new reader containing the same data
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// Log the error with body content
		e.logger.WarnContext(req.Context(), "HTTP error response",
			"method", req.Method,
			"url", req.URL.String(),
			"status_code", resp.StatusCode,
			"status", resp.Status,
			"response_body", string(bodyBytes),
		)

		// Return both the error and the response so downstream code can inspect the body
		return resp, httpErr
	}

	// Success case - pass through unchanged
	return resp, nil
}
