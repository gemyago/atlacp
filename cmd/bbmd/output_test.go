package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		key := faker.Word()
		val := faker.Sentence()
		input := map[string]string{key: val}

		var buf bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&buf)

		err := writeJSON(cmd, input)
		require.NoError(t, err)

		var got map[string]string
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		assert.Equal(t, input, got)
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
		cmd.SetOut(errWriter{err: errors.New(faker.Sentence())})

		err := writeJSON(cmd, map[string]int{faker.Word(): int(faker.RandomUnixTime() % 1000)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "write JSON")
	})
}

func TestExecAndWrite(t *testing.T) {
	t.Parallel()

	t.Run("noop skips target and output", func(t *testing.T) {
		t.Parallel()
		param := faker.Sentence()
		var outBuf bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetOut(&outBuf)
		rootParams := &rootCommandParams{Noop: true}
		called := false
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execAndWrite(cmd, deps, execArgs[string, string]{
			rootParams: rootParams,
			params:     param,
			target: func(_ context.Context, _ string) (string, error) {
				called = true
				return faker.Word(), nil
			},
		})
		require.NoError(t, err)
		assert.False(t, called)
		assert.Empty(t, outBuf.String())
	})

	t.Run("success writes JSON", func(t *testing.T) {
		t.Parallel()
		n := int(faker.RandomUnixTime() % 10000)
		key := faker.Word()
		want := map[string]int{key: n}

		var buf bytes.Buffer
		cmd := &cobra.Command{}
		cmd.SetContext(t.Context())
		cmd.SetOut(&buf)
		rootParams := &rootCommandParams{}
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execAndWrite(cmd, deps, execArgs[int, map[string]int]{
			rootParams: rootParams,
			params:     n,
			target: func(ctx context.Context, p int) (map[string]int, error) {
				assert.Equal(t, n, p)
				assert.Equal(t, cmd.Context(), ctx)
				return want, nil
			},
		})
		require.NoError(t, err)

		var got map[string]int
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		assert.Equal(t, want, got)
	})

	t.Run("target error", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		rootParams := &rootCommandParams{}
		wantErr := errors.New(faker.Sentence())
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execAndWrite(cmd, deps, execArgs[struct{}, int]{
			rootParams: rootParams,
			params:     struct{}{},
			target: func(_ context.Context, _ struct{}) (int, error) {
				return 0, wantErr
			},
		})
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("write JSON error", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		cmd.SetOut(errWriter{err: errors.New(faker.Sentence())})
		rootParams := &rootCommandParams{}
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execAndWrite(cmd, deps, execArgs[struct{}, map[string]int]{
			rootParams: rootParams,
			params:     struct{}{},
			target: func(_ context.Context, _ struct{}) (map[string]int, error) {
				return map[string]int{faker.Word(): int(faker.RandomUnixTime() % 1000)}, nil
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "write JSON")
	})

	t.Run("marshal JSON error", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		rootParams := &rootCommandParams{}
		var ch chan int
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execAndWrite(cmd, deps, execArgs[struct{}, chan int]{
			rootParams: rootParams,
			params:     struct{}{},
			target: func(_ context.Context, _ struct{}) (chan int, error) {
				return ch, nil
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "marshal JSON")
	})
}

func TestExecNoOutput(t *testing.T) {
	t.Parallel()

	t.Run("noop skips run", func(t *testing.T) {
		t.Parallel()
		param := faker.Sentence()
		cmd := &cobra.Command{}
		called := false
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execNoOutput(cmd, deps, &rootCommandParams{Noop: true}, param, func(_ context.Context, _ string) error {
			called = true
			return errors.New(faker.Sentence())
		})
		require.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("success runs handler", func(t *testing.T) {
		t.Parallel()
		want := faker.Sentence()
		cmd := &cobra.Command{}
		cmd.SetContext(t.Context())
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		var got string
		err := execNoOutput(cmd, deps, &rootCommandParams{}, want, func(ctx context.Context, p string) error {
			assert.Equal(t, want, p)
			assert.Equal(t, cmd.Context(), ctx)
			got = p
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("handler error propagates", func(t *testing.T) {
		t.Parallel()
		cmd := &cobra.Command{}
		wantErr := errors.New(faker.Sentence())
		deps := execDeps{RootLogger: diag.RootTestLogger()}
		err := execNoOutput(cmd, deps, &rootCommandParams{}, struct{}{}, func(_ context.Context, _ struct{}) error {
			return wantErr
		})
		require.ErrorIs(t, err, wantErr)
	})
}

type errWriter struct {
	err error
}

func (w errWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}
