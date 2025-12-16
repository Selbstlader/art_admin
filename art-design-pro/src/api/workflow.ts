import request from '@/utils/http'

// ============================================================================
// 类型定义
// ============================================================================

// 流程定义相关类型
export interface ProcessDefListParams {
  page?: number
  pageSize?: number
  name?: string
  code?: string
  category?: string
  status?: string
}

export interface ProcessGraph {
  nodes: ProcessNode[]
  edges: ProcessEdge[]
  globalProps?: GlobalProperties
}

export interface ProcessNode {
  id: string
  type: string // start | end | approval | condition | parallel
  name: string
  position: Position
  properties?: NodeProperties
}

export interface Position {
  x: number
  y: number
}

export interface NodeProperties {
  assigneeRule?: AssigneeRule
  approvalMode?: string // or_sign | and_sign
  timeoutHours?: number
  rejectAction?: string // terminate | return_prev
  fieldPermissions?: Record<string, string>
}

export interface AssigneeRule {
  type: string // user | role | dept_leader | initiator_leader
  values: number[]
}

export interface ProcessEdge {
  id: string
  source: string
  target: string
  condition?: string
}

export interface GlobalProperties {
  allowWithdraw: boolean
}

export interface CreateProcessDefParams {
  name: string
  code: string
  description?: string
  category?: string
  formTemplateId?: number
  graph: ProcessGraph
}

export interface UpdateProcessDefParams {
  name: string
  description?: string
  category?: string
  formTemplateId?: number
  graph: ProcessGraph
}

export interface ProcessDefResponse {
  id: number
  name: string
  code: string
  description: string
  category: string
  formTemplateId?: number
  version: number
  status: string
  createdBy: number
  createdByName?: string
  createdAt: string
  updatedAt: string
}

export interface ProcessDefDetailResponse extends ProcessDefResponse {
  graph?: ProcessGraph
  formTemplateName?: string
}

// 流程实例相关类型
export interface StartProcessParams {
  processDefId: number
  title: string
  formData?: Record<string, any>
}

export interface CompleteTaskParams {
  taskId: number
  action: 'approve' | 'reject'
  comment?: string
}

export interface DelegateTaskParams {
  taskId: number
  toUserId: number
  reason?: string
}

export interface TransferTaskParams {
  taskId: number
  toUserId: number
  reason?: string
}

export interface ProcessInstResponse {
  id: number
  processDefId: number
  processDefVersion: number
  processName?: string
  processCode?: string
  title: string
  initiatorId: number
  initiatorName: string
  status: string
  currentNodeId: string
  currentNodeName?: string
  startedAt: string
  completedAt?: string
}

export interface ProcessInstDetailResponse extends ProcessInstResponse {
  formData?: Record<string, any>
  approvalTrail: ApprovalRecord[]
  currentTasks: WfTaskResponse[]
}

export interface ProcessInstSummaryResponse {
  id: number
  processName: string
  processCode?: string
  title: string
  initiatorId: number
  initiatorName: string
  status: string
  startedAt: string
  completedAt?: string
}

export interface WfTaskResponse {
  id: number
  processInstId: number
  processTitle?: string
  processName?: string
  processCode?: string
  nodeId: string
  nodeName: string
  assigneeId: number
  assigneeName: string
  status: string
  comment: string
  delegatedFrom?: number
  transferredFrom?: number
  dueAt?: string
  createdAt: string
  completedAt?: string
}

export interface WfTaskSummaryResponse {
  id: number
  processInstId: number
  processTitle: string
  processName: string
  processCode?: string
  nodeName: string
  initiatorId?: number
  initiatorName?: string
  status: string
  dueAt?: string
  createdAt: string
  completedAt?: string
}

export interface ApprovalRecord {
  taskId: number
  nodeId?: string
  nodeName?: string
  operatorId: number
  operatorName: string
  action: string
  comment: string
  createdAt: string
}

// 表单模板相关类型
export interface FormTemplateListParams {
  page?: number
  pageSize?: number
  name?: string
  code?: string
  status?: string
}

export interface FormSchema {
  fields: FormField[]
}

export interface FormField {
  key: string
  label: string
  type: string // text | number | date | select | file | textarea
  required: boolean
  placeholder?: string
  defaultValue?: any
  options?: SelectOption[]
  validation?: FieldValidation
}

export interface SelectOption {
  label: string
  value: any
}

export interface FieldValidation {
  minLength?: number
  maxLength?: number
  min?: number
  max?: number
  pattern?: string
}

export interface CreateFormTemplateParams {
  name: string
  code: string
  description?: string
  schema: FormSchema
}

export interface UpdateFormTemplateParams {
  name: string
  description?: string
  schema: FormSchema
}

export interface FormTemplateResponse {
  id: number
  name: string
  code: string
  description: string
  status: string
  createdBy: number
  createdByName?: string
  createdAt: string
  updatedAt: string
}

export interface FormTemplateDetailResponse extends FormTemplateResponse {
  schema?: FormSchema
}

// 查询相关类型
export interface QueryListParams {
  page?: number
  pageSize?: number
  title?: string
  processCode?: string
  status?: string
  startDate?: string
  endDate?: string
}

export interface StatisticsParams {
  processCode?: string
  category?: string
  startDate?: string
  endDate?: string
}

export interface ProcessStatisticsResponse {
  totalInstances: number
  runningInstances: number
  completedInstances: number
  rejectedInstances: number
  withdrawnInstances: number
  avgProcessingTime: number
  byCategory?: CategoryStatistics[]
  byProcess?: ProcessTypeStatistics[]
}

export interface CategoryStatistics {
  category: string
  count: number
}

export interface ProcessTypeStatistics {
  processCode?: string
  processName: string
  totalCount: number
  completedCount: number
  avgProcessingTime: number
}

// ============================================================================
// 流程定义 API
// ============================================================================

export const processDefApi = {
  /**
   * 获取流程定义列表
   */
  getList: (params: ProcessDefListParams) => {
    return request.get<{ list: ProcessDefResponse[]; total: number }>({
      url: '/api/workflow/process-def/list',
      params
    })
  },

  /**
   * 获取流程定义详情
   */
  getDetail: (id: number) => {
    return request.get<ProcessDefDetailResponse>({
      url: `/api/workflow/process-def/${id}`
    })
  },

  /**
   * 创建流程定义
   */
  create: (params: CreateProcessDefParams) => {
    return request.post<ProcessDefResponse>({
      url: '/api/workflow/process-def',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 更新流程定义
   */
  update: (id: number, params: UpdateProcessDefParams) => {
    return request.put<ProcessDefResponse>({
      url: `/api/workflow/process-def/${id}`,
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 删除流程定义
   */
  delete: (id: number) => {
    return request.del({
      url: `/api/workflow/process-def/${id}`,
      showSuccessMessage: true
    })
  },

  /**
   * 发布流程定义
   */
  publish: (id: number) => {
    return request.post<ProcessDefResponse>({
      url: `/api/workflow/process-def/${id}/publish`,
      showSuccessMessage: true
    })
  }
}

// ============================================================================
// 流程实例 API
// ============================================================================

export const processInstApi = {
  /**
   * 启动流程
   */
  start: (params: StartProcessParams) => {
    return request.post<ProcessInstResponse>({
      url: '/api/workflow/process-inst/start',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 完成任务（通过/拒绝）
   */
  completeTask: (params: CompleteTaskParams) => {
    return request.post({
      url: '/api/workflow/process-inst/complete-task',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 委托任务
   */
  delegateTask: (params: DelegateTaskParams) => {
    return request.post({
      url: '/api/workflow/process-inst/delegate-task',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 转办任务
   */
  transferTask: (params: TransferTaskParams) => {
    return request.post({
      url: '/api/workflow/process-inst/transfer-task',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 撤回流程
   */
  withdraw: (id: number) => {
    return request.post({
      url: `/api/workflow/process-inst/${id}/withdraw`,
      showSuccessMessage: true
    })
  },

  /**
   * 获取流程实例详情
   */
  getDetail: (id: number) => {
    return request.get<ProcessInstDetailResponse>({
      url: `/api/workflow/process-inst/${id}`
    })
  }
}

// ============================================================================
// 表单模板 API
// ============================================================================

export const formTemplateApi = {
  /**
   * 获取表单模板列表
   */
  getList: (params: FormTemplateListParams) => {
    return request.get<{ records: FormTemplateResponse[]; total: number }>({
      url: '/api/workflow/form-template/list',
      params
    })
  },

  /**
   * 获取表单模板详情
   */
  getDetail: (id: number) => {
    return request.get<FormTemplateDetailResponse>({
      url: `/api/workflow/form-template/${id}`
    })
  },

  /**
   * 创建表单模板
   */
  create: (params: CreateFormTemplateParams) => {
    return request.post<FormTemplateResponse>({
      url: '/api/workflow/form-template',
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 更新表单模板
   */
  update: (id: number, params: UpdateFormTemplateParams) => {
    return request.put<FormTemplateResponse>({
      url: `/api/workflow/form-template/${id}`,
      data: params,
      showSuccessMessage: true
    })
  },

  /**
   * 删除表单模板
   */
  delete: (id: number) => {
    return request.del({
      url: `/api/workflow/form-template/${id}`,
      showSuccessMessage: true
    })
  }
}

// ============================================================================
// 工作流查询 API
// ============================================================================

export const workflowQueryApi = {
  /**
   * 获取我发起的流程
   */
  getMyInitiated: (params: QueryListParams) => {
    return request.get<{ list: ProcessInstSummaryResponse[]; total: number }>({
      url: '/api/workflow/query/my-initiated',
      params
    })
  },

  /**
   * 获取我的待办
   */
  getMyTodo: (params: QueryListParams) => {
    return request.get<{ list: WfTaskSummaryResponse[]; total: number }>({
      url: '/api/workflow/query/my-todo',
      params
    })
  },

  /**
   * 获取我的已办
   */
  getMyDone: (params: QueryListParams) => {
    return request.get<{ list: WfTaskSummaryResponse[]; total: number }>({
      url: '/api/workflow/query/my-done',
      params
    })
  },

  /**
   * 获取审批轨迹
   */
  getApprovalTrail: (instanceId: number) => {
    return request.get<ApprovalRecord[]>({
      url: `/api/workflow/query/approval-trail/${instanceId}`
    })
  },

  /**
   * 获取流程统计
   */
  getStatistics: (params?: StatisticsParams) => {
    return request.get<ProcessStatisticsResponse>({
      url: '/api/workflow/query/statistics',
      params
    })
  }
}

// ============================================================================
// 导出统一的工作流 API
// ============================================================================

export const workflowApi = {
  processDef: processDefApi,
  processInst: processInstApi,
  formTemplate: formTemplateApi,
  query: workflowQueryApi
}

export default workflowApi
