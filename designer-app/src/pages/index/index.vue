<template>
  <view class="page-container">
    <!-- 搜索框 - Search bar -->
    <view class="search-bar">
      <view class="search-input-wrapper glass-card">
        <text class="search-icon">🔍</text>
        <input
          class="search-input"
          type="text"
          placeholder="搜索项目名称"
          v-model="searchKeyword"
          @input="handleSearch"
        />
        <text v-if="searchKeyword" class="clear-icon" @click="clearSearch">✕</text>
      </view>
    </view>

    <!-- 项目列表 - Project list -->
    <scroll-view
      class="project-list"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="handleRefresh"
      @scrolltolower="handleLoadMore"
    >
      <!-- 加载中状态 - Loading state -->
      <view v-if="loading && projects.length === 0" class="loading-container">
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 项目卡片列表 - Project cards -->
      <view v-else-if="projects.length > 0" class="project-cards">
        <view
          v-for="project in filteredProjects"
          :key="project.id"
          class="project-card float-card"
          @click="goToDetail(project.id)"
        >
          <view class="card-header">
            <text class="project-name">{{ project.name }}</text>
            <view :class="['status-tag', `status-${project.status}`]">
              {{ getStatusText(project.status) }}
            </view>
          </view>
          <view class="card-body">
            <view class="info-row">
              <text class="info-label">面积</text>
              <text class="info-value">{{ project.area }}㎡</text>
            </view>
            <view class="info-row">
              <text class="info-label">更新时间</text>
              <text class="info-value">{{ formatDate(project.updatedAt) }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <text class="empty-text">暂无项目</text>
      </view>

      <!-- 加载更多 - Load more -->
      <view v-if="hasMore && projects.length > 0" class="load-more">
        <text class="load-more-text">{{ loadingMore ? '加载中...' : '上拉加载更多' }}</text>
      </view>
    </scroll-view>

    <!-- 离线模式提示 - Offline mode tip -->
    <view v-if="isOffline" class="offline-tip">
      <text class="offline-text">当前为离线模式</text>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** Project list page - displays user's projects with search and pagination ***/
import { ref, computed, onMounted } from 'vue'
import { onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'

// 项目数据类型 - Project data type
interface Project {
  id: number
  name: string
  status: 'draft' | 'in_progress' | 'completed' | 'archived'
  area: number
  budget: number
  style: string
  updatedAt: string
}

// 响应式数据 - Reactive data
const projects = ref<Project[]>([])
const searchKeyword = ref('')
const loading = ref(false)
const loadingMore = ref(false)
const isRefreshing = ref(false)
const hasMore = ref(true)
const isOffline = ref(false)
const currentPage = ref(1)
const pageSize = 10

// 过滤后的项目列表 - Filtered projects
const filteredProjects = computed(() => {
  if (!searchKeyword.value) return projects.value
  const keyword = searchKeyword.value.toLowerCase()
  return projects.value.filter(p => p.name.toLowerCase().includes(keyword))
})

// 获取状态文本 - Get status text
const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    draft: '草稿',
    in_progress: '进行中',
    completed: '已完成',
    archived: '已归档'
  }
  return statusMap[status] || status
}

// 格式化日期 - Format date
const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// 搜索处理 - Handle search
const handleSearch = () => {
  // 实时过滤，无需额外处理
}

// 清除搜索 - Clear search
const clearSearch = () => {
  searchKeyword.value = ''
}

// 下拉刷新 - Pull down refresh
const handleRefresh = async () => {
  isRefreshing.value = true
  currentPage.value = 1
  hasMore.value = true
  await loadProjects(true)
  isRefreshing.value = false
}

// 加载更多 - Load more
const handleLoadMore = async () => {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  currentPage.value++
  await loadProjects(false)
  loadingMore.value = false
}

// 加载项目列表 - Load projects
const loadProjects = async (refresh: boolean = false) => {
  loading.value = true
  try {
    // TODO: 调用后端接口获取项目列表
    // const res = await projectApi.getList({ current: currentPage.value, size: pageSize })
    
    // 模拟数据 - Mock data for development
    const mockProjects: Project[] = [
      { id: 1, name: '星巴克咖啡厅装修项目', status: 'in_progress', area: 150, budget: 500000, style: '现代简约', updatedAt: '2024-12-20T10:00:00Z' },
      { id: 2, name: '华为体验店设计', status: 'completed', area: 300, budget: 1200000, style: '科技风', updatedAt: '2024-12-19T15:30:00Z' },
      { id: 3, name: '万达广场餐饮区', status: 'draft', area: 500, budget: 2000000, style: '工业风', updatedAt: '2024-12-18T09:00:00Z' }
    ]
    
    if (refresh) {
      projects.value = mockProjects
    } else {
      projects.value = [...projects.value, ...mockProjects]
    }
    
    // 判断是否还有更多数据
    hasMore.value = mockProjects.length >= pageSize
  } catch (error) {
    console.error('加载项目列表失败:', error)
    uni.showToast({ title: '加载失败，请重试', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// 跳转到项目详情 - Navigate to project detail
const goToDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/project/detail?id=${id}` })
}

// 页面加载 - Page mounted
onMounted(() => {
  loadProjects(true)
})

// 下拉刷新生命周期 - Pull down refresh lifecycle
onPullDownRefresh(async () => {
  await handleRefresh()
  uni.stopPullDownRefresh()
})

// 触底加载更多 - Reach bottom load more
onReachBottom(() => {
  handleLoadMore()
})
</script>

<style lang="scss" scoped>
.page-container {
  min-height: 100vh;
  background-color: #F5F7FA;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.search-bar {
  padding: 24rpx 32rpx;
  background: #FFFFFF;
  position: sticky;
  top: 0;
  z-index: 100;
}

.search-input-wrapper {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background: #F5F7FA;
  border-radius: 16rpx;
}

.search-icon {
  font-size: 32rpx;
  margin-right: 16rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #2C3E50;
}

.clear-icon {
  font-size: 28rpx;
  color: #94A3B8;
  padding: 8rpx;
}

.project-list {
  height: calc(100vh - 120rpx);
  padding: 24rpx 32rpx;
}

.project-cards {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.project-card {
  padding: 32rpx;
  background: #FFFFFF;
  border-radius: 16rpx;
  box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.project-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #2C3E50;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-tag {
  font-size: 24rpx;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  margin-left: 16rpx;
}

.status-draft {
  background: #F1F5F9;
  color: #64748B;
}

.status-in_progress {
  background: #DBEAFE;
  color: #3B82F6;
}

.status-completed {
  background: #D1FAE5;
  color: #10B981;
}

.status-archived {
  background: #FEE2E2;
  color: #EF4444;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 26rpx;
  color: #64748B;
}

.info-value {
  font-size: 26rpx;
  color: #2C3E50;
}

.loading-container,
.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 100rpx 0;
}

.loading-text,
.empty-text {
  font-size: 28rpx;
  color: #94A3B8;
}

.load-more {
  display: flex;
  justify-content: center;
  padding: 32rpx 0;
}

.load-more-text {
  font-size: 26rpx;
  color: #94A3B8;
}

.offline-tip {
  position: fixed;
  bottom: 120rpx;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.7);
  padding: 16rpx 32rpx;
  border-radius: 32rpx;
}

.offline-text {
  font-size: 26rpx;
  color: #FFFFFF;
}
</style>
""