package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	documentService "art_admin_backend/internal/service/document"
	"strconv"

	"github.com/gin-gonic/gin"
)

var documentSvc *documentService.DocumentService

// SetDocumentService 设置文档服务
// Set document service
func SetDocumentService(svc *documentService.DocumentService) {
	documentSvc = svc
}

// UploadDocument 上传文档
// @Summary 上传项目文档
// @Description 上传项目文档，支持PDF/Word/图片格式
// @Tags 文档分析
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文档文件"
// @Param projectId formData int true "项目ID"
// @Success 200 {object} response.Response{data=response.DocumentUploadResponse}
// @Router /api/designer/documents/upload [post]
func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	projectIDStr := c.PostForm("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil || projectID == 0 {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := documentSvc.Upload(file, uint(projectID), userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "上传成功", result)
}

// AnalyzeDocument 分析文档
// @Summary 分析项目文档
// @Description 使用AI分析文档内容，提取关键信息
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param request body request.AnalyzeDocumentRequest true "分析请求"
// @Success 200 {object} response.Response{data=response.DocumentAnalysisResult}
// @Router /api/designer/documents/analyze [post]
func AnalyzeDocument(c *gin.Context) {
	var req request.AnalyzeDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := documentSvc.Analyze(req.DocumentID, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "分析完成", result)
}

// GetDocumentList 获取文档列表
// @Summary 获取项目文档列表
// @Description 分页获取指定项目的文档列表
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int true "当前页码" default(1)
// @Param size query int true "每页数量" default(10)
// @Param analysisStatus query string false "分析状态筛选"
// @Param fileType query string false "文件类型筛选"
// @Success 200 {object} response.Response{data=response.DocumentListResponse}
// @Router /api/designer/documents [get]
func GetDocumentList(c *gin.Context) {
	var req request.DocumentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := documentSvc.List(&req, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDocumentDetail 获取文档详情
// @Summary 获取文档详情
// @Description 获取指定文档的详细信息和分析结果
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param id path int true "文档ID"
// @Success 200 {object} response.Response{data=response.DocumentResponse}
// @Router /api/designer/documents/{id} [get]
func GetDocumentDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := documentSvc.GetByID(uint(id), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDocumentKeywords 获取文档关键字
// @Summary 获取文档关键字
// @Description 获取已分析文档的关键字列表
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param id path int true "文档ID"
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/designer/documents/{id}/keywords [get]
func GetDocumentKeywords(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	keywords, err := documentSvc.GetKeywords(uint(id), userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, keywords)
}

// GetDocumentSummary 获取文档摘要
// @Summary 获取文档摘要
// @Description 获取已分析文档的摘要内容
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param id path int true "文档ID"
// @Success 200 {object} response.Response{data=response.DocumentSummaryResponse}
// @Router /api/designer/documents/{id}/summary [get]
func GetDocumentSummary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	summary, err := documentSvc.GetSummary(uint(id), userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, summary)
}

// DeleteDocument 删除文档
// @Summary 删除文档
// @Description 删除指定的项目文档
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param id path int true "文档ID"
// @Success 200 {object} response.Response
// @Router /api/designer/documents/{id} [delete]
func DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := documentSvc.Delete(uint(id), userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteDocuments 批量删除文档
// @Summary 批量删除文档
// @Description 批量删除多个项目文档
// @Tags 文档分析
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteDocumentRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/documents/batch-delete [post]
func BatchDeleteDocuments(c *gin.Context) {
	var req request.BatchDeleteDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := documentSvc.BatchDelete(req.IDs, userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}
