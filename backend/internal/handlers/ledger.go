package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	service *services.LedgerService
}

func NewLedgerHandler(service *services.LedgerService) *LedgerHandler {
	return &LedgerHandler{service: service}
}

func (h *LedgerHandler) ListCosts(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)

	total, items, err := h.service.ListCosts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": items,
	})
}

func (h *LedgerHandler) CreateCost(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Note   string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.CreateCost(req.Amount, req.Note, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *LedgerHandler) DeleteCost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cost id"})
		return
	}

	if err := h.service.DeleteCost(uint(id)); err != nil {
		if err.Error() == "cost record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cost record deleted"})
}

func (h *LedgerHandler) ListCustomers(c *gin.Context) {
	items, err := h.service.ListCustomers(c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *LedgerHandler) CreateCustomer(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.CreateCustomer(req.Name, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *LedgerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.UpdateCustomer(uint(id), req.Name)
	if err != nil {
		switch err.Error() {
		case "customer not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *LedgerHandler) ListTransactions(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)

	var customerID *uint
	if raw := c.Query("customer_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer_id"})
			return
		}
		converted := uint(id)
		customerID = &converted
	}

	var isNewCustomer *bool
	if raw := c.Query("is_new_customer"); raw != "" {
		switch raw {
		case "1", "true":
			value := true
			isNewCustomer = &value
		case "0", "false":
			value := false
			isNewCustomer = &value
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "is_new_customer must be 0/1 or true/false"})
			return
		}
	}

	total, items, err := h.service.ListTransactions(page, limit, customerID, isNewCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": items,
	})
}

func (h *LedgerHandler) CreateTransaction(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		CustomerID    *uint   `json:"customer_id"`
		CustomerName  string  `json:"customer_name"`
		Amount        float64 `json:"amount" binding:"required"`
		Channel       string  `json:"channel" binding:"required"`
		IsNewCustomer bool    `json:"is_new_customer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.CreateTransaction(services.CreateTransactionInput{
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		Amount:        req.Amount,
		Channel:       req.Channel,
		IsNewCustomer: req.IsNewCustomer,
		RecordedBy:    userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *LedgerHandler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func getCurrentUserID(c *gin.Context) (uint, error) {
	raw, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("missing user_id in token context")
	}
	userID, ok := raw.(uint)
	if !ok || userID == 0 {
		return 0, errors.New("invalid user_id in token context")
	}
	return userID, nil
}

func parsePositiveInt(raw string, defaultValue int) int {
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
