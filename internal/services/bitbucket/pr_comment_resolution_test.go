package bitbucket

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvedStateFromResolutionJSON(t *testing.T) {
	t.Run("nil or empty raw is unknown", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(nil)
		assert.False(t, resolved)
		assert.False(t, known)

		resolved, known = ResolvedStateFromResolutionJSON(json.RawMessage{})
		assert.False(t, resolved)
		assert.False(t, known)
	})

	t.Run("JSON null is known not resolved", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage("null"))
		assert.False(t, resolved)
		assert.True(t, known)
	})

	t.Run("whitespace around null", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage("  null  "))
		assert.False(t, resolved)
		assert.True(t, known)
	})

	t.Run("empty object is ambiguous", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage("{}"))
		assert.False(t, resolved)
		assert.False(t, known)
	})

	t.Run("whitespace around empty object", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage("  {}  "))
		assert.False(t, resolved)
		assert.False(t, known)
	})

	t.Run("resolved true", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage(`{"resolved":true}`))
		assert.True(t, resolved)
		assert.True(t, known)
	})

	t.Run("resolved false", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage(`{"resolved":false}`))
		assert.False(t, resolved)
		assert.True(t, known)
	})

	t.Run("resolved_by implies resolved", func(t *testing.T) {
		raw := json.RawMessage(`{"resolved_by":{"type":"user","display_name":"Alice"}}`)
		resolved, known := ResolvedStateFromResolutionJSON(raw)
		assert.True(t, resolved)
		assert.True(t, known)
	})

	t.Run("resolved_on implies resolved", func(t *testing.T) {
		raw := json.RawMessage(`{"resolved_on":"2024-01-02T15:04:05Z"}`)
		resolved, known := ResolvedStateFromResolutionJSON(raw)
		assert.True(t, resolved)
		assert.True(t, known)
	})

	t.Run("resolved_on whitespace only is not metadata", func(t *testing.T) {
		raw := json.RawMessage(`{"resolved_on":"   "}`)
		resolved, known := ResolvedStateFromResolutionJSON(raw)
		assert.False(t, resolved)
		assert.True(t, known)
	})

	t.Run("invalid JSON is unknown", func(t *testing.T) {
		resolved, known := ResolvedStateFromResolutionJSON(json.RawMessage(`{`))
		assert.False(t, resolved)
		assert.False(t, known)
	})

	t.Run("non empty object without resolution signals is known not resolved", func(t *testing.T) {
		raw := json.RawMessage(`{"other":"field"}`)
		resolved, known := ResolvedStateFromResolutionJSON(raw)
		assert.False(t, resolved)
		assert.True(t, known)
	})
}
