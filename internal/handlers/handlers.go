package handlers

import (
	"github.com/Cheasezz/authSrvc/internal/app"
	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func New(env *app.Env) *Handlers {
	return &Handlers{}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "err",
		})
	})

	{
		api := router.Group("/api")
	}
	return router
}
