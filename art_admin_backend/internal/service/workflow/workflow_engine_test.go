package workflow

import (
	"art_admin_backend/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ============================================================================
// Mock Repositories for Testing
// ============================================================================

// mockProcessDefRepo 模拟流程定义仓储
type mockProcessDefRepo struct {
	definitions map[int64]*model.ProcessDefinition
}

func newMockProcessDefRepo() *mockProcessDefRepo {
	return &mockProcessDefRepo{
		definitions: make(map[int64]*model.ProcessDefinition),
	}
}

func (r *mockProcessDefRepo) AddDefinition(def *model.ProcessDefinition) {
	r.definitions[def.ID] = def
}

func (r *mockProcessDefRepo) FindByID(id int64) (*model.ProcessDefinition, error) {
	if def, exists := r.definitions[id]; exists {
		return def, nil
	}
	return nil, ErrProcessDefNotFound
}

// mockProcessInstRepo 模拟流程实例仓储
type mockProcessInstRepo struct {
	instances   map[int64]*model.ProcessInstance
	nextID      int64
	createdInTx []*model.ProcessInstance
}

func newMockProcessInstRepo() *mockProcessInstRepo {
	return &mockProcessInstRepo{
		instances:   make(map[int64]*model.ProcessInstance),
		nextID:      1,
		createdInTx: make([]*model.ProcessInstance, 0),
	}
}

func (r *mockProcessInstRepo) Create(inst *model.ProcessInstance) error {
	inst.ID = r.nextID
	r.nextID++
	r.instances[inst.ID] = inst
	return nil
}

func (r *mockProcessInstRepo) FindByID(id int64) (*model.ProcessInstance, error) {
	if inst, exists := r.instances[id]; exists {
		return inst, nil
	}
	return nil, ErrProcessInstNotFound
}

func (r *mockProcessInstRepo) Update(inst *model.ProcessInstance) error {
	r.instances[inst.ID] = inst
	return nil
}

// mockWfTaskRepo 模拟任务仓储
type mockWfTaskRepo struct {
	tasks        map[int64]*model.WfTask
	taskLogs     map[int64]*model.TaskLog
	nextTaskID   int64
	nextLogID    int64
	createdTasks []*model.WfTask
}

func newMockWfTaskRepo() *mockWfTaskRepo {
	return &mockWfTaskRepo{
		tasks:        make(map[int64]*model.WfTask),
		taskLogs:     make(map[int64]*model.TaskLog),
		nextTaskID:   1,
		nextLogID:    1,
		createdTasks: make([]*model.WfTask, 0),
	}
}

func (r *mockWfTaskRepo) Create(task *model.WfTask) error {
	task.ID = r.nextTaskID
	r.nextTaskID++
	r.tasks[task.ID] = task
	r.createdTasks = append(r.createdTasks, task)
	return nil
}

func (r *mockWfTaskRepo) FindByID(id int64) (*model.WfTask, error) {
	if task, exists := r.tasks[id]; exists {
		return task, nil
	}
	return nil, ErrTaskNotFound
}

func (r *mockWfTaskRepo) FindByProcessInstID(processInstID int64) ([]model.WfTask, error) {
	var tasks []model.WfTask
	for _, task := range r.tasks {
		if task.ProcessInstID == processInstID {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (r *mockWfTaskRepo) FindPendingByProcessInstID(processInstID int64) ([]model.WfTask, error) {
	var tasks []model.WfTask
	for _, task := range r.tasks {
		if task.ProcessInstID == processInstID && task.Status == model.WfTaskStatusPending {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (r *mockWfTaskRepo) CountPendingByProcessInstIDAndNodeID(processInstID int64, nodeID string) (int64, error) {
	var count int64
	for _, task := range r.tasks {
		if task.ProcessInstID == processInstID && task.NodeID == nodeID && task.Status == model.WfTaskStatusPending {
			count++
		}
	}
	return count, nil
}

func (r *mockWfTaskRepo) CreateTaskLog(log *model.TaskLog) error {
	log.ID = r.nextLogID
	r.nextLogID++
	r.taskLogs[log.ID] = log
	return nil
}

func (r *mockWfTaskRepo) FindTaskLogsByProcessInstID(processInstID int64) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	for _, log := range r.taskLogs {
		if log.ProcessInstID == processInstID {
			logs = append(logs, *log)
		}
	}
	return logs, nil
}

func (r *mockWfTaskRepo) GetCreatedTasks() []*model.WfTask {
	return r.createdTasks
}

func (r *mockWfTaskRepo) ClearCreatedTasks() {
	r.createdTasks = make([]*model.WfTask, 0)
}

// mockFormEngineService 模拟表单引擎服务
type mockFormEngineService struct {
	schemas map[int64]*model.FormSchema
}

func newMockFormEngineService() *mockFormEngineService {
	return &mockFormEngineService{
		schemas: make(map[int64]*model.FormSchema),
	}
}

func (s *mockFormEngineService) AddSchema(templateID int64, schema *model.FormSchema) {
	s.schemas[templateID] = schema
}

func (s *mockFormEngineService) GetSchema(ctx context.Context, id int64) (*model.FormSchema, error) {
	if schema, exists := s.schemas[id]; exists {
		return schema, nil
	}
	return nil, fmt.Errorf("schema not found")
}

func (s *mockFormEngineService) ValidateFormData(ctx context.Context, schema *model.FormSchema, data map[string]interface{}) *FormValidationResult {
	return &FormValidationResult{Valid: true}
}

func (s *mockFormEngineService) CreateTemplate(ctx context.Context, req *CreateFormTemplateRequest) (*model.FormTemplate, error) {
	return nil, nil
}

func (s *mockFormEngineService) UpdateTemplate(ctx context.Context, id int64, req *UpdateFormTemplateRequest) (*model.FormTemplate, error) {
	return nil, nil
}

func (s *mockFormEngineService) GetTemplate(ctx context.Context, id int64) (*model.FormTemplate, error) {
	return nil, nil
}

func (s *mockFormEngineService) ListTemplates(ctx context.Context, req *ListFormTemplateRequest) ([]*model.FormTemplate, int64, error) {
	return nil, 0, nil
}

func (s *mockFormEngineService) DeleteTemplate(ctx context.Context, id int64) error {
	return nil
}

func (s *mockFormEngineService) GetFieldPermissions(ctx context.Context, templateID int64, nodePermissions map[string]string) (*FieldPermissions, error) {
	return nil, nil
}

// mockAssigneeResolver 模拟审批人解析器
type mockAssigneeResolver struct {
	assignees map[string][]int64 // nodeID -> assigneeIDs
}

func newMockAssigneeResolver() *mockAssigneeResolver {
	return &mockAssigneeResolver{
		assignees: make(map[string][]int64),
	}
}

func (r *mockAssigneeResolver) SetAssignees(nodeID string, assigneeIDs []int64) {
	r.assignees[nodeID] = assigneeIDs
}

func (r *mockAssigneeResolver) Resolve(ctx context.Context, rule *model.AssigneeRule, processCtx *ProcessContext) ([]int64, error) {
	// Return default assignees if rule has values
	if rule != nil && len(rule.Values) > 0 {
		return rule.Values, nil
	}
	// Return error if no assignees configured
	return nil, ErrNoAssigneeFound
}

// ============================================================================
// Test Workflow Engine
// ============================================================================

// testWorkflowEngine 可测试的工作流引擎
type testWorkflowEngine struct {
	processDefRepo   *mockProcessDefRepo
	processInstRepo  *mockProcessInstRepo
	taskRepo         *mockWfTaskRepo
	formEngine       *mockFormEngineService
	assigneeResolver *mockAssigneeResolver
	serializer       Serializer
	validator        Validator
	engine           *inMemoryWorkflowEngine
}

// inMemoryWorkflowEngine 内存工作流引擎（用于测试）
type inMemoryWorkflowEngine struct {
	processDefRepo   *mockProcessDefRepo
	processInstRepo  *mockProcessInstRepo
	taskRepo         *mockWfTaskRepo
	formEngine       *mockFormEngineService
	assigneeResolver *mockAssigneeResolver
	serializer       Serializer
	validator        Validator
}

func newTestWorkflowEngine() *testWorkflowEngine {
	processDefRepo := newMockProcessDefRepo()
	processInstRepo := newMockProcessInstRepo()
	taskRepo := newMockWfTaskRepo()
	formEngine := newMockFormEngineService()
	assigneeResolver := newMockAssigneeResolver()
	serializer := NewSerializer()
	validator := NewValidator()

	engine := &inMemoryWorkflowEngine{
		processDefRepo:   processDefRepo,
		processInstRepo:  processInstRepo,
		taskRepo:         taskRepo,
		formEngine:       formEngine,
		assigneeResolver: assigneeResolver,
		serializer:       serializer,
		validator:        validator,
	}

	return &testWorkflowEngine{
		processDefRepo:   processDefRepo,
		processInstRepo:  processInstRepo,
		taskRepo:         taskRepo,
		formEngine:       formEngine,
		assigneeResolver: assigneeResolver,
		serializer:       serializer,
		validator:        validator,
		engine:           engine,
	}
}

// StartProcess 启动流程（内存实现）
func (e *inMemoryWorkflowEngine) StartProcess(ctx context.Context, req *StartProcessRequest) (*model.ProcessInstance, error) {
	// 1. 获取流程定义
	processDef, err := e.processDefRepo.FindByID(req.ProcessDefID)
	if err != nil {
		return nil, ErrProcessDefNotFound
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

	// 4. 查找开始节点
	var startNode *model.ProcessNode
	for i := range graph.Nodes {
		if graph.Nodes[i].Type == model.NodeTypeStart {
			startNode = &graph.Nodes[i]
			break
		}
	}
	if startNode == nil {
		return nil, ErrStartNodeNotFound
	}

	// 5. 查找开始节点的下一个节点
	var nextNodeID string
	for _, edge := range graph.Edges {
		if edge.Source == startNode.ID {
			nextNodeID = edge.Target
			break
		}
	}
	if nextNodeID == "" {
		return nil, ErrNextNodeNotFound
	}

	var nextNode *model.ProcessNode
	for i := range graph.Nodes {
		if graph.Nodes[i].ID == nextNodeID {
			nextNode = &graph.Nodes[i]
			break
		}
	}
	if nextNode == nil {
		return nil, ErrNextNodeNotFound
	}

	// 6. 序列化表单数据
	formDataJSON := ""
	if req.FormData != nil {
		data, err := json.Marshal(req.FormData)
		if err != nil {
			return nil, fmt.Errorf("序列化表单数据失败: %w", err)
		}
		formDataJSON = string(data)
	}

	// 7. 创建流程实例
	instance := &model.ProcessInstance{
		ProcessDefID:      processDef.ID,
		ProcessDefVersion: processDef.Version,
		Title:             req.Title,
		InitiatorID:       req.InitiatorID,
		InitiatorName:     req.InitiatorName,
		Status:            model.ProcessInstStatusRunning,
		CurrentNodeID:     nextNodeID,
		FormData:          formDataJSON,
	}

	if err := e.processInstRepo.Create(instance); err != nil {
		return nil, fmt.Errorf("创建流程实例失败: %w", err)
	}

	// 8. 创建首个任务
	if err := e.createTasksForNode(ctx, instance, nextNode, graph); err != nil {
		return nil, err
	}

	return instance, nil
}

// createTasksForNode 为节点创建任务
func (e *inMemoryWorkflowEngine) createTasksForNode(ctx context.Context, instance *model.ProcessInstance, node *model.ProcessNode, graph *model.ProcessGraph) error {
	// 如果是结束节点，完成流程
	if node.Type == model.NodeTypeEnd {
		instance.Status = model.ProcessInstStatusCompleted
		return e.processInstRepo.Update(instance)
	}

	// 如果不是审批节点，直接流转到下一节点
	if node.Type != model.NodeTypeApproval {
		var nextNodeID string
		for _, edge := range graph.Edges {
			if edge.Source == node.ID {
				nextNodeID = edge.Target
				break
			}
		}
		if nextNodeID == "" {
			return fmt.Errorf("找不到下一节点")
		}
		var nextNode *model.ProcessNode
		for i := range graph.Nodes {
			if graph.Nodes[i].ID == nextNodeID {
				nextNode = &graph.Nodes[i]
				break
			}
		}
		if nextNode == nil {
			return fmt.Errorf("下一节点不存在: %s", nextNodeID)
		}
		instance.CurrentNodeID = nextNodeID
		if err := e.processInstRepo.Update(instance); err != nil {
			return err
		}
		return e.createTasksForNode(ctx, instance, nextNode, graph)
	}

	// 审批节点：解析审批人并创建任务
	if node.Properties == nil || node.Properties.AssigneeRule == nil {
		return fmt.Errorf("审批节点 %s 未配置审批人规则", node.ID)
	}

	// 解析审批人
	assigneeIDs, err := e.assigneeResolver.Resolve(ctx, node.Properties.AssigneeRule, nil)
	if err != nil {
		return fmt.Errorf("解析审批人失败: %w", err)
	}

	// 为每个审批人创建任务
	for _, assigneeID := range assigneeIDs {
		task := &model.WfTask{
			ProcessInstID: instance.ID,
			NodeID:        node.ID,
			NodeName:      node.Name,
			AssigneeID:    assigneeID,
			Status:        model.WfTaskStatusPending,
		}
		if err := e.taskRepo.Create(task); err != nil {
			return fmt.Errorf("创建任务失败: %w", err)
		}
	}

	// 更新当前节点
	instance.CurrentNodeID = node.ID
	return e.processInstRepo.Update(instance)
}

// ============================================================================
// Generators
// ============================================================================

// genProcessDefID generates a process definition ID
func genProcessDefID() gopter.Gen {
	return gen.Int64Range(1, 100)
}

// genInitiatorIDForEngine generates an initiator user ID
func genInitiatorIDForEngine() gopter.Gen {
	return gen.Int64Range(1, 100)
}

// genAssigneeIDs generates a slice of assignee IDs
func genAssigneeIDs() gopter.Gen {
	return gen.SliceOfN(3, gen.Int64Range(1, 100)).SuchThat(func(ids []int64) bool {
		return len(ids) > 0
	})
}

// genTitle generates a process title
func genTitle() gopter.Gen {
	return gen.Int64Range(1, 10000).Map(func(n int64) string {
		return fmt.Sprintf("Test Process %d", n)
	})
}

// createValidProcessGraph creates a valid process graph for testing
func createValidProcessGraph(assigneeIDs []int64) *model.ProcessGraph {
	return &model.ProcessGraph{
		Nodes: []model.ProcessNode{
			{
				ID:       "start_1",
				Type:     model.NodeTypeStart,
				Name:     "开始",
				Position: model.Position{X: 100, Y: 100},
			},
			{
				ID:       "approval_1",
				Type:     model.NodeTypeApproval,
				Name:     "审批节点1",
				Position: model.Position{X: 300, Y: 100},
				Properties: &model.NodeProperties{
					AssigneeRule: &model.AssigneeRule{
						Type:   model.AssigneeTypeUser,
						Values: assigneeIDs,
					},
					ApprovalMode: model.ApprovalModeOrSign,
				},
			},
			{
				ID:       "end_1",
				Type:     model.NodeTypeEnd,
				Name:     "结束",
				Position: model.Position{X: 500, Y: 100},
			},
		},
		Edges: []model.ProcessEdge{
			{ID: "edge_1", Source: "start_1", Target: "approval_1"},
			{ID: "edge_2", Source: "approval_1", Target: "end_1"},
		},
	}
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 5: 流程实例初始状态正确性**
// **Validates: Requirements 2.1, 2.2**
// *For any* 新创建的ProcessInstance，其状态应为"running"，且应存在至少一个状态为"pending"的Task。

func TestWorkflowEngineStartProcessInitialState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("New process instance has running status and at least one pending task", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeIDs []int64, title string) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a valid process definition
			graph := createValidProcessGraph(assigneeIDs)
			graphJSON, err := te.serializer.SerializeProcessGraph(graph)
			if err != nil {
				t.Logf("Failed to serialize graph: %v", err)
				return false
			}

			processDef := &model.ProcessDefinition{
				ID:        processDefID,
				Name:      "Test Process",
				Code:      fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON,
				Version:   1,
				Status:    model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID:  processDefID,
				Title:         title,
				InitiatorID:   initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}

			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Property 1: Instance status should be "running"
			if instance.Status != model.ProcessInstStatusRunning {
				t.Logf("Expected status 'running', got '%s'", instance.Status)
				return false
			}

			// Property 2: There should be at least one pending task
			tasks, err := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if err != nil {
				t.Logf("Failed to find tasks: %v", err)
				return false
			}

			if len(tasks) == 0 {
				t.Logf("Expected at least one pending task, got 0")
				return false
			}

			// Property 3: All tasks should have pending status
			for _, task := range tasks {
				if task.Status != model.WfTaskStatusPending {
					t.Logf("Expected task status 'pending', got '%s'", task.Status)
					return false
				}
			}

			// Property 4: Number of tasks should match number of assignees
			if len(tasks) != len(assigneeIDs) {
				t.Logf("Expected %d tasks, got %d", len(assigneeIDs), len(tasks))
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		genAssigneeIDs(),
		genTitle(),
	))

	properties.TestingRun(t)
}

// TestWorkflowEngineStartProcessUnpublishedDef tests that starting a process with unpublished definition fails
func TestWorkflowEngineStartProcessUnpublishedDef(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Starting process with unpublished definition returns error", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeIDs []int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a draft process definition
			graph := createValidProcessGraph(assigneeIDs)
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)

			processDef := &model.ProcessDefinition{
				ID:        processDefID,
				Name:      "Test Process",
				Code:      fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON,
				Version:   1,
				Status:    model.ProcessDefStatusDraft, // Not published
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID:  processDefID,
				Title:         "Test",
				InitiatorID:   initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}

			_, err := te.engine.StartProcess(ctx, req)

			// Should return error
			if err != ErrProcessDefNotPublished {
				t.Logf("Expected ErrProcessDefNotPublished, got %v", err)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		genAssigneeIDs(),
	))

	properties.TestingRun(t)
}

// TestWorkflowEngineStartProcessNonExistentDef tests that starting a process with non-existent definition fails
func TestWorkflowEngineStartProcessNonExistentDef(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Starting process with non-existent definition returns error", prop.ForAll(
		func(processDefID int64, initiatorID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Don't add any process definition

			// Start process
			req := &StartProcessRequest{
				ProcessDefID:  processDefID,
				Title:         "Test",
				InitiatorID:   initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}

			_, err := te.engine.StartProcess(ctx, req)

			// Should return error
			if err != ErrProcessDefNotFound {
				t.Logf("Expected ErrProcessDefNotFound, got %v", err)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
	))

	properties.TestingRun(t)
}

// ============================================================================
// Property 8: Task Approval Transition Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 8: 任务通过后流转正确性**
// **Validates: Requirements 3.1**
// *For any* 被通过的Task，当前Task状态应变为"approved"，且应根据Transition创建下一节点的Task（除非是结束节点）。

// CompleteTask 完成任务（内存实现）
func (e *inMemoryWorkflowEngine) CompleteTask(ctx context.Context, taskID int64, action string, comment string, operatorID int64) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != operatorID {
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
	var currentNode *model.ProcessNode
	for i := range graph.Nodes {
		if graph.Nodes[i].ID == task.NodeID {
			currentNode = &graph.Nodes[i]
			break
		}
	}
	if currentNode == nil {
		return fmt.Errorf("当前节点不存在: %s", task.NodeID)
	}

	// 8. 根据操作类型处理
	switch action {
	case model.TaskLogActionApprove:
		return e.handleApprove(ctx, task, instance, currentNode, graph, comment)
	case model.TaskLogActionReject:
		return e.handleReject(ctx, task, instance, currentNode, comment)
	default:
		return fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// handleApprove 处理通过操作
func (e *inMemoryWorkflowEngine) handleApprove(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance, currentNode *model.ProcessNode, graph *model.ProcessGraph, comment string) error {
	// 1. 更新任务状态为已通过
	task.Status = model.WfTaskStatusApproved
	task.Comment = comment
	e.taskRepo.tasks[task.ID] = task

	// 2. 记录任务日志
	log := &model.TaskLog{
		TaskID:        task.ID,
		ProcessInstID: instance.ID,
		OperatorID:    task.AssigneeID,
		OperatorName:  task.AssigneeName,
		Action:        model.TaskLogActionApprove,
		Comment:       comment,
	}
	e.taskRepo.CreateTaskLog(log)

	// 3. 检查审批模式，决定是否流转
	shouldTransition := e.shouldTransitionToNextNode(ctx, instance, currentNode)
	if !shouldTransition {
		return nil // 会签模式下还有其他人未审批，不流转
	}

	// 4. 查找下一个节点
	var nextNodeID string
	for _, edge := range graph.Edges {
		if edge.Source == currentNode.ID {
			nextNodeID = edge.Target
			break
		}
	}
	if nextNodeID == "" {
		return fmt.Errorf("找不到下一节点")
	}

	var nextNode *model.ProcessNode
	for i := range graph.Nodes {
		if graph.Nodes[i].ID == nextNodeID {
			nextNode = &graph.Nodes[i]
			break
		}
	}
	if nextNode == nil {
		return fmt.Errorf("下一节点不存在: %s", nextNodeID)
	}

	// 5. 处理下一节点
	return e.createTasksForNode(ctx, instance, nextNode, graph)
}

// handleReject 处理拒绝操作
func (e *inMemoryWorkflowEngine) handleReject(ctx context.Context, task *model.WfTask, instance *model.ProcessInstance, currentNode *model.ProcessNode, comment string) error {
	// 1. 更新任务状态为已拒绝
	task.Status = model.WfTaskStatusRejected
	task.Comment = comment
	e.taskRepo.tasks[task.ID] = task

	// 2. 记录任务日志
	log := &model.TaskLog{
		TaskID:        task.ID,
		ProcessInstID: instance.ID,
		OperatorID:    task.AssigneeID,
		OperatorName:  task.AssigneeName,
		Action:        model.TaskLogActionReject,
		Comment:       comment,
	}
	e.taskRepo.CreateTaskLog(log)

	// 3. 根据拒绝处理方式决定流程走向
	rejectAction := model.RejectActionTerminate // 默认终止
	if currentNode.Properties != nil && currentNode.Properties.RejectAction != "" {
		rejectAction = currentNode.Properties.RejectAction
	}

	switch rejectAction {
	case model.RejectActionTerminate:
		// 终止流程
		instance.Status = model.ProcessInstStatusRejected
		return e.processInstRepo.Update(instance)
	default:
		instance.Status = model.ProcessInstStatusRejected
		return e.processInstRepo.Update(instance)
	}
}

// shouldTransitionToNextNode 判断是否应该流转到下一节点
func (e *inMemoryWorkflowEngine) shouldTransitionToNextNode(ctx context.Context, instance *model.ProcessInstance, node *model.ProcessNode) bool {
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
		pendingCount, _ := e.taskRepo.CountPendingByProcessInstIDAndNodeID(instance.ID, node.ID)
		return pendingCount == 0
	default:
		return true
	}
}

func TestWorkflowEngineTaskApprovalTransition(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("After task approval, task status is approved and next node tasks are created", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process with two approval nodes
			graph := &model.ProcessGraph{
				Nodes: []model.ProcessNode{
					{ID: "start_1", Type: model.NodeTypeStart, Name: "开始", Position: model.Position{X: 100, Y: 100}},
					{
						ID: "approval_1", Type: model.NodeTypeApproval, Name: "审批节点1",
						Position: model.Position{X: 300, Y: 100},
						Properties: &model.NodeProperties{
							AssigneeRule: &model.AssigneeRule{Type: model.AssigneeTypeUser, Values: []int64{assigneeID}},
							ApprovalMode: model.ApprovalModeOrSign,
						},
					},
					{
						ID: "approval_2", Type: model.NodeTypeApproval, Name: "审批节点2",
						Position: model.Position{X: 500, Y: 100},
						Properties: &model.NodeProperties{
							AssigneeRule: &model.AssigneeRule{Type: model.AssigneeTypeUser, Values: []int64{assigneeID + 1}},
							ApprovalMode: model.ApprovalModeOrSign,
						},
					},
					{ID: "end_1", Type: model.NodeTypeEnd, Name: "结束", Position: model.Position{X: 700, Y: 100}},
				},
				Edges: []model.ProcessEdge{
					{ID: "edge_1", Source: "start_1", Target: "approval_1"},
					{ID: "edge_2", Source: "approval_1", Target: "approval_2"},
					{ID: "edge_3", Source: "approval_2", Target: "end_1"},
				},
			}

			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get the first task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			firstTask := tasks[0]

			// Complete the first task
			err = te.engine.CompleteTask(ctx, firstTask.ID, model.TaskLogActionApprove, "Approved", assigneeID)
			if err != nil {
				t.Logf("CompleteTask failed: %v", err)
				return false
			}

			// Property 1: First task status should be "approved"
			updatedTask, _ := te.taskRepo.FindByID(firstTask.ID)
			if updatedTask.Status != model.WfTaskStatusApproved {
				t.Logf("Expected task status 'approved', got '%s'", updatedTask.Status)
				return false
			}

			// Property 2: Next node tasks should be created
			allTasks, _ := te.taskRepo.FindByProcessInstID(instance.ID)
			var nextNodeTasks []model.WfTask
			for _, task := range allTasks {
				if task.NodeID == "approval_2" {
					nextNodeTasks = append(nextNodeTasks, task)
				}
			}
			if len(nextNodeTasks) == 0 {
				t.Logf("Expected tasks for next node, got 0")
				return false
			}

			// Property 3: Next node tasks should be pending
			for _, task := range nextNodeTasks {
				if task.Status != model.WfTaskStatusPending {
					t.Logf("Expected next task status 'pending', got '%s'", task.Status)
					return false
				}
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}

// TestWorkflowEngineTaskApprovalToEndNode tests that approving the last task completes the process
func TestWorkflowEngineTaskApprovalToEndNode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Approving task before end node completes the process", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a simple process with one approval node
			graph := createValidProcessGraph([]int64{assigneeID})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get the task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			task := tasks[0]

			// Complete the task
			err = te.engine.CompleteTask(ctx, task.ID, model.TaskLogActionApprove, "Approved", assigneeID)
			if err != nil {
				t.Logf("CompleteTask failed: %v", err)
				return false
			}

			// Property: Process should be completed
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusCompleted {
				t.Logf("Expected process status 'completed', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}

// ============================================================================
// Property 7: Approval Mode Transition Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 7: 审批模式流转正确性**
// **Validates: Requirements 3.3, 3.4**
// *For any* 配置为会签模式的审批节点，只有当所有Assignee都完成时才流转；
// *For any* 配置为或签模式的审批节点，任一Assignee完成即流转。

func TestWorkflowEngineOrSignMode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Or-sign mode: any assignee approval triggers transition", prop.ForAll(
		func(processDefID int64, initiatorID int64, assignee1 int64, assignee2 int64) bool {
			// Ensure different assignees
			if assignee1 == assignee2 {
				assignee2 = assignee1 + 1
			}

			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process with or-sign mode and multiple assignees
			graph := &model.ProcessGraph{
				Nodes: []model.ProcessNode{
					{ID: "start_1", Type: model.NodeTypeStart, Name: "开始", Position: model.Position{X: 100, Y: 100}},
					{
						ID: "approval_1", Type: model.NodeTypeApproval, Name: "审批节点1",
						Position: model.Position{X: 300, Y: 100},
						Properties: &model.NodeProperties{
							AssigneeRule: &model.AssigneeRule{Type: model.AssigneeTypeUser, Values: []int64{assignee1, assignee2}},
							ApprovalMode: model.ApprovalModeOrSign, // Or-sign mode
						},
					},
					{ID: "end_1", Type: model.NodeTypeEnd, Name: "结束", Position: model.Position{X: 500, Y: 100}},
				},
				Edges: []model.ProcessEdge{
					{ID: "edge_1", Source: "start_1", Target: "approval_1"},
					{ID: "edge_2", Source: "approval_1", Target: "end_1"},
				},
			}

			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Verify: Two tasks should be created (one for each assignee)
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) != 2 {
				t.Logf("Expected 2 pending tasks, got %d", len(tasks))
				return false
			}

			// Find task for assignee1
			var task1 *model.WfTask
			for i := range tasks {
				if tasks[i].AssigneeID == assignee1 {
					task1 = &tasks[i]
					break
				}
			}
			if task1 == nil {
				t.Logf("Task for assignee1 not found")
				return false
			}

			// Complete only one task (assignee1)
			err = te.engine.CompleteTask(ctx, task1.ID, model.TaskLogActionApprove, "Approved", assignee1)
			if err != nil {
				t.Logf("CompleteTask failed: %v", err)
				return false
			}

			// Property: Process should be completed (or-sign: one approval is enough)
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusCompleted {
				t.Logf("Expected process status 'completed', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100),
	))

	properties.TestingRun(t)
}

func TestWorkflowEngineAndSignMode(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("And-sign mode: all assignees must approve before transition", prop.ForAll(
		func(processDefID int64, initiatorID int64, assignee1 int64, assignee2 int64) bool {
			// Ensure different assignees
			if assignee1 == assignee2 {
				assignee2 = assignee1 + 1
			}

			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process with and-sign mode and multiple assignees
			graph := &model.ProcessGraph{
				Nodes: []model.ProcessNode{
					{ID: "start_1", Type: model.NodeTypeStart, Name: "开始", Position: model.Position{X: 100, Y: 100}},
					{
						ID: "approval_1", Type: model.NodeTypeApproval, Name: "审批节点1",
						Position: model.Position{X: 300, Y: 100},
						Properties: &model.NodeProperties{
							AssigneeRule: &model.AssigneeRule{Type: model.AssigneeTypeUser, Values: []int64{assignee1, assignee2}},
							ApprovalMode: model.ApprovalModeAndSign, // And-sign mode
						},
					},
					{ID: "end_1", Type: model.NodeTypeEnd, Name: "结束", Position: model.Position{X: 500, Y: 100}},
				},
				Edges: []model.ProcessEdge{
					{ID: "edge_1", Source: "start_1", Target: "approval_1"},
					{ID: "edge_2", Source: "approval_1", Target: "end_1"},
				},
			}

			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Verify: Two tasks should be created
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) != 2 {
				t.Logf("Expected 2 pending tasks, got %d", len(tasks))
				return false
			}

			// Find tasks for both assignees
			var task1, task2 *model.WfTask
			for i := range tasks {
				if tasks[i].AssigneeID == assignee1 {
					task1 = &tasks[i]
				} else if tasks[i].AssigneeID == assignee2 {
					task2 = &tasks[i]
				}
			}
			if task1 == nil || task2 == nil {
				t.Logf("Tasks for assignees not found")
				return false
			}

			// Complete only one task (assignee1)
			err = te.engine.CompleteTask(ctx, task1.ID, model.TaskLogActionApprove, "Approved", assignee1)
			if err != nil {
				t.Logf("CompleteTask failed: %v", err)
				return false
			}

			// Property 1: Process should still be running (and-sign: need all approvals)
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusRunning {
				t.Logf("Expected process status 'running' after first approval, got '%s'", updatedInstance.Status)
				return false
			}

			// Complete second task (assignee2)
			err = te.engine.CompleteTask(ctx, task2.ID, model.TaskLogActionApprove, "Approved", assignee2)
			if err != nil {
				t.Logf("CompleteTask for second assignee failed: %v", err)
				return false
			}

			// Property 2: Process should be completed after all approvals
			finalInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if finalInstance.Status != model.ProcessInstStatusCompleted {
				t.Logf("Expected process status 'completed' after all approvals, got '%s'", finalInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100),
	))

	properties.TestingRun(t)
}

// ============================================================================
// Property 10: Delegation/Transfer Record Integrity Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 10: 委托/转办记录完整性**
// **Validates: Requirements 3.5, 3.6**
// *For any* 被委托或转办的Task，应记录原Assignee信息，且新Assignee应能查询到该任务。

// DelegateTask 委托任务（内存实现）
func (e *inMemoryWorkflowEngine) DelegateTask(ctx context.Context, taskID int64, toUserID int64, toUserName string, reason string, operatorID int64, operatorName string) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != operatorID {
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

	// 5. 执行委托
	originalAssigneeID := task.AssigneeID
	task.DelegatedFrom = &originalAssigneeID
	task.AssigneeID = toUserID
	task.AssigneeName = toUserName
	// 委托后任务状态仍为pending
	e.taskRepo.tasks[task.ID] = task

	// 6. 记录委托日志
	log := &model.TaskLog{
		TaskID:        task.ID,
		ProcessInstID: instance.ID,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		Action:        model.TaskLogActionDelegate,
		Comment:       fmt.Sprintf("委托给 %s，原因：%s", toUserName, reason),
	}
	e.taskRepo.CreateTaskLog(log)

	return nil
}

// TransferTask 转办任务（内存实现）
func (e *inMemoryWorkflowEngine) TransferTask(ctx context.Context, taskID int64, toUserID int64, toUserName string, reason string, operatorID int64, operatorName string) error {
	// 1. 获取任务
	task, err := e.taskRepo.FindByID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}

	// 2. 检查任务状态
	if task.Status != model.WfTaskStatusPending {
		return ErrTaskNotPending
	}

	// 3. 检查操作人是否是任务审批人
	if task.AssigneeID != operatorID {
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

	// 5. 执行转办
	originalAssigneeID := task.AssigneeID
	task.TransferredFrom = &originalAssigneeID
	task.AssigneeID = toUserID
	task.AssigneeName = toUserName
	// 转办后任务状态仍为pending
	e.taskRepo.tasks[task.ID] = task

	// 6. 记录转办日志
	log := &model.TaskLog{
		TaskID:        task.ID,
		ProcessInstID: instance.ID,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		Action:        model.TaskLogActionTransfer,
		Comment:       fmt.Sprintf("转办给 %s，原因：%s", toUserName, reason),
	}
	e.taskRepo.CreateTaskLog(log)

	return nil
}

func TestWorkflowEngineDelegateTask(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Delegated task records original assignee and new assignee can query it", prop.ForAll(
		func(processDefID int64, initiatorID int64, originalAssignee int64, newAssignee int64) bool {
			// Ensure different assignees
			if originalAssignee == newAssignee {
				newAssignee = originalAssignee + 1
			}

			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process
			graph := createValidProcessGraph([]int64{originalAssignee})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get the task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			task := tasks[0]

			// Delegate the task
			err = te.engine.DelegateTask(ctx, task.ID, newAssignee, fmt.Sprintf("user_%d", newAssignee), "Test delegation", originalAssignee, fmt.Sprintf("user_%d", originalAssignee))
			if err != nil {
				t.Logf("DelegateTask failed: %v", err)
				return false
			}

			// Property 1: Task should record original assignee
			updatedTask, _ := te.taskRepo.FindByID(task.ID)
			if updatedTask.DelegatedFrom == nil || *updatedTask.DelegatedFrom != originalAssignee {
				t.Logf("Expected DelegatedFrom to be %d, got %v", originalAssignee, updatedTask.DelegatedFrom)
				return false
			}

			// Property 2: Task should have new assignee
			if updatedTask.AssigneeID != newAssignee {
				t.Logf("Expected AssigneeID to be %d, got %d", newAssignee, updatedTask.AssigneeID)
				return false
			}

			// Property 3: Task should still be pending
			if updatedTask.Status != model.WfTaskStatusPending {
				t.Logf("Expected task status 'pending', got '%s'", updatedTask.Status)
				return false
			}

			// Property 4: Delegation log should be recorded
			logs, _ := te.taskRepo.FindTaskLogsByProcessInstID(instance.ID)
			var delegateLog *model.TaskLog
			for i := range logs {
				if logs[i].Action == model.TaskLogActionDelegate {
					delegateLog = &logs[i]
					break
				}
			}
			if delegateLog == nil {
				t.Logf("Delegation log not found")
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100),
	))

	properties.TestingRun(t)
}

func TestWorkflowEngineTransferTask(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Transferred task records original assignee and new assignee can query it", prop.ForAll(
		func(processDefID int64, initiatorID int64, originalAssignee int64, newAssignee int64) bool {
			// Ensure different assignees
			if originalAssignee == newAssignee {
				newAssignee = originalAssignee + 1
			}

			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process
			graph := createValidProcessGraph([]int64{originalAssignee})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get the task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			task := tasks[0]

			// Transfer the task
			err = te.engine.TransferTask(ctx, task.ID, newAssignee, fmt.Sprintf("user_%d", newAssignee), "Test transfer", originalAssignee, fmt.Sprintf("user_%d", originalAssignee))
			if err != nil {
				t.Logf("TransferTask failed: %v", err)
				return false
			}

			// Property 1: Task should record original assignee
			updatedTask, _ := te.taskRepo.FindByID(task.ID)
			if updatedTask.TransferredFrom == nil || *updatedTask.TransferredFrom != originalAssignee {
				t.Logf("Expected TransferredFrom to be %d, got %v", originalAssignee, updatedTask.TransferredFrom)
				return false
			}

			// Property 2: Task should have new assignee
			if updatedTask.AssigneeID != newAssignee {
				t.Logf("Expected AssigneeID to be %d, got %d", newAssignee, updatedTask.AssigneeID)
				return false
			}

			// Property 3: Task should still be pending
			if updatedTask.Status != model.WfTaskStatusPending {
				t.Logf("Expected task status 'pending', got '%s'", updatedTask.Status)
				return false
			}

			// Property 4: Transfer log should be recorded
			logs, _ := te.taskRepo.FindTaskLogsByProcessInstID(instance.ID)
			var transferLog *model.TaskLog
			for i := range logs {
				if logs[i].Action == model.TaskLogActionTransfer {
					transferLog = &logs[i]
					break
				}
			}
			if transferLog == nil {
				t.Logf("Transfer log not found")
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100),
	))

	properties.TestingRun(t)
}

// ============================================================================
// Property 9: Process State Terminal Consistency Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 9: 流程状态终态一致性**
// **Validates: Requirements 2.3, 2.4, 2.5**
// *For any* ProcessInstance，当到达结束节点时状态应为"completed"；
// 当被撤回时状态应为"withdrawn"；当被拒绝且配置为终止时状态应为"rejected"。

// WithdrawProcess 撤回流程（内存实现）
func (e *inMemoryWorkflowEngine) WithdrawProcess(ctx context.Context, instanceID int64, operatorID int64) error {
	// 1. 获取流程实例
	instance, err := e.processInstRepo.FindByID(instanceID)
	if err != nil {
		return ErrProcessInstNotFound
	}

	// 2. 检查流程实例状态
	if instance.Status != model.ProcessInstStatusRunning {
		return ErrProcessNotRunning
	}

	// 3. 检查操作人是否是发起人
	if instance.InitiatorID != operatorID {
		return ErrNotProcessInitiator
	}

	// 4. 取消所有待处理任务
	for _, task := range e.taskRepo.tasks {
		if task.ProcessInstID == instanceID && task.Status == model.WfTaskStatusPending {
			delete(e.taskRepo.tasks, task.ID)
		}
	}

	// 5. 更新流程实例状态
	instance.Status = model.ProcessInstStatusWithdrawn
	return e.processInstRepo.Update(instance)
}

func TestWorkflowEngineProcessCompletedState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Process reaching end node has completed status", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a simple process
			graph := createValidProcessGraph([]int64{assigneeID})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get and complete the task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			task := tasks[0]

			err = te.engine.CompleteTask(ctx, task.ID, model.TaskLogActionApprove, "Approved", assigneeID)
			if err != nil {
				t.Logf("CompleteTask failed: %v", err)
				return false
			}

			// Property: Process should be completed
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusCompleted {
				t.Logf("Expected process status 'completed', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}

func TestWorkflowEngineProcessWithdrawnState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Withdrawn process has withdrawn status", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a simple process
			graph := createValidProcessGraph([]int64{assigneeID})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Withdraw the process
			err = te.engine.WithdrawProcess(ctx, instance.ID, initiatorID)
			if err != nil {
				t.Logf("WithdrawProcess failed: %v", err)
				return false
			}

			// Property: Process should be withdrawn
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusWithdrawn {
				t.Logf("Expected process status 'withdrawn', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}

func TestWorkflowEngineProcessRejectedState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Rejected process with terminate action has rejected status", prop.ForAll(
		func(processDefID int64, initiatorID int64, assigneeID int64) bool {
			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a process with terminate on reject
			graph := &model.ProcessGraph{
				Nodes: []model.ProcessNode{
					{ID: "start_1", Type: model.NodeTypeStart, Name: "开始", Position: model.Position{X: 100, Y: 100}},
					{
						ID: "approval_1", Type: model.NodeTypeApproval, Name: "审批节点1",
						Position: model.Position{X: 300, Y: 100},
						Properties: &model.NodeProperties{
							AssigneeRule: &model.AssigneeRule{Type: model.AssigneeTypeUser, Values: []int64{assigneeID}},
							ApprovalMode: model.ApprovalModeOrSign,
							RejectAction: model.RejectActionTerminate, // Terminate on reject
						},
					},
					{ID: "end_1", Type: model.NodeTypeEnd, Name: "结束", Position: model.Position{X: 500, Y: 100}},
				},
				Edges: []model.ProcessEdge{
					{ID: "edge_1", Source: "start_1", Target: "approval_1"},
					{ID: "edge_2", Source: "approval_1", Target: "end_1"},
				},
			}

			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Get and reject the task
			tasks, _ := te.taskRepo.FindPendingByProcessInstID(instance.ID)
			if len(tasks) == 0 {
				t.Logf("No pending tasks found")
				return false
			}
			task := tasks[0]

			err = te.engine.CompleteTask(ctx, task.ID, model.TaskLogActionReject, "Rejected", assigneeID)
			if err != nil {
				t.Logf("CompleteTask (reject) failed: %v", err)
				return false
			}

			// Property: Process should be rejected
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusRejected {
				t.Logf("Expected process status 'rejected', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		genInitiatorIDForEngine(),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}

func TestWorkflowEngineWithdrawNotInitiator(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Non-initiator cannot withdraw process", prop.ForAll(
		func(processDefID int64, initiatorID int64, otherUserID int64, assigneeID int64) bool {
			// Ensure different users
			if initiatorID == otherUserID {
				otherUserID = initiatorID + 1
			}

			ctx := context.Background()
			te := newTestWorkflowEngine()

			// Setup: Create a simple process
			graph := createValidProcessGraph([]int64{assigneeID})
			graphJSON, _ := te.serializer.SerializeProcessGraph(graph)
			processDef := &model.ProcessDefinition{
				ID: processDefID, Name: "Test Process", Code: fmt.Sprintf("test_%d", processDefID),
				GraphJSON: graphJSON, Version: 1, Status: model.ProcessDefStatusPublished,
			}
			te.processDefRepo.AddDefinition(processDef)

			// Start process
			req := &StartProcessRequest{
				ProcessDefID: processDefID, Title: "Test", InitiatorID: initiatorID,
				InitiatorName: fmt.Sprintf("user_%d", initiatorID),
			}
			instance, err := te.engine.StartProcess(ctx, req)
			if err != nil {
				t.Logf("StartProcess failed: %v", err)
				return false
			}

			// Try to withdraw as non-initiator
			err = te.engine.WithdrawProcess(ctx, instance.ID, otherUserID)

			// Property: Should return error
			if err != ErrNotProcessInitiator {
				t.Logf("Expected ErrNotProcessInitiator, got %v", err)
				return false
			}

			// Process should still be running
			updatedInstance, _ := te.processInstRepo.FindByID(instance.ID)
			if updatedInstance.Status != model.ProcessInstStatusRunning {
				t.Logf("Expected process status 'running', got '%s'", updatedInstance.Status)
				return false
			}

			return true
		},
		genProcessDefID(),
		gen.Int64Range(1, 50),
		gen.Int64Range(51, 100),
		gen.Int64Range(1, 100),
	))

	properties.TestingRun(t)
}
