/***
 * Designer Compliance API
 * 设计规范合规检查API
 * Requirements: 11.1, 11.2, 11.3, 11.4, 11.5
 ***/
import request from '@/utils/http'

/*** Design Standard Interfaces ***/
export interface DesignStandard {
  id: number
  code: string
  name: string
  category: string
  content: string
  applicableTypes: string[]
  version: string
  effectiveDate: string
  source: string
  interpretation: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface CreateDesignStandardRequest {
  code: string
  name: string
  category: string
  content: string
  applicableTypes?: string[]
  version?: string
  effectiveDate?: string
  source?: string
  interpretation?: string
}

export interface UpdateDesignStandardRequest extends CreateDesignStandardRequest {
  id: number
  status?: string
}

export interface DesignStandardListRequest {
  current: number
  size: number
  category?: string
  status?: string
  keyword?: string
}

export interface DesignStandardListResponse {
  records: DesignStandard[]
  current: number
  size: number
  total: number
}

export interface DesignStandardStats {
  totalCount: number
  activeCount: number
  categoryStats: { category: string; count: number }[]
}

/*** Compliance Check Interfaces ***/
export interface PassedItem {
  standardId: number
  standardCode: string
  standardName: string
  category: string
  description: string
}

export interface FailedItem {
  standardId: number
  standardCode: string
  standardName: string
  category: string
  violationContent: string
  originalRequirement: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  location: string
  suggestion: string
}

export interface ComplianceSuggestion {
  content: string
  priority: 'low' | 'medium' | 'high'
  category: string
  costImpact: string
  reference: string
}

export interface ComplianceCheckResult {
  id: number
  projectId: number
  projectName: string
  checkType: string
  passedItems: PassedItem[]
  failedItems: FailedItem[]
  suggestions: ComplianceSuggestion[]
  overallScore: number
  checkStatus: string
  errorMessage: string
  checkedStandards: number[]
  userId: number
  createdAt: string
  updatedAt: string
}

export interface ComplianceCheckRequest {
  projectId: number
  categories?: string[]
  checkType?: 'full' | 'partial'
}

export interface ComplianceCheckListRequest {
  projectId: number
  current: number
  size: number
  checkStatus?: string
}

export interface ComplianceCheckListResponse {
  records: ComplianceCheckResult[]
  current: number
  size: number
  total: number
}

export interface ComplianceCheckSummary {
  totalChecks: number
  passedCount: number
  failedCount: number
  averageScore: number
  latestCheckAt: string
}

/*** Design Standard APIs ***/

// 创建设计规范 / Create design standard
export function createDesignStandard(data: CreateDesignStandardRequest) {
  return request.post<Http.BaseResponse<DesignStandard>>({
    url: '/api/designer/compliance/standards',
    data,
    _fullResponse: true
  })
}

// 更新设计规范 / Update design standard
export function updateDesignStandard(data: UpdateDesignStandardRequest) {
  return request.put<Http.BaseResponse<DesignStandard>>({
    url: '/api/designer/compliance/standards',
    data,
    _fullResponse: true
  })
}

// 获取设计规范列表 / Get design standard list
export function getDesignStandardList(params: DesignStandardListRequest) {
  return request.get<Http.BaseResponse<DesignStandardListResponse>>({
    url: '/api/designer/compliance/standards',
    params,
    _fullResponse: true
  })
}

// 获取设计规范详情 / Get design standard detail
export function getDesignStandardDetail(id: number) {
  return request.get<Http.BaseResponse<DesignStandard>>({
    url: `/api/designer/compliance/standards/${id}`,
    _fullResponse: true
  })
}

// 删除设计规范 / Delete design standard
export function deleteDesignStandard(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/compliance/standards/${id}`,
    _fullResponse: true
  })
}

// 批量删除设计规范 / Batch delete design standards
export function batchDeleteDesignStandards(ids: number[]) {
  return request.post<Http.BaseResponse<null>>({
    url: '/api/designer/compliance/standards/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

// 获取设计规范统计 / Get design standard stats
export function getDesignStandardStats() {
  return request.get<Http.BaseResponse<DesignStandardStats>>({
    url: '/api/designer/compliance/standards/stats',
    _fullResponse: true
  })
}

/*** Compliance Check APIs ***/

// 执行合规检查 / Perform compliance check
export function performComplianceCheck(data: ComplianceCheckRequest) {
  return request.post<Http.BaseResponse<ComplianceCheckResult>>({
    url: '/api/designer/compliance/check',
    data,
    _fullResponse: true
  })
}

// 获取合规检查结果列表 / Get compliance check result list
export function getComplianceCheckList(params: ComplianceCheckListRequest) {
  return request.get<Http.BaseResponse<ComplianceCheckListResponse>>({
    url: '/api/designer/compliance/results',
    params,
    _fullResponse: true
  })
}

// 获取合规检查结果详情 / Get compliance check result detail
export function getComplianceCheckDetail(id: number) {
  return request.get<Http.BaseResponse<ComplianceCheckResult>>({
    url: `/api/designer/compliance/results/${id}`,
    _fullResponse: true
  })
}

// 获取项目最新合规检查结果 / Get latest compliance check result
export function getLatestComplianceCheck(projectId: number) {
  return request.get<Http.BaseResponse<ComplianceCheckResult>>({
    url: '/api/designer/compliance/latest',
    params: { projectId },
    _fullResponse: true
  })
}

// 获取项目合规检查摘要 / Get compliance check summary
export function getComplianceCheckSummary(projectId: number) {
  return request.get<Http.BaseResponse<ComplianceCheckSummary>>({
    url: '/api/designer/compliance/summary',
    params: { projectId },
    _fullResponse: true
  })
}

// 删除合规检查结果 / Delete compliance check result
export function deleteComplianceCheck(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/compliance/results/${id}`,
    _fullResponse: true
  })
}

// 批量删除合规检查结果 / Batch delete compliance check results
export function batchDeleteComplianceChecks(ids: number[]) {
  return request.post<Http.BaseResponse<null>>({
    url: '/api/designer/compliance/results/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

/*** Category Constants ***/
export const STANDARD_CATEGORIES = [
  { value: 'fire', label: '消防规范' },
  { value: 'accessibility', label: '无障碍规范' },
  { value: 'environmental', label: '环保规范' },
  { value: 'safety', label: '安全规范' },
  { value: 'other', label: '其他规范' }
]

export const SEVERITY_OPTIONS = [
  { value: 'low', label: '低', type: 'info' },
  { value: 'medium', label: '中', type: 'warning' },
  { value: 'high', label: '高', type: 'danger' },
  { value: 'critical', label: '严重', type: 'danger' }
]

export const PRIORITY_OPTIONS = [
  { value: 'low', label: '低', type: 'info' },
  { value: 'medium', label: '中', type: 'warning' },
  { value: 'high', label: '高', type: 'danger' }
]
