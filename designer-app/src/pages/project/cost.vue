<template>
  <!-- 
    成本报告页面 - C4D风格高级UI
    Cost Report Page - C4D style premium UI with 3D data visualization
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
        <text class="header-title">成本报告</text>
      </view>
    </view>

    <!-- 加载状态 - Loading state -->
    <view v-if="loading" class="loading-container">
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
      <view class="retry-btn" @click="loadCostData">
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
    <view v-else-if="!costEstimate" class="empty-container">
      <view class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="1" x2="12" y2="23"></line>
          <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
        </svg>
      </view>
      <text class="empty-title">暂无成本数据</text>
      <text class="empty-desc">该项目尚未生成成本估算报告</text>
    </view>

    <!-- 成本报告内容 - Cost report content -->
    <scroll-view v-else class="content-scroll" scroll-y>
      <!-- 成本汇总卡片 - Cost summary card -->
      <CostSummaryCard 
        :total-cost="costEstimate.totalCost" 
        :budget-limit="costEstimate.budgetLimit"
      />

      <!-- 成本分类明细 - Cost category details -->
      <view class="category-section">
        <text class="section-title">成本明细</text>
        
        <!-- 分类卡片列表 - Category card list -->
        <view 
          v-for="(category, index) in groupedCategories" 
          :key="category.name"
          class="category-card"
        >
          <!-- 分类头部 - Category header (clickable) -->
          <view class="category-header" @click="toggleCategory(index)">
            <view class="category-left">
              <view :class="['category-icon', `icon-${getCategoryKey(category.name)}`]">
                <svg v-if="category.name === '材料费'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m21 16-4 4-4-4"></path>
                  <path d="M17 20V4"></path>
                  <path d="m3 8 4-4 4 4"></path>
                  <path d="M7 4v16"></path>
                </svg>
                <svg v-else-if="category.name === '人工费'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
                <svg v-else-if="category.name === '设备费'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 6V2H8"></path>
                  <path d="m8 18-4 4V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2Z"></path>
                  <path d="M2 12h2"></path>
                  <path d="M9 11v2"></path>
                  <path d="M15 11v2"></path>
                  <path d="M20 12h2"></path>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M16 20V4a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
                  <rect width="20" height="14" x="2" y="6" rx="2"></rect>
                </svg>
              </view>
              <view class="category-info">
                <text class="category-name">{{ category.name }}</text>
                <text class="category-count">{{ category.items.length }} 项</text>
              </view>
            </view>
            <view class="category-right">
              <text class="category-total">{{ formatMoney(category.total) }}</text>
              <view :class="['expand-icon', { 'expanded': expandedCategories.includes(index) }]">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m6 9 6 6 6-6"></path>
                </svg>
              </view>
            </view>
          </view>

          <!-- 分类明细列表 - Category detail list (expandable) -->
          <view v-if="expandedCategories.includes(index)" class="category-details">
            <view 
              v-for="item in category.items" 
              :key="item.name"
              class="detail-item"
            >
              <view class="detail-left">
                <text class="detail-name">{{ item.name }}</text>
                <text class="detail-spec">{{ item.quantity }} × {{ formatMoney(item.unitPrice) }}</text>
              </view>
              <text class="detail-price">{{ formatMoney(item.totalPrice) }}</text>
            </view>
          </view>

          <!-- 卡片3D阴影 - Card 3D shadow -->
          <view class="category-shadow"></view>
        </view>
      </view>

      <!-- 成本占比图表 - Cost proportion chart -->
      <view class="chart-section">
        <text class="section-title">成本占比</text>
        <view class="chart-card">
          <view class="chart-bars">
            <view 
              v-for="category in groupedCategories" 
              :key="category.name"
              class="chart-bar-item"
            >
              <view class="bar-label">
                <text class="bar-name">{{ category.name }}</text>
                <text class="bar-percentage">{{ getPercentage(category.total) }}%</text>
              </view>
              <view class="bar-container">
                <view 
                  :class="['bar-fill', `bar-${getCategoryKey(category.name)}`]"
                  :style="{ width: getPercentage(category.total) + '%' }"
                ></view>
              </view>
              <text class="bar-value">{{ formatMoney(category.total) }}</text>
            </view>
          </view>
          <!-- 卡片3D阴影 - Card 3D shadow -->
          <view class="chart-shadow"></view>
        </view>
      </view>

      <!-- 底部安全区 - Bottom safe area -->
      <view class="safe-area-bottom"></view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Cost Report Page - C4D style premium UI
 * 成本报告页面 - C4D风格高级UI，3D立体数据图表
 ***/
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { costApi } from '@/api/cost'
import { CostSummaryCard } from '@/components'
import type { CostEstimate, CostItem } from '@/types/api'

/*** 分组后的成本类别 - Grouped cost category ***/
interface GroupedCategory {
  name: string
  items: CostItem[]
  total: number
}

/*** Local state - 本地状态 ***/
const projectId = ref<number>(0)
const costEstimate = ref<CostEstimate | null>(null)
const loading = ref(true)
const loadError = ref(false)
const errorMessage = ref('网络异常，请稍后重试')
const expandedCategories = ref<number[]>([])

/*** 按类别分组的成本数据 - Grouped cost data by category ***/
const groupedCategories = computed((): GroupedCategory[] => {
  if (!costEstimate.value?.items) return []
  
  const groups: Record<string, CostItem[]> = {}
  
  // 按类别分组 - Group by category
  costEstimate.value.items.forEach((item: CostItem) => {
    const category = item.category || '其他'
    if (!groups[category]) {
      groups[category] = []
    }
    groups[category].push(item)
  })
  
  // 转换为数组并计算小计 - Convert to array and calculate subtotal
  const categoryOrder = ['材料费', '人工费', '设备费', '管理费', '其他']
  
  return categoryOrder
    .filter(name => groups[name])
    .map(name => ({
      name,
      items: groups[name],
      total: groups[name].reduce((sum, item) => sum + item.totalPrice, 0)
    }))
})

/*** 加载成本数据 - Load cost data ***/
const loadCostData = async (): Promise<void> => {
  if (!projectId.value) {
    loadError.value = true
    errorMessage.value = '项目ID无效'
    return
  }

  loading.value = true
  loadError.value = false

  try {
    const res = await costApi.getByProjectId(projectId.value)
    if (res.code === 200 && res.data) {
      costEstimate.value = res.data
    } else {
      // 没有成本数据时不显示错误，显示空状态
      costEstimate.value = null
    }
  } catch (error: unknown) {
    console.error('加载成本数据失败:', error)
    loadError.value = true
    
    if (error instanceof Error) {
      errorMessage.value = error.message || '网络异常，请稍后重试'
    } else {
      errorMessage.value = '网络异常，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}

/*** 切换分类展开状态 - Toggle category expansion ***/
const toggleCategory = (index: number): void => {
  const idx = expandedCategories.value.indexOf(index)
  if (idx > -1) {
    expandedCategories.value.splice(idx, 1)
  } else {
    expandedCategories.value.push(index)
  }
}

/*** 获取分类键名 - Get category key for styling ***/
const getCategoryKey = (name: string): string => {
  const keyMap: Record<string, string> = {
    '材料费': 'material',
    '人工费': 'labor',
    '设备费': 'equipment',
    '管理费': 'management'
  }
  return keyMap[name] || 'other'
}

/*** 格式化金额 - Format money with unit ***/
const formatMoney = (amount: number): string => {
  if (typeof amount !== 'number' || isNaN(amount)) return '¥0'
  
  if (amount >= 100000000) {
    return '¥' + (amount / 100000000).toFixed(2) + '亿'
  } else if (amount >= 10000) {
    return '¥' + (amount / 10000).toFixed(2) + '万'
  } else {
    return '¥' + amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }
}

/*** 获取百分比 - Get percentage ***/
const getPercentage = (amount: number): string => {
  if (!costEstimate.value?.totalCost || costEstimate.value.totalCost === 0) return '0'
  return ((amount / costEstimate.value.totalCost) * 100).toFixed(1)
}

/*** 返回上一页 - Go back to previous page ***/
const goBack = (): void => {
  uni.navigateBack()
}

/*** 页面加载时获取参数 - Get params on page load ***/
onLoad((options: { id?: string }) => {
  if (options?.id) {
    projectId.value = parseInt(options.id, 10)
  }
})

/*** 组件挂载时加载数据 - Load data on component mounted ***/
onMounted(() => {
  loadCostData()
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
  background: linear-gradient(135deg, #10B981 0%, #059669 50%, #047857 100%);
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
  border-top-color: #10B981;
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
  background: linear-gradient(135deg, #10B981 0%, #059669 100%);
  border-radius: 40rpx;
  box-shadow: 0 8rpx 24rpx rgba(16, 185, 129, 0.3);
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
  background: #F1F5F9;
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
  color: #64748B;
  text-align: center;
}

/*** 内容滚动区 - Content scroll area ***/
.content-scroll {
  height: calc(100vh - 180rpx);
  padding: 0 32rpx;
}

/*** 分类区域 - Category section ***/
.category-section {
  margin-top: 32rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 20rpx;
  display: block;
}

/*** 分类卡片 - Category card ***/
.category-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 24rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28rpx;
  transition: background 0.2s ease;
  
  &:active {
    background: #F8FAFC;
  }
}

.category-left {
  display: flex;
  align-items: center;
}

.category-icon {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16rpx;
  margin-right: 16rpx;
}

.icon-material {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.icon-labor {
  background: linear-gradient(135deg, #FCE7F3 0%, #FBCFE8 100%);
  color: #EC4899;
}

.icon-equipment {
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  color: #F59E0B;
}

.icon-management {
  background: linear-gradient(135deg, #E0E7FF 0%, #C7D2FE 100%);
  color: #6366F1;
}

.icon-other {
  background: linear-gradient(135deg, #F1F5F9 0%, #E2E8F0 100%);
  color: #64748B;
}

.category-info {
  display: flex;
  flex-direction: column;
}

.category-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 4rpx;
}

.category-count {
  font-size: 22rpx;
  color: #94A3B8;
}

.category-right {
  display: flex;
  align-items: center;
}

.category-total {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
  margin-right: 12rpx;
}

.expand-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  transition: transform 0.3s ease;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

/*** 分类明细 - Category details ***/
.category-details {
  border-top: 1rpx solid #F1F5F9;
  padding: 0 28rpx;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #F8FAFC;
  
  &:last-child {
    border-bottom: none;
  }
}

.detail-left {
  display: flex;
  flex-direction: column;
}

.detail-name {
  font-size: 28rpx;
  color: #1E293B;
  margin-bottom: 4rpx;
}

.detail-spec {
  font-size: 22rpx;
  color: #94A3B8;
}

.detail-price {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
}

.category-shadow {
  position: absolute;
  bottom: -8rpx;
  left: 24rpx;
  right: 24rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 图表区域 - Chart section ***/
.chart-section {
  margin-top: 32rpx;
}

.chart-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.chart-bars {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.chart-bar-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.bar-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.bar-name {
  font-size: 26rpx;
  color: #1E293B;
}

.bar-percentage {
  font-size: 24rpx;
  font-weight: 600;
  color: #64748B;
}

.bar-container {
  height: 16rpx;
  background: #F1F5F9;
  border-radius: 8rpx;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 8rpx;
  transition: width 0.5s ease;
}

.bar-material {
  background: linear-gradient(90deg, #60A5FA 0%, #3B82F6 100%);
}

.bar-labor {
  background: linear-gradient(90deg, #F472B6 0%, #EC4899 100%);
}

.bar-equipment {
  background: linear-gradient(90deg, #FBBF24 0%, #F59E0B 100%);
}

.bar-management {
  background: linear-gradient(90deg, #818CF8 0%, #6366F1 100%);
}

.bar-other {
  background: linear-gradient(90deg, #94A3B8 0%, #64748B 100%);
}

.bar-value {
  font-size: 24rpx;
  color: #64748B;
  text-align: right;
}

.chart-shadow {
  position: absolute;
  bottom: -8rpx;
  left: 24rpx;
  right: 24rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}
</style>
