package repository

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// JournalAnalytics 日记分析统计
type JournalAnalytics struct {
	UserID             int64     `json:"user_id"`
	TotalEntries       int       `json:"total_entries"`
	AverageWordCount   float64   `json:"average_word_count"`
	MostCommonMoodType string    `json:"most_common_mood_type"`
	AnalysisDate       time.Time `json:"analysis_date"`
	CreatedAt          time.Time `json:"created_at"`
}

// JournalEntryRepository 日记条目仓库
type JournalEntryRepository struct {
	db *gorm.DB
}

// NewJournalEntryRepository 创建日记条目仓库
func NewJournalEntryRepository() *JournalEntryRepository {
	return &JournalEntryRepository{
		db: database.GetDB(),
	}
}

// Create 创建日记条目
func (r *JournalEntryRepository) Create(entry *model.JournalEntry) error {
	return r.db.Create(entry).Error
}

// FindByID 根据ID查找日记条目
func (r *JournalEntryRepository) FindByID(id int64) (*model.JournalEntry, error) {
	var entry model.JournalEntry
	err := r.db.First(&entry, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("日记条目不存在")
		}
		return nil, err
	}

	return &entry, nil
}

// Update 更新日记条目
func (r *JournalEntryRepository) Update(entry *model.JournalEntry) error {
	return r.db.Save(entry).Error
}

// Delete 删除日记条目
func (r *JournalEntryRepository) Delete(id int64) error {
	return r.db.Delete(&model.JournalEntry{}, id).Error
}

// FindWithPagination 分页查询日记条目
func (r *JournalEntryRepository) FindWithPagination(query map[string]interface{}, page, pageSize int) ([]model.JournalEntry, int64, error) {
	var entries []model.JournalEntry
	var total int64

	// 构建查询
	db := r.db.Model(&model.JournalEntry{})

	// 添加查询条件
	if userID, ok := query["user_id"].(int64); ok && userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if moodRecordID, ok := query["mood_record_id"].(int64); ok && moodRecordID > 0 {
		db = db.Where("mood_record_id = ?", moodRecordID)
	}
	if isPublic, ok := query["is_public"].(bool); ok {
		db = db.Where("is_public = ?", isPublic)
	}
	if startDate, ok := query["start_date"].(time.Time); ok && !startDate.IsZero() {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate, ok := query["end_date"].(time.Time); ok && !endDate.IsZero() {
		db = db.Where("created_at <= ?", endDate)
	}
	if hasAttachment, ok := query["has_attachment"].(bool); ok {
		if hasAttachment {
			db = db.Where("attachment_url IS NOT NULL AND attachment_url != ''")
		} else {
			db = db.Where("attachment_url IS NULL OR attachment_url = ''")
		}
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// FindByUserID 根据用户ID查找日记条目
func (r *JournalEntryRepository) FindByUserID(userID int64, page, pageSize int) ([]model.JournalEntry, int64, error) {
	query := map[string]interface{}{
		"user_id": userID,
	}
	return r.FindWithPagination(query, page, pageSize)
}

// FindByDateRange 根据日期范围查找日记条目
func (r *JournalEntryRepository) FindByDateRange(userID int64, startDate, endDate time.Time, page, pageSize int) ([]model.JournalEntry, int64, error) {
	query := map[string]interface{}{
		"user_id":    userID,
		"start_date": startDate,
		"end_date":   endDate,
	}
	return r.FindWithPagination(query, page, pageSize)
}

// FindByMoodRecordID 根据情绪记录ID查找日记条目
func (r *JournalEntryRepository) FindByMoodRecordID(moodRecordID int64) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry
	err := r.db.Where("mood_record_id = ?", moodRecordID).Order("created_at DESC").Find(&entries).Error
	return entries, err
}

// FindRecent 查找最近的日记条目
func (r *JournalEntryRepository) FindRecent(userID int64, limit int) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry
	query := r.db.Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&entries).Error
	return entries, err
}

// Search 搜索日记条目
func (r *JournalEntryRepository) Search(userID int64, keyword string, page, pageSize int) ([]model.JournalEntry, int64, error) {
	var entries []model.JournalEntry
	var total int64

	// 构建搜索查询
	db := r.db.Model(&model.JournalEntry{}).Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// GetJournalStats 获取日记统计信息
func (r *JournalEntryRepository) GetJournalStats(userID int64, days int) (*JournalAnalytics, error) {
	var stats JournalAnalytics

	// 计算日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 查询指定时间范围内的日记
	var entries []model.JournalEntry
	if err := r.db.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).Find(&entries).Error; err != nil {
		return nil, err
	}

	// 填充统计数据
	stats.UserID = userID
	stats.TotalEntries = len(entries)
	stats.AverageWordCount = 0
	stats.MostCommonMoodType = ""

	if len(entries) > 0 {
		totalWordCount := 0
		moodCount := make(map[string]int)

		for _, entry := range entries {
			// 计算字数（简单统计）
			totalWordCount += len([]rune(entry.Content))

			// 统计情绪类型（如果有关联的情绪记录）
			if entry.MoodRecordID != nil && *entry.MoodRecordID > 0 {
				var moodRecord model.MoodRecord
				if err := r.db.First(&moodRecord, *entry.MoodRecordID).Error; err == nil {
					moodCount[moodRecord.MoodType]++
				}
			}
		}

		stats.AverageWordCount = float64(totalWordCount) / float64(len(entries))

		// 找出最常见的情绪类型
		maxCount := 0
		for moodType, count := range moodCount {
			if count > maxCount {
				maxCount = count
				stats.MostCommonMoodType = moodType
			}
		}
	}

	stats.AnalysisDate = endDate
	stats.CreatedAt = time.Now()

	return &stats, nil
}

// UpdateSentimentScore 更新情感得分
func (r *JournalEntryRepository) UpdateSentimentScore(id int64, score float64) error {
	return r.db.Model(&model.JournalEntry{}).Where("id = ?", id).Update("sentiment_score", score).Error
}

// CountByUserID 统计用户的日记数量
func (r *JournalEntryRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.JournalEntry{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetWritingStreak 获取连续写作天数
func (r *JournalEntryRepository) GetWritingStreak(userID int64) (int, error) {
	var entries []model.JournalEntry

	// 获取最近30天的日记条目，按日期降序排列
	startDate := time.Now().AddDate(0, 0, -30)
	if err := r.db.Where("user_id = ? AND created_at >= ?", userID, startDate).Order("created_at DESC").Find(&entries).Error; err != nil {
		return 0, err
	}

	if len(entries) == 0 {
		return 0, nil
	}

	streak := 0
	currentDate := time.Now().Truncate(24 * time.Hour)

	for _, entry := range entries {
		entryDate := entry.CreatedAt.Truncate(24 * time.Hour)

		// 如果是今天或昨天的日记，开始计算连续天数
		if entryDate.Equal(currentDate) || entryDate.Equal(currentDate.AddDate(0, 0, -1)) {
			streak++
			currentDate = entryDate.AddDate(0, 0, -1)
		} else if entryDate.Before(currentDate) {
			// 如果有断层，停止计算
			break
		}
	}

	return streak, nil
}

// JournalStatistics 日记统计信息结构
type JournalStatistics struct {
	TotalEntries     int64     `json:"total_entries"`
	AverageWordCount float64   `json:"average_word_count"`
	MostCommonMood   string    `json:"most_common_mood"`
	WritingStreak    int       `json:"writing_streak"`
	AnalysisDate     time.Time `json:"analysis_date"`
}

// GetJournalStatistics 获取日记统计信息
func (r *JournalEntryRepository) GetJournalStatistics(userID int64, days int) (*JournalStatistics, error) {
	db := database.GetDB()

	// 计算起始时间
	startTime := time.Now().AddDate(0, 0, -days)

	// 基础统计查询
	var totalCount int64
	if err := db.Model(&model.JournalEntry{}).Where("user_id = ? AND created_at >= ?", userID, startTime).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 平均情感分数
	var avgSentiment float64
	db.Model(&model.JournalEntry{}).Where("user_id = ? AND created_at >= ?", userID, startTime).
		Select("AVG(sentiment_score)").Row().Scan(&avgSentiment)

	// 获取连续天数
	streak, _ := r.GetWritingStreak(userID)

	// 构建统计结果
	stats := &JournalStatistics{
		TotalEntries:     totalCount,
		AverageWordCount: avgSentiment, // 临时使用情感分数
		MostCommonMood:   "neutral",    // 临时值
		WritingStreak:    streak,
		AnalysisDate:     time.Now(),
	}

	return stats, nil
}

// UpdateTags 更新日记标签
func (r *JournalEntryRepository) UpdateTags(entryID int64, tags []string) error {
	db := database.GetDB()
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	return db.Model(&model.JournalEntry{}).Where("id = ?", entryID).Update("tags", string(tagsJSON)).Error
}

// AddImages 添加日记图片
func (r *JournalEntryRepository) AddImages(entryID int64, newImages []string) error {
	db := database.GetDB()

	// 获取现有图片
	var entry model.JournalEntry
	if err := db.First(&entry, entryID).Error; err != nil {
		return err
	}

	// 解析现有图片
	var existingImages []string
	if entry.Images != "" {
		json.Unmarshal([]byte(entry.Images), &existingImages)
	}

	// 合并图片
	existingImages = append(existingImages, newImages...)

	// 更新
	imagesJSON, err := json.Marshal(existingImages)
	if err != nil {
		return err
	}
	return db.Model(&model.JournalEntry{}).Where("id = ?", entryID).Update("images", string(imagesJSON)).Error
}

// RemoveImage 删除日记图片
func (r *JournalEntryRepository) RemoveImage(entryID int64, imageURL string) error {
	db := database.GetDB()

	// 获取现有图片
	var entry model.JournalEntry
	if err := db.First(&entry, entryID).Error; err != nil {
		return err
	}

	// 解析现有图片
	var existingImages []string
	if entry.Images != "" {
		json.Unmarshal([]byte(entry.Images), &existingImages)
	}

	// 移除指定图片
	var newImages []string
	for _, img := range existingImages {
		if img != imageURL {
			newImages = append(newImages, img)
		}
	}

	// 更新
	imagesJSON, err := json.Marshal(newImages)
	if err != nil {
		return err
	}
	return db.Model(&model.JournalEntry{}).Where("id = ?", entryID).Update("images", string(imagesJSON)).Error
}

// GetAllUserTags 获取用户所有标签
func (r *JournalEntryRepository) GetAllUserTags(userID int64) ([]string, error) {
	db := database.GetDB()

	var entries []model.JournalEntry
	if err := db.Where("user_id = ? AND tags IS NOT NULL AND tags != ''", userID).
		Select("tags").Find(&entries).Error; err != nil {
		return nil, err
	}

	// 收集所有标签并去重
	tagSet := make(map[string]bool)
	for _, entry := range entries {
		var tags []string
		if err := json.Unmarshal([]byte(entry.Tags), &tags); err == nil {
			for _, tag := range tags {
				tagSet[tag] = true
			}
		}
	}

	// 转换为切片
	result := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		result = append(result, tag)
	}

	return result, nil
}

// FindByTag 按标签查找日记
func (r *JournalEntryRepository) FindByTag(userID int64, tag string, page, pageSize int) ([]model.JournalEntry, int64, error) {
	db := database.GetDB()

	var entries []model.JournalEntry
	var total int64

	// 使用JSON_CONTAINS查询包含指定标签的日记
	query := db.Model(&model.JournalEntry{}).
		Where("user_id = ? AND JSON_CONTAINS(tags, ?)", userID, `"`+tag+`"`)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}
