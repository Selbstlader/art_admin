<template>
  <div class="document-upload art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <!-- 上传区域 -->
      <el-upload
        ref="uploadRef"
        class="upload-area"
        drag
        action="#"
        :auto-upload="false"
        :on-change="handleFileChange"
        :before-upload="beforeUpload"
        :file-list="fileList"
        :accept="acceptTypes"
        :limit="10"
        :http-request="customUpload"
        multiple
      >
        <el-icon class="upload-icon"><Upload /></el-icon>
        <div class="upload-text">
          <span>将文件拖到此处，或</span>
          <em>点击上传</em>
        </div>
        <template #tip>
          <div class="upload-tip"> 支持 PDF、Word、JPG、PNG、WEBP 格式，单个文件不超过 50MB </div>
        </template>
      </el-upload>

    <!-- 待上传文件列表 -->
    <div v-if="fileList.length > 0" class="file-list">
      <div class="file-list-header">
        <span>待上传文件 ({{ fileList.length }})</span>
        <el-button type="primary" size="small" @click="handleUploadAll" :loading="uploading">
          开始上传
        </el-button>
      </div>
      <div v-for="(file, index) in fileList" :key="index" class="file-item">
        <el-icon class="file-icon"><Document /></el-icon>
        <span class="file-name">{{ file.name }}</span>
        <span class="file-size">{{ formatFileSize(file.size || 0) }}</span>
        <el-tag v-if="file.status === 'success'" type="success" size="small">已上传</el-tag>
        <el-tag v-else-if="file.status === 'uploading'" type="warning" size="small">上传中</el-tag>
        <el-icon v-else class="remove-icon" @click="handleRemove(index)"><Close /></el-icon>
      </div>
    </div>

    <!-- 已上传文档列表 -->
    <div class="uploaded-list">
      <div class="uploaded-list-header">
        <span>已上传文档 ({{ uploadedDocs.length }})</span>
        <el-button
          text
          type="primary"
          size="small"
          @click="() => fetchUploadedDocs()"
          :loading="loading"
        >
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <el-skeleton :loading="loading" animated :count="3">
        <template #default>
          <el-empty
            v-if="uploadedDocs.length === 0"
            description="暂无已上传文档"
            :image-size="80"
          />
          <div v-else class="uploaded-items">
            <div v-for="doc in uploadedDocs" :key="doc.id" class="uploaded-item-wrapper">
              <div class="uploaded-item">
                <el-icon class="file-icon"><Document /></el-icon>
                <div class="file-info">
                  <span class="file-name">{{ doc.fileName }}</span>
                  <span class="file-meta">
                    {{ formatFileSize(doc.fileSize) }} · {{ formatDate(doc.createdAt) }}
                  </span>
                </div>
                <el-tag :type="getStatusType(doc.analysisStatus)" size="small">
                  <template v-if="doc.analysisStatus === 'processing'">
                    <el-icon class="is-loading"><Loading /></el-icon>
                  </template>
                  {{ getStatusText(doc.analysisStatus) }}
                </el-tag>
                <el-button
                  v-if="doc.analysisStatus === 'pending' || doc.analysisStatus === 'failed'"
                  type="primary"
                  size="small"
                  text
                  @click="handleAnalyze(doc.id)"
                >
                  {{ doc.analysisStatus === 'failed' ? '重新分析' : '分析' }}
                </el-button>
                <el-button type="danger" size="small" text @click="handleDelete(doc.id)">
                  删除
                </el-button>
              </div>
              <!-- 状态信息展示区域 -->
              <div v-if="doc.errorMessage || doc.summary" class="doc-status-info">
                <!-- 分析失败显示错误信息 -->
                <div v-if="doc.analysisStatus === 'failed'" class="error-info">
                  <el-icon><WarningFilled /></el-icon>
                  <span>{{ doc.errorMessage || '分析失败，请重试' }}</span>
                </div>
                <!-- 分析完成显示摘要和建议 -->
                <div v-else-if="doc.analysisStatus === 'completed'" class="success-info">
                  <div v-if="doc.summary" class="summary">
                    <span class="label">摘要：</span>
                    <span
                      >{{ doc.summary.slice(0, 100)
                      }}{{ doc.summary.length > 100 ? '...' : '' }}</span
                    >
                  </div>
                  <div v-if="doc.errorMessage" class="suggestion">
                    <el-icon><InfoFilled /></el-icon>
                    <span>{{ doc.errorMessage }}</span>
                  </div>
                </div>
                <!-- 分析中显示提示 -->
                <div v-else-if="doc.analysisStatus === 'processing'" class="processing-info">
                  <span>正在分析文档内容，请稍候...</span>
                </div>
              </div>
            </div>
          </div>
        </template>
      </el-skeleton>
    </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /*** Document Upload Component - 文档上传组件 ***/
  import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    Upload,
    Document,
    Close,
    Refresh,
    Loading,
    WarningFilled,
    InfoFilled
  } from '@element-plus/icons-vue'
  import {
    uploadDocument,
    getDocumentList,
    deleteDocument,
    analyzeDocument,
    type DocumentResponse
  } from '@/api/designer-document'
  import type { UploadFile } from 'element-plus'
  import { useRoute } from 'vue-router'
  import { useNotificationStore } from '@/store/modules/notification'

  const route = useRoute()
  const notificationStore = useNotificationStore()

  const emit = defineEmits<{
    (e: 'upload-success', data: any): void
    (e: 'upload-error', error: any): void
  }>()

  const uploadRef = ref()
  const fileList = ref<UploadFile[]>([])
  const uploading = ref(false)
  const loading = ref(false)
  const uploadedDocs = ref<DocumentResponse[]>([])
  const acceptTypes = '.pdf,.doc,.docx,.jpg,.jpeg,.png,.webp'
  let pollingTimer: ReturnType<typeof setInterval> | null = null

  /*** 检查是否有正在处理的文档 ***/
  const hasProcessingDocs = computed(() =>
    uploadedDocs.value.some((doc) => doc.analysisStatus === 'processing')
  )

  /*** 页面加载时获取已上传文档 ***/
  onMounted(() => {
    if (route.query.projectId) {
      fetchUploadedDocs()
    }
  })

  /*** 页面卸载时清除轮询 ***/
  onUnmounted(() => {
    stopPolling()
  })

  /*** 监听 projectId 变化 ***/
  watch(
    () => route.query.projectId,
    (newId) => {
      if (newId) {
        fetchUploadedDocs()
      }
    }
  )

  /*** 监听处理中的文档，启动/停止轮询 ***/
  watch(hasProcessingDocs, (hasProcessing) => {
    if (hasProcessing) {
      startPolling()
    } else {
      stopPolling()
    }
  })

  /*** 启动轮询 ***/
  // 修复：将轮询间隔从3秒增加到5秒，减少服务器压力
  function startPolling() {
    if (pollingTimer) return
    pollingTimer = setInterval(() => {
      fetchUploadedDocs(true) // 静默刷新
    }, 5000) // 每5秒轮询一次
  }

  /*** 停止轮询 ***/
  function stopPolling() {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  /*** 获取已上传文档列表 ***/
  async function fetchUploadedDocs(silent = false) {
    if (!route.query.projectId) return

    if (!silent) loading.value = true
    try {
      const res = await getDocumentList({
        projectId: Number(route.query.projectId),
        current: 1,
        size: 100
      })
      const newDocs = res?.data?.records || []

      // 检查状态变化，发送通知
      checkStatusChanges(uploadedDocs.value, newDocs)

      uploadedDocs.value = newDocs
    } catch (error: any) {
      console.error('获取文档列表失败:', error)
    } finally {
      if (!silent) loading.value = false
    }
  }

  /*** 检查文档状态变化并发送通知 ***/
  function checkStatusChanges(oldDocs: DocumentResponse[], newDocs: DocumentResponse[]) {
    for (const newDoc of newDocs) {
      const oldDoc = oldDocs.find((d) => d.id === newDoc.id)
      if (!oldDoc) continue

      // 从 processing 变为 completed
      if (oldDoc.analysisStatus === 'processing' && newDoc.analysisStatus === 'completed') {
        notificationStore.addNotification({
          title: `文档分析完成: ${newDoc.fileName}`,
          type: 'notice',
          time: new Date().toLocaleString('zh-CN')
        })
        ElMessage.success(`文档 "${newDoc.fileName}" 分析完成`)
      }

      // 从 processing 变为 failed
      if (oldDoc.analysisStatus === 'processing' && newDoc.analysisStatus === 'failed') {
        notificationStore.addNotification({
          title: `文档分析失败: ${newDoc.fileName}`,
          type: 'email',
          time: new Date().toLocaleString('zh-CN')
        })
        ElMessage.error(`文档 "${newDoc.fileName}" 分析失败`)
      }
    }
  }

  /*** Custom upload handler - 空实现，阻止默认上传行为 ***/
  function customUpload() {
    return Promise.resolve()
  }

  /*** Format file size ***/
  function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  /*** Format date ***/
  function formatDate(dateStr: string): string {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    return date.toLocaleDateString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  /*** Get status tag type ***/
  function getStatusType(status: string): 'success' | 'warning' | 'info' | 'danger' {
    const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
      completed: 'success',
      processing: 'warning',
      pending: 'info',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  /*** Get status text ***/
  function getStatusText(status: string): string {
    const map: Record<string, string> = {
      completed: '已分析',
      processing: '分析中',
      pending: '待分析',
      failed: '分析失败'
    }
    return map[status] || status
  }

  /*** Before upload validation - 阻止自动上传 ***/
  function beforeUpload(file: File): boolean {
    const maxSize = 50 * 1024 * 1024 // 50MB
    if (file.size > maxSize) {
      ElMessage.error(`文件 ${file.name} 超过50MB限制`)
      return false
    }
    return false
  }

  /*** Handle file change ***/
  function handleFileChange(file: UploadFile, files: UploadFile[]) {
    fileList.value = files
  }

  /*** Handle remove file ***/
  function handleRemove(index: number) {
    fileList.value.splice(index, 1)
  }

  /*** Upload all files - 手动触发上传 ***/
  async function handleUploadAll() {
    if (fileList.value.length === 0) {
      ElMessage.warning('请先选择要上传的文件')
      return
    }

    if (!route.query.projectId) {
      ElMessage.error('项目ID无效，请先选择项目')
      return
    }

    uploading.value = true
    const results: any[] = []
    const errors: string[] = []

    for (const file of fileList.value) {
      if (file.status === 'success') continue

      if (!file.raw) {
        errors.push(`${file.name}: 文件数据无效`)
        continue
      }

      try {
        file.status = 'uploading'
        const res = await uploadDocument(route.query.projectId as string, file.raw)
        file.status = 'success'
        results.push(res)
      } catch (error: any) {
        file.status = 'fail'
        errors.push(`${file.name}: ${error.message || '上传失败'}`)
      }
    }

    uploading.value = false

    if (results.length > 0) {
      ElMessage.success(`成功上传 ${results.length} 个文件`)
      emit('upload-success', results)
      // 刷新已上传列表
      fetchUploadedDocs()
    }

    if (errors.length > 0) {
      ElMessage.error(errors.join('\n'))
      emit('upload-error', errors)
    }

    fileList.value = fileList.value.filter((f) => f.status !== 'success')
  }

  /*** Handle analyze document ***/
  async function handleAnalyze(docId: number) {
    try {
      await analyzeDocument(docId)
      ElMessage.success('已开始分析文档')
      fetchUploadedDocs()
    } catch (error: any) {
      ElMessage.error(error.message || '分析失败')
    }
  }

  /*** Handle delete document ***/
  async function handleDelete(docId: number) {
    try {
      await ElMessageBox.confirm('确定要删除该文档吗？', '提示', { type: 'warning' })
      await deleteDocument(docId)
      ElMessage.success('删除成功')
      fetchUploadedDocs()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }
</script>

<style scoped lang="scss">
  .document-upload {
    .upload-area {
      width: 100%;

      :deep(.el-upload-dragger) {
        padding: 40px 20px;
        border-radius: 8px;
        background: var(--el-fill-color-lighter);

        &:hover {
          border-color: var(--el-color-primary);
        }
      }
    }

    .upload-icon {
      font-size: 48px;
      color: var(--el-text-color-secondary);
      margin-bottom: 16px;
    }

    .upload-text {
      color: var(--el-text-color-regular);
      em {
        color: var(--el-color-primary);
        font-style: normal;
      }
    }

    .upload-tip {
      margin-top: 8px;
      color: var(--el-text-color-secondary);
      font-size: 12px;
    }

    .file-list,
    .uploaded-list {
      margin-top: 20px;
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 8px;
      padding: 16px;

      .file-list-header,
      .uploaded-list-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        font-weight: 500;
      }

      .file-item {
        display: flex;
        align-items: center;
        padding: 8px 12px;
        background: var(--el-fill-color-lighter);
        border-radius: 4px;
        margin-bottom: 8px;

        &:last-child {
          margin-bottom: 0;
        }

        .file-icon {
          font-size: 20px;
          color: var(--el-color-primary);
          margin-right: 8px;
        }

        .file-name {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .file-size {
          color: var(--el-text-color-secondary);
          font-size: 12px;
          margin: 0 12px;
        }

        .remove-icon {
          cursor: pointer;
          color: var(--el-text-color-secondary);
          &:hover {
            color: var(--el-color-danger);
          }
        }
      }
    }

    .uploaded-items {
      .uploaded-item-wrapper {
        margin-bottom: 12px;

        &:last-child {
          margin-bottom: 0;
        }
      }

      .uploaded-item {
        display: flex;
        align-items: center;
        padding: 12px;
        background: var(--el-fill-color-lighter);
        border-radius: 4px;
        gap: 8px;

        .file-icon {
          font-size: 24px;
          color: var(--el-color-primary);
          flex-shrink: 0;
        }

        .file-info {
          flex: 1;
          min-width: 0;
          display: flex;
          flex-direction: column;
          gap: 4px;

          .file-name {
            font-size: 14px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .file-meta {
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }

        .is-loading {
          animation: rotating 1.5s linear infinite;
          margin-right: 4px;
        }
      }

      .doc-status-info {
        margin-top: 8px;
        padding: 8px 12px;
        border-radius: 4px;
        font-size: 12px;
        line-height: 1.5;

        .error-info {
          display: flex;
          align-items: flex-start;
          gap: 6px;
          color: var(--el-color-danger);
          background: var(--el-color-danger-light-9);
          padding: 8px;
          border-radius: 4px;

          .el-icon {
            flex-shrink: 0;
            margin-top: 2px;
          }
        }

        .success-info {
          background: var(--el-fill-color-lighter);
          padding: 8px;
          border-radius: 4px;

          .summary {
            color: var(--el-text-color-regular);
            margin-bottom: 6px;

            .label {
              color: var(--el-text-color-secondary);
            }
          }

          .suggestion {
            display: flex;
            align-items: flex-start;
            gap: 6px;
            color: var(--el-color-warning);
            margin-top: 6px;

            .el-icon {
              flex-shrink: 0;
              margin-top: 2px;
            }
          }
        }

        .processing-info {
          color: var(--el-color-primary);
          background: var(--el-color-primary-light-9);
          padding: 8px;
          border-radius: 4px;
        }
      }
    }
  }

  @keyframes rotating {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
