package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func (h *Handler) MiddlewareLogger() gin.HandlerFunc {
	log := h.log
	return func(c *gin.Context) {
		queryStTime := time.Now()

		log.Info(
			"received request",
		)

		c.Next()

		log.Info(
			"finishing working with request in",
			zapcore.Field{
				Key:       "time",
				Interface: time.Now().Sub(queryStTime),
			},
			zapcore.Field{
				Key:       "request url",
				Interface: c.Request.URL,
			},
			zapcore.Field{
				Key:       "request method",
				Interface: c.Request.Method,
			},
			zapcore.Field{
				Key:       "user IP",
				Interface: c.Request.RemoteAddr,
			},
		)

	}
}
