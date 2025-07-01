package app

import (
	"github.com/Cheasezz/authSrvc/config"
	"github.com/Cheasezz/authSrvc/pkg/logger"
)

type Env struct {
	Logger logger.Logger
}

func NewEnv(cfg *config.Config) (*Env, error) {
	logger := logger.New(cfg.Log.Level)

	env := Env{
		Logger: logger,
	}

	return &env, nil
}

func (env *Env) Close() {}
