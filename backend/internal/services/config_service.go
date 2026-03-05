package services

import (
	"errors"

	"github.com/baixuejie/key-management-tool/backend/internal/models"
	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

// GetCopyTemplate retrieves the copy template from the configs table.
// Returns the default template "{{key}}" if not found.
// Returns error only for database errors (not for not found).
func (s *ConfigService) GetCopyTemplate() (string, error) {
	var config models.Config
	err := s.db.Where("`key` = ?", "copy_template").First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return default template if not found
			return "{{key}}", nil
		}
		// Return error for other database errors
		return "", err
	}

	return config.Value, nil
}

// UpdateCopyTemplate updates or creates the copy template configuration.
// Uses upsert pattern to handle both create and update cases.
func (s *ConfigService) UpdateCopyTemplate(template string) error {
	config := models.Config{
		Key:         "copy_template",
		Value:       template,
		Description: "Global copy template for key values",
	}

	// Use Save which will update if exists or create if not
	// First check if record exists
	var existing models.Config
	err := s.db.Where("`key` = ?", "copy_template").First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new record
			return s.db.Create(&config).Error
		}
		// Return error for other database errors
		return err
	}

	// Update existing record
	existing.Value = template
	existing.Description = config.Description
	return s.db.Save(&existing).Error
}
