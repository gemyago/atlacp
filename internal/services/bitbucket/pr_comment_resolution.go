package bitbucket

import (
	"encoding/json"
	"strings"
)

// ResolvedStateFromResolutionJSON derives whether a PR comment is resolved from Bitbucket's
// resolution JSON. When the API returns an empty object ("{}") or omits resolution, known is false
// and callers should fetch the single comment if they need a definitive resolved flag.
func ResolvedStateFromResolutionJSON(raw json.RawMessage) (bool, bool) {
	if len(raw) == 0 {
		return false, false
	}
	s := strings.TrimSpace(string(raw))
	if s == "null" {
		return false, true
	}
	if s == "{}" {
		return false, false
	}
	var r struct {
		Resolved   *bool    `json:"resolved,omitempty"`
		ResolvedBy *Account `json:"resolved_by,omitempty"`
		ResolvedOn string   `json:"resolved_on,omitempty"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return false, false
	}
	if r.Resolved != nil {
		return *r.Resolved, true
	}
	if r.ResolvedBy != nil || strings.TrimSpace(r.ResolvedOn) != "" {
		return true, true
	}
	return false, true
}
