package service

import (
	"context"
	"dashboard/internal/models"
)

type Service interface {
	GetLogs(ctx context.Context, filter models.LogFilter) ([]*models.LogEntry, error)

	GetLogLevelStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LevelCount, error)
	GetServiceStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ServiceCount, error)
	GetTopErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ErrorCount, error)
	GetRecentErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LogEntry, error)

	GetAllLogLevels(ctx context.Context) ([]string, error)
	GetAllServiceNames(ctx context.Context) ([]string, error)
}
