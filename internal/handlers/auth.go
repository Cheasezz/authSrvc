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
	ErrSignup        = errors.New("signup errror: error on server side or user already exist")
	ErrGetUserId     = errors.New("getUserId error: error on server side")
)

// @Tags auth
// @Summary create account
// @Description create account in db and return access token in JSON and refresh token in cookies
// @ID create-account
// @Accept  json
// @Produce json
// @Param 	uuid query string true "User uuid" example(fb62aa81-1172-4c73-8fc3-cd5a446346bf)
// @Success 200 {object} TokenResponse
// @Header 	200 {string} Set-Cookie "JWT refreshToken Example: refreshToken=9838c5.9cf.f93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None"
// @Failure 400 {object} errBadRequestResp
// @Failure 500 {object} errSignupResp
// @Router 	/session [post]
func (h *Handlers) tokenIssuance(c *gin.Context) {
	id := c.Query("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(apperrors.New(err, ErrUncorrectUuid))
		return
	}

	ua := c.Request.UserAgent()
	ip := c.ClientIP()

	tkns, err := h.services.IssueTokens(c, uuid, ua, ip)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrSignup))
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTknCookie, tkns.Refresh, int(tkns.RefreshTTL.Seconds()), "", "", false, true)
	c.JSON(http.StatusOK, TokenResponse{AccessToken: tkns.Access})
}

// @Tags auth
// @Summary return curent user id
// @Description chek Authorization header and extract user id from claims in jwt.
// Set user id in gin context and return json with user id.
// @ID getuserid
// @Produce  json
// @Success 200 {object} UserIdResponse
// @Failure 500 {object} errGetUserIdResp
// @Security		bearerAuth
// @Router 	/session/me [get]
func (h *Handlers) me(c *gin.Context) {
	userId, err := getUserIdCtx(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrGetUserId))
		return
	}

	c.JSON(http.StatusOK, UserIdResponse{UserId: userId.String()})

}

func (h *Handlers) refresh(c *gin.Context) {
	sessionId, err := getSessionIdCtx(c)
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

	ua := c.Request.UserAgent()
	ip := c.ClientIP()

	h.services.Refresh(c, refreshT, sessionId, ua, ip)
}
