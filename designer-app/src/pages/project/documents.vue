<template>
  <!-- 
    文档分析列表页面 - C4D风格高级UI
    Document analysis list page - C4D style premium UI with 3D floating cards
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
        <text class="header-title">文档分析</text>
      </view>
    </view>

    <!-- 加载状态 - Loading state -->
    <view v-if="loading && !documents.length" class="loading-container">
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
      <view class="retry-btn" @click="loadDocuments">
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
    <view v-else-if="!documents.length" class="empty-container">
      <view class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path>
          <path d="M14 2v4a2 2 0 0 0 2 2h4"></path>
        </svg>
      </view>
      <text class="empty-title">暂无文档</text>
      <text class="empty-desc">该项目还没有上传任何文档</text>
    </view>

    <!-- 文档列表 - Document list -->
    <scroll-view 
      v-else 
      class="content-scroll" 
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
    >
      <!-- 统计卡片 - Statistics card -->
      <view class="stats-card">
        <view class="stats-item">
          <text class="stats-value">{{ documents.length }}</text>
          <text class="stats-label">全部文档</text>
        </view>
        <view class="stats-divider"></view>
        <view class="stats-item">
          <text class="stats-value completed-value">{{ completedCount }}</text>
          <text class="stats-label">已分析</text>
        </view>
        <view class="stats-divider"></view>
        <view class="stats-item">
          <text class="stats-value pending-value">{{ pendingCount }}</text>
          <text class="stats-label">待分析</text>
        </view>
      </view>

      <!-- 文档列表 - Document list items -->
      <view class="document-list">
        <view 
          v-for="doc in documents" 
          :key="doc.id" 
          class="document-card"
          @click="handleDocumentClick(doc)"
        >
          <!-- 文档图标 - Document icon -->
          <view :class="['doc-icon-wrapper', getIconClass(doc.fileType)]">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path>
              <path d="M14 2v4a2 2 0 0 0 2 2h4"></path>
              <path d="M10 9H8"></path>
              <path d="M16 13H8"></path>
              <path d="M16 17H8"></path>
            </svg>
          </view>

          <!-- 文档信息 - Document info -->
          <view class="doc-info">
            <text class="doc-name">{{ doc.fileName }}</text>
            <view class="doc-meta">
              <text class="doc-size">{{ formatFileSize(doc.fileSize) }}</text>
              <text class="doc-time">{{ formatDate(doc.createdAt, 'YYYY-MM-DD HH:mm') }}</text>
            </view>
          </view>

          <!-- 分析状态 - Analysis status -->
          <view :class="['status-badge', `status-${doc.analysisStatus}`]">
            <view v-if="doc.analysisStatus === 'analyzing'" class="status-spinner"></view>
            <view v-else class="status-dot"></view>
            <text class="status-text">{{ getStatusText(doc.analysisStatus) }}</text>
          </view>

          <!-- 箭头指示 - Arrow indicator (only for completed) -->
          <view v-if="doc.analysisStatus === 'completed'" class="doc-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>
      </view>

      <!-- 底部安全区 - Bottom safe area -->
      <view class="safe-area-bottom"></view>
    </scroll-view>

    <!-- 文档详情弹窗 - Document detail modal -->
    <view v-if="showDetailModal" class="modal-overlay" @click="closeDetailModal">
      <view class="modal-container" @click.stop>
        <!-- 弹窗头部 - Modal header -->
        <view class="modal-header">
          <text class="modal-title">文档分析结果</text>
          <view class="modal-close" @click="closeDetailModal">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </view>
        </view>

        <!-- 弹窗内容 - Modal content -->
        <scroll-view class="modal-content" scroll-y>
          <!-- 加载状态 - Loading state -->
          <view v-if="detailLoading" class="detail-loading">
            <view class="loading-spinner small"></view>
            <text class="loading-text">加载分析结果...</text>
          </view>

          <!-- 分析结果 - Analysis result -->
          <view v-else-if="currentSummary" class="analysis-result">
            <!-- 文档信息 - Document info -->
            <view class="result-doc-info">
              <text class="result-doc-name">{{ selectedDocument?.fileName }}</text>
            </view>

            <!-- 关键字标签 - Keywords tags -->
            <view v-if="currentKeywords.length" class="keywords-section">
              <text class="section-label">关键字提取</text>
              <view class="keywords-list">
                <view v-for="(keyword, index) in currentKeywords" :key="index" class="keyword-tag">
                  {{ keyword }}
                </view>
              </view>
            </view>

            <!-- 项目概述 - Project overview -->
            <view class="summary-section">
              <view class="section-header">
                <view class="section-icon overview-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <path d="M12 16v-4"></path>
                    <path d="M12 8h.01"></path>
                  </svg>
                </view>
                <text class="section-title">项目概述</text>
              </view>
              <text class="section-content">{{ currentSummary.overview || '暂无概述信息' }}</text>
            </view>

            <!-- 核心需求 - Core requirements -->
            <view class="summary-section">
              <view class="section-header">
                <view class="section-icon requirements-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"></path>
                    <path d="m9 12 2 2 4-4"></path>
                  </svg>
                </view>
                <text class="section-title">核心需求</text>
              </view>
              <text class="section-content">{{ currentSummary.requirements || '暂无需求信息' }}</text>
            </view>

            <!-- 特殊要求 - Special requirements -->
            <view class="summary-section">
              <view class="section-header">
                <view class="section-icon special-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"></path>
                  </svg>
                </view>
                <text class="section-title">特殊要求</text>
              </view>
              <text class="section-content">{{ currentSummary.special || '暂无特殊要求' }}</text>
            </view>
          </view>

          <!-- 加载失败 - Load failed -->
          <view v-else class="detail-error">
            <text class="detail-error-text">加载分析结果失败</text>
            <view class="detail-retry-btn" @click="loadDocumentDetail">
              <text>重试</text>
            </view>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * Document analysis list page - C4D style premium UI
 * 文档分析列表页面 - C4D风格高级UI，3D立体数据图表
 * Requirements: 5.1, 5.2, 5.3, 5.4
 ***/
import { ref, computed, onMounted } from 'vue'
import { onLoad, onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { documentApi } from '@/api/document'
import type { Document, DocumentSummary } from '@/types/api'
import { formatDate, formatFileSize } from '@/utils/format'
import { ErrorMessages, DocumentAnalysisStatus } from '@/types/common'

/*** 页面状态 - Page state ***/
const projectId = ref<number>(0)
const documents = ref<Document[]>([])
const loading = ref(true)
const loadError = ref(false)
const errorMessage = ref('网络异常，请稍后重试')
const refreshing = ref(false)

/*** 详情弹窗状态 - Detail modal state ***/
const showDetailModal = ref(false)
const selectedDocument = ref<Document | null>(null)
const currentSummary = ref<DocumentSummary | null>(null)
const currentKeywords = ref<string[]>([])
const detailLoading = ref(false)

/*** 计算属性 - Computed properties ***/
/*** 已分析文档数量 - Completed document count ***/
const completedCount = computed((): number => {
  return documents.value.filter(doc => doc.analysisStatus === 'completed').length
})

/*** 待分析文档数量 - Pending document count ***/
const pendingCount = computed((): number => {
  return documents.value.filter(doc => 
    doc.analysisStatus === 'pending' || doc.analysisStatus === 'analyzing'
  ).length
})

/*** 获取状态文本 - Get status text ***/
const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    [DocumentAnalysisStatus.PENDING]: '待分析',
    [DocumentAnalysisStatus.ANALYZING]: '分析中',
    [DocumentAnalysisStatus.COMPLETED]: '已分析',
    [DocumentAnalysisStatus.FAILED]: '分析失败'
  }
  return statusMap[status] || status
}

/*** 获取文件图标样式 - Get file icon class ***/
const getIconClass = (fileType: string): string => {
  const typeMap: Record<string, string> = {
    'pdf': 'pdf-icon',
    'doc': 'word-icon',
    'docx': 'word-icon',
    'xls': 'excel-icon',
    'xlsx': 'excel-icon',
    'ppt': 'ppt-icon',
    'pptx': 'ppt-icon',
    'txt': 'txt-icon'
  }
  return typeMap[fileType?.toLowerCase()] || 'default-icon'
}

/*** 加载文档列表 - Load document list ***/
const loadDocuments = async (): Promise<void> => {
  if (!projectId.value) {
    loadError.value = true
    errorMessage.value = '项目ID无效'
    return
  }

  loading.value = true
  loadError.value = false

  try {
    const res = await documentApi.getList(projectId.value, { current: 1, size: 100 })
    if (res.code === 200 && res.data) {
      documents.value = res.data.records || []
    } else {
      loadError.value = true
      errorMessage.value = res.msg || res.message || '加载失败'
    }
  } catch (error: unknown) {
    console.error('加载文档列表失败:', error)
    loadError.value = true
    if (error instanceof Error) {
      errorMessage.value = error.message || '网络异常，请稍后重试'
    } else {
      errorMessage.value = '网络异常，请稍后重试'
    }
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

/*** 下拉刷新 - Pull down refresh ***/
const onRefresh = async (): Promise<void> => {
  refreshing.value = true
  await loadDocuments()
}

/*** 处理文档点击 - Handle document click ***/
const handleDocumentClick = (doc: Document): void => {
  if (doc.analysisStatus === 'completed') {
    // 已分析文档，显示详情弹窗
    selectedDocument.value = doc
    showDetailModal.value = true
    loadDocumentDetail()
  } else if (doc.analysisStatus === 'pending' || doc.analysisStatus === 'failed') {
    // 未分析或分析失败的文档，显示提示
    uni.showToast({
      title: ErrorMessages[40002] || '该文档尚未分析，请在 Web 端进行分析',
      icon: 'none',
      duration: 2500
    })
  } else if (doc.analysisStatus === 'analyzing') {
    // 分析中的文档
    uni.showToast({
      title: '文档正在分析中，请稍后查看',
      icon: 'none',
      duration: 2000
    })
  }
}

/*** 加载文档详情 - Load document detail ***/
const loadDocumentDetail = async (): Promise<void> => {
  if (!selectedDocument.value) return

  detailLoading.value = true
  currentSummary.value = null
  currentKeywords.value = []

  try {
    // 并行请求关键字和摘要
    const [keywordsRes, summaryRes] = await Promise.all([
      documentApi.getKeywords(selectedDocument.value.id),
      documentApi.getSummary(selectedDocument.value.id)
    ])

    if (keywordsRes.code === 200 && keywordsRes.data) {
      currentKeywords.value = keywordsRes.data
    }

    if (summaryRes.code === 200 && summaryRes.data) {
      currentSummary.value = summaryRes.data
    }
  } catch (error: unknown) {
    console.error('加载文档详情失败:', error)
    currentSummary.value = null
  } finally {
    detailLoading.value = false
  }
}

/*** 关闭详情弹窗 - Close detail modal ***/
const closeDetailModal = (): void => {
  showDetailModal.value = false
  selectedDocument.value = null
  currentSummary.value = null
  currentKeywords.value = []
}

/*** 返回上一页 - Go back to previous page ***/
const goBack = (): void => {
  uni.navigateBack()
}

/*** 页面加载时获取参数 - Get params on page load ***/
onLoad((options) => {
  if (options?.id) {
    projectId.value = parseInt(options.id as string, 10)
  }
})

/*** 组件挂载时加载数据 - Load data on component mounted ***/
onMounted(() => {
  loadDocuments()
})

/*** 下拉刷新事件 - Pull down refresh event ***/
onPullDownRefresh(async () => {
  await loadDocuments()
  uni.stopPullDownRefresh()
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
  
  &.small {
    width: 48rpx;
    height: 48rpx;
    border-width: 3rpx;
  }
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

/*** 统计卡片 - Statistics card ***/
.stats-card {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: #FFFFFF;
  border-radius: 28rpx;
  padding: 32rpx 24rpx;
  margin-top: -20rpx;
  margin-bottom: 32rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.stats-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.stats-value {
  font-size: 44rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 8rpx;
}

.completed-value {
  color: #10B981;
}

.pending-value {
  color: #F59E0B;
}

.stats-label {
  font-size: 24rpx;
  color: #94A3B8;
}

.stats-divider {
  width: 2rpx;
  height: 60rpx;
  background: #E2E8F0;
}

/*** 文档列表 - Document list ***/
.document-list {
  margin-bottom: 32rpx;
}

.document-card {
  display: flex;
  align-items: center;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

/*** 文档图标 - Document icon ***/
.doc-icon-wrapper {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 18rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.pdf-icon {
  background: linear-gradient(135deg, #FEE2E2 0%, #FECACA 100%);
  color: #EF4444;
}

.word-icon {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.excel-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.ppt-icon {
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  color: #F59E0B;
}

.txt-icon {
  background: linear-gradient(135deg, #E0E7FF 0%, #C7D2FE 100%);
  color: #6366F1;
}

.default-icon {
  background: linear-gradient(135deg, #F1F5F9 0%, #E2E8F0 100%);
  color: #64748B;
}

/*** 文档信息 - Document info ***/
.doc-info {
  flex: 1;
  min-width: 0;
  margin-right: 16rpx;
}

.doc-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 8rpx;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-meta {
  display: flex;
  align-items: center;
}

.doc-size {
  font-size: 22rpx;
  color: #94A3B8;
  margin-right: 16rpx;
}

.doc-time {
  font-size: 22rpx;
  color: #94A3B8;
}

/*** 状态徽章 - Status badge ***/
.status-badge {
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
  margin-right: 8rpx;
}

.status-pending {
  background: rgba(245, 158, 11, 0.12);
}

.status-analyzing {
  background: rgba(59, 130, 246, 0.12);
}

.status-completed {
  background: rgba(16, 185, 129, 0.12);
}

.status-failed {
  background: rgba(239, 68, 68, 0.12);
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  margin-right: 8rpx;
}

.status-spinner {
  width: 12rpx;
  height: 12rpx;
  border: 2rpx solid #BFDBFE;
  border-top-color: #3B82F6;
  border-radius: 50%;
  margin-right: 8rpx;
  animation: spin 0.8s linear infinite;
}

.status-pending .status-dot {
  background: #F59E0B;
}

.status-completed .status-dot {
  background: #10B981;
}

.status-failed .status-dot {
  background: #EF4444;
}

.status-text {
  font-size: 22rpx;
  font-weight: 500;
}

.status-pending .status-text {
  color: #F59E0B;
}

.status-analyzing .status-text {
  color: #3B82F6;
}

.status-completed .status-text {
  color: #10B981;
}

.status-failed .status-text {
  color: #EF4444;
}

/*** 箭头指示 - Arrow indicator ***/
.doc-arrow {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}

/*** 弹窗样式 - Modal styles ***/
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  width: 100%;
  max-height: 80vh;
  background: #FFFFFF;
  border-radius: 32rpx 32rpx 0 0;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx;
  border-bottom: 2rpx solid #F1F5F9;
}

.modal-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
}

.modal-close {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F1F5F9;
  border-radius: 50%;
  color: #64748B;
  
  &:active {
    background: #E2E8F0;
  }
}

.modal-content {
  max-height: calc(80vh - 120rpx);
  padding: 32rpx;
}

/*** 详情加载状态 - Detail loading state ***/
.detail-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
}

/*** 分析结果 - Analysis result ***/
.analysis-result {
  padding-bottom: 40rpx;
}

.result-doc-info {
  margin-bottom: 32rpx;
}

.result-doc-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
}

/*** 关键字区域 - Keywords section ***/
.keywords-section {
  margin-bottom: 32rpx;
}

.section-label {
  font-size: 24rpx;
  color: #64748B;
  margin-bottom: 16rpx;
  display: block;
}

.keywords-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.keyword-tag {
  padding: 10rpx 20rpx;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #3B82F6;
}

/*** 摘要区域 - Summary section ***/
.summary-section {
  background: #F8FAFC;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;
}

.section-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12rpx;
  margin-right: 12rpx;
}

.overview-icon {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  color: #3B82F6;
}

.requirements-icon {
  background: linear-gradient(135deg, #D1FAE5 0%, #A7F3D0 100%);
  color: #10B981;
}

.special-icon {
  background: linear-gradient(135deg, #FEF3C7 0%, #FDE68A 100%);
  color: #F59E0B;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
}

.section-content {
  font-size: 26rpx;
  color: #475569;
  line-height: 1.8;
}

/*** 详情错误状态 - Detail error state ***/
.detail-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
}

.detail-error-text {
  font-size: 28rpx;
  color: #64748B;
  margin-bottom: 24rpx;
}

.detail-retry-btn {
  padding: 16rpx 40rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-radius: 32rpx;
  color: #FFFFFF;
  font-size: 26rpx;
  
  &:active {
    opacity: 0.9;
  }
}
</style>
