package handler

import (
	"dashboard/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FormLogFilterFromURL(c *gin.Context) (models.LogFilter, error) {
	const op = "handler.FormLogFilterFromURL"

	var filter models.LogFilter

	if level := c.Query(""); level != "" {

	}
}
