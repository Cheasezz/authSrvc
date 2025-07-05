package handlers

import (
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/app"
	"github.com/Cheasezz/authSrvc/internal/services"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	logger   logger.Logger
	services services.Services
}

func New(env *app.Env) *Handlers {
	return &Handlers{
		logger:   env.Logger,
		services: env.Services,
	}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), h.errMiddleware)

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
