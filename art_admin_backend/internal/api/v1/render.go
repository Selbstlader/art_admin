package v1

import (
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/cad"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var renderService *cad.RenderService
var renderRecordService *cad.RenderRecordService

// SetRenderService 设置渲染服务
func SetRenderService(service *cad.RenderService) {
	renderService = service
}

// SetRenderRecordService 设置效果图记录服务
func SetRenderRecordService(service *cad.RenderRecordService) {
	renderRecordService = service
}

// GenerateRenderRequest 生成效果图请求
type GenerateRenderRequest struct {
	CadFileID   uint   `json:"cadFileId" binding:"required"`
	ProjectID   uint   `json:"projectId"`
	DocumentIDs []uint `json:"documentIds"` // 参考文档ID列表
	Style       string `json:"style"`       // 设计风格（用户输入）
	RoomType    string `json:"roomType"`    // 空间类型（用户输入）
	Description string `json:"description"` // 设计要求描述
}

// GenerateCadRender 提交CAD效果图生成任务（异步）
// @Summary 提交CAD效果图生成任务
// @Description 异步生成AI效果图，立即返回任务ID，结果在"我的效果图"中查看
// @Tags CAD渲染
// @Accept json
// @Produce json
// @Param request body GenerateRenderRequest true "渲染请求"
// @Success 200 {object} response.Response
// @Router /api/designer/cad/render [post]
func GenerateCadRender(c *gin.Context) {
	var req GenerateRenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取用户ID / Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 提交异步任务 / Submit async task
	result, err := renderRecordService.SubmitRenderTask(&cad.SubmitRenderRequest{
		UserID:      userID.(int64),
		ProjectID:   req.ProjectID,
		CadFileID:   req.CadFileID,
		DocumentIDs: req.DocumentIDs,
		Style:       req.Style,
		RoomType:    req.RoomType,
		Description: req.Description,
	})

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// GetRenderHistory 获取渲染历史
// @Summary 获取CAD渲染历史
// @Description 获取指定CAD文件的渲染历史记录
// @Tags CAD渲染
// @Accept json
// @Produce json
// @Param cadFileId query int true "CAD文件ID"
// @Success 200 {object} response.Response{data=response.RenderHistoryResponse}
// @Router /api/designer/cad/render/history [get]
func GetRenderHistory(c *gin.Context) {
	cadFileIDStr := c.Query("cadFileId")
	cadFileID, err := strconv.ParseUint(cadFileIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的CAD文件ID")
		return
	}

	records, err := renderService.GetRenderHistory(uint(cadFileID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"cadFileId": cadFileID,
		"records":   records,
		"total":     len(records),
	})
}

// GetRenderStyles 获取可用的渲染风格
// @Summary 获取渲染风格列表
// @Description 获取所有可用的室内设计风格
// @Tags CAD渲染
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/designer/cad/render/styles [get]
func GetRenderStyles(c *gin.Context) {
	styles := []map[string]string{
		{"value": "modern", "label": "现代简约", "description": "简洁线条，中性色调，大面积留白"},
		{"value": "chinese", "label": "新中式", "description": "木质元素，传统纹样，温暖色调"},
		{"value": "european", "label": "欧式古典", "description": "精致雕花，华丽装饰，暖色灯光"},
		{"value": "minimalist", "label": "极简主义", "description": "纯白空间，几何造型，自然光线"},
		{"value": "nordic", "label": "北欧风格", "description": "原木色调，简约家具，明亮通透"},
		{"value": "japanese", "label": "日式风格", "description": "榻榻米，障子门，禅意空间"},
	}

	roomTypes := []map[string]string{
		{"value": "living", "label": "客厅"},
		{"value": "bedroom", "label": "卧室"},
		{"value": "kitchen", "label": "厨房"},
		{"value": "bathroom", "label": "卫生间"},
		{"value": "study", "label": "书房"},
		{"value": "dining", "label": "餐厅"},
	}

	viewAngles := []map[string]string{
		{"value": "perspective", "label": "透视图"},
		{"value": "top", "label": "俯视图"},
		{"value": "front", "label": "正视图"},
	}

	response.Success(c, map[string]interface{}{
		"styles":     styles,
		"roomTypes":  roomTypes,
		"viewAngles": viewAngles,
	})
}

// GetUserRenderQuota 获取用户效果图配额
// @Summary 获取用户效果图配额
// @Description 获取当前用户的效果图保存配额信息
// @Tags CAD渲染
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/designer/cad/render/quota [get]
func GetUserRenderQuota(c *gin.Context) {
	// 从上下文获取用户ID / Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	quota, err := renderRecordService.GetUserQuota(userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"dailyLimit":   quota.DailyLimit,
		"usedToday":    quota.UsedToday,
		"remainingUse": quota.DailyLimit - quota.UsedToday,
	})
}

// GetUserRenderQuotaByAdmin 管理员获取指定用户效果图配额
// @Summary 管理员获取指定用户效果图配额
// @Description 管理员获取指定用户的效果图配额信息
// @Tags 用户管理
// @Produce json
// @Param userId query int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/system/user/render-quota [get]
func GetUserRenderQuotaByAdmin(c *gin.Context) {
	userIDStr := c.Query("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	quota, err := renderRecordService.GetUserQuota(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"userId":       userID,
		"dailyLimit":   quota.DailyLimit,
		"usedToday":    quota.UsedToday,
		"remainingUse": quota.DailyLimit - quota.UsedToday,
	})
}

// ListRenderRecords 获取效果图记录列表
// @Summary 获取效果图记录列表
// @Description 获取当前用户保存的效果图记录
// @Tags CAD渲染
// @Produce json
// @Param current query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param projectId query int false "项目ID"
// @Param cadFileId query int false "CAD文件ID"
// @Success 200 {object} response.Response
// @Router /api/designer/cad/render/records [get]
func ListRenderRecords(c *gin.Context) {
	// 从上下文获取用户ID / Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 32)
	cadFileID, _ := strconv.ParseUint(c.Query("cadFileId"), 10, 32)

	records, total, err := renderRecordService.ListRenderRecords(current, size, userID.(int64), uint(projectID), uint(cadFileID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"records": records,
		"total":   total,
		"current": current,
		"size":    size,
	})
}

// DeleteRenderRecord 删除效果图记录
// @Summary 删除效果图记录
// @Description 删除指定的效果图记录
// @Tags CAD渲染
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response
// @Router /api/designer/cad/render/records/{id} [delete]
func DeleteRenderRecord(c *gin.Context) {
	// 从上下文获取用户ID / Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	if err := renderRecordService.DeleteRenderRecord(uint(id), userID.(int64)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserRenderQuota 更新用户效果图配额（管理员）
// @Summary 更新用户效果图配额
// @Description 管理员更新指定用户的每日效果图保存限制
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body UpdateQuotaRequest true "更新请求"
// @Success 200 {object} response.Response
// @Router /api/system/user/render-quota [put]
func UpdateUserRenderQuota(c *gin.Context) {
	var req struct {
		UserID     int64 `json:"userId" binding:"required"`
		DailyLimit int   `json:"dailyLimit" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := renderRecordService.UpdateUserQuota(req.UserID, req.DailyLimit); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// ResetUserRenderQuota 重置用户今日使用次数（管理员）
// @Summary 重置用户今日使用次数
// @Description 管理员重置指定用户的今日效果图保存次数
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param userId query int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/system/user/render-quota/reset [post]
func ResetUserRenderQuota(c *gin.Context) {
	userIDStr := c.Query("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := renderRecordService.ResetUserQuota(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
