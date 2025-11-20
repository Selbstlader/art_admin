package project

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var taskService = service.NewTaskService()

// GetTaskList 获取任务列表
func GetTaskList(c *gin.Context) {
	var req request.GetTaskListRequest
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

	tasks, total, err := taskService.GetTaskList(&req)
	if err != nil {
		response.ServerError(c, "查询任务列表失败")
		return
	}

	response.SuccessWithPagination(c, tasks, req.Page, req.PageSize, total)
}

// GetTaskDetail 获取任务详情
func GetTaskDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	task, err := taskService.GetTaskDetail(id)
	if err != nil {
		response.ServerError(c, "查询任务详情失败")
		return
	}

	response.Success(c, task)
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	var req request.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	userName := "管理员" // 暂时使用固定值,后续可以从用户表查询

	err := taskService.CreateTask(&req, userID.(int64), userName)
	if err != nil {
		response.ServerError(c, "创建任务失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	var req request.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := taskService.UpdateTask(&req)
	if err != nil {
		response.ServerError(c, "更新任务失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchUpdateTask 批量更新任务(甘特图拖拽)
func BatchUpdateTask(c *gin.Context) {
	var req request.BatchUpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := taskService.BatchUpdateTask(&req)
	if err != nil {
		response.ServerError(c, "批量更新任务失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// QuickUpdateTask 快速更新任务(甘特图实时更新)
func QuickUpdateTask(c *gin.Context) {
	var req request.QuickUpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := taskService.QuickUpdateTask(&req)
	if err != nil {
		response.ServerError(c, "快速更新任务失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	err = taskService.DeleteTask(id)
	if err != nil {
		response.ServerError(c, "删除任务失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetTaskGanttData 获取甘特图数据
func GetTaskGanttData(c *gin.Context) {
	var req request.GetTaskGanttDataRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	data, err := taskService.GetTaskGanttData(req.ProjectID)
	if err != nil {
		response.ServerError(c, "获取甘特图数据失败")
		return
	}

	response.Success(c, data)
}

// CreateTaskDependency 创建任务依赖
func CreateTaskDependency(c *gin.Context) {
	var req request.CreateTaskDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := taskService.CreateDependency(&req)
	if err != nil {
		response.ServerError(c, "创建依赖关系失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteTaskDependency 删除任务依赖
func DeleteTaskDependency(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的依赖ID")
		return
	}

	err = taskService.DeleteDependency(id)
	if err != nil {
		response.ServerError(c, "删除依赖关系失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateTaskComment 创建任务评论
func CreateTaskComment(c *gin.Context) {
	var req request.CreateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户信息
	userID, _ := c.Get("userId")
	userName := "管理员" // 暂时使用固定值,后续可以从用户表查询

	err := taskService.CreateComment(&req, userID.(int64), userName)
	if err != nil {
		response.ServerError(c, "创建评论失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
