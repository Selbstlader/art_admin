/**
 * 旅游规划模块类型定义
 */

// ==================== 基础枚举 ====================

/** 路书可见性 */
export enum RoadbookVisibility {
  /** 公开 */
  PUBLIC = 1,
  /** 私有 */
  PRIVATE = 2,
  /** 指定用户 */
  SPECIFIED = 3
}

/** 路书状态 */
export enum RoadbookStatus {
  /** 草稿 */
  DRAFT = 0,
  /** 已发布 */
  PUBLISHED = 1,
  /** 已下架 */
  OFFLINE = 2
}

/** 出行方式 */
export enum TravelMode {
  /** 驾车 */
  DRIVING = 'driving',
  /** 步行 */
  WALKING = 'walking',
  /** 骑行 */
  CYCLING = 'cycling',
  /** 公交 */
  TRANSIT = 'transit'
}

/** 途经点类型 */
export enum WaypointType {
  /** 起点 */
  START = 1,
  /** 途经点 */
  WAYPOINT = 2,
  /** 终点 */
  END = 3
}

/** 模板分类 */
export enum TemplateCategory {
  /** 周边游 */
  NEARBY = 1,
  /** 国内游 */
  DOMESTIC = 2,
  /** 出境游 */
  OVERSEAS = 3,
  /** 主题游 */
  THEME = 4
}

/** 分享类型 */
export enum ShareType {
  /** 链接分享 */
  LINK = 1,
  /** 二维码分享 */
  QRCODE = 2
}

// ==================== 核心数据模型 ====================

/** 路书 */
export interface Roadbook {
  id: number
  userId: number
  title: string
  description: string
  coverUrl: string
  startDate: string
  endDate: string
  visibility: RoadbookVisibility
  travelMode: TravelMode
  totalDistance: number
  totalDuration: number
  totalBudget: number
  viewCount: number
  favoriteCount: number
  likeCount: number
  commentCount: number
  shareCount: number
  status: RoadbookStatus
  createdAt: string
  updatedAt: string
  // 关联数据
  waypoints?: Waypoint[]
  tags?: Tag[]
  author?: RoadbookAuthor
}

/** 路书作者信息 */
export interface RoadbookAuthor {
  id: number
  nickname: string
  avatar: string
}

/** 途经点 */
export interface Waypoint {
  id: number
  roadbookId: number
  name: string
  address: string
  longitude: number
  latitude: number
  poiId: string
  poiType: string
  dayIndex: number
  sortOrder: number
  stayDuration: number
  budget: number
  notes: string
  images: string[]
  waypointType: WaypointType
  createdAt: string
  updatedAt: string
}

/** 标签 */
export interface Tag {
  id: number
  name: string
  color: string
  useCount: number
  isSystem: boolean
  createdAt: string
}

/** 评论 */
export interface Comment {
  id: number
  roadbookId: number
  userId: number
  parentId: number
  content: string
  likeCount: number
  createdAt: string
  // 关联数据
  user?: CommentUser
  replies?: Comment[]
}

/** 评论用户信息 */
export interface CommentUser {
  id: number
  nickname: string
  avatar: string
}

/** 收藏记录 */
export interface Favorite {
  id: number
  userId: number
  roadbookId: number
  createdAt: string
  roadbook?: Roadbook
}

/** 分享记录 */
export interface Share {
  id: number
  roadbookId: number
  userId: number
  shareCode: string
  shareType: ShareType
  visitCount: number
  expiresAt: string | null
  createdAt: string
}

/** 路书模板 */
export interface RoadbookTemplate {
  id: number
  roadbookId: number
  userId: number
  title: string
  description: string
  coverUrl: string
  destination: string
  days: number
  category: TemplateCategory
  useCount: number
  status: number
  createdAt: string
}

/** 每日统计 */
export interface DailyStat {
  id: number
  roadbookId: number
  statDate: string
  viewCount: number
  favoriteCount: number
  shareCount: number
}

// ==================== API 请求类型 ====================

/** 创建路书请求 */
export interface CreateRoadbookRequest {
  title: string
  description?: string
  coverUrl?: string
  startDate?: string
  endDate?: string
  visibility?: RoadbookVisibility
  travelMode?: TravelMode
  totalBudget?: number
  tagIds?: number[]
}

/** 更新路书请求 */
export interface UpdateRoadbookRequest {
  title?: string
  description?: string
  coverUrl?: string
  startDate?: string
  endDate?: string
  visibility?: RoadbookVisibility
  travelMode?: TravelMode
  totalBudget?: number
  status?: RoadbookStatus
  tagIds?: number[]
}

/** 路书列表查询参数 */
export interface RoadbookListParams {
  page?: number
  pageSize?: number
  keyword?: string
  tagId?: number
  visibility?: RoadbookVisibility
  status?: RoadbookStatus
  userId?: number
  sortBy?: 'created_at' | 'view_count' | 'favorite_count'
  sortOrder?: 'asc' | 'desc'
}

/** 添加途经点请求 */
export interface AddWaypointRequest {
  name: string
  address?: string
  longitude: number
  latitude: number
  poiId?: string
  poiType?: string
  dayIndex?: number
  sortOrder?: number
  stayDuration?: number
  budget?: number
  notes?: string
  images?: string[]
  waypointType?: WaypointType
}

/** 批量更新途经点请求 */
export interface BatchUpdateWaypointsRequest {
  waypoints: UpdateWaypointItem[]
}

/** 更新途经点项 */
export interface UpdateWaypointItem {
  id: number
  dayIndex?: number
  sortOrder?: number
  stayDuration?: number
  budget?: number
  notes?: string
}

/** 发表评论请求 */
export interface CreateCommentRequest {
  content: string
  parentId?: number
}

/** 评论列表查询参数 */
export interface CommentListParams {
  page?: number
  pageSize?: number
}

/** 生成分享链接请求 */
export interface CreateShareRequest {
  shareType?: ShareType
  expiresAt?: string
}

/** 生成导航链接请求 */
export interface NavLinkRequest {
  mode?: TravelMode
  waypointIds?: number[]
}

/** 路书搜索参数 */
export interface SearchRoadbookParams {
  keyword?: string
  tagIds?: number[]
  destination?: string
  minDays?: number
  maxDays?: number
  page?: number
  pageSize?: number
  sortBy?: 'relevance' | 'created_at' | 'view_count' | 'favorite_count'
}

/** AI规划请求 */
export interface AIPlanRequest {
  destination: string
  days: number
  preferences?: string[]
  budget?: number
  travelMode?: TravelMode
}

/** AI对话请求 */
export interface AIChatRequest {
  message: string
  roadbookId?: number
  context?: string
}

/** 使用模板请求 */
export interface UseTemplateRequest {
  startDate?: string
}

// ==================== API 响应类型 ====================

/** 分页响应 */
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

/** 路书详情响应 */
export interface RoadbookDetailResponse extends Roadbook {
  waypoints: Waypoint[]
  tags: Tag[]
  author: RoadbookAuthor
  isFavorited: boolean
  isLiked: boolean
}

/** 导航链接响应 */
export interface NavLinkResponse {
  url: string
  segments?: NavSegment[]
  needSegment: boolean
}

/** 导航分段 */
export interface NavSegment {
  index: number
  url: string
  waypointNames: string[]
}

/** 分享链接响应 */
export interface ShareResponse {
  shareCode: string
  shareUrl: string
  qrcodeUrl?: string
}

/** AI规划响应 */
export interface AIPlanResponse {
  plan: AIPlanDay[]
  summary: string
}

/** AI规划每日行程 */
export interface AIPlanDay {
  dayIndex: number
  date: string
  waypoints: AIPlanWaypoint[]
  meals: AIPlanMeal[]
  accommodation?: AIPlanAccommodation
}

/** AI规划途经点 */
export interface AIPlanWaypoint {
  name: string
  address: string
  longitude: number
  latitude: number
  description: string
  stayDuration: number
  tips: string
}

/** AI规划餐饮 */
export interface AIPlanMeal {
  type: 'breakfast' | 'lunch' | 'dinner'
  name: string
  address: string
  recommendation: string
}

/** AI规划住宿 */
export interface AIPlanAccommodation {
  name: string
  address: string
  priceRange: string
}

/** AI对话响应 */
export interface AIChatResponse {
  reply: string
  suggestions?: string[]
}

/** 统计数据响应 */
export interface StatsResponse {
  totalViews: number
  totalFavorites: number
  totalShares: number
  totalComments: number
  dailyStats: DailyStat[]
}

// ==================== 地图相关类型 ====================

/** POI搜索结果 */
export interface POISearchResult {
  id: string
  name: string
  address: string
  location: {
    lng: number
    lat: number
  }
  type: string
  distance?: number
}

/** 路线规划结果 */
export interface RouteResult {
  distance: number
  duration: number
  polyline: [number, number][]
  steps: RouteStep[]
}

/** 路线步骤 */
export interface RouteStep {
  instruction: string
  distance: number
  duration: number
  polyline: [number, number][]
}

/** 地图标记 */
export interface MapMarker {
  id: string | number
  position: [number, number]
  title: string
  icon?: string
  label?: string
  data?: any
}

export default {}
