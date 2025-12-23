/**
 * 通知状态管理（仅 App 端）
 * Notification state management (App only)
 */

import { defineStore } from 'pinia'
import type { NotificationState } from '@/types/store'
import type { Notification, PageParams } from '@/types/api'
import { notificationApi } from '@/api/notification'

/*** TabBar 消息图标索引 - Message tab index in tabBar ***/
const MESSAGE_TAB_INDEX = 3

/*** Notification store - manages push notifications (App only) ***/
export const useNotificationStore = defineStore('notification', {
  state: (): NotificationState => ({
    list: [],
    unreadCount: 0,
    hasPermission: false,
    loading: false,
    currentPage: 1,
    hasMore: true,
    total: 0
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
    async fetchList(refresh: boolean = false): Promise<void> {
      if (this.loading) return
      
      this.loading = true
      
      /*** Reset pagination on refresh ***/
      if (refresh) {
        this.currentPage = 1
        this.hasMore = true
      }
      
      try {
        const params: PageParams = {
          current: this.currentPage,
          size: 10
        }
        
        const res = await notificationApi.getList(params)
        
        if (res.code === 200 && res.data) {
          const { records, total } = res.data
          
          /*** Replace or append based on refresh flag ***/
          if (refresh) {
            this.list = records
          } else {
            this.list = [...this.list, ...records]
          }
          
          this.total = total
          this.hasMore = this.list.length < total
          this.currentPage++
          
          /*** Update unread count from list ***/
          this.updateUnreadCount()
        }
      } catch (error) {
        console.error('获取通知列表失败:', error)
        uni.showToast({
          title: '获取通知失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
    },

    /*** 加载更多通知 - Load more notifications ***/
    async loadMore(): Promise<void> {
      if (!this.hasMore || this.loading) return
      await this.fetchList(false)
    },

    /*** 标记为已读 - Mark as read ***/
    async markAsRead(id: number): Promise<boolean> {
      const notification = this.list.find(n => n.id === id)
      if (!notification || notification.isRead) {
        return true
      }
      
      try {
        const res = await notificationApi.markAsRead(id)
        
        if (res.code === 200) {
          /*** Update local state ***/
          notification.isRead = true
          this.updateUnreadCount()
          return true
        }
        return false
      } catch (error) {
        console.error('标记已读失败:', error)
        return false
      }
    },

    /*** 标记全部已读 - Mark all as read ***/
    async markAllAsRead(): Promise<boolean> {
      try {
        const res = await notificationApi.markAllAsRead()
        
        if (res.code === 200) {
          /*** Update all local notifications to read ***/
          this.list.forEach(n => {
            n.isRead = true
          })
          this.unreadCount = 0
          /*** Update tabBar badge ***/
          this.updateTabBarBadge()
          return true
        }
        return false
      } catch (error) {
        console.error('标记全部已读失败:', error)
        return false
      }
    },

    /*** 更新未读数量 - Update unread count from local list ***/
    updateUnreadCount(): void {
      this.unreadCount = this.list.filter(n => !n.isRead).length
      /*** Update tabBar badge ***/
      this.updateTabBarBadge()
    },

    /*** 获取未读数量 - Get unread count from server ***/
    async getUnreadCount(): Promise<number> {
      try {
        const res = await notificationApi.getUnreadCount()
        
        if (res.code === 200 && res.data !== undefined) {
          this.unreadCount = res.data
          /*** Update tabBar badge ***/
          this.updateTabBarBadge()
        }
      } catch (error) {
        console.error('获取未读数量失败:', error)
      }
      return this.unreadCount
    },

    /*** 更新 tabBar 红点/角标 - Update tabBar badge ***/
    updateTabBarBadge(): void {
      // #ifdef APP-PLUS
      if (this.unreadCount > 0) {
        /*** Show badge with count ***/
        uni.setTabBarBadge({
          index: MESSAGE_TAB_INDEX,
          text: this.unreadCount > 99 ? '99+' : String(this.unreadCount)
        })
      } else {
        /*** Remove badge when no unread messages ***/
        uni.removeTabBarBadge({
          index: MESSAGE_TAB_INDEX
        })
      }
      // #endif
    },

    /*** 显示红点（无数字）- Show red dot without number ***/
    showTabBarRedDot(): void {
      // #ifdef APP-PLUS
      if (this.unreadCount > 0) {
        uni.showTabBarRedDot({
          index: MESSAGE_TAB_INDEX
        })
      }
      // #endif
    },

    /*** 隐藏红点 - Hide red dot ***/
    hideTabBarRedDot(): void {
      // #ifdef APP-PLUS
      uni.hideTabBarRedDot({
        index: MESSAGE_TAB_INDEX
      })
      // #endif
    },

    /*** 设置推送权限状态 - Set push permission status ***/
    setPermission(hasPermission: boolean): void {
      this.hasPermission = hasPermission
    },

    /*** 清空通知列表 - Clear notification list ***/
    clearList(): void {
      this.list = []
      this.unreadCount = 0
      this.currentPage = 1
      this.hasMore = true
      this.total = 0
      /*** Remove tabBar badge ***/
      this.updateTabBarBadge()
    },

    /*** 根据通知类型获取跳转路径 - Get navigation path by notification type ***/
    getNavigationPath(notification: Notification): string {
      switch (notification.type) {
        case 'project':
          return notification.relatedId 
            ? `/pages/project/detail?id=${notification.relatedId}` 
            : '/pages/index/index'
        case 'approval':
          return notification.relatedId 
            ? `/pages/project/detail?id=${notification.relatedId}` 
            : '/pages/index/index'
        case 'system':
        default:
          return '/pages/index/index'
      }
    }
  }
})
