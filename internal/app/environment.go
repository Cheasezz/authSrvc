package app

import (
	"errors"

	"github.com/Cheasezz/authSrvc/config"
	"github.com/Cheasezz/authSrvc/internal/services"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

var (
	ErrTokensNew = errors.New("error when tokens.New in NewEnv")
)

type Env struct {
	Logger     logger.Logger
	TknManager tokens.Manager
	Services   services.Services
}

func NewEnv(cfg *config.Config) (*Env, error) {
	logger := logger.New(cfg.Log.Level)

	manager, err := tokens.New(cfg.Auth.SigningKey, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return nil, errors.Join(ErrTokensNew, err)
	}

	services := services.New(manager)
	env := Env{
		Logger:     logger,
		TknManager: manager,
		Services:   services,
	}

	return &env, nil
}

func (env *Env) Close() {}
