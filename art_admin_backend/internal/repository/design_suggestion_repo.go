package repository

import (
	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// DesignSuggestionRepository 设计建议仓库
// Design suggestion repository
type DesignSuggestionRepository struct {
	db *gorm.DB
}

// NewDesignSuggestionRepository 创建设计建议仓库
// Create design suggestion repository
func NewDesignSuggestionRepository(db *gorm.DB) *DesignSuggestionRepository {
	return &DesignSuggestionRepository{db: db}
}

// SuggestionFilter 建议筛选条件
// Suggestion filter criteria
type SuggestionFilter struct {
	ProjectID uint
	Category  string
	Status    string
	UserID    uint
}

// Create 创建设计建议
// Create design suggestion
func (r *DesignSuggestionRepository) Create(suggestion *model.DesignSuggestion) error {
	return r.db.Create(suggestion).Error
}

// BatchCreate 批量创建设计建议
// Batch create design suggestions
func (r *DesignSuggestionRepository) BatchCreate(suggestions []model.DesignSuggestion) error {
	return r.db.Create(&suggestions).Error
}

// GetByID 根据ID获取设计建议
// Get design suggestion by ID
func (r *DesignSuggestionRepository) GetByID(id uint) (*model.DesignSuggestion, error) {
	var suggestion model.DesignSuggestion
	err := r.db.First(&suggestion, id).Error
	if err != nil {
		return nil, err
	}
	return &suggestion, nil
}

// GetByIDWithProject 根据ID获取设计建议（包含项目信息）
// Get design suggestion by ID with project info
func (r *DesignSuggestionRepository) GetByIDWithProject(id uint) (*model.DesignSuggestion, error) {
	var suggestion model.DesignSuggestion
	err := r.db.Preload("Project").First(&suggestion, id).Error
	if err != nil {
		return nil, err
	}
	return &suggestion, nil
}

// Update 更新设计建议
// Update design suggestion
func (r *DesignSuggestionRepository) Update(suggestion *model.DesignSuggestion) error {
	return r.db.Save(suggestion).Error
}

// UpdateStatus 更新建议状态
// Update suggestion status
func (r *DesignSuggestionRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.DesignSuggestion{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除设计建议
// Delete design suggestion
func (r *DesignSuggestionRepository) Delete(id uint) error {
	return r.db.Delete(&model.DesignSuggestion{}, id).Error
}

// BatchDelete 批量删除设计建议
// Batch delete design suggestions
func (r *DesignSuggestionRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.DesignSuggestion{}, ids).Error
}

// List 获取设计建议列表
// Get design suggestion list
func (r *DesignSuggestionRepository) List(page, pageSize int, filter SuggestionFilter) ([]model.DesignSuggestion, int64, error) {
	var suggestions []model.DesignSuggestion
	var total int64

	query := r.db.Model(&model.DesignSuggestion{})

	// 应用筛选条件 / Apply filter criteria
	if filter.ProjectID > 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	// 获取总数 / Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 / Paginated query
	offset := (page - 1) * pageSize
	err := query.Order("priority DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&suggestions).Error

	return suggestions, total, err
}

// GetByProjectID 根据项目ID获取所有建议
// Get all suggestions by project ID
func (r *DesignSuggestionRepository) GetByProjectID(projectID uint) ([]model.DesignSuggestion, error) {
	var suggestions []model.DesignSuggestion
	err := r.db.Where("project_id = ?", projectID).
		Order("priority DESC, created_at DESC").
		Find(&suggestions).Error
	return suggestions, err
}

// CountByProjectID 统计项目建议数量
// Count suggestions by project ID
func (r *DesignSuggestionRepository) CountByProjectID(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.DesignSuggestion{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// GetStats 获取建议统计
// Get suggestion statistics
func (r *DesignSuggestionRepository) GetStats(projectID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总数 / Total count
	var total int64
	query := r.db.Model(&model.DesignSuggestion{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	query.Count(&total)
	stats["total"] = total

	// 各状态统计 / Status statistics
	statuses := []string{"pending", "adopted", "ignored"}
	for _, status := range statuses {
		var count int64
		q := r.db.Model(&model.DesignSuggestion{}).Where("status = ?", status)
		if projectID > 0 {
			q = q.Where("project_id = ?", projectID)
		}
		q.Count(&count)
		stats[status] = count
	}

	return stats, nil
}

// GetCategoryStats 获取分类统计
// Get category statistics
func (r *DesignSuggestionRepository) GetCategoryStats(projectID uint) (map[string]int64, error) {
	type Result struct {
		Category string
		Count    int64
	}
	var results []Result

	query := r.db.Model(&model.DesignSuggestion{}).
		Select("category, count(*) as count").
		Group("category")
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Category] = r.Count
	}
	return stats, nil
}

// CreatePreference 创建偏好记录
// Create preference record
func (r *DesignSuggestionRepository) CreatePreference(pref *model.SuggestionPreference) error {
	return r.db.Create(pref).Error
}

// GetPreferencesBySuggestionID 获取建议的偏好记录
// Get preferences by suggestion ID
func (r *DesignSuggestionRepository) GetPreferencesBySuggestionID(suggestionID uint) ([]model.SuggestionPreference, error) {
	var prefs []model.SuggestionPreference
	err := r.db.Where("suggestion_id = ?", suggestionID).
		Order("created_at DESC").
		Find(&prefs).Error
	return prefs, err
}

// DeleteByProjectID 根据项目ID删除所有建议
// Delete all suggestions by project ID
func (r *DesignSuggestionRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.DesignSuggestion{}).Error
}
