package handlers

import (
	"net/http"
	"strconv"

	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KeySpecHandler struct {
	service *services.KeySpecService
}

func NewKeySpecHandler(service *services.KeySpecService) *KeySpecHandler {
	return &KeySpecHandler{service: service}
}

// CreateKeySpec handles POST /api/key-specs
func (h *KeySpecHandler) CreateKeySpec(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keySpec, err := h.service.CreateKeySpec(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          keySpec.ID,
		"name":        keySpec.Name,
		"description": keySpec.Description,
		"created_at":  keySpec.CreatedAt,
	})
}

// ListKeySpecs handles GET /api/key-specs
func (h *KeySpecHandler) ListKeySpecs(c *gin.Context) {
	keySpecs, err := h.service.ListKeySpecs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]gin.H, len(keySpecs))
	for i, spec := range keySpecs {
		response[i] = gin.H{
			"id":          spec.ID,
			"name":        spec.Name,
			"description": spec.Description,
			"created_at":  spec.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetKeySpec handles GET /api/key-specs/:id
func (h *KeySpecHandler) GetKeySpec(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	keySpec, err := h.service.GetKeySpec(uint(id))
	if err != nil {
		if err.Error() == "key spec with ID "+idStr+" not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          keySpec.ID,
		"name":        keySpec.Name,
		"description": keySpec.Description,
		"created_at":  keySpec.CreatedAt,
	})
}

// UpdateKeySpec handles PUT /api/key-specs/:id
func (h *KeySpecHandler) UpdateKeySpec(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keySpec, err := h.service.UpdateKeySpec(uint(id), req.Name, req.Description)
	if err != nil {
		if err.Error() == "key spec with ID "+idStr+" not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "name cannot be empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          keySpec.ID,
		"name":        keySpec.Name,
		"description": keySpec.Description,
		"updated_at":  keySpec.UpdatedAt,
	})
}

// DeleteKeySpec handles DELETE /api/key-specs/:id
func (h *KeySpecHandler) DeleteKeySpec(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.service.DeleteKeySpec(uint(id))
	if err != nil {
		if err.Error() == "key spec with ID "+idStr+" not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key spec deleted successfully"})
}
