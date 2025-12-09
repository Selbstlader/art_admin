# Requirements Document

## Introduction

本文档定义了OA流程引擎的功能需求，该引擎将为Art Admin系统提供完整的工作流审批能力。流程引擎支持可视化流程设计、多种审批模式、动态表单配置，并与现有的用户权限体系无缝集成，满足企业日常办公中的请假、报销、采购等审批场景。

## Glossary

- **Workflow_Engine（流程引擎）**: 负责流程定义解析、实例创建、任务流转的核心服务组件
- **Process_Definition（流程定义）**: 描述审批流程结构的模板，包含节点、连线、条件等元素
- **Process_Instance（流程实例）**: 基于流程定义创建的具体审批单据
- **Task（任务）**: 流程实例中分配给特定用户的待办事项
- **Node（节点）**: 流程中的处理单元，包括开始节点、审批节点、条件节点、结束节点等
- **Transition（流转）**: 节点之间的连接，定义流程走向
- **Assignee（审批人）**: 被分配处理任务的用户
- **Approval_Mode（审批模式）**: 任务处理方式，包括或签（任一人通过）、会签（所有人通过）
- **Form_Template（表单模板）**: 可复用的表单结构定义，包含字段配置、验证规则等，可被多个流程定义引用
- **Form_Instance（表单实例）**: 基于Form_Template创建的具体表单数据，与Process_Instance关联
- **Delegation（委托）**: 将待办任务转交给其他用户处理

## Requirements

### Requirement 1: 流程定义管理

**User Story:** As a 系统管理员, I want to 创建和管理审批流程模板, so that 可以为不同业务场景配置标准化的审批流程。

#### Acceptance Criteria

1. WHEN 管理员在流程设计器中添加节点并连接 THEN Workflow_Engine SHALL 验证流程结构完整性并保存Process_Definition
2. WHEN 管理员配置审批节点的Assignee规则 THEN Workflow_Engine SHALL 支持指定用户、指定角色、部门负责人、发起人上级四种分配方式
3. WHEN 管理员设置条件分支节点 THEN Workflow_Engine SHALL 根据表达式计算结果决定流转路径
4. WHEN 管理员发布流程定义 THEN Workflow_Engine SHALL 将Process_Definition状态更新为已发布并生成版本号
5. WHEN 管理员修改已发布的流程 THEN Workflow_Engine SHALL 创建新版本而保留原版本供运行中实例使用
6. WHEN 流程定义被序列化存储 THEN Workflow_Engine SHALL 使用JSON格式编码流程结构
7. WHEN 从存储加载流程定义 THEN Workflow_Engine SHALL 解析JSON并还原完整的Process_Definition对象

### Requirement 2: 流程实例生命周期

**User Story:** As a 普通用户, I want to 发起审批申请并跟踪进度, so that 可以完成日常办公审批事务。

#### Acceptance Criteria

1. WHEN 用户提交审批申请 THEN Workflow_Engine SHALL 基于Process_Definition创建Process_Instance并生成首个Task
2. WHEN Process_Instance创建成功 THEN Workflow_Engine SHALL 将实例状态设置为运行中并记录发起时间
3. WHEN 流程到达结束节点 THEN Workflow_Engine SHALL 将Process_Instance状态更新为已完成
4. WHEN 用户撤回申请 THEN Workflow_Engine SHALL 终止Process_Instance并将状态更新为已撤回
5. WHEN 审批被拒绝且流程配置为拒绝即终止 THEN Workflow_Engine SHALL 将Process_Instance状态更新为已拒绝
6. IF Process_Instance创建时表单数据不符合Form_Definition约束 THEN Workflow_Engine SHALL 拒绝创建并返回验证错误信息

### Requirement 3: 任务处理与流转

**User Story:** As a 审批人, I want to 处理分配给我的审批任务, so that 可以及时完成审批工作。

#### Acceptance Criteria

1. WHEN 审批人通过任务 THEN Workflow_Engine SHALL 完成当前Task并根据Transition创建下一节点的Task
2. WHEN 审批人拒绝任务 THEN Workflow_Engine SHALL 根据流程配置决定终止流程或退回上一节点
3. WHEN 审批节点配置为会签模式 THEN Workflow_Engine SHALL 等待所有Assignee完成后才流转到下一节点
4. WHEN 审批节点配置为或签模式 THEN Workflow_Engine SHALL 在任一Assignee通过后立即流转到下一节点
5. WHEN 审批人将任务委托给他人 THEN Workflow_Engine SHALL 更新Task的Assignee并记录Delegation历史
6. WHEN 审批人转办任务 THEN Workflow_Engine SHALL 将Task完全转移给新Assignee并记录转办原因
7. WHEN 任务超过配置的时限 THEN Workflow_Engine SHALL 触发超时提醒通知

### Requirement 4: 表单模板管理

**User Story:** As a 系统管理员, I want to 创建和管理可复用的表单模板, so that 不同流程可以选择合适的表单模板来收集业务数据。

#### Acceptance Criteria

1. WHEN 管理员创建表单模板 THEN Form_Engine SHALL 保存Form_Template并生成唯一标识
2. WHEN 管理员定义表单字段 THEN Form_Engine SHALL 支持文本、数字、日期、下拉选择、文件上传等字段类型
3. WHEN 管理员配置字段验证规则 THEN Form_Engine SHALL 在表单提交时执行必填、格式、范围等校验
4. WHEN 管理员将表单模板关联到流程定义 THEN Workflow_Engine SHALL 记录Process_Definition与Form_Template的绑定关系
5. WHEN 管理员设置节点字段权限 THEN Form_Engine SHALL 根据当前节点控制字段的可见性和可编辑性
6. WHEN 表单模板被序列化存储 THEN Form_Engine SHALL 使用JSON格式编码表单结构
7. WHEN 从存储加载表单模板 THEN Form_Engine SHALL 解析JSON并还原完整的Form_Template对象
8. WHEN 用户填写表单 THEN Form_Engine SHALL 实时验证输入并提供错误提示
9. WHEN 管理员查询表单模板列表 THEN Form_Engine SHALL 返回所有可用模板供流程配置时选择

### Requirement 5: 审批人规则解析

**User Story:** As a 流程设计者, I want to 灵活配置审批人规则, so that 系统可以自动确定每个节点的审批人。

#### Acceptance Criteria

1. WHEN 审批人规则配置为指定用户 THEN Workflow_Engine SHALL 直接分配Task给配置的用户列表
2. WHEN 审批人规则配置为指定角色 THEN Workflow_Engine SHALL 查询拥有该角色的所有用户并分配Task
3. WHEN 审批人规则配置为部门负责人 THEN Workflow_Engine SHALL 根据发起人所属部门查找负责人并分配Task
4. WHEN 审批人规则配置为发起人上级 THEN Workflow_Engine SHALL 根据组织架构查找发起人的直接上级并分配Task
5. IF 审批人规则解析结果为空 THEN Workflow_Engine SHALL 记录异常并通知流程管理员处理

### Requirement 6: 消息通知集成

**User Story:** As a 审批参与者, I want to 及时收到审批相关通知, so that 可以快速响应待办事项。

#### Acceptance Criteria

1. WHEN 新Task创建 THEN Notification_Service SHALL 通过WebSocket向Assignee推送待办提醒
2. WHEN Process_Instance状态变更 THEN Notification_Service SHALL 通知发起人审批进度更新
3. WHEN 任务即将超时 THEN Notification_Service SHALL 提前向Assignee发送催办提醒
4. WHEN 审批完成 THEN Notification_Service SHALL 通知发起人最终审批结果

### Requirement 7: 审批记录与查询

**User Story:** As a 用户, I want to 查看审批历史和统计数据, so that 可以追溯审批过程和分析效率。

#### Acceptance Criteria

1. WHEN 用户查询我发起的申请 THEN Query_Service SHALL 返回该用户作为发起人的所有Process_Instance列表
2. WHEN 用户查询我的待办 THEN Query_Service SHALL 返回分配给该用户且状态为待处理的Task列表
3. WHEN 用户查询我的已办 THEN Query_Service SHALL 返回该用户已处理的Task历史记录
4. WHEN 用户查看流程详情 THEN Query_Service SHALL 返回Process_Instance的完整审批轨迹和表单数据
5. WHEN 管理员查询流程统计 THEN Query_Service SHALL 返回按流程类型、时间段分组的审批数量和平均耗时

### Requirement 8: 流程设计器前端

**User Story:** As a 系统管理员, I want to 通过可视化界面设计流程, so that 无需编写代码即可配置审批流程。

#### Acceptance Criteria

1. WHEN 管理员打开流程设计器 THEN Designer_UI SHALL 显示画布区域和节点工具栏
2. WHEN 管理员从工具栏拖拽节点到画布 THEN Designer_UI SHALL 在画布上创建对应类型的节点元素
3. WHEN 管理员连接两个节点 THEN Designer_UI SHALL 创建Transition并显示连接线
4. WHEN 管理员点击节点 THEN Designer_UI SHALL 显示该节点的属性配置面板
5. WHEN 管理员保存流程设计 THEN Designer_UI SHALL 将画布内容序列化为JSON并调用后端API保存
