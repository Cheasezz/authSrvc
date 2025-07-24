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
	PG   `yaml:"pg"`
	App  `yaml:"app"`
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

type PG struct {
	URL string `env-required:"false" yaml:"pg_url"   env:"PG_URL"`
}

type App struct {
	DevMod     bool   `env-required:"false" yaml:"dev_mod" env:"MOD"`
	WebhookUrl string `env-required:"false" yaml:"webhook_url" env:"WEBHOOK_URL"`
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
