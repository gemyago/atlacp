//go:build !release

package app

import (
	"github.com/go-faker/faker/v4"
)

// AtlassianAccountOpt is a function that configures an AtlassianAccount.
type AtlassianAccountOpt func(*AtlassianAccount)

// WithAtlassianAccountDefault sets the account as default.
func WithAtlassianAccountDefault(isDefault bool) AtlassianAccountOpt {
	return func(a *AtlassianAccount) {
		a.Default = isDefault
	}
}

// WithAtlassianAccountName sets the account name.
func WithAtlassianAccountName(name string) AtlassianAccountOpt {
	return func(a *AtlassianAccount) {
		a.Name = name
	}
}

// WithAtlassianAccountBitbucket adds a Bitbucket configuration to the account.
func WithAtlassianAccountBitbucket(opts ...AtlassianTokenOpt) AtlassianAccountOpt {
	return func(a *AtlassianAccount) {
		a.Bitbucket = NewRandomAtlassianToken(opts...)
	}
}

// WithAtlassianAccountJira adds a Jira configuration to the account.
func WithAtlassianAccountJira(opts ...AtlassianTokenOpt) AtlassianAccountOpt {
	return func(a *AtlassianAccount) {
		a.Jira = NewRandomAtlassianToken(opts...)
	}
}

// BitbucketAccountOpt is a function that configures a BitbucketAccount.
type AtlassianTokenOpt func(*AtlassianToken)

// WithBitbucketAccountToken sets the Bitbucket token.
func WithBitbucketAccountToken(token string) AtlassianTokenOpt {
	return func(b *AtlassianToken) {
		b.Value = token
	}
}

// WithBitbucketAccountTokenType sets the Bitbucket token type.
func WithBitbucketAccountTokenType(tokenType string) AtlassianTokenOpt {
	return func(b *AtlassianToken) {
		b.Type = tokenType
	}
}

// NewRandomAtlassianAccount generates a random AtlassianAccount for testing.
// Options can be used to customize the account.
func NewRandomAtlassianAccount(opts ...AtlassianAccountOpt) AtlassianAccount {
	account := AtlassianAccount{
		Name:      faker.Name(),
		Bitbucket: NewRandomAtlassianToken(),
		Jira:      NewRandomAtlassianToken(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(&account)
	}

	return account
}

// NewRandomBitbucketAccount generates a random BitbucketAccount for testing.
func NewRandomAtlassianToken(opts ...AtlassianTokenOpt) *AtlassianToken {
	account := &AtlassianToken{
		Type:  faker.Word(),
		Value: faker.UUIDHyphenated(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(account)
	}

	return account
}
