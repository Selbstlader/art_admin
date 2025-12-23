<template>
  <!-- 
    设计图列表页面 - C4D风格高级UI
    Design images page - C4D style premium UI with 3D floating cards
  -->
  <view class="page-container">
    <!-- 顶部区域 - Header area with blue gradient -->
    <view class="header-area">
      <view class="header-bg"></view>
      <view class="header-content">
        <!-- 返回按钮 - Back button -->
        <view class="back-btn" @click="goBack">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6"></path>
          </svg>
        </view>
        <text class="header-title">设计图</text>
        <view class="header-count" v-if="!loading && !loadError">
          <text class="count-text">共 {{ total }} 张</text>
        </view>
      </view>
    </view>

    <!-- 加载状态 - Loading state -->
    <view v-if="loading && designList.length === 0" class="loading-container">
      <view class="loading-spinner"></view>
      <text class="loading-text">加载中...</text>
    </view>

    <!-- 错误状态 - Error state with retry -->
    <view v-else-if="loadError" class="error-container">
      <view class="error-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" x2="12" y1="8" y2="12"></line>
          <line x1="12" x2="12.01" y1="16" y2="16"></line>
        </svg>
      </view>
      <text class="error-title">加载失败</text>
      <text class="error-desc">{{ errorMessage }}</text>
      <view class="retry-btn" @click="loadDesignList">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
          <path d="M3 3v5h5"></path>
          <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
          <path d="M16 16h5v5"></path>
        </svg>
        <text class="retry-text">重新加载</text>
      </view>
    </view>

    <!-- 空状态 - Empty state -->
    <view v-else-if="designList.length === 0" class="empty-container">
      <view class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
          <circle cx="9" cy="9" r="2"></circle>
          <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
        </svg>
      </view>
      <text class="empty-title">暂无设计图</text>
      <text class="empty-desc">该项目还没有上传设计图</text>
    </view>

    <!-- 设计图网格列表 - Design image grid -->
    <scroll-view 
      v-else 
      class="content-scroll" 
      scroll-y 
      @scrolltolower="loadMore"
      :refresher-enabled="true"
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
    >
      <view class="design-grid">
        <view 
          v-for="(design, index) in designList" 
          :key="design.id" 
          class="design-card"
          @click="openPreview(index)"
        >
          <!-- 图片容器 - Image container -->
          <view class="image-wrapper">
            <!-- 加载进度指示器 - Loading progress indicator -->
            <view v-if="loadingStates[design.id]" class="image-loading">
              <view class="loading-ring"></view>
              <text class="loading-percent">{{ loadingProgress[design.id] || 0 }}%</text>
            </view>
            
            <!-- 加载失败占位图 - Error placeholder -->
            <view v-else-if="errorStates[design.id]" class="image-error" @click.stop="retryLoadImage(design)">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
                <circle cx="9" cy="9" r="2"></circle>
                <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
              </svg>
              <text class="error-hint">加载失败</text>
              <text class="retry-hint">点击重试</text>
            </view>
            
            <!-- 设计图图片 - Design image -->
            <image 
              v-else
              class="design-image"
              :src="design.thumbnailUrl || design.imageUrl"
              mode="aspectFill"
              :lazy-load="true"
              @load="onImageLoad(design)"
              @error="onImageError(design)"
            />
            
            <!-- 3D 悬浮效果阴影 - 3D floating shadow -->
            <view class="card-shadow"></view>
          </view>
          
          <!-- 设计图信息 - Design info -->
          <view class="design-info">
            <text class="design-name">{{ design.name || '未命名设计图' }}</text>
            <text class="design-date">{{ formatDate(design.createdAt) }}</text>
          </view>
        </view>
      </view>

      <!-- 加载更多状态 - Load more state -->
      <view v-if="hasMore && designList.length > 0" class="load-more">
        <view v-if="loadingMore" class="loading-more-spinner"></view>
        <text class="load-more-text">{{ loadingMore ? '加载中...' : '上拉加载更多' }}</text>
      </view>
      
      <!-- 没有更多数据 - No more data -->
      <view v-else-if="designList.length > 0" class="no-more">
        <text class="no-more-text">已加载全部设计图</text>
      </view>

      <!-- 底部安全区 - Bottom safe area -->
      <view class="safe-area-bottom"></view>
    </scroll-view>

    <!-- 全屏预览弹窗 - Fullscreen preview modal -->
    <view v-if="showPreview" class="preview-modal" @touchmove.stop.prevent>
      <!-- 预览背景 - Preview background -->
      <view class="preview-bg" @click="closePreview"></view>
      
      <!-- 预览头部 - Preview header -->
      <view class="preview-header">
        <view class="preview-close" @click="closePreview">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6 6 18"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </view>
        <text class="preview-title">{{ currentDesign?.name || '设计图预览' }}</text>
        <text class="preview-counter">{{ currentIndex + 1 }} / {{ designList.length }}</text>
      </view>
      
      <!-- 预览内容区 - Preview content area -->
      <swiper 
        class="preview-swiper"
        :current="currentIndex"
        @change="onSwiperChange"
        :circular="false"
      >
        <swiper-item v-for="(design, index) in designList" :key="design.id">
          <movable-area class="movable-area">
            <movable-view 
              class="movable-view"
              direction="all"
              :scale="true"
              :scale-min="1"
              :scale-max="4"
              :scale-value="scaleValue"
              @scale="onScale"
              @change="onMoveChange"
            >
              <image 
                class="preview-image"
                :src="design.imageUrl"
                mode="aspectFit"
                @load="onPreviewImageLoad(design)"
                @error="onPreviewImageError(design)"
              />
            </movable-view>
          </movable-area>
          
          <!-- 预览加载状态 - Preview loading state -->
          <view v-if="previewLoadingStates[design.id]" class="preview-loading">
            <view class="preview-loading-spinner"></view>
            <text class="preview-loading-text">加载中...</text>
          </view>
          
          <!-- 预览加载失败 - Preview load error -->
          <view v-if="previewErrorStates[design.id]" class="preview-error">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
              <circle cx="9" cy="9" r="2"></circle>
              <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
            </svg>
            <text class="preview-error-text">图片加载失败</text>
          </view>
        </swiper-item>
      </swiper>
      
      <!-- 预览底部操作提示 - Preview bottom hint -->
      <view class="preview-footer">
        <text class="preview-hint">双指缩放 · 左右滑动切换</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Design images page - C4D style premium UI
 * 设计图列表页面 - C4D风格高级UI，支持缩略图网格和全屏预览
 ***/
import { ref, reactive, computed, onMounted } from 'vue'
import { onLoad, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { designApi } from '@/api/design'
import type { DesignImage } from '@/types/api'
import { formatDate as formatDateUtil } from '@/utils/format'
import { DEFAULT_PAGE_SIZE } from '@/types/common'

/*** 页面参数 - Page params ***/
const projectId = ref<number>(0)

/*** 列表状态 - List state ***/
const designList = ref<DesignImage[]>([])
const total = ref(0)
const currentPage = ref(1)
const loading = ref(true)
const loadingMore = ref(false)
const refreshing = ref(false)
const loadError = ref(false)
const errorMessage = ref('网络异常，请稍后重试')

/*** 图片加载状态 - Image loading states ***/
const loadingStates = reactive<Record<number, boolean>>({})
const errorStates = reactive<Record<number, boolean>>({})
const loadingProgress = reactive<Record<number, number>>({})

/*** 预览状态 - Preview state ***/
const showPreview = ref(false)
const currentIndex = ref(0)
const scaleValue = ref(1)
const previewLoadingStates = reactive<Record<number, boolean>>({})
const previewErrorStates = reactive<Record<number, boolean>>({})

/*** 是否有更多数据 - Has more data ***/
const hasMore = computed(() => designList.value.length < total.value)

/*** 当前预览的设计图 - Current preview design ***/
const currentDesign = computed(() => designList.value[currentIndex.value])

/*** 格式化日期 - Format date ***/
const formatDate = (dateStr: string): string => {
  return formatDateUtil(dateStr, 'YYYY-MM-DD')
}

/*** 加载设计图列表 - Load design list ***/
const loadDesignList = async (isRefresh: boolean = false): Promise<void> => {
  if (!projectId.value) {
    loadError.value = true
    errorMessage.value = '项目ID无效'
    return
  }

  if (isRefresh) {
    currentPage.value = 1
    refreshing.value = true
  } else {
    loading.value = true
  }
  loadError.value = false

  try {
    const res = await designApi.getList({
      projectId: projectId.value,
      current: currentPage.value,
      size: DEFAULT_PAGE_SIZE
    })

    if (res.code === 200 && res.data) {
      if (isRefresh) {
        designList.value = res.data.records || []
      } else {
        designList.value = res.data.records || []
      }
      total.value = res.data.total || 0
      
      // 初始化图片加载状态 - Initialize image loading states
      res.data.records?.forEach((design: DesignImage) => {
        loadingStates[design.id] = true
        errorStates[design.id] = false
        loadingProgress[design.id] = 0
        previewLoadingStates[design.id] = true
        previewErrorStates[design.id] = false
      })
    } else {
      if (designList.value.length === 0) {
        loadError.value = true
        errorMessage.value = res.msg || res.message || '加载失败'
      }
    }
  } catch (error: unknown) {
    console.error('加载设计图列表失败:', error)
    if (designList.value.length === 0) {
      loadError.value = true
      if (error instanceof Error) {
        errorMessage.value = error.message || '网络异常，请稍后重试'
      } else {
        errorMessage.value = '网络异常，请稍后重试'
      }
    }
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

/*** 加载更多 - Load more ***/
const loadMore = async (): Promise<void> => {
  if (!hasMore.value || loadingMore.value) return

  loadingMore.value = true
  currentPage.value++

  try {
    const res = await designApi.getList({
      projectId: projectId.value,
      current: currentPage.value,
      size: DEFAULT_PAGE_SIZE
    })

    if (res.code === 200 && res.data) {
      const newRecords = res.data.records || []
      designList.value = [...designList.value, ...newRecords]
      total.value = res.data.total || 0
      
      // 初始化新图片的加载状态 - Initialize new images loading states
      newRecords.forEach((design: DesignImage) => {
        loadingStates[design.id] = true
        errorStates[design.id] = false
        loadingProgress[design.id] = 0
        previewLoadingStates[design.id] = true
        previewErrorStates[design.id] = false
      })
    }
  } catch (error) {
    console.error('加载更多失败:', error)
    currentPage.value--
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loadingMore.value = false
  }
}

/*** 下拉刷新 - Pull down refresh ***/
const onRefresh = async (): Promise<void> => {
  await loadDesignList(true)
}

/*** 图片加载成功 - Image load success ***/
const onImageLoad = (design: DesignImage): void => {
  loadingStates[design.id] = false
  errorStates[design.id] = false
  loadingProgress[design.id] = 100
}

/*** 图片加载失败 - Image load error ***/
const onImageError = (design: DesignImage): void => {
  loadingStates[design.id] = false
  errorStates[design.id] = true
}

/*** 重试加载图片 - Retry load image ***/
const retryLoadImage = (design: DesignImage): void => {
  loadingStates[design.id] = true
  errorStates[design.id] = false
  loadingProgress[design.id] = 0
  
  // 强制重新加载图片 - Force reload image
  const index = designList.value.findIndex(d => d.id === design.id)
  if (index !== -1) {
    const updatedDesign = { ...design, imageUrl: design.imageUrl + '?t=' + Date.now() }
    designList.value[index] = updatedDesign
  }
}

/*** 打开预览 - Open preview ***/
const openPreview = (index: number): void => {
  currentIndex.value = index
  scaleValue.value = 1
  showPreview.value = true
}

/*** 关闭预览 - Close preview ***/
const closePreview = (): void => {
  showPreview.value = false
  scaleValue.value = 1
}

/*** 轮播切换 - Swiper change ***/
const onSwiperChange = (e: any): void => {
  const newIndex = e.detail.current
  // 边界处理 - Boundary handling (Property 16)
  if (newIndex >= 0 && newIndex < designList.value.length) {
    currentIndex.value = newIndex
  }
  scaleValue.value = 1
}

/*** 缩放事件 - Scale event ***/
const onScale = (e: any): void => {
  scaleValue.value = e.detail.scale
}

/*** 移动事件 - Move event ***/
const onMoveChange = (_e: any): void => {
  // 可用于记录位置状态 - Can be used to record position state
}

/*** 预览图片加载成功 - Preview image load success ***/
const onPreviewImageLoad = (design: DesignImage): void => {
  previewLoadingStates[design.id] = false
  previewErrorStates[design.id] = false
}

/*** 预览图片加载失败 - Preview image load error ***/
const onPreviewImageError = (design: DesignImage): void => {
  previewLoadingStates[design.id] = false
  previewErrorStates[design.id] = true
}

/*** 返回上一页 - Go back ***/
const goBack = (): void => {
  uni.navigateBack()
}

/*** 页面加载时获取参数 - Get params on page load ***/
onLoad((options) => {
  if (options?.id) {
    projectId.value = parseInt(options.id as string, 10)
  }
})

/*** 组件挂载时加载数据 - Load data on mounted ***/
onMounted(() => {
  loadDesignList()
})

/*** 下拉刷新事件 - Pull down refresh event ***/
onPullDownRefresh(async () => {
  await loadDesignList(true)
  uni.stopPullDownRefresh()
})

/*** 触底加载更多 - Reach bottom load more ***/
onReachBottom(() => {
  loadMore()
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
  padding: 60rpx 32rpx 40rpx;
  overflow: hidden;
}

.header-bg {
  position: absolute;
  top: -100rpx;
  left: -50rpx;
  right: -50rpx;
  height: 350rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 50%, #1D4ED8 100%);
  border-radius: 0 0 60rpx 60rpx;
  transform: rotate(-3deg);
}

.header-content {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
}

.back-btn {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
  color: #FFFFFF;
  margin-right: 24rpx;
  transition: all 0.2s ease;
  
  &:active {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(0.95);
  }
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #FFFFFF;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  flex: 1;
}

.header-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
}

.count-text {
  font-size: 24rpx;
  color: #FFFFFF;
}

/*** 加载状态 - Loading state ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.loading-spinner {
  width: 64rpx;
  height: 64rpx;
  border: 4rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 24rpx;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 28rpx;
  color: #94A3B8;
}

/*** 错误状态 - Error state ***/
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 160rpx 48rpx;
}

.error-icon {
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FEF2F2;
  border-radius: 50%;
  margin-bottom: 32rpx;
  color: #EF4444;
}

.error-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.error-desc {
  font-size: 26rpx;
  color: #64748B;
  text-align: center;
  margin-bottom: 40rpx;
}

.retry-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20rpx 48rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-radius: 40rpx;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.3);
  transition: all 0.2s ease;
  color: #FFFFFF;
  
  &:active {
    transform: scale(0.95);
    opacity: 0.9;
  }
}

.retry-text {
  font-size: 28rpx;
  font-weight: 500;
  color: #FFFFFF;
  margin-left: 12rpx;
}

/*** 空状态 - Empty state ***/
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 160rpx 48rpx;
}

.empty-icon {
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F1F5F9 0%, #E2E8F0 100%);
  border-radius: 50%;
  margin-bottom: 32rpx;
  color: #94A3B8;
}

.empty-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: #94A3B8;
}

/*** 内容滚动区 - Content scroll area ***/
.content-scroll {
  height: calc(100vh - 180rpx);
  padding: 0 24rpx;
}

/*** 设计图网格 - Design grid ***/
.design-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
  padding-top: 16rpx;
}

/*** 设计图卡片 - Design card with 3D effect ***/
.design-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    box-shadow: 
      0 4rpx 16rpx rgba(59, 130, 246, 0.12),
      0 1rpx 4rpx rgba(0, 0, 0, 0.06);
  }
}

/*** 图片容器 - Image wrapper ***/
.image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%; /* 1:1 比例 - 1:1 aspect ratio */
  background: linear-gradient(135deg, #F8FAFC 0%, #F1F5F9 100%);
  overflow: hidden;
}

/*** 设计图图片 - Design image ***/
.design-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/*** 图片加载状态 - Image loading state ***/
.image-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F8FAFC 0%, #F1F5F9 100%);
}

.loading-ring {
  width: 48rpx;
  height: 48rpx;
  border: 4rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 12rpx;
}

.loading-percent {
  font-size: 22rpx;
  color: #94A3B8;
}

/*** 图片加载失败 - Image error state ***/
.image-error {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #FEF2F2 0%, #FEE2E2 100%);
  color: #EF4444;
}

.error-hint {
  font-size: 24rpx;
  color: #EF4444;
  margin-top: 12rpx;
}

.retry-hint {
  font-size: 20rpx;
  color: #F87171;
  margin-top: 8rpx;
}

/*** 卡片阴影 - Card shadow for 3D effect ***/
.card-shadow {
  position: absolute;
  bottom: -8rpx;
  left: 16rpx;
  right: 16rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.08) 0%, transparent 70%);
  border-radius: 50%;
}

/*** 设计图信息 - Design info ***/
.design-info {
  padding: 20rpx;
}

.design-name {
  font-size: 26rpx;
  font-weight: 600;
  color: #1E293B;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 8rpx;
}

.design-date {
  font-size: 22rpx;
  color: #94A3B8;
}

/*** 加载更多 - Load more ***/
.load-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 0;
}

.loading-more-spinner {
  width: 32rpx;
  height: 32rpx;
  border: 3rpx solid #E2E8F0;
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 12rpx;
}

.load-more-text {
  font-size: 24rpx;
  color: #94A3B8;
}

/*** 没有更多 - No more ***/
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
  height: 60rpx;
}

/*** 预览弹窗 - Preview modal ***/
.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.preview-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.95);
}

/*** 预览头部 - Preview header ***/
.preview-header {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  padding: 60rpx 32rpx 24rpx;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.6) 0%, transparent 100%);
}

.preview-close {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  color: #FFFFFF;
  margin-right: 20rpx;
  
  &:active {
    background: rgba(255, 255, 255, 0.2);
  }
}

.preview-title {
  flex: 1;
  font-size: 32rpx;
  font-weight: 600;
  color: #FFFFFF;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-counter {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-left: 16rpx;
}

/*** 预览轮播 - Preview swiper ***/
.preview-swiper {
  flex: 1;
  width: 100%;
}

.movable-area {
  width: 100%;
  height: 100%;
}

.movable-view {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-image {
  width: 100%;
  height: 100%;
}

/*** 预览加载状态 - Preview loading state ***/
.preview-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.preview-loading-spinner {
  width: 64rpx;
  height: 64rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.2);
  border-top-color: #FFFFFF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16rpx;
}

.preview-loading-text {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
}

/*** 预览加载失败 - Preview error state ***/
.preview-error {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  color: rgba(255, 255, 255, 0.5);
}

.preview-error-text {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 16rpx;
}

/*** 预览底部 - Preview footer ***/
.preview-footer {
  position: relative;
  z-index: 10;
  padding: 24rpx 32rpx 60rpx;
  background: linear-gradient(0deg, rgba(0, 0, 0, 0.6) 0%, transparent 100%);
  text-align: center;
}

.preview-hint {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
}
</style>
