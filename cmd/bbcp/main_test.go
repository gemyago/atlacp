package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gemyago/atlacp/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("default logs file path", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		accountsPath := filepath.Join(dir, "accounts.json")

		rootCmd := setupCommands()
		rootCmd.SetArgs([]string{
			"stdio",
			"--noop",
			"--atlassian-accounts-file",
			accountsPath,
		})
		require.NoError(t, rootCmd.Execute())

		expectedLogPath, err := services.NewAtlacpPathResolver().DefaultLogPath("bbcp.log")
		require.NoError(t, err)
		if !filepath.IsAbs(expectedLogPath) {
			expectedLogPath = filepath.Join(dir, expectedLogPath)
		}

		_, err = os.Stat(expectedLogPath)
		require.NoError(t, err)
	})

	t.Run("http", func(t *testing.T) {
		t.Run("should initialize app", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{
				"http",
				"--noop",
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("should initialize app with default accounts file path when flag omitted", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{
				"http",
				"--noop",
				"--logs-file",
				"../../test.log",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("should fail if bad log level", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{
				"http",
				"--noop",
				"-l",
				faker.Word(),
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			assert.Error(t, rootCmd.Execute())
		})
		t.Run("should fail if unexpected env", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{
				"http",
				"--noop",
				"-e",
				faker.Word(),
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			gotErr := rootCmd.Execute()
			assert.ErrorContains(t, gotErr, "failed to read config")
		})
	})
	t.Run("stdio", func(t *testing.T) {
		t.Run("should initialize app", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{
				"stdio",
				"--noop",
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("should initialize app with default accounts file path when flag omitted", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{
				"stdio",
				"--noop",
				"--logs-file",
				"../../test.log",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("should fail if bad log level", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{
				"stdio",
				"--noop",
				"-l",
				faker.Word(),
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			assert.Error(t, rootCmd.Execute())
		})
		t.Run("should fail if unexpected env", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetArgs([]string{
				"stdio",
				"--noop",
				"-e",
				faker.Word(),
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../quick-start/atlassian-accounts-stub.json",
			})
			gotErr := rootCmd.Execute()
			assert.ErrorContains(t, gotErr, "failed to read config")
		})
	})
}
