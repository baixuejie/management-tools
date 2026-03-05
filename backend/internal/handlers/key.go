package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KeyHandler struct {
	service *services.KeyService
}

func NewKeyHandler(service *services.KeyService) *KeyHandler {
	return &KeyHandler{service: service}
}

// BatchUploadKeys handles POST /api/keys/batch
func (h *KeyHandler) BatchUploadKeys(c *gin.Context) {
	var req struct {
		SpecID uint     `json:"spec_id" binding:"required"`
		Keys   []string `json:"keys" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Keys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keys array cannot be empty"})
		return
	}

	err := h.service.BatchCreateKeys(req.SpecID, req.Keys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%d keys uploaded successfully", len(req.Keys)),
	})
}

// GetAvailableKey handles GET /api/keys/available/:spec_id
func (h *KeyHandler) GetAvailableKey(c *gin.Context) {
	specIDStr := c.Param("spec_id")
	specID, err := strconv.ParseUint(specIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_id"})
		return
	}

	key, err := h.service.GetAvailableKey(uint(specID))
	if err != nil {
		if err.Error() == "no keys found for this spec" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        key.ID,
		"spec_id":   key.SpecID,
		"key_value": key.KeyValue,
		"is_used":   key.IsUsed,
		"used_at":   key.UsedAt,
	})
}

// MarkKeyAsUsed handles PUT /api/keys/:id/use
func (h *KeyHandler) MarkKeyAsUsed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.service.MarkKeyAsUsed(uint(id))
	if err != nil {
		if err.Error() == "key not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key marked as used"})
}

// ListKeys handles GET /api/keys?spec_id=X&only_unused=true&limit=20&offset=0
func (h *KeyHandler) ListKeys(c *gin.Context) {
	specIDStr := c.Query("spec_id")
	if specIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "spec_id is required"})
		return
	}

	specID, err := strconv.ParseUint(specIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spec_id"})
		return
	}

	onlyUnused := c.Query("only_unused") == "true"

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	keys, total, err := h.service.ListKeys(uint(specID), onlyUnused, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keys":  keys,
		"total": total,
	})
}

// DeleteKey handles DELETE /api/keys/:id
func (h *KeyHandler) DeleteKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.service.DeleteKey(uint(id))
	if err != nil {
		if err.Error() == "key not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}
