package handlers

import (
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/app"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	logger logger.Logger
}

func New(env *app.Env) *Handlers {
	return &Handlers{
		logger: env.Logger,
	}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "err",
		})
	})

	{
		api := router.Group("/api")
		api.POST("/signup", h.signup)
	}
	return router
}
