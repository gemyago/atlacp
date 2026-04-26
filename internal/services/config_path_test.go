package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtlacpPathResolver(t *testing.T) {
	t.Run("defaultBaseDir", func(t *testing.T) {
		home := t.TempDir()
		r := &AtlacpPathResolver{
			Home: home,
		}
		got := r.defaultBaseDir()
		assert.Equal(t, filepath.Join(home, ".atlacp"), got)
	})

	t.Run("DefaultPath", func(t *testing.T) {
		home := t.TempDir()
		r := &AtlacpPathResolver{
			Home: home,
		}
		path, err := r.DefaultAccountsFilePath()
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(home, ".atlacp", "accounts.json"), path)

		t.Run("NewAtlacpPathResolver uses runtime environment", func(t *testing.T) {
			envHome, ok := os.LookupEnv("HOME")
			require.True(t, ok)
			require.NotEmpty(t, envHome)

			resolvedPath, runtimeErr := NewAtlacpPathResolver().DefaultAccountsFilePath()
			require.NoError(t, runtimeErr)
			assert.Equal(t, filepath.Join(envHome, ".atlacp", "accounts.json"), resolvedPath)
		})
	})

	t.Run("DefaultLogPath", func(t *testing.T) {
		t.Run("returns local filename for go run executable path", func(t *testing.T) {
			home := t.TempDir()
			r := &AtlacpPathResolver{
				Home:           home,
				ExecutablePath: filepath.Join(string(filepath.Separator), "tmp", "go-build1234", "b001", "exe", "bbmd"),
			}
			path, err := r.DefaultLogPath("bbmd.log")
			require.NoError(t, err)
			assert.Equal(t, "bbmd.log", path)
		})

		t.Run("returns HOME/.atlacp/<file> for compiled binary path", func(t *testing.T) {
			home := t.TempDir()
			r := &AtlacpPathResolver{
				Home:           home,
				ExecutablePath: filepath.Join(home, ".atlacp", "bin", "bbmd"),
			}
			path, err := r.DefaultLogPath("bbmd.log")
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(home, ".atlacp", "bbmd.log"), path)
		})

		t.Run("empty log file name errors", func(t *testing.T) {
			r := NewAtlacpPathResolver()
			_, err := r.DefaultLogPath("")
			require.Error(t, err)
			require.ErrorContains(t, err, "log file name is empty")
		})
	})

	t.Run("EnsureParentDirsForFile", func(t *testing.T) {
		t.Run("empty path errors", func(t *testing.T) {
			r := NewAtlacpPathResolver()
			err := r.EnsureParentDirsForFile("")
			require.Error(t, err)
		})

		t.Run("creates missing parents", func(t *testing.T) {
			base := t.TempDir()
			nested := filepath.Join(base, "a", "b", "c")
			filePath := filepath.Join(nested, "accounts.json")
			r := NewAtlacpPathResolver()
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
			st, err := os.Stat(nested)
			require.NoError(t, err)
			assert.True(t, st.IsDir())
		})

		t.Run("idempotent when dir exists", func(t *testing.T) {
			base := t.TempDir()
			filePath := filepath.Join(base, "atlacp", "accounts.json")
			r := NewAtlacpPathResolver()
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
			require.NoError(t, r.EnsureParentDirsForFile(filePath))
		})

		t.Run("returns wrapped error when parent directory cannot be created", func(t *testing.T) {
			base := t.TempDir()
			blockingPath := filepath.Join(base, "not-a-dir")
			require.NoError(t, os.WriteFile(blockingPath, []byte("x"), 0o600))

			filePath := filepath.Join(blockingPath, "accounts.json")
			r := NewAtlacpPathResolver()
			err := r.EnsureParentDirsForFile(filePath)
			require.Error(t, err)
			require.ErrorContains(t, err, "create file parent directory")
		})
	})
}
