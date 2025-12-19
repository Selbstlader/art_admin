/***
 * Designer Compare API
 * 设计比对接口（支持多图片、多文档、CAD文件）
 * Requirements: 2.1, 2.2, 2.3, 2.4, 2.5
 ***/
import request from '@/utils/http'

/*** Design Image Info Response ***/
export interface DesignImageInfoResponse {
  fileName: string
  filePath: string
  fileUrl: string
  fileType: 'image' | 'pdf' | 'cad'
  fileSize: number
}

/*** Match Item Response ***/
export interface MatchItemResponse {
  requirement: string
  designMatch: string
  score: number
  sourceImage?: string
  sourceDoc?: string
}

/*** Deviation Item Response ***/
export interface DeviationItemResponse {
  location: string
  content: string
  originalRequirement: string
  severity: 'low' | 'medium' | 'high'
  sourceImage?: string
  sourceDoc?: string
}

/*** Suggestion Item Response ***/
export interface SuggestionItemResponse {
  content: string
  priority: 'low' | 'medium' | 'high'
  costImpact: string
}

/*** Design Compare Upload Response ***/
export interface DesignCompareUploadResponse {
  id: number
  projectId: number
  name: string
  documentIds: number[]
  designImages: DesignImageInfoResponse[]
  analysisStatus: string
}

/*** Design Compare Response ***/
export interface DesignCompareResponse {
  id: number
  projectId: number
  name: string
  documentIds: number[]
  designImages: DesignImageInfoResponse[]
  matchItems: MatchItemResponse[]
  deviationItems: DeviationItemResponse[]
  suggestions: SuggestionItemResponse[]
  overallScore: number
  analysisStatus: string
  errorMessage: string
  createdAt: string
  updatedAt: string
}

/*** Design Compare Analysis Result ***/
export interface DesignCompareAnalysisResult {
  compareId: number
  matchItems: MatchItemResponse[]
  deviationItems: DeviationItemResponse[]
  suggestions: SuggestionItemResponse[]
  overallScore: number
}

/*** Design Compare List Response ***/
export interface DesignCompareListResponse {
  records: DesignCompareResponse[]
  current: number
  size: number
  total: number
}

/*** Design Compare List Request ***/
export interface DesignCompareListRequest {
  projectId: number
  current: number
  size: number
  analysisStatus?: string
}

/*** Upload design images for comparison (supports multiple files and CAD) ***/
export function uploadDesignImages(
  projectId: number,
  documentIds: number[],
  files: File[],
  name?: string
) {
  const formData = new FormData()
  formData.append('projectId', projectId.toString())
  formData.append('documentIds', documentIds.join(','))
  if (name) {
    formData.append('name', name)
  }
  // 支持多文件上传
  files.forEach((file) => {
    formData.append('files', file)
  })

  return request.post<Http.BaseResponse<DesignCompareUploadResponse>>({
    url: '/api/designer/compare/upload',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
    _fullResponse: true
  })
}

/*** Analyze design comparison ***/
export function analyzeDesignCompare(compareId: number) {
  return request.post<Http.BaseResponse<DesignCompareAnalysisResult>>({
    url: '/api/designer/compare/analyze',
    data: { compareId },
    _fullResponse: true
  })
}

/*** Get design compare list ***/
export function getDesignCompareList(params: DesignCompareListRequest) {
  return request.get<Http.BaseResponse<DesignCompareListResponse>>({
    url: '/api/designer/compare',
    params,
    _fullResponse: true
  })
}

/*** Get design compare detail ***/
export function getDesignCompareDetail(id: number) {
  return request.get<Http.BaseResponse<DesignCompareResponse>>({
    url: `/api/designer/compare/${id}`,
    _fullResponse: true
  })
}

/*** Delete design compare record ***/
export function deleteDesignCompare(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/compare/${id}`,
    _fullResponse: true
  })
}

/*** Batch delete design compare records ***/
export function batchDeleteDesignCompare(ids: number[]) {
  return request.post<Http.BaseResponse<null>>({
    url: '/api/designer/compare/batch-delete',
    data: { ids },
    _fullResponse: true
  })
}
