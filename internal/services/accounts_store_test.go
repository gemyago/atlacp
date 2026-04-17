package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountsStore(t *testing.T) {
	createTempAccountsFile := func(t *testing.T, accounts []app.AtlassianAccount) string {
		t.Helper()

		config := atlassianAccountsConfig{
			Accounts: accounts,
		}

		tempFile := filepath.Join(t.TempDir(), "accounts.json")

		data, err := json.Marshal(config)
		require.NoError(t, err, "Failed to marshal config data")

		err = os.WriteFile(tempFile, data, 0o600)
		require.NoError(t, err, "Failed to write config file")

		return tempFile
	}

	t.Run("LoadFromFile and queries", func(t *testing.T) {
		t.Run("happy path loads and returns default and account by name", func(t *testing.T) {
			defaultName := "default-" + faker.Username()
			otherName := "other-" + faker.Username()

			defaultAccount := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountDefault(true),
				app.WithAtlassianAccountName(defaultName),
			)
			other := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(otherName))

			path := createTempAccountsFile(t, []app.AtlassianAccount{other, defaultAccount})

			store := NewAccountsStore()
			err := store.LoadFromFile(path)
			require.NoError(t, err)

			gotDefault, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			require.NotNil(t, gotDefault)
			assert.Equal(t, defaultAccount, *gotDefault)

			gotByName, err := store.GetAccountByName(t.Context(), otherName)
			require.NoError(t, err)
			require.NotNil(t, gotByName)
			assert.Equal(t, other, *gotByName)
		})

		t.Run("returns not found when file path does not exist", func(t *testing.T) {
			store := NewAccountsStore()
			missing := filepath.Join(t.TempDir(), "missing-"+faker.Username()+".json")

			err := store.LoadFromFile(missing)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "not found")
		})

		t.Run("fails on invalid JSON", func(t *testing.T) {
			tempFile := filepath.Join(t.TempDir(), "accounts.json")
			err := os.WriteFile(tempFile, []byte("not json {"), 0o600)
			require.NoError(t, err)

			store := NewAccountsStore()
			err = store.LoadFromFile(tempFile)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to parse accounts configuration")
		})

		t.Run("fails validation when accounts list is empty", func(t *testing.T) {
			path := createTempAccountsFile(t, nil)

			store := NewAccountsStore()
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "no accounts configured")
		})

		t.Run("fails validation on duplicate names", func(t *testing.T) {
			name := "dup-" + faker.Username()
			a := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(name),
				app.WithAtlassianAccountDefault(true),
			)
			b := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(name))

			path := createTempAccountsFile(t, []app.AtlassianAccount{a, b})

			store := NewAccountsStore()
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "duplicate account name")
		})

		t.Run("fails validation when two defaults", func(t *testing.T) {
			a := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			b := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))

			path := createTempAccountsFile(t, []app.AtlassianAccount{a, b})

			store := NewAccountsStore()
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "multiple default accounts defined")
		})

		t.Run("fails validation when no default among accounts", func(t *testing.T) {
			a := app.NewRandomAtlassianAccount()
			b := app.NewRandomAtlassianAccount()
			a.Default = false
			b.Default = false

			path := createTempAccountsFile(t, []app.AtlassianAccount{a, b})

			store := NewAccountsStore()
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "no default account specified")
		})

		t.Run("failed load does not replace prior valid state", func(t *testing.T) {
			good := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			goodPath := createTempAccountsFile(t, []app.AtlassianAccount{good})

			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(goodPath))

			badPath := createTempAccountsFile(t, nil)
			err := store.LoadFromFile(badPath)
			require.Error(t, err)

			acc, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, good, *acc)
		})
	})

	t.Run("GetDefaultAccount", func(t *testing.T) {
		t.Run("returns ErrNoDefaultAccount when store is empty", func(t *testing.T) {
			store := NewAccountsStore()
			acc, err := store.GetDefaultAccount(t.Context())
			assert.Nil(t, acc)
			require.ErrorIs(t, err, app.ErrNoDefaultAccount)
		})
	})

	t.Run("GetAccountByName", func(t *testing.T) {
		t.Run("returns ErrAccountNotFound when store is empty", func(t *testing.T) {
			store := NewAccountsStore()
			name := "missing-" + faker.Username()
			acc, err := store.GetAccountByName(t.Context(), name)
			assert.Nil(t, acc)
			require.ErrorIs(t, err, app.ErrAccountNotFound)
			assert.Contains(t, err.Error(), name)
		})
	})
}
