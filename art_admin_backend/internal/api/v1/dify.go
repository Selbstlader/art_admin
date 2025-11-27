package v1

import (
	"art_admin_backend/internal/dto/request"
	difyResponse "art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/response"
	difySvc "art_admin_backend/internal/service/dify"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var difyService *difySvc.DifyService

func getDifyService() *difySvc.DifyService {
	if difyService == nil {
		difyService = difySvc.NewDifyService()
	}
	return difyService
}

// GetDatasetList 获取知识库列表
// @Summary 获取知识库列表
// @Description 获取Dify知识库列表,支持分页和搜索
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" minimum(1)
// @Param limit query int false "每页数量" minimum(1) maximum(100)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/dataset/list [get]
func GetDatasetList(c *gin.Context) {
	var req request.DifyDatasetListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := getDifyService().GetDatasetList(&req)
	if err != nil {
		response.ServerError(c, "获取知识库列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetDatasetDetail 获取知识库详情
// @Summary 获取知识库详情
// @Description 获取指定知识库的详细信息
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "知识库ID"
// @Success 200 {object} response.Response "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/dataset/{id} [get]
func GetDatasetDetail(c *gin.Context) {
	datasetID := c.Param("id")
	if datasetID == "" {
		response.BadRequest(c, "知识库ID不能为空")
		return
	}

	result, err := getDifyService().GetDatasetDetail(datasetID)
	if err != nil {
		response.ServerError(c, "获取知识库详情失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// UploadFileToDataset 上传文件到知识库
// @Summary 上传文件到知识库
// @Description 上传文件到指定的Dify知识库
// @Tags Dify
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param dataset_id formData string true "知识库ID"
// @Param file formData file true "上传的文件"
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/dataset/upload [post]
func UploadFileToDataset(c *gin.Context) {
	// 获取知识库ID
	datasetID := c.PostForm("dataset_id")
	if datasetID == "" {
		response.BadRequest(c, "知识库ID不能为空")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "获取上传文件失败: "+err.Error())
		return
	}

	// 创建临时目录
	tempDir := "./temp/uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		response.ServerError(c, "创建临时目录失败: "+err.Error())
		return
	}

	// 保存文件到临时目录
	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		response.ServerError(c, "保存文件失败: "+err.Error())
		return
	}

	// 上传到Dify
	result, err := getDifyService().UploadFile(datasetID, tempFilePath)

	// 删除临时文件
	defer os.Remove(tempFilePath)

	if err != nil {
		response.ServerError(c, "上传文件到知识库失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "上传成功", result)
}

// ChatWithAI AI对话(非流式)
// @Summary AI对话
// @Description 与Dify AI进行对话(非流式)
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.DifyChatRequest true "对话请求"
// @Success 200 {object} response.Response "对话成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/chat [post]
func ChatWithAI(c *gin.Context) {
	var req request.DifyChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置为非流式模式
	req.ResponseMode = "blocking"

	result, err := getDifyService().Chat(&req)
	if err != nil {
		response.ServerError(c, "AI对话失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// ChatWithAIStreaming AI对话(流式)
// @Summary AI对话(流式)
// @Description 与Dify AI进行流式对话
// @Tags Dify
// @Accept json
// @Produce text/event-stream
// @Security BearerAuth
// @Param request body request.DifyChatRequest true "对话请求"
// @Success 200 {string} string "流式响应"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/chat/stream [post]
func ChatWithAIStreaming(c *gin.Context) {
	var req request.DifyChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置为流式模式
	req.ResponseMode = "streaming"

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 创建一个通道来传递错误
	errChan := make(chan error, 1)

	// 启动流式对话
	go func() {
		err := getDifyService().ChatStreaming(&req, func(event *difyResponse.DifyChatStreamEvent) error {
			// 将事件发送给客户端
			eventData := fmt.Sprintf("data: %s\n\n", mustMarshalJSON(event))

			// 写入响应
			if _, err := io.WriteString(c.Writer, eventData); err != nil {
				return err
			}

			// 刷新缓冲区
			c.Writer.Flush()

			return nil
		})

		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	// 等待完成或错误
	if err := <-errChan; err != nil {
		// 发送错误事件
		errorEvent := fmt.Sprintf("data: {\"event\":\"error\",\"message\":\"%s\"}\n\n", err.Error())
		io.WriteString(c.Writer, errorEvent)
		c.Writer.Flush()
	}
}

// mustMarshalJSON 将对象序列化为JSON字符串,失败时panic
func mustMarshalJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// StopChatMessage 停止消息生成
// @Summary 停止消息生成
// @Description 停止正在进行的AI对话消息生成(仅流式模式支持)
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "任务ID"
// @Param user query string true "用户标识"
// @Success 200 {object} response.Response "停止成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/chat/stop/{task_id} [post]
func StopChatMessage(c *gin.Context) {
	taskID := c.Param("task_id")
	user := c.Query("user")

	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}
	if user == "" {
		response.BadRequest(c, "用户标识不能为空")
		return
	}

	result, err := getDifyService().StopChatMessage(taskID, user)
	if err != nil {
		response.ServerError(c, "停止消息生成失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetSuggestedQuestions 获取建议问题
// @Summary 获取建议问题
// @Description 获取当前消息的下一轮建议问题列表
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message_id path string true "消息ID"
// @Param user query string true "用户标识"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/messages/{message_id}/suggested [get]
func GetSuggestedQuestions(c *gin.Context) {
	messageID := c.Param("message_id")
	user := c.Query("user")

	if messageID == "" {
		response.BadRequest(c, "消息ID不能为空")
		return
	}
	if user == "" {
		response.BadRequest(c, "用户标识不能为空")
		return
	}

	result, err := getDifyService().GetSuggestedQuestions(messageID, user)
	if err != nil {
		// 如果是建议问题功能被禁用，返回空数组而不是错误
		if strings.Contains(err.Error(), "Suggested Questions Is Disabled") {
			response.Success(c, &difyResponse.DifySuggestedQuestionsResponse{
				Data: []string{},
			})
			return
		}
		response.ServerError(c, "获取建议问题失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetConversations 获取会话列表
// @Summary 获取会话列表
// @Description 获取当前用户的会话列表
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user query string true "用户标识"
// @Param last_id query string false "最后一条记录ID(分页)"
// @Param limit query int false "每页数量(默认20,最大100)"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/conversations [get]
func GetConversations(c *gin.Context) {
	user := c.Query("user")
	lastID := c.Query("last_id")
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if user == "" {
		response.BadRequest(c, "用户标识不能为空")
		return
	}

	result, err := getDifyService().GetConversations(user, lastID, limit)
	if err != nil {
		response.ServerError(c, "获取会话列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetConversationMessages 获取会话消息历史
// @Summary 获取会话消息历史
// @Description 获取指定会话的历史消息记录
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation_id query string true "会话ID"
// @Param user query string true "用户标识"
// @Param first_id query string false "第一条记录ID(分页)"
// @Param limit query int false "每页数量(默认20,最大100)"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/messages [get]
func GetConversationMessages(c *gin.Context) {
	conversationID := c.Query("conversation_id")
	user := c.Query("user")
	firstID := c.Query("first_id")
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if conversationID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}
	if user == "" {
		response.BadRequest(c, "用户标识不能为空")
		return
	}

	result, err := getDifyService().GetConversationMessages(conversationID, user, firstID, limit)
	if err != nil {
		response.ServerError(c, "获取消息历史失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteConversation 删除会话
// @Summary 删除会话
// @Description 删除指定的会话
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation_id path string true "会话ID"
// @Param user query string true "用户标识"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/conversations/{conversation_id} [delete]
func DeleteConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	user := c.Query("user")

	if conversationID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}
	if user == "" {
		response.BadRequest(c, "用户标识不能为空")
		return
	}

	err := getDifyService().DeleteConversation(conversationID, user)
	if err != nil {
		response.ServerError(c, "删除会话失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// RenameConversation 重命名会话
// @Summary 重命名会话
// @Description 修改会话的名称
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation_id path string true "会话ID"
// @Param request body map[string]string true "请求参数"
// @Success 200 {object} response.Response "重命名成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/conversations/{conversation_id}/name [post]
func RenameConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req struct {
		User string `json:"user" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if conversationID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}

	result, err := getDifyService().RenameConversation(conversationID, req.User, req.Name)
	if err != nil {
		response.ServerError(c, "重命名会话失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteDataset 删除知识库
// @Summary 删除知识库
// @Description 删除指定的知识库
// @Tags Dify
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "知识库ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/dify/dataset/{id} [delete]
func DeleteDataset(c *gin.Context) {
	datasetID := c.Param("id")
	if datasetID == "" {
		response.BadRequest(c, "知识库ID不能为空")
		return
	}

	err := getDifyService().DeleteDataset(datasetID)
	if err != nil {
		response.ServerError(c, "删除知识库失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}
