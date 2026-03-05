package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/baixuejie/key-management-tool/backend/internal/database"
	"github.com/baixuejie/key-management-tool/backend/internal/handlers"
	"github.com/baixuejie/key-management-tool/backend/internal/middleware"
	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := resolveConfigPath()
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config (%s): %v", configPath, err)
	}

	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := database.SeedDefaults(); err != nil {
		log.Fatal("Failed to seed default data:", err)
	}

	authService := services.NewAuthService(database.DB)
	keySpecService := services.NewKeySpecService(database.DB)
	keyService := services.NewKeyService(database.DB)
	configService := services.NewConfigService(database.DB)
	ledgerService := services.NewLedgerService(database.DB)

	authHandler := handlers.NewAuthHandler(authService)
	keySpecHandler := handlers.NewKeySpecHandler(keySpecService)
	keyHandler := handlers.NewKeyHandler(keyService)
	configHandler := handlers.NewConfigHandler(configService)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/api/login", authHandler.Login)
	router.POST("/api/auth/login", authHandler.Login)

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/auth/me", authHandler.Me)

		api.POST("/key-specs", keySpecHandler.CreateKeySpec)
		api.GET("/key-specs", keySpecHandler.ListKeySpecs)
		api.PUT("/key-specs/reorder", keySpecHandler.ReorderKeySpecs)
		api.GET("/key-specs/:id", keySpecHandler.GetKeySpec)
		api.PUT("/key-specs/:id", keySpecHandler.UpdateKeySpec)
		api.DELETE("/key-specs/:id", keySpecHandler.DeleteKeySpec)

		api.POST("/keys/batch", keyHandler.BatchUploadKeys)
		api.GET("/keys/available/:spec_id", keyHandler.GetAvailableKey)
		api.PUT("/keys/:id/use", keyHandler.MarkKeyAsUsed)
		api.GET("/keys", keyHandler.ListKeys)
		api.DELETE("/keys/:id", keyHandler.DeleteKey)

		api.GET("/config/copy-template", configHandler.GetCopyTemplate)
		api.PUT("/config/copy-template", configHandler.UpdateCopyTemplate)

		api.GET("/ledger/costs", ledgerHandler.ListCosts)
		api.POST("/ledger/costs", ledgerHandler.CreateCost)
		api.DELETE("/ledger/costs/:id", ledgerHandler.DeleteCost)

		api.GET("/ledger/customers", ledgerHandler.ListCustomers)
		api.POST("/ledger/customers", ledgerHandler.CreateCustomer)
		api.PUT("/ledger/customers/:id", ledgerHandler.UpdateCustomer)

		api.GET("/ledger/transactions", ledgerHandler.ListTransactions)
		api.POST("/ledger/transactions", ledgerHandler.CreateTransaction)

		api.GET("/ledger/statistics", ledgerHandler.GetStatistics)
	}

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func resolveConfigPath() string {
	if path := os.Getenv("KM_CONFIG_PATH"); path != "" {
		return path
	}

	candidates := []string{
		"config.yaml",
		filepath.Join("backend", "config.yaml"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "config.yaml"
}
