package grpcapp

import (
	"collector-service/internal/grpc/collector"
	"collector-service/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net"

	_ "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	service    service.LogService
}

func NewApp(
	log *slog.Logger,
	service service.LogService,
	port int,
) *App {

	zapLogger, _ := zap.NewDevelopment()

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(
			func(p interface{}) error {
				log.Error("recovered from panic", slog.Any("panic", p))

				return status.Errorf(codes.Internal, "internal error")
			},
		),
	}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(zapLogger),
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(zapLogger),
			recovery.StreamServerInterceptor(recoveryOpts...),
			logging.StreamServerInterceptor(InterceptorLogger(log), loggingOpts...),
		),
	)

	collector.Register(gRPCServer, service, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		service:    service,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcapp.Run"

	errChan := make(chan error)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: failed to create listener: %w", op, err)
	}

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			errChan <- fmt.Errorf("%s: error while running server: %w", op, err)
		}
	}()

	select {
	case <-ctx.Done():
		a.Stop()
		return nil
	case err := <-errChan:
		return err
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.Info("stopping gRPC server", slog.String("op", op), slog.Int("port", a.port))

	doneCh := make(chan struct{})

	go func() {
		if err := a.service.Stop(); err != nil {
			a.log.Error("failed to stop service", slog.String("op", op), slog.Any("error", err))
		}
		close(doneCh)
	}()

	a.gRPCServer.GracefulStop()
	<-doneCh
}
