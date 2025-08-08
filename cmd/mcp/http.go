package main

import (
	"context"
	"errors"
	"log/slog"
	"os/signal"
	"time"

	mcpserver "github.com/gemyago/atlacp/internal/api/mcp/server"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"golang.org/x/sys/unix"
)

type startHTTPServerParams struct {
	dig.In `ignore-unexported:"true"`

	noop bool

	RootLogger *slog.Logger

	MCPServer *mcpserver.MCPServer

	*services.ShutdownHooks
}

func startHTTPServer(rootCtx context.Context, params startHTTPServerParams) error {
	rootLogger := params.RootLogger
	httpServer := params.MCPServer.NewHTTPServer()

	shutdown := func() error {
		rootLogger.InfoContext(rootCtx, "Trying to shut down gracefully")
		ts := time.Now()

		err := params.ShutdownHooks.PerformShutdown(rootCtx)
		if err != nil {
			rootLogger.ErrorContext(rootCtx, "Failed to shut down gracefully", diag.ErrAttr(err))
		}

		rootLogger.InfoContext(rootCtx, "Service stopped",
			slog.Duration("duration", time.Since(ts)),
		)
		return err
	}

	signalCtx, cancel := signal.NotifyContext(rootCtx, unix.SIGINT, unix.SIGTERM)
	defer cancel()

	const startedComponents = 2
	startupErrors := make(chan error, startedComponents)
	go func() {
		if params.noop {
			rootLogger.InfoContext(signalCtx, "NOOP: Starting http server")
			startupErrors <- nil
			return
		}
		startupErrors <- httpServer.Start(signalCtx)
	}()

	var startupErr error
	select {
	case startupErr = <-startupErrors:
		if startupErr != nil {
			rootLogger.ErrorContext(rootCtx, "Server startup failed", "err", startupErr)
		}
	case <-signalCtx.Done(): // coverage-ignore
		// We will attempt to shut down in both cases
		// so doing it once on a next line
	}
	return errors.Join(startupErr, shutdown())
}

func newHTTPCmd(container *dig.Container) *cobra.Command {
	noop := false
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Start MCP server",
		Long:  "Start MCP server via HTTP transport",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(p startHTTPServerParams) error {
				p.noop = noop
				return startHTTPServer(cmd.Context(), p)
			})
		},
	}
	cmd.Flags().BoolVar(
		&noop,
		"noop",
		false,
		"Run in noop mode. Useful for testing if setup is all working.",
	)
	return cmd
}
