/**
 * 旅游规划模块 API 接口
 */
import request from '@/utils/http'
import type {
  CreateRoadbookRequest,
  UpdateRoadbookRequest,
  RoadbookListParams,
  AddWaypointRequest,
  BatchUpdateWaypointsRequest,
  CreateCommentRequest,
  CommentListParams,
  CreateShareRequest,
  NavLinkRequest,
  SearchRoadbookParams,
  AIPlanRequest,
  AIChatRequest,
  UseTemplateRequest,
  PaginatedResponse,
  Roadbook,
  RoadbookDetailResponse,
  Waypoint,
  Comment,
  Tag,
  Share,
  ShareResponse,
  NavLinkResponse,
  RoadbookTemplate,
  AIPlanResponse,
  AIChatResponse,
  StatsResponse
} from '../types'

const BASE_URL = '/api/v1/travel'

// ==================== 路书相关接口 ====================

export const roadbookApi = {
  /** 创建路书 */
  create: (data: CreateRoadbookRequest) =>
    request.post<Roadbook>({ url: `${BASE_URL}/roadbooks`, data }),

  /** 获取路书列表 */
  getList: (params?: RoadbookListParams) =>
    request.get<PaginatedResponse<Roadbook>>({ url: `${BASE_URL}/roadbooks`, params }),

  /** 获取路书详情 */
  getDetail: (id: number) =>
    request.get<RoadbookDetailResponse>({ url: `${BASE_URL}/roadbooks/${id}` }),

  /** 更新路书 */
  update: (id: number, data: UpdateRoadbookRequest) =>
    request.put<Roadbook>({ url: `${BASE_URL}/roadbooks/${id}`, data }),

  /** 删除路书 */
  delete: (id: number) => request.del({ url: `${BASE_URL}/roadbooks/${id}` }),

  /** 获取我的路书列表 */
  getMyList: (params?: RoadbookListParams) =>
    request.get<PaginatedResponse<Roadbook>>({ url: `${BASE_URL}/roadbooks/my`, params }),

  /** 获取收藏的路书列表 */
  getFavoriteList: (params?: RoadbookListParams) =>
    request.get<PaginatedResponse<Roadbook>>({ url: `${BASE_URL}/roadbooks/favorites`, params })
}

// ==================== 途经点相关接口 ====================

export const waypointApi = {
  /** 添加途经点 */
  add: (roadbookId: number, data: AddWaypointRequest) =>
    request.post<Waypoint>({ url: `${BASE_URL}/roadbooks/${roadbookId}/waypoints`, data }),

  /** 获取途经点列表 */
  getList: (roadbookId: number) =>
    request.get<Waypoint[]>({ url: `${BASE_URL}/roadbooks/${roadbookId}/waypoints` }),

  /** 更新途经点 */
  update: (roadbookId: number, waypointId: number, data: Partial<AddWaypointRequest>) =>
    request.put<Waypoint>({
      url: `${BASE_URL}/roadbooks/${roadbookId}/waypoints/${waypointId}`,
      data
    }),

  /** 删除途经点 */
  delete: (roadbookId: number, waypointId: number) =>
    request.del({ url: `${BASE_URL}/roadbooks/${roadbookId}/waypoints/${waypointId}` }),

  /** 批量更新途经点 */
  batchUpdate: (roadbookId: number, data: BatchUpdateWaypointsRequest) =>
    request.put<Waypoint[]>({ url: `${BASE_URL}/roadbooks/${roadbookId}/waypoints`, data })
}

// ==================== 评论相关接口 ====================

export const commentApi = {
  /** 发表评论 */
  create: (roadbookId: number, data: CreateCommentRequest) =>
    request.post<Comment>({ url: `${BASE_URL}/roadbooks/${roadbookId}/comments`, data }),

  /** 获取评论列表 */
  getList: (roadbookId: number, params?: CommentListParams) =>
    request.get<PaginatedResponse<Comment>>({
      url: `${BASE_URL}/roadbooks/${roadbookId}/comments`,
      params
    }),

  /** 删除评论 */
  delete: (roadbookId: number, commentId: number) =>
    request.del({ url: `${BASE_URL}/roadbooks/${roadbookId}/comments/${commentId}` }),

  /** 点赞评论 */
  like: (roadbookId: number, commentId: number) =>
    request.post({ url: `${BASE_URL}/roadbooks/${roadbookId}/comments/${commentId}/like` }),

  /** 取消点赞评论 */
  unlike: (roadbookId: number, commentId: number) =>
    request.del({ url: `${BASE_URL}/roadbooks/${roadbookId}/comments/${commentId}/like` })
}

// ==================== 收藏相关接口 ====================

export const favoriteApi = {
  /** 收藏路书 */
  add: (roadbookId: number) =>
    request.post({ url: `${BASE_URL}/roadbooks/${roadbookId}/favorite` }),

  /** 取消收藏 */
  remove: (roadbookId: number) =>
    request.del({ url: `${BASE_URL}/roadbooks/${roadbookId}/favorite` }),

  /** 检查是否已收藏 */
  check: (roadbookId: number) =>
    request.get<{ isFavorited: boolean }>({
      url: `${BASE_URL}/roadbooks/${roadbookId}/favorite/check`
    })
}

// ==================== 点赞相关接口 ====================

export const likeApi = {
  /** 点赞路书 */
  add: (roadbookId: number) => request.post({ url: `${BASE_URL}/roadbooks/${roadbookId}/like` }),

  /** 取消点赞 */
  remove: (roadbookId: number) => request.del({ url: `${BASE_URL}/roadbooks/${roadbookId}/like` })
}

// ==================== 分享相关接口 ====================

export const shareApi = {
  /** 生成分享链接 */
  create: (roadbookId: number, data?: CreateShareRequest) =>
    request.post<ShareResponse>({ url: `${BASE_URL}/roadbooks/${roadbookId}/share`, data }),

  /** 获取分享记录列表 */
  getList: (roadbookId: number) =>
    request.get<Share[]>({ url: `${BASE_URL}/roadbooks/${roadbookId}/shares` }),

  /** 通过分享码获取路书 */
  getByCode: (shareCode: string) =>
    request.get<RoadbookDetailResponse>({ url: `${BASE_URL}/share/${shareCode}` })
}

// ==================== 导航相关接口 ====================

export const navApi = {
  /** 生成高德导航链接 */
  getNavLink: (roadbookId: number, params?: NavLinkRequest) =>
    request.get<NavLinkResponse>({ url: `${BASE_URL}/roadbooks/${roadbookId}/nav-link`, params })
}

// ==================== 标签相关接口 ====================

export const tagApi = {
  /** 获取标签列表 */
  getList: (params?: { keyword?: string; isSystem?: boolean }) =>
    request.get<Tag[]>({ url: `${BASE_URL}/tags`, params }),

  /** 获取热门标签 */
  getHot: (limit?: number) =>
    request.get<Tag[]>({ url: `${BASE_URL}/tags/hot`, params: { limit } }),

  /** 创建自定义标签 */
  create: (data: { name: string; color?: string }) =>
    request.post<Tag>({ url: `${BASE_URL}/tags`, data })
}

// ==================== 路书广场/搜索相关接口 ====================

export const exploreApi = {
  /** 获取路书广场列表 */
  getExploreList: (params?: { page?: number; pageSize?: number; sortBy?: string }) =>
    request.get<PaginatedResponse<Roadbook>>({ url: `${BASE_URL}/explore`, params }),

  /** 搜索路书 */
  search: (params: SearchRoadbookParams) =>
    request.get<PaginatedResponse<Roadbook>>({ url: `${BASE_URL}/search`, params }),

  /** 获取热门路书 */
  getHot: (limit?: number) =>
    request.get<Roadbook[]>({ url: `${BASE_URL}/explore/hot`, params: { limit } }),

  /** 获取推荐路书 */
  getRecommended: (limit?: number) =>
    request.get<Roadbook[]>({ url: `${BASE_URL}/explore/recommended`, params: { limit } })
}

// ==================== AI相关接口 ====================

export const aiApi = {
  /** AI生成行程规划 */
  plan: (data: AIPlanRequest) => request.post<AIPlanResponse>({ url: `${BASE_URL}/ai/plan`, data }),

  /** AI助手对话 */
  chat: (data: AIChatRequest) => request.post<AIChatResponse>({ url: `${BASE_URL}/ai/chat`, data }),

  /** 导入AI规划为路书草稿 */
  importPlan: (data: { plan: AIPlanResponse; title: string }) =>
    request.post<Roadbook>({ url: `${BASE_URL}/ai/import`, data })
}

// ==================== 模板相关接口 ====================

export const templateApi = {
  /** 获取模板列表 */
  getList: (params?: {
    category?: number
    destination?: string
    page?: number
    pageSize?: number
  }) => request.get<PaginatedResponse<RoadbookTemplate>>({ url: `${BASE_URL}/templates`, params }),

  /** 获取模板详情 */
  getDetail: (id: number) => request.get<RoadbookTemplate>({ url: `${BASE_URL}/templates/${id}` }),

  /** 使用模板创建路书 */
  use: (id: number, data?: UseTemplateRequest) =>
    request.post<Roadbook>({ url: `${BASE_URL}/templates/${id}/use`, data })
}

// ==================== 统计相关接口 ====================

export const statsApi = {
  /** 获取路书统计数据 */
  getRoadbookStats: (roadbookId: number, params?: { startDate?: string; endDate?: string }) =>
    request.get<StatsResponse>({ url: `${BASE_URL}/roadbooks/${roadbookId}/stats`, params }),

  /** 获取个人数据中心统计 */
  getMyStats: () => request.get<StatsResponse>({ url: `${BASE_URL}/stats/my` })
}

export default {
  roadbookApi,
  waypointApi,
  commentApi,
  favoriteApi,
  likeApi,
  shareApi,
  navApi,
  tagApi,
  exploreApi,
  aiApi,
  templateApi,
  statsApi
}
