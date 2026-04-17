package services

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAccountsFilePath(t *testing.T) {
	t.Run("defaultAccountsBaseDir", func(t *testing.T) {
		t.Run("darwin uses HOME/.config and does not call UserConfigDir", func(t *testing.T) {
			home := t.TempDir()
			userConfigDirCalled := false
			got, err := defaultAccountsBaseDir(goosDarwin, home, func() (string, error) {
				userConfigDirCalled = true
				return "", nil
			})
			require.NoError(t, err)
			assert.False(t, userConfigDirCalled, "UserConfigDir must not be called on darwin")
			assert.Equal(t, filepath.Join(home, ".config"), got)
		})

		t.Run("non-darwin uses UserConfigDir callback", func(t *testing.T) {
			got, err := defaultAccountsBaseDir("linux", "/unused", func() (string, error) {
				return "/fake/config", nil
			})
			require.NoError(t, err)
			assert.Equal(t, "/fake/config", got)
		})

		t.Run("non-darwin propagates UserConfigDir error", func(t *testing.T) {
			sentinel := errors.New("user config dir error")
			_, err := defaultAccountsBaseDir("linux", "/unused", func() (string, error) {
				return "", sentinel
			})
			require.ErrorIs(t, err, sentinel)
		})
	})

	t.Run("defaultAccountsFilePath", func(t *testing.T) {
		t.Run("joins atlacp segment after UserConfigDir", func(t *testing.T) {
			path, err := defaultAccountsFilePath("linux", "/unused", func() (string, error) {
				return "/fake/config", nil
			})
			require.NoError(t, err)
			assert.Equal(t, filepath.Join("/fake/config", "atlacp", "accounts.json"), path)
		})

		t.Run("propagates UserConfigDir error", func(t *testing.T) {
			sentinel := errors.New("user config dir error")
			_, err := defaultAccountsFilePath("linux", "/unused", func() (string, error) {
				return "", sentinel
			})
			require.ErrorIs(t, err, sentinel)
		})
	})

	t.Run("returns non-empty path ending in atlacp/accounts.json", func(t *testing.T) {
		path, err := DefaultAccountsFilePath()
		require.NoError(t, err)
		require.NotEmpty(t, path)

		suffix := filepath.Join("atlacp", "accounts.json")
		assert.True(t, strings.HasSuffix(path, suffix), "got %q", path)
	})

	t.Run("on macOS uses XDG-style path under HOME", func(t *testing.T) {
		if runtime.GOOS != goosDarwin {
			t.Skip("macOS-only behavior")
		}

		home := t.TempDir()
		t.Setenv("HOME", home)

		path, err := DefaultAccountsFilePath()
		require.NoError(t, err)

		want := filepath.Join(home, ".config", "atlacp", "accounts.json")
		assert.Equal(t, want, path)
	})

	t.Run("on non-macOS uses os.UserConfigDir", func(t *testing.T) {
		if runtime.GOOS == goosDarwin {
			t.Skip("non-macOS behavior")
		}

		base, err := os.UserConfigDir()
		require.NoError(t, err)

		path, err := DefaultAccountsFilePath()
		require.NoError(t, err)

		want := filepath.Join(base, "atlacp", "accounts.json")
		assert.Equal(t, want, path)
	})
}
