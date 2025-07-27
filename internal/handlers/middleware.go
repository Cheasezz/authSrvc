package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Cheasezz/authSrvc/internal/apperrors"
	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	refreshTknCookie    = "refreshToken"
)

var (
	ErrEmptyAuthHeader         = errors.New("error empty auth header")
	ErrInvalidAuthHeader       = errors.New("error invalid auth header")
	ErrParseAccessToken        = errors.New("error parse access token from header")
	ErrAuth                    = errors.New("error authorization")
	ErrEmptyRefreshCookie      = errors.New("error empty refresh cookie")
	ErrTypeAssertAccessClaims  = errors.New("error type assertion access token claims")
	ErrTypeAssertRefreshClaims = errors.New("error type assertion refresh token claims")
	ErrNotEqualSessionId       = errors.New("error not equal sessionId in claims")
)

func (h *Handlers) errMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		var appErr *apperrors.AppError
		err := c.Errors.Last().Err

		if errors.As(err, &appErr) {
			h.logger.WithField(gin.H{
				"path":   c.FullPath(),
				"userIp": c.ClientIP(),
			}).Error(appErr.LogErr)

			c.AbortWithStatusJSON(c.Writer.Status(), ErrorResponse{
				Success: false,
				Message: appErr.Error(),
			})
		}
	}
}

// Проверяем корректность заголовка Authorization
// И записываем токен от туда в контекст
func extractAccessToken(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrEmptyAuthHeader, ErrEmptyAuthHeader))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrInvalidAuthHeader, ErrInvalidAuthHeader))
		return
	}
	c.Set(accessTknCtx, headerParts[1])

	c.Next()
}

// Проверяем наличие рефреш куки.
// Елсли есть, то записываем в контекст.
func extractRefreshToken(c *gin.Context) {
	token, err := c.Cookie(refreshTknCookie)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrEmptyRefreshCookie, ErrEmptyRefreshCookie))
		return
	}
	c.Set(refreshTknCtx, token)

	c.Next()
}

// Проверяем токен доступа из контекста
// Парсим и записываем userId из claims в контекст
func (h *Handlers) checkUserAccess(c *gin.Context) {
	token, err := getAccessTknCtx(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
		return
	}

	parsedToken, err := h.tm.ParseAccessToken(token, &core.AccessTokenClaims{})
	if err != nil {
		if errors.Is(err, tokens.ErrAccessTokenExpired) {
			c.Status(http.StatusUnauthorized)
			c.Error(apperrors.New(err, tokens.ErrAccessTokenExpired))
			return
		}

		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrParseAccessToken))
		return
	}
	// В sub находится userId
	userId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
		return
	}

	c.Set(userIdCtx, userId)

	c.Next()
}

// Достаем токен доступа и рефреш из контекста.
// Парсим их (истекший токен доступа - не ошибка).
// Сравниваем sessionId из claims
// Если одинаковый, то записываем его в текущий контекст.
func (h *Handlers) checkTokenPair(c *gin.Context) {
	accessT, err := getAccessTknCtx(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
		return
	}

	refreshT, err := getRefreshTknCtx(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
		return
	}

	parsedTkns, err := h.tm.ParseTokenPair(accessT, refreshT, &core.AccessTokenClaims{}, &core.RefreshTokenClaims{})
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
		return
	}

	accessClaims, ok := parsedTkns.AccessToken.Claims.(*core.AccessTokenClaims)
	if !ok {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrTypeAssertAccessClaims, ErrAuth))
		return
	}

	refreshClaims, ok := parsedTkns.RefreshToken.Claims.(*core.RefreshTokenClaims)
	if !ok {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrTypeAssertRefreshClaims, ErrAuth))
		return
	}

	if accessClaims.SessionId != refreshClaims.SessionId {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(ErrNotEqualSessionId, ErrAuth))
		return
	}

	c.Set(sessionIdCtx, accessClaims.SessionId)

	c.Next()
}
