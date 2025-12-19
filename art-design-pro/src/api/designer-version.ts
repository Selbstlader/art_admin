/***
 * Designer Version API
 * 设计版本接口
 * Requirements: 7.1, 7.2, 7.3, 7.4
 ***/
import request from '@/utils/http'

/*** Design Image Response ***/
export interface DesignImageResponse {
  url: string
  name: string
  description: string
  type: string
}

/*** Layout Info Response ***/
export interface LayoutInfoResponse {
  zone: string
  area: number
  position: string
}

/*** Area Info Response ***/
export interface AreaInfoResponse {
  name: string
  value: number
  unit: string
}

/*** Style Info Response ***/
export interface StyleInfoResponse {
  category: string
  value: string
  details: string
}

/*** Material Info Response ***/
export interface MaterialInfoResponse {
  name: string
  category: string
  quantity: number
  unit: string
}

/*** CAD File Brief Response ***/
export interface CadFileBriefResponse {
  id: number
  fileName: string
  fileFormat: string
  filePath: string
  isAiGenerated: boolean
}

/*** Design Version Response ***/
export interface DesignVersionResponse {
  id: number
  projectId: number
  versionNumber: number
  versionName: string
  description: string
  designImages: DesignImageResponse[]
  cadFileIds: number[]
  cadFiles: CadFileBriefResponse[]
  layoutInfo: LayoutInfoResponse[]
  areaInfo: AreaInfoResponse[]
  styleInfo: StyleInfoResponse[]
  materialInfo: MaterialInfoResponse[]
  status: string
  createdBy: number
  createdAt: string
  updatedAt: string
}

/*** Design Version Brief Response ***/
export interface DesignVersionBriefResponse {
  id: number
  versionNumber: number
  versionName: string
  status: string
}

/*** Change Item Response ***/
export interface ChangeItemResponse {
  field: string
  oldValue: string
  newValue: string
  changeType: 'added' | 'removed' | 'modified'
  description: string
}

/*** Version Compare Response ***/
export interface VersionCompareResponse {
  id: number
  projectId: number
  versionAId: number
  versionBId: number
  versionA: DesignVersionBriefResponse | null
  versionB: DesignVersionBriefResponse | null
  layoutChanges: ChangeItemResponse[]
  areaChanges: ChangeItemResponse[]
  elementChanges: ChangeItemResponse[]
  styleChanges: ChangeItemResponse[]
  materialChanges: ChangeItemResponse[]
  summary: string
  compareStatus: string
  createdAt: string
  updatedAt: string
}

/*** Version Diff Response ***/
export interface VersionDiffResponse {
  versionA: DesignVersionResponse | null
  versionB: DesignVersionResponse | null
  layoutChanges: ChangeItemResponse[]
  areaChanges: ChangeItemResponse[]
  elementChanges: ChangeItemResponse[]
  styleChanges: ChangeItemResponse[]
  materialChanges: ChangeItemResponse[]
  summary: string
}

/*** Design Version List Response ***/
export interface DesignVersionListResponse {
  records: DesignVersionResponse[]
  current: number
  size: number
  total: number
}

/*** Version Compare List Response ***/
export interface VersionCompareListResponse {
  records: VersionCompareResponse[]
  current: number
  size: number
  total: number
}

/*** Design Version List Request ***/
export interface DesignVersionListRequest {
  projectId: number
  current: number
  size: number
  status?: string
}

/*** Create Design Version Request ***/
export interface CreateDesignVersionRequest {
  projectId: number
  versionName?: string
  description?: string
  designImages?: string[]
  cadFileIds?: number[]
  layoutInfo?: string
  areaInfo?: string
  styleInfo?: string
  materialInfo?: string
}

/*** Update Design Version Request ***/
export interface UpdateDesignVersionRequest {
  id: number
  versionName?: string
  description?: string
  designImages?: string[]
  cadFileIds?: number[]
  layoutInfo?: string
  areaInfo?: string
  styleInfo?: string
  materialInfo?: string
  status?: string
}

/*** Compare Design Versions Request ***/
export interface CompareDesignVersionsRequest {
  projectId: number
  versionAId: number
  versionBId: number
}

/*** Version Compare List Request ***/
export interface VersionCompareListRequest {
  projectId: number
  current: number
  size: number
  compareStatus?: string
}

/*** Get design version list ***/
export function getDesignVersionList(params: DesignVersionListRequest) {
  return request.get<Http.BaseResponse<DesignVersionListResponse>>({
    url: '/api/designer/versions',
    params,
    _fullResponse: true
  })
}

/*** Get design version detail ***/
export function getDesignVersionDetail(id: number) {
  return request.get<Http.BaseResponse<DesignVersionResponse>>({
    url: `/api/designer/versions/${id}`,
    _fullResponse: true
  })
}

/*** Create design version ***/
export function createDesignVersion(data: CreateDesignVersionRequest) {
  return request.post<Http.BaseResponse<DesignVersionResponse>>({
    url: '/api/designer/versions',
    data,
    _fullResponse: true
  })
}

/*** Update design version ***/
export function updateDesignVersion(data: UpdateDesignVersionRequest) {
  return request.put<Http.BaseResponse<DesignVersionResponse>>({
    url: `/api/designer/versions/${data.id}`,
    data,
    _fullResponse: true
  })
}

/*** Delete design version ***/
export function deleteDesignVersion(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/versions/${id}`,
    _fullResponse: true
  })
}

/*** Compare two design versions ***/
export function compareDesignVersions(data: CompareDesignVersionsRequest) {
  return request.post<Http.BaseResponse<VersionCompareResponse>>({
    url: '/api/designer/versions/compare',
    data,
    _fullResponse: true
  })
}

/*** Get version compare list ***/
export function getVersionCompareList(params: VersionCompareListRequest) {
  return request.get<Http.BaseResponse<VersionCompareListResponse>>({
    url: '/api/designer/versions/compares',
    params,
    _fullResponse: true
  })
}

/*** Get version compare detail ***/
export function getVersionCompareDetail(id: number) {
  return request.get<Http.BaseResponse<VersionCompareResponse>>({
    url: `/api/designer/versions/compares/${id}`,
    _fullResponse: true
  })
}

/*** Get version diff for side-by-side display ***/
export function getVersionDiff(versionAId: number, versionBId: number) {
  return request.get<Http.BaseResponse<VersionDiffResponse>>({
    url: '/api/designer/versions/diff',
    params: { versionAId, versionBId },
    _fullResponse: true
  })
}

/*** Delete version compare record ***/
export function deleteVersionCompare(id: number) {
  return request.del<Http.BaseResponse<null>>({
    url: `/api/designer/versions/compares/${id}`,
    _fullResponse: true
  })
}
