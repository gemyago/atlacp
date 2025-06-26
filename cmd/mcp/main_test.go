package main

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		t.Run("should initialize app", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SetArgs([]string{
				"http",
				"--noop",
				"--logs-file",
				"../../test.log",
				"--atlassian-accounts-file",
				"../../examples/atlassian-accounts-stub.json",
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
				"../../examples/atlassian-accounts-stub.json",
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
				"../../examples/atlassian-accounts-stub.json",
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
				"../../examples/atlassian-accounts-stub.json",
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
				"../../examples/atlassian-accounts-stub.json",
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
				"../../examples/atlassian-accounts-stub.json",
			})
			gotErr := rootCmd.Execute()
			assert.ErrorContains(t, gotErr, "failed to read config")
		})
	})
}
