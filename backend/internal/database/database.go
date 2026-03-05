package database

import (
	"fmt"
	"log"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/baixuejie/key-management-tool/backend/internal/models"
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

// Seed inserts default data if not already present.
func Seed() error {
	var count int64
	DB.Model(&models.Config{}).Where("`key` = ?", "copy_template").Count(&count)

	if count == 0 {
		return DB.Create(&models.Config{
			Key:         "copy_template",
			Value:       "{{key}}",
			Description: "Global copy template for key values",
		}).Error
	}
	return nil
}
