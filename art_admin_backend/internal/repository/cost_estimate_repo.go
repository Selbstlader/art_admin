package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// CostEstimateRepository 成本估算仓库接口
// Cost estimate repository interface
type CostEstimateRepository interface {
	Create(estimate *model.CostEstimate) error
	Update(estimate *model.CostEstimate) error
	Delete(id uint) error
	FindByID(id uint) (*model.CostEstimate, error)
	FindByProjectID(projectID uint) (*model.CostEstimate, error)
	FindAll(params CostEstimateQueryParams) ([]model.CostEstimate, int64, error)
	GetSummary(userID uint) (*CostSummary, error)
}

// CostEstimateQueryParams 成本估算查询参数
// Cost estimate query parameters
type CostEstimateQueryParams struct {
	Current   int
	Size      int
	UserID    uint
	ProjectID uint
	MinCost   float64
	MaxCost   float64
}

// CostSummary 成本汇总统计
// Cost summary statistics
type CostSummary struct {
	TotalProjects     int64
	TotalBudget       float64
	TotalCost         float64
	AvgCostPerProject float64
	OverBudgetCount   int64
}

// costEstimateRepository 成本估算仓库实现
// Cost estimate repository implementation
type costEstimateRepository struct {
	db *gorm.DB
}

// NewCostEstimateRepository 创建成本估算仓库实例
// Create cost estimate repository instance
func NewCostEstimateRepository() CostEstimateRepository {
	return &costEstimateRepository{
		db: database.DB,
	}
}

// Create 创建成本估算
// Create cost estimate
func (r *costEstimateRepository) Create(estimate *model.CostEstimate) error {
	return r.db.Create(estimate).Error
}

// Update 更新成本估算
// Update cost estimate
func (r *costEstimateRepository) Update(estimate *model.CostEstimate) error {
	return r.db.Save(estimate).Error
}

// Delete 删除成本估算
// Delete cost estimate
func (r *costEstimateRepository) Delete(id uint) error {
	return r.db.Delete(&model.CostEstimate{}, id).Error
}

// FindByID 根据ID查找成本估算
// Find cost estimate by ID
func (r *costEstimateRepository) FindByID(id uint) (*model.CostEstimate, error) {
	var estimate model.CostEstimate
	err := r.db.Preload("Project").First(&estimate, id).Error
	if err != nil {
		return nil, err
	}
	return &estimate, nil
}

// FindByProjectID 根据项目ID查找成本估算
// Find cost estimate by project ID
func (r *costEstimateRepository) FindByProjectID(projectID uint) (*model.CostEstimate, error) {
	var estimate model.CostEstimate
	err := r.db.Preload("Project").Where("project_id = ?", projectID).First(&estimate).Error
	if err != nil {
		return nil, err
	}
	return &estimate, nil
}

// FindAll 查询成本估算列表
// Find all cost estimates with pagination
func (r *costEstimateRepository) FindAll(params CostEstimateQueryParams) ([]model.CostEstimate, int64, error) {
	var estimates []model.CostEstimate
	var total int64

	query := r.db.Model(&model.CostEstimate{}).Preload("Project")

	// 按用户筛选(通过项目关联) / Filter by user (through project association)
	if params.UserID > 0 {
		query = query.Joins("JOIN designer_projects ON designer_projects.id = cost_estimates.project_id").
			Where("designer_projects.user_id = ?", params.UserID)
	}

	// 按项目筛选 / Filter by project
	if params.ProjectID > 0 {
		query = query.Where("project_id = ?", params.ProjectID)
	}

	// 成本区间筛选 / Cost range filter
	if params.MinCost > 0 {
		query = query.Where("total_cost >= ?", params.MinCost)
	}
	if params.MaxCost > 0 {
		query = query.Where("total_cost <= ?", params.MaxCost)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Pagination query
	offset := (params.Current - 1) * params.Size
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Size).Find(&estimates).Error; err != nil {
		return nil, 0, err
	}

	return estimates, total, nil
}

// GetSummary 获取成本汇总统计
// Get cost summary statistics
func (r *costEstimateRepository) GetSummary(userID uint) (*CostSummary, error) {
	var summary CostSummary

	query := r.db.Model(&model.CostEstimate{})
	if userID > 0 {
		query = query.Joins("JOIN designer_projects ON designer_projects.id = cost_estimates.project_id").
			Where("designer_projects.user_id = ?", userID)
	}

	// 获取项目总数和成本统计 / Get total projects and cost statistics
	type Stats struct {
		TotalProjects int64
		TotalBudget   float64
		TotalCost     float64
	}
	var stats Stats
	if err := query.Select("COUNT(*) as total_projects, COALESCE(SUM(budget_limit), 0) as total_budget, COALESCE(SUM(total_cost), 0) as total_cost").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	summary.TotalProjects = stats.TotalProjects
	summary.TotalBudget = stats.TotalBudget
	summary.TotalCost = stats.TotalCost

	// 计算平均成本 / Calculate average cost
	if summary.TotalProjects > 0 {
		summary.AvgCostPerProject = summary.TotalCost / float64(summary.TotalProjects)
	}

	// 获取超预算项目数 / Get over budget project count
	overBudgetQuery := r.db.Model(&model.CostEstimate{}).
		Where("budget_limit > 0 AND total_cost > budget_limit")
	if userID > 0 {
		overBudgetQuery = overBudgetQuery.Joins("JOIN designer_projects ON designer_projects.id = cost_estimates.project_id").
			Where("designer_projects.user_id = ?", userID)
	}
	if err := overBudgetQuery.Count(&summary.OverBudgetCount).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}
