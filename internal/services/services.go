package services

import (
	"time"

	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
)

type services struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	webhookUrl string

	tokenManager tokens.Manager
	repo         core.AuthRepo
	log          logger.Logger
}

func New(tm tokens.Manager, repo core.AuthRepo, log logger.Logger, attl, rttl time.Duration, webhookUrl string) *services {
	return &services{
		accessTTL:    attl,
		refreshTTL:   rttl,
		webhookUrl:   webhookUrl,
		tokenManager: tm,
		repo:         repo,
		log:          log,
	}
}
