package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/gemyago/atlacp/internal/app"
	"go.uber.org/dig"
)

// atlassianAccountsConfig represents the structure of the accounts configuration file.
type atlassianAccountsConfig struct {
	// List of Atlassian accounts
	Accounts []app.AtlassianAccount `json:"accounts"`
}

// atlassianAccountsRepository implements the app.AtlassianAccountsRepository interface.
type atlassianAccountsRepository struct {
	config *atlassianAccountsConfig
	logger *slog.Logger
}

// AtlassianAccountsRepositoryDeps contains dependencies for the accounts repository.
type AtlassianAccountsRepositoryDeps struct {
	dig.In

	RootLogger *slog.Logger
	ConfigPath string `name:"config.atlassian.accountsFilePath"`
}

// NewAtlassianAccountsRepository creates a new Atlassian accounts repository.
//
//nolint:ireturn // Hides repository implementation behind app port interface.
func NewAtlassianAccountsRepository(deps AtlassianAccountsRepositoryDeps) (app.AtlassianAccountsRepository, error) {
	logger := deps.RootLogger.WithGroup("atlassian-accounts")
	configPath := deps.ConfigPath

	// Use default path if not specified
	if configPath == "" {
		return nil, errors.New("accounts configuration path not specified")
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("accounts configuration file not found at %s", configPath)
	}

	// Read and parse configuration
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts configuration: %w", err)
	}

	var config atlassianAccountsConfig
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse accounts configuration: %w", unmarshalErr)
	}

	// Validate configuration
	if validateErr := app.ValidateAtlassianAccounts(config.Accounts); validateErr != nil {
		return nil, fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	return &atlassianAccountsRepository{
		config: &config,
		logger: logger,
	}, nil
}

// GetDefaultAccount returns the default Atlassian account configuration.
func (r *atlassianAccountsRepository) GetDefaultAccount(_ context.Context) (*app.AtlassianAccount, error) {
	for i, account := range r.config.Accounts {
		if account.Default {
			return &r.config.Accounts[i], nil
		}
	}
	return nil, app.ErrNoDefaultAccount
}

// GetAccountByName returns an account with the specified name.
func (r *atlassianAccountsRepository) GetAccountByName(_ context.Context, name string) (*app.AtlassianAccount, error) {
	for i, account := range r.config.Accounts {
		if account.Name == name {
			return &r.config.Accounts[i], nil
		}
	}
	return nil, fmt.Errorf("%w: %s", app.ErrAccountNotFound, name)
}
