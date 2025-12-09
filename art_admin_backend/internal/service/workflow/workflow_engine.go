package workflow

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ============================================================================
// Errors
// ============================================================================

var (
	// ErrProcessDefNotFound 流程定义不存在
	ErrProcessDefNotFound = errors.New("流程定义不存在")
	// ErrProcessDefNotPublished 流程定义未发布
	ErrProcessDefNotPublished = errors.New("流程定义未发布")
	// ErrProcessInstNotFound 流程实例不存在
	ErrProcessInstNotFound = errors.New("流程实例不存在")
	// ErrTaskNotFound 任务不存在
	ErrTaskNotFound = errors.New("任务不存在")
	// ErrTaskNotPending 任务非待处理状态
	ErrTaskNotPending = errors.New("任务非待处理状态")
	// ErrNotTaskAssignee 非任务审批人
	ErrNotTaskAssignee = errors.New("非任务审批人")
	// ErrProcessNotRunning 流程非运行中状态
	ErrProcessNotRunning = errors.New("流程非运行中状态")
	// ErrNotProcessInitiator 非流程发起人
	ErrNotProcessInitiator = errors.New("非流程发起人")
	// ErrFormValidationFailed 表单验证失败
	ErrFormValidationFailed = errors.New("表单数据验证失败")
	// ErrStartNodeNotFound 开始节点不存在
	ErrStartNodeNotFound = errors.New("开始节点不存在")
	// ErrNextNodeNotFound 下一节点不存在
	ErrNextNodeNotFound = errors.New("下一节点不存在")
)

// ============================================================================
// Request/Response Types
// ============================================================================

// StartProcessRequest 启动流程请求
type StartProcessRequest struct {
	ProcessDefID  int64                  `json:"processDefId" binding:"required"`
	Title         string                 `json:"title" binding:"required"`
	InitiatorID   int64                  `json:"initiatorId" binding:"required"`
	InitiatorName string                 `json:"initiatorName"`
	DeptID        *int64                 `json:"deptId"`
	FormData      map[string]interface{} `json:"formData"`
}

// CompleteTaskRequest 完成任务请求
type CompleteTaskRequest struct {
	TaskID     int64  `json:"taskId" binding:"required"`
	Action     string `json:"action" binding:"required"` // approve/reject
	Comment    string `json:"comment"`
	OperatorID int64  `json:"operatorId" binding:"required"`
}

// DelegateTaskRequest 委托任务请求
type DelegateTaskRequest struct {
	TaskID       int64  `json:"taskId" binding:"required"`
	ToUserID     int64  `json:"toUserId" binding:"required"`
	ToUserName   string `json:"toUserName"`
	Reason       string `json:"reason"`
	OperatorID   int64  `json:"operatorId" binding:"required"`
	OperatorName string `json:"operatorName"`
}

// TransferTaskRequest 转办任务请求
type TransferTaskRequest struct {
	TaskID       int64  `json:"taskId" binding:"required"`
	ToUserID     int64  `json:"toUserId" binding:"required"`
	ToUserName   string `json:"toUserName"`
	Reason       string `json:"reason"`
	OperatorID   int64  `json:"operatorId" binding:"required"`
	OperatorName string `json:"operatorName"`
}

// ProcessInstanceDetail 流程实例详情
type ProcessInstanceDetail struct {
	Instance     *model.ProcessInstance `json:"instance"`
	Tasks        []*model.WfTask        `json:"tasks"`
	ApprovalLogs []*model.TaskLog       `json:"approvalLogs"`
}

// ============================================================================
// Interface
// ============================================================================

// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
	// StartProcess 启动流程实例
	// Requirements: 2.1 - WHEN 用户提交审批申请 THEN Workflow_Engine SHALL 基于Process_Definition创建Process_Instance并生成首个Task
	// Requirements: 2.2 - WHEN Process_Instance创建成功 THEN Workflow_Engine SHALL 将实例状态设置为运行中并记录发起时间
	StartProcess(ctx context.Context, req *StartProcessRequest) (*model.ProcessInstance, error)

	// CompleteTask 完成任务（通过/拒绝）
	// Requirements: 3.1 - WHEN 审批人通过任务 THEN Workflow_Engine SHALL 完成当前Task并根据Transition创建下一节点的Task
	// Requirements: 3.2 - WHEN 审批人拒绝任务 THEN Workflow_Engine SHALL 根据流程配置决定终止流程或退回上一节点
	CompleteTask(ctx context.Context, req *CompleteTaskRequest) error

	// DelegateTask 委托任务
	// Requirements: 3.5 - WHEN 审批人将任务委托给他人 THEN Workflow_Engine SHALL 更新Task的Assignee并记录Delegation历史
	DelegateTask(ctx context.Context, req *DelegateTaskRequest) error

	// TransferTask 转办任务
	// Requirements: 3.6 - WHEN 审批人转办任务 THEN Workflow_Engine SHALL 将Task完全转移给新Assignee并记录转办原因
	TransferTask(ctx context.Context, req *TransferTaskRequest) error

	// WithdrawProcess 撤回流程
	// Requirements: 2.4 - WHEN 用户撤回申请 THEN Workflow_Engine SHALL 终止Process_Instance并将状态更新为已撤回
	WithdrawProcess(ctx context.Context, instanceID int64, operatorID int64) error

	// GetProcessInstance 获取流程实例详情
	GetProcessInstance(ctx context.Context, instanceID int64) (*ProcessInstanceDetail, error)
}

// ============================================================================
// Implementation
// ============================================================================

// workflowEngine 工作流引擎实现
type workflowEngine struct {
	processDefRepo      *repository.ProcessDefinitionRepository
	processInstRepo     *repository.ProcessInstanceRepository
	taskRepo            *repository.WfTaskRepository
	formTemplateRepo    *repository.FormTemplateRepository
	serializer          Serializer
	validator           Validator
	formEngine          FormEngineService
	assigneeResolver    AssigneeResolver
	notificationService NotificationService
}

// NewWorkflowEngine 创建工作流引擎实例
func NewWorkflowEngine(
	processDefRepo *repository.ProcessDefinitionRepository,
	processInstRepo *repository.ProcessInstanceRepository,
	taskRepo *repository.WfTaskRepository,
	formTemplateRepo *repository.FormTemplateRepository,
	serializer Serializer,
	validator Validator,
	formEngine FormEngineService,
	assigneeResolver AssigneeResolver,
	notificationService NotificationService,
) WorkflowEngine {
	return &workflowEngine{
		processDefRepo:      processDefRepo,
		processInstRepo:     processInstRepo,
		taskRepo:            taskRepo,
		formTemplateRepo:    formTemplateRepo,
		serializer:          serializer,
		validator:           validator,
		formEngine:          formEngine,
		assigneeResolver:    assigneeResolver,
		notificationService: notificationService,
	}
}

// StartProcess 启动流程实例
// **Feature: oa-workflow-engine, Property 5: 流程实例初始状态正确性**
// **Validates: Requirements 2.1, 2.2**
func (e *workflowEngine) StartProcess(ctx context.Context, req *StartProcessRequest) (*model.ProcessInstance, error) {
	// 1. 获取流程定义
	processDef, err := e.processDefRepo.FindByID(req.ProcessDefID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProcessDefNotFound
		}
		return nil, fmt.Errorf("查询流程定义失败: %w", err)
	}

	// 2. 检查流程定义状态
	if processDef.Status != model.ProcessDefStatusPublished {
		return nil, ErrProcessDefNotPublished
	}

	// 3. 解析流程图
	graph, err := e.serializer.DeserializeProcessGraph(processDef.GraphJSON)
	if err != nil {
		return nil, fmt.Errorf("解析流程图失败: %w", err)
	}

	// 4. 验证表单数据（如果有表单模板）
	if processDef.FormTemplateID != nil && req.FormData != nil {
		schema, err := e.formEngine.GetSchema(ctx, *processDef.FormTemplateID)
		if err != nil {
			return nil, fmt.Errorf("获取表单结构失败: %w", err)
		}
		validationResult := e.formEngine.ValidateFormData(ctx, schema, req.FormData)
		if !validationResult.Valid {
			return nil, fmt.Errorf("%w: %v", ErrFormValidationFailed, validationResult.Errors)
		}
	}

	// 5. 查找开始节点
	startNode := e.findNodeByType(graph, model.NodeTypeStart)
	if startNode == nil {
		return nil, ErrStartNodeNotFound
	}

	// 6. 查找开始节点的下一个节点
	nextNodeID := e.findNextNodeID(graph, startNode.ID)
	if nextNodeID == "" {
		return nil, ErrNextNodeNotFound
	}

	nextNode := e.findNodeByID(graph, nextNodeID)
	if nextNode == nil {
		return nil, ErrNextNodeNotFound
	}

	// 7. 序列化表单数据
	formDataJSON := ""
	if req.FormData != nil {
		data, err := json.Marshal(req.FormData)
		if err != nil {
			return nil, fmt.Errorf("序列化表单数据失败: %w", err)
		}
		formDataJSON = string(data)
	}

	// 8. 创建流程实例和首个任务（事务）
	var instance *model.ProcessInstance
	err = e.processInstRepo.Transaction(func(tx *gorm.DB) error {
		// 创建流程实例
		instance = &model.ProcessInstance{
			ProcessDefID:      processDef.ID,
			ProcessDefVersion: processDef.Version,
			Title:             req.Title,
			InitiatorID:       req.InitiatorID,
			InitiatorName:     req.InitiatorName,
			Status:            model.ProcessInstStatusRunning,
			CurrentNodeID:     nextNodeID,
			FormData:          formDataJSON,
		}

		if err := e.processInstRepo.CreateWithTx(tx, instance); err != nil {
			return fmt.Errorf("创建流程实例失败: %w", err)
		}

		// 创建首个任务
		processCtx := &ProcessContext{
			InitiatorID:   req.InitiatorID,
			InitiatorName: req.InitiatorName,
			DeptID:        req.DeptID,
			FormData:      req.FormData,
		}

		if err := e.createTasksForNode(ctx, tx, instance, nextNode, graph, processCtx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// CompleteTask 完成任务（通过/拒绝）
// **Feature: oa-workflow-engine, Property 8: 任务通过后流转正确性**
// **Validates: Requirements 3.1**
func (e *workflowEngine) CompleteTask(ctx context.Context, req *CompleteTaskRequest) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(req.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != req.OperatorID {
		return ErrNotTaskAssignee
	}

	// 4. 获取流程实例
	instance, err := e.processInstRepo.FindByID(task.ProcessInstID)
	if err != nil {
		return fmt.Errorf("查询流程实例失败: %w", err)
	}

	// 5. 检查流程实例状态
	if instance.Status != model.ProcessInstStatusRunning {
		return ErrProcessNotRunning
	}

	// 6. 获取流程定义和流程图
	processDef, err := e.processDefRepo.FindByID(instance.ProcessDefID)
	if err != nil {
		return fmt.Errorf("查询流程定义失败: %w", err)
	}

	graph, err := e.serializer.DeserializeProcessGraph(processDef.GraphJSON)
	if err != nil {
		return fmt.Errorf("解析流程图失败: %w", err)
	}

	// 7. 获取当前节点
	currentNode := e.findNodeByID(graph, task.NodeID)
	if currentNode == nil {
		return fmt.Errorf("当前节点不存在: %s", task.NodeID)
	}

	// 8. 根据操作类型处理
	switch req.Action {
	case model.TaskLogActionApprove:
		return e.handleApprove(ctx, task, instance, currentNode, graph, req.Comment, processDef)
	case model.TaskLogActionReject:
		return e.handleReject(ctx, task, instance, currentNode, graph, req.Comment, processDef)
	default:
		return fmt.Errorf("不支持的操作类型: %s", req.Action)
	}
}

// handleApprove 处理通过操作
// **Feature: oa-workflow-engine, Property 7: 审批模式流转正确性**
// **Validates: Requirements 3.3, 3.4**
func (e *workflowEngine) handleApprove(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance, currentNode *model.ProcessNode, graph *model.ProcessGraph, comment string, processDef *model.ProcessDefinition) error {
	return e.taskRepo.Transaction(func(tx *gorm.DB) error {
		// 1. 更新任务状态为已通过
		if err := e.taskRepo.UpdateStatusWithTx(tx, task.ID, model.WfTaskStatusApproved, comment); err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 2. 记录任务日志
		log := &model.TaskLog{
			TaskID:        task.ID,
			ProcessInstID: instance.ID,
			OperatorID:    task.AssigneeID,
			OperatorName:  task.AssigneeName,
			Action:        model.TaskLogActionApprove,
			Comment:       comment,
		}
		if err := e.taskRepo.CreateTaskLogWithTx(tx, log); err != nil {
			return fmt.Errorf("创建任务日志失败: %w", err)
		}

		// 3. 检查审批模式，决定是否流转
		shouldTransition := e.shouldTransitionToNextNode(ctx, tx, instance, currentNode)
		if !shouldTransition {
			return nil // 会签模式下还有其他人未审批，不流转
		}

		// 4. 查找下一个节点
		nextNodeID := e.findNextNodeID(graph, currentNode.ID)
		if nextNodeID == "" {
			return fmt.Errorf("找不到下一节点")
		}

		nextNode := e.findNodeByID(graph, nextNodeID)
		if nextNode == nil {
			return fmt.Errorf("下一节点不存在: %s", nextNodeID)
		}

		// 5. 处理下一节点
		return e.processNextNode(ctx, tx, instance, nextNode, graph, processDef)
	})
}

// handleReject 处理拒绝操作
// **Feature: oa-workflow-engine, Property 9: 流程状态终态一致性**
// **Validates: Requirements 2.3, 2.4, 2.5**
func (e *workflowEngine) handleReject(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance, currentNode *model.ProcessNode, graph *model.ProcessGraph, comment string, processDef *model.ProcessDefinition) error {
	err := e.taskRepo.Transaction(func(tx *gorm.DB) error {
		// 1. 更新任务状态为已拒绝
		if err := e.taskRepo.UpdateStatusWithTx(tx, task.ID, model.WfTaskStatusRejected, comment); err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 2. 记录任务日志
		log := &model.TaskLog{
			TaskID:        task.ID,
			ProcessInstID: instance.ID,
			OperatorID:    task.AssigneeID,
			OperatorName:  task.AssigneeName,
			Action:        model.TaskLogActionReject,
			Comment:       comment,
		}
		if err := e.taskRepo.CreateTaskLogWithTx(tx, log); err != nil {
			return fmt.Errorf("创建任务日志失败: %w", err)
		}

		// 3. 根据拒绝处理方式决定流程走向
		rejectAction := model.RejectActionTerminate // 默认终止
		if currentNode.Properties != nil && currentNode.Properties.RejectAction != "" {
			rejectAction = currentNode.Properties.RejectAction
		}

		switch rejectAction {
		case model.RejectActionTerminate:
			// 终止流程
			return e.terminateProcess(tx, instance, model.ProcessInstStatusRejected)
		case model.RejectActionReturnPrev:
			// 退回上一节点 - 简化实现：终止流程
			// TODO: 实现退回上一节点逻辑
			return e.terminateProcess(tx, instance, model.ProcessInstStatusRejected)
		default:
			return e.terminateProcess(tx, instance, model.ProcessInstStatusRejected)
		}
	})

	// 发送流程拒绝通知（事务成功后）
	// Requirements: 6.2, 6.4 - 状态变更时通知发起人
	if err == nil && e.notificationService != nil {
		go func(inst *model.ProcessInstance, rejectComment string) {
			_ = e.notificationService.NotifyProcessRejected(ctx, inst, rejectComment)
		}(instance, comment)
	}

	return err
}

// DelegateTask 委托任务
// **Feature: oa-workflow-engine, Property 10: 委托/转办记录完整性**
// **Validates: Requirements 3.5, 3.6**
func (e *workflowEngine) DelegateTask(ctx context.Context, req *DelegateTaskRequest) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(req.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != req.OperatorID {
		return ErrNotTaskAssignee
	}

	// 4. 获取流程实例并检查状态
	instance, err := e.processInstRepo.FindByID(task.ProcessInstID)
	if err != nil {
		return fmt.Errorf("查询流程实例失败: %w", err)
	}
	if instance.Status != model.ProcessInstStatusRunning {
		return ErrProcessNotRunning
	}

	// 5. 执行委托（事务）
	err = e.taskRepo.Transaction(func(tx *gorm.DB) error {
		// 记录原审批人ID
		originalAssigneeID := task.AssigneeID

		// 更新任务：设置委托来源，更新审批人，状态保持pending
		task.DelegatedFrom = &originalAssigneeID
		task.AssigneeID = req.ToUserID
		task.AssigneeName = req.ToUserName
		// 委托后任务状态仍为pending，新审批人可以继续处理
		if err := e.taskRepo.UpdateWithTx(tx, task); err != nil {
			return fmt.Errorf("更新任务失败: %w", err)
		}

		// 记录委托日志
		log := &model.TaskLog{
			TaskID:        task.ID,
			ProcessInstID: instance.ID,
			OperatorID:    req.OperatorID,
			OperatorName:  req.OperatorName,
			Action:        model.TaskLogActionDelegate,
			Comment:       fmt.Sprintf("委托给 %s，原因：%s", req.ToUserName, req.Reason),
		}
		if err := e.taskRepo.CreateTaskLogWithTx(tx, log); err != nil {
			return fmt.Errorf("创建任务日志失败: %w", err)
		}

		return nil
	})

	// 发送新任务通知给被委托人（事务成功后）
	// Requirements: 6.1 - 新任务创建时推送通知
	if err == nil && e.notificationService != nil {
		go func(t *model.WfTask, inst *model.ProcessInstance) {
			_ = e.notificationService.NotifyNewTask(ctx, t, inst)
		}(task, instance)
	}

	return err
}

// TransferTask 转办任务
// **Feature: oa-workflow-engine, Property 10: 委托/转办记录完整性**
// **Validates: Requirements 3.5, 3.6**
func (e *workflowEngine) TransferTask(ctx context.Context, req *TransferTaskRequest) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(req.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != req.OperatorID {
		return ErrNotTaskAssignee
	}

	// 4. 获取流程实例并检查状态
	instance, err := e.processInstRepo.FindByID(task.ProcessInstID)
	if err != nil {
		return fmt.Errorf("查询流程实例失败: %w", err)
	}
	if instance.Status != model.ProcessInstStatusRunning {
		return ErrProcessNotRunning
	}

	// 5. 执行转办（事务）
	err = e.taskRepo.Transaction(func(tx *gorm.DB) error {
		// 记录原审批人ID
		originalAssigneeID := task.AssigneeID

		// 更新任务：设置转办来源，更新审批人，状态保持pending
		task.TransferredFrom = &originalAssigneeID
		task.AssigneeID = req.ToUserID
		task.AssigneeName = req.ToUserName
		// 转办后任务状态仍为pending，新审批人可以继续处理
		if err := e.taskRepo.UpdateWithTx(tx, task); err != nil {
			return fmt.Errorf("更新任务失败: %w", err)
		}

		// 记录转办日志
		log := &model.TaskLog{
			TaskID:        task.ID,
			ProcessInstID: instance.ID,
			OperatorID:    req.OperatorID,
			OperatorName:  req.OperatorName,
			Action:        model.TaskLogActionTransfer,
			Comment:       fmt.Sprintf("转办给 %s，原因：%s", req.ToUserName, req.Reason),
		}
		if err := e.taskRepo.CreateTaskLogWithTx(tx, log); err != nil {
			return fmt.Errorf("创建任务日志失败: %w", err)
		}

		return nil
	})

	// 发送新任务通知给被转办人（事务成功后）
	// Requirements: 6.1 - 新任务创建时推送通知
	if err == nil && e.notificationService != nil {
		go func(t *model.WfTask, inst *model.ProcessInstance) {
			_ = e.notificationService.NotifyNewTask(ctx, t, inst)
		}(task, instance)
	}

	return err
}

// WithdrawProcess 撤回流程
// **Feature: oa-workflow-engine, Property 9: 流程状态终态一致性**
// **Validates: Requirements 2.3, 2.4, 2.5**
func (e *workflowEngine) WithdrawProcess(ctx context.Context, instanceID int64, operatorID int64) error {
	// 1. 获取流程实例
	instance, err := e.processInstRepo.FindByID(instanceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProcessInstNotFound
		}
		return fmt.Errorf("查询流程实例失败: %w", err)
	}

	// 2. 检查流程实例状态
	if instance.Status != model.ProcessInstStatusRunning {
		return ErrProcessNotRunning
	}

	// 3. 检查操作人是否是发起人
	if instance.InitiatorID != operatorID {
		return ErrNotProcessInitiator
	}

	// 4. 执行撤回（事务）
	return e.processInstRepo.Transaction(func(tx *gorm.DB) error {
		// 取消所有待处理任务
		if err := e.taskRepo.CancelPendingByProcessInstIDWithTx(tx, instanceID); err != nil {
			return fmt.Errorf("取消待处理任务失败: %w", err)
		}

		// 更新流程实例状态
		return e.terminateProcessWithTx(tx, instance, model.ProcessInstStatusWithdrawn)
	})
}

// GetProcessInstance 获取流程实例详情
func (e *workflowEngine) GetProcessInstance(ctx context.Context, instanceID int64) (*ProcessInstanceDetail, error) {
	// 1. 获取流程实例
	instance, err := e.processInstRepo.FindByID(instanceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProcessInstNotFound
		}
		return nil, fmt.Errorf("查询流程实例失败: %w", err)
	}

	// 2. 获取所有任务
	tasks, err := e.taskRepo.FindByProcessInstID(instanceID)
	if err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}

	// 3. 获取审批日志
	logs, err := e.taskRepo.FindTaskLogsByProcessInstID(instanceID)
	if err != nil {
		return nil, fmt.Errorf("查询审批日志失败: %w", err)
	}

	// 转换为指针切片
	taskPtrs := make([]*model.WfTask, len(tasks))
	for i := range tasks {
		taskPtrs[i] = &tasks[i]
	}

	logPtrs := make([]*model.TaskLog, len(logs))
	for i := range logs {
		logPtrs[i] = &logs[i]
	}

	return &ProcessInstanceDetail{
		Instance:     instance,
		Tasks:        taskPtrs,
		ApprovalLogs: logPtrs,
	}, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// findNodeByType 根据类型查找节点
func (e *workflowEngine) findNodeByType(graph *model.ProcessGraph, nodeType model.NodeType) *model.ProcessNode {
	for i := range graph.Nodes {
		if graph.Nodes[i].Type == nodeType {
			return &graph.Nodes[i]
		}
	}
	return nil
}

// findNodeByID 根据ID查找节点
func (e *workflowEngine) findNodeByID(graph *model.ProcessGraph, nodeID string) *model.ProcessNode {
	for i := range graph.Nodes {
		if graph.Nodes[i].ID == nodeID {
			return &graph.Nodes[i]
		}
	}
	return nil
}

// findNextNodeID 查找下一个节点ID
func (e *workflowEngine) findNextNodeID(graph *model.ProcessGraph, sourceNodeID string) string {
	for _, edge := range graph.Edges {
		if edge.Source == sourceNodeID {
			return edge.Target
		}
	}
	return ""
}

// findPrevNodeID 查找上一个节点ID
func (e *workflowEngine) findPrevNodeID(graph *model.ProcessGraph, targetNodeID string) string {
	for _, edge := range graph.Edges {
		if edge.Target == targetNodeID {
			return edge.Source
		}
	}
	return ""
}

// createTasksForNode 为节点创建任务
func (e *workflowEngine) createTasksForNode(ctx context.Context, tx *gorm.DB, instance *model.ProcessInstance, node *model.ProcessNode, graph *model.ProcessGraph, processCtx *ProcessContext) error {
	// 如果是结束节点，完成流程
	if node.Type == model.NodeTypeEnd {
		if err := e.completeProcessWithTx(tx, instance); err != nil {
			return err
		}
		// 发送流程完成通知（异步，不影响事务）
		// Requirements: 6.4 - WHEN 审批完成 THEN Notification_Service SHALL 通知发起人最终审批结果
		if e.notificationService != nil {
			go func(inst *model.ProcessInstance) {
				_ = e.notificationService.NotifyProcessCompleted(ctx, inst)
			}(instance)
		}
		return nil
	}

	// 如果不是审批节点，直接流转到下一节点
	if node.Type != model.NodeTypeApproval {
		nextNodeID := e.findNextNodeID(graph, node.ID)
		if nextNodeID == "" {
			return fmt.Errorf("找不到下一节点")
		}
		nextNode := e.findNodeByID(graph, nextNodeID)
		if nextNode == nil {
			return fmt.Errorf("下一节点不存在: %s", nextNodeID)
		}
		// 更新当前节点
		instance.CurrentNodeID = nextNodeID
		if err := e.processInstRepo.UpdateWithTx(tx, instance); err != nil {
			return fmt.Errorf("更新流程实例失败: %w", err)
		}
		return e.createTasksForNode(ctx, tx, instance, nextNode, graph, processCtx)
	}

	// 审批节点：解析审批人并创建任务
	if node.Properties == nil || node.Properties.AssigneeRule == nil {
		return fmt.Errorf("审批节点 %s 未配置审批人规则", node.ID)
	}

	// 解析审批人
	assigneeIDs, err := e.assigneeResolver.Resolve(ctx, node.Properties.AssigneeRule, processCtx)
	if err != nil {
		return fmt.Errorf("解析审批人失败: %w", err)
	}

	// 计算超时时间
	var dueAt *time.Time
	if node.Properties.TimeoutHours > 0 {
		due := time.Now().Add(time.Duration(node.Properties.TimeoutHours) * time.Hour)
		dueAt = &due
	}

	// 为每个审批人创建任务
	tasks := make([]*model.WfTask, 0, len(assigneeIDs))
	for _, assigneeID := range assigneeIDs {
		task := &model.WfTask{
			ProcessInstID: instance.ID,
			NodeID:        node.ID,
			NodeName:      node.Name,
			AssigneeID:    assigneeID,
			Status:        model.WfTaskStatusPending,
			DueAt:         dueAt,
		}
		tasks = append(tasks, task)
	}

	if err := e.taskRepo.BatchCreateWithTx(tx, tasks); err != nil {
		return fmt.Errorf("创建任务失败: %w", err)
	}

	// 更新当前节点
	instance.CurrentNodeID = node.ID
	if err := e.processInstRepo.UpdateWithTx(tx, instance); err != nil {
		return fmt.Errorf("更新流程实例失败: %w", err)
	}

	// 发送新任务通知（异步，不影响事务）
	// Requirements: 6.1 - WHEN 新Task创建 THEN Notification_Service SHALL 通过WebSocket向Assignee推送待办提醒
	if e.notificationService != nil {
		go func(tasks []*model.WfTask, inst *model.ProcessInstance) {
			_ = e.notificationService.NotifyNewTasks(ctx, tasks, inst)
		}(tasks, instance)
	}

	return nil
}

// shouldTransitionToNextNode 判断是否应该流转到下一节点
// 根据审批模式（会签/或签）决定
func (e *workflowEngine) shouldTransitionToNextNode(ctx context.Context, tx *gorm.DB, instance *model.ProcessInstance, node *model.ProcessNode) bool {
	// 获取审批模式，默认为或签
	approvalMode := model.ApprovalModeOrSign
	if node.Properties != nil && node.Properties.ApprovalMode != "" {
		approvalMode = node.Properties.ApprovalMode
	}

	switch approvalMode {
	case model.ApprovalModeOrSign:
		// 或签：任一人通过即流转
		return true
	case model.ApprovalModeAndSign:
		// 会签：所有人通过才流转
		// 检查是否还有待处理任务
		pendingCount, err := e.taskRepo.CountPendingByProcessInstIDAndNodeID(instance.ID, node.ID)
		if err != nil {
			return false
		}
		// 如果没有待处理任务了，说明所有人都已处理
		return pendingCount == 0
	default:
		return true
	}
}

// processNextNode 处理下一节点
func (e *workflowEngine) processNextNode(ctx context.Context, tx *gorm.DB, instance *model.ProcessInstance, nextNode *model.ProcessNode, graph *model.ProcessGraph, processDef *model.ProcessDefinition) error {
	// 构建流程上下文
	var formData map[string]interface{}
	if instance.FormData != "" {
		if err := json.Unmarshal([]byte(instance.FormData), &formData); err != nil {
			return fmt.Errorf("解析表单数据失败: %w", err)
		}
	}

	processCtx := &ProcessContext{
		InitiatorID:   instance.InitiatorID,
		InitiatorName: instance.InitiatorName,
		FormData:      formData,
	}

	return e.createTasksForNode(ctx, tx, instance, nextNode, graph, processCtx)
}

// terminateProcess 终止流程
func (e *workflowEngine) terminateProcess(tx *gorm.DB, instance *model.ProcessInstance, status string) error {
	return e.terminateProcessWithTx(tx, instance, status)
}

// terminateProcessWithTx 在事务中终止流程
func (e *workflowEngine) terminateProcessWithTx(tx *gorm.DB, instance *model.ProcessInstance, status string) error {
	now := time.Now()
	instance.Status = status
	instance.CompletedAt = &now
	return e.processInstRepo.UpdateWithTx(tx, instance)
}

// completeProcessWithTx 在事务中完成流程
func (e *workflowEngine) completeProcessWithTx(tx *gorm.DB, instance *model.ProcessInstance) error {
	now := time.Now()
	instance.Status = model.ProcessInstStatusCompleted
	instance.CompletedAt = &now
	return e.processInstRepo.UpdateWithTx(tx, instance)
}
