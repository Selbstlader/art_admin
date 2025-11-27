package v1

import (
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/pkg/response"
	operationLogSvc "art_admin_backend/internal/service/operation_log"
	"strconv"

	"github.com/gin-gonic/gin"
)

var operationLogService = operationLogSvc.NewOperationLogService()

// GetOperationLogList 获取操作日志列表
// @Summary 获取操作日志列表
// @Description 分页查询操作日志列表
// @Tags 操作日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param current query int true "当前页码" minimum(1)
// @Param size query int true "每页条数" minimum(1) maximum(100)
// @Param module query string false "操作模块"
// @Param businessType query string false "业务类型"
// @Param operatorName query string false "操作人员"
// @Param status query int false "操作状态"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} response.Response{data=response.PaginatedData} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/operation-log/list [get]
func GetOperationLogList(c *gin.Context) {
	var req request.OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	logs, total, err := operationLogService.GetOperationLogList(&req)
	if err != nil {
		response.ServerError(c, "查询操作日志列表失败")
		return
	}

	response.SuccessWithPagination(c, logs, req.Current, req.Size, total)
}

// GetOperationLogDetail 获取操作日志详情
// @Summary 获取操作日志详情
// @Description 根据ID获取操作日志详情
// @Tags 操作日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日志ID"
// @Success 200 {object} response.Response{data=response.OperationLogDetail} "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/operation-log/{id} [get]
func GetOperationLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	log, err := operationLogService.GetOperationLogDetail(id)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, log)
}

// DeleteOperationLog 删除操作日志
// @Summary 删除操作日志
// @Description 删除操作日志
// @Tags 操作日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日志ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/operation-log/{id} [delete]
func DeleteOperationLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := operationLogService.DeleteOperationLog(id); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除操作日志成功", nil)
}

// BatchDeleteOperationLog 批量删除操作日志
// @Summary 批量删除操作日志
// @Description 批量删除操作日志
// @Tags 操作日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ids body []int64 true "日志ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/operation-log/batch-delete [delete]
func BatchDeleteOperationLog(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := operationLogService.BatchDeleteOperationLog(req.IDs); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除操作日志成功", nil)
}

// CleanOperationLog 清理操作日志
// @Summary 清理操作日志
// @Description 清理指定天数之前的操作日志
// @Tags 操作日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days body int true "保留天数"
// @Success 200 {object} response.Response "清理成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /api/operation-log/clean [post]
func CleanOperationLog(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := operationLogService.CleanOperationLog(req.Days); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.SuccessWithMsg(c, "清理操作日志成功", nil)
}
