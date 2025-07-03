package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	errUncorrectUuid = errors.New("signup error: uncorrect uuid")
)

func (h *Handlers) signup(c *gin.Context) {
	id := c.Query("uuid")
	_, err := uuid.Parse(id)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err, errUncorrectUuid)
		return
	}

	// TODO: generate tokens
}
