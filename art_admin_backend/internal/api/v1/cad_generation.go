package v1

import (
	"art_admin_backend/internal/dto/request"
	dtoResponse "art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/cad_generation"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CadGenerationAPI CAD生成API控制器
// CAD generation API controller
type CadGenerationAPI struct {
	service *cad_generation.CadGenerationService
}

var cadGenerationAPI *CadGenerationAPI

// initCadGenerationAPI 初始化CAD生成API
// Initialize CAD generation API
func initCadGenerationAPI() {
	if cadGenerationAPI == nil {
		db := database.GetDB()
		cadGenerationAPI = &CadGenerationAPI{
			service: cad_generation.NewCadGenerationService(db),
		}
	}
}

// NewCadGenerationAPI 创建CAD生成API实例
// Create CAD generation API instance
func NewCadGenerationAPI() *CadGenerationAPI {
	initCadGenerationAPI()
	return cadGenerationAPI
}

// Create 创建CAD生成任务
// @Summary 创建CAD生成任务
// @Description 基于项目文档AI生成CAD文件
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param data body request.CreateCadGenerationReq true "创建请求"
// @Success 200 {object} dtoResponse.CadGenerationResp
// @Router /api/designer/cad-generations [post]
func (api *CadGenerationAPI) Create(c *gin.Context) {
	var req request.CreateCadGenerationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID / Get current user ID
	userID := getCadGenUserID(c)

	resp, err := api.service.Create(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// getCadGenUserID 从上下文获取用户ID
// Get user ID from context
func getCadGenUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// GetByID 获取任务详情
// @Summary 获取CAD生成任务详情
// @Description 根据ID获取CAD生成任务详情
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} dtoResponse.CadGenerationResp
// @Router /api/designer/cad-generations/{id} [get]
func (api *CadGenerationAPI) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	resp, err := api.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// List 获取任务列表
// @Summary 获取CAD生成任务列表
// @Description 分页获取CAD生成任务列表
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param projectId query int false "项目ID"
// @Param generationType query string false "生成类型"
// @Param status query string false "状态"
// @Param current query int true "当前页"
// @Param size query int true "每页数量"
// @Success 200 {object} dtoResponse.CadGenerationListResp
// @Router /api/designer/cad-generations [get]
func (api *CadGenerationAPI) List(c *gin.Context) {
	var req request.CadGenerationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := api.service.List(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// Delete 删除任务
// @Summary 删除CAD生成任务
// @Description 根据ID删除CAD生成任务
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/designer/cad-generations/{id} [delete]
func (api *CadGenerationAPI) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	if err := api.service.Delete(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, "删除成功")
}

// Retry 重试生成任务
// @Summary 重试CAD生成任务
// @Description 重试失败的CAD生成任务
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} dtoResponse.CadGenerationResp
// @Router /api/designer/cad-generations/{id}/retry [post]
func (api *CadGenerationAPI) Retry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的ID")
		return
	}

	resp, err := api.service.Retry(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// Confirm 确认并保存CAD文件
// @Summary 确认并保存CAD文件
// @Description 确认生成的CAD文件并保存到项目
// @Tags CAD生成
// @Accept json
// @Produce json
// @Param data body request.ConfirmCadFileReq true "确认请求"
// @Success 200 {object} dtoResponse.CadGenerationResp
// @Router /api/designer/cad-generations/confirm [post]
func (api *CadGenerationAPI) Confirm(c *gin.Context) {
	var req request.ConfirmCadFileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getCadGenUserID(c)

	resp, err := api.service.ConfirmAndSave(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetTypeOptions 获取生成类型选项
// @Summary 获取CAD生成类型选项
// @Description 获取所有可用的CAD生成类型
// @Tags CAD生成
// @Accept json
// @Produce json
// @Success 200 {object} []dtoResponse.CadGenerationTypeOption
// @Router /api/designer/cad-generations/type-options [get]
func (api *CadGenerationAPI) GetTypeOptions(c *gin.Context) {
	response.Success(c, dtoResponse.GetGenerationTypeOptions())
}
