package service

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// JournalEntryService 日记条目服务
type JournalEntryService struct {
	journalEntryRepo  *repository.JournalEntryRepository
	moodRecordRepo    *repository.MoodRecordRepository
	aiAnalysisService *AIAnalysisService
}

// NewJournalEntryService 创建日记条目服务
func NewJournalEntryService() *JournalEntryService {
	return &JournalEntryService{
		journalEntryRepo:  repository.NewJournalEntryRepository(),
		moodRecordRepo:    repository.NewMoodRecordRepository(),
		aiAnalysisService: NewAIAnalysisService(),
	}
}

// CreateJournalEntry 创建日记条目
func (s *JournalEntryService) CreateJournalEntry(userID int64, req *request.CreateJournalEntryRequest) error {
	// 验证标题和内容
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("日记标题不能为空")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("日记内容不能为空")
	}

	// 验证字数限制
	if len([]rune(req.Content)) > 10000 {
		return errors.New("日记内容不能超过10000字")
	}

	// 验证关联的情绪记录是否存在（如果提供）
	if req.MoodRecordID != nil && *req.MoodRecordID > 0 {
		_, err := s.moodRecordRepo.FindByID(*req.MoodRecordID)
		if err != nil {
			return errors.New("关联的情绪记录不存在")
		}
	}

	// 创建模型
	entry := &model.JournalEntry{
		UserID:       userID,
		Title:        req.Title,
		Content:      req.Content,
		MoodRecordID: req.MoodRecordID,
		IsPrivate:    req.IsPrivate,
	}

	// 保存到数据库
	if err := s.journalEntryRepo.Create(entry); err != nil {
		return errors.New("创建日记条目失败")
	}

	// 异步触发AI分析
	go s.triggerAsyncAnalysis(userID)

	return nil
}

// UpdateJournalEntry 更新日记条目
func (s *JournalEntryService) UpdateJournalEntry(userID int64, req *request.UpdateJournalEntryRequest) error {
	// 获取现有条目
	existingEntry, err := s.journalEntryRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限（只有创建者可以修改）
	if existingEntry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 验证标题和内容（如果提供）
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return errors.New("日记标题不能为空")
	}
	if req.Content != nil && strings.TrimSpace(*req.Content) == "" {
		return errors.New("日记内容不能为空")
	}

	// 验证字数限制（如果提供）
	if req.Content != nil && len([]rune(*req.Content)) > 10000 {
		return errors.New("日记内容不能超过10000字")
	}

	// 验证关联的情绪记录（如果提供）
	if req.MoodRecordID != nil {
		if *req.MoodRecordID > 0 {
			_, err := s.moodRecordRepo.FindByID(*req.MoodRecordID)
			if err != nil {
				return errors.New("关联的情绪记录不存在")
			}
		}
		existingEntry.MoodRecordID = getIntPtr(req.MoodRecordID)
	}

	// 更新字段
	if req.Title != nil {
		existingEntry.Title = *req.Title
	}
	if req.Content != nil {
		existingEntry.Content = *req.Content
	}
	if req.MoodRecordID != nil {
		existingEntry.MoodRecordID = req.MoodRecordID
	}
	if req.IsPrivate != nil {
		existingEntry.IsPrivate = *req.IsPrivate
	}

	existingEntry.UpdatedAt = time.Now()

	// 更新数据库
	return s.journalEntryRepo.Update(existingEntry)
}

// DeleteJournalEntry 删除日记条目
func (s *JournalEntryService) DeleteJournalEntry(userID, id int64) error {
	// 获取条目
	entry, err := s.journalEntryRepo.FindByID(id)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限（只有创建者可以删除）
	if entry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 删除条目
	return s.journalEntryRepo.Delete(id)
}

// GetJournalEntryList 获取日记条目列表
func (s *JournalEntryService) GetJournalEntryList(userID int64, req *request.JournalEntryListRequest) ([]response.JournalEntryItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	query["user_id"] = userID
	if req.MoodRecordID != nil {
		query["mood_record_id"] = *req.MoodRecordID
	}
	if req.IsPrivate != nil {
		query["is_private"] = *req.IsPrivate
	}
	if req.StartDate != nil {
		query["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		query["end_date"] = *req.EndDate
	}
	if req.HasAttachment != nil {
		query["has_attachment"] = *req.HasAttachment
	}

	// 查询条目
	entries, total, err := s.journalEntryRepo.FindWithPagination(query, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]response.JournalEntryItem, len(entries))
	for i, entry := range entries {
		items[i] = s.convertToJournalEntryItem(&entry)
	}

	return items, total, nil
}

// GetJournalEntry 获取单个日记条目
func (s *JournalEntryService) GetJournalEntry(userID, entryID int64) (*response.JournalEntryItem, error) {
	// 获取条目
	entry, err := s.journalEntryRepo.FindByID(entryID)
	if err != nil {
		return nil, errors.New("日记条目不存在")
	}

	// 验证权限
	if entry.UserID != userID {
		return nil, errors.New("无权限操作此记录")
	}

	// 转换为响应格式
	item := s.convertToJournalEntryItem(entry)
	return &item, nil
}

// GetJournalEntryDetail 获取日记条目详情
func (s *JournalEntryService) GetJournalEntryDetail(userID, id int64) (*response.JournalEntryDetail, error) {
	// 获取条目
	entry, err := s.journalEntryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 验证权限（只有创建者可以查看私有日记）
	if entry.UserID != userID && entry.IsPrivate {
		return nil, errors.New("无权限查看此日记")
	}

	// 转换为响应格式
	detail := &response.JournalEntryDetail{
		JournalEntryItem: s.convertToJournalEntryItem(entry),
		AttachmentURL:    "",        // 模型中没有此字段
		PrivacySetting:   "private", // 模型中没有此字段，使用默认值
		UpdatedAt:        entry.UpdatedAt,
	}

	return detail, nil
}

// GetJournalEntriesByDateRange 根据日期范围获取日记条目
func (s *JournalEntryService) GetJournalEntriesByDateRange(userID int64, startDate, endDate time.Time, page, pageSize int) ([]response.JournalEntryItem, int64, error) {
	entries, total, err := s.journalEntryRepo.FindByDateRange(userID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.JournalEntryItem, len(entries))
	for i, entry := range entries {
		items[i] = s.convertToJournalEntryItem(&entry)
	}

	return items, total, nil
}

// GetRecentJournalEntries 获取最近的日记条目
func (s *JournalEntryService) GetRecentJournalEntries(userID int64, limit int) ([]response.JournalEntryItem, error) {
	entries, err := s.journalEntryRepo.FindRecent(userID, limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.JournalEntryItem, len(entries))
	for i, entry := range entries {
		items[i] = s.convertToJournalEntryItem(&entry)
	}

	return items, nil
}

// SearchJournalEntries 搜索日记条目
func (s *JournalEntryService) SearchJournalEntries(userID int64, keyword string, page, pageSize int) ([]response.JournalEntryItem, int64, error) {
	entries, total, err := s.journalEntryRepo.Search(userID, keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.JournalEntryItem, len(entries))
	for i, entry := range entries {
		items[i] = s.convertToJournalEntryItem(&entry)
	}

	return items, total, nil
}

// GetJournalStatistics 获取日记统计信息
func (s *JournalEntryService) GetJournalStatistics(userID int64, days int) (*response.JournalStatistics, error) {
	stats, err := s.journalEntryRepo.GetJournalStats(userID, days)
	if err != nil {
		return nil, err
	}

	// 获取连续写作天数
	streak, err := s.journalEntryRepo.GetWritingStreak(userID)
	if err != nil {
		streak = 0
	}

	// 获取总日记数量
	totalCount, err := s.journalEntryRepo.CountByUserID(userID)
	if err != nil {
		totalCount = 0
	}

	// 转换为响应格式
	result := &response.JournalStatistics{
		TotalEntries:       stats.TotalEntries,
		TotalCount:         totalCount,
		AverageWordCount:   stats.AverageWordCount,
		MostCommonMoodType: stats.MostCommonMoodType,
		WritingStreak:      streak,
		AnalysisDate:       stats.AnalysisDate,
	}

	return result, nil
}

// UpdateSentimentScore 更新情感得分
func (s *JournalEntryService) UpdateSentimentScore(userID, id int64, score float64) error {
	// 获取条目
	entry, err := s.journalEntryRepo.FindByID(id)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限（只有创建者可以修改）
	if entry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 更新情感得分
	return s.journalEntryRepo.UpdateSentimentScore(id, score)
}

// convertToJournalEntryItem 转换模型为响应项
func (s *JournalEntryService) convertToJournalEntryItem(entry *model.JournalEntry) response.JournalEntryItem {
	// 转换为响应格式
	var sentimentScore float64
	if entry.SentimentScore != nil {
		sentimentScore = *entry.SentimentScore
	}

	item := &response.JournalEntryItem{
		ID:             entry.ID,
		UserID:         entry.UserID,
		Title:          entry.Title,
		Content:        truncateContent(entry.Content, 200), // 只显示前200字
		MoodRecordID:   entry.MoodRecordID,
		SentimentScore: sentimentScore,
		IsPublic:       !entry.IsPrivate, // 响应中使用IsPublic，模型中使用IsPrivate
		CreatedAt:      entry.CreatedAt,
	}

	return *item
}

// getIntPtr 获取int64指针的辅助函数
func getIntPtr(ptr *int64) *int64 {
	if ptr == nil || *ptr == 0 {
		return nil
	}
	return ptr
}

// truncateContent 截取内容预览
func truncateContent(content string, maxLen int) string {
	if len([]rune(content)) <= maxLen {
		return content
	}
	runes := []rune(content)
	return string(runes[:maxLen]) + "..."
}

// AnalyzeJournalEntry 分析日记条目
func (s *JournalEntryService) AnalyzeJournalEntry(ctx context.Context, userID, entryID int64) (*response.JournalEntryAnalysis, error) {
	// 获取日记条目
	entry, err := s.journalEntryRepo.FindByID(entryID)
	if err != nil {
		return nil, errors.New("日记条目不存在")
	}

	// 验证权限
	if entry.UserID != userID {
		return nil, errors.New("无权限分析此日记")
	}

	// 构建分析请求
	req := &EmotionAnalysisRequest{
		UserID:       userID,
		StartDate:    entry.CreatedAt,
		EndDate:      entry.CreatedAt.Add(24 * time.Hour),
		AnalysisType: "journal",
	}

	// 调用AI分析服务
	result, err := s.aiAnalysisService.AnalyzeUserEmotions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 转换为响应格式
	analysis := &response.JournalEntryAnalysis{
		EntryID:         entryID,
		Title:           entry.Title,
		AnalysisType:    result.AnalysisType,
		Summary:         result.Summary,
		Insights:        result.Insights,
		Recommendations: result.Recommendations,
		AnalyzedAt:      result.AnalyzedAt,
	}

	return analysis, nil
}

// GetJournalInsights 获取日记洞察
func (s *JournalEntryService) GetJournalInsights(ctx context.Context, userID int64, days int) (*response.JournalInsights, error) {
	// 构建分析请求
	req := &EmotionAnalysisRequest{
		UserID:       userID,
		StartDate:    time.Now().AddDate(0, 0, -days),
		EndDate:      time.Now(),
		AnalysisType: "insights",
	}

	// 调用AI分析服务
	result, err := s.aiAnalysisService.AnalyzeUserEmotions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 获取日记统计信息
	stats, err := s.journalEntryRepo.GetJournalStatistics(userID, days)
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %w", err)
	}

	// 转换为响应格式
	insights := &response.JournalInsights{
		UserID:           userID,
		AnalysisType:     result.AnalysisType,
		StartDate:        result.StartDate,
		EndDate:          result.EndDate,
		TotalEntries:     int(stats.TotalEntries),
		AverageWordCount: stats.AverageWordCount,
		MostCommonMood:   stats.MostCommonMood,
		WritingStreak:    stats.WritingStreak,
		Summary:          result.Summary,
		Patterns:         s.convertEmotionPatterns(result.Patterns),
		Trends:           s.convertEmotionTrends(result.Trends),
		Insights:         result.Insights,
		Recommendations:  result.Recommendations,
		AnalyzedAt:       result.AnalyzedAt,
		DataPoints:       result.DataPoints,
	}

	return insights, nil
}

// GetJournalTrends 获取日记趋势
func (s *JournalEntryService) GetJournalTrends(userID int64, days int) ([]response.JournalTrendData, error) {
	// 获取指定天数内的日记记录
	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	entries, _, err := s.journalEntryRepo.FindByDateRange(userID, startDate, endDate, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	trendMap := make(map[string]response.JournalTrendData)
	for _, entry := range entries {
		date := entry.CreatedAt.Format("2006-01-02")
		wordCount := len([]rune(entry.Content))

		if existing, exists := trendMap[date]; exists {
			// 累加数据
			existing.EntryCount++
			existing.TotalWordCount += wordCount
			if entry.SentimentScore != nil {
				existing.AverageSentiment = (existing.AverageSentiment + *entry.SentimentScore) / 2
			}
			trendMap[date] = existing
		} else {
			// 创建新数据
			avgSentiment := 0.0
			if entry.SentimentScore != nil {
				avgSentiment = *entry.SentimentScore
			}

			trendMap[date] = response.JournalTrendData{
				Date:             date,
				EntryCount:       1,
				TotalWordCount:   wordCount,
				AverageWordCount: wordCount,
				AverageSentiment: avgSentiment,
			}
		}
	}

	// 转换为数组并按日期排序
	var trends []response.JournalTrendData
	for _, trend := range trendMap {
		// 计算平均字数
		if trend.EntryCount > 0 {
			trend.AverageWordCount = trend.TotalWordCount / trend.EntryCount
		}
		trends = append(trends, trend)
	}

	// 按日期排序
	for i := 0; i < len(trends)-1; i++ {
		for j := i + 1; j < len(trends); j++ {
			if trends[i].Date > trends[j].Date {
				trends[i], trends[j] = trends[j], trends[i]
			}
		}
	}

	return trends, nil
}

// triggerAsyncAnalysis 触发异步AI分析
func (s *JournalEntryService) triggerAsyncAnalysis(userID int64) {
	// 检查防抖机制
	if !GetAnalysisDebouncer().ShouldCheck(userID) {
		return // 跳过本次分析
	}

	ctx := context.Background()

	// 执行AI分析
	_, err := s.GetJournalInsights(ctx, userID, 7) // 分析最近7天
	if err != nil {
		fmt.Printf("异步日记分析失败，用户ID: %d, 错误: %v\n", userID, err)
	}
}

// convertEmotionPatterns 转换情绪模式
func (s *JournalEntryService) convertEmotionPatterns(patterns []EmotionPattern) []response.EmotionPattern {
	result := make([]response.EmotionPattern, len(patterns))
	for i, pattern := range patterns {
		result[i] = response.EmotionPattern{
			Type:        pattern.Type,
			Frequency:   pattern.Frequency,
			Percentage:  pattern.Percentage,
			Description: pattern.Description,
		}
	}
	return result
}

// convertEmotionTrends 转换情绪趋势
func (s *JournalEntryService) convertEmotionTrends(trends []EmotionTrend) []response.EmotionTrend {
	result := make([]response.EmotionTrend, len(trends))
	for i, trend := range trends {
		result[i] = response.EmotionTrend{
			Period:     trend.Period,
			Direction:  trend.Direction,
			ChangeRate: trend.ChangeRate,
			Confidence: trend.Confidence,
		}
	}
	return result
}

// UpdateJournalTags 更新日记标签
func (s *JournalEntryService) UpdateJournalTags(userID, entryID int64, tags []string) error {
	// 获取日记条目
	entry, err := s.journalEntryRepo.FindByID(entryID)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限
	if entry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 更新标签
	return s.journalEntryRepo.UpdateTags(entryID, tags)
}

// AddJournalImages 添加日记图片
func (s *JournalEntryService) AddJournalImages(userID, entryID int64, images []string) error {
	// 获取日记条目
	entry, err := s.journalEntryRepo.FindByID(entryID)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限
	if entry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 添加图片
	return s.journalEntryRepo.AddImages(entryID, images)
}

// RemoveJournalImage 删除日记图片
func (s *JournalEntryService) RemoveJournalImage(userID, entryID int64, imageURL string) error {
	// 获取日记条目
	entry, err := s.journalEntryRepo.FindByID(entryID)
	if err != nil {
		return errors.New("日记条目不存在")
	}

	// 验证权限
	if entry.UserID != userID {
		return errors.New("无权限操作此日记")
	}

	// 删除图片
	return s.journalEntryRepo.RemoveImage(entryID, imageURL)
}

// GetAllUserTags 获取用户所有标签
func (s *JournalEntryService) GetAllUserTags(userID int64) ([]string, error) {
	return s.journalEntryRepo.GetAllUserTags(userID)
}

// GetJournalsByTag 按标签获取日记
func (s *JournalEntryService) GetJournalsByTag(userID int64, tag string, page, pageSize int) ([]response.JournalEntryItem, int64, error) {
	entries, total, err := s.journalEntryRepo.FindByTag(userID, tag, page, pageSize)
	if err != nil {
		return nil, 0, errors.New("查询日记失败")
	}

	items := make([]response.JournalEntryItem, len(entries))
	for i, entry := range entries {
		var sentimentScore float64
		if entry.SentimentScore != nil {
			sentimentScore = *entry.SentimentScore
		}
		items[i] = response.JournalEntryItem{
			ID:             entry.ID,
			UserID:         entry.UserID,
			Title:          entry.Title,
			Content:        entry.Content,
			MoodRecordID:   entry.MoodRecordID,
			SentimentScore: sentimentScore,
			IsPublic:       !entry.IsPrivate,
			AnalysisStatus: entry.AnalysisStatus,
			CreatedAt:      entry.CreatedAt,
		}
	}

	return items, total, nil
}
