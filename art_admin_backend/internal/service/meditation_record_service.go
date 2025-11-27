package service

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/dify"
	"art_admin_backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// MeditationRecordService 冥想记录服务
type MeditationRecordService struct {
	meditationRecordRepo  *repository.MeditationRecordRepository
	meditationContentRepo *repository.MeditationContentRepository
	moodRecordRepo        *repository.MoodRecordRepository
	aiAnalysisService     *AIAnalysisService
}

// NewMeditationRecordService 创建冥想记录服务
func NewMeditationRecordService() *MeditationRecordService {
	return &MeditationRecordService{
		meditationRecordRepo:  repository.NewMeditationRecordRepository(),
		meditationContentRepo: repository.NewMeditationContentRepository(),
		moodRecordRepo:        repository.NewMoodRecordRepository(),
		aiAnalysisService:     NewAIAnalysisService(),
	}
}

// ValidateMeditationRecord 验证冥想记录数据
func (s *MeditationRecordService) ValidateMeditationRecord(req *request.CreateMeditationRecordRequest) error {
	// 验证冥想时长
	if req.Duration <= 0 || req.Duration > 480 { // 最长8小时
		return errors.New("冥想时长必须在1-480分钟之间")
	}

	// 验证冥想类型
	if req.MeditationType != "" && !s.isValidMeditationType(req.MeditationType) {
		return errors.New("无效的冥想类型")
	}

	// 验证难度等级
	if req.DifficultyLevel < 1 || req.DifficultyLevel > 5 {
		return errors.New("难度等级必须在1-5之间")
	}

	// 验证情绪强度
	if req.PreMoodIntensity < 1 || req.PreMoodIntensity > 10 {
		return errors.New("冥想前情绪强度必须在1-10之间")
	}

	if req.PostMoodIntensity < 1 || req.PostMoodIntensity > 10 {
		return errors.New("冥想后情绪强度必须在1-10之间")
	}

	// 验证分心程度
	if req.DistractionLevel < 1 || req.DistractionLevel > 10 {
		return errors.New("分心程度必须在1-10之间")
	}

	// 验证技巧数组
	if len(req.Techniques) > 10 {
		return errors.New("技巧数量不能超过10个")
	}

	return nil
}

// CreateMeditationRecord 创建冥想记录
func (s *MeditationRecordService) CreateMeditationRecord(userID int64, req *request.CreateMeditationRecordRequest) error {
	// 验证数据
	if err := s.ValidateMeditationRecord(req); err != nil {
		return err
	}

	// 验证关联的冥想内容是否存在
	if req.ContentID != nil && *req.ContentID > 0 {
		_, err := s.meditationContentRepo.FindByID(*req.ContentID)
		if err != nil {
			return errors.New("关联的冥想内容不存在")
		}
	}

	// 转换技巧数组为JSON字符串
	techniquesJSON := ""
	if len(req.Techniques) > 0 {
		techniquesBytes, err := json.Marshal(req.Techniques)
		if err != nil {
			return errors.New("冥想技巧序列化失败")
		}
		techniquesJSON = string(techniquesBytes)
	}

	// 创建模型
	record := &model.MeditationRecord{
		UserID:            userID,
		ContentID:         req.ContentID,
		Duration:          req.Duration,
		MeditationType:    req.MeditationType,
		DifficultyLevel:   req.DifficultyLevel,
		PreMoodType:       req.PreMoodType,
		PostMoodType:      req.PostMoodType,
		PreMoodIntensity:  req.PreMoodIntensity,
		PostMoodIntensity: req.PostMoodIntensity,
		Experience:        req.Experience,
		DistractionLevel:  req.DistractionLevel,
		Techniques:        techniquesJSON,
		Notes:             req.Notes,
		AnalysisStatus:    string(model.AnalysisStatusPending),
	}

	// 保存到数据库
	if err := s.meditationRecordRepo.Create(record); err != nil {
		return errors.New("创建冥想记录失败")
	}

	// 异步触发AI分析
	go s.triggerAsyncAnalysis(userID, record.ID)

	return nil
}

// UpdateMeditationRecord 更新冥想记录
func (s *MeditationRecordService) UpdateMeditationRecord(userID int64, req *request.UpdateMeditationRecordRequest) error {
	// 获取现有记录
	existingRecord, err := s.meditationRecordRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("冥想记录不存在")
	}

	// 验证权限
	if existingRecord.UserID != userID {
		return errors.New("无权限操作此记录")
	}

	// 验证冥想类型（如果提供）
	if req.MeditationType != nil && !s.isValidMeditationType(*req.MeditationType) {
		return errors.New("无效的冥想类型")
	}

	// 验证时长（如果提供）
	if req.Duration != nil && *req.Duration <= 0 {
		return errors.New("冥想时长必须大于0")
	}

	// 验证难度等级（如果提供）
	if req.DifficultyLevel != nil && (*req.DifficultyLevel < 1 || *req.DifficultyLevel > 5) {
		return errors.New("难度等级必须在1-5之间")
	}

	// 验证情绪强度（如果提供）
	if req.PreMoodIntensity != nil && (*req.PreMoodIntensity < 1 || *req.PreMoodIntensity > 10) {
		return errors.New("冥想前情绪强度必须在1-10之间")
	}
	if req.PostMoodIntensity != nil && (*req.PostMoodIntensity < 1 || *req.PostMoodIntensity > 10) {
		return errors.New("冥想后情绪强度必须在1-10之间")
	}

	// 验证分心程度（如果提供）
	if req.DistractionLevel != nil && (*req.DistractionLevel < 0 || *req.DistractionLevel > 10) {
		return errors.New("分心程度必须在0-10之间")
	}

	// 转换技巧数组为JSON字符串（如果提供）
	var techniquesJSON string
	if req.Techniques != nil {
		techniquesBytes, err := json.Marshal(req.Techniques)
		if err != nil {
			return errors.New("冥想技巧序列化失败")
		}
		techniquesJSON = string(techniquesBytes)
	}

	// 更新字段
	if req.ContentID != nil {
		existingRecord.ContentID = req.ContentID
	}
	if req.Duration != nil {
		existingRecord.Duration = *req.Duration
	}
	if req.MeditationType != nil {
		existingRecord.MeditationType = *req.MeditationType
	}
	if req.DifficultyLevel != nil {
		existingRecord.DifficultyLevel = *req.DifficultyLevel
	}
	if req.PreMoodType != nil {
		existingRecord.PreMoodType = *req.PreMoodType
	}
	if req.PostMoodType != nil {
		existingRecord.PostMoodType = *req.PostMoodType
	}
	if req.PreMoodIntensity != nil {
		existingRecord.PreMoodIntensity = *req.PreMoodIntensity
	}
	if req.PostMoodIntensity != nil {
		existingRecord.PostMoodIntensity = *req.PostMoodIntensity
	}
	if req.Experience != nil {
		existingRecord.Experience = *req.Experience
	}
	if req.DistractionLevel != nil {
		existingRecord.DistractionLevel = *req.DistractionLevel
	}
	if req.Techniques != nil {
		existingRecord.Techniques = techniquesJSON
	}
	if req.Notes != nil {
		existingRecord.Notes = *req.Notes
	}

	existingRecord.UpdatedAt = time.Now()

	// 更新数据库
	return s.meditationRecordRepo.Update(existingRecord)
}

// GetMeditationRecord 获取单个冥想记录
func (s *MeditationRecordService) GetMeditationRecord(userID, recordID int64) (*response.MeditationRecordItem, error) {
	// 获取记录
	record, err := s.meditationRecordRepo.FindByID(recordID)
	if err != nil {
		return nil, errors.New("冥想记录不存在")
	}

	// 验证权限
	if record.UserID != userID {
		return nil, errors.New("无权限操作此记录")
	}

	// 转换为响应格式
	item := &response.MeditationRecordItem{
		ID:                  record.ID,
		UserID:              record.UserID,
		ContentID:           record.ContentID,
		Duration:            record.Duration,
		MeditationType:      record.MeditationType,
		MeditationTypeLabel: model.GetMeditationTypeLabel(model.MeditationType(record.MeditationType)),
		DifficultyLevel:     record.DifficultyLevel,
		PreMoodType:         record.PreMoodType,
		PostMoodType:        record.PostMoodType,
		PreMoodIntensity:    record.PreMoodIntensity,
		PostMoodIntensity:   record.PostMoodIntensity,
		Experience:          record.Experience,
		DistractionLevel:    record.DistractionLevel,
		Notes:               record.Notes,
		AnalysisStatus:      record.AnalysisStatus,
		MoodImprovement:     record.MoodImprovement,
		CreatedAt:           record.CreatedAt,
	}

	// 解析技巧JSON
	if record.Techniques != "" {
		var techniques []string
		if err := json.Unmarshal([]byte(record.Techniques), &techniques); err == nil {
			item.Techniques = techniques
		}
	}

	return item, nil
}

// DeleteMeditationRecord 删除冥想记录
func (s *MeditationRecordService) DeleteMeditationRecord(userID, id int64) error {
	// 获取记录
	record, err := s.meditationRecordRepo.FindByID(id)
	if err != nil {
		return errors.New("冥想记录不存在")
	}

	// 验证权限
	if record.UserID != userID {
		return errors.New("无权限操作此记录")
	}

	// 删除记录
	return s.meditationRecordRepo.Delete(id)
}

// GetMeditationRecordList 获取冥想记录列表
func (s *MeditationRecordService) GetMeditationRecordList(userID int64, req *request.MeditationRecordListRequest) ([]response.MeditationRecordItem, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	query["user_id"] = userID
	if req.ContentID != nil {
		query["content_id"] = *req.ContentID
	}
	if req.MeditationType != "" {
		query["meditation_type"] = req.MeditationType
	}
	if req.DifficultyLevel != nil {
		query["difficulty_level"] = *req.DifficultyLevel
	}
	if req.AnalysisStatus != "" {
		query["analysis_status"] = req.AnalysisStatus
	}
	if req.MinDuration != nil {
		query["min_duration"] = *req.MinDuration
	}
	if req.MaxDuration != nil {
		query["max_duration"] = *req.MaxDuration
	}
	if req.StartDate != nil {
		query["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		query["end_date"] = *req.EndDate
	}

	// 查询记录
	records, total, err := s.meditationRecordRepo.FindWithPagination(query, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	items := make([]response.MeditationRecordItem, len(records))
	for i, record := range records {
		items[i] = s.convertToMeditationRecordItem(&record)
	}

	return items, total, nil
}

// GetMeditationRecordDetail 获取冥想记录详情
func (s *MeditationRecordService) GetMeditationRecordDetail(userID, id int64) (*response.MeditationRecordDetail, error) {
	// 获取记录
	record, err := s.meditationRecordRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 验证权限
	if record.UserID != userID {
		return nil, errors.New("无权限查看此记录")
	}

	// 转换为响应格式
	detail := &response.MeditationRecordDetail{
		MeditationRecordItem: s.convertToMeditationRecordItem(record),
		UpdatedAt:            record.UpdatedAt,
	}

	return detail, nil
}

// GetMeditationRecordsByDateRange 根据日期范围获取冥想记录
func (s *MeditationRecordService) GetMeditationRecordsByDateRange(userID int64, startDate, endDate time.Time, page, pageSize int) ([]response.MeditationRecordItem, int64, error) {
	records, total, err := s.meditationRecordRepo.FindByUserIDAndDateRange(userID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.MeditationRecordItem, len(records))
	for i, record := range records {
		items[i] = s.convertToMeditationRecordItem(&record)
	}

	return items, total, nil
}

// GetMeditationRecordsByDate 根据具体日期获取冥想记录
func (s *MeditationRecordService) GetMeditationRecordsByDate(userID int64, date time.Time) ([]response.MeditationRecordItem, error) {
	records, err := s.meditationRecordRepo.FindByUserIDAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	items := make([]response.MeditationRecordItem, len(records))
	for i, record := range records {
		items[i] = s.convertToMeditationRecordItem(&record)
	}

	return items, nil
}

// GetMeditationStatistics 获取冥想统计信息
func (s *MeditationRecordService) GetMeditationStatistics(userID int64, days int) (*response.MeditationStatistics, error) {
	// 获取基础统计
	stats, err := s.meditationRecordRepo.GetMeditationStats(userID, days)
	if err != nil {
		return nil, err
	}

	// 获取连续天数
	streak, err := s.meditationRecordRepo.GetMeditationStreak(userID)
	if err != nil {
		streak = 0
	}

	// 获取总记录数
	totalCount, err := s.meditationRecordRepo.CountByUserID(userID)
	if err != nil {
		totalCount = 0
	}

	// 分析冥想效果与情绪改善的关联性
	correlation, err := s.analyzeMoodMeditationCorrelation(userID, days)
	if err != nil {
		correlation = 0.0
	}

	// 转换为响应格式
	result := &response.MeditationStatistics{
		TotalSessions:      stats.TotalSessions,
		TotalDuration:      stats.TotalDuration,
		AverageDuration:    int(stats.AverageDuration),
		MostUsedType:       stats.MostUsedType,
		MostUsedTypeLabel:  model.GetMeditationTypeLabel(model.MeditationType(stats.MostUsedType)),
		AverageImprovement: stats.AverageImprovement,
		RecordingStreak:    streak,
		TotalCount:         int(totalCount),
		MoodCorrelation:    correlation,
		TypeDistribution:   s.parseDistribution(stats.Distribution),
		AnalysisDate:       stats.AnalyzeDate,
	}

	return result, nil
}

// AnalyzeMeditationSession 分析单次冥想效果
func (s *MeditationRecordService) AnalyzeMeditationSession(ctx context.Context, userID, recordID int64) (map[string]interface{}, error) {
	// TODO: 临时使用map代替结构体解决类型解析问题
	// 获取冥想记录
	record, err := s.meditationRecordRepo.FindByID(recordID)
	if err != nil {
		return nil, errors.New("冥想记录不存在")
	}

	// 验证权限
	if record.UserID != userID {
		return nil, errors.New("无权限分析此记录")
	}

	// 更新分析状态
	if err := s.meditationRecordRepo.UpdateAnalysisStatus(recordID, string(model.AnalysisStatusAnalyzing)); err != nil {
		return nil, fmt.Errorf("更新分析状态失败: %w", err)
	}

	// 构建分析提示词
	prompt := s.buildMeditationAnalysisPrompt(record)

	// 调用AI分析
	systemPrompt := `你是一位专业的冥想指导师和心理健康专家，擅长分析冥想练习效果。
请基于提供的冥想记录数据进行深入分析，并以JSON格式返回结构化结果。
分析结果应该包括：
1. 冥想效果评估
2. 情绪改善分析
3. 个性化建议
4. 下次练习推荐
请确保分析结果专业、准确且具有建设性。`

	aiResponse, err := s.aiAnalysisService.deepSeekClient.SimpleChat(ctx, systemPrompt, prompt)
	if err != nil {
		s.meditationRecordRepo.UpdateAnalysisStatus(recordID, string(model.AnalysisStatusFailed))
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 解析AI响应
	analysis, err := s.parseMeditationAnalysisResponse(aiResponse, record)
	if err != nil {
		s.meditationRecordRepo.UpdateAnalysisStatus(recordID, string(model.AnalysisStatusFailed))
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 存储分析结果到数据库
	analysisJSON, _ := json.Marshal(analysis)
	s.meditationRecordRepo.UpdateAnalysisResult(recordID, string(analysisJSON), "")

	// 异步存储到Dify知识库
	go s.storeToDifyKnowledge(ctx, record, analysis)

	return analysis, nil
}

// GetMeditationTrends 获取冥想趋势
func (s *MeditationRecordService) GetMeditationTrends(userID int64, days int) ([]response.MeditationTrendData, error) {
	// 获取指定天数内的记录
	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	records, _, err := s.meditationRecordRepo.FindByUserIDAndDateRange(userID, startDate, endDate, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	trendMap := make(map[string]response.MeditationTrendData)
	for _, record := range records {
		date := record.CreatedAt.Format("2006-01-02")

		if existing, exists := trendMap[date]; exists {
			// 累加数据
			existing.TotalDuration += record.Duration
			existing.SessionCount++
			existing.AverageImprovement = (existing.AverageImprovement + s.calculateImprovement(&record)) / 2
			trendMap[date] = existing
		} else {
			// 创建新数据
			trendMap[date] = response.MeditationTrendData{
				Date:               date,
				TotalDuration:      record.Duration,
				SessionCount:       1,
				AverageImprovement: s.calculateImprovement(&record),
				MostUsedType:       record.MeditationType,
			}
		}
	}

	// 转换为数组并按日期排序
	var trends []response.MeditationTrendData
	for _, trend := range trendMap {
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

// GetMeditationRecommendations 获取冥想推荐
func (s *MeditationRecordService) GetMeditationRecommendations(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	// TODO: 临时使用map代替结构体解决类型解析问题
	// 获取用户最近的情绪记录
	moodRecords, err := s.moodRecordRepo.FindByUserIDAndDateRange(userID, 7)
	if err != nil {
		return nil, err
	}

	// 获取用户最近的冥想记录
	meditationRecords, err := s.meditationRecordRepo.FindRecent(userID, 10)
	if err != nil {
		return nil, err
	}

	// 构建推荐提示词
	prompt := s.buildRecommendationPrompt(moodRecords, meditationRecords)

	systemPrompt := `你是一位专业的冥想指导师，请基于用户的情绪状态和冥想练习历史，提供个性化的冥想推荐。
返回JSON格式，包含推荐类型、原因和预期效果。`

	// 调用AI分析
	aiResponse, err := s.aiAnalysisService.deepSeekClient.SimpleChat(ctx, systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI推荐失败: %w", err)
	}

	// 解析推荐结果
	recommendations, err := s.parseRecommendationResponse(aiResponse)
	if err != nil {
		return nil, fmt.Errorf("解析推荐失败: %w", err)
	}

	return recommendations, nil
}

// triggerAsyncAnalysis 触发异步AI分析
func (s *MeditationRecordService) triggerAsyncAnalysis(userID, recordID int64) {
	// 检查防抖机制
	if !GetAnalysisDebouncer().ShouldCheck(userID) {
		return // 跳过本次分析
	}

	ctx := context.Background()

	// 执行AI分析
	_, err := s.AnalyzeMeditationSession(ctx, userID, recordID)
	if err != nil {
		fmt.Printf("异步冥想分析失败，用户ID: %d, 记录ID: %d, 错误: %v\n", userID, recordID, err)
	}
}

// isValidMeditationType 验证冥想类型是否有效
func (s *MeditationRecordService) isValidMeditationType(meditationType string) bool {
	validTypes := model.GetMeditationTypes()
	for _, validType := range validTypes {
		if string(validType) == meditationType {
			return true
		}
	}
	return false
}

// isValidJSON 验证字符串是否为有效JSON
func (s *MeditationRecordService) isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// convertToMeditationRecordItem 转换模型为响应项
func (s *MeditationRecordService) convertToMeditationRecordItem(record *model.MeditationRecord) response.MeditationRecordItem {
	// 解析技巧
	var techniques []string
	if record.Techniques != "" {
		json.Unmarshal([]byte(record.Techniques), &techniques)
	}

	item := response.MeditationRecordItem{
		ID:                  record.ID,
		UserID:              record.UserID,
		ContentID:           record.ContentID,
		Duration:            record.Duration,
		MeditationType:      record.MeditationType,
		MeditationTypeLabel: model.GetMeditationTypeLabel(model.MeditationType(record.MeditationType)),
		DifficultyLevel:     record.DifficultyLevel,
		PreMoodType:         record.PreMoodType,
		PostMoodType:        record.PostMoodType,
		PreMoodIntensity:    record.PreMoodIntensity,
		PostMoodIntensity:   record.PostMoodIntensity,
		Experience:          record.Experience,
		DistractionLevel:    record.DistractionLevel,
		Techniques:          techniques,
		Notes:               record.Notes,
		AnalysisStatus:      record.AnalysisStatus,
		CreatedAt:           record.CreatedAt,
	}

	// 计算情绪改善程度
	item.MoodImprovement = s.calculateImprovement(record)

	return item
}

// calculateImprovement 计算情绪改善程度
func (s *MeditationRecordService) calculateImprovement(record *model.MeditationRecord) float64 {
	if record.PreMoodIntensity == 0 || record.PostMoodIntensity == 0 {
		return 0
	}
	improvement := float64(record.PreMoodIntensity-record.PostMoodIntensity) / 10.0
	if improvement < 0 {
		improvement = 0 // 负改善视为0
	}
	return improvement
}

// analyzeMoodMeditationCorrelation 分析冥想与情绪改善的关联性
func (s *MeditationRecordService) analyzeMoodMeditationCorrelation(userID int64, days int) (float64, error) {
	// 获取冥想记录
	meditationRecords, _, err := s.meditationRecordRepo.FindByUserIDAndDateRange(
		userID, time.Now().AddDate(0, 0, -days), time.Now(), 1, 1000)
	if err != nil {
		return 0, err
	}

	if len(meditationRecords) == 0 {
		return 0, nil
	}

	// 计算平均改善程度
	totalImprovement := 0.0
	validRecords := 0

	for _, record := range meditationRecords {
		improvement := s.calculateImprovement(&record)
		if improvement > 0 {
			totalImprovement += improvement
			validRecords++
		}
	}

	if validRecords == 0 {
		return 0, nil
	}

	return totalImprovement / float64(validRecords), nil
}

// parseDistribution 解析分布JSON
func (s *MeditationRecordService) parseDistribution(distributionJSON string) map[string]int {
	result := make(map[string]int)
	if distributionJSON != "" {
		json.Unmarshal([]byte(distributionJSON), &result)
	}
	return result
}

// buildMeditationAnalysisPrompt 构建冥想分析提示词
func (s *MeditationRecordService) buildMeditationAnalysisPrompt(record *model.MeditationRecord) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("请分析以下冥想练习记录：\n\n"))
	prompt.WriteString(fmt.Sprintf("=== 基本信息 ===\n"))
	prompt.WriteString(fmt.Sprintf("冥想类型: %s\n", record.MeditationType))
	prompt.WriteString(fmt.Sprintf("时长: %d 分钟\n", record.Duration))
	prompt.WriteString(fmt.Sprintf("难度等级: %d\n", record.DifficultyLevel))
	prompt.WriteString(fmt.Sprintf("练习时间: %s\n\n", record.CreatedAt.Format("2006-01-02 15:04")))

	prompt.WriteString(fmt.Sprintf("=== 情绪状态 ===\n"))
	prompt.WriteString(fmt.Sprintf("冥想前情绪: %s (强度: %d)\n", record.PreMoodType, record.PreMoodIntensity))
	prompt.WriteString(fmt.Sprintf("冥想后情绪: %s (强度: %d)\n\n", record.PostMoodType, record.PostMoodIntensity))

	if record.Experience != "" {
		prompt.WriteString(fmt.Sprintf("=== 练习体验 ===\n%s\n\n", record.Experience))
	}

	prompt.WriteString(fmt.Sprintf("=== 其他信息 ===\n"))
	prompt.WriteString(fmt.Sprintf("分心程度: %d/10\n", record.DistractionLevel))

	if record.Techniques != "" {
		prompt.WriteString(fmt.Sprintf("使用技巧: %s\n", record.Techniques))
	}

	if record.Notes != "" {
		prompt.WriteString(fmt.Sprintf("备注: %s\n", record.Notes))
	}

	prompt.WriteString(fmt.Sprintf("\n请以JSON格式返回分析结果，包含以下字段：\n"))
	prompt.WriteString("{\n")
	prompt.WriteString("  \"effectiveness\": \"效果评估\",\n")
	prompt.WriteString("  \"moodImprovement\": \"情绪改善分析\",\n")
	prompt.WriteString("  \"techniquesFeedback\": \"技巧反馈\",\n")
	prompt.WriteString("  \"recommendations\": [\"建议1\", \"建议2\"],\n")
	prompt.WriteString("  \"nextPractice\": \"下次练习推荐\"\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

// parseMeditationAnalysisResponse 解析冥想分析响应
func (s *MeditationRecordService) parseMeditationAnalysisResponse(response string, record *model.MeditationRecord) (map[string]interface{}, error) {
	// 提取JSON部分
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("响应中未找到有效的JSON数据")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var aiResponse struct {
		Effectiveness      string   `json:"effectiveness"`
		MoodImprovement    string   `json:"moodImprovement"`
		TechniquesFeedback string   `json:"techniquesFeedback"`
		Recommendations    []string `json:"recommendations"`
		NextPractice       string   `json:"nextPractice"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	// 构建最终结果
	// TODO: 临时注释掉类型解析问题，后续修复
	// analysis := &response.MeditationSessionAnalysis{
	// 	RecordID:           record.ID,
	// 	Effectiveness:      aiResponse.Effectiveness,
	// 	MoodImprovement:    aiResponse.MoodImprovement,
	// 	TechniquesFeedback: aiResponse.TechniquesFeedback,
	// 	Recommendations:    aiResponse.Recommendations,
	// 	NextPractice:       aiResponse.NextPractice,
	// 	AnalyzedAt:         time.Now(),
	// }

	// 临时返回基本分析结果
	analysis := map[string]interface{}{
		"record_id":           record.ID,
		"effectiveness":       aiResponse.Effectiveness,
		"mood_improvement":    aiResponse.MoodImprovement,
		"techniques_feedback": aiResponse.TechniquesFeedback,
		"recommendations":     aiResponse.Recommendations,
		"next_practice":       aiResponse.NextPractice,
		"analyzed_at":         time.Now(),
	}

	return analysis, nil
}

// buildRecommendationPrompt 构建推荐提示词
func (s *MeditationRecordService) buildRecommendationPrompt(moodRecords []model.MoodRecord, meditationRecords []model.MeditationRecord) string {
	var prompt strings.Builder

	prompt.WriteString("基于以下用户数据，提供个性化冥想推荐：\n\n")

	if len(moodRecords) > 0 {
		prompt.WriteString("=== 最近情绪状态 ===\n")
		for _, record := range moodRecords {
			prompt.WriteString(fmt.Sprintf("%s: %s (强度: %d)\n",
				record.CreatedAt.Format("01-02"), record.MoodType, record.Intensity))
		}
		prompt.WriteString("\n")
	}

	if len(meditationRecords) > 0 {
		prompt.WriteString("=== 最近冥想练习 ===\n")
		for _, record := range meditationRecords {
			prompt.WriteString(fmt.Sprintf("%s: %s冥想 %d分钟\n",
				record.CreatedAt.Format("01-02"), record.MeditationType, record.Duration))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("请以JSON格式返回推荐，包含：\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"recommendations\": [{\n")
	prompt.WriteString("    \"type\": \"推荐类型\",\n")
	prompt.WriteString("    \"reason\": \"推荐原因\",\n")
	prompt.WriteString("    \"expectedEffect\": \"预期效果\",\n")
	prompt.WriteString("    \"duration\": \"建议时长\"\n")
	prompt.WriteString("  }]\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

// parseRecommendationResponse 解析推荐响应
func (s *MeditationRecordService) parseRecommendationResponse(response string) ([]map[string]interface{}, error) {
	// TODO: 临时使用map代替结构体解决类型解析问题
	// 提取JSON部分
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, errors.New("无效的JSON响应")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var aiResponse struct {
		Recommendations []struct {
			Type           string `json:"type"`
			Reason         string `json:"reason"`
			ExpectedEffect string `json:"expectedEffect"`
			Duration       string `json:"duration"`
		} `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	// 转换为响应格式
	result := make([]map[string]interface{}, len(aiResponse.Recommendations))
	for i, rec := range aiResponse.Recommendations {
		result[i] = map[string]interface{}{
			"type":            rec.Type,
			"reason":          rec.Reason,
			"expected_effect": rec.ExpectedEffect,
			"duration":        rec.Duration,
		}
	}

	return result, nil
}

// storeToDifyKnowledge 存储到Dify知识库
func (s *MeditationRecordService) storeToDifyKnowledge(ctx context.Context, record *model.MeditationRecord, analysis map[string]interface{}) error {
	// 安全获取分析结果字段
	effectiveness := s.safeGetString(analysis, "effectiveness")
	moodImprovement := s.safeGetString(analysis, "mood_improvement")
	techniquesFeedback := s.safeGetString(analysis, "techniques_feedback")
	recommendations := s.safeGetStringSlice(analysis, "recommendations")
	nextPractice := s.safeGetString(analysis, "next_practice")

	// 构造知识库文档内容
	documentContent := fmt.Sprintf(`# 冥想练习分析报告

## 基本信息
- 用户ID: %d
- 冥想类型: %s
- 练习时长: %d 分钟
- 练习时间: %s

## 情绪变化
- 冥想前: %s (强度: %d)
- 冥想后: %s (强度: %d)
- 改善程度: %.1f

## 练习体验
%s

## AI分析结果
### 效果评估
%s

### 情绪改善分析
%s

### 技巧反馈
%s

### 个性化建议
%s

### 下次练习推荐
%s

---
*本报告由AI系统自动生成，仅供参考*
`,
		record.UserID,
		record.MeditationType,
		record.Duration,
		record.CreatedAt.Format("2006-01-02 15:04:05"),
		record.PreMoodType, record.PreMoodIntensity,
		record.PostMoodType, record.PostMoodIntensity,
		s.calculateImprovement(record),
		record.Experience,
		effectiveness,
		moodImprovement,
		techniquesFeedback,
		s.formatRecommendations(recommendations),
		nextPractice,
	)

	// 调用Dify知识库存储API
	datasetID := config.GlobalConfig.Dify.DatasetID
	if datasetID == "" {
		fmt.Printf("Dify数据集ID未配置，跳过知识库存储\n")
		return nil
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("meditation_analysis_%d_*.md", record.UserID))
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入内容到临时文件
	if _, err := tmpFile.WriteString(documentContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("写入文件内容失败: %w", err)
	}
	tmpFile.Close()

	// 上传到Dify知识库
	difyClient := dify.NewClient(
		config.GlobalConfig.Dify.DatasetAPIKey,
		config.GlobalConfig.Dify.BaseURL,
		config.GlobalConfig.Dify.Timeout,
	)

	uploadResponse, err := difyClient.UploadFile(datasetID, tmpFile.Name())
	if err != nil {
		return fmt.Errorf("上传到Dify知识库失败: %w", err)
	}

	// 更新记录的Dify文档ID
	if uploadResponse != nil && uploadResponse.Document.ID != "" {
		s.meditationRecordRepo.UpdateAnalysisResult(record.ID, "", uploadResponse.Document.ID)
		fmt.Printf("成功上传冥想分析到Dify知识库 - 数据集: %s, 文档ID: %s\n", datasetID, uploadResponse.Document.ID)
	}

	return nil
}

// safeGetString 安全获取字符串值
func (s *MeditationRecordService) safeGetString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

// safeGetStringSlice 安全获取字符串切片
func (s *MeditationRecordService) safeGetStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key]; ok {
		if slice, ok := v.([]string); ok {
			return slice
		}
		// 处理 []interface{} 类型
		if slice, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

// formatRecommendations 格式化建议
func (s *MeditationRecordService) formatRecommendations(recommendations []string) string {
	if len(recommendations) == 0 {
		return "暂无特别建议"
	}

	var result strings.Builder
	for _, rec := range recommendations {
		result.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	return result.String()
}
