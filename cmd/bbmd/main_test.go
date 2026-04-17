package main

import (
	"path/filepath"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBBMD(t *testing.T) {
	t.Run("pr", func(t *testing.T) {
		t.Run("read noop exercises DI", func(t *testing.T) {
			rootCmd := setupCommands()
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "read",
				"--noop",
				"--logs-file", logFile,
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--pr-id", "1",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("read noop with explicit accounts file", func(t *testing.T) {
			rootCmd := setupCommands()
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "read",
				"--noop",
				"--logs-file", logFile,
				"--atlassian-accounts-file", "../../quick-start/atlassian-accounts-stub.json",
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--pr-id", "1",
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("read without noop returns stub error", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "read",
				"--logs-file", logFile,
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--pr-id", "1",
			})
			gotErr := rootCmd.Execute()
			require.Error(t, gotErr)
			assert.ErrorContains(t, gotErr, "not yet implemented")
		})
	})
	t.Run("root", func(t *testing.T) {
		t.Run("should fail if bad log level", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "read",
				"--noop",
				"-l", faker.Word(),
				"--logs-file", logFile,
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--pr-id", "1",
			})
			assert.Error(t, rootCmd.Execute())
		})
		t.Run("should fail if unexpected env", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "read",
				"--noop",
				"-e", faker.Word(),
				"--logs-file", logFile,
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--pr-id", "1",
			})
			gotErr := rootCmd.Execute()
			assert.ErrorContains(t, gotErr, "failed to read config")
		})
	})
}
