package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServer(t *testing.T) {
	makeMockDeps := func() MCPServerDeps {
		return MCPServerDeps{
			RootLogger:    diag.RootTestLogger(),
			ShutdownHooks: services.NewTestShutdownHooks(),
			Controllers:   []ToolsFactory{},
		}
	}

	makeToolCallRequest := func() mcp.CallToolRequest {
		return mcp.CallToolRequest{
			Request: mcp.Request{
				Method: "tools/call",
			},
			Params: mcp.CallToolParams{
				Name: "tool-1-" + faker.Word(),
			},
			Header: http.Header{},
		}
	}

	newToolCallResult := func() *mcp.CallToolResult {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(faker.Sentence()),
			},
		}
	}

	newToolsFactories := func(
		toolName string,
		handler mcpserver.ToolHandlerFunc,
	) []ToolsFactory {
		return []ToolsFactory{
			ToolsFactoryFunc(func() []mcpserver.ServerTool {
				return []mcpserver.ServerTool{
					{
						Tool:    mcp.Tool{Name: toolName},
						Handler: handler,
					},
				}
			}),
		}
	}

	t.Run("middleware", func(t *testing.T) {
		t.Run("should process success tool call", func(t *testing.T) {
			deps := makeMockDeps()

			wantCall := makeToolCallRequest()
			wantResult := newToolCallResult()

			deps.Controllers = newToolsFactories(
				wantCall.Params.Name,
				func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					assert.NotNil(t, ctx)
					assert.Equal(t, wantCall, req)
					return wantResult, nil
				})
			srv := NewMCPServer(deps)
			ctx := t.Context()
			testServer := newTestMCPServer()
			err := testServer.Start(ctx, srv.mcpServer)
			require.NoError(t, err)

			client := testServer.Client()

			gotResult, err := client.CallTool(ctx, wantCall)
			require.NoError(t, err)
			assert.Equal(t, wantResult, gotResult)
		})

		t.Run("should setup correlation id in context", func(t *testing.T) {
			deps := makeMockDeps()

			wantCall := makeToolCallRequest()

			contextChecked := false
			deps.Controllers = newToolsFactories(
				wantCall.Params.Name,
				func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					diagCtx := diag.GetLogAttributesFromContext(ctx)
					assert.NotEmpty(t, diagCtx.CorrelationID)
					contextChecked = true
					return newToolCallResult(), nil
				})
			srv := NewMCPServer(deps)
			ctx := t.Context()
			testServer := newTestMCPServer()
			err := testServer.Start(ctx, srv.mcpServer)
			require.NoError(t, err)

			client := testServer.Client()

			_, err = client.CallTool(ctx, wantCall)
			require.NoError(t, err)
			assert.True(t, contextChecked)
		})

		t.Run("should reuse correlation id from context", func(t *testing.T) {
			deps := makeMockDeps()

			wantCorrelationID := faker.UUIDHyphenated()
			wantCall := makeToolCallRequest()

			callCtx := diag.SetLogAttributesToContext(t.Context(), diag.LogAttributes{
				CorrelationID: slog.StringValue(wantCorrelationID),
			})

			contextChecked := false
			deps.Controllers = newToolsFactories(
				wantCall.Params.Name,
				func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					diagCtx := diag.GetLogAttributesFromContext(ctx)
					assert.Equal(t, wantCorrelationID, diagCtx.CorrelationID.String())
					contextChecked = true
					return newToolCallResult(), nil
				})
			srv := NewMCPServer(deps)
			testServer := newTestMCPServer()
			err := testServer.Start(callCtx, srv.mcpServer)
			require.NoError(t, err)

			client := testServer.Client()

			_, err = client.CallTool(callCtx, wantCall)
			require.NoError(t, err)
			assert.True(t, contextChecked)
		})

		t.Run("should respond with error if tool call fails", func(t *testing.T) {
			deps := makeMockDeps()

			wantCall := makeToolCallRequest()
			errorMsg := faker.Sentence()
			wantError := errors.New(errorMsg)

			deps.Controllers = newToolsFactories(
				wantCall.Params.Name,
				func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					assert.NotNil(t, ctx)
					assert.Equal(t, wantCall, req)
					return nil, wantError
				})
			srv := NewMCPServer(deps)
			ctx := t.Context()
			testServer := newTestMCPServer()
			err := testServer.Start(ctx, srv.mcpServer)
			require.NoError(t, err)

			client := testServer.Client()

			result, err := client.CallTool(ctx, wantCall)
			require.NoError(t, err) // Now expecting no error
			require.NotNil(t, result)
			assert.True(t, result.IsError)

			if len(result.Content) > 0 {
				content, ok := mcp.AsTextContent(result.Content[0])
				require.True(t, ok, "Error content should be text content")
				assert.Contains(t, content.Text, errorMsg)
			}
		})

		t.Run("should include correlation id in error message", func(t *testing.T) {
			deps := makeMockDeps()

			wantCorrelationID := faker.UUIDHyphenated()
			wantCall := makeToolCallRequest()
			errorMsg := faker.Sentence()
			wantError := errors.New(errorMsg)

			callCtx := diag.SetLogAttributesToContext(t.Context(), diag.LogAttributes{
				CorrelationID: slog.StringValue(wantCorrelationID),
			})

			deps.Controllers = newToolsFactories(
				wantCall.Params.Name,
				func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					assert.NotNil(t, ctx)
					assert.Equal(t, wantCall, req)
					return nil, wantError
				})
			srv := NewMCPServer(deps)
			testServer := newTestMCPServer()
			err := testServer.Start(callCtx, srv.mcpServer)
			require.NoError(t, err)

			client := testServer.Client()

			result, err := client.CallTool(callCtx, wantCall)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsError)

			if len(result.Content) > 0 {
				content, ok := mcp.AsTextContent(result.Content[0])
				require.True(t, ok, "Error content should be text content")
				assert.Contains(t, content.Text, errorMsg)
				assert.Contains(t, content.Text, "Error details:")
				assert.Contains(t, content.Text, "CorrelationID:")
				assert.Contains(t, content.Text, wantCorrelationID)
			}
		})
	})
}
