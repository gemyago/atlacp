package server

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/gemyago/golang-backend-boilerplate/internal/diag"
	"github.com/gemyago/golang-backend-boilerplate/internal/services"
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
			wantError := errors.New(faker.Sentence())

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

			_, err = client.CallTool(ctx, wantCall)
			require.Error(t, err)
			assert.Equal(t, wantError, err)
		})
	})
}
