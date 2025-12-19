package designer_chat

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/volcengine"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

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

	// 如果有项目上下文，添加项目信息 / Add project context if available
	if projectID > 0 {
		var project model.DesignerProject
		if err := s.db.First(&project, projectID).Error; err == nil {
			projectContext := fmt.Sprintf(`

当前项目信息：
- 项目名称：%s
- 项目描述：%s
- 面积：%.2f 平方米
- 预算：%.2f 元
- 设计风格：%s
- 项目状态：%s

请基于以上项目信息，为用户提供针对性的设计建议和解答。`,
				project.Name,
				project.Description,
				project.Area,
				project.Budget,
				project.Style,
				project.Status,
			)
			basePrompt += projectContext
		}
	}

	return basePrompt
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
