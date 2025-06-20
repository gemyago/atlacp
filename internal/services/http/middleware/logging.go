package middleware

import (
	"log/slog"
	"net/http"
	"time"
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

// NewLoggingMiddleware creates a new logging middleware.
func NewLoggingMiddleware(transport http.RoundTripper, deps LoggingMiddlewareDeps) http.RoundTripper {
	return &LoggingMiddleware{
		transport: transport,
		logger:    deps.RootLogger.WithGroup("http-logging-middleware"),
	}
}

// RoundTrip implements http.RoundTripper interface.
// Logs request and response details with structured logging.
func (l *LoggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log request
	l.logger.DebugContext(req.Context(), "HTTP request started",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.String("host", req.Host),
	)

	// Call next transport
	resp, err := l.transport.RoundTrip(req)
	duration := time.Since(start)

	// Log response
	if err != nil {
		l.logger.ErrorContext(req.Context(), "HTTP request failed",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Duration("duration", duration),
			slog.Any("error", err),
		)
		return nil, err
	}

	l.logger.DebugContext(req.Context(), "HTTP request completed",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("status_code", resp.StatusCode),
		slog.Duration("duration", duration),
	)

	return resp, nil
}
