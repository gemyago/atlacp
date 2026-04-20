package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const accountsFileParentDirPerm = 0o750

// AccountsFilePathResolver resolves the default Atlassian accounts JSON path for the current
// platform and can create missing parent directories for any accounts file path.
// Use [NewAccountsFilePathResolver] for production defaults; tests may set fields directly.
type AccountsFilePathResolver struct {
	Home string
}

// NewAccountsFilePathResolver returns a resolver using $HOME.
func NewAccountsFilePathResolver() *AccountsFilePathResolver {
	return &AccountsFilePathResolver{
		Home: os.Getenv("HOME"),
	}
}

func (r *AccountsFilePathResolver) defaultBaseDir() string {
	return filepath.Join(r.Home, ".atlacp")
}

// DefaultPath returns the default filesystem path for the Atlassian accounts configuration file.
func (r *AccountsFilePathResolver) DefaultPath() (string, error) {
	return filepath.Join(r.defaultBaseDir(), "accounts.json"), nil
}

// EnsureParentDirsForFile creates the parent directory of filePath if needed (idempotent).
// It MUST be called before writing the accounts file when the directory may not exist yet.
func (r *AccountsFilePathResolver) EnsureParentDirsForFile(filePath string) error {
	if filePath == "" {
		return errors.New("accounts file path is empty")
	}

	dir := filepath.Dir(filepath.Clean(filePath))
	if err := os.MkdirAll(dir, accountsFileParentDirPerm); err != nil {
		return fmt.Errorf("create accounts file parent directory %q: %w", dir, err)
	}

	return nil
}
