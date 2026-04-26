package main

import (
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type commandParitySpec struct {
	CommandPath []string
	Params      any
	FieldToFlag map[string]string
	Exceptions  map[string]string
}

func TestCommandParityWithBitbucketAppParams(t *testing.T) {
	t.Parallel()

	specs := []commandParitySpec{
		{
			CommandPath: []string{"file", "content"},
			Params:      app.BitbucketGetFileContentParams{},
			FieldToFlag: map[string]string{
				"AccountName": "account",
				"RepoOwner":   "repo-owner",
				"RepoName":    "repo-name",
				"Commit":      "commit",
				"Path":        "path",
			},
		},
		{
			CommandPath: []string{"pr", "create"},
			Params:      app.BitbucketCreatePRParams{},
			FieldToFlag: map[string]string{
				"AccountName":       "account",
				"RepoOwner":         "repo-owner",
				"RepoName":          "repo-name",
				"Title":             "title",
				"Description":       "description",
				"SourceBranch":      "source-branch",
				"DestBranch":        "target-branch",
				"CloseSourceBranch": "close-source-branch",
				"Reviewers":         "reviewer",
				"Draft":             "draft",
			},
		},
		{
			CommandPath: []string{"pr", "read"},
			Params:      app.BitbucketReadPRParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
			},
		},
		{
			CommandPath: []string{"pr", "update"},
			Params:      app.BitbucketUpdatePRParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"Title":         "title",
				"Description":   "description",
				"Draft":         "draft",
			},
		},
		{
			CommandPath: []string{"pr", "approve"},
			Params:      app.BitbucketApprovePRParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
			},
		},
		{
			CommandPath: []string{"pr", "request-changes"},
			Params:      app.BitbucketRequestPRChangesParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
			},
		},
		{
			CommandPath: []string{"pr", "merge"},
			Params:      app.BitbucketMergePRParams{},
			FieldToFlag: map[string]string{
				"AccountName":       "account",
				"RepoOwner":         "repo-owner",
				"RepoName":          "repo-name",
				"PullRequestID":     "pr-id",
				"Message":           "message",
				"CloseSourceBranch": "close-source-branch",
				"MergeStrategy":     "strategy",
			},
		},
		{
			CommandPath: []string{"pr", "list-tasks"},
			Params:      app.BitbucketListTasksParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"Query":         "query",
				"Sort":          "sort",
				"PageLen":       "pagelen",
			},
		},
		{
			CommandPath: []string{"pr", "create-task"},
			Params:      app.BitbucketCreateTaskParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"Content":       "content",
				"CommentID":     "comment-id",
				"State":         "state",
			},
		},
		{
			CommandPath: []string{"pr", "update-task"},
			Params:      app.BitbucketUpdateTaskParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"TaskID":        "task-id",
				"Content":       "content",
				"State":         "state",
			},
		},
		{
			CommandPath: []string{"pr", "diffstat"},
			Params:      app.BitbucketGetPRDiffStatParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
			},
		},
		{
			CommandPath: []string{"pr", "diff"},
			Params:      app.BitbucketGetPRDiffParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"FilePaths":     "path",
				"ContextLines":  "context",
			},
		},
		{
			CommandPath: []string{"pr", "add-comment"},
			Params:      app.BitbucketAddPRCommentParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"Content":       "content",
				"FilePath":      "file-path",
				"LineFrom":      "line-from",
				"LineTo":        "line-to",
				"Pending":       "pending",
				"ParentID":      "parent-comment-id",
			},
		},
		{
			CommandPath: []string{"pr", "list-comments"},
			Params:      app.BitbucketListPRCommentsParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"Page":          "page",
				"PageLen":       "pagelen",
			},
		},
		{
			CommandPath: []string{"pr", "resolve-comment"},
			Params:      app.BitbucketResolvePRCommentParams{},
			FieldToFlag: map[string]string{
				"AccountName":   "account",
				"RepoOwner":     "repo-owner",
				"RepoName":      "repo-name",
				"PullRequestID": "pr-id",
				"CommentID":     "comment-id",
			},
		},
	}

	for _, spec := range specs {
		testName := strings.Join(spec.CommandPath, " ")
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd := mustFindCommand(t, setupCommands(), spec.CommandPath)
			fieldNames := exportedFieldNames(spec.Params)

			for fieldName, flagName := range spec.FieldToFlag {
				assert.Containsf(t, fieldNames, fieldName, "mapping references unknown field %s", fieldName)
				assert.NotNilf(t, cmd.Flags().Lookup(flagName), "missing flag --%s on %s", flagName, cmd.CommandPath())
			}

			for _, fieldName := range fieldNames {
				if _, found := spec.Exceptions[fieldName]; found {
					continue
				}
				_, mapped := spec.FieldToFlag[fieldName]
				assert.Truef(t, mapped, "field %s is not mapped for %s", fieldName, cmd.CommandPath())
			}
		})
	}
}

func TestCommandExamplesOnlyReferenceExistingFlags(t *testing.T) {
	t.Parallel()

	flagPattern := regexp.MustCompile(`--[a-z0-9-]+`)
	for _, cmd := range allCommands(setupCommands()) {
		if strings.TrimSpace(cmd.Example) == "" {
			continue
		}
		for line := range strings.SplitSeq(cmd.Example, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			target := commandForExampleLine(setupCommands(), line)
			exampleFlags := flagPattern.FindAllString(line, -1)
			for _, flag := range exampleFlags {
				flagName := strings.TrimPrefix(flag, "--")
				assert.NotNilf(
					t,
					lookupFlag(target, flagName),
					"example for %s references unknown flag %s in line %q",
					target.CommandPath(),
					flag,
					line,
				)
			}
		}
	}
}

func commandForExampleLine(root *cobra.Command, line string) *cobra.Command {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return root
	}

	start := 0
	if fields[0] == root.Name() {
		start = 1
	}

	current := root
	for _, token := range fields[start:] {
		if strings.HasPrefix(token, "-") {
			break
		}
		found := false
		for _, child := range current.Commands() {
			if child.Name() == token {
				current = child
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return current
}

func exportedFieldNames(v any) []string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	names := make([]string, 0, t.NumField())
	for field := range t.Fields() {
		if field.IsExported() {
			names = append(names, field.Name)
		}
	}
	slices.Sort(names)
	return names
}

func mustFindCommand(t *testing.T, root *cobra.Command, path []string) *cobra.Command {
	t.Helper()

	current := root
	for _, part := range path {
		found := false
		for _, child := range current.Commands() {
			if child.Name() == part {
				current = child
				found = true
				break
			}
		}
		require.Truef(t, found, "command not found: %s", strings.Join(path, " "))
	}
	return current
}

func allCommands(root *cobra.Command) []*cobra.Command {
	commands := make([]*cobra.Command, 0)
	var walk func(cmd *cobra.Command)
	walk = func(cmd *cobra.Command) {
		commands = append(commands, cmd)
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)
	return commands
}

func lookupFlag(cmd *cobra.Command, name string) *pflag.Flag {
	if f := cmd.Flags().Lookup(name); f != nil {
		return f
	}
	return cmd.InheritedFlags().Lookup(name)
}
