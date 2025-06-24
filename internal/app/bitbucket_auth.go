package app

import (
	"context"
	"errors"
	"log/slog"

	"go.uber.org/dig"
)

// bitbucketAuthFactory is a factory for creating bitbucketAuth implementations.
type bitbucketAuthFactory interface {
	// getTokenProvider returns a TokenProvider for the specified account name.
	// If accountName is empty, uses the default account.
	getTokenProvider(ctx context.Context, accountName string) (TokenProvider, error)
}

// bitbucketAuthFactoryImpl provides authentication for Bitbucket operations by resolving
// account information and providing tokens for API requests.
type bitbucketAuthFactoryImpl struct {
	accountsRepo AtlassianAccountsRepository
	logger       *slog.Logger
}

// BitbucketAuthFactoryDeps contains dependencies for the Bitbucket account auth.
type BitbucketAuthFactoryDeps struct {
	dig.In

	AccountsRepo AtlassianAccountsRepository
	RootLogger   *slog.Logger
}

// newBitbucketAuthFactory creates a new Bitbucket account auth component.
func newBitbucketAuthFactory(deps BitbucketAuthFactoryDeps) *bitbucketAuthFactoryImpl {
	return &bitbucketAuthFactoryImpl{
		accountsRepo: deps.AccountsRepo,
		logger:       deps.RootLogger.WithGroup("app.bitbucket-account-auth"),
	}
}

// GetTokenProvider returns a TokenProvider for the specified account name.
// If accountName is empty, uses the default account.
func (a *bitbucketAuthFactoryImpl) GetTokenProvider(ctx context.Context, accountName string) (TokenProvider, error) {
	var account *AtlassianAccount
	var err error

	if accountName == "" {
		account, err = a.accountsRepo.GetDefaultAccount(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		account, err = a.accountsRepo.GetAccountByName(ctx, accountName)
		if err != nil {
			return nil, err
		}
	}

	// Validate account has Bitbucket configuration
	if account.Bitbucket == nil {
		return nil, errors.New("bitbucket configuration not found for account: " + account.Name)
	}

	// Create token provider using account token
	tokenProvider := newStaticTokenProvider(account.Bitbucket.Token)

	return tokenProvider, nil
}
