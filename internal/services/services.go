package services

import (
	"context"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type Services interface {
	Signup(ctx context.Context, userId string) (tokens.TokensPair, error)
}

type services struct {
	tokenManager tokens.Manager
}

func New(tm tokens.Manager) *services {
	return &services{
		tokenManager: tm,
	}
}
