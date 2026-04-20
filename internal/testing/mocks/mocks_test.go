package mocks

import "testing"

func TestGetMock(t *testing.T) {
	t.Run("returns typed mock", func(t *testing.T) {
		type mockType struct {
			name string
		}

		expected := &mockType{name: "mock"}
		actual := GetMock[*mockType](t, expected)

		if actual != expected {
			t.Fatalf("unexpected mock: got %v, want %v", actual, expected)
		}
	})
}
