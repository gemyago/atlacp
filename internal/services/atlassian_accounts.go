package services

import (
	"context"
	"errors"

	"github.com/gemyago/atlacp/internal/app"
)

// atlassianAccountsRepository implements the app.AtlassianAccountsRepository interface.
type atlassianAccountsRepository struct {
	// Will be implemented later
}

// NewAtlassianAccountsRepository creates a new Atlassian accounts repository.
func NewAtlassianAccountsRepository() app.AtlassianAccountsRepository {
	return &atlassianAccountsRepository{}
}

// GetDefaultAccount returns the default Atlassian account configuration.
func (r *atlassianAccountsRepository) GetDefaultAccount(_ context.Context) (*app.AtlassianAccount, error) {
	return nil, errors.New("not implemented")
}

// GetAccountByName returns an account with the specified name.
func (r *atlassianAccountsRepository) GetAccountByName(_ context.Context, _ string) (*app.AtlassianAccount, error) {
	return nil, errors.New("not implemented")
}
