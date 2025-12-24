<script setup lang="ts">
/*** Root App component - handles app lifecycle events ***/
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useNotificationStore } from '@/store/notification'
import { useNetworkStore } from '@/store/network'
import { pushManager, subscribeManager } from '@/utils/push'
import { platform } from '@/utils/platform'

/*** 初始化网络状态监听 - Initialize network status listener ***/
const initNetworkListener = () => {
  const networkStore = useNetworkStore()
  /*** Initialize network listener on app launch ***/
  networkStore.initNetworkListener()
}

/*** 初始化推送通知 - Initialize push notifications ***/
const initPushNotification = async () => {
  const notificationStore = useNotificationStore()
  
  // #ifdef APP-PLUS
  /*** Request push permission on App ***/
  const hasPermission = await pushManager.requestPermission()
  notificationStore.setPermission(hasPermission)
  
  if (hasPermission) {
    /*** Listen to push messages ***/
    pushManager.onPushMessage((message) => {
      console.log('收到推送消息:', message)
      
      if (message.type === 'click') {
        /*** Handle push message click - navigate to related page ***/
        handlePushClick(message.payload)
      } else if (message.type === 'receive') {
        /*** Handle push message receive - refresh unread count ***/
        notificationStore.getUnreadCount()
      }
    })
  }
  // #endif
  
  // #ifdef MP-WEIXIN
  /*** Request subscribe message permission on Mini Program ***/
  /*** Note: Subscribe message templates need to be configured in WeChat backend ***/
  const subscribeResult = await subscribeManager.requestSubscribe([
    // TODO: 添加实际的订阅消息模板 ID
    // 'template_id_1',
    // 'template_id_2'
  ])
  notificationStore.setPermission(subscribeResult.success)
  // #endif
}

/*** 处理推送消息点击 - Handle push message click ***/
const handlePushClick = (payload: any) => {
  if (!payload) return
  
  try {
    const data = typeof payload === 'string' ? JSON.parse(payload) : payload
    
    /*** Navigate based on notification type ***/
    if (data.type === 'project' && data.projectId) {
      uni.navigateTo({
        url: `/pages/project/detail?id=${data.projectId}`
      })
    } else if (data.type === 'approval' && data.projectId) {
      uni.navigateTo({
        url: `/pages/project/detail?id=${data.projectId}`
      })
    } else {
      /*** Default: go to message center ***/
      uni.switchTab({
        url: '/pages/message/index'
      })
    }
  } catch (error) {
    console.error('解析推送消息失败:', error)
    /*** Fallback: go to message center ***/
    uni.switchTab({
      url: '/pages/message/index'
    })
  }
}

/*** 初始化未读消息数量 - Initialize unread message count ***/
const initUnreadCount = async () => {
  // #ifdef APP-PLUS
  const notificationStore = useNotificationStore()
  /*** getUnreadCount will automatically update tabBar badge ***/
  await notificationStore.getUnreadCount()
  // #endif
}

onLaunch(() => {
  console.log('App Launch')
  
  /*** Initialize network listener on app launch ***/
  initNetworkListener()
  
  /*** Initialize push notification on app launch ***/
  initPushNotification()
})

onShow(() => {
  console.log('App Show')
  
  /*** Check network status when app becomes visible ***/
  const networkStore = useNetworkStore()
  networkStore.checkNetworkStatus()
  
  /*** Refresh unread count when app becomes visible ***/
  initUnreadCount()
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
/**
 * 全局样式 - Global styles
 * C4D风格高级UI设计系统
 * C4D style premium UI design system
 */

/* 引入全局通用样式 - Import global common styles */
@import '@/styles/global.scss';

/* ============================================
   页面基础样式 - Page Base Styles
   ============================================ */

page {
  background-color: #F5F7FA;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 28rpx;
  color: #1E293B;
  box-sizing: border-box;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* ============================================
   蓝色渐变主题色 - Blue Gradient Theme
   ============================================ */

/*** 主题色渐变背景 - Primary gradient background ***/
.gradient-bg {
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
}

.gradient-blue {
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
}

.gradient-blue-light {
  background: linear-gradient(135deg, #60A5FA 0%, #3B82F6 100%);
}

.gradient-blue-vertical {
  background: linear-gradient(180deg, #60A5FA 0%, #3B82F6 100%);
}

/*** 页面渐变背景 - Page gradient background ***/
.gradient-page {
  background: linear-gradient(180deg, #F0F7FF 0%, #F5F7FA 100%);
}

/* ============================================
   玻璃拟态样式 - Glass Morphism Styles
   ============================================ */

/*** 玻璃拟态卡片 - Glass morphism card ***/
.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 2rpx solid rgba(255, 255, 255, 0.3);
  border-radius: 24rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(0, 0, 0, 0.06),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

/*** 玻璃拟态卡片 - 浅色 - Light glass card ***/
.glass-card-light {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-radius: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

/*** 玻璃拟态按钮 - Glass morphism button ***/
.glass-btn {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 2rpx solid rgba(255, 255, 255, 0.3);
  border-radius: 20rpx;
  box-shadow: 
    0 4rpx 16rpx rgba(0, 0, 0, 0.06),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
  transition: all 0.2s ease;
  
  &:active {
    opacity: 0.8;
    transform: scale(0.98);
  }
}

/* ============================================
   悬浮卡片样式 - Floating Card Styles
   ============================================ */

/*** 悬浮卡片 - Floating card ***/
.float-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 28rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}

.float-card:active {
  transform: scale(0.98);
  box-shadow: 
    0 4rpx 16rpx rgba(59, 130, 246, 0.06),
    0 2rpx 6rpx rgba(0, 0, 0, 0.03);
}

/*** 悬浮卡片 - 带3D阴影 - Floating card with 3D shadow ***/
.float-card-3d {
  position: relative;
  background: #FFFFFF;
  border-radius: 28rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  
  &::after {
    content: '';
    position: absolute;
    bottom: -8rpx;
    left: 24rpx;
    right: 24rpx;
    height: 16rpx;
    background: radial-gradient(ellipse, rgba(0, 0, 0, 0.08) 0%, transparent 70%);
    border-radius: 50%;
    z-index: -1;
  }
}

/* ============================================
   按钮样式 - Button Styles
   ============================================ */

/*** 主要按钮 - Primary button ***/
.btn-primary {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  color: #FFFFFF;
  border: none;
  border-radius: 24rpx;
  padding: 24rpx 48rpx;
  font-size: 32rpx;
  font-weight: 600;
  box-shadow: 
    0 8rpx 24rpx rgba(59, 130, 246, 0.25),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.2);
  transition: all 0.2s ease;
}

.btn-primary:active {
  opacity: 0.8;
  transform: scale(0.98);
}

/*** 次要按钮 - Secondary button ***/
.btn-secondary {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F1F5F9;
  color: #64748B;
  border: none;
  border-radius: 24rpx;
  padding: 24rpx 48rpx;
  font-size: 32rpx;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-secondary:active {
  opacity: 0.8;
  transform: scale(0.98);
  background: #E2E8F0;
}

/*** 轮廓按钮 - Outline button ***/
.btn-outline {
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  color: #3B82F6;
  border: 2rpx solid #3B82F6;
  border-radius: 24rpx;
  padding: 22rpx 46rpx;
  font-size: 32rpx;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-outline:active {
  opacity: 0.8;
  transform: scale(0.98);
  background: rgba(59, 130, 246, 0.05);
}

/* ============================================
   文本样式 - Text Styles
   ============================================ */

.text-primary {
  color: #3B82F6;
}

.text-secondary {
  color: #64748B;
}

.text-tertiary {
  color: #94A3B8;
}

.text-danger {
  color: #EF4444;
}

.text-success {
  color: #10B981;
}

.text-warning {
  color: #F59E0B;
}

.text-dark {
  color: #1E293B;
}

/*** 渐变文本 - Gradient text ***/
.text-gradient {
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ============================================
   间距工具类 - Spacing Utilities
   ============================================ */

.p-8 { padding: 16rpx; }
.p-16 { padding: 32rpx; }
.p-24 { padding: 48rpx; }
.p-32 { padding: 64rpx; }

.px-8 { padding-left: 16rpx; padding-right: 16rpx; }
.px-16 { padding-left: 32rpx; padding-right: 32rpx; }
.px-24 { padding-left: 48rpx; padding-right: 48rpx; }

.py-8 { padding-top: 16rpx; padding-bottom: 16rpx; }
.py-16 { padding-top: 32rpx; padding-bottom: 32rpx; }
.py-24 { padding-top: 48rpx; padding-bottom: 48rpx; }

.m-8 { margin: 16rpx; }
.m-16 { margin: 32rpx; }
.m-24 { margin: 48rpx; }

.mx-8 { margin-left: 16rpx; margin-right: 16rpx; }
.mx-16 { margin-left: 32rpx; margin-right: 32rpx; }
.mx-24 { margin-left: 48rpx; margin-right: 48rpx; }

.my-8 { margin-top: 16rpx; margin-bottom: 16rpx; }
.my-16 { margin-top: 32rpx; margin-bottom: 32rpx; }
.my-24 { margin-top: 48rpx; margin-bottom: 48rpx; }

.mb-8 { margin-bottom: 16rpx; }
.mb-16 { margin-bottom: 32rpx; }
.mb-24 { margin-bottom: 48rpx; }

.mt-8 { margin-top: 16rpx; }
.mt-16 { margin-top: 32rpx; }
.mt-24 { margin-top: 48rpx; }

.ml-8 { margin-left: 16rpx; }
.ml-16 { margin-left: 32rpx; }

.mr-8 { margin-right: 16rpx; }
.mr-16 { margin-right: 32rpx; }

/* ============================================
   Flex 布局工具类 - Flex Layout Utilities
   ============================================ */

.flex { display: flex; }
.flex-col { flex-direction: column; }
.flex-row { flex-direction: row; }
.items-center { align-items: center; }
.items-start { align-items: flex-start; }
.items-end { align-items: flex-end; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.justify-around { justify-content: space-around; }
.justify-start { justify-content: flex-start; }
.justify-end { justify-content: flex-end; }
.flex-1 { flex: 1; }
.flex-wrap { flex-wrap: wrap; }
.flex-shrink-0 { flex-shrink: 0; }

/* ============================================
   圆角工具类 - Border Radius Utilities
   ============================================ */

.rounded-sm { border-radius: 8rpx; }
.rounded-md { border-radius: 12rpx; }
.rounded-lg { border-radius: 16rpx; }
.rounded-xl { border-radius: 24rpx; }
.rounded-2xl { border-radius: 28rpx; }
.rounded-full { border-radius: 9999rpx; }

/* ============================================
   阴影工具类 - Shadow Utilities
   ============================================ */

.shadow-sm { box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08); }
.shadow-md { box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.1); }
.shadow-lg { box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.12); }
.shadow-xl { box-shadow: 0 12rpx 32rpx rgba(0, 0, 0, 0.15); }
.shadow-primary { box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.2); }
.shadow-card { 
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

/* ============================================
   安全区域适配 - Safe Area Adaptation
   ============================================ */

.safe-area-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.safe-area-top {
  padding-top: constant(safe-area-inset-top);
  padding-top: env(safe-area-inset-top);
}

.safe-area-inset {
  padding-top: constant(safe-area-inset-top);
  padding-top: env(safe-area-inset-top);
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

/* ============================================
   动画工具类 - Animation Utilities
   ============================================ */

.transition-all {
  transition: all 0.2s ease;
}

.transition-fast {
  transition: all 0.15s ease;
}

.transition-slow {
  transition: all 0.3s ease;
}

/* ============================================
   其他工具类 - Other Utilities
   ============================================ */

.overflow-hidden { overflow: hidden; }
.overflow-auto { overflow: auto; }
.text-center { text-align: center; }
.text-left { text-align: left; }
.text-right { text-align: right; }
.font-medium { font-weight: 500; }
.font-semibold { font-weight: 600; }
.font-bold { font-weight: 700; }
.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.line-clamp-2 {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
