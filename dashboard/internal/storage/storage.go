package storage

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func NewStorage(addr, database, user, password string) (*driver.Conn, error) {
	const op = "storage.NewStorage"

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		TLS: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open connect: %w", op, err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: failed to ping db: %w", op, err)
	}

	return &conn, nil
}
