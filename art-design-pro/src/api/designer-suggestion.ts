/***
 * Designer Suggestion API
 * 设计建议API接口
 * Requirements: 3.1, 3.2, 3.3, 3.4
 ***/
import request from '@/utils/http'

/*** Design Suggestion Interface ***/
export interface DesignSuggestion {
  id: number
  projectId: number
  content: string
  applicableScene: string
  costImpact: string
  category: string
  priority: number
  status: string
  detailInfo: Record<string, unknown> | null
  referenceImages: string[]
  createdAt: string
  updatedAt: string
}

/*** Generate Suggestion Request ***/
export interface GenerateSuggestionRequest {
  projectId: number
  category?: string
  count?: number
}

/*** Generate Suggestion Response (sync) ***/
export interface GenerateSuggestionResponse {
  suggestions: DesignSuggestion[]
  projectId: number
  count: number
}

/*** Generate Suggestion Async Response ***/
export interface GenerateSuggestionAsyncResponse {
  taskId: number
  projectId: number
  status: string
  message: string
}

/*** Suggestion Task Status Response ***/
export interface SuggestionTaskStatusResponse {
  taskId: number
  projectId: number
  status: string
  generatedCount: number
  errorMessage: string
  createdAt: string
}

/*** Suggestion List Request ***/
export interface SuggestionListRequest {
  projectId: number
  current: number
  size: number
  category?: string
  status?: string
}

/*** Suggestion List Response ***/
export interface SuggestionListResponse {
  records: DesignSuggestion[]
  current: number
  size: number
  total: number
}

/*** Update Status Request ***/
export interface UpdateStatusRequest {
  id: number
  status: 'pending' | 'adopted' | 'ignored'
  feedback?: string
}

/*** Suggestion Stats Response ***/
export interface SuggestionStatsResponse {
  totalCount: number
  adoptedCount: number
  ignoredCount: number
  pendingCount: number
  categoryStats: Record<string, number>
}

/*** Suggestion Detail Response ***/
export interface SuggestionDetailResponse extends DesignSuggestion {
  projectName: string
  projectInfo: {
    area: number
    budget: number
    style: string
  }
}

/***
 * Generate design suggestions (async)
 * 异步生成设计建议
 * Requirements: 3.1, 3.2
 ***/
export function generateSuggestions(data: GenerateSuggestionRequest) {
  return request.post<GenerateSuggestionAsyncResponse>({
    url: '/api/designer/suggestions/generate',
    data,
    _fullResponse: true
  })
}

/***
 * Get suggestion task status
 * 获取建议生成任务状态
 ***/
export function getSuggestionTaskStatus(taskId: number) {
  return request.get<SuggestionTaskStatusResponse>({
    url: `/api/designer/suggestions/task/${taskId}`,
    _fullResponse: true
  })
}

/***
 * Get suggestion list
 * 获取建议列表
 ***/
export function getSuggestionList(params: SuggestionListRequest) {
  return request.get<SuggestionListResponse>({
    url: '/api/designer/suggestions',
    params,
    _fullResponse: true
  })
}

/***
 * Get suggestion detail
 * 获取建议详情
 * Requirements: 3.3
 ***/
export function getSuggestionDetail(id: number) {
  return request.get<SuggestionDetailResponse>({
    url: `/api/designer/suggestions/${id}`,
    _fullResponse: true
  })
}

/***
 * Update suggestion status
 * 更新建议状态
 * Requirements: 3.4
 ***/
export function updateSuggestionStatus(data: UpdateStatusRequest) {
  return request.put({
    url: '/api/designer/suggestions/status',
    data,
    _fullResponse: true
  })
}

/***
 * Batch update suggestion status
 * 批量更新建议状态
 ***/
export function batchUpdateSuggestionStatus(data: {
  ids: number[]
  status: string
  feedback?: string
}) {
  return request.put({
    url: '/api/designer/suggestions/batch-status',
    data,
    _fullResponse: true
  })
}

/***
 * Delete suggestion
 * 删除建议
 ***/
export function deleteSuggestion(id: number) {
  return request.del({
    url: `/api/designer/suggestions/${id}`,
    _fullResponse: true
  })
}

/***
 * Batch delete suggestions
 * 批量删除建议
 ***/
export function batchDeleteSuggestions(ids: number[]) {
  return request.post({
    url: '/api/designer/suggestions/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

/***
 * Get suggestion stats
 * 获取建议统计
 ***/
export function getSuggestionStats(projectId?: number) {
  return request.get<SuggestionStatsResponse>({
    url: '/api/designer/suggestions/stats',
    params: { projectId },
    _fullResponse: true
  })
}
