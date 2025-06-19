//go:build !release

package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// testMCPServer encapsulates an MCP server and manages resources like pipes and context.
type testMCPServer struct {
	serverReader *io.PipeReader
	serverWriter *io.PipeWriter
	clientReader *io.PipeReader
	clientWriter *io.PipeWriter

	logBuffer bytes.Buffer

	transport transport.Interface
	client    *client.Client

	wg sync.WaitGroup
}

// newTestMCPServer creates a new MCP server instance with the given name, but does not start the server.
// Useful for tests where you need to add tools before starting the server.
func newTestMCPServer() *testMCPServer {
	server := &testMCPServer{}

	// Set up pipes for client-server communication
	server.serverReader, server.clientWriter = io.Pipe()
	server.clientReader, server.serverWriter = io.Pipe()

	// Return the configured server
	return server
}

// Start starts the server in a goroutine. Make sure to defer Close() after Start().
// When using NewServer(), the returned server is already started.
func (s *testMCPServer) Start(
	ctx context.Context,
	mcpServer *mcpserver.MCPServer,
) error {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		logger := log.New(&s.logBuffer, "", 0)

		stdioServer := mcpserver.NewStdioServer(mcpServer)
		stdioServer.SetErrorLogger(logger)

		if err := stdioServer.Listen(ctx, s.serverReader, s.serverWriter); err != nil {
			logger.Println("StdioServer.Listen failed:", err)
		}
	}()

	s.transport = transport.NewIO(s.clientReader, s.clientWriter, io.NopCloser(&s.logBuffer))
	if err := s.transport.Start(ctx); err != nil {
		return fmt.Errorf("transport.Start(): %w", err)
	}

	s.client = client.NewClient(s.transport)

	var initReq mcp.InitializeRequest
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	if _, err := s.client.Initialize(ctx, initReq); err != nil {
		return fmt.Errorf("client.Initialize(): %w", err)
	}

	return nil
}

// Client returns an MCP client connected to the server.
// The client is already initialized, i.e. you do _not_ need to call Client.Initialize().
func (s *testMCPServer) Client() *client.Client {
	return s.client
}
