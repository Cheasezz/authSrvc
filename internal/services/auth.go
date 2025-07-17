package services

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCreateTokens = errors.New("error when create tokenst pair in Signup")
	ErrRepoSignup   = errors.New("error when repo.Signup in Signup")
)

func (s *services) Signup(ctx context.Context, userId uuid.UUID, userAgent, ip string) (core.SignupResult, error) {
	var tp core.SignupResult
	sessionId := uuid.New()

	claimsA := core.AccressTokenClaims{
		SessionId: sessionId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
		},
	}
	claimsB := core.RefreshTokenClaims{
		SessionId: sessionId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
		},
	}
	tknPair, err := s.tokenManager.NewTokensPair(claimsA, claimsB)
	if err != nil {
		return tp, errors.Join(ErrCreateTokens, err)
	}

	bcryptRT, err := hashRT(tknPair.RefreshToken)
	if err != nil {
		return tp, err
	}

	sessionInfo := core.Session{
		UserId:       userId,
		RefreshToken: string(bcryptRT),
		UserAgent:    userAgent,
		Ip:           ip,
		ExpriresAt:   time.Now().Add(s.refreshTTL),
	}

	err = s.repo.Signup(ctx, sessionInfo)
	if err != nil {
		return tp, errors.Join(ErrRepoSignup, err)
	}
	tp = core.SignupResult{
		Access:     tknPair.AccessToken,
		Refresh:    tknPair.RefreshToken,
		RefreshTTL: s.refreshTTL,
	}
	return tp, nil
}

func hashRT(token string) (string, error) {
	// У меня рефрешь токен - jwt, он может быть длиннее 72 байт,
	// а bcrypt хеширует первые 72 байта.
	// Получается входеые данные усекаются, хеши могут часто повторяться.
	// Поэтому сначала хеширую до 32 байт.
	hashRT := sha256.Sum256([]byte(token))

	bcryptRT, err := bcrypt.GenerateFromPassword(hashRT[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptRT), nil
}

func compareHashRT(token, hash string) bool {
	hashRT := sha256.Sum256([]byte(token))

	err := bcrypt.CompareHashAndPassword([]byte(hash), hashRT[:])
	return err == nil
}
