/**
 * 通知管理API
 * Notification Management API
 */
import http from '@/utils/http'

// 通知项类型 / Notification item type
export interface NotificationItem {
  id: number
  userId: number
  title: string
  content: string
  type: 'notice' | 'message' | 'email'
  isRead: boolean
  relatedId: number
  createdAt: string
}

// 通知列表响应 / Notification list response
export interface NotificationListResponse {
  records: NotificationItem[]
  total: number
  current: number
  size: number
}

// 未读数量响应 / Unread count response
export interface UnreadCountResponse {
  count: number
}

/**
 * 获取通知列表
 * Get notification list
 */
export function getNotificationList(params: { current?: number; size?: number }) {
  return http.get<NotificationListResponse>({
    url: '/api/notification/list',
    params,
    _fullResponse: true
  })
}

/**
 * 获取未读通知数量
 * Get unread notification count
 */
export function getUnreadCount() {
  return http.get<UnreadCountResponse>({
    url: '/api/notification/unread-count',
    _fullResponse: true
  })
}

/**
 * 标记通知为已读
 * Mark notification as read
 */
export function markNotificationRead(id: number) {
  return http.post({
    url: `/api/notification/read/${id}`,
    _fullResponse: true
  })
}

/**
 * 标记全部通知为已读
 * Mark all notifications as read
 */
export function markAllNotificationsRead() {
  return http.post({
    url: '/api/notification/read-all',
    _fullResponse: true
  })
}

/**
 * 删除通知
 * Delete notification
 */
export function deleteNotification(id: number) {
  return http.del({
    url: `/api/notification/${id}`,
    _fullResponse: true
  })
}
