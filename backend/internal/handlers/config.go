package handlers

import (
	"net/http"

	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	service *services.ConfigService
}

func NewConfigHandler(service *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: service}
}

// GetCopyTemplate handles GET /api/config/copy-template
func (h *ConfigHandler) GetCopyTemplate(c *gin.Context) {
	template, err := h.service.GetCopyTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// UpdateCopyTemplate handles PUT /api/config/copy-template
func (h *ConfigHandler) UpdateCopyTemplate(c *gin.Context) {
	var req struct {
		Template string `json:"template" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateCopyTemplate(req.Template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Copy template updated successfully"})
}
