package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// DesignVersionCompareRepository 版本对比仓库
// Design version compare repository
type DesignVersionCompareRepository struct {
	db *gorm.DB
}

// NewDesignVersionCompareRepository 创建版本对比仓库
// Create design version compare repository
func NewDesignVersionCompareRepository(db *gorm.DB) *DesignVersionCompareRepository {
	return &DesignVersionCompareRepository{db: db}
}

// Create 创建版本对比记录
// Create version compare record
func (r *DesignVersionCompareRepository) Create(compare *model.DesignVersionCompare) error {
	return r.db.Create(compare).Error
}

// GetByID 根据ID获取版本对比记录
// Get version compare record by ID
func (r *DesignVersionCompareRepository) GetByID(id uint) (*model.DesignVersionCompare, error) {
	var compare model.DesignVersionCompare
	err := r.db.First(&compare, id).Error
	if err != nil {
		return nil, err
	}
	return &compare, nil
}

// GetByIDWithVersions 根据ID获取版本对比记录（包含版本信息）
// Get version compare record by ID with version info
func (r *DesignVersionCompareRepository) GetByIDWithVersions(id uint) (*model.DesignVersionCompare, error) {
	var compare model.DesignVersionCompare
	err := r.db.Preload("Project").Preload("VersionA").Preload("VersionB").First(&compare, id).Error
	if err != nil {
		return nil, err
	}
	return &compare, nil
}

// GetByVersionPair 根据版本对获取对比记录
// Get compare record by version pair
func (r *DesignVersionCompareRepository) GetByVersionPair(versionAID, versionBID uint) (*model.DesignVersionCompare, error) {
	var compare model.DesignVersionCompare
	err := r.db.Where("(version_a_id = ? AND version_b_id = ?) OR (version_a_id = ? AND version_b_id = ?)",
		versionAID, versionBID, versionBID, versionAID).First(&compare).Error
	if err != nil {
		return nil, err
	}
	return &compare, nil
}

// VersionCompareFilter 版本对比筛选条件
// Version compare filter
type VersionCompareFilter struct {
	ProjectID     uint
	CompareStatus string
}

// List 获取版本对比列表
// Get version compare list
func (r *DesignVersionCompareRepository) List(page, pageSize int, filter VersionCompareFilter) ([]model.DesignVersionCompare, int64, error) {
	var compares []model.DesignVersionCompare
	var total int64

	query := r.db.Model(&model.DesignVersionCompare{})

	// 应用筛选条件 / Apply filters
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.CompareStatus != "" {
		query = query.Where("compare_status = ?", filter.CompareStatus)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Preload("VersionA").Preload("VersionB").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&compares).Error
	if err != nil {
		return nil, 0, err
	}

	return compares, total, nil
}

// GetByProjectID 根据项目ID获取所有对比记录
// Get all compare records by project ID
func (r *DesignVersionCompareRepository) GetByProjectID(projectID uint) ([]model.DesignVersionCompare, error) {
	var compares []model.DesignVersionCompare
	err := r.db.Preload("VersionA").Preload("VersionB").
		Where("project_id = ?", projectID).Order("created_at DESC").Find(&compares).Error
	return compares, err
}

// Update 更新版本对比记录
// Update version compare record
func (r *DesignVersionCompareRepository) Update(compare *model.DesignVersionCompare) error {
	return r.db.Save(compare).Error
}

// UpdateFields 更新指定字段
// Update specific fields
func (r *DesignVersionCompareRepository) UpdateFields(id uint, updates map[string]any) error {
	return r.db.Model(&model.DesignVersionCompare{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除版本对比记录
// Delete version compare record
func (r *DesignVersionCompareRepository) Delete(id uint) error {
	return r.db.Delete(&model.DesignVersionCompare{}, id).Error
}

// DeleteByProjectID 根据项目ID删除所有对比记录
// Delete all compare records by project ID
func (r *DesignVersionCompareRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.DesignVersionCompare{}).Error
}

// DeleteByVersionID 根据版本ID删除相关对比记录
// Delete compare records by version ID
func (r *DesignVersionCompareRepository) DeleteByVersionID(versionID uint) error {
	return r.db.Where("version_a_id = ? OR version_b_id = ?", versionID, versionID).
		Delete(&model.DesignVersionCompare{}).Error
}
