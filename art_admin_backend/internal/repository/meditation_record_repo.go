package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// MeditationRecordRepository 冥想记录仓库
type MeditationRecordRepository struct {
	db *gorm.DB
}

// NewMeditationRecordRepository 创建冥想记录仓库
func NewMeditationRecordRepository() *MeditationRecordRepository {
	return &MeditationRecordRepository{
		db: database.GetDB(),
	}
}

// Create 创建冥想记录
func (r *MeditationRecordRepository) Create(record *model.MeditationRecord) error {
	return r.db.Create(record).Error
}

// FindByID 根据ID查找冥想记录
func (r *MeditationRecordRepository) FindByID(id int64) (*model.MeditationRecord, error) {
	var record model.MeditationRecord
	err := r.db.Preload("Content").First(&record, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("冥想记录不存在")
		}
		return nil, err
	}
	return &record, nil
}

// Update 更新冥想记录
func (r *MeditationRecordRepository) Update(record *model.MeditationRecord) error {
	return r.db.Save(record).Error
}

// Delete 删除冥想记录
func (r *MeditationRecordRepository) Delete(id int64) error {
	return r.db.Delete(&model.MeditationRecord{}, id).Error
}

// FindWithPagination 分页查询冥想记录
func (r *MeditationRecordRepository) FindWithPagination(query map[string]interface{}, page, pageSize int) ([]model.MeditationRecord, int64, error) {
	var records []model.MeditationRecord
	var total int64

	// 构建查询
	db := r.db.Model(&model.MeditationRecord{}).Preload("Content")

	// 添加查询条件
	if userID, ok := query["user_id"].(int64); ok && userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if contentID, ok := query["content_id"].(int64); ok && contentID > 0 {
		db = db.Where("content_id = ?", contentID)
	}
	if meditationType, ok := query["meditation_type"].(string); ok && meditationType != "" {
		db = db.Where("meditation_type = ?", meditationType)
	}
	if difficultyLevel, ok := query["difficulty_level"].(int); ok && difficultyLevel > 0 {
		db = db.Where("difficulty_level = ?", difficultyLevel)
	}
	if analysisStatus, ok := query["analysis_status"].(string); ok && analysisStatus != "" {
		db = db.Where("analysis_status = ?", analysisStatus)
	}
	if minDuration, ok := query["min_duration"].(int); ok && minDuration > 0 {
		db = db.Where("duration >= ?", minDuration)
	}
	if maxDuration, ok := query["max_duration"].(int); ok && maxDuration > 0 {
		db = db.Where("duration <= ?", maxDuration)
	}
	if startDate, ok := query["start_date"].(time.Time); ok && !startDate.IsZero() {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate, ok := query["end_date"].(time.Time); ok && !endDate.IsZero() {
		db = db.Where("created_at <= ?", endDate)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// FindByUserID 根据用户ID查找冥想记录
func (r *MeditationRecordRepository) FindByUserID(userID int64, page, pageSize int) ([]model.MeditationRecord, int64, error) {
	query := map[string]interface{}{
		"user_id": userID,
	}
	return r.FindWithPagination(query, page, pageSize)
}

// FindByUserIDAndDateRange 根据用户ID和日期范围查找冥想记录
func (r *MeditationRecordRepository) FindByUserIDAndDateRange(userID int64, startDate, endDate time.Time, page, pageSize int) ([]model.MeditationRecord, int64, error) {
	query := map[string]interface{}{
		"user_id":    userID,
		"start_date": startDate,
		"end_date":   endDate,
	}
	return r.FindWithPagination(query, page, pageSize)
}

// FindByUserIDAndDate 根据用户ID和具体日期查找冥想记录
func (r *MeditationRecordRepository) FindByUserIDAndDate(userID int64, date time.Time) ([]model.MeditationRecord, error) {
	var records []model.MeditationRecord
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Preload("Content").Where("user_id = ? AND created_at >= ? AND created_at < ?",
		userID, startOfDay, endOfDay).
		Order("created_at DESC").
		Find(&records).Error

	return records, err
}

// FindRecent 查找最近的冥想记录
func (r *MeditationRecordRepository) FindRecent(userID int64, limit int) ([]model.MeditationRecord, error) {
	var records []model.MeditationRecord
	query := r.db.Preload("Content").Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&records).Error
	return records, err
}

// FindByType 根据冥想类型查找记录
func (r *MeditationRecordRepository) FindByType(userID int64, meditationType string, limit int) ([]model.MeditationRecord, error) {
	var records []model.MeditationRecord
	query := r.db.Preload("Content").Where("user_id = ? AND meditation_type = ?", userID, meditationType)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&records).Error
	return records, err
}

// FindPendingAnalysis 查找等待分析的记录
func (r *MeditationRecordRepository) FindPendingAnalysis(limit int) ([]model.MeditationRecord, error) {
	var records []model.MeditationRecord
	query := r.db.Where("analysis_status = ?", model.AnalysisStatusPending)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at ASC").Find(&records).Error
	return records, err
}

// UpdateAnalysisStatus 更新分析状态
func (r *MeditationRecordRepository) UpdateAnalysisStatus(id int64, status string) error {
	return r.db.Model(&model.MeditationRecord{}).Where("id = ?", id).Update("analysis_status", status).Error
}

// UpdateAnalysisResult 更新分析结果
func (r *MeditationRecordRepository) UpdateAnalysisResult(id int64, result string, difyDocID string) error {
	updates := map[string]interface{}{
		"analysis_result": result,
		"analysis_status": model.AnalysisStatusCompleted,
	}
	if difyDocID != "" {
		updates["dify_document_id"] = difyDocID
	}
	return r.db.Model(&model.MeditationRecord{}).Where("id = ?", id).Updates(updates).Error
}

// GetMeditationStats 获取冥想统计信息
func (r *MeditationRecordRepository) GetMeditationStats(userID int64, days int) (*model.MeditationAnalytics, error) {
	var stats model.MeditationAnalytics

	// 计算日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 查询指定时间范围内的冥想记录
	var records []model.MeditationRecord
	if err := r.db.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).Find(&records).Error; err != nil {
		return nil, err
	}

	// 填充统计数据
	stats.UserID = userID
	stats.TotalSessions = len(records)
	stats.TotalDuration = 0
	stats.AverageDuration = 0
	stats.MostUsedType = ""
	stats.AverageImprovement = 0

	if len(records) > 0 {
		totalDuration := 0
		typeCount := make(map[string]int)
		totalImprovement := 0.0
		improvementCount := 0

		for _, record := range records {
			totalDuration += record.Duration
			typeCount[record.MeditationType]++

			// 计算情绪改善程度
			if record.PreMoodIntensity > 0 && record.PostMoodIntensity > 0 {
				improvement := float64(record.PreMoodIntensity-record.PostMoodIntensity) / 10.0
				totalImprovement += improvement
				improvementCount++
			}
		}

		stats.TotalDuration = totalDuration
		stats.AverageDuration = float64(totalDuration) / float64(len(records))

		// 找出最常用的冥想类型
		maxCount := 0
		for meditationType, count := range typeCount {
			if count > maxCount {
				maxCount = count
				stats.MostUsedType = meditationType
			}
		}

		// 计算平均改善程度
		if improvementCount > 0 {
			stats.AverageImprovement = totalImprovement / float64(improvementCount)
		}

		// 构建类型分布JSON
		distribution := make(map[string]int)
		for meditationType, count := range typeCount {
			distribution[meditationType] = count
		}
		if distributionJSON, err := r.formatJSON(distribution); err == nil {
			stats.Distribution = distributionJSON
		}
	}

	stats.AnalyzeType = "daily"
	stats.AnalyzeDate = endDate
	stats.CreatedAt = time.Now()

	return &stats, nil
}

// GetMeditationStreak 获取用户冥想连续天数
func (r *MeditationRecordRepository) GetMeditationStreak(userID int64) (int, error) {
	type DateRecord struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var records []DateRecord

	// 获取最近30天的记录，按日期分组
	err := r.db.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM meditation_records 
		WHERE user_id = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, userID).Scan(&records).Error

	if err != nil {
		return 0, err
	}

	if len(records) == 0 {
		return 0, nil
	}

	// 计算连续天数
	streak := 0
	expectedDate := time.Now().Truncate(24 * time.Hour)

	for _, record := range records {
		recordDate := record.Date.Truncate(24 * time.Hour)

		if recordDate.Equal(expectedDate) {
			streak++
			expectedDate = expectedDate.AddDate(0, 0, -1)
		} else if recordDate.Before(expectedDate) {
			// 如果记录日期早于期望日期，说明连续记录中断
			break
		}
	}

	return streak, nil
}

// CountByUserID 统计用户的冥想记录数量
func (r *MeditationRecordRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.MeditationRecord{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Search 搜索冥想记录
func (r *MeditationRecordRepository) Search(userID int64, keyword string, page, pageSize int) ([]model.MeditationRecord, int64, error) {
	var records []model.MeditationRecord
	var total int64

	// 构建搜索查询
	db := r.db.Model(&model.MeditationRecord{}).Preload("Content").
		Where("user_id = ? AND (experience LIKE ? OR notes LIKE ?)",
			userID, "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// formatJSON 格式化JSON的辅助函数
func (r *MeditationRecordRepository) formatJSON(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
