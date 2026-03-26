package bitbucket

import (
	"encoding/json"
	"strings"
)

// ResolvedStateFromResolutionJSON derives whether a PR comment is resolved from Bitbucket's
// resolution JSON. It returns true when raw is a JSON object (including "{}" or any object with
// keys). It returns false for null, non-objects, invalid JSON, or empty input.
func ResolvedStateFromResolutionJSON(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	s := strings.TrimSpace(string(raw))
	if s == "" {
		return false
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil || m == nil {
		return false
	}
	return true
}
