package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handlers) errResponse(c *gin.Context, statusCode int, logErr, UserErr error) {
	h.logger.Error("path: %s, logErr: %s, UserErr: %s", c.FullPath(), logErr, UserErr)
	c.AbortWithStatusJSON(statusCode, gin.H{
		"success": false,
		"message": UserErr.Error(),
	})
}
