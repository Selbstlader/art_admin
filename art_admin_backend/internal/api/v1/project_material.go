package v1

import (
	"net/http"
	"strconv"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/project_material"

	"github.com/gin-gonic/gin"
)

// ProjectMaterialAPI 项目材料清单API
// Project material list API controller
type ProjectMaterialAPI struct {
	service *project_material.ProjectMaterialService
}

// NewProjectMaterialAPI 创建项目材料清单API实例
// Create project material API instance
func NewProjectMaterialAPI() *ProjectMaterialAPI {
	return &ProjectMaterialAPI{
		service: project_material.NewProjectMaterialService(),
	}
}

// GetProjectMaterials 获取项目材料清单
// @Summary 获取项目材料清单
// @Description 获取指定项目的材料清单
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} response.Response{data=response.ProjectMaterialListResponse}
// @Router /api/designer/projects/{id}/materials [get]
func (api *ProjectMaterialAPI) GetProjectMaterials(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	result, err := api.service.GetProjectMaterials(uint(projectID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, result)
}

// SaveProjectMaterials 保存项目材料清单
// @Summary 保存项目材料清单
// @Description 保存项目的材料清单(替换全部)
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.SaveProjectMaterialRequest true "保存材料清单请求"
// @Success 200 {object} response.Response{data=response.ProjectMaterialListResponse}
// @Router /api/designer/projects/{id}/materials [post]
func (api *ProjectMaterialAPI) SaveProjectMaterials(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	var req request.SaveProjectMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	req.ProjectID = uint(projectID)

	result, err := api.service.SaveProjectMaterials(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddProjectMaterial 添加单个材料到项目
// @Summary 添加单个材料到项目
// @Description 添加单个材料到项目材料清单
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.ProjectMaterialItem true "材料项"
// @Success 200 {object} response.Response{data=response.ProjectMaterialItem}
// @Router /api/designer/projects/{id}/materials/add [post]
func (api *ProjectMaterialAPI) AddProjectMaterial(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	var req request.ProjectMaterialItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.service.AddProjectMaterial(uint(projectID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateProjectMaterial 更新项目材料项
// @Summary 更新项目材料项
// @Description 更新项目材料清单中的某一项
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param itemId path int true "材料项ID"
// @Param request body request.ProjectMaterialItem true "材料项"
// @Success 200 {object} response.Response{data=response.ProjectMaterialItem}
// @Router /api/designer/projects/{id}/materials/{itemId} [put]
func (api *ProjectMaterialAPI) UpdateProjectMaterial(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的材料项ID / Invalid item ID")
		return
	}

	var req request.ProjectMaterialItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.service.UpdateProjectMaterial(uint(projectID), uint(itemID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteProjectMaterial 删除项目材料项
// @Summary 删除项目材料项
// @Description 删除项目材料清单中的某一项
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param itemId path int true "材料项ID"
// @Success 200 {object} response.Response
// @Router /api/designer/projects/{id}/materials/{itemId} [delete]
func (api *ProjectMaterialAPI) DeleteProjectMaterial(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的材料项ID / Invalid item ID")
		return
	}

	if err := api.service.DeleteProjectMaterial(uint(projectID), uint(itemID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetProjectCostSummary 获取项目成本汇总
// @Summary 获取项目成本汇总
// @Description 获取指定项目的成本汇总
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} response.Response{data=response.ProjectCostSummary}
// @Router /api/designer/projects/{id}/cost [get]
func (api *ProjectMaterialAPI) GetProjectCostSummary(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	result, err := api.service.GetProjectCostSummary(uint(projectID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, result)
}

// SaveProjectCost 保存项目成本配置
// @Summary 保存项目成本配置
// @Description 保存项目的成本配置(预算、人工费等)
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.SaveProjectCostRequest true "保存成本配置请求"
// @Success 200 {object} response.Response{data=response.ProjectCostSummary}
// @Router /api/designer/projects/{id}/cost [post]
func (api *ProjectMaterialAPI) SaveProjectCost(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	var req request.SaveProjectCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	req.ProjectID = uint(projectID)

	result, err := api.service.SaveProjectCost(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ImportMaterialsToProject 从材料库导入材料到项目
// @Summary 从材料库导入材料到项目
// @Description 从材料库选择材料导入到项目材料清单
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.ImportMaterialsRequest true "导入材料请求"
// @Success 200 {object} response.Response{data=response.ProjectMaterialListResponse}
// @Router /api/designer/projects/{id}/materials/import [post]
func (api *ProjectMaterialAPI) ImportMaterialsToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID / Invalid project ID")
		return
	}

	var req request.ImportMaterialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误 / Invalid parameters: "+err.Error())
		return
	}

	result, err := api.service.ImportMaterialsToProject(uint(projectID), req.MaterialIDs, req.DefaultQuantity)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ExportProjectMaterials 导出项目材料清单
// @Summary 导出项目材料清单
// @Description 导出项目材料清单为Excel或PDF
// @Tags 项目材料管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.ExportProjectMaterialsRequest true "导出请求"
// @Success 200 {object} response.Response{data=response.ExportProjectMaterialsResponse}
// @Router /api/designer/projects/{id}/materials/export [post]
func (api *ProjectMaterialAPI) ExportProjectMaterials(c *gin.Context) {
	// TODO: 实现导出功能 / Implement export functionality
	response.Error(c, http.StatusNotImplemented, "导出功能开发中 / Export feature in development")
}
