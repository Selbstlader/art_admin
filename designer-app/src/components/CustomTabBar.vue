<template>
  <!-- 
    自定义底部导航栏组件 - iOS 26 Liquid Glass风格
    Custom TabBar Component - iOS 26 Liquid Glass style with smooth slide animation
    App端显示5个tab，小程序端显示3个tab（隐藏AI对话和消息）
  -->
  <view class="tabbar-container">
    <!-- 背景模糊层 - Background blur layer -->
    <view class="tabbar-blur"></view>
    
    <!-- 玻璃主体 - Glass body -->
    <view class="custom-tabbar safe-area-bottom" :class="[`tab-count-${visibleTabs.length}`]">
      <!-- 玻璃高光 - Glass highlight -->
      <view class="glass-highlight"></view>
      <view class="glass-highlight-bottom"></view>
      
      <!-- 滑动指示器 - Sliding indicator (pure CSS) -->
      <view 
        class="slide-indicator" 
        :class="[`active-${activeIndex}`, `tab-count-${visibleTabs.length}`]"
      >
        <view class="indicator-glow"></view>
        <view class="indicator-pill"></view>
      </view>
      
      <!-- 导航项 - Navigation items -->
      <view 
        v-for="(item, index) in visibleTabs" 
        :key="item.pagePath"
        :class="['tabbar-item', { 'active': activeIndex === index }]"
        @click="handleTabClick(item, index)"
      >
        <!-- 图标容器 - Icon container -->
        <view class="icon-container">
          <!-- 首页图标 - Home icon -->
          <svg v-if="item.icon === 'home'" class="tabbar-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
            <polyline points="9 22 9 12 15 12 15 22"></polyline>
          </svg>
          
          <!-- 材料库图标 - Material icon -->
          <svg v-else-if="item.icon === 'material'" class="tabbar-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
            <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
            <line x1="12" y1="22.08" x2="12" y2="12"></line>
          </svg>
          
          <!-- AI对话图标 - Chat icon -->
          <svg v-else-if="item.icon === 'chat'" class="tabbar-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 8V4H8"></path>
            <rect width="16" height="12" x="4" y="8" rx="2"></rect>
            <path d="M2 14h2"></path>
            <path d="M20 14h2"></path>
            <path d="M15 13v2"></path>
            <path d="M9 13v2"></path>
          </svg>
          
          <!-- 消息图标 - Message icon -->
          <svg v-else-if="item.icon === 'message'" class="tabbar-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"></path>
            <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"></path>
          </svg>
          
          <!-- 我的图标 - Mine icon -->
          <svg v-else-if="item.icon === 'mine'" class="tabbar-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
          
          <!-- 未读消息红点 - Unread message badge -->
          <view v-if="item.icon === 'message' && unreadCount > 0" class="badge-dot"></view>
        </view>
        
        <!-- 文字标签 - Text label -->
        <text class="tabbar-text">{{ item.text }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * CustomTabBar Component
 * 自定义底部导航栏 - iOS 26 Liquid Glass风格，纯CSS流畅滑动动画
 * Custom bottom navigation bar with iOS 26 style and smooth CSS-only animation
 ***/
import { ref, computed, onMounted, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { platform } from '@/utils/platform'
import { useNotificationStore } from '@/store/notification'

/*** Tab 项配置接口 - Tab item configuration interface ***/
interface TabItem {
  pagePath: string
  text: string
  icon: string
  appOnly?: boolean  // 仅 App 端显示
}

/*** Component Props - 组件属性定义 ***/
interface Props {
  current?: number
}

const props = withDefaults(defineProps<Props>(), {
  current: 0
})

/*** Component Emits - 组件事件定义 ***/
const emit = defineEmits<{
  (e: 'change', index: number): void
}>()

/*** 获取通知 Store - Get notification store ***/
const notificationStore = useNotificationStore()

/*** 未读消息数量 - Unread message count ***/
const unreadCount = computed(() => notificationStore.unreadCount)

/*** 所有 Tab 项配置 - All tab items configuration ***/
const allTabs: TabItem[] = [
  { pagePath: '/pages/index/index', text: '首页', icon: 'home' },
  { pagePath: '/pages/material/index', text: '材料库', icon: 'material' },
  { pagePath: '/pages/chat/index', text: 'AI对话', icon: 'chat', appOnly: true },
  { pagePath: '/pages/message/index', text: '消息', icon: 'message', appOnly: true },
  { pagePath: '/pages/mine/index', text: '我的', icon: 'mine' }
]

/*** 可见的 Tab 项 - Visible tab items based on platform ***/
const visibleTabs = computed((): TabItem[] => {
  const isApp = platform.isApp()
  if (isApp) {
    return allTabs
  } else {
    return allTabs.filter(tab => !tab.appOnly)
  }
})

/*** 根据路径获取在可见tabs中的索引 - Get visible index from path ***/
const getVisibleIndexFromPath = (path: string): number => {
  const index = visibleTabs.value.findIndex(tab => tab.pagePath === path)
  return index !== -1 ? index : 0
}

/*** 根据当前路由获取索引 - Get index from current route ***/
const getActiveIndexFromRoute = (): number => {
  const pages = getCurrentPages()
  if (pages.length > 0) {
    const currentPage = pages[pages.length - 1]
    const currentPath = '/' + currentPage.route
    return getVisibleIndexFromPath(currentPath)
  }
  return 0
}

/*** 当前选中索引 - Current active index ***/
const activeIndex = ref(0)

/*** 是否已初始化 - Whether initialized ***/
const isInitialized = ref(false)

/*** 初始化索引 - Initialize index ***/
const initIndex = (): void => {
  // 优先使用路由判断，确保准确
  activeIndex.value = getActiveIndexFromRoute()
  isInitialized.value = true
}

/*** 处理 Tab 点击 - Handle tab click ***/
const handleTabClick = (item: TabItem, index: number): void => {
  /*** 如果点击当前已选中的tab，不做任何操作 ***/
  if (activeIndex.value === index) return
  
  /*** 立即更新索引触发动画 - Immediately update index to trigger animation ***/
  activeIndex.value = index
  
  /*** 延迟跳转页面，让动画先执行 - Delay navigation to let animation play first ***/
  setTimeout(() => {
    uni.switchTab({
      url: item.pagePath,
      fail: () => {
        uni.navigateTo({ url: item.pagePath })
      }
    })
  }, 50)
}

onMounted(() => {
  /*** 使用 nextTick 确保路由已更新 - Use nextTick to ensure route is updated ***/
  setTimeout(() => {
    initIndex()
  }, 0)
  
  // #ifdef APP-PLUS
  notificationStore.getUnreadCount()
  // #endif
})

/*** 页面显示时同步索引 - Sync index when page shows ***/
onShow(() => {
  /*** 每次页面显示时重新从路由获取索引 ***/
  const newIndex = getActiveIndexFromRoute()
  if (activeIndex.value !== newIndex) {
    activeIndex.value = newIndex
  }
})
</script>

<style lang="scss" scoped>
/*** TabBar 外层容器 - TabBar outer container ***/
.tabbar-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 999;
  padding: 0 24rpx 16rpx;
}

/*** 背景模糊层 - Background blur layer for iOS 26 effect ***/
.tabbar-blur {
  position: absolute;
  top: -20rpx;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, transparent 0%, rgba(255, 255, 255, 0.6) 30%);
  pointer-events: none;
}

/*** 自定义 TabBar - iOS 26 Liquid Glass style ***/
.custom-tabbar {
  position: relative;
  height: 120rpx;
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.75) 0%,
    rgba(255, 255, 255, 0.55) 20%,
    rgba(240, 245, 255, 0.65) 100%
  );
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border-radius: 32rpx;
  border: 1.5rpx solid rgba(255, 255, 255, 0.6);
  box-shadow: 
    0 8rpx 32rpx rgba(0, 0, 0, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04),
    inset 0 1rpx 1rpx rgba(255, 255, 255, 0.9),
    inset 0 -1rpx 1rpx rgba(0, 0, 0, 0.02);
  overflow: hidden;
}

/*** 玻璃顶部高光 - Glass top highlight ***/
.glass-highlight {
  position: absolute;
  top: 0;
  left: 20rpx;
  right: 20rpx;
  height: 1rpx;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.9) 20%,
    rgba(255, 255, 255, 1) 50%,
    rgba(255, 255, 255, 0.9) 80%,
    transparent 100%
  );
  opacity: 0.8;
}

/*** 玻璃底部微光 - Glass bottom subtle glow ***/
.glass-highlight-bottom {
  position: absolute;
  bottom: 0;
  left: 40rpx;
  right: 40rpx;
  height: 1rpx;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(200, 220, 255, 0.3) 30%,
    rgba(200, 220, 255, 0.4) 50%,
    rgba(200, 220, 255, 0.3) 70%,
    transparent 100%
  );
}

/*** 滑动指示器 - Sliding indicator (纯CSS定位) ***/
.slide-indicator {
  position: absolute;
  top: 50%;
  height: 88rpx;
  pointer-events: none;
  z-index: 0;
  /*** 流畅弹性动画 - Smooth spring animation ***/
  transition: left 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateY(-50%);
}

/*** 3个Tab时的指示器位置 - Indicator positions for 3 tabs ***/
.slide-indicator.tab-count-3 {
  width: calc(100% / 3);
  
  &.active-0 { left: 0; }
  &.active-1 { left: calc(100% / 3); }
  &.active-2 { left: calc(100% / 3 * 2); }
}

/*** 5个Tab时的指示器位置 - Indicator positions for 5 tabs ***/
.slide-indicator.tab-count-5 {
  width: calc(100% / 5);
  
  &.active-0 { left: 0; }
  &.active-1 { left: calc(100% / 5); }
  &.active-2 { left: calc(100% / 5 * 2); }
  &.active-3 { left: calc(100% / 5 * 3); }
  &.active-4 { left: calc(100% / 5 * 4); }
}

/*** 指示器光晕 - Indicator glow ***/
.indicator-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  height: 100%;
  background: radial-gradient(
    ellipse at center,
    rgba(59, 130, 246, 0.2) 0%,
    rgba(59, 130, 246, 0.08) 50%,
    transparent 80%
  );
  border-radius: 50%;
  filter: blur(4rpx);
}

/*** 指示器胶囊 - Indicator pill ***/
.indicator-pill {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 70%;
  height: 72rpx;
  background: linear-gradient(
    135deg,
    rgba(59, 130, 246, 0.15) 0%,
    rgba(37, 99, 235, 0.1) 50%,
    rgba(59, 130, 246, 0.12) 100%
  );
  border-radius: 20rpx;
  border: 1rpx solid rgba(59, 130, 246, 0.15);
  box-shadow:
    0 4rpx 12rpx rgba(59, 130, 246, 0.12),
    inset 0 1rpx 2rpx rgba(255, 255, 255, 0.5);
}

/*** TabBar 项 - TabBar item ***/
.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 12rpx 0;
  position: relative;
  z-index: 1;
  /*** 快速响应的过渡 - Fast responsive transition ***/
  transition: transform 0.2s ease-out;
  
  &:active {
    transform: scale(0.9);
  }
}

/*** 图标容器 - Icon container ***/
.icon-container {
  position: relative;
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4rpx;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/*** TabBar 图标 - TabBar icon ***/
.tabbar-icon {
  width: 48rpx;
  height: 48rpx;
  color: rgba(100, 116, 139, 0.7);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/*** TabBar 文字 - TabBar text ***/
.tabbar-text {
  font-size: 22rpx;
  color: rgba(100, 116, 139, 0.7);
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  letter-spacing: 0.5rpx;
}

/*** 选中状态 - Active state ***/
.tabbar-item.active {
  .icon-container {
    transform: translateY(-4rpx) scale(1.05);
  }
  
  .tabbar-icon {
    color: #2563EB;
    filter: drop-shadow(0 2rpx 6rpx rgba(37, 99, 235, 0.35));
  }
  
  .tabbar-text {
    color: #2563EB;
    font-weight: 600;
  }
}

/*** 未读消息红点 - Unread message badge ***/
.badge-dot {
  position: absolute;
  top: 0;
  right: 0;
  width: 16rpx;
  height: 16rpx;
  background: linear-gradient(135deg, #F87171 0%, #EF4444 100%);
  border-radius: 50%;
  border: 2rpx solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 2rpx 6rpx rgba(239, 68, 68, 0.4);
  animation: badge-pulse 2s ease-in-out infinite;
}

@keyframes badge-pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.15); opacity: 0.9; }
}
</style>
