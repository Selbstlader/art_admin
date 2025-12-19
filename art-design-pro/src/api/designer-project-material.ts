/***
 * Designer Project Material API
 * 项目材料清单管理 API 接口
 * 关联项目与材料，支持成本估算
 ***/
import request from '@/utils/http'

/*** 项目材料清单项 / Project material item ***/
export interface ProjectMaterialItem {
  id?: number
  projectId: number
  materialId?: number // 关联材料库，可为空（自定义材料）
  name: string
  category: string
  specification: string
  unit: string
  unitPrice: number
  quantity: number
  totalPrice: number
  brand?: string
  supplier?: string
  remark?: string
  createdAt?: string
  updatedAt?: string
}

/*** 项目材料清单响应 / Project material list response ***/
export interface ProjectMaterialListResponse {
  projectId: number
  projectName: string
  items: ProjectMaterialItem[]
  materialCost: number
  itemCount: number
}

/*** 保存项目材料清单请求 / Save project material request ***/
export interface SaveProjectMaterialRequest {
  projectId: number
  items: Omit<ProjectMaterialItem, 'id' | 'projectId' | 'createdAt' | 'updatedAt'>[]
}

/*** 自定义费用项 / Custom cost item ***/
export interface CustomCostItem {
  name: string
  amount: number // 支持负数进行扣减 / Supports negative for deduction
}

/*** 项目成本汇总 / Project cost summary ***/
export interface ProjectCostSummary {
  projectId: number
  materialCost: number
  laborCost: number
  equipmentCost: number
  managementCost: number
  customCosts: CustomCostItem[] // 自定义费用项 / Custom cost items
  customCostTotal: number // 自定义费用小计 / Custom cost subtotal
  totalCost: number
  budgetLimit: number
  budgetExceeded: boolean
  exceededAmount: number
}

/*** 保存项目成本请求 / Save project cost request ***/
export interface SaveProjectCostRequest {
  projectId: number
  budgetLimit: number
  laborCost: number
  equipmentCost: number
  managementCost: number
  customCosts?: CustomCostItem[] // 自定义费用项 / Custom cost items
}

/***
 * Get project material list
 * 获取项目材料清单
 ***/
export function getProjectMaterials(projectId: number) {
  return request.get<ProjectMaterialListResponse>({
    url: `/api/designer/projects/${projectId}/materials`,
    _fullResponse: true
  })
}

/***
 * Save project material list
 * 保存项目材料清单
 ***/
export function saveProjectMaterials(data: SaveProjectMaterialRequest) {
  return request.post<ProjectMaterialListResponse>({
    url: `/api/designer/projects/${data.projectId}/materials`,
    data,
    _fullResponse: true
  })
}

/***
 * Add single material to project
 * 添加单个材料到项目
 ***/
export function addProjectMaterial(projectId: number, item: Omit<ProjectMaterialItem, 'id' | 'projectId' | 'createdAt' | 'updatedAt'>) {
  return request.post<ProjectMaterialItem>({
    url: `/api/designer/projects/${projectId}/materials/add`,
    data: item,
    _fullResponse: true
  })
}

/***
 * Update project material item
 * 更新项目材料项
 ***/
export function updateProjectMaterial(projectId: number, itemId: number, item: Partial<ProjectMaterialItem>) {
  return request.put<ProjectMaterialItem>({
    url: `/api/designer/projects/${projectId}/materials/${itemId}`,
    data: item,
    _fullResponse: true
  })
}

/***
 * Delete project material item
 * 删除项目材料项
 ***/
export function deleteProjectMaterial(projectId: number, itemId: number) {
  return request.del({
    url: `/api/designer/projects/${projectId}/materials/${itemId}`,
    _fullResponse: true
  })
}

/***
 * Get project cost summary
 * 获取项目成本汇总
 ***/
export function getProjectCostSummary(projectId: number) {
  return request.get<ProjectCostSummary>({
    url: `/api/designer/projects/${projectId}/cost`,
    _fullResponse: true
  })
}

/***
 * Save project cost config
 * 保存项目成本配置
 ***/
export function saveProjectCost(data: SaveProjectCostRequest) {
  return request.post<ProjectCostSummary>({
    url: `/api/designer/projects/${data.projectId}/cost`,
    data,
    _fullResponse: true
  })
}

/***
 * Import materials from library to project
 * 从材料库导入材料到项目
 ***/
export function importMaterialsToProject(projectId: number, materialIds: number[], defaultQuantity = 1) {
  return request.post<ProjectMaterialListResponse>({
    url: `/api/designer/projects/${projectId}/materials/import`,
    data: { materialIds, defaultQuantity },
    _fullResponse: true
  })
}

/***
 * Export project material list
 * 导出项目材料清单
 ***/
export function exportProjectMaterials(projectId: number, format: 'excel' | 'pdf' = 'excel') {
  return request.post<{ fileUrl: string; fileName: string }>({
    url: `/api/designer/projects/${projectId}/materials/export`,
    data: { format },
    _fullResponse: true
  })
}
