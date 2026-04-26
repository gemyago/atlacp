package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	commandSkipBootstrapAnnotation    = "bbmd.skipBootstrap"
	commandIncludeInSkillAnnotation   = "bbmd.includeInSkill"
	commandExcludeFromSkillAnnotation = "bbmd.excludeFromSkill"
	commandHelpName                   = "help"
	commandCompletionName             = "completion"
)

const skillRootSubcommandLevel = 2

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Generate a markdown skill guide for bbmd",
		Long: "Generate a markdown guide based on the current bbmd command structure. " +
			"Use this output for agents and operators that need a single, stable reference.",
		Example: `bbmd skill`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSkill(cmd)
		},
		Args: cobra.NoArgs,
	}
	cmd.Annotations = map[string]string{
		commandSkipBootstrapAnnotation:    "true",
		commandExcludeFromSkillAnnotation: "true",
	}
	return cmd
}

func runSkill(cmd *cobra.Command) error {
	return writeSkillGuide(cmd.Root(), cmd.OutOrStdout())
}

func writeSkillGuide(root *cobra.Command, out io.Writer) error {
	if root == nil {
		return errors.New("missing root command")
	}

	var b strings.Builder
	b.WriteString("# bbmd CLI Skill\n\n")
	b.WriteString("The following instructions are generated from the live Cobra command tree.\n\n")

	_, _ = b.WriteString("## `" + root.CommandPath() + "`\n\n")
	_, _ = b.WriteString(commandDescription(root) + "\n\n")
	_, _ = b.WriteString(renderCommandUsage(root) + "\n")

	renderSkillFlags(&b, "Global Flags", root.PersistentFlags())

	renderSkillCommand(&b, root, skillRootSubcommandLevel)

	_, err := out.Write([]byte(b.String()))
	if err != nil {
		return fmt.Errorf("write skill output: %w", err)
	}
	return nil
}

func commandDescription(cmd *cobra.Command) string {
	desc := strings.TrimSpace(cmd.Long)
	if desc == "" {
		desc = strings.TrimSpace(cmd.Short)
	}
	return desc
}

func renderCommandUsage(cmd *cobra.Command) string {
	usage := strings.TrimSpace(cmd.UseLine())
	if usage == "" {
		return ""
	}
	return fmt.Sprintf("Usage: `%s`", usage)
}

func renderSkillCommand(out *strings.Builder, cmd *cobra.Command, level int) {
	for _, child := range cmd.Commands() {
		if shouldRenderInSkill(child) {
			_, _ = out.WriteString(strings.Repeat("#", level+1))
			_, _ = out.WriteString(" `" + child.CommandPath() + "`\n\n")
			_, _ = out.WriteString(commandDescription(child) + "\n\n")
			_, _ = out.WriteString(renderCommandUsage(child) + "\n\n")

			renderSkillFlags(out, "Flags", child.NonInheritedFlags())

			example := strings.TrimSpace(child.Example)
			if example != "" {
				_, _ = out.WriteString("Examples:\n\n")
				_, _ = out.WriteString("```bash\n")
				_, _ = out.WriteString(example + "\n")
				_, _ = out.WriteString("```\n\n")
			}
		}
		renderSkillCommand(out, child, level+1)
	}
}

func renderSkillFlags(out *strings.Builder, heading string, fs *pflag.FlagSet) {
	if fs == nil || !containsRenderableFlags(fs) {
		return
	}
	_, _ = out.WriteString("### " + heading + "\n\n")
	fs.VisitAll(func(flag *pflag.Flag) {
		if !shouldRenderFlag(flag) {
			return
		}
		_, _ = out.WriteString("- " + renderFlag(flag) + "\n")
	})
	_, _ = out.WriteString("\n")
}

func renderFlag(flag *pflag.Flag) string {
	var name string
	if flag.Shorthand != "" {
		name = fmt.Sprintf("`--%s`, `-%s`", flag.Name, flag.Shorthand)
	} else {
		name = fmt.Sprintf("`--%s`", flag.Name)
	}

	desc := strings.TrimSpace(flag.Usage)
	if flag.NoOptDefVal != "" {
		desc = strings.TrimSpace(desc + " (optional)")
	}
	return name + ": " + desc
}

func containsRenderableFlags(fs *pflag.FlagSet) bool {
	renderable := false
	fs.VisitAll(func(flag *pflag.Flag) {
		renderable = renderable || shouldRenderFlag(flag)
	})
	return renderable
}

func shouldRenderFlag(flag *pflag.Flag) bool {
	return !flag.Hidden && flag.Name != commandHelpName
}

func shouldRenderInSkill(cmd *cobra.Command) bool {
	if cmd == nil || cmd.Hidden {
		return false
	}

	if cmd.Name() == commandHelpName || cmd.Name() == commandCompletionName {
		return false
	}

	if hasCommandAnnotation(cmd, commandIncludeInSkillAnnotation) {
		return true
	}
	for current := cmd; current != nil; current = current.Parent() {
		if current.Name() == commandHelpName || current.Name() == commandCompletionName {
			return false
		}
		if hasCommandAnnotation(current, commandExcludeFromSkillAnnotation) {
			return false
		}
	}
	return true
}

func shouldSkipBootstrap(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if hasCommandAnnotation(current, commandSkipBootstrapAnnotation) {
			return true
		}
	}
	return false
}

func hasCommandAnnotation(cmd *cobra.Command, annotation string) bool {
	return lo.If(
		cmd != nil && cmd.Annotations != nil,
		cmd.Annotations[annotation] == "true",
	).Else(false)
}
