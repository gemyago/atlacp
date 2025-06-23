package app

import (
	"context"
	"errors"
)

// AtlassianAccountsRepository defines the port for accessing Atlassian account information.
// This is an outbound port that will be implemented by the infrastructure layer.
type AtlassianAccountsRepository interface {
	// GetDefaultAccount returns the default Atlassian account configuration.
	// Returns an error if no default account is found.
	GetDefaultAccount(ctx context.Context) (*AtlassianAccount, error)

	// GetAccountByName returns an account with the specified name.
	// Returns an error if no account with the name is found.
	GetAccountByName(ctx context.Context, name string) (*AtlassianAccount, error)
}

// TokenProvider provides authentication tokens for API requests.
type TokenProvider interface {
	// GetToken returns an authentication token for API requests.
	GetToken(ctx context.Context) (string, error)
}

// Error types for account-related operations.
var (
	// ErrNoDefaultAccount is returned when no default account is configured.
	ErrNoDefaultAccount = errors.New("no default Atlassian account configured")

	// ErrAccountNotFound is returned when a specific named account is not found.
	ErrAccountNotFound = errors.New("atlassian account not found")

	// ErrAccountConfigInvalid is returned when account configuration is invalid.
	ErrAccountConfigInvalid = errors.New("atlassian account configuration is invalid")
)
