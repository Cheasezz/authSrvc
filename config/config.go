package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
	Auth `yaml:"auth"`
}

type HTTP struct {
	Host string `env-required:"false" yaml:"host" env:"HOST"`
	Port string `env-required:"false" yaml:"port" env:"PORT"`
}

type Log struct {
	Level string `env-required:"false" yaml:"log_level" env:"LOG_LEVEL"`
}

type Auth struct {
	SigningKey      string        `env-required:"false" yaml:"signing_key" env:"SIGNING_KEY"`
	AccessTokenTTL  time.Duration `env-required:"false" yaml:"accessttl" env:"ATTL"`
	RefreshTokenTTL time.Duration `env-required:"false" yaml:"refreshttl" env:"RTTL"`
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
