package repository

import (
	"context"
	"dashboard/internal/models"
)

const (
	telemetryTable = "telemetry"
	topErrorsCount = 5
	errorLevel     = "ERROR"
)

type Repository interface {
	GetLogs(ctx context.Context, filter models.LogFilter) ([]*models.LogEntry, error)

	//aggregations
	GetLogLevelStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LevelCount, error)
	GetServiceStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ServiceCount, error)
	GetTopErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ErrorCount, error)
	GetRecentErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LogEntry, error)

	//for ui filters
	GetAllLogLevels(ctx context.Context) ([]string, error)
	GetAllServiceNames(ctx context.Context) ([]string, error)
}
