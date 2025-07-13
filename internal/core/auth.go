package core

import (
	"context"
	"time"

	"github.com/Cheasezz/authSrvc/pkg/tokens"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(ctx context.Context, userId uuid.UUID, userAgent, ip string) (tokens.TokensPair, error)
}

type AuthRepo interface {
	Signup(ctx context.Context, session Session) error
}

type Session struct {
	UserId       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	Ip           string    `db:"ip"`
	CreatedAt    time.Time `db:"created_at"`
	ExpriresAt   time.Time `db:"expires_at"`
}
