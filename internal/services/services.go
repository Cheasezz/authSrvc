package services

import (
	"time"

	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type services struct {
	accessTTL    time.Duration
	refreshTTL   time.Duration
	tokenManager tokens.Manager
	repo         core.AuthRepo
}

func New(tm tokens.Manager, repo core.AuthRepo, attl, rttl time.Duration) *services {
	return &services{
		accessTTL:    attl,
		refreshTTL:   rttl,
		tokenManager: tm,
		repo:         repo,
	}
}
