# Implementation Plan

## Phase 1: 基础设施与数据模型

- [x] 1. 创建数据库模型和迁移文件
  - [x] 1.1 创建流程定义模型 (ProcessDefinition)
    - 创建 `internal/model/workflow.go`
    - 定义 ProcessDefinition, ProcessGraph, ProcessNode, ProcessEdge 结构
    - _Requirements: 1.1, 1.6_
  - [x] 1.2 创建流程实例和任务模型
    - 定义 ProcessInstance, Task, TaskLog 结构
    - _Requirements: 2.1, 3.1_
  - [x] 1.3 创建表单模板模型 (FormTemplate)
    - 定义 FormTemplate, FormSchema, FormField 结构
    - _Requirements: 4.1, 4.6_
  - [x] 1.4 创建数据库迁移文件
    - 创建 `migrations/xxx_create_workflow_tables.sql`
    - 包含所有工作流相关表
    - _Requirements: 1.1, 2.1, 4.1_

- [x] 2. 创建Repository层
  - [x] 2.1 实现 ProcessDefinitionRepository
    - 创建 `internal/repository/process_def_repo.go`
    - 实现 CRUD 操作
    - _Requirements: 1.1, 1.4, 1.5_
  - [x] 2.2 实现 ProcessInstanceRepository
    - 创建 `internal/repository/process_inst_repo.go`
    - 实现实例和任务查询
    - _Requirements: 2.1, 7.1, 7.2, 7.3_
  - [x] 2.3 实现 FormTemplateRepository
    - 创建 `internal/repository/form_template_repo.go`
    - 实现模板 CRUD
    - _Requirements: 4.1, 4.9_
  - [x] 2.4 实现 TaskRepository
    - 创建 `internal/repository/task_repo.go`
    - 实现任务查询和更新
    - _Requirements: 3.1, 3.5, 3.6_

## Phase 2: 核心服务层

- [x] 3. 实现JSON序列化/反序列化
  - [x] 3.1 实现流程定义序列化服务
    - 创建 `internal/service/workflow/serializer.go`
    - 实现 ProcessGraph 的 JSON 序列化和反序列化
    - _Requirements: 1.6, 1.7_
  - [x] 3.2 编写属性测试：流程定义序列化往返
    - **Property 1: 流程定义序列化往返一致性**
    - **Validates: Requirements 1.6, 1.7**
  - [x] 3.3 实现表单模板序列化服务
    - 实现 FormSchema 的 JSON 序列化和反序列化
    - _Requirements: 4.6, 4.7_
  - [x] 3.4 编写属性测试：表单模板序列化往返
    - **Property 2: 表单模板序列化往返一致性**
    - **Validates: Requirements 4.6, 4.7**

- [x] 4. 实现流程定义服务
  - [x] 4.1 实现流程结构验证
    - 创建 `internal/service/workflow/validator.go`
    - 验证开始/结束节点、节点可达性
    - _Requirements: 1.1_
  - [x] 4.2 编写属性测试：流程结构完整性验证
    - **Property 3: 流程结构完整性验证**
    - **Validates: Requirements 1.1**
  - [x] 4.3 实现流程定义CRUD服务
    - 创建 `internal/service/workflow/process_def_service.go`
    - 实现创建、更新、删除、查询
    - _Requirements: 1.1, 1.4_
  - [x] 4.4 实现流程发布和版本控制
    - 实现 Publish 方法，版本号递增
    - _Requirements: 1.4, 1.5_
  - [x] 4.5 编写属性测试：流程发布版本递增
    - **Property 4: 流程发布版本递增**
    - **Validates: Requirements 1.4, 1.5**

- [x] 5. Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.

- [x] 6. 实现表单引擎服务
  - [x] 6.1 实现表单模板CRUD服务
    - 创建 `internal/service/workflow/form_engine_service.go`
    - 实现模板管理
    - _Requirements: 4.1, 4.4, 4.9_
  - [x] 6.2 实现表单数据验证
    - 实现必填、格式、范围校验
    - _Requirements: 4.3_
  - [x] 6.3 编写属性测试：表单数据验证正确性
    - **Property 6: 表单数据验证正确性**
    - **Validates: Requirements 2.6, 4.3**
  - [x] 6.4 实现节点字段权限控制
    - 根据节点配置返回字段可见性/可编辑性
    - _Requirements: 4.5_

- [x] 7. 实现审批人解析器
  - [x] 7.1 实现审批人规则解析服务
    - 创建 `internal/service/workflow/assignee_resolver.go`
    - 实现四种规则类型解析
    - _Requirements: 5.1, 5.2, 5.3, 5.4_
  - [x] 7.2 编写属性测试：审批人规则解析正确性
    - **Property 11: 审批人规则解析正确性**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.4**
  - [x] 7.3 实现空结果错误处理
    - 解析结果为空时返回错误
    - _Requirements: 5.5_
  - [x] 7.4 编写属性测试：审批人规则空结果处理
    - **Property 12: 审批人规则空结果处理**
    - **Validates: Requirements 5.5**

- [x] 8. 实现工作流引擎核心
  - [x] 8.1 实现流程启动服务
    - 创建 `internal/service/workflow/workflow_engine.go`
    - 实现 StartProcess 方法
    - _Requirements: 2.1, 2.2_
  - [x] 8.2 编写属性测试：流程实例初始状态正确性
    - **Property 5: 流程实例初始状态正确性**
    - **Validates: Requirements 2.1, 2.2**
  - [x] 8.3 实现任务完成和流转
    - 实现 CompleteTask 方法（通过/拒绝）
    - _Requirements: 3.1, 3.2_
  - [x] 8.4 编写属性测试：任务通过后流转正确性
    - **Property 8: 任务通过后流转正确性**
    - **Validates: Requirements 3.1**
  - [x] 8.5 实现会签/或签逻辑
    - 根据审批模式决定流转时机
    - _Requirements: 3.3, 3.4_
  - [x] 8.6 编写属性测试：审批模式流转正确性
    - **Property 7: 审批模式流转正确性**
    - **Validates: Requirements 3.3, 3.4**
  - [x] 8.7 实现委托和转办
    - 实现 DelegateTask, TransferTask 方法
    - _Requirements: 3.5, 3.6_
  - [x] 8.8 编写属性测试：委托/转办记录完整性
    - **Property 10: 委托/转办记录完整性**
    - **Validates: Requirements 3.5, 3.6**
  - [x] 8.9 实现流程撤回
    - 实现 WithdrawProcess 方法
    - _Requirements: 2.4_
  - [x] 8.10 实现流程状态终态处理
    - 处理完成、拒绝、撤回状态
    - _Requirements: 2.3, 2.4, 2.5_
  - [x] 8.11 编写属性测试：流程状态终态一致性
    - **Property 9: 流程状态终态一致性**
    - **Validates: Requirements 2.3, 2.4, 2.5**

- [x] 9. Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.

## Phase 3: 查询服务与通知

- [x] 10. 实现查询服务
  - [x] 10.1 实现待办/已办/我发起的查询
    - 创建 `internal/service/workflow/query_service.go`
    - 实现 ListMyTodo, ListMyDone, ListMyInitiated
    - _Requirements: 7.1, 7.2, 7.3_
  - [x] 10.2 编写属性测试：查询结果过滤正确性
    - **Property 13: 查询结果过滤正确性**
    - **Validates: Requirements 7.1, 7.2, 7.3**
  - [x] 10.3 实现审批轨迹查询
    - 实现 GetApprovalTrail 方法
    - _Requirements: 7.4_
  - [x] 10.4 实现流程统计查询
    - 实现 GetStatistics 方法
    - _Requirements: 7.5_

- [x] 11. 实现通知服务
  - [x] 11.1 创建通知服务
    - 创建 `internal/service/workflow/notification_service.go`
    - 复用现有 WebSocket Hub
    - _Requirements: 6.1, 6.2_
  - [x] 11.2 实现待办通知推送
    - 新任务创建时推送通知
    - _Requirements: 6.1_
  - [x] 11.3 实现流程状态变更通知
    - 状态变更时通知发起人
    - _Requirements: 6.2, 6.4_

## Phase 4: API层

- [x] 12. 创建DTO定义
  - [x] 12.1 创建请求/响应DTO
    - 创建 `internal/dto/request/workflow_request.go`
    - 创建 `internal/dto/response/workflow_response.go`
    - _Requirements: 1.1, 2.1, 4.1_

- [x] 13. 实现后端API
  - [x] 13.1 实现流程定义API
    - 创建 `internal/api/v1/process_def_api.go`
    - 实现 CRUD 和发布接口
    - _Requirements: 1.1, 1.4, 1.5_
  - [x] 13.2 实现流程实例API
    - 创建 `internal/api/v1/process_inst_api.go`
    - 实现发起、审批、撤回接口
    - _Requirements: 2.1, 3.1, 3.5, 3.6_
  - [x] 13.3 实现表单模板API
    - 创建 `internal/api/v1/form_template_api.go`
    - 实现模板 CRUD 接口
    - _Requirements: 4.1, 4.9_
  - [x] 13.4 实现查询API
    - 创建 `internal/api/v1/workflow_query_api.go`
    - 实现待办、已办、统计接口
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_
  - [x] 13.5 注册路由
    - 更新 `internal/router/` 添加工作流路由
    - _Requirements: 1.1, 2.1, 4.1_

- [x] 14. Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.

## Phase 5: 前端实现

- [x] 15. 安装前端依赖
  - [x] 15.1 安装流程设计器依赖
    - 安装 @vue-flow/core, @vue-flow/background, @vue-flow/controls, @vue-flow/minimap
    - _Requirements: 8.1_
  - [x] 15.2 安装表单设计器依赖
    - 安装 vuedraggable@4.x
    - _Requirements: 4.1_

- [x] 16. 实现前端API层
  - [x] 16.1 创建工作流API模块
    - 创建 `src/api/workflow.ts`
    - 封装所有工作流相关接口
    - _Requirements: 1.1, 2.1, 4.1_

- [x] 17. 实现表单设计器
  - [x] 17.1 创建表单设计器组件
    - 创建 `src/components/workflow/FormDesigner.vue`
    - 实现三栏布局：组件面板、画布、属性面板
    - _Requirements: 4.1, 4.2_
  - [x] 17.2 实现组件拖拽功能
    - 实现从组件面板拖拽到画布
    - 实现画布内字段排序
    - _Requirements: 4.1_
  - [x] 17.3 实现字段属性配置
    - 实现属性面板字段配置
    - 支持验证规则配置
    - _Requirements: 4.2, 4.3_

- [x] 18. 实现流程设计器
  - [x] 18.1 创建流程设计器组件
    - 创建 `src/components/workflow/ProcessDesigner.vue`
    - 集成 Vue Flow
    - _Requirements: 8.1, 8.2_
  - [x] 18.2 实现自定义节点组件
    - 创建开始、审批、条件、结束节点组件
    - _Requirements: 8.2, 8.3_
  - [x] 18.3 实现节点属性配置面板
    - 实现审批人规则、审批模式等配置
    - _Requirements: 8.4_
  - [x] 18.4 实现流程保存和发布
    - 序列化画布内容为JSON
    - 调用后端API保存
    - _Requirements: 8.5_

- [x] 19. 实现表单模板管理页面
  - [x] 19.1 创建表单模板列表页
    - 创建 `src/views/workflow/form-template/index.vue`
    - 实现列表展示和搜索
    - _Requirements: 4.9_
  - [x] 19.2 创建表单模板编辑页
    - 创建 `src/views/workflow/form-template/edit.vue`
    - 集成表单设计器
    - _Requirements: 4.1_

- [x] 20. 实现流程定义管理页面
  - [x] 20.1 创建流程定义列表页
    - 创建 `src/views/workflow/process-def/index.vue`
    - 实现列表、搜索、发布操作
    - _Requirements: 1.1, 1.4_
  - [x] 20.2 创建流程定义编辑页
    - 创建 `src/views/workflow/process-def/edit.vue`
    - 集成流程设计器
    - _Requirements: 8.1, 8.5_

- [x] 21. 实现审批工作台
  - [x] 21.1 创建待办列表页
    - 创建 `src/views/workflow/workbench/todo.vue`
    - 实现待办任务列表
    - _Requirements: 7.2_
  - [x] 21.2 创建已办列表页
    - 创建 `src/views/workflow/workbench/done.vue`
    - 实现已办任务列表
    - _Requirements: 7.3_
  - [x] 21.3 创建我发起的列表页
    - 创建 `src/views/workflow/workbench/initiated.vue`
    - 实现我发起的流程列表
    - _Requirements: 7.1_
  - [x] 21.4 创建审批详情页
    - 创建 `src/views/workflow/workbench/detail.vue`
    - 显示表单数据和审批轨迹
    - 实现审批操作（通过/拒绝/委托/转办）
    - _Requirements: 3.1, 3.5, 3.6, 7.4_

- [x] 22. 实现发起申请页面
  - [x] 22.1 创建发起申请页
    - 创建 `src/views/workflow/apply/index.vue`
    - 显示可用流程列表
    - _Requirements: 2.1_
  - [x] 22.2 创建申请表单页
    - 创建 `src/views/workflow/apply/form.vue`
    - 动态渲染表单并提交
    - _Requirements: 2.1, 4.8_

- [ ] 23. 配置路由和菜单
  - [x] 23.1 添加工作流路由配置
    - 更新 `src/router/` 添加工作流相关路由
    - _Requirements: 8.1_
  - [-] 23.2 添加菜单配置
    - 添加工作流管理、审批工作台菜单
    - _Requirements: 8.1_

- [ ] 24. Final Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.
