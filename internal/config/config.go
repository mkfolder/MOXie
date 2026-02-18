package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Makefolder/cynero/internal/common"
	"github.com/caarlos0/env/v10"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
	HTTP     HTTP     `yaml:"http"`
}

type Server struct {
	Environment common.Environment `yaml:"environment" env:"ENVIRONMENT"`
	Port        string             `yaml:"port" env:"SERVER_PORT"`
}

type Postgres struct {
	DSN string `yaml:"dsn" env:"DATABASE_URL"`
}

type HTTP struct {
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
}

func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("error parsing env vars: %w", err)
	}

	if !config.Server.Environment.IsProduction() && !config.Server.Environment.IsDevelopment() {
		return nil, errors.New("invalid environment")
	}

	return &config, nil
}
