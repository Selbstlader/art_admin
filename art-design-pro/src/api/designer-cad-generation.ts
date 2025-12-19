/***
 * Designer CAD Generation API
 * AI生成CAD接口
 * Requirements: AI based CAD generation from project documents
 ***/
import request from '@/utils/http'

/*** CAD Generation Parameters ***/
export interface CadGenerationParams {
  width?: number
  height?: number
  scale?: string
  style?: string
  roomTypes?: string[]
  requirements?: string
}

/*** Create CAD Generation Request ***/
export interface CreateCadGenerationRequest {
  projectId: number
  documentIds?: number[]
  generationType: string
  prompt?: string
  parameters?: CadGenerationParams
}

/*** CAD File Brief Response ***/
export interface CadFileBriefResponse {
  id: number
  fileName: string
  fileFormat: string
  filePath: string
}

/*** CAD Generation Response ***/
export interface CadGenerationResponse {
  id: number
  projectId: number
  projectName?: string
  documentIds: number[]
  generationType: string
  generationLabel: string
  prompt: string
  parameters?: CadGenerationParams
  status: string
  statusLabel: string
  progress: number
  resultFileId?: number
  resultFile?: CadFileBriefResponse
  resultFilePath?: string
  previewImageUrl?: string
  errorMessage?: string
  aiModel?: string
  processingTime: number
  createdBy: number
  createdAt: string
  updatedAt: string
}

/*** CAD Generation List Request ***/
export interface CadGenerationListRequest {
  projectId?: number
  generationType?: string
  status?: string
  current: number
  size: number
}

/*** CAD Generation List Response ***/
export interface CadGenerationListResponse {
  records: CadGenerationResponse[]
  current: number
  size: number
  total: number
}

/*** Generation Type Option ***/
export interface GenerationTypeOption {
  value: string
  label: string
}

/*** Confirm CAD File Request ***/
export interface ConfirmCadFileRequest {
  id: number
  fileName?: string
}

/*** Create CAD generation task ***/
export function createCadGeneration(data: CreateCadGenerationRequest) {
  return request.post<Http.BaseResponse<CadGenerationResponse>>({
    url: '/api/designer/cad-generations',
    data,
    _fullResponse: true
  })
}

/*** Get CAD generation task detail ***/
export function getCadGenerationDetail(id: number) {
  return request.get<Http.BaseResponse<CadGenerationResponse>>({
    url: `/api/designer/cad-generations/${id}`,
    _fullResponse: true
  })
}

/*** Get CAD generation task list ***/
export function getCadGenerationList(params: CadGenerationListRequest) {
  return request.get<Http.BaseResponse<CadGenerationListResponse>>({
    url: '/api/designer/cad-generations',
    params,
    _fullResponse: true
  })
}

/*** Delete CAD generation task ***/
export function deleteCadGeneration(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/cad-generations/${id}`,
    _fullResponse: true
  })
}

/*** Retry CAD generation task ***/
export function retryCadGeneration(id: number) {
  return request.post<Http.BaseResponse<CadGenerationResponse>>({
    url: `/api/designer/cad-generations/${id}/retry`,
    _fullResponse: true
  })
}

/*** Confirm and save CAD file ***/
export function confirmCadFile(data: ConfirmCadFileRequest) {
  return request.post<Http.BaseResponse<CadGenerationResponse>>({
    url: '/api/designer/cad-generations/confirm',
    data,
    _fullResponse: true
  })
}

/*** Get generation type options ***/
export function getGenerationTypeOptions() {
  return request.get<Http.BaseResponse<GenerationTypeOption[]>>({
    url: '/api/designer/cad-generations/type-options',
    _fullResponse: true
  })
}
