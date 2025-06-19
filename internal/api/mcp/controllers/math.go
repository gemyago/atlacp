package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/golang-backend-boilerplate/internal/app"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/dig"
)

// MathController provides MCP tools for mathematical operations.

// MathControllerDeps contains dependencies for the math MCP controller.
type MathControllerDeps struct {
	dig.In

	MathService *app.MathService
	RootLogger  *slog.Logger
}

// MathController implements MCP tools for mathematical operations.
type MathController struct {
	mathService *app.MathService
	logger      *slog.Logger
}

// NewMathController creates a new math MCP controller.
func NewMathController(deps MathControllerDeps) *MathController {
	return &MathController{
		mathService: deps.MathService,
		logger:      deps.RootLogger.WithGroup("mcp.math-controller"),
	}
}

// newCalculateServerTool returns a server tool for generic calculations.
func (mc *MathController) newCalculateServerTool() mcpserver.ServerTool {
	tool := mcp.NewTool(
		"calculate",
		mcp.WithDescription("Perform mathematical calculations (add, subtract, multiply, divide)"),
		mcp.WithString("operation",
			mcp.Description("Mathematical operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("a", mcp.Description("First number")),
		mcp.WithNumber("b", mcp.Description("Second number")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		operation, a, b, err := mc.extractCalculateParams(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
		}

		mathRequest := &app.MathRequest{
			Operation: app.MathOperation(operation),
			A:         a,
			B:         b,
		}

		response, err := mc.mathService.Calculate(ctx, mathRequest)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Calculation failed: %v", err)), nil
		}

		resultText := fmt.Sprintf("Result: %g (operation: %s, a: %g, b: %g)",
			response.Result, response.Operation, response.A, response.B)

		return mcp.NewToolResultText(resultText), nil
	}

	return mcpserver.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newAddServerTool returns a server tool for addition.
func (mc *MathController) newAddServerTool() mcpserver.ServerTool {
	tool := mcp.NewTool(
		"add",
		mcp.WithDescription("Add two numbers together"),
		mcp.WithNumber("a", mcp.Description("First number to add")),
		mcp.WithNumber("b", mcp.Description("Second number to add")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, b, err := mc.extractNumberParams(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
		}

		response, err := mc.mathService.Add(ctx, a, b)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Addition failed: %v", err)), nil
		}

		resultText := fmt.Sprintf("Result: %g + %g = %g", response.A, response.B, response.Result)

		return mcp.NewToolResultText(resultText), nil
	}

	return mcpserver.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newSubtractServerTool returns a server tool for subtraction.
func (mc *MathController) newSubtractServerTool() mcpserver.ServerTool {
	tool := mcp.NewTool(
		"subtract",
		mcp.WithDescription("Subtract second number from first number"),
		mcp.WithNumber("a", mcp.Description("Number to subtract from")),
		mcp.WithNumber("b", mcp.Description("Number to subtract")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, b, err := mc.extractNumberParams(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
		}

		response, err := mc.mathService.Subtract(ctx, a, b)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Subtraction failed: %v", err)), nil
		}

		resultText := fmt.Sprintf("Result: %g - %g = %g", response.A, response.B, response.Result)

		return mcp.NewToolResultText(resultText), nil
	}

	return mcpserver.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newMultiplyServerTool returns a server tool for multiplication.
func (mc *MathController) newMultiplyServerTool() mcpserver.ServerTool {
	tool := mcp.NewTool(
		"multiply",
		mcp.WithDescription("Multiply two numbers together"),
		mcp.WithNumber("a", mcp.Description("First number to multiply")),
		mcp.WithNumber("b", mcp.Description("Second number to multiply")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, b, err := mc.extractNumberParams(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
		}

		response, err := mc.mathService.Multiply(ctx, a, b)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Multiplication failed: %v", err)), nil
		}

		resultText := fmt.Sprintf("Result: %g × %g = %g", response.A, response.B, response.Result)

		return mcp.NewToolResultText(resultText), nil
	}

	return mcpserver.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newDivideServerTool returns a server tool for division.
func (mc *MathController) newDivideServerTool() mcpserver.ServerTool {
	tool := mcp.NewTool(
		"divide",
		mcp.WithDescription("Divide first number by second number"),
		mcp.WithNumber("a", mcp.Description("Dividend (number to be divided)")),
		mcp.WithNumber("b", mcp.Description("Divisor (number to divide by)")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, b, err := mc.extractNumberParams(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid parameters: %v", err)), nil
		}

		response, err := mc.mathService.Divide(ctx, a, b)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Division failed: %v", err)), nil
		}

		resultText := fmt.Sprintf("Result: %g ÷ %g = %g", response.A, response.B, response.Result)

		return mcp.NewToolResultText(resultText), nil
	}

	return mcpserver.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// extractCalculateParams extracts operation, a, and b parameters from arguments.
func (mc *MathController) extractCalculateParams(args interface{}) (string, float64, float64, error) {
	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return "", 0, 0, errors.New("arguments must be an object")
	}

	operation, ok := argsMap["operation"].(string)
	if !ok {
		return "", 0, 0, errors.New("operation parameter is required and must be a string")
	}

	a, err := mc.extractNumberParam(argsMap, "a")
	if err != nil {
		return "", 0, 0, err
	}

	b, err := mc.extractNumberParam(argsMap, "b")
	if err != nil {
		return "", 0, 0, err
	}

	return operation, a, b, nil
}

// extractNumberParams extracts a and b number parameters from arguments.
func (mc *MathController) extractNumberParams(args interface{}) (float64, float64, error) {
	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return 0, 0, errors.New("arguments must be an object")
	}

	a, err := mc.extractNumberParam(argsMap, "a")
	if err != nil {
		return 0, 0, err
	}

	b, err := mc.extractNumberParam(argsMap, "b")
	if err != nil {
		return 0, 0, err
	}

	return a, b, nil
}

// extractNumberParam extracts and validates a number parameter from args.
func (mc *MathController) extractNumberParam(
	args map[string]interface{},
	paramName string,
) (float64, error) {
	value, exists := args[paramName]
	if !exists {
		return 0, fmt.Errorf("%s parameter is required", paramName)
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("%s parameter must be a number", paramName)
	}
}

// NewTools returns all math tools.
// Satisfies the ToolsFactory interface.
func (mc *MathController) NewTools() []mcpserver.ServerTool {
	return []mcpserver.ServerTool{
		mc.newCalculateServerTool(),
		mc.newAddServerTool(),
		mc.newSubtractServerTool(),
		mc.newMultiplyServerTool(),
		mc.newDivideServerTool(),
	}
}
