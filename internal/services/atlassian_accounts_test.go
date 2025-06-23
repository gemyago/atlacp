package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtlassianAccountsRepository(t *testing.T) {
	// Helper to create temporary account config file
	createTempAccountsFile := func(t *testing.T, accounts []app.AtlassianAccount) string {
		t.Helper()

		// Create configuration object
		config := atlassianAccountsConfig{
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

	t.Run("GetDefaultAccount", func(t *testing.T) {
		t.Run("should return default account when configuration is valid", func(t *testing.T) {
			// Arrange
			defaultAccount := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true))
			nonDefaultAccount := NewRandomAtlassianAccount()

			accounts := []app.AtlassianAccount{nonDefaultAccount, defaultAccount}

			configPath := createTempAccountsFile(t, accounts)
			deps := makeMockDeps(configPath)

			repository, err := NewAtlassianAccountsRepository(deps)
			require.NoError(t, err, "Failed to create repository")

			// Act
			account, err := repository.GetDefaultAccount(t.Context())

			// Assert
			require.NoError(t, err, "GetDefaultAccount should not return an error")
			require.NotNil(t, account, "Default account should not be nil")

			// Compare the entire account
			assert.Equal(t, defaultAccount, *account, "Should return the correct default account")
		})

		t.Run("should return error when no default account exists", func(t *testing.T) {
			// Arrange - create a repository with accounts but no default
			// Generate two accounts without the default flag set
			nonDefaultAccounts := []app.AtlassianAccount{
				NewRandomAtlassianAccount(),
				NewRandomAtlassianAccount(),
			}

			// Create repository manually, bypassing validation
			repository := &atlassianAccountsRepository{
				config: &atlassianAccountsConfig{
					Accounts: nonDefaultAccounts,
				},
				logger: diag.RootTestLogger().WithGroup("atlassian-accounts"),
			}

			// Act
			account, err := repository.GetDefaultAccount(t.Context())

			// Assert
			assert.Nil(t, account, "Account should be nil when no default account exists")
			require.ErrorIs(t, err, app.ErrNoDefaultAccount, "Error should be ErrNoDefaultAccount")
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
	})

	t.Run("GetAccountByName", func(t *testing.T) {
		t.Run("should return account when name exists", func(t *testing.T) {
			// Arrange
			// Generate random account names
			defaultName := "default-" + faker.Username()
			userName := "user-" + faker.Username()
			botName := "bot-" + faker.Username()
			adminName := "admin-" + faker.Username()

			defaultAccount := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true), WithAtlassianAccountName(defaultName))
			account1 := NewRandomAtlassianAccount(WithAtlassianAccountName(userName))
			account2 := NewRandomAtlassianAccount(WithAtlassianAccountName(botName))
			account3 := NewRandomAtlassianAccount(WithAtlassianAccountName(adminName))

			accounts := []app.AtlassianAccount{defaultAccount, account1, account2, account3}

			configPath := createTempAccountsFile(t, accounts)
			deps := makeMockDeps(configPath)

			repository, err := NewAtlassianAccountsRepository(deps)
			require.NoError(t, err, "Failed to create repository")

			// Act
			result, err := repository.GetAccountByName(t.Context(), botName)

			// Assert
			require.NoError(t, err, "GetAccountByName should not return an error")
			require.NotNil(t, result, "Account should not be nil")

			// Compare the entire account
			assert.Equal(t, account2, *result, "Should return the correct account")
		})

		t.Run("should return error when account name doesn't exist", func(t *testing.T) {
			// Arrange
			// Generate random account names
			defaultName := "default-" + faker.Username()
			userName := "user-" + faker.Username()
			botName := "bot-" + faker.Username()

			// Generate a non-existent name that's guaranteed to be different
			nonExistentName := "nonexistent-" + faker.Username()

			defaultAccount := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true), WithAtlassianAccountName(defaultName))
			account1 := NewRandomAtlassianAccount(WithAtlassianAccountName(userName))
			account2 := NewRandomAtlassianAccount(WithAtlassianAccountName(botName))

			accounts := []app.AtlassianAccount{defaultAccount, account1, account2}

			configPath := createTempAccountsFile(t, accounts)
			deps := makeMockDeps(configPath)

			repository, err := NewAtlassianAccountsRepository(deps)
			require.NoError(t, err, "Failed to create repository")

			// Act
			result, err := repository.GetAccountByName(t.Context(), nonExistentName)

			// Assert
			assert.Nil(t, result, "Account should be nil when name doesn't exist")
			require.ErrorIs(t, err, app.ErrAccountNotFound, "Error should be ErrAccountNotFound")
			assert.Contains(t, err.Error(), nonExistentName, "Error should contain the account name")
		})
	})
}
