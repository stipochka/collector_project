package database

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectToDB(addr, database, user, password string) (*driver.Conn, error) {
	const op = "database.ConnectToDB"

	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true, // Set to true for self-signed certificates; false for production,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to database: %w", op, err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return &db, nil
}
