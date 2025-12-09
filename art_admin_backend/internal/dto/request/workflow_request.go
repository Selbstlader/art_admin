package request

import "time"

// ============================================================================
// 流程定义相关请求
// ============================================================================

// CreateProcessDefRequest 创建流程定义请求
type CreateProcessDefRequest struct {
	Name           string           `json:"name" binding:"required,max=100"`
	Code           string           `json:"code" binding:"required,max=50"`
	Description    string           `json:"description" binding:"max=500"`
	Category       string           `json:"category" binding:"max=50"`
	FormTemplateID *int64           `json:"formTemplateId"`
	Graph          *ProcessGraphDTO `json:"graph" binding:"required"`
}

// UpdateProcessDefRequest 更新流程定义请求
type UpdateProcessDefRequest struct {
	Name           string           `json:"name" binding:"required,max=100"`
	Description    string           `json:"description" binding:"max=500"`
	Category       string           `json:"category" binding:"max=50"`
	FormTemplateID *int64           `json:"formTemplateId"`
	Graph          *ProcessGraphDTO `json:"graph" binding:"required"`
}

// ListProcessDefRequest 获取流程定义列表请求
type ListProcessDefRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Category string `form:"category"`
	Status   string `form:"status"`
}

// ProcessGraphDTO 流程图DTO
type ProcessGraphDTO struct {
	Nodes       []ProcessNodeDTO     `json:"nodes"`
	Edges       []ProcessEdgeDTO     `json:"edges"`
	GlobalProps *GlobalPropertiesDTO `json:"globalProps,omitempty"`
}

// ProcessNodeDTO 流程节点DTO
type ProcessNodeDTO struct {
	ID         string             `json:"id" binding:"required"`
	Type       string             `json:"type" binding:"required"`
	Name       string             `json:"name" binding:"required"`
	Position   PositionDTO        `json:"position"`
	Properties *NodePropertiesDTO `json:"properties,omitempty"`
}

// PositionDTO 节点位置DTO
type PositionDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodePropertiesDTO 节点属性DTO
type NodePropertiesDTO struct {
	AssigneeRule     *AssigneeRuleDTO  `json:"assigneeRule,omitempty"`
	ApprovalMode     string            `json:"approvalMode,omitempty"`
	TimeoutHours     int               `json:"timeoutHours,omitempty"`
	RejectAction     string            `json:"rejectAction,omitempty"`
	FieldPermissions map[string]string `json:"fieldPermissions,omitempty"`
}

// AssigneeRuleDTO 审批人规则DTO
type AssigneeRuleDTO struct {
	Type   string  `json:"type" binding:"required"`
	Values []int64 `json:"values"`
}

// ProcessEdgeDTO 流程边DTO
type ProcessEdgeDTO struct {
	ID        string  `json:"id" binding:"required"`
	Source    string  `json:"source" binding:"required"`
	Target    string  `json:"target" binding:"required"`
	Condition *string `json:"condition,omitempty"`
}

// GlobalPropertiesDTO 全局属性DTO
type GlobalPropertiesDTO struct {
	AllowWithdraw bool `json:"allowWithdraw"`
}

// ============================================================================
// 流程实例相关请求
// ============================================================================

// StartProcessRequest 启动流程请求
type StartProcessRequest struct {
	ProcessDefID int64                  `json:"processDefId" binding:"required"`
	Title        string                 `json:"title" binding:"required,max=200"`
	FormData     map[string]interface{} `json:"formData"`
}

// CompleteTaskRequest 完成任务请求
type CompleteTaskRequest struct {
	TaskID  int64  `json:"taskId" binding:"required"`
	Action  string `json:"action" binding:"required,oneof=approve reject"`
	Comment string `json:"comment" binding:"max=500"`
}

// DelegateTaskRequest 委托任务请求
type DelegateTaskRequest struct {
	TaskID   int64  `json:"taskId" binding:"required"`
	ToUserID int64  `json:"toUserId" binding:"required"`
	Reason   string `json:"reason" binding:"max=500"`
}

// TransferTaskRequest 转办任务请求
type TransferTaskRequest struct {
	TaskID   int64  `json:"taskId" binding:"required"`
	ToUserID int64  `json:"toUserId" binding:"required"`
	Reason   string `json:"reason" binding:"max=500"`
}

// WithdrawProcessRequest 撤回流程请求
type WithdrawProcessRequest struct {
	InstanceID int64  `json:"instanceId" binding:"required"`
	Reason     string `json:"reason" binding:"max=500"`
}

// ListProcessInstRequest 获取流程实例列表请求
type ListProcessInstRequest struct {
	Page        int        `form:"page" binding:"min=1"`
	PageSize    int        `form:"pageSize" binding:"min=1,max=100"`
	Title       string     `form:"title"`
	ProcessCode string     `form:"processCode"`
	Status      string     `form:"status"`
	InitiatorID *int64     `form:"initiatorId"`
	StartDate   *time.Time `form:"startDate"`
	EndDate     *time.Time `form:"endDate"`
}

// ============================================================================
// 任务查询相关请求
// ============================================================================

// ListTaskRequest 获取任务列表请求
type ListTaskRequest struct {
	Page        int        `form:"page" binding:"min=1"`
	PageSize    int        `form:"pageSize" binding:"min=1,max=100"`
	ProcessCode string     `form:"processCode"`
	Status      string     `form:"status"`
	StartDate   *time.Time `form:"startDate"`
	EndDate     *time.Time `form:"endDate"`
}

// ============================================================================
// 表单模板相关请求
// ============================================================================

// CreateFormTemplateRequest 创建表单模板请求
type CreateFormTemplateRequest struct {
	Name        string         `json:"name" binding:"required,max=100"`
	Code        string         `json:"code" binding:"required,max=50"`
	Description string         `json:"description" binding:"max=500"`
	Schema      *FormSchemaDTO `json:"schema" binding:"required"`
}

// UpdateFormTemplateRequest 更新表单模板请求
type UpdateFormTemplateRequest struct {
	Name        string         `json:"name" binding:"required,max=100"`
	Description string         `json:"description" binding:"max=500"`
	Schema      *FormSchemaDTO `json:"schema" binding:"required"`
}

// ListFormTemplateRequest 获取表单模板列表请求
type ListFormTemplateRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   string `form:"status"`
}

// FormSchemaDTO 表单结构DTO
type FormSchemaDTO struct {
	Fields []FormFieldDTO `json:"fields"`
}

// FormFieldDTO 表单字段DTO
type FormFieldDTO struct {
	Key          string              `json:"key" binding:"required"`
	Label        string              `json:"label" binding:"required"`
	Type         string              `json:"type" binding:"required"`
	Required     bool                `json:"required"`
	Placeholder  string              `json:"placeholder,omitempty"`
	DefaultValue interface{}         `json:"defaultValue,omitempty"`
	Options      []SelectOptionDTO   `json:"options,omitempty"`
	Validation   *FieldValidationDTO `json:"validation,omitempty"`
}

// SelectOptionDTO 下拉选项DTO
type SelectOptionDTO struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

// FieldValidationDTO 字段验证规则DTO
type FieldValidationDTO struct {
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Pattern   *string  `json:"pattern,omitempty"`
}

// ============================================================================
// 统计查询相关请求
// ============================================================================

// StatisticsRequest 流程统计请求
type StatisticsRequest struct {
	ProcessCode string     `form:"processCode"`
	Category    string     `form:"category"`
	StartDate   *time.Time `form:"startDate"`
	EndDate     *time.Time `form:"endDate"`
}
