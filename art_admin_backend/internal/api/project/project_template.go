package project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	projectSvc "art_admin_backend/internal/service/project"
	"strconv"

	"github.com/gin-gonic/gin"
)

var templateService = projectSvc.NewProjectTemplateService()

// GetTemplateList 获取项目模板列表
func GetTemplateList(c *gin.Context) {
	var req request.GetProjectTemplateListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	templates, total, err := templateService.GetTemplateList(&req)
	if err != nil {
		response.ServerError(c, "查询模板列表失败")
		return
	}

	response.SuccessWithPagination(c, templates, req.Page, req.PageSize, total)
}

// GetTemplateDetail 获取项目模板详情
func GetTemplateDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的模板ID")
		return
	}

	template, err := templateService.GetTemplateDetail(id)
	if err != nil {
		response.ServerError(c, "查询模板详情失败")
		return
	}

	response.Success(c, template)
}

// CreateTemplate 创建项目模板
func CreateTemplate(c *gin.Context) {
	var req request.CreateProjectTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	userName := "管理员" // 暂时使用固定值,后续可以从用户表查询

	err := templateService.CreateTemplate(&req, userID.(int64), userName)
	if err != nil {
		response.ServerError(c, "创建模板失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateTemplate 更新项目模板
func UpdateTemplate(c *gin.Context) {
	var req request.UpdateProjectTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := templateService.UpdateTemplate(&req)
	if err != nil {
		response.ServerError(c, "更新模板失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteTemplate 删除项目模板
func DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的模板ID")
		return
	}

	err = templateService.DeleteTemplate(id)
	if err != nil {
		response.ServerError(c, "删除模板失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateTemplateDependency 创建模板任务依赖
func CreateTemplateDependency(c *gin.Context) {
	var req request.CreateTemplateDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := templateService.CreateDependency(&req)
	if err != nil {
		response.ServerError(c, "创建依赖关系失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteTemplateDependency 删除模板任务依赖
func DeleteTemplateDependency(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的依赖ID")
		return
	}

	err = templateService.DeleteDependency(id)
	if err != nil {
		response.ServerError(c, "删除依赖关系失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
