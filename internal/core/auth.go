package core

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(ctx context.Context, userId uuid.UUID, userAgent, ip string) (SignupResult, error)
}

type AuthRepo interface {
	Signup(ctx context.Context, session Session) error
}

type AccressTokenClaims struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}

type SignupResult struct {
	Access     string
	Refresh    string
	RefreshTTL time.Duration
}

type Session struct {
	UserId       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	Ip           string    `db:"ip"`
	CreatedAt    time.Time `db:"created_at"`
	ExpriresAt   time.Time `db:"expires_at"`
}
