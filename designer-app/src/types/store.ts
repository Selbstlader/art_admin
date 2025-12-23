/**
 * Store 状态类型定义
 * Store state type definitions
 */

import type { 
  UserInfo, 
  Project, 
  ProjectDetail, 
  ChatMessage, 
  ChatSession, 
  Notification 
} from './api'

/*** 用户状态 - User state ***/
export interface UserState {
  token: string
  refreshToken: string
  tokenExpireTime: number
  userInfo: UserInfo | null
  isLoggedIn: boolean
}

/*** 项目状态 - Project state ***/
export interface ProjectState {
  list: Project[]
  currentProject: ProjectDetail | null
  total: number
  loading: boolean
  searchKeyword: string
  currentPage: number
  hasMore: boolean
}

/*** 聊天状态 - Chat state (App only) ***/
export interface ChatState {
  messages: ChatMessage[]
  currentSessionId: string
  isStreaming: boolean
  sessions: ChatSession[]
  projectContext: number | null
}

/*** 通知状态 - Notification state (App only) ***/
export interface NotificationState {
  list: Notification[]
  unreadCount: number
  hasPermission: boolean
  loading: boolean
}

/*** 缓存数据结构 - Cache data structure ***/
export interface CacheData<T> {
  data: T
  timestamp: number
  expireTime: number
}

/*** 项目列表缓存 - Project list cache ***/
export interface ProjectListCache extends CacheData<Project[]> {
  total: number
  lastPage: number
}

/*** 离线状态 - Offline state ***/
export interface OfflineState {
  isOffline: boolean
  lastOnlineTime: number
}
