package tests

import (
	"collector-service/internal/logagent"
	"context"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/brianvoe/gofakeit"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	serverAddress      = "localhost:9001"
	contextTimeout     = time.Second * 10
	ackStatus          = "ok"
	logBatches         = 4
	rowsCount          = 2000
	clickhouseAddr     = "localhost:9000"
	clickhouseDatabase = "telemetry_db"
	clickhouseUser     = "telemetry_user"
	clickhousePassword = "telemetry_pass"
)

func generateRandomLogs() []*logagent.LogEntry {
	var logs []*logagent.LogEntry

	for len(logs) < 500 {
		var log logagent.LogEntry

		log.Timestamp = timestamppb.Now()
		log.Level = logagent.LogLevel(gofakeit.Number(0, 4))
		log.Message = gofakeit.Sentence(4)
		log.Op = gofakeit.BeerStyle()
		log.ServiceName = gofakeit.FirstName()

		logs = append(logs, &log)
	}

	return logs
}

func TestSuccessLogStreaming_EndToEnd(t *testing.T) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	testContext, cancel := context.WithDeadline(context.Background(), time.Now().Add(contextTimeout))
	defer cancel()

	client := logagent.NewLogCollectorClient(conn)

	stream, err := client.SendLogs(testContext)
	require.NoError(t, err)

	for _ = range logBatches {
		batch := generateRandomLogs()

		for _, log := range batch {
			err := stream.Send(log)
			require.NoError(t, err)
		}
	}

	ack, err := stream.CloseAndRecv()
	require.NoError(t, err)

	require.Equal(t, ackStatus, ack.Status)

	clickhouseConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseAddr},
		Auth: clickhouse.Auth{
			Database: clickhouseDatabase,
			Username: clickhouseUser,
			Password: clickhousePassword,
		},
		TLS: nil,
	})
	require.NoError(t, err)

	var count uint64

	err = clickhouseConn.QueryRow(testContext, "SELECT count(*) FROM telemetry;").Scan(&count)
	require.NoError(t, err)

	rows, err := clickhouseConn.Query(testContext, "SELECT level, message, service_name FROM telemetry LIMIT 10;")
	require.NoError(t, err)
	for rows.Next() {
		var (
			level       string
			message     string
			serviceName string
		)

		err = rows.Scan(&level, &message, &serviceName)
		require.NoError(t, err)

		t.Logf("received log parts: level=%s, message=%s, service_name=%s", level, message, serviceName)
	}

}
