package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	designerProjectService "art_admin_backend/internal/service/designer_project"
	"strconv"

	"github.com/gin-gonic/gin"
)

var designerProjectSvc *designerProjectService.DesignerProjectService

// SetDesignerProjectService 设置设计师项目服务
// Set designer project service
func SetDesignerProjectService(svc *designerProjectService.DesignerProjectService) {
	designerProjectSvc = svc
}

// CreateDesignerProject 创建设计师项目
// @Summary 创建设计师项目
// @Description 创建新的工装设计项目
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param request body request.CreateDesignerProjectRequest true "创建项目请求"
// @Success 200 {object} response.Response{data=response.DesignerProjectResponse}
// @Router /api/designer/projects [post]
func CreateDesignerProject(c *gin.Context) {
	var req request.CreateDesignerProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designerProjectSvc.Create(&req, userID)
	if err != nil {
		response.ServerError(c, "创建项目失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", result)
}

// GetDesignerProjectList 获取设计师项目列表
// @Summary 获取设计师项目列表
// @Description 分页获取当前用户的工装设计项目列表，支持多条件筛选
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param current query int true "当前页码" default(1)
// @Param size query int true "每页数量" default(10)
// @Param name query string false "项目名称"
// @Param status query string false "状态(draft/in_progress/completed/archived)"
// @Param style query string false "设计风格"
// @Param minBudget query number false "最小预算"
// @Param maxBudget query number false "最大预算"
// @Param startDate query string false "开始日期(YYYY-MM-DD)"
// @Param endDate query string false "结束日期(YYYY-MM-DD)"
// @Param keyword query string false "关键字搜索"
// @Success 200 {object} response.Response{data=response.DesignerProjectListResponse}
// @Router /api/designer/projects [get]
func GetDesignerProjectList(c *gin.Context) {
	var req request.DesignerProjectListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designerProjectSvc.List(&req, userID)
	if err != nil {
		response.ServerError(c, "获取项目列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignerProjectDetail 获取设计师项目详情
// @Summary 获取设计师项目详情
// @Description 获取指定项目的详细信息，包含关联的文档、CAD文件、成本估算等
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} response.Response{data=response.DesignerProjectDetailResponse}
// @Router /api/designer/projects/{id} [get]
func GetDesignerProjectDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designerProjectSvc.GetDetailByID(uint(id), userID)
	if err != nil {
		if err.Error() == "无权访问该项目" {
			response.Forbidden(c, err.Error())
			return
		}
		response.NotFound(c, "项目不存在")
		return
	}

	response.Success(c, result)
}

// UpdateDesignerProject 更新设计师项目
// @Summary 更新设计师项目
// @Description 更新指定项目的信息
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param request body request.UpdateDesignerProjectRequest true "更新项目请求"
// @Success 200 {object} response.Response{data=response.DesignerProjectResponse}
// @Router /api/designer/projects/{id} [put]
func UpdateDesignerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	var req request.UpdateDesignerProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 确保路径ID和请求体ID一致 / Ensure path ID matches request body ID
	req.ID = uint(id)

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designerProjectSvc.Update(&req, userID)
	if err != nil {
		if err.Error() == "无权修改该项目" {
			response.Forbidden(c, err.Error())
			return
		}
		response.ServerError(c, "更新项目失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", result)
}

// DeleteDesignerProject 删除设计师项目
// @Summary 删除设计师项目
// @Description 删除指定项目及其所有关联数据（文档、CAD文件、比对结果、成本估算）
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} response.Response
// @Router /api/designer/projects/{id} [delete]
func DeleteDesignerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := designerProjectSvc.Delete(uint(id), userID); err != nil {
		if err.Error() == "无权删除该项目" {
			response.Forbidden(c, err.Error())
			return
		}
		response.ServerError(c, "删除项目失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchDeleteDesignerProjects 批量删除设计师项目
// @Summary 批量删除设计师项目
// @Description 批量删除多个项目及其所有关联数据
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param request body request.BatchDeleteDesignerProjectRequest true "批量删除请求"
// @Success 200 {object} response.Response
// @Router /api/designer/projects/batch-delete [post]
func BatchDeleteDesignerProjects(c *gin.Context) {
	var req request.BatchDeleteDesignerProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	if err := designerProjectSvc.BatchDelete(req.IDs, userID); err != nil {
		response.ServerError(c, "批量删除失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// SearchDesignerProjects 搜索设计师项目
// @Summary 搜索设计师项目
// @Description 根据关键字搜索项目（搜索项目名称和描述）
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键字"
// @Param current query int false "当前页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.DesignerProjectListResponse}
// @Router /api/designer/projects/search [get]
func SearchDesignerProjects(c *gin.Context) {
	keyword := c.Query("keyword")
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if current < 1 {
		current = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := designerProjectSvc.Search(keyword, userID, current, size)
	if err != nil {
		response.ServerError(c, "搜索项目失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetDesignerProjectStats 获取设计师项目统计
// @Summary 获取设计师项目统计
// @Description 获取当前用户的项目统计信息
// @Tags 设计师项目管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/designer/projects/stats [get]
func GetDesignerProjectStats(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	// 获取各状态的项目数量 / Get project count by status
	draftReq := &request.DesignerProjectListRequest{Current: 1, Size: 1, Status: "draft"}
	inProgressReq := &request.DesignerProjectListRequest{Current: 1, Size: 1, Status: "in_progress"}
	completedReq := &request.DesignerProjectListRequest{Current: 1, Size: 1, Status: "completed"}
	archivedReq := &request.DesignerProjectListRequest{Current: 1, Size: 1, Status: "archived"}
	allReq := &request.DesignerProjectListRequest{Current: 1, Size: 1}

	draftResult, _ := designerProjectSvc.List(draftReq, userID)
	inProgressResult, _ := designerProjectSvc.List(inProgressReq, userID)
	completedResult, _ := designerProjectSvc.List(completedReq, userID)
	archivedResult, _ := designerProjectSvc.List(archivedReq, userID)
	allResult, _ := designerProjectSvc.List(allReq, userID)

	stats := gin.H{
		"total":      allResult.Total,
		"draft":      draftResult.Total,
		"inProgress": inProgressResult.Total,
		"completed":  completedResult.Total,
		"archived":   archivedResult.Total,
	}

	response.Success(c, stats)
}
