package app

import (
	grpcapp "collector-service/internal/app/grpc"
	"collector-service/internal/repository"
	"collector-service/internal/service"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, port int, repo repository.Repository) *App {
	service := service.NewLogService(repo)

	srv := grpcapp.NewApp(log, service, port)

	return &App{
		GRPCServer: srv,
	}
}
