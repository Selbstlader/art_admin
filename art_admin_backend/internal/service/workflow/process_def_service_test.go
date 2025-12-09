package workflow

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"art_admin_backend/internal/model"
)

// ============================================================================
// Mock Repository for Testing
// ============================================================================

// mockProcessDefinitionRepository 模拟流程定义仓储
type mockProcessDefinitionRepository struct {
	mu          sync.RWMutex
	definitions map[int64]*model.ProcessDefinition
	nextID      int64
}

// newMockProcessDefinitionRepository 创建模拟仓储
func newMockProcessDefinitionRepository() *mockProcessDefinitionRepository {
	return &mockProcessDefinitionRepository{
		definitions: make(map[int64]*model.ProcessDefinition),
		nextID:      1,
	}
}

func (r *mockProcessDefinitionRepository) Create(def *model.ProcessDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	def.ID = r.nextID
	r.nextID++
	// Deep copy to avoid reference issues
	copied := *def
	r.definitions[def.ID] = &copied
	return nil
}

func (r *mockProcessDefinitionRepository) Update(def *model.ProcessDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.definitions[def.ID]; !exists {
		return fmt.Errorf("not found")
	}
	copied := *def
	r.definitions[def.ID] = &copied
	return nil
}

func (r *mockProcessDefinitionRepository) FindByID(id int64) (*model.ProcessDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if def, exists := r.definitions[id]; exists {
		copied := *def
		return &copied, nil
	}
	return nil, fmt.Errorf("not found")
}

func (r *mockProcessDefinitionRepository) FindByCode(code string) (*model.ProcessDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, def := range r.definitions {
		if def.Code == code {
			copied := *def
			return &copied, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *mockProcessDefinitionRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.definitions, id)
	return nil
}

func (r *mockProcessDefinitionRepository) FindWithPagination(query map[string]interface{}, current, size int) ([]model.ProcessDefinition, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.ProcessDefinition
	for _, def := range r.definitions {
		result = append(result, *def)
	}
	return result, int64(len(result)), nil
}

func (r *mockProcessDefinitionRepository) FindPublishedByCode(code string) (*model.ProcessDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, def := range r.definitions {
		if def.Code == code && def.Status == model.ProcessDefStatusPublished {
			copied := *def
			return &copied, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *mockProcessDefinitionRepository) FindLatestVersionByCode(code string) (*model.ProcessDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var latest *model.ProcessDefinition
	for _, def := range r.definitions {
		if def.Code == code {
			if latest == nil || def.Version > latest.Version {
				copied := *def
				latest = &copied
			}
		}
	}
	if latest == nil {
		return nil, fmt.Errorf("not found")
	}
	return latest, nil
}

func (r *mockProcessDefinitionRepository) GetMaxVersionByCode(code string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	maxVersion := 0
	for _, def := range r.definitions {
		if def.Code == code && def.Version > maxVersion {
			maxVersion = def.Version
		}
	}
	return maxVersion, nil
}

func (r *mockProcessDefinitionRepository) UpdateStatus(id int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if def, exists := r.definitions[id]; exists {
		def.Status = status
		return nil
	}
	return fmt.Errorf("not found")
}

func (r *mockProcessDefinitionRepository) ExistsByCode(code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, def := range r.definitions {
		if def.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockProcessDefinitionRepository) ExistsByCodeExcludeID(code string, excludeID int64) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, def := range r.definitions {
		if def.Code == code && def.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockProcessDefinitionRepository) FindByFormTemplateID(formTemplateID int64) ([]model.ProcessDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []model.ProcessDefinition
	for _, def := range r.definitions {
		if def.FormTemplateID != nil && *def.FormTemplateID == formTemplateID {
			result = append(result, *def)
		}
	}
	return result, nil
}

// ============================================================================
// Test Service Factory
// ============================================================================

// testableProcessDefinitionService 可测试的流程定义服务
type testableProcessDefinitionService struct {
	repo       *mockProcessDefinitionRepository
	serializer Serializer
	validator  Validator
}

func newTestableProcessDefinitionService() *testableProcessDefinitionService {
	return &testableProcessDefinitionService{
		repo:       newMockProcessDefinitionRepository(),
		serializer: NewSerializer(),
		validator:  NewValidator(),
	}
}

// Create 创建流程定义
func (s *testableProcessDefinitionService) Create(ctx context.Context, req *CreateProcessDefRequest) (*model.ProcessDefinition, error) {
	// 1. 验证流程图结构
	validationResult := s.validator.ValidateProcessGraph(req.Graph)
	if !validationResult.Valid {
		return nil, fmt.Errorf("流程结构验证失败: %v", validationResult.Errors)
	}

	// 2. 检查编码是否已存在
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, fmt.Errorf("检查编码失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("流程编码已存在")
	}

	// 3. 序列化流程图
	graphJSON, err := s.serializer.SerializeProcessGraph(req.Graph)
	if err != nil {
		return nil, fmt.Errorf("序列化流程图失败: %w", err)
	}

	// 4. 创建流程定义
	def := &model.ProcessDefinition{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		Category:       req.Category,
		FormTemplateID: req.FormTemplateID,
		GraphJSON:      graphJSON,
		Version:        1,
		Status:         model.ProcessDefStatusDraft,
		CreatedBy:      req.CreatedBy,
	}

	if err := s.repo.Create(def); err != nil {
		return nil, fmt.Errorf("创建流程定义失败: %w", err)
	}

	return def, nil
}

// Publish 发布流程定义
func (s *testableProcessDefinitionService) Publish(ctx context.Context, id int64) (*model.ProcessDefinition, error) {
	// 1. 获取流程定义
	def, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("流程定义不存在: %w", err)
	}

	// 2. 验证流程图结构
	graph, err := s.serializer.DeserializeProcessGraph(def.GraphJSON)
	if err != nil {
		return nil, fmt.Errorf("解析流程图失败: %w", err)
	}

	validationResult := s.validator.ValidateProcessGraph(graph)
	if !validationResult.Valid {
		return nil, fmt.Errorf("流程结构验证失败: %v", validationResult.Errors)
	}

	// 3. 如果是草稿状态，直接发布
	if def.Status == model.ProcessDefStatusDraft {
		def.Status = model.ProcessDefStatusPublished
		if err := s.repo.Update(def); err != nil {
			return nil, fmt.Errorf("发布流程定义失败: %w", err)
		}
		return def, nil
	}

	// 4. 如果是已发布状态，创建新版本
	if def.Status == model.ProcessDefStatusPublished {
		// 获取当前最大版本号
		maxVersion, err := s.repo.GetMaxVersionByCode(def.Code)
		if err != nil {
			return nil, fmt.Errorf("获取最大版本号失败: %w", err)
		}

		// 创建新版本
		newDef := &model.ProcessDefinition{
			Name:           def.Name,
			Code:           def.Code,
			Description:    def.Description,
			Category:       def.Category,
			FormTemplateID: def.FormTemplateID,
			GraphJSON:      def.GraphJSON,
			Version:        maxVersion + 1,
			Status:         model.ProcessDefStatusPublished,
			CreatedBy:      def.CreatedBy,
		}

		if err := s.repo.Create(newDef); err != nil {
			return nil, fmt.Errorf("创建新版本失败: %w", err)
		}

		return newDef, nil
	}

	return nil, fmt.Errorf("禁用状态的流程定义不能发布")
}

// ============================================================================
// Generators
// ============================================================================

// genValidProcessGraphForService generates a valid ProcessGraph for service testing
func genValidProcessGraphForService() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(0, 3), // number of middle nodes
	).Map(func(vals []interface{}) *model.ProcessGraph {
		middleNodeCount := vals[0].(int)

		nodes := make([]model.ProcessNode, 0)
		edges := make([]model.ProcessEdge, 0)

		// Add start node
		nodes = append(nodes, model.ProcessNode{
			ID:       "start",
			Type:     model.NodeTypeStart,
			Name:     "开始",
			Position: model.Position{X: 0, Y: 0},
		})

		// Add middle nodes
		for i := 0; i < middleNodeCount; i++ {
			nodes = append(nodes, model.ProcessNode{
				ID:       fmt.Sprintf("node_%d", i),
				Type:     model.NodeTypeApproval,
				Name:     fmt.Sprintf("审批%d", i),
				Position: model.Position{X: float64(i+1) * 100, Y: 0},
			})
		}

		// Add end node
		nodes = append(nodes, model.ProcessNode{
			ID:       "end",
			Type:     model.NodeTypeEnd,
			Name:     "结束",
			Position: model.Position{X: float64(middleNodeCount+1) * 100, Y: 0},
		})

		// Create edges
		if middleNodeCount == 0 {
			edges = append(edges, model.ProcessEdge{
				ID:     "edge_start_end",
				Source: "start",
				Target: "end",
			})
		} else {
			edges = append(edges, model.ProcessEdge{
				ID:     "edge_start_0",
				Source: "start",
				Target: "node_0",
			})
			for i := 0; i < middleNodeCount-1; i++ {
				edges = append(edges, model.ProcessEdge{
					ID:     fmt.Sprintf("edge_%d_%d", i, i+1),
					Source: fmt.Sprintf("node_%d", i),
					Target: fmt.Sprintf("node_%d", i+1),
				})
			}
			edges = append(edges, model.ProcessEdge{
				ID:     fmt.Sprintf("edge_%d_end", middleNodeCount-1),
				Source: fmt.Sprintf("node_%d", middleNodeCount-1),
				Target: "end",
			})
		}

		return &model.ProcessGraph{
			Nodes: nodes,
			Edges: edges,
		}
	})
}

// ============================================================================
// Property-Based Tests
// ============================================================================

// **Feature: oa-workflow-engine, Property 4: 流程发布版本递增**
// **Validates: Requirements 1.4, 1.5**
// *For any* 已发布的ProcessDefinition，再次发布修改后的版本时，新版本号应大于原版本号，且原版本保持不变。
func TestPublishVersionIncrement(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Publishing a published definition creates a new version with incremented version number", prop.ForAll(
		func(graph *model.ProcessGraph, publishCount int) bool {
			ctx := context.Background()
			svc := newTestableProcessDefinitionService()

			// Create a process definition
			code := fmt.Sprintf("test_code_%d", publishCount)
			req := &CreateProcessDefRequest{
				Name:      "Test Process",
				Code:      code,
				Graph:     graph,
				CreatedBy: 1,
			}

			def, err := svc.Create(ctx, req)
			if err != nil {
				t.Logf("Create failed: %v", err)
				return false
			}

			// Initial version should be 1
			if def.Version != 1 {
				t.Logf("Initial version should be 1, got %d", def.Version)
				return false
			}

			// First publish (draft -> published)
			published, err := svc.Publish(ctx, def.ID)
			if err != nil {
				t.Logf("First publish failed: %v", err)
				return false
			}

			if published.Status != model.ProcessDefStatusPublished {
				t.Logf("Status should be published after first publish")
				return false
			}

			// Version should still be 1 after first publish
			if published.Version != 1 {
				t.Logf("Version should be 1 after first publish, got %d", published.Version)
				return false
			}

			// Publish again (should create new version)
			for i := 0; i < publishCount; i++ {
				previousVersion := published.Version
				newPublished, err := svc.Publish(ctx, published.ID)
				if err != nil {
					t.Logf("Subsequent publish failed: %v", err)
					return false
				}

				// New version should be greater than previous
				if newPublished.Version <= previousVersion {
					t.Logf("New version %d should be greater than previous %d", newPublished.Version, previousVersion)
					return false
				}

				// Original definition should still exist with original version
				original, err := svc.repo.FindByID(published.ID)
				if err != nil {
					t.Logf("Original definition should still exist: %v", err)
					return false
				}

				if original.Version != previousVersion {
					t.Logf("Original version should be preserved: expected %d, got %d", previousVersion, original.Version)
					return false
				}

				published = newPublished
			}

			return true
		},
		genValidProcessGraphForService(),
		gen.IntRange(1, 5), // Number of additional publishes
	))

	properties.TestingRun(t)
}

// Test that draft to published doesn't change version
func TestPublishDraftToPublished(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	parameters.Rng.Seed(42)

	properties := gopter.NewProperties(parameters)

	properties.Property("Publishing a draft definition changes status but keeps version 1", prop.ForAll(
		func(graph *model.ProcessGraph, codeNum int) bool {
			ctx := context.Background()
			svc := newTestableProcessDefinitionService()

			// Create a process definition
			code := fmt.Sprintf("draft_test_%d", codeNum)
			req := &CreateProcessDefRequest{
				Name:      "Test Process",
				Code:      code,
				Graph:     graph,
				CreatedBy: 1,
			}

			def, err := svc.Create(ctx, req)
			if err != nil {
				t.Logf("Create failed: %v", err)
				return false
			}

			// Should be draft with version 1
			if def.Status != model.ProcessDefStatusDraft {
				t.Logf("Initial status should be draft")
				return false
			}
			if def.Version != 1 {
				t.Logf("Initial version should be 1")
				return false
			}

			// Publish
			published, err := svc.Publish(ctx, def.ID)
			if err != nil {
				t.Logf("Publish failed: %v", err)
				return false
			}

			// Should be published with version 1
			if published.Status != model.ProcessDefStatusPublished {
				t.Logf("Status should be published")
				return false
			}
			if published.Version != 1 {
				t.Logf("Version should still be 1 after first publish")
				return false
			}

			return true
		},
		genValidProcessGraphForService(),
		gen.IntRange(1, 1000), // Just for unique codes
	))

	properties.TestingRun(t)
}
