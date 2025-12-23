<template>
  <!-- 
    材料卡片组件 - C4D风格高级UI
    Material Card Component - C4D style premium UI with 3D floating effect
  -->
  <view class="material-card" @click="handleClick">
    <!-- 材料图片区域 - Material image area -->
    <view class="card-image-wrapper">
      <image 
        class="card-image" 
        :src="material.imageUrl || defaultImage"
        mode="aspectFill"
        @error="handleImageError"
      />
      <!-- 图片加载失败占位 - Image error placeholder -->
      <view v-if="imageError" class="image-placeholder">
        <view class="placeholder-icon">
          <!-- Material icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
            <circle cx="9" cy="9" r="2"></circle>
            <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
          </svg>
        </view>
      </view>
    </view>
    
    <!-- 卡片内容区 - Card content area -->
    <view class="card-content">
      <!-- 材料名称 - Material name -->
      <text class="material-name">{{ material.name }}</text>
      
      <!-- 品牌信息 - Brand info -->
      <view class="brand-row">
        <view class="brand-icon">
          <!-- Brand icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"></path>
            <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"></path>
            <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"></path>
            <path d="M10 6h4"></path>
            <path d="M10 10h4"></path>
            <path d="M10 14h4"></path>
            <path d="M10 18h4"></path>
          </svg>
        </view>
        <text class="brand-text">{{ material.brand || '未知品牌' }}</text>
      </view>
      
      <!-- 价格信息 - Price info -->
      <view class="price-row">
        <text class="price-symbol">¥</text>
        <text class="price-value">{{ formattedPrice }}</text>
        <text class="price-unit">/{{ material.unit || '件' }}</text>
      </view>
    </view>
    
    <!-- 卡片3D阴影效果 - Card 3D shadow effect -->
    <view class="card-shadow"></view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * MaterialCard Component
 * 材料卡片组件 - 悬浮卡片设计，显示材料图片、名称、品牌、单价
 * Floating card design showing material image, name, brand, and unit price
 ***/
import { ref, computed } from 'vue'
import type { Material } from '@/types/api'

/*** Component Props - 组件属性定义 ***/
interface Props {
  material: Material
}

const props = defineProps<Props>()

/*** Component Emits - 组件事件定义 ***/
const emit = defineEmits<{
  (e: 'click', id: number): void
}>()

/*** 默认图片 - Default placeholder image ***/
const defaultImage = '/static/images/material-placeholder.png'

/*** 图片加载错误状态 - Image error state ***/
const imageError = ref(false)

/*** 格式化价格 - Format price with thousand separators ***/
const formattedPrice = computed((): string => {
  const price = props.material.price
  if (typeof price !== 'number' || isNaN(price)) return '0.00'
  
  // 保留两位小数并添加千分位分隔符
  return price.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
})

/*** 处理图片加载错误 - Handle image load error ***/
const handleImageError = (): void => {
  imageError.value = true
}

/*** 处理点击事件 - Handle card click event ***/
const handleClick = (): void => {
  emit('click', props.material.id)
}
</script>

<style lang="scss" scoped>
/*** 材料卡片 - Material card with 3D floating effect ***/
.material-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 20rpx;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

/*** 卡片图片区域 - Card image wrapper ***/
.card-image-wrapper {
  position: relative;
  width: 100%;
  height: 240rpx;
  background: #F1F5F9;
  overflow: hidden;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/*** 图片占位符 - Image placeholder ***/
.image-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F1F5F9 0%, #E2E8F0 100%);
}

.placeholder-icon {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
}

/*** 卡片内容区 - Card content area ***/
.card-content {
  padding: 20rpx 24rpx 24rpx;
}

/*** 材料名称 - Material name ***/
.material-name {
  display: block;
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
  line-height: 1.4;
  margin-bottom: 12rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/*** 品牌行 - Brand row ***/
.brand-row {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;
}

.brand-icon {
  width: 28rpx;
  height: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 8rpx;
}

.brand-text {
  font-size: 24rpx;
  color: #64748B;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/*** 价格行 - Price row ***/
.price-row {
  display: flex;
  align-items: baseline;
}

.price-symbol {
  font-size: 24rpx;
  font-weight: 600;
  color: #3B82F6;
  margin-right: 2rpx;
}

.price-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #3B82F6;
}

.price-unit {
  font-size: 22rpx;
  color: #94A3B8;
  margin-left: 4rpx;
}

/*** 卡片3D阴影 - Card 3D shadow ***/
.card-shadow {
  position: absolute;
  bottom: -6rpx;
  left: 16rpx;
  right: 16rpx;
  height: 12rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}
</style>
