package app

import (
	"github.com/Cheasezz/authSrvc/config"
)

type Env struct{}

func NewEnv(cfg *config.Config) (*Env, error) {
	env := Env{}

	return &env, nil
}

func (env *Env) Close() {}
