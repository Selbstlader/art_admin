package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// DesignVersionRepository 设计版本仓库
// Design version repository
type DesignVersionRepository struct {
	db *gorm.DB
}

// NewDesignVersionRepository 创建设计版本仓库
// Create design version repository
func NewDesignVersionRepository(db *gorm.DB) *DesignVersionRepository {
	return &DesignVersionRepository{db: db}
}

// Create 创建设计版本
// Create design version
func (r *DesignVersionRepository) Create(version *model.DesignVersion) error {
	return r.db.Create(version).Error
}

// GetByID 根据ID获取设计版本
// Get design version by ID
func (r *DesignVersionRepository) GetByID(id uint) (*model.DesignVersion, error) {
	var version model.DesignVersion
	err := r.db.First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetByIDWithProject 根据ID获取设计版本（包含项目信息）
// Get design version by ID with project info
func (r *DesignVersionRepository) GetByIDWithProject(id uint) (*model.DesignVersion, error) {
	var version model.DesignVersion
	err := r.db.Preload("Project").First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// DesignVersionFilter 设计版本筛选条件
// Design version filter
type DesignVersionFilter struct {
	ProjectID uint
	Status    string
	CreatedBy uint
}

// List 获取设计版本列表
// Get design version list
func (r *DesignVersionRepository) List(page, pageSize int, filter DesignVersionFilter) ([]model.DesignVersion, int64, error) {
	var versions []model.DesignVersion
	var total int64

	query := r.db.Model(&model.DesignVersion{})

	// 应用筛选条件 / Apply filters
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CreatedBy > 0 {
		query = query.Where("created_by = ?", filter.CreatedBy)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Order("version_number DESC").Offset(offset).Limit(pageSize).Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetByProjectID 根据项目ID获取所有版本
// Get all versions by project ID
func (r *DesignVersionRepository) GetByProjectID(projectID uint) ([]model.DesignVersion, error) {
	var versions []model.DesignVersion
	err := r.db.Where("project_id = ?", projectID).Order("version_number DESC").Find(&versions).Error
	return versions, err
}

// GetLatestVersionNumber 获取项目最新版本号
// Get latest version number for project
func (r *DesignVersionRepository) GetLatestVersionNumber(projectID uint) (int, error) {
	var maxVersion int
	err := r.db.Model(&model.DesignVersion{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion).Error
	return maxVersion, err
}

// Update 更新设计版本
// Update design version
func (r *DesignVersionRepository) Update(version *model.DesignVersion) error {
	return r.db.Save(version).Error
}

// UpdateFields 更新指定字段
// Update specific fields
func (r *DesignVersionRepository) UpdateFields(id uint, updates map[string]any) error {
	return r.db.Model(&model.DesignVersion{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除设计版本
// Delete design version
func (r *DesignVersionRepository) Delete(id uint) error {
	return r.db.Delete(&model.DesignVersion{}, id).Error
}

// DeleteByProjectID 根据项目ID删除所有版本
// Delete all versions by project ID
func (r *DesignVersionRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.DesignVersion{}).Error
}

// CountByProjectID 统计项目版本数量
// Count versions by project ID
func (r *DesignVersionRepository) CountByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.DesignVersion{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}
