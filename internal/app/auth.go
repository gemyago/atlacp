package app

import (
	"context"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// staticTokenProvider provides a static token for authentication.
type staticTokenProvider struct {
	token string
}

// newStaticTokenProvider creates a new provider that returns a static token.
func newStaticTokenProvider(token string) *staticTokenProvider {
	return &staticTokenProvider{token: token}
}

// GetToken returns the static token.
func (p *staticTokenProvider) GetToken(_ context.Context) (middleware.Token, error) {
	return middleware.Token{Type: "Bearer", Value: p.token}, nil
}

type tokenProviderFunc func(ctx context.Context) (middleware.Token, error)

func (f tokenProviderFunc) GetToken(ctx context.Context) (middleware.Token, error) {
	return f(ctx)
}
