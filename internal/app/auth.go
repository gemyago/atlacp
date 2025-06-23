package app

import (
	"context"
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
func (p *staticTokenProvider) GetToken(_ context.Context) (string, error) {
	return p.token, nil
}
