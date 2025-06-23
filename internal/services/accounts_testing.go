//go:build !release

package services

import (
	"github.com/gemyago/atlacp/internal/app"
	"github.com/go-faker/faker/v4"
)

// AtlassianAccountOpt is a function that configures an AtlassianAccount.
type AtlassianAccountOpt func(*app.AtlassianAccount)

// WithAtlassianAccountDefault sets the account as default.
func WithAtlassianAccountDefault(isDefault bool) AtlassianAccountOpt {
	return func(a *app.AtlassianAccount) {
		a.Default = isDefault
	}
}

// WithAtlassianAccountName sets the account name.
func WithAtlassianAccountName(name string) AtlassianAccountOpt {
	return func(a *app.AtlassianAccount) {
		a.Name = name
	}
}

// WithAtlassianAccountBitbucket adds a Bitbucket configuration to the account.
func WithAtlassianAccountBitbucket(opts ...BitbucketAccountOpt) AtlassianAccountOpt {
	return func(a *app.AtlassianAccount) {
		a.Bitbucket = NewRandomBitbucketAccount(opts...)
	}
}

// WithAtlassianAccountJira adds a Jira configuration to the account.
func WithAtlassianAccountJira(opts ...JiraAccountOpt) AtlassianAccountOpt {
	return func(a *app.AtlassianAccount) {
		a.Jira = NewRandomJiraAccount(opts...)
	}
}

// BitbucketAccountOpt is a function that configures a BitbucketAccount.
type BitbucketAccountOpt func(*app.BitbucketAccount)

// WithBitbucketAccountToken sets the Bitbucket token.
func WithBitbucketAccountToken(token string) BitbucketAccountOpt {
	return func(b *app.BitbucketAccount) {
		b.Token = token
	}
}

// WithBitbucketAccountWorkspace sets the Bitbucket workspace.
func WithBitbucketAccountWorkspace(workspace string) BitbucketAccountOpt {
	return func(b *app.BitbucketAccount) {
		b.Workspace = workspace
	}
}

// JiraAccountOpt is a function that configures a JiraAccount.
type JiraAccountOpt func(*app.JiraAccount)

// WithJiraAccountToken sets the Jira token.
func WithJiraAccountToken(token string) JiraAccountOpt {
	return func(j *app.JiraAccount) {
		j.Token = token
	}
}

// WithJiraAccountDomain sets the Jira domain.
func WithJiraAccountDomain(domain string) JiraAccountOpt {
	return func(j *app.JiraAccount) {
		j.Domain = domain
	}
}

// NewRandomAtlassianAccount generates a random AtlassianAccount for testing.
// Options can be used to customize the account.
func NewRandomAtlassianAccount(opts ...AtlassianAccountOpt) app.AtlassianAccount {
	account := app.AtlassianAccount{
		Name:      faker.Name(),
		Bitbucket: NewRandomBitbucketAccount(),
		Jira:      NewRandomJiraAccount(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(&account)
	}

	return account
}

// NewRandomBitbucketAccount generates a random BitbucketAccount for testing.
func NewRandomBitbucketAccount(opts ...BitbucketAccountOpt) *app.BitbucketAccount {
	account := &app.BitbucketAccount{
		Token:     faker.UUIDHyphenated(),
		Workspace: faker.Username(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(account)
	}

	return account
}

// NewRandomJiraAccount generates a random JiraAccount for testing.
func NewRandomJiraAccount(opts ...JiraAccountOpt) *app.JiraAccount {
	account := &app.JiraAccount{
		Token:  faker.UUIDHyphenated(),
		Domain: faker.DomainName(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(account)
	}

	return account
}
