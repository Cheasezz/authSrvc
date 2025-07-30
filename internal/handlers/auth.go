package handlers

import (
	"errors"
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/apperrors"
	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrUncorrectUuid = errors.New("create session error: uncorrect uuid")
	ErrTokenIssuance = errors.New("create session error: error on server side or user already exist")
	ErrGetUserId     = errors.New("me error: error on server side")
)

// @Tags auth
// @Summary create session
// @Description create session in db with ip and user agent.
// @Description Return access token in JSON and refresh token in cookies
// @ID create-session
// @Produce json
// @Param 	uuid query string true "User uuid" example(fb62aa81-1172-4c73-8fc3-cd5a446346bf)
// @Success 200 {object} TokenResponse
// @Header 	200 {string} Set-Cookie "JWT refreshToken Example: refreshToken=9838c5.9cf.f93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None"
// @Failure 400 {object} errBadRequestResp
// @Failure 500 {object} errTokenIssuanceResp
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
		c.Error(apperrors.New(err, ErrTokenIssuance))
		return
	}

	newTokensResponse(c, tkns)
}

// @Tags auth
// @Summary return curent user id
// @Description chek Authorization header and extract user id from claims in jwt.
// @Description Set user id in gin context and return json with user id.
// @ID me
// @Produce  json
// @Success 200 {object} MeResponse
// @Failure 500 {object} errMeResp
// @Security		bearerAuth
// @Router 	/session/me [get]
func (h *Handlers) me(c *gin.Context) {
	userId, err := getUserIdCtx(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrGetUserId))
		return
	}

	c.JSON(http.StatusOK, MeResponse{UserId: userId.String()})

}

// @Tags auth
// @Summary refresh session
// @Description check access and refresh tokens.
// @Description Abort request if user agent not like in db and delete session from db.
// @Description Send post request on webhook if ip not like in db.
// @Description If everything ok, then delete old session from db and generete new with new tokens.
// @ID refresh-session
// @Produce json
// @Param Cookie header string true "Refersh token cookie"
// @Success 200 {object} TokenResponse
// @Header 	200 {string} Set-Cookie "JWT refreshToken Example: refreshToken=9838c5.9cf.f93e21; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None"
// @Failure 401 {object} errAuthResp
// @Failure 500 {object} errTokenIssuanceResp
// @Security		bearerAuth
// @Router 	/session/refresh [post]
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

	tkns, err := h.services.Refresh(c, refreshT, sessionId, ua, ip)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(apperrors.New(err, ErrTokenIssuance))
		return
	}

	newTokensResponse(c, tkns)
}

// @Tags auth
// @Summary logout
// @Description check access and refresh tokens.
// @Description Delete session by id and set empty tokens in json response and cookie.
// @Description Even if there are errors in the handler, it will still set empty tokens.
// @ID delete-session
// @Produce json
// @Success 200 {object} TokenResponse
// @Header 	200 {string} Set-Cookie "JWT refreshToken Example: refreshToken=; Path=/; Max-Age=2628000; HttpOnly; Secure; SameSite=None"
// @Failure 401 {object} errAuthResp
// @Security		bearerAuth
// @Router 	/session [delete]
func (h *Handlers) logout(c *gin.Context) {
	var emptyTkns *core.TokenPairResult

	sessionId, err := getSessionIdCtx(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
	}

	emptyTkns, err = h.services.DeleteSession(c, sessionId)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		c.Error(apperrors.New(err, ErrAuth))
	}

	newTokensResponse(c, emptyTkns)
}
