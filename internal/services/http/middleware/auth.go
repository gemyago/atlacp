package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

// authTokenKey is the context key for storing authentication tokens.
type authTokenKey struct{}

// Token represents an authentication token with a type and value.
type Token struct {
	Type  string
	Value string
}

// WithAuthTokenV2 adds a Token struct to the context.
func WithAuthTokenV2(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, authTokenKey{}, token)
}

// WithAuthToken adds a Bearer authentication token to the context (legacy API).
func WithAuthToken(ctx context.Context, token string) context.Context {
	return WithAuthTokenV2(ctx, Token{Type: "Bearer", Value: token})
}

// AuthTokenFromContext extracts the Token from the context.
func AuthTokenFromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(authTokenKey{}).(Token)
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
// Extracts token from context and adds Authorization header.
func (a *AuthenticationMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extract token from context
	token, hasToken := AuthTokenFromContext(req.Context())

	// If no token in context, log and pass request through unchanged
	if !hasToken || token.Value == "" {
		a.logger.DebugContext(req.Context(), "No authentication token found in context, passing request through unchanged")
		return a.transport.RoundTrip(req)
	}

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add Authorization header
	clonedReq.Header.Set("Authorization", token.Type+" "+token.Value)

	// Pass the modified request to the next transport
	return a.transport.RoundTrip(clonedReq)
}
