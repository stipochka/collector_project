package collector

import (
	"collector-service/internal/logagent"
	"collector-service/internal/models"
	"collector-service/internal/service"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const maxBatchSize = 1000 // Maximum number of log entries per batch

const flushTimeout = time.Second * 3 // flush logs timeout

type serverApi struct {
	logagent.UnimplementedLogCollectorServer
	service service.LogService
	log     *slog.Logger
}

func Register(gRPCServer *grpc.Server, service service.LogService, log *slog.Logger) {
	logagent.RegisterLogCollectorServer(gRPCServer, &serverApi{service: service, log: log})
}

func (s *serverApi) SendLogs(stream logagent.LogCollector_SendLogsServer) error {
	const op = "collector.serverApi.SendLogs"

	log := s.log.With(slog.String("op", op))

	log.Info("started logs stream")

	ctx := stream.Context()

	var (
		batch []models.TelemetryRecord
		mu    sync.Mutex
	)

	flush := func() error {
		mu.Lock()
		defer mu.Unlock()

		if len(batch) == 0 {
			return nil
		}

		if err := s.service.StoreBatch(ctx, batch); err != nil {
			log.Error("failed to insert batch", slog.Any("error", err))
			return err
		}

		batch = batch[:0]
		log.Info("successfully inserted batch")

		return nil
	}

	flushTimer := time.NewTimer(flushTimeout)
	defer flushTimer.Stop()

	stopCh := make(chan struct{})

	resetCh := make(chan struct{})

	// background flusher
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-flushTimer.C:
				_ = flush()
				flushTimer.Reset(flushTimeout)
			case <-stopCh:
				return
			case <-resetCh:
				if !flushTimer.Stop() {
					<-flushTimer.C
				}
				flushTimer.Reset(flushTimeout)
			}
		}
	}()

	for {
		entry, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Info("stream ended")
				close(stopCh)
				if err = flush(); err != nil {
					return err
				}

				return stream.SendAndClose(&logagent.Ack{Status: "ok"})
			}

			log.Error("failed to receive message", slog.Any("error", err))

			return fmt.Errorf("%s: failed to receive log: %w", op, err)
		}

		record := models.TelemetryRecord{
			Timestamp:   entry.GetTimestamp().AsTime(),
			Level:       entry.GetLevel().String(),
			ServiceName: entry.GetServiceName(),
			Message:     entry.GetMessage(),
			Op:          entry.GetOp(),
		}
		mu.Lock()
		batch = append(batch, record)
		shouldFlush := len(batch) >= maxBatchSize
		mu.Unlock()

		// flush the batch if its reaches maximum size
		if shouldFlush {
			if err := flush(); err != nil {
				close(stopCh)
				return err
			}
			select {
			case resetCh <- struct{}{}:
			default:
			}
		}

	}

}
