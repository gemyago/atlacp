package main

import (
	"context"
	"fmt"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newAuthCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Atlassian account management",
	}
	auth.AddCommand(
		newAuthStatusCmd(container, rootParams),
		newAuthAddCmd(container, rootParams),
		newAuthRemoveCmd(container, rootParams),
		newAuthSetDefaultCmd(container, rootParams),
	)
	return auth
}

type authAddOpts struct {
	Name       string
	Default    bool
	TokenType  string
	TokenValue string
}

func newAuthStatusCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "List configured accounts (tokens partially redacted)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAuthStatus(cmd, container, rootParams)
		},
	}
}

func newAuthAddCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var opts authAddOpts
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add or replace an account",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAuthAdd(cmd, container, rootParams, opts)
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "", "Account name")
	cmd.Flags().BoolVar(&opts.Default, "default", false, "Mark this account as the default")
	cmd.Flags().StringVar(&opts.TokenType, "token-type", "Bearer", "Token type (e.g. Bearer)")
	cmd.Flags().StringVar(&opts.TokenValue, "token-value", "", "Token value")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("token-value")
	return cmd
}

func newAuthRemoveCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an account by name",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAuthRemove(cmd, container, rootParams, name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Account name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newAuthSetDefaultCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "set-default",
		Short: "Set the default account by name",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAuthSetDefault(cmd, container, rootParams, name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Account name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

// authStatusAccount mirrors app.AtlassianAccount JSON with redacted token values for display.
type authStatusAccount struct {
	Name      string           `json:"name"`
	Default   bool             `json:"default"`
	Bitbucket *authStatusToken `json:"bitbucket,omitempty"`
	Jira      *authStatusToken `json:"jira,omitempty"`
}

type authStatusToken struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value"`
}

func runAuthStatus(cmd *cobra.Command, container *dig.Container, rootParams *rootCommandParams) error {
	return container.Invoke(func(deps execDeps, store *services.AccountsStore) error {
		return execAndWrite(cmd, deps, execArgs[struct{}, []authStatusAccount]{
			rootParams: rootParams,
			params:     struct{}{},
			target: func(_ context.Context, _ struct{}) ([]authStatusAccount, error) {
				raw := store.ListAccounts()
				out := make([]authStatusAccount, 0, len(raw))
				for _, a := range raw {
					out = append(out, redactAccountForStatus(a))
				}
				return out, nil
			},
		})
	})
}

func redactAccountForStatus(a app.AtlassianAccount) authStatusAccount {
	res := authStatusAccount{Name: a.Name, Default: a.Default}
	if a.Bitbucket != nil {
		t := *a.Bitbucket
		t.Value = redactTokenValue(t.Value)
		res.Bitbucket = &authStatusToken{Type: t.Type, Value: t.Value}
	}
	if a.Jira != nil {
		t := *a.Jira
		t.Value = redactTokenValue(t.Value)
		res.Jira = &authStatusToken{Type: t.Type, Value: t.Value}
	}
	return res
}

const redactTokenMaxPlainLen = 8

// redactTokenValue shows the first and last four characters with "***" between; strings of length
// eight or fewer are replaced entirely with "***".
func redactTokenValue(s string) string {
	if len(s) <= redactTokenMaxPlainLen {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

func runAuthAdd(cmd *cobra.Command, container *dig.Container, rootParams *rootCommandParams, opts authAddOpts) error {
	return container.Invoke(func(deps execDeps, store *services.AccountsStore) error {
		return execNoOutput(cmd, deps, rootParams, opts, func(_ context.Context, o authAddOpts) error {
			account := app.AtlassianAccount{
				Name:    o.Name,
				Default: o.Default,
				Bitbucket: &app.AtlassianToken{
					Type:  o.TokenType,
					Value: o.TokenValue,
				},
			}
			if upsertErr := store.Upsert(account); upsertErr != nil {
				return upsertErr
			}
			if saveErr := store.SaveToFile(rootParams.ResolvedAccountsFilePath); saveErr != nil {
				return fmt.Errorf("save accounts: %w", saveErr)
			}
			return nil
		})
	})
}

func runAuthRemove(cmd *cobra.Command, container *dig.Container, rootParams *rootCommandParams, name string) error {
	return container.Invoke(func(deps execDeps, store *services.AccountsStore) error {
		return execNoOutput(cmd, deps, rootParams, name, func(_ context.Context, n string) error {
			if removeErr := store.Remove(n); removeErr != nil {
				return removeErr
			}
			if saveErr := store.SaveToFile(rootParams.ResolvedAccountsFilePath); saveErr != nil {
				return fmt.Errorf("save accounts: %w", saveErr)
			}
			return nil
		})
	})
}

func runAuthSetDefault(cmd *cobra.Command, container *dig.Container, rootParams *rootCommandParams, name string) error {
	return container.Invoke(func(deps execDeps, store *services.AccountsStore) error {
		return execNoOutput(cmd, deps, rootParams, name, func(_ context.Context, n string) error {
			if setErr := store.SetDefault(n); setErr != nil {
				return setErr
			}
			if saveErr := store.SaveToFile(rootParams.ResolvedAccountsFilePath); saveErr != nil {
				return fmt.Errorf("save accounts: %w", saveErr)
			}
			return nil
		})
	})
}
