package main

import (
	"collector-service/internal/app"
	"collector-service/internal/config"
	"collector-service/internal/database"
	"collector-service/internal/migrate"
	"collector-service/internal/prtettylogger"
	"collector-service/internal/repository"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfgPath := flag.String("path", "", "used to get path to config")
	flag.Parse()

	cfg := config.MustLoadConfig[config.Config](*cfgPath)

	log := setupLogger(cfg.Env)

	log.Info(
		"starting collector service",
		slog.String("env", cfg.Env),
	)

	ctx, cancel := context.WithCancel(context.Background())

	db, err := database.ConnectToDB(cfg.Storage.Address, cfg.Database, cfg.Storage.Username, cfg.Storage.Password)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = migrate.Migrate(*db, cfg.Storage.MigrationPath)
	if err != nil {
		log.Error("failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo := repository.NewClickHouseRepository(*db)

	application := app.NewApp(log, cfg.GRPCPort, repo)
	go func() {
		application.GRPCServer.MustRun(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Info("Gracefully stopping server")
	cancel()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := prtettylogger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
