package workflow

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"context"
	"time"
)

// ============================================================================
// Request/Response Types
// ============================================================================

// QueryRequest 查询请求
type QueryRequest struct {
	Current      int                    `json:"current" form:"current"`           // 当前页码
	Size         int                    `json:"size" form:"size"`                 // 每页大小
	Title        string                 `json:"title" form:"title"`               // 标题搜索
	Status       string                 `json:"status" form:"status"`             // 状态筛选
	ProcessDefID *int64                 `json:"processDefId" form:"processDefId"` // 流程定义ID
	NodeName     string                 `json:"nodeName" form:"nodeName"`         // 节点名称
	StartTime    *time.Time             `json:"startTime" form:"startTime"`       // 开始时间
	EndTime      *time.Time             `json:"endTime" form:"endTime"`           // 结束时间
	Extra        map[string]interface{} `json:"extra"`                            // 额外查询条件
}

// ProcessInstanceSummary 流程实例摘要
type ProcessInstanceSummary struct {
	ID                int64      `json:"id"`
	ProcessDefID      int64      `json:"processDefId"`
	ProcessDefVersion int        `json:"processDefVersion"`
	ProcessDefName    string     `json:"processDefName"`
	Title             string     `json:"title"`
	InitiatorID       int64      `json:"initiatorId"`
	InitiatorName     string     `json:"initiatorName"`
	Status            string     `json:"status"`
	CurrentNodeID     string     `json:"currentNodeId"`
	CurrentNodeName   string     `json:"currentNodeName"`
	StartedAt         time.Time  `json:"startedAt"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
}

// TaskSummary 任务摘要
type TaskSummary struct {
	ID              int64      `json:"id"`
	ProcessInstID   int64      `json:"processInstId"`
	ProcessTitle    string     `json:"processTitle"`
	ProcessDefName  string     `json:"processDefName"`
	NodeID          string     `json:"nodeId"`
	NodeName        string     `json:"nodeName"`
	AssigneeID      int64      `json:"assigneeId"`
	AssigneeName    string     `json:"assigneeName"`
	Status          string     `json:"status"`
	Comment         string     `json:"comment"`
	DelegatedFrom   *int64     `json:"delegatedFrom,omitempty"`
	TransferredFrom *int64     `json:"transferredFrom,omitempty"`
	DueAt           *time.Time `json:"dueAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
}

// ApprovalRecord 审批记录
type ApprovalRecord struct {
	ID           int64     `json:"id"`
	TaskID       int64     `json:"taskId"`
	NodeID       string    `json:"nodeId"`
	NodeName     string    `json:"nodeName"`
	OperatorID   int64     `json:"operatorId"`
	OperatorName string    `json:"operatorName"`
	Action       string    `json:"action"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"createdAt"`
}

// StatisticsRequest 统计请求
type StatisticsRequest struct {
	ProcessDefID *int64     `json:"processDefId" form:"processDefId"` // 流程定义ID
	StartTime    *time.Time `json:"startTime" form:"startTime"`       // 开始时间
	EndTime      *time.Time `json:"endTime" form:"endTime"`           // 结束时间
}

// ProcessStatistics 流程统计
type ProcessStatistics struct {
	TotalCount     int64                  `json:"totalCount"`     // 总数
	RunningCount   int64                  `json:"runningCount"`   // 运行中数量
	CompletedCount int64                  `json:"completedCount"` // 已完成数量
	RejectedCount  int64                  `json:"rejectedCount"`  // 已拒绝数量
	WithdrawnCount int64                  `json:"withdrawnCount"` // 已撤回数量
	AvgDuration    float64                `json:"avgDuration"`    // 平均耗时(小时)
	ByProcessDef   []ProcessDefStatistics `json:"byProcessDef"`   // 按流程定义分组统计
}

// ProcessDefStatistics 按流程定义分组的统计
type ProcessDefStatistics struct {
	ProcessDefID   int64   `json:"processDefId"`
	ProcessDefName string  `json:"processDefName"`
	TotalCount     int64   `json:"totalCount"`
	CompletedCount int64   `json:"completedCount"`
	AvgDuration    float64 `json:"avgDuration"` // 平均耗时(小时)
}

// ============================================================================
// Interface
// ============================================================================

// QueryService 查询服务接口
type QueryService interface {
	// ListMyInitiated 查询我发起的流程
	// Requirements: 7.1 - WHEN 用户查询我发起的申请 THEN Query_Service SHALL 返回该用户作为发起人的所有Process_Instance列表
	ListMyInitiated(ctx context.Context, userID int64, req *QueryRequest) ([]*ProcessInstanceSummary, int64, error)

	// ListMyTodo 查询我的待办
	// Requirements: 7.2 - WHEN 用户查询我的待办 THEN Query_Service SHALL 返回分配给该用户且状态为待处理的Task列表
	ListMyTodo(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error)

	// ListMyDone 查询我的已办
	// Requirements: 7.3 - WHEN 用户查询我的已办 THEN Query_Service SHALL 返回该用户已处理的Task历史记录
	ListMyDone(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error)

	// GetApprovalTrail 获取审批轨迹
	// Requirements: 7.4 - WHEN 用户查看流程详情 THEN Query_Service SHALL 返回Process_Instance的完整审批轨迹和表单数据
	GetApprovalTrail(ctx context.Context, instanceID int64) ([]*ApprovalRecord, error)

	// GetStatistics 获取流程统计
	// Requirements: 7.5 - WHEN 管理员查询流程统计 THEN Query_Service SHALL 返回按流程类型、时间段分组的审批数量和平均耗时
	GetStatistics(ctx context.Context, req *StatisticsRequest) (*ProcessStatistics, error)
}

// ============================================================================
// Implementation
// ============================================================================

// queryService 查询服务实现
type queryService struct {
	processDefRepo  *repository.ProcessDefinitionRepository
	processInstRepo *repository.ProcessInstanceRepository
	taskRepo        *repository.WfTaskRepository
}

// NewQueryService 创建查询服务实例
func NewQueryService(
	processDefRepo *repository.ProcessDefinitionRepository,
	processInstRepo *repository.ProcessInstanceRepository,
	taskRepo *repository.WfTaskRepository,
) QueryService {
	return &queryService{
		processDefRepo:  processDefRepo,
		processInstRepo: processInstRepo,
		taskRepo:        taskRepo,
	}
}

// ListMyInitiated 查询我发起的流程
// **Feature: oa-workflow-engine, Property 13: 查询结果过滤正确性**
// **Validates: Requirements 7.1, 7.2, 7.3**
func (s *queryService) ListMyInitiated(ctx context.Context, userID int64, req *QueryRequest) ([]*ProcessInstanceSummary, int64, error) {
	// 设置默认分页参数
	current, size := s.normalizePagination(req.Current, req.Size)

	// 构建查询条件
	query := make(map[string]interface{})
	if req.Title != "" {
		query["title"] = req.Title
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	if req.ProcessDefID != nil {
		query["processDefId"] = *req.ProcessDefID
	}

	// 查询流程实例
	instances, total, err := s.processInstRepo.FindByInitiator(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为摘要
	summaries := make([]*ProcessInstanceSummary, 0, len(instances))
	for _, inst := range instances {
		summary := &ProcessInstanceSummary{
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

		// 获取流程定义名称
		if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
			summary.ProcessDefName = processDef.Name
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// ListMyTodo 查询我的待办
// **Feature: oa-workflow-engine, Property 13: 查询结果过滤正确性**
// **Validates: Requirements 7.1, 7.2, 7.3**
func (s *queryService) ListMyTodo(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error) {
	// 设置默认分页参数
	current, size := s.normalizePagination(req.Current, req.Size)

	// 构建查询条件
	query := make(map[string]interface{})
	if req.NodeName != "" {
		query["nodeName"] = req.NodeName
	}
	if req.ProcessDefID != nil {
		query["processInstId"] = *req.ProcessDefID
	}

	// 查询待办任务
	tasks, total, err := s.taskRepo.FindPendingByAssignee(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为摘要
	summaries := make([]*TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		summary := s.taskToSummary(&task)

		// 获取流程实例信息
		if inst, err := s.processInstRepo.FindByID(task.ProcessInstID); err == nil {
			summary.ProcessTitle = inst.Title
			// 获取流程定义名称
			if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
				summary.ProcessDefName = processDef.Name
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// ListMyDone 查询我的已办
// **Feature: oa-workflow-engine, Property 13: 查询结果过滤正确性**
// **Validates: Requirements 7.1, 7.2, 7.3**
func (s *queryService) ListMyDone(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error) {
	// 设置默认分页参数
	current, size := s.normalizePagination(req.Current, req.Size)

	// 构建查询条件
	query := make(map[string]interface{})
	if req.NodeName != "" {
		query["nodeName"] = req.NodeName
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	if req.ProcessDefID != nil {
		query["processInstId"] = *req.ProcessDefID
	}

	// 查询已办任务
	tasks, total, err := s.taskRepo.FindCompletedByAssignee(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

	// 转换为摘要
	summaries := make([]*TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		summary := s.taskToSummary(&task)

		// 获取流程实例信息
		if inst, err := s.processInstRepo.FindByID(task.ProcessInstID); err == nil {
			summary.ProcessTitle = inst.Title
			// 获取流程定义名称
			if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
				summary.ProcessDefName = processDef.Name
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// GetApprovalTrail 获取审批轨迹
// Requirements: 7.4
func (s *queryService) GetApprovalTrail(ctx context.Context, instanceID int64) ([]*ApprovalRecord, error) {
	// 获取任务日志
	logs, err := s.taskRepo.FindTaskLogsByProcessInstID(instanceID)
	if err != nil {
		return nil, err
	}

	// 获取所有任务以获取节点信息
	tasks, err := s.taskRepo.FindByProcessInstID(instanceID)
	if err != nil {
		return nil, err
	}

	// 构建任务ID到节点信息的映射
	taskNodeMap := make(map[int64]struct {
		NodeID   string
		NodeName string
	})
	for _, task := range tasks {
		taskNodeMap[task.ID] = struct {
			NodeID   string
			NodeName string
		}{
			NodeID:   task.NodeID,
			NodeName: task.NodeName,
		}
	}

	// 转换为审批记录
	records := make([]*ApprovalRecord, 0, len(logs))
	for _, log := range logs {
		record := &ApprovalRecord{
			ID:           log.ID,
			TaskID:       log.TaskID,
			OperatorID:   log.OperatorID,
			OperatorName: log.OperatorName,
			Action:       log.Action,
			Comment:      log.Comment,
			CreatedAt:    log.CreatedAt,
		}

		// 填充节点信息
		if nodeInfo, ok := taskNodeMap[log.TaskID]; ok {
			record.NodeID = nodeInfo.NodeID
			record.NodeName = nodeInfo.NodeName
		}

		records = append(records, record)
	}

	return records, nil
}

// GetStatistics 获取流程统计
// Requirements: 7.5
func (s *queryService) GetStatistics(ctx context.Context, req *StatisticsRequest) (*ProcessStatistics, error) {
	stats := &ProcessStatistics{
		ByProcessDef: make([]ProcessDefStatistics, 0),
	}

	// 统计各状态数量
	runningCount, err := s.processInstRepo.CountByStatus(model.ProcessInstStatusRunning)
	if err != nil {
		return nil, err
	}
	stats.RunningCount = runningCount

	completedCount, err := s.processInstRepo.CountByStatus(model.ProcessInstStatusCompleted)
	if err != nil {
		return nil, err
	}
	stats.CompletedCount = completedCount

	rejectedCount, err := s.processInstRepo.CountByStatus(model.ProcessInstStatusRejected)
	if err != nil {
		return nil, err
	}
	stats.RejectedCount = rejectedCount

	withdrawnCount, err := s.processInstRepo.CountByStatus(model.ProcessInstStatusWithdrawn)
	if err != nil {
		return nil, err
	}
	stats.WithdrawnCount = withdrawnCount

	stats.TotalCount = runningCount + completedCount + rejectedCount + withdrawnCount

	// 按流程定义分组统计
	processDefs, _, err := s.processDefRepo.FindWithPagination(nil, 1, 100)
	if err != nil {
		return nil, err
	}

	for _, def := range processDefs {
		defStats := ProcessDefStatistics{
			ProcessDefID:   def.ID,
			ProcessDefName: def.Name,
		}

		// 统计该流程定义的实例数
		count, err := s.processInstRepo.CountByProcessDefID(def.ID)
		if err != nil {
			continue
		}
		defStats.TotalCount = count

		if count > 0 {
			stats.ByProcessDef = append(stats.ByProcessDef, defStats)
		}
	}

	return stats, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// normalizePagination 规范化分页参数
func (s *queryService) normalizePagination(current, size int) (int, int) {
	if current <= 0 {
		current = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return current, size
}

// taskToSummary 将任务转换为摘要
func (s *queryService) taskToSummary(task *model.WfTask) *TaskSummary {
	return &TaskSummary{
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
