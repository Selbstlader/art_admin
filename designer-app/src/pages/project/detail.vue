<template>
  <!-- 
    项目详情页面 - C4D风格高级UI
    Project detail page - C4D style premium UI with 3D floating cards
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
        <text class="header-title">项目详情</text>
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
      <view class="retry-btn" @click="loadProjectDetail">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
          <path d="M3 3v5h5"></path>
          <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
          <path d="M16 16h5v5"></path>
        </svg>
        <text class="retry-text">重新加载</text>
      </view>
    </view>

    <!-- 项目详情内容 - Project detail content -->
    <scroll-view v-else-if="project" class="content-scroll" scroll-y>
      <!-- 项目基本信息卡片 - Project basic info card -->
      <view class="info-card main-card">
        <view class="card-header">
          <view class="project-title-row">
            <text class="project-name">{{ project.name }}</text>
            <view :class="['status-badge', `status-${project.status}`]">
              <view class="status-dot"></view>
              <text class="status-text">{{ statusText }}</text>
            </view>
          </view>
          <text v-if="project.description" class="project-desc">{{ project.description }}</text>
        </view>

        <!-- 项目数据网格 - Project data grid with 3D effect -->
        <view class="data-grid">
          <view class="data-cell">
            <view class="data-icon-wrapper area-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect width="18" height="18" x="3" y="3" rx="2"></rect>
                <path d="M3 9h18"></path>
                <path d="M9 21V9"></path>
              </svg>
            </view>
            <view class="data-info">
              <text class="data-label">面积</text>
              <text class="data-value">{{ project.area }}<text class="data-unit">㎡</text></text>
            </view>
          </view>

          <view class="data-cell">
            <view class="data-icon-wrapper budget-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
              </svg>
            </view>
            <view class="data-info">
              <text class="data-label">预算</text>
              <text class="data-value budget-value">{{ formattedBudget }}</text>
            </view>
          </view>

          <view class="data-cell">
            <view class="data-icon-wrapper style-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"></path>
              </svg>
            </view>
            <view class="data-info">
              <text class="data-label">风格</text>
              <text class="data-value">{{ project.style || '未设置' }}</text>
            </view>
          </view>

          <view class="data-cell">
            <view class="data-icon-wrapper time-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
            </view>
            <view class="data-info">
              <text class="data-label">更新时间</text>
              <text class="data-value">{{ formattedDate }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 项目统计卡片 - Project statistics cards with 3D effect -->
      <view class="stats-section">
        <text class="section-title">项目统计</text>
        <view class="stats-grid">
          <!-- 文档数量 - Document count -->
          <view class="stat-card">
            <view class="stat-icon-wrapper doc-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path>
                <path d="M14 2v4a2 2 0 0 0 2 2h4"></path>
                <path d="M10 9H8"></path>
                <path d="M16 13H8"></path>
                <path d="M16 17H8"></path>
              </svg>
            </view>
            <text class="stat-value">{{ project.documentCount || 0 }}</text>
            <text class="stat-label">文档数量</text>
            <view class="stat-shadow"></view>
          </view>

          <!-- 设计图数量 - Design count -->
          <view class="stat-card">
            <view class="stat-icon-wrapper design-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
                <circle cx="9" cy="9" r="2"></circle>
                <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
              </svg>
            </view>
            <text class="stat-value">{{ project.designCount || 0 }}</text>
            <text class="stat-label">设计图</text>
            <view class="stat-shadow"></view>
          </view>

          <!-- 最近分析时间 - Last analysis time -->
          <view class="stat-card analysis-card">
            <view class="stat-icon-wrapper analysis-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 3v18h18"></path>
                <path d="m19 9-5 5-4-4-3 3"></path>
              </svg>
            </view>
            <text class="stat-value-small">{{ formattedAnalysisTime }}</text>
            <text class="stat-label">最近分析</text>
            <view class="stat-shadow"></view>
          </view>
        </view>
      </view>

      <!-- 快捷操作卡片 - Quick action cards -->
      <view class="actions-section">
        <text class="section-title">快捷操作</text>
        
        <!-- 查看文档分析 - View document analysis -->
        <view class="action-card" @click="goToDocuments">
          <view class="action-left">
            <view class="action-icon-wrapper doc-action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
                <polyline points="10 9 9 9 8 9"></polyline>
              </svg>
            </view>
            <view class="action-info">
              <text class="action-title">查看文档分析</text>
              <text class="action-desc">查看项目文档的AI分析结果</text>
            </view>
          </view>
          <view class="action-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>

        <!-- 查看成本预算 - View cost estimate -->
        <view class="action-card" @click="goToCost">
          <view class="action-left">
            <view class="action-icon-wrapper cost-action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="12" y1="1" x2="12" y2="23"></line>
                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
              </svg>
            </view>
            <view class="action-info">
              <text class="action-title">查看成本预算</text>
              <text class="action-desc">查看项目成本估算报告</text>
            </view>
          </view>
          <view class="action-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>

        <!-- 查看设计图 - View designs -->
        <view class="action-card" @click="goToDesigns">
          <view class="action-left">
            <view class="action-icon-wrapper design-action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
                <circle cx="9" cy="9" r="2"></circle>
                <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
              </svg>
            </view>
            <view class="action-info">
              <text class="action-title">查看设计图</text>
              <text class="action-desc">浏览项目设计图纸</text>
            </view>
          </view>
          <view class="action-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>
      </view>

      <!-- 底部安全区 - Bottom safe area -->
      <view class="safe-area-bottom"></view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Project detail page - C4D style premium UI
 * 项目详情页面 - C4D风格高级UI，3D立体数据图表
 ***/
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useProjectStore } from '@/store/project'
import type { ProjectDetail } from '@/types/api'
import { ProjectStatusText, ProjectStatus, ErrorMessages } from '@/types/common'
import { formatDate, formatRelativeTime } from '@/utils/format'

/*** Store instance - 状态管理实例 ***/
const projectStore = useProjectStore()

/*** Local state - 本地状态 ***/
const projectId = ref<number>(0)
const project = ref<ProjectDetail | null>(null)
const loading = ref(true)
const loadError = ref(false)
const errorMessage = ref('网络异常，请稍后重试')

/*** 状态文本 - Status text mapping ***/
const statusText = computed((): string => {
  if (!project.value) return ''
  return ProjectStatusText[project.value.status as ProjectStatus] || project.value.status
})

/*** 格式化预算 - Format budget with unit ***/
const formattedBudget = computed((): string => {
  if (!project.value) return '¥0'
  const budget = project.value.budget
  if (budget >= 10000) {
    return '¥' + (budget / 10000).toFixed(1) + '万'
  }
  return '¥' + budget.toLocaleString()
})

/*** 格式化日期 - Format date ***/
const formattedDate = computed((): string => {
  if (!project.value?.updatedAt) return '-'
  return formatDate(project.value.updatedAt, 'YYYY-MM-DD')
})

/*** 格式化分析时间 - Format analysis time ***/
const formattedAnalysisTime = computed((): string => {
  if (!project.value?.lastAnalysisTime) return '暂无'
  return formatRelativeTime(project.value.lastAnalysisTime)
})

/*** 加载项目详情 - Load project detail ***/
const loadProjectDetail = async (): Promise<void> => {
  if (!projectId.value) {
    loadError.value = true
    errorMessage.value = '项目ID无效'
    return
  }

  loading.value = true
  loadError.value = false

  try {
    const detail = await projectStore.getDetail(projectId.value)
    if (detail) {
      project.value = detail
    } else {
      loadError.value = true
      errorMessage.value = ErrorMessages[40001] || '项目不存在或已被删除'
    }
  } catch (error: unknown) {
    console.error('加载项目详情失败:', error)
    loadError.value = true
    
    // 处理不同类型的错误 - Handle different error types
    if (error instanceof Error) {
      errorMessage.value = error.message || '网络异常，请稍后重试'
    } else {
      errorMessage.value = '网络异常，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}

/*** 返回上一页 - Go back to previous page ***/
const goBack = (): void => {
  uni.navigateBack()
}

/*** 跳转到文档分析页 - Navigate to documents page ***/
const goToDocuments = (): void => {
  uni.navigateTo({ url: `/pages/project/documents?id=${projectId.value}` })
}

/*** 跳转到成本报告页 - Navigate to cost page ***/
const goToCost = (): void => {
  uni.navigateTo({ url: `/pages/project/cost?id=${projectId.value}` })
}

/*** 跳转到设计图页 - Navigate to designs page ***/
const goToDesigns = (): void => {
  uni.navigateTo({ url: `/pages/project/designs?id=${projectId.value}` })
}

/*** 页面加载时获取参数 - Get params on page load ***/
onLoad((options) => {
  if (options?.id) {
    projectId.value = parseInt(options.id as string, 10)
  }
})

/*** 组件挂载时加载数据 - Load data on component mounted ***/
onMounted(() => {
  loadProjectDetail()
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

/*** 内容滚动区 - Content scroll area ***/
.content-scroll {
  height: calc(100vh - 180rpx);
  padding: 0 32rpx;
}

/*** 信息卡片 - Info card with glass morphism ***/
.info-card {
  background: #FFFFFF;
  border-radius: 28rpx;
  padding: 32rpx;
  margin-bottom: 32rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.main-card {
  margin-top: -20rpx;
  position: relative;
  z-index: 10;
}

.card-header {
  margin-bottom: 24rpx;
}

.project-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.project-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #1E293B;
  flex: 1;
  line-height: 1.4;
  margin-right: 16rpx;
}

.project-desc {
  font-size: 26rpx;
  color: #64748B;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/*** 状态徽章 - Status badge ***/
.status-badge {
  display: flex;
  align-items: center;
  padding: 10rpx 20rpx;
  border-radius: 24rpx;
  flex-shrink: 0;
}

.status-draft {
  background: rgba(148, 163, 184, 0.15);
}

.status-in_progress {
  background: rgba(59, 130, 246, 0.12);
}

.status-completed {
  background: rgba(16, 185, 129, 0.12);
}

.status-archived {
  background: rgba(239, 68, 68, 0.12);
}

.status-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  margin-right: 10rpx;
}

.status-draft .status-dot {
  background: #64748B;
}

.status-in_progress .status-dot {
  background: #3B82F6;
}

.status-completed .status-dot {
  background: #10B981;
}

.status-archived .status-dot {
  background: #EF4444;
}

.status-text {
  font-size: 24rpx;
  font-weight: 500;
}

.status-draft .status-text {
  color: #64748B;
}

.status-in_progress .status-text {
  color: #3B82F6;
}

.status-completed .status-text {
  color: #10B981;
}

.status-archived .status-text {
  color: #EF4444;
}

/*** 数据网格 - Data grid with 3D icons ***/
.data-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
}

.data-cell {
  display: flex;
  align-items: center;
  padding: 20rpx;
  background: #F8FAFC;
  border-radius: 20rpx;
}

.data-icon-wrapper {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16rpx;
  margin-right: 16rpx;
}

.area-icon {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.budget-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.style-icon {
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  color: #F59E0B;
}

.time-icon {
  background: linear-gradient(135deg, #E0E7FF 0%, #C7D2FE 100%);
  color: #6366F1;
}

.data-info {
  display: flex;
  flex-direction: column;
}

.data-label {
  font-size: 22rpx;
  color: #94A3B8;
  margin-bottom: 4rpx;
}

.data-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
}

.data-unit {
  font-size: 22rpx;
  font-weight: 400;
  color: #64748B;
}

.budget-value {
  color: #10B981;
}

/*** 统计区域 - Statistics section ***/
.stats-section {
  margin-bottom: 32rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 20rpx;
  display: block;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
}

.stat-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.stat-icon-wrapper {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20rpx;
  margin-bottom: 16rpx;
}

.doc-icon {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.design-icon {
  background: linear-gradient(135deg, #FCE7F3 0%, #FBCFE8 100%);
  color: #EC4899;
}

.analysis-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.stat-value {
  font-size: 40rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 8rpx;
}

.stat-value-small {
  font-size: 26rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 8rpx;
  text-align: center;
}

.stat-label {
  font-size: 22rpx;
  color: #94A3B8;
}

.stat-shadow {
  position: absolute;
  bottom: -6rpx;
  left: 16rpx;
  right: 16rpx;
  height: 12rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 操作区域 - Actions section ***/
.actions-section {
  margin-bottom: 32rpx;
}

.action-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

.action-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.action-icon-wrapper {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20rpx;
  margin-right: 20rpx;
}

.doc-action-icon {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.cost-action-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.design-action-icon {
  background: linear-gradient(135deg, #FCE7F3 0%, #FBCFE8 100%);
  color: #EC4899;
}

.action-info {
  display: flex;
  flex-direction: column;
}

.action-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 6rpx;
}

.action-desc {
  font-size: 24rpx;
  color: #94A3B8;
}

.action-arrow {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}
</style>
