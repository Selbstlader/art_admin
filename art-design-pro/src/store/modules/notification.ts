/***
 * Notification Store - 通知状态管理
 * Manages notification state for the application
 ***/
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  getNotificationList,
  getUnreadCount,
  markNotificationRead,
  markAllNotificationsRead,
  type NotificationItem,
  type NotificationRelatedType
} from '@/api/notification'

export type NoticeType = 'email' | 'message' | 'collection' | 'user' | 'notice'

export interface NotificationDisplayItem {
  id?: number
  title: string
  time: string
  type: NoticeType
  read?: boolean
  content?: string
  relatedId?: number
  relatedType?: NotificationRelatedType
}

export const useNotificationStore = defineStore('notificationStore', () => {
  /*** 通知列表 ***/
  const notifications = ref<NotificationDisplayItem[]>([])

  /*** 未读通知数量 ***/
  const unreadCount = ref(0)

  /*** 加载状态 ***/
  const loading = ref(false)

  /*** 从接口加载通知列表 ***/
  async function fetchNotifications() {
    loading.value = true
    try {
      const res = (await getNotificationList({ current: 1, size: 50 })) as any
      if (res.data?.records) {
        notifications.value = res.data.records.map((item: NotificationItem) => ({
          id: item.id,
          title: item.title,
          time: formatTime(item.createdAt),
          type: item.type as NoticeType,
          read: item.isRead,
          content: item.content,
          relatedId: item.relatedId,
          relatedType: item.relatedType
        }))
      }
    } catch (error) {
      console.error('获取通知列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  /*** 从接口获取未读数量 ***/
  async function fetchUnreadCount() {
    try {
      const res = (await getUnreadCount()) as any
      if (res.data) {
        unreadCount.value = res.data.count || 0
      }
    } catch (error) {
      console.error('获取未读数量失败:', error)
    }
  }

  /*** 标记为已读 ***/
  async function markAsRead(id: number) {
    try {
      await markNotificationRead(id)
      const item = notifications.value.find((n) => n.id === id)
      if (item) {
        item.read = true
      }
      // 更新未读数量
      await fetchUnreadCount()
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }

  /*** 标记全部已读 ***/
  async function markAllAsRead() {
    try {
      await markAllNotificationsRead()
      notifications.value.forEach((n) => (n.read = true))
      unreadCount.value = 0
    } catch (error) {
      console.error('标记全部已读失败:', error)
    }
  }

  /*** 刷新通知数据 ***/
  async function refresh() {
    await Promise.all([fetchNotifications(), fetchUnreadCount()])
  }

  /*** 格式化时间 ***/
  function formatTime(dateStr: string): string {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)

    if (minutes < 1) return '刚刚'
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    return date.toLocaleString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
  }

  return {
    notifications,
    unreadCount,
    loading,
    fetchNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    refresh
  }
})
