package app

import (
	"errors"

	"github.com/Cheasezz/authSrvc/config"
	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/internal/repo"
	"github.com/Cheasezz/authSrvc/internal/services"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/Cheasezz/authSrvc/pkg/pgx5"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

var (
	ErrTokensNew = errors.New("error when tokens.New in NewEnv")
	ErrPoolNew   = errors.New("error when pgx5.New in NewEnv")
)

type Env struct {
	Logger   logger.Logger
	Services core.AuthService
	db       *pgx5.Pgx5
	TM       tokens.Manager
}

func NewEnv(cfg *config.Config) (*Env, error) {
	logger := logger.New(cfg.Log.Level)

	manager, err := tokens.New(cfg.Auth.SigningKey, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	if err != nil {
		return nil, errors.Join(ErrTokensNew, err)
	}

	db, err := pgx5.New(cfg.PG.URL)
	if err != nil {
		return nil, errors.Join(ErrPoolNew, err)
	}

	repo := repo.New(db)

	services := services.New(manager, repo)
	env := Env{
		Logger:   logger,
		Services: services,
		db:       db,
		TM:       manager,
	}

	return &env, nil
}

func (env *Env) Close() {
	env.db.Close()
}
