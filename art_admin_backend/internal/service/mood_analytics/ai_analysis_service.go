package mood_analytics

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/deepseek"
	"art_admin_backend/internal/pkg/dify"
	"art_admin_backend/internal/pkg/utils"
	"art_admin_backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// AIAnalysisService AI分析服务
type AIAnalysisService struct {
	deepSeekClient   *deepseek.Client
	difyKnowledge    *dify.KnowledgeClient
	moodRecordRepo   *repository.MoodRecordRepository
	journalEntryRepo *repository.JournalEntryRepository
}

// NewAIAnalysisService 创建AI分析服务
func NewAIAnalysisService() *AIAnalysisService {
	// 从配置获取DeepSeek客户端配置
	deepSeekConfig := deepseek.Config{
		APIKey:      config.GlobalConfig.DeepSeek.APIKey,
		BaseURL:     config.GlobalConfig.DeepSeek.BaseURL,
		Model:       config.GlobalConfig.DeepSeek.Model,
		Timeout:     config.GlobalConfig.DeepSeek.Timeout,
		MaxTokens:   config.GlobalConfig.DeepSeek.MaxTokens,
		Temperature: config.GlobalConfig.DeepSeek.Temperature,
	}

	// 从配置获取Dify知识库客户端配置
	difyConfig := dify.KnowledgeConfig{
		APIKey:  config.GlobalConfig.Dify.DatasetAPIKey,
		BaseURL: config.GlobalConfig.Dify.BaseURL,
		Timeout: config.GlobalConfig.Dify.Timeout,
	}

	return &AIAnalysisService{
		deepSeekClient:   deepseek.NewClient(deepSeekConfig),
		difyKnowledge:    dify.NewKnowledgeClient(difyConfig),
		moodRecordRepo:   repository.NewMoodRecordRepository(),
		journalEntryRepo: repository.NewJournalEntryRepository(),
	}
}

// EmotionAnalysisRequest 情绪分析请求
type EmotionAnalysisRequest struct {
	UserID       int64     `json:"user_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	AnalysisType string    `json:"analysis_type"` // "pattern", "trend", "insight"
}

// EmotionAnalysisResult 情绪分析结果
type EmotionAnalysisResult struct {
	UserID          int64            `json:"user_id"`
	AnalysisType    string           `json:"analysis_type"`
	StartDate       time.Time        `json:"start_date"`
	EndDate         time.Time        `json:"end_date"`
	Summary         string           `json:"summary"`
	Patterns        []EmotionPattern `json:"patterns"`
	Trends          []EmotionTrend   `json:"trends"`
	Insights        []string         `json:"insights"`
	Recommendations []string         `json:"recommendations"`
	AnalyzedAt      time.Time        `json:"analyzed_at"`
	DataPoints      int              `json:"data_points"`
}

// EmotionPattern 情绪模式
type EmotionPattern struct {
	Type        string  `json:"type"`
	Frequency   int     `json:"frequency"`
	Percentage  float64 `json:"percentage"`
	Description string  `json:"description"`
}

// EmotionTrend 情绪趋势
type EmotionTrend struct {
	Period     string  `json:"period"`
	Direction  string  `json:"direction"` // "improving", "declining", "stable"
	ChangeRate float64 `json:"change_rate"`
	Confidence float64 `json:"confidence"`
}

// AnalyzeUserEmotions 分析用户情绪模式
func (s *AIAnalysisService) AnalyzeUserEmotions(ctx context.Context, req *EmotionAnalysisRequest) (*EmotionAnalysisResult, error) {
	// 获取用户的情绪记录和日记数据
	moodRecords, err := s.moodRecordRepo.FindByUserIDAndDateRange(req.UserID, int(req.EndDate.Sub(req.StartDate).Hours()/24))
	if err != nil {
		return nil, fmt.Errorf("获取情绪记录失败: %w", err)
	}

	journalEntries, _, err := s.journalEntryRepo.FindByDateRange(req.UserID, req.StartDate, req.EndDate, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("获取日记条目失败: %w", err)
	}

	if len(moodRecords) == 0 && len(journalEntries) == 0 {
		return nil, fmt.Errorf("指定时间范围内没有数据可供分析")
	}

	// 构造分析提示词
	prompt := s.buildAnalysisPrompt(moodRecords, journalEntries, req.AnalysisType)

	// 调用DeepSeek进行分析
	systemPrompt := `你是一位专业的心理健康分析师，擅长分析情绪模式和提供个性化建议。
请基于提供的数据进行深入的情绪分析，并以JSON格式返回结构化结果。
分析结果应该包括：
1. 情绪模式识别
2. 情绪趋势分析
3. 深度洞察
4. 个性化建议
请确保分析结果专业、准确且具有建设性。`

	response, err := s.deepSeekClient.SimpleChat(ctx, systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 解析AI响应
	result, err := s.parseAIResponse(response, req)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 存储到Dify知识库
	if err := s.storeToDifyKnowledge(ctx, result); err != nil {
		// 记录错误但不影响返回结果
		fmt.Printf("存储到Dify知识库失败: %v\n", err)
	}

	return result, nil
}

// buildAnalysisPrompt 构建分析提示词
func (s *AIAnalysisService) buildAnalysisPrompt(moodRecords []model.MoodRecord, journalEntries []model.JournalEntry, analysisType string) string {
	var prompt strings.Builder

	// 确定时间范围
	var startDate, endDate string
	if len(moodRecords) > 0 {
		startDate = moodRecords[0].CreatedAt.Format("2006-01-02")
		endDate = moodRecords[len(moodRecords)-1].CreatedAt.Format("2006-01-02")
	} else if len(journalEntries) > 0 {
		startDate = journalEntries[0].CreatedAt.Format("2006-01-02")
		endDate = journalEntries[len(journalEntries)-1].CreatedAt.Format("2006-01-02")
	} else {
		startDate = "未知"
		endDate = "未知"
	}

	prompt.WriteString(fmt.Sprintf("请分析以下用户在%s到%s期间的情绪数据：\n\n", startDate, endDate))

	// 添加情绪记录数据
	if len(moodRecords) > 0 {
		prompt.WriteString("=== 情绪记录 ===\n")
		for _, record := range moodRecords {
			prompt.WriteString(fmt.Sprintf("日期: %s, 情绪: %s, 强度: %d",
				record.CreatedAt.Format("2006-01-02 15:04"),
				record.MoodType,
				record.Intensity))

			if record.Note != "" {
				prompt.WriteString(fmt.Sprintf(", 备注: %s", record.Note))
			}

			if record.Triggers != "" {
				triggers, _ := utils.JSONUnmarshal(record.Triggers)
				prompt.WriteString(fmt.Sprintf(", 触发因素: %v", triggers))
			}

			prompt.WriteString("\n")
		}
		prompt.WriteString("\n")
	}

	// 添加日记数据
	if len(journalEntries) > 0 {
		prompt.WriteString("=== 日记内容 ===\n")
		for _, entry := range journalEntries {
			prompt.WriteString(fmt.Sprintf("日期: %s, 标题: %s\n内容: %s\n\n",
				entry.CreatedAt.Format("2006-01-02 15:04"),
				entry.Title,
				entry.Content))
		}
	}

	// 添加分析要求
	prompt.WriteString(fmt.Sprintf("\n=== 分析要求 ===\n分析类型: %s\n", analysisType))
	prompt.WriteString("请以JSON格式返回分析结果，包含以下字段：\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"summary\": \"整体情绪状况总结\",\n")
	prompt.WriteString("  \"patterns\": [{\"type\": \"情绪类型\", \"frequency\": 频率, \"percentage\": 百分比, \"description\": \"模式描述\"}],\n")
	prompt.WriteString("  \"trends\": [{\"period\": \"时间段\", \"direction\": \"趋势方向\", \"changeRate\": 变化率, \"confidence\": 置信度}],\n")
	prompt.WriteString("  \"insights\": [\"洞察1\", \"洞察2\", \"洞察3\"],\n")
	prompt.WriteString("  \"recommendations\": [\"建议1\", \"建议2\", \"建议3\"]\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

// parseAIResponse 解析AI响应
func (s *AIAnalysisService) parseAIResponse(response string, req *EmotionAnalysisRequest) (*EmotionAnalysisResult, error) {
	// 尝试提取JSON部分
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("响应中未找到有效的JSON数据")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var aiResponse struct {
		Summary         string           `json:"summary"`
		Patterns        []EmotionPattern `json:"patterns"`
		Trends          []EmotionTrend   `json:"trends"`
		Insights        []string         `json:"insights"`
		Recommendations []string         `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	// 构建最终结果
	result := &EmotionAnalysisResult{
		UserID:          req.UserID,
		AnalysisType:    req.AnalysisType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Summary:         aiResponse.Summary,
		Patterns:        aiResponse.Patterns,
		Trends:          aiResponse.Trends,
		Insights:        aiResponse.Insights,
		Recommendations: aiResponse.Recommendations,
		AnalyzedAt:      time.Now(),
		DataPoints:      0, // 将在调用处设置
	}

	return result, nil
}

// storeToDifyKnowledge 存储到Dify知识库
func (s *AIAnalysisService) storeToDifyKnowledge(ctx context.Context, result *EmotionAnalysisResult) error {
	// 构造知识库文档内容
	documentContent := fmt.Sprintf(`# 用户情绪分析报告

## 基本信息
- 用户ID: %d
- 分析类型: %s
- 分析时间范围: %s 至 %s
- 分析生成时间: %s
- 数据点数量: %d

## 整体总结
%s

## 情绪模式识别
%s

## 情绪趋势分析
%s

## 深度洞察
%s

## 个性化建议
%s

---
*本报告由AI系统自动生成，仅供参考*
`,
		result.UserID,
		result.AnalysisType,
		result.StartDate.Format("2006-01-02"),
		result.EndDate.Format("2006-01-02"),
		result.AnalyzedAt.Format("2006-01-02 15:04:05"),
		result.DataPoints,
		result.Summary,
		s.formatPatterns(result.Patterns),
		s.formatTrends(result.Trends),
		s.formatInsights(result.Insights),
		s.formatRecommendations(result.Recommendations),
	)

	// 构造文档标题
	documentTitle := fmt.Sprintf("情绪分析报告_用户%d_%s_%s",
		result.UserID,
		result.StartDate.Format("20060102"),
		result.EndDate.Format("20060102"))

	// 上传到Dify知识库
	datasetID := config.GlobalConfig.Dify.DatasetID
	if datasetID == "" {
		return fmt.Errorf("Dify数据集ID未配置")
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("emotion_analysis_%d_*.md", result.UserID))
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // 确保临时文件被删除

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

	fmt.Printf("成功上传到Dify知识库 - 数据集: %s, 文档ID: %s, 标题: %s\n",
		datasetID, uploadResponse.Document.ID, documentTitle)

	return nil
}

// formatPatterns 格式化情绪模式
func (s *AIAnalysisService) formatPatterns(patterns []EmotionPattern) string {
	if len(patterns) == 0 {
		return "暂无明显模式"
	}

	var result strings.Builder
	for _, pattern := range patterns {
		result.WriteString(fmt.Sprintf("- **%s**: 出现频率 %d 次 (%.1f%%), %s\n",
			pattern.Type, pattern.Frequency, pattern.Percentage, pattern.Description))
	}
	return result.String()
}

// formatTrends 格式化情绪趋势
func (s *AIAnalysisService) formatTrends(trends []EmotionTrend) string {
	if len(trends) == 0 {
		return "暂无明显趋势"
	}

	var result strings.Builder
	for _, trend := range trends {
		result.WriteString(fmt.Sprintf("- **%s**: 趋势 %s, 变化率 %.1f%%, 置信度 %.1f%%\n",
			trend.Period, trend.Direction, trend.ChangeRate*100, trend.Confidence*100))
	}
	return result.String()
}

// formatInsights 格式化洞察
func (s *AIAnalysisService) formatInsights(insights []string) string {
	if len(insights) == 0 {
		return "暂无特别洞察"
	}

	var result strings.Builder
	for _, insight := range insights {
		result.WriteString(fmt.Sprintf("- %s\n", insight))
	}
	return result.String()
}

// formatRecommendations 格式化建议
func (s *AIAnalysisService) formatRecommendations(recommendations []string) string {
	if len(recommendations) == 0 {
		return "暂无特别建议"
	}

	var result strings.Builder
	for _, rec := range recommendations {
		result.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	return result.String()
}
