package repository

import (
	"context"

	"art_admin_backend/internal/model"

	"gorm.io/gorm"
)

// LearningMaterialRepository 教材仓储
type LearningMaterialRepository struct {
	db *gorm.DB
}

// NewLearningMaterialRepository 创建教材仓储
func NewLearningMaterialRepository(db *gorm.DB) *LearningMaterialRepository {
	return &LearningMaterialRepository{db: db}
}

// Create 创建教材
func (r *LearningMaterialRepository) Create(ctx context.Context, material *model.LearningMaterial) error {
	return r.db.WithContext(ctx).Create(material).Error
}

// GetByID 根据 ID 获取教材
func (r *LearningMaterialRepository) GetByID(ctx context.Context, id int64) (*model.LearningMaterial, error) {
	var material model.LearningMaterial
	if err := r.db.WithContext(ctx).Where("id = ? AND status = 1", id).First(&material).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// List 获取教材列表
func (r *LearningMaterialRepository) List(ctx context.Context, userID, subjectID int64, grade string, offset, limit int) ([]*model.LearningMaterial, int64, error) {
	var materials []*model.LearningMaterial
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LearningMaterial{}).Where("status = 1")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if subjectID > 0 {
		query = query.Where("subject_id = ?", subjectID)
	}

	if grade != "" {
		query = query.Where("grade = ?", grade)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

// Update 更新教材
func (r *LearningMaterialRepository) Update(ctx context.Context, material *model.LearningMaterial) error {
	return r.db.WithContext(ctx).Save(material).Error
}

// Delete 删除教材（软删除）
func (r *LearningMaterialRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.LearningMaterial{}).
		Where("id = ?", id).
		Update("status", 0).Error
}

// IncrementViewCount 增加浏览次数
func (r *LearningMaterialRepository) IncrementViewCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.LearningMaterial{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementFavoriteCount 增加收藏次数
func (r *LearningMaterialRepository) IncrementFavoriteCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.LearningMaterial{}).
		Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

// DecrementFavoriteCount 减少收藏次数
func (r *LearningMaterialRepository) DecrementFavoriteCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.LearningMaterial{}).
		Where("id = ? AND favorite_count > 0", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
}

// GetGradesBySubject 获取某学科下的年级列表及课程数量
func (r *LearningMaterialRepository) GetGradesBySubject(ctx context.Context, userID, subjectID int64) (map[string]int64, error) {
	type GradeCount struct {
		Grade string
		Count int64
	}

	var results []GradeCount
	query := r.db.WithContext(ctx).Model(&model.LearningMaterial{}).
		Select("grade, COUNT(*) as count").
		Where("status = 1")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if subjectID > 0 {
		query = query.Where("subject_id = ?", subjectID)
	}

	if err := query.Group("grade").Find(&results).Error; err != nil {
		return nil, err
	}

	gradeMap := make(map[string]int64)
	for _, r := range results {
		gradeMap[r.Grade] = r.Count
	}

	return gradeMap, nil
}
