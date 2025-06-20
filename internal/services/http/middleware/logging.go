package middleware

import (
	"log/slog"
	"net/http"
)

// LoggingMiddlewareDeps contains dependencies for the logging middleware.
type LoggingMiddlewareDeps struct {
	RootLogger *slog.Logger
}

// LoggingMiddleware wraps an http.RoundTripper to add structured logging.
type LoggingMiddleware struct {
	transport http.RoundTripper
	logger    *slog.Logger
}

// This is a stub NOOP implementation for TDD.
func NewLoggingMiddleware(transport http.RoundTripper, deps LoggingMiddlewareDeps) http.RoundTripper {
	return &LoggingMiddleware{
		transport: transport,
		logger:    deps.RootLogger,
	}
}

// This is a stub NOOP implementation for TDD.
func (l *LoggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// NOOP: Just pass through to next transport without logging
	// This should make tests pass but logging will be added in actual implementation
	return l.transport.RoundTrip(req)
}
