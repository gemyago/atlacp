package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gemyago/atlacp/internal/app"
	"go.uber.org/dig"
)

// atlassianAccountsConfig represents the structure of the accounts configuration file.
type atlassianAccountsConfig struct {
	// List of Atlassian accounts
	Accounts []atlassianAccountConfig `json:"accounts"`
}

// atlassianAccountConfig represents configuration for a single Atlassian account.
type atlassianAccountConfig struct {
	// Friendly name of the account
	Name string `json:"name"`

	// Is this the default account
	Default bool `json:"default"`

	// Bitbucket-specific configuration
	Bitbucket *bitbucketAccountConfig `json:"bitbucket,omitempty"`

	// Jira-specific configuration
	Jira *jiraAccountConfig `json:"jira,omitempty"`
}

// bitbucketAccountConfig contains Bitbucket-specific account configuration.
type bitbucketAccountConfig struct {
	// API token for authentication
	Token string `json:"token"`

	// Workspace is the Bitbucket workspace/username for this account
	Workspace string `json:"workspace"`
}

// jiraAccountConfig contains Jira-specific account configuration.
type jiraAccountConfig struct {
	// API token for authentication
	Token string `json:"token"`

	// Domain is the Jira cloud instance domain (e.g., "mycompany" for mycompany.atlassian.net)
	Domain string `json:"domain"`
}

// atlassianAccountsRepository implements the app.AtlassianAccountsRepository interface.
type atlassianAccountsRepository struct {
	// Will be implemented later
}

// AtlassianAccountsRepositoryDeps contains dependencies for the accounts repository.
type AtlassianAccountsRepositoryDeps struct {
	dig.In

	RootLogger *slog.Logger
	ConfigPath string `name:"config.atlassian.accountsFilePath" optional:"true"`
}

// NewAtlassianAccountsRepository creates a new Atlassian accounts repository.
func NewAtlassianAccountsRepository() app.AtlassianAccountsRepository {
	return &atlassianAccountsRepository{}
}

// GetDefaultAccount returns the default Atlassian account configuration.
func (r *atlassianAccountsRepository) GetDefaultAccount(ctx context.Context) (*app.AtlassianAccount, error) {
	return nil, errors.New("not implemented")
}

// GetAccountByName returns an account with the specified name.
func (r *atlassianAccountsRepository) GetAccountByName(ctx context.Context, name string) (*app.AtlassianAccount, error) {
	return nil, errors.New("not implemented")
}

// validateAccountsConfig validates the accounts configuration.
func validateAccountsConfig(config *atlassianAccountsConfig) error {
	if len(config.Accounts) == 0 {
		return errors.New("no accounts configured")
	}

	foundDefault := false
	accountNames := make(map[string]bool)

	for _, account := range config.Accounts {
		// Check for duplicate names
		if accountNames[account.Name] {
			return fmt.Errorf("duplicate account name: %s", account.Name)
		}
		accountNames[account.Name] = true

		// Check that name is specified
		if account.Name == "" {
			return errors.New("account missing name")
		}

		// Ensure at least one service is configured
		if account.Bitbucket == nil && account.Jira == nil {
			return fmt.Errorf("account %s must have at least one service configured", account.Name)
		}

		// Validate Bitbucket configuration if provided
		if account.Bitbucket != nil {
			if account.Bitbucket.Token == "" {
				return fmt.Errorf("account %s is missing Bitbucket token", account.Name)
			}
			if account.Bitbucket.Workspace == "" {
				return fmt.Errorf("account %s is missing Bitbucket workspace", account.Name)
			}
		}

		// Validate Jira configuration if provided
		if account.Jira != nil {
			if account.Jira.Token == "" {
				return fmt.Errorf("account %s is missing Jira token", account.Name)
			}
			if account.Jira.Domain == "" {
				return fmt.Errorf("account %s is missing Jira domain", account.Name)
			}
		}

		// Track if we found a default account
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

// 2. Current directory: ./accounts.json.
func getDefaultConfigPath() string {
	// Try home directory first
	if homeDir, err := os.UserHomeDir(); err == nil {
		configDir := filepath.Join(homeDir, ".config", "atlacp")
		path := filepath.Join(configDir, "accounts.json")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Fall back to current directory
	return "accounts.json"
}

// convertToAppAccount converts the internal account configuration to the application layer type.
func convertToAppAccount(config atlassianAccountConfig) *app.AtlassianAccount {
	appAccount := &app.AtlassianAccount{
		Name:    config.Name,
		Default: config.Default,
	}

	if config.Bitbucket != nil {
		appAccount.Bitbucket = &app.BitbucketAccount{
			Token:     config.Bitbucket.Token,
			Workspace: config.Bitbucket.Workspace,
		}
	}

	if config.Jira != nil {
		appAccount.Jira = &app.JiraAccount{
			Token:  config.Jira.Token,
			Domain: config.Jira.Domain,
		}
	}

	return appAccount
}
