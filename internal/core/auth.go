package core

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	IssueTokens(ctx context.Context, userId uuid.UUID, userAgent, ip string) (*TokenPairResult, error)
	Refresh(ctx context.Context, refreshTkn, sessionId, userAgent, ip string) (*TokenPairResult, error)
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionById(ctx context.Context, sessionId string) (*Session, error)
	DeleteSessionById(ctx context.Context, sessionId string) error
}

type AccessTokenClaims struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPairResult struct {
	Access     string
	Refresh    string
	RefreshTTL time.Duration
}

type Session struct {
	Id               uuid.UUID `db:"id"`
	UserId           uuid.UUID `db:"user_id"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	UserAgent        string    `db:"user_agent"`
	Ip               string    `db:"ip"`
	CreatedAt        time.Time `db:"created_at"`
	ExpriresAt       time.Time `db:"expires_at"`
}

// Структура для отапрвки данных на вебхук
type RefreshChangeIpPayload struct {
	SessionId string    `json:"session_id"`
	UserId    string    `json:"user_id"`
	OldIp     string    `json:"old_ip"`
	NewIp     string    `json:"new_ip"`
	Timestamp time.Time `json:"timestamp"`
}
