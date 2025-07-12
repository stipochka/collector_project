package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string `mapstructure:"env" default:"local"` //local, dev, prod
	Storage    `mapstructure:"storage" required:"true"`
	HTTPServer `mapstructure:"http_server" required:"true"`
}

type Storage struct {
	Type     string `mapstructure:"type" required:"true"`
	Address  string `mapstructure:"address" required:"true"`
	Database string `mapstructure:"database"`
	Table    string `mapstructure:"table"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type HTTPServer struct {
	Address      string        `mapstructure:"address"`
	Port         int           `mapstructure:"port"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func MustLoadConfig[T any](path string) T {
	var config T
	if err := LoadConfig(path, &config); err != nil {
		panic(err)
	}
	return config
}

func LoadConfig[T any](path string, config *T) error {
	const op = "config.LoadConfig"

	reader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: failed to open config file %w", op, err)
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(reader); err != nil {
		return fmt.Errorf("%s: failed to read config file: %w", op, err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("%s: failed to Unmarshal config: %w", op, err)
	}

	return nil
}
