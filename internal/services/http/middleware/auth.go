package middleware

import (
	"context"
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

// AuthenticationMiddleware wraps an http.RoundTripper to add authentication headers.
type AuthenticationMiddleware struct {
	transport http.RoundTripper
}

// Tokens are injected via context using WithAuthToken.
func NewAuthenticationMiddleware(transport http.RoundTripper) http.RoundTripper {
	return &AuthenticationMiddleware{
		transport: transport,
	}
}

// Extracts token from context and adds Bearer Authorization header.
func (a *AuthenticationMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extract token from context
	token, hasToken := AuthTokenFromContext(req.Context())

	// If no token in context, pass request through unchanged
	if !hasToken {
		return a.transport.RoundTrip(req)
	}

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add Bearer authorization header
	clonedReq.Header.Set("Authorization", "Bearer "+token)

	// Pass the modified request to the next transport
	return a.transport.RoundTrip(clonedReq)
}
