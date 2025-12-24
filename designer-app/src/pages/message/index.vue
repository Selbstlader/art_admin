<template>
  <!-- #ifdef APP-PLUS -->
  <view class="message-page">
    <!-- 页面头部 - Page header -->
    <view class="page-header">
      <view class="header-content">
        <text class="header-title">消息中心</text>
        <view 
          v-if="notificationStore.hasUnread" 
          class="mark-all-btn"
          @click="handleMarkAllRead"
        >
          <text class="mark-all-text">全部已读</text>
        </view>
      </view>
    </view>

    <!-- 通知列表 - Notification list -->
    <scroll-view
      class="notification-list"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="handleRefresh"
      @scrolltolower="handleLoadMore"
    >
      <!-- 加载中状态 - Loading state -->
      <view v-if="notificationStore.loading && notificationStore.list.length === 0" class="loading-container">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 通知列表内容 - Notification list content -->
      <view v-else-if="notificationStore.sortedList.length > 0" class="list-content">
        <view
          v-for="notification in notificationStore.sortedList"
          :key="notification.id"
          class="notification-item float-card"
          :class="{ 'is-read': notification.isRead }"
          @click="handleNotificationClick(notification)"
        >
          <!-- 通知图标 - Notification icon -->
          <view class="notification-icon" :class="getIconClass(notification.type)">
            <text class="icon-text">{{ getIconText(notification.type) }}</text>
          </view>

          <!-- 通知内容 - Notification content -->
          <view class="notification-content">
            <view class="notification-header">
              <text class="notification-title">{{ notification.title }}</text>
              <view v-if="!notification.isRead" class="unread-dot"></view>
            </view>
            <text class="notification-body">{{ notification.content }}</text>
            <text class="notification-time">{{ formatTime(notification.createdAt) }}</text>
          </view>

          <!-- 箭头指示 - Arrow indicator -->
          <view class="notification-arrow">
            <text class="arrow-icon">›</text>
          </view>
        </view>

        <!-- 加载更多 - Load more -->
        <view v-if="notificationStore.hasMore" class="load-more">
          <text v-if="notificationStore.loading" class="load-more-text">加载中...</text>
          <text v-else class="load-more-text">上拉加载更多</text>
        </view>

        <!-- 没有更多 - No more -->
        <view v-else class="no-more">
          <text class="no-more-text">没有更多消息了</text>
        </view>
      </view>

      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <view class="empty-icon">📭</view>
        <text class="empty-text">暂无消息</text>
        <text class="empty-hint">项目更新、审批提醒等消息将在这里显示</text>
      </view>
    </scroll-view>
  </view>
  <!-- #endif -->

  <!-- #ifdef MP-WEIXIN -->
  <view class="not-supported-page">
    <view class="not-supported-content">
      <view class="not-supported-icon">📱</view>
      <text class="not-supported-title">功能暂不可用</text>
      <text class="not-supported-desc">消息通知功能仅在 App 端可用</text>
      <text class="not-supported-hint">请下载 App 体验完整功能</text>
    </view>
    
    <!-- 自定义底部导航栏 - Custom TabBar (小程序端不显示此页面，但保留组件) -->
    <CustomTabBar />
  </view>
  <!-- #endif -->
</template>

<script setup lang="ts">
/*** Message center page - displays notifications (App only) ***/
import { ref, onMounted } from 'vue'
import { onShow, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { useNotificationStore } from '@/store/notification'
import type { Notification } from '@/types/api'
import CustomTabBar from '@/components/CustomTabBar.vue'

const notificationStore = useNotificationStore()
const isRefreshing = ref(false)

/*** 初始化加载 - Initial load ***/
onMounted(async () => {
  // #ifdef APP-PLUS
  await notificationStore.fetchList(true)
  // #endif
})

/*** 页面显示时刷新 - Refresh on page show ***/
onShow(() => {
  // #ifdef APP-PLUS
  notificationStore.getUnreadCount()
  // #endif
})

/*** 下拉刷新 - Pull down refresh ***/
onPullDownRefresh(async () => {
  // #ifdef APP-PLUS
  await notificationStore.fetchList(true)
  uni.stopPullDownRefresh()
  // #endif
})

/*** 触底加载更多 - Load more on reach bottom ***/
onReachBottom(() => {
  // #ifdef APP-PLUS
  notificationStore.loadMore()
  // #endif
})

/*** 处理下拉刷新 - Handle refresh ***/
const handleRefresh = async () => {
  isRefreshing.value = true
  await notificationStore.fetchList(true)
  isRefreshing.value = false
}

/*** 处理加载更多 - Handle load more ***/
const handleLoadMore = () => {
  if (!notificationStore.loading && notificationStore.hasMore) {
    notificationStore.loadMore()
  }
}

/*** 处理通知点击 - Handle notification click ***/
const handleNotificationClick = async (notification: Notification) => {
  /*** Mark as read first ***/
  await notificationStore.markAsRead(notification.id)
  
  /*** Navigate to related page ***/
  const path = notificationStore.getNavigationPath(notification)
  uni.navigateTo({
    url: path,
    fail: () => {
      /*** Fallback to switchTab if navigateTo fails ***/
      uni.switchTab({
        url: '/pages/index/index'
      })
    }
  })
}

/*** 处理全部已读 - Handle mark all as read ***/
const handleMarkAllRead = async () => {
  uni.showModal({
    title: '提示',
    content: '确定将所有消息标记为已读吗？',
    success: async (res) => {
      if (res.confirm) {
        const success = await notificationStore.markAllAsRead()
        if (success) {
          uni.showToast({
            title: '已全部标记为已读',
            icon: 'success'
          })
        }
      }
    }
  })
}

/*** 获取图标样式类 - Get icon class by type ***/
const getIconClass = (type: string): string => {
  switch (type) {
    case 'project':
      return 'icon-project'
    case 'approval':
      return 'icon-approval'
    case 'system':
    default:
      return 'icon-system'
  }
}

/*** 获取图标文字 - Get icon text by type ***/
const getIconText = (type: string): string => {
  switch (type) {
    case 'project':
      return '📋'
    case 'approval':
      return '✅'
    case 'system':
    default:
      return '🔔'
  }
}

/*** 格式化时间 - Format time ***/
const formatTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
}
</script>

<style lang="scss" scoped>
/*** 消息页面容器 - Message page container ***/
.message-page {
  min-height: 100vh;
  background: #F5F7FA;
  display: flex;
  flex-direction: column;
}

/*** 页面头部 - Page header ***/
.page-header {
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  padding: 32rpx;
  padding-top: calc(32rpx + var(--status-bar-height, 0px));
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #FFFFFF;
}

.mark-all-btn {
  padding: 12rpx 24rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 24rpx;
  transition: all 0.2s ease;
  
  &:active {
    opacity: 0.8;
    transform: scale(0.98);
  }
}

.mark-all-text {
  font-size: 24rpx;
  color: #FFFFFF;
}

/*** 通知列表 - Notification list ***/
.notification-list {
  flex: 1;
  height: calc(100vh - 120rpx - var(--status-bar-height, 0px));
}

.list-content {
  padding: 24rpx;
}

/*** 通知项 - Notification item ***/
.notification-item {
  display: flex;
  align-items: flex-start;
  padding: 32rpx;
  margin-bottom: 24rpx;
  background: #FFFFFF;
  border-radius: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  }
  
  &.is-read {
    opacity: 0.7;
    
    .notification-title {
      font-weight: 400;
    }
  }
}

/*** 通知图标 - Notification icon ***/
.notification-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
  flex-shrink: 0;
  
  &.icon-project {
    background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  }
  
  &.icon-approval {
    background: linear-gradient(135deg, #10B981 0%, #059669 100%);
  }
  
  &.icon-system {
    background: linear-gradient(135deg, #F59E0B 0%, #D97706 100%);
  }
}

.icon-text {
  font-size: 36rpx;
}

/*** 通知内容 - Notification content ***/
.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-header {
  display: flex;
  align-items: center;
  margin-bottom: 8rpx;
}

.notification-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #2C3E50;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.unread-dot {
  width: 16rpx;
  height: 16rpx;
  background: #EF4444;
  border-radius: 50%;
  margin-left: 12rpx;
  flex-shrink: 0;
}

.notification-body {
  font-size: 26rpx;
  color: #64748B;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 12rpx;
}

.notification-time {
  font-size: 22rpx;
  color: #94A3B8;
}

/*** 箭头指示 - Arrow indicator ***/
.notification-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 16rpx;
}

.arrow-icon {
  font-size: 32rpx;
  color: #CBD5E1;
}

/*** 加载状态 - Loading state ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
}

.loading-spinner {
  width: 60rpx;
  height: 60rpx;
  border: 4rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 24rpx;
  font-size: 26rpx;
  color: #64748B;
}

/*** 加载更多 - Load more ***/
.load-more,
.no-more {
  display: flex;
  justify-content: center;
  padding: 32rpx 0;
}

.load-more-text,
.no-more-text {
  font-size: 24rpx;
  color: #94A3B8;
}

/*** 空状态 - Empty state ***/
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 48rpx;
}

.empty-icon {
  font-size: 120rpx;
  margin-bottom: 32rpx;
}

.empty-text {
  font-size: 32rpx;
  font-weight: 500;
  color: #2C3E50;
  margin-bottom: 16rpx;
}

.empty-hint {
  font-size: 26rpx;
  color: #94A3B8;
  text-align: center;
}

/*** 不支持页面 - Not supported page ***/
.not-supported-page {
  min-height: 100vh;
  background: #F5F7FA;
  display: flex;
  align-items: center;
  justify-content: center;
}

.not-supported-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48rpx;
}

.not-supported-icon {
  font-size: 120rpx;
  margin-bottom: 32rpx;
}

.not-supported-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #2C3E50;
  margin-bottom: 16rpx;
}

.not-supported-desc {
  font-size: 28rpx;
  color: #64748B;
  margin-bottom: 12rpx;
}

.not-supported-hint {
  font-size: 24rpx;
  color: #94A3B8;
}
</style>
