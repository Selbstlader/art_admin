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
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6"></path>
          </svg>
        </view>
        <text class="header-title">项目详情</text>
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
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" x2="12" y1="8" y2="12"></line>
          <line x1="12" x2="12.01" y1="16" y2="16"></line>
        </svg>
      </view>
      <text class="error-title">加载失败</text>
      <text class="error-desc">{{ errorMessage }}</text>
      <view class="retry-btn" @click="loadDocuments">
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

    <!-- 空状态 - Empty state -->
    <view v-else-if="!documents.length" class="empty-container">
      <view class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path>
          <path d="M14 2v4a2 2 0 0 0 2 2h4"></path>
        </svg>
      </view>
      <text class="empty-title">暂无文档</text>
      <text class="empty-desc">该项目还没有上传任何文档</text>
    </view>

    <!-- 文档列表 - Document list -->
    <scroll-view v-else class="content-scroll" scroll-y :refresher-enabled="true" :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh">
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
        <view v-for="doc in documents" :key="doc.id" class="document-card" @click="handleDocumentClick(doc)">
          <!-- 文档图标 - Document icon -->
          <view style="display: flex;">
            <view :class="['doc-icon-wrapper', getIconClass(doc.fileType)]">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </view>
          </view>
          <!-- 操作按钮区域 - Action buttons area -->
          <view class="doc-actions">
            <view class="preview-btn" @click.stop="previewDocument(doc)">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
              <text class="action-text">预览</text>
            </view>
          </view>
          <view v-if="doc.errorMessage" class="doc-errrr">{{ doc.errorMessage }}</view>
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
          <text class="modal-title">{{ selectedDocument?.fileName }}</text>
          <view class="modal-close" @click="closeDetailModal">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
          <view v-else-if="documentDetail" class="analysis-result">
            <!-- 提取信息卡片 - Extracted info cards -->
            <view class="info-cards">
              <!-- 项目名称 - Project name -->
              <view v-if="documentDetail.projectName" class="info-card">
                <view class="info-icon project-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                    <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                  </svg>
                </view>
                <view class="info-content">
                  <text class="info-label">项目名称</text>
                  <text class="info-value">{{ documentDetail.projectName }}</text>
                </view>
              </view>

              <!-- 面积 - Area -->
              <view v-if="documentDetail.extractedArea" class="info-card">
                <view class="info-icon area-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect width="18" height="18" x="3" y="3" rx="2"></rect>
                    <path d="M3 9h18"></path>
                    <path d="M9 21V9"></path>
                  </svg>
                </view>
                <view class="info-content">
                  <text class="info-label">面积</text>
                  <text class="info-value">{{ documentDetail.extractedArea }} ㎡</text>
                </view>
              </view>

              <!-- 预算 - Budget -->
              <view v-if="documentDetail.extractedBudget" class="info-card">
                <view class="info-icon budget-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
                  </svg>
                </view>
                <view class="info-content">
                  <text class="info-label">预算</text>
                  <text class="info-value">{{ formatBudget(documentDetail.extractedBudget) }}</text>
                </view>
              </view>

              <!-- 风格 - Style -->
              <view v-if="documentDetail.extractedStyle" class="info-card">
                <view class="info-icon style-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path>
                  </svg>
                </view>
                <view class="info-content">
                  <text class="info-label">设计风格</text>
                  <text class="info-value">{{ documentDetail.extractedStyle }}</text>
                </view>
              </view>
            </view>

            <!-- 功能分区 - Functional zones -->
            <view v-if="documentDetail.functionalZones?.length" class="tags-section">
              <text class="section-label">功能分区</text>
              <view class="tags-list">
                <view v-for="(zone, index) in documentDetail.functionalZones" :key="index" class="zone-tag">
                  {{ zone }}
                </view>
              </view>
            </view>

            <!-- 关键字标签 - Keywords tags -->
            <view v-if="documentDetail.keywords?.length" class="tags-section">
              <text class="section-label">关键字</text>
              <view class="tags-list">
                <view v-for="(keyword, index) in documentDetail.keywords" :key="index" class="keyword-tag">
                  {{ keyword }}
                </view>
              </view>
            </view>

            <!-- 文档摘要 - Document summary -->
            <view v-if="documentDetail.summary" class="summary-section">
              <view class="section-header">
                <view class="section-icon overview-icon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <path d="M14 2v6h6"></path>
                    <path d="M16 13H8"></path>
                    <path d="M16 17H8"></path>
                    <path d="M10 9H8"></path>
                  </svg>
                </view>
                <text class="section-title">文档摘要</text>
              </view>
              <text class="section-content">{{ documentDetail.summary }}</text>
            </view>

            <!-- 无数据提示 - No data hint -->
            <view v-if="!documentDetail.projectName && !documentDetail.extractedArea && !documentDetail.summary"
              class="no-data-hint">
              <text class="no-data-text">暂无分析数据</text>
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
import type { Document } from '@/types/api'
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
const documentDetail = ref<Document | null>(null)
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

/*** 格式化预算金额 - Format budget amount ***/
const formatBudget = (budget: number): string => {
  if (budget >= 10000) {
    return `${(budget / 10000).toFixed(1)}万元`
  }
  return `${budget.toFixed(0)}元`
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
  documentDetail.value = null

  try {
    // 获取文档详情（包含所有分析数据）
    const res = await documentApi.getDetail(selectedDocument.value.id)
    if (res.code === 200 && res.data) {
      documentDetail.value = res.data
    }
  } catch (error: unknown) {
    console.error('加载文档详情失败:', error)
    documentDetail.value = null
  } finally {
    detailLoading.value = false
  }
}

/*** 关闭详情弹窗 - Close detail modal ***/
const closeDetailModal = (): void => {
  showDetailModal.value = false
  selectedDocument.value = null
  documentDetail.value = null
}

/*** 返回上一页 - Go back to previous page ***/
const goBack = (): void => {
  uni.navigateBack()
}

/*** 获取 API 基础 URL - Get API base URL ***/
const getApiBaseUrl = (): string => {
  // @ts-ignore - Vite env types
  const envBaseUrl = typeof import.meta !== 'undefined' && (import.meta as any).env?.VITE_API_BASE_URL
  return envBaseUrl || 'http://localhost:48080'
}

/*** 预览文档 - Preview document ***/
const previewDocument = (doc: Document): void => {
  if (!doc.filePath) {
    uni.showToast({
      title: '文件路径不存在',
      icon: 'none',
      duration: 2000
    })
    return
  }

  // 构建文件预览 URL - Build file preview URL
  const baseUrl = getApiBaseUrl()
  const fileUrl = `${baseUrl}/uploads/${doc.filePath}`
  const token = uni.getStorageSync('token') || ''

  uni.showLoading({ title: '加载文件中...', mask: true })

  // 根据文件类型选择预览方式 - Choose preview method based on file type
  const fileType = doc.fileType?.toLowerCase() || ''
  
  // 支持的文档类型 - Supported document types for openDocument
  const supportedTypes = ['doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'pdf']
  
  if (supportedTypes.includes(fileType)) {
    // 下载文件后使用 openDocument 预览 - Download and preview with openDocument
    uni.downloadFile({
      url: fileUrl,
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (downloadRes) => {
        uni.hideLoading()
        if (downloadRes.statusCode === 200) {
          uni.openDocument({
            filePath: downloadRes.tempFilePath,
            fileType: fileType as any,
            showMenu: true,
            success: () => {
              console.log('文档打开成功')
            },
            fail: (err) => {
              console.error('打开文档失败:', err)
              uni.showToast({
                title: '无法打开此文件类型',
                icon: 'none',
                duration: 2000
              })
            }
          })
        } else {
          uni.showToast({
            title: '文件下载失败',
            icon: 'none',
            duration: 2000
          })
        }
      },
      fail: (err) => {
        uni.hideLoading()
        console.error('下载文件失败:', err)
        uni.showToast({
          title: '文件下载失败，请检查网络',
          icon: 'none',
          duration: 2000
        })
      }
    })
  } else if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(fileType)) {
    // 图片类型使用 previewImage - Image types use previewImage
    uni.hideLoading()
    uni.previewImage({
      urls: [fileUrl],
      current: fileUrl
    })
  } else {
    // 其他类型尝试用 webview 打开或提示不支持 - Other types try webview or show unsupported
    uni.hideLoading()
    // #ifdef H5
    window.open(fileUrl, '_blank')
    // #endif
    // #ifndef H5
    uni.showToast({
      title: '暂不支持预览此文件类型',
      icon: 'none',
      duration: 2000
    })
    // #endif
  }
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
/*** 页面容器 - Page container with pure white background ***/
.page-container {
  min-height: 100vh;
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

  &.small {
    width: 48rpx;
    height: 48rpx;
    border-width: 3rpx;
    margin-bottom: 16rpx;
    box-shadow: none;
  }
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
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
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

/*** 空状态 - Empty state with glass card ***/
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 140rpx 48rpx;
  margin: 0 32rpx;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  border-radius: 32rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.8);
  box-shadow:
    0 8rpx 32rpx rgba(0, 0, 0, 0.04),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

.empty-icon {
  width: 140rpx;
  height: 140rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(241, 245, 249, 0.9) 0%, rgba(226, 232, 240, 0.7) 100%);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  margin-bottom: 28rpx;
  color: #94A3B8;
  box-shadow: 0 8rpx 24rpx rgba(148, 163, 184, 0.15);
}

.empty-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: #64748B;
  text-align: center;
  line-height: 1.6;
}

/*** 内容滚动区 - Content scroll area ***/
.content-scroll {
  height: calc(100vh - 140rpx);
  padding: 0 32rpx;
  box-sizing: border-box;
}

/*** 统计卡片 - Statistics card with glass morphism ***/
.stats-card {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-radius: 32rpx;
  padding: 36rpx 24rpx;
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

.stats-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  position: relative;
  z-index: 1;
}

.stats-value {
  font-size: 48rpx;
  font-weight: 800;
  color: #1E293B;
  margin-bottom: 8rpx;
  text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.05);
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
  font-weight: 500;
}

.stats-divider {
  width: 2rpx;
  height: 64rpx;
  background: linear-gradient(180deg, transparent 0%, rgba(226, 232, 240, 0.8) 50%, transparent 100%);
}

/*** 文档列表 - Document list ***/
.document-list {
  margin-bottom: 32rpx;
}

.document-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 28rpx;
  padding: 24rpx;
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

  .doc-errrr {
    margin-top: 20rpx;
    margin-left: 20rpx;
    font-size: 0.6875rem;
    color: #94A3B8;
    font-weight: 500;
  }
}

/*** 文档图标 - Document icon with 3D effect ***/
.doc-icon-wrapper {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
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
    opacity: 0.45;
    border-radius: 50%;
    z-index: -1;
  }
}

.pdf-icon {
  background: linear-gradient(145deg, #F87171 0%, #EF4444 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(239, 68, 68, 0.4);
}

.word-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(59, 130, 246, 0.4);
}

.excel-icon {
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(16, 185, 129, 0.4);
}

.ppt-icon {
  background: linear-gradient(145deg, #FBBF24 0%, #F59E0B 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(245, 158, 11, 0.4);
}

.txt-icon {
  background: linear-gradient(145deg, #818CF8 0%, #6366F1 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(99, 102, 241, 0.4);
}

.default-icon {
  background: linear-gradient(145deg, #94A3B8 0%, #64748B 100%);
  color: #FFFFFF;
  box-shadow: 0 6rpx 20rpx rgba(100, 116, 139, 0.4);
}

/*** 文档信息 - Document info ***/
.doc-info {
  flex: 1;
  min-width: 0;
  margin-right: 16rpx;
  position: relative;
  z-index: 1;
}

.doc-name {
  font-size: 30rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 10rpx;
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
  font-weight: 500;
}

.doc-time {
  font-size: 22rpx;
  color: #94A3B8;
  font-weight: 500;
}

/*** 错误信息 - Error message for failed documents ***/
.doc-error {
  font-size: 22rpx;
  color: #EF4444;
  margin-top: 8rpx;
  display: block;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/*** 状态徽章 - Status badge with glass effect ***/
.status-badge {
  display: flex;
  align-items: center;
  padding: 10rpx 18rpx;
  border-radius: 24rpx;
  flex-shrink: 0;
  margin-right: 8rpx;
  backdrop-filter: blur(10px);
  border: 1rpx solid rgba(255, 255, 255, 0.5);
  position: relative;
  z-index: 1;
}

.status-pending {
  background: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.25);
}

.status-analyzing {
  background: rgba(59, 130, 246, 0.12);
  border-color: rgba(59, 130, 246, 0.25);
}

.status-completed {
  background: rgba(16, 185, 129, 0.12);
  border-color: rgba(16, 185, 129, 0.25);
}

.status-failed {
  background: rgba(239, 68, 68, 0.12);
  border-color: rgba(239, 68, 68, 0.25);
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  margin-right: 8rpx;
  box-shadow: 0 0 8rpx currentColor;
}

.status-spinner {
  width: 12rpx;
  height: 12rpx;
  border: 2rpx solid rgba(59, 130, 246, 0.3);
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
  font-weight: 600;
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
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 14rpx;
  position: relative;
  z-index: 1;
  transition: all 0.2s ease;
}

.document-card:active .doc-arrow {
  color: #3B82F6;
  background: rgba(59, 130, 246, 0.1);
}

/*** 操作按钮区域 - Action buttons area ***/
.doc-actions {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid rgba(226, 232, 240, 0.5);
  width: 100%;
}

.preview-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12rpx 24rpx;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.15) 100%);
  border-radius: 16rpx;
  border: 1rpx solid rgba(59, 130, 246, 0.2);
  color: #3B82F6;
  transition: all 0.2s ease;

  &:active {
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(37, 99, 235, 0.25) 100%);
    transform: scale(0.96);
  }
}

.action-text {
  font-size: 24rpx;
  font-weight: 600;
  margin-left: 8rpx;
}

/*** 底部安全区 - Bottom safe area ***/
.safe-area-bottom {
  height: 120rpx;
}

/*** 弹窗样式 - Modal styles with glass morphism ***/
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  width: 100%;
  max-height: 80vh;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(248, 250, 252, 0.95) 100%);
  backdrop-filter: blur(24px);
  border-radius: 40rpx 40rpx 0 0;
  overflow: hidden;
  box-shadow: 0 -8rpx 40rpx rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 36rpx 32rpx;
  border-bottom: 1rpx solid rgba(226, 232, 240, 0.6);
  background: rgba(255, 255, 255, 0.8);
}

.modal-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #1E293B;
  letter-spacing: 1rpx;
}

.modal-close {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(241, 245, 249, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  color: #64748B;
  transition: all 0.2s ease;

  &:active {
    background: rgba(226, 232, 240, 0.9);
    transform: scale(0.92);
  }
}

.modal-content {
  max-height: calc(80vh - 130rpx);
  padding: 32rpx;
  box-sizing: border-box;
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

/*** 信息卡片网格 - Info cards grid ***/
.info-cards {
  display: flex;
  // grid-template-columns: repeat(2, 1fr);
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 28rpx;
}

.info-card {
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, rgba(248, 250, 252, 0.9) 0%, rgba(241, 245, 249, 0.8) 100%);
  backdrop-filter: blur(12px);
  border-radius: 20rpx;
  padding: 20rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.03);
  flex: 1 1 auto;
  width: 35%;
}

.info-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12rpx;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.project-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 4rpx 12rpx rgba(59, 130, 246, 0.35);
}

.area-icon {
  background: linear-gradient(145deg, #34D399 0%, #10B981 100%);
  color: #FFFFFF;
  box-shadow: 0 4rpx 12rpx rgba(16, 185, 129, 0.35);
}

.budget-icon {
  background: linear-gradient(145deg, #FBBF24 0%, #F59E0B 100%);
  color: #FFFFFF;
  box-shadow: 0 4rpx 12rpx rgba(245, 158, 11, 0.35);
}

.style-icon {
  background: linear-gradient(145deg, #A78BFA 0%, #8B5CF6 100%);
  color: #FFFFFF;
  box-shadow: 0 4rpx 12rpx rgba(139, 92, 246, 0.35);
}

.info-content {
  flex: 1;
  min-width: 0;
}

.info-label {
  font-size: 22rpx;
  color: #94A3B8;
  display: block;
  margin-bottom: 4rpx;
}

.info-value {
  font-size: 26rpx;
  font-weight: 600;
  color: #1E293B;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/*** 标签区域 - Tags section ***/
.tags-section {
  margin-bottom: 24rpx;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.zone-tag {
  padding: 10rpx 20rpx;
  background: linear-gradient(135deg, rgba(236, 253, 245, 0.9) 0%, rgba(209, 250, 229, 0.8) 100%);
  backdrop-filter: blur(10px);
  border-radius: 20rpx;
  font-size: 24rpx;
  font-weight: 600;
  color: #10B981;
  border: 1rpx solid rgba(16, 185, 129, 0.2);
}

/*** 无数据提示 - No data hint ***/
.no-data-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60rpx 0;
}

.no-data-text {
  font-size: 28rpx;
  color: #94A3B8;
}

.result-doc-info {
  margin-bottom: 32rpx;
}

.result-doc-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #1E293B;
}

/*** 关键字区域 - Keywords section with glass tags ***/
.keywords-section {
  margin-bottom: 32rpx;
}

.section-label {
  font-size: 24rpx;
  color: #64748B;
  margin-bottom: 16rpx;
  display: block;
  font-weight: 600;
}

.keywords-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.keyword-tag {
  padding: 12rpx 24rpx;
  background: linear-gradient(135deg, rgba(239, 246, 255, 0.9) 0%, rgba(219, 234, 254, 0.8) 100%);
  backdrop-filter: blur(10px);
  border-radius: 24rpx;
  font-size: 24rpx;
  font-weight: 600;
  color: #3B82F6;
  border: 1rpx solid rgba(59, 130, 246, 0.2);
  box-shadow: 0 2rpx 8rpx rgba(59, 130, 246, 0.1);
}

/*** 摘要区域 - Summary section with glass card ***/
.summary-section {
  background: linear-gradient(135deg, rgba(248, 250, 252, 0.9) 0%, rgba(241, 245, 249, 0.8) 100%);
  backdrop-filter: blur(12px);
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.9);
  box-shadow:
    0 4rpx 16rpx rgba(0, 0, 0, 0.03),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;
}

.section-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 14rpx;
  margin-right: 12rpx;
  position: relative;

  /*** 图标3D阴影 - Icon 3D shadow ***/
  &::after {
    content: '';
    position: absolute;
    bottom: -4rpx;
    left: 8rpx;
    right: 8rpx;
    height: 6rpx;
    background: inherit;
    filter: blur(6rpx);
    opacity: 0.4;
    border-radius: 50%;
    z-index: -1;
  }
}

.overview-icon {
  background: linear-gradient(145deg, #60A5FA 0%, #3B82F6 100%);
  color: #FFFFFF;
  box-shadow: 0 4rpx 12rpx rgba(59, 130, 246, 0.35);
}

.section-title {
  font-size: 28rpx;
  font-weight: 700;
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
  font-weight: 500;
}

.detail-retry-btn {
  padding: 20rpx 48rpx;
  background: linear-gradient(135deg, #4F8EF7 0%, #3B82F6 50%, #2563EB 100%);
  border-radius: 36rpx;
  color: #FFFFFF;
  font-size: 26rpx;
  font-weight: 600;
  box-shadow:
    0 6rpx 20rpx rgba(59, 130, 246, 0.35),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.3);
  transition: all 0.2s ease;

  &:active {
    transform: scale(0.95);
    box-shadow: 0 4rpx 12rpx rgba(59, 130, 246, 0.3);
  }
}
</style>
