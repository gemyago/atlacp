package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
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

		t.Run("fails read when path is not a regular file", func(t *testing.T) {
			store := NewAccountsStore()
			err := store.LoadFromFile(t.TempDir())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to read accounts configuration")
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

	t.Run("mutations and save", func(t *testing.T) {
		t.Run("SaveToFile round-trip preserves accounts", func(t *testing.T) {
			defaultName := "default-" + faker.Username()
			acc := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(defaultName),
				app.WithAtlassianAccountDefault(true),
			)

			src := createTempAccountsFile(t, []app.AtlassianAccount{acc})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(src))

			outPath := filepath.Join(t.TempDir(), "saved.json")
			require.NoError(t, store.SaveToFile(outPath))

			reloaded := NewAccountsStore()
			require.NoError(t, reloaded.LoadFromFile(outPath))

			got, err := reloaded.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, acc, *got)
		})

		t.Run("SetDefault switches which account is default", func(t *testing.T) {
			nameA := "a-" + faker.Username()
			nameB := "b-" + faker.Username()
			a := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(nameA),
				app.WithAtlassianAccountDefault(true),
			)
			b := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(nameB),
				app.WithAtlassianAccountDefault(false),
			)

			path := createTempAccountsFile(t, []app.AtlassianAccount{a, b})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			require.NoError(t, store.SetDefault(nameB))
			def, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, nameB, def.Name)
			assert.True(t, def.Default)

			require.NoError(t, store.SetDefault(nameA))
			def, err = store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, nameA, def.Name)
		})

		t.Run("Remove sole account fails validation and leaves state unchanged", func(t *testing.T) {
			only := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{only})

			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			err := store.Remove(only.Name)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "no accounts configured")

			got, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, only, *got)
		})

		t.Run("Remove and SetDefault return ErrAccountNotFound for missing name", func(t *testing.T) {
			acc := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{acc})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			missing := "missing-" + faker.Username()
			err := store.Remove(missing)
			require.ErrorIs(t, err, app.ErrAccountNotFound)

			err = store.SetDefault(missing)
			require.ErrorIs(t, err, app.ErrAccountNotFound)
		})

		t.Run("invalid Upsert does not replace prior state", func(t *testing.T) {
			good := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{good})

			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			badUpsert := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName("orphan-"+faker.Username()),
				app.WithAtlassianAccountDefault(false),
			)
			badUpsert.Bitbucket = nil
			badUpsert.Jira = nil

			err := store.Upsert(badUpsert)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")

			got, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, good, *got)
		})

		t.Run("SaveToFile rejects empty path", func(t *testing.T) {
			store := NewAccountsStore()
			err := store.SaveToFile("")
			require.Error(t, err)
			assert.Equal(t, "accounts save path is empty", err.Error())
		})

		t.Run("SaveToFile fails when store state is invalid", func(t *testing.T) {
			store := NewAccountsStore()
			err := store.SaveToFile(filepath.Join(t.TempDir(), "out.json"))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
		})

		t.Run("SaveToFile fails when parent directory for temp file does not exist", func(t *testing.T) {
			acc := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{acc})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			missingParent := filepath.Join(t.TempDir(), "nope", "out.json")
			err := store.SaveToFile(missingParent)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to create temp file")
		})

		t.Run("Upsert replaces existing account by name", func(t *testing.T) {
			name := "stable-" + faker.Username()
			first := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(name),
				app.WithAtlassianAccountDefault(true),
			)
			path := createTempAccountsFile(t, []app.AtlassianAccount{first})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			replacement := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(name),
				app.WithAtlassianAccountDefault(true),
			)
			require.NoError(t, store.Upsert(replacement))

			got, err := store.GetAccountByName(t.Context(), name)
			require.NoError(t, err)
			assert.Equal(t, replacement, *got)
		})

		t.Run("Remove drops a non-default account and keeps validation", func(t *testing.T) {
			defName := "def-" + faker.Username()
			extraName := "extra-" + faker.Username()
			defAcc := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(defName),
				app.WithAtlassianAccountDefault(true),
			)
			extra := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(extraName),
				app.WithAtlassianAccountDefault(false),
			)
			path := createTempAccountsFile(t, []app.AtlassianAccount{defAcc, extra})
			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			require.NoError(t, store.Remove(extraName))
			_, err := store.GetAccountByName(t.Context(), extraName)
			require.ErrorIs(t, err, app.ErrAccountNotFound)

			got, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, defName, got.Name)
		})

		t.Run("concurrent Upsert keeps store valid", func(t *testing.T) {
			base := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{base})

			store := NewAccountsStore()
			require.NoError(t, store.LoadFromFile(path))

			const n = 32
			var seq atomic.Int64
			var upsertErrs atomic.Int32
			var wg sync.WaitGroup
			wg.Add(n)
			for range n {
				go func() {
					defer wg.Done()
					id := seq.Add(1)
					extra := app.NewRandomAtlassianAccount(
						app.WithAtlassianAccountName(fmt.Sprintf("extra-%d", id)),
						app.WithAtlassianAccountDefault(false),
					)
					if err := store.Upsert(extra); err != nil {
						upsertErrs.Add(1)
					}
				}()
			}
			wg.Wait()

			require.EqualValues(t, 0, upsertErrs.Load(), "expected all concurrent Upserts to succeed")
			_, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			byName, err := store.GetAccountByName(t.Context(), base.Name)
			require.NoError(t, err)
			assert.Equal(t, base.Name, byName.Name)
		})
	})
}
