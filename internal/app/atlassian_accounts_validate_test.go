package app

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAtlassianAccounts(t *testing.T) {
	t.Run("should fail when no accounts are configured", func(t *testing.T) {
		err := ValidateAtlassianAccounts([]AtlassianAccount{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no accounts configured")
	})

	t.Run("should fail with duplicate account names", func(t *testing.T) {
		name := "duplicate-" + faker.Username()
		account1 := NewRandomAtlassianAccount(WithAtlassianAccountName(name))
		account2 := NewRandomAtlassianAccount(WithAtlassianAccountName(name))

		err := ValidateAtlassianAccounts([]AtlassianAccount{account1, account2})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate account name")
	})

	t.Run("should fail with multiple default accounts", func(t *testing.T) {
		account1 := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true))
		account2 := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true))

		err := ValidateAtlassianAccounts([]AtlassianAccount{account1, account2})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple default accounts defined")
	})

	t.Run("should fail with no default account", func(t *testing.T) {
		account1 := NewRandomAtlassianAccount()
		account1.Default = false
		account2 := NewRandomAtlassianAccount()
		account2.Default = false

		err := ValidateAtlassianAccounts([]AtlassianAccount{account1, account2})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no default account specified")
	})

	t.Run("should fail with empty account name", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Name = ""

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "account missing name")
	})

	t.Run("should fail with no services configured", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Bitbucket = nil
		account.Jira = nil

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must have at least one service configured")
	})

	t.Run("should fail when Bitbucket token value empty", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Bitbucket.Value = ""

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing Bitbucket token value")
	})

	t.Run("should fail when Bitbucket token type empty", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Bitbucket.Type = ""

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing Bitbucket token type")
	})

	t.Run("should fail when Jira token value empty", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Jira.Value = ""

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing Jira token value")
	})

	t.Run("should fail when Jira token type empty", func(t *testing.T) {
		account := NewRandomAtlassianAccount()
		account.Jira.Type = ""

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing Jira token type")
	})

	t.Run("should succeed for valid single default account", func(t *testing.T) {
		account := NewRandomAtlassianAccount(WithAtlassianAccountDefault(true))

		err := ValidateAtlassianAccounts([]AtlassianAccount{account})
		require.NoError(t, err)
	})
}
