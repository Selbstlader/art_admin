<template>
  <div class="compare-upload-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">设计比对</span>
          <span class="subtitle">上传多张设计图与多份需求文档进行智能比对分析</span>
        </div>
      </template>

      <!-- 上传区域 -->
      <ElRow :gutter="24">
        <ElCol :span="12">
          <div class="upload-box">
            <h4>
              设计图来源
              <ElTag size="small" type="info">支持多选</ElTag>
            </h4>

            <!-- 上传新文件 -->
            <div class="upload-section">
              <div class="section-label">上传新文件</div>
              <ElUpload
                class="upload-area"
                drag
                multiple
                :auto-upload="false"
                :on-change="handleFileChange"
                :on-remove="handleFileRemove"
                :file-list="fileList"
                accept=".jpg,.jpeg,.png,.webp,.pdf"
              >
                <ElIcon class="el-icon--upload"><UploadFilled /></ElIcon>
                <div class="el-upload__text">将设计图拖到此处，或<em>点击上传</em></div>
                <template #tip>
                  <div class="el-upload__tip">支持 JPG、PNG、WEBP、PDF 格式，单文件最大 50MB</div>
                </template>
              </ElUpload>
              <div v-if="designFiles.length > 0" class="file-summary">
                <span>已选择 {{ designFiles.length }} 个新文件</span>
                <ElButton type="primary" link size="small" @click="clearFiles">清空</ElButton>
              </div>
            </div>

            <!-- 选择 AI 生成效果图 -->
            <div class="existing-section">
              <div class="section-label">选择 AI 生成效果图</div>
              <ElSelect
                v-model="selectedRenderImages"
                placeholder="从已生成的效果图中选择"
                style="width: 100%"
                :loading="loadingRenderImages"
                filterable
                multiple
                collapse-tags
                collapse-tags-tooltip
              >
                <ElOption
                  v-for="item in renderRecordList"
                  :key="item.id"
                  :label="getImageLabel(item)"
                  :value="getFullImageUrl(item.imageUrl)"
                >
                  <div class="img-option">
                    <span class="option-style">{{ item.style || '默认风格' }}</span>
                    <span class="option-time">{{ formatTime(item.createdAt) }}</span>
                  </div>
                </ElOption>
              </ElSelect>
              <div v-if="selectedRenderImages.length > 0" class="file-summary">
                <span>已选择 {{ selectedRenderImages.length }} 张效果图</span>
                <ElButton type="primary" link size="small" @click="selectedRenderImages = []">
                  清空
                </ElButton>
              </div>
            </div>

            <!-- 汇总显示 -->
            <div v-if="totalSelectedCount > 0" class="total-summary">
              <ElIcon><Picture /></ElIcon>
              <span>共选择 {{ totalSelectedCount }} 个设计图</span>
            </div>
          </div>
        </ElCol>
        <ElCol :span="12">
          <div class="upload-box">
            <h4>
              选择需求文档
              <ElTag size="small" type="info">支持多选</ElTag>
            </h4>
            <ElSelect
              v-model="selectedDocIds"
              placeholder="请选择已分析的需求文档"
              style="width: 100%"
              :loading="loadingDocs"
              filterable
              multiple
              collapse-tags
              collapse-tags-tooltip
            >
              <ElOption
                v-for="doc in documentList"
                :key="doc.id"
                :label="doc.fileName"
                :value="doc.id"
              >
                <div class="doc-option">
                  <span>{{ doc.fileName }}</span>
                  <ElTag size="small" :type="getStatusType(doc.analysisStatus)">
                    {{ getStatusText(doc.analysisStatus) }}
                  </ElTag>
                </div>
              </ElOption>
            </ElSelect>
            <div v-if="selectedDocs.length > 0" class="selected-docs-info">
              <div v-for="doc in selectedDocs" :key="doc.id" class="doc-info-item">
                <span class="doc-name">{{ doc.fileName }}</span>
                <span v-if="doc.extractedStyle" class="doc-tag">{{ doc.extractedStyle }}</span>
              </div>
            </div>
          </div>
        </ElCol>
      </ElRow>

      <div class="task-name-area">
        <ElInput
          v-model="taskName"
          placeholder="输入比对任务名称（可选）"
          :prefix-icon="Edit"
          clearable
        />
      </div>

      <div class="action-area">
        <ElButton
          type="primary"
          :loading="uploading"
          :disabled="!canCompare"
          @click="handleCompare"
        >
          <ElIcon v-if="!uploading"><VideoPlay /></ElIcon>
          {{ uploading ? '分析中...' : '开始比对' }}
        </ElButton>
        <div v-if="analysisStatus === 'processing'" class="status-tip">
          <ElIcon class="is-loading"><Loading /></ElIcon>
          <span>AI 正在分析 {{ designFiles.length }} 张设计图，请稍候...</span>
        </div>
      </div>
    </ElCard>

    <!-- 历史比对记录列表 -->
    <ElCard shadow="never" class="history-card">
      <template #header>
        <div class="card-header">
          <span class="title">历史比对记录</span>
          <ElButton type="primary" link @click="loadCompareList">
            <ElIcon><Refresh /></ElIcon>
            刷新
          </ElButton>
        </div>
      </template>

      <ElTable :data="compareList" v-loading="loadingList" stripe>
        <ElTableColumn prop="name" label="任务名称" min-width="150" show-overflow-tooltip />
        <ElTableColumn label="设计图" width="100">
          <template #default="{ row }">
            <span>{{ row.designImages?.length || 0 }} 张</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="匹配度" width="100">
          <template #default="{ row }">
            <ElTag
              v-if="row.analysisStatus === 'completed'"
              :type="getScoreTagType(row.overallScore)"
            >
              {{ Math.round(row.overallScore) }}%
            </ElTag>
            <span v-else>-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.analysisStatus)">
              {{ getStatusText(row.analysisStatus) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" link size="small" @click="viewDetail(row)">
              查看详情
            </ElButton>
            <ElButton type="danger" link size="small" @click="handleDelete(row)"> 删除 </ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <div class="pagination-area">
        <ElPagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadCompareList"
          @current-change="loadCompareList"
        />
      </div>
    </ElCard>

    <!-- 详情抽屉 -->
    <ElDrawer v-model="drawerVisible" title="比对详情" size="50%" destroy-on-close>
      <template v-if="currentDetail">
        <div class="drawer-content">
          <!-- 基本信息 -->
          <div class="info-section">
            <h4>基本信息</h4>
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="任务名称">{{ currentDetail.name }}</ElDescriptionsItem>
              <ElDescriptionsItem label="状态">
                <ElTag :type="getStatusType(currentDetail.analysisStatus)">
                  {{ getStatusText(currentDetail.analysisStatus) }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="设计图数量">
                {{ currentDetail.designImages?.length || 0 }} 张
              </ElDescriptionsItem>
              <ElDescriptionsItem label="关联文档数">
                {{ currentDetail.documentIds?.length || 0 }} 份
              </ElDescriptionsItem>
              <ElDescriptionsItem label="创建时间">
                {{ formatTime(currentDetail.createdAt) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="更新时间">
                {{ formatTime(currentDetail.updatedAt) }}
              </ElDescriptionsItem>
            </ElDescriptions>
          </div>

          <!-- 设计图列表 -->
          <div class="info-section">
            <h4>设计图文件</h4>
            <div class="image-grid">
              <div
                v-for="(img, index) in currentDetail.designImages"
                :key="index"
                class="image-item"
              >
                <ElImage
                  :src="getFullImageUrl(img.filePath)"
                  fit="cover"
                  class="image-preview"
                  :preview-src-list="getPreviewList()"
                  :initial-index="index"
                >
                  <template #error>
                    <div class="image-error">
                      <ElIcon><Document /></ElIcon>
                      <span>{{ img.fileType }}</span>
                    </div>
                  </template>
                </ElImage>
                <div class="image-info">
                  <span class="file-name">{{ img.fileName }}</span>
                  <span class="file-meta">
                    <ElTag size="small">{{ img.fileType }}</ElTag>
                    <span v-if="img.fileSize > 0">{{ formatFileSize(img.fileSize) }}</span>
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分析结果 -->
          <div v-if="currentDetail.analysisStatus === 'completed'" class="info-section">
            <h4>分析结果</h4>
            <div class="score-display">
              <ElProgress
                type="dashboard"
                :percentage="Math.round(currentDetail.overallScore)"
                :color="getScoreColor(currentDetail.overallScore)"
                :width="120"
              />
              <span class="score-label">整体匹配度</span>
            </div>

            <ElCollapse v-model="drawerCollapse">
              <ElCollapseItem name="match">
                <template #title>
                  <span>匹配项 ({{ currentDetail.matchItems?.length || 0 }})</span>
                </template>
                <div v-for="(item, i) in currentDetail.matchItems" :key="i" class="result-item">
                  <p><strong>需求:</strong> {{ item.requirement }}</p>
                  <p><strong>匹配:</strong> {{ item.designMatch }}</p>
                  <ElProgress :percentage="item.score" :stroke-width="6" />
                </div>
                <ElEmpty
                  v-if="!currentDetail.matchItems?.length"
                  description="暂无"
                  :image-size="40"
                />
              </ElCollapseItem>

              <ElCollapseItem name="deviation">
                <template #title>
                  <span>偏差项 ({{ currentDetail.deviationItems?.length || 0 }})</span>
                </template>
                <ElAlert
                  v-for="(item, i) in currentDetail.deviationItems"
                  :key="i"
                  :type="getSeverityType(item.severity)"
                  :title="item.content"
                  :closable="false"
                  class="mb-2"
                >
                  <p>位置: {{ item.location }}</p>
                </ElAlert>
                <ElEmpty
                  v-if="!currentDetail.deviationItems?.length"
                  description="暂无"
                  :image-size="40"
                />
              </ElCollapseItem>

              <ElCollapseItem name="suggestion">
                <template #title>
                  <span>建议 ({{ currentDetail.suggestions?.length || 0 }})</span>
                </template>
                <div v-for="(item, i) in currentDetail.suggestions" :key="i" class="result-item">
                  <p>{{ item.content }}</p>
                  <div class="suggestion-meta">
                    <ElTag size="small" :type="getPriorityType(item.priority)">
                      {{ getPriorityText(item.priority) }}
                    </ElTag>
                    <span>成本影响: {{ item.costImpact }}</span>
                  </div>
                </div>
                <ElEmpty
                  v-if="!currentDetail.suggestions?.length"
                  description="暂无"
                  :image-size="40"
                />
              </ElCollapseItem>
            </ElCollapse>
          </div>

          <!-- 错误信息 -->
          <ElAlert
            v-if="currentDetail.analysisStatus === 'failed' && currentDetail.errorMessage"
            type="error"
            :title="currentDetail.errorMessage"
            :closable="false"
          />
        </div>
      </template>

      <template #footer>
        <ElButton @click="drawerVisible = false">关闭</ElButton>
        <ElButton type="danger" @click="handleDelete(currentDetail)">删除</ElButton>
      </template>
    </ElDrawer>
  </div>
</template>

<script setup lang="ts">
  /***
   * Compare Upload Component
   * 设计比对上传页面组件 - 支持多图片、多文档、CAD文件、历史记录
   ***/
  import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    UploadFilled,
    VideoPlay,
    Loading,
    Edit,
    Refresh,
    Document,
    Picture
  } from '@element-plus/icons-vue'
  import type { UploadFile, UploadUserFile } from 'element-plus'
  import { getDocumentList, type DocumentResponse } from '@/api/designer-document'
  import {
    uploadDesignImages,
    uploadDesignImagesWithExisting,
    analyzeDesignCompare,
    getDesignCompareDetail,
    getDesignCompareList,
    deleteDesignCompare,
    type DesignCompareResponse
  } from '@/api/designer-compare'
  import { getRenderRecords, type RenderRecordItem } from '@/api/designer-cad'
  import { useNotificationStore } from '@/store/modules/notification'

  const route = useRoute()
  const notificationStore = useNotificationStore()

  // 上传相关状态
  const uploading = ref(false)
  const loadingDocs = ref(false)
  const selectedDocIds = ref<number[]>([])
  const designFiles = ref<File[]>([])
  const fileList = ref<UploadUserFile[]>([])
  const documentList = ref<DocumentResponse[]>([])
  const taskName = ref('')
  const analysisStatus = ref('')
  const pollingTimer = ref<number | null>(null)

  // 设计图来源 - 同时支持上传和选择
  const selectedRenderImages = ref<string[]>([]) // AI 效果图 URL 列表
  const loadingRenderImages = ref(false)
  const renderRecordList = ref<RenderRecordItem[]>([]) // AI 效果图记录列表

  // 计算总选择数量
  const totalSelectedCount = computed(
    () => designFiles.value.length + selectedRenderImages.value.length
  )

  // 列表相关状态
  const loadingList = ref(false)
  const compareList = ref<DesignCompareResponse[]>([])
  const pagination = ref({ current: 1, size: 10, total: 0 })

  // 抽屉相关状态
  const drawerVisible = ref(false)
  const currentDetail = ref<DesignCompareResponse | null>(null)
  const drawerCollapse = ref(['match', 'deviation', 'suggestion'])

  // 计算属性
  const projectId = computed(() => Number(route.query.projectId) || 0)
  const canCompare = computed(() => {
    const hasImages = designFiles.value.length > 0 || selectedRenderImages.value.length > 0
    return hasImages && selectedDocIds.value.length > 0
  })
  const selectedDocs = computed(() =>
    documentList.value.filter((d) => selectedDocIds.value.includes(d.id))
  )

  // 加载文档列表
  const loadDocuments = async () => {
    if (!projectId.value) return
    loadingDocs.value = true
    try {
      const res = await getDocumentList({
        projectId: projectId.value,
        current: 1,
        size: 100,
        analysisStatus: 'completed'
      })
      if (res.code === 200 && res.data) {
        documentList.value = res.data.records || []
      }
    } catch (error) {
      console.error('加载文档列表失败:', error)
    } finally {
      loadingDocs.value = false
    }
  }

  // 加载比对记录列表
  const loadCompareList = async () => {
    if (!projectId.value) return
    loadingList.value = true
    try {
      const res = await getDesignCompareList({
        projectId: projectId.value,
        current: pagination.value.current,
        size: pagination.value.size
      })
      if (res.code === 200 && res.data) {
        compareList.value = res.data.records || []
        pagination.value.total = res.data.total
      }
    } catch (error) {
      console.error('加载比对列表失败:', error)
    } finally {
      loadingList.value = false
    }
  }

  // 加载 AI 生成效果图列表
  const loadRenderImages = async () => {
    if (!projectId.value) return
    loadingRenderImages.value = true
    try {
      const res = (await getRenderRecords({
        projectId: projectId.value,
        current: 1,
        size: 50
      })) as any
      if (res.code === 200 && res.data) {
        // 只显示已完成的效果图
        renderRecordList.value = (res.data.records || []).filter(
          (r: RenderRecordItem) => r.status === 'completed' && r.imageUrl
        )
      }
    } catch (error) {
      console.error('加载效果图列表失败:', error)
    } finally {
      loadingRenderImages.value = false
    }
  }

  // 获取效果图标签
  const getImageLabel = (item: RenderRecordItem) => {
    const style = item.style || '默认风格'
    const time = formatTime(item.createdAt)
    return `${style} - ${time}`
  }

  // 获取完整图片URL
  const getFullImageUrl = (imageUrl: string) => {
    if (!imageUrl) return ''
    // 如果已经是完整URL，直接返回
    if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://')) {
      return imageUrl
    }
    // 拼接后端基础URL
    const baseUrl = import.meta.env.VITE_API_BASE_URL || ''
    return `${baseUrl}/uploads/${imageUrl}`
  }

  // 获取预览图片列表
  const getPreviewList = () => {
    if (!currentDetail.value?.designImages) return []
    return currentDetail.value.designImages.map((img) => getFullImageUrl(img.filePath))
  }

  // 文件操作
  const handleFileChange = (file: UploadFile) => {
    if (file.raw) {
      if (file.raw.size > 50 * 1024 * 1024) {
        ElMessage.warning(`文件 ${file.name} 超过50MB限制`)
        return
      }
      designFiles.value.push(file.raw)
    }
  }

  const handleFileRemove = (file: UploadFile) => {
    const index = designFiles.value.findIndex((f) => f.name === file.name)
    if (index > -1) designFiles.value.splice(index, 1)
  }

  const clearFiles = () => {
    designFiles.value = []
    fileList.value = []
  }

  // 开始比对
  const handleCompare = async () => {
    const hasImages = designFiles.value.length > 0 || selectedRenderImages.value.length > 0

    if (!hasImages || !selectedDocIds.value.length || !projectId.value) {
      ElMessage.warning('请选择设计文件并选择需求文档')
      return
    }

    uploading.value = true
    analysisStatus.value = ''

    try {
      let uploadRes

      if (designFiles.value.length > 0 && selectedRenderImages.value.length === 0) {
        // 只有上传文件
        uploadRes = await uploadDesignImages(
          projectId.value,
          selectedDocIds.value,
          designFiles.value,
          taskName.value
        )
      } else if (designFiles.value.length === 0 && selectedRenderImages.value.length > 0) {
        // 只有选择的效果图
        uploadRes = await uploadDesignImagesWithExisting(
          projectId.value,
          selectedDocIds.value,
          selectedRenderImages.value,
          taskName.value
        )
      } else {
        // 同时有上传文件和选择的效果图，先上传文件，再追加效果图
        uploadRes = await uploadDesignImages(
          projectId.value,
          selectedDocIds.value,
          designFiles.value,
          taskName.value,
          selectedRenderImages.value // 传递效果图 URL
        )
      }

      if (uploadRes.code !== 200 || !uploadRes.data) {
        throw new Error(uploadRes.msg || '创建比对任务失败')
      }

      const compareId = uploadRes.data.id
      const imageCount = uploadRes.data.designImages?.length || 0
      ElMessage.success(`已准备 ${imageCount} 个文件，正在分析...`)

      const analyzeRes = await analyzeDesignCompare(compareId)
      if (analyzeRes.code !== 200) {
        throw new Error(analyzeRes.msg || '启动分析失败')
      }

      analysisStatus.value = 'processing'
      startPolling(compareId)
    } catch (error: any) {
      ElMessage.error(error.message || '比对失败')
      uploading.value = false
    }
  }

  // 轮询
  const startPolling = (compareId: number) => {
    stopPolling()
    pollingTimer.value = window.setInterval(async () => {
      try {
        const res = await getDesignCompareDetail(compareId)
        if (res.code === 200 && res.data) {
          const data = res.data
          analysisStatus.value = data.analysisStatus

          if (data.analysisStatus === 'completed') {
            stopPolling()
            uploading.value = false
            analysisStatus.value = ''
            ElMessage.success('比对分析完成')
            // 添加本地通知
            ;(notificationStore as any).addNotification({
              title: `设计比对完成，匹配度: ${Math.round(data.overallScore)}%`,
              time: new Date().toLocaleString(),
              type: 'notice'
            })
            // 重新加载列表
            await loadCompareList()
            // 清空表单
            clearFiles()
            selectedRenderImages.value = []
            taskName.value = ''
          } else if (data.analysisStatus === 'failed') {
            stopPolling()
            uploading.value = false
            analysisStatus.value = ''
            ElMessage.error(data.errorMessage || '分析失败')
            await loadCompareList()
          }
        }
      } catch (error) {
        console.error('轮询失败:', error)
      }
    }, 2000)

    setTimeout(
      () => {
        if (pollingTimer.value) {
          stopPolling()
          uploading.value = false
          if (analysisStatus.value === 'processing') {
            ElMessage.warning('分析超时，请稍后刷新查看结果')
          }
        }
      },
      10 * 60 * 1000
    )
  }

  const stopPolling = () => {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
      pollingTimer.value = null
    }
  }

  // 查看详情
  const viewDetail = (row: DesignCompareResponse) => {
    currentDetail.value = row
    drawerVisible.value = true
  }

  // 删除
  const handleDelete = async (row: DesignCompareResponse | null) => {
    if (!row) return
    try {
      await ElMessageBox.confirm('确定要删除该比对记录吗？', '提示', { type: 'warning' })
      const res = await deleteDesignCompare(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        drawerVisible.value = false
        loadCompareList()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch {
      // 取消删除
    }
  }

  // 辅助函数
  const formatTime = (time: string) => {
    if (!time) return '-'
    return new Date(time).toLocaleString()
  }

  const formatFileSize = (size: number) => {
    if (size < 1024) return size + ' B'
    if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
    return (size / 1024 / 1024).toFixed(1) + ' MB'
  }

  const getScoreColor = (score: number) => {
    if (score >= 80) return '#67c23a'
    if (score >= 60) return '#e6a23c'
    return '#f56c6c'
  }

  const getScoreTagType = (score: number): 'success' | 'warning' | 'danger' => {
    if (score >= 80) return 'success'
    if (score >= 60) return 'warning'
    return 'danger'
  }

  const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' => {
    const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
      completed: 'success',
      processing: 'warning',
      pending: 'info',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      completed: '已完成',
      processing: '分析中',
      pending: '待分析',
      failed: '失败'
    }
    return map[status] || status
  }

  const getSeverityType = (severity: string): 'success' | 'warning' | 'info' | 'error' => {
    const map: Record<string, 'success' | 'warning' | 'info' | 'error'> = {
      low: 'info',
      medium: 'warning',
      high: 'error'
    }
    return map[severity] || 'info'
  }

  const getPriorityType = (priority: string): 'success' | 'warning' | 'info' | 'danger' => {
    const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
      low: 'info',
      medium: 'warning',
      high: 'danger'
    }
    return map[priority] || 'info'
  }

  const getPriorityText = (priority: string) => {
    const map: Record<string, string> = { low: '低', medium: '中', high: '高' }
    return map[priority] || priority
  }

  watch(projectId, () => {
    loadDocuments()
    loadCompareList()
    loadRenderImages()
  })

  onMounted(() => {
    loadDocuments()
    loadCompareList()
    loadRenderImages()
  })

  onUnmounted(() => {
    stopPolling()
  })
</script>

<style scoped lang="scss">
  .compare-upload-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }

      .subtitle {
        font-size: 12px;
        color: var(--el-text-color-secondary);
        margin-left: 8px;
      }
    }

    .upload-box {
      h4 {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 12px;
        font-size: 14px;
        font-weight: 500;
      }

      .upload-area {
        width: 100%;

        :deep(.el-upload-dragger) {
          width: 100%;
          height: 120px;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
        }
      }

      .section-label {
        font-size: 13px;
        color: var(--el-text-color-secondary);
        margin-bottom: 8px;
      }

      .upload-section {
        margin-bottom: 16px;
      }

      .existing-section {
        margin-bottom: 12px;
      }

      .total-summary {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px;
        background: var(--el-color-primary-light-9);
        border-radius: 4px;
        color: var(--el-color-primary);
        font-size: 14px;
        font-weight: 500;
      }

      .doc-option,
      .img-option {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        .option-style {
          font-size: 13px;
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .option-time {
          font-size: 12px;
          color: var(--el-text-color-secondary);
          margin-left: 12px;
          flex-shrink: 0;
        }
      }

      .file-summary,
      .selected-docs-info {
        margin-top: 12px;
        padding: 8px 12px;
        background: var(--el-fill-color-light);
        border-radius: 4px;
        font-size: 13px;
      }

      .file-summary {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .selected-docs-info {
        max-height: 120px;
        overflow-y: auto;

        .doc-info-item {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 4px 0;

          .doc-name {
            flex: 1;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .doc-tag {
            color: var(--el-color-primary);
            font-size: 12px;
          }
        }
      }
    }

    .task-name-area {
      margin-top: 16px;
      max-width: 400px;
    }

    .action-area {
      margin-top: 24px;
      text-align: center;

      .status-tip {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        margin-left: 16px;
        color: var(--el-color-primary);
        font-size: 14px;

        .is-loading {
          animation: rotating 2s linear infinite;
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

    .history-card {
      margin-top: 16px;
    }

    .pagination-area {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }
  }

  // 抽屉样式
  .drawer-content {
    .info-section {
      margin-bottom: 24px;

      h4 {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 12px;
        padding-bottom: 8px;
        border-bottom: 1px solid var(--el-border-color-lighter);
      }
    }

    .file-list {
      .file-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 12px;
        background: var(--el-fill-color-light);
        border-radius: 4px;
        margin-bottom: 8px;

        .file-name {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .file-size {
          color: var(--el-text-color-secondary);
          font-size: 12px;
        }
      }
    }

    .image-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;

      .image-item {
        border-radius: 8px;
        overflow: hidden;
        background: var(--el-fill-color-light);

        .image-preview {
          width: 100%;
          height: 120px;
          display: block;
          cursor: pointer;
        }

        .image-error {
          width: 100%;
          height: 120px;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          gap: 8px;
          color: var(--el-text-color-secondary);
          background: var(--el-fill-color);

          .el-icon {
            font-size: 32px;
          }
        }

        .image-info {
          padding: 8px 12px;

          .file-name {
            display: block;
            font-size: 13px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            margin-bottom: 4px;
          }

          .file-meta {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }
      }
    }

    .score-display {
      text-align: center;
      margin-bottom: 16px;

      .score-label {
        display: block;
        margin-top: 8px;
        color: var(--el-text-color-secondary);
        font-size: 13px;
      }
    }

    .result-item {
      padding: 12px;
      background: var(--el-fill-color-light);
      border-radius: 4px;
      margin-bottom: 8px;

      p {
        margin: 4px 0;
        font-size: 13px;
      }

      .suggestion-meta {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-top: 8px;
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }

    .mb-2 {
      margin-bottom: 8px;
    }
  }
</style>
