/***
 * Construction Annotation API
 * 施工图标注相关接口
 * Requirements: 9.1, 9.2, 9.3, 9.4, 9.5
 ***/
import request from '@/utils/http'

// 分析施工图请求参数
// Analyze construction drawing request parameters
export interface AnalyzeConstructionDrawingRequest {
  projectId: number
  versionId?: number // 关联设计版本ID / Associated design version ID
  cadFileId: number
  imagePath: string
}

// 标注项数据结构
// Annotation item data structure
export interface AnnotationItem {
  id: string
  type: string // dimension, material, process
  position: {
    x: number
    y: number
    anchor: string
  }
  content: string
  style: {
    fontSize: number
    fontColor: string
    lineColor: string
    lineWidth: number
    background: string
  }
  properties: Record<string, any>
  createdAt: string
  updatedAt: string
  isGenerated: boolean
}

// 检测到的元素
// Detected element
export interface DetectedElement {
  type: string
  confidence: number
  boundingBox: {
    x: number
    y: number
    width: number
    height: number
  }
  properties: Record<string, any>
}

// 施工图标注响应
// Construction annotation response
export interface ConstructionAnnotationResponse {
  id: number
  projectId: number
  versionId?: number // 关联设计版本ID / Associated design version ID
  cadFileId: number
  imagePath: string
  analysisStatus: string
  elementsDetected: DetectedElement[]
  annotations: AnnotationItem[]
  errorMessage: string
  createdAt: string
  updatedAt: string
  project?: {
    id: number
    name: string
    description: string
  }
  version?: {
    id: number
    versionNumber: number
    versionName: string
  }
  cadFile?: {
    id: number
    fileName: string
    fileFormat: string
  }
}

// 标注列表响应
// Annotation list response
export interface AnnotationListResponse {
  list: ConstructionAnnotationResponse[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// 导出配置
// Export configuration
export interface AnnotationExportConfig {
  format: string // pdf, png, jpg
  quality: number // 1-100
  scale: number
  showLayers: boolean
}

// 导出结果
// Export result
export interface ExportResultResponse {
  fileName: string
  filePath: string
  fileSize: number
  format: string
  downloadUrl: string
}

// 导出格式信息
// Export format information
export interface ExportFormat {
  format: string
  name: string
  description: string
  mimeType: string
}

// 标注操作
// Annotation operation
export interface AnnotationOperation {
  operation: string // add, update, delete
  itemId?: string
  item?: Partial<AnnotationItem>
}

/***
 * Analyze construction drawing
 * 分析施工图
 ***/
export const analyzeConstructionDrawing = (data: AnalyzeConstructionDrawingRequest) => {
  return request.post<ConstructionAnnotationResponse>({
    url: '/api/designer/construction-annotation/analyze',
    data,
    _fullResponse: true
  })
}

/***
 * Get annotation details
 * 获取标注详情
 ***/
export const getAnnotation = (id: number) => {
  return request.get<ConstructionAnnotationResponse>({
    url: `/api/designer/construction-annotation/${id}`,
    _fullResponse: true
  })
}

/***
 * Get annotation list
 * 获取标注列表
 ***/
export const getAnnotationsList = (params: {
  projectId?: number
  versionId?: number // 关联设计版本ID / Associated design version ID
  cadFileId?: number
  status?: string
  page?: number
  pageSize?: number
}) => {
  return request.get<AnnotationListResponse>({
    url: '/api/designer/construction-annotation/list',
    params,
    _fullResponse: true
  })
}

/***
 * Update annotations
 * 更新标注
 ***/
export const updateAnnotations = (id: number, annotations: AnnotationItem[]) => {
  return request.put({
    url: `/api/designer/construction-annotation/${id}/annotations`,
    data: { annotations },
    _fullResponse: true
  })
}

/***
 * Delete annotation
 * 删除标注
 ***/
export const deleteAnnotation = (id: number) => {
  return request.del({
    url: `/api/designer/construction-annotation/${id}`,
    _fullResponse: true
  })
}

/***
 * Add annotation item
 * 添加标注项
 ***/
export const addAnnotationItem = (id: number, item: Partial<AnnotationItem>) => {
  return request.post({
    url: `/api/designer/construction-annotation/${id}/items`,
    data: item,
    _fullResponse: true
  })
}

/***
 * Update annotation item
 * 更新标注项
 ***/
export const updateAnnotationItem = (id: number, itemId: string, item: Partial<AnnotationItem>) => {
  return request.put({
    url: `/api/designer/construction-annotation/${id}/items/${itemId}`,
    data: item,
    _fullResponse: true
  })
}

/***
 * Delete annotation item
 * 删除标注项
 ***/
export const deleteAnnotationItem = (id: number, itemId: string) => {
  return request.del({
    url: `/api/designer/construction-annotation/${id}/items/${itemId}`,
    _fullResponse: true
  })
}

/***
 * Get annotation item
 * 获取标注项
 ***/
export const getAnnotationItem = (id: number, itemId: string) => {
  return request.get<AnnotationItem>({
    url: `/api/designer/construction-annotation/${id}/items/${itemId}`,
    _fullResponse: true
  })
}

/***
 * Batch update annotations
 * 批量更新标注
 ***/
export const batchUpdateAnnotations = (id: number, operations: AnnotationOperation[]) => {
  return request.post({
    url: `/api/designer/construction-annotation/${id}/batch`,
    data: operations,
    _fullResponse: true
  })
}

/***
 * Export annotation
 * 导出标注
 ***/
export const exportAnnotation = (annotationId: number, config: AnnotationExportConfig) => {
  return request.post<ExportResultResponse>({
    url: '/api/designer/construction-annotation/export',
    data: { annotationId, config },
    _fullResponse: true
  })
}

/***
 * Get supported export formats
 * 获取支持的导出格式
 ***/
export const getExportFormats = () => {
  return request.get<ExportFormat[]>({
    url: '/api/designer/construction-annotation/export-formats',
    _fullResponse: true
  })
}

/***
 * Get export history
 * 获取导出历史
 ***/
export const getExportHistory = (id: number) => {
  return request.get({
    url: `/api/designer/construction-annotation/${id}/export-history`,
    _fullResponse: true
  })
}
