<template>
  <!-- 
    材料详情页面 - C4D风格高级UI
    Material Detail Page - C4D style premium UI with image gallery
  -->
  <view class="page-container">
    <!-- 加载状态 - Loading state -->
    <view v-if="loading" class="loading-container">
      <view class="loading-spinner"></view>
      <text class="loading-text">加载中...</text>
    </view>
    
    <!-- 错误状态 - Error state -->
    <view v-else-if="loadError" class="error-container">
      <view class="error-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" x2="12" y1="8" y2="12"></line>
          <line x1="12" x2="12.01" y1="16" y2="16"></line>
        </svg>
      </view>
      <text class="error-text">{{ loadError }}</text>
      <view class="retry-btn" @click="loadMaterialDetail">
        <text class="retry-text">点击重试</text>
      </view>
    </view>
    
    <!-- 材料详情内容 - Material detail content -->
    <view v-else-if="material" class="detail-content">
      <!-- 图片轮播区域 - Image swiper area -->
      <view class="image-section">
        <swiper 
          class="image-swiper"
          :indicator-dots="images.length > 1"
          indicator-color="rgba(255,255,255,0.5)"
          indicator-active-color="#FFFFFF"
          :autoplay="false"
          :current="currentImageIndex"
          @change="handleSwiperChange"
        >
          <swiper-item 
            v-for="(img, index) in images" 
            :key="index"
            @click="previewImage(index)"
          >
            <image 
              class="swiper-image" 
              :src="img"
              mode="aspectFill"
              @error="handleImageError(index)"
            />
          </swiper-item>
        </swiper>
        
        <!-- 图片计数器 - Image counter -->
        <view v-if="images.length > 1" class="image-counter">
          <text class="counter-text">{{ currentImageIndex + 1 }}/{{ images.length }}</text>
        </view>
        
        <!-- 点击放大提示 - Click to zoom hint -->
        <view class="zoom-hint">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
            <path d="M11 8v6"></path>
            <path d="M8 11h6"></path>
          </svg>
          <text class="hint-text">点击放大</text>
        </view>
      </view>
      
      <!-- 基本信息卡片 - Basic info card -->
      <view class="info-card">
        <!-- 材料名称和价格 - Name and price -->
        <view class="name-price-row">
          <view class="name-section">
            <text class="material-name">{{ material.name }}</text>
            <view class="category-tag">
              <text class="tag-text">{{ material.category }}</text>
            </view>
          </view>
          <view class="price-section">
            <text class="price-symbol">¥</text>
            <text class="price-value">{{ formattedPrice }}</text>
            <text class="price-unit">/{{ material.unit || '件' }}</text>
          </view>
        </view>
        
        <!-- 品牌信息 - Brand info -->
        <view class="brand-row">
          <view class="brand-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"></path>
              <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"></path>
              <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"></path>
            </svg>
          </view>
          <text class="brand-label">品牌：</text>
          <text class="brand-value">{{ material.brand || '未知品牌' }}</text>
        </view>
      </view>
      
      <!-- 规格信息卡片 - Specification card -->
      <view class="spec-card">
        <view class="card-header">
          <view class="header-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" x2="8" y1="13" y2="13"></line>
              <line x1="16" x2="8" y1="17" y2="17"></line>
              <line x1="10" x2="8" y1="9" y2="9"></line>
            </svg>
          </view>
          <text class="header-title">规格参数</text>
        </view>
        
        <view class="spec-content">
          <view class="spec-item">
            <text class="spec-label">规格型号</text>
            <text class="spec-value">{{ material.specification || '-' }}</text>
          </view>
          <view class="spec-item">
            <text class="spec-label">计量单位</text>
            <text class="spec-value">{{ material.unit || '-' }}</text>
          </view>
        </view>
      </view>
      
      <!-- 供应商信息卡片 - Supplier card -->
      <view class="supplier-card">
        <view class="card-header">
          <view class="header-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 9h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9Z"></path>
              <path d="m3 9 2.45-4.9A2 2 0 0 1 7.24 3h9.52a2 2 0 0 1 1.8 1.1L21 9"></path>
              <path d="M12 3v6"></path>
            </svg>
          </view>
          <text class="header-title">供应商信息</text>
        </view>
        
        <view class="supplier-content">
          <text class="supplier-name">{{ material.supplier || '暂无供应商信息' }}</text>
        </view>
      </view>
      
      <!-- 适用场景卡片 - Applicable scenes card -->
      <view v-if="material.applicableScenes && material.applicableScenes.length > 0" class="scenes-card">
        <view class="card-header">
          <view class="header-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
          </view>
          <text class="header-title">适用场景</text>
        </view>
        
        <view class="scenes-content">
          <view 
            v-for="(scene, index) in material.applicableScenes" 
            :key="index"
            class="scene-tag"
          >
            <text class="scene-text">{{ scene }}</text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 底部安全区域 - Bottom safe area -->
    <view class="bottom-safe-area"></view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Material Detail Page
 * 材料详情页面 - 展示材料详情，图片支持点击放大
 * Display material details with image preview support
 ***/
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { materialApi } from '@/api/material'
import type { MaterialDetail } from '@/types/api'

/*** 材料ID - Material ID from route params ***/
const materialId = ref<number>(0)

/*** 材料详情 - Material detail data ***/
const material = ref<MaterialDetail | null>(null)

/*** 加载状态 - Loading state ***/
const loading = ref(false)

/*** 加载错误 - Load error message ***/
const loadError = ref('')

/*** 当前图片索引 - Current image index ***/
const currentImageIndex = ref(0)

/*** 图片加载错误状态 - Image error states ***/
const imageErrors = ref<Set<number>>(new Set())

/*** 图片列表 - Image list (combine main image and additional images) ***/
const images = computed((): string[] => {
  if (!material.value) return []
  
  const imgList: string[] = []
  
  // 添加主图
  if (material.value.imageUrl) {
    imgList.push(material.value.imageUrl)
  }
  
  // 添加附加图片
  if (material.value.images && material.value.images.length > 0) {
    material.value.images.forEach((img: string) => {
      if (img && !imgList.includes(img)) {
        imgList.push(img)
      }
    })
  }
  
  // 如果没有图片，返回占位图
  if (imgList.length === 0) {
    imgList.push('/static/images/material-placeholder.png')
  }
  
  return imgList
})

/*** 格式化价格 - Format price with thousand separators ***/
const formattedPrice = computed((): string => {
  if (!material.value) return '0.00'
  const price = material.value.price
  if (typeof price !== 'number' || isNaN(price)) return '0.00'
  return price.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
})

/*** 页面加载 - Page load with route params ***/
onLoad((options) => {
  const params = options as { id?: string } | undefined
  if (params?.id) {
    materialId.value = parseInt(params.id, 10)
    loadMaterialDetail()
  } else {
    loadError.value = '材料ID无效'
  }
})

/*** 加载材料详情 - Load material detail ***/
const loadMaterialDetail = async (): Promise<void> => {
  if (!materialId.value) {
    loadError.value = '材料ID无效'
    return
  }
  
  loading.value = true
  loadError.value = ''
  
  try {
    const res = await materialApi.getDetail(materialId.value)
    
    if (res.code === 200 && res.data) {
      material.value = res.data
      // 设置导航栏标题
      uni.setNavigationBarTitle({
        title: res.data.name || '材料详情'
      })
    } else {
      loadError.value = res.msg || '加载失败'
    }
  } catch (error) {
    console.error('Failed to load material detail:', error)
    loadError.value = '网络错误，请稍后重试'
  } finally {
    loading.value = false
  }
}

/*** 处理轮播图切换 - Handle swiper change ***/
const handleSwiperChange = (e: { detail: { current: number } }): void => {
  currentImageIndex.value = e.detail.current
}

/*** 处理图片加载错误 - Handle image load error ***/
const handleImageError = (index: number): void => {
  imageErrors.value.add(index)
}

/*** 预览图片 - Preview image with full screen ***/
const previewImage = (index: number): void => {
  if (images.value.length === 0) return
  
  uni.previewImage({
    current: index,
    urls: images.value,
    indicator: 'number',
    loop: true
  })
}
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container ***/
.page-container {
  min-height: 100vh;
  background: #F5F7FA;
}

/*** 加载状态 - Loading state ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
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

/*** 错误状态 - Error state ***/
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 48rpx;
}

.error-icon {
  width: 128rpx;
  height: 128rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #F87171;
  margin-bottom: 24rpx;
}

.error-text {
  font-size: 28rpx;
  color: #64748B;
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

/*** 详情内容 - Detail content ***/
.detail-content {
  padding-bottom: 32rpx;
}

/*** 图片区域 - Image section ***/
.image-section {
  position: relative;
  width: 100%;
  height: 560rpx;
  background: #FFFFFF;
}

.image-swiper {
  width: 100%;
  height: 100%;
}

.swiper-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/*** 图片计数器 - Image counter ***/
.image-counter {
  position: absolute;
  bottom: 24rpx;
  right: 24rpx;
  padding: 8rpx 16rpx;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 20rpx;
}

.counter-text {
  font-size: 22rpx;
  color: #FFFFFF;
}

/*** 放大提示 - Zoom hint ***/
.zoom-hint {
  position: absolute;
  bottom: 24rpx;
  left: 24rpx;
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 20rpx;
  color: #FFFFFF;
}

.hint-text {
  font-size: 22rpx;
  margin-left: 8rpx;
}

/*** 基本信息卡片 - Info card ***/
.info-card {
  margin: 24rpx;
  padding: 28rpx;
  background: #FFFFFF;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.name-price-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20rpx;
}

.name-section {
  flex: 1;
  margin-right: 24rpx;
}

.material-name {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
  color: #1E293B;
  line-height: 1.4;
  margin-bottom: 12rpx;
}

.category-tag {
  display: inline-flex;
  padding: 6rpx 16rpx;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border-radius: 16rpx;
}

.tag-text {
  font-size: 22rpx;
  color: #3B82F6;
  font-weight: 500;
}

.price-section {
  display: flex;
  align-items: baseline;
  flex-shrink: 0;
}

.price-symbol {
  font-size: 28rpx;
  font-weight: 600;
  color: #EF4444;
}

.price-value {
  font-size: 44rpx;
  font-weight: 700;
  color: #EF4444;
}

.price-unit {
  font-size: 24rpx;
  color: #94A3B8;
  margin-left: 4rpx;
}

.brand-row {
  display: flex;
  align-items: center;
  padding-top: 20rpx;
  border-top: 1rpx solid #F1F5F9;
}

.brand-icon {
  width: 36rpx;
  height: 36rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 8rpx;
}

.brand-label {
  font-size: 26rpx;
  color: #64748B;
}

.brand-value {
  font-size: 26rpx;
  color: #1E293B;
  font-weight: 500;
}

/*** 规格卡片 - Spec card ***/
.spec-card,
.supplier-card,
.scenes-card {
  margin: 0 24rpx 24rpx;
  padding: 24rpx;
  background: #FFFFFF;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #F1F5F9;
}

.header-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3B82F6;
  margin-right: 12rpx;
}

.header-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
}

.spec-content {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.spec-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.spec-label {
  font-size: 26rpx;
  color: #64748B;
}

.spec-value {
  font-size: 26rpx;
  color: #1E293B;
  font-weight: 500;
}

/*** 供应商卡片 - Supplier card ***/
.supplier-content {
  padding: 4rpx 0;
}

.supplier-name {
  font-size: 28rpx;
  color: #1E293B;
  line-height: 1.6;
}

/*** 适用场景卡片 - Scenes card ***/
.scenes-content {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.scene-tag {
  padding: 12rpx 24rpx;
  background: #F1F5F9;
  border-radius: 24rpx;
}

.scene-text {
  font-size: 24rpx;
  color: #475569;
}

/*** 底部安全区域 - Bottom safe area ***/
.bottom-safe-area {
  height: env(safe-area-inset-bottom);
}
</style>
