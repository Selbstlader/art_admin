package v1

import (
	"net/http"
	"strconv"
	"strings"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/construction_annotation"

	"github.com/gin-gonic/gin"
)

// ConstructionAnnotationController 施工图标注控制器
// Construction annotation controller
type ConstructionAnnotationController struct {
	annotationService *construction_annotation.AnnotationService
}

// NewConstructionAnnotationController 创建施工图标注控制器
// Create construction annotation controller
func NewConstructionAnnotationController(annotationService *construction_annotation.AnnotationService) *ConstructionAnnotationController {
	return &ConstructionAnnotationController{
		annotationService: annotationService,
	}
}

// AnalyzeConstructionDrawing 分析施工图
// @Summary 分析施工图并识别元素
// @Description 上传施工图进行AI分析，识别尺寸线、材料区域、设备位置等元素
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param request body request.AnalyzeConstructionDrawingRequest true "分析请求"
// @Success 200 {object} response.Response "分析成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/designer/construction-annotation/analyze [post]
func (c *ConstructionAnnotationController) AnalyzeConstructionDrawing(ctx *gin.Context) {
	var req request.AnalyzeConstructionDrawingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 验证图片格式 / Validate image format
	if err := c.annotationService.ValidateImageFormat(req.ImagePath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 开始分析 / Start analysis
	annotation, err := c.annotationService.AnalyzeConstructionDrawing(req.ProjectID, req.CadFileID, req.ImagePath)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "分析失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "分析已开始，请稍后查看结果", annotation)
}

// GetAnnotation 获取标注详情
// @Summary 获取标注详情
// @Description 根据ID获取施工图标注的详细信息
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id} [get]
func (c *ConstructionAnnotationController) GetAnnotation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	annotation, err := c.annotationService.GetAnnotationByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "标注不存在")
		return
	}

	response.Success(ctx, annotation)
}

// GetAnnotationsList 获取标注列表
// @Summary 获取标注列表
// @Description 分页获取施工图标注列表，支持按项目和CAD文件筛选
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param projectId query int false "项目ID"
// @Param cadFileId query int false "CAD文件ID"
// @Param status query string false "分析状态"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/designer/construction-annotation/list [get]
func (c *ConstructionAnnotationController) GetAnnotationsList(ctx *gin.Context) {
	var req request.GetAnnotationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	req.Validate()

	var annotations []*model.ConstructionAnnotation
	var total int64
	var err error

	// 根据筛选条件获取数据 / Get data based on filter conditions
	if req.ProjectID != nil {
		annotations, err = c.annotationService.GetAnnotationsByProject(*req.ProjectID)
		total = int64(len(annotations))
		// 手动分页 / Manual pagination
		start := req.GetOffset()
		end := start + req.PageSize
		if start > len(annotations) {
			annotations = []*model.ConstructionAnnotation{}
		} else if end > len(annotations) {
			annotations = annotations[start:]
		} else {
			annotations = annotations[start:end]
		}
	} else if req.CadFileID != nil {
		annotations, err = c.annotationService.GetAnnotationsByCadFile(*req.CadFileID)
		total = int64(len(annotations))
		// 手动分页 / Manual pagination
		start := req.GetOffset()
		end := start + req.PageSize
		if start > len(annotations) {
			annotations = []*model.ConstructionAnnotation{}
		} else if end > len(annotations) {
			annotations = annotations[start:]
		} else {
			annotations = annotations[start:end]
		}
	} else {
		// 没有筛选条件时，获取所有标注（分页）/ Get all annotations when no filter
		annotations, total, err = c.annotationService.GetAnnotationsList(req.GetOffset(), req.PageSize)
	}

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "获取标注列表失败: "+err.Error())
		return
	}

	response.SuccessWithPagination(ctx, annotations, req.Page, req.PageSize, total)
}

// UpdateAnnotations 更新标注
// @Summary 更新标注数据
// @Description 更新施工图的标注信息，支持增删改标注项
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param request body request.UpdateAnnotationsRequest true "更新请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id}/annotations [put]
func (c *ConstructionAnnotationController) UpdateAnnotations(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	var req request.UpdateAnnotationsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 检查标注是否存在 / Check if annotation exists
	_, err = c.annotationService.GetAnnotationByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "标注不存在")
		return
	}

	// 更新标注 / Update annotations
	if err := c.annotationService.UpdateAnnotations(uint(id), req.Annotations); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "更新标注失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// DeleteAnnotation 删除标注
// @Summary 删除标注
// @Description 删除指定的施工图标注记录
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id} [delete]
func (c *ConstructionAnnotationController) DeleteAnnotation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	// 检查标注是否存在 / Check if annotation exists
	_, err = c.annotationService.GetAnnotationByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "标注不存在")
		return
	}

	// 删除标注 / Delete annotation
	if err := c.annotationService.DeleteAnnotation(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "删除标注失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// AddAnnotationItem 添加标注项
// @Summary 添加标注项
// @Description 向施工图标注中添加新的标注项
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param request body request.CreateAnnotationItemRequest true "标注项数据"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id}/items [post]
func (c *ConstructionAnnotationController) AddAnnotationItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	var req request.CreateAnnotationItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 构建标注项 / Build annotation item
	item := model.AnnotationItem{
		Type:       req.Type,
		Position:   req.Position,
		Content:    req.Content,
		Style:      req.Style,
		Properties: req.Properties,
	}

	// 添加标注项 / Add annotation item
	if err := c.annotationService.AddAnnotationItem(uint(id), item); err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.Error(ctx, http.StatusNotFound, err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "添加标注项失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMsg(ctx, "添加成功", nil)
}

// UpdateAnnotationItem 更新标注项
// @Summary 更新标注项
// @Description 更新施工图标注中的指定标注项
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param itemId path string true "标注项ID"
// @Param request body request.CreateAnnotationItemRequest true "标注项数据"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注或标注项不存在"
// @Router /api/designer/construction-annotation/{id}/items/{itemId} [put]
func (c *ConstructionAnnotationController) UpdateAnnotationItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	itemID := ctx.Param("itemId")
	if itemID == "" {
		response.Error(ctx, http.StatusBadRequest, "标注项ID不能为空")
		return
	}

	var req request.CreateAnnotationItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 构建标注项 / Build annotation item
	item := model.AnnotationItem{
		Type:       req.Type,
		Position:   req.Position,
		Content:    req.Content,
		Style:      req.Style,
		Properties: req.Properties,
	}

	// 更新标注项 / Update annotation item
	if err := c.annotationService.UpdateAnnotationItem(uint(id), itemID, item); err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.Error(ctx, http.StatusNotFound, err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "更新标注项失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// DeleteAnnotationItem 删除标注项
// @Summary 删除标注项
// @Description 删除施工图标注中的指定标注项
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param itemId path string true "标注项ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注或标注项不存在"
// @Router /api/designer/construction-annotation/{id}/items/{itemId} [delete]
func (c *ConstructionAnnotationController) DeleteAnnotationItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	itemID := ctx.Param("itemId")
	if itemID == "" {
		response.Error(ctx, http.StatusBadRequest, "标注项ID不能为空")
		return
	}

	// 删除标注项 / Delete annotation item
	if err := c.annotationService.DeleteAnnotationItem(uint(id), itemID); err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.Error(ctx, http.StatusNotFound, err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "删除标注项失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetAnnotationItem 获取标注项
// @Summary 获取标注项
// @Description 获取施工图标注中的指定标注项详情
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param itemId path string true "标注项ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注或标注项不存在"
// @Router /api/designer/construction-annotation/{id}/items/{itemId} [get]
func (c *ConstructionAnnotationController) GetAnnotationItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	itemID := ctx.Param("itemId")
	if itemID == "" {
		response.Error(ctx, http.StatusBadRequest, "标注项ID不能为空")
		return
	}

	// 获取标注项 / Get annotation item
	item, err := c.annotationService.GetAnnotationItem(uint(id), itemID)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.Error(ctx, http.StatusNotFound, err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "获取标注项失败: "+err.Error())
		}
		return
	}

	response.Success(ctx, item)
}

// BatchUpdateAnnotations 批量更新标注
// @Summary 批量更新标注
// @Description 批量执行标注项的增删改操作
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Param request body []construction_annotation.AnnotationOperation true "批量操作数据"
// @Success 200 {object} response.Response "操作成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id}/batch [post]
func (c *ConstructionAnnotationController) BatchUpdateAnnotations(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	var operations []construction_annotation.AnnotationOperation
	if err := ctx.ShouldBindJSON(&operations); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 批量更新标注 / Batch update annotations
	if err := c.annotationService.BatchUpdateAnnotations(uint(id), operations); err != nil {
		if strings.Contains(err.Error(), "不存在") {
			response.Error(ctx, http.StatusNotFound, err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "批量更新失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMsg(ctx, "批量操作成功", nil)
}

// ExportAnnotation 导出标注
// @Summary 导出标注
// @Description 将施工图标注导出为PDF或图片格式
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param request body request.ExportAnnotationRequest true "导出请求"
// @Success 200 {object} response.Response "导出成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/export [post]
func (c *ConstructionAnnotationController) ExportAnnotation(ctx *gin.Context) {
	var req request.ExportAnnotationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 验证导出配置 / Validate export configuration
	if err := c.annotationService.ValidateExportConfig(req.Config); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 检查标注是否存在 / Check if annotation exists
	_, err := c.annotationService.GetAnnotationByID(req.AnnotationID)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "标注不存在")
		return
	}

	// 导出标注 / Export annotation
	result, err := c.annotationService.ExportAnnotation(req.AnnotationID, req.Config)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "导出失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "导出成功", result)
}

// GetExportFormats 获取支持的导出格式
// @Summary 获取支持的导出格式
// @Description 获取系统支持的标注导出格式列表
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "获取成功"
// @Router /api/designer/construction-annotation/export-formats [get]
func (c *ConstructionAnnotationController) GetExportFormats(ctx *gin.Context) {
	formats := c.annotationService.GetExportFormats()
	response.Success(ctx, formats)
}

// GetExportHistory 获取导出历史
// @Summary 获取导出历史
// @Description 获取指定标注的导出历史记录
// @Tags 施工图标注
// @Accept json
// @Produce json
// @Param id path int true "标注ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "标注不存在"
// @Router /api/designer/construction-annotation/{id}/export-history [get]
func (c *ConstructionAnnotationController) GetExportHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "无效的标注ID")
		return
	}

	// 检查标注是否存在 / Check if annotation exists
	_, err = c.annotationService.GetAnnotationByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "标注不存在")
		return
	}

	// 获取导出历史 / Get export history
	history, err := c.annotationService.GetExportHistory(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "获取导出历史失败: "+err.Error())
		return
	}

	response.Success(ctx, history)
}
