package repository

import (
	"context"
	"dashboard/internal/models"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseRepo struct {
	db driver.Conn
}

func NewClickhouseRepo(db driver.Conn) *ClickHouseRepo {
	return &ClickHouseRepo{
		db: db,
	}
}

func (c *ClickHouseRepo) GetLogs(ctx context.Context, filter models.LogFilter) ([]*models.LogEntry, error) {
	const op = "repository.GetLogs"

	var logs []*models.LogEntry

	cond, args := buildCondition(filter)

	offsetQuery, limitArgs := buildLimitOffsetQuery(filter)

	args = append(args, limitArgs...)

	query := fmt.Sprintf("SELECT timestamp, level, service_name, message, op FROM %s %s %s;", telemetryTable, cond, offsetQuery)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	for rows.Next() {
		var log models.LogEntry
		if err := rows.Scan(
			&log.TimeStamp,
			&log.Level,
			&log.ServiceName,
			&log.Message,
			&log.Op,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan: %w", op, err)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (c *ClickHouseRepo) GetLogLevelStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LevelCount, error) {
	const op = "repository.GetLogLevelStats"

	var levels []*models.LevelCount

	cond, args := buildTimeRangeCondition(filter)

	query := fmt.Sprintf("SELECT level, count(level) FROM %s %s GROUP BY level", telemetryTable, cond)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	for rows.Next() {
		var levelCount models.LevelCount
		if err := rows.Scan(&levelCount.Level, &levelCount.Count); err != nil {
			return nil, fmt.Errorf("%s: failed to scan rows: %w", op, err)
		}

		levels = append(levels, &levelCount)
	}

	return levels, nil
}

func (c *ClickHouseRepo) GetServiceStats(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ServiceCount, error) {
	const op = "repository.GetServiceStats"

	var services []*models.ServiceCount

	cond, args := buildTimeRangeCondition(filter)

	query := fmt.Sprintf("SELECT service_name, count(service_name) FROM %s %s GROUP BY service_name;", telemetryTable, cond)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	for rows.Next() {
		var service models.ServiceCount
		if err := rows.Scan(&service.ServiceName, &service.Count); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		services = append(services, &service)
	}

	return services, nil
}

func (c *ClickHouseRepo) GetTopErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.ErrorCount, error) {
	const op = "repository.GetTopErrors"

	queryFilter := models.LogFilter{
		Level: errorLevel,
		From:  filter.TimeFrom,
		To:    filter.TimeTo,
		Limit: topErrorsCount,
	}

	cond, args := buildCondition(queryFilter)
	limitCond, limArgs := buildLimitOffsetQuery(queryFilter)
	args = append(args, limArgs...)

	query := fmt.Sprintf("SELECT message, count(message) FROM %s %s GROUP BY message %s", telemetryTable, cond, limitCond)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	var topErrors []*models.ErrorCount
	for rows.Next() {
		var topError models.ErrorCount

		if err := rows.Scan(&topError.Message, &topError.Count); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		topErrors = append(topErrors, &topError)
	}

	return topErrors, nil
}

func (c *ClickHouseRepo) GetRecentErrors(ctx context.Context, filter models.TimeRangeFilter) ([]*models.LogEntry, error) {
	const op = "repository.GetRecentErrors"

	errFilter := models.LogFilter{
		Level: errorLevel,
		From:  filter.TimeFrom,
		To:    filter.TimeTo,
	}

	cond, args := buildCondition(errFilter)

	query := fmt.Sprintf("SELECT level, timestamp, service_name, message, op FROM %s %s", telemetryTable, cond)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	var logs []*models.LogEntry

	for rows.Next() {
		var log models.LogEntry

		if err := rows.Scan(&log.Level, &log.TimeStamp, &log.ServiceName, &log.Message, &log.Op); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (c *ClickHouseRepo) GetAllLogLevels(ctx context.Context) ([]string, error) {
	const op = "repository.GetAllLogLevels"

	var levels []string

	query := fmt.Sprintf("SELECT level FROM %s GROUP BY level", telemetryTable)
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	for rows.Next() {
		var level string

		if err := rows.Scan(&level); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		levels = append(levels, level)
	}

	return levels, nil
}

func (c *ClickHouseRepo) GetAllServiceNames(ctx context.Context) ([]string, error) {
	const op = "repository.GetAllServiceNames"

	var services []string

	query := fmt.Sprintf("SELECT service_name FROM %s GROUP BY service_name", telemetryTable)
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	for rows.Next() {
		var service string

		if err := rows.Scan(&service); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		services = append(services, service)
	}

	return services, nil
}
