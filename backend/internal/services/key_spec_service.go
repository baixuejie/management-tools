package services

import (
	"errors"
	"fmt"

	"github.com/baixuejie/key-management-tool/backend/internal/models"
	"gorm.io/gorm"
)

type KeySpecService struct {
	db *gorm.DB
}

func NewKeySpecService(db *gorm.DB) *KeySpecService {
	return &KeySpecService{db: db}
}

// CreateKeySpec creates a new key specification
func (s *KeySpecService) CreateKeySpec(name, description string) (*models.KeySpec, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	keySpec := &models.KeySpec{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(keySpec).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("key spec with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to create key spec: %w", err)
	}

	return keySpec, nil
}

// GetKeySpec retrieves a key specification by ID
func (s *KeySpecService) GetKeySpec(id uint) (*models.KeySpec, error) {
	var keySpec models.KeySpec
	if err := s.db.First(&keySpec, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("key spec with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get key spec: %w", err)
	}

	return &keySpec, nil
}

// ListKeySpecs retrieves all key specifications (excluding soft deleted)
func (s *KeySpecService) ListKeySpecs() ([]models.KeySpec, error) {
	var keySpecs []models.KeySpec
	if err := s.db.Order("created_at DESC").Find(&keySpecs).Error; err != nil {
		return nil, fmt.Errorf("failed to list key specs: %w", err)
	}

	return keySpecs, nil
}

// UpdateKeySpec updates an existing key specification
func (s *KeySpecService) UpdateKeySpec(id uint, name, description string) (*models.KeySpec, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	var keySpec models.KeySpec
	if err := s.db.First(&keySpec, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("key spec with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find key spec: %w", err)
	}

	keySpec.Name = name
	keySpec.Description = description

	if err := s.db.Save(&keySpec).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("key spec with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to update key spec: %w", err)
	}

	return &keySpec, nil
}

// DeleteKeySpec soft deletes a key specification
func (s *KeySpecService) DeleteKeySpec(id uint) error {
	result := s.db.Delete(&models.KeySpec{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete key spec: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("key spec with ID %d not found", id)
	}

	return nil
}
