package handler

import (
	"dashboard/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log     *zap.Logger
	service service.Service
}

func NewHandler(service service.Service, log *zap.Logger) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	api.Use(h.MiddlewareLogger())

	{
		logs := api.Group("/logs")
		{
			logs.GET("", h.GetLogs)

			errors := logs.Group("/errors")
			{
				errors.GET("/top", h.GetTopErrors)
				errors.GET("/recent", h.GetRecentErrors)
			}

			levels := logs.Group("/levels")
			{
				levels.GET("", h.GetAllLogLevels)
				levels.GET("/stats", h.GetLogLevelStats)
			}
		}

		services := api.Group("/services")

		{
			services.GET("", h.GetAllServiceNames)
			services.GET("/stats", h.GetServiceStats)
		}
	}

	return router
}
