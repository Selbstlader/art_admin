package workflow

import (
	"context"
	"errors"
	"fmt"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"

	"gorm.io/gorm"
)

// ProcessDefinitionService 流程定义服务接口
type ProcessDefinitionService interface {
	// Create 创建流程定义
	// Requirements: 1.1 - WHEN 管理员在流程设计器中添加节点并连接 THEN Workflow_Engine SHALL 验证流程结构完整性并保存Process_Definition
	Create(ctx context.Context, req *CreateProcessDefRequest) (*model.ProcessDefinition, error)

	// Update 更新流程定义（仅草稿状态可更新）
	Update(ctx context.Context, id int64, req *UpdateProcessDefRequest) (*model.ProcessDefinition, error)

	// GetByID 根据ID获取流程定义
	GetByID(ctx context.Context, id int64) (*model.ProcessDefinition, error)

	// List 分页查询流程定义列表
	List(ctx context.Context, req *ListProcessDefRequest) ([]*model.ProcessDefinition, int64, error)

	// Delete 删除流程定义（仅草稿状态可删除）
	Delete(ctx context.Context, id int64) error

	// Publish 发布流程定义
	// Requirements: 1.4 - WHEN 管理员发布流程定义 THEN Workflow_Engine SHALL 将Process_Definition状态更新为已发布并生成版本号
	// Requirements: 1.5 - WHEN 管理员修改已发布的流程 THEN Workflow_Engine SHALL 创建新版本而保留原版本供运行中实例使用
	Publish(ctx context.Context, id int64) (*model.ProcessDefinition, error)

	// GetGraph 获取流程图结构
	GetGraph(ctx context.Context, id int64) (*model.ProcessGraph, error)
}

// CreateProcessDefRequest 创建流程定义请求
type CreateProcessDefRequest struct {
	Name           string              `json:"name" binding:"required"`
	Code           string              `json:"code" binding:"required"`
	Description    string              `json:"description"`
	Category       string              `json:"category"`
	FormTemplateID *int64              `json:"formTemplateId"`
	Graph          *model.ProcessGraph `json:"graph" binding:"required"`
	CreatedBy      int64               `json:"createdBy"`
}

// UpdateProcessDefRequest 更新流程定义请求
type UpdateProcessDefRequest struct {
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Category       string              `json:"category"`
	FormTemplateID *int64              `json:"formTemplateId"`
	Graph          *model.ProcessGraph `json:"graph"`
}

// ListProcessDefRequest 查询流程定义列表请求
type ListProcessDefRequest struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Category  string `json:"category"`
	Status    string `json:"status"`
	CreatedBy *int64 `json:"createdBy"`
	Current   int    `json:"current"`
	Size      int    `json:"size"`
}

// processDefinitionService 流程定义服务实现
type processDefinitionService struct {
	repo       *repository.ProcessDefinitionRepository
	serializer Serializer
	validator  Validator
}

// NewProcessDefinitionService 创建流程定义服务实例
func NewProcessDefinitionService(
	repo *repository.ProcessDefinitionRepository,
	serializer Serializer,
	validator Validator,
) ProcessDefinitionService {
	return &processDefinitionService{
		repo:       repo,
		serializer: serializer,
		validator:  validator,
	}
}

// Create 创建流程定义
func (s *processDefinitionService) Create(ctx context.Context, req *CreateProcessDefRequest) (*model.ProcessDefinition, error) {
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
		return nil, errors.New("流程编码已存在")
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

// Update 更新流程定义
func (s *processDefinitionService) Update(ctx context.Context, id int64, req *UpdateProcessDefRequest) (*model.ProcessDefinition, error) {
	// 1. 获取现有流程定义
	def, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("流程定义不存在")
		}
		return nil, fmt.Errorf("查询流程定义失败: %w", err)
	}

	// 2. 检查状态，只有草稿状态可以更新
	if def.Status != model.ProcessDefStatusDraft {
		return nil, errors.New("只有草稿状态的流程定义可以更新")
	}

	// 3. 如果提供了新的流程图，验证并序列化
	if req.Graph != nil {
		validationResult := s.validator.ValidateProcessGraph(req.Graph)
		if !validationResult.Valid {
			return nil, fmt.Errorf("流程结构验证失败: %v", validationResult.Errors)
		}

		graphJSON, err := s.serializer.SerializeProcessGraph(req.Graph)
		if err != nil {
			return nil, fmt.Errorf("序列化流程图失败: %w", err)
		}
		def.GraphJSON = graphJSON
	}

	// 4. 更新其他字段
	if req.Name != "" {
		def.Name = req.Name
	}
	if req.Description != "" {
		def.Description = req.Description
	}
	if req.Category != "" {
		def.Category = req.Category
	}
	if req.FormTemplateID != nil {
		def.FormTemplateID = req.FormTemplateID
	}

	// 5. 保存更新
	if err := s.repo.Update(def); err != nil {
		return nil, fmt.Errorf("更新流程定义失败: %w", err)
	}

	return def, nil
}

// GetByID 根据ID获取流程定义
func (s *processDefinitionService) GetByID(ctx context.Context, id int64) (*model.ProcessDefinition, error) {
	def, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("流程定义不存在")
		}
		return nil, fmt.Errorf("查询流程定义失败: %w", err)
	}
	return def, nil
}

// List 分页查询流程定义列表
func (s *processDefinitionService) List(ctx context.Context, req *ListProcessDefRequest) ([]*model.ProcessDefinition, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.Name != "" {
		query["name"] = req.Name
	}
	if req.Code != "" {
		query["code"] = req.Code
	}
	if req.Category != "" {
		query["category"] = req.Category
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	if req.CreatedBy != nil {
		query["createdBy"] = *req.CreatedBy
	}

	// 设置默认分页参数
	current := req.Current
	if current <= 0 {
		current = 1
	}
	size := req.Size
	if size <= 0 {
		size = 10
	}

	// 查询
	defs, total, err := s.repo.FindWithPagination(query, current, size)
	if err != nil {
		return nil, 0, fmt.Errorf("查询流程定义列表失败: %w", err)
	}

	// 转换为指针切片
	result := make([]*model.ProcessDefinition, len(defs))
	for i := range defs {
		result[i] = &defs[i]
	}

	return result, total, nil
}

// Delete 删除流程定义
func (s *processDefinitionService) Delete(ctx context.Context, id int64) error {
	// 1. 获取流程定义
	def, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("流程定义不存在")
		}
		return fmt.Errorf("查询流程定义失败: %w", err)
	}

	// 2. 检查状态，只有草稿状态可以删除
	if def.Status != model.ProcessDefStatusDraft {
		return errors.New("只有草稿状态的流程定义可以删除")
	}

	// 3. 删除
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除流程定义失败: %w", err)
	}

	return nil
}

// Publish 发布流程定义
// Requirements: 1.4 - WHEN 管理员发布流程定义 THEN Workflow_Engine SHALL 将Process_Definition状态更新为已发布并生成版本号
// Requirements: 1.5 - WHEN 管理员修改已发布的流程 THEN Workflow_Engine SHALL 创建新版本而保留原版本供运行中实例使用
func (s *processDefinitionService) Publish(ctx context.Context, id int64) (*model.ProcessDefinition, error) {
	// 1. 获取流程定义
	def, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("流程定义不存在")
		}
		return nil, fmt.Errorf("查询流程定义失败: %w", err)
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

	return nil, errors.New("禁用状态的流程定义不能发布")
}

// GetGraph 获取流程图结构
func (s *processDefinitionService) GetGraph(ctx context.Context, id int64) (*model.ProcessGraph, error) {
	def, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("流程定义不存在")
		}
		return nil, fmt.Errorf("查询流程定义失败: %w", err)
	}

	graph, err := s.serializer.DeserializeProcessGraph(def.GraphJSON)
	if err != nil {
		return nil, fmt.Errorf("解析流程图失败: %w", err)
	}

	return graph, nil
}
