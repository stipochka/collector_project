package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) MiddlewareLogger() gin.HandlerFunc {
	log := h.log
	return func(c *gin.Context) {
		start := time.Now()

		log.Info("received request",
			zap.String("method", c.Request.Method),
			zap.String("client ip", c.ClientIP()),
			zap.String("url", c.Request.URL.String()),
		)

		c.Next()

		duration := time.Since(start)

		log.Info("request completed",
			zap.Duration("completed in", duration),
			zap.Int("status", c.Writer.Status()),
		)
	}

}
