package model

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================================
// 流程定义相关模型
// ============================================================================

// ProcessDefinition 流程定义模型
type ProcessDefinition struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`             // 流程名称
	Code           string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`           // 流程编码
	Description    string         `gorm:"type:varchar(500)" json:"description"`               // 流程描述
	Category       string         `gorm:"type:varchar(50);index" json:"category"`             // 流程分类
	FormTemplateID *int64         `gorm:"index" json:"formTemplateId"`                        // 关联的表单模板ID
	GraphJSON      string         `gorm:"type:text;not null" json:"-"`                        // 流程图JSON
	Version        int            `gorm:"default:1" json:"version"`                           // 版本号
	Status         string         `gorm:"type:varchar(20);default:draft;index" json:"status"` // 状态: draft/published/disabled
	CreatedBy      int64          `gorm:"index" json:"createdBy"`                             // 创建人ID
	CreatedAt      time.Time      `gorm:"type:datetime;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"type:datetime;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ProcessDefinition) TableName() string {
	return "wf_process_definition"
}

// ProcessDefinitionStatus 流程定义状态常量
const (
	ProcessDefStatusDraft     = "draft"     // 草稿
	ProcessDefStatusPublished = "published" // 已发布
	ProcessDefStatusDisabled  = "disabled"  // 已禁用
)

// ============================================================================
// 流程图结构 (JSON存储)
// ============================================================================

// ProcessGraph 流程图结构
type ProcessGraph struct {
	Nodes       []ProcessNode     `json:"nodes"`                 // 节点列表
	Edges       []ProcessEdge     `json:"edges"`                 // 边列表
	GlobalProps *GlobalProperties `json:"globalProps,omitempty"` // 全局属性
}

// ProcessNode 流程节点
type ProcessNode struct {
	ID         string          `json:"id"`                   // 节点ID
	Type       NodeType        `json:"type"`                 // 节点类型
	Name       string          `json:"name"`                 // 节点名称
	Position   Position        `json:"position"`             // 节点位置
	Properties *NodeProperties `json:"properties,omitempty"` // 节点属性
}

// Position 节点位置
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeType 节点类型
type NodeType string

const (
	NodeTypeStart     NodeType = "start"     // 开始节点
	NodeTypeEnd       NodeType = "end"       // 结束节点
	NodeTypeApproval  NodeType = "approval"  // 审批节点
	NodeTypeCondition NodeType = "condition" // 条件节点
	NodeTypeParallel  NodeType = "parallel"  // 并行节点
)

// NodeProperties 节点属性
type NodeProperties struct {
	AssigneeRule     *AssigneeRule     `json:"assigneeRule,omitempty"`     // 审批人规则
	ApprovalMode     ApprovalMode      `json:"approvalMode,omitempty"`     // 审批模式
	TimeoutHours     int               `json:"timeoutHours,omitempty"`     // 超时时间(小时)
	RejectAction     RejectAction      `json:"rejectAction,omitempty"`     // 拒绝处理方式
	FieldPermissions map[string]string `json:"fieldPermissions,omitempty"` // 字段权限: fieldKey -> visible/editable/hidden
}

// ApprovalMode 审批模式
type ApprovalMode string

const (
	ApprovalModeOrSign  ApprovalMode = "or_sign"  // 或签(任一人通过)
	ApprovalModeAndSign ApprovalMode = "and_sign" // 会签(所有人通过)
)

// RejectAction 拒绝处理方式
type RejectAction string

const (
	RejectActionTerminate  RejectAction = "terminate"   // 终止流程
	RejectActionReturnPrev RejectAction = "return_prev" // 退回上一节点
)

// AssigneeRule 审批人规则
type AssigneeRule struct {
	Type   AssigneeType `json:"type"`   // 规则类型
	Values []int64      `json:"values"` // 用户ID或角色ID列表
}

// AssigneeType 审批人规则类型
type AssigneeType string

const (
	AssigneeTypeUser            AssigneeType = "user"             // 指定用户
	AssigneeTypeRole            AssigneeType = "role"             // 指定角色
	AssigneeTypeDeptLeader      AssigneeType = "dept_leader"      // 部门负责人
	AssigneeTypeInitiatorLeader AssigneeType = "initiator_leader" // 发起人上级
)

// ProcessEdge 流程边(连线)
type ProcessEdge struct {
	ID        string  `json:"id"`                  // 边ID
	Source    string  `json:"source"`              // 源节点ID
	Target    string  `json:"target"`              // 目标节点ID
	Condition *string `json:"condition,omitempty"` // 条件表达式
}

// GlobalProperties 全局属性
type GlobalProperties struct {
	AllowWithdraw bool `json:"allowWithdraw"` // 是否允许撤回
}

// ============================================================================
// 流程实例相关模型
// ============================================================================

// ProcessInstance 流程实例模型
type ProcessInstance struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProcessDefID      int64          `gorm:"index;not null" json:"processDefId"`                   // 流程定义ID
	ProcessDefVersion int            `gorm:"not null" json:"processDefVersion"`                    // 流程定义版本
	Title             string         `gorm:"type:varchar(200);not null" json:"title"`              // 流程标题
	InitiatorID       int64          `gorm:"index;not null" json:"initiatorId"`                    // 发起人ID
	InitiatorName     string         `gorm:"type:varchar(50)" json:"initiatorName"`                // 发起人名称
	Status            string         `gorm:"type:varchar(20);default:running;index" json:"status"` // 状态: running/completed/rejected/withdrawn
	CurrentNodeID     string         `gorm:"type:varchar(50)" json:"currentNodeId"`                // 当前节点ID
	FormData          string         `gorm:"type:text" json:"-"`                                   // 表单数据JSON
	StartedAt         time.Time      `gorm:"type:datetime;not null;autoCreateTime" json:"startedAt"`
	CompletedAt       *time.Time     `gorm:"type:datetime" json:"completedAt,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ProcessInstance) TableName() string {
	return "wf_process_instance"
}

// ProcessInstanceStatus 流程实例状态常量
const (
	ProcessInstStatusRunning   = "running"   // 运行中
	ProcessInstStatusCompleted = "completed" // 已完成
	ProcessInstStatusRejected  = "rejected"  // 已拒绝
	ProcessInstStatusWithdrawn = "withdrawn" // 已撤回
)

// ============================================================================
// 任务相关模型
// ============================================================================

// WfTask 工作流任务模型 (避免与项目管理Task冲突)
type WfTask struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProcessInstID   int64          `gorm:"index;not null" json:"processInstId"`                  // 流程实例ID
	NodeID          string         `gorm:"type:varchar(50);not null" json:"nodeId"`              // 节点ID
	NodeName        string         `gorm:"type:varchar(100)" json:"nodeName"`                    // 节点名称
	AssigneeID      int64          `gorm:"index;not null" json:"assigneeId"`                     // 审批人ID
	AssigneeName    string         `gorm:"type:varchar(50)" json:"assigneeName"`                 // 审批人名称
	Status          string         `gorm:"type:varchar(20);default:pending;index" json:"status"` // 状态: pending/approved/rejected/delegated/transferred
	Comment         string         `gorm:"type:varchar(500)" json:"comment"`                     // 审批意见
	DelegatedFrom   *int64         `json:"delegatedFrom,omitempty"`                              // 委托来源用户ID
	TransferredFrom *int64         `json:"transferredFrom,omitempty"`                            // 转办来源用户ID
	DueAt           *time.Time     `gorm:"type:datetime" json:"dueAt,omitempty"`                 // 截止时间
	CreatedAt       time.Time      `gorm:"type:datetime;not null;autoCreateTime" json:"createdAt"`
	CompletedAt     *time.Time     `gorm:"type:datetime" json:"completedAt,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (WfTask) TableName() string {
	return "wf_task"
}

// WfTaskStatus 任务状态常量
const (
	WfTaskStatusPending     = "pending"     // 待处理
	WfTaskStatusApproved    = "approved"    // 已通过
	WfTaskStatusRejected    = "rejected"    // 已拒绝
	WfTaskStatusDelegated   = "delegated"   // 已委托
	WfTaskStatusTransferred = "transferred" // 已转办
)

// TaskLog 任务日志模型
type TaskLog struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID        int64     `gorm:"index;not null" json:"taskId"`            // 任务ID
	ProcessInstID int64     `gorm:"index;not null" json:"processInstId"`     // 流程实例ID
	OperatorID    int64     `gorm:"index" json:"operatorId"`                 // 操作人ID
	OperatorName  string    `gorm:"type:varchar(50)" json:"operatorName"`    // 操作人名称
	Action        string    `gorm:"type:varchar(20);not null" json:"action"` // 操作类型: approve/reject/delegate/transfer/withdraw
	Comment       string    `gorm:"type:varchar(500)" json:"comment"`        // 操作备注
	CreatedAt     time.Time `gorm:"type:datetime;not null;autoCreateTime" json:"createdAt"`
}

// TableName 表名
func (TaskLog) TableName() string {
	return "wf_task_log"
}

// TaskLogAction 任务日志操作类型常量
const (
	TaskLogActionApprove  = "approve"  // 通过
	TaskLogActionReject   = "reject"   // 拒绝
	TaskLogActionDelegate = "delegate" // 委托
	TaskLogActionTransfer = "transfer" // 转办
	TaskLogActionWithdraw = "withdraw" // 撤回
)

// ============================================================================
// 表单模板相关模型
// ============================================================================

// FormTemplate 表单模板模型
type FormTemplate struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`              // 模板名称
	Code        string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`            // 模板编码
	Description string         `gorm:"type:varchar(500)" json:"description"`                // 模板描述
	SchemaJSON  string         `gorm:"type:text;not null" json:"-"`                         // 表单结构JSON
	Status      string         `gorm:"type:varchar(20);default:active;index" json:"status"` // 状态: active/disabled
	CreatedBy   int64          `gorm:"index" json:"createdBy"`                              // 创建人ID
	CreatedAt   time.Time      `gorm:"type:datetime;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (FormTemplate) TableName() string {
	return "wf_form_template"
}

// FormTemplateStatus 表单模板状态常量
const (
	FormTemplateStatusActive   = "active"   // 启用
	FormTemplateStatusDisabled = "disabled" // 禁用
)

// ============================================================================
// 表单结构 (JSON存储)
// ============================================================================

// FormSchema 表单结构
type FormSchema struct {
	Fields []FormField `json:"fields"` // 字段列表
}

// FormField 表单字段
type FormField struct {
	Key          string           `json:"key"`                    // 字段标识
	Label        string           `json:"label"`                  // 字段标签
	Type         FieldType        `json:"type"`                   // 字段类型
	Required     bool             `json:"required"`               // 是否必填
	Placeholder  string           `json:"placeholder,omitempty"`  // 占位提示
	DefaultValue interface{}      `json:"defaultValue,omitempty"` // 默认值
	Options      []SelectOption   `json:"options,omitempty"`      // 选项(下拉选择用)
	Validation   *FieldValidation `json:"validation,omitempty"`   // 验证规则
}

// FieldType 字段类型
type FieldType string

const (
	FieldTypeText     FieldType = "text"     // 文本输入
	FieldTypeNumber   FieldType = "number"   // 数字输入
	FieldTypeDate     FieldType = "date"     // 日期选择
	FieldTypeSelect   FieldType = "select"   // 下拉选择
	FieldTypeFile     FieldType = "file"     // 文件上传
	FieldTypeTextarea FieldType = "textarea" // 多行文本
)

// SelectOption 下拉选项
type SelectOption struct {
	Label string      `json:"label"` // 选项标签
	Value interface{} `json:"value"` // 选项值
}

// FieldValidation 字段验证规则
type FieldValidation struct {
	MinLength *int     `json:"minLength,omitempty"` // 最小长度
	MaxLength *int     `json:"maxLength,omitempty"` // 最大长度
	Min       *float64 `json:"min,omitempty"`       // 最小值
	Max       *float64 `json:"max,omitempty"`       // 最大值
	Pattern   *string  `json:"pattern,omitempty"`   // 正则表达式
}

// FieldPermission 字段权限
type FieldPermission string

const (
	FieldPermissionVisible  FieldPermission = "visible"  // 可见
	FieldPermissionEditable FieldPermission = "editable" // 可编辑
	FieldPermissionHidden   FieldPermission = "hidden"   // 隐藏
)
