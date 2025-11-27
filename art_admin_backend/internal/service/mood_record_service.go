package service

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type MoodRecordService struct {
	moodRecordRepo    *repository.MoodRecordRepository
	aiAnalysisService *AIAnalysisService
}

func NewMoodRecordService() *MoodRecordService {
	return &MoodRecordService{
		moodRecordRepo:    repository.NewMoodRecordRepository(),
		aiAnalysisService: NewAIAnalysisService(),
	}
}

// CreateMoodRecord 创建情绪记录
func (s *MoodRecordService) CreateMoodRecord(userID int64, req *request.CreateMoodRecordRequest) error {
	// 验证情绪类型
	if !utils.ValidateMoodType(req.MoodType) {
		return errors.New("无效的情绪类型")
	}

	// 验证强度
	if req.Intensity < 1 || req.Intensity > 10 {
		return errors.New("情绪强度必须在1-10之间")
	}

	// 创建模型
	record := &model.MoodRecord{
		UserID:    userID,
		MoodType:  req.MoodType,
		Intensity: req.Intensity,
		Note:      req.Note,
		Location:  req.Location,
	}

	// 处理JSON字段 - 只有当字段不为空时才设置
	if len(req.Triggers) > 0 {
		triggersJSON, err := utils.JSONMarshal(req.Triggers)
		if err != nil {
			return errors.New("触发因素序列化失败")
		}
		record.Triggers = triggersJSON
	} else {
		// 设置空JSON数组而不是空字符串
		record.Triggers = "[]"
	}

	if len(req.Activities) > 0 {
		activitiesJSON, err := utils.JSONMarshal(req.Activities)
		if err != nil {
			return errors.New("活动序列化失败")
		}
		record.Activities = activitiesJSON
	} else {
		// 设置空JSON数组而不是空字符串
		record.Activities = "[]"
	}

	// 保存到数据库
	if err := s.moodRecordRepo.Create(record); err != nil {
		fmt.Printf("创建情绪记录数据库错误: %v\n", err)
		return errors.New("创建情绪记录失败")
	}

	// 异步触发AI分析
	go s.triggerAsyncAnalysis(userID)

	return nil
}

// UpdateMoodRecord 更新情绪记录
func (s *MoodRecordService) UpdateMoodRecord(userID int64, req *request.UpdateMoodRecordRequest) error {
	// 获取现有记录
	existingRecord, err := s.moodRecordRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("记录不存在")
	}

	// 验证所有权
	if existingRecord.UserID != userID {
		return errors.New("无权限操作此记录")
	}

	// 验证情绪类型（如果提供）
	if req.MoodType != "" && !utils.ValidateMoodType(req.MoodType) {
		return errors.New("无效的情绪类型")
	}

	// 验证情绪强度（如果提供）
	if req.Intensity != nil {
		if *req.Intensity < 1 || *req.Intensity > 10 {
			return errors.New("情绪强度必须在1-10之间")
		}
	}

	// 更新字段
	if req.MoodType != "" {
		existingRecord.MoodType = req.MoodType
	}
	if req.Intensity != nil {
		existingRecord.Intensity = *req.Intensity
	}
	if req.Triggers != nil {
		triggersJSON, err := utils.JSONMarshal(req.Triggers)
		if err != nil {
			return errors.New("触发因素序列化失败: " + err.Error())
		}
		existingRecord.Triggers = triggersJSON
	}
	if req.Activities != nil {
		activitiesJSON, err := utils.JSONMarshal(req.Activities)
		if err != nil {
			return errors.New("活动序列化失败: " + err.Error())
		}
		existingRecord.Activities = activitiesJSON
	}
	if req.Note != "" {
		existingRecord.Note = req.Note
	}
	if req.Location != "" {
		existingRecord.Location = req.Location
	}

	existingRecord.UpdatedAt = time.Now()

	return s.moodRecordRepo.Update(existingRecord)
}

// GetMoodRecord 获取单个情绪记录
func (s *MoodRecordService) GetMoodRecord(userID, recordID int64) (*response.MoodRecordItem, error) {
	// 获取记录
	record, err := s.moodRecordRepo.FindByID(recordID)
	if err != nil {
		return nil, errors.New("情绪记录不存在")
	}

	// 验证权限
	if record.UserID != userID {
		return nil, errors.New("无权限操作此记录")
	}

	// 转换为响应格式
	item := &response.MoodRecordItem{
		ID:             record.ID,
		UserID:         record.UserID,
		MoodType:       record.MoodType,
		MoodTypeLabel:  model.GetMoodTypeLabel(model.MoodType(record.MoodType)),
		Intensity:      record.Intensity,
		Triggers:       []string{},
		Activities:     []string{},
		Note:           record.Note,
		Location:       record.Location,
		AnalysisStatus: record.AnalysisStatus,
		DifyDocumentID: record.DifyDocumentID,
		CreatedAt:      record.CreatedAt,
	}

	// 解析触发因素JSON
	if record.Triggers != "" {
		triggers, err := utils.JSONUnmarshal(record.Triggers)
		if err == nil {
			item.Triggers = triggers
		}
	}

	// 解析活动JSON
	if record.Activities != "" {
		activities, err := utils.JSONUnmarshal(record.Activities)
		if err == nil {
			item.Activities = activities
		}
	}

	return item, nil
}

// GetMoodRecordsByDate 按日期获取情绪记录
func (s *MoodRecordService) GetMoodRecordsByDate(userID int64, dateStr string) ([]response.MoodRecordItem, error) {
	records, err := s.moodRecordRepo.FindByUserIDAndDateString(userID, dateStr)
	if err != nil {
		return nil, errors.New("查询情绪记录失败")
	}

	items := make([]response.MoodRecordItem, len(records))
	for i, record := range records {
		items[i] = response.MoodRecordItem{
			ID:             record.ID,
			UserID:         record.UserID,
			MoodType:       record.MoodType,
			MoodLabel:      model.GetMoodTypeLabel(model.MoodType(record.MoodType)),
			Intensity:      record.Intensity,
			Triggers:       []string{},
			Activities:     []string{},
			Note:           record.Note,
			Location:       record.Location,
			AnalysisStatus: record.AnalysisStatus,
			DifyDocumentID: record.DifyDocumentID,
			CreatedAt:      record.CreatedAt,
		}

		// 解析触发因素JSON
		if record.Triggers != "" {
			triggers, err := utils.JSONUnmarshal(record.Triggers)
			if err == nil {
				items[i].Triggers = triggers
			}
		}

		// 解析活动JSON
		if record.Activities != "" {
			activities, err := utils.JSONUnmarshal(record.Activities)
			if err == nil {
				items[i].Activities = activities
			}
		}
	}

	return items, nil
}

// GetMoodCalendar 获取情绪日历数据
func (s *MoodRecordService) GetMoodCalendar(userID int64, year, month int) ([]map[string]interface{}, error) {
	// 验证参数
	if year < 2000 || year > 2100 {
		return nil, errors.New("无效的年份")
	}
	if month < 1 || month > 12 {
		return nil, errors.New("无效的月份")
	}

	calendarData, err := s.moodRecordRepo.GetCalendarData(userID, year, month)
	if err != nil {
		return nil, errors.New("获取日历数据失败")
	}

	return calendarData, nil
}

// DeleteMoodRecord 删除情绪记录
func (s *MoodRecordService) DeleteMoodRecord(userID, id int64) error {
	// 获取记录
	record, err := s.moodRecordRepo.FindByID(id)
	if err != nil {
		return errors.New("记录不存在")
	}

	// 验证所有权
	if record.UserID != userID {
		return errors.New("无权限操作此记录")
	}

	// 删除记录
	return s.moodRecordRepo.Delete(id)
}

// GetMoodRecordList 获取情绪记录列表
func (s *MoodRecordService) GetMoodRecordList(userID int64, req *request.MoodRecordListRequest) ([]response.MoodRecordItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	query["user_id"] = userID // 强制使用认证用户的ID
	if req.MoodType != "" {
		query["mood_type"] = req.MoodType
	}
	if req.StartDate != nil {
		query["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		query["end_date"] = *req.EndDate
	}
	if req.MinIntensity != nil {
		query["min_intensity"] = *req.MinIntensity
	}
	if req.MaxIntensity != nil {
		query["max_intensity"] = *req.MaxIntensity
	}

	// 查询记录
	records, total, err := s.moodRecordRepo.FindWithPagination(query, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]response.MoodRecordItem, len(records))
	for i, record := range records {
		items[i] = s.convertToMoodRecordItem(&record)
	}

	return items, total, nil
}

// GetMoodStatistics 获取情绪统计
func (s *MoodRecordService) GetMoodStatistics(userID int64, days int) (*response.MoodStatistics, error) {
	// 获取指定天数内的记录
	records, err := s.moodRecordRepo.FindByUserIDAndDateRange(userID, days)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &response.MoodStatistics{
			TotalRecords:     0,
			MoodDistribution: make(map[string]int),
		}, nil
	}

	// 计算统计数据
	stats := &response.MoodStatistics{
		TotalRecords:     len(records),
		MoodDistribution: make(map[string]int),
	}

	totalIntensity := 0
	moodCount := make(map[string]int)

	for _, record := range records {
		totalIntensity += record.Intensity
		moodCount[record.MoodType]++
	}

	// 计算平均强度
	stats.AverageIntensity = float64(totalIntensity) / float64(len(records))

	// 找出最常见的情绪
	maxCount := 0
	for moodType, count := range moodCount {
		stats.MoodDistribution[moodType] = count
		if count > maxCount {
			maxCount = count
			stats.MostCommonMood = moodType
		}
	}
	stats.MostCommonMoodLabel = model.GetMoodTypeLabel(model.MoodType(stats.MostCommonMood))

	// 获取最近趋势（最近7天）
	recentRecords, _ := s.moodRecordRepo.FindByUserIDAndDateRange(userID, 7)
	stats.RecentTrend = make([]response.MoodRecordItem, len(recentRecords))
	for i, record := range recentRecords {
		stats.RecentTrend[i] = s.convertToMoodRecordItem(&record)
	}

	return stats, nil
}

// GetMoodAnalytics 获取情绪分析
func (s *MoodRecordService) GetMoodAnalytics(userID int64, req *request.MoodAnalyticsRequest) (*response.MoodAnalyticsItem, error) {
	// 这里可以实现更复杂的分析逻辑
	// 暂时返回基本统计信息
	stats, err := s.GetMoodStatistics(userID, 30) // 默认30天
	if err != nil {
		return nil, err
	}

	analytics := &response.MoodAnalyticsItem{
		UserID:              userID,
		AnalyzeType:         req.AnalyzeType,
		AnalyzeDate:         time.Now(),
		MoodDistribution:    stats.MoodDistribution,
		AverageIntensity:    stats.AverageIntensity,
		TotalRecords:        stats.TotalRecords,
		MostCommonMood:      stats.MostCommonMood,
		MostCommonMoodLabel: stats.MostCommonMoodLabel,
		CreatedAt:           time.Now(),
	}

	return analytics, nil
}

// convertToMoodRecordItem 转换模型为响应项
func (s *MoodRecordService) convertToMoodRecordItem(record *model.MoodRecord) response.MoodRecordItem {
	item := response.MoodRecordItem{
		ID:        record.ID,
		UserID:    record.UserID,
		MoodType:  record.MoodType,
		MoodLabel: model.GetMoodTypeLabel(model.MoodType(record.MoodType)),
		Intensity: record.Intensity,
		Note:      record.Note,
		Location:  record.Location,
		CreatedAt: record.CreatedAt,
	}

	// 解析JSON字段
	if record.Triggers != "" {
		triggers, _ := utils.JSONUnmarshal(record.Triggers)
		item.Triggers = triggers
	}

	if record.Activities != "" {
		activities, _ := utils.JSONUnmarshal(record.Activities)
		item.Activities = activities
	}

	return item
}

// AnalyzeMoodPatterns 分析情绪模式
func (s *MoodRecordService) AnalyzeMoodPatterns(ctx context.Context, userID int64, days int) (*response.MoodPatternAnalysis, error) {
	// 构建分析请求
	req := &EmotionAnalysisRequest{
		UserID:       userID,
		StartDate:    time.Now().AddDate(0, 0, -days),
		EndDate:      time.Now(),
		AnalysisType: "pattern",
	}

	// 调用AI分析服务
	result, err := s.aiAnalysisService.AnalyzeUserEmotions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 转换为响应格式
	analysis := &response.MoodPatternAnalysis{
		UserID:          result.UserID,
		AnalysisType:    result.AnalysisType,
		StartDate:       result.StartDate,
		EndDate:         result.EndDate,
		Summary:         result.Summary,
		Patterns:        s.convertEmotionPatterns(result.Patterns),
		Trends:          s.convertEmotionTrends(result.Trends),
		Insights:        result.Insights,
		Recommendations: result.Recommendations,
		AnalyzedAt:      result.AnalyzedAt,
		DataPoints:      result.DataPoints,
	}

	return analysis, nil
}

// GetMoodTrends 获取情绪趋势
func (s *MoodRecordService) GetMoodTrends(userID int64, days int) ([]response.MoodTrendData, error) {
	// 获取指定天数内的记录
	records, err := s.moodRecordRepo.FindByUserIDAndDateRange(userID, days)
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	trendMap := make(map[string][]response.MoodTrendData)
	for _, record := range records {
		date := record.CreatedAt.Format("2006-01-02")
		trend := response.MoodTrendData{
			Date:      date,
			MoodType:  record.MoodType,
			Intensity: record.Intensity,
		}
		trendMap[date] = append(trendMap[date], trend)
	}

	// 转换为数组格式
	var trends []response.MoodTrendData
	for date, dayTrends := range trendMap {
		// 计算当天的平均强度
		totalIntensity := 0
		for _, trend := range dayTrends {
			totalIntensity += trend.Intensity
		}
		avgIntensity := float64(totalIntensity) / float64(len(dayTrends))

		// 使用当天最主要的情绪类型
		mainMood := dayTrends[0].MoodType

		trends = append(trends, response.MoodTrendData{
			Date:      date,
			MoodType:  mainMood,
			Intensity: int(avgIntensity + 0.5), // 四舍五入
		})
	}

	return trends, nil
}

// triggerAsyncAnalysis 触发异步AI分析
func (s *MoodRecordService) triggerAsyncAnalysis(userID int64) {
	// 检查防抖机制
	if !GetAnalysisDebouncer().ShouldCheck(userID) {
		return // 跳过本次分析
	}

	ctx := context.Background()

	// 执行AI分析
	_, err := s.AnalyzeMoodPatterns(ctx, userID, 7) // 分析最近7天
	if err != nil {
		fmt.Printf("异步AI分析失败，用户ID: %d, 错误: %v\n", userID, err)
	}
}

// convertEmotionPatterns 转换情绪模式
func (s *MoodRecordService) convertEmotionPatterns(patterns []EmotionPattern) []response.EmotionPattern {
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
func (s *MoodRecordService) convertEmotionTrends(trends []EmotionTrend) []response.EmotionTrend {
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
