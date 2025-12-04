package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// FileRepository 文件仓库
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create 创建文件记录
func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *FileRepository) GetByID(id uint) (*model.File, error) {
	var file model.File
	err := r.db.First(&file, id).Error
	return &file, err
}

// GetByStorageName 根据存储名获取文件
func (r *FileRepository) GetByStorageName(storageName string) (*model.File, error) {
	var file model.File
	err := r.db.Where("storage_name = ?", storageName).First(&file).Error
	return &file, err
}

// List 获取文件列表
func (r *FileRepository) List(page, pageSize int, category string, userID uint) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Model(&model.File{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if userID > 0 {
		query = query.Where("upload_user_id = ?", userID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}

// Delete 删除文件记录
func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&model.File{}, id).Error
}

// BatchDelete 批量删除
func (r *FileRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.File{}, ids).Error
}
