package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type C struct {
	Endpoint string `yaml:"endpoint" env:"ENDPOINT"`
	PgDSN    string `yaml:"pg_dsn" env:"PG_DSN"`
}

func New() (*C, error) {
	cfg := &C{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config.New: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
