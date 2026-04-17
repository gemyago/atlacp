package main

import (
	"encoding/json"
	"fmt"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newAuthCmd(container *dig.Container) *cobra.Command {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Atlassian account management",
	}
	auth.AddCommand(
		newAuthStatusCmd(container),
		newAuthAddCmd(container),
		newAuthRemoveCmd(container),
		newAuthSetDefaultCmd(container),
	)
	return auth
}

func newAuthStatusCmd(container *dig.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "List configured accounts (tokens partially redacted)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAuthStatus(cmd, container)
		},
	}
}

func newAuthAddCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add or replace an account",
		RunE: func(c *cobra.Command, _ []string) error {
			return runAuthAdd(container, c)
		},
	}
	cmd.Flags().String("name", "", "Account name")
	cmd.Flags().Bool("default", false, "Mark this account as the default")
	cmd.Flags().String("token-type", "Bearer", "Token type (e.g. Bearer)")
	cmd.Flags().String("token-value", "", "Token value")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("token-value")
	return cmd
}

func newAuthRemoveCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an account by name",
		RunE: func(c *cobra.Command, _ []string) error {
			return runAuthRemove(container, c)
		},
	}
	cmd.Flags().String("name", "", "Account name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newAuthSetDefaultCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default",
		Short: "Set the default account by name",
		RunE: func(c *cobra.Command, _ []string) error {
			return runAuthSetDefault(container, c)
		},
	}
	cmd.Flags().String("name", "", "Account name")
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

func runAuthStatus(cmd *cobra.Command, container *dig.Container) error {
	return container.Invoke(func(store *services.AccountsStore) error {
		if noop {
			return nil
		}
		raw := store.ListAccounts()
		out := make([]authStatusAccount, 0, len(raw))
		for _, a := range raw {
			out = append(out, redactAccountForStatus(a))
		}
		data, marshalErr := json.MarshalIndent(out, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("marshal auth status: %w", marshalErr)
		}
		if _, writeErr := fmt.Fprintln(cmd.OutOrStdout(), string(data)); writeErr != nil {
			return fmt.Errorf("write auth status: %w", writeErr)
		}
		return nil
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

func runAuthAdd(container *dig.Container, cmd *cobra.Command) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return fmt.Errorf("get name flag: %w", err)
	}
	isDefault, err := cmd.Flags().GetBool("default")
	if err != nil {
		return fmt.Errorf("get default flag: %w", err)
	}
	tokenType, err := cmd.Flags().GetString("token-type")
	if err != nil {
		return fmt.Errorf("get token-type flag: %w", err)
	}
	tokenValue, err := cmd.Flags().GetString("token-value")
	if err != nil {
		return fmt.Errorf("get token-value flag: %w", err)
	}
	account := app.AtlassianAccount{
		Name:    name,
		Default: isDefault,
		Bitbucket: &app.AtlassianToken{
			Type:  tokenType,
			Value: tokenValue,
		},
	}
	return container.Invoke(func(store *services.AccountsStore) error {
		if noop {
			return nil
		}
		if upsertErr := store.Upsert(account); upsertErr != nil {
			return upsertErr
		}
		if saveErr := store.SaveToFile(resolvedAccountsFilePath); saveErr != nil {
			return fmt.Errorf("save accounts: %w", saveErr)
		}
		return nil
	})
}

func runAuthRemove(container *dig.Container, cmd *cobra.Command) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return fmt.Errorf("get name flag: %w", err)
	}
	return container.Invoke(func(store *services.AccountsStore) error {
		if noop {
			return nil
		}
		if removeErr := store.Remove(name); removeErr != nil {
			return removeErr
		}
		if saveErr := store.SaveToFile(resolvedAccountsFilePath); saveErr != nil {
			return fmt.Errorf("save accounts: %w", saveErr)
		}
		return nil
	})
}

func runAuthSetDefault(container *dig.Container, cmd *cobra.Command) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return fmt.Errorf("get name flag: %w", err)
	}
	return container.Invoke(func(store *services.AccountsStore) error {
		if noop {
			return nil
		}
		if setErr := store.SetDefault(name); setErr != nil {
			return setErr
		}
		if saveErr := store.SaveToFile(resolvedAccountsFilePath); saveErr != nil {
			return fmt.Errorf("save accounts: %w", saveErr)
		}
		return nil
	})
}
