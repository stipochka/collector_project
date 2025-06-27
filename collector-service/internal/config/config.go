package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Storage    `yaml:"storage" env-required:"true"`
	GRPCServer `yaml:"grpc_server"`
}

type Storage struct {
	Type          string        `yaml:"type"`
	Address       string        `yaml:"address" env-required:"true"`
	Database      string        `yaml:"database" env-required:"true"`
	Table         string        `yaml:"table" env-required:"true"`
	Username      string        `yaml:"username" env-required:"true"`
	Password      string        `yaml:"password" env-required:"true"`
	MigrationPath string        `yaml:"migration_path" env-default:"config/config.yaml"`
	BatchSize     int           `yaml:"batch_size" env-default:"1000"`
	FlushTimeout  time.Duration `yaml:"flush_timeout" env-default:"1"`
}

type GRPCServer struct {
	GRPCPort     int `yaml:"grpc_port" env-default:"9000"`
	MaxMsgSizeMb int `yaml:"max_msg_size_mb" env-default:"10"`
}

func MustLoadConfig[T any](path string) *T {
	cfg, err := loadConfig[T](path)
	if err != nil {
		log.Fatalf("failed to load config %v", err)
		return nil
	}
	return &cfg

}

func loadConfig[T any](path string) (T, error) {
	var cfg T

	if _, err := os.Stat(path); err != nil {
		return cfg, err
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
