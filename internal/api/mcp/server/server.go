package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	httpserver "github.com/gemyago/atlacp/internal/api/http/server"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/gofrs/uuid/v5"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/dig"
)

// Constants for server configuration.
const (
	httpReadTimeout  = 30 * time.Second
	httpWriteTimeout = 30 * time.Second
	httpIdleTimeout  = 120 * time.Second
	shutdownTimeout  = 10 * time.Second
)

type ToolsFactory interface {
	NewTools() []mcpserver.ServerTool
}

type ToolsFactoryFunc func() []mcpserver.ServerTool

func (f ToolsFactoryFunc) NewTools() []mcpserver.ServerTool {
	return f()
}

// MCPServerDeps contains dependencies for creating the MCP server.
type MCPServerDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	Name     string `name:"config.mcpServer.name"`
	Version  string `name:"config.mcpServer.version"`
	HTTPHost string `name:"config.mcpServer.httpHost"`
	HTTPPort int    `name:"config.mcpServer.httpPort"`

	// services
	*services.ShutdownHooks

	// controllers
	Controllers []ToolsFactory `group:"mcp-controllers"`
}

// ToolInfo contains information about a registered tool.
type ToolInfo struct {
	Tool    mcp.Tool
	Handler mcpserver.ToolHandlerFunc
}

// MCPServer wraps the mcp-go server with additional functionality.
type MCPServer struct {
	mcpServer     *mcpserver.MCPServer
	deps          MCPServerDeps
	logger        *slog.Logger
	shutdownHooks *services.ShutdownHooks
}

// NewMCPServer creates a new MCP server instance.
func NewMCPServer(deps MCPServerDeps) *MCPServer {
	logger := deps.RootLogger.WithGroup("mcp-server")

	mcpServer := mcpserver.NewMCPServer(
		deps.Name,
		deps.Version,
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithToolHandlerMiddleware(
			func(next mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
				return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					nextCtx := ctx
					diagCtx := diag.GetLogAttributesFromContext(nextCtx)

					// We may need to revisit this. It may be so that the diag context is always set
					// for the stdio transport at least.
					if diagCtx.CorrelationID.Kind() != slog.KindString {
						diagCtx.CorrelationID = slog.StringValue(uuid.Must(uuid.NewV4()).String())
						nextCtx = diag.SetLogAttributesToContext(nextCtx, diagCtx)
					}

					// It may be quite verbose and we may want to log just the "processed" part.
					logger.InfoContext(nextCtx, "Processing tool call",
						slog.String("tool", req.Params.Name),
						slog.Any("params", req.Params),
						slog.Any("meta", req.Params.Meta),
					)

					res, err := next(nextCtx, req)
					if err != nil {
						diagCtx = diag.GetLogAttributesFromContext(nextCtx)
						correlationID := diagCtx.CorrelationID.String()

						logger.ErrorContext(nextCtx, "Error processing tool call",
							slog.String("tool", req.Params.Name),
							slog.Any("error", err),
						)

						errorMessage := fmt.Sprintf("Error details: %s; CorrelationID: %s", err.Error(), correlationID)
						return mcp.NewToolResultError(errorMessage), nil
					}

					logger.InfoContext(nextCtx, "Tool call processed",
						slog.String("tool", req.Params.Name),
					)
					return res, nil
				}
			},
		),
		mcpserver.WithRecovery(),
	)

	mcpSrv := &MCPServer{
		deps:          deps,
		mcpServer:     mcpServer,
		logger:        deps.RootLogger.WithGroup("mcp-server"),
		shutdownHooks: deps.ShutdownHooks,
	}

	for _, controller := range deps.Controllers {
		tools := controller.NewTools()
		mcpSrv.mcpServer.AddTools(tools...)
	}

	return mcpSrv
}

// ListenStdioServer starts the MCP server with stdio transport.
func (s *MCPServer) ListenStdioServer(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
) error { // coverage-ignore -- Challenging to test this
	stdioSrv := mcpserver.NewStdioServer(s.mcpServer)
	s.logger.InfoContext(ctx, "Starting MCP server with stdio transport",
		slog.String("name", s.deps.Name),
		slog.String("version", s.deps.Version))

	return stdioSrv.Listen(ctx, stdin, stdout)
}

// NewStreamableHTTPServer creates a new streamable HTTP server.
func (s *MCPServer) NewStreamableHTTPServer() *httpserver.HTTPServer {
	return httpserver.NewHTTPServer(httpserver.HTTPServerDeps{
		RootLogger: s.logger,

		Host:              s.deps.HTTPHost,
		Port:              s.deps.HTTPPort,
		IdleTimeout:       httpIdleTimeout,
		ReadHeaderTimeout: httpReadTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,

		ShutdownHooks: s.shutdownHooks,
		Handler: mcpserver.NewStreamableHTTPServer(
			s.mcpServer,
			mcpserver.WithStateLess(true),
		),
	})
}
