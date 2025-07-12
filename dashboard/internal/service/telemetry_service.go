package service

import (
	"context"
	"dashboard/internal/cache"
	"dashboard/internal/models"
	"dashboard/internal/repository"
	"fmt"
)

const (
	serviceNamesKey = "service_names"
	logLevelsKey    = "log_levels"
)

type TelemetryService struct {
	repo  repository.Repository
	cache cache.TTLCache
}

func NewTelemetryService(repo repository.Repository, cache cache.TTLCache) *TelemetryService {
	return &TelemetryService{
		repo:  repo,
		cache: cache,
	}
}

func (ts *TelemetryService) GetLogs(ctx context.Context, filter models.LogFilter) ([]*models.LogEntry, error) {
	const op = "service.GetLogs"

	logs, err := ts.repo.GetLogs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get logs: %w", op, err)
	}

	return logs, nil
}

func (ts *TelemetryService) GetLogLevelStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LevelCount, error) {
	const op = "service.GetLogLevelStats"

	stats, err := ts.repo.GetLogLevelStats(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get log level stats: %w", op, err)
	}

	return stats, nil
}

func (ts *TelemetryService) GetServiceStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ServiceCount, error) {
	const op = "service.GetServiceStats"

	stats, err := ts.repo.GetServiceStats(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get service stats: %w", op, err)
	}

	return stats, nil
}

func (ts *TelemetryService) GetTopErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ErrorCount, error) {
	const op = "service.GetTopErrors"

	topErrors, err := ts.repo.GetTopErrors(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get top errors: %w", op, err)
	}

	return topErrors, nil
}

func (ts *TelemetryService) GetRecentErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LogEntry, error) {
	const op = "service.GetRecentErrors"

	recentErrors, err := ts.repo.GetRecentErrors(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get recent errors: %w", op, err)
	}

	return recentErrors, nil
}

func (ts *TelemetryService) GetAllServiceNames(ctx context.Context, filter models.TimeRangeFilter) ([]string, error) {
	const op = "service.GetAllServiceNames"

	if names, exists := ts.cache.Get(ctx, serviceNamesKey); exists {
		cNames, _ := names.Value.([]string)
		return cNames, nil
	}

	names, err := ts.repo.GetAllServiceNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get service names: %w", op, err)
	}

	ts.cache.Set(ctx, serviceNamesKey, names)

	return names, nil
}

func (ts *TelemetryService) GetAllLogLevels(ctx context.Context) ([]string, error) {
	const op = "service.GetAllLogLevels"

	if levels, exists := ts.cache.Get(ctx, logLevelsKey); exists {
		cLevels, _ := levels.Value.([]string)
		return cLevels, nil
	}

	levels, err := ts.repo.GetAllLogLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get log levels: %w", op, err)
	}

	ts.cache.Set(ctx, logLevelsKey, levels)

	return levels, nil
}
