package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func writeJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	if _, werr := fmt.Fprintln(cmd.OutOrStdout(), string(data)); werr != nil {
		return fmt.Errorf("write JSON: %w", werr)
	}
	return nil
}

type execDeps struct {
	dig.In

	RootLogger *slog.Logger
}

type execArgs[TParams any, TResult any] struct {
	rootParams *rootCommandParams
	params     TParams
	target     func(ctx context.Context, params TParams) (TResult, error)
}

func execAndWrite[TParams any, TResult any](
	cmd *cobra.Command,
	deps execDeps,
	args execArgs[TParams, TResult],
) error {
	if args.rootParams.Noop {
		deps.RootLogger.Info("noop: skipping execution", "command", cmd.Name())
		return nil
	}
	result, err := args.target(cmd.Context(), args.params)
	if err != nil {
		return err
	}
	return writeJSON(cmd, result)
}
