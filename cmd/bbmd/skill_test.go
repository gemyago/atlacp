package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorWriter struct{}

func (errorWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestSkillCommand(t *testing.T) {
	t.Parallel()

	t.Run("renders skill guide from cobra command metadata", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		root := setupCommands()
		root.SetOut(&out)
		root.SetErr(io.Discard)
		logFile := filepath.Join(t.TempDir(), "bbmd.log")
		root.SetArgs([]string{
			"skill",
			"--logs-file", logFile,
		})

		require.NoError(t, root.Execute())
		guide := out.String()

		assert.Contains(t, guide, "# bbmd CLI Skill")
		assert.Contains(t, guide, "## Global Flags")
		assert.Contains(t, guide, "### `bbmd pr`")
		assert.Contains(t, guide, "## `bbmd`")
		assert.Contains(t, guide, "bbmd pr add-comment")
		assert.Contains(t, guide, "bbmd auth status --account")
		assert.Contains(t, guide, "#### `bbmd auth status`")
		assert.Contains(t, guide, "Examples:")
		assert.NotContains(t, guide, "### `bbmd completion`")
		assert.NotContains(t, guide, "bbmd auth add")
		assert.NotContains(t, guide, "bbmd auth remove")
		assert.NotContains(t, guide, "bbmd auth set-default")
		assert.NotContains(t, guide, "### `bbmd skill`")
	})

	t.Run("skill command skips bootstrap while operational commands still bootstrap", func(t *testing.T) {
		t.Parallel()

		restrictedDir := t.TempDir()
		require.NoError(t, os.Chmod(restrictedDir, 0o500))
		accountsPath := filepath.Join(restrictedDir, "nested", "accounts.json")
		skillOut := bytes.Buffer{}

		skillCmd := setupCommands()
		skillCmd.SetOut(&skillOut)
		skillCmd.SetErr(io.Discard)
		skillCmd.SetArgs([]string{
			"skill",
			"--logs-file", filepath.Join(t.TempDir(), "bbmd.log"),
			"--atlassian-accounts-file", accountsPath,
		})
		require.NoError(t, skillCmd.Execute())
		assert.NotEmpty(t, strings.TrimSpace(skillOut.String()))

		authStatus := setupCommands()
		authStatus.SetOut(io.Discard)
		authStatus.SetErr(io.Discard)
		authStatus.SetArgs([]string{
			"auth", "status",
			"--logs-file", filepath.Join(t.TempDir(), "bbmd.log"),
			"--atlassian-accounts-file", accountsPath,
		})
		err := authStatus.Execute()
		require.Error(t, err)
		assert.ErrorContains(t, err, "ensure atlassian accounts file parent directories")
	})
}

func TestWriteSkillGuide(t *testing.T) {
	t.Parallel()

	t.Run("returns error on nil root", func(t *testing.T) {
		t.Parallel()

		err := writeSkillGuide(nil, io.Discard)
		require.Error(t, err)
		assert.ErrorContains(t, err, "missing root command")
	})

	t.Run("returns error on write failure", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: "bbmd"}
		err := writeSkillGuide(cmd, errorWriter{})
		require.Error(t, err)
		require.ErrorContains(t, err, "write skill output")
		require.ErrorContains(t, err, "write failed")
	})
}

func TestSkillCommandMetadataHelpers(t *testing.T) {
	t.Parallel()

	t.Run("renders empty usage when command use line is empty", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: ""}
		assert.Empty(t, renderCommandUsage(cmd))
	})

	t.Run("writes nothing when flagset has no renderable flags", func(t *testing.T) {
		t.Parallel()

		var out strings.Builder
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		_ = fs.BoolP("help", "h", false, "help")
		fs.Lookup("help").Hidden = true

		renderSkillFlags(&out, "Flags", fs)
		assert.Empty(t, strings.TrimSpace(out.String()))
	})

	t.Run("renders rendered flags and shorthand formatting", func(t *testing.T) {
		t.Parallel()

		var out strings.Builder
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		_ = fs.Int("count", 0, "count")
		fs.Lookup("count").Shorthand = "c"
		fs.Lookup("count").NoOptDefVal = "optional"
		_ = fs.BoolP("help", "h", false, "help")

		renderSkillFlags(&out, "Flags", fs)
		outStr := out.String()
		assert.Contains(t, outStr, "### Flags")
		assert.Contains(t, outStr, "`--count`, `-c`: count (optional)")
		assert.NotContains(t, outStr, "`--help`")
	})

	t.Run("renders command descriptions and nested examples", func(t *testing.T) {
		t.Parallel()

		var out strings.Builder
		root := &cobra.Command{
			Use:   "bbmd",
			Short: "short",
		}
		child := &cobra.Command{
			Use:     "child",
			Short:   "short child",
			Example: "bbmd child",
		}
		root.AddCommand(child)
		renderSkillCommand(&out, root, skillRootSubcommandLevel)
		result := out.String()
		assert.Contains(t, result, "### `bbmd child`")
		assert.Contains(t, result, "short child")
		assert.Contains(t, result, "Usage: `bbmd child`")
		assert.Contains(t, result, "Examples:\n\n```bash\nbbmd child")
	})

	t.Run("respects skill render annotations", func(t *testing.T) {
		t.Parallel()

		root := &cobra.Command{Use: "root"}
		auth := &cobra.Command{
			Use: "auth",
			Annotations: map[string]string{
				commandExcludeFromSkillAnnotation: "true",
			},
		}
		authStatus := &cobra.Command{
			Use: "status",
			Annotations: map[string]string{
				commandIncludeInSkillAnnotation: "true",
			},
		}
		hidden := &cobra.Command{
			Use:    "hidden",
			Hidden: true,
		}
		excluded := &cobra.Command{
			Use: "exclude",
			Annotations: map[string]string{
				commandExcludeFromSkillAnnotation: "true",
			},
		}
		included := &cobra.Command{
			Use: "include",
			Annotations: map[string]string{
				commandIncludeInSkillAnnotation: "true",
			},
		}
		helpCmd := &cobra.Command{Use: "help"}
		completion := &cobra.Command{Use: "completion"}
		auth.AddCommand(authStatus)
		root.AddCommand(hidden, excluded, included, helpCmd, completion, auth)

		assert.False(t, shouldRenderInSkill(hidden))
		assert.False(t, shouldRenderInSkill(excluded))
		assert.True(t, shouldRenderInSkill(included))
		assert.False(t, shouldRenderInSkill(helpCmd))
		assert.False(t, shouldRenderInSkill(completion))
		assert.False(t, shouldRenderInSkill(auth))
		assert.True(t, shouldRenderInSkill(authStatus))
	})

	t.Run("checks bootstrap skip annotation on ancestors", func(t *testing.T) {
		t.Parallel()

		root := &cobra.Command{Use: "root"}
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{Use: "child"}
		root.AddCommand(parent)
		parent.AddCommand(child)

		assert.False(t, shouldSkipBootstrap(child))
		parent.Annotations = map[string]string{
			commandSkipBootstrapAnnotation: "true",
		}
		assert.True(t, shouldSkipBootstrap(child))
		assert.True(t, shouldSkipBootstrap(parent))
		assert.True(t, hasCommandAnnotation(parent, commandSkipBootstrapAnnotation))
		assert.False(t, hasCommandAnnotation(&cobra.Command{Use: "cmd"}, commandSkipBootstrapAnnotation))
	})
}
