package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"errors"
	"time"

	"gorm.io/gorm"
)

// MeditationFavoriteRepository 冥想收藏仓库
type MeditationFavoriteRepository struct{}

// NewMeditationFavoriteRepository 创建冥想收藏仓库
func NewMeditationFavoriteRepository() *MeditationFavoriteRepository {
	return &MeditationFavoriteRepository{}
}

// AddFavorite 添加收藏
func (r *MeditationFavoriteRepository) AddFavorite(userID, contentID int64) error {
	db := database.GetDB()

	// 检查是否已收藏
	var existing model.MeditationFavorite
	err := db.Where("user_id = ? AND content_id = ?", userID, contentID).First(&existing).Error
	if err == nil {
		return errors.New("已收藏该内容")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建收藏
	favorite := &model.MeditationFavorite{
		UserID:    userID,
		ContentID: contentID,
		CreatedAt: time.Now(),
	}
	return db.Create(favorite).Error
}

// RemoveFavorite 取消收藏
func (r *MeditationFavoriteRepository) RemoveFavorite(userID, contentID int64) error {
	db := database.GetDB()
	result := db.Where("user_id = ? AND content_id = ?", userID, contentID).Delete(&model.MeditationFavorite{})
	if result.RowsAffected == 0 {
		return errors.New("未收藏该内容")
	}
	return result.Error
}

// IsFavorite 检查是否已收藏
func (r *MeditationFavoriteRepository) IsFavorite(userID, contentID int64) (bool, error) {
	db := database.GetDB()
	var count int64
	err := db.Model(&model.MeditationFavorite{}).Where("user_id = ? AND content_id = ?", userID, contentID).Count(&count).Error
	return count > 0, err
}

// GetUserFavorites 获取用户收藏列表
func (r *MeditationFavoriteRepository) GetUserFavorites(userID int64, page, pageSize int) ([]model.MeditationFavorite, int64, error) {
	db := database.GetDB()
	var favorites []model.MeditationFavorite
	var total int64

	// 获取总数
	if err := db.Model(&model.MeditationFavorite{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载内容
	offset := (page - 1) * pageSize
	if err := db.Where("user_id = ?", userID).
		Preload("Content").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// GetFavoriteCount 获取用户收藏数量
func (r *MeditationFavoriteRepository) GetFavoriteCount(userID int64) (int64, error) {
	db := database.GetDB()
	var count int64
	err := db.Model(&model.MeditationFavorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// MeditationPlayRecordRepository 冥想播放记录仓库
type MeditationPlayRecordRepository struct{}

// NewMeditationPlayRecordRepository 创建冥想播放记录仓库
func NewMeditationPlayRecordRepository() *MeditationPlayRecordRepository {
	return &MeditationPlayRecordRepository{}
}

// CreateOrUpdatePlayRecord 创建或更新播放记录
func (r *MeditationPlayRecordRepository) CreateOrUpdatePlayRecord(userID, contentID int64, duration int, progress float64, completed bool) error {
	db := database.GetDB()

	var existing model.MeditationPlayRecord
	err := db.Where("user_id = ? AND content_id = ?", userID, contentID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		record := &model.MeditationPlayRecord{
			UserID:       userID,
			ContentID:    contentID,
			Duration:     duration,
			Progress:     progress,
			Completed:    completed,
			LastPlayedAt: time.Now(),
			PlayCount:    1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		return db.Create(record).Error
	}

	if err != nil {
		return err
	}

	// 更新现有记录
	updates := map[string]interface{}{
		"duration":       duration,
		"progress":       progress,
		"completed":      completed || existing.Completed, // 一旦完成就保持完成状态
		"last_played_at": time.Now(),
		"play_count":     existing.PlayCount + 1,
		"updated_at":     time.Now(),
	}
	return db.Model(&existing).Updates(updates).Error
}

// GetPlayRecord 获取播放记录
func (r *MeditationPlayRecordRepository) GetPlayRecord(userID, contentID int64) (*model.MeditationPlayRecord, error) {
	db := database.GetDB()
	var record model.MeditationPlayRecord
	err := db.Where("user_id = ? AND content_id = ?", userID, contentID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetUserPlayHistory 获取用户播放历史
func (r *MeditationPlayRecordRepository) GetUserPlayHistory(userID int64, page, pageSize int) ([]model.MeditationPlayRecord, int64, error) {
	db := database.GetDB()
	var records []model.MeditationPlayRecord
	var total int64

	// 获取总数
	if err := db.Model(&model.MeditationPlayRecord{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载内容
	offset := (page - 1) * pageSize
	if err := db.Where("user_id = ?", userID).
		Preload("Content").
		Order("last_played_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetRecentlyPlayed 获取最近播放
func (r *MeditationPlayRecordRepository) GetRecentlyPlayed(userID int64, limit int) ([]model.MeditationPlayRecord, error) {
	db := database.GetDB()
	var records []model.MeditationPlayRecord

	if err := db.Where("user_id = ?", userID).
		Preload("Content").
		Order("last_played_at DESC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

// GetContinuePlaying 获取继续播放列表（未完成的）
func (r *MeditationPlayRecordRepository) GetContinuePlaying(userID int64, limit int) ([]model.MeditationPlayRecord, error) {
	db := database.GetDB()
	var records []model.MeditationPlayRecord

	if err := db.Where("user_id = ? AND completed = false AND progress > 0", userID).
		Preload("Content").
		Order("last_played_at DESC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

// GetPlayStats 获取播放统计
func (r *MeditationPlayRecordRepository) GetPlayStats(userID int64) (map[string]interface{}, error) {
	db := database.GetDB()

	var totalPlays int64
	var totalDuration int64
	var completedCount int64

	// 总播放次数
	db.Model(&model.MeditationPlayRecord{}).Where("user_id = ?", userID).
		Select("SUM(play_count)").Row().Scan(&totalPlays)

	// 总播放时长
	db.Model(&model.MeditationPlayRecord{}).Where("user_id = ?", userID).
		Select("SUM(duration)").Row().Scan(&totalDuration)

	// 完成次数
	db.Model(&model.MeditationPlayRecord{}).Where("user_id = ? AND completed = true", userID).
		Count(&completedCount)

	return map[string]interface{}{
		"total_plays":     totalPlays,
		"total_duration":  totalDuration,
		"completed_count": completedCount,
	}, nil
}
