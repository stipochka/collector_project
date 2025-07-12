package main

import (
	"dashboard/internal/config"
	"dashboard/internal/storage"
	"log"
	"os"

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

	lgr.With(zap.String("service-name", "dashboard"))
	lgr.Info("started dashboard service")

	db, err := storage.NewStorage(cfg.A)
	//TODO: init storage

	//TODO: init repo

	//TODO: init router -- go-chi
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
