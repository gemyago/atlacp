package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&buf)

		err := writeJSON(cmd, map[string]string{"hello": "world"})
		require.NoError(t, err)

		out := buf.String()
		assert.True(t, strings.HasPrefix(out, "{\n"))
		assert.Contains(t, out, `"hello": "world"`)
		assert.True(t, strings.HasSuffix(out, "}\n"))
	})

	t.Run("marshal error", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		var ch chan int
		err := writeJSON(cmd, ch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "marshal JSON")
	})

	t.Run("write error", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		cmd.SetOut(errWriter{err: errors.New("disk full")})

		err := writeJSON(cmd, map[string]int{"a": 1})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "write JSON")
	})
}

type errWriter struct {
	err error
}

func (w errWriter) Write(p []byte) (int, error) {
	return 0, w.err
}
