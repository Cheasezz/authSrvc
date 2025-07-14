package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Cheasezz/authSrvc/internal/apperrors"
	"github.com/Cheasezz/authSrvc/pkg/tokens"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

var (
	ErrEmptyAuthHeader   = errors.New("error empty auth header")
	ErrInvalidAuthHeader = errors.New("error invalid auth header")
	ErrParseToken        = errors.New("error parse token from header")
	ErrUserCtxNotFound   = errors.New("error userCtx not found")
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

func (m *Handlers) userIdentity(c *gin.Context) {
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

	userId, err := m.tm.ParseAccessToken(headerParts[1])
	if err != nil {
		if errors.Is(err, tokens.ErrTokenExpired) {
			c.Status(http.StatusUnauthorized)
			c.Error(apperrors.New(err, tokens.ErrTokenExpired))
			return
		}

		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrParseToken))
		return
	}

	c.Set(userCtx, userId)
}

func getUserCtx(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return uuid.UUID{}, ErrUserCtxNotFound
	}
	parsedId, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedId, nil
}
