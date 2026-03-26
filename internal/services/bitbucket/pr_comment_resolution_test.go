package bitbucket

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvedStateFromResolutionJSON(t *testing.T) {
	tests := []struct {
		name string
		raw  json.RawMessage
		want bool
	}{
		{name: "nil", raw: nil, want: false},
		{name: "empty_raw", raw: json.RawMessage{}, want: false},
		{name: "whitespace_only", raw: json.RawMessage("   \t  "), want: false},
		{name: "null", raw: json.RawMessage("null"), want: false},
		{name: "null_with_padding", raw: json.RawMessage("  null  "), want: false},
		{name: "invalid_json", raw: json.RawMessage(`{`), want: false},
		{name: "array", raw: json.RawMessage(`[]`), want: false},
		{name: "string_primitive", raw: json.RawMessage(`"x"`), want: false},
		{name: "empty_object", raw: json.RawMessage("{}"), want: true},
		{name: "empty_object_with_padding", raw: json.RawMessage("  {}  "), want: true},
		{name: "non_empty_object", raw: json.RawMessage(`{"resolved":false}`), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ResolvedStateFromResolutionJSON(tt.raw))
		})
	}
}
