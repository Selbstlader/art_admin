package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"strings"

	"gorm.io/gorm"
)

// AITagRepository handles database operations for AI tags
type AITagRepository struct {
	db *gorm.DB
}

// NewAITagRepository creates a new AITagRepository instance
func NewAITagRepository() *AITagRepository {
	return &AITagRepository{
		db: database.DB,
	}
}

// NewAITagRepositoryWithDB creates a new AITagRepository with a custom DB (for testing)
func NewAITagRepositoryWithDB(db *gorm.DB) *AITagRepository {
	return &AITagRepository{
		db: db,
	}
}

// getDB returns the database connection
func (r *AITagRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return database.DB
}

// Create creates a new AI tag in the database
// Requirements: 1.1 - Create new AI_Tag record
func (r *AITagRepository) Create(tag *model.AITag) error {
	return r.getDB().Create(tag).Error
}

// Update updates an existing AI tag
// Requirements: 3.1 - Update AI_Tag record
func (r *AITagRepository) Update(tag *model.AITag) error {
	return r.getDB().Save(tag).Error
}

// Delete performs a soft delete on an AI tag by setting status to inactive
// Requirements: 4.1 - Soft delete by setting status to inactive
func (r *AITagRepository) Delete(id int64) error {
	return r.getDB().Model(&model.AITag{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": 0,
		}).Error
}

// HardDelete performs a hard delete (with GORM soft delete) on an AI tag
func (r *AITagRepository) HardDelete(id int64) error {
	return r.getDB().Delete(&model.AITag{}, id).Error
}

// FindByID finds an AI tag by its ID (excludes soft-deleted records)
// Requirements: 1.1 - Return created tag details
func (r *AITagRepository) FindByID(id int64) (*model.AITag, error) {
	var tag model.AITag
	err := r.getDB().First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName finds an AI tag by its name (excludes soft-deleted records)
// Requirements: 1.2 - Check for duplicate names
func (r *AITagRepository) FindByName(name string) (*model.AITag, error) {
	var tag model.AITag
	err := r.getDB().Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindActiveByID finds an active AI tag by its ID (status=1, excludes soft-deleted)
// Requirements: 5.1 - Get tag config for chat integration
func (r *AITagRepository) FindActiveByID(id int64) (*model.AITag, error) {
	var tag model.AITag
	err := r.getDB().Where("id = ? AND status = ?", id, 1).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// List returns paginated AI tags with optional keyword search
// Requirements: 2.1 - Return all active AI_Tag records with pagination
// Requirements: 2.2 - Search tags by keyword in name or description
func (r *AITagRepository) List(page, pageSize int, keyword string) ([]*model.AITag, int64, error) {
	var tags []*model.AITag
	var total int64

	db := r.getDB().Model(&model.AITag{})

	// Apply keyword search filter (case-insensitive search in name or description)
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// FindByNameExcludingID finds a tag by name excluding a specific ID (for update validation)
// Requirements: 3.2 - Check for name conflicts during update
func (r *AITagRepository) FindByNameExcludingID(name string, excludeID int64) (*model.AITag, error) {
	var tag model.AITag
	err := r.getDB().Where("name = ? AND id != ?", name, excludeID).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindIncludingDeleted finds a tag by ID including soft-deleted records
// Used for verifying soft delete behavior in tests
func (r *AITagRepository) FindIncludingDeleted(id int64) (*model.AITag, error) {
	var tag model.AITag
	err := r.getDB().Unscoped().First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// CountAll returns the total count of all non-deleted tags
func (r *AITagRepository) CountAll() (int64, error) {
	var count int64
	err := r.getDB().Model(&model.AITag{}).Count(&count).Error
	return count, err
}

// GetDB returns the underlying database connection for testing purposes
func (r *AITagRepository) GetDB() *gorm.DB {
	return r.getDB()
}
