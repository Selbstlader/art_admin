<template>
  <!-- 
    首页项目列表 - C4D风格高级UI
    Home page project list - C4D style premium UI
  -->
  <view class="page-container">
    <!-- 顶部区域 - Header area with blue gradient -->
    <view class="header-area">
      <view class="header-bg"></view>
      <view class="header-content">
        <text class="header-title">项目列表</text>
        <text class="header-subtitle">管理您的设计项目</text>
      </view>
    </view>

    <!-- 搜索框 - Search bar with glass morphism effect -->
    <view class="search-section">
      <view class="search-wrapper">
        <view class="search-icon-wrapper">
          <!-- Search icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </svg>
        </view>
        <input
          class="search-input"
          type="text"
          placeholder="搜索项目名称"
          v-model="searchKeyword"
          @input="handleSearch"
          @confirm="handleSearchConfirm"
        />
        <view v-if="searchKeyword" class="clear-btn" @click="clearSearch">
          <!-- Clear icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="m15 9-6 6"></path>
            <path d="m9 9 6 6"></path>
          </svg>
        </view>
      </view>
    </view>

    <!-- 项目列表 - Project list with scroll view -->
    <scroll-view
      class="project-list"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="handleRefresh"
      @scrolltolower="handleLoadMore"
    >
      <!-- 加载中状态 - Loading state -->
      <view v-if="projectStore.loading && projectStore.list.length === 0" class="loading-container">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 项目卡片列表 - Project cards with 3D effect -->
      <view v-else-if="displayProjects.length > 0" class="project-cards">
        <ProjectCard
          v-for="project in displayProjects"
          :key="project.id"
          :project="project"
          @click="goToDetail"
        />
      </view>

      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <view class="empty-icon">
          <!-- Folder icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z"></path>
          </svg>
        </view>
        <text class="empty-title">{{ searchKeyword ? '未找到匹配项目' : '暂无项目' }}</text>
        <text class="empty-desc">{{ searchKeyword ? '请尝试其他关键字' : '您还没有创建任何项目' }}</text>
      </view>

      <!-- 加载更多 - Load more indicator -->
      <view v-if="projectStore.hasMore && displayProjects.length > 0" class="load-more">
        <view v-if="loadingMore" class="loading-spinner-small"></view>
        <text class="load-more-text">{{ loadingMore ? '加载中...' : '上拉加载更多' }}</text>
      </view>
      
      <!-- 没有更多数据 - No more data -->
      <view v-else-if="!projectStore.hasMore && displayProjects.length > 0" class="no-more">
        <text class="no-more-text">— 已加载全部项目 —</text>
      </view>
      
      <!-- 底部安全区 - Bottom safe area -->
      <view class="safe-area-bottom"></view>
    </scroll-view>

    <!-- 离线模式提示 - Offline mode tip -->
    <view v-if="isOffline" class="offline-tip">
      <view class="offline-icon">
        <!-- Wifi off icon -->
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12.01 21.49 23.64 7c-.45-.34-4.93-4-11.64-4C5.28 3 .81 6.66.36 7l11.63 14.49.01.01.01-.01z"></path>
          <line x1="2" x2="22" y1="2" y2="22"></line>
        </svg>
      </view>
      <text class="offline-text">当前为离线模式</text>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Project list page - C4D style premium UI with 3D cards
 * 首页项目列表页面 - C4D风格高级UI，3D立体卡片设计
 ***/
import { ref, computed, onMounted, watch } from 'vue'
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import { useProjectStore } from '@/store/project'
import { ProjectCard } from '@/components'

/*** Store instance - 状态管理实例 ***/
const projectStore = useProjectStore()

/*** Local state - 本地状态 ***/
const searchKeyword = ref('')
const isRefreshing = ref(false)
const loadingMore = ref(false)
const isOffline = ref(false)
let searchTimer: ReturnType<typeof setTimeout> | null = null

/*** 显示的项目列表 - Display projects (filtered or all) ***/
const displayProjects = computed(() => {
  if (searchKeyword.value) {
    return projectStore.filteredList
  }
  return projectStore.list
})

/*** 检查网络状态 - Check network status ***/
const checkNetworkStatus = (): void => {
  uni.getNetworkType({
    success: (res) => {
      isOffline.value = res.networkType === 'none'
    }
  })
}

/*** 监听网络状态变化 - Listen to network status changes ***/
const setupNetworkListener = (): void => {
  uni.onNetworkStatusChange((res) => {
    const wasOffline = isOffline.value
    isOffline.value = !res.isConnected
    
    // 网络恢复时自动刷新 - Auto refresh when network recovers
    if (wasOffline && res.isConnected) {
      uni.showToast({ title: '网络已恢复', icon: 'none' })
      handleRefresh()
    }
  })
}

/*** 处理搜索输入 - Handle search input with debounce ***/
const handleSearch = (): void => {
  // 清除之前的定时器 - Clear previous timer
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  // 设置本地过滤关键字 - Set local filter keyword
  projectStore.setSearchKeyword(searchKeyword.value)
  
  // 防抖处理远程搜索 - Debounce remote search
  searchTimer = setTimeout(() => {
    if (searchKeyword.value.length >= 2) {
      projectStore.search(searchKeyword.value)
    }
  }, 500)
}

/*** 处理搜索确认 - Handle search confirm (keyboard enter) ***/
const handleSearchConfirm = (): void => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  if (searchKeyword.value) {
    projectStore.search(searchKeyword.value)
  }
}

/*** 清除搜索 - Clear search ***/
const clearSearch = (): void => {
  searchKeyword.value = ''
  projectStore.setSearchKeyword('')
  
  // 重新加载完整列表 - Reload full list
  const cachedData = projectStore.loadFromCache()
  if (cachedData && cachedData.length > 0) {
    // 使用缓存数据 - Use cached data
  } else {
    projectStore.fetchList(true)
  }
}

/*** 处理下拉刷新 - Handle pull down refresh ***/
const handleRefresh = async (): Promise<void> => {
  if (isRefreshing.value) return
  
  isRefreshing.value = true
  
  try {
    await projectStore.fetchList(true)
  } catch (error) {
    console.error('刷新失败:', error)
    uni.showToast({ title: '刷新失败，请重试', icon: 'none' })
  } finally {
    isRefreshing.value = false
  }
}

/*** 处理上拉加载更多 - Handle scroll to bottom load more ***/
const handleLoadMore = async (): Promise<void> => {
  if (loadingMore.value || !projectStore.hasMore || projectStore.loading) return
  
  loadingMore.value = true
  
  try {
    await projectStore.loadMore()
  } catch (error) {
    console.error('加载更多失败:', error)
    uni.showToast({ title: '加载失败，请重试', icon: 'none' })
  } finally {
    loadingMore.value = false
  }
}

/*** 跳转到项目详情 - Navigate to project detail ***/
const goToDetail = (id: number): void => {
  uni.navigateTo({ url: `/pages/project/detail?id=${id}` })
}

/*** 初始化加载 - Initial load ***/
const initLoad = async (): Promise<void> => {
  checkNetworkStatus()
  
  // 先尝试从缓存加载 - Try to load from cache first
  const cachedData = projectStore.loadFromCache()
  
  if (cachedData && cachedData.length > 0 && !projectStore.isCacheExpired()) {
    // 使用缓存数据，后台静默刷新 - Use cache, silent refresh in background
    if (!isOffline.value) {
      projectStore.fetchList(true).catch(() => {
        // 静默失败，已有缓存数据 - Silent fail, already have cached data
      })
    }
  } else {
    // 无缓存或已过期，直接请求 - No cache or expired, fetch directly
    try {
      await projectStore.fetchList(true)
    } catch (error) {
      if (isOffline.value) {
        uni.showToast({ title: '网络不可用，请检查网络连接', icon: 'none' })
      }
    }
  }
}

/*** 页面显示时 - On page show ***/
onShow(() => {
  checkNetworkStatus()
})

/*** 组件挂载时 - On component mounted ***/
onMounted(() => {
  setupNetworkListener()
  initLoad()
})

/*** 页面下拉刷新 - Page pull down refresh ***/
onPullDownRefresh(async () => {
  await handleRefresh()
  uni.stopPullDownRefresh()
})

/*** 页面触底 - Page reach bottom ***/
onReachBottom(() => {
  handleLoadMore()
})

/*** 监听搜索关键字变化 - Watch search keyword changes ***/
watch(searchKeyword, (newVal) => {
  if (!newVal) {
    projectStore.setSearchKeyword('')
  }
})
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container with gradient background ***/
.page-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #E8F4FD 0%, #FFFFFF 30%, #F8FAFC 100%);
}

/*** 顶部区域 - Header area with blue gradient ***/
.header-area {
  position: relative;
  padding: 60rpx 40rpx 80rpx;
  overflow: hidden;
}

.header-bg {
  position: absolute;
  top: -100rpx;
  left: -50rpx;
  right: -50rpx;
  height: 400rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 50%, #1D4ED8 100%);
  border-radius: 0 0 60rpx 60rpx;
  transform: rotate(-3deg);
}

.header-content {
  position: relative;
  z-index: 2;
}

.header-title {
  font-size: 48rpx;
  font-weight: 700;
  color: #FFFFFF;
  display: block;
  margin-bottom: 8rpx;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
}

.header-subtitle {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.85);
}

/*** 搜索区域 - Search section with glass morphism ***/
.search-section {
  padding: 0 32rpx;
  margin-top: -50rpx;
  position: relative;
  z-index: 10;
}

.search-wrapper {
  display: flex;
  align-items: center;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 0 28rpx;
  height: 96rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.12),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.search-icon-wrapper {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 16rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #1E293B;
  height: 100%;
}

.search-input::placeholder {
  color: #94A3B8;
}

.clear-btn {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  transition: opacity 0.2s ease;
  
  &:active {
    opacity: 0.7;
  }
}

/*** 项目列表 - Project list scroll view ***/
.project-list {
  height: calc(100vh - 280rpx);
  padding: 32rpx;
}

.project-cards {
  display: flex;
  flex-direction: column;
  gap: 28rpx;
}

/*** 加载状态 - Loading state ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-spinner {
  width: 48rpx;
  height: 48rpx;
  border: 4rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 20rpx;
}

.loading-spinner-small {
  width: 32rpx;
  height: 32rpx;
  border: 3rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 12rpx;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 26rpx;
  color: #94A3B8;
}

/*** 空状态 - Empty state ***/
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.empty-icon {
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F1F5F9;
  border-radius: 50%;
  margin-bottom: 32rpx;
  color: #CBD5E1;
}

.empty-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #64748B;
  margin-bottom: 12rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: #94A3B8;
}

/*** 加载更多 - Load more ***/
.load-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 0;
}

.load-more-text {
  font-size: 26rpx;
  color: #94A3B8;
}

/*** 没有更多 - No more data ***/
.no-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 0;
}

.no-more-text {
  font-size: 24rpx;
  color: #CBD5E1;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}

/*** 离线提示 - Offline tip ***/
.offline-tip {
  position: fixed;
  bottom: 180rpx;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  background: rgba(30, 41, 59, 0.9);
  padding: 16rpx 28rpx;
  border-radius: 40rpx;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.15);
  z-index: 100;
}

.offline-icon {
  width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FBBF24;
  margin-right: 12rpx;
}

.offline-text {
  font-size: 24rpx;
  color: #FFFFFF;
}
</style>
