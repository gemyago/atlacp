package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gemyago/atlacp/internal/app"
)

// AccountsStore holds validated Atlassian accounts in memory. It can be loaded from
// the same JSON file shape as the file-backed repository ({ "accounts": [...] }).
// Mutations and persistence are not implemented in this iteration.
type AccountsStore struct {
	mu       sync.RWMutex
	accounts []app.AtlassianAccount
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
