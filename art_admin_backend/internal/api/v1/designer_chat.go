package v1

import (
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/designer_chat"
	"strconv"

	"github.com/gin-gonic/gin"
)

var designerChatService *designer_chat.DesignerChatService

// getDesignerChatService 获取设计师对话服务实例
// Get designer chat service instance
func getDesignerChatService() *designer_chat.DesignerChatService {
	if designerChatService == nil {
		designerChatService = designer_chat.NewDesignerChatService()
	}
	return designerChatService
}

// DesignerChatRequest 设计师对话请求
// Designer chat request
type DesignerChatRequest struct {
	Query     string `json:"query" binding:"required"`
	ProjectID uint   `json:"projectId"`
	SessionID string `json:"sessionId"`
}

// DesignerChat 设计师AI对话
// @Summary 设计师AI对话
// @Description 与设计师AI助手进行对话，支持项目上下文注入
// @Tags Designer-Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DesignerChatRequest true "对话请求"
// @Success 200 {object} response.Response "对话成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/designer/chat [post]
func DesignerChat(c *gin.Context) {
	var req DesignerChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID / Get current user ID
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	chatReq := &designer_chat.ChatRequest{
		Query:     req.Query,
		ProjectID: req.ProjectID,
		SessionID: req.SessionID,
	}

	result, err := getDesignerChatService().Chat(chatReq, userID)
	if err != nil {
		response.ServerError(c, "AI对话失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignerChatHistory 获取对话历史
// @Summary 获取对话历史
// @Description 获取设计师AI对话历史记录
// @Tags Designer-Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId query string false "会话ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/designer/chat/history [get]
func GetDesignerChatHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	sessionID := c.Query("sessionId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := getDesignerChatService().GetChatHistory(sessionID, userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取对话历史失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"data":     messages,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetDesignerChatSessions 获取会话列表
// @Summary 获取会话列表
// @Description 获取设计师AI对话会话列表
// @Tags Designer-Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId query int false "项目ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/designer/chat/sessions [get]
func GetDesignerChatSessions(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 32)

	sessions, err := getDesignerChatService().GetSessions(userID, uint(projectID))
	if err != nil {
		response.ServerError(c, "获取会话列表失败: "+err.Error())
		return
	}

	response.Success(c, sessions)
}

// DeleteDesignerChatSession 删除会话
// @Summary 删除会话
// @Description 删除指定的对话会话
// @Tags Designer-Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "会话ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/designer/chat/sessions/{sessionId} [delete]
func DeleteDesignerChatSession(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}

	err := getDesignerChatService().DeleteSession(sessionID, userID)
	if err != nil {
		response.ServerError(c, "删除会话失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// ExportDesignerChatHistory 导出对话历史
// @Summary 导出对话历史
// @Description 导出指定会话的对话历史为Markdown或JSON格式
// @Tags Designer-Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "会话ID"
// @Param format query string false "导出格式(md/json)" default(md)
// @Success 200 {object} response.Response "导出成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/designer/chat/export/{sessionId} [get]
func ExportDesignerChatHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}

	format := c.DefaultQuery("format", "md")

	if format == "json" {
		data, err := getDesignerChatService().ExportChatHistoryJSON(sessionID, userID)
		if err != nil {
			response.ServerError(c, "导出失败: "+err.Error())
			return
		}
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=chat_history.json")
		c.Data(200, "application/json", data)
		return
	}

	// 默认导出Markdown格式 / Default export as Markdown
	content, err := getDesignerChatService().ExportChatHistory(sessionID, userID)
	if err != nil {
		response.ServerError(c, "导出失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"content":  content,
		"format":   "markdown",
		"filename": "chat_history.md",
	})
}

// getUserIDFromContext 从上下文获取用户ID
// Get user ID from context
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}
