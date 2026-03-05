# Key Management Tool Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a personal tools website with key management as the first tool, supporting batch upload, custom copy templates, and usage tracking.

**Architecture:** Microservices with Docker Compose (MySQL, Go backend API, Vue frontend, Nginx gateway). Mobile-first responsive design with long session support.

**Tech Stack:** Go 1.25 + Gin + GORM + JWT, Vue 3 + Vant UI, MySQL 8.0, Docker Compose

---

## Phase 1: Project Structure & Database

### Task 1: Initialize Backend Structure

**Files:**
- Create: `backend/cmd/main.go`
- Create: `backend/go.mod`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/models/models.go`

**Step 1: Create backend directory structure**

Run: `mkdir -p backend/cmd backend/internal/{config,models,handlers,middleware,services,utils}`

**Step 2: Initialize Go module**

Run: `cd backend && go mod init github.com/baixuejie/key-management-tool`

**Step 3: Install dependencies**

Run: `cd backend && go get github.com/gin-gonic/gin github.com/golang-jwt/jwt/v5 gorm.io/gorm gorm.io/driver/mysql github.com/spf13/viper golang.org/x/crypto/bcrypt`

**Step 4: Create basic main.go**

File: `backend/cmd/main.go`
```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

**Step 5: Test server starts**

Run: `cd backend && go run cmd/main.go`
Expected: Server starts on :8080

**Step 6: Commit**

```bash
git add backend/
git commit -m "feat: initialize backend structure with basic server"
```

### Task 2: Database Models

**Files:**
- Create: `backend/internal/models/key_spec.go`
- Create: `backend/internal/models/key.go`
- Create: `backend/internal/models/config.go`

**Step 1: Create KeySpec model**

File: `backend/internal/models/key_spec.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type KeySpec struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Step 2: Create Key model**

File: `backend/internal/models/key.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Key struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	SpecID    uint           `gorm:"not null;index:idx_spec_used_time" json:"spec_id"`
	KeyValue  string         `gorm:"type:text;not null" json:"key_value"`
	IsUsed    bool           `gorm:"default:false;index:idx_spec_used_time" json:"is_used"`
	UsedAt    *time.Time     `gorm:"index:idx_spec_used_time" json:"used_at"`
	CreatedAt time.Time      `gorm:"index:idx_spec_created" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Spec      KeySpec        `gorm:"foreignKey:SpecID" json:"spec,omitempty"`
}
```

**Step 3: Create Config model**

File: `backend/internal/models/config.go`
```go
package models

import "time"

type Config struct {
	Key         string    `gorm:"primarykey;size:100" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Description string    `gorm:"size:255" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

**Step 4: Commit**

```bash
git add backend/internal/models/
git commit -m "feat: add database models for key management"
```

### Task 3: Configuration Management

**Files:**
- Create: `backend/config.yaml`
- Modify: `backend/internal/config/config.go`

**Step 1: Create config.yaml template**

File: `backend/config.yaml`
```yaml
server:
  port: 8080

database:
  host: mysql
  port: 3306
  user: root
  password: rootpassword
  dbname: key_management

jwt:
  secret: "change-this-to-a-random-secret-key-in-production"
  expiry_hours: 168  # 7 days

encryption:
  key: "change-this-to-32-byte-aes-key-prod"  # Must be 32 bytes for AES-256

auth:
  username: admin
  # bcrypt hash of "admin123"
  password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
```

**Step 2: Create config loader**

File: `backend/internal/config/config.go`
```go
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	Auth       AuthConfig       `mapstructure:"auth"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpiryHours int    `mapstructure:"expiry_hours"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type AuthConfig struct {
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

**Step 3: Commit**

```bash
git add backend/config.yaml backend/internal/config/
git commit -m "feat: add configuration management with YAML support"
```

### Task 4: Database Connection

**Files:**
- Create: `backend/internal/database/database.go`
- Modify: `backend/cmd/main.go`

**Step 1: Create database connection**

File: `backend/internal/database/database.go`
```go
package database

import (
	"fmt"
	"log"

	"github.com/baixuejie/key-management-tool/internal/config"
	"github.com/baixuejie/key-management-tool/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.KeySpec{},
		&models.Key{},
		&models.Config{},
	)
}
```

**Step 2: Update main.go to connect database**

File: `backend/cmd/main.go`
```go
package main

import (
	"log"

	"github.com/baixuejie/key-management-tool/internal/config"
	"github.com/baixuejie/key-management-tool/internal/database"
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

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

**Step 3: Commit**

```bash
git add backend/internal/database/ backend/cmd/main.go
git commit -m "feat: add database connection and auto-migration"
```

## Phase 2: Authentication & Security

### Task 5: JWT Utilities

**Files:**
- Create: `backend/internal/utils/jwt.go`

**Step 1: Create JWT utility functions**

File: `backend/internal/utils/jwt.go`
```go
package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, secret string, expiryHours int) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
```

**Step 2: Commit**

```bash
git add backend/internal/utils/jwt.go
git commit -m "feat: add JWT token generation and validation"
```

### Task 6: Encryption Utilities

**Files:**
- Create: `backend/internal/utils/crypto.go`

**Step 1: Create AES encryption utilities**

File: `backend/internal/utils/crypto.go`
```go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(plaintext string, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
```

**Step 2: Commit**

```bash
git add backend/internal/utils/crypto.go
git commit -m "feat: add AES-256 encryption utilities"
```

### Task 7: Authentication Handler

**Files:**
- Create: `backend/internal/handlers/auth.go`
- Create: `backend/internal/middleware/auth.go`

**Step 1: Create auth handler**

File: `backend/internal/handlers/auth.go`
```go
package handlers

import (
	"net/http"

	"github.com/baixuejie/key-management-tool/internal/config"
	"github.com/baixuejie/key-management-tool/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify username
	if req.Username != h.cfg.Auth.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(h.cfg.Auth.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	expiryHours := h.cfg.JWT.ExpiryHours
	if req.RememberMe {
		expiryHours = 720 // 30 days
	}

	token, err := utils.GenerateToken(req.Username, h.cfg.JWT.Secret, expiryHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
```

**Step 2: Create auth middleware**

File: `backend/internal/middleware/auth.go`
