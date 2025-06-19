package app

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"go.uber.org/dig"
)

// Time service provides time-related functionality for MCP tools.

// TimeFormat represents supported time formats.
type TimeFormat string

const (
	TimeFormatISO     TimeFormat = "iso"
	TimeFormatRFC3339 TimeFormat = "rfc3339"
	TimeFormatUnix    TimeFormat = "unix"
)

// TimeRequest represents a request for current time.
type TimeRequest struct {
	Format TimeFormat `json:"format,omitempty"`
}

// TimeResponse represents a time response.
type TimeResponse struct {
	Time   string `json:"time"`
	Format string `json:"format"`
}

// TimeServiceDeps contains dependencies for the time service.
type TimeServiceDeps struct {
	dig.In

	RootLogger *slog.Logger
}

// TimeService provides time operations.
type TimeService struct {
	logger *slog.Logger
}

// GetCurrentTime returns the current time in the specified format.
func (svc *TimeService) GetCurrentTime(ctx context.Context, req *TimeRequest) (*TimeResponse, error) {
	svc.logger.InfoContext(ctx, "Getting current time", slog.String("format", string(req.Format)))

	now := time.Now()
	var timeStr string
	format := req.Format

	// Default to ISO format if not specified
	if format == "" {
		format = TimeFormatISO
	}

	switch format {
	case TimeFormatISO:
		timeStr = now.Format(time.RFC3339)
	case TimeFormatRFC3339:
		timeStr = now.Format(time.RFC3339)
	case TimeFormatUnix:
		timeStr = strconv.FormatInt(now.Unix(), 10)
	default:
		timeStr = now.Format(time.RFC3339)
		format = TimeFormatISO
	}

	response := &TimeResponse{
		Time:   timeStr,
		Format: string(format),
	}

	svc.logger.InfoContext(ctx, "Current time retrieved",
		slog.String("time", response.Time),
		slog.String("format", response.Format))

	return response, nil
}

// NewTimeService creates a new time service instance.
func NewTimeService(deps TimeServiceDeps) *TimeService {
	return &TimeService{
		logger: deps.RootLogger.WithGroup("app.time-service"),
	}
}
