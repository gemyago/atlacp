package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
	"go.uber.org/dig"
)

// bitbucketAuthFactory is a factory for creating bitbucketAuth implementations.
type bitbucketAuthFactory interface {
	// getTokenProvider returns a TokenProvider for the specified account name.
	// If accountName is empty, uses the default account.
	getTokenProvider(ctx context.Context, accountName string) TokenProvider
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
func newBitbucketAuthFactory(deps BitbucketAuthFactoryDeps) bitbucketAuthFactory {
	return &bitbucketAuthFactoryImpl{
		accountsRepo: deps.AccountsRepo,
		logger:       deps.RootLogger.WithGroup("app.bitbucket-account-auth"),
	}
}

// getTokenProvider returns a TokenProvider for the specified account name.
// If accountName is empty, uses the default account.
func (a *bitbucketAuthFactoryImpl) getTokenProvider(_ context.Context, accountName string) TokenProvider {
	return tokenProviderFunc(func(ctx context.Context) (middleware.Token, error) {
		var account *AtlassianAccount
		var err error

		if accountName == "" {
			account, err = a.accountsRepo.GetDefaultAccount(ctx)
			if err != nil {
				return middleware.Token{}, err
			}
		} else {
			account, err = a.accountsRepo.GetAccountByName(ctx, accountName)
			if err != nil {
				return middleware.Token{}, err
			}
		}

		// Validate account has Bitbucket configuration
		if account.Bitbucket == nil {
			return middleware.Token{}, errors.New("bitbucket configuration not found for account: " + account.Name)
		}

		return middleware.Token{
			Type:  "Bearer",
			Value: account.Bitbucket.Token,
		}, nil
	})
}
