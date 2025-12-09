package v1

import (
	"art_admin_backend/internal/api/middleware"
	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	pkgResponse "art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/repository"
	"art_admin_backend/internal/service/workflow"
	"context"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

// workflowEngine 工作流引擎实例
var workflowEngine workflow.WorkflowEngine

// initWorkflowEngine 初始化工作流引擎
func initWorkflowEngine() {
	if workflowEngine == nil {
		processDefRepo := repository.NewProcessDefinitionRepository()
		processInstRepo := repository.NewProcessInstanceRepository()
		taskRepo := repository.NewWfTaskRepository()
		formTemplateRepo := repository.NewFormTemplateRepository()
		serializer := workflow.NewSerializer()
		validator := workflow.NewValidator()
		formEngine := workflow.NewFormEngineService(formTemplateRepo, serializer)
		// 创建数据提供者
		userProvider := NewUserProviderImpl()
		deptProvider := NewDepartmentProviderImpl()
		orgProvider := NewOrganizationProviderImpl()
		assigneeResolver := workflow.NewAssigneeResolver(userProvider, deptProvider, orgProvider)
		// 通知服务可以为nil，后续可以注入
		workflowEngine = workflow.NewWorkflowEngine(
			processDefRepo,
			processInstRepo,
			taskRepo,
			formTemplateRepo,
			serializer,
			validator,
			formEngine,
			assigneeResolver,
			nil, // notificationService
		)
	}
}

// StartProcess 启动流程
// @Summary 启动流程
// @Description 发起新的审批流程
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.StartProcessRequest true "启动流程请求"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessInstResponse}
// @Router /api/workflow/process-inst/start [post]
func StartProcess(c *gin.Context) {
	initWorkflowEngine()

	var req request.StartProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户信息
	userID := middleware.GetUserID(c)
	userName := middleware.GetUsername(c)

	// 转换为服务层请求
	svcReq := &workflow.StartProcessRequest{
		ProcessDefID:  req.ProcessDefID,
		Title:         req.Title,
		InitiatorID:   userID,
		InitiatorName: userName,
		FormData:      req.FormData,
	}

	instance, err := workflowEngine.StartProcess(context.Background(), svcReq)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "流程启动成功", toProcessInstResponse(instance))
}

// CompleteTask 完成任务
// @Summary 完成任务
// @Description 审批通过或拒绝任务
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CompleteTaskRequest true "完成任务请求"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/process-inst/complete-task [post]
func CompleteTask(c *gin.Context) {
	initWorkflowEngine()

	var req request.CompleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	// 转换为服务层请求
	svcReq := &workflow.CompleteTaskRequest{
		TaskID:     req.TaskID,
		Action:     req.Action,
		Comment:    req.Comment,
		OperatorID: userID,
	}

	if err := workflowEngine.CompleteTask(context.Background(), svcReq); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "任务处理成功", nil)
}

// DelegateTask 委托任务
// @Summary 委托任务
// @Description 将任务委托给其他用户处理
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.DelegateTaskRequest true "委托任务请求"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/process-inst/delegate-task [post]
func DelegateTask(c *gin.Context) {
	initWorkflowEngine()

	var req request.DelegateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户信息
	userID := middleware.GetUserID(c)
	userName := middleware.GetUsername(c)

	// 转换为服务层请求
	svcReq := &workflow.DelegateTaskRequest{
		TaskID:       req.TaskID,
		ToUserID:     req.ToUserID,
		Reason:       req.Reason,
		OperatorID:   userID,
		OperatorName: userName,
	}

	if err := workflowEngine.DelegateTask(context.Background(), svcReq); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "任务委托成功", nil)
}

// TransferTask 转办任务
// @Summary 转办任务
// @Description 将任务转办给其他用户
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TransferTaskRequest true "转办任务请求"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/process-inst/transfer-task [post]
func TransferTask(c *gin.Context) {
	initWorkflowEngine()

	var req request.TransferTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户信息
	userID := middleware.GetUserID(c)
	userName := middleware.GetUsername(c)

	// 转换为服务层请求
	svcReq := &workflow.TransferTaskRequest{
		TaskID:       req.TaskID,
		ToUserID:     req.ToUserID,
		Reason:       req.Reason,
		OperatorID:   userID,
		OperatorName: userName,
	}

	if err := workflowEngine.TransferTask(context.Background(), svcReq); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "任务转办成功", nil)
}

// WithdrawProcess 撤回流程
// @Summary 撤回流程
// @Description 撤回自己发起的流程
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程实例ID"
// @Success 200 {object} pkgResponse.Response
// @Router /api/workflow/process-inst/{id}/withdraw [post]
func WithdrawProcess(c *gin.Context) {
	initWorkflowEngine()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserID(c)

	if err := workflowEngine.WithdrawProcess(context.Background(), id, userID); err != nil {
		pkgResponse.Error(c, pkgResponse.CodeBadRequest, err.Error())
		return
	}

	pkgResponse.SuccessWithMsg(c, "流程撤回成功", nil)
}

// GetProcessInstDetail 获取流程实例详情
// @Summary 获取流程实例详情
// @Description 获取流程实例详情，包含审批轨迹
// @Tags 工作流-流程实例
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流程实例ID"
// @Success 200 {object} pkgResponse.Response{data=response.ProcessInstDetailResponse}
// @Router /api/workflow/process-inst/{id} [get]
func GetProcessInstDetail(c *gin.Context) {
	initWorkflowEngine()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkgResponse.BadRequest(c, "无效的ID参数")
		return
	}

	detail, err := workflowEngine.GetProcessInstance(context.Background(), id)
	if err != nil {
		pkgResponse.Error(c, pkgResponse.CodeNotFound, err.Error())
		return
	}

	pkgResponse.Success(c, toProcessInstDetailResponse(detail))
}

// ============================================================================
// 辅助函数：DTO转换
// ============================================================================

// toProcessInstResponse 转换为流程实例响应
func toProcessInstResponse(inst *model.ProcessInstance) response.ProcessInstResponse {
	return response.ProcessInstResponse{
		ID:                inst.ID,
		ProcessDefID:      inst.ProcessDefID,
		ProcessDefVersion: inst.ProcessDefVersion,
		Title:             inst.Title,
		InitiatorID:       inst.InitiatorID,
		InitiatorName:     inst.InitiatorName,
		Status:            inst.Status,
		CurrentNodeID:     inst.CurrentNodeID,
		StartedAt:         inst.StartedAt,
		CompletedAt:       inst.CompletedAt,
	}
}

// toProcessInstDetailResponse 转换为流程实例详情响应
func toProcessInstDetailResponse(detail *workflow.ProcessInstanceDetail) response.ProcessInstDetailResponse {
	resp := response.ProcessInstDetailResponse{
		ProcessInstResponse: toProcessInstResponse(detail.Instance),
		CurrentTasks:        make([]response.WfTaskResponse, 0),
		ApprovalTrail:       make([]response.ApprovalRecordResponse, 0),
	}

	// 解析表单数据
	if detail.Instance.FormData != "" {
		// FormData 已经是 JSON 字符串，需要解析
		var formData map[string]interface{}
		if err := json.Unmarshal([]byte(detail.Instance.FormData), &formData); err == nil {
			resp.FormData = formData
		}
	}

	// 转换当前任务
	for _, task := range detail.Tasks {
		if task.Status == model.WfTaskStatusPending {
			resp.CurrentTasks = append(resp.CurrentTasks, toWfTaskResponse(task))
		}
	}

	// 转换审批轨迹
	for _, log := range detail.ApprovalLogs {
		resp.ApprovalTrail = append(resp.ApprovalTrail, toApprovalRecordResponse(log))
	}

	return resp
}

// toWfTaskResponse 转换为任务响应
func toWfTaskResponse(task *model.WfTask) response.WfTaskResponse {
	return response.WfTaskResponse{
		ID:              task.ID,
		ProcessInstID:   task.ProcessInstID,
		NodeID:          task.NodeID,
		NodeName:        task.NodeName,
		AssigneeID:      task.AssigneeID,
		AssigneeName:    task.AssigneeName,
		Status:          task.Status,
		Comment:         task.Comment,
		DelegatedFrom:   task.DelegatedFrom,
		TransferredFrom: task.TransferredFrom,
		DueAt:           task.DueAt,
		CreatedAt:       task.CreatedAt,
		CompletedAt:     task.CompletedAt,
	}
}

// toApprovalRecordResponse 转换为审批记录响应
func toApprovalRecordResponse(log *model.TaskLog) response.ApprovalRecordResponse {
	return response.ApprovalRecordResponse{
		TaskID:       log.TaskID,
		OperatorID:   log.OperatorID,
		OperatorName: log.OperatorName,
		Action:       log.Action,
		Comment:      log.Comment,
		CreatedAt:    log.CreatedAt,
	}
}

// ============================================================================
// Provider Implementations
// ============================================================================

// UserProviderImpl 用户数据提供者实现
type UserProviderImpl struct{}

// NewUserProviderImpl 创建用户数据提供者
func NewUserProviderImpl() *UserProviderImpl {
	return &UserProviderImpl{}
}

// GetUserByID 根据ID获取用户
func (p *UserProviderImpl) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByIDs 根据ID列表获取用户
func (p *UserProviderImpl) GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*model.User, error) {
	var users []*model.User
	if err := database.GetDB().Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByRoleID 根据角色ID获取用户列表
func (p *UserProviderImpl) GetUsersByRoleID(ctx context.Context, roleID int64) ([]*model.User, error) {
	var users []*model.User
	if err := database.GetDB().Where("role_id = ?", roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// DepartmentProviderImpl 部门数据提供者实现
type DepartmentProviderImpl struct{}

// NewDepartmentProviderImpl 创建部门数据提供者
func NewDepartmentProviderImpl() *DepartmentProviderImpl {
	return &DepartmentProviderImpl{}
}

// GetDepartmentByID 根据ID获取部门
func (p *DepartmentProviderImpl) GetDepartmentByID(ctx context.Context, deptID int64) (*model.Department, error) {
	var dept model.Department
	if err := database.GetDB().First(&dept, deptID).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// GetDepartmentLeader 获取部门负责人用户ID
func (p *DepartmentProviderImpl) GetDepartmentLeader(ctx context.Context, deptID int64) (int64, error) {
	var dept model.Department
	if err := database.GetDB().First(&dept, deptID).Error; err != nil {
		return 0, err
	}
	// 通过Leader名称查找用户ID
	if dept.Leader != "" {
		var user model.User
		if err := database.GetDB().Where("user_name = ? OR nick_name = ?", dept.Leader, dept.Leader).First(&user).Error; err == nil {
			return user.ID, nil
		}
	}
	return 0, nil
}

// OrganizationProviderImpl 组织架构数据提供者实现
type OrganizationProviderImpl struct{}

// NewOrganizationProviderImpl 创建组织架构数据提供者
func NewOrganizationProviderImpl() *OrganizationProviderImpl {
	return &OrganizationProviderImpl{}
}

// GetUserSuperior 获取用户的直接上级
func (p *OrganizationProviderImpl) GetUserSuperior(ctx context.Context, userID int64) (int64, error) {
	// 通过用户所属部门的负责人作为上级
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		return 0, err
	}
	if user.DeptID == nil {
		return 0, nil
	}
	var dept model.Department
	if err := database.GetDB().First(&dept, *user.DeptID).Error; err != nil {
		return 0, err
	}
	// 通过Leader名称查找用户ID
	if dept.Leader != "" {
		var leader model.User
		if err := database.GetDB().Where("user_name = ? OR nick_name = ?", dept.Leader, dept.Leader).First(&leader).Error; err == nil {
			if leader.ID != userID {
				return leader.ID, nil
			}
		}
	}
	return 0, nil
}
