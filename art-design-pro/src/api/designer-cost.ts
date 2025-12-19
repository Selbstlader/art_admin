/***
 * Designer Cost API
 * 设计师成本估算管理 API 接口
 * Requirements: 10.1, 10.2, 10.3, 10.4, 10.5
 ***/
import request from '@/utils/http'

/*** Cost API Types ***/

// 成本明细项 / Cost item
export interface CostItem {
  materialId: number
  materialName: string
  quantity: number
  unitPrice: number
  totalPrice: number
  category: string
  unit: string
}

// 成本估算响应 / Cost estimate response
export interface CostEstimateResponse {
  id: number
  projectId: number
  projectName: string
  materialCost: number
  laborCost: number
  equipmentCost: number
  managementCost: number
  totalCost: number
  budgetLimit: number
  budgetExceeded: boolean
  exceededAmount: number
  items: CostItem[]
  createdAt: string
  updatedAt: string
}

// 创建成本估算请求 / Create cost estimate request
export interface CreateCostEstimateRequest {
  projectId: number
  budgetLimit?: number
  items?: CostItem[]
  laborRate?: number
  equipRate?: number
  mgmtRate?: number
}

// 更新成本估算请求 / Update cost estimate request
export interface UpdateCostEstimateRequest {
  id: number
  budgetLimit?: number
  items?: CostItem[]
  laborRate?: number
  equipRate?: number
  mgmtRate?: number
}

// 计算成本请求 / Calculate cost request
export interface CalculateCostRequest {
  projectId: number
  area?: number
  items?: CostItem[]
  laborRate?: number
  equipRate?: number
  mgmtRate?: number
  budgetLimit?: number
}

// 计算成本响应 / Calculate cost response
export interface CostCalculateResponse {
  materialCost: number
  laborCost: number
  equipmentCost: number
  managementCost: number
  totalCost: number
  budgetLimit: number
  budgetExceeded: boolean
  exceededAmount: number
  budgetWarning: string
  items: CostItem[]
}

// 导出成本报告请求 / Export cost report request
export interface ExportCostReportRequest {
  projectId: number
  format: 'excel' | 'pdf'
}

// 导出成本报告响应 / Export cost report response
export interface CostReportExportResponse {
  fileUrl: string
  fileName: string
  fileSize: number
}

// 成本汇总响应 / Cost summary response
export interface CostSummaryResponse {
  totalProjects: number
  totalBudget: number
  totalCost: number
  avgCostPerProject: number
  overBudgetCount: number
}

/*** Cost API Functions ***/

/***
 * Create cost estimate
 * 创建成本估算
 ***/
export function createCostEstimate(data: CreateCostEstimateRequest) {
  return request.post<CostEstimateResponse>({
    url: '/api/designer/cost',
    data,
    _fullResponse: true
  })
}

/***
 * Update cost estimate
 * 更新成本估算
 ***/
export function updateCostEstimate(data: UpdateCostEstimateRequest) {
  return request.put<CostEstimateResponse>({
    url: '/api/designer/cost',
    data,
    _fullResponse: true
  })
}

/***
 * Delete cost estimate
 * 删除成本估算
 ***/
export function deleteCostEstimate(id: number) {
  return request.del({
    url: `/api/designer/cost/${id}`,
    _fullResponse: true
  })
}

/***
 * Get cost estimate by ID
 * 根据ID获取成本估算
 ***/
export function getCostEstimateById(id: number) {
  return request.get<CostEstimateResponse>({
    url: `/api/designer/cost/${id}`,
    _fullResponse: true
  })
}

/***
 * Get cost estimate by project ID
 * 根据项目ID获取成本估算
 ***/
export function getCostEstimateByProjectId(projectId: number) {
  return request.get<CostEstimateResponse>({
    url: '/api/designer/cost/project',
    params: { projectId },
    _fullResponse: true
  })
}

/***
 * Calculate cost (without saving)
 * 计算成本（不保存）
 ***/
export function calculateCost(data: CalculateCostRequest) {
  return request.post<CostCalculateResponse>({
    url: '/api/designer/cost/calculate',
    data,
    _fullResponse: true
  })
}

/***
 * Export cost report
 * 导出成本报告
 ***/
export function exportCostReport(data: ExportCostReportRequest) {
  return request.post<CostReportExportResponse>({
    url: '/api/designer/cost/export',
    data,
    _fullResponse: true
  })
}

/***
 * Get cost summary
 * 获取成本汇总
 ***/
export function getCostSummary() {
  return request.get<CostSummaryResponse>({
    url: '/api/designer/cost/summary',
    _fullResponse: true
  })
}
