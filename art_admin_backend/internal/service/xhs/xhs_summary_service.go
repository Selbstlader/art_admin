package xhs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/baiduocr"
	"art_admin_backend/internal/pkg/config"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/pkg/xhsfetcher"
	"art_admin_backend/internal/repository"
)

// SummaryService 小红书内容总结服务
type SummaryService struct {
	volcClient  *volcengine.Client
	ocrClient   *baiduocr.Client
	fetcher     *xhsfetcher.Fetcher
	summaryRepo *repository.XHSSummaryRepository
	folderRepo  *repository.XHSFolderRepository
}

// NewSummaryService 创建总结服务
func NewSummaryService() *SummaryService {
	cfg := config.GlobalConfig

	volcClient := volcengine.NewClient(volcengine.Config{
		APIKey:      cfg.VolcEngine.APIKey,
		BaseURL:     cfg.VolcEngine.BaseURL,
		Model:       cfg.VolcEngine.Model,
		Timeout:     cfg.VolcEngine.Timeout,
		MaxTokens:   cfg.VolcEngine.MaxTokens,
		Temperature: cfg.VolcEngine.Temperature,
	})

	ocrClient := baiduocr.NewClient(baiduocr.Config{
		AppID:     cfg.BaiduOCR.AppID,
		APIKey:    cfg.BaiduOCR.APIKey,
		SecretKey: cfg.BaiduOCR.SecretKey,
		Timeout:   cfg.BaiduOCR.Timeout,
	})

	return &SummaryService{
		volcClient:  volcClient,
		ocrClient:   ocrClient,
		fetcher:     xhsfetcher.NewFetcher(),
		summaryRepo: repository.NewXHSSummaryRepository(),
		folderRepo:  repository.NewXHSFolderRepository(),
	}
}

// SummaryRequest 总结请求
type SummaryRequest struct {
	URL            string   `json:"url"`             // 小红书链接
	Content        string   `json:"content"`         // 笔记文字内容（可选，如果无法解析链接则需手动提供）
	ImageURLs      []string `json:"image_urls"`      // 笔记图片URL列表
	Style          string   `json:"style"`           // 总结风格: concise, detailed, casual
	MaxLength      int      `json:"max_length"`      // 最大字数
	EnableOCR      bool     `json:"enable_ocr"`      // 是否开启图片OCR
	EnableAnalysis bool     `json:"enable_analysis"` // 是否开启专业分析
}

// SummaryResponse 总结响应
type SummaryResponse struct {
	Title           string                `json:"title"`                      // 总结标题
	Summary         string                `json:"summary"`                    // 核心观点概述
	KeyPoints       []string              `json:"key_points"`                 // 关键信息列表
	Highlights      []string              `json:"highlights,omitempty"`       // 亮点/优点
	Concerns        []string              `json:"concerns,omitempty"`         // 注意事项/缺点
	Tags            []string              `json:"tags"`                       // 内容标签
	Sentiment       string                `json:"sentiment"`                  // 情感倾向: positive, neutral, negative
	ReadingTime     string                `json:"reading_time,omitempty"`     // 预计阅读时间
	TargetAudience  string                `json:"target_audience,omitempty"`  // 目标受众
	ContentType     string                `json:"content_type,omitempty"`     // 内容类型
	ActionItems     []string              `json:"action_items,omitempty"`     // 可执行建议
	RelatedTopics   []string              `json:"related_topics,omitempty"`   // 相关话题
	CredibilityNote string                `json:"credibility_note,omitempty"` // 可信度说明
	QuickFacts      map[string]string     `json:"quick_facts,omitempty"`      // 快速事实
	OCRTexts        []string              `json:"ocr_texts,omitempty"`        // OCR识别的文字
	Analysis        *ProfessionalAnalysis `json:"analysis,omitempty"`         // 专业分析
	OriginalURL     string                `json:"original_url,omitempty"`     // 原始链接
}

// ProfessionalAnalysis 专业分析结果
type ProfessionalAnalysis struct {
	Category        string            `json:"category"`        // 内容品类: beauty, digital, food, travel等
	MainConclusion  string            `json:"main_conclusion"` // 核心结论
	DetailedData    map[string]string `json:"detailed_data"`   // 详细数据
	RiskWarnings    []string          `json:"risk_warnings"`   // 风险提示
	Recommendations []string          `json:"recommendations"` // 建议
}

// Summarize 生成总结
func (s *SummaryService) Summarize(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
	// 0. 解析分享文本，提取标题和URL
	var noteTitle, noteAuthor, noteURL string
	if req.URL != "" {
		noteTitle, noteAuthor, noteURL = s.fetcher.ParseShareText(req.URL)
		if noteURL != "" {
			req.URL = noteURL // 使用提取的干净URL
		}
	}

	// 1. 如果有URL，尝试抓取笔记内容
	var fetchedContent string
	var fetchedImages []string
	if req.URL != "" && req.Content == "" {
		fmt.Printf("正在抓取小红书内容: %s\n", req.URL)

		// 设置抓取超时
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		note, err := s.fetcher.FetchNote(req.URL)
		if err != nil {
			fmt.Printf("抓取失败: %v\n", err)
		} else if note != nil {
			// 构建抓取到的内容 - 简化处理
			var contentParts []string
			if note.Title != "" {
				contentParts = append(contentParts, fmt.Sprintf("【标题】%s", note.Title))
				if noteTitle == "" {
					noteTitle = note.Title
				}
			}
			if note.Author != "" {
				contentParts = append(contentParts, fmt.Sprintf("【作者】%s", note.Author))
				if noteAuthor == "" {
					noteAuthor = note.Author
				}
			}
			if note.Content != "" {
				// 限制内容长度，避免过长
				content := note.Content
				if len(content) > 2000 {
					content = content[:2000] + "..."
				}
				contentParts = append(contentParts, fmt.Sprintf("【正文】\n%s", content))
			}
			// 只保留最重要的信息
			if len(contentParts) > 0 {
				fetchedContent = strings.Join(contentParts, "\n")
				fmt.Printf("成功抓取内容，长度: %d\n", len(fetchedContent))
			}
			// 只取前3张图片，减少OCR时间
			if len(note.Images) > 3 {
				fetchedImages = note.Images[:3]
			} else {
				fetchedImages = note.Images
			}
		}
	}

	// 2. 处理图片OCR
	var ocrTexts []string
	allImages := append(req.ImageURLs, fetchedImages...)
	if req.EnableOCR && len(allImages) > 0 {
		for _, imgURL := range allImages {
			text, err := s.ocrClient.ExtractTextFromURL(ctx, imgURL)
			if err != nil {
				fmt.Printf("OCR识别失败 [%s]: %v\n", imgURL, err)
				continue
			}
			if text != "" {
				ocrTexts = append(ocrTexts, text)
			}
		}
	}

	// 3. 合并内容
	var fullContent string

	// 优先使用抓取的内容
	if fetchedContent != "" {
		fullContent = fetchedContent
	} else if req.Content != "" {
		fullContent = req.Content
	}

	// 如果还是没有内容，使用分享文本中的标题
	if fullContent == "" && noteTitle != "" {
		fullContent = fmt.Sprintf("【笔记标题】%s\n【作者】%s\n\n请根据标题推测并分析可能的笔记内容。", noteTitle, noteAuthor)
	}

	// 添加OCR文本
	if len(ocrTexts) > 0 {
		fullContent += "\n\n【图片中的文字内容】\n" + strings.Join(ocrTexts, "\n---\n")
	}

	// 添加原始URL供参考
	if req.URL != "" && fullContent != "" {
		fullContent = fmt.Sprintf("【笔记链接】%s\n\n%s", req.URL, fullContent)
	}

	if fullContent == "" {
		return nil, fmt.Errorf("没有可总结的内容，请提供笔记链接、内容或图片")
	}

	// 3. 设置默认值
	style := req.Style
	if style == "" {
		style = "concise"
	}
	maxLength := req.MaxLength
	if maxLength == 0 {
		maxLength = 300
	}

	// 4. 调用AI生成总结
	summaryResult, err := s.generateSummary(ctx, fullContent, style, maxLength)
	if err != nil {
		return nil, fmt.Errorf("生成总结失败: %w", err)
	}

	// 5. 专业分析（如果开启）
	var analysis *ProfessionalAnalysis
	if req.EnableAnalysis {
		analysis, err = s.generateProfessionalAnalysis(ctx, fullContent)
		if err != nil {
			// 记录错误但不影响主流程
			fmt.Printf("专业分析失败: %v\n", err)
		}
	}

	// 6. 组装响应
	response := &SummaryResponse{
		Title:           summaryResult.Title,
		Summary:         summaryResult.Summary,
		KeyPoints:       summaryResult.KeyPoints,
		Highlights:      summaryResult.Highlights,
		Concerns:        summaryResult.Concerns,
		Tags:            summaryResult.Tags,
		Sentiment:       summaryResult.Sentiment,
		ReadingTime:     summaryResult.ReadingTime,
		TargetAudience:  summaryResult.TargetAudience,
		ContentType:     summaryResult.ContentType,
		ActionItems:     summaryResult.ActionItems,
		RelatedTopics:   summaryResult.RelatedTopics,
		CredibilityNote: summaryResult.CredibilityNote,
		QuickFacts:      summaryResult.QuickFacts,
		OCRTexts:        ocrTexts,
		Analysis:        analysis,
		OriginalURL:     req.URL,
	}

	return response, nil
}

// summaryResult AI总结结果结构
type summaryResult struct {
	Title           string            `json:"title"`            // 总结标题
	Summary         string            `json:"summary"`          // 核心观点概述
	KeyPoints       []string          `json:"key_points"`       // 关键信息点
	Highlights      []string          `json:"highlights"`       // 亮点/优点
	Concerns        []string          `json:"concerns"`         // 注意事项/缺点
	Tags            []string          `json:"tags"`             // 内容标签
	Sentiment       string            `json:"sentiment"`        // 情感倾向
	ReadingTime     string            `json:"reading_time"`     // 预计阅读时间
	TargetAudience  string            `json:"target_audience"`  // 目标受众
	ContentType     string            `json:"content_type"`     // 内容类型
	ActionItems     []string          `json:"action_items"`     // 可执行建议
	RelatedTopics   []string          `json:"related_topics"`   // 相关话题
	CredibilityNote string            `json:"credibility_note"` // 可信度说明
	QuickFacts      map[string]string `json:"quick_facts"`      // 快速事实（价格、时间、地点等）
}

// generateSummary 生成AI总结
func (s *SummaryService) generateSummary(ctx context.Context, content, style string, maxLength int) (*summaryResult, error) {
	stylePrompts := map[string]string{
		"concise":  "请用简洁凝练的语言总结，仅保留核心信息，适合快速浏览。",
		"detailed": "请详细全面地总结，包含完整逻辑链及细节信息，适合深度了解。",
		"casual":   "请用口语化、生活化的语言转述，让新手用户也能轻松理解。",
	}

	styleHint := stylePrompts["concise"]
	if hint, ok := stylePrompts[style]; ok {
		styleHint = hint
	}

	systemPrompt := fmt.Sprintf(`你是一位专业的小红书内容分析师，擅长深度分析和全面总结笔记内容。
%s

处理规则：
1. 如果提供了完整的笔记正文内容，请详细分析该内容
2. 如果只提供了笔记标题和作者，请根据标题推测笔记可能的内容方向，并基于你对小红书内容的理解进行分析
3. 标题通常能反映笔记的核心主题，请充分利用标题信息进行有价值的分析

请严格以JSON格式返回全面的分析结果：
{
  "title": "总结标题（10字以内，精准概括主题）",
  "summary": "核心观点概述（%d字以内，完整表达核心价值）",
  "key_points": ["关键信息1", "关键信息2", "关键信息3", "关键信息4", "关键信息5"],
  "highlights": ["亮点1", "亮点2", "亮点3"],
  "concerns": ["注意事项1", "注意事项2"],
  "tags": ["标签1", "标签2", "标签3", "标签4", "标签5"],
  "sentiment": "positive/neutral/negative",
  "reading_time": "预计阅读时间，如：2分钟",
  "target_audience": "目标受众描述，如：护肤新手、数码爱好者",
  "content_type": "内容类型，如：好物推荐/教程攻略/经验分享/测评对比/避坑指南",
  "action_items": ["可执行建议1", "可执行建议2", "可执行建议3"],
  "related_topics": ["相关话题1", "相关话题2", "相关话题3"],
  "credibility_note": "内容可信度评估说明",
  "quick_facts": {
    "价格": "如有提及",
    "品牌": "如有提及",
    "时间": "如有提及",
    "地点": "如有提及",
    "规格": "如有提及"
  }
}

字段说明：
1. key_points: 提取5-8个关键信息点，包括数据、价格、规格、时间、地点、效果等具体信息
2. highlights: 提取2-4个内容亮点或优势，让读者快速了解价值
3. concerns: 提取1-3个需要注意的事项、潜在风险或缺点
4. tags: 提取4-6个内容标签，便于分类和搜索
5. sentiment: 根据整体内容情感倾向判断（positive积极/neutral中性/negative消极）
6. reading_time: 根据内容长度估算阅读时间
7. target_audience: 分析内容最适合的目标人群
8. content_type: 判断内容属于哪种类型
9. action_items: 提取读者可以立即执行的具体建议
10. related_topics: 推荐相关的话题或延伸阅读方向
11. credibility_note: 评估内容的可信度，是否有数据支撑、是否为广告等
12. quick_facts: 提取关键事实数据，没有的字段可以省略

注意：
- 所有字段都要认真填写，确保分析全面深入
- 如果某些信息在原文中没有，可以标注"未提及"或省略该字段
- 如果无法访问链接或解析内容，请在summary中说明原因`, styleHint, maxLength)

	userMessage := fmt.Sprintf("请总结以下小红书笔记内容：\n\n%s", content)

	response, err := s.volcClient.SimpleChat(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var result summaryResult
	if err := parseJSONResponse(response, &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	return &result, nil
}

// generateProfessionalAnalysis 生成专业分析
func (s *SummaryService) generateProfessionalAnalysis(ctx context.Context, content string) (*ProfessionalAnalysis, error) {
	systemPrompt := `你是一位专业的内容分析专家，擅长对小红书笔记进行深度专业分析。

请根据内容自动识别品类，并进行针对性分析：
- 美妆护肤类：分析成分功效、安全评级、适配肤质、风险成分
- 数码产品类：分析核心参数、性能评分、场景适配、竞品对比
- 食品类：分析营养成分、热量、适用人群
- 旅游攻略类：分析行程安排、预算建议、注意事项

请严格以JSON格式返回结果：
{
  "category": "beauty/digital/food/travel/lifestyle/other",
  "main_conclusion": "核心结论（一句话概括）",
  "detailed_data": {
    "数据项1": "值1",
    "数据项2": "值2"
  },
  "risk_warnings": ["风险提示1", "风险提示2"],
  "recommendations": ["建议1", "建议2"]
}`

	userMessage := fmt.Sprintf("请对以下内容进行专业分析：\n\n%s", content)

	response, err := s.volcClient.SimpleChat(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, err
	}

	var result ProfessionalAnalysis
	if err := parseJSONResponse(response, &result); err != nil {
		return nil, fmt.Errorf("解析专业分析响应失败: %w", err)
	}

	return &result, nil
}

// ExtractTextFromImage 从图片提取文字
func (s *SummaryService) ExtractTextFromImage(ctx context.Context, imageURL string) (string, error) {
	return s.ocrClient.ExtractTextFromURL(ctx, imageURL)
}

// ExtractTextFromImageBytes 从图片字节提取文字
func (s *SummaryService) ExtractTextFromImageBytes(ctx context.Context, imageData []byte) (string, error) {
	return s.ocrClient.ExtractTextFromBytes(ctx, imageData)
}

// ExtractTextFromImageBase64 从图片Base64提取文字
func (s *SummaryService) ExtractTextFromImageBase64(ctx context.Context, imageBase64 string) (string, error) {
	return s.ocrClient.ExtractText(ctx, imageBase64)
}

// parseJSONResponse 解析JSON响应（处理可能的markdown代码块）
func parseJSONResponse(response string, result interface{}) error {
	// 尝试提取JSON内容（可能被markdown代码块包裹）
	jsonContent := response

	// 移除可能的markdown代码块标记
	if strings.Contains(response, "```json") {
		re := regexp.MustCompile("(?s)```json\\s*(.+?)\\s*```")
		matches := re.FindStringSubmatch(response)
		if len(matches) > 1 {
			jsonContent = matches[1]
		}
	} else if strings.Contains(response, "```") {
		re := regexp.MustCompile("(?s)```\\s*(.+?)\\s*```")
		matches := re.FindStringSubmatch(response)
		if len(matches) > 1 {
			jsonContent = matches[1]
		}
	}

	// 清理内容
	jsonContent = strings.TrimSpace(jsonContent)

	return json.Unmarshal([]byte(jsonContent), result)
}

// ValidateXHSURL 验证小红书链接格式
func ValidateXHSURL(url string) bool {
	// 支持的小红书链接格式
	patterns := []string{
		`^https?://www\.xiaohongshu\.com/explore/[a-zA-Z0-9]+`,
		`^https?://www\.xiaohongshu\.com/discovery/item/[a-zA-Z0-9]+`,
		`^https?://www\.xiaohongshu\.com/user/profile/[a-zA-Z0-9]+`,
		`^https?://xhslink\.com/[a-zA-Z0-9]+`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}

	return false
}

// ExtractNoteID 从URL提取笔记ID
func ExtractNoteID(url string) string {
	patterns := []struct {
		regex   string
		groupID int
	}{
		{`xiaohongshu\.com/explore/([a-zA-Z0-9]+)`, 1},
		{`xiaohongshu\.com/discovery/item/([a-zA-Z0-9]+)`, 1},
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.regex)
		matches := re.FindStringSubmatch(url)
		if len(matches) > p.groupID {
			return matches[p.groupID]
		}
	}

	return ""
}

// ========== 历史记录管理 ==========

// SummarizeAndSave 生成总结并保存
func (s *SummaryService) SummarizeAndSave(ctx context.Context, userID int64, req *SummaryRequest, noteTitle string) (*model.XHSSummary, error) {
	// 1. 生成总结
	resp, err := s.Summarize(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. 转换分析结果
	var analysis *model.AnalysisJSON
	if resp.Analysis != nil {
		analysis = &model.AnalysisJSON{
			Category:        resp.Analysis.Category,
			MainConclusion:  resp.Analysis.MainConclusion,
			DetailedData:    resp.Analysis.DetailedData,
			RiskWarnings:    resp.Analysis.RiskWarnings,
			Recommendations: resp.Analysis.Recommendations,
		}
	}

	// 3. 创建记录
	summary := &model.XHSSummary{
		UserID:         userID,
		OriginalURL:    req.URL,
		NoteTitle:      noteTitle,
		NoteContent:    req.Content,
		ImageURLs:      req.ImageURLs,
		SummaryTitle:   resp.Title,
		SummaryContent: resp.Summary,
		KeyPoints:      resp.KeyPoints,
		Tags:           resp.Tags,
		Sentiment:      resp.Sentiment,
		OCRTexts:       resp.OCRTexts,
		Analysis:       analysis,
		Style:          req.Style,
	}

	if err := s.summaryRepo.Create(summary); err != nil {
		return nil, fmt.Errorf("保存总结失败: %w", err)
	}

	return summary, nil
}

// GetSummaryByID 获取总结详情
func (s *SummaryService) GetSummaryByID(id int64) (*model.XHSSummary, error) {
	return s.summaryRepo.FindByID(id)
}

// GetUserSummaries 获取用户的总结历史
func (s *SummaryService) GetUserSummaries(userID int64, page, pageSize int) ([]model.XHSSummary, int64, error) {
	return s.summaryRepo.FindByUserID(userID, page, pageSize)
}

// DeleteSummary 删除总结
func (s *SummaryService) DeleteSummary(id, userID int64) error {
	return s.summaryRepo.Delete(id, userID)
}

// UpdateFavorite 更新收藏状态
func (s *SummaryService) UpdateFavorite(id, userID int64, isFavorite bool, folderID int64) error {
	return s.summaryRepo.UpdateFavorite(id, userID, isFavorite, folderID)
}

// GetFavorites 获取收藏列表
func (s *SummaryService) GetFavorites(userID int64, folderID int64, page, pageSize int) ([]model.XHSSummary, int64, error) {
	return s.summaryRepo.FindFavorites(userID, folderID, page, pageSize)
}

// SearchSummaries 搜索总结
func (s *SummaryService) SearchSummaries(userID int64, keyword string, page, pageSize int) ([]model.XHSSummary, int64, error) {
	return s.summaryRepo.Search(userID, keyword, page, pageSize)
}

// ========== 收藏夹管理 ==========

// CreateFolder 创建收藏夹
func (s *SummaryService) CreateFolder(userID int64, name string) (*model.XHSFolder, error) {
	folder := &model.XHSFolder{
		UserID: userID,
		Name:   name,
	}
	if err := s.folderRepo.Create(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// GetUserFolders 获取用户的收藏夹列表
func (s *SummaryService) GetUserFolders(userID int64) ([]model.XHSFolder, error) {
	return s.folderRepo.FindByUserID(userID)
}

// UpdateFolder 更新收藏夹
func (s *SummaryService) UpdateFolder(folder *model.XHSFolder) error {
	return s.folderRepo.Update(folder)
}

// DeleteFolder 删除收藏夹
func (s *SummaryService) DeleteFolder(id, userID int64) error {
	return s.folderRepo.Delete(id, userID)
}
