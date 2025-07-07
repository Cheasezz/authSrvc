package services

import (
	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type services struct {
	tokenManager tokens.Manager
	repo         core.AuthRepo
}

func New(tm tokens.Manager, repo core.AuthRepo) *services {
	return &services{
		tokenManager: tm,
		repo:         repo,
	}
}
