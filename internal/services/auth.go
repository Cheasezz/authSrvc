package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCreateTokens      = errors.New("error when create tokens pair in IssueTokens")
	ErrRepoSignup        = errors.New("error when repo.CreateSession in IssueTokens")
	ErrRefresh           = errors.New("error in Refresh")
	ErrTokenshashCompare = errors.New("error when compare tokens hash in Refresh")
	ErrUserAgent         = errors.New("error when compare user agent in Refresh (not equal)")
	ErrWebhookStatusCode = errors.New("error webhook returned not 2xx status code")
)

func (s *services) IssueTokens(ctx context.Context, userId uuid.UUID, userAgent, ip string) (*core.TokenPairResult, error) {
	var tp core.TokenPairResult
	sessionId := uuid.New()

	claimsA := core.AccessTokenClaims{
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
		return nil, errors.Join(ErrCreateTokens, err)
	}

	bcryptRT, err := hashRT(tknPair.RefreshToken)
	if err != nil {
		return nil, err
	}

	sessionInfo := core.Session{
		Id:               sessionId,
		UserId:           userId,
		RefreshTokenHash: string(bcryptRT),
		UserAgent:        userAgent,
		Ip:               ip,
		ExpriresAt:       time.Now().Add(s.refreshTTL),
	}

	err = s.repo.CreateSession(ctx, &sessionInfo)
	if err != nil {
		return nil, errors.Join(ErrRepoSignup, err)
	}
	tp = core.TokenPairResult{
		Access:     tknPair.AccessToken,
		Refresh:    tknPair.RefreshToken,
		RefreshTTL: s.refreshTTL,
	}
	return &tp, nil
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

func compareHashRT(token, hash string) error {
	hashRT := sha256.Sum256([]byte(token))

	err := bcrypt.CompareHashAndPassword([]byte(hash), hashRT[:])
	if err != nil {
		return ErrTokenshashCompare
	}

	return nil
}

// При обновлении токенов сравниваем данные из окружения пользователя.
// Если все корректно, то удаляем старую сессию, и создаем новую, с новыми токенами.
func (s *services) Refresh(ctx context.Context, refreshTkn, sessionId, userAgent, ip string) (*core.TokenPairResult, error) {
	session, err := s.repo.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, errors.Join(ErrRefresh, err)
	}

	// Если юзер агент не совпадает,
	// то деавторизуем пользователя, удаляя сессию из бд.
	if userAgent != session.UserAgent {
		err := s.repo.DeleteSessionById(ctx, sessionId)
		if err != nil {
			return nil, errors.Join(ErrRefresh, err)
		}
		return nil, ErrUserAgent
	}

	if err := compareHashRT(refreshTkn, session.RefreshTokenHash); err != nil {
		return nil, errors.Join(ErrRefresh, err)
	}

	// У меня ip в БД - это тип inet, он храниться с маской.
	// Перед сравнением нужео убрать маску.
	parsIp, _, err := net.ParseCIDR(session.Ip)
	if err != nil {
		return nil, err
	}

	// При попытке обновления токена с нового ip,
	// делаем не блокирующий post зарос на вебхук.
	if ip != parsIp.String() {
		go func() {
			err := s.sendWebhook(&core.RefreshChangeIpPayload{
				SessionId: sessionId,
				UserId:    session.UserId.String(),
				OldIp:     session.Ip,
				NewIp:     ip,
				Timestamp: time.Now(),
			})
			if err != nil {
				s.log.Info("Webhook alert:%s", err)
			}
		}()
	}

	err = s.repo.DeleteSessionById(ctx, sessionId)
	if err != nil {
		return nil, errors.Join(ErrRefresh, err)
	}

	newTokens, err := s.IssueTokens(ctx, session.UserId, userAgent, ip)
	if err != nil {
		return nil, errors.Join(ErrRefresh, err)
	}

	return newTokens, nil
}

func (s *services) sendWebhook(payload *core.RefreshChangeIpPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.webhookUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return ErrWebhookStatusCode
	}

	return nil
}
