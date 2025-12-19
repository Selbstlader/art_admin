package v1

import (
	"strconv"

	"art_admin_backend/internal/dto/request"
	dtoResponse "art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/response"
	designCompareService "art_admin_backend/internal/service/design_compare"

	"github.com/gin-gonic/gin"
)

var designCompareSvc *designCompareService.DesignCompareService

// SetDesignCompareService 设置设计比对服务
// Set design compare service
func SetDesignCompareService(svc *designCompareService.DesignCompareService) {
	designCompareSvc = svc
}

// UploadDesignImages 批量上传设计图（支持多文件、CAD）
// @Summary 批量上传设计图
// @Description 批量上传设计图用于与多份需求文档进行比对分析，支持图片、PDF、CAD文件
// @Tags 设计比对
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "设计图文件（支持多个）"
// @Param projectId formData int true "项目ID"
// @Param documentIds formData string true "需求文档ID列表（逗号分隔）"
// @Param name formData string false "比对任务名称"
// @Success 200 {object} response.Response{data=response.DesignCompareUploadResponse}
// @Router /api/designer/compare/upload [post]
func UploadDesignImages(c *gin.Context) {
	// 获取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "请选择要上传的设计文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		// 兼容单文件上传
		files = form.File["file"]
	}
	if len(files) == 0 {
		response.BadRequest(c, "请选择要上传的设计文件")
		return
	}

	var req request.UploadDesignImagesRequest
	req.Name = c.PostForm("name")
	req.DocumentIDs = c.PostForm("documentIds")

	projectIDStr := c.PostForm("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil || projectID == 0 {
		response.BadRequest(c, "无效的项目ID")
		return
	}
	req.ProjectID = uint(projectID)

	if req.DocumentIDs == "" {
		response.BadRequest(c, "请选择需求文档")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designCompareSvc.Upload(files, &req, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "上传成功", result)
}

// AnalyzeDesignCompare 分析设计比对
// @Summary 分析设计比对
// @Description 使用AI分析设计图与需求文档的匹配程度
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param request body request.AnalyzeDesignCompareRequest true "分析请求"
// @Success 200 {object} response.Response{data=response.DesignCompareAnalysisResult}
// @Router /api/designer/compare/analyze [post]
func AnalyzeDesignCompare(c *gin.Context) {
	var req request.AnalyzeDesignCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designCompareSvc.Analyze(req.CompareID, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "分析完成", result)
}

// GetDesignCompareList 获取设计比对列表
// @Summary 获取设计比对列表
// @Description 分页获取指定项目的设计比对记录列表
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int true "当前页码" default(1)
// @Param size query int true "每页数量" default(10)
// @Param analysisStatus query string false "分析状态筛选"
// @Success 200 {object} response.Response{data=response.DesignCompareListResponse}
// @Router /api/designer/compare [get]
func GetDesignCompareList(c *gin.Context) {
	var req request.DesignCompareListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designCompareSvc.List(&req, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignCompareDetail 获取设计比对详情
// @Summary 获取设计比对详情
// @Description 获取指定设计比对记录的详细信息和分析结果
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param id path int true "比对记录ID"
// @Success 200 {object} response.Response{data=response.DesignCompareResponse}
// @Router /api/designer/compare/{id} [get]
func GetDesignCompareDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的比对记录ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designCompareSvc.GetByID(uint(id), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignCompareResult 获取设计比对结果
// @Summary 获取设计比对结果
// @Description 获取指定设计比对记录的分析结果
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param id path int true "比对记录ID"
// @Success 200 {object} response.Response{data=response.DesignCompareAnalysisResult}
// @Router /api/designer/compare/{id}/result [get]
func GetDesignCompareResult(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的比对记录ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designCompareSvc.GetByID(uint(id), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// 转换为分析结果格式 / Convert to analysis result format
	analysisResult := &struct {
		CompareID      uint                                 `json:"compareId"`
		MatchItems     []dtoResponse.MatchItemResponse      `json:"matchItems"`
		DeviationItems []dtoResponse.DeviationItemResponse  `json:"deviationItems"`
		Suggestions    []dtoResponse.SuggestionItemResponse `json:"suggestions"`
		OverallScore   float64                              `json:"overallScore"`
	}{
		CompareID:      result.ID,
		MatchItems:     result.MatchItems,
		DeviationItems: result.DeviationItems,
		Suggestions:    result.Suggestions,
		OverallScore:   result.OverallScore,
	}

	response.Success(c, analysisResult)
}

// DeleteDesignCompare 删除设计比对记录
// @Summary 删除设计比对记录
// @Description 删除指定的设计比对记录
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param id path int true "比对记录ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compare/{id} [delete]
func DeleteDesignCompare(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的比对记录ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := designCompareSvc.Delete(uint(id), userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteDesignCompare 批量删除设计比对记录
// @Summary 批量删除设计比对记录
// @Description 批量删除多个设计比对记录
// @Tags 设计比对
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteDesignCompareRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compare/batch-delete [post]
func BatchDeleteDesignCompare(c *gin.Context) {
	var req request.BatchDeleteDesignCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := designCompareSvc.BatchDelete(req.IDs, userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}
