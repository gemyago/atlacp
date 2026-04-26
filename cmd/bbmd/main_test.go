package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gemyago/atlacp/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// faultyWriter fails every Write; used to cover auth status stdout error handling.
type faultyWriter struct{}

func (faultyWriter) Write([]byte) (int, error) {
	return 0, errors.New("injected write error")
}

func TestBBMD(t *testing.T) {
	t.Run("default logs file path", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		accountsPath := filepath.Join(dir, "accounts.json")

		rootCmd := setupCommands()
		rootCmd.SetArgs([]string{
			"auth", "status",
			"--noop",
			"--atlassian-accounts-file", accountsPath,
		})
		require.NoError(t, rootCmd.Execute())

		expectedLogPath, err := services.NewAtlacpPathResolver().DefaultLogPath("bbmd.log")
		require.NoError(t, err)
		if !filepath.IsAbs(expectedLogPath) {
			expectedLogPath = filepath.Join(dir, expectedLogPath)
		}

		_, err = os.Stat(expectedLogPath)
		require.NoError(t, err)
	})

	t.Run("auth", func(t *testing.T) {
		t.Run("status noop exercises DI", func(t *testing.T) {
			rootCmd := setupCommands()
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"auth", "status",
				"--noop",
				"--logs-file", logFile,
			})
			require.NoError(t, rootCmd.Execute())
		})
		t.Run("status without noop prints empty array", func(t *testing.T) {
			var stdout bytes.Buffer
			rootCmd := setupCommands()
			rootCmd.SetOut(&stdout)
			rootCmd.SetErr(io.Discard)
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			rootCmd.SetArgs([]string{
				"auth", "status",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, rootCmd.Execute())
			assert.JSONEq(t, "[]", strings.TrimSpace(stdout.String()))
		})
		t.Run("add first account without --default persists and marks default", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			base := []string{"--logs-file", logFile, "--atlassian-accounts-file", accountsPath}

			add := setupCommands()
			add.SetOut(io.Discard)
			add.SetErr(io.Discard)
			add.SetArgs(append([]string{
				"auth", "add",
				"--name", "solo",
				"--token-value", "zzzzzzzzzzzzzzzz",
			}, base...))
			require.NoError(t, add.Execute())

			var st bytes.Buffer
			stCmd := setupCommands()
			stCmd.SetOut(&st)
			stCmd.SetErr(io.Discard)
			stCmd.SetArgs(append([]string{"auth", "status"}, base...))
			require.NoError(t, stCmd.Execute())

			var rows []map[string]any
			require.NoError(t, json.Unmarshal(bytes.TrimSpace(st.Bytes()), &rows))
			require.Len(t, rows, 1)
			assert.Equal(t, "solo", rows[0]["name"])
			assert.Equal(t, true, rows[0]["default"])
		})
		t.Run("status redacts jira token from file", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			accountsJSON := `{
  "accounts": [
    {
      "name": "jin",
      "default": true,
      "jira": { "type": "Bearer", "value": "jjjjjjjjjjjjjjjj" }
    }
  ]
}`
			require.NoError(t, os.WriteFile(accountsPath, []byte(accountsJSON), 0o600))

			var stdout bytes.Buffer
			rootCmd := setupCommands()
			rootCmd.SetOut(&stdout)
			rootCmd.SetErr(io.Discard)
			rootCmd.SetArgs([]string{
				"auth", "status",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, rootCmd.Execute())

			var rows []map[string]any
			require.NoError(t, json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &rows))
			require.Len(t, rows, 1)
			jr := rows[0]["jira"].(map[string]any)
			assert.Equal(t, "jjjj***jjjj", jr["value"])
		})
		t.Run("status returns error when stdout write fails", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			rootCmd := setupCommands()
			rootCmd.SetOut(faultyWriter{})
			rootCmd.SetErr(io.Discard)
			rootCmd.SetArgs([]string{
				"auth", "status",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := rootCmd.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "write JSON")
		})
		t.Run("mutating commands persist and status redacts", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			base := []string{"--logs-file", logFile, "--atlassian-accounts-file", accountsPath}

			runExec := func(args ...string) {
				t.Helper()
				c := setupCommands()
				c.SetOut(io.Discard)
				c.SetErr(io.Discard)
				c.SetArgs(append(args, base...))
				require.NoError(t, c.Execute())
			}

			runExec(
				"auth", "add",
				"--name", "a1",
				"--default",
				"--token-value", "aaaaaaaaaaaaaaaa",
			)
			runExec(
				"auth", "add",
				"--name", "b2",
				"--token-value", "bbbbbbbbbbbbbbbb",
			)

			var st bytes.Buffer
			stCmd := setupCommands()
			stCmd.SetOut(&st)
			stCmd.SetErr(io.Discard)
			stCmd.SetArgs(append([]string{"auth", "status"}, base...))
			require.NoError(t, stCmd.Execute())

			var rows []map[string]any
			require.NoError(t, json.Unmarshal(bytes.TrimSpace(st.Bytes()), &rows))
			require.Len(t, rows, 2)

			runExec("auth", "set-default", "--name", "b2")
			runExec("auth", "remove", "--name", "a1")

			st2 := bytes.Buffer{}
			stCmd2 := setupCommands()
			stCmd2.SetOut(&st2)
			stCmd2.SetErr(io.Discard)
			stCmd2.SetArgs(append([]string{"auth", "status"}, base...))
			require.NoError(t, stCmd2.Execute())
			require.NoError(t, json.Unmarshal(bytes.TrimSpace(st2.Bytes()), &rows))
			require.Len(t, rows, 1)
			bb := rows[0]["bitbucket"].(map[string]any)
			assert.Equal(t, "bbbb***bbbb", bb["value"])
		})
		t.Run("add remove set-default noop", func(t *testing.T) {
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			for _, args := range [][]string{
				{"auth", "add", "--noop", "--name", "n", "--token-value", "zzzzzzzzzzzzzzzz"},
				{"auth", "remove", "--noop", "--name", "n"},
				{"auth", "set-default", "--noop", "--name", "n"},
			} {
				rootCmd := setupCommands()
				rootCmd.SetOut(io.Discard)
				rootCmd.SetErr(io.Discard)
				rootCmd.SetArgs(append(args, "--logs-file", logFile))
				require.NoError(t, rootCmd.Execute())
			}
		})
		t.Run("add returns error when save fails", func(t *testing.T) {
			roDir := filepath.Join(t.TempDir(), "ro")
			require.NoError(t, os.Mkdir(roDir, 0o555))
			accountsPath := filepath.Join(roDir, "accounts.json")
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)
			rootCmd.SetArgs([]string{
				"auth", "add",
				"--name", "n1",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := rootCmd.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "save accounts")
		})
		t.Run("remove returns error when save fails", func(t *testing.T) {
			dir := t.TempDir()
			dataDir := filepath.Join(dir, "data")
			require.NoError(t, os.MkdirAll(dataDir, 0o755))
			accountsPath := filepath.Join(dataDir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd-test.log")

			add1 := setupCommands()
			add1.SetOut(io.Discard)
			add1.SetErr(io.Discard)
			add1.SetArgs([]string{
				"auth", "add",
				"--name", "keep",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add1.Execute())

			add2 := setupCommands()
			add2.SetOut(io.Discard)
			add2.SetErr(io.Discard)
			add2.SetArgs([]string{
				"auth", "add",
				"--name", "drop",
				"--token-value", "yyyyyyyyyyyyyyyy",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add2.Execute())

			require.NoError(t, os.Chmod(dataDir, 0o555))
			t.Cleanup(func() { _ = os.Chmod(dataDir, 0o755) })

			rm := setupCommands()
			rm.SilenceErrors = true
			rm.SilenceUsage = true
			rm.SetOut(io.Discard)
			rm.SetErr(io.Discard)
			rm.SetArgs([]string{
				"auth", "remove",
				"--name", "drop",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := rm.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "save accounts")
		})
		t.Run("set-default returns error when save fails", func(t *testing.T) {
			dir := t.TempDir()
			dataDir := filepath.Join(dir, "data")
			require.NoError(t, os.MkdirAll(dataDir, 0o755))
			accountsPath := filepath.Join(dataDir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd-test.log")

			add := setupCommands()
			add.SetOut(io.Discard)
			add.SetErr(io.Discard)
			add.SetArgs([]string{
				"auth", "add",
				"--name", "only",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add.Execute())

			require.NoError(t, os.Chmod(dataDir, 0o555))
			t.Cleanup(func() { _ = os.Chmod(dataDir, 0o755) })

			sd := setupCommands()
			sd.SilenceErrors = true
			sd.SilenceUsage = true
			sd.SetOut(io.Discard)
			sd.SetErr(io.Discard)
			sd.SetArgs([]string{
				"auth", "set-default",
				"--name", "only",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := sd.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "save accounts")
		})
		t.Run("remove returns error when account not found", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			add := setupCommands()
			add.SetOut(io.Discard)
			add.SetErr(io.Discard)
			add.SetArgs([]string{
				"auth", "add",
				"--name", "only",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add.Execute())

			rm := setupCommands()
			rm.SilenceErrors = true
			rm.SetOut(io.Discard)
			rm.SetErr(io.Discard)
			rm.SetArgs([]string{
				"auth", "remove",
				"--name", "missing",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := rm.Execute()
			require.Error(t, err)
		})
		t.Run("set-default returns error when account not found", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			add := setupCommands()
			add.SetOut(io.Discard)
			add.SetErr(io.Discard)
			add.SetArgs([]string{
				"auth", "add",
				"--name", "only",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add.Execute())

			sd := setupCommands()
			sd.SilenceErrors = true
			sd.SetOut(io.Discard)
			sd.SetErr(io.Discard)
			sd.SetArgs([]string{
				"auth", "set-default",
				"--name", "missing",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := sd.Execute()
			require.Error(t, err)
		})
		t.Run("add returns error when upsert validation fails", func(t *testing.T) {
			dir := t.TempDir()
			accountsPath := filepath.Join(dir, "accounts.json")
			logFile := filepath.Join(dir, "bbmd.log")
			add1 := setupCommands()
			add1.SetOut(io.Discard)
			add1.SetErr(io.Discard)
			add1.SetArgs([]string{
				"auth", "add",
				"--name", "a1",
				"--default",
				"--token-value", "zzzzzzzzzzzzzzzz",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			require.NoError(t, add1.Execute())

			add2 := setupCommands()
			add2.SilenceErrors = true
			add2.SetOut(io.Discard)
			add2.SetErr(io.Discard)
			add2.SetArgs([]string{
				"auth", "add",
				"--name", "a2",
				"--default",
				"--token-value", "yyyyyyyyyyyyyyyy",
				"--logs-file", logFile,
				"--atlassian-accounts-file", accountsPath,
			})
			err := add2.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "invalid accounts configuration")
		})
	})
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
		t.Run("all subcommands noop exercise DI", func(t *testing.T) {
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			base := []string{"--noop", "--logs-file", logFile}
			repoOnly := []string{"--repo-owner", "o", "--repo-name", "n"}
			repoPR := []string{"--repo-owner", "o", "--repo-name", "n", "--pr-id", "1"}
			cases := []struct {
				name string
				args []string
			}{
				{
					name: "create",
					args: append(
						[]string{"pr", "create", "--title", "t", "--source-branch", "s", "--target-branch", "d"},
						repoOnly...,
					),
				},
				{name: "read", args: append([]string{"pr", "read"}, repoPR...)},
				{
					name: "update",
					args: append([]string{
						"pr", "update",
						"--title", "x",
						"--description", "d",
						"--draft",
					}, repoPR...),
				},
				{name: "approve", args: append([]string{"pr", "approve"}, repoPR...)},
				{name: "request-changes", args: append([]string{"pr", "request-changes"}, repoPR...)},
				{name: "merge", args: append([]string{"pr", "merge"}, repoPR...)},
				{name: "list-tasks", args: append([]string{"pr", "list-tasks"}, repoPR...)},
				{
					name: "create-task",
					args: append([]string{"pr", "create-task", "--content", "c", "--comment-id", "9"}, repoPR...),
				},
				{
					name: "update-task",
					args: append([]string{"pr", "update-task", "--task-id", "1", "--content", "u"}, repoPR...),
				},
				{name: "diffstat", args: append([]string{"pr", "diffstat"}, repoPR...)},
				{
					name: "diff",
					args: append([]string{"pr", "diff", "--path", "a.go", "--context", "2"}, repoPR...),
				},
				{
					name: "add-comment",
					args: append(
						[]string{
							"pr", "add-comment",
							"--content", "hi",
							"--line-from", "10",
							"--line-to", "10",
							"--parent-comment-id", "9",
						},
						repoPR...,
					),
				},
				{
					name: "add-comment range",
					args: append(
						[]string{
							"pr", "add-comment",
							"--content", "hi",
							"--file-path", "a.go",
							"--line-from", "10",
							"--line-to", "12",
							"--parent-comment-id", "9",
						},
						repoPR...,
					),
				},
				{
					name: "list-comments",
					args: append([]string{
						"pr", "list-comments",
						"--include-resolved",
						"--page", "2",
						"--pagelen", "5",
					}, repoPR...),
				},
				{
					name: "resolve-comment",
					args: append([]string{"pr", "resolve-comment", "--comment-id", "1"}, repoPR...),
				},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					rootCmd := setupCommands()
					rootCmd.SetOut(io.Discard)
					rootCmd.SetErr(io.Discard)
					rootCmd.SetArgs(append(tc.args, base...))
					require.NoError(t, rootCmd.Execute())
				})
			}
		})
		t.Run("list-comments help documents pagination flags", func(t *testing.T) {
			var buf bytes.Buffer
			rootCmd := setupCommands()
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(io.Discard)
			rootCmd.SetArgs([]string{"pr", "list-comments", "--help"})
			require.NoError(t, rootCmd.Execute())
			out := buf.String()
			assert.Contains(t, out, "--page")
			assert.Contains(t, out, "--pagelen")
		})
		t.Run("list-comments rejects negative page", func(t *testing.T) {
			rootCmd := setupCommands()
			rootCmd.SilenceErrors = true
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"pr", "list-comments",
				"--noop",
				"--logs-file", logFile,
				"--repo-owner", "o",
				"--repo-name", "n",
				"--pr-id", "1",
				"--page", "-1",
			})
			err := rootCmd.Execute()
			require.Error(t, err)
			assert.ErrorContains(t, err, "invalid --page")
		})
	})
	t.Run("file", func(t *testing.T) {
		t.Run("content noop exercises DI", func(t *testing.T) {
			rootCmd := setupCommands()
			logFile := filepath.Join(t.TempDir(), "bbmd-test.log")
			rootCmd.SetArgs([]string{
				"file", "content",
				"--noop",
				"--logs-file", logFile,
				"--repo-owner", "dummy-owner",
				"--repo-name", "dummy-repo",
				"--commit", "deadbeef",
				"--path", "README.md",
			})
			require.NoError(t, rootCmd.Execute())
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
