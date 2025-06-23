package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtlassianAccountsRepository_GetDefaultAccount(t *testing.T) {
	// Helper to create temporary account config file
	createTempAccountsFile := func(t *testing.T, accounts []app.AtlassianAccount) string {
		t.Helper()

		// Create configuration object
		config := app.AtlassianAccountsConfig{
			Accounts: accounts,
		}

		// Create temporary file
		tempFile := filepath.Join(t.TempDir(), "accounts.json")

		// Write config to file
		data, err := json.Marshal(config)
		require.NoError(t, err, "Failed to marshal config data")

		err = os.WriteFile(tempFile, data, 0600)
		require.NoError(t, err, "Failed to write config file")

		return tempFile
	}

	makeMockDeps := func(configPath string) AtlassianAccountsRepositoryDeps {
		return AtlassianAccountsRepositoryDeps{
			RootLogger: diag.RootTestLogger(),
			ConfigPath: configPath,
		}
	}

	t.Run("should return default account when configuration is valid", func(t *testing.T) {
		// Arrange
		accounts := []app.AtlassianAccount{
			{
				Name:    "non-default-account",
				Default: false,
				Bitbucket: &app.BitbucketAccount{
					Token:     "token1",
					Workspace: "workspace1",
				},
			},
			{
				Name:    "default-account",
				Default: true,
				Jira: &app.JiraAccount{
					Token:  "token2",
					Domain: "example",
				},
			},
		}

		configPath := createTempAccountsFile(t, accounts)
		deps := makeMockDeps(configPath)

		repository, err := NewAtlassianAccountsRepository(deps)
		require.NoError(t, err, "Failed to create repository")

		// Act
		account, err := repository.GetDefaultAccount(t.Context())

		// Assert
		require.NoError(t, err, "GetDefaultAccount should not return an error")
		require.NotNil(t, account, "Default account should not be nil")
		assert.Equal(t, "default-account", account.Name, "Should return the correct default account")
		assert.True(t, account.Default, "Default account should have Default=true")
		assert.Nil(t, account.Bitbucket, "Bitbucket should be nil for this account")
		require.NotNil(t, account.Jira, "Jira should not be nil for this account")
		assert.Equal(t, "token2", account.Jira.Token, "Should have correct Jira token")
		assert.Equal(t, "example", account.Jira.Domain, "Should have correct Jira domain")
	})

	t.Run("should return error when no default account exists", func(t *testing.T) {
		// Arrange - create a repository with accounts but no default
		accounts := []app.AtlassianAccount{
			{
				Name:    "account1",
				Default: false,
				Bitbucket: &app.BitbucketAccount{
					Token:     "token1",
					Workspace: "workspace1",
				},
			},
			{
				Name:    "account2",
				Default: false,
				Jira: &app.JiraAccount{
					Token:  "token2",
					Domain: "example",
				},
			},
		}

		// Create repository manually, bypassing validation
		repository := &atlassianAccountsRepository{
			config: &app.AtlassianAccountsConfig{
				Accounts: accounts,
			},
			logger: diag.RootTestLogger().WithGroup("atlassian-accounts"),
		}

		// Act
		account, err := repository.GetDefaultAccount(t.Context())

		// Assert
		assert.Nil(t, account, "Account should be nil when no default account exists")
		require.Error(t, err, "Should return error when no default account exists")
		assert.ErrorIs(t, err, app.ErrNoDefaultAccount, "Error should be ErrNoDefaultAccount")
	})

	t.Run("should return error when config file doesn't exist", func(t *testing.T) {
		// Arrange
		deps := AtlassianAccountsRepositoryDeps{
			RootLogger: diag.RootTestLogger(),
			ConfigPath: "/path/that/does/not/exist",
		}

		// Act
		_, err := NewAtlassianAccountsRepository(deps)

		// Assert
		require.Error(t, err, "Should return error when config file doesn't exist")
		assert.Contains(t, err.Error(), "not found", "Error should mention file not found")
	})
}
