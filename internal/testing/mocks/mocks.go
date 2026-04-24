//go:build !release

package mocks

import (
	"testing"
)

// GetMock is a helper function to get a mock from a given instance.
// Note: This should only be used internally in tests.
//
//nolint:ireturn
func GetMock[TOutput any](t *testing.T, input any) TOutput {
	mock, ok := input.(TOutput)
	if !ok {
		t.Fatalf("input is not a %T", input)
	}
	return mock
}
