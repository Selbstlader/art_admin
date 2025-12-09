# OA流程引擎设计文档

## Overview

OA流程引擎是一个轻量级的工作流审批系统，基于现有Art Admin架构构建。采用状态机模式驱动流程流转，支持可视化流程设计、动态表单模板、多种审批模式，并与现有用户权限体系深度集成。

### 设计目标

1. **轻量级**: 自研核心引擎，避免引入重量级BPM框架
2. **可扩展**: 支持自定义节点类型、审批人规则、表单字段
3. **易集成**: 复用现有User/Role/Department模型和WebSocket通知
4. **高可用**: 流程定义版本化，运行中实例不受定义修改影响

### 技术选型

- **后端**: Go + Gin + GORM（与现有架构一致）
- **流程存储**: JSON格式存储流程定义和表单模板
- **前端流程设计器**: Vue Flow（基于Vue 3的流程图库）
- **前端表单设计器**: VueDraggable + 自定义表单组件（拖拽式表单设计）
- **消息推送**: 复用现有WebSocket Hub

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                         │
├─────────────────┬─────────────────┬─────────────────────────────┤
│  流程设计器      │   审批工作台     │      表单模板管理            │
│  (Vue Flow)     │   (待办/已办)    │      (动态表单)              │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Layer (Gin)                            │
├─────────────────┬─────────────────┬─────────────────────────────┤
│ ProcessDefAPI   │ ProcessInstAPI  │ FormTemplateAPI             │
│ - CRUD定义      │ - 发起/撤回     │ - CRUD模板                  │
│ - 发布/版本     │ - 审批/委托     │ - 字段配置                  │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Service Layer                                │
├─────────────────┬─────────────────┬─────────────────────────────┤
│ ProcessDefSvc   │ WorkflowEngine  │ FormEngineSvc               │
│ - 定义管理      │ - 实例创建      │ - 模板管理                  │
│ - 版本控制      │ - 任务流转      │ - 数据验证                  │
│ - JSON序列化    │ - 状态机驱动    │ - 权限控制                  │
├─────────────────┼─────────────────┼─────────────────────────────┤
│ AssigneeResolver│ NotificationSvc │ QueryService                │
│ - 规则解析      │ - WebSocket推送 │ - 待办/已办查询             │
│ - 用户查找      │ - 消息通知      │ - 统计分析                  │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Repository Layer                              │
├─────────────────┬─────────────────┬─────────────────────────────┤
│ ProcessDefRepo  │ ProcessInstRepo │ FormTemplateRepo            │
│ TaskRepo        │ TaskLogRepo     │ FormInstanceRepo            │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Database (MySQL)                           │
│  wf_process_def | wf_process_inst | wf_task | wf_form_template  │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. 流程定义服务 (ProcessDefinitionService)

```go
type ProcessDefinitionService interface {
    // 创建流程定义
    Create(ctx context.Context, req *CreateProcessDefRequest) (*ProcessDefinition, error)
    // 更新流程定义（草稿状态）
    Update(ctx context.Context, id int64, req *UpdateProcessDefRequest) (*ProcessDefinition, error)
    // 发布流程定义（生成新版本）
    Publish(ctx context.Context, id int64) (*ProcessDefinition, error)
    // 获取流程定义详情
    GetByID(ctx context.Context, id int64) (*ProcessDefinition, error)
    // 获取流程定义列表
    List(ctx context.Context, req *ListProcessDefRequest) ([]*ProcessDefinition, int64, error)
    // 删除流程定义（仅草稿可删）
    Delete(ctx context.Context, id int64) error
    // 序列化流程定义为JSON
    SerializeDefinition(def *ProcessDefinition) (string, error)
    // 从JSON反序列化流程定义
    DeserializeDefinition(jsonStr string) (*ProcessGraph, error)
}
```

### 2. 工作流引擎 (WorkflowEngine)

```go
type WorkflowEngine interface {
    // 启动流程实例
    StartProcess(ctx context.Context, req *StartProcessRequest) (*ProcessInstance, error)
    // 完成任务（通过/拒绝）
    CompleteTask(ctx context.Context, req *CompleteTaskRequest) error
    // 委托任务
    DelegateTask(ctx context.Context, taskID int64, toUserID int64, reason string) error
    // 转办任务
    TransferTask(ctx context.Context, taskID int64, toUserID int64, reason string) error
    // 撤回流程
    WithdrawProcess(ctx context.Context, instanceID int64) error
    // 获取流程实例详情
    GetProcessInstance(ctx context.Context, instanceID int64) (*ProcessInstanceDetail, error)
}
```

### 3. 审批人解析器 (AssigneeResolver)

```go
type AssigneeResolver interface {
    // 解析审批人规则，返回用户ID列表
    Resolve(ctx context.Context, rule *AssigneeRule, context *ProcessContext) ([]int64, error)
}

type AssigneeRule struct {
    Type   AssigneeType `json:"type"`   // user/role/dept_leader/initiator_leader
    Values []int64      `json:"values"` // 用户ID或角色ID列表
}
```

### 4. 表单引擎服务 (FormEngineService)

```go
type FormEngineService interface {
    // 创建表单模板
    CreateTemplate(ctx context.Context, req *CreateFormTemplateRequest) (*FormTemplate, error)
    // 更新表单模板
    UpdateTemplate(ctx context.Context, id int64, req *UpdateFormTemplateRequest) (*FormTemplate, error)
    // 获取表单模板
    GetTemplate(ctx context.Context, id int64) (*FormTemplate, error)
    // 获取表单模板列表
    ListTemplates(ctx context.Context, req *ListFormTemplateRequest) ([]*FormTemplate, int64, error)
    // 删除表单模板
    DeleteTemplate(ctx context.Context, id int64) error
    // 验证表单数据
    ValidateFormData(ctx context.Context, templateID int64, data map[string]interface{}) error
    // 序列化表单模板为JSON
    SerializeTemplate(template *FormTemplate) (string, error)
    // 从JSON反序列化表单模板
    DeserializeTemplate(jsonStr string) (*FormSchema, error)
    // 获取节点字段权限
    GetFieldPermissions(ctx context.Context, templateID int64, nodeID string) (*FieldPermissions, error)
}
```

### 5. 通知服务 (NotificationService)

```go
type NotificationService interface {
    // 发送待办通知
    NotifyNewTask(ctx context.Context, task *Task) error
    // 发送流程状态变更通知
    NotifyProcessStatusChange(ctx context.Context, instance *ProcessInstance) error
    // 发送超时提醒
    NotifyTaskTimeout(ctx context.Context, task *Task) error
}
```

### 6. 查询服务 (QueryService)

```go
type QueryService interface {
    // 查询我发起的流程
    ListMyInitiated(ctx context.Context, userID int64, req *QueryRequest) ([]*ProcessInstanceSummary, int64, error)
    // 查询我的待办
    ListMyTodo(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error)
    // 查询我的已办
    ListMyDone(ctx context.Context, userID int64, req *QueryRequest) ([]*TaskSummary, int64, error)
    // 获取审批轨迹
    GetApprovalTrail(ctx context.Context, instanceID int64) ([]*ApprovalRecord, error)
    // 获取流程统计
    GetStatistics(ctx context.Context, req *StatisticsRequest) (*ProcessStatistics, error)
}
```

## Data Models

### 流程定义 (ProcessDefinition)

```go
type ProcessDefinition struct {
    ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    Name           string         `gorm:"type:varchar(100);not null" json:"name"`
    Code           string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`
    Description    string         `gorm:"type:varchar(500)" json:"description"`
    Category       string         `gorm:"type:varchar(50);index" json:"category"`
    FormTemplateID *int64         `gorm:"index" json:"formTemplateId"`
    GraphJSON      string         `gorm:"type:text;not null" json:"-"`
    Version        int            `gorm:"default:1" json:"version"`
    Status         string         `gorm:"type:varchar(20);default:draft;index" json:"status"` // draft/published/disabled
    CreatedBy      int64          `gorm:"index" json:"createdBy"`
    CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProcessDefinition) TableName() string {
    return "wf_process_definition"
}
```

### 流程图结构 (ProcessGraph - JSON存储)

```go
type ProcessGraph struct {
    Nodes       []ProcessNode      `json:"nodes"`
    Edges       []ProcessEdge      `json:"edges"`
    GlobalProps *GlobalProperties  `json:"globalProps,omitempty"`
}

type ProcessNode struct {
    ID           string            `json:"id"`
    Type         NodeType          `json:"type"`         // start/end/approval/condition/parallel
    Name         string            `json:"name"`
    Position     Position          `json:"position"`
    Properties   *NodeProperties   `json:"properties,omitempty"`
}

type NodeType string
const (
    NodeTypeStart     NodeType = "start"
    NodeTypeEnd       NodeType = "end"
    NodeTypeApproval  NodeType = "approval"
    NodeTypeCondition NodeType = "condition"
    NodeTypeParallel  NodeType = "parallel"
)

type NodeProperties struct {
    AssigneeRule     *AssigneeRule     `json:"assigneeRule,omitempty"`
    ApprovalMode     ApprovalMode      `json:"approvalMode,omitempty"`     // or_sign/and_sign
    TimeoutHours     int               `json:"timeoutHours,omitempty"`
    RejectAction     RejectAction      `json:"rejectAction,omitempty"`     // terminate/return_prev
    FieldPermissions map[string]string `json:"fieldPermissions,omitempty"` // fieldKey -> visible/editable/hidden
}

type ProcessEdge struct {
    ID        string  `json:"id"`
    Source    string  `json:"source"`
    Target    string  `json:"target"`
    Condition *string `json:"condition,omitempty"` // 条件表达式
}
```

### 流程实例 (ProcessInstance)

```go
type ProcessInstance struct {
    ID                 int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    ProcessDefID       int64          `gorm:"index;not null" json:"processDefId"`
    ProcessDefVersion  int            `gorm:"not null" json:"processDefVersion"`
    Title              string         `gorm:"type:varchar(200);not null" json:"title"`
    InitiatorID        int64          `gorm:"index;not null" json:"initiatorId"`
    Status             string         `gorm:"type:varchar(20);default:running;index" json:"status"` // running/completed/rejected/withdrawn
    CurrentNodeID      string         `gorm:"type:varchar(50)" json:"currentNodeId"`
    FormData           string         `gorm:"type:text" json:"-"`
    StartedAt          time.Time      `gorm:"autoCreateTime" json:"startedAt"`
    CompletedAt        *time.Time     `json:"completedAt,omitempty"`
    DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProcessInstance) TableName() string {
    return "wf_process_instance"
}
```

### 任务 (Task)

```go
type Task struct {
    ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    ProcessInstID   int64          `gorm:"index;not null" json:"processInstId"`
    NodeID          string         `gorm:"type:varchar(50);not null" json:"nodeId"`
    NodeName        string         `gorm:"type:varchar(100)" json:"nodeName"`
    AssigneeID      int64          `gorm:"index;not null" json:"assigneeId"`
    Status          string         `gorm:"type:varchar(20);default:pending;index" json:"status"` // pending/approved/rejected/delegated/transferred
    Comment         string         `gorm:"type:varchar(500)" json:"comment"`
    DelegatedFrom   *int64         `json:"delegatedFrom,omitempty"`
    TransferredFrom *int64         `json:"transferredFrom,omitempty"`
    DueAt           *time.Time     `json:"dueAt,omitempty"`
    CreatedAt       time.Time      `gorm:"autoCreateTime" json:"createdAt"`
    CompletedAt     *time.Time     `json:"completedAt,omitempty"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
    return "wf_task"
}
```

### 任务日志 (TaskLog)

```go
type TaskLog struct {
    ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    TaskID        int64     `gorm:"index;not null" json:"taskId"`
    ProcessInstID int64     `gorm:"index;not null" json:"processInstId"`
    OperatorID    int64     `gorm:"index" json:"operatorId"`
    Action        string    `gorm:"type:varchar(20);not null" json:"action"` // approve/reject/delegate/transfer/withdraw
    Comment       string    `gorm:"type:varchar(500)" json:"comment"`
    CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (TaskLog) TableName() string {
    return "wf_task_log"
}
```

### 表单模板 (FormTemplate)

```go
type FormTemplate struct {
    ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string         `gorm:"type:varchar(100);not null" json:"name"`
    Code        string         `gorm:"type:varchar(50);uniqueIndex" json:"code"`
    Description string         `gorm:"type:varchar(500)" json:"description"`
    SchemaJSON  string         `gorm:"type:text;not null" json:"-"`
    Status      string         `gorm:"type:varchar(20);default:active;index" json:"status"` // active/disabled
    CreatedBy   int64          `gorm:"index" json:"createdBy"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FormTemplate) TableName() string {
    return "wf_form_template"
}
```

### 表单结构 (FormSchema - JSON存储)

```go
type FormSchema struct {
    Fields []FormField `json:"fields"`
}

type FormField struct {
    Key          string            `json:"key"`
    Label        string            `json:"label"`
    Type         FieldType         `json:"type"`         // text/number/date/select/file/textarea
    Required     bool              `json:"required"`
    Placeholder  string            `json:"placeholder,omitempty"`
    DefaultValue interface{}       `json:"defaultValue,omitempty"`
    Options      []SelectOption    `json:"options,omitempty"`    // for select type
    Validation   *FieldValidation  `json:"validation,omitempty"`
}

type FieldType string
const (
    FieldTypeText     FieldType = "text"
    FieldTypeNumber   FieldType = "number"
    FieldTypeDate     FieldType = "date"
    FieldTypeSelect   FieldType = "select"
    FieldTypeFile     FieldType = "file"
    FieldTypeTextarea FieldType = "textarea"
)

type FieldValidation struct {
    MinLength *int     `json:"minLength,omitempty"`
    MaxLength *int     `json:"maxLength,omitempty"`
    Min       *float64 `json:"min,omitempty"`
    Max       *float64 `json:"max,omitempty"`
    Pattern   *string  `json:"pattern,omitempty"`
}
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: 流程定义序列化往返一致性

*For any* 有效的ProcessGraph对象，序列化为JSON后再反序列化，应产生与原始对象等价的ProcessGraph。

**Validates: Requirements 1.6, 1.7**

### Property 2: 表单模板序列化往返一致性

*For any* 有效的FormSchema对象，序列化为JSON后再反序列化，应产生与原始对象等价的FormSchema。

**Validates: Requirements 4.6, 4.7**

### Property 3: 流程结构完整性验证

*For any* ProcessGraph，验证函数应正确识别：必须有且仅有一个开始节点、至少一个结束节点、所有节点可达、无孤立节点。

**Validates: Requirements 1.1**

### Property 4: 流程发布版本递增

*For any* 已发布的ProcessDefinition，再次发布修改后的版本时，新版本号应大于原版本号，且原版本保持不变。

**Validates: Requirements 1.4, 1.5**

### Property 5: 流程实例初始状态正确性

*For any* 新创建的ProcessInstance，其状态应为"running"，且应存在至少一个状态为"pending"的Task。

**Validates: Requirements 2.1, 2.2**

### Property 6: 表单数据验证正确性

*For any* FormSchema和表单数据，如果数据不满足必填字段、格式或范围约束，验证函数应返回错误；如果满足所有约束，应返回成功。

**Validates: Requirements 2.6, 4.3**

### Property 7: 审批模式流转正确性

*For any* 配置为会签模式的审批节点，只有当所有Assignee都完成时才流转；*For any* 配置为或签模式的审批节点，任一Assignee完成即流转。

**Validates: Requirements 3.3, 3.4**

### Property 8: 任务通过后流转正确性

*For any* 被通过的Task，当前Task状态应变为"approved"，且应根据Transition创建下一节点的Task（除非是结束节点）。

**Validates: Requirements 3.1**

### Property 9: 流程状态终态一致性

*For any* ProcessInstance，当到达结束节点时状态应为"completed"；当被撤回时状态应为"withdrawn"；当被拒绝且配置为终止时状态应为"rejected"。

**Validates: Requirements 2.3, 2.4, 2.5**

### Property 10: 委托/转办记录完整性

*For any* 被委托或转办的Task，应记录原Assignee信息，且新Assignee应能查询到该任务。

**Validates: Requirements 3.5, 3.6**

### Property 11: 审批人规则解析正确性

*For any* AssigneeRule：
- 指定用户规则应返回配置的用户ID列表
- 指定角色规则应返回拥有该角色的所有用户
- 部门负责人规则应返回发起人所属部门的负责人
- 发起人上级规则应返回发起人的直接上级

**Validates: Requirements 5.1, 5.2, 5.3, 5.4**

### Property 12: 审批人规则空结果处理

*For any* AssigneeRule解析结果为空时，系统应返回错误而非创建无Assignee的Task。

**Validates: Requirements 5.5**

### Property 13: 查询结果过滤正确性

*For any* 用户查询：
- "我发起的"查询结果中所有实例的initiatorId应等于查询用户ID
- "我的待办"查询结果中所有任务的assigneeId应等于查询用户ID且状态为pending
- "我的已办"查询结果中所有任务应由查询用户处理过

**Validates: Requirements 7.1, 7.2, 7.3**

### Property 14: 条件表达式求值正确性

*For any* 条件表达式和上下文数据，表达式求值结果应与预期的布尔值一致。

**Validates: Requirements 1.3**

## Frontend Components

### 表单设计器 (FormDesigner)

采用拖拽式设计，用户可从左侧组件面板拖拽字段到中间画布区域，右侧显示字段属性配置。

```
┌─────────────────────────────────────────────────────────────────┐
│                      表单设计器                                  │
├──────────────┬─────────────────────────┬────────────────────────┤
│   组件面板    │       画布区域          │      属性面板          │
│              │                         │                        │
│  📝 文本输入  │  ┌─────────────────┐   │  字段标识: name        │
│  🔢 数字输入  │  │ 申请人姓名      │   │  字段名称: 申请人姓名   │
│  📅 日期选择  │  │ [文本输入]      │   │  是否必填: ✓           │
│  📋 下拉选择  │  └─────────────────┘   │  占位提示: 请输入...   │
│  📄 多行文本  │  ┌─────────────────┐   │  验证规则:             │
│  📎 文件上传  │  │ 请假天数        │   │    最小长度: 2         │
│              │  │ [数字输入]      │   │    最大长度: 50        │
│              │  └─────────────────┘   │                        │
│              │  ┌─────────────────┐   │                        │
│              │  │ 开始日期        │   │                        │
│              │  │ [日期选择]      │   │                        │
│              │  └─────────────────┘   │                        │
│              │                         │                        │
│              │  [拖拽字段到此处]       │                        │
└──────────────┴─────────────────────────┴────────────────────────┘
```

#### 技术实现

```typescript
// 表单设计器组件结构
interface FormDesignerProps {
  modelValue: FormSchema;        // 表单结构
  readonly?: boolean;            // 只读模式
}

// 组件面板配置
const componentPalette: ComponentItem[] = [
  { type: 'text', label: '文本输入', icon: 'edit' },
  { type: 'number', label: '数字输入', icon: 'calculator' },
  { type: 'date', label: '日期选择', icon: 'calendar' },
  { type: 'select', label: '下拉选择', icon: 'list' },
  { type: 'textarea', label: '多行文本', icon: 'document' },
  { type: 'file', label: '文件上传', icon: 'upload' },
];

// 拖拽事件处理
const handleDrop = (event: DragEvent, index: number) => {
  const fieldType = event.dataTransfer?.getData('fieldType');
  const newField: FormField = createDefaultField(fieldType);
  formSchema.fields.splice(index, 0, newField);
};
```

#### 依赖库

- `vuedraggable@4.x`: 基于Sortable.js的Vue 3拖拽组件
- `@vueuse/core`: 拖拽相关hooks

### 流程设计器 (ProcessDesigner)

基于Vue Flow实现的可视化流程编排器。

```
┌─────────────────────────────────────────────────────────────────┐
│  工具栏: [保存] [发布] [预览] [撤销] [重做]                       │
├──────────────┬──────────────────────────────────────────────────┤
│   节点面板    │              画布区域                            │
│              │                                                  │
│  ⚪ 开始节点  │      ⚪ ──────► 🔲 ──────► 🔲 ──────► ⚫        │
│  🔲 审批节点  │     开始       审批1      审批2      结束        │
│  ◇ 条件节点  │                                                  │
│  ⚫ 结束节点  │                                                  │
│              │                                                  │
└──────────────┴──────────────────────────────────────────────────┤
│  属性面板: 审批节点配置                                          │
│  ├─ 节点名称: 部门经理审批                                       │
│  ├─ 审批人规则: [部门负责人 ▼]                                   │
│  ├─ 审批模式: ○ 或签  ● 会签                                    │
│  ├─ 超时时间: 24 小时                                           │
│  └─ 拒绝处理: ○ 终止流程  ● 退回上一节点                         │
└─────────────────────────────────────────────────────────────────┘
```

#### 依赖库

- `@vue-flow/core@1.x`: Vue 3流程图核心库
- `@vue-flow/background`: 背景网格
- `@vue-flow/controls`: 缩放控制
- `@vue-flow/minimap`: 小地图导航

## Error Handling

### 流程定义错误

| 错误场景 | 错误码 | 处理方式 |
|---------|--------|---------|
| 流程结构无效（无开始/结束节点） | WF_DEF_INVALID_STRUCTURE | 返回验证错误详情 |
| 流程定义不存在 | WF_DEF_NOT_FOUND | 返回404 |
| 已发布流程不可删除 | WF_DEF_CANNOT_DELETE | 返回400 |
| JSON解析失败 | WF_DEF_PARSE_ERROR | 返回解析错误位置 |

### 流程实例错误

| 错误场景 | 错误码 | 处理方式 |
|---------|--------|---------|
| 流程定义未发布 | WF_INST_DEF_NOT_PUBLISHED | 返回400 |
| 表单数据验证失败 | WF_INST_FORM_INVALID | 返回字段级错误 |
| 流程实例不存在 | WF_INST_NOT_FOUND | 返回404 |
| 无权撤回（非发起人） | WF_INST_NO_PERMISSION | 返回403 |

### 任务处理错误

| 错误场景 | 错误码 | 处理方式 |
|---------|--------|---------|
| 任务不存在 | WF_TASK_NOT_FOUND | 返回404 |
| 非任务Assignee | WF_TASK_NOT_ASSIGNEE | 返回403 |
| 任务已处理 | WF_TASK_ALREADY_COMPLETED | 返回400 |
| 审批人解析为空 | WF_TASK_NO_ASSIGNEE | 记录日志，通知管理员 |

### 表单模板错误

| 错误场景 | 错误码 | 处理方式 |
|---------|--------|---------|
| 模板不存在 | WF_FORM_NOT_FOUND | 返回404 |
| 模板已被流程引用不可删除 | WF_FORM_IN_USE | 返回400 |
| Schema解析失败 | WF_FORM_PARSE_ERROR | 返回解析错误 |

## Testing Strategy

### 测试框架选择

- **单元测试**: Go标准库 `testing` + `testify/assert`
- **属性测试**: `github.com/leanovate/gopter` (Go Property-Based Testing)
- **集成测试**: 使用SQLite内存数据库

### 单元测试覆盖

1. **流程定义服务**
   - 创建/更新/删除流程定义
   - 发布版本控制
   - 流程结构验证

2. **工作流引擎**
   - 流程启动
   - 任务流转（通过/拒绝）
   - 会签/或签逻辑
   - 委托/转办

3. **表单引擎**
   - 模板CRUD
   - 字段验证规则
   - 节点权限控制

4. **审批人解析器**
   - 四种规则类型解析
   - 空结果处理

5. **查询服务**
   - 待办/已办/我发起的查询
   - 分页和过滤

### 属性测试要求

- 每个属性测试配置运行 **100次** 迭代
- 每个属性测试必须注释引用设计文档中的正确性属性
- 注释格式: `// **Feature: oa-workflow-engine, Property {number}: {property_text}**`

### 属性测试覆盖

| 属性 | 测试文件 | 生成器 |
|-----|---------|--------|
| Property 1 | `process_def_test.go` | 随机ProcessGraph生成器 |
| Property 2 | `form_template_test.go` | 随机FormSchema生成器 |
| Property 3 | `process_def_test.go` | 有效/无效流程图生成器 |
| Property 6 | `form_engine_test.go` | 随机表单数据+Schema生成器 |
| Property 7 | `workflow_engine_test.go` | 会签/或签场景生成器 |
| Property 11 | `assignee_resolver_test.go` | 随机规则+组织架构生成器 |
| Property 13 | `query_service_test.go` | 随机实例/任务数据生成器 |
