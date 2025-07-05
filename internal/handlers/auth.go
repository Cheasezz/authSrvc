package handlers

import (
	"errors"
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrUncorrectUuid = errors.New("signup error: uncorrect uuid")
	ErrorServer      = errors.New("signup errror: error on server side")
)

// Проверяем параметр запроса "uuid".
// Создаем пару токеном. Refresh в куки, access в json ответе.
func (h *Handlers) signup(c *gin.Context) {
	id := c.Query("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(apperrors.New(err, ErrUncorrectUuid))
		return
	}

	tkns, err := h.services.Signup(c, uuid.String())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrorServer))
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refreshToken", tkns.RefreshToken, int(tkns.RefreshTokenTtl.Seconds()), "", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"accessToken": tkns.AccessToken,
	})
}
