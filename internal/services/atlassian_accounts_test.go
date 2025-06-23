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
			defaultAccount := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			nonDefaultAccount := app.NewRandomAtlassianAccount()

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
				app.NewRandomAtlassianAccount(),
				app.NewRandomAtlassianAccount(),
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

			defaultAccount := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true), app.WithAtlassianAccountName(defaultName))
			account1 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(userName))
			account2 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(botName))
			account3 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(adminName))

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

			defaultAccount := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true), app.WithAtlassianAccountName(defaultName))
			account1 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(userName))
			account2 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(botName))

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

	t.Run("NewAtlassianAccountsRepository", func(t *testing.T) {
		t.Run("should fail when file read fails", func(t *testing.T) {
			// Create a directory instead of a file to cause read failure
			tempDir := filepath.Join(t.TempDir(), "accounts.json")
			err := os.Mkdir(tempDir, 0700)
			require.NoError(t, err, "Failed to create directory")

			deps := makeMockDeps(tempDir)

			// Act
			_, err = NewAtlassianAccountsRepository(deps)

			// Assert
			require.Error(t, err, "Should error when file read fails")
			assert.Contains(t, err.Error(), "failed to read accounts configuration", "Error should mention read failure")
		})

		t.Run("should fail when JSON parsing fails", func(t *testing.T) {
			// Create invalid JSON file
			tempFile := filepath.Join(t.TempDir(), "accounts.json")
			err := os.WriteFile(tempFile, []byte("invalid json"), 0600)
			require.NoError(t, err, "Failed to write config file")

			deps := makeMockDeps(tempFile)

			// Act
			_, err = NewAtlassianAccountsRepository(deps)

			// Assert
			require.Error(t, err, "Should error when JSON parsing fails")
			assert.Contains(t, err.Error(), "failed to parse accounts configuration", "Error should mention parse failure")
		})
	})

	t.Run("validateAccountsConfig", func(t *testing.T) {
		t.Run("should fail when no accounts are configured", func(t *testing.T) {
			// Arrange
			config := &atlassianAccountsConfig{
				Accounts: []app.AtlassianAccount{},
			}

			// Act
			err := validateAccountsConfig(config)

			// Assert
			require.Error(t, err, "Should fail with empty accounts")
			assert.Contains(t, err.Error(), "no accounts configured", "Error should mention no accounts")
		})

		t.Run("should fail with duplicate account names", func(t *testing.T) {
			// Arrange
			name := "duplicate-" + faker.Username()
			account1 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(name))
			account2 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(name))

			config := &atlassianAccountsConfig{
				Accounts: []app.AtlassianAccount{account1, account2},
			}

			// Act
			err := validateAccountsConfig(config)

			// Assert
			require.Error(t, err, "Should fail with duplicate names")
			assert.Contains(t, err.Error(), "duplicate account name", "Error should mention duplicate name")
		})

		t.Run("should fail with multiple default accounts", func(t *testing.T) {
			// Arrange
			account1 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			account2 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))

			config := &atlassianAccountsConfig{
				Accounts: []app.AtlassianAccount{account1, account2},
			}

			// Act
			err := validateAccountsConfig(config)

			// Assert
			require.Error(t, err, "Should fail with multiple default accounts")
			assert.Contains(t, err.Error(), "multiple default accounts defined", "Error should mention multiple defaults")
		})

		t.Run("should fail with no default account", func(t *testing.T) {
			// Arrange
			account1 := app.NewRandomAtlassianAccount()
			account1.Default = false
			account2 := app.NewRandomAtlassianAccount()
			account2.Default = false

			config := &atlassianAccountsConfig{
				Accounts: []app.AtlassianAccount{account1, account2},
			}

			// Act
			err := validateAccountsConfig(config)

			// Assert
			require.Error(t, err, "Should fail with no default account")
			assert.Contains(t, err.Error(), "no default account specified", "Error should mention no default")
		})
	})

	t.Run("validateBasicAccountProperties", func(t *testing.T) {
		t.Run("should fail with empty account name", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Name = ""
			existingNames := make(map[string]bool)

			// Act
			err := validateBasicAccountProperties(account, existingNames)

			// Assert
			require.Error(t, err, "Should fail with empty name")
			assert.Contains(t, err.Error(), "account missing name", "Error should mention missing name")
		})

		t.Run("should fail with no services configured", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Bitbucket = nil
			account.Jira = nil
			existingNames := make(map[string]bool)

			// Act
			err := validateBasicAccountProperties(account, existingNames)

			// Assert
			require.Error(t, err, "Should fail with no services")
			assert.Contains(t, err.Error(), "must have at least one service configured",
				"Error should mention service requirement")
		})
	})

	t.Run("validateBitbucketConfig", func(t *testing.T) {
		t.Run("should fail with empty token", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Bitbucket.Token = ""

			// Act
			err := validateBitbucketConfig(account)

			// Assert
			require.Error(t, err, "Should fail with empty token")
			assert.Contains(t, err.Error(), "missing Bitbucket token", "Error should mention missing token")
		})

		t.Run("should fail with empty workspace", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Bitbucket.Workspace = ""

			// Act
			err := validateBitbucketConfig(account)

			// Assert
			require.Error(t, err, "Should fail with empty workspace")
			assert.Contains(t, err.Error(), "missing Bitbucket workspace", "Error should mention missing workspace")
		})
	})

	t.Run("validateJiraConfig", func(t *testing.T) {
		t.Run("should fail with empty token", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Jira.Token = ""

			// Act
			err := validateJiraConfig(account)

			// Assert
			require.Error(t, err, "Should fail with empty token")
			assert.Contains(t, err.Error(), "missing Jira token", "Error should mention missing token")
		})

		t.Run("should fail with empty domain", func(t *testing.T) {
			// Arrange
			account := app.NewRandomAtlassianAccount()
			account.Jira.Domain = ""

			// Act
			err := validateJiraConfig(account)

			// Assert
			require.Error(t, err, "Should fail with empty domain")
			assert.Contains(t, err.Error(), "missing Jira domain", "Error should mention missing domain")
		})
	})
}
