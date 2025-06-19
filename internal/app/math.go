package app

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/dig"
)

// Math service provides mathematical operations for MCP tools.

// MathOperation represents supported mathematical operations.
type MathOperation string

const (
	MathOperationAdd      MathOperation = "add"
	MathOperationSubtract MathOperation = "subtract"
	MathOperationMultiply MathOperation = "multiply"
	MathOperationDivide   MathOperation = "divide"
)

// MathRequest represents a request for mathematical operation.
type MathRequest struct {
	Operation MathOperation `json:"operation"`
	A         float64       `json:"a"`
	B         float64       `json:"b"`
}

// MathResponse represents a mathematical operation response.
type MathResponse struct {
	Result    float64       `json:"result"`
	Operation MathOperation `json:"operation"`
	A         float64       `json:"a"`
	B         float64       `json:"b"`
}

// MathServiceDeps contains dependencies for the math service.
type MathServiceDeps struct {
	dig.In

	RootLogger *slog.Logger
}

// MathService provides mathematical operations.
type MathService struct {
	logger *slog.Logger
}

// Add performs addition operation.
func (svc *MathService) Add(ctx context.Context, a, b float64) (*MathResponse, error) {
	svc.logger.InfoContext(ctx, "Performing addition operation",
		slog.Float64("a", a),
		slog.Float64("b", b))

	result := a + b

	response := &MathResponse{
		Result:    result,
		Operation: MathOperationAdd,
		A:         a,
		B:         b,
	}

	svc.logger.InfoContext(ctx, "Addition operation completed",
		slog.Float64("result", result))

	return response, nil
}

// Subtract performs subtraction operation.
func (svc *MathService) Subtract(ctx context.Context, a, b float64) (*MathResponse, error) {
	svc.logger.InfoContext(ctx, "Performing subtraction operation",
		slog.Float64("a", a),
		slog.Float64("b", b))

	result := a - b

	response := &MathResponse{
		Result:    result,
		Operation: MathOperationSubtract,
		A:         a,
		B:         b,
	}

	svc.logger.InfoContext(ctx, "Subtraction operation completed",
		slog.Float64("result", result))

	return response, nil
}

// Multiply performs multiplication operation.
func (svc *MathService) Multiply(ctx context.Context, a, b float64) (*MathResponse, error) {
	svc.logger.InfoContext(ctx, "Performing multiplication operation",
		slog.Float64("a", a),
		slog.Float64("b", b))

	result := a * b

	response := &MathResponse{
		Result:    result,
		Operation: MathOperationMultiply,
		A:         a,
		B:         b,
	}

	svc.logger.InfoContext(ctx, "Multiplication operation completed",
		slog.Float64("result", result))

	return response, nil
}

// Divide performs division operation with zero-division protection.
func (svc *MathService) Divide(ctx context.Context, a, b float64) (*MathResponse, error) {
	svc.logger.InfoContext(ctx, "Performing division operation",
		slog.Float64("a", a),
		slog.Float64("b", b))

	if b == 0 {
		svc.logger.WarnContext(ctx, "Division by zero attempted",
			slog.Float64("a", a),
			slog.Float64("b", b))
		return nil, fmt.Errorf("division by zero: cannot divide %f by %f", a, b)
	}

	result := a / b

	response := &MathResponse{
		Result:    result,
		Operation: MathOperationDivide,
		A:         a,
		B:         b,
	}

	svc.logger.InfoContext(ctx, "Division operation completed",
		slog.Float64("result", result))

	return response, nil
}

// Calculate performs any mathematical operation based on the request.
func (svc *MathService) Calculate(ctx context.Context, req *MathRequest) (*MathResponse, error) {
	svc.logger.InfoContext(ctx, "Processing math calculation request",
		slog.String("operation", string(req.Operation)),
		slog.Float64("a", req.A),
		slog.Float64("b", req.B))

	switch req.Operation {
	case MathOperationAdd:
		return svc.Add(ctx, req.A, req.B)
	case MathOperationSubtract:
		return svc.Subtract(ctx, req.A, req.B)
	case MathOperationMultiply:
		return svc.Multiply(ctx, req.A, req.B)
	case MathOperationDivide:
		return svc.Divide(ctx, req.A, req.B)
	default:
		svc.logger.ErrorContext(ctx, "Unsupported math operation",
			slog.String("operation", string(req.Operation)))
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}
}

// NewMathService creates a new math service instance.
func NewMathService(deps MathServiceDeps) *MathService {
	return &MathService{
		logger: deps.RootLogger.WithGroup("app.math-service"),
	}
}
