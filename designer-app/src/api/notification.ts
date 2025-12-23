/**
 * 通知相关 API（仅 App 端）
 * Notification related API (App only)
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, Notification, PageParams } from '@/types/api'

/*** Notification API - handles notification operations (App only) ***/
export const notificationApi = {
  /*** 获取通知列表 - Get notification list ***/
  getList(params?: PageParams): Promise<ApiResponse<PageResponse<Notification>>> {
    return request({
      url: '/api/notifications',
      method: 'GET',
      data: params
    })
  },

  /*** 标记为已读 - Mark as read ***/
  markAsRead(id: number): Promise<ApiResponse<null>> {
    return request({
      url: `/api/notifications/${id}/read`,
      method: 'PUT'
    })
  },

  /*** 获取未读数量 - Get unread count ***/
  getUnreadCount(): Promise<ApiResponse<number>> {
    return request({
      url: '/api/notifications/unread-count',
      method: 'GET'
    })
  },

  /*** 标记全部已读 - Mark all as read ***/
  markAllAsRead(): Promise<ApiResponse<null>> {
    return request({
      url: '/api/notifications/read-all',
      method: 'PUT'
    })
  }
}
