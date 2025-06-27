package repository

import (
	"collector-service/internal/models"
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Repository interface {
	InsertBatch(ctx context.Context, records []models.TelemetryRecord) error
	Stop() error
}

type ClickHouseRepository struct {
	db driver.Conn
}

func NewClickHouseRepository(db driver.Conn) *ClickHouseRepository {
	return &ClickHouseRepository{db: db}
}

func (c *ClickHouseRepository) InsertBatch(ctx context.Context, logEntry []models.TelemetryRecord) error {
	const op = "repository.ClickHouseRepository.InsertBatch"

	if len(logEntry) == 0 {
		return nil // No records to insert
	}

	batch, err := c.db.PrepareBatch(ctx, `INSERT INTO telemetry(timestamp, level, service_name, message, op) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare batch: %w", op, err)
	}
	defer batch.Close()

	for _, log := range logEntry {
		if err := batch.Append(
			log.Timestamp,
			log.Level,
			log.ServiceName,
			log.Message,
			log.Op,
		); err != nil {
			// Log the error with the service name for better debugging
			return fmt.Errorf("%s: failed to append log entry in batch from service: %s: %w", op, log.ServiceName, err)
		}
	}

	return batch.Send()
}

func (c *ClickHouseRepository) Stop() error {
	if err := c.db.Close(); err != nil {
		return err
	}

	return nil
}
