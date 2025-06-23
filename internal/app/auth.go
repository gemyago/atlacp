package app

import (
	"context"
)

// StaticTokenProvider provides a static token for authentication.
type StaticTokenProvider struct {
	token string
}

// NewStaticTokenProvider creates a new provider that returns a static token.
func NewStaticTokenProvider(token string) *StaticTokenProvider {
	return &StaticTokenProvider{token: token}
}

// GetToken returns the static token.
func (p *StaticTokenProvider) GetToken(_ context.Context) (string, error) {
	return p.token, nil
}

// BitbucketTokenProviderAdapter adapts TokenProvider to bitbucket.TokenProvider.
type BitbucketTokenProviderAdapter struct {
	provider TokenProvider
}

// NewBitbucketTokenProviderAdapter creates a new adapter.
func NewBitbucketTokenProviderAdapter(provider TokenProvider) *BitbucketTokenProviderAdapter {
	return &BitbucketTokenProviderAdapter{provider: provider}
}

// GetToken delegates to the underlying provider.
func (a *BitbucketTokenProviderAdapter) GetToken(ctx context.Context) (string, error) {
	return a.provider.GetToken(ctx)
}
