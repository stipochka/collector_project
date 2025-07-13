package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}

func (h *Handler) GetLogs(c *gin.Context) {
	const handlerName = "GetLogs"

	log := h.log.With(zap.String("handler_name", handlerName))

	filter, err := h.FormLogFilterFromURL(c)
	if err != nil {
		log.Error("failed to get filter from request",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusBadRequest, "invalid request params")
		return
	}

	logsRes, err := h.service.GetLogs(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to get logs",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	log.Info("successfully get logs",
		zap.Int("length", len(logsRes)),
	)

	c.JSON(http.StatusOK, logsRes)
}

func (h *Handler) GetLogLevelStats(c *gin.Context) {
	const handlerName = "GetLogLevelStats"

	log := h.log.With(zap.String("handler_name", handlerName))

	filter, err := h.FormTimeRangeFilter(c)
	if err != nil {
		log.Error("failed to parse params",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	stats, err := h.service.GetLogLevelStats(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to get logs stats",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	log.Info("successfully go log stats")

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetServiceStats(c *gin.Context) {
	const handlerName = "GetServiceStats"

	log := h.log.With(zap.String("handler_name", handlerName))

	filter, err := h.FormTimeRangeFilter(c)
	if err != nil {
		log.Error("failed to parse params",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	stats, err := h.service.GetServiceStats(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to get service stats",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	log.Info("successfully get service stats")

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetTopErrors(c *gin.Context) {
	const handlerName = "GetTopErrors"

	log := h.log.With(zap.String("handler_name", handlerName))

	filter, err := h.FormTimeRangeFilter(c)
	if err != nil {
		log.Error("failed to parse params",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	topErrors, err := h.service.GetTopErrors(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to get top errors",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")

		return
	}

	log.Info("successfully get top errors")

	c.JSON(http.StatusOK, topErrors)
}

func (h *Handler) GetRecentErrors(c *gin.Context) {
	const handlerName = "GetRecentErrors"

	log := h.log.With(zap.String("handler_name", handlerName))

	filter, err := h.FormTimeRangeFilter(c)
	if err != nil {
		log.Error("failed to parse params",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	recentErrors, err := h.service.GetRecentErrors(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to get recent errors",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")

		return
	}

	log.Info("successfully get recent errors")

	c.JSON(http.StatusOK, recentErrors)
}

func (h *Handler) GetAllLogLevels(c *gin.Context) {
	const handlerName = "GetAllLogLevels"

	log := h.log.With(zap.String("handler_name", handlerName))

	levels, err := h.service.GetAllLogLevels(c.Request.Context())
	if err != nil {
		log.Error("failed to get log levels",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	log.Info("successfully get all log levels")

	c.JSON(http.StatusOK, levels)
}

func (h *Handler) GetAllServiceNames(c *gin.Context) {
	const handlerName = "GetAllServiceNames"

	log := h.log.With(zap.String("handler_name", handlerName))

	services, err := h.service.GetAllServiceNames(c.Request.Context())
	if err != nil {
		log.Error("failed to get service names",
			zap.Any("error", err),
		)

		newErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	log.Info("successfully get all service names")

	c.JSON(http.StatusOK, services)
}
