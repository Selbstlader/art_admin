/**
 * 通知状态管理（仅 App 端）
 * Notification state management (App only)
 */

import { defineStore } from 'pinia'
import type { NotificationState } from '@/types/store'
import type { Notification } from '@/types/api'

/*** Notification store - manages push notifications (App only) ***/
export const useNotificationStore = defineStore('notification', {
  state: (): NotificationState => ({
    list: [],
    unreadCount: 0,
    hasPermission: false,
    loading: false
  }),

  getters: {
    /*** 是否有未读消息 - Has unread messages ***/
    hasUnread: (state) => state.unreadCount > 0,
    
    /*** 按时间倒序排列的通知列表 - Notifications sorted by time desc ***/
    sortedList: (state) => {
      return [...state.list].sort((a, b) => 
        new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      )
    }
  },

  actions: {
    /*** 获取通知列表 - Fetch notification list ***/
    async fetchList(): Promise<void> {
      this.loading = true
      try {
        // TODO: 调用后端接口
        // const res = await notificationApi.getList()
        // this.list = res.data
        // this.updateUnreadCount()
      } catch (error) {
        console.error('获取通知列表失败:', error)
      } finally {
        this.loading = false
      }
    },

    /*** 标记为已读 - Mark as read ***/
    async markAsRead(id: number): Promise<void> {
      const notification = this.list.find(n => n.id === id)
      if (notification && !notification.isRead) {
        notification.isRead = true
        this.updateUnreadCount()
        
        // TODO: 调用后端接口
        // await notificationApi.markAsRead(id)
      }
    },

    /*** 更新未读数量 - Update unread count ***/
    updateUnreadCount(): void {
      this.unreadCount = this.list.filter(n => !n.isRead).length
    },

    /*** 获取未读数量 - Get unread count ***/
    async getUnreadCount(): Promise<number> {
      // TODO: 调用后端接口
      // const res = await notificationApi.getUnreadCount()
      // this.unreadCount = res.data
      return this.unreadCount
    },

    /*** 设置推送权限状态 - Set push permission status ***/
    setPermission(hasPermission: boolean): void {
      this.hasPermission = hasPermission
    }
  }
})
