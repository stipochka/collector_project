package handler

import (
	"dashboard/internal/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLogs(c *gin.Context) {

}

func ParseRequestURLFilter(c *gin.Context) (models.LogFilter, error) {
	const op = "handler.ParseRequestURLFilter"

	var filter models.LogFilter

	if levelRaw, exist := c.Get("level"); exist {
		if level, ok := levelRaw.(string); ok {
			filter.Level = level
		} else {
			return filter, fmt.Errorf("%s: failed to convert level to string", op)
		}

	}

	if serviceNameRaw, exists := c.Get("service_name"); exists {
		if serviceName, ok := serviceNameRaw.(string); ok {
			filter.ServiceName = serviceName
		} else {
			return filter, fmt.Errorf("%s: failed to convert service name", op)
		}
	}

	if opRaw, exists := c.Get("op"); exists {
		if op, ok := opRaw.(string); ok {
			filter.Op = op
		} else {
			return filter, fmt.Errorf("%s: failed to convert op", op)
		}
	}

	if messageLikeRaw, exists := c.Get("message_like"); exists {
		if messageLike, ok := messageLikeRaw.(string); ok {
			filter.MessageLike = messageLike
		} else {
			return filter, fmt.Errorf("%s: failed to convert message like", op)
		}
	}

	if fromRaw, exists := c.Get("from"); exists {
		if from, ok := fromRaw.(time.Time); ok {
			filter.From = from
		} else {
			return filter, fmt.Errorf("%s: failed to convert from", op)
		}
	}

	if toRaw, exists := c.Get("to"); exists {
		if to, ok := toRaw.(time.Time); ok {
			filter.To = to
		} else {
			return filter, fmt.Errorf("%s: failed to convert to", op)
		}
	}

	if limitRaw, exists := c.Get("limit"); exists {
		if limit, ok := limitRaw.(int); ok {
			filter.Limit = limit
		} else {
			return filter, fmt.Errorf("%s: failed to convert limit", op)
		}
	}

	if offsetRaw, exists := c.Get("from"); exists {
		if offset, ok := offsetRaw.(int); ok {
			filter.Offset = offset
		} else {
			return filter, fmt.Errorf("%s: failed to convert from", op)
		}
	}

	return filter, nil
}
