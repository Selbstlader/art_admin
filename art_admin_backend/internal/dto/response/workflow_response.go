package response

import "time"

// ============================================================================
// 流程定义相关响应
// ============================================================================

// ProcessDefResponse 流程定义响应
type ProcessDefResponse struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	FormTemplateID *int64    `json:"formTemplateId"`
	Version        int       `json:"version"`
	Status         string    `json:"status"`
	CreatedBy      int64     `json:"createdBy"`
	CreatedByName  string    `json:"createdByName"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ProcessDefDetailResponse 流程定义详情响应
type ProcessDefDetailResponse struct {
	ProcessDefResponse
	Graph            *ProcessGraphResponse `json:"graph"`
	FormTemplateName string                `json:"formTemplateName,omitempty"`
}

// ProcessGraphResponse 流程图响应
type ProcessGraphResponse struct {
	Nodes       []ProcessNodeResponse     `json:"nodes"`
	Edges       []ProcessEdgeResponse     `json:"edges"`
	GlobalProps *GlobalPropertiesResponse `json:"globalProps,omitempty"`
}

// ProcessNodeResponse 流程节点响应
type ProcessNodeResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Name       string                  `json:"name"`
	Position   PositionResponse        `json:"position"`
	Properties *NodePropertiesResponse `json:"properties,omitempty"`
}

// PositionResponse 节点位置响应
type PositionResponse struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodePropertiesResponse 节点属性响应
type NodePropertiesResponse struct {
	AssigneeRule     *AssigneeRuleResponse `json:"assigneeRule,omitempty"`
	ApprovalMode     string                `json:"approvalMode,omitempty"`
	TimeoutHours     int                   `json:"timeoutHours,omitempty"`
	RejectAction     string                `json:"rejectAction,omitempty"`
	FieldPermissions map[string]string     `json:"fieldPermissions,omitempty"`
}

// AssigneeRuleResponse 审批人规则响应
type AssigneeRuleResponse struct {
	Type   string  `json:"type"`
	Values []int64 `json:"values"`
}

// ProcessEdgeResponse 流程边响应
type ProcessEdgeResponse struct {
	ID        string  `json:"id"`
	Source    string  `json:"source"`
	Target    string  `json:"target"`
	Condition *string `json:"condition,omitempty"`
}

// GlobalPropertiesResponse 全局属性响应
type GlobalPropertiesResponse struct {
	AllowWithdraw bool `json:"allowWithdraw"`
}

// ============================================================================
// 流程实例相关响应
// ============================================================================

// ProcessInstResponse 流程实例响应
type ProcessInstResponse struct {
	ID                int64      `json:"id"`
	ProcessDefID      int64      `json:"processDefId"`
	ProcessDefVersion int        `json:"processDefVersion"`
	ProcessName       string     `json:"processName"`
	ProcessCode       string     `json:"processCode"`
	Title             string     `json:"title"`
	InitiatorID       int64      `json:"initiatorId"`
	InitiatorName     string     `json:"initiatorName"`
	Status            string     `json:"status"`
	CurrentNodeID     string     `json:"currentNodeId"`
	CurrentNodeName   string     `json:"currentNodeName"`
	StartedAt         time.Time  `json:"startedAt"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
}

// ProcessInstDetailResponse 流程实例详情响应
type ProcessInstDetailResponse struct {
	ProcessInstResponse
	FormData      map[string]interface{}   `json:"formData"`
	ApprovalTrail []ApprovalRecordResponse `json:"approvalTrail"`
	CurrentTasks  []WfTaskResponse         `json:"currentTasks"`
}

// ProcessInstSummaryResponse 流程实例摘要响应(用于列表)
type ProcessInstSummaryResponse struct {
	ID            int64      `json:"id"`
	ProcessName   string     `json:"processName"`
	ProcessCode   string     `json:"processCode"`
	Title         string     `json:"title"`
	InitiatorID   int64      `json:"initiatorId"`
	InitiatorName string     `json:"initiatorName"`
	Status        string     `json:"status"`
	StartedAt     time.Time  `json:"startedAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

// ============================================================================
// 任务相关响应
// ============================================================================

// WfTaskResponse 工作流任务响应
type WfTaskResponse struct {
	ID              int64      `json:"id"`
	ProcessInstID   int64      `json:"processInstId"`
	ProcessTitle    string     `json:"processTitle"`
	ProcessName     string     `json:"processName"`
	ProcessCode     string     `json:"processCode"`
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

// WfTaskSummaryResponse 任务摘要响应(用于待办/已办列表)
type WfTaskSummaryResponse struct {
	ID            int64      `json:"id"`
	ProcessInstID int64      `json:"processInstId"`
	ProcessTitle  string     `json:"processTitle"`
	ProcessName   string     `json:"processName"`
	ProcessCode   string     `json:"processCode"`
	NodeName      string     `json:"nodeName"`
	InitiatorID   int64      `json:"initiatorId"`
	InitiatorName string     `json:"initiatorName"`
	Status        string     `json:"status"`
	DueAt         *time.Time `json:"dueAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

// ApprovalRecordResponse 审批记录响应
type ApprovalRecordResponse struct {
	TaskID       int64     `json:"taskId"`
	NodeID       string    `json:"nodeId"`
	NodeName     string    `json:"nodeName"`
	OperatorID   int64     `json:"operatorId"`
	OperatorName string    `json:"operatorName"`
	Action       string    `json:"action"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ============================================================================
// 表单模板相关响应
// ============================================================================

// FormTemplateResponse 表单模板响应
type FormTemplateResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedBy     int64     `json:"createdBy"`
	CreatedByName string    `json:"createdByName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// FormTemplateDetailResponse 表单模板详情响应
type FormTemplateDetailResponse struct {
	FormTemplateResponse
	Schema *FormSchemaResponse `json:"schema"`
}

// FormSchemaResponse 表单结构响应
type FormSchemaResponse struct {
	Fields []FormFieldResponse `json:"fields"`
}

// FormFieldResponse 表单字段响应
type FormFieldResponse struct {
	Key          string                   `json:"key"`
	Label        string                   `json:"label"`
	Type         string                   `json:"type"`
	Required     bool                     `json:"required"`
	Placeholder  string                   `json:"placeholder,omitempty"`
	DefaultValue interface{}              `json:"defaultValue,omitempty"`
	Options      []SelectOptionResponse   `json:"options,omitempty"`
	Validation   *FieldValidationResponse `json:"validation,omitempty"`
}

// SelectOptionResponse 下拉选项响应
type SelectOptionResponse struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

// FieldValidationResponse 字段验证规则响应
type FieldValidationResponse struct {
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Pattern   *string  `json:"pattern,omitempty"`
}

// FieldPermissionsResponse 字段权限响应
type FieldPermissionsResponse struct {
	Fields map[string]string `json:"fields"` // fieldKey -> visible/editable/hidden
}

// ============================================================================
// 统计相关响应
// ============================================================================

// ProcessStatisticsResponse 流程统计响应
type ProcessStatisticsResponse struct {
	TotalInstances     int64                   `json:"totalInstances"`
	RunningInstances   int64                   `json:"runningInstances"`
	CompletedInstances int64                   `json:"completedInstances"`
	RejectedInstances  int64                   `json:"rejectedInstances"`
	WithdrawnInstances int64                   `json:"withdrawnInstances"`
	AvgProcessingTime  float64                 `json:"avgProcessingTime"` // 平均处理时长(小时)
	ByCategory         []CategoryStatistics    `json:"byCategory,omitempty"`
	ByProcess          []ProcessTypeStatistics `json:"byProcess,omitempty"`
}

// CategoryStatistics 按分类统计
type CategoryStatistics struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

// ProcessTypeStatistics 按流程类型统计
type ProcessTypeStatistics struct {
	ProcessCode       string  `json:"processCode"`
	ProcessName       string  `json:"processName"`
	TotalCount        int64   `json:"totalCount"`
	CompletedCount    int64   `json:"completedCount"`
	AvgProcessingTime float64 `json:"avgProcessingTime"`
}

// ============================================================================
// 通用分页响应
// ============================================================================

// ProcessDefListResponse 流程定义列表响应
type ProcessDefListResponse struct {
	List  []ProcessDefResponse `json:"list"`
	Total int64                `json:"total"`
}

// ProcessInstListResponse 流程实例列表响应
type ProcessInstListResponse struct {
	List  []ProcessInstSummaryResponse `json:"list"`
	Total int64                        `json:"total"`
}

// WfTaskListResponse 任务列表响应
type WfTaskListResponse struct {
	List  []WfTaskSummaryResponse `json:"list"`
	Total int64                   `json:"total"`
}

// FormTemplateListResponse 表单模板列表响应
type FormTemplateListResponse struct {
	List  []FormTemplateResponse `json:"list"`
	Total int64                  `json:"total"`
}

// ============================================================================
// 简单选项响应(用于下拉选择)
// ============================================================================

// ProcessDefOptionResponse 流程定义选项响应
type ProcessDefOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// FormTemplateOptionResponse 表单模板选项响应
type FormTemplateOptionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
