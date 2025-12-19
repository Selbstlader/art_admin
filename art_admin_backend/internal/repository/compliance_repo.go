package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// ComplianceRepository 合规检查仓库
// Compliance repository for design standards and check results
type ComplianceRepository struct {
	db *gorm.DB
}

// NewComplianceRepository 创建合规检查仓库
// Create compliance repository
func NewComplianceRepository(db *gorm.DB) *ComplianceRepository {
	return &ComplianceRepository{db: db}
}

// ========== 设计规范相关方法 / Design Standard Methods ==========

// CreateStandard 创建设计规范
// Create design standard
func (r *ComplianceRepository) CreateStandard(standard *model.DesignStandard) error {
	return r.db.Create(standard).Error
}

// GetStandardByID 根据ID获取设计规范
// Get design standard by ID
func (r *ComplianceRepository) GetStandardByID(id uint) (*model.DesignStandard, error) {
	var standard model.DesignStandard
	err := r.db.First(&standard, id).Error
	if err != nil {
		return nil, err
	}
	return &standard, nil
}

// GetStandardByCode 根据编号获取设计规范
// Get design standard by code
func (r *ComplianceRepository) GetStandardByCode(code string) (*model.DesignStandard, error) {
	var standard model.DesignStandard
	err := r.db.Where("code = ?", code).First(&standard).Error
	if err != nil {
		return nil, err
	}
	return &standard, nil
}

// StandardFilter 设计规范筛选条件
// Design standard filter
type StandardFilter struct {
	Category string
	Status   string
	Keyword  string
}

// ListStandards 获取设计规范列表
// Get design standard list
func (r *ComplianceRepository) ListStandards(page, pageSize int, filter StandardFilter) ([]model.DesignStandard, int64, error) {
	var standards []model.DesignStandard
	var total int64

	query := r.db.Model(&model.DesignStandard{})

	// 应用筛选条件 / Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR content LIKE ?", keyword, keyword, keyword)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Order("category ASC, code ASC").Offset(offset).Limit(pageSize).Find(&standards).Error
	if err != nil {
		return nil, 0, err
	}

	return standards, total, nil
}

// GetStandardsByCategory 根据类别获取设计规范
// Get design standards by category
func (r *ComplianceRepository) GetStandardsByCategory(category string) ([]model.DesignStandard, error) {
	var standards []model.DesignStandard
	err := r.db.Where("category = ? AND status = ?", category, model.StandardStatusActive).
		Order("code ASC").Find(&standards).Error
	return standards, err
}

// GetStandardsByCategories 根据多个类别获取设计规范
// Get design standards by multiple categories
func (r *ComplianceRepository) GetStandardsByCategories(categories []string) ([]model.DesignStandard, error) {
	var standards []model.DesignStandard
	err := r.db.Where("category IN ? AND status = ?", categories, model.StandardStatusActive).
		Order("category ASC, code ASC").Find(&standards).Error
	return standards, err
}

// GetAllActiveStandards 获取所有有效的设计规范
// Get all active design standards
func (r *ComplianceRepository) GetAllActiveStandards() ([]model.DesignStandard, error) {
	var standards []model.DesignStandard
	err := r.db.Where("status = ?", model.StandardStatusActive).
		Order("category ASC, code ASC").Find(&standards).Error
	return standards, err
}

// UpdateStandard 更新设计规范
// Update design standard
func (r *ComplianceRepository) UpdateStandard(standard *model.DesignStandard) error {
	return r.db.Save(standard).Error
}

// DeleteStandard 删除设计规范
// Delete design standard
func (r *ComplianceRepository) DeleteStandard(id uint) error {
	return r.db.Delete(&model.DesignStandard{}, id).Error
}

// BatchDeleteStandards 批量删除设计规范
// Batch delete design standards
func (r *ComplianceRepository) BatchDeleteStandards(ids []uint) error {
	return r.db.Delete(&model.DesignStandard{}, ids).Error
}

// GetStandardCategoryStats 获取规范类别统计
// Get standard category statistics
func (r *ComplianceRepository) GetStandardCategoryStats() ([]struct {
	Category string
	Count    int64
}, error) {
	var stats []struct {
		Category string
		Count    int64
	}
	err := r.db.Model(&model.DesignStandard{}).
		Select("category, COUNT(*) as count").
		Where("status = ?", model.StandardStatusActive).
		Group("category").
		Scan(&stats).Error
	return stats, err
}

// GetStandardStats 获取规范统计信息
// Get standard statistics
func (r *ComplianceRepository) GetStandardStats() (total, active int64, err error) {
	err = r.db.Model(&model.DesignStandard{}).Count(&total).Error
	if err != nil {
		return
	}
	err = r.db.Model(&model.DesignStandard{}).Where("status = ?", model.StandardStatusActive).Count(&active).Error
	return
}

// ========== 合规检查结果相关方法 / Compliance Check Result Methods ==========

// CreateCheckResult 创建合规检查结果
// Create compliance check result
func (r *ComplianceRepository) CreateCheckResult(result *model.ComplianceCheckResult) error {
	return r.db.Create(result).Error
}

// GetCheckResultByID 根据ID获取合规检查结果
// Get compliance check result by ID
func (r *ComplianceRepository) GetCheckResultByID(id uint) (*model.ComplianceCheckResult, error) {
	var result model.ComplianceCheckResult
	err := r.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCheckResultByIDWithProject 根据ID获取合规检查结果（包含项目信息）
// Get compliance check result by ID with project
func (r *ComplianceRepository) GetCheckResultByIDWithProject(id uint) (*model.ComplianceCheckResult, error) {
	var result model.ComplianceCheckResult
	err := r.db.Preload("Project").First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckResultFilter 合规检查结果筛选条件
// Compliance check result filter
type CheckResultFilter struct {
	ProjectID   uint
	CheckStatus string
	UserID      uint
}

// ListCheckResults 获取合规检查结果列表
// Get compliance check result list
func (r *ComplianceRepository) ListCheckResults(page, pageSize int, filter CheckResultFilter) ([]model.ComplianceCheckResult, int64, error) {
	var results []model.ComplianceCheckResult
	var total int64

	query := r.db.Model(&model.ComplianceCheckResult{})

	// 应用筛选条件 / Apply filters
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.CheckStatus != "" {
		query = query.Where("check_status = ?", filter.CheckStatus)
	}
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	// 统计总数 / Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Preload("Project").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetLatestCheckResultByProjectID 获取项目最新的合规检查结果
// Get latest compliance check result by project ID
func (r *ComplianceRepository) GetLatestCheckResultByProjectID(projectID uint) (*model.ComplianceCheckResult, error) {
	var result model.ComplianceCheckResult
	err := r.db.Where("project_id = ? AND check_status = ?", projectID, model.CheckStatusCompleted).
		Order("created_at DESC").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCheckResult 更新合规检查结果
// Update compliance check result
func (r *ComplianceRepository) UpdateCheckResult(result *model.ComplianceCheckResult) error {
	return r.db.Save(result).Error
}

// UpdateCheckResultStatus 更新检查结果状态
// Update check result status
func (r *ComplianceRepository) UpdateCheckResultStatus(id uint, updates map[string]any) error {
	return r.db.Model(&model.ComplianceCheckResult{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCheckResult 删除合规检查结果
// Delete compliance check result
func (r *ComplianceRepository) DeleteCheckResult(id uint) error {
	return r.db.Delete(&model.ComplianceCheckResult{}, id).Error
}

// BatchDeleteCheckResults 批量删除合规检查结果
// Batch delete compliance check results
func (r *ComplianceRepository) BatchDeleteCheckResults(ids []uint) error {
	return r.db.Delete(&model.ComplianceCheckResult{}, ids).Error
}

// DeleteCheckResultsByProjectID 根据项目ID删除所有检查结果
// Delete all check results by project ID
func (r *ComplianceRepository) DeleteCheckResultsByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.ComplianceCheckResult{}).Error
}

// GetCheckResultCountByProjectID 获取项目的检查结果数量
// Get check result count by project ID
func (r *ComplianceRepository) GetCheckResultCountByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ComplianceCheckResult{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}
