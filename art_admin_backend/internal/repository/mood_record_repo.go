package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"time"
)

type MoodRecordRepository struct{}

func NewMoodRecordRepository() *MoodRecordRepository {
	return &MoodRecordRepository{}
}

// FindByID 根据ID查询情绪记录
func (r *MoodRecordRepository) FindByID(id int64) (*model.MoodRecord, error) {
	var record model.MoodRecord
	err := database.DB.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Create 创建情绪记录
func (r *MoodRecordRepository) Create(record *model.MoodRecord) error {
	return database.DB.Create(record).Error
}

// Update 更新情绪记录
func (r *MoodRecordRepository) Update(record *model.MoodRecord) error {
	return database.DB.Save(record).Error
}

// Delete 删除情绪记录
func (r *MoodRecordRepository) Delete(id int64) error {
	return database.DB.Delete(&model.MoodRecord{}, id).Error
}

// FindWithPagination 分页查询情绪记录
func (r *MoodRecordRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.MoodRecord, int64, error) {
	var records []model.MoodRecord
	var total int64

	db := database.DB.Model(&model.MoodRecord{})

	// 添加查询条件
	if userID, ok := query["user_id"]; ok {
		db = db.Where("user_id = ?", userID)
	}
	if moodType, ok := query["mood_type"]; ok {
		db = db.Where("mood_type = ?", moodType)
	}
	if startDate, ok := query["start_date"]; ok {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate, ok := query["end_date"]; ok {
		db = db.Where("created_at <= ?", endDate)
	}
	if minIntensity, ok := query["min_intensity"]; ok {
		db = db.Where("intensity >= ?", minIntensity)
	}
	if maxIntensity, ok := query["max_intensity"]; ok {
		db = db.Where("intensity <= ?", maxIntensity)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (current - 1) * size
	err = db.Order("created_at DESC").Offset(offset).Limit(size).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// FindByUserID 根据用户ID查询情绪记录
func (r *MoodRecordRepository) FindByUserID(userID int64) ([]model.MoodRecord, error) {
	var records []model.MoodRecord
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&records).Error
	return records, err
}

// FindByUserIDAndDateRange 根据用户ID和日期范围查询情绪记录
func (r *MoodRecordRepository) FindByUserIDAndDateRange(userID int64, days int) ([]model.MoodRecord, error) {
	var records []model.MoodRecord
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	err := database.DB.Where("user_id = ? AND created_at BETWEEN ? AND ?",
		userID, startDate, endDate).
		Order("created_at DESC").
		Find(&records).Error

	return records, err
}

// FindByUserIDAndDate 根据用户ID和具体日期查询情绪记录
func (r *MoodRecordRepository) FindByUserIDAndDate(userID int64, date time.Time) ([]model.MoodRecord, error) {
	var records []model.MoodRecord
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := database.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?",
		userID, startOfDay, endOfDay).
		Order("created_at DESC").
		Find(&records).Error

	return records, err
}

// GetMoodCountByType 获取用户各类型情绪的数量统计
func (r *MoodRecordRepository) GetMoodCountByType(userID int64, days int) (map[string]int, error) {
	var results []struct {
		MoodType string `json:"mood_type"`
		Count    int    `json:"count"`
	}

	startDate := time.Now().AddDate(0, 0, -days)

	err := database.DB.Model(&model.MoodRecord{}).
		Select("mood_type, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("mood_type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, result := range results {
		countMap[result.MoodType] = result.Count
	}

	return countMap, nil
}

// GetAverageIntensity 获取用户平均情绪强度
func (r *MoodRecordRepository) GetAverageIntensity(userID int64, days int) (float64, error) {
	var avgIntensity float64

	startDate := time.Now().AddDate(0, 0, -days)

	err := database.DB.Model(&model.MoodRecord{}).
		Select("AVG(intensity)").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Scan(&avgIntensity).Error

	return avgIntensity, err
}

// GetLatestMoodRecord 获取用户最新的情绪记录
func (r *MoodRecordRepository) GetLatestMoodRecord(userID int64) (*model.MoodRecord, error) {
	var record model.MoodRecord
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByUserIDAndDateString 根据用户ID和日期字符串查询情绪记录
func (r *MoodRecordRepository) FindByUserIDAndDateString(userID int64, dateStr string) ([]model.MoodRecord, error) {
	var records []model.MoodRecord

	// 解析日期字符串 (格式: 2006-01-02)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err = database.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?",
		userID, startOfDay, endOfDay).
		Order("created_at DESC").
		Find(&records).Error

	return records, err
}

// GetCalendarData 获取日历视图数据（按月统计每天的情绪记录）
func (r *MoodRecordRepository) GetCalendarData(userID int64, year int, month int) ([]map[string]interface{}, error) {
	var results []struct {
		Date         string  `gorm:"column:date"`
		RecordCount  int     `gorm:"column:record_count"`
		AvgIntensity float64 `gorm:"column:avg_intensity"`
	}

	// 计算月份的开始和结束日期
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	err := database.DB.Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m-%d') as date,
			COUNT(*) as record_count,
			AVG(intensity) as avg_intensity
		FROM mood_records 
		WHERE user_id = ? AND created_at >= ? AND created_at < ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, userID, startDate, endDate).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为map切片
	calendarData := make([]map[string]interface{}, len(results))
	for i, result := range results {
		calendarData[i] = map[string]interface{}{
			"date":          result.Date,
			"record_count":  result.RecordCount,
			"avg_intensity": result.AvgIntensity,
		}
	}

	return calendarData, nil
}

// GetMoodStreak 获取用户情绪记录连续天数
func (r *MoodRecordRepository) GetMoodStreak(userID int64) (int, error) {
	type DateRecord struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var records []DateRecord

	// 获取最近30天的记录，按日期分组
	err := database.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM mood_records 
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
