package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/baixuejie/key-management-tool/backend/internal/models"
	"github.com/baixuejie/key-management-tool/backend/internal/utils"
)

type KeyService struct {
	db *gorm.DB
}

func NewKeyService(db *gorm.DB) *KeyService {
	return &KeyService{db: db}
}

// BatchCreateKeys creates multiple keys for a spec
// All keys start with IsUsed=false, UsedAt=nil
func (s *KeyService) BatchCreateKeys(specID uint, keyValues []string) error {
	// Validate specID exists
	var spec models.KeySpec
	if err := s.db.First(&spec, specID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("key spec not found")
		}
		return err
	}

	// Encrypt each key value and prepare for batch insert
	keys := make([]models.Key, 0, len(keyValues))
	for _, keyValue := range keyValues {
		encrypted, err := utils.Encrypt(keyValue)
		if err != nil {
			return err
		}

		keys = append(keys, models.Key{
			SpecID:   specID,
			KeyValue: encrypted,
			IsUsed:   false,
			UsedAt:   nil,
		})
	}

	// Batch insert keys
	if err := s.db.Create(&keys).Error; err != nil {
		return err
	}

	return nil
}

// GetAvailableKey gets one unused key for the spec
// If no unused keys, get most recently used key
func (s *KeyService) GetAvailableKey(specID uint) (*models.Key, error) {
	var key models.Key

	// Try to get an unused key first (ordered by created_at ASC)
	err := s.db.Where("spec_id = ? AND is_used = ?", specID, false).
		Order("created_at ASC").
		First(&key).Error

	if err == nil {
		// Found an unused key, decrypt and return
		decrypted, err := utils.Decrypt(key.KeyValue)
		if err != nil {
			return nil, err
		}
		key.KeyValue = decrypted
		return &key, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// No unused keys, get most recently used key
	err = s.db.Where("spec_id = ? AND is_used = ?", specID, true).
		Order("used_at DESC").
		First(&key).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no keys found for this spec")
		}
		return nil, err
	}

	// Decrypt key value before returning
	decrypted, err := utils.Decrypt(key.KeyValue)
	if err != nil {
		return nil, err
	}
	key.KeyValue = decrypted
	return &key, nil
}

// MarkKeyAsUsed marks a key as used with current timestamp
func (s *KeyService) MarkKeyAsUsed(keyID uint) error {
	now := time.Now()
	result := s.db.Model(&models.Key{}).
		Where("id = ?", keyID).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("key not found")
	}

	return nil
}

// ListKeys lists keys for a spec with pagination
func (s *KeyService) ListKeys(specID uint, onlyUnused bool, limit, offset int) ([]models.Key, int64, error) {
	var keys []models.Key
	var total int64

	query := s.db.Model(&models.Key{}).Where("spec_id = ?", specID)

	if onlyUnused {
		query = query.Where("is_used = ?", false)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with ordering
	// Order: unused first, then by used_at DESC, then by created_at DESC
	err := query.Preload("Spec").
		Order("is_used ASC").
		Order("used_at DESC").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&keys).Error

	if err != nil {
		return nil, 0, err
	}

	// Decrypt all key values
	for i := range keys {
		decrypted, err := utils.Decrypt(keys[i].KeyValue)
		if err != nil {
			return nil, 0, err
		}
		keys[i].KeyValue = decrypted
	}

	return keys, total, nil
}

// DeleteKey hard deletes a key from database
func (s *KeyService) DeleteKey(keyID uint) error {
	result := s.db.Unscoped().Delete(&models.Key{}, keyID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("key not found")
	}

	return nil
}
