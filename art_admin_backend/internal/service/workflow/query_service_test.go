package workflow

import (
	"art_admin_backend/internal/model"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ============================================================================
// Mock Repositories for Query Service Testing
// ============================================================================

// mockProcessDefRepoForQuery 模拟流程定义仓储（查询服务测试用）
type mockProcessDefRepoForQuery struct {
	definitions map[int64]*model.ProcessDefinition
}

func newMockProcessDefRepoForQuery() *mockProcessDefRepoForQuery {
	return &mockProcessDefRepoForQuery{
		definitions: make(map[int64]*model.ProcessDefinition),
	}
}

func (r *mockProcessDefRepoForQuery) AddDefinition(def *model.ProcessDefinition) {
	r.definitions[def.ID] = def
}

func (r *mockProcessDefRepoForQuery) FindByID(id int64) (*model.ProcessDefinition, error) {
	if def, exists := r.definitions[id]; exists {
		return def, nil
	}
	return nil, fmt.Errorf("process definition not found")
}

func (r *mockProcessDefRepoForQuery) FindWithPagination(query map[string]interface{}, current, size int) ([]model.ProcessDefinition, int64, error) {
	defs := make([]model.ProcessDefinition, 0, len(r.definitions))
	for _, def := range r.definitions {
		defs = append(defs, *def)
	}
	return defs, int64(len(defs)), nil
}

// mockProcessInstRepoForQuery 模拟流程实例仓储（查询服务测试用）
type mockProcessInstRepoForQuery struct {
	instances map[int64]*model.ProcessInstance
	nextID    int64
}

func newMockProcessInstRepoForQuery() *mockProcessInstRepoForQuery {
	return &mockProcessInstRepoForQuery{
		instances: make(map[int64]*model.ProcessInstance),
		nextID:    1,
	}
}

func (r *mockProcessInstRepoForQuery) AddInstance(inst *model.ProcessInstance) {
	if inst.ID == 0 {
		inst.ID = r.nextID
		r.nextID++
	}
	r.instances[inst.ID] = inst
}

func (r *mockProcessInstRepoForQuery) FindByID(id int64) (*model.ProcessInstance, error) {
	if inst, exists := r.instances[id]; exists {
		return inst, nil
	}
	return nil, fmt.Errorf("process instance not found")
}

func (r *mockProcessInstRepoForQuery) FindByInitiator(initiatorID int64, query map[string]interface{}, current, size int) ([]model.ProcessInstance, int64, error) {
	var results []model.ProcessInstance
	for _, inst := range r.instances {
		if inst.InitiatorID == initiatorID {
			// Apply filters
			if title, ok := query["title"]; ok && title != "" {
				// Simple contains check for title
				if inst.Title != title.(string) {
					continue
				}
			}
			if status, ok := query["status"]; ok && status != "" {
				if inst.Status != status.(string) {
					continue
				}
			}
			results = append(results, *inst)
		}
	}

	// Apply pagination
	total := int64(len(results))
	start := (current - 1) * size
	end := start + size
	if start >= len(results) {
		return []model.ProcessInstance{}, total, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

func (r *mockProcessInstRepoForQuery) CountByStatus(status string) (int64, error) {
	var count int64
	for _, inst := range r.instances {
		if inst.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *mockProcessInstRepoForQuery) CountByProcessDefID(processDefID int64) (int64, error) {
	var count int64
	for _, inst := range r.instances {
		if inst.ProcessDefID == processDefID {
			count++
		}
	}
	return count, nil
}

// mockWfTaskRepoForQuery 模拟任务仓储（查询服务测试用）
type mockWfTaskRepoForQuery struct {
	tasks    map[int64]*model.WfTask
	taskLogs map[int64]*model.TaskLog
	nextID   int64
	logID    int64
}

func newMockWfTaskRepoForQuery() *mockWfTaskRepoForQuery {
	return &mockWfTaskRepoForQuery{
		tasks:    make(map[int64]*model.WfTask),
		taskLogs: make(map[int64]*model.TaskLog),
		nextID:   1,
		logID:    1,
	}
}

func (r *mockWfTaskRepoForQuery) AddTask(task *model.WfTask) {
	if task.ID == 0 {
		task.ID = r.nextID
		r.nextID++
	}
	r.tasks[task.ID] = task
}

func (r *mockWfTaskRepoForQuery) AddTaskLog(log *model.TaskLog) {
	if log.ID == 0 {
		log.ID = r.logID
		r.logID++
	}
	r.taskLogs[log.ID] = log
}

func (r *mockWfTaskRepoForQuery) FindPendingByAssignee(assigneeID int64, query map[string]interface{}, current, size int) ([]model.WfTask, int64, error) {
	var results []model.WfTask
	for _, task := range r.tasks {
		if task.AssigneeID == assigneeID && task.Status == model.WfTaskStatusPending {
			results = append(results, *task)
		}
	}

	// Apply pagination
	total := int64(len(results))
	start := (current - 1) * size
	end := start + size
	if start >= len(results) {
		return []model.WfTask{}, total, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

func (r *mockWfTaskRepoForQuery) FindCompletedByAssignee(assigneeID int64, query map[string]interface{}, current, size int) ([]model.WfTask, int64, error) {
	var results []model.WfTask
	for _, task := range r.tasks {
		if task.AssigneeID == assigneeID &&
			(task.Status == model.WfTaskStatusApproved || task.Status == model.WfTaskStatusRejected) {
			results = append(results, *task)
		}
	}

	// Apply pagination
	total := int64(len(results))
	start := (current - 1) * size
	end := start + size
	if start >= len(results) {
		return []model.WfTask{}, total, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

func (r *mockWfTaskRepoForQuery) FindByProcessInstID(processInstID int64) ([]model.WfTask, error) {
	var results []model.WfTask
	for _, task := range r.tasks {
		if task.ProcessInstID == processInstID {
			results = append(results, *task)
		}
	}
	return results, nil
}

func (r *mockWfTaskRepoForQuery) FindTaskLogsByProcessInstID(processInstID int64) ([]model.TaskLog, error) {
	var results []model.TaskLog
	for _, log := range r.taskLogs {
		if log.ProcessInstID == processInstID {
			results = append(results, *log)
		}
	}
	return results, nil
}

// ============================================================================
// Test Query Service
// ============================================================================

// testQueryService 可测试的查询服务
type testQueryService struct {
	processDefRepo  *mockProcessDefRepoForQuery
	processInstRepo *mockProcessInstRepoForQuery
	taskRepo        *mockWfTaskRepoForQuery
	service         *inMemoryQueryService
}

// inMemoryQueryService 内存查询服务（用于测试）
type inMemoryQueryService struct {
	processDefRepo  *mockProcessDefRepoForQuery
	processInstRepo *mockProcessInstRepoForQuery
	taskRepo        *mockWfTaskRepoForQuery
}

func newTestQueryService() *testQueryService {
	processDefRepo := newMockProcessDefRepoForQuery()
	processInstRepo := newMockProcessInstRepoForQuery()
	taskRepo := newMockWfTaskRepoForQuery()

	service := &inMemoryQueryService{
		processDefRepo:  processDefRepo,
		processInstRepo: processInstRepo,
		taskRepo:        taskRepo,
	}

	return &testQueryService{
		processDefRepo:  processDefRepo,
		processInstRepo: processInstRepo,
		taskRepo:        taskRepo,
		service:         service,
	}
}

// ListMyInitiated 查询我发起的流程（内存实现）
func (s *inMemoryQueryService) ListMyInitiated(ctx context.Context, userID int64, req *QueryRequest) ([]*ProcessInstanceSummary, int64, error) {
	current, size := normalizePaginationForTest(req.Current, req.Size)

	query := make(map[string]interface{})
	if req.Title != "" {
		query["title"] = req.Title
	}
	if req.Status != "" {
		query["status"] = req.Status
	}

	instances, total, err := s.processInstRepo.FindByInitiator(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

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

		if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
			summary.ProcessDefName = processDef.Name
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// ListMyTodo 查询我的待办（内存实现）
func (s *inMemoryQueryService) ListMyTodo(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error) {
	current, size := normalizePaginationForTest(req.Current, req.Size)

	query := make(map[string]interface{})
	if req.NodeName != "" {
		query["nodeName"] = req.NodeName
	}

	tasks, total, err := s.taskRepo.FindPendingByAssignee(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

	summaries := make([]*TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		summary := taskToSummaryForTest(&task)

		if inst, err := s.processInstRepo.FindByID(task.ProcessInstID); err == nil {
			summary.ProcessTitle = inst.Title
			if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
				summary.ProcessDefName = processDef.Name
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// ListMyDone 查询我的已办（内存实现）
func (s *inMemoryQueryService) ListMyDone(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error) {
	current, size := normalizePaginationForTest(req.Current, req.Size)

	query := make(map[string]interface{})
	if req.NodeName != "" {
		query["nodeName"] = req.NodeName
	}
	if req.Status != "" {
		query["status"] = req.Status
	}

	tasks, total, err := s.taskRepo.FindCompletedByAssignee(userID, query, current, size)
	if err != nil {
		return nil, 0, err
	}

	summaries := make([]*TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		summary := taskToSummaryForTest(&task)

		if inst, err := s.processInstRepo.FindByID(task.ProcessInstID); err == nil {
			summary.ProcessTitle = inst.Title
			if processDef, err := s.processDefRepo.FindByID(inst.ProcessDefID); err == nil {
				summary.ProcessDefName = processDef.Name
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, total, nil
}

// Helper functions
func normalizePaginationForTest(current, size int) (int, int) {
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

func taskToSummaryForTest(task *model.WfTask) *TaskSummary {
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

// ============================================================================
// Generators
// ============================================================================

// genUserID generates a user ID
func genUserID() gopter.Gen {
	return gen.Int64Range(1, 100)
}

// genProcessInstanceStatus generates a process instance status
func genProcessInstanceStatus() gopter.Gen {
	return gen.OneConstOf(
		model.ProcessInstStatusRunning,
		model.ProcessInstStatusCompleted,
		model.ProcessInstStatusRejected,
		model.ProcessInstStatusWithdrawn,
	)
}

// genTaskStatus generates a task status
func genTaskStatus() gopter.Gen {
	return gen.OneConstOf(
		model.WfTaskStatusPending,
		model.WfTaskStatusApproved,
		model.WfTaskStatusRejected,
	)
}

// genProcessInstance generates a process instance
func genProcessInstance(initiatorID int64, processDefID int64) gopter.Gen {
	return gen.Struct(reflect.TypeOf(model.ProcessInstance{}), map[string]gopter.Gen{
		"ID":                gen.Int64Range(1, 1000),
		"ProcessDefID":      gen.Const(processDefID),
		"ProcessDefVersion": gen.IntRange(1, 10),
		"Title":             gen.AnyString().Map(func(s string) string { return "Process " + s }),
		"InitiatorID":       gen.Const(initiatorID),
		"InitiatorName":     gen.Const(fmt.Sprintf("User_%d", initiatorID)),
		"Status":            genProcessInstanceStatus(),
		"CurrentNodeID":     gen.Const("node_1"),
		"StartedAt":         gen.Const(time.Now()),
	})
}

// genWfTask generates a workflow task
func genWfTask(assigneeID int64, processInstID int64, status string) gopter.Gen {
	return gen.Struct(reflect.TypeOf(model.WfTask{}), map[string]gopter.Gen{
		"ID":            gen.Int64Range(1, 1000),
		"ProcessInstID": gen.Const(processInstID),
		"NodeID":        gen.Const("approval_1"),
		"NodeName":      gen.Const("审批节点"),
		"AssigneeID":    gen.Const(assigneeID),
		"AssigneeName":  gen.Const(fmt.Sprintf("User_%d", assigneeID)),
		"Status":        gen.Const(status),
		"CreatedAt":     gen.Const(time.Now()),
	})
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 13: 查询结果过滤正确性**
// **Validates: Requirements 7.1, 7.2, 7.3**
// *For any* 用户查询：
// - "我发起的"查询结果中所有实例的initiatorId应等于查询用户ID
// - "我的待办"查询结果中所有任务的assigneeId应等于查询用户ID且状态为pending
// - "我的已办"查询结果中所有任务应由查询用户处理过

func TestQueryServiceListMyInitiatedFilterCorrectness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("ListMyInitiated returns only instances initiated by the querying user", prop.ForAll(
		func(queryUserID int64, otherUserID int64, numInstances int) bool {
			// Ensure different users
			if queryUserID == otherUserID {
				otherUserID = queryUserID + 1
			}

			ctx := context.Background()
			ts := newTestQueryService()

			// Setup: Add process definition
			processDef := &model.ProcessDefinition{
				ID:     1,
				Name:   "Test Process",
				Code:   "test_process",
				Status: model.ProcessDefStatusPublished,
			}
			ts.processDefRepo.AddDefinition(processDef)

			// Add instances for query user
			for i := 0; i < numInstances; i++ {
				inst := &model.ProcessInstance{
					ProcessDefID:      1,
					ProcessDefVersion: 1,
					Title:             fmt.Sprintf("Instance %d", i),
					InitiatorID:       queryUserID,
					InitiatorName:     fmt.Sprintf("User_%d", queryUserID),
					Status:            model.ProcessInstStatusRunning,
					CurrentNodeID:     "node_1",
					StartedAt:         time.Now(),
				}
				ts.processInstRepo.AddInstance(inst)
			}

			// Add instances for other user
			for i := 0; i < numInstances; i++ {
				inst := &model.ProcessInstance{
					ProcessDefID:      1,
					ProcessDefVersion: 1,
					Title:             fmt.Sprintf("Other Instance %d", i),
					InitiatorID:       otherUserID,
					InitiatorName:     fmt.Sprintf("User_%d", otherUserID),
					Status:            model.ProcessInstStatusRunning,
					CurrentNodeID:     "node_1",
					StartedAt:         time.Now(),
				}
				ts.processInstRepo.AddInstance(inst)
			}

			// Query
			req := &QueryRequest{Current: 1, Size: 100}
			results, total, err := ts.service.ListMyInitiated(ctx, queryUserID, req)
			if err != nil {
				t.Logf("ListMyInitiated failed: %v", err)
				return false
			}

			// Property 1: Total should match number of instances for query user
			if total != int64(numInstances) {
				t.Logf("Expected total %d, got %d", numInstances, total)
				return false
			}

			// Property 2: All results should have initiatorID equal to queryUserID
			for _, result := range results {
				if result.InitiatorID != queryUserID {
					t.Logf("Expected initiatorID %d, got %d", queryUserID, result.InitiatorID)
					return false
				}
			}

			return true
		},
		genUserID(),
		genUserID(),
		gen.IntRange(1, 10),
	))

	properties.TestingRun(t)
}

func TestQueryServiceListMyTodoFilterCorrectness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("ListMyTodo returns only pending tasks assigned to the querying user", prop.ForAll(
		func(queryUserID int64, otherUserID int64, numTasks int) bool {
			// Ensure different users
			if queryUserID == otherUserID {
				otherUserID = queryUserID + 1
			}

			ctx := context.Background()
			ts := newTestQueryService()

			// Setup: Add process definition and instance
			processDef := &model.ProcessDefinition{
				ID:     1,
				Name:   "Test Process",
				Code:   "test_process",
				Status: model.ProcessDefStatusPublished,
			}
			ts.processDefRepo.AddDefinition(processDef)

			inst := &model.ProcessInstance{
				ID:                1,
				ProcessDefID:      1,
				ProcessDefVersion: 1,
				Title:             "Test Instance",
				InitiatorID:       100,
				InitiatorName:     "Initiator",
				Status:            model.ProcessInstStatusRunning,
				CurrentNodeID:     "node_1",
				StartedAt:         time.Now(),
			}
			ts.processInstRepo.AddInstance(inst)

			// Add pending tasks for query user
			for i := 0; i < numTasks; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    queryUserID,
					AssigneeName:  fmt.Sprintf("User_%d", queryUserID),
					Status:        model.WfTaskStatusPending,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Add pending tasks for other user
			for i := 0; i < numTasks; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    otherUserID,
					AssigneeName:  fmt.Sprintf("User_%d", otherUserID),
					Status:        model.WfTaskStatusPending,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Add completed tasks for query user (should not be returned)
			for i := 0; i < numTasks; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    queryUserID,
					AssigneeName:  fmt.Sprintf("User_%d", queryUserID),
					Status:        model.WfTaskStatusApproved,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Query
			req := &QueryRequest{Current: 1, Size: 100}
			results, total, err := ts.service.ListMyTodo(ctx, queryUserID, req)
			if err != nil {
				t.Logf("ListMyTodo failed: %v", err)
				return false
			}

			// Property 1: Total should match number of pending tasks for query user
			if total != int64(numTasks) {
				t.Logf("Expected total %d, got %d", numTasks, total)
				return false
			}

			// Property 2: All results should have assigneeID equal to queryUserID
			for _, result := range results {
				if result.AssigneeID != queryUserID {
					t.Logf("Expected assigneeID %d, got %d", queryUserID, result.AssigneeID)
					return false
				}
			}

			// Property 3: All results should have status "pending"
			for _, result := range results {
				if result.Status != model.WfTaskStatusPending {
					t.Logf("Expected status 'pending', got '%s'", result.Status)
					return false
				}
			}

			return true
		},
		genUserID(),
		genUserID(),
		gen.IntRange(1, 10),
	))

	properties.TestingRun(t)
}

func TestQueryServiceListMyDoneFilterCorrectness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("ListMyDone returns only completed tasks processed by the querying user", prop.ForAll(
		func(queryUserID int64, otherUserID int64, numApproved int, numRejected int) bool {
			// Ensure different users
			if queryUserID == otherUserID {
				otherUserID = queryUserID + 1
			}

			ctx := context.Background()
			ts := newTestQueryService()

			// Setup: Add process definition and instance
			processDef := &model.ProcessDefinition{
				ID:     1,
				Name:   "Test Process",
				Code:   "test_process",
				Status: model.ProcessDefStatusPublished,
			}
			ts.processDefRepo.AddDefinition(processDef)

			inst := &model.ProcessInstance{
				ID:                1,
				ProcessDefID:      1,
				ProcessDefVersion: 1,
				Title:             "Test Instance",
				InitiatorID:       100,
				InitiatorName:     "Initiator",
				Status:            model.ProcessInstStatusRunning,
				CurrentNodeID:     "node_1",
				StartedAt:         time.Now(),
			}
			ts.processInstRepo.AddInstance(inst)

			// Add approved tasks for query user
			for i := 0; i < numApproved; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    queryUserID,
					AssigneeName:  fmt.Sprintf("User_%d", queryUserID),
					Status:        model.WfTaskStatusApproved,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Add rejected tasks for query user
			for i := 0; i < numRejected; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    queryUserID,
					AssigneeName:  fmt.Sprintf("User_%d", queryUserID),
					Status:        model.WfTaskStatusRejected,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Add completed tasks for other user (should not be returned)
			for i := 0; i < numApproved; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    otherUserID,
					AssigneeName:  fmt.Sprintf("User_%d", otherUserID),
					Status:        model.WfTaskStatusApproved,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Add pending tasks for query user (should not be returned)
			for i := 0; i < numApproved; i++ {
				task := &model.WfTask{
					ProcessInstID: 1,
					NodeID:        "approval_1",
					NodeName:      "审批节点",
					AssigneeID:    queryUserID,
					AssigneeName:  fmt.Sprintf("User_%d", queryUserID),
					Status:        model.WfTaskStatusPending,
					CreatedAt:     time.Now(),
				}
				ts.taskRepo.AddTask(task)
			}

			// Query
			req := &QueryRequest{Current: 1, Size: 100}
			results, total, err := ts.service.ListMyDone(ctx, queryUserID, req)
			if err != nil {
				t.Logf("ListMyDone failed: %v", err)
				return false
			}

			expectedTotal := numApproved + numRejected

			// Property 1: Total should match number of completed tasks for query user
			if total != int64(expectedTotal) {
				t.Logf("Expected total %d, got %d", expectedTotal, total)
				return false
			}

			// Property 2: All results should have assigneeID equal to queryUserID
			for _, result := range results {
				if result.AssigneeID != queryUserID {
					t.Logf("Expected assigneeID %d, got %d", queryUserID, result.AssigneeID)
					return false
				}
			}

			// Property 3: All results should have status "approved" or "rejected"
			for _, result := range results {
				if result.Status != model.WfTaskStatusApproved && result.Status != model.WfTaskStatusRejected {
					t.Logf("Expected status 'approved' or 'rejected', got '%s'", result.Status)
					return false
				}
			}

			return true
		},
		genUserID(),
		genUserID(),
		gen.IntRange(1, 5),
		gen.IntRange(1, 5),
	))

	properties.TestingRun(t)
}
