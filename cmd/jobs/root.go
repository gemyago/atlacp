package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/config"
	"github.com/gemyago/atlacp/internal/di"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newRootCmd(container *dig.Container) *cobra.Command {
	logsOutputFile := ""

	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Command to run jobs",
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().StringP("log-level", "l", "", "Produce logs with given level. Default is env specific.")
	cmd.PersistentFlags().StringVar(
		&logsOutputFile,
		"logs-file",
		"",
		"Produce logs to file instead of stdout. Used for tests only.",
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
	cfg := config.New()
	lo.Must0(cfg.BindPFlag("jsonLogs", cmd.PersistentFlags().Lookup("json-logs")))
	lo.Must0(cfg.BindPFlag("defaultLogLevel", cmd.PersistentFlags().Lookup("log-level")))
	lo.Must0(cfg.BindPFlag("env", cmd.PersistentFlags().Lookup("env")))
	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		err := config.Load(cfg, config.NewLoadOpts().WithEnv(cfg.GetString("env")))
		if err != nil {
			return err
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
