package handlers

import (
	"net/http"

	"github.com/Cheasezz/authSrvc/internal/app"
	"github.com/Cheasezz/authSrvc/internal/core"
	"github.com/Cheasezz/authSrvc/pkg/logger"
	"github.com/gin-gonic/gin"

	_ "github.com/Cheasezz/authSrvc/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	logger   logger.Logger
	services core.AuthService
}

func New(env *app.Env) *Handlers {
	return &Handlers{
		logger:   env.Logger,
		services: env.Services,
	}
}

func (h *Handlers) Init(devMod bool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), h.errMiddleware)
	// GIN-debug и Swagger только во время режима разработки
	if devMod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(-1)))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

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
