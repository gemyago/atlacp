package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gemyago/atlacp/internal/app"
	"go.uber.org/dig"
)

// atlassianAccountsRepository implements the app.AtlassianAccountsRepository interface.
type atlassianAccountsRepository struct {
	config *app.AtlassianAccountsConfig
	logger *slog.Logger
}

// AtlassianAccountsRepositoryDeps contains dependencies for the accounts repository.
type AtlassianAccountsRepositoryDeps struct {
	dig.In

	RootLogger *slog.Logger
	ConfigPath string `name:"config.atlassian.accountsFilePath" optional:"true"`
}

// NewAtlassianAccountsRepository creates a new Atlassian accounts repository.
func NewAtlassianAccountsRepository(deps AtlassianAccountsRepositoryDeps) (app.AtlassianAccountsRepository, error) {
	logger := deps.RootLogger.WithGroup("atlassian-accounts")
	configPath := deps.ConfigPath

	// Use default path if not specified
	if configPath == "" {
		configPath = getDefaultConfigPath()
		logger.Debug("Using default accounts configuration path", "path", configPath)
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

	var config app.AtlassianAccountsConfig
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse accounts configuration: %w", unmarshalErr)
	}

	// Validate configuration
	if validateErr := validateAccountsConfig(&config); validateErr != nil {
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

// validateAccountsConfig validates the accounts configuration.
func validateAccountsConfig(config *app.AtlassianAccountsConfig) error {
	if len(config.Accounts) == 0 {
		return errors.New("no accounts configured")
	}

	accountNames := make(map[string]bool)
	foundDefault := false

	for _, account := range config.Accounts {
		// Validate basic account properties
		if err := validateBasicAccountProperties(account, accountNames); err != nil {
			return err
		}
		accountNames[account.Name] = true

		// Validate service-specific configuration
		if err := validateServiceConfigs(account); err != nil {
			return err
		}

		// Track default account
		if account.Default {
			if foundDefault {
				return errors.New("multiple default accounts defined")
			}
			foundDefault = true
		}
	}

	// Ensure at least one default account exists
	if !foundDefault {
		return errors.New("no default account specified")
	}

	return nil
}

// validateBasicAccountProperties validates non-service-specific account properties.
func validateBasicAccountProperties(account app.AtlassianAccount, existingNames map[string]bool) error {
	// Check for duplicate names
	if existingNames[account.Name] {
		return fmt.Errorf("duplicate account name: %s", account.Name)
	}

	// Check that name is specified
	if account.Name == "" {
		return errors.New("account missing name")
	}

	// Ensure at least one service is configured
	if account.Bitbucket == nil && account.Jira == nil {
		return fmt.Errorf("account %s must have at least one service configured", account.Name)
	}

	return nil
}

// validateServiceConfigs validates Bitbucket and Jira configurations for an account.
func validateServiceConfigs(account app.AtlassianAccount) error {
	// Validate Bitbucket configuration if provided
	if account.Bitbucket != nil {
		if err := validateBitbucketConfig(account); err != nil {
			return err
		}
	}

	// Validate Jira configuration if provided
	if account.Jira != nil {
		if err := validateJiraConfig(account); err != nil {
			return err
		}
	}

	return nil
}

// validateBitbucketConfig validates Bitbucket-specific configuration.
func validateBitbucketConfig(account app.AtlassianAccount) error {
	if account.Bitbucket.Token == "" {
		return fmt.Errorf("account %s is missing Bitbucket token", account.Name)
	}
	if account.Bitbucket.Workspace == "" {
		return fmt.Errorf("account %s is missing Bitbucket workspace", account.Name)
	}
	return nil
}

// validateJiraConfig validates Jira-specific configuration.
func validateJiraConfig(account app.AtlassianAccount) error {
	if account.Jira.Token == "" {
		return fmt.Errorf("account %s is missing Jira token", account.Name)
	}
	if account.Jira.Domain == "" {
		return fmt.Errorf("account %s is missing Jira domain", account.Name)
	}
	return nil
}

// getDefaultConfigPath returns the default location for the accounts configuration file.
// It tries to find the configuration in common locations:
// 1. $HOME/.config/atlacp/accounts.json
// 2. Current directory: ./accounts.json.
func getDefaultConfigPath() string {
	// Try home directory first
	if homeDir, err := os.UserHomeDir(); err == nil {
		configDir := filepath.Join(homeDir, ".config", "atlacp")
		path := filepath.Join(configDir, "accounts.json")
		if _, statErr := os.Stat(path); statErr == nil {
			return path
		}
	}

	// Fall back to current directory
	return "accounts.json"
}
