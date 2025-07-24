package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	userIdCtx     = "userId"
	accessTknCtx  = "access"
	refreshTknCtx = "refresh"
	sessionIdCtx  = "session"
)

var (
	ErrUserCtxNotFound       = errors.New("error userIdCtx not found")
	ErrAccessTknCtxNotFound  = errors.New("error accessTknCtx not found")
	ErrRefreshTknCtxNotFound = errors.New("error refreshTknCtx not found")
	ErrSessionIdCtxNotFound  = errors.New("error sessionIdCtx not found")
)

// Достать accessToken из контекста.
func getAccessTknCtx(c *gin.Context) (string, error) {
	token, ok := c.Get(accessTknCtx)
	if !ok {
		return "", ErrAccessTknCtxNotFound
	}

	return token.(string), nil
}

// Достать refreshToken из контекста.
func getRefreshTknCtx(c *gin.Context) (string, error) {
	token, ok := c.Get(refreshTknCtx)
	if !ok {
		return "", ErrRefreshTknCtxNotFound
	}

	return token.(string), nil
}

// Достать userId из контекста.
func getUserIdCtx(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userIdCtx)
	if !ok {
		return uuid.UUID{}, ErrUserCtxNotFound
	}
	parsedId, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedId, nil
}

// Достать sessionId из контекста.
func getSessionIdCtx(c *gin.Context) (string, error) {
	sessionId, ok := c.Get(sessionIdCtx)
	if !ok {
		return "", ErrSessionIdCtxNotFound
	}

	return sessionId.(string), nil
}
