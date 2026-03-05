package services

import (
	"errors"
	"fmt"
	"strings"

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
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	var maxOrder int
	if err := s.db.Model(&models.KeySpec{}).Select("COALESCE(MAX(display_order), 0)").Scan(&maxOrder).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate display order: %w", err)
	}

	keySpec := &models.KeySpec{
		Name:         name,
		Description:  strings.TrimSpace(description),
		DisplayOrder: maxOrder + 1,
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
	if err := s.db.Order("display_order ASC").Order("created_at DESC").Find(&keySpecs).Error; err != nil {
		return nil, fmt.Errorf("failed to list key specs: %w", err)
	}

	return keySpecs, nil
}

// UpdateKeySpec updates an existing key specification
func (s *KeySpecService) UpdateKeySpec(id uint, name, description string) (*models.KeySpec, error) {
	name = strings.TrimSpace(name)
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
	keySpec.Description = strings.TrimSpace(description)

	if err := s.db.Save(&keySpec).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("key spec with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to update key spec: %w", err)
	}

	return &keySpec, nil
}

// ReorderKeySpecs updates key spec display order by ID sequence.
func (s *KeySpecService) ReorderKeySpecs(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("ids cannot be empty")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var keySpecs []models.KeySpec
		if err := tx.Order("display_order ASC").Order("created_at DESC").Find(&keySpecs).Error; err != nil {
			return fmt.Errorf("failed to list key specs: %w", err)
		}
		if len(keySpecs) == 0 {
			return nil
		}

		existing := make(map[uint]models.KeySpec, len(keySpecs))
		for _, spec := range keySpecs {
			existing[spec.ID] = spec
		}

		used := make(map[uint]bool, len(ids))
		order := 1
		for _, id := range ids {
			if used[id] {
				continue
			}
			spec, ok := existing[id]
			if !ok {
				return fmt.Errorf("key spec with ID %d not found", id)
			}
			if err := tx.Model(&models.KeySpec{}).Where("id = ?", spec.ID).Update("display_order", order).Error; err != nil {
				return fmt.Errorf("failed to update key spec order: %w", err)
			}
			used[id] = true
			order++
		}

		for _, spec := range keySpecs {
			if used[spec.ID] {
				continue
			}
			if err := tx.Model(&models.KeySpec{}).Where("id = ?", spec.ID).Update("display_order", order).Error; err != nil {
				return fmt.Errorf("failed to update key spec order: %w", err)
			}
			order++
		}

		return nil
	})
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
