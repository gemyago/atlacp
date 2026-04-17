// AccountsStore (defined below) provides in-memory management of validated Atlassian accounts
// with optional JSON file persistence (load/save using the same {"accounts":[...]} envelope as
// the file-backed repository). This standalone component is not registered in application
// dependency injection; construct it directly where needed until wiring is added.

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/gemyago/atlacp/internal/app"
	"go.uber.org/dig"
)

// AccountsStore holds validated Atlassian accounts in memory. It can be loaded from and saved
// to the same JSON file shape as the file-backed repository ({ "accounts": [...] }).
// Mutations replace the in-memory list only after full re-validation; SaveToFile writes atomically.
type AccountsStore struct {
	mu       sync.RWMutex
	accounts []app.AtlassianAccount
	logger   *slog.Logger // set by NewAccountsStoreWithDeps; reserved for future logging parity with the former repository
}

// AccountsStoreDeps contains dependencies for constructing an AccountsStore from the configured accounts file (DI).
type AccountsStoreDeps struct {
	dig.In

	RootLogger *slog.Logger
	ConfigPath string `name:"config.atlassian.accountsFilePath"`
}

// NewAccountsStoreWithDeps loads validated accounts from ConfigPath into a new store. It mirrors startup
// behavior of the former file-backed repository: empty path and missing file fail before load; other errors
// come from LoadFromFile.
func NewAccountsStoreWithDeps(deps AccountsStoreDeps) (*AccountsStore, error) {
	logger := deps.RootLogger.WithGroup("atlassian-accounts")
	configPath := deps.ConfigPath

	if configPath == "" {
		return nil, errors.New("accounts configuration path not specified")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("accounts configuration file not found at %s", configPath)
	}

	store := &AccountsStore{logger: logger}
	if err := store.LoadFromFile(configPath); err != nil {
		return nil, err
	}

	return store, nil
}

// NewAccountsStore returns an empty store. Call LoadFromFile to populate it.
func NewAccountsStore() *AccountsStore {
	return &AccountsStore{}
}

// LoadFromFile reads JSON from path, validates with app.ValidateAtlassianAccounts,
// and replaces the in-memory account list. On any error the previous state is unchanged.
func (s *AccountsStore) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("accounts configuration file not found at %s: %w", path, err)
		}

		return fmt.Errorf("failed to read accounts configuration: %w", err)
	}

	var config atlassianAccountsConfig
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return fmt.Errorf("failed to parse accounts configuration: %w", unmarshalErr)
	}

	if validateErr := app.ValidateAtlassianAccounts(config.Accounts); validateErr != nil {
		return fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	copied := make([]app.AtlassianAccount, len(config.Accounts))
	copy(copied, config.Accounts)

	s.mu.Lock()
	s.accounts = copied
	s.mu.Unlock()

	return nil
}

// GetDefaultAccount returns the default Atlassian account.
func (s *AccountsStore) GetDefaultAccount(_ context.Context) (*app.AtlassianAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.accounts {
		if s.accounts[i].Default {
			return &s.accounts[i], nil
		}
	}

	return nil, app.ErrNoDefaultAccount
}

// GetAccountByName returns the account with the given name.
func (s *AccountsStore) GetAccountByName(_ context.Context, name string) (*app.AtlassianAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.accounts {
		if s.accounts[i].Name == name {
			return &s.accounts[i], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", app.ErrAccountNotFound, name)
}

// Upsert inserts or replaces an account by name. The full resulting configuration is validated
// before it replaces the stored state.
func (s *AccountsStore) Upsert(account app.AtlassianAccount) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneAtlassianAccounts(s.accounts)
	found := false
	for i := range next {
		if next[i].Name == account.Name {
			next[i] = account
			found = true
			break
		}
	}
	if !found {
		next = append(next, account)
	}

	if validateErr := app.ValidateAtlassianAccounts(next); validateErr != nil {
		return fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	s.accounts = cloneAtlassianAccounts(next)
	return nil
}

// Remove deletes the account with the given name. The full resulting configuration is validated
// before committing (e.g. removing the last account or the only default fails validation).
func (s *AccountsStore) Remove(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneAtlassianAccounts(s.accounts)
	found := false
	for i := range next {
		if next[i].Name == name {
			next = append(next[:i], next[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: %s", app.ErrAccountNotFound, name)
	}

	if validateErr := app.ValidateAtlassianAccounts(next); validateErr != nil {
		return fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	s.accounts = cloneAtlassianAccounts(next)
	return nil
}

// SetDefault marks the named account as the sole default. The full resulting configuration is validated
// before committing.
func (s *AccountsStore) SetDefault(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneAtlassianAccounts(s.accounts)
	found := false
	for i := range next {
		if next[i].Name == name {
			found = true
		}
		next[i].Default = next[i].Name == name
	}
	if !found {
		return fmt.Errorf("%w: %s", app.ErrAccountNotFound, name)
	}

	if validateErr := app.ValidateAtlassianAccounts(next); validateErr != nil {
		return fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	s.accounts = cloneAtlassianAccounts(next)
	return nil
}

// SaveToFile writes the current accounts to path as JSON ({ "accounts": [...] }) using a
// temporary file in the same directory followed by rename so the target file is not left partial.
func (s *AccountsStore) SaveToFile(path string) error {
	if path == "" {
		return errors.New("accounts save path is empty")
	}

	s.mu.RLock()
	snapshot := cloneAtlassianAccounts(s.accounts)
	s.mu.RUnlock()

	if validateErr := app.ValidateAtlassianAccounts(snapshot); validateErr != nil {
		return fmt.Errorf("invalid accounts configuration: %w", validateErr)
	}

	config := atlassianAccountsConfig{Accounts: snapshot}
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal accounts configuration: %w", err)
	}

	if writeErr := writeAccountsFileAtomically(path, data); writeErr != nil {
		return writeErr
	}

	return nil
}

func cloneAtlassianAccounts(src []app.AtlassianAccount) []app.AtlassianAccount {
	if len(src) == 0 {
		return nil
	}
	dst := make([]app.AtlassianAccount, len(src))
	copy(dst, src)
	return dst
}

func writeAccountsFileAtomically(path string, data []byte) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file for accounts save: %w", err)
	}

	tmpName := f.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if _, writeErr := f.Write(data); writeErr != nil {
		_ = f.Close()
		return fmt.Errorf("failed to write accounts temp file: %w", writeErr)
	}
	if syncErr := f.Sync(); syncErr != nil {
		_ = f.Close()
		return fmt.Errorf("failed to sync accounts temp file: %w", syncErr)
	}
	if closeErr := f.Close(); closeErr != nil {
		return fmt.Errorf("failed to close accounts temp file: %w", closeErr)
	}

	//nolint:gosec // caller-chosen save path; temp file is in the same directory as destination
	if renameErr := os.Rename(tmpName, path); renameErr != nil {
		return fmt.Errorf("failed to replace accounts file: %w", renameErr)
	}
	cleanup = false
	return nil
}
