/**
 * 通用类型定义
 * Common type definitions
 */

/*** 平台类型 - Platform type ***/
export type PlatformType = 'app' | 'mp-weixin' | 'h5'

/*** 网络错误类型 - Network error type ***/
export enum NetworkErrorType {
  TIMEOUT = 'TIMEOUT',
  OFFLINE = 'OFFLINE',
  SERVER_ERROR = 'SERVER_ERROR',
  UNAUTHORIZED = 'UNAUTHORIZED'
}

/*** 业务错误码 - Business error codes ***/
export const BusinessErrorCode = {
  PROJECT_NOT_FOUND: 40001,
  DOCUMENT_NOT_ANALYZED: 40002,
  AI_SERVICE_UNAVAILABLE: 50001
} as const

/*** 错误提示映射 - Error message mapping ***/
export const ErrorMessages: Record<number, string> = {
  40001: '项目不存在或已被删除',
  40002: '该文档尚未分析，请在 Web 端进行分析',
  50001: 'AI 服务暂时不可用，请稍后重试'
}

/*** 项目状态枚举 - Project status enum ***/
export enum ProjectStatus {
  DRAFT = 'draft',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  ARCHIVED = 'archived'
}

/*** 项目状态文本映射 - Project status text mapping ***/
export const ProjectStatusText: Record<ProjectStatus, string> = {
  [ProjectStatus.DRAFT]: '草稿',
  [ProjectStatus.IN_PROGRESS]: '进行中',
  [ProjectStatus.COMPLETED]: '已完成',
  [ProjectStatus.ARCHIVED]: '已归档'
}

/*** 文档分析状态枚举 - Document analysis status enum ***/
export enum DocumentAnalysisStatus {
  PENDING = 'pending',
  ANALYZING = 'analyzing',
  COMPLETED = 'completed',
  FAILED = 'failed'
}

/*** 成本类别 - Cost categories ***/
export const CostCategories = {
  MATERIAL: '材料费',
  LABOR: '人工费',
  EQUIPMENT: '设备费',
  MANAGEMENT: '管理费'
} as const

/*** 材料分类 - Material categories ***/
export const MaterialCategories = [
  { name: '地板', icon: '🪵' },
  { name: '墙面', icon: '🧱' },
  { name: '天花', icon: '💡' },
  { name: '家具', icon: '🪑' },
  { name: '灯具', icon: '💡' },
  { name: '卫浴', icon: '🚿' },
  { name: '门窗', icon: '🚪' },
  { name: '其他', icon: '📦' }
] as const

/*** 缓存过期时间（7天，毫秒） - Cache expiration time (7 days in ms) ***/
export const CACHE_EXPIRE_TIME = 7 * 24 * 60 * 60 * 1000

/*** 分页默认大小 - Default page size ***/
export const DEFAULT_PAGE_SIZE = 10

/*** 请求超时时间（毫秒） - Request timeout (ms) ***/
export const REQUEST_TIMEOUT = 30000

/*** AI 请求超时时间（毫秒） - AI request timeout (ms) ***/
export const AI_REQUEST_TIMEOUT = 60000
