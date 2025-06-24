package controllers

import (
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The testing.T parameter is used for logging in test context.
func makeBitbucketControllerDeps(t *testing.T) BitbucketControllerDeps {
	// Create a mock bitbucketService for testing
	mockBitbucketService := NewMockbitbucketService(t)

	// Use t for test logging context
	logger := diag.RootTestLogger().With("test", t.Name())

	return BitbucketControllerDeps{
		RootLogger:       logger,
		BitbucketService: mockBitbucketService,
	}
}

func TestBitbucketController(t *testing.T) {
	t.Run("should create Bitbucket controller with dependencies", func(t *testing.T) {
		deps := makeBitbucketControllerDeps(t)

		controller := NewBitbucketController(deps)

		require.NotNil(t, controller)
		require.NotNil(t, controller.logger)
		require.NotNil(t, controller.bitbucketService)
	})

	t.Run("tool definitions", func(t *testing.T) {
		t.Run("should define CreatePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newCreatePRServerTool()

			assert.Equal(t, "bitbucket_create_pr", serverTool.Tool.Name)
			assert.Equal(t, "Create a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ReadPR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newReadPRServerTool()

			assert.Equal(t, "bitbucket_read_pr", serverTool.Tool.Name)
			assert.Equal(t, "Get pull request details from Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define UpdatePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newUpdatePRServerTool()

			assert.Equal(t, "bitbucket_update_pr", serverTool.Tool.Name)
			assert.Equal(t, "Update a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ApprovePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newApprovePRServerTool()

			assert.Equal(t, "bitbucket_approve_pr", serverTool.Tool.Name)
			assert.Equal(t, "Approve a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define MergePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newMergePRServerTool()

			assert.Equal(t, "bitbucket_merge_pr", serverTool.Tool.Name)
			assert.Equal(t, "Merge a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})
	})

	t.Run("should register all tools", func(t *testing.T) {
		deps := makeBitbucketControllerDeps(t)
		controller := NewBitbucketController(deps)

		tools := controller.NewTools()

		require.Len(t, tools, 5) // Expect 5 tools: create, read, update, approve, merge
		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = tool.Tool.Name
		}
		assert.Contains(t, toolNames, "bitbucket_create_pr")
		assert.Contains(t, toolNames, "bitbucket_read_pr")
		assert.Contains(t, toolNames, "bitbucket_update_pr")
		assert.Contains(t, toolNames, "bitbucket_approve_pr")
		assert.Contains(t, toolNames, "bitbucket_merge_pr")
	})

	// Example of how to use mocks.GetMock for retrieving the mock instance:
	// t.Run("should handle CreatePR call", func(t *testing.T) {
	//     deps := makeBitbucketControllerDeps(t)
	//     mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
	//     controller := NewBitbucketController(deps)
	//
	//     // Setup mock expectations
	//     mockService.EXPECT().CreatePR(...).Return(...)
	//
	//     // Test the handler
	//     // ...
	// })

	// Future tests for each handler implementation will go here
	// These will be added once we implement the actual handlers
}
