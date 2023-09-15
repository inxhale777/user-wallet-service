package config

import "fmt"
import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Addr  string `yaml:"addr" env:"ADDR"`
	PgURL string `yaml:"pg_url" env:"PG_URL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
