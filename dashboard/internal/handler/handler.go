package handler

import (
	"dashboard/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	log     *zap.Logger
	service service.TelemetryService
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	api.Use(h.MiddlewareLogger())
	{
		logs := api.Group("/logs", h.GetLogs)
		{
			errors := logs.Group("/errors")
			{
				errors.GET("/top", h.)
				errors.GET("/recent", h.)
			}

			levels := logs.Group("/levels", h.) 
			{
				levels.GET("/stats", h)
			}
		}
		services := api.Group("/services", h.)
		{
			services.GET("/stats", h.)
		}
	}


}
