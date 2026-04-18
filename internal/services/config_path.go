package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const goosDarwin = "darwin"

const accountsFileParentDirPerm = 0o750

// AccountsFilePathResolver resolves the default Atlassian accounts JSON path for the current
// platform and can create missing parent directories for any accounts file path.
// Use [NewAccountsFilePathResolver] for production defaults; tests may set fields directly.
type AccountsFilePathResolver struct {
	GOOS          string
	Home          string
	UserConfigDir func() (string, error)
}

// NewAccountsFilePathResolver returns a resolver using [runtime.GOOS], $HOME, and [os.UserConfigDir].
func NewAccountsFilePathResolver() *AccountsFilePathResolver {
	return &AccountsFilePathResolver{
		GOOS:          runtime.GOOS,
		Home:          os.Getenv("HOME"),
		UserConfigDir: os.UserConfigDir,
	}
}

func (r *AccountsFilePathResolver) defaultBaseDir() (string, error) {
	if r.GOOS == goosDarwin {
		return filepath.Join(r.Home, ".config"), nil
	}

	return r.UserConfigDir()
}

// DefaultPath returns the default filesystem path for the Atlassian accounts configuration file.
// On macOS the base directory is $HOME/.config (XDG-style) instead of [os.UserConfigDir], which on
// darwin resolves to ~/Library/Application Support.
func (r *AccountsFilePathResolver) DefaultPath() (string, error) {
	baseDir, err := r.defaultBaseDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, "atlacp", "accounts.json"), nil
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
