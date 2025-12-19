/***
 * Designer Document API
 * 文档分析接口
 * Requirements: 1.1, 1.2, 1.3, 1.4, 1.5
 ***/
import request from '@/utils/http'

/*** Document Upload Response ***/
export interface DocumentUploadResponse {
  id: number
  fileName: string
  fileType: string
  fileSize: number
  projectId: number
  status: string
}

/*** Document Analysis Result ***/
export interface DocumentAnalysisResult {
  documentId: number
  projectName: string
  area: number
  budget: number
  style: string
  functionalZones: string[]
  keywords: string[]
  summary: string
  missingFields: string[]
  suggestions: string[]
}

/*** Document Response ***/
export interface DocumentResponse {
  id: number
  projectId: number
  fileName: string
  filePath: string
  fileType: string
  fileSize: number
  analysisStatus: string
  errorMessage: string // 分析错误信息或建议
  keywords: string[]
  summary: string
  projectName: string
  extractedArea: number
  extractedBudget: number
  extractedStyle: string
  functionalZones: string[]
  createdAt: string
  updatedAt: string
}

/*** Document List Response ***/
export interface DocumentListResponse {
  records: DocumentResponse[]
  current: number
  size: number
  total: number
}

/*** Document List Request ***/
export interface DocumentListRequest {
  projectId: number
  current: number
  size: number
  analysisStatus?: string
  fileType?: string
}

/*** Upload document to project ***/
export function uploadDocument(projectId: string, file: File) {
  if (!projectId) {
    return Promise.reject(new Error('项目ID不能为空'))
  }
  if (!file) {
    return Promise.reject(new Error('文件不能为空'))
  }
  const formData = new FormData()
  formData.append('file', file)
  formData.append('projectId', String(projectId))

  return request.post<DocumentUploadResponse>({
    url: '/api/designer/documents/upload',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
    _fullResponse: true
  })
}

/*** Analyze document ***/
export function analyzeDocument(documentId: number) {
  return request.post<DocumentAnalysisResult>({
    url: '/api/designer/documents/analyze',
    data: { documentId },
    _fullResponse: true
  })
}

/*** Get document list ***/
export function getDocumentList(params: DocumentListRequest) {
  return request.get<Http.BaseResponse<DocumentListResponse>>({
    url: '/api/designer/documents',
    params,
    _fullResponse: true
  })
}

/*** Get document detail ***/
export function getDocumentDetail(id: number) {
  return request.get<DocumentResponse>({
    url: `/api/designer/documents/${id}`,
    _fullResponse: true
  })
}

/*** Get document keywords ***/
export function getDocumentKeywords(id: number) {
  return request.get<string[]>({
    url: `/api/designer/documents/${id}/keywords`,
    _fullResponse: true
  })
}

/*** Get document summary ***/
export function getDocumentSummary(id: number) {
  return request.get<string>({
    url: `/api/designer/documents/${id}/summary`,
    _fullResponse: true
  })
}

/*** Delete document ***/
export function deleteDocument(id: number) {
  return request.del({
    url: `/api/designer/documents/${id}`,
    _fullResponse: true
  })
}

/*** Batch delete documents ***/
export function batchDeleteDocuments(ids: number[]) {
  return request.post({
    url: '/api/designer/documents/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}
