package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// MeditationContentRepository 冥想内容仓库
type MeditationContentRepository struct {
	db *gorm.DB
}

// NewMeditationContentRepository 创建冥想内容仓库
func NewMeditationContentRepository() *MeditationContentRepository {
	return &MeditationContentRepository{
		db: database.GetDB(),
	}
}

// Create 创建冥想内容
func (r *MeditationContentRepository) Create(content *model.MeditationContent) error {
	if r.db == nil {
		fmt.Printf("MeditationContentRepository: database connection is nil\n")
		return errors.New("数据库连接未初始化")
	}
	return r.db.Create(content).Error
}

// FindByID 根据ID查找冥想内容
func (r *MeditationContentRepository) FindByID(id int64) (*model.MeditationContent, error) {
	var content model.MeditationContent
	err := r.db.First(&content, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("冥想内容不存在")
		}
		return nil, err
	}
	return &content, nil
}

// Update 更新冥想内容
func (r *MeditationContentRepository) Update(content *model.MeditationContent) error {
	return r.db.Save(content).Error
}

// Delete 删除冥想内容
func (r *MeditationContentRepository) Delete(id int64) error {
	return r.db.Delete(&model.MeditationContent{}, id).Error
}

// FindWithPagination 分页查询冥想内容
func (r *MeditationContentRepository) FindWithPagination(query map[string]interface{}, page, pageSize int) ([]model.MeditationContent, int64, error) {
	var contents []model.MeditationContent
	var total int64

	// 构建查询
	db := r.db.Model(&model.MeditationContent{})

	// 添加查询条件
	if category, ok := query["category"].(string); ok && category != "" {
		db = db.Where("category = ?", category)
	}
	if difficulty, ok := query["difficulty_level"].(int); ok && difficulty > 0 {
		db = db.Where("difficulty_level = ?", difficulty)
	}
	if isPublic, ok := query["is_public"].(bool); ok {
		db = db.Where("is_public = ?", isPublic)
	}
	if minDuration, ok := query["min_duration"].(int); ok && minDuration > 0 {
		db = db.Where("duration >= ?", minDuration)
	}
	if maxDuration, ok := query["max_duration"].(int); ok && maxDuration > 0 {
		db = db.Where("duration <= ?", maxDuration)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

// FindByCategory 根据分类查找冥想内容
func (r *MeditationContentRepository) FindByCategory(category string, limit int) ([]model.MeditationContent, error) {
	var contents []model.MeditationContent
	query := r.db.Where("category = ? AND is_public = ?", category, true)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("view_count DESC, created_at DESC").Find(&contents).Error
	return contents, err
}

// FindByDifficulty 根据难度等级查找冥想内容
func (r *MeditationContentRepository) FindByDifficulty(difficulty int, limit int) ([]model.MeditationContent, error) {
	var contents []model.MeditationContent
	query := r.db.Where("difficulty_level = ? AND is_public = ?", difficulty, true)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("view_count DESC, created_at DESC").Find(&contents).Error
	return contents, err
}

// FindPopular 查找热门冥想内容
func (r *MeditationContentRepository) FindPopular(limit int) ([]model.MeditationContent, error) {
	var contents []model.MeditationContent
	query := r.db.Where("is_public = ?", true)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("view_count DESC, like_count DESC, created_at DESC").Find(&contents).Error
	return contents, err
}

// IncrementViewCount 增加查看次数
func (r *MeditationContentRepository) IncrementViewCount(id int64) error {
	return r.db.Model(&model.MeditationContent{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *MeditationContentRepository) IncrementLikeCount(id int64) error {
	return r.db.Model(&model.MeditationContent{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (r *MeditationContentRepository) DecrementLikeCount(id int64) error {
	return r.db.Model(&model.MeditationContent{}).Where("id = ? AND like_count > 0", id).Update("like_count", gorm.Expr("like_count - 1")).Error
}

// GetCategories 获取所有分类
func (r *MeditationContentRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&model.MeditationContent{}).Where("is_public = ?", true).Distinct("category").Pluck("category", &categories).Error
	return categories, err
}

// GetDifficultyLevels 获取所有难度等级
func (r *MeditationContentRepository) GetDifficultyLevels() ([]int, error) {
	var levels []int
	err := r.db.Model(&model.MeditationContent{}).Where("is_public = ?", true).Distinct("difficulty_level").Pluck("difficulty_level", &levels).Error
	return levels, err
}

// Search 搜索冥想内容
func (r *MeditationContentRepository) Search(keyword string, page, pageSize int) ([]model.MeditationContent, int64, error) {
	var contents []model.MeditationContent
	var total int64

	// 构建搜索查询
	db := r.db.Model(&model.MeditationContent{}).Where("is_public = ? AND (title LIKE ? OR description LIKE ?)", true, "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("view_count DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}
