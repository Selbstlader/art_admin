<template>
  <!-- 
    首页项目列表 - C4D风格高级UI
    Home page project list - C4D style premium UI with 3D effects
  -->
  <view class="page-container">
    <!-- 背景装饰元素 - Background decorative elements -->
    <view class="bg-decoration">
      <view class="bg-orb bg-orb-1"></view>
      <view class="bg-orb bg-orb-2"></view>
      <view class="bg-orb bg-orb-3"></view>
    </view>

    <!-- 顶部区域 - Header area with 3D blue gradient -->
    <view class="header-area">
      <view class="header-bg">
        <view class="header-bg-layer-1"></view>
        <view class="header-bg-layer-2"></view>
        <view class="header-bg-shine"></view>
      </view>
      <view class="header-content">
        <view class="header-badge">
          <view class="badge-dot"></view>
          <text class="badge-text">设计工作台</text>
        </view>
        <text class="header-title">项目列表</text>
        <text class="header-subtitle">管理您的设计项目</text>
      </view>
      <!-- 3D装饰球体 - 3D decorative spheres -->
      <view class="header-sphere header-sphere-1"></view>
      <view class="header-sphere header-sphere-2"></view>
    </view>

    <!-- 搜索框 - Search bar with glass morphism effect -->
    <view class="search-section">
      <view class="search-wrapper">
        <view class="search-glow"></view>
        <view class="search-inner">
          <view class="search-icon-wrapper">
            <!-- Search icon -->
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.3-4.3"></path>
            </svg>
          </view>
          <input class="search-input" type="text" placeholder="搜索项目名称" v-model="searchKeyword" @input="handleSearch"
            @confirm="handleSearchConfirm" />
          <view v-if="searchKeyword" class="clear-btn" @click="clearSearch">
            <!-- Clear icon -->
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="m15 9-6 6"></path>
              <path d="m9 9 6 6"></path>
            </svg>
          </view>
        </view>
      </view>
    </view>

    <!-- 项目列表 - Project list with scroll view -->
    <scroll-view class="project-list" scroll-y :refresher-enabled="true" :refresher-triggered="isRefreshing"
      @refresherrefresh="handleRefresh" @scrolltolower="handleLoadMore" refresher-default-style="none">
      <!-- 加载中状态 - Loading state (仅首次加载时显示，排除下拉刷新场景) -->
      <view v-if="projectStore.loading && projectStore.list.length === 0 && !isRefreshing" class="loading-container">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 项目卡片列表 - Project cards with 3D effect -->
      <view v-else-if="displayProjects.length > 0" class="project-cards">
        <ProjectCard v-for="project in displayProjects" :key="project.id" :project="project" @click="goToDetail" />
      </view>

      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <view class="empty-icon">
          <!-- Folder icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <path
              d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z">
            </path>
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

    <!-- 离线模式提示 - Offline mode tip (using global component) -->
    <OfflineTip :bottom="180" />

    <!-- 自定义底部导航栏 - Custom TabBar -->
    <CustomTabBar />
  </view>
</template>

<script setup lang="ts">
/*** 
 * Project list page - C4D style premium UI with 3D cards
 * 首页项目列表页面 - C4D风格高级UI，3D立体卡片设计
 ***/
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import { useProjectStore } from '@/store/project'
import { useNetworkStore } from '@/store/network'
import { ProjectCard, OfflineTip } from '@/components'
import CustomTabBar from '@/components/CustomTabBar.vue'

/*** Store instance - 状态管理实例 ***/
const projectStore = useProjectStore()
const networkStore = useNetworkStore()

/*** Local state - 本地状态 ***/
const searchKeyword = ref('')
const isRefreshing = ref(false)
const loadingMore = ref(false)
let searchTimer: ReturnType<typeof setTimeout> | null = null

/*** 显示的项目列表 - Display projects (filtered or all) ***/
const displayProjects = computed(() => {
  if (searchKeyword.value) {
    return projectStore.filteredList
  }
  return projectStore.list
})

/*** 处理网络恢复事件 - Handle network recovery event ***/
const handleNetworkRecovered = (): void => {
  /*** Auto refresh data when network recovers ***/
  handleRefresh()
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
  // 先尝试从缓存加载 - Try to load from cache first
  const cachedData = projectStore.loadFromCache()

  if (cachedData && cachedData.length > 0 && !projectStore.isCacheExpired()) {
    // 使用缓存数据，后台静默刷新 - Use cache, silent refresh in background
    if (networkStore.isOnline) {
      projectStore.fetchList(true).catch(() => {
        // 静默失败，已有缓存数据 - Silent fail, already have cached data
      })
    }
  } else {
    // 无缓存或已过期，直接请求 - No cache or expired, fetch directly
    try {
      await projectStore.fetchList(true)
    } catch (error) {
      if (networkStore.isOffline) {
        uni.showToast({ title: '网络不可用，请检查网络连接', icon: 'none' })
      }
    }
  }
}

/*** 页面显示时 - On page show ***/
onShow(() => {
  /*** Check network status when page shows ***/
  networkStore.checkNetworkStatus()
})

/*** 组件挂载时 - On component mounted ***/
onMounted(() => {
  /*** Listen to network recovery event ***/
  uni.$on('network:recovered', handleNetworkRecovered)

  initLoad()
})

/*** 组件卸载时 - On component unmounted ***/
onUnmounted(() => {
  /*** Remove network recovery event listener ***/
  uni.$off('network:recovered', handleNetworkRecovered)
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
/*** 页面容器 - Page container with subtle gradient ***/
.page-container {
  min-height: calc(100vh - 50px);
  background: linear-gradient(180deg, #FFFFFF 0%, #F8FAFC 100%);
  position: relative;
  overflow: hidden;
}

/*** 背景装饰元素 - Background decorative orbs with glass effect ***/
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100rpx);
  opacity: 0.35;
}

.bg-orb-1 {
  width: 500rpx;
  height: 500rpx;
  background: linear-gradient(135deg, #93C5FD 0%, #3B82F6 100%);
  top: -150rpx;
  right: -150rpx;
  animation: float-orb 12s ease-in-out infinite;
}

.bg-orb-2 {
  width: 400rpx;
  height: 400rpx;
  background: linear-gradient(135deg, #BFDBFE 0%, #60A5FA 100%);
  top: 500rpx;
  left: -200rpx;
  animation: float-orb 15s ease-in-out infinite reverse;
}

.bg-orb-3 {
  width: 300rpx;
  height: 300rpx;
  background: linear-gradient(135deg, #DBEAFE 0%, #93C5FD 100%);
  bottom: 400rpx;
  right: -100rpx;
  animation: float-orb 10s ease-in-out infinite 2s;
}

@keyframes float-orb {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-30rpx) scale(1.08); }
}

/*** 顶部区域 - Header area with 3D blue gradient ***/
.header-area {
  position: relative;
  padding: 80rpx 40rpx 110rpx;
  overflow: visible;
}

.header-bg {
  position: absolute;
  top: -80rpx;
  left: -60rpx;
  right: -60rpx;
  height: 460rpx;
  overflow: hidden;
}

.header-bg-layer-1 {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    155deg,
    #60A5FA 0%,
    #3B82F6 30%,
    #2563EB 60%,
    #1D4ED8 100%
  );
  border-radius: 0 0 100rpx 100rpx;
  transform: rotate(-3deg) translateY(-30rpx);
}

.header-bg-layer-2 {
  position: absolute;
  top: 30rpx;
  left: 30rpx;
  right: 30rpx;
  bottom: -30rpx;
  background: linear-gradient(155deg, #3B82F6 0%, #1E40AF 100%);
  border-radius: 0 0 80rpx 80rpx;
  transform: rotate(-1.5deg);
  opacity: 0.5;
}

.header-bg-shine {
  position: absolute;
  top: 50rpx;
  left: 80rpx;
  width: 280rpx;
  height: 140rpx;
  background: linear-gradient(
    180deg,
    rgba(255, 255, 255, 0.35) 0%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 100%
  );
  border-radius: 140rpx;
  transform: rotate(-20deg);
}

/*** 3D装饰球体 - 3D decorative glass spheres ***/
.header-sphere {
  position: absolute;
  border-radius: 50%;
  z-index: 1;
}

.header-sphere-1 {
  width: 140rpx;
  height: 140rpx;
  top: 50rpx;
  right: 30rpx;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.5) 0%,
    rgba(255, 255, 255, 0.15) 40%,
    rgba(255, 255, 255, 0.05) 100%
  );
  box-shadow:
    inset -15rpx -15rpx 40rpx rgba(0, 0, 0, 0.08),
    inset 15rpx 15rpx 40rpx rgba(255, 255, 255, 0.4),
    0 25rpx 50rpx rgba(30, 64, 175, 0.35);
  animation: float-sphere 5s ease-in-out infinite;
}

.header-sphere-2 {
  width: 70rpx;
  height: 70rpx;
  top: 190rpx;
  right: 150rpx;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.6) 0%,
    rgba(255, 255, 255, 0.2) 50%,
    transparent 100%
  );
  box-shadow:
    inset -8rpx -8rpx 20rpx rgba(0, 0, 0, 0.06),
    inset 8rpx 8rpx 20rpx rgba(255, 255, 255, 0.5),
    0 15rpx 30rpx rgba(30, 64, 175, 0.25);
  animation: float-sphere 4s ease-in-out infinite 0.8s;
}

@keyframes float-sphere {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20rpx) rotate(5deg); }
}

.header-content {
  position: relative;
  z-index: 2;
}

.header-badge {
  display: inline-flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  padding: 10rpx 24rpx;
  border-radius: 32rpx;
  margin-bottom: 20rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.3);
}

.badge-dot {
  width: 14rpx;
  height: 14rpx;
  background: linear-gradient(135deg, #4ADE80 0%, #22C55E 100%);
  border-radius: 50%;
  margin-right: 12rpx;
  box-shadow: 0 0 12rpx rgba(34, 197, 94, 0.6);
  animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; box-shadow: 0 0 12rpx rgba(34, 197, 94, 0.6); }
  50% { opacity: 0.7; box-shadow: 0 0 20rpx rgba(34, 197, 94, 0.8); }
}

.badge-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.95);
  font-weight: 600;
  letter-spacing: 1rpx;
}

.header-title {
  font-size: 60rpx;
  font-weight: 800;
  color: #FFFFFF;
  display: block;
  margin-bottom: 14rpx;
  text-shadow: 0 6rpx 20rpx rgba(0, 0, 0, 0.2);
  letter-spacing: 3rpx;
}

.header-subtitle {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
  letter-spacing: 1rpx;
}

/*** 搜索区域 - Search section with premium glass morphism ***/
.search-section {
  padding: 0 32rpx;
  margin-top: -65rpx;
  position: relative;
  z-index: 10;
}

.search-wrapper {
  position: relative;
}

.search-glow {
  position: absolute;
  top: 15rpx;
  left: 25rpx;
  right: 25rpx;
  bottom: -15rpx;
  background: linear-gradient(
    135deg,
    rgba(59, 130, 246, 0.25) 0%,
    rgba(37, 99, 235, 0.15) 100%
  );
  border-radius: 32rpx;
  filter: blur(25rpx);
}

.search-inner {
  position: relative;
  display: flex;
  align-items: center;
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.98) 0%,
    rgba(255, 255, 255, 0.92) 100%
  );
  backdrop-filter: blur(30px);
  -webkit-backdrop-filter: blur(30px);
  border-radius: 32rpx;
  padding: 0 36rpx;
  height: 112rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 12rpx 40rpx rgba(59, 130, 246, 0.18),
    0 4rpx 12rpx rgba(0, 0, 0, 0.04),
    inset 0 2rpx 6rpx rgba(255, 255, 255, 1);
}

.search-icon-wrapper {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3B82F6;
  margin-right: 18rpx;
}

.search-input {
  flex: 1;
  font-size: 32rpx;
  color: #1E293B;
  height: 100%;
  background: transparent;
  font-weight: 500;
}

.search-input::placeholder {
  color: #94A3B8;
  font-weight: 400;
}

.clear-btn {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748B;
  transition: all 0.25s ease;
  border-radius: 50%;
  background: linear-gradient(145deg, #F1F5F9 0%, #E2E8F0 100%);
  box-shadow: inset 0 1rpx 2rpx rgba(255, 255, 255, 0.8);

  &:active {
    opacity: 0.7;
    transform: scale(0.9);
  }
}

/*** 项目列表 - Project list scroll view ***/
.project-list {
  width: 100%;
  height: calc(100% - 400rpx);
  padding: 48rpx 32rpx;
  box-sizing: border-box;
}

.project-cards {
  display: flex;
  flex-direction: column;
  gap: 36rpx;
}

/*** 加载状态 - Loading state with glass spinner ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 180rpx 0;
}

.loading-spinner {
  width: 72rpx;
  height: 72rpx;
  border: 6rpx solid rgba(59, 130, 246, 0.15);
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.9s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  margin-bottom: 28rpx;
  box-shadow:
    0 6rpx 20rpx rgba(59, 130, 246, 0.25),
    inset 0 0 20rpx rgba(59, 130, 246, 0.05);
}

.loading-spinner-small {
  width: 40rpx;
  height: 40rpx;
  border: 4rpx solid rgba(59, 130, 246, 0.15);
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.9s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  margin-right: 14rpx;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 30rpx;
  color: #64748B;
  font-weight: 600;
  letter-spacing: 1rpx;
}

/*** 空状态 - Empty state with 3D glass icon ***/
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 180rpx 0;
}

.empty-icon {
  width: 200rpx;
  height: 200rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(
    145deg,
    rgba(255, 255, 255, 0.9) 0%,
    rgba(241, 245, 249, 0.8) 100%
  );
  border-radius: 50%;
  margin-bottom: 48rpx;
  color: #94A3B8;
  box-shadow:
    0 25rpx 50rpx rgba(59, 130, 246, 0.1),
    0 10rpx 20rpx rgba(0, 0, 0, 0.04),
    inset 0 -6rpx 12rpx rgba(0, 0, 0, 0.03),
    inset 0 6rpx 12rpx rgba(255, 255, 255, 1);
  border: 2rpx solid rgba(255, 255, 255, 0.8);
}

.empty-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #475569;
  margin-bottom: 14rpx;
  letter-spacing: 1rpx;
}

.empty-desc {
  font-size: 28rpx;
  color: #94A3B8;
  font-weight: 500;
}

/*** 加载更多 - Load more ***/
.load-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48rpx 0;
}

.load-more-text {
  font-size: 28rpx;
  color: #64748B;
  font-weight: 500;
}

/*** 没有更多 - No more data ***/
.no-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48rpx 0;
}

.no-more-text {
  font-size: 26rpx;
  color: #CBD5E1;
  letter-spacing: 3rpx;
  font-weight: 500;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 160rpx;
}
</style>
