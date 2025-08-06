package main

import (
	"context"
	"errors"
	"log/slog"
	"os/signal"
	"time"

	httpserver "github.com/gemyago/atlacp/internal/api/http/server"
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
	sse  bool

	RootLogger *slog.Logger

	MCPServer *mcpserver.MCPServer

	*services.ShutdownHooks
}

func startHTTPServer(rootCtx context.Context, params startHTTPServerParams) error {
	rootLogger := params.RootLogger
	var httpServer *httpserver.HTTPServer
	if params.sse {
		httpServer = params.MCPServer.NewSSEServer()
	} else {
		httpServer = params.MCPServer.NewStreamableHTTPServer()
	}

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
	sse := false
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Start MCP server (Streamable HTTP by default)",
		Long:  "Start MCP server using using Streamable HTTP (default) or SSE (run with `--sse` flag)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(p startHTTPServerParams) error {
				p.noop = noop
				p.sse = sse
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
	cmd.Flags().BoolVar(
		&sse,
		"sse",
		false,
		"Start SSE server instead of Streamable HTTP server.",
	)
	return cmd
}
