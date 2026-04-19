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

func TestAccountsFilePathResolver(t *testing.T) {
	t.Run("defaultBaseDir", func(t *testing.T) {
		t.Run("darwin uses HOME/.config and does not call UserConfigDir", func(t *testing.T) {
			home := t.TempDir()
			userConfigDirCalled := false
			r := &AccountsFilePathResolver{
				GOOS: goosDarwin,
				Home: home,
				UserConfigDir: func() (string, error) {
					userConfigDirCalled = true
					return "", nil
				},
			}
			got, err := r.defaultBaseDir()
			require.NoError(t, err)
			assert.False(t, userConfigDirCalled, "UserConfigDir must not be called on darwin")
			assert.Equal(t, filepath.Join(home, ".config"), got)
		})

		t.Run("non-darwin uses UserConfigDir callback", func(t *testing.T) {
			r := &AccountsFilePathResolver{
				GOOS: "linux",
				Home: "/unused",
				UserConfigDir: func() (string, error) {
					return "/fake/config", nil
				},
			}
			got, err := r.defaultBaseDir()
			require.NoError(t, err)
			assert.Equal(t, "/fake/config", got)
		})

		t.Run("non-darwin propagates UserConfigDir error", func(t *testing.T) {
			sentinel := errors.New("user config dir error")
			r := &AccountsFilePathResolver{
				GOOS: "linux",
				Home: "/unused",
				UserConfigDir: func() (string, error) {
					return "", sentinel
				},
			}
			_, err := r.defaultBaseDir()
			require.ErrorIs(t, err, sentinel)
		})
	})

	t.Run("DefaultPath", func(t *testing.T) {
		t.Run("joins atlacp segment after UserConfigDir", func(t *testing.T) {
			r := &AccountsFilePathResolver{
				GOOS: "linux",
				Home: "/unused",
				UserConfigDir: func() (string, error) {
					return "/fake/config", nil
				},
			}
			path, err := r.DefaultPath()
			require.NoError(t, err)
			assert.Equal(t, filepath.Join("/fake/config", "atlacp", "accounts.json"), path)
		})

		t.Run("propagates UserConfigDir error", func(t *testing.T) {
			sentinel := errors.New("user config dir error")
			r := &AccountsFilePathResolver{
				GOOS: "linux",
				Home: "/unused",
				UserConfigDir: func() (string, error) {
					return "", sentinel
				},
			}
			_, err := r.DefaultPath()
			require.ErrorIs(t, err, sentinel)
		})

		t.Run("NewAccountsFilePathResolver uses runtime environment", func(t *testing.T) {
			t.Run("returns non-empty path ending in atlacp/accounts.json", func(t *testing.T) {
				path, err := NewAccountsFilePathResolver().DefaultPath()
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

				r := &AccountsFilePathResolver{
					GOOS:          goosDarwin,
					Home:          home,
					UserConfigDir: os.UserConfigDir,
				}
				path, err := r.DefaultPath()
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

				path, err := NewAccountsFilePathResolver().DefaultPath()
				require.NoError(t, err)

				want := filepath.Join(base, "atlacp", "accounts.json")
				assert.Equal(t, want, path)
			})
		})
	})

	t.Run("EnsureParentDirsForFile", func(t *testing.T) {
		t.Run("empty path errors", func(t *testing.T) {
			r := NewAccountsFilePathResolver()
			err := r.EnsureParentDirsForFile("")
			require.Error(t, err)
		})

		t.Run("creates missing parents", func(t *testing.T) {
			base := t.TempDir()
			nested := filepath.Join(base, "a", "b", "c")
			filePath := filepath.Join(nested, "accounts.json")
			r := NewAccountsFilePathResolver()
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
			st, err := os.Stat(nested)
			require.NoError(t, err)
			assert.True(t, st.IsDir())
		})

		t.Run("idempotent when dir exists", func(t *testing.T) {
			base := t.TempDir()
			filePath := filepath.Join(base, "atlacp", "accounts.json")
			r := NewAccountsFilePathResolver()
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
		})
	})
}
