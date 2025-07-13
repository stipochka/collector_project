package main

import (
	"dashboard/internal/cache"
	"dashboard/internal/config"
	"dashboard/internal/handler"
	"dashboard/internal/repository"
	"dashboard/internal/service"
	"dashboard/internal/storage"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env files")
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("failed to get config path variable")
	}

	cfg := config.MustLoadConfig[config.Config](cfgPath)

	lgr := setupLogger(cfg.Env)

	fmt.Println("User:", cfg.Storage.User)
	fmt.Println("Password:", cfg.Storage.Password)
	fmt.Println("Address:", cfg.Storage.Address)
	lgr.With(zap.String("service-name", "dashboard"))
	lgr.Info("started dashboard service")

	db, err := storage.NewStorage(cfg.Storage.Address, cfg.Storage.Database, cfg.Storage.User, cfg.Storage.Password)
	if err != nil {
		lgr.Fatal("failed to open db connect", zap.Any("error", err))
	}

	repo := repository.NewClickhouseRepo(*db)

	telemetryCache := cache.NewInMemoryCache(time.Minute * 5)

	service := service.NewTelemetryService(repo, telemetryCache)

	handler := handler.NewHandler(service, lgr)

	router := handler.InitRoutes()

	lgr.Info("started dashboard http server", zap.String("address", cfg.HTTPServer.Address))

	if err := router.Run(fmt.Sprintf("%s:%d", cfg.HTTPServer.Address, cfg.HTTPServer.Port)); err != nil {
		lgr.Error("error while running dashboard server", zap.Any("error", err))
	}
}

func setupLogger(env string) *zap.Logger {
	var lgr *zap.Logger
	switch env {
	case envDev:
		lgr = zap.Must(zap.NewDevelopment())
	case envProd:
		lgr = zap.Must(zap.NewProduction())
	default:
		lgr = zap.NewExample()
	}

	return lgr
}
