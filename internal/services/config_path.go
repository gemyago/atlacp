package services

import (
	"os"
	"path/filepath"
	"runtime"
)

const goosDarwin = "darwin"

// defaultAccountsBaseDir returns the base config directory (parent of the atlacp segment).
func defaultAccountsBaseDir(goos string, home string, userConfigDir func() (string, error)) (string, error) {
	if goos == goosDarwin {
		return filepath.Join(home, ".config"), nil
	}

	return userConfigDir()
}

func defaultAccountsFilePath(goos string, home string, userConfigDir func() (string, error)) (string, error) {
	baseDir, err := defaultAccountsBaseDir(goos, home, userConfigDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, "atlacp", "accounts.json"), nil
}

// DefaultAccountsFilePath returns the default filesystem path for the Atlassian accounts
// configuration file. On macOS the base directory is $HOME/.config (XDG-style) instead of
// [os.UserConfigDir], which on darwin resolves to ~/Library/Application Support.
func DefaultAccountsFilePath() (string, error) {
	return defaultAccountsFilePath(runtime.GOOS, os.Getenv("HOME"), os.UserConfigDir)
}
