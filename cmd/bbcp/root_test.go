package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gemyago/atlacp/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareMCPLogsOutputFile(t *testing.T) {
	newCmd := func(t *testing.T) *cobra.Command {
		t.Helper()
		cmd := &cobra.Command{Use: "bbcp"}
		cmd.Flags().String("logs-file", "", "")
		return cmd
	}

	t.Run("uses resolver default and creates parent directory", func(t *testing.T) {
		home := t.TempDir()
		cmd := newCmd(t)
		resolver := &services.AtlacpPathResolver{
			Home:           home,
			ExecutablePath: filepath.Join(home, ".atlacp", "bin", "bbcp"),
		}

		logsOutputFile := ""
		err := prepareMCPLogsOutputFile(cmd, resolver, &logsOutputFile)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(home, ".atlacp", "bbcp.log"), logsOutputFile)

		_, statErr := os.Stat(filepath.Join(home, ".atlacp"))
		require.NoError(t, statErr)
	})

	t.Run("explicit empty logs path keeps stdout mode", func(t *testing.T) {
		cmd := newCmd(t)
		require.NoError(t, cmd.Flags().Set("logs-file", ""))

		logsOutputFile := ""
		resolver := &services.AtlacpPathResolver{
			Home:           t.TempDir(),
			ExecutablePath: filepath.Join(string(filepath.Separator), "tmp", "go-build1234", "b001", "exe", "bbcp"),
		}

		err := prepareMCPLogsOutputFile(cmd, resolver, &logsOutputFile)
		require.NoError(t, err)
		assert.Empty(t, logsOutputFile)
	})

	t.Run("returns error when parent directory cannot be created", func(t *testing.T) {
		base := t.TempDir()
		blockingPath := filepath.Join(base, "not-a-dir")
		require.NoError(t, os.WriteFile(blockingPath, []byte("x"), 0o600))

		cmd := newCmd(t)
		logsOutputFile := filepath.Join(blockingPath, "bbcp.log")
		require.NoError(t, cmd.Flags().Set("logs-file", logsOutputFile))

		resolver := &services.AtlacpPathResolver{
			Home:           base,
			ExecutablePath: filepath.Join(base, ".atlacp", "bin", "bbcp"),
		}
		err := prepareMCPLogsOutputFile(cmd, resolver, &logsOutputFile)
		require.Error(t, err)
		require.ErrorContains(t, err, "ensure logs output file parent directories")
	})
}

func TestPrepareMCPAccountsFilePath(t *testing.T) {
	t.Run("returns error when parent directory cannot be created", func(t *testing.T) {
		base := t.TempDir()
		blockingPath := filepath.Join(base, "not-a-dir")
		require.NoError(t, os.WriteFile(blockingPath, []byte("x"), 0o600))

		cfg := viper.New()
		resolver := &services.AtlacpPathResolver{
			Home:           base,
			ExecutablePath: filepath.Join(base, ".atlacp", "bin", "bbcp"),
		}
		err := prepareMCPAccountsFilePath(cfg, resolver, filepath.Join(blockingPath, "accounts.json"))
		require.Error(t, err)
		require.ErrorContains(t, err, "ensure atlassian accounts file parent directories")
	})
}
