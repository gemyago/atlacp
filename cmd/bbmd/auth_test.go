package main

import (
	"testing"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestRedactTokenValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "***"},
		{name: "short eight", in: "12345678", want: "***"},
		{name: "nine chars", in: "123456789", want: "1234***6789"},
		{name: "long token", in: "abcdefghijklmnop", want: "abcd***mnop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, redactTokenValue(tt.in))
		})
	}
}

func TestRedactAccountForStatus(t *testing.T) {
	t.Parallel()
	acc := app.AtlassianAccount{
		Name:    "a1",
		Default: true,
		Bitbucket: &app.AtlassianToken{
			Type:  "Bearer",
			Value: "abcdefghijklmnop",
		},
		Jira: &app.AtlassianToken{
			Type:  "Bearer",
			Value: "short",
		},
	}
	got := redactAccountForStatus(acc)
	assert.Equal(t, "abcd***mnop", got.Bitbucket.Value)
	assert.Equal(t, "***", got.Jira.Value)
	assert.Equal(t, "a1", got.Name)
	assert.True(t, got.Default)
}
