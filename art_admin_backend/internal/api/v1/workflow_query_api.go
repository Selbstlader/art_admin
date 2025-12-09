package v1

import (
	"context"
	"strconv"

	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	pkgResponse "art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service/workflow"

	"github.com/gin-gonic/gin"
)

// queryService 查询服务实例
var queryService workflow.QueryService

// initQueryService 初始化查询服务
func initQueryService() {
	if queryService == nil {
		processDefRepo := repository.NewProcessDefinitionRepository()
		processInstRepo := repository.NewProcessInstanceRepository()
		taskRepo := repository.NewWfTaskRepository()
		queryService = workflow.NewQueryService(processDefRepo, processInstRepo, taskRepo)
	}
}

// GetMyInitiated 获取我发起的流程
// @Summary 获取我发起的流程
// @Description 查询当前用户发起的流程实例列表
// @Tags 工作流-查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param title query string false "标题"
// @Param status query string false "状态"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessInstListResponse}
// @Router /api/workflow/query/my-initiated [get]
func GetMyInitiated(c *gin.Context) {
	initQueryService()

	var req request.ListProcessInstRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.QueryRequest{
		Current: req.Page,
		Size:    req.PageSize,
		Title:   req.Title,
		Status:  req.Status,
	}

	summaries, total, err := queryService.ListMyInitiated(context.Background(), userID, svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.ProcessInstSummaryResponse, 0, len(summaries))
	for _, s := range summaries {
		list = append(list, toProcessInstSummaryResponse(s))
	}

	pkgResponse.SuccessWithPagination(c, list, req.Page, req.PageSize, total)
}

// GetMyTodo 获取我的待办
// @Summary 获取我的待办
// @Description 查询当前用户的待办任务列表
// @Tags 工作流-查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param processCode query string false "流程编码"
// @Success 200 {object} pkgResponse.Response{data=response.WfTaskListResponse}
// @Router /api/workflow/query/my-todo [get]
func GetMyTodo(c *gin.Context) {
	initQueryService()

	var req request.ListTaskRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.QueryRequest{
		Current: req.Page,
		Size:    req.PageSize,
	}

	summaries, total, err := queryService.ListMyTodo(context.Background(), userID, svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.WfTaskSummaryResponse, 0, len(summaries))
	for _, s := range summaries {
		list = append(list, toWfTaskSummaryResponse(s))
	}

	pkgResponse.SuccessWithPagination(c, list, req.Page, req.PageSize, total)
}

// GetMyDone 获取我的已办
// @Summary 获取我的已办
// @Description 查询当前用户的已办任务列表
// @Tags 工作流-查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页条数" default(10)
// @Param processCode query string false "流程编码"
// @Param status query string false "状态"
// @Success 200 {object} pkgResponse.Response{data=response.WfTaskListResponse}
// @Router /api/workflow/query/my-done [get]
func GetMyDone(c *gin.Context) {
	initQueryService()

	var req request.ListTaskRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.QueryRequest{
		Current: req.Page,
		Size:    req.PageSize,
		Status:  req.Status,
	}

	summaries, total, err := queryService.ListMyDone(context.Background(), userID, svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.WfTaskSummaryResponse, 0, len(summaries))
	for _, s := range summaries {
		list = append(list, toWfTaskSummaryResponse(s))
	}

	pkgResponse.SuccessWithPagination(c, list, req.Page, req.PageSize, total)
}

// GetApprovalTrail 获取审批轨迹
// @Summary 获取审批轨迹
// @Description 获取流程实例的完整审批轨迹
// @Tags 工作流-查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程实例ID"
// @Success 200 {object} pkgResponse.Response{data=[]response.ApprovalRecordResponse}
// @Router /api/workflow/query/approval-trail/{id} [get]
func GetApprovalTrail(c *gin.Context) {
	initQueryService()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	records, err := queryService.GetApprovalTrail(context.Background(), id)
	if err != nil {
		pkgResponse.ServerError(c, "查询失败: "+err.Error())
		return
	}

	// 转换为响应格式
	list := make([]response.ApprovalRecordResponse, 0, len(records))
	for _, r := range records {
		list = append(list, toApprovalRecordResponseFromQuery(r))
	}

	pkgResponse.Success(c, list)
}

// GetWorkflowStatistics 获取流程统计
// @Summary 获取流程统计
// @Description 获取流程统计数据
// @Tags 工作流-查询
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param processDefId query int false "流程定义ID"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessStatisticsResponse}
// @Router /api/workflow/query/statistics [get]
func GetWorkflowStatistics(c *gin.Context) {
	initQueryService()

	var req request.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 转换为服务层请求
	svcReq := &workflow.StatisticsRequest{
		ProcessDefID: nil,
		StartTime:    req.StartDate,
		EndTime:      req.EndDate,
	}

	stats, err := queryService.GetStatistics(context.Background(), svcReq)
	if err != nil {
		pkgResponse.ServerError(c, "查询失败: "+err.Error())
		return
	}

	pkgResponse.Success(c, toProcessStatisticsResponse(stats))
}

// ============================================================================
// 辅助函数：DTO转换
// ============================================================================

// toProcessInstSummaryResponse 转换为流程实例摘要响应
func toProcessInstSummaryResponse(s *workflow.ProcessInstanceSummary) response.ProcessInstSummaryResponse {
	return response.ProcessInstSummaryResponse{
		ID:            s.ID,
		ProcessName:   s.ProcessDefName,
		Title:         s.Title,
		InitiatorID:   s.InitiatorID,
		InitiatorName: s.InitiatorName,
		Status:        s.Status,
		StartedAt:     s.StartedAt,
		CompletedAt:   s.CompletedAt,
	}
}

// toWfTaskSummaryResponse 转换为任务摘要响应
func toWfTaskSummaryResponse(s *workflow.TaskSummary) response.WfTaskSummaryResponse {
	return response.WfTaskSummaryResponse{
		ID:            s.ID,
		ProcessInstID: s.ProcessInstID,
		ProcessTitle:  s.ProcessTitle,
		ProcessName:   s.ProcessDefName,
		NodeName:      s.NodeName,
		Status:        s.Status,
		DueAt:         s.DueAt,
		CreatedAt:     s.CreatedAt,
		CompletedAt:   s.CompletedAt,
	}
}

// toApprovalRecordResponseFromQuery 转换为审批记录响应
func toApprovalRecordResponseFromQuery(r *workflow.ApprovalRecord) response.ApprovalRecordResponse {
	return response.ApprovalRecordResponse{
		TaskID:       r.TaskID,
		NodeID:       r.NodeID,
		NodeName:     r.NodeName,
		OperatorID:   r.OperatorID,
		OperatorName: r.OperatorName,
		Action:       r.Action,
		Comment:      r.Comment,
		CreatedAt:    r.CreatedAt,
	}
}

// toProcessStatisticsResponse 转换为流程统计响应
func toProcessStatisticsResponse(s *workflow.ProcessStatistics) response.ProcessStatisticsResponse {
	resp := response.ProcessStatisticsResponse{
		TotalInstances:     s.TotalCount,
		RunningInstances:   s.RunningCount,
		CompletedInstances: s.CompletedCount,
		RejectedInstances:  s.RejectedCount,
		WithdrawnInstances: s.WithdrawnCount,
		AvgProcessingTime:  s.AvgDuration,
		ByProcess:          make([]response.ProcessTypeStatistics, 0, len(s.ByProcessDef)),
	}

	for _, p := range s.ByProcessDef {
		resp.ByProcess = append(resp.ByProcess, response.ProcessTypeStatistics{
			ProcessName:       p.ProcessDefName,
			TotalCount:        p.TotalCount,
			CompletedCount:    p.CompletedCount,
			AvgProcessingTime: p.AvgDuration,
		})
	}

	return resp
}
