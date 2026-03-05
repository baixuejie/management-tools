package database

import (
	"fmt"
	"log"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/baixuejie/key-management-tool/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		&models.User{},
		&models.KeySpec{},
		&models.Key{},
		&models.Config{},
		&models.CostRecord{},
		&models.Customer{},
		&models.Transaction{},
	)
}

func SeedDefaults() error {
	defaultUsers := []models.User{
		{
			Username:     "admin",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			DisplayName:  "White",
		},
		{
			Username:     "fanchen",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			DisplayName:  "Fanchen",
		},
	}

	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoNothing: true,
	}).Create(&defaultUsers).Error; err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	defaultTemplate := models.Config{
		Key:         "copy_template",
		Value:       "{{key}}",
		Description: "Global copy template for key values",
	}
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoNothing: true,
	}).Create(&defaultTemplate).Error; err != nil {
		return fmt.Errorf("failed to seed config defaults: %w", err)
	}

	return nil
}
