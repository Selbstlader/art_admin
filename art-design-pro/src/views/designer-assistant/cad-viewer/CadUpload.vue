<template>
  <div class="cad-upload-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">CAD图纸预览</span>
          <ElSelect
            v-if="projectList.length > 0"
            v-model="selectedProjectId"
            placeholder="选择项目"
            style="width: 200px"
          >
            <ElOption
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </ElSelect>
        </div>
      </template>

      <!-- 上传区域 -->
      <div v-if="!currentCadFile" class="upload-section">
        <ElUpload
          class="upload-area"
          drag
          :auto-upload="false"
          :on-change="handleFileChange"
          :disabled="!selectedProjectId"
          accept=".dwg,.dxf"
        >
          <ElIcon class="el-icon--upload"><UploadFilled /></ElIcon>
          <div class="el-upload__text"> 将CAD文件拖到此处，或<em>点击上传</em> </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 DWG、DXF 格式，文件大小不超过100MB
              <span v-if="!selectedProjectId" class="warning-text">（请先选择项目）</span>
              <br />
              <span class="tip-text">提示：DXF格式可直接解析，DWG格式需要服务器配置ODA转换器</span>
            </div>
          </template>
        </ElUpload>

        <!-- 已上传文件列表 -->
        <div v-if="cadFileList.length > 0" class="file-list">
          <h4>已上传的CAD文件</h4>
          <ElTable :data="cadFileList" size="small">
            <ElTableColumn prop="fileName" label="文件名" />
            <ElTableColumn prop="fileFormat" label="格式" width="80">
              <template #default="{ row }">
                <ElTag size="small">{{ row.fileFormat.toUpperCase() }}</ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn prop="parseStatus" label="状态" width="100">
              <template #default="{ row }">
                <ElTag :type="getStatusType(row.parseStatus)" size="small">
                  {{ getStatusText(row.parseStatus) }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn label="操作" width="150">
              <template #default="{ row }">
                <ElButton type="primary" link size="small" @click="handleViewCad(row)">
                  查看
                </ElButton>
                <ElButton type="danger" link size="small" @click="handleDeleteCad(row)">
                  删除
                </ElButton>
              </template>
            </ElTableColumn>
          </ElTable>
        </div>
      </div>

      <!-- CAD预览区域 -->
      <div v-else class="viewer-section">
        <div class="toolbar">
          <div class="toolbar-left">
            <ElButtonGroup>
              <ElButton @click="handleZoomIn" :icon="ZoomIn" />
              <ElButton @click="handleZoomOut" :icon="ZoomOut" />
              <ElButton @click="handleReset" :icon="Refresh" />
            </ElButtonGroup>
            <span class="file-info">{{ currentCadFile.fileName }}</span>
            <ElTag v-if="parseResult?.has3d" type="success" size="small">包含3D</ElTag>
          </div>
          <div class="toolbar-right">
            <ElButton type="primary" @click="openRenderDialog" :icon="Picture">
              AI生成效果图
            </ElButton>
            <ElButton type="success" @click="openRecordsDialog"> 我的效果图 </ElButton>
            <ElButton @click="handleBackToList">返回列表</ElButton>
          </div>
        </div>

        <div class="viewer-container">
          <!-- 图层控制面板 -->
          <div class="layer-panel">
            <h4>图层控制 ({{ parseResult?.layerCount || 0 }})</h4>
            <div class="layer-actions">
              <ElButton size="small" link @click="selectAllLayers">全选</ElButton>
              <ElButton size="small" link @click="deselectAllLayers">全不选</ElButton>
            </div>
            <div class="layer-list">
              <ElCheckbox
                v-for="layer in layerList"
                :key="layer.name"
                v-model="layer.visible"
                @change="handleLayerChange"
              >
                <span class="layer-name">{{ layer.name }}</span>
                <span class="layer-count">({{ layer.entityCount }})</span>
              </ElCheckbox>
            </div>
          </div>

          <!-- Three.js 渲染区域 -->
          <div class="canvas-container">
            <CadViewer
              ref="cadViewerRef"
              :parse-result="parseResult"
              :visible-layers="visibleLayers"
              @entity-click="handleEntityClick"
            />
          </div>

          <!-- 元素属性面板 -->
          <div class="property-panel">
            <h4>元素属性</h4>
            <template v-if="selectedEntity">
              <ElDescriptions :column="1" border size="small">
                <ElDescriptionsItem label="类型">{{ selectedEntity.type }}</ElDescriptionsItem>
                <ElDescriptionsItem label="图层">{{ selectedEntity.layer }}</ElDescriptionsItem>
                <ElDescriptionsItem label="颜色">
                  <div
                    class="color-preview"
                    :style="{ backgroundColor: getColorHex(selectedEntity.color) }"
                  ></div>
                  {{ selectedEntity.color }}
                </ElDescriptionsItem>
                <ElDescriptionsItem v-if="selectedEntity.lineType" label="线型">
                  {{ selectedEntity.lineType }}
                </ElDescriptionsItem>
                <template v-for="(value, key) in selectedEntity.properties" :key="key">
                  <ElDescriptionsItem :label="formatPropertyKey(key)">
                    {{ formatPropertyValue(value) }}
                  </ElDescriptionsItem>
                </template>
              </ElDescriptions>
            </template>
            <ElEmpty v-else description="点击图纸元素查看属性" :image-size="60" />
          </div>
        </div>
      </div>
    </ElCard>

    <!-- AI渲染对话框 -->
    <ElDialog v-model="showRenderDialog" title="AI生成效果图" width="700px">
      <ElForm :model="renderForm" label-width="100px">
        <ElFormItem label="参考文档">
          <div class="document-select">
            <ElCheckboxGroup v-model="renderForm.selectedDocIds">
              <ElCheckbox
                v-for="doc in projectDocuments"
                :key="doc.id"
                :value="doc.id"
                :label="doc.id"
              >
                {{ doc.fileName }}
                <ElTag size="small" type="info">{{ doc.fileType }}</ElTag>
              </ElCheckbox>
            </ElCheckboxGroup>
            <ElEmpty
              v-if="projectDocuments.length === 0"
              description="暂无项目文档"
              :image-size="40"
            />
          </div>
        </ElFormItem>
        <ElFormItem label="设计风格">
          <ElInput
            v-model="renderForm.style"
            placeholder="请输入设计风格，如：现代简约、新中式、北欧风格、工业风等"
          />
        </ElFormItem>
        <ElFormItem label="空间类型">
          <ElInput
            v-model="renderForm.roomType"
            placeholder="请输入空间类型，如：客厅、卧室、办公室、餐厅、会议室等"
          />
        </ElFormItem>
        <ElFormItem label="设计要求">
          <ElInput
            v-model="renderForm.description"
            type="textarea"
            :rows="4"
            placeholder="详细描述您的设计要求，如：&#10;- 整体色调偏暖，以米白色和原木色为主&#10;- 需要大面积落地窗，采光充足&#10;- 家具风格简约现代，注重功能性&#10;- 添加绿植装饰，营造自然氛围"
          />
        </ElFormItem>
      </ElForm>

      <!-- 配额信息 -->
      <div v-if="userQuota" class="quota-info">
        <span>今日剩余生成次数：{{ userQuota.remainingUse }} / {{ userQuota.dailyLimit }}</span>
      </div>

      <template #footer>
        <ElButton @click="showRenderDialog = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="renderLoading"
          :disabled="!userQuota || userQuota.remainingUse <= 0"
          @click="handleGenerateRender"
        >
          {{ renderLoading ? '提交中...' : '提交生成任务' }}
        </ElButton>
      </template>
    </ElDialog>

    <!-- 效果图记录列表对话框 -->
    <ElDialog
      v-model="showRecordsDialog"
      :title="currentCadFile ? `效果图 - ${currentCadFile.fileName}` : '我的效果图'"
      width="950px"
      @open="loadRenderRecords"
    >
      <div class="records-toolbar">
        <ElButton type="primary" size="small" @click="loadRenderRecords">刷新</ElButton>
      </div>
      <ElTable :data="renderRecords" v-loading="recordsLoading">
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="效果图" width="150">
          <template #default="{ row }">
            <img
              v-if="row.status === 'completed' && row.imageUrl"
              :src="getImageUrl(row.imageUrl)"
              alt="效果图"
              class="record-thumbnail"
              @click="previewImage(row.imageUrl)"
            />
            <span
              v-else-if="row.status === 'pending' || row.status === 'processing'"
              class="status-text"
              >生成中...</span
            >
            <span v-else-if="row.status === 'failed'" class="status-text error">生成失败</span>
            <span v-else class="status-text">-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="style" label="设计风格" width="120" />
        <ElTableColumn prop="roomType" label="空间类型" />
        <ElTableColumn prop="createdAt" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="120">
          <template #default="{ row }">
            <ElButton
              v-if="row.status === 'completed'"
              type="primary"
              link
              size="small"
              @click="viewRecordDetail(row)"
            >
              查看
            </ElButton>
            <ElButton type="danger" link size="small" @click="handleDeleteRecord(row)">
              删除
            </ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
      <div class="pagination-wrapper">
        <ElPagination
          v-model:current-page="recordsPagination.current"
          v-model:page-size="recordsPagination.size"
          :total="recordsPagination.total"
          layout="total, prev, pager, next"
          @current-change="loadRenderRecords"
        />
      </div>
    </ElDialog>

    <!-- 效果图详情对话框 -->
    <ElDialog v-model="showRecordDetailDialog" title="效果图详情" width="800px">
      <div v-if="currentRecord" class="record-detail">
        <div class="detail-image">
          <img :src="getImageUrl(currentRecord.imageUrl)" alt="效果图" />
        </div>
        <div class="detail-info">
          <ElDescriptions :column="2" border>
            <ElDescriptionsItem label="设计风格">{{
              currentRecord.style || '-'
            }}</ElDescriptionsItem>
            <ElDescriptionsItem label="空间类型">{{
              currentRecord.roomType || '-'
            }}</ElDescriptionsItem>
            <ElDescriptionsItem label="保存时间" :span="2">{{
              formatDate(currentRecord.createdAt)
            }}</ElDescriptionsItem>
          </ElDescriptions>
        </div>
        <div v-if="currentRecord.designProposal" class="detail-proposal">
          <h4>设计方案</h4>
          <div class="proposal-content" v-html="formatProposal(currentRecord.designProposal)"></div>
        </div>
      </div>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * CAD Upload and Viewer Page Component
   * CAD文件上传和预览页面组件
   * Requirements: 8.1, 8.2, 8.3, 8.5, 8.6
   ***/
  import { ref, computed, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { UploadFilled, ZoomIn, ZoomOut, Refresh, Picture } from '@element-plus/icons-vue'
  import type { UploadFile } from 'element-plus'
  import CadViewer from './CadViewer.vue'
  import {
    uploadCadFile,
    parseCadFile,
    getCadFileList,
    getCadParseResult,
    deleteCadFile,
    generateCadRender,
    getUserRenderQuota,
    getRenderRecords,
    deleteRenderRecord,
    type CadFileResponse,
    type CadParseResult,
    type CadEntityInfo,
    type CadLayerInfo,
    type UserRenderQuota,
    type RenderRecordItem
  } from '@/api/designer-cad'
  import { getDesignerProjects, type ProjectResponse } from '@/api/designer-project'
  import { getDocumentList } from '@/api/designer-document'

  // 路由 / Route
  const route = useRoute()

  // 项目列表 / Project list
  const projectList = ref<{ id: number; name: string }[]>([])
  const selectedProjectId = ref<number | null>(null)

  // 路由参数中的CAD文件ID / CAD file ID from route params
  const routeCadFileId = ref<number | null>(null)

  // CAD文件相关 / CAD file related
  const cadFileList = ref<CadFileResponse[]>([])
  const currentCadFile = ref<CadFileResponse | null>(null)
  const parseResult = ref<CadParseResult | null>(null)
  const selectedEntity = ref<CadEntityInfo | null>(null)

  // 图层列表 / Layer list
  const layerList = ref<(CadLayerInfo & { visible: boolean })[]>([])

  // 组件引用 / Component ref
  const cadViewerRef = ref<InstanceType<typeof CadViewer>>()

  // AI渲染相关 / AI render related
  const showRenderDialog = ref(false)
  const renderLoading = ref(false)
  const projectDocuments = ref<{ id: number; fileName: string; fileType: string }[]>([])
  const renderForm = ref({
    selectedDocIds: [] as number[],
    style: '',
    roomType: '',
    description: ''
  })

  // 用户配额 / User quota
  const userQuota = ref<UserRenderQuota | null>(null)

  // 效果图记录相关 / Render records related
  const showRecordsDialog = ref(false)
  const showRecordDetailDialog = ref(false)
  const recordsLoading = ref(false)
  const renderRecords = ref<RenderRecordItem[]>([])
  const currentRecord = ref<RenderRecordItem | null>(null)
  const recordsPagination = ref({
    current: 1,
    size: 10,
    total: 0
  })

  // 计算可见图层 / Compute visible layers
  const visibleLayers = computed(() => {
    return layerList.value.filter((l) => l.visible).map((l) => l.name)
  })

  // 加载项目列表 / Load project list
  const loadProjectList = async () => {
    try {
      const res = (await getDesignerProjects({
        current: 1,
        size: 100
      })) as any
      projectList.value = (res.data?.records || []).map((p: ProjectResponse) => ({
        id: p.id,
        name: p.name
      }))

      // 优先使用路由参数中的projectId / Prefer projectId from route params
      const routeProjectId = route.query.projectId ? Number(route.query.projectId) : null
      if (routeProjectId && projectList.value.some((p) => p.id === routeProjectId)) {
        selectedProjectId.value = routeProjectId
      } else if (projectList.value.length > 0) {
        selectedProjectId.value = projectList.value[0].id
      }

      // 加载CAD文件列表 / Load CAD file list
      await loadCadFileList()

      // 如果路由参数中有CAD文件ID，自动打开预览 / Auto open preview if route has CAD file ID
      const routeFileId = route.query.id ? Number(route.query.id) : null
      if (routeFileId) {
        routeCadFileId.value = routeFileId
        // 在文件列表中查找并打开 / Find and open in file list
        const targetFile = cadFileList.value.find((f) => f.id === routeFileId)
        if (targetFile) {
          handleViewCad(targetFile)
        } else {
          // 如果在当前项目列表中找不到，尝试直接获取解析结果 / If not found in current list, try to get parse result directly
          try {
            const res = (await getCadParseResult(routeFileId)) as any
            if (res.data) {
              // 创建一个临时的文件对象用于显示 / Create a temp file object for display
              currentCadFile.value = {
                id: routeFileId,
                projectId: routeProjectId || 0,
                fileName: res.data.fileName || 'CAD文件',
                fileFormat: 'dxf',
                parseStatus: 'completed',
                createdAt: '',
                updatedAt: ''
              } as CadFileResponse
              parseResult.value = res.data
              layerList.value = (res.data?.layers || []).map((layer: CadLayerInfo) => ({
                ...layer,
                visible: true
              }))
            }
          } catch {
            ElMessage.warning('未找到指定的CAD文件')
          }
        }
      }
    } catch {
      console.error('加载项目列表失败')
    }
  }

  // 加载CAD文件列表 / Load CAD file list
  const loadCadFileList = async () => {
    if (!selectedProjectId.value) return
    try {
      const res = (await getCadFileList({
        projectId: selectedProjectId.value,
        current: 1,
        size: 50
      })) as any
      cadFileList.value = res.data?.records || []
    } catch {
      console.error('加载CAD文件列表失败')
    }
  }

  // 文件变化处理 / Handle file change
  const handleFileChange = async (file: UploadFile) => {
    if (!file.raw || !selectedProjectId.value) return

    try {
      // 上传文件 / Upload file
      const formData = new FormData()
      formData.append('file', file.raw)
      formData.append('projectId', String(selectedProjectId.value))

      ElMessage.info('正在上传CAD文件...')
      const uploadRes = (await uploadCadFile(formData)) as any
      ElMessage.success('CAD文件上传成功')

      // 解析文件 / Parse file
      ElMessage.info('正在解析CAD文件...')
      await parseCadFile(uploadRes.data.id)
      ElMessage.success('CAD文件解析成功')

      // 刷新列表并查看 / Refresh list and view
      await loadCadFileList()

      // 查找并查看刚上传的文件 / Find and view uploaded file
      const uploadedFile = cadFileList.value.find((f) => f.id === uploadRes.data.id)
      if (uploadedFile) {
        handleViewCad(uploadedFile)
      }
    } catch (error: any) {
      ElMessage.error(error.message || 'CAD文件处理失败')
    }
  }

  // 查看CAD文件 / View CAD file
  const handleViewCad = async (file: CadFileResponse) => {
    currentCadFile.value = file
    selectedEntity.value = null

    if (file.parseStatus === 'completed') {
      try {
        const res = (await getCadParseResult(file.id)) as any
        if (!res.data) {
          ElMessage.error('解析数据不存在，请尝试重新解析')
          // 重置状态，允许重新解析 / Reset status, allow re-parse
          file.parseStatus = 'pending'
          return
        }
        parseResult.value = res.data

        // 初始化图层列表 / Initialize layer list
        layerList.value = (res.data?.layers || []).map((layer: CadLayerInfo) => ({
          ...layer,
          visible: true
        }))
      } catch (error: any) {
        ElMessage.error(error.message || '获取解析结果失败')
        // 返回列表 / Back to list
        currentCadFile.value = null
      }
    } else if (file.parseStatus === 'pending') {
      // 尝试解析 / Try to parse
      try {
        ElMessage.info('正在解析CAD文件...')
        const res = (await parseCadFile(file.id)) as any
        parseResult.value = res.data

        layerList.value = (res.data?.layers || []).map((layer: CadLayerInfo) => ({
          ...layer,
          visible: true
        }))

        ElMessage.success('解析成功')
        loadCadFileList()
      } catch (error: any) {
        ElMessage.error(error.message || '解析失败')
      }
    } else {
      ElMessage.error(`文件状态异常: ${file.errorMessage || file.parseStatus}`)
    }
  }

  // 删除CAD文件 / Delete CAD file
  const handleDeleteCad = async (file: CadFileResponse) => {
    try {
      await ElMessageBox.confirm('确定要删除该CAD文件吗？', '提示', {
        type: 'warning'
      })
      await deleteCadFile(file.id)
      ElMessage.success('删除成功')
      loadCadFileList()
    } catch {
      // 用户取消
    }
  }

  // 返回列表 / Back to list
  const handleBackToList = () => {
    currentCadFile.value = null
    parseResult.value = null
    selectedEntity.value = null
    layerList.value = []
  }

  // 加载项目文档 / Load project documents
  const loadProjectDocuments = async () => {
    if (!selectedProjectId.value) {
      projectDocuments.value = []
      return
    }
    try {
      const res = (await getDocumentList({
        projectId: selectedProjectId.value,
        current: 1,
        size: 50
      })) as any
      projectDocuments.value = res.data?.records || []
    } catch {
      console.error('加载项目文档失败')
      projectDocuments.value = []
    }
  }

  // 打开渲染对话框 / Open render dialog
  const openRenderDialog = () => {
    showRenderDialog.value = true
    loadProjectDocuments()
    loadUserQuota()
  }

  // 打开效果图记录对话框 / Open records dialog
  const openRecordsDialog = () => {
    // 重置分页 / Reset pagination
    recordsPagination.value.current = 1
    showRecordsDialog.value = true
    loadRenderRecords()
  }

  // 加载用户配额 / Load user quota
  const loadUserQuota = async () => {
    try {
      const res = (await getUserRenderQuota()) as any
      userQuota.value = res.data
    } catch {
      console.error('加载用户配额失败')
    }
  }

  // 生成AI效果图（异步提交任务）/ Generate AI render (async task submission)
  const handleGenerateRender = async () => {
    if (!currentCadFile.value) return

    if (!renderForm.value.style && !renderForm.value.roomType && !renderForm.value.description) {
      ElMessage.warning('请至少填写设计风格、空间类型或设计要求中的一项')
      return
    }

    renderLoading.value = true

    try {
      const res = (await generateCadRender({
        cadFileId: currentCadFile.value.id,
        projectId: selectedProjectId.value || undefined,
        documentIds: renderForm.value.selectedDocIds,
        style: renderForm.value.style,
        roomType: renderForm.value.roomType,
        description: renderForm.value.description
      })) as any

      if (res.data) {
        ElMessage.success(res.data.message || '任务已提交，请在"我的效果图"中查看结果')
        // 刷新配额 / Refresh quota
        loadUserQuota()
        // 关闭对话框 / Close dialog
        showRenderDialog.value = false
        // 重置表单 / Reset form
        renderForm.value = {
          selectedDocIds: [],
          style: '',
          roomType: '',
          description: ''
        }
      }
    } catch (error: any) {
      ElMessage.error(error.message || '提交任务失败')
    } finally {
      renderLoading.value = false
    }
  }

  // 加载效果图记录 / Load render records
  // 当有currentCadFile时，只加载该CAD文件的效果图
  const loadRenderRecords = async () => {
    recordsLoading.value = true
    try {
      const res = (await getRenderRecords({
        current: recordsPagination.value.current,
        size: recordsPagination.value.size,
        projectId: selectedProjectId.value || undefined,
        cadFileId: currentCadFile.value?.id || undefined
      })) as any

      renderRecords.value = res.data?.records || []
      recordsPagination.value.total = res.data?.total || 0
    } catch {
      console.error('加载效果图记录失败')
    } finally {
      recordsLoading.value = false
    }
  }

  // 查看效果图记录详情 / View record detail
  const viewRecordDetail = (record: RenderRecordItem) => {
    currentRecord.value = record
    showRecordDetailDialog.value = true
  }

  // 删除效果图记录 / Delete render record
  const handleDeleteRecord = async (record: RenderRecordItem) => {
    try {
      await ElMessageBox.confirm('确定要删除该效果图记录吗？', '提示', {
        type: 'warning'
      })
      await deleteRenderRecord(record.id)
      ElMessage.success('删除成功')
      loadRenderRecords()
    } catch {
      // 用户取消
    }
  }

  // 预览图片 / Preview image
  const previewImage = (url: string) => {
    window.open(getImageUrl(url), '_blank')
  }

  // 格式化日期 / Format date
  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleString('zh-CN')
  }

  // 缩放操作 / Zoom operations
  const handleZoomIn = () => cadViewerRef.value?.zoomIn()
  const handleZoomOut = () => cadViewerRef.value?.zoomOut()
  const handleReset = () => cadViewerRef.value?.resetView()

  // 图层操作 / Layer operations
  const selectAllLayers = () => {
    layerList.value.forEach((l) => (l.visible = true))
  }

  const deselectAllLayers = () => {
    layerList.value.forEach((l) => (l.visible = false))
  }

  const handleLayerChange = () => {
    // 触发visibleLayers计算属性更新 / Trigger visibleLayers computed update
  }

  // 元素点击 / Entity click
  const handleEntityClick = (entity: CadEntityInfo) => {
    selectedEntity.value = entity
  }

  // 状态相关 / Status related
  const getStatusType = (status: string): 'info' | 'warning' | 'success' | 'danger' => {
    const map: Record<string, 'info' | 'warning' | 'success' | 'danger'> = {
      pending: 'info',
      processing: 'warning',
      completed: 'success',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      pending: '待解析',
      processing: '解析中',
      completed: '已完成',
      failed: '失败'
    }
    return map[status] || status
  }

  // 颜色转换 / Color conversion
  const getColorHex = (colorIndex: number) => {
    const colorMap: Record<number, string> = {
      1: '#ff0000',
      2: '#ffff00',
      3: '#00ff00',
      4: '#00ffff',
      5: '#0000ff',
      6: '#ff00ff',
      7: '#ffffff',
      8: '#808080',
      9: '#c0c0c0'
    }
    return colorMap[colorIndex] || '#ffffff'
  }

  // 格式化属性 / Format properties
  const formatPropertyKey = (key: string) => {
    const keyMap: Record<string, string> = {
      x0: 'X起点',
      y0: 'Y起点',
      z0: 'Z起点',
      x1: 'X终点',
      y1: 'Y终点',
      z1: 'Z终点',
      radius: '半径',
      startAngle: '起始角度',
      endAngle: '结束角度',
      text: '文本内容'
    }
    return keyMap[key] || key
  }

  const formatPropertyValue = (value: any) => {
    if (typeof value === 'number') {
      return value.toFixed(2)
    }
    return String(value)
  }

  // 格式化设计方案（Markdown转HTML）/ Format design proposal (Markdown to HTML)
  const formatProposal = (text: string) => {
    if (!text) return ''
    // 简单的Markdown转换 / Simple Markdown conversion
    return text
      .replace(/## (.*)/g, '<h3>$1</h3>')
      .replace(/### (.*)/g, '<h4>$1</h4>')
      .replace(/- (.*)/g, '<li>$1</li>')
      .replace(/\n\n/g, '</p><p>')
      .replace(/\n/g, '<br/>')
  }

  // 获取图片完整URL / Get full image URL
  // 修复：使用环境变量替代硬编码地址
  const getImageUrl = (path: string) => {
    if (!path) return ''
    if (path.startsWith('http')) return path
    const baseUrl = import.meta.env.VITE_API_URL || import.meta.env.VITE_BASE_URL || ''
    return `${baseUrl}/uploads/${path}`
  }

  // 生命周期 / Lifecycle
  onMounted(() => {
    loadProjectList()
  })
</script>

<style scoped lang="scss">
  .cad-upload-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .upload-section {
      .upload-area {
        width: 100%;

        :deep(.el-upload-dragger) {
          width: 100%;
        }
      }

      .warning-text {
        color: var(--el-color-warning);
      }

      .file-list {
        margin-top: 24px;

        h4 {
          margin-bottom: 12px;
          font-size: 14px;
          font-weight: 500;
        }
      }
    }

    .viewer-section {
      .toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;

        .toolbar-left {
          display: flex;
          align-items: center;
          gap: 12px;

          .file-info {
            font-size: 14px;
            color: var(--el-text-color-secondary);
          }
        }
      }

      .viewer-container {
        display: flex;
        gap: 16px;
        height: 600px;

        .layer-panel {
          width: 200px;
          padding: 16px;
          background: var(--el-fill-color-light);
          border-radius: 4px;
          overflow-y: auto;

          h4 {
            margin-bottom: 8px;
            font-size: 14px;
            font-weight: 500;
          }

          .layer-actions {
            margin-bottom: 12px;
            display: flex;
            gap: 8px;
          }

          .layer-list {
            display: flex;
            flex-direction: column;
            gap: 8px;

            .layer-name {
              margin-right: 4px;
            }

            .layer-count {
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }
        }

        .canvas-container {
          flex: 1;
          border-radius: 4px;
          overflow: hidden;
        }

        .property-panel {
          width: 280px;
          padding: 16px;
          background: var(--el-fill-color-light);
          border-radius: 4px;
          overflow-y: auto;

          h4 {
            margin-bottom: 12px;
            font-size: 14px;
            font-weight: 500;
          }

          .color-preview {
            display: inline-block;
            width: 16px;
            height: 16px;
            border-radius: 2px;
            margin-right: 8px;
            vertical-align: middle;
            border: 1px solid #ddd;
          }
        }
      }
    }
  }

  .render-preview {
    margin-top: 16px;
    padding: 16px;
    background: var(--el-fill-color-light);
    border-radius: 8px;
    max-height: 400px;
    overflow-y: auto;

    h4 {
      margin-bottom: 12px;
      font-size: 16px;
      font-weight: 600;
      color: var(--el-color-primary);
    }

    h5 {
      margin: 16px 0 8px;
      font-size: 14px;
      font-weight: 500;
    }

    .design-proposal {
      .proposal-content {
        font-size: 14px;
        line-height: 1.8;
        color: var(--el-text-color-regular);

        :deep(h3) {
          font-size: 15px;
          font-weight: 600;
          margin: 16px 0 8px;
          color: var(--el-text-color-primary);
          border-bottom: 1px solid var(--el-border-color-light);
          padding-bottom: 4px;
        }

        :deep(h4) {
          font-size: 14px;
          font-weight: 500;
          margin: 12px 0 6px;
        }

        :deep(li) {
          margin-left: 16px;
          list-style-type: disc;
        }

        :deep(p) {
          margin: 8px 0;
        }
      }
    }

    .image-preview {
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid var(--el-border-color-light);

      img {
        width: 100%;
        max-width: 600px;
        border-radius: 4px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      }
    }
  }

  .toolbar-right {
    display: flex;
    gap: 8px;
  }

  .document-select {
    max-height: 200px;
    overflow-y: auto;
    padding: 12px;
    background: var(--el-fill-color-lighter);
    border-radius: 4px;
    border: 1px solid var(--el-border-color-light);

    :deep(.el-checkbox-group) {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    :deep(.el-checkbox) {
      display: flex;
      align-items: center;
      margin-right: 0;

      .el-tag {
        margin-left: 8px;
      }
    }
  }

  .tip-text {
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }

  .quota-info {
    margin-top: 16px;
    padding: 8px 12px;
    background: var(--el-fill-color-lighter);
    border-radius: 4px;
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }

  .record-thumbnail {
    width: 120px;
    height: 80px;
    object-fit: cover;
    border-radius: 4px;
    cursor: pointer;
    transition: transform 0.2s;

    &:hover {
      transform: scale(1.05);
    }
  }

  .records-toolbar {
    margin-bottom: 12px;
    display: flex;
    justify-content: flex-end;
  }

  .status-text {
    font-size: 13px;
    color: var(--el-text-color-secondary);

    &.error {
      color: var(--el-color-danger);
    }
  }

  .pagination-wrapper {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }

  .record-detail {
    .detail-image {
      text-align: center;
      margin-bottom: 16px;

      img {
        max-width: 100%;
        max-height: 400px;
        border-radius: 8px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
      }
    }

    .detail-info {
      margin-bottom: 16px;
    }

    .detail-proposal {
      h4 {
        margin-bottom: 12px;
        font-size: 15px;
        font-weight: 600;
        color: var(--el-color-primary);
      }

      .proposal-content {
        padding: 16px;
        background: var(--el-fill-color-lighter);
        border-radius: 8px;
        font-size: 14px;
        line-height: 1.8;
        max-height: 300px;
        overflow-y: auto;
      }
    }
  }
</style>
