package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/pkg/volcengine"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service/design_suggestion"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var designSuggestionService *design_suggestion.DesignSuggestionService

// initDesignSuggestionService 初始化设计建议服务
// Initialize design suggestion service
func initDesignSuggestionService() {
	if designSuggestionService == nil {
		db := database.GetDB()
		repo := repository.NewDesignSuggestionRepository(db)
		projectRepo := repository.NewDesignerProjectRepository(db)
		aiClient := createVolcengineClient()
		designSuggestionService = design_suggestion.NewDesignSuggestionService(repo, projectRepo, aiClient)
	}
}

// createVolcengineClient 创建火山AI客户端
// Create VolcEngine AI client
func createVolcengineClient() *volcengine.Client {
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

	return volcengine.NewClient(apiKey, baseURL, model, timeout, maxTokens, temperature)
}

// getUserID 从上下文获取用户ID
// Get user ID from context
func getUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GenerateDesignSuggestions 异步生成设计建议
// @Summary 异步生成设计建议
// @Description 基于项目需求异步调用AI生成设计建议，完成后通过通知提醒
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param request body request.GenerateDesignSuggestionRequest true "生成建议请求"
// @Success 200 {object} response.GenerateSuggestionAsyncResponse
// @Router /api/designer/suggestions/generate [post]
func GenerateDesignSuggestions(c *gin.Context) {
	initDesignSuggestionService()

	var req request.GenerateDesignSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	result, err := designSuggestionService.GenerateSuggestionsAsync(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetSuggestionTaskStatus 获取建议生成任务状态
// @Summary 获取建议生成任务状态
// @Description 查询异步生成任务的当前状态
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param taskId path uint true "任务ID"
// @Success 200 {object} response.SuggestionTaskStatusResponse
// @Router /api/designer/suggestions/task/{taskId} [get]
func GetSuggestionTaskStatus(c *gin.Context) {
	initDesignSuggestionService()

	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的任务ID")
		return
	}

	userID := getUserID(c)
	result, err := designSuggestionService.GetTaskStatus(uint(taskID), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignSuggestionList 获取设计建议列表
// @Summary 获取设计建议列表
// @Description 获取项目的设计建议列表
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param projectId query uint true "项目ID"
// @Param current query int true "当前页"
// @Param size query int true "每页数量"
// @Param category query string false "类别筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} response.DesignSuggestionListResponse
// @Router /api/designer/suggestions [get]
func GetDesignSuggestionList(c *gin.Context) {
	initDesignSuggestionService()

	var req request.GetDesignSuggestionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	result, err := designSuggestionService.List(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignSuggestionDetail 获取设计建议详情
// @Summary 获取设计建议详情
// @Description 获取单个设计建议的详细信息
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param id path uint true "建议ID"
// @Success 200 {object} response.SuggestionDetailResponse
// @Router /api/designer/suggestions/{id} [get]
func GetDesignSuggestionDetail(c *gin.Context) {
	initDesignSuggestionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的建议ID")
		return
	}

	userID := getUserID(c)
	result, err := designSuggestionService.GetByID(uint(id), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateDesignSuggestionStatus 更新设计建议状态
// @Summary 更新设计建议状态
// @Description 更新建议状态（采纳/忽略）
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param request body request.UpdateSuggestionStatusRequest true "更新状态请求"
// @Success 200 {object} response.Response
// @Router /api/designer/suggestions/status [put]
func UpdateDesignSuggestionStatus(c *gin.Context) {
	initDesignSuggestionService()

	var req request.UpdateSuggestionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	if err := designSuggestionService.UpdateStatus(&req, userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchUpdateDesignSuggestionStatus 批量更新设计建议状态
// @Summary 批量更新设计建议状态
// @Description 批量更新多个建议的状态
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param request body request.BatchUpdateSuggestionStatusRequest true "批量更新状态请求"
// @Success 200 {object} response.Response
// @Router /api/designer/suggestions/batch-status [put]
func BatchUpdateDesignSuggestionStatus(c *gin.Context) {
	initDesignSuggestionService()

	var req request.BatchUpdateSuggestionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	if err := designSuggestionService.BatchUpdateStatus(&req, userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteDesignSuggestion 删除设计建议
// @Summary 删除设计建议
// @Description 删除单个设计建议
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param id path uint true "建议ID"
// @Success 200 {object} response.Response
// @Router /api/designer/suggestions/{id} [delete]
func DeleteDesignSuggestion(c *gin.Context) {
	initDesignSuggestionService()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的建议ID")
		return
	}

	userID := getUserID(c)
	if err := designSuggestionService.Delete(uint(id), userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDeleteDesignSuggestions 批量删除设计建议
// @Summary 批量删除设计建议
// @Description 批量删除多个设计建议
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteDesignSuggestionRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/suggestions/batch-delete [post]
func BatchDeleteDesignSuggestions(c *gin.Context) {
	initDesignSuggestionService()

	var req request.BatchDeleteDesignSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	if err := designSuggestionService.BatchDelete(req.IDs, userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetDesignSuggestionStats 获取设计建议统计
// @Summary 获取设计建议统计
// @Description 获取项目的设计建议统计信息
// @Tags 设计建议
// @Accept json
// @Produce json
// @Param projectId query uint false "项目ID"
// @Success 200 {object} response.SuggestionStatsResponse
// @Router /api/designer/suggestions/stats [get]
func GetDesignSuggestionStats(c *gin.Context) {
	initDesignSuggestionService()

	projectIDStr := c.Query("projectId")
	var projectID uint
	if projectIDStr != "" {
		id, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err != nil {
			response.Error(c, 400, "无效的项目ID")
			return
		}
		projectID = uint(id)
	}

	userID := getUserID(c)
	result, err := designSuggestionService.GetStats(projectID, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}
