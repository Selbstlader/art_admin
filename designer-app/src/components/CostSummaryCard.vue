<template>
  <!-- 
    成本汇总卡片组件 - C4D风格高级UI
    Cost Summary Card Component - C4D style premium UI with 3D data visualization
  -->
  <view class="cost-summary-card">
    <!-- 卡片头部 - Card header with title -->
    <view class="card-header">
      <view class="header-left">
        <view class="header-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="1" x2="12" y2="23"></line>
            <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
          </svg>
        </view>
        <text class="header-title">成本汇总</text>
      </view>
      <!-- 超支警告标签 - Overspend warning badge -->
      <view v-if="isOverBudget" class="warning-badge">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
          <line x1="12" y1="9" x2="12" y2="13"></line>
          <line x1="12" y1="17" x2="12.01" y2="17"></line>
        </svg>
        <text class="warning-text">超支</text>
      </view>
    </view>

    <!-- 3D 立体数据展示区 - 3D data visualization area -->
    <view class="data-visualization">
      <!-- 总成本展示 - Total cost display with 3D effect -->
      <view class="total-cost-section">
        <view class="cost-ring-container">
          <!-- 环形进度条背景 - Ring progress background -->
          <view class="ring-bg"></view>
          <!-- 环形进度条 - Ring progress bar -->
          <view class="ring-progress" :style="ringProgressStyle"></view>
          <!-- 中心数据 - Center data -->
          <view class="ring-center">
            <text class="cost-label">总成本</text>
            <text :class="['cost-value', { 'over-budget': isOverBudget }]">{{ formattedTotalCost }}</text>
          </view>
          <!-- 3D 阴影效果 - 3D shadow effect -->
          <view class="ring-shadow"></view>
        </view>
      </view>

      <!-- 预算对比区 - Budget comparison area -->
      <view class="budget-comparison">
        <view class="comparison-item">
          <view class="comparison-icon budget-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 5c-1.5 0-2.8 1.4-3 2-3.5-1.5-11-.3-11 5 0 1.8 0 3 2 4.5V20h4v-2h3v2h4v-4c1-.5 1.7-1 2-2h2v-4h-2c0-1-.5-1.5-1-2h0V5z"></path>
              <path d="M2 9v1c0 1.1.9 2 2 2h1"></path>
              <path d="M16 11h0"></path>
            </svg>
          </view>
          <view class="comparison-info">
            <text class="comparison-label">预算上限</text>
            <text class="comparison-value">{{ formattedBudgetLimit }}</text>
          </view>
        </view>

        <view class="comparison-divider"></view>

        <view class="comparison-item">
          <view :class="['comparison-icon', isOverBudget ? 'overspend-icon' : 'remain-icon']">
            <svg v-if="isOverBudget" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m3 16 4 4 4-4"></path>
              <path d="M7 20V4"></path>
              <path d="m21 8-4-4-4 4"></path>
              <path d="M17 4v16"></path>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"></path>
              <path d="m9 12 2 2 4-4"></path>
            </svg>
          </view>
          <view class="comparison-info">
            <text class="comparison-label">{{ isOverBudget ? '超支金额' : '剩余预算' }}</text>
            <text :class="['comparison-value', { 'over-budget-value': isOverBudget, 'remain-value': !isOverBudget }]">
              {{ isOverBudget ? '+' : '' }}{{ formattedDifference }}
            </text>
          </view>
        </view>
      </view>
    </view>

    <!-- 进度条展示 - Progress bar display -->
    <view class="progress-section">
      <view class="progress-bar-container">
        <view class="progress-bar-bg"></view>
        <view 
          :class="['progress-bar-fill', { 'over-budget-bar': isOverBudget }]" 
          :style="{ width: progressWidth }"
        ></view>
        <!-- 预算线标记 - Budget line marker -->
        <view v-if="isOverBudget" class="budget-line"></view>
      </view>
      <view class="progress-labels">
        <text class="progress-label">0</text>
        <text class="progress-label">{{ formattedBudgetLimit }}</text>
      </view>
    </view>

    <!-- 卡片3D阴影效果 - Card 3D shadow effect -->
    <view class="card-shadow"></view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * CostSummaryCard Component
 * 成本汇总卡片组件 - 3D立体数据图表设计
 * 3D data chart design with glass morphism effect
 ***/
import { computed } from 'vue'

/*** Component Props - 组件属性定义 ***/
interface Props {
  totalCost: number
  budgetLimit: number
}

const props = withDefaults(defineProps<Props>(), {
  totalCost: 0,
  budgetLimit: 0
})

/*** 是否超支 - Check if over budget ***/
const isOverBudget = computed((): boolean => {
  return props.totalCost > props.budgetLimit && props.budgetLimit > 0
})

/*** 预算差额 - Budget difference ***/
const budgetDifference = computed((): number => {
  return Math.abs(props.totalCost - props.budgetLimit)
})

/*** 格式化总成本 - Format total cost ***/
const formattedTotalCost = computed((): string => {
  return formatMoney(props.totalCost)
})

/*** 格式化预算上限 - Format budget limit ***/
const formattedBudgetLimit = computed((): string => {
  return formatMoney(props.budgetLimit)
})

/*** 格式化差额 - Format difference ***/
const formattedDifference = computed((): string => {
  return formatMoney(budgetDifference.value)
})

/*** 进度百分比 - Progress percentage ***/
const progressPercentage = computed((): number => {
  if (props.budgetLimit <= 0) return 0
  const percentage = (props.totalCost / props.budgetLimit) * 100
  return Math.min(percentage, 100)
})

/*** 进度条宽度 - Progress bar width ***/
const progressWidth = computed((): string => {
  return `${progressPercentage.value}%`
})

/*** 环形进度条样式 - Ring progress style ***/
const ringProgressStyle = computed(() => {
  const percentage = Math.min((props.totalCost / props.budgetLimit) * 100, 100)
  const degree = (percentage / 100) * 360
  const color = isOverBudget.value 
    ? 'linear-gradient(135deg, #EF4444 0%, #DC2626 100%)' 
    : 'linear-gradient(135deg, #3B82F6 0%, #2563EB 100%)'
  
  return {
    background: `conic-gradient(${isOverBudget.value ? '#EF4444' : '#3B82F6'} ${degree}deg, transparent ${degree}deg)`
  }
})

/*** 格式化金额 - Format money with unit ***/
function formatMoney(amount: number): string {
  if (typeof amount !== 'number' || isNaN(amount)) return '¥0'
  
  if (amount >= 100000000) {
    return '¥' + (amount / 100000000).toFixed(2) + '亿'
  } else if (amount >= 10000) {
    return '¥' + (amount / 10000).toFixed(2) + '万'
  } else {
    return '¥' + amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }
}
</script>


<style lang="scss" scoped>
/*** 成本汇总卡片 - Cost summary card with 3D effect ***/
.cost-summary-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 28rpx;
  padding: 32rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

/*** 卡片头部 - Card header ***/
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32rpx;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-icon {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  border-radius: 14rpx;
  color: #3B82F6;
  margin-right: 16rpx;
}

.header-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
}

/*** 超支警告标签 - Warning badge ***/
.warning-badge {
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 20rpx;
  color: #EF4444;
}

.warning-text {
  font-size: 24rpx;
  font-weight: 500;
  color: #EF4444;
  margin-left: 6rpx;
}

/*** 数据可视化区域 - Data visualization area ***/
.data-visualization {
  display: flex;
  align-items: center;
  margin-bottom: 32rpx;
}

/*** 总成本展示区 - Total cost section ***/
.total-cost-section {
  flex: 1;
  display: flex;
  justify-content: center;
}

.cost-ring-container {
  position: relative;
  width: 200rpx;
  height: 200rpx;
}

.ring-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #F1F5F9;
}

.ring-progress {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  mask: radial-gradient(transparent 60%, black 60%);
  -webkit-mask: radial-gradient(transparent 60%, black 60%);
}

.ring-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 120rpx;
  height: 120rpx;
  background: #FFFFFF;
  border-radius: 50%;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.cost-label {
  font-size: 20rpx;
  color: #94A3B8;
  margin-bottom: 4rpx;
}

.cost-value {
  font-size: 24rpx;
  font-weight: 700;
  color: #3B82F6;
  text-align: center;
}

.cost-value.over-budget {
  color: #EF4444;
}

.ring-shadow {
  position: absolute;
  bottom: -12rpx;
  left: 20rpx;
  right: 20rpx;
  height: 24rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.08) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 预算对比区 - Budget comparison area ***/
.budget-comparison {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding-left: 32rpx;
}

.comparison-item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
}

.comparison-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12rpx;
  margin-right: 16rpx;
}

.budget-icon {
  background: linear-gradient(135deg, #E0E7FF 0%, #C7D2FE 100%);
  color: #6366F1;
}

.remain-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.overspend-icon {
  background: linear-gradient(135deg, #FEE2E2 0%, #FECACA 100%);
  color: #EF4444;
}

.comparison-info {
  display: flex;
  flex-direction: column;
}

.comparison-label {
  font-size: 22rpx;
  color: #94A3B8;
  margin-bottom: 4rpx;
}

.comparison-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
}

.over-budget-value {
  color: #EF4444;
}

.remain-value {
  color: #10B981;
}

.comparison-divider {
  height: 1rpx;
  background: #F1F5F9;
  margin: 8rpx 0;
}

/*** 进度条区域 - Progress bar section ***/
.progress-section {
  margin-top: 8rpx;
}

.progress-bar-container {
  position: relative;
  height: 16rpx;
  border-radius: 8rpx;
  overflow: visible;
}

.progress-bar-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #F1F5F9;
  border-radius: 8rpx;
}

.progress-bar-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #60A5FA 0%, #3B82F6 100%);
  border-radius: 8rpx;
  transition: width 0.3s ease;
}

.progress-bar-fill.over-budget-bar {
  background: linear-gradient(90deg, #F87171 0%, #EF4444 100%);
}

.budget-line {
  position: absolute;
  top: -4rpx;
  right: 0;
  width: 4rpx;
  height: 24rpx;
  background: #1E293B;
  border-radius: 2rpx;
}

.progress-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 8rpx;
}

.progress-label {
  font-size: 20rpx;
  color: #94A3B8;
}

/*** 卡片3D阴影 - Card 3D shadow ***/
.card-shadow {
  position: absolute;
  bottom: -10rpx;
  left: 32rpx;
  right: 32rpx;
  height: 20rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}
</style>
