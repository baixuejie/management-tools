package main

import (
	"fmt"
	"log"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/baixuejie/key-management-tool/backend/internal/database"
	"github.com/baixuejie/key-management-tool/backend/internal/handlers"
	"github.com/baixuejie/key-management-tool/backend/internal/middleware"
	"github.com/baixuejie/key-management-tool/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	keySpecService := services.NewKeySpecService(database.DB)
	keyService := services.NewKeyService(database.DB)
	configService := services.NewConfigService(database.DB)

	// Initialize handlers
	keySpecHandler := handlers.NewKeySpecHandler(keySpecService)
	keyHandler := handlers.NewKeyHandler(keyService)
	configHandler := handlers.NewConfigHandler(configService)

	// Setup router
	router := gin.Default()

	// CORS middleware
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

	// Public routes
	router.POST("/api/login", handlers.Login)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Key specs
		api.POST("/key-specs", keySpecHandler.CreateKeySpec)
		api.GET("/key-specs", keySpecHandler.ListKeySpecs)
		api.GET("/key-specs/:id", keySpecHandler.GetKeySpec)
		api.PUT("/key-specs/:id", keySpecHandler.UpdateKeySpec)
		api.DELETE("/key-specs/:id", keySpecHandler.DeleteKeySpec)

		// Keys
		api.POST("/keys/batch", keyHandler.BatchUploadKeys)
		api.GET("/keys/available/:spec_id", keyHandler.GetAvailableKey)
		api.PUT("/keys/:id/use", keyHandler.MarkKeyAsUsed)
		api.GET("/keys", keyHandler.ListKeys)
		api.DELETE("/keys/:id", keyHandler.DeleteKey)

		// Config
		api.GET("/config/copy-template", configHandler.GetCopyTemplate)
		api.PUT("/config/copy-template", configHandler.UpdateCopyTemplate)
	}

	// Start server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
