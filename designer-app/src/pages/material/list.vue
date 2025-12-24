<template>
  <!-- 
    材料列表页面 - C4D风格高级UI
    Material List Page - C4D style premium UI with search and grid layout
  -->
  <view class="page-container">
    <!-- 搜索栏 - Search bar with glass morphism -->
    <view class="search-section">
      <view class="search-bar">
        <view class="search-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </svg>
        </view>
        <input 
          class="search-input"
          type="text"
          v-model="searchKeyword"
          placeholder="搜索材料名称、品牌..."
          confirm-type="search"
          @input="handleSearchInput"
          @confirm="handleSearch"
        />
        <view v-if="searchKeyword" class="clear-icon" @click="clearSearch">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="m15 9-6 6"></path>
            <path d="m9 9 6 6"></path>
          </svg>
        </view>
      </view>
    </view>
    
    <!-- 分类标签 - Category tag -->
    <view class="category-tag-section" v-if="categoryName">
      <view class="category-tag">
        <text class="tag-text">{{ categoryName }}</text>
        <view class="tag-close" @click="clearCategory">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </view>
      </view>
      <text class="result-count">共 {{ total }} 种材料</text>
    </view>
    
    <!-- 材料列表 - Material list -->
    <scroll-view 
      class="material-list-container"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="refreshing"
      @refresherrefresh="handleRefresh"
      @scrolltolower="handleLoadMore"
    >
      <!-- 加载状态 - Loading state -->
      <view v-if="loading && materials.length === 0" class="loading-container">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>
      
      <!-- 材料网格 - Material grid -->
      <view v-else-if="materials.length > 0" class="material-grid">
        <view 
          v-for="(material, index) in materials" 
          :key="material.id"
          class="material-item"
          :style="{ animationDelay: ((index as number) % 10 * 0.05) + 's' }"
        >
          <MaterialCard 
            :material="material" 
            @click="handleMaterialClick"
          />
        </view>
      </view>
      
      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <view class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </svg>
        </view>
        <text class="empty-text">{{ emptyText }}</text>
        <view v-if="!searchKeyword" class="retry-btn" @click="() => loadMaterials()">
          <text class="retry-text">点击重试</text>
        </view>
      </view>
      
      <!-- 加载更多状态 - Load more state -->
      <view v-if="materials.length > 0" class="load-more-container">
        <view v-if="loadingMore" class="loading-more">
          <view class="loading-spinner-small"></view>
          <text class="loading-more-text">加载中...</text>
        </view>
        <text v-else-if="!hasMore" class="no-more-text">没有更多了</text>
      </view>
    </scroll-view>
    
    <!-- 错误提示 - Error toast -->
    <view v-if="errorMsg" class="error-toast">
      <text class="error-text">{{ errorMsg }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Material List Page
 * 材料列表页面 - 展示分类下的材料列表，支持搜索功能
 * Display materials in category with search functionality
 ***/
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { onLoad, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { materialApi } from '@/api/material'
import { MaterialCard } from '@/components'
import { useNetworkStore } from '@/store/network'
import type { Material, MaterialListParams } from '@/types/api'
import { DEFAULT_PAGE_SIZE } from '@/types/common'

/*** Store instance - 状态管理实例 ***/
const networkStore = useNetworkStore()

/*** 分类名称 - Category name from route params ***/
const categoryName = ref('')

/*** 搜索关键字 - Search keyword ***/
const searchKeyword = ref('')

/*** 材料列表 - Material list ***/
const materials = ref<Material[]>([])

/*** 分页参数 - Pagination params ***/
const currentPage = ref(1)
const pageSize = ref(DEFAULT_PAGE_SIZE)
const total = ref(0)

/*** 加载状态 - Loading states ***/
const loading = ref(false)
const loadingMore = ref(false)
const refreshing = ref(false)

/*** 错误信息 - Error message ***/
const errorMsg = ref('')

/*** 搜索防抖定时器 - Search debounce timer ***/
let searchTimer: ReturnType<typeof setTimeout> | null = null

/*** 是否还有更多数据 - Has more data ***/
const hasMore = computed((): boolean => {
  return materials.value.length < total.value
})

/*** 空状态文本 - Empty state text ***/
const emptyText = computed((): string => {
  if (searchKeyword.value) {
    return `未找到"${searchKeyword.value}"相关材料`
  }
  if (networkStore.isOffline) {
    return '网络不可用，请检查网络连接'
  }
  return '暂无材料数据'
})

/*** 处理网络恢复事件 - Handle network recovery event ***/
const handleNetworkRecovered = (): void => {
  /*** Auto reload data when network recovers ***/
  if (materials.value.length === 0) {
    loadMaterials()
  }
}

/*** 页面加载 - Page load with route params ***/
onLoad((options) => {
  const params = options as { category?: string } | undefined
  if (params?.category) {
    categoryName.value = decodeURIComponent(params.category)
    // 设置导航栏标题
    uni.setNavigationBarTitle({
      title: categoryName.value
    })
  }
  loadMaterials()
})

/*** 组件挂载时 - On component mounted ***/
onMounted(() => {
  /*** Listen to network recovery event ***/
  uni.$on('network:recovered', handleNetworkRecovered)
})

/*** 组件卸载时 - On component unmounted ***/
onUnmounted(() => {
  /*** Remove network recovery event listener ***/
  uni.$off('network:recovered', handleNetworkRecovered)
})

/*** 下拉刷新 - Pull down refresh ***/
onPullDownRefresh(() => {
  handleRefresh()
})

/*** 上拉加载更多 - Reach bottom load more ***/
onReachBottom(() => {
  handleLoadMore()
})

/*** 加载材料列表 - Load material list ***/
const loadMaterials = async (isLoadMore: boolean = false): Promise<void> => {
  if (loading.value || loadingMore.value) return
  
  if (isLoadMore) {
    if (!hasMore.value) return
    loadingMore.value = true
  } else {
    loading.value = true
    currentPage.value = 1
    materials.value = []
  }
  
  errorMsg.value = ''
  
  try {
    const params: MaterialListParams = {
      current: currentPage.value,
      size: pageSize.value
    }
    
    // 添加分类筛选
    if (categoryName.value) {
      params.category = categoryName.value
    }
    
    // 添加搜索关键字
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    
    const res = await materialApi.getList(params)
    
    if (res.code === 200 && res.data) {
      const { records, total: totalCount } = res.data
      
      if (isLoadMore) {
        materials.value = [...materials.value, ...records]
      } else {
        materials.value = records
      }
      
      total.value = totalCount
      currentPage.value++
    } else {
      showError(res.msg || '加载失败')
    }
  } catch (error) {
    console.error('Failed to load materials:', error)
    if (networkStore.isOffline) {
      showError('网络不可用，请检查网络连接')
    } else {
      showError('网络错误，请稍后重试')
    }
  } finally {
    loading.value = false
    loadingMore.value = false
    refreshing.value = false
    uni.stopPullDownRefresh()
  }
}

/*** 处理下拉刷新 - Handle pull down refresh ***/
const handleRefresh = (): void => {
  refreshing.value = true
  loadMaterials()
}

/*** 处理加载更多 - Handle load more ***/
const handleLoadMore = (): void => {
  if (hasMore.value && !loadingMore.value) {
    loadMaterials(true)
  }
}

/*** 处理搜索输入（防抖） - Handle search input with debounce ***/
const handleSearchInput = (): void => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  searchTimer = setTimeout(() => {
    loadMaterials()
  }, 500)
}

/*** 处理搜索确认 - Handle search confirm ***/
const handleSearch = (): void => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  loadMaterials()
}

/*** 清除搜索 - Clear search ***/
const clearSearch = (): void => {
  searchKeyword.value = ''
  loadMaterials()
}

/*** 清除分类筛选 - Clear category filter ***/
const clearCategory = (): void => {
  categoryName.value = ''
  uni.setNavigationBarTitle({
    title: '材料列表'
  })
  loadMaterials()
}

/*** 处理材料点击 - Handle material click ***/
const handleMaterialClick = (id: number): void => {
  uni.navigateTo({
    url: `/pages/material/detail?id=${id}`
  })
}

/*** 显示错误提示 - Show error message ***/
const showError = (msg: string): void => {
  errorMsg.value = msg
  setTimeout(() => {
    errorMsg.value = ''
  }, 3000)
}
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container ***/
.page-container {
  min-height: 100vh;
  background: #F5F7FA;
  display: flex;
  flex-direction: column;
}

/*** 搜索区域 - Search section ***/
.search-section {
  padding: 16rpx 24rpx;
  background: #FFFFFF;
  position: sticky;
  top: 0;
  z-index: 10;
}

.search-bar {
  display: flex;
  align-items: center;
  background: #F1F5F9;
  border-radius: 40rpx;
  padding: 16rpx 24rpx;
}

.search-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #1E293B;
  background: transparent;
  
  &::placeholder {
    color: #94A3B8;
  }
}

.clear-icon {
  width: 36rpx;
  height: 36rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-left: 12rpx;
  flex-shrink: 0;
  
  &:active {
    opacity: 0.7;
  }
}

/*** 分类标签区域 - Category tag section ***/
.category-tag-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 24rpx;
  background: #FFFFFF;
  border-top: 1rpx solid #F1F5F9;
}

.category-tag {
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border-radius: 20rpx;
}

.tag-text {
  font-size: 24rpx;
  color: #3B82F6;
  font-weight: 500;
}

.tag-close {
  width: 28rpx;
  height: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3B82F6;
  margin-left: 8rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.result-count {
  font-size: 24rpx;
  color: #94A3B8;
}

/*** 材料列表容器 - Material list container ***/
.material-list-container {
  flex: 1;
  padding: 16rpx 24rpx;
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
  width: 64rpx;
  height: 64rpx;
  border: 4rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 24rpx;
  font-size: 26rpx;
  color: #94A3B8;
}

/*** 材料网格 - Material grid (2 columns) ***/
.material-grid {
  display: flex;
  flex-wrap: wrap;
  margin: -10rpx;
}

.material-item {
  width: 50%;
  padding: 10rpx;
  box-sizing: border-box;
  animation: fadeInUp 0.4s ease forwards;
  opacity: 0;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
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
  width: 128rpx;
  height: 128rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #94A3B8;
  margin-bottom: 32rpx;
  text-align: center;
}

.retry-btn {
  padding: 16rpx 48rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  border-radius: 40rpx;
  
  &:active {
    opacity: 0.8;
  }
}

.retry-text {
  font-size: 26rpx;
  color: #FFFFFF;
  font-weight: 500;
}

/*** 加载更多状态 - Load more state ***/
.load-more-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 0 48rpx;
}

.loading-more {
  display: flex;
  align-items: center;
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

.loading-more-text {
  font-size: 24rpx;
  color: #94A3B8;
}

.no-more-text {
  font-size: 24rpx;
  color: #CBD5E1;
}

/*** 错误提示 - Error toast ***/
.error-toast {
  position: fixed;
  bottom: 200rpx;
  left: 50%;
  transform: translateX(-50%);
  padding: 20rpx 40rpx;
  background: rgba(0, 0, 0, 0.75);
  border-radius: 40rpx;
  z-index: 999;
}

.error-text {
  font-size: 26rpx;
  color: #FFFFFF;
}
</style>
