package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"time"
)

// AchievementRepository 成就仓库
type AchievementRepository struct{}

// NewAchievementRepository 创建成就仓库
func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

// FindByID 根据ID查询成就
func (r *AchievementRepository) FindByID(id int64) (*model.Achievement, error) {
	var achievement model.Achievement
	err := database.DB.First(&achievement, id).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// FindByType 根据类型查询成就列表
func (r *AchievementRepository) FindByType(achievementType string) ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := database.DB.Where("type = ? AND is_active = ?", achievementType, true).Find(&achievements).Error
	return achievements, err
}

// FindActive 查询所有激活的成就
func (r *AchievementRepository) FindActive() ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := database.DB.Where("is_active = ?", true).Find(&achievements).Error
	return achievements, err
}

// Create 创建成就
func (r *AchievementRepository) Create(achievement *model.Achievement) error {
	return database.DB.Create(achievement).Error
}

// Update 更新成就
func (r *AchievementRepository) Update(achievement *model.Achievement) error {
	return database.DB.Save(achievement).Error
}

// Delete 删除成就
func (r *AchievementRepository) Delete(id int64) error {
	return database.DB.Delete(&model.Achievement{}, id).Error
}

// UserAchievementRepository 用户成就仓库
type UserAchievementRepository struct{}

// NewUserAchievementRepository 创建用户成就仓库
func NewUserAchievementRepository() *UserAchievementRepository {
	return &UserAchievementRepository{}
}

// FindByUserID 根据用户ID查找用户成就
func (r *UserAchievementRepository) FindByUserID(userID int64) ([]model.UserAchievement, error) {
	var userAchievements []model.UserAchievement
	err := database.DB.Preload("Achievement").Where("user_id = ?", userID).Find(&userAchievements).Error
	return userAchievements, err
}

// FindByUserIDWithPagination 根据用户ID分页查找用户成就
func (r *UserAchievementRepository) FindByUserIDWithPagination(userID int64, page, pageSize int) ([]model.UserAchievement, int64, error) {
	var userAchievements []model.UserAchievement
	var total int64

	// 获取总数
	if err := database.DB.Model(&model.UserAchievement{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := database.DB.Preload("Achievement").Where("user_id = ?", userID).
		Offset(offset).Limit(pageSize).Find(&userAchievements).Error

	return userAchievements, total, err
}

// FindByUserAndAchievement 根据用户ID和成就ID查询用户成就
func (r *UserAchievementRepository) FindByUserAndAchievement(userID, achievementID int64) (*model.UserAchievement, error) {
	var userAchievement model.UserAchievement
	err := database.DB.Preload("Achievement").Where("user_id = ? AND achievement_id = ?", userID, achievementID).First(&userAchievement).Error
	if err != nil {
		return nil, err
	}
	return &userAchievement, nil
}

// Create 创建用户成就记录
func (r *UserAchievementRepository) Create(userAchievement *model.UserAchievement) error {
	return database.DB.Create(userAchievement).Error
}

// Update 更新用户成就记录
func (r *UserAchievementRepository) Update(userAchievement *model.UserAchievement) error {
	return database.DB.Save(userAchievement).Error
}

// UnlockAchievement 解锁成就
func (r *UserAchievementRepository) UnlockAchievement(userID, achievementID int64) error {
	now := time.Now()
	return database.DB.Model(&model.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Updates(map[string]interface{}{
			"is_unlocked": true,
			"unlocked_at": &now,
			"updated_at":  now,
		}).Error
}

// GetUnlockedCount 获取用户已解锁成就数量
func (r *UserAchievementRepository) GetUnlockedCount(userID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.UserAchievement{}).
		Where("user_id = ? AND is_unlocked = ?", userID, true).
		Count(&count).Error
	return count, err
}

// GetUnlockedByTier 获取用户按等级分组的已解锁成就数量
func (r *UserAchievementRepository) GetUnlockedByTier(userID int64) (map[string]int64, error) {
	var results []struct {
		Tier  string `json:"tier"`
		Count int64  `json:"count"`
	}

	err := database.DB.Model(&model.UserAchievement{}).
		Select("a.tier, COUNT(*) as count").
		Joins("LEFT JOIN achievements a ON user_achievements.achievement_id = a.id").
		Where("user_achievements.user_id = ? AND user_achievements.is_unlocked = ?", userID, true).
		Group("a.tier").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	tierCount := make(map[string]int64)
	for _, result := range results {
		tierCount[result.Tier] = result.Count
	}

	return tierCount, nil
}

// UserPointRepository 用户积分仓库
type UserPointRepository struct{}

// NewUserPointRepository 创建用户积分仓库
func NewUserPointRepository() *UserPointRepository {
	return &UserPointRepository{}
}

// FindByUserID 根据用户ID查询用户积分
func (r *UserPointRepository) FindByUserID(userID int64) (*model.UserPoint, error) {
	var userPoint model.UserPoint
	err := database.DB.Where("user_id = ?", userID).First(&userPoint).Error
	if err != nil {
		return nil, err
	}
	return &userPoint, nil
}

// Create 创建用户积分记录
func (r *UserPointRepository) Create(userPoint *model.UserPoint) error {
	return database.DB.Create(userPoint).Error
}

// Update 更新用户积分记录
func (r *UserPointRepository) Update(userPoint *model.UserPoint) error {
	return database.DB.Save(userPoint).Error
}

// AddPoints 增加用户积分
func (r *UserPointRepository) AddPoints(userID int64, points int, category string) error {
	now := time.Now()

	// 构建更新字段
	updates := map[string]interface{}{
		"total_points": database.DB.Raw("total_points + ?", points),
		"last_updated": now,
		"updated_at":   now,
	}

	// 根据分类更新对应字段
	switch category {
	case "mood":
		updates["mood_points"] = database.DB.Raw("mood_points + ?", points)
	case "meditation":
		updates["meditation_points"] = database.DB.Raw("meditation_points + ?", points)
	case "journal":
		updates["journal_points"] = database.DB.Raw("journal_points + ?", points)
	case "analysis":
		updates["analysis_points"] = database.DB.Raw("analysis_points + ?", points)
	case "bonus":
		updates["bonus_points"] = database.DB.Raw("bonus_points + ?", points)
	}

	return database.DB.Model(&model.UserPoint{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

// GetTopUsers 获取积分排行榜
func (r *UserPointRepository) GetTopUsers(limit int) ([]model.UserPoint, error) {
	var userPoints []model.UserPoint
	err := database.DB.Order("total_points DESC, COALESCE(updated_at, created_at) DESC").
		Limit(limit).
		Find(&userPoints).Error
	return userPoints, err
}

// AchievementEventRepository 成就事件仓库
type AchievementEventRepository struct{}

// NewAchievementEventRepository 创建成就事件仓库
func NewAchievementEventRepository() *AchievementEventRepository {
	return &AchievementEventRepository{}
}

// Create 创建成就事件
func (r *AchievementEventRepository) Create(event *model.AchievementEvent) error {
	return database.DB.Create(event).Error
}

// FindUnprocessed 查找未处理的事件
func (r *AchievementEventRepository) FindUnprocessed(limit int) ([]model.AchievementEvent, error) {
	var events []model.AchievementEvent
	err := database.DB.Where("processed_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// MarkAsProcessed 标记事件为已处理
func (r *AchievementEventRepository) MarkAsProcessed(eventID int64) error {
	now := time.Now()
	return database.DB.Model(&model.AchievementEvent{}).
		Where("id = ?", eventID).
		Update("processed_at", &now).Error
}

// LeaderboardRepository 排行榜仓库
type LeaderboardRepository struct{}

// NewLeaderboardRepository 创建排行榜仓库
func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{}
}

// FindByTypeAndPeriod 根据类型和周期查询排行榜
func (r *LeaderboardRepository) FindByTypeAndPeriod(leaderboardType, period string) (*model.Leaderboard, error) {
	var leaderboard model.Leaderboard
	err := database.DB.Where("type = ? AND period = ?", leaderboardType, period).First(&leaderboard).Error
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// Update 更新排行榜
func (r *LeaderboardRepository) Update(leaderboard *model.Leaderboard) error {
	return database.DB.Save(leaderboard).Error
}

// Upsert 创建或更新排行榜
func (r *LeaderboardRepository) Upsert(leaderboard *model.Leaderboard) error {
	var existing model.Leaderboard
	err := database.DB.Where("type = ? AND period = ?", leaderboard.Type, leaderboard.Period).First(&existing).Error

	if err != nil {
		// 不存在则创建
		return database.DB.Create(leaderboard).Error
	}

	// 存在则更新
	leaderboard.ID = existing.ID
	return database.DB.Save(leaderboard).Error
}
