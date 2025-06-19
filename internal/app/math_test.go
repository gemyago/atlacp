package app

import (
	"context"
	"testing"

	"github.com/gemyago/golang-backend-boilerplate/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMathServiceDeps() MathServiceDeps {
	return MathServiceDeps{
		RootLogger: diag.RootTestLogger(),
	}
}

func TestNewMathService(t *testing.T) {
	t.Run("should create math service with dependencies", func(t *testing.T) {
		deps := makeMathServiceDeps()

		service := NewMathService(deps)

		require.NotNil(t, service)
		require.NotNil(t, service.logger)
	})
}

func TestMathService_Add(t *testing.T) {
	t.Run("should perform basic addition operation", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 5.0
		b := 3.0

		response, err := service.Add(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 8.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationAdd, response.Operation)
		assert.InEpsilon(t, a, response.A, 0.0001)
		assert.InEpsilon(t, b, response.B, 0.0001)
	})

	t.Run("should handle negative numbers in addition", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := -5.0
		b := 3.0

		response, err := service.Add(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, -2.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationAdd, response.Operation)
	})

	t.Run("should handle decimal numbers in addition", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 2.5
		b := 3.7

		response, err := service.Add(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 6.2, response.Result, 0.001)
		assert.Equal(t, MathOperationAdd, response.Operation)
	})
}

func TestMathService_Subtract(t *testing.T) {
	t.Run("should perform basic subtraction operation", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 8.0
		b := 3.0

		response, err := service.Subtract(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 5.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationSubtract, response.Operation)
		assert.InEpsilon(t, a, response.A, 0.0001)
		assert.InEpsilon(t, b, response.B, 0.0001)
	})

	t.Run("should handle negative result in subtraction", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 3.0
		b := 8.0

		response, err := service.Subtract(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, -5.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationSubtract, response.Operation)
	})

	t.Run("should handle decimal subtraction", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 10.5
		b := 4.3

		response, err := service.Subtract(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 6.2, response.Result, 0.001)
		assert.Equal(t, MathOperationSubtract, response.Operation)
	})
}

func TestMathService_Multiply(t *testing.T) {
	t.Run("should perform basic multiplication operation", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 4.0
		b := 5.0

		response, err := service.Multiply(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 20.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationMultiply, response.Operation)
		assert.InEpsilon(t, a, response.A, 0.0001)
		assert.InEpsilon(t, b, response.B, 0.0001)
	})

	t.Run("should handle multiplication by zero", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 7.0
		b := 0.0

		response, err := service.Multiply(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InDelta(t, 0.0, response.Result, 0)
		assert.Equal(t, MathOperationMultiply, response.Operation)
	})

	t.Run("should handle negative multiplication", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := -3.0
		b := 4.0

		response, err := service.Multiply(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, -12.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationMultiply, response.Operation)
	})

	t.Run("should handle decimal multiplication", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 2.5
		b := 3.2

		response, err := service.Multiply(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 8.0, response.Result, 0.001)
		assert.Equal(t, MathOperationMultiply, response.Operation)
	})
}

func TestMathService_Divide(t *testing.T) {
	t.Run("should perform basic division operation", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 15.0
		b := 3.0

		response, err := service.Divide(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 5.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationDivide, response.Operation)
		assert.InEpsilon(t, a, response.A, 0.0001)
		assert.InEpsilon(t, b, response.B, 0.0001)
	})

	t.Run("should handle division by zero with error", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 10.0
		b := 0.0

		response, err := service.Divide(ctx, a, b)

		require.Error(t, err)
		require.Nil(t, response)
		assert.Contains(t, err.Error(), "division by zero")
		assert.Contains(t, err.Error(), "cannot divide 10.000000 by 0.000000")
	})

	t.Run("should handle decimal division", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 7.5
		b := 2.5

		response, err := service.Divide(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 3.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationDivide, response.Operation)
	})

	t.Run("should handle division resulting in decimal", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := 10.0
		b := 3.0

		response, err := service.Divide(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 3.333333, response.Result, 0.000001)
		assert.Equal(t, MathOperationDivide, response.Operation)
	})

	t.Run("should handle negative division", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		a := -12.0
		b := 4.0

		response, err := service.Divide(ctx, a, b)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, -3.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationDivide, response.Operation)
	})
}

func TestMathService_Calculate(t *testing.T) {
	t.Run("should calculate addition using generic method", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperationAdd,
			A:         7.0,
			B:         3.0,
		}

		response, err := service.Calculate(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 10.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationAdd, response.Operation)
	})

	t.Run("should calculate subtraction using generic method", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperationSubtract,
			A:         10.0,
			B:         4.0,
		}

		response, err := service.Calculate(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 6.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationSubtract, response.Operation)
	})

	t.Run("should calculate multiplication using generic method", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperationMultiply,
			A:         6.0,
			B:         7.0,
		}

		response, err := service.Calculate(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 42.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationMultiply, response.Operation)
	})

	t.Run("should calculate division using generic method", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperationDivide,
			A:         20.0,
			B:         4.0,
		}

		response, err := service.Calculate(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.InEpsilon(t, 5.0, response.Result, 0.0001)
		assert.Equal(t, MathOperationDivide, response.Operation)
	})

	t.Run("should handle division by zero in generic method", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperationDivide,
			A:         15.0,
			B:         0.0,
		}

		response, err := service.Calculate(ctx, request)

		require.Error(t, err)
		require.Nil(t, response)
		assert.Contains(t, err.Error(), "division by zero")
	})

	t.Run("should handle unsupported operation", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx := t.Context()

		request := &MathRequest{
			Operation: MathOperation(faker.Word()),
			A:         5.0,
			B:         3.0,
		}

		response, err := service.Calculate(ctx, request)

		require.Error(t, err)
		require.Nil(t, response)
		assert.Contains(t, err.Error(), "unsupported operation")
	})
}

func TestMathService_ContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation gracefully in all operations", func(t *testing.T) {
		deps := makeMathServiceDeps()
		service := NewMathService(deps)
		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel context immediately

		// Test all operations with cancelled context
		// Math operations are synchronous and don't depend on context for computation
		addResp, addErr := service.Add(ctx, 5.0, 3.0)
		require.NoError(t, addErr)
		require.NotNil(t, addResp)
		assert.InEpsilon(t, 8.0, addResp.Result, 0.0001)

		subResp, subErr := service.Subtract(ctx, 10.0, 3.0)
		require.NoError(t, subErr)
		require.NotNil(t, subResp)
		assert.InEpsilon(t, 7.0, subResp.Result, 0.0001)

		mulResp, mulErr := service.Multiply(ctx, 4.0, 5.0)
		require.NoError(t, mulErr)
		require.NotNil(t, mulResp)
		assert.InEpsilon(t, 20.0, mulResp.Result, 0.0001)

		divResp, divErr := service.Divide(ctx, 15.0, 3.0)
		require.NoError(t, divErr)
		require.NotNil(t, divResp)
		assert.InEpsilon(t, 5.0, divResp.Result, 0.0001)
	})
}
