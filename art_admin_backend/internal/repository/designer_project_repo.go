package repository

import (
	"art_admin_backend/internal/model"
	"time"

	"gorm.io/gorm"
)

// DesignerProjectRepository 设计师项目仓库
// Designer project repository
type DesignerProjectRepository struct {
	db *gorm.DB
}

// NewDesignerProjectRepository 创建设计师项目仓库
// Create designer project repository
func NewDesignerProjectRepository(db *gorm.DB) *DesignerProjectRepository {
	return &DesignerProjectRepository{db: db}
}

// ProjectFilter 项目筛选条件
// Project filter criteria
type ProjectFilter struct {
	Name      string    // 项目名称模糊搜索 / Project name fuzzy search
	Status    string    // 状态筛选 / Status filter
	Style     string    // 风格筛选 / Style filter
	MinBudget float64   // 最小预算 / Minimum budget
	MaxBudget float64   // 最大预算 / Maximum budget
	StartDate time.Time // 开始日期 / Start date
	EndDate   time.Time // 结束日期 / End date
	Keyword   string    // 关键字搜索 / Keyword search
	UserID    uint      // 用户ID / User ID
}

// Create 创建项目
// Create project
func (r *DesignerProjectRepository) Create(project *model.DesignerProject) error {
	return r.db.Create(project).Error
}

// GetByID 根据ID获取项目
// Get project by ID
func (r *DesignerProjectRepository) GetByID(id uint) (*model.DesignerProject, error) {
	var project model.DesignerProject
	err := r.db.First(&project, id).Error
	return &project, err
}

// GetByIDWithAssociations 根据ID获取项目及关联数据
// Get project by ID with associations
func (r *DesignerProjectRepository) GetByIDWithAssociations(id uint) (*model.DesignerProject, error) {
	var project model.DesignerProject
	err := r.db.Preload("Documents").
		Preload("CadFiles").
		Preload("CompareResults").
		Preload("CostEstimate").
		First(&project, id).Error
	return &project, err
}

// Update 更新项目
// Update project
func (r *DesignerProjectRepository) Update(project *model.DesignerProject) error {
	return r.db.Save(project).Error
}

// Delete 删除项目（软删除）
// Delete project (soft delete)
func (r *DesignerProjectRepository) Delete(id uint) error {
	return r.db.Delete(&model.DesignerProject{}, id).Error
}

// DeleteWithAssociations 删除项目及关联数据（级联删除）
// Delete project with associations (cascade delete)
func (r *DesignerProjectRepository) DeleteWithAssociations(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的文档 / Delete associated documents
		if err := tx.Where("project_id = ?", id).Delete(&model.ProjectDocument{}).Error; err != nil {
			return err
		}

		// 删除关联的CAD文件 / Delete associated CAD files
		if err := tx.Where("project_id = ?", id).Delete(&model.CadFile{}).Error; err != nil {
			return err
		}

		// 删除关联的比对结果 / Delete associated compare results
		if err := tx.Where("project_id = ?", id).Delete(&model.DesignCompareResult{}).Error; err != nil {
			return err
		}

		// 删除关联的成本估算 / Delete associated cost estimate
		if err := tx.Where("project_id = ?", id).Delete(&model.CostEstimate{}).Error; err != nil {
			return err
		}

		// 删除项目本身 / Delete the project itself
		if err := tx.Delete(&model.DesignerProject{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// List 获取项目列表
// Get project list
func (r *DesignerProjectRepository) List(page, pageSize int, filter ProjectFilter) ([]model.DesignerProject, int64, error) {
	var projects []model.DesignerProject
	var total int64

	query := r.db.Model(&model.DesignerProject{})

	// 应用筛选条件 / Apply filter criteria
	query = r.applyFilter(query, filter)

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

// Search 搜索项目（支持关键字搜索）
// Search projects (supports keyword search)
func (r *DesignerProjectRepository) Search(keyword string, userID uint, page, pageSize int) ([]model.DesignerProject, int64, error) {
	var projects []model.DesignerProject
	var total int64

	query := r.db.Model(&model.DesignerProject{})

	// 用户ID筛选 / User ID filter
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 关键字搜索（名称或描述） / Keyword search (name or description)
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

// GetByUserID 根据用户ID获取项目列表
// Get projects by user ID
func (r *DesignerProjectRepository) GetByUserID(userID uint) ([]model.DesignerProject, error) {
	var projects []model.DesignerProject
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&projects).Error
	return projects, err
}

// CountByUserID 统计用户项目数量
// Count projects by user ID
func (r *DesignerProjectRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.DesignerProject{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ExistsByID 检查项目是否存在
// Check if project exists by ID
func (r *DesignerProjectRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.DesignerProject{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// CountDocumentsByProjectID 统计项目文档数量
// Count documents by project ID
func (r *DesignerProjectRepository) CountDocumentsByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ProjectDocument{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// CountCadFilesByProjectID 统计项目CAD文件数量
// Count CAD files by project ID
func (r *DesignerProjectRepository) CountCadFilesByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.CadFile{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// applyFilter 应用筛选条件
// Apply filter criteria
func (r *DesignerProjectRepository) applyFilter(query *gorm.DB, filter ProjectFilter) *gorm.DB {
	// 用户ID筛选 / User ID filter
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	// 名称模糊搜索 / Name fuzzy search
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// 状态筛选 / Status filter
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 风格筛选 / Style filter
	if filter.Style != "" {
		query = query.Where("style = ?", filter.Style)
	}

	// 预算范围筛选 / Budget range filter
	if filter.MinBudget > 0 {
		query = query.Where("budget >= ?", filter.MinBudget)
	}
	if filter.MaxBudget > 0 {
		query = query.Where("budget <= ?", filter.MaxBudget)
	}

	// 日期范围筛选 / Date range filter
	if !filter.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// 关键字搜索（名称或描述） / Keyword search (name or description)
	if filter.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	return query
}
