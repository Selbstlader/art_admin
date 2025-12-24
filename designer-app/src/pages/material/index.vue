<template>
  <!-- 
    材料分类页面 - C4D风格高级UI
    Material Categories Page - C4D style premium UI
  -->
  <view class="page-container">
    <!-- 页面头部 - Page header with gradient background -->
    <view class="page-header">
      <view class="header-content">
        <text class="header-title">材料库</text>
        <text class="header-subtitle">快速查询材料信息</text>
      </view>
      <!-- 装饰元素 - Decorative elements -->
      <view class="header-decoration">
        <view class="decoration-circle circle-1"></view>
        <view class="decoration-circle circle-2"></view>
      </view>
    </view>
    
    <!-- 分类列表 - Category list -->
    <view class="category-section">
      <view class="section-title">
        <text class="title-text">材料分类</text>
        <text class="title-count">共 {{ categories.length }} 个分类</text>
      </view>
      
      <!-- 加载状态 - Loading state -->
      <view v-if="loading" class="loading-container">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>
      
      <!-- 分类网格 - Category grid -->
      <view v-else-if="categories.length > 0" class="category-grid">
        <view 
          v-for="(category, index) in categories" 
          :key="category.name"
          class="category-card"
          :style="{ animationDelay: `${index * 0.05}s` }"
          @click="handleCategoryClick(category)"
        >
          <!-- 卡片图标 - Card icon -->
          <view class="card-icon-wrapper">
            <view class="card-icon">
              <component :is="getCategoryIcon(category.name)" />
            </view>
          </view>
          
          <!-- 卡片内容 - Card content -->
          <view class="card-info">
            <text class="card-name">{{ category.name }}</text>
            <text class="card-count">{{ category.count }} 种材料</text>
          </view>
          
          <!-- 箭头图标 - Arrow icon -->
          <view class="card-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
          
          <!-- 卡片3D阴影 - Card 3D shadow -->
          <view class="card-shadow"></view>
        </view>
      </view>
      
      <!-- 空状态 - Empty state -->
      <view v-else class="empty-container">
        <view class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48"></path>
          </svg>
        </view>
        <text class="empty-text">暂无材料分类</text>
        <view class="retry-btn" @click="loadCategories">
          <text class="retry-text">点击重试</text>
        </view>
      </view>
    </view>
    
    <!-- 错误提示 - Error toast -->
    <view v-if="errorMsg" class="error-toast">
      <text class="error-text">{{ errorMsg }}</text>
    </view>
    
    <!-- 自定义底部导航栏 - Custom TabBar -->
    <CustomTabBar />
  </view>
</template>

<script setup lang="ts">
/*** 
 * Material Categories Page
 * 材料分类页面 - 展示材料分类列表，点击分类跳转材料列表
 * Display material category list, click to navigate to material list
 ***/
import { ref, onMounted } from 'vue'
import { materialApi } from '@/api/material'
import type { MaterialCategory } from '@/types/api'
import CustomTabBar from '@/components/CustomTabBar.vue'

/*** 分类列表 - Category list ***/
const categories = ref<MaterialCategory[]>([])

/*** 加载状态 - Loading state ***/
const loading = ref(false)

/*** 错误信息 - Error message ***/
const errorMsg = ref('')

/*** 加载分类数据 - Load category data ***/
const loadCategories = async (): Promise<void> => {
  loading.value = true
  errorMsg.value = ''
  
  try {
    const res = await materialApi.getCategories()
    if (res.code === 200 && res.data) {
      categories.value = res.data
    } else {
      // 使用默认分类数据
      categories.value = getDefaultCategories()
    }
  } catch (error) {
    console.error('Failed to load categories:', error)
    // 加载失败时使用默认分类
    categories.value = getDefaultCategories()
    showError('加载分类失败，显示默认分类')
  } finally {
    loading.value = false
  }
}

/*** 获取默认分类 - Get default categories ***/
const getDefaultCategories = (): MaterialCategory[] => {
  return [
    { name: '地板', count: 0, icon: 'floor' },
    { name: '墙面', count: 0, icon: 'wall' },
    { name: '天花', count: 0, icon: 'ceiling' },
    { name: '家具', count: 0, icon: 'furniture' },
    { name: '灯具', count: 0, icon: 'light' },
    { name: '卫浴', count: 0, icon: 'bathroom' },
    { name: '门窗', count: 0, icon: 'door' },
    { name: '其他', count: 0, icon: 'other' }
  ]
}

/*** 获取分类图标组件 - Get category icon component ***/
const getCategoryIcon = (name: string): object => {
  const iconMap: Record<string, object> = {
    '地板': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"></rect><path d="M3 9h18"></path><path d="M3 15h18"></path><path d="M9 3v18"></path><path d="M15 3v18"></path></svg>`
    },
    '墙面': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect><line x1="3" x2="21" y1="9" y2="9"></line><line x1="3" x2="21" y1="15" y2="15"></line><line x1="12" x2="12" y1="3" y2="9"></line><line x1="8" x2="8" y1="9" y2="15"></line><line x1="16" x2="16" y1="9" y2="15"></line><line x1="12" x2="12" y1="15" y2="21"></line></svg>`
    },
    '天花': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v8"></path><path d="m4.93 10.93 1.41 1.41"></path><path d="M2 18h2"></path><path d="M20 18h2"></path><path d="m19.07 10.93-1.41 1.41"></path><path d="M22 22H2"></path><path d="m8 22 4-10 4 10"></path></svg>`
    },
    '家具': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M19 9V6a2 2 0 0 0-2-2H7a2 2 0 0 0-2 2v3"></path><path d="M3 16a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5a2 2 0 0 0-4 0v1.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V11a2 2 0 0 0-4 0z"></path><path d="M5 18v2"></path><path d="M19 18v2"></path></svg>`
    },
    '灯具': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"></path><path d="M9 18h6"></path><path d="M10 22h4"></path></svg>`
    },
    '卫浴': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12h16a1 1 0 0 1 1 1v3a4 4 0 0 1-4 4H7a4 4 0 0 1-4-4v-3a1 1 0 0 1 1-1z"></path><path d="M6 12V5a2 2 0 0 1 2-2h3v2.25"></path><path d="m4 21 1-1.5"></path><path d="m20 21-1-1.5"></path></svg>`
    },
    '门窗': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M13 4h3a2 2 0 0 1 2 2v14"></path><path d="M2 20h3"></path><path d="M13 20h9"></path><path d="M10 12v.01"></path><path d="M13 4.562v16.157a1 1 0 0 1-1.242.97L5 20V5.562a2 2 0 0 1 1.515-1.94l4-1A2 2 0 0 1 13 4.561Z"></path></svg>`
    },
    '其他': {
      template: `<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"></path><path d="m3.3 7 8.7 5 8.7-5"></path><path d="M12 22V12"></path></svg>`
    }
  }
  
  return iconMap[name] || iconMap['其他']
}

/*** 处理分类点击 - Handle category click ***/
const handleCategoryClick = (category: MaterialCategory): void => {
  uni.navigateTo({
    url: `/pages/material/list?category=${encodeURIComponent(category.name)}`
  })
}

/*** 显示错误提示 - Show error message ***/
const showError = (msg: string): void => {
  errorMsg.value = msg
  setTimeout(() => {
    errorMsg.value = ''
  }, 3000)
}

/*** 页面加载 - Page mounted ***/
onMounted(() => {
  loadCategories()
})
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container ***/
.page-container {
  min-height: 100vh;
  background: #F5F7FA;
  padding-bottom: env(safe-area-inset-bottom);
}

/*** 页面头部 - Page header with gradient ***/
.page-header {
  position: relative;
  padding: 48rpx 32rpx 64rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 50%, #93C5FD 100%);
  overflow: hidden;
}

.header-content {
  position: relative;
  z-index: 1;
}

.header-title {
  display: block;
  font-size: 44rpx;
  font-weight: 700;
  color: #FFFFFF;
  margin-bottom: 8rpx;
}

.header-subtitle {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.85);
}

/*** 装饰元素 - Decorative elements ***/
.header-decoration {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 200rpx;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.circle-1 {
  width: 160rpx;
  height: 160rpx;
  top: -40rpx;
  right: -40rpx;
}

.circle-2 {
  width: 100rpx;
  height: 100rpx;
  bottom: 20rpx;
  right: 60rpx;
}

/*** 分类区域 - Category section ***/
.category-section {
  margin-top: -32rpx;
  padding: 0 24rpx;
  position: relative;
  z-index: 2;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 8rpx;
}

.title-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
}

.title-count {
  font-size: 24rpx;
  color: #94A3B8;
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

/*** 分类网格 - Category grid ***/
.category-grid {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

/*** 分类卡片 - Category card ***/
.category-card {
  position: relative;
  display: flex;
  align-items: center;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx 24rpx;
  animation: fadeInUp 0.4s ease forwards;
  opacity: 0;
  transition: transform 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
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

/*** 卡片图标 - Card icon wrapper ***/
.card-icon-wrapper {
  width: 88rpx;
  height: 88rpx;
  border-radius: 20rpx;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
  flex-shrink: 0;
}

.card-icon {
  color: #3B82F6;
  display: flex;
  align-items: center;
  justify-content: center;
}

/*** 卡片信息 - Card info ***/
.card-info {
  flex: 1;
  min-width: 0;
}

.card-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 6rpx;
}

.card-count {
  display: block;
  font-size: 24rpx;
  color: #94A3B8;
}

/*** 箭头图标 - Arrow icon ***/
.card-arrow {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
  flex-shrink: 0;
  transition: transform 0.2s ease;
}

.category-card:active .card-arrow {
  transform: translateX(4rpx);
}

/*** 卡片3D阴影 - Card 3D shadow ***/
.card-shadow {
  position: absolute;
  bottom: -6rpx;
  left: 24rpx;
  right: 24rpx;
  height: 12rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
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
