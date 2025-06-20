package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

// authTokenKey is the context key for storing authentication tokens.
type authTokenKey struct{}

// WithAuthToken adds an authentication token to the context.
func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey{}, token)
}

// AuthTokenFromContext extracts the authentication token from the context.
func AuthTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authTokenKey{}).(string)
	return token, ok
}

// AuthenticationMiddlewareDeps contains dependencies for the authentication middleware.
type AuthenticationMiddlewareDeps struct {
	RootLogger *slog.Logger
}

// AuthenticationMiddleware wraps an http.RoundTripper to add authentication headers.
type AuthenticationMiddleware struct {
	transport http.RoundTripper
	logger    *slog.Logger
}

// NewAuthenticationMiddleware creates a new authentication middleware
// Tokens are injected via context using WithAuthToken.
func NewAuthenticationMiddleware(transport http.RoundTripper, deps AuthenticationMiddlewareDeps) http.RoundTripper {
	return &AuthenticationMiddleware{
		transport: transport,
		logger:    deps.RootLogger.WithGroup("http-auth-middleware"),
	}
}

// RoundTrip implements http.RoundTripper interface
// Extracts token from context and adds Bearer Authorization header.
func (a *AuthenticationMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extract token from context
	token, hasToken := AuthTokenFromContext(req.Context())

	// If no token in context, log and pass request through unchanged
	if !hasToken {
		a.logger.DebugContext(req.Context(), "No authentication token found in context, passing request through unchanged")
		return a.transport.RoundTrip(req)
	}

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add Bearer authorization header
	clonedReq.Header.Set("Authorization", "Bearer "+token)

	// Pass the modified request to the next transport
	return a.transport.RoundTrip(clonedReq)
}
