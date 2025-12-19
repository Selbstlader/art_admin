/***
 * Designer Material API
 * 设计师材料库管理 API 接口
 * Requirements: 6.2, 6.3, 6.4, 6.5
 ***/
import request from '@/utils/http'

/*** Material API Types ***/

// 材料列表请求参数 / Material list request params
export interface MaterialListParams {
  current: number
  size: number
  name?: string
  category?: string
  brand?: string
  minPrice?: number
  maxPrice?: number
  status?: string
  keyword?: string
}

// 材料响应 / Material response
export interface MaterialResponse {
  id: number
  name: string
  category: string
  specification: string
  unit: string
  unitPrice: number
  brand: string
  supplier: string
  description: string
  imageUrl: string
  applicableScenes: string[]
  status: string
  createdAt: string
  updatedAt: string
}

// 材料列表响应 / Material list response
export interface MaterialListResponse {
  records: MaterialResponse[]
  current: number
  size: number
  total: number
}

// 创建材料请求 / Create material request
export interface CreateMaterialRequest {
  name: string
  category: string
  specification?: string
  unit: string
  unitPrice: number
  brand?: string
  supplier?: string
  description?: string
  imageUrl?: string
  applicableScenes?: string[]
  status?: string
}

// 更新材料请求 / Update material request
export interface UpdateMaterialRequest extends CreateMaterialRequest {
  id: number
}

// 材料分类响应 / Material category response
export interface MaterialCategoryResponse {
  category: string
  count: number
}

// 材料品牌响应 / Material brand response
export interface MaterialBrandResponse {
  brand: string
  count: number
}

// 材料统计响应 / Material stats response
export interface MaterialStatsResponse {
  totalCount: number
  activeCount: number
  inactiveCount: number
  categoryStats: MaterialCategoryResponse[]
  brandStats: MaterialBrandResponse[]
  avgPrice: number
  maxPrice: number
  minPrice: number
}

// 材料推荐请求 / Material recommend request
export interface MaterialRecommendRequest {
  projectId: number
  spaceType?: string
  style?: string
  budget?: number
  area?: number
  category?: string
  limit?: number
}

// 材料推荐响应 / Material recommend response
export interface MaterialRecommendResponse {
  materials: MaterialResponse[]
  recommendReason: string
  totalCount: number
}

// 材料用量计算请求 / Material calculate request
export interface MaterialCalculateRequest {
  materialId: number
  area: number
  lossRate?: number
}

// 材料用量计算响应 / Material calculate response
export interface MaterialCalculateResponse {
  materialId: number
  materialName: string
  unit: string
  unitPrice: number
  area: number
  lossRate: number
  quantity: number
  baseQuantity: number
  totalCost: number
}

// 批量计算响应 / Batch calculate response
export interface MaterialBatchCalculateResponse {
  items: MaterialCalculateResponse[]
  totalCost: number
  totalQuantity: number
}

/*** Material API Functions ***/

/***
 * Get material list
 * 获取材料列表
 ***/
export function getMaterialList(params: MaterialListParams) {
  return request.get<MaterialListResponse>({
    url: '/api/designer/materials',
    params,
    _fullResponse: true
  })
}

/***
 * Get material detail
 * 获取材料详情
 ***/
export function getMaterialDetail(id: number) {
  return request.get<MaterialResponse>({
    url: `/api/designer/materials/${id}`,
    _fullResponse: true
  })
}

/***
 * Create material
 * 创建材料
 ***/
export function createMaterial(data: CreateMaterialRequest) {
  return request.post<MaterialResponse>({
    url: '/api/designer/materials',
    data,
    _fullResponse: true
  })
}

/***
 * Update material
 * 更新材料
 ***/
export function updateMaterial(data: UpdateMaterialRequest) {
  return request.put<MaterialResponse>({
    url: '/api/designer/materials',
    data,
    _fullResponse: true
  })
}

/***
 * Delete material
 * 删除材料
 ***/
export function deleteMaterial(id: number) {
  return request.del({
    url: `/api/designer/materials/${id}`,
    _fullResponse: true
  })
}

/***
 * Batch delete materials
 * 批量删除材料
 ***/
export function batchDeleteMaterials(ids: number[]) {
  return request.post({
    url: '/api/designer/materials/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

/***
 * Get material categories
 * 获取材料分类列表
 ***/
export function getMaterialCategories() {
  return request.get<MaterialCategoryResponse[]>({
    url: '/api/designer/materials/categories',
    _fullResponse: true
  })
}

/***
 * Get material brands
 * 获取材料品牌列表
 ***/
export function getMaterialBrands() {
  return request.get<MaterialBrandResponse[]>({
    url: '/api/designer/materials/brands',
    _fullResponse: true
  })
}

/***
 * Get material stats
 * 获取材料统计信息
 ***/
export function getMaterialStats() {
  return request.get<MaterialStatsResponse>({
    url: '/api/designer/materials/stats',
    _fullResponse: true
  })
}

/***
 * Recommend materials
 * 智能推荐材料
 ***/
export function recommendMaterials(data: MaterialRecommendRequest) {
  return request.post<MaterialRecommendResponse>({
    url: '/api/designer/materials/recommend',
    data,
    _fullResponse: true
  })
}

/***
 * Calculate material usage
 * 计算材料用量
 ***/
export function calculateMaterialUsage(data: MaterialCalculateRequest) {
  return request.post<MaterialCalculateResponse>({
    url: '/api/designer/materials/calculate',
    data,
    _fullResponse: true
  })
}

/***
 * Batch calculate material usage
 * 批量计算材料用量
 ***/
export function batchCalculateMaterialUsage(items: MaterialCalculateRequest[]) {
  return request.post<MaterialBatchCalculateResponse>({
    url: '/api/designer/materials/batch-calculate',
    data: { items },
    _fullResponse: true
  })
}

/***
 * Batch import materials
 * 批量导入材料
 ***/
export function batchImportMaterials(materials: CreateMaterialRequest[]) {
  return request.post({
    url: '/api/designer/materials/batch-import',
    data: { materials },
    _fullResponse: true
  })
}
