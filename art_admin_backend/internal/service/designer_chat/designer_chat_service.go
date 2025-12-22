package designer_chat

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/volcengine"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// ProjectContext 项目完整上下文信息
// Full project context for AI chat
type ProjectContext struct {
	Project       *model.DesignerProject
	Documents     []DocumentInfo
	CadFiles      []CadFileInfo
	RenderRecords []RenderInfo
	Versions      []VersionInfo
	Materials     []MaterialInfo
	CostSummary   *CostSummaryInfo
}

// DocumentInfo 文档信息
type DocumentInfo struct {
	FileName       string `json:"fileName"`
	FileType       string `json:"fileType"`
	Summary        string `json:"summary"`
	Keywords       string `json:"keywords"`
	AnalysisStatus string `json:"analysisStatus"`
}

// CadFileInfo CAD文件信息
type CadFileInfo struct {
	FileName    string   `json:"fileName"`
	FileFormat  string   `json:"fileFormat"`
	LayerCount  int      `json:"layerCount"`
	Layers      []string `json:"layers"`
	Has3D       bool     `json:"has3d"`
	ParseStatus string   `json:"parseStatus"`
}

// RenderInfo 效果图信息
type RenderInfo struct {
	FileName   string `json:"fileName"`
	Style      string `json:"style"`
	Status     string `json:"status"`
	ResultPath string `json:"resultPath"`
}

// VersionInfo 设计版本信息
type VersionInfo struct {
	VersionNumber string `json:"versionNumber"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	CadFileCount  int    `json:"cadFileCount"`
	RenderCount   int    `json:"renderCount"`
}

// MaterialInfo 材料信息
type MaterialInfo struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

// CostSummaryInfo 成本汇总信息
type CostSummaryInfo struct {
	MaterialCost   float64 `json:"materialCost"`
	LaborCost      float64 `json:"laborCost"`
	EquipmentCost  float64 `json:"equipmentCost"`
	ManagementCost float64 `json:"managementCost"`
	TotalCost      float64 `json:"totalCost"`
	BudgetLimit    float64 `json:"budgetLimit"`
	BudgetExceeded bool    `json:"budgetExceeded"`
}

// DesignerChatService 设计师AI对话服务
// Designer AI chat service using VolcEngine
type DesignerChatService struct {
	db       *gorm.DB
	aiClient *volcengine.Client
}

// ChatMessage 对话消息
// Chat message structure
type ChatMessage struct {
	ID        uint      `json:"id"`
	SessionID string    `json:"sessionId"`
	ProjectID uint      `json:"projectId"`
	Role      string    `json:"role"` // user/assistant
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// ChatSession 对话会话
// Chat session structure
type ChatSession struct {
	ID        string    `json:"id"`
	ProjectID uint      `json:"projectId"`
	UserID    uint      `json:"userId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatRequest 对话请求
// Chat request structure
type ChatRequest struct {
	Query     string `json:"query" binding:"required"`
	ProjectID uint   `json:"projectId"`
	SessionID string `json:"sessionId"`
}

// ChatResponse 对话响应
// Chat response structure
type ChatResponse struct {
	Answer       string `json:"answer"`
	SessionID    string `json:"sessionId"`
	MessageID    uint   `json:"messageId"`
	TokensUsed   int    `json:"tokensUsed"`
	ResponseTime int64  `json:"responseTime"` // 响应时间(毫秒)
}

// StreamChunkCallback 流式输出回调
// Stream chunk callback
type StreamChunkCallback func(content string, done bool) error

// NewDesignerChatService 创建设计师对话服务
// Create designer chat service
func NewDesignerChatService() *DesignerChatService {
	apiKey := viper.GetString("volcengine.apiKey")
	baseURL := viper.GetString("volcengine.baseUrl")
	model := viper.GetString("volcengine.model")
	timeout := viper.GetInt("volcengine.timeout")
	maxTokens := viper.GetInt("volcengine.maxTokens")
	temperature := viper.GetFloat64("volcengine.temperature")

	if timeout == 0 {
		timeout = 120
	}
	if maxTokens == 0 {
		maxTokens = 4096
	}
	if temperature == 0 {
		temperature = 0.7
	}

	return &DesignerChatService{
		db:       database.GetDB(),
		aiClient: volcengine.NewClient(apiKey, baseURL, model, timeout, maxTokens, temperature),
	}
}

// Chat 发送对话消息
// Send chat message
// Requirements: 5.1 - 基于当前项目信息提供上下文相关的回答
// Requirements: 5.2 - 引用相关国家标准或行业规范进行回答
func (s *DesignerChatService) Chat(req *ChatRequest, userID uint) (*ChatResponse, error) {
	startTime := time.Now()

	// 构建系统提示词 / Build system prompt
	systemPrompt := s.buildSystemPrompt(req.ProjectID)

	// 获取历史消息作为上下文 / Get history messages as context
	historyMessages := s.getHistoryMessages(req.SessionID, 10)

	// 构建消息列表 / Build message list
	messages := []volcengine.ChatMessage{
		{
			Role: "system",
			Content: []volcengine.ContentPart{
				{Type: "text", Text: systemPrompt},
			},
		},
	}

	// 添加历史消息 / Add history messages
	for _, msg := range historyMessages {
		messages = append(messages, volcengine.ChatMessage{
			Role: msg.Role,
			Content: []volcengine.ContentPart{
				{Type: "text", Text: msg.Content},
			},
		})
	}

	// 添加当前用户消息 / Add current user message
	messages = append(messages, volcengine.ChatMessage{
		Role: "user",
		Content: []volcengine.ContentPart{
			{Type: "text", Text: req.Query},
		},
	})

	// 调用AI服务 / Call AI service
	resp, err := s.aiClient.Chat(messages)
	if err != nil {
		return nil, fmt.Errorf("AI服务调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI服务未返回结果")
	}

	answer := resp.Choices[0].Message.Content
	responseTime := time.Since(startTime).Milliseconds()

	// 生成或使用会话ID / Generate or use session ID
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d_%d", userID, time.Now().UnixNano())
	}

	// 保存对话记录 / Save chat history
	userMsg := &model.DesignerChatMessage{
		SessionID: sessionID,
		ProjectID: req.ProjectID,
		UserID:    userID,
		Role:      "user",
		Content:   req.Query,
	}
	s.db.Create(userMsg)

	assistantMsg := &model.DesignerChatMessage{
		SessionID:    sessionID,
		ProjectID:    req.ProjectID,
		UserID:       userID,
		Role:         "assistant",
		Content:      answer,
		TokensUsed:   resp.Usage.TotalTokens,
		ResponseTime: responseTime,
	}
	s.db.Create(assistantMsg)

	// 记录Token使用日志 / Log token usage
	// Requirements: 13.4 - 记录每次AI调用的token消耗和响应时间
	s.logTokenUsage(userID, req.ProjectID, resp.Usage.TotalTokens, responseTime)

	return &ChatResponse{
		Answer:       answer,
		SessionID:    sessionID,
		MessageID:    assistantMsg.ID,
		TokensUsed:   resp.Usage.TotalTokens,
		ResponseTime: responseTime,
	}, nil
}

// ChatStream 流式对话
// Stream chat message
// Requirements: 5.1 - 基于当前项目信息提供上下文相关的回答
func (s *DesignerChatService) ChatStream(req *ChatRequest, userID uint, callback StreamChunkCallback) (*ChatResponse, error) {
	startTime := time.Now()

	// 构建系统提示词 / Build system prompt
	systemPrompt := s.buildSystemPrompt(req.ProjectID)

	// 获取历史消息作为上下文 / Get history messages as context
	historyMessages := s.getHistoryMessages(req.SessionID, 10)

	// 构建消息列表 / Build message list
	messages := []volcengine.ChatMessage{
		{
			Role: "system",
			Content: []volcengine.ContentPart{
				{Type: "text", Text: systemPrompt},
			},
		},
	}

	// 添加历史消息 / Add history messages
	for _, msg := range historyMessages {
		messages = append(messages, volcengine.ChatMessage{
			Role: msg.Role,
			Content: []volcengine.ContentPart{
				{Type: "text", Text: msg.Content},
			},
		})
	}

	// 添加当前用户消息 / Add current user message
	messages = append(messages, volcengine.ChatMessage{
		Role: "user",
		Content: []volcengine.ContentPart{
			{Type: "text", Text: req.Query},
		},
	})

	// 生成或使用会话ID / Generate or use session ID
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d_%d", userID, time.Now().UnixNano())
	}

	// 保存用户消息 / Save user message
	userMsg := &model.DesignerChatMessage{
		SessionID: sessionID,
		ProjectID: req.ProjectID,
		UserID:    userID,
		Role:      "user",
		Content:   req.Query,
	}
	s.db.Create(userMsg)

	// 调用流式AI服务 / Call streaming AI service
	resp, err := s.aiClient.ChatStream(messages, func(chunk volcengine.StreamChunk) error {
		return callback(chunk.Content, chunk.Done)
	})
	if err != nil {
		return nil, fmt.Errorf("AI服务调用失败: %w", err)
	}

	answer := ""
	if len(resp.Choices) > 0 {
		answer = resp.Choices[0].Message.Content
	}

	responseTime := time.Since(startTime).Milliseconds()

	// 保存AI回复 / Save AI response
	assistantMsg := &model.DesignerChatMessage{
		SessionID:    sessionID,
		ProjectID:    req.ProjectID,
		UserID:       userID,
		Role:         "assistant",
		Content:      answer,
		TokensUsed:   resp.Usage.TotalTokens,
		ResponseTime: responseTime,
	}
	s.db.Create(assistantMsg)

	// 记录Token使用日志 / Log token usage
	s.logTokenUsage(userID, req.ProjectID, resp.Usage.TotalTokens, responseTime)

	return &ChatResponse{
		Answer:       answer,
		SessionID:    sessionID,
		MessageID:    assistantMsg.ID,
		TokensUsed:   resp.Usage.TotalTokens,
		ResponseTime: responseTime,
	}, nil
}

// buildSystemPrompt 构建系统提示词
// Build system prompt with project context
// Requirements: 5.1 - 基于当前项目信息提供上下文相关的回答
func (s *DesignerChatService) buildSystemPrompt(projectID uint) string {
	basePrompt := `你是一位专业的工装设计师AI助手，具备以下专业能力：
1. 熟悉工装设计（工业厂房、办公空间、商业空间等）的设计规范和标准
2. 了解国家建筑设计规范、消防规范、无障碍设计规范等
3. 熟悉各类装修材料的特性、价格和适用场景
4. 能够进行成本估算和预算分析
5. 了解施工工艺和流程

请基于用户的问题提供专业、准确、实用的建议。如果涉及具体数值计算，请提供计算过程。
如果涉及设计规范，请引用相关国家标准或行业规范。`

	// 如果有项目上下文，添加完整项目信息 / Add full project context if available
	if projectID > 0 {
		ctx := s.getProjectContext(projectID)
		if ctx != nil && ctx.Project != nil {
			basePrompt += s.buildProjectContextPrompt(ctx)
		}
	}

	return basePrompt
}

// getProjectContext 获取项目完整上下文
// Get full project context for AI chat
func (s *DesignerChatService) getProjectContext(projectID uint) *ProjectContext {
	var project model.DesignerProject
	if err := s.db.Preload("Documents").Preload("CadFiles").First(&project, projectID).Error; err != nil {
		return nil
	}

	ctx := &ProjectContext{
		Project: &project,
	}

	// 获取文档信息 / Get documents info
	for _, doc := range project.Documents {
		ctx.Documents = append(ctx.Documents, DocumentInfo{
			FileName:       doc.FileName,
			FileType:       doc.FileType,
			Summary:        doc.Summary,
			Keywords:       doc.Keywords,
			AnalysisStatus: doc.AnalysisStatus,
		})
	}

	// 获取CAD文件信息 / Get CAD files info
	for _, cad := range project.CadFiles {
		var layers []string
		if cad.Layers != "" {
			json.Unmarshal([]byte(cad.Layers), &layers)
		}
		ctx.CadFiles = append(ctx.CadFiles, CadFileInfo{
			FileName:    cad.FileName,
			FileFormat:  cad.FileFormat,
			LayerCount:  cad.LayerCount,
			Layers:      layers,
			Has3D:       cad.Has3D,
			ParseStatus: cad.ParseStatus,
		})
	}

	// 获取效果图记录 / Get render records
	var renders []model.RenderRecord
	if err := s.db.Where("project_id = ? AND status = ?", projectID, "completed").Find(&renders).Error; err == nil {
		for _, r := range renders {
			ctx.RenderRecords = append(ctx.RenderRecords, RenderInfo{
				FileName:   r.RoomType, // 使用RoomType作为名称
				Style:      r.Style,
				Status:     r.Status,
				ResultPath: r.ImageURL, // 使用ImageURL
			})
		}
	}

	// 获取设计版本 / Get design versions
	var versions []model.DesignVersion
	if err := s.db.Where("project_id = ?", projectID).Order("created_at DESC").Limit(5).Find(&versions).Error; err == nil {
		for _, v := range versions {
			// 统计版本关联的CAD和效果图数量
			var cadCount, renderCount int64
			s.db.Model(&model.CadFile{}).Where("project_id = ?", projectID).Count(&cadCount)
			s.db.Model(&model.RenderRecord{}).Where("project_id = ? AND status = ?", projectID, "completed").Count(&renderCount)

			ctx.Versions = append(ctx.Versions, VersionInfo{
				VersionNumber: fmt.Sprintf("V%d", v.VersionNumber), // 转换为字符串
				Description:   v.Description,
				Status:        v.Status,
				CadFileCount:  int(cadCount),
				RenderCount:   int(renderCount),
			})
		}
	}

	// 获取材料清单 / Get materials
	var materials []model.ProjectMaterial
	if err := s.db.Where("project_id = ?", projectID).Find(&materials).Error; err == nil {
		for _, m := range materials {
			ctx.Materials = append(ctx.Materials, MaterialInfo{
				Name:       m.Name,
				Category:   m.Category,
				Quantity:   m.Quantity,
				Unit:       m.Unit,
				UnitPrice:  m.UnitPrice,
				TotalPrice: m.TotalPrice,
			})
		}
	}

	// 获取成本汇总 / Get cost summary
	var costEstimate model.CostEstimate
	if err := s.db.Where("project_id = ?", projectID).First(&costEstimate).Error; err == nil {
		var materialCost float64
		s.db.Model(&model.ProjectMaterial{}).Where("project_id = ?", projectID).
			Select("COALESCE(SUM(total_price), 0)").Scan(&materialCost)

		totalCost := materialCost + costEstimate.LaborCost + costEstimate.EquipmentCost + costEstimate.ManagementCost
		ctx.CostSummary = &CostSummaryInfo{
			MaterialCost:   materialCost,
			LaborCost:      costEstimate.LaborCost,
			EquipmentCost:  costEstimate.EquipmentCost,
			ManagementCost: costEstimate.ManagementCost,
			TotalCost:      totalCost,
			BudgetLimit:    costEstimate.BudgetLimit,
			BudgetExceeded: costEstimate.BudgetLimit > 0 && totalCost > costEstimate.BudgetLimit,
		}
	}

	return ctx
}

// buildProjectContextPrompt 构建项目上下文提示词
// Build project context prompt for AI
func (s *DesignerChatService) buildProjectContextPrompt(ctx *ProjectContext) string {
	var sb strings.Builder
	project := ctx.Project

	sb.WriteString("\n\n=== 当前项目信息 ===\n")
	sb.WriteString(fmt.Sprintf("项目名称：%s\n", project.Name))
	sb.WriteString(fmt.Sprintf("项目描述：%s\n", project.Description))
	sb.WriteString(fmt.Sprintf("面积：%.2f 平方米\n", project.Area))
	sb.WriteString(fmt.Sprintf("预算：%.2f 元\n", project.Budget))
	sb.WriteString(fmt.Sprintf("设计风格：%s\n", project.Style))
	sb.WriteString(fmt.Sprintf("项目状态：%s\n", project.Status))

	// 添加文档分析信息 / Add document analysis info
	if len(ctx.Documents) > 0 {
		sb.WriteString("\n=== 项目文档（用户上传的需求文档）===\n")
		for _, doc := range ctx.Documents {
			sb.WriteString(fmt.Sprintf("- 文档【%s】(%s)", doc.FileName, doc.FileType))
			if doc.AnalysisStatus == "completed" {
				if doc.Summary != "" {
					sb.WriteString(fmt.Sprintf("\n  摘要: %s", doc.Summary))
				}
				if doc.Keywords != "" && doc.Keywords != "[]" {
					sb.WriteString(fmt.Sprintf("\n  关键词: %s", doc.Keywords))
				}
			} else {
				sb.WriteString(fmt.Sprintf(" [分析状态: %s]", doc.AnalysisStatus))
			}
			sb.WriteString("\n")
		}
	}

	// 添加CAD文件信息 / Add CAD files info
	if len(ctx.CadFiles) > 0 {
		sb.WriteString("\n=== CAD图纸文件 ===\n")
		for _, cad := range ctx.CadFiles {
			sb.WriteString(fmt.Sprintf("- 图纸【%s】(格式:%s, 图层数:%d", cad.FileName, cad.FileFormat, cad.LayerCount))
			if cad.Has3D {
				sb.WriteString(", 包含3D信息")
			}
			sb.WriteString(")\n")
			if len(cad.Layers) > 0 {
				if len(cad.Layers) <= 10 {
					sb.WriteString(fmt.Sprintf("  图层: %s\n", strings.Join(cad.Layers, ", ")))
				} else {
					sb.WriteString(fmt.Sprintf("  主要图层: %s 等%d个图层\n", strings.Join(cad.Layers[:10], ", "), len(cad.Layers)))
				}
			}
		}
	}

	// 添加效果图信息 / Add render info
	if len(ctx.RenderRecords) > 0 {
		sb.WriteString("\n=== 已生成效果图 ===\n")
		sb.WriteString(fmt.Sprintf("共 %d 张效果图\n", len(ctx.RenderRecords)))
		for i, r := range ctx.RenderRecords {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("... 等共%d张效果图\n", len(ctx.RenderRecords)))
				break
			}
			sb.WriteString(fmt.Sprintf("- %s (风格: %s)\n", r.FileName, r.Style))
		}
	}

	// 添加设计版本信息 / Add version info
	if len(ctx.Versions) > 0 {
		sb.WriteString("\n=== 设计版本 ===\n")
		for _, v := range ctx.Versions {
			sb.WriteString(fmt.Sprintf("- 版本 %s: %s (状态: %s)\n", v.VersionNumber, v.Description, v.Status))
		}
	}

	// 添加材料清单信息 / Add material list info
	if len(ctx.Materials) > 0 {
		sb.WriteString("\n=== 项目材料清单 ===\n")
		// 按类别分组统计 / Group by category
		categoryMap := make(map[string][]MaterialInfo)
		for _, m := range ctx.Materials {
			cat := m.Category
			if cat == "" {
				cat = "其他"
			}
			categoryMap[cat] = append(categoryMap[cat], m)
		}
		for cat, materials := range categoryMap {
			var totalCost float64
			for _, m := range materials {
				totalCost += m.TotalPrice
			}
			sb.WriteString(fmt.Sprintf("- %s类材料: %d种, 合计%.2f元\n", cat, len(materials), totalCost))
			// 列出前3种材料 / List top 3 materials
			for i, m := range materials {
				if i >= 3 {
					sb.WriteString(fmt.Sprintf("  ... 等%d种材料\n", len(materials)-3))
					break
				}
				sb.WriteString(fmt.Sprintf("  · %s: %.2f%s × %.2f元 = %.2f元\n",
					m.Name, m.Quantity, m.Unit, m.UnitPrice, m.TotalPrice))
			}
		}
	}

	// 添加成本汇总信息 / Add cost summary info
	if ctx.CostSummary != nil {
		sb.WriteString("\n=== 成本汇总 ===\n")
		sb.WriteString(fmt.Sprintf("- 材料费: %.2f元\n", ctx.CostSummary.MaterialCost))
		sb.WriteString(fmt.Sprintf("- 人工费: %.2f元\n", ctx.CostSummary.LaborCost))
		sb.WriteString(fmt.Sprintf("- 设备费: %.2f元\n", ctx.CostSummary.EquipmentCost))
		sb.WriteString(fmt.Sprintf("- 管理费: %.2f元\n", ctx.CostSummary.ManagementCost))
		sb.WriteString(fmt.Sprintf("- 总成本: %.2f元\n", ctx.CostSummary.TotalCost))
		if ctx.CostSummary.BudgetLimit > 0 {
			sb.WriteString(fmt.Sprintf("- 预算上限: %.2f元\n", ctx.CostSummary.BudgetLimit))
			if ctx.CostSummary.BudgetExceeded {
				exceeded := ctx.CostSummary.TotalCost - ctx.CostSummary.BudgetLimit
				sb.WriteString(fmt.Sprintf("⚠️ 警告：当前成本已超出预算 %.2f元\n", exceeded))
			}
		}
	}

	sb.WriteString("\n请基于以上项目完整信息，为用户提供针对性的设计建议和解答。")

	return sb.String()
}

// getHistoryMessages 获取历史消息
// Get history messages for context
func (s *DesignerChatService) getHistoryMessages(sessionID string, limit int) []ChatMessage {
	if sessionID == "" {
		return nil
	}

	var messages []model.DesignerChatMessage
	s.db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages)

	result := make([]ChatMessage, len(messages))
	for i, msg := range messages {
		result[i] = ChatMessage{
			ID:        msg.ID,
			SessionID: msg.SessionID,
			ProjectID: msg.ProjectID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	return result
}

// logTokenUsage 记录Token使用日志
// Log token usage for monitoring
// Requirements: 13.4 - 记录每次AI调用的token消耗和响应时间
func (s *DesignerChatService) logTokenUsage(userID, projectID uint, tokens int, responseTime int64) {
	log := &model.AIUsageLog{
		UserID:       userID,
		ProjectID:    projectID,
		ServiceType:  "designer_chat",
		TokensUsed:   tokens,
		ResponseTime: responseTime,
	}
	s.db.Create(log)
}

// GetChatHistory 获取对话历史
// Get chat history
// Requirements: 5.4 - 保存对话历史记录
func (s *DesignerChatService) GetChatHistory(sessionID string, userID uint, page, pageSize int) ([]ChatMessage, int64, error) {
	var messages []model.DesignerChatMessage
	var total int64

	query := s.db.Model(&model.DesignerChatMessage{}).Where("user_id = ?", userID)
	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]ChatMessage, len(messages))
	for i, msg := range messages {
		result[i] = ChatMessage{
			ID:        msg.ID,
			SessionID: msg.SessionID,
			ProjectID: msg.ProjectID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	return result, total, nil
}

// GetSessions 获取会话列表
// Get session list
func (s *DesignerChatService) GetSessions(userID uint, projectID uint) ([]ChatSession, error) {
	var sessions []ChatSession

	query := `
		SELECT 
			session_id as id,
			project_id,
			user_id,
			MIN(content) as title,
			MIN(created_at) as created_at,
			MAX(created_at) as updated_at
		FROM designer_chat_messages
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	if projectID > 0 {
		query += " AND project_id = ?"
		args = append(args, projectID)
	}

	query += " GROUP BY session_id, project_id, user_id ORDER BY updated_at DESC"

	if err := s.db.Raw(query, args...).Scan(&sessions).Error; err != nil {
		return nil, err
	}

	// 截取标题 / Truncate title
	for i := range sessions {
		if len(sessions[i].Title) > 50 {
			sessions[i].Title = sessions[i].Title[:50] + "..."
		}
	}

	return sessions, nil
}

// DeleteSession 删除会话
// Delete session
func (s *DesignerChatService) DeleteSession(sessionID string, userID uint) error {
	return s.db.Where("session_id = ? AND user_id = ?", sessionID, userID).
		Delete(&model.DesignerChatMessage{}).Error
}

// ExportChatHistory 导出对话历史
// Export chat history to document format
// Requirements: 5.4 - 支持导出为文档格式
func (s *DesignerChatService) ExportChatHistory(sessionID string, userID uint) (string, error) {
	var messages []model.DesignerChatMessage
	err := s.db.Where("session_id = ? AND user_id = ?", sessionID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return "", err
	}

	if len(messages) == 0 {
		return "", fmt.Errorf("没有找到对话记录")
	}

	// 构建导出内容 / Build export content
	var content string
	content += "# 设计师AI对话记录\n\n"
	content += fmt.Sprintf("导出时间：%s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	content += "---\n\n"

	for _, msg := range messages {
		roleLabel := "用户"
		if msg.Role == "assistant" {
			roleLabel = "AI助手"
		}
		content += fmt.Sprintf("### %s (%s)\n\n", roleLabel, msg.CreatedAt.Format("2006-01-02 15:04:05"))
		content += msg.Content + "\n\n"
		content += "---\n\n"
	}

	return content, nil
}

// ExportChatHistoryJSON 导出对话历史为JSON格式
// Export chat history to JSON format
func (s *DesignerChatService) ExportChatHistoryJSON(sessionID string, userID uint) ([]byte, error) {
	var messages []model.DesignerChatMessage
	err := s.db.Where("session_id = ? AND user_id = ?", sessionID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	exportData := struct {
		SessionID  string                      `json:"sessionId"`
		ExportTime string                      `json:"exportTime"`
		Messages   []model.DesignerChatMessage `json:"messages"`
	}{
		SessionID:  sessionID,
		ExportTime: time.Now().Format("2006-01-02 15:04:05"),
		Messages:   messages,
	}

	return json.MarshalIndent(exportData, "", "  ")
}
