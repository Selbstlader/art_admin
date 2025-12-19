/***
 * Designer Project API
 * 设计师项目管理 API 接口
 * Requirements: 4.1, 4.2, 4.3, 4.4
 ***/
import request from '@/utils/http'

// 项目列表请求参数
export interface ProjectListParams {
  current: number
  size: number
  name?: string
  status?: string
  style?: string
  minBudget?: number
  maxBudget?: number
  startDate?: string
  endDate?: string
  keyword?: string
}

// 项目响应数据
export interface ProjectResponse {
  id: number
  name: string
  description: string
  area: number
  budget: number
  style: string
  status: string
  userId: number
  createdAt: string
  updatedAt: string
}

// 项目详情响应数据
export interface ProjectDetailResponse extends ProjectResponse {
  documentCount: number
  cadFileCount: number
  versionCount?: number
  documents: DocumentBrief[]
  cadFiles: CadFileBrief[]
  costEstimate: CostEstimateBrief | null
  compareResults: CompareResultBrief[]
}

// 文档简要信息
export interface DocumentBrief {
  id: number
  fileName: string
  fileType: string
  fileSize: number
  analysisStatus: string
  createdAt: string
}

// CAD文件简要信息
export interface CadFileBrief {
  id: number
  fileName: string
  fileFormat: string
  parseStatus: string
  layerCount: number
  has3d: boolean
  createdAt: string
}

// 成本估算简要信息
export interface CostEstimateBrief {
  id: number
  materialCost: number
  laborCost: number
  equipmentCost: number
  managementCost: number
  totalCost: number
  budgetLimit: number
  updatedAt: string
}

// 比对结果简要信息
export interface CompareResultBrief {
  id: number
  overallScore: number
  createdAt: string
}

// 创建项目请求参数
export interface CreateProjectParams {
  name: string
  description?: string
  area?: number
  budget?: number
  style?: string
  status?: string
}

// 更新项目请求参数
export interface UpdateProjectParams extends CreateProjectParams {
  id: number
}

// 项目列表响应
export interface ProjectListResponse {
  records: ProjectResponse[]
  current: number
  size: number
  total: number
}

// 项目统计响应
export interface ProjectStatsResponse {
  total: number
  draft: number
  inProgress: number
  completed: number
  archived: number
}

/***
 * Get designer project list
 * 获取设计师项目列表
 ***/
export function getDesignerProjects(params: ProjectListParams) {
  return request.get<ProjectListResponse>({
    url: '/api/designer/projects',
    params,
    _fullResponse: true
  })
}

/***
 * Get designer project detail
 * 获取设计师项目详情
 ***/
export function getDesignerProject(id: number) {
  return request.get<ProjectDetailResponse>({
    url: `/api/designer/projects/${id}`,
    _fullResponse: true
  })
}

/***
 * Create designer project
 * 创建设计师项目
 ***/
export function createDesignerProject(data: CreateProjectParams) {
  return request.post<ProjectResponse>({
    url: '/api/designer/projects',
    data,
    _fullResponse: true
  })
}

/***
 * Update designer project
 * 更新设计师项目
 ***/
export function updateDesignerProject(id: number, data: UpdateProjectParams) {
  return request.put<ProjectResponse>({
    url: `/api/designer/projects/${id}`,
    data,
    _fullResponse: true
  })
}

/***
 * Delete designer project
 * 删除设计师项目
 ***/
export function deleteDesignerProject(id: number) {
  return request.del({
    url: `/api/designer/projects/${id}`,
    _fullResponse: true
  })
}

/***
 * Batch delete designer projects
 * 批量删除设计师项目
 ***/
export function batchDeleteDesignerProjects(ids: number[]) {
  return request.post({
    url: '/api/designer/projects/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

/***
 * Search designer projects
 * 搜索设计师项目
 ***/
export function searchDesignerProjects(keyword: string, current = 1, size = 10) {
  return request.get<ProjectListResponse>({
    url: '/api/designer/projects/search',
    params: { keyword, current, size },
    _fullResponse: true
  })
}

/***
 * Get designer project stats
 * 获取设计师项目统计
 ***/
export function getDesignerProjectStats() {
  return request.get<ProjectStatsResponse>({
    url: '/api/designer/projects/stats',
    _fullResponse: true
  })
}
