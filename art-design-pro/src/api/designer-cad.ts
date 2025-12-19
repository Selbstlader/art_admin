/**
 * CAD文件管理API
 * CAD File Management API
 * Requirements: 8.1, 8.2, 8.3, 8.5, 8.6
 */
import http from '@/utils/http'

// CAD文件响应类型 / CAD file response type
export interface CadFileResponse {
  id: number
  projectId: number
  fileName: string
  originalPath: string
  parsedPath: string
  fileFormat: string
  parseStatus: string
  layerCount: number
  layers: string[]
  has3d: boolean
  fileSize: number
  errorMessage: string
  createdAt: string
  updatedAt: string
}

// CAD文件上传响应 / CAD file upload response
export interface CadFileUploadResponse {
  id: number
  fileName: string
  fileFormat: string
  fileSize: number
  projectId: number
  status: string
}

// CAD图层信息 / CAD layer information
export interface CadLayerInfo {
  name: string
  color: number
  lineType: string
  visible: boolean
  frozen: boolean
  locked: boolean
  entityCount: number
}

// CAD实体信息 / CAD entity information
export interface CadEntityInfo {
  type: string
  layer: string
  color: number
  lineType: string
  properties: Record<string, any>
}

// 边界框 / Bounding box
export interface BoundingBox {
  minX: number
  minY: number
  minZ: number
  maxX: number
  maxY: number
  maxZ: number
}

// CAD解析结果 / CAD parse result
export interface CadParseResult {
  cadFileId: number
  fileName: string
  fileFormat: string
  layers: CadLayerInfo[]
  layerCount: number
  has3d: boolean
  boundingBox: BoundingBox
  entities: CadEntityInfo[]
  entityCount: number
  parsedData: string
}

// CAD文件列表响应 / CAD file list response
export interface CadFileListResponse {
  records: CadFileResponse[]
  current: number
  size: number
  total: number
}

// CAD图层列表响应 / CAD layer list response
export interface CadLayerListResponse {
  cadFileId: number
  fileName: string
  layerCount: number
  layers: CadLayerInfo[]
}

/**
 * 上传CAD文件
 * Upload CAD file
 */
export function uploadCadFile(formData: FormData) {
  return http.post<CadFileUploadResponse>({
    url: '/api/designer/cad/upload',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    _fullResponse: true
  })
}

/**
 * 解析CAD文件
 * Parse CAD file
 */
export function parseCadFile(cadFileId: number) {
  return http.post<CadParseResult>({
    url: '/api/designer/cad/parse',
    data: { cadFileId },
    _fullResponse: true
  })
}

/**
 * 获取CAD文件列表
 * Get CAD file list
 */
export function getCadFileList(params: {
  projectId: number
  current?: number
  size?: number
  parseStatus?: string
  fileFormat?: string
}) {
  return http.get<CadFileListResponse>({
    url: '/api/designer/cad',
    params,
    _fullResponse: true
  })
}

/**
 * 获取CAD文件详情
 * Get CAD file detail
 */
export function getCadFileDetail(id: number) {
  return http.get<CadFileResponse>({
    url: `/api/designer/cad/${id}`,
    _fullResponse: true
  })
}

/**
 * 获取CAD解析结果
 * Get CAD parse result
 */
export function getCadParseResult(id: number) {
  return http.get<CadParseResult>({
    url: `/api/designer/cad/${id}/parse`,
    _fullResponse: true
  })
}

/**
 * 获取CAD图层列表
 * Get CAD layer list
 */
export function getCadLayers(id: number) {
  return http.get<CadLayerListResponse>({
    url: `/api/designer/cad/${id}/layers`,
    _fullResponse: true
  })
}

/**
 * 删除CAD文件
 * Delete CAD file
 */
export function deleteCadFile(id: number) {
  return http.del({
    url: `/api/designer/cad/${id}`,
    _fullResponse: true
  })
}

/**
 * 批量删除CAD文件
 * Batch delete CAD files
 */
export function batchDeleteCadFiles(ids: number[]) {
  return http.post({
    url: '/api/designer/cad/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}

// 渲染结果类型 / Render result type
export interface RenderResult {
  cadFileId: number
  imageUrl: string
  designProposal?: string // AI生成的设计方案 / AI-generated design proposal
  style: string
  roomType: string
  prompt: string
  generatedAt: string
}

// 渲染风格选项 / Render style options
export interface RenderStyleOption {
  value: string
  label: string
  description: string
}

// 渲染配置响应 / Render config response
export interface RenderConfigResponse {
  styles: RenderStyleOption[]
  roomTypes: { value: string; label: string }[]
  viewAngles: { value: string; label: string }[]
}

/**
 * 生成CAD效果图
 * Generate CAD rendering
 */
export function generateCadRender(data: {
  cadFileId: number
  projectId?: number
  documentIds?: number[]
  style?: string
  roomType?: string
  description?: string
}) {
  return http.post<RenderResult>({
    url: '/api/designer/cad/render',
    data,
    _fullResponse: true
  })
}

/**
 * 获取渲染历史
 * Get render history
 */
export function getRenderHistory(cadFileId: number) {
  return http.get<{ cadFileId: number; records: RenderResult[]; total: number }>({
    url: '/api/designer/cad/render/history',
    params: { cadFileId },
    _fullResponse: true
  })
}

/**
 * 获取渲染风格配置
 * Get render style config
 */
export function getRenderStyles() {
  return http.get<RenderConfigResponse>({
    url: '/api/designer/cad/render/styles',
    _fullResponse: true
  })
}

// 保存效果图响应 / Save render response
export interface SaveRenderResponse {
  id: number
  imageUrl: string
  remainingUse: number
}

// 用户配额信息 / User quota info
export interface UserRenderQuota {
  dailyLimit: number
  usedToday: number
  remainingUse: number
}

// 效果图记录 / Render record
export interface RenderRecordItem {
  id: number
  userId: number
  projectId: number
  cadFileId: number
  status: string // pending | processing | completed | failed
  imageUrl: string
  designProposal: string
  style: string
  roomType: string
  prompt: string
  errorMessage?: string // 失败时的错误信息
  createdAt: string
}

// saveRenderImage 已移除 - 异步模式下自动保存到记录中
// saveRenderImage removed - auto-saved to records in async mode

/**
 * 获取用户效果图配额
 * Get user render quota
 */
export function getUserRenderQuota() {
  return http.get<UserRenderQuota>({
    url: '/api/designer/cad/render/quota',
    _fullResponse: true
  })
}

/**
 * 获取效果图记录列表
 * Get render records list
 */
export function getRenderRecords(params: {
  current?: number
  size?: number
  projectId?: number
  cadFileId?: number
}) {
  return http.get<{ records: RenderRecordItem[]; total: number; current: number; size: number }>({
    url: '/api/designer/cad/render/records',
    params,
    _fullResponse: true
  })
}

/**
 * 删除效果图记录
 * Delete render record
 */
export function deleteRenderRecord(id: number) {
  return http.del({
    url: `/api/designer/cad/render/records/${id}`,
    _fullResponse: true
  })
}
