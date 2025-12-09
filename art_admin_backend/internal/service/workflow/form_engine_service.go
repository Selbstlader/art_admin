package workflow

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"

	"gorm.io/gorm"
)

// FormEngineService 表单引擎服务接口
type FormEngineService interface {
	// CreateTemplate 创建表单模板
	// Requirements: 4.1 - WHEN 管理员创建表单模板 THEN Form_Engine SHALL 保存Form_Template并生成唯一标识
	CreateTemplate(ctx context.Context, req *CreateFormTemplateRequest) (*model.FormTemplate, error)

	// UpdateTemplate 更新表单模板
	UpdateTemplate(ctx context.Context, id int64, req *UpdateFormTemplateRequest) (*model.FormTemplate, error)

	// GetTemplate 获取表单模板
	GetTemplate(ctx context.Context, id int64) (*model.FormTemplate, error)

	// ListTemplates 分页查询表单模板列表
	// Requirements: 4.9 - WHEN 管理员查询表单模板列表 THEN Form_Engine SHALL 返回所有可用模板供流程配置时选择
	ListTemplates(ctx context.Context, req *ListFormTemplateRequest) ([]*model.FormTemplate, int64, error)

	// DeleteTemplate 删除表单模板
	DeleteTemplate(ctx context.Context, id int64) error

	// GetSchema 获取表单结构
	GetSchema(ctx context.Context, id int64) (*model.FormSchema, error)

	// ValidateFormData 验证表单数据
	// Requirements: 4.3 - WHEN 管理员配置字段验证规则 THEN Form_Engine SHALL 在表单提交时执行必填、格式、范围等校验
	// Requirements: 2.6 - IF Process_Instance创建时表单数据不符合Form_Definition约束 THEN Workflow_Engine SHALL 拒绝创建并返回验证错误信息
	ValidateFormData(ctx context.Context, schema *model.FormSchema, data map[string]interface{}) *FormValidationResult

	// GetFieldPermissions 获取节点字段权限
	// Requirements: 4.5 - WHEN 管理员设置节点字段权限 THEN Form_Engine SHALL 根据当前节点控制字段的可见性和可编辑性
	GetFieldPermissions(ctx context.Context, templateID int64, nodePermissions map[string]string) (*FieldPermissions, error)
}

// CreateFormTemplateRequest 创建表单模板请求
type CreateFormTemplateRequest struct {
	Name        string            `json:"name" binding:"required"`
	Code        string            `json:"code" binding:"required"`
	Description string            `json:"description"`
	Schema      *model.FormSchema `json:"schema" binding:"required"`
	CreatedBy   int64             `json:"createdBy"`
}

// UpdateFormTemplateRequest 更新表单模板请求
type UpdateFormTemplateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Schema      *model.FormSchema `json:"schema"`
}

// ListFormTemplateRequest 查询表单模板列表请求
type ListFormTemplateRequest struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	CreatedBy *int64 `json:"createdBy"`
	Current   int    `json:"current"`
	Size      int    `json:"size"`
}

// FormValidationResult 表单验证结果
type FormValidationResult struct {
	Valid  bool                  `json:"valid"`
	Errors []FormValidationError `json:"errors,omitempty"`
}

// FormValidationError 表单验证错误
type FormValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FieldPermissions 字段权限结果
type FieldPermissions struct {
	Fields []FieldPermissionItem `json:"fields"`
}

// FieldPermissionItem 单个字段权限
type FieldPermissionItem struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	Permission string `json:"permission"` // visible/editable/hidden
}

// formEngineService 表单引擎服务实现
type formEngineService struct {
	repo       *repository.FormTemplateRepository
	serializer Serializer
}

// NewFormEngineService 创建表单引擎服务实例
func NewFormEngineService(
	repo *repository.FormTemplateRepository,
	serializer Serializer,
) FormEngineService {
	return &formEngineService{
		repo:       repo,
		serializer: serializer,
	}
}

// CreateTemplate 创建表单模板
func (s *formEngineService) CreateTemplate(ctx context.Context, req *CreateFormTemplateRequest) (*model.FormTemplate, error) {
	// 1. 检查编码是否已存在
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, fmt.Errorf("检查编码失败: %w", err)
	}
	if exists {
		return nil, errors.New("表单模板编码已存在")
	}

	// 2. 序列化表单结构
	schemaJSON, err := s.serializer.SerializeFormSchema(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("序列化表单结构失败: %w", err)
	}

	// 3. 创建表单模板
	template := &model.FormTemplate{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		SchemaJSON:  schemaJSON,
		Status:      model.FormTemplateStatusActive,
		CreatedBy:   req.CreatedBy,
	}

	if err := s.repo.Create(template); err != nil {
		return nil, fmt.Errorf("创建表单模板失败: %w", err)
	}

	return template, nil
}

// UpdateTemplate 更新表单模板
func (s *formEngineService) UpdateTemplate(ctx context.Context, id int64, req *UpdateFormTemplateRequest) (*model.FormTemplate, error) {
	// 1. 获取现有模板
	template, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单模板不存在")
		}
		return nil, fmt.Errorf("查询表单模板失败: %w", err)
	}

	// 2. 如果提供了新的表单结构，序列化
	if req.Schema != nil {
		schemaJSON, err := s.serializer.SerializeFormSchema(req.Schema)
		if err != nil {
			return nil, fmt.Errorf("序列化表单结构失败: %w", err)
		}
		template.SchemaJSON = schemaJSON
	}

	// 3. 更新其他字段
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}

	// 4. 保存更新
	if err := s.repo.Update(template); err != nil {
		return nil, fmt.Errorf("更新表单模板失败: %w", err)
	}

	return template, nil
}

// GetTemplate 获取表单模板
func (s *formEngineService) GetTemplate(ctx context.Context, id int64) (*model.FormTemplate, error) {
	template, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单模板不存在")
		}
		return nil, fmt.Errorf("查询表单模板失败: %w", err)
	}
	return template, nil
}

// ListTemplates 分页查询表单模板列表
func (s *formEngineService) ListTemplates(ctx context.Context, req *ListFormTemplateRequest) ([]*model.FormTemplate, int64, error) {
	// 构建查询条件
	query := make(map[string]interface{})
	if req.Name != "" {
		query["name"] = req.Name
	}
	if req.Code != "" {
		query["code"] = req.Code
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
	templates, total, err := s.repo.FindWithPagination(query, current, size)
	if err != nil {
		return nil, 0, fmt.Errorf("查询表单模板列表失败: %w", err)
	}

	// 转换为指针切片
	result := make([]*model.FormTemplate, len(templates))
	for i := range templates {
		result[i] = &templates[i]
	}

	return result, total, nil
}

// DeleteTemplate 删除表单模板
func (s *formEngineService) DeleteTemplate(ctx context.Context, id int64) error {
	// 1. 检查模板是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("表单模板不存在")
		}
		return fmt.Errorf("查询表单模板失败: %w", err)
	}

	// 2. 检查是否被流程定义引用
	inUse, err := s.repo.IsInUse(id)
	if err != nil {
		return fmt.Errorf("检查模板引用失败: %w", err)
	}
	if inUse {
		return errors.New("表单模板已被流程定义引用，无法删除")
	}

	// 3. 删除
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除表单模板失败: %w", err)
	}

	return nil
}

// GetSchema 获取表单结构
func (s *formEngineService) GetSchema(ctx context.Context, id int64) (*model.FormSchema, error) {
	template, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表单模板不存在")
		}
		return nil, fmt.Errorf("查询表单模板失败: %w", err)
	}

	schema, err := s.serializer.DeserializeFormSchema(template.SchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("解析表单结构失败: %w", err)
	}

	return schema, nil
}

// ValidateFormData 验证表单数据
// Requirements: 4.3 - WHEN 管理员配置字段验证规则 THEN Form_Engine SHALL 在表单提交时执行必填、格式、范围等校验
// Requirements: 2.6 - IF Process_Instance创建时表单数据不符合Form_Definition约束 THEN Workflow_Engine SHALL 拒绝创建并返回验证错误信息
func (s *formEngineService) ValidateFormData(ctx context.Context, schema *model.FormSchema, data map[string]interface{}) *FormValidationResult {
	result := &FormValidationResult{
		Valid:  true,
		Errors: []FormValidationError{},
	}

	if schema == nil || len(schema.Fields) == 0 {
		return result
	}

	for _, field := range schema.Fields {
		value, exists := data[field.Key]

		// 1. 必填校验
		if field.Required {
			if !exists || isEmptyValue(value) {
				result.Valid = false
				result.Errors = append(result.Errors, FormValidationError{
					Field:   field.Key,
					Message: fmt.Sprintf("字段 %s 为必填项", field.Label),
				})
				continue
			}
		}

		// 如果字段不存在或为空，且非必填，跳过后续验证
		if !exists || isEmptyValue(value) {
			continue
		}

		// 2. 根据字段类型进行验证
		fieldErrors := validateFieldValue(field, value)
		if len(fieldErrors) > 0 {
			result.Valid = false
			result.Errors = append(result.Errors, fieldErrors...)
		}
	}

	return result
}

// GetFieldPermissions 获取节点字段权限
// Requirements: 4.5 - WHEN 管理员设置节点字段权限 THEN Form_Engine SHALL 根据当前节点控制字段的可见性和可编辑性
func (s *formEngineService) GetFieldPermissions(ctx context.Context, templateID int64, nodePermissions map[string]string) (*FieldPermissions, error) {
	// 1. 获取表单结构
	schema, err := s.GetSchema(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// 2. 构建字段权限列表
	permissions := &FieldPermissions{
		Fields: make([]FieldPermissionItem, 0, len(schema.Fields)),
	}

	for _, field := range schema.Fields {
		permission := string(model.FieldPermissionEditable) // 默认可编辑

		// 如果节点配置了该字段的权限，使用节点配置
		if nodePermissions != nil {
			if perm, ok := nodePermissions[field.Key]; ok {
				permission = perm
			}
		}

		permissions.Fields = append(permissions.Fields, FieldPermissionItem{
			Key:        field.Key,
			Label:      field.Label,
			Permission: permission,
		})
	}

	return permissions, nil
}

// ============================================================================
// 辅助函数
// ============================================================================

// isEmptyValue 检查值是否为空
func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// validateFieldValue 验证字段值
func validateFieldValue(field model.FormField, value interface{}) []FormValidationError {
	var errors []FormValidationError

	switch field.Type {
	case model.FieldTypeText, model.FieldTypeTextarea:
		errors = validateTextValue(field, value)
	case model.FieldTypeNumber:
		errors = validateNumberValue(field, value)
	case model.FieldTypeDate:
		errors = validateDateValue(field, value)
	case model.FieldTypeSelect:
		errors = validateSelectValue(field, value)
	case model.FieldTypeFile:
		// 文件类型暂不做额外验证
	}

	return errors
}

// validateTextValue 验证文本值
func validateTextValue(field model.FormField, value interface{}) []FormValidationError {
	var errors []FormValidationError

	strValue, ok := value.(string)
	if !ok {
		errors = append(errors, FormValidationError{
			Field:   field.Key,
			Message: fmt.Sprintf("字段 %s 必须是文本类型", field.Label),
		})
		return errors
	}

	if field.Validation != nil {
		// 最小长度校验
		if field.Validation.MinLength != nil && len(strValue) < *field.Validation.MinLength {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 长度不能小于 %d", field.Label, *field.Validation.MinLength),
			})
		}

		// 最大长度校验
		if field.Validation.MaxLength != nil && len(strValue) > *field.Validation.MaxLength {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 长度不能大于 %d", field.Label, *field.Validation.MaxLength),
			})
		}

		// 正则表达式校验
		if field.Validation.Pattern != nil && *field.Validation.Pattern != "" {
			matched, err := regexp.MatchString(*field.Validation.Pattern, strValue)
			if err != nil || !matched {
				errors = append(errors, FormValidationError{
					Field:   field.Key,
					Message: fmt.Sprintf("字段 %s 格式不正确", field.Label),
				})
			}
		}
	}

	return errors
}

// validateNumberValue 验证数字值
func validateNumberValue(field model.FormField, value interface{}) []FormValidationError {
	var errors []FormValidationError

	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case float32:
		numValue = float64(v)
	case int:
		numValue = float64(v)
	case int64:
		numValue = float64(v)
	case string:
		// 尝试将字符串转换为数字
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 必须是数字类型", field.Label),
			})
			return errors
		}
		numValue = parsed
	default:
		errors = append(errors, FormValidationError{
			Field:   field.Key,
			Message: fmt.Sprintf("字段 %s 必须是数字类型", field.Label),
		})
		return errors
	}

	if field.Validation != nil {
		// 最小值校验
		if field.Validation.Min != nil && numValue < *field.Validation.Min {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 不能小于 %v", field.Label, *field.Validation.Min),
			})
		}

		// 最大值校验
		if field.Validation.Max != nil && numValue > *field.Validation.Max {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 不能大于 %v", field.Label, *field.Validation.Max),
			})
		}
	}

	return errors
}

// validateDateValue 验证日期值
func validateDateValue(field model.FormField, value interface{}) []FormValidationError {
	var errors []FormValidationError

	// 日期通常以字符串形式传递
	_, ok := value.(string)
	if !ok {
		errors = append(errors, FormValidationError{
			Field:   field.Key,
			Message: fmt.Sprintf("字段 %s 必须是日期格式", field.Label),
		})
	}

	// 可以添加更多日期格式验证逻辑

	return errors
}

// validateSelectValue 验证下拉选择值
func validateSelectValue(field model.FormField, value interface{}) []FormValidationError {
	var errors []FormValidationError

	// 检查值是否在选项列表中
	if len(field.Options) > 0 {
		found := false
		for _, opt := range field.Options {
			if opt.Value == value {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, FormValidationError{
				Field:   field.Key,
				Message: fmt.Sprintf("字段 %s 的值不在有效选项中", field.Label),
			})
		}
	}

	return errors
}
