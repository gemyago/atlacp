package app

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTimeServiceDeps() TimeServiceDeps {
	return TimeServiceDeps{
		RootLogger: diag.RootTestLogger(),
	}
}

func TestNewTimeService(t *testing.T) {
	t.Run("should create time service with dependencies", func(t *testing.T) {
		deps := makeTimeServiceDeps()

		service := NewTimeService(deps)

		require.NotNil(t, service)
		require.NotNil(t, service.logger)
	})
}

func TestTimeService_GetCurrentTime_ISO(t *testing.T) {
	t.Run("should return current time in ISO format", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		request := &TimeRequest{Format: TimeFormatISO}

		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "iso", response.Format)
		assert.NotEmpty(t, response.Time)

		// Verify the time is in ISO format (can be parsed as RFC3339)
		_, parseErr := time.Parse(time.RFC3339, response.Time)
		assert.NoError(t, parseErr, "Time should be in valid ISO/RFC3339 format")
	})

	t.Run("should default to ISO format when format not specified", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		// Empty format should default to ISO
		request := &TimeRequest{Format: ""}

		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "iso", response.Format)
		assert.NotEmpty(t, response.Time)
	})
}

func TestTimeService_GetCurrentTime_RFC3339(t *testing.T) {
	t.Run("should return current time in RFC3339 format", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		request := &TimeRequest{Format: TimeFormatRFC3339}

		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "rfc3339", response.Format)
		assert.NotEmpty(t, response.Time)

		// Verify the time is in RFC3339 format
		_, parseErr := time.Parse(time.RFC3339, response.Time)
		assert.NoError(t, parseErr, "Time should be in valid RFC3339 format")
	})
}

func TestTimeService_GetCurrentTime_Unix(t *testing.T) {
	t.Run("should return current time in Unix format", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		request := &TimeRequest{Format: TimeFormatUnix}

		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "unix", response.Format)
		assert.NotEmpty(t, response.Time)

		// Verify the time is a valid Unix timestamp (parsable as integer)
		_, parseErr := strconv.ParseInt(response.Time, 10, 64)
		assert.NoError(t, parseErr, "Time should be a valid Unix timestamp")
	})
}

func TestTimeService_GetCurrentTime_InvalidFormat(t *testing.T) {
	t.Run("should default to ISO format for invalid format", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		// Use an invalid format
		request := &TimeRequest{Format: TimeFormat(faker.Word())}

		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "iso", response.Format) // Should default to ISO
		assert.NotEmpty(t, response.Time)
	})
}

func TestTimeService_GetCurrentTime_ContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel context immediately

		request := &TimeRequest{Format: TimeFormatISO}

		// Service should still work as it doesn't depend on context for time operations
		response, err := service.GetCurrentTime(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "iso", response.Format)
		assert.NotEmpty(t, response.Time)
	})
}

func TestTimeService_GetCurrentTime_TimeAccuracy(t *testing.T) {
	t.Run("should return recent time", func(t *testing.T) {
		deps := makeTimeServiceDeps()
		service := NewTimeService(deps)
		ctx := t.Context()

		before := time.Now()
		request := &TimeRequest{Format: TimeFormatISO}

		response, err := service.GetCurrentTime(ctx, request)
		after := time.Now()

		require.NoError(t, err)
		require.NotNil(t, response)

		// Parse the returned time and verify it's between before and after
		parsedTime, parseErr := time.Parse(time.RFC3339, response.Time)
		require.NoError(t, parseErr)

		assert.True(t, parsedTime.After(before.Add(-time.Second)), "Returned time should be after test start")
		assert.True(t, parsedTime.Before(after.Add(time.Second)), "Returned time should be before test end")
	})
}
