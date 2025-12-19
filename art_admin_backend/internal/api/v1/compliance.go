package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service/compliance"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ComplianceAPI 合规检查API
// Compliance check API controller
type ComplianceAPI struct {
	service *compliance.ComplianceService
}

// NewComplianceAPI 创建合规检查API
// Create compliance API
func NewComplianceAPI() *ComplianceAPI {
	db := database.GetDB()
	repo := repository.NewComplianceRepository(db)
	projectRepo := repository.NewDesignerProjectRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	aiClient := createVolcengineClient() // 使用已有的函数 / Use existing function

	service := compliance.NewComplianceService(repo, projectRepo, docRepo, aiClient)
	return &ComplianceAPI{service: service}
}

// ========== 设计规范管理接口 / Design Standard Management APIs ==========

// CreateDesignStandard 创建设计规范
// @Summary 创建设计规范
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param request body request.CreateDesignStandardRequest true "创建规范请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards [post]
func (api *ComplianceAPI) CreateDesignStandard(c *gin.Context) {
	var req request.CreateDesignStandardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	result, err := api.service.CreateStandard(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateDesignStandard 更新设计规范
// @Summary 更新设计规范
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param request body request.UpdateDesignStandardRequest true "更新规范请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards [put]
func (api *ComplianceAPI) UpdateDesignStandard(c *gin.Context) {
	var req request.UpdateDesignStandardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	result, err := api.service.UpdateStandard(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignStandardList 获取设计规范列表
// @Summary 获取设计规范列表
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param current query int true "当前页码"
// @Param size query int true "每页数量"
// @Param category query string false "类别筛选"
// @Param status query string false "状态筛选"
// @Param keyword query string false "关键字搜索"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards [get]
func (api *ComplianceAPI) GetDesignStandardList(c *gin.Context) {
	var req request.DesignStandardListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	result, err := api.service.ListStandards(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignStandardDetail 获取设计规范详情
// @Summary 获取设计规范详情
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param id path int true "规范ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards/{id} [get]
func (api *ComplianceAPI) GetDesignStandardDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的规范ID")
		return
	}

	result, err := api.service.GetStandardByID(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteDesignStandard 删除设计规范
// @Summary 删除设计规范
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param id path int true "规范ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards/{id} [delete]
func (api *ComplianceAPI) DeleteDesignStandard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的规范ID")
		return
	}

	if err := api.service.DeleteStandard(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteDesignStandards 批量删除设计规范
// @Summary 批量删除设计规范
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteDesignStandardRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards/batch-delete [post]
func (api *ComplianceAPI) BatchDeleteDesignStandards(c *gin.Context) {
	var req request.BatchDeleteDesignStandardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := api.service.BatchDeleteStandards(req.IDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// GetDesignStandardStats 获取设计规范统计
// @Summary 获取设计规范统计
// @Tags 合规检查
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/standards/stats [get]
func (api *ComplianceAPI) GetDesignStandardStats(c *gin.Context) {
	result, err := api.service.GetStandardStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// ========== 合规检查接口 / Compliance Check APIs ==========

// PerformComplianceCheck 执行合规检查
// @Summary 执行合规检查
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param request body request.ComplianceCheckRequest true "合规检查请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/check [post]
func (api *ComplianceAPI) PerformComplianceCheck(c *gin.Context) {
	var req request.ComplianceCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	result, err := api.service.PerformComplianceCheck(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetComplianceCheckList 获取合规检查结果列表
// @Summary 获取合规检查结果列表
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Param current query int true "当前页码"
// @Param size query int true "每页数量"
// @Param checkStatus query string false "检查状态筛选"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/results [get]
func (api *ComplianceAPI) GetComplianceCheckList(c *gin.Context) {
	var req request.ComplianceCheckListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	result, err := api.service.ListCheckResults(&req, userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetComplianceCheckDetail 获取合规检查结果详情
// @Summary 获取合规检查结果详情
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param id path int true "检查结果ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/results/{id} [get]
func (api *ComplianceAPI) GetComplianceCheckDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的检查结果ID")
		return
	}

	userID := getUserID(c)
	result, err := api.service.GetCheckResultByID(uint(id), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetLatestComplianceCheck 获取项目最新合规检查结果
// @Summary 获取项目最新合规检查结果
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/latest [get]
func (api *ComplianceAPI) GetLatestComplianceCheck(c *gin.Context) {
	projectIDStr := c.Query("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	userID := getUserID(c)
	result, err := api.service.GetLatestCheckResult(uint(projectID), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// GetComplianceCheckSummary 获取项目合规检查摘要
// @Summary 获取项目合规检查摘要
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/summary [get]
func (api *ComplianceAPI) GetComplianceCheckSummary(c *gin.Context) {
	projectIDStr := c.Query("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的项目ID")
		return
	}

	userID := getUserID(c)
	result, err := api.service.GetCheckSummary(uint(projectID), userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteComplianceCheck 删除合规检查结果
// @Summary 删除合规检查结果
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param id path int true "检查结果ID"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/results/{id} [delete]
func (api *ComplianceAPI) DeleteComplianceCheck(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的检查结果ID")
		return
	}

	userID := getUserID(c)
	if err := api.service.DeleteCheckResult(uint(id), userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteComplianceChecks 批量删除合规检查结果
// @Summary 批量删除合规检查结果
// @Tags 合规检查
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteComplianceCheckRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/compliance/results/batch-delete [post]
func (api *ComplianceAPI) BatchDeleteComplianceChecks(c *gin.Context) {
	var req request.BatchDeleteComplianceCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)
	if err := api.service.BatchDeleteCheckResults(req.IDs, userID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}
