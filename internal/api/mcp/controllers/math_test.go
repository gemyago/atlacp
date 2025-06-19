package controllers

import (
	"log/slog"
	"testing"

	"github.com/gemyago/golang-backend-boilerplate/internal/app"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMathControllerDeps() MathControllerDeps {
	mathServiceDeps := app.MathServiceDeps{
		RootLogger: slog.Default(),
	}
	mathService := app.NewMathService(mathServiceDeps)

	return MathControllerDeps{
		MathService: mathService,
		RootLogger:  slog.Default(),
	}
}

func TestNewMathController(t *testing.T) {
	t.Run("should create math controller with dependencies", func(t *testing.T) {
		deps := makeMathControllerDeps()

		controller := NewMathController(deps)

		require.NotNil(t, controller)
		require.NotNil(t, controller.mathService)
		require.NotNil(t, controller.logger)
	})
}

func TestMathController_ToolDefinitions(t *testing.T) {
	t.Run("should return calculate tool with correct schema", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		serverTool := controller.newCalculateServerTool()

		assert.Equal(t, "calculate", serverTool.Tool.Name)
		assert.Equal(t, "Perform mathematical calculations (add, subtract, multiply, divide)", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})

	t.Run("should return add tool with correct schema", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		serverTool := controller.newAddServerTool()

		assert.Equal(t, "add", serverTool.Tool.Name)
		assert.Equal(t, "Add two numbers together", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})

	t.Run("should return subtract tool with correct schema", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		serverTool := controller.newSubtractServerTool()

		assert.Equal(t, "subtract", serverTool.Tool.Name)
		assert.Equal(t, "Subtract second number from first number", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})

	t.Run("should return multiply tool with correct schema", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		serverTool := controller.newMultiplyServerTool()

		assert.Equal(t, "multiply", serverTool.Tool.Name)
		assert.Equal(t, "Multiply two numbers together", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})

	t.Run("should return divide tool with correct schema", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		serverTool := controller.newDivideServerTool()

		assert.Equal(t, "divide", serverTool.Tool.Name)
		assert.Equal(t, "Divide first number by second number", serverTool.Tool.Description)
		assert.NotNil(t, serverTool.Tool.InputSchema)
		assert.NotNil(t, serverTool.Handler)
	})
}

func TestMathController_HandleCalculate(t *testing.T) {
	t.Run("should handle calculate add request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"operation": "add",
					"a":         5.0,
					"b":         3.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Result: 8")
			assert.Contains(t, content.Text, "operation: add")
		}
	})

	t.Run("should handle calculate multiply request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"operation": "multiply",
					"a":         6.0,
					"b":         7.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "Result: 42")
			assert.Contains(t, content.Text, "operation: multiply")
		}
	})

	t.Run("should handle invalid operation parameter", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"operation": 123, // Invalid type
					"a":         5.0,
					"b":         3.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})

	t.Run("should handle missing parameters", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"operation": "add",
					"a":         5.0,
					// Missing "b" parameter
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})

	t.Run("should handle non-object arguments for calculate", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      "calculate",
				Arguments: "not_an_object", // Invalid arguments type
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})

	t.Run("should handle missing operation parameter", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"a": 5.0,
					"b": 3.0,
					// Missing "operation" parameter
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})

	t.Run("should handle division by zero error", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newCalculateServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "calculate",
				Arguments: map[string]interface{}{
					"operation": "divide",
					"a":         5.0,
					"b":         0.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Calculation failed")
		}
	})
}

func TestMathController_HandleAdd(t *testing.T) {
	t.Run("should handle add request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newAddServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add",
				Arguments: map[string]interface{}{
					"a": 7.0,
					"b": 3.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "7 + 3 = 10")
		}
	})

	t.Run("should handle invalid parameters", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newAddServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add",
				Arguments: map[string]interface{}{
					"a": "invalid", // Invalid type
					"b": 3.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})
}

func TestMathController_HandleSubtract(t *testing.T) {
	t.Run("should handle subtract request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newSubtractServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "subtract",
				Arguments: map[string]interface{}{
					"a": 10.0,
					"b": 4.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "10 - 4 = 6")
		}
	})

	t.Run("should handle invalid parameters", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newSubtractServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "subtract",
				Arguments: map[string]interface{}{
					"a": "invalid", // Invalid type
					"b": 3.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})
}

func TestMathController_HandleMultiply(t *testing.T) {
	t.Run("should handle multiply request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newMultiplyServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "multiply",
				Arguments: map[string]interface{}{
					"a": 6.0,
					"b": 7.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "6 × 7 = 42")
		}
	})

	t.Run("should handle invalid parameters", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newMultiplyServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "multiply",
				Arguments: map[string]interface{}{
					"a": 6.0,
					"b": "invalid", // Invalid type
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})
}

func TestMathController_HandleDivide(t *testing.T) {
	t.Run("should handle divide request successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newDivideServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "divide",
				Arguments: map[string]interface{}{
					"a": 20.0,
					"b": 4.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "First content should be text content")
			assert.Contains(t, content.Text, "20 ÷ 4 = 5")
		}
	})

	t.Run("should handle division by zero error", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newDivideServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "divide",
				Arguments: map[string]interface{}{
					"a": 10.0,
					"b": 0.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Division failed")
		}
	})

	t.Run("should handle invalid parameters", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newDivideServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "divide",
				Arguments: map[string]interface{}{
					"a": "invalid", // Invalid type
					"b": 4.0,
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})
}

// Test for HandleAdd missing parameter error cases.
func TestMathController_HandleAdd_ParameterErrors(t *testing.T) {
	t.Run("should handle missing 'a' parameter", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newAddServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "add",
				Arguments: map[string]interface{}{
					"b": 3.0,
					// Missing "a" parameter
				},
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})

	t.Run("should handle non-object arguments", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)
		ctx := t.Context()

		serverTool := controller.newAddServerTool()
		handler := serverTool.Handler

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      "add",
				Arguments: "invalid_arguments", // Not an object
			},
		}

		result, err := handler(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.NotEmpty(t, result.Content)

		if len(result.Content) > 0 {
			content, ok := mcp.AsTextContent(result.Content[0])
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Invalid parameters")
		}
	})
}

func TestMathController_ParameterExtraction(t *testing.T) {
	t.Run("should extract number parameters correctly", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		// Test with float64
		args := map[string]interface{}{
			"a": 5.5,
			"b": 3.2,
		}
		a, b, err := controller.extractNumberParams(args)
		require.NoError(t, err)
		assert.InEpsilon(t, 5.5, a, 0.0001)
		assert.InEpsilon(t, 3.2, b, 0.0001)

		// Test with int
		args = map[string]interface{}{
			"a": 5,
			"b": 3,
		}
		a, b, err = controller.extractNumberParams(args)
		require.NoError(t, err)
		assert.InEpsilon(t, 5.0, a, 0.0001)
		assert.InEpsilon(t, 3.0, b, 0.0001)

		// Test with int64
		args = map[string]interface{}{
			"a": int64(5),
			"b": int64(3),
		}
		a, b, err = controller.extractNumberParams(args)
		require.NoError(t, err)
		assert.InEpsilon(t, 5.0, a, 0.0001)
		assert.InEpsilon(t, 3.0, b, 0.0001)
	})

	t.Run("should handle invalid argument types", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		// Test with invalid arguments object
		_, _, err := controller.extractNumberParams("invalid")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "arguments must be an object")

		// Test with missing parameters
		args := map[string]interface{}{
			"a": 5.0,
			// Missing "b"
		}
		_, _, err = controller.extractNumberParams(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "b parameter is required")

		// Test with invalid parameter type
		args = map[string]interface{}{
			"a": "invalid",
			"b": 3.0,
		}
		_, _, err = controller.extractNumberParams(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "a parameter must be a number")
	})

	t.Run("should extract calculate parameters correctly", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		args := map[string]interface{}{
			"operation": "multiply",
			"a":         6.0,
			"b":         7.0,
		}

		operation, a, b, err := controller.extractCalculateParams(args)
		require.NoError(t, err)
		assert.Equal(t, "multiply", operation)
		assert.InEpsilon(t, 6.0, a, 0.0001)
		assert.InEpsilon(t, 7.0, b, 0.0001)
	})
}

func TestMathController_NewTools(t *testing.T) {
	t.Run("should return all tools successfully", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		tools := controller.NewTools()

		require.NotEmpty(t, tools)
	})
}

func TestMathController_ParameterTypeSafety(t *testing.T) {
	t.Run("should validate parameter types strictly", func(t *testing.T) {
		deps := makeMathControllerDeps()
		controller := NewMathController(deps)

		testCases := []struct {
			name      string
			value     interface{}
			expected  float64
			shouldErr bool
		}{
			{"float64", 5.5, 5.5, false},
			{"int", 5, 5.0, false},
			{"int64", int64(5), 5.0, false},
			{"string", "5", 0, true},
			{"bool", true, 0, true},
			{"nil", nil, 0, true},
			{"array", []int{1, 2, 3}, 0, true},
			{"map", map[string]int{"a": 1}, 0, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				args := map[string]interface{}{
					"test_param": tc.value,
				}

				result, err := controller.extractNumberParam(args, "test_param")

				if tc.shouldErr {
					require.Error(t, err)
					assert.Contains(t, err.Error(), "parameter must be a number")
				} else {
					require.NoError(t, err)
					assert.InEpsilon(t, tc.expected, result, 0.0001)
				}
			})
		}
	})
}
