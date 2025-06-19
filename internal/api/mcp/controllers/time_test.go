package controllers

import (
	"context"
	"log/slog"
	"testing"

	"github.com/gemyago/golang-backend-boilerplate/internal/app"
	"github.com/go-faker/faker/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTimeControllerDeps() TimeControllerDeps {
	timeServiceDeps := app.TimeServiceDeps{
		RootLogger: slog.Default(),
	}
	timeService := app.NewTimeService(timeServiceDeps)

	return TimeControllerDeps{
		RootLogger:  slog.Default(),
		TimeService: timeService,
	}
}

func TestNewTimeController(t *testing.T) {
	t.Run("should create time controller with dependencies", func(t *testing.T) {
		deps := makeTimeControllerDeps()

		controller := NewTimeController(deps)

		require.NotNil(t, controller)
		require.NotNil(t, controller.logger)
		require.NotNil(t, controller.timeService)
	})
}

func TestTimeController_ToolDefinition(t *testing.T) {
	t.Run("should return valid time tool definition", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)

		serverTool := controller.newGetCurrentTimeServerTool()

		assert.Equal(t, "get_current_time", serverTool.Tool.Name)
		assert.Equal(t, "Get the current date and time in various formats", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})
}

func TestTimeController_HandleGetCurrentTime_ISO(t *testing.T) {
	t.Run("should handle ISO format request", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_time",
				Arguments: map[string]interface{}{
					"format": "iso",
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// Check that the result contains expected text
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: iso")
		}
	})
}

func TestTimeController_HandleGetCurrentTime_RFC3339(t *testing.T) {
	t.Run("should handle RFC3339 format request", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_time",
				Arguments: map[string]interface{}{
					"format": "rfc3339",
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// Check that the result contains expected text
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: rfc3339")
		}
	})
}

func TestTimeController_HandleGetCurrentTime_Unix(t *testing.T) {
	t.Run("should handle Unix format request", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_time",
				Arguments: map[string]interface{}{
					"format": "unix",
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// Check that the result contains expected text
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: unix")
		}
	})
}

func TestTimeController_HandleGetCurrentTime_DefaultFormat(t *testing.T) {
	t.Run("should default to ISO format when no format specified", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      "get_current_time",
				Arguments: nil, // No arguments
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// Check that the result contains expected text with ISO format
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: iso")
		}
	})
}

func TestTimeController_HandleGetCurrentTime_InvalidFormat(t *testing.T) {
	t.Run("should default to ISO format for invalid format", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_time",
				Arguments: map[string]interface{}{
					"format": faker.Word(), // Random invalid format
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// Check that the result contains expected text with ISO format (default)
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: iso")
		}
	})
}

func TestTimeController_HandleGetCurrentTime_ContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel context immediately

		serverTool := controller.newGetCurrentTimeServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_time",
				Arguments: map[string]interface{}{
					"format": "iso",
				},
			},
		}

		// Service should still work as it doesn't depend on context for time operations
		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)
	})
}

func TestTimeController_RegisterWithServer(t *testing.T) {
	t.Run("should register time tool with server", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)

		tools := controller.NewTools()

		require.NotEmpty(t, tools)
	})
}

func TestTimeController_Integration(t *testing.T) {
	t.Run("should work end-to-end", func(t *testing.T) {
		deps := makeTimeControllerDeps()
		controller := NewTimeController(deps)
		ctx := t.Context()

		// Get the server tool
		serverTool := controller.newGetCurrentTimeServerTool()
		assert.Equal(t, "get_current_time", serverTool.Tool.Name)

		// Simulate an MCP tool call
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: serverTool.Tool.Name,
				Arguments: map[string]interface{}{
					"format": "iso",
				},
			},
		}

		// Handle the tool call
		result, err := serverTool.Handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		// The result should contain a valid time
		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Current time:")
			assert.Contains(t, content.Text, "format: iso")
		}
	})
}
