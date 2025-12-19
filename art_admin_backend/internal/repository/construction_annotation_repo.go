package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// ConstructionAnnotationRepository 施工图标注数据仓库接口
// Construction annotation repository interface
type ConstructionAnnotationRepository interface {
	Create(annotation *model.ConstructionAnnotation) error
	GetByID(id uint) (*model.ConstructionAnnotation, error)
	GetByProjectID(projectID uint) ([]*model.ConstructionAnnotation, error)
	GetByCadFileID(cadFileID uint) ([]*model.ConstructionAnnotation, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.ConstructionAnnotation, int64, error)
}

// constructionAnnotationRepo 施工图标注数据仓库实现
// Construction annotation repository implementation
type constructionAnnotationRepo struct {
	db *gorm.DB
}

// NewConstructionAnnotationRepository 创建施工图标注数据仓库
// Create construction annotation repository
func NewConstructionAnnotationRepository(db *gorm.DB) ConstructionAnnotationRepository {
	return &constructionAnnotationRepo{db: db}
}

// Create 创建标注记录
// Create annotation record
func (r *constructionAnnotationRepo) Create(annotation *model.ConstructionAnnotation) error {
	return r.db.Create(annotation).Error
}

// GetByID 根据ID获取标注记录
// Get annotation record by ID
func (r *constructionAnnotationRepo) GetByID(id uint) (*model.ConstructionAnnotation, error) {
	var annotation model.ConstructionAnnotation
	err := r.db.Preload("Project").Preload("CadFile").First(&annotation, id).Error
	if err != nil {
		return nil, err
	}
	return &annotation, nil
}

// GetByProjectID 根据项目ID获取标注记录
// Get annotation records by project ID
func (r *constructionAnnotationRepo) GetByProjectID(projectID uint) ([]*model.ConstructionAnnotation, error) {
	var annotations []*model.ConstructionAnnotation
	err := r.db.Where("project_id = ?", projectID).
		Preload("Project").
		Preload("CadFile").
		Order("created_at DESC").
		Find(&annotations).Error
	return annotations, err
}

// GetByCadFileID 根据CAD文件ID获取标注记录
// Get annotation records by CAD file ID
func (r *constructionAnnotationRepo) GetByCadFileID(cadFileID uint) ([]*model.ConstructionAnnotation, error) {
	var annotations []*model.ConstructionAnnotation
	err := r.db.Where("cad_file_id = ?", cadFileID).
		Preload("Project").
		Preload("CadFile").
		Order("created_at DESC").
		Find(&annotations).Error
	return annotations, err
}

// Update 更新标注记录
// Update annotation record
func (r *constructionAnnotationRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.ConstructionAnnotation{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除标注记录
// Delete annotation record
func (r *constructionAnnotationRepo) Delete(id uint) error {
	return r.db.Delete(&model.ConstructionAnnotation{}, id).Error
}

// List 分页获取标注记录列表
// Get paginated list of annotation records
func (r *constructionAnnotationRepo) List(offset, limit int) ([]*model.ConstructionAnnotation, int64, error) {
	var annotations []*model.ConstructionAnnotation
	var total int64

	// 获取总数 / Get total count
	if err := r.db.Model(&model.ConstructionAnnotation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据 / Get paginated data
	err := r.db.Preload("Project").
		Preload("CadFile").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&annotations).Error

	return annotations, total, err
}
