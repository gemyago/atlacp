package services

import (
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
		t.Run("returns HOME/.atlacp on darwin", func(t *testing.T) {
			home := t.TempDir()
			r := &AccountsFilePathResolver{
				Home: home,
			}
			got := r.defaultBaseDir()
			assert.Equal(t, filepath.Join(home, ".atlacp"), got)
		})

		t.Run("returns HOME/.atlacp on non-darwin", func(t *testing.T) {
			home := t.TempDir()
			r := &AccountsFilePathResolver{
				Home: home,
			}
			got := r.defaultBaseDir()
			assert.Equal(t, filepath.Join(home, ".atlacp"), got)
		})
	})

	t.Run("DefaultPath", func(t *testing.T) {
		t.Run("returns HOME/.atlacp/accounts.json on linux", func(t *testing.T) {
			home := t.TempDir()
			r := &AccountsFilePathResolver{
				Home: home,
			}
			path, err := r.DefaultPath()
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(home, ".atlacp", "accounts.json"), path)
		})

		t.Run("returns HOME/.atlacp/accounts.json on darwin", func(t *testing.T) {
			home := t.TempDir()
			r := &AccountsFilePathResolver{
				Home: home,
			}
			path, err := r.DefaultPath()
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(home, ".atlacp", "accounts.json"), path)
		})

		t.Run("NewAccountsFilePathResolver uses runtime environment", func(t *testing.T) {
			t.Run("returns non-empty path ending in atlacp/accounts.json", func(t *testing.T) {
				path, err := NewAccountsFilePathResolver().DefaultPath()
				require.NoError(t, err)
				require.NotEmpty(t, path)

				suffix := filepath.Join("atlacp", "accounts.json")
				assert.True(t, strings.HasSuffix(path, suffix), "got %q", path)
			})

			t.Run("on macOS uses ~/.atlacp path under HOME", func(t *testing.T) {
				if runtime.GOOS != "darwin" {
					t.Skip("macOS-only behavior")
				}

				home := t.TempDir()
				t.Setenv("HOME", home)

				r := &AccountsFilePathResolver{
					Home: home,
				}
				path, err := r.DefaultPath()
				require.NoError(t, err)

				want := filepath.Join(home, ".atlacp", "accounts.json")
				assert.Equal(t, want, path)
			})

			t.Run("on non-macOS uses ~/.atlacp under HOME", func(t *testing.T) {
				if runtime.GOOS == "darwin" {
					t.Skip("non-macOS behavior")
				}

				home, ok := os.LookupEnv("HOME")
				require.True(t, ok)
				require.NotEmpty(t, home)

				path, err := NewAccountsFilePathResolver().DefaultPath()
				require.NoError(t, err)

				want := filepath.Join(home, ".atlacp", "accounts.json")
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

		t.Run("returns wrapped error when parent directory cannot be created", func(t *testing.T) {
			base := t.TempDir()
			blockingPath := filepath.Join(base, "not-a-dir")
			require.NoError(t, os.WriteFile(blockingPath, []byte("x"), 0o600))

			filePath := filepath.Join(blockingPath, "accounts.json")
			r := NewAccountsFilePathResolver()
			err := r.EnsureParentDirsForFile(filePath)
			require.Error(t, err)
			require.ErrorContains(t, err, "create accounts file parent directory")
		})
	})
}
