package controllers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/golang-backend-boilerplate/internal/app"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/dig"
)

// TimeControllerDeps contains dependencies for the time MCP controller.
type TimeControllerDeps struct {
	dig.In

	RootLogger  *slog.Logger
	TimeService *app.TimeService
}

// TimeController provides MCP time tool functionality.
type TimeController struct {
	logger      *slog.Logger
	timeService *app.TimeService
}

// NewTimeController creates a new time MCP controller.
func NewTimeController(deps TimeControllerDeps) *TimeController {
	return &TimeController{
		logger:      deps.RootLogger.WithGroup("mcp.time-controller"),
		timeService: deps.TimeService,
	}
}

// newGetCurrentTimeServerTool returns a server tool for getting current time.
func (tc *TimeController) newGetCurrentTimeServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"get_current_time",
		mcp.WithDescription("Get the current date and time in various formats"),
		mcp.WithString("format",
			mcp.Description("Time format to return (iso, rfc3339, or unix)"),
			mcp.Enum("iso", "rfc3339", "unix"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		format := app.TimeFormatISO
		if request.Params.Arguments != nil {
			if args, argsOk := request.Params.Arguments.(map[string]interface{}); argsOk {
				if formatStr, formatOk := args["format"].(string); formatOk {
					switch formatStr {
					case string(app.TimeFormatRFC3339):
						format = app.TimeFormatRFC3339
					case string(app.TimeFormatUnix):
						format = app.TimeFormatUnix
					case string(app.TimeFormatISO):
						format = app.TimeFormatISO
					default:
						// Invalid format, fallback to ISO
						format = app.TimeFormatISO
					}
				}
			}
		}

		timeRequest := &app.TimeRequest{Format: format}
		timeResponse, err := tc.timeService.GetCurrentTime(ctx, timeRequest)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current time: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Current time: %s (format: %s)",
			timeResponse.Time, timeResponse.Format)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// NewTools returns all time tools.
// Satisfies the ToolsFactory interface.
func (tc *TimeController) NewTools() []server.ServerTool {
	return []server.ServerTool{
		tc.newGetCurrentTimeServerTool(),
	}
}
