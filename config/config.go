package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
}

type HTTP struct {
	Host string `env-required:"false" yaml:"host" env:"HOST"`
	Port string `env-required:"false" yaml:"port" env:"PORT"`
}

type Log struct {
	Level string `env-required:"false" yaml:"log_level" env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
