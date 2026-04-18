package main

import (
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestApplyPRCommentsResolvedFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil unchanged", func(t *testing.T) {
		t.Parallel()
		applyPRCommentsResolvedFilter(nil, false)
	})

	t.Run("include resolved leaves values", func(t *testing.T) {
		t.Parallel()
		result := &app.BitbucketListPRCommentsResult{
			Values: []app.BitbucketPRComment{{Resolved: true}, {Resolved: false}},
		}
		applyPRCommentsResolvedFilter(result, true)
		assert.Len(t, result.Values, 2)
	})

	t.Run("filters resolved when false", func(t *testing.T) {
		t.Parallel()
		result := &app.BitbucketListPRCommentsResult{
			Values: []app.BitbucketPRComment{
				{ID: 1, Resolved: true},
				{ID: 2, Resolved: false},
				{ID: 3, Resolved: true},
			},
		}
		applyPRCommentsResolvedFilter(result, false)
		assert.Len(t, result.Values, 1)
		assert.Equal(t, int64(2), result.Values[0].ID)
	})
}
