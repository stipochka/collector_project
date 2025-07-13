package handler

import (
	"dashboard/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const timeLayout = "2006-01-02"

func (h *Handler) FormLogFilterFromURL(c *gin.Context) (models.LogFilter, error) {
	const op = "handler.FormLogFilterFromURL"

	var filter models.LogFilter

	if level := c.Query("level"); level != "" {
		filter.Level = level
	}

	if serviceName := c.Query("service_name"); serviceName != "" {
		filter.ServiceName = serviceName
	}

	if opName := c.Query("op"); op != "" {
		filter.Op = opName
	}

	if messageLike := c.Query("message_like"); messageLike != "" {
		filter.MessageLike = messageLike
	}

	if fromStr := c.Query("from"); fromStr != "" {
		from, err := time.Parse(timeLayout, fromStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		filter.From = from
	}

	if toStr := c.Query("to"); toStr != "" {
		to, err := time.Parse(timeLayout, toStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		filter.To = to
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		if limit < 0 {
			return filter, fmt.Errorf("%s: limit out of range", op)
		}

		filter.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		if offset < 0 {
			return filter, fmt.Errorf("%s: offset out of range", op)
		}

		filter.Offset = offset
	}

	return filter, nil
}

func (h *Handler) FormTimeRangeFilter(c *gin.Context) (models.TimeRangeFilter, error) {
	const op = "handler.FormTimeRangeFilter"

	var filter models.TimeRangeFilter

	if fromStr := c.Query("from"); fromStr != "" {
		from, err := time.Parse(timeLayout, fromStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		filter.TimeFrom = from
	}

	if toStr := c.Query("to"); toStr != "" {
		to, err := time.Parse(timeLayout, toStr)
		if err != nil {
			return filter, fmt.Errorf("%s: %w", op, err)
		}

		filter.TimeTo = to
	}

	return filter, nil
}
