package project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var projectService = service.NewProjectService()

// GetProjectList 获取项目列表
func GetProjectList(c *gin.Context) {
	var req request.GetProjectListRequest
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

	projects, total, err := projectService.GetProjectList(&req)
	if err != nil {
		response.ServerError(c, "查询项目列表失败")
		return
	}

	response.SuccessWithPagination(c, projects, req.Page, req.PageSize, total)
}

// GetProjectDetail 获取项目详情
func GetProjectDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	project, err := projectService.GetProjectDetail(id)
	if err != nil {
		response.ServerError(c, "查询项目详情失败")
		return
	}

	response.Success(c, project)
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	userName := "管理员" // 暂时使用固定值,后续可以从用户表查询

	err := projectService.CreateProject(&req, userID.(int64), userName)
	if err != nil {
		response.ServerError(c, "创建项目失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateProjectFromTemplate 从模板创建项目
func CreateProjectFromTemplate(c *gin.Context) {
	var req request.CreateProjectFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	userName := "管理员" // 暂时使用固定值,后续可以从用户表查询

	err := projectService.CreateProjectFromTemplate(&req, userID.(int64), userName)
	if err != nil {
		response.ServerError(c, "从模板创建项目失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := projectService.UpdateProject(&req)
	if err != nil {
		response.ServerError(c, "更新项目失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	err = projectService.DeleteProject(id)
	if err != nil {
		response.ServerError(c, "删除项目失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetProjectStatistics 获取项目统计
func GetProjectStatistics(c *gin.Context) {
	stats, err := projectService.GetProjectStatistics()
	if err != nil {
		response.ServerError(c, "获取项目统计失败")
		return
	}

	response.Success(c, stats)
}
