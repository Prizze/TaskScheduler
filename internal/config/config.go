package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "internal/config/config.yaml"

type Config struct {
	JwtSecret string
	HTTPAddr  string `yaml:"addr"`
	DBURL     string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.DBURL = os.Getenv("DB_URL")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	if cfg.HTTPAddr == "" {
		return nil, errors.New("http addr is required")
	}

	if cfg.DBURL == "" {
		return nil, errors.New("DB_URL is required")
	}

	if cfg.JwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}
