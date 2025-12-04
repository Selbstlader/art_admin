package v1

import (
	"art_admin_backend/internal/pkg/response"
	fileService "art_admin_backend/internal/service/file"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var fileSvc *fileService.FileService

// SetFileService 设置文件服务
func SetFileService(svc *fileService.FileService) {
	fileSvc = svc
}

// UploadFile 上传文件
// @Summary 上传文件
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "文件"
// @Param category formData string false "分类(image/document/video/audio/other)"
// @Success 200 {object} response.Response
// @Router /api/file/upload [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	// 限制文件大小 50MB
	if file.Size > 50*1024*1024 {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过50MB")
		return
	}

	category := c.PostForm("category")
	userID := c.GetUint("userID")

	result, err := fileSvc.Upload(file, userID, category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UploadMultipleFiles 批量上传文件
// @Summary 批量上传文件
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件列表"
// @Param category formData string false "分类"
// @Success 200 {object} response.Response
// @Router /api/file/upload-multiple [post]
func UploadMultipleFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	category := c.PostForm("category")
	userID := c.GetUint("userID")

	var results []interface{}
	var errors []string

	for _, file := range files {
		if file.Size > 50*1024*1024 {
			errors = append(errors, file.Filename+": 文件大小超过50MB")
			continue
		}

		result, err := fileSvc.Upload(file, userID, category)
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}
		results = append(results, result)
	}

	response.Success(c, gin.H{
		"success": results,
		"errors":  errors,
	})
}

// GetFileList 获取文件列表
// @Summary 获取文件列表
// @Tags 文件管理
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param category query string false "分类"
// @Success 200 {object} response.Response
// @Router /api/file/list [get]
func GetFileList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	userID := c.GetUint("userID")

	files, total, err := fileSvc.List(page, pageSize, category, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":     files,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetFileDetail 获取文件详情
// @Summary 获取文件详情
// @Tags 文件管理
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response
// @Router /api/file/{id} [get]
func GetFileDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	file, err := fileSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	response.Success(c, file)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Tags 文件管理
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} response.Response
// @Router /api/file/{id} [delete]
func DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	if err := fileSvc.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchDeleteFiles 批量删除文件
// @Summary 批量删除文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param ids body []uint true "文件ID列表"
// @Success 200 {object} response.Response
// @Router /api/file/batch-delete [post]
func BatchDeleteFiles(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := fileSvc.BatchDelete(req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
