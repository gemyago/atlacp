package main

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/config"
	"github.com/gemyago/atlacp/internal/di"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type rootCommandParams struct {
	Noop                     bool
	ResolvedAccountsFilePath string
	LogsOutputFile           string
}

func prepareBBMDAccountsFilePath(cfg *viper.Viper, rootParams *rootCommandParams, accountsFile string) error {
	pathResolver := services.NewAccountsFilePathResolver()
	var resolved string
	if accountsFile == "" {
		defaultPath, pathErr := pathResolver.DefaultPath()
		if pathErr != nil {
			return fmt.Errorf("resolve default atlassian accounts file path: %w", pathErr)
		}
		resolved = defaultPath
		cfg.Set("atlassian.accountsFilePath", defaultPath)
	} else {
		resolved = filepath.Clean(accountsFile)
		cfg.Set("atlassian.accountsFilePath", resolved)
	}
	rootParams.ResolvedAccountsFilePath = resolved
	if mkdirErr := pathResolver.EnsureParentDirsForFile(resolved); mkdirErr != nil {
		return fmt.Errorf("ensure atlassian accounts file parent directories: %w", mkdirErr)
	}
	return nil
}

func newRootCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bbmd",
		Short: "Bitbucket CLI — direct access to Bitbucket and account operations",
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().StringP("log-level", "l", "", "Produce logs with given level. Default is env specific.")
	cmd.PersistentFlags().StringVar(
		&rootParams.LogsOutputFile,
		"logs-file",
		"bbmd.log",
		"Write logs to this file (default bbmd.log in the current directory).",
	)
	cmd.PersistentFlags().Bool(
		"json-logs",
		false,
		"Indicates if logs should be in JSON format or text (default)",
	)
	cmd.PersistentFlags().StringP(
		"env",
		"e",
		"",
		"Env that the process is running in.",
	)
	cmd.PersistentFlags().StringP(
		"atlassian-accounts-file",
		"a",
		"",
		"Path to the Atlassian accounts file.",
	)
	cmd.PersistentFlags().BoolVar(
		&rootParams.Noop,
		"noop",
		false,
		"Dry-run: wire dependencies and skip real Bitbucket/account side effects.",
	)
	cfg := config.New()
	lo.Must0(cfg.BindPFlag("atlassian.accountsFilePath", cmd.PersistentFlags().Lookup("atlassian-accounts-file")))
	lo.Must0(cfg.BindPFlag("jsonLogs", cmd.PersistentFlags().Lookup("json-logs")))
	lo.Must0(cfg.BindPFlag("defaultLogLevel", cmd.PersistentFlags().Lookup("log-level")))
	lo.Must0(cfg.BindPFlag("env", cmd.PersistentFlags().Lookup("env")))
	cmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		err := config.Load(cfg, config.NewLoadOpts().WithEnv(cfg.GetString("env")))
		if err != nil {
			return err
		}

		accountsFile, err := cmd.Flags().GetString("atlassian-accounts-file")
		if err != nil {
			return fmt.Errorf("get atlassian-accounts-file flag: %w", err)
		}
		if prepErr := prepareBBMDAccountsFilePath(cfg, rootParams, accountsFile); prepErr != nil {
			return prepErr
		}

		var logLevel slog.Level
		if err = logLevel.UnmarshalText([]byte(cfg.GetString("defaultLogLevel"))); err != nil {
			return err
		}

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(cfg.GetBool("jsonLogs")).
				WithLogLevel(logLevel).
				WithOptionalOutputFile(rootParams.LogsOutputFile),
		)

		err = errors.Join(
			config.Provide(container, cfg),
			app.Register(container),
			services.Register(container),
			di.ProvideAll(container,
				di.ProvideValue(rootLogger),
			),
		)

		return lo.
			If(err != nil, fmt.Errorf("failed to inject dependencies: %w", err)).
			Else(nil)
	}

	cmd.AddCommand(
		newPRCmd(container, rootParams),
		newFileCmd(container, rootParams),
		newAuthCmd(container, rootParams),
	)

	return cmd
}
