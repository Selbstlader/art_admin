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
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" x2="12" y1="8" y2="12"></line>
          <line x1="12" x2="12.01" y1="16" y2="16"></line>
        </svg>
      </view>
      <text class="error-title">加载失败</text>
      <text class="error-desc">{{ errorMessage }}</text>
      <view class="retry-btn" @click="loadProjectDetail">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path
                  d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z">
                </path>
              </svg>
            </view>
            <view class="data-info">
              <text class="data-label">风格</text>
              <text class="data-value">{{ project.style || '未设置' }}</text>
            </view>
          </view>

          <view class="data-cell">
            <view class="data-icon-wrapper time-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
        <!-- <text class="section-title">项目统计</text> -->
        <view class="stats-grid">
          <!-- 文档数量 - Document count -->
          <view class="stat-card">
            <view class="stat-icon-wrapper doc-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>

        <!-- 查看成本预算 - View cost estimate -->
        <view class="action-card" @click="goToCost">
          <view class="action-left">
            <view class="action-icon-wrapper cost-action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>

        <!-- 查看设计图 - View designs -->
        <view class="action-card" @click="goToDesigns">
          <view class="action-left">
            <view class="action-icon-wrapper design-action-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useProjectStore } from '@/store/project'
import { useNetworkStore } from '@/store/network'
import type { ProjectDetail } from '@/types/api'
import { ProjectStatusText, ProjectStatus, ErrorMessages } from '@/types/common'
import { formatDate, formatRelativeTime } from '@/utils/format'

/*** Store instance - 状态管理实例 ***/
const projectStore = useProjectStore()
const networkStore = useNetworkStore()

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

/*** 处理网络恢复事件 - Handle network recovery event ***/
const handleNetworkRecovered = (): void => {
  /*** Auto reload data when network recovers ***/
  if (loadError.value) {
    loadProjectDetail()
  }
}

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
    if (networkStore.isOffline) {
      errorMessage.value = '网络不可用，请检查网络连接'
    } else if (error instanceof Error) {
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
  /*** Listen to network recovery event ***/
  uni.$on('network:recovered', handleNetworkRecovered)

  loadProjectDetail()
})

/*** 组件卸载时 - On component unmounted ***/
onUnmounted(() => {
  /*** Remove network recovery event listener ***/
  uni.$off('network:recovered', handleNetworkRecovered)
})
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container with pure white background ***/
.page-container {
  min-height: calc(100vh - 140rpx);
  background: #FFFFFF;
  position: relative;
  overflow: hidden;

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
  background: linear-gradient(135deg, #4F8EF7 0%, #3B7BF6 30%, #2563EB 70%, #1D4ED8 100%);
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
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.1);

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
  text-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.15);
}

/*** 加载状态 - Loading state with glass effect ***/
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
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
  box-shadow: 0 0 24rpx rgba(59, 130, 246, 0.2);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 28rpx;
  color: #64748B;
  font-weight: 500;
}

/*** 错误状态 - Error state with glass card ***/
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 140rpx 48rpx;
  margin: 0 32rpx;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border-radius: 32rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.8);
  box-shadow:
    0 8rpx 32rpx rgba(0, 0, 0, 0.04),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

.error-icon {
  width: 140rpx;
  height: 140rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(254, 226, 226, 0.8) 0%, rgba(254, 202, 202, 0.6) 100%);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  margin-bottom: 28rpx;
  color: #EF4444;
  box-shadow: 0 8rpx 24rpx rgba(239, 68, 68, 0.15);
}

.error-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.error-desc {
  font-size: 26rpx;
  color: #64748B;
  text-align: center;
  margin-bottom: 40rpx;
  line-height: 1.6;
}

.retry-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx 56rpx;
  background: linear-gradient(135deg, #4F8EF7 0%, #3B82F6 50%, #2563EB 100%);
  border-radius: 48rpx;
  box-shadow:
    0 8rpx 24rpx rgba(59, 130, 246, 0.35),
    0 2rpx 8rpx rgba(59, 130, 246, 0.2),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.3);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  color: #FFFFFF;

  &:active {
    transform: scale(0.95) translateY(2rpx);
    box-shadow: 0 4rpx 16rpx rgba(59, 130, 246, 0.3);
  }
}

.retry-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #FFFFFF;
  margin-left: 12rpx;
}

/*** 内容滚动区 - Content scroll area ***/
.content-scroll {
  height: calc(100vh - 140rpx);
  padding: 0 32rpx;
  box-sizing: border-box;
}

/*** 信息卡片 - Info card with glass morphism effect ***/
.info-card {
  // background: rgba(255, 255, 255, 0.75);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);

  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-radius: 32rpx;
  padding: 36rpx;
  margin-bottom: 28rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 12rpx 40rpx rgba(59, 130, 246, 0.08),
    0 4rpx 12rpx rgba(0, 0, 0, 0.03),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
  position: relative;
  overflow: hidden;

  /*** 卡片内部光泽 - Inner card glossy effect ***/
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 50%;
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.4) 0%, transparent 100%);
    border-radius: 32rpx 32rpx 0 0;
    pointer-events: none;
  }
}

.main-card {
  // margin-top: -24rpx;
  position: relative;
  z-index: 10;
}

.card-header {
  margin-bottom: 28rpx;
  position: relative;
  z-index: 1;
}

.project-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.project-name {
  font-size: 40rpx;
  font-weight: 800;
  color: #1E293B;
  flex: 1;
  line-height: 1.35;
  margin-right: 16rpx;
  letter-spacing: 1rpx;
}

.project-desc {
  font-size: 26rpx;
  color: #64748B;
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/*** 状态徽章 - Status badge with glass effect ***/
.status-badge {
  display: flex;
  align-items: center;
  padding: 12rpx 24rpx;
  border-radius: 28rpx;
  flex-shrink: 0;
  backdrop-filter: blur(10px);
  border: 1rpx solid rgba(255, 255, 255, 0.5);
}

.status-draft {
  background: rgba(148, 163, 184, 0.12);
  border-color: rgba(148, 163, 184, 0.3);
}

.status-in_progress {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.25);
}

.status-completed {
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.25);
}

.status-archived {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.25);
}

.status-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  margin-right: 10rpx;
  box-shadow: 0 0 8rpx currentColor;
}

.status-draft .status-dot {
  background: #64748B;
}

.status-in_progress .status-dot {
  background: #3B82F6;
  animation: pulse 2s ease-in-out infinite;
}

.status-completed .status-dot {
  background: #10B981;
}

.status-archived .status-dot {
  background: #EF4444;
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
    box-shadow: 0 0 8rpx #3B82F6;
  }

  50% {
    opacity: 0.6;
    box-shadow: 0 0 16rpx #3B82F6;
  }
}

.status-text {
  font-size: 24rpx;
  font-weight: 600;
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

/*** 数据网格 - Data grid with 3D glass cards ***/
.data-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  position: relative;
  z-index: 1;
}

.data-cell {
  display: flex;
  align-items: center;
  padding: 24rpx;
  background: rgba(248, 250, 252, 0.8);
  backdrop-filter: blur(12px);
  border-radius: 24rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 4rpx 16rpx rgba(0, 0, 0, 0.03),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);

  &:active {
    transform: scale(0.98);
  }
}

.data-icon-wrapper {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20rpx;
  margin-right: 16rpx;
  position: relative;

  /*** 3D立体效果 - 3D effect ***/
  &::after {
    content: '';
    position: absolute;
    bottom: -4rpx;
    left: 8rpx;
    right: 8rpx;
    height: 8rpx;
    background: inherit;
    filter: blur(8rpx);
    opacity: 0.4;
    border-radius: 50%;
  }
}

.area-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(59, 130, 246, 0.35);
}

.budget-icon {
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(16, 185, 129, 0.35);
}

.style-icon {
  background: linear-gradient(145deg, #FBBF24 0%, #F59E0B 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(245, 158, 11, 0.35);
}

.time-icon {
  background: linear-gradient(145deg, #818CF8 0%, #6366F1 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(99, 102, 241, 0.35);
}

.data-info {
  display: flex;
  flex-direction: column;
}

.data-label {
  font-size: 22rpx;
  color: #94A3B8;
  margin-bottom: 6rpx;
  font-weight: 500;
}

.data-value {
  font-size: 30rpx;
  font-weight: 700;
  color: #1E293B;
}

.data-unit {
  font-size: 22rpx;
  font-weight: 500;
  color: #64748B;
}

.budget-value {
  color: #10B981;
}

/*** 统计区域 - Statistics section with 3D floating cards ***/
.stats-section {
  margin-bottom: 28rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 24rpx;
  display: block;
  letter-spacing: 1rpx;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
}

.stat-card {
  position: relative;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 28rpx;
  padding: 28rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 8rpx 32rpx rgba(0, 0, 0, 0.04),
    0 2rpx 8rpx rgba(0, 0, 0, 0.02),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;

  /*** 卡片顶部光泽 - Top glossy effect ***/
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 45%;
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.5) 0%, transparent 100%);
    border-radius: 28rpx 28rpx 0 0;
    pointer-events: none;
  }

  &:active {
    transform: scale(0.96) translateY(4rpx);
  }
}

.stat-icon-wrapper {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22rpx;
  margin-bottom: 16rpx;
  position: relative;
  z-index: 1;

  /*** 图标3D阴影 - Icon 3D shadow ***/
  &::after {
    content: '';
    position: absolute;
    bottom: -6rpx;
    left: 12rpx;
    right: 12rpx;
    height: 10rpx;
    background: inherit;
    filter: blur(10rpx);
    opacity: 0.5;
    border-radius: 50%;
    z-index: -1;
  }
}

.doc-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(59, 130, 246, 0.4);
}

.design-icon {
  background: linear-gradient(145deg, #F472B6 0%, #EC4899 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(236, 72, 153, 0.4);
}

.analysis-icon {
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(16, 185, 129, 0.4);
}

.stat-value {
  font-size: 44rpx;
  font-weight: 800;
  color: #1E293B;
  margin-bottom: 8rpx;
  position: relative;
  z-index: 1;
  text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.05);
}

.stat-value-small {
  font-size: 26rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 8rpx;
  text-align: center;
  position: relative;
  z-index: 1;
}

.stat-label {
  font-size: 22rpx;
  color: #94A3B8;
  font-weight: 500;
  position: relative;
  z-index: 1;
}

.stat-shadow {
  position: absolute;
  bottom: -8rpx;
  left: 20rpx;
  right: 20rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(59, 130, 246, 0.12) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 操作区域 - Actions section with glass cards ***/
.actions-section {
  margin-bottom: 28rpx;
}

.action-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 28rpx;
  padding: 28rpx 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 8rpx 32rpx rgba(0, 0, 0, 0.04),
    0 2rpx 8rpx rgba(0, 0, 0, 0.02),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;

  /*** 卡片光泽效果 - Card glossy effect ***/
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

  &:active {
    transform: scale(0.98) translateY(2rpx);
    box-shadow:
      0 4rpx 16rpx rgba(0, 0, 0, 0.03),
      inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
  }
}

.action-left {
  display: flex;
  align-items: center;
  flex: 1;
  position: relative;
  z-index: 1;
}

.action-icon-wrapper {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 24rpx;
  margin-right: 20rpx;
  position: relative;

  /*** 图标3D阴影 - Icon 3D shadow ***/
  &::after {
    content: '';
    position: absolute;
    bottom: -6rpx;
    left: 12rpx;
    right: 12rpx;
    height: 10rpx;
    background: inherit;
    filter: blur(10rpx);
    opacity: 0.45;
    border-radius: 50%;
    z-index: -1;
  }
}

.doc-action-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(59, 130, 246, 0.4);
}

.cost-action-icon {
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(16, 185, 129, 0.4);
}

.design-action-icon {
  background: linear-gradient(145deg, #F472B6 0%, #EC4899 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(236, 72, 153, 0.4);
}

.action-info {
  display: flex;
  flex-direction: column;
}

.action-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 8rpx;
}

.action-desc {
  font-size: 24rpx;
  color: #94A3B8;
  font-weight: 500;
}

.action-arrow {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 16rpx;
  position: relative;
  z-index: 1;
  transition: all 0.2s ease;
}

.action-card:active .action-arrow {
  color: #3B82F6;
  background: rgba(59, 130, 246, 0.1);
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}
</style>
