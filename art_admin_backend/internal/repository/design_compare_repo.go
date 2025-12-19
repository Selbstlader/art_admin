package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// DesignCompareRepository 设计比对仓库
// Design compare repository
type DesignCompareRepository struct {
	db *gorm.DB
}

// NewDesignCompareRepository 创建设计比对仓库
// Create design compare repository
func NewDesignCompareRepository(db *gorm.DB) *DesignCompareRepository {
	return &DesignCompareRepository{db: db}
}

// Create 创建设计比对记录
// Create design compare record
func (r *DesignCompareRepository) Create(compare *model.DesignCompareResult) error {
	return r.db.Create(compare).Error
}

// GetByID 根据ID获取设计比对记录
// Get design compare record by ID
func (r *DesignCompareRepository) GetByID(id uint) (*model.DesignCompareResult, error) {
	var compare model.DesignCompareResult
	err := r.db.First(&compare, id).Error
	if err != nil {
		return nil, err
	}
	return &compare, nil
}

// GetByIDWithRelations 根据ID获取设计比对记录（包含关联）
// Get design compare record by ID with relations
func (r *DesignCompareRepository) GetByIDWithRelations(id uint) (*model.DesignCompareResult, error) {
	var compare model.DesignCompareResult
	// 只预加载 Project，Document 关联已移除（改用 DocumentIDs JSON 数组）
	err := r.db.Preload("Project").First(&compare, id).Error
	if err != nil {
		return nil, err
	}
	return &compare, nil
}

// DesignCompareFilter 设计比对筛选条件
// Design compare filter
type DesignCompareFilter struct {
	ProjectID      uint
	DocumentID     uint
	AnalysisStatus string
}

// List 获取设计比对列表
// Get design compare list
func (r *DesignCompareRepository) List(page, pageSize int, filter DesignCompareFilter) ([]model.DesignCompareResult, int64, error) {
	var compares []model.DesignCompareResult
	var total int64

	query := r.db.Model(&model.DesignCompareResult{})

	// 应用筛选条件 / Apply filters
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.DocumentID > 0 {
		query = query.Where("document_id = ?", filter.DocumentID)
	}
	if filter.AnalysisStatus != "" {
		query = query.Where("analysis_status = ?", filter.AnalysisStatus)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&compares).Error
	if err != nil {
		return nil, 0, err
	}

	return compares, total, nil
}

// Update 更新设计比对记录
// Update design compare record
func (r *DesignCompareRepository) Update(compare *model.DesignCompareResult) error {
	return r.db.Save(compare).Error
}

// UpdateAnalysisResult 更新分析结果
// Update analysis result
func (r *DesignCompareRepository) UpdateAnalysisResult(id uint, updates map[string]any) error {
	return r.db.Model(&model.DesignCompareResult{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除设计比对记录
// Delete design compare record
func (r *DesignCompareRepository) Delete(id uint) error {
	return r.db.Delete(&model.DesignCompareResult{}, id).Error
}

// BatchDelete 批量删除设计比对记录
// Batch delete design compare records
func (r *DesignCompareRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.DesignCompareResult{}, ids).Error
}

// GetByProjectID 根据项目ID获取所有比对记录
// Get all compare records by project ID
func (r *DesignCompareRepository) GetByProjectID(projectID uint) ([]model.DesignCompareResult, error) {
	var compares []model.DesignCompareResult
	err := r.db.Where("project_id = ?", projectID).Find(&compares).Error
	return compares, err
}

// DeleteByProjectID 根据项目ID删除所有比对记录
// Delete all compare records by project ID
func (r *DesignCompareRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.DesignCompareResult{}).Error
}
