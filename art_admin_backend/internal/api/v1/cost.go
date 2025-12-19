package v1

import (
	"net/http"
	"strconv"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/cost"

	"github.com/gin-gonic/gin"
)

// CostAPI 成本估算API
// Cost estimate API controller
type CostAPI struct {
	costService       cost.CostService
	costExportService cost.CostExportService
}

// NewCostAPI 创建成本估算API实例
// Create cost estimate API instance
func NewCostAPI() *CostAPI {
	return &CostAPI{
		costService:       cost.NewCostService(),
		costExportService: cost.NewCostExportService(),
	}
}

// CreateCostEstimate 创建成本估算
// @Summary 创建成本估算
// @Description 为项目创建成本估算
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param request body request.CreateCostEstimateRequest true "创建成本估算请求"
// @Success 200 {object} response.Response{data=response.CostEstimateResponse}
// @Router /api/designer/cost [post]
func (api *CostAPI) CreateCostEstimate(c *gin.Context) {
	var req request.CreateCostEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.costService.CreateCostEstimate(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateCostEstimate 更新成本估算
// @Summary 更新成本估算
// @Description 更新项目的成本估算
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param request body request.UpdateCostEstimateRequest true "更新成本估算请求"
// @Success 200 {object} response.Response{data=response.CostEstimateResponse}
// @Router /api/designer/cost [put]
func (api *CostAPI) UpdateCostEstimate(c *gin.Context) {
	var req request.UpdateCostEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.costService.UpdateCostEstimate(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteCostEstimate 删除成本估算
// @Summary 删除成本估算
// @Description 删除指定的成本估算
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param id path int true "成本估算ID"
// @Success 200 {object} response.Response
// @Router /api/designer/cost/{id} [delete]
func (api *CostAPI) DeleteCostEstimate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID / Invalid ID")
		return
	}

	if err := api.costService.DeleteCostEstimate(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetCostEstimateByID 根据ID获取成本估算
// @Summary 根据ID获取成本估算
// @Description 获取指定ID的成本估算详情
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param id path int true "成本估算ID"
// @Success 200 {object} response.Response{data=response.CostEstimateResponse}
// @Router /api/designer/cost/{id} [get]
func (api *CostAPI) GetCostEstimateByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID / Invalid ID")
		return
	}

	result, err := api.costService.GetCostEstimateByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCostEstimateByProjectID 根据项目ID获取成本估算
// @Summary 根据项目ID获取成本估算
// @Description 获取指定项目的成本估算详情
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param projectId query int true "项目ID"
// @Success 200 {object} response.Response{data=response.CostEstimateResponse}
// @Router /api/designer/cost/project [get]
func (api *CostAPI) GetCostEstimateByProjectID(c *gin.Context) {
	projectIDStr := c.Query("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	result, err := api.costService.GetCostEstimateByProjectID(uint(projectID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, result)
}

// CalculateCost 计算成本(不保存)
// @Summary 计算成本
// @Description 根据参数计算成本，不保存到数据库
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param request body request.CalculateCostRequest true "计算成本请求"
// @Success 200 {object} response.Response{data=response.CostCalculateResponse}
// @Router /api/designer/cost/calculate [post]
func (api *CostAPI) CalculateCost(c *gin.Context) {
	var req request.CalculateCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.costService.CalculateCost(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCostSummary 获取成本汇总
// @Summary 获取成本汇总
// @Description 获取用户的成本汇总统计
// @Tags 成本估算
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.CostSummaryResponse}
// @Router /api/designer/cost/summary [get]
func (api *CostAPI) GetCostSummary(c *gin.Context) {
	// 从上下文获取用户ID / Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		userID = uint(0)
	}

	result, err := api.costService.GetCostSummary(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ExportCostReport 导出成本报告
// @Summary 导出成本报告
// @Description 导出项目的成本报告为Excel或PDF格式
// @Tags 成本估算
// @Accept json
// @Produce json
// @Param request body request.ExportCostReportRequest true "导出成本报告请求"
// @Success 200 {object} response.Response{data=response.CostReportExportResponse}
// @Router /api/designer/cost/export [post]
func (api *CostAPI) ExportCostReport(c *gin.Context) {
	var req request.ExportCostReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.costExportService.ExportCostReport(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
