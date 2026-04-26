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
	"github.com/gemyago/atlacp/internal/diag"
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

	t.Run("NewAccountsStoreWithDeps", func(t *testing.T) {
		makeDeps := func(configPath string) AccountsStoreDeps {
			return AccountsStoreDeps{
				RootLogger: diag.RootTestLogger(),
				ConfigPath: configPath,
			}
		}

		t.Run("fails when config path is empty", func(t *testing.T) {
			store, err := NewAccountsStoreWithDeps(makeDeps(""))
			require.Error(t, err)
			assert.Nil(t, store)
			assert.Equal(t, "accounts configuration path not specified", err.Error())
		})

		t.Run("succeeds with empty store when config file does not exist", func(t *testing.T) {
			missing := filepath.Join(t.TempDir(), "nonexistent-"+faker.Username()+".json")
			store, err := NewAccountsStoreWithDeps(makeDeps(missing))
			require.NoError(t, err)
			require.NotNil(t, store)
			assert.Empty(t, store.ListAccounts())
		})

		t.Run("loads valid config from file", func(t *testing.T) {
			defaultAccount := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			other := app.NewRandomAtlassianAccount()
			path := createTempAccountsFile(t, []app.AtlassianAccount{other, defaultAccount})

			store, err := NewAccountsStoreWithDeps(makeDeps(path))
			require.NoError(t, err)
			require.NotNil(t, store)

			got, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, defaultAccount, *got)
		})

		t.Run("fails when file read fails because path is a directory", func(t *testing.T) {
			tempDir := filepath.Join(t.TempDir(), "accounts.json")
			err := os.Mkdir(tempDir, 0o700)
			require.NoError(t, err)

			store, err := NewAccountsStoreWithDeps(makeDeps(tempDir))
			require.Error(t, err)
			assert.Nil(t, store)
			assert.Contains(t, err.Error(), "failed to read accounts configuration")
		})

		t.Run("fails when JSON parsing fails", func(t *testing.T) {
			tempFile := filepath.Join(t.TempDir(), "accounts.json")
			err := os.WriteFile(tempFile, []byte("invalid json"), 0o600)
			require.NoError(t, err)

			store, err := NewAccountsStoreWithDeps(makeDeps(tempFile))
			require.Error(t, err)
			assert.Nil(t, store)
			assert.Contains(t, err.Error(), "failed to parse accounts configuration")
		})

		t.Run("rejects invalid accounts via shared validation", func(t *testing.T) {
			invalid := []app.AtlassianAccount{
				app.NewRandomAtlassianAccount(),
				app.NewRandomAtlassianAccount(),
			}
			invalid[0].Default = false
			invalid[1].Default = false

			path := createTempAccountsFile(t, invalid)
			store, err := NewAccountsStoreWithDeps(makeDeps(path))
			require.Error(t, err)
			assert.Nil(t, store)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "no default account specified")
		})
	})

	t.Run("ListAccounts", func(t *testing.T) {
		t.Run("returns empty slice when no accounts loaded", func(t *testing.T) {
			store := &AccountsStore{}
			got := store.ListAccounts()
			require.NotNil(t, got)
			assert.Empty(t, got)
		})

		t.Run("returns all accounts after LoadFromFile", func(t *testing.T) {
			defaultAccount := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			other := app.NewRandomAtlassianAccount()
			path := createTempAccountsFile(t, []app.AtlassianAccount{other, defaultAccount})

			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(path))

			got := store.ListAccounts()
			require.Len(t, got, 2)
			assert.Equal(t, []app.AtlassianAccount{other, defaultAccount}, got)
		})

		t.Run("returns updated list after Upsert", func(t *testing.T) {
			initial := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{initial})

			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(path))

			added := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(false))
			require.NoError(t, store.Upsert(added))

			got := store.ListAccounts()
			require.Len(t, got, 2)
			assert.Equal(t, []app.AtlassianAccount{initial, added}, got)
		})

		t.Run("snapshot is independent of later mutations", func(t *testing.T) {
			acc := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{acc})

			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(path))

			snapshot := store.ListAccounts()
			require.Len(t, snapshot, 1)

			replacement := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(acc.Name),
				app.WithAtlassianAccountDefault(true),
			)
			require.NoError(t, store.Upsert(replacement))

			assert.Equal(t, []app.AtlassianAccount{acc}, snapshot)
			assert.Equal(t, replacement, store.ListAccounts()[0])
		})
	})

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

			store := &AccountsStore{}
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
			store := &AccountsStore{}
			missing := filepath.Join(t.TempDir(), "missing-"+faker.Username()+".json")

			err := store.LoadFromFile(missing)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "not found")
		})

		t.Run("fails read when path is not a regular file", func(t *testing.T) {
			store := &AccountsStore{}
			err := store.LoadFromFile(t.TempDir())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to read accounts configuration")
		})

		t.Run("fails on invalid JSON", func(t *testing.T) {
			tempFile := filepath.Join(t.TempDir(), "accounts.json")
			err := os.WriteFile(tempFile, []byte("not json {"), 0o600)
			require.NoError(t, err)

			store := &AccountsStore{}
			err = store.LoadFromFile(tempFile)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to parse accounts configuration")
		})

		t.Run("fails validation when accounts list is empty", func(t *testing.T) {
			path := createTempAccountsFile(t, nil)

			store := &AccountsStore{}
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

			store := &AccountsStore{}
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "duplicate account name")
		})

		t.Run("fails validation when two defaults", func(t *testing.T) {
			a := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			b := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))

			path := createTempAccountsFile(t, []app.AtlassianAccount{a, b})

			store := &AccountsStore{}
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

			store := &AccountsStore{}
			err := store.LoadFromFile(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
			assert.Contains(t, err.Error(), "no default account specified")
		})

		t.Run("failed load does not replace prior valid state", func(t *testing.T) {
			good := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			goodPath := createTempAccountsFile(t, []app.AtlassianAccount{good})

			store := &AccountsStore{}
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
			store := &AccountsStore{}
			acc, err := store.GetDefaultAccount(t.Context())
			assert.Nil(t, acc)
			require.ErrorIs(t, err, app.ErrNoDefaultAccount)
		})

		t.Run("returns ErrNoDefaultAccount when accounts exist but none is default", func(t *testing.T) {
			a := app.NewRandomAtlassianAccount()
			b := app.NewRandomAtlassianAccount()
			a.Default = false
			b.Default = false

			store := &AccountsStore{accounts: []app.AtlassianAccount{a, b}}
			acc, err := store.GetDefaultAccount(t.Context())
			assert.Nil(t, acc)
			require.ErrorIs(t, err, app.ErrNoDefaultAccount)
		})
	})

	t.Run("GetAccountByName", func(t *testing.T) {
		t.Run("returns ErrAccountNotFound when store is empty", func(t *testing.T) {
			store := &AccountsStore{}
			name := "missing-" + faker.Username()
			acc, err := store.GetAccountByName(t.Context(), name)
			assert.Nil(t, acc)
			require.ErrorIs(t, err, app.ErrAccountNotFound)
			assert.Contains(t, err.Error(), name)
		})

		t.Run("returns ErrAccountNotFound when name does not exist among loaded accounts", func(t *testing.T) {
			defaultName := "default-" + faker.Username()
			userName := "user-" + faker.Username()
			botName := "bot-" + faker.Username()

			defaultAccount := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountDefault(true),
				app.WithAtlassianAccountName(defaultName),
			)
			account1 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(userName))
			account2 := app.NewRandomAtlassianAccount(app.WithAtlassianAccountName(botName))

			path := createTempAccountsFile(t, []app.AtlassianAccount{defaultAccount, account1, account2})
			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(path))

			nonExistentName := "nonexistent-" + faker.Username()
			result, err := store.GetAccountByName(t.Context(), nonExistentName)
			assert.Nil(t, result)
			require.ErrorIs(t, err, app.ErrAccountNotFound)
			assert.Contains(t, err.Error(), nonExistentName)
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
			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(src))

			outPath := filepath.Join(t.TempDir(), "saved.json")
			require.NoError(t, store.SaveToFile(outPath))

			reloaded := &AccountsStore{}
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
			store := &AccountsStore{}
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

			store := &AccountsStore{}
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
			store := &AccountsStore{}
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

			store := &AccountsStore{}
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
			store := &AccountsStore{}
			err := store.SaveToFile("")
			require.Error(t, err)
			assert.Equal(t, "accounts save path is empty", err.Error())
		})

		t.Run("SaveToFile fails when store state is invalid", func(t *testing.T) {
			store := &AccountsStore{}
			err := store.SaveToFile(filepath.Join(t.TempDir(), "out.json"))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid accounts configuration")
		})

		t.Run("SaveToFile round-trips to nested path when parent dirs prepared", func(t *testing.T) {
			acc := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(true))
			path := createTempAccountsFile(t, []app.AtlassianAccount{acc})
			store := &AccountsStore{}
			require.NoError(t, store.LoadFromFile(path))

			base := t.TempDir()
			outPath := filepath.Join(base, "deep", "nested", "accounts.json")
			require.NoError(t, NewAtlacpPathResolver().EnsureParentDirsForFile(outPath))
			require.NoError(t, store.SaveToFile(outPath))

			reloaded := &AccountsStore{}
			require.NoError(t, reloaded.LoadFromFile(outPath))
			got, err := reloaded.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, acc, *got)
		})

		t.Run("Upsert replaces existing account by name", func(t *testing.T) {
			name := "stable-" + faker.Username()
			first := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountName(name),
				app.WithAtlassianAccountDefault(true),
			)
			path := createTempAccountsFile(t, []app.AtlassianAccount{first})
			store := &AccountsStore{}
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

		t.Run("Upsert on empty store sets default when sole account has Default false", func(t *testing.T) {
			acc := app.NewRandomAtlassianAccount(app.WithAtlassianAccountDefault(false))
			store := &AccountsStore{}
			require.NoError(t, store.Upsert(acc))

			got := store.ListAccounts()
			require.Len(t, got, 1)
			want := acc
			want.Default = true
			assert.Equal(t, want, got[0])

			def, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.True(t, def.Default)
		})

		t.Run("Upsert second account with Default false keeps first as sole default", func(t *testing.T) {
			name1 := "first-" + faker.Username()
			name2 := "second-" + faker.Username()
			first := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountDefault(false),
				app.WithAtlassianAccountName(name1),
			)
			store := &AccountsStore{}
			require.NoError(t, store.Upsert(first))

			second := app.NewRandomAtlassianAccount(
				app.WithAtlassianAccountDefault(false),
				app.WithAtlassianAccountName(name2),
			)
			require.NoError(t, store.Upsert(second))

			require.Len(t, store.ListAccounts(), 2)
			def, err := store.GetDefaultAccount(t.Context())
			require.NoError(t, err)
			assert.Equal(t, name1, def.Name)
			assert.True(t, def.Default)
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
			store := &AccountsStore{}
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

			store := &AccountsStore{}
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
