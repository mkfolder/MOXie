package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/mkfolder/moxie/internal/common"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
	HTTP     HTTP     `yaml:"http"`
	Workers  Workers  `yaml:"workers"`
	Auth     Auth     `yaml:"auth"`
}

type Server struct {
	Environment common.Environment `yaml:"environment" env:"ENVIRONMENT"`
	Port        string             `yaml:"port" env:"SERVER_PORT"`
	WebhookURL  string             `yaml:"webhook_url" env:"WEBHOOK_URL"`
}

type Postgres struct {
	DSN string `yaml:"dsn" env:"DATABASE_URL"`
}

type HTTP struct {
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
}

type Workers struct {
	CleanerInterval time.Duration `yaml:"cleaner_interval" env:"CLEANER_INTERVAL"`
	OrderExpiration time.Duration `yaml:"order_expiration" env:"ORDER_EXPIRATION"`
}

type Auth struct {
	JWTSecret       string        `yaml:"jwt_secret" env:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL"`
	EncryptionKey   string        `yaml:"encryption_key" env:"ENCRYPTION_KEY"`
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

	if config.Server.Port == "" {
		return nil, errors.New("invalid server port")
	}

	if config.Server.WebhookURL == "" {
		return nil, errors.New("invalid webhook url")
	}

	if config.Postgres.DSN == "" {
		return nil, errors.New("invalid postgres dsn")
	}

	if config.HTTP.Timeout == 0 {
		return nil, errors.New("invalid http timeout")
	}

	if config.Auth.JWTSecret == "" {
		return nil, errors.New("invalid jwt secret")
	}

	if config.Auth.AccessTokenTTL == 0 {
		return nil, errors.New("invalid access token ttl")
	}

	if config.Auth.RefreshTokenTTL == 0 {
		return nil, errors.New("invalid refresh token ttl")
	}

	if config.Auth.EncryptionKey == "" {
		return nil, errors.New("invalid encryption key")
	}

	return &config, nil
}
