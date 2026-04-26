package main

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/gemyago/atlacp/internal/api/mcp/controllers"
	"github.com/gemyago/atlacp/internal/api/mcp/server"
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

func prepareMCPAccountsFilePath(
	cfg *viper.Viper,
	pathResolver *services.AtlacpPathResolver,
	accountsFile string,
) error {
	var resolved string
	if accountsFile == "" {
		defaultPath, pathErr := pathResolver.DefaultAccountsFilePath()
		if pathErr != nil {
			return fmt.Errorf("resolve default atlassian accounts file path: %w", pathErr)
		}
		resolved = defaultPath
		cfg.Set("atlassian.accountsFilePath", defaultPath)
	} else {
		resolved = filepath.Clean(accountsFile)
		cfg.Set("atlassian.accountsFilePath", resolved)
	}
	if mkdirErr := pathResolver.EnsureParentDirsForFile(resolved); mkdirErr != nil {
		return fmt.Errorf("ensure atlassian accounts file parent directories: %w", mkdirErr)
	}
	return nil
}

func prepareMCPLogsOutputFile(
	cmd *cobra.Command,
	pathResolver *services.AtlacpPathResolver,
	logsOutputFile *string,
) error {
	resolvedLogsPath := *logsOutputFile
	switch {
	case cmd.Flags().Changed("logs-file"):
		if resolvedLogsPath != "" {
			resolvedLogsPath = filepath.Clean(resolvedLogsPath)
		}
	case cmd.Name() != "stdio":
		resolvedLogsPath = ""
	default:
		defaultPath, pathErr := pathResolver.DefaultLogPath("bbcp.log")
		if pathErr != nil {
			return fmt.Errorf("resolve default logs output file path: %w", pathErr)
		}
		resolvedLogsPath = defaultPath
	}

	*logsOutputFile = resolvedLogsPath
	if resolvedLogsPath == "" {
		return nil
	}

	if mkdirErr := pathResolver.EnsureParentDirsForFile(resolvedLogsPath); mkdirErr != nil {
		return fmt.Errorf("ensure logs output file parent directories: %w", mkdirErr)
	}

	return nil
}

func newRootCmd(container *dig.Container) *cobra.Command {
	logsOutputFile := ""

	cmd := &cobra.Command{
		Use:   "bbcp",
		Short: "MCP (Model Context Protocol) server command",
		Long:  "Start MCP server with stdio or HTTP transport for providing tools to MCP clients",
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().StringP("log-level", "l", "", "Produce logs with given level. Default is env specific.")
	cmd.PersistentFlags().StringVar(
		&logsOutputFile,
		"logs-file",
		"",
		"Write logs to this file. If omitted, stdio uses its default log path and http logs to stdout.",
	)
	cmd.PersistentFlags().Bool(
		"json-logs",
		true,
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

		pathResolver := services.NewAtlacpPathResolver()
		accountsFile, err := cmd.Flags().GetString("atlassian-accounts-file")
		if err != nil {
			return fmt.Errorf("get atlassian-accounts-file flag: %w", err)
		}
		if prepErr := prepareMCPAccountsFilePath(cfg, pathResolver, accountsFile); prepErr != nil {
			return prepErr
		}
		if prepErr := prepareMCPLogsOutputFile(cmd, pathResolver, &logsOutputFile); prepErr != nil {
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
				WithOptionalOutputFile(logsOutputFile),
		)

		err = errors.Join(
			config.Provide(container, cfg),

			// app layer
			app.Register(container),

			// services
			services.Register(container),

			// mcp components
			controllers.Register(container),
			di.ProvideAll(container,
				server.NewMCPServer,
			),

			di.ProvideAll(container,
				di.ProvideValue(rootLogger),
			),
		)

		return lo.
			If(err != nil, fmt.Errorf("failed to inject dependencies: %w", err)).
			Else(nil)
	}
	return cmd
}
