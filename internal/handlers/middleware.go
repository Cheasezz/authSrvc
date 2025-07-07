package handlers

import (
	"errors"

	"github.com/Cheasezz/authSrvc/internal/apperrors"
	"github.com/gin-gonic/gin"
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
