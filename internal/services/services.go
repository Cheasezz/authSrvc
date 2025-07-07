package services

import (
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type services struct {
	tokenManager tokens.Manager
}

func New(tm tokens.Manager) *services {
	return &services{
		tokenManager: tm,
	}
}
