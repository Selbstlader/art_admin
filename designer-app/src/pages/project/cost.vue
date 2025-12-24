<template>
  <!-- 
    成本报告页面 - C4D风格高级UI
    Cost Report Page - C4D style premium UI with glass morphism
  -->
  <view class="page-container">
    <!-- 顶部区域 - Header area with blue gradient -->
    <view class="header-area">
      <view class="header-bg"></view>
      <view class="header-content">
        <view class="back-btn" @click="goBack">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6"></path>
          </svg>
        </view>
        <text class="header-title">成本报告</text>
      </view>
    </view>

    <!-- 加载状态 -->
    <view v-if="loading" class="loading-container">
      <view class="loading-spinner"></view>
      <text class="loading-text">加载中...</text>
    </view>

    <!-- 错误状态 -->
    <view v-else-if="loadError" class="state-card">
      <view class="state-icon error">!</view>
      <text class="state-title">加载失败</text>
      <text class="state-desc">{{ errorMessage }}</text>
      <view class="retry-btn" @click="loadCostData">
        <text>重新加载</text>
      </view>
    </view>

    <!-- 空状态 -->
    <view v-else-if="!costEstimate" class="state-card">
      <view class="state-icon empty">$</view>
      <text class="state-title">暂无成本数据</text>
      <text class="state-desc">该项目尚未生成成本估算报告</text>
    </view>

    <!-- 成本报告内容 -->
    <scroll-view v-else class="content-scroll" scroll-y>
      <!-- 总成本卡片 -->
      <view class="glass-card total-card">
        <view class="total-header">
          <text class="total-label">总成本</text>
          <view v-if="costEstimate.budgetExceeded" class="exceed-badge">超支</view>
        </view>
        <text :class="['total-amount', { exceed: costEstimate.budgetExceeded }]">
          {{ formatMoney(costEstimate.totalCost) }}
        </text>
        <view class="budget-row">
          <text class="budget-label">预算上限</text>
          <text class="budget-value">{{ formatMoney(costEstimate.budgetLimit) }}</text>
        </view>
        <view class="progress-bar">
          <view :class="['progress-fill', { exceed: costEstimate.budgetExceeded }]" :style="{ width: progressWidth }">
          </view>
        </view>
        <view class="progress-labels">
          <text>¥0</text>
          <text>{{ formatMoney(costEstimate.budgetLimit) }}</text>
        </view>
      </view>

      <!-- 超预算警告 -->
      <view v-if="costEstimate.budgetExceeded" class="warning-card">
        <view class="warning-icon">!</view>
        <view class="warning-text">
          <text class="warning-title">预算超支提醒</text>
          <text class="warning-desc">超出预算 {{ formatMoney(costEstimate.exceededAmount) }}</text>
        </view>
      </view>

      <!-- 费用分类 -->
      <view class="category-grid">
        <view class="glass-card category-item">
          <view class="cat-icon material">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <path d="m7.5 4.27 9 5.15" />
              <path
                d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z" />
              <path d="m3.3 7 8.7 5 8.7-5" />
              <path d="M12 22V12" />
            </svg>
          </view>
          <text class="cat-label">材料费</text>
          <text class="cat-value">{{ formatMoney(costEstimate.materialCost) }}</text>
        </view>
        <view class="glass-card category-item">
          <view class="cat-icon labor">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </view>
          <text class="cat-label">人工费</text>
          <text class="cat-value">{{ formatMoney(costEstimate.laborCost) }}</text>
        </view>
        <view class="glass-card category-item">
          <view class="cat-icon equipment">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <path
                d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
            </svg>
          </view>
          <text class="cat-label">设备费</text>
          <text class="cat-value">{{ formatMoney(costEstimate.equipmentCost) }}</text>
        </view>
        <view class="glass-card category-item">
          <view class="cat-icon management">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <rect width="20" height="14" x="2" y="5" rx="2" />
              <line x1="2" x2="22" y1="10" y2="10" />
            </svg>
          </view>
          <text class="cat-label">管理费</text>
          <text class="cat-value">{{ formatMoney(costEstimate.managementCost) }}</text>
        </view>
      </view>

      <!-- 自定义费用 -->
      <view v-if="costEstimate.customCosts && costEstimate.customCosts.length > 0" class="glass-card section-card">
        <view class="section-header">
          <text class="section-title">自定义费用</text>
          <text :class="['section-total', { negative: costEstimate.customCostTotal < 0 }]">{{
            formatMoney(costEstimate.customCostTotal) }}</text>
        </view>
        <view class="list-divider"></view>
        <view v-for="(item, idx) in costEstimate.customCosts" :key="idx" class="list-item">
          <text class="item-name">{{ item.name }}</text>
          <text :class="['item-value', { negative: item.amount < 0 }]">{{ formatMoney(item.amount) }}</text>
        </view>
      </view>

      <!-- 材料明细 -->
      <view v-if="costEstimate.materialItems && costEstimate.materialItems.length > 0" class="glass-card section-card">
        <view class="section-header" @click="toggleMaterialList">
          <view class="section-left">
            <text class="section-title">材料明细</text>
            <text class="section-count">{{ costEstimate.materialItems.length }}项</text>
          </view>
          <view :class="['arrow-icon', { expanded: showMaterialList }]">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <path d="m6 9 6 6 6-6" />
            </svg>
          </view>
        </view>
        <view v-if="showMaterialList" class="material-list">
          <view class="list-divider"></view>
          <view v-for="item in costEstimate.materialItems" :key="item.id" class="list-item material">
            <view class="material-info">
              <text class="material-name">{{ item.name }}</text>
              <text class="material-spec">{{ item.quantity }}{{ item.unit }} × {{ formatMoney(item.unitPrice) }}</text>
            </view>
            <text class="item-value">{{ formatMoney(item.totalPrice) }}</text>
          </view>
        </view>
      </view>

      <!-- 成本占比 -->
      <view class="glass-card section-card">
        <text class="section-title">成本占比</text>
        <view class="chart-list">
          <view v-for="item in costBreakdown" :key="item.name" class="chart-item">
            <view class="chart-row">
              <text class="chart-name">{{ item.name }}</text>
              <text :class="['chart-percent', { negative: item.percent < 0 }]">{{ item.percent }}%</text>
            </view>
            <view class="chart-bar-bg">
              <view :class="['chart-bar-fill', item.colorClass]" :style="{ width: Math.abs(item.percent) + '%' }">
              </view>
            </view>
            <text :class="['chart-value', { negative: item.amount < 0 }]">{{ formatMoney(item.amount) }}</text>
          </view>
        </view>
      </view>

      <view class="safe-bottom"></view>
    </scroll-view>
  </view>
</template>


<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { costApi } from '@/api/cost'
import type { CostEstimate } from '@/types/api'

interface CostBreakdownItem {
  name: string
  amount: number
  percent: number
  colorClass: string
}

const projectId = ref<number>(0)
const costEstimate = ref<CostEstimate | null>(null)
const loading = ref(true)
const loadError = ref(false)
const errorMessage = ref('网络异常，请稍后重试')
const showMaterialList = ref(false)

const progressWidth = computed((): string => {
  if (!costEstimate.value || costEstimate.value.budgetLimit <= 0) return '0%'
  const percent = (costEstimate.value.totalCost / costEstimate.value.budgetLimit) * 100
  return Math.min(percent, 100) + '%'
})

const costBreakdown = computed((): CostBreakdownItem[] => {
  if (!costEstimate.value) return []
  const total = costEstimate.value.totalCost || 1
  const items: CostBreakdownItem[] = []

  if (costEstimate.value.materialCost > 0) {
    items.push({ name: '材料费', amount: costEstimate.value.materialCost, percent: Math.round((costEstimate.value.materialCost / total) * 1000) / 10, colorClass: 'bar-material' })
  }
  if (costEstimate.value.laborCost > 0) {
    items.push({ name: '人工费', amount: costEstimate.value.laborCost, percent: Math.round((costEstimate.value.laborCost / total) * 1000) / 10, colorClass: 'bar-labor' })
  }
  if (costEstimate.value.equipmentCost > 0) {
    items.push({ name: '设备费', amount: costEstimate.value.equipmentCost, percent: Math.round((costEstimate.value.equipmentCost / total) * 1000) / 10, colorClass: 'bar-equipment' })
  }
  if (costEstimate.value.managementCost > 0) {
    items.push({ name: '管理费', amount: costEstimate.value.managementCost, percent: Math.round((costEstimate.value.managementCost / total) * 1000) / 10, colorClass: 'bar-management' })
  }
  if (costEstimate.value.customCostTotal !== 0) {
    items.push({ name: '自定义费用', amount: costEstimate.value.customCostTotal, percent: Math.round((costEstimate.value.customCostTotal / total) * 1000) / 10, colorClass: 'bar-custom' })
  }
  return items
})

const loadCostData = async (): Promise<void> => {
  if (!projectId.value) { loadError.value = true; errorMessage.value = '项目ID无效'; return }
  loading.value = true
  loadError.value = false
  try {
    const res = await costApi.getByProjectId(projectId.value)
    costEstimate.value = res.code === 200 && res.data ? res.data : null
  } catch (error: unknown) {
    loadError.value = true
    errorMessage.value = error instanceof Error ? error.message : '网络异常，请稍后重试'
  } finally {
    loading.value = false
  }
}

const toggleMaterialList = (): void => { showMaterialList.value = !showMaterialList.value }

const formatMoney = (amount: number): string => {
  if (typeof amount !== 'number' || isNaN(amount)) return '¥0.00'
  const prefix = amount < 0 ? '-¥' : '¥'
  return prefix + Math.abs(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const goBack = (): void => { uni.navigateBack() }

onLoad((options) => {
  const params = options as { id?: string } | undefined
  if (params?.id) projectId.value = parseInt(params.id, 10)
})

onMounted(() => { loadCostData() })
</script>


<style lang="scss" scoped>
/*** 页面容器 - Page container with pure white background ***/
.page-container {
  min-height: 100vh;
  background: #FFFFFF;
  position: relative;
  overflow: hidden;
  box-sizing: border-box;

  /*** 背景装饰光晕 - Background decorative glow ***/
  &::before {
    content: '';
    position: absolute;
    top: -200rpx;
    right: -150rpx;
    width: 500rpx;
    height: 500rpx;
    background: radial-gradient(circle, rgba(59, 130, 246, 0.08) 0%, transparent 70%);
    border-radius: 50%;
    pointer-events: none;
  }

  &::after {
    content: '';
    position: absolute;
    bottom: 200rpx;
    left: -100rpx;
    width: 400rpx;
    height: 400rpx;
    background: radial-gradient(circle, rgba(16, 185, 129, 0.06) 0%, transparent 70%);
    border-radius: 50%;
    pointer-events: none;
  }
}

/*** 顶部区域 - Header area with blue gradient ***/
.header-area {
  position: relative;
  padding: 30rpx 30rpx 30rpx;
  overflow: hidden;
}

.header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 320rpx;
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  border-radius: 0 0 48rpx 48rpx;

  /*** 3D光泽效果 - 3D glossy effect ***/
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 50%;
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.15) 0%, transparent 100%);
    border-radius: 0 0 48rpx 48rpx;
  }

  /*** 装饰光点 - Decorative light spots ***/
  &::after {
    content: '';
    position: absolute;
    top: 40rpx;
    right: 60rpx;
    width: 120rpx;
    height: 120rpx;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 0%, transparent 60%);
    border-radius: 50%;
  }
}

.header-content {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
}

.back-btn {
  width: 76rpx;
  height: 76rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  border-radius: 22rpx;
  color: #FFFFFF;
  margin-right: 24rpx;

  &:active {
    background: rgba(255, 255, 255, 0.35);
    transform: scale(0.92);
  }
}

.header-title {
  font-size: 38rpx;
  font-weight: 700;
  color: #FFFFFF;
  letter-spacing: 2rpx;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 200rpx 0;
}

.loading-spinner {
  width: 72rpx;
  height: 72rpx;
  border: 4rpx solid rgba(59, 130, 246, 0.15);
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 28rpx;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 28rpx;
  color: #64748B;
}

.state-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 48rpx;
  margin: 0 32rpx;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  border-radius: 28rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.04), inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

.state-icon {
  width: 120rpx;
  height: 120rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 48rpx;
  font-weight: 700;
  margin-bottom: 24rpx;

  &.error {
    background: linear-gradient(135deg, rgba(254, 226, 226, 0.8), rgba(254, 202, 202, 0.6));
    color: #EF4444;
  }

  &.empty {
    background: linear-gradient(135deg, rgba(219, 234, 254, 0.8), rgba(191, 219, 254, 0.6));
    color: #3B82F6;
  }
}

.state-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.state-desc {
  font-size: 26rpx;
  color: #64748B;
}

.retry-btn {
  margin-top: 32rpx;
  padding: 20rpx 48rpx;
  background: linear-gradient(135deg, #4F8EF7 0%, #2563EB 100%);
  border-radius: 40rpx;
  color: #FFFFFF;
  font-size: 28rpx;
  font-weight: 600;

  &:active {
    transform: scale(0.95);
  }
}

.content-scroll {
  height: calc(100vh - 140rpx);
  padding: 0 32rpx;
  box-sizing: border-box;
}

/* 玻璃拟态卡片 */
.glass-card {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 28rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.04), 0 2rpx 8rpx rgba(0, 0, 0, 0.02), inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 50%;
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.4) 0%, transparent 100%);
    border-radius: 28rpx 28rpx 0 0;
    pointer-events: none;
  }
}

.total-card {
  padding: 32rpx;
  margin-bottom: 24rpx;
}

.total-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.total-label {
  font-size: 28rpx;
  color: #64748B;
  font-weight: 500;
}

.exceed-badge {
  padding: 6rpx 16rpx;
  background: linear-gradient(135deg, #FEE2E2, #FECACA);
  border-radius: 20rpx;
  font-size: 22rpx;
  font-weight: 600;
  color: #EF4444;
}

.total-amount {
  font-size: 56rpx;
  font-weight: 800;
  color: #1E293B;

  &.exceed {
    color: #EF4444;
  }
}

.budget-row {
  display: flex;
  align-items: center;
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid rgba(241, 245, 249, 0.8);
}

.budget-label {
  font-size: 24rpx;
  color: #94A3B8;
  margin-right: 12rpx;
}

.budget-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #64748B;
}

.progress-bar {
  height: 12rpx;
  background: rgba(241, 245, 249, 0.8);
  border-radius: 6rpx;
  margin-top: 20rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4F8EF7 0%, #2563EB 100%);
  border-radius: 6rpx;
  transition: width 0.5s ease;

  &.exceed {
    background: linear-gradient(90deg, #F87171 0%, #EF4444 100%);
  }
}

.progress-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #94A3B8;
}

.warning-card {
  display: flex;
  align-items: center;
  padding: 24rpx;
  margin-bottom: 24rpx;
  background: linear-gradient(135deg, rgba(254, 242, 242, 0.9) 0%, rgba(254, 226, 226, 0.8) 100%);
  border-radius: 20rpx;
  border: 1rpx solid rgba(239, 68, 68, 0.2);
}

.warning-icon {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #EF4444;
  border-radius: 50%;
  color: #FFFFFF;
  font-size: 28rpx;
  font-weight: 700;
  margin-right: 16rpx;
}

.warning-text {
  flex: 1;
}

.warning-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #DC2626;
  display: block;
}

.warning-desc {
  font-size: 24rpx;
  color: #B91C1C;
  margin-top: 4rpx;
  display: block;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.category-item {
  padding: 24rpx;
  display: flex;
  flex-direction: column;
}

.cat-icon {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 14rpx;
  margin-bottom: 12rpx;

  &.material {
    background: linear-gradient(135deg, #DBEAFE, #BFDBFE);
    color: #3B82F6;
  }

  &.labor {
    background: linear-gradient(135deg, #FCE7F3, #FBCFE8);
    color: #EC4899;
  }

  &.equipment {
    background: linear-gradient(135deg, #FEF3C7, #FDE68A);
    color: #F59E0B;
  }

  &.management {
    background: linear-gradient(135deg, #E0E7FF, #C7D2FE);
    color: #6366F1;
  }
}

.cat-label {
  font-size: 24rpx;
  color: #64748B;
  margin-bottom: 8rpx;
}

.cat-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
}

.section-card {
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-left {
  display: flex;
  align-items: center;
}

.section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1E293B;
}

.section-count {
  font-size: 24rpx;
  color: #94A3B8;
  margin-left: 12rpx;
  padding: 4rpx 12rpx;
  background: rgba(241, 245, 249, 0.8);
  border-radius: 12rpx;
}

.section-total {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;

  &.negative {
    color: #EF4444;
  }
}

.arrow-icon {
  width: 36rpx;
  height: 36rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  transition: transform 0.3s ease;

  &.expanded {
    transform: rotate(180deg);
  }
}

.list-divider {
  height: 1rpx;
  background: rgba(241, 245, 249, 0.8);
  margin: 16rpx 0;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid rgba(248, 250, 252, 0.8);

  &:last-child {
    border-bottom: none;
  }

  &.material {
    align-items: flex-start;
  }
}

.item-name {
  font-size: 28rpx;
  color: #475569;
}

.item-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;

  &.negative {
    color: #EF4444;
  }
}

.material-list {
  /* 移除内部滚动，让外层scroll-view统一滚动 */
}

.material-info {
  flex: 1;
  min-width: 0;
}

.material-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #1E293B;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.material-spec {
  font-size: 24rpx;
  color: #94A3B8;
  margin-top: 4rpx;
  display: block;
}

.chart-list {
  margin-top: 20rpx;
}

.chart-item {
  margin-bottom: 20rpx;

  &:last-child {
    margin-bottom: 0;
  }
}

.chart-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8rpx;
}

.chart-name {
  font-size: 26rpx;
  color: #475569;
}

.chart-percent {
  font-size: 26rpx;
  font-weight: 600;
  color: #3B82F6;

  &.negative {
    color: #EF4444;
  }
}

.chart-bar-bg {
  height: 16rpx;
  background: rgba(241, 245, 249, 0.8);
  border-radius: 8rpx;
  overflow: hidden;
}

.chart-bar-fill {
  height: 100%;
  border-radius: 8rpx;
  transition: width 0.5s ease;

  &.bar-material {
    background: linear-gradient(90deg, #60A5FA, #3B82F6);
  }

  &.bar-labor {
    background: linear-gradient(90deg, #F472B6, #EC4899);
  }

  &.bar-equipment {
    background: linear-gradient(90deg, #FBBF24, #F59E0B);
  }

  &.bar-management {
    background: linear-gradient(90deg, #818CF8, #6366F1);
  }

  &.bar-custom {
    background: linear-gradient(90deg, #A78BFA, #8B5CF6);
  }
}

.chart-value {
  font-size: 24rpx;
  color: #64748B;
  text-align: right;
  display: block;
  margin-top: 6rpx;

  &.negative {
    color: #EF4444;
  }
}

.safe-bottom {
  height: 100rpx;
}
</style>
