package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AnalyzeUserEmotions 分析用户情绪模式
// @Summary 分析用户情绪模式
// @Description 基于用户的情绪记录和日记数据进行AI情绪分析
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.EmotionAnalysisRequest true "情绪分析请求"
// @Success 200 {object} response.Response{data=service.EmotionAnalysisResult} "分析成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/emotions [post]
func AnalyzeUserEmotions(c *gin.Context) {
	var req request.EmotionAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// 验证日期范围
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		response.BadRequest(c, "开始日期和结束日期不能为空")
		return
	}

	if req.StartDate.After(req.EndDate) {
		response.BadRequest(c, "开始日期不能晚于结束日期")
		return
	}

	// 限制分析时间范围不超过90天
	if req.EndDate.Sub(req.StartDate).Hours() > 90*24 {
		response.BadRequest(c, "分析时间范围不能超过90天")
		return
	}

	// 验证分析类型
	validTypes := map[string]bool{"pattern": true, "trend": true, "insight": true}
	if !validTypes[req.AnalysisType] {
		response.BadRequest(c, "分析类型必须是: pattern, trend, insight")
		return
	}

	// 构造分析请求
	analysisReq := &service.EmotionAnalysisRequest{
		UserID:       userID,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		AnalysisType: req.AnalysisType,
	}

	// 执行AI分析
	result, err := aiAnalysisService.AnalyzeUserEmotions(c.Request.Context(), analysisReq)
	if err != nil {
		response.ServerError(c, "情绪分析失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetEmotionAnalysisHistory 获取情绪分析历史
// @Summary 获取情绪分析历史
// @Description 获取用户的历史情绪分析记录
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "当前页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param analysisType query string false "分析类型"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/history [get]
func GetEmotionAnalysisHistory(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	_ = c.Query("analysisType") // 暂时忽略此参数

	// TODO: 实现从数据库获取分析历史
	// 这里暂时返回空列表
	response.SuccessWithPagination(c, []interface{}{}, page, pageSize, 0)
}

// GetEmotionAnalysisDetail 获取情绪分析详情
// @Summary 获取情绪分析详情
// @Description 获取特定情绪分析的详细信息
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分析ID"
// @Success 200 {object} response.Response{data=service.EmotionAnalysisResult} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/{id} [get]
func GetEmotionAnalysisDetail(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	idStr := c.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: 实现从数据库获取分析详情
	// 这里暂时返回错误
	response.ServerError(c, "分析详情功能暂未实现")
}

// DeleteEmotionAnalysis 删除情绪分析记录
// @Summary 删除情绪分析记录
// @Description 删除特定的情绪分析记录
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分析ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/{id} [delete]
func DeleteEmotionAnalysis(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	idStr := c.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: 实现删除分析记录
	response.SuccessWithMsg(c, "删除分析记录成功", nil)
}

// ExportEmotionAnalysis 导出情绪分析报告
// @Summary 导出情绪分析报告
// @Description 导出情绪分析报告为PDF或JSON格式
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分析ID"
// @Param format query string false "导出格式" Enums(json,pdf) default(json)
// @Success 200 {object} response.Response "导出成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/{id}/export [get]
func ExportEmotionAnalysis(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	idStr := c.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "pdf" {
		response.BadRequest(c, "导出格式必须是json或pdf")
		return
	}

	// TODO: 实现导出功能
	response.ServerError(c, "导出功能暂未实现")
}

// BatchAnalyzeEmotions 批量分析情绪
// @Summary 批量分析情绪
// @Description 批量分析多个用户的情绪数据（管理员功能）
// @Tags AI分析
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.BatchEmotionAnalysisRequest true "批量分析请求"
// @Success 200 {object} response.Response "分析成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/ai-analysis/batch [post]
func BatchAnalyzeEmotions(c *gin.Context) {
	// 从JWT token中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "无效的用户认证")
		return
	}

	// TODO: 验证管理员权限
	// TODO: 实现批量分析功能
	response.ServerError(c, "批量分析功能暂未实现")
}
