package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	cadService "art_admin_backend/internal/service/cad"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var cadSvc *cadService.CadService

// SetCadService 设置CAD服务
// Set CAD service instance
func SetCadService(svc *cadService.CadService) {
	cadSvc = svc
}

// UploadCadFile 上传CAD文件
// @Summary 上传CAD文件
// @Description 上传DWG或DXF格式的CAD文件
// @Tags CAD图纸预览
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CAD文件(DWG/DXF)"
// @Param projectId formData int true "项目ID"
// @Success 200 {object} response.Response{data=response.CadFileUploadResponse}
// @Router /api/designer/cad/upload [post]
func UploadCadFile(c *gin.Context) {
	// 获取上传的文件 / Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的CAD文件")
		return
	}

	// 限制文件大小 100MB / Limit file size to 100MB
	if file.Size > 100*1024*1024 {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过100MB")
		return
	}

	// 获取项目ID / Get project ID
	projectIDStr := c.PostForm("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil || projectID == 0 {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 上传文件 / Upload file
	result, err := cadSvc.Upload(file, uint(projectID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ParseCadFile 解析CAD文件
// @Summary 解析CAD文件
// @Description 解析已上传的CAD文件，提取图层和实体信息
// @Tags CAD图纸预览
// @Accept json
// @Produce json
// @Param request body request.ParseCadFileRequest true "解析请求"
// @Success 200 {object} response.Response{data=response.CadParseResult}
// @Router /api/designer/cad/parse [post]
func ParseCadFile(c *gin.Context) {
	var req request.ParseCadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := cadSvc.Parse(req.CadFileID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCadFileList 获取CAD文件列表
// @Summary 获取CAD文件列表
// @Description 获取项目下的CAD文件列表
// @Tags CAD图纸预览
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int false "当前页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param parseStatus query string false "解析状态筛选(pending/processing/completed/failed)"
// @Param fileFormat query string false "文件格式筛选(dwg/dxf)"
// @Success 200 {object} response.Response{data=response.CadFileListResponse}
// @Router /api/designer/cad [get]
func GetCadFileList(c *gin.Context) {
	var req request.CadFileListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := cadSvc.List(req.ProjectID, req.Current, req.Size, req.ParseStatus, req.FileFormat)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCadFileDetail 获取CAD文件详情
// @Summary 获取CAD文件详情
// @Description 获取CAD文件的详细信息
// @Tags CAD图纸预览
// @Produce json
// @Param id path int true "CAD文件ID"
// @Success 200 {object} response.Response{data=response.CadFileResponse}
// @Router /api/designer/cad/{id} [get]
func GetCadFileDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的CAD文件ID")
		return
	}

	result, err := cadSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "CAD文件不存在")
		return
	}

	response.Success(c, result)
}

// GetCadParseResult 获取CAD解析结果
// @Summary 获取CAD解析结果
// @Description 获取CAD文件的解析结果，包含图层和实体信息
// @Tags CAD图纸预览
// @Produce json
// @Param id path int true "CAD文件ID"
// @Success 200 {object} response.Response{data=response.CadParseResult}
// @Router /api/designer/cad/{id}/parse [get]
func GetCadParseResult(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的CAD文件ID")
		return
	}

	result, err := cadSvc.GetParseResult(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCadLayers 获取CAD图层列表
// @Summary 获取CAD图层列表
// @Description 获取CAD文件的图层列表
// @Tags CAD图纸预览
// @Produce json
// @Param id path int true "CAD文件ID"
// @Success 200 {object} response.Response{data=response.CadLayerListResponse}
// @Router /api/designer/cad/{id}/layers [get]
func GetCadLayers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的CAD文件ID")
		return
	}

	result, err := cadSvc.GetLayers(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteCadFile 删除CAD文件
// @Summary 删除CAD文件
// @Description 删除指定的CAD文件
// @Tags CAD图纸预览
// @Produce json
// @Param id path int true "CAD文件ID"
// @Success 200 {object} response.Response
// @Router /api/designer/cad/{id} [delete]
func DeleteCadFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的CAD文件ID")
		return
	}

	if err := cadSvc.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDeleteCadFiles 批量删除CAD文件
// @Summary 批量删除CAD文件
// @Description 批量删除多个CAD文件
// @Tags CAD图纸预览
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteCadFileRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/cad/batch-delete [post]
func BatchDeleteCadFiles(c *gin.Context) {
	var req request.BatchDeleteCadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := cadSvc.BatchDelete(req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
