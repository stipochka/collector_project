package migrate

import (
	"context"
	"fmt"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func Migrate(db driver.Conn, path string) error {
	const op = "migrate.Migrate"

	query, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Exec(context.Background(), string(query)); err != nil {
		return fmt.Errorf("%s: failed to execute migration query: %w", op, err)
	}

	return nil
}
