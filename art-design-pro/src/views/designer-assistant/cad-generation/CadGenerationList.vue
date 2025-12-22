<template>
  <div class="cad-generation-page">
    <!-- 页面标题 / Page Header -->
    <div class="page-header">
      <div class="left">
        <ElButton link @click="handleBack" v-if="projectId">
          <ElIcon>
            <ArrowLeft />
          </ElIcon>
          返回项目
        </ElButton>
        <h2>AI生成CAD</h2>
      </div>
      <ElButton type="primary" @click="handleCreate">
        <ElIcon>
          <Plus />
        </ElIcon>
        新建生成任务
      </ElButton>
    </div>

    <!-- 筛选条件 / Filter -->
    <ElCard class="filter-card" shadow="never">
      <ElForm :inline="true" :model="filterForm">
        <ElFormItem label="项目" v-if="!projectId">
          <ElSelect v-model="filterForm.projectId" placeholder="选择项目" clearable>
            <ElOption v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="生成类型">
          <ElInput
            v-model="filterForm.generationType"
            placeholder="输入类型筛选"
            clearable
            style="width: 180px"
          />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSelect v-model="filterForm.status" placeholder="选择状态" clearable>
            <ElOption label="待处理" value="pending" />
            <ElOption label="处理中" value="processing" />
            <ElOption label="已完成" value="completed" />
            <ElOption label="失败" value="failed" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch">
            <ElIcon>
              <Search />
            </ElIcon>
            搜索
          </ElButton>
          <ElButton @click="handleReset">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 任务列表 / Task List -->
    <ElCard shadow="never">
      <ElTable :data="taskList" v-loading="loading" stripe>
        <ElTableColumn prop="id" label="ID" width="80" />
        <ElTableColumn label="项目" min-width="150" v-if="!projectId">
          <template #default="{ row }">
            {{ row.projectName || '-' }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="生成类型" width="120">
          <template #default="{ row }">
            <ElTag>{{ row.generationLabel }}</ElTag>
          </template>
        </ElTableColumn>

        <ElTableColumn label="状态" width="120">
          <template #default="{ row }">
            <div class="status-cell">
              <ElTag :type="getStatusType(row.status)">{{ row.statusLabel }}</ElTag>
              <!-- <ElProgress
                v-if="row.status === 'processing'"
                :percentage="row.progress"
                :stroke-width="4"
                style="width: 60px; margin-left: 8px"
              /> -->
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn label="详情">
          <template #default="{ row }">
            {{ row.prompt }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="文件" width="150">
          <template #default="{ row }">
            <template v-if="row.status === 'completed' && row.resultFilePath">
              <ElLink
                type="primary"
                :href="getFullImageUrl(row.resultFilePath)"
                target="_blank"
                :underline="false"
              >
                <ElIcon>
                  <Document />
                </ElIcon>
                下载DXF
              </ElLink>
            </template>
            <span v-else class="text-secondary">-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="耗时" width="100">
          <template #default="{ row }">
            {{ row.processingTime ? `${row.processingTime}s` : '-' }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>

        <ElTableColumn label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <ElButton
              v-if="row.status === 'completed' && !row.resultFileId"
              type="primary"
              size="small"
              @click="handleConfirm(row)"
            >
              确认保存
            </ElButton>
            <ElButton
              v-if="row.status === 'completed' && row.resultFileId"
              type="success"
              size="small"
              @click="handlePreview(row)"
            >
              查看CAD
            </ElButton>
            <ElButton
              v-if="row.status === 'failed'"
              type="warning"
              size="small"
              @click="handleRetry(row)"
            >
              重试
            </ElButton>
            <ElButton type="danger" size="small" @click="handleDelete(row)">删除</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <!-- 分页 / Pagination -->
      <div class="pagination-wrapper">
        <ElPagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTasks"
          @current-change="loadTasks"
        />
      </div>
    </ElCard>

    <!-- 新建任务对话框 / Create Task Dialog -->
    <ElDialog v-model="createDialogVisible" title="新建CAD生成任务" width="600px">
      <ElForm ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <ElFormItem label="项目" prop="projectId" v-if="!projectId">
          <ElSelect
            v-model="createForm.projectId"
            placeholder="选择项目"
            @change="handleProjectChange"
            style="width: 100%"
          >
            <ElOption v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="关联文档" prop="documentIds">
          <ElSelect
            v-model="createForm.documentIds"
            placeholder="选择关联文档（可多选，将自动提取设计信息）"
            multiple
            style="width: 100%"
            @change="handleDocumentChange"
          >
            <ElOption
              v-for="doc in projectDocuments"
              :key="doc.id"
              :label="doc.fileName"
              :value="doc.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="生成类型" prop="generationType">
          <ElInput
            v-model="createForm.generationType"
            placeholder="输入生成类型，如：平面布局图、立面图、天花图、电气图等"
            style="width: 100%"
          />
        </ElFormItem>
        <ElFormItem label="设计要求" prop="prompt">
          <ElInput
            v-model="createForm.prompt"
            type="textarea"
            :rows="4"
            placeholder="请描述您的设计要求，如：现代简约风格，开放式厨房，主卧带独立卫生间..."
          />
        </ElFormItem>
        <ElDivider>高级参数（可选）</ElDivider>
        <ElRow :gutter="16">
          <ElCol :span="12">
            <ElFormItem label="宽度(米)">
              <ElInputNumber
                v-model="createForm.parameters.width"
                :min="1"
                :max="100"
                style="width: 100%"
              />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="高度(米)">
              <ElInputNumber
                v-model="createForm.parameters.height"
                :min="1"
                :max="100"
                style="width: 100%"
              />
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElFormItem label="比例">
          <ElSelect
            v-model="createForm.parameters.scale"
            placeholder="选择比例"
            style="width: 100%"
          >
            <ElOption label="1:50" value="1:50" />
            <ElOption label="1:100" value="1:100" />
            <ElOption label="1:200" value="1:200" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="设计风格">
          <ElInput
            v-model="createForm.parameters.style"
            placeholder="输入设计风格，如：现代简约、中式、欧式、北欧等"
            style="width: 100%"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="createDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="submitCreate" :loading="creating">开始生成</ElButton>
      </template>
    </ElDialog>

    <!-- 确认保存对话框 / Confirm Save Dialog -->
    <ElDialog v-model="confirmDialogVisible" title="确认保存CAD文件" width="400px">
      <ElForm :model="confirmForm" label-width="80px">
        <ElFormItem label="文件名">
          <ElInput v-model="confirmForm.fileName" placeholder="输入文件名（可选）" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="confirmDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="submitConfirm" :loading="confirming">确认保存</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * CAD Generation List Component
   * AI生成CAD任务列表组件
   ***/
  import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
  import { ArrowLeft, Plus, Search, Document } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import {
    getCadGenerationList,
    createCadGeneration,
    deleteCadGeneration,
    retryCadGeneration,
    confirmCadFile,
    type CadGenerationResponse,
    type CadGenerationListResponse
  } from '@/api/designer-cad-generation'
  import { getDesignerProjects } from '@/api/designer-project'
  import {
    getDocumentList,
    type DocumentListResponse,
    type DocumentResponse
  } from '@/api/designer-document'

  /*** Type definitions for API responses ***/
  interface BaseResponse<T> {
    code: number
    msg?: string
    data?: T
  }

  interface ProjectListResponse {
    records: { id: number; name: string }[]
    total: number
  }

  // 获取完整图片URL / Get full image URL
  const baseUrl = import.meta.env.VITE_BASE_URL || ''
  const getFullImageUrl = (url: string | undefined): string => {
    if (!url) return ''
    // 如果已经是完整URL则直接返回 / Return directly if already full URL
    if (url.startsWith('http://') || url.startsWith('https://')) {
      return url
    }
    // 拼接基础URL / Concatenate base URL
    return `${baseUrl}${url}`
  }

  const router = useRouter()
  const route = useRoute()

  /*** State ***/
  const loading = ref(false)
  const creating = ref(false)
  const confirming = ref(false)
  const createDialogVisible = ref(false)
  const confirmDialogVisible = ref(false)
  const taskList = ref<CadGenerationResponse[]>([])
  const projects = ref<{ id: number; name: string }[]>([])
  const projectDocuments = ref<DocumentResponse[]>([])
  const createFormRef = ref<FormInstance>()
  const currentTask = ref<CadGenerationResponse | null>(null)

  // 从路由获取项目ID / Get project ID from route
  const projectId = ref<number | undefined>(
    route.query.projectId ? Number(route.query.projectId) : undefined
  )

  const filterForm = reactive({
    projectId: projectId.value,
    generationType: '',
    status: ''
  })

  const pagination = reactive({
    current: 1,
    size: 10,
    total: 0
  })

  const createForm = reactive({
    projectId: projectId.value,
    documentIds: [] as number[],
    generationType: '',
    prompt: '',
    parameters: {
      width: undefined as number | undefined,
      height: undefined as number | undefined,
      scale: '1:100',
      style: ''
    }
  })

  const confirmForm = reactive({
    fileName: ''
  })

  const createRules: FormRules = {
    projectId: [{ required: true, message: '请选择项目', trigger: 'change' }],
    generationType: [{ required: true, message: '请输入生成类型', trigger: 'blur' }]
  }

  /*** Methods ***/
  // 加载任务列表 / Load task list
  const loadTasks = async () => {
    loading.value = true
    try {
      const res = (await getCadGenerationList({
        projectId: filterForm.projectId,
        generationType: filterForm.generationType || undefined,
        status: filterForm.status || undefined,
        current: pagination.current,
        size: pagination.size
      })) as unknown as BaseResponse<CadGenerationListResponse>

      if (res.code === 200 && res.data) {
        taskList.value = res.data.records || []
        pagination.total = res.data.total || 0
      }
    } catch {
      ElMessage.error('加载任务列表失败')
    } finally {
      loading.value = false
    }
  }

  // 加载项目列表 / Load project list
  const loadProjects = async () => {
    try {
      const res = (await getDesignerProjects({
        current: 1,
        size: 100
      })) as unknown as BaseResponse<ProjectListResponse>
      if (res.code === 200 && res.data) {
        projects.value = res.data.records || []
      }
    } catch {
      console.error('加载项目列表失败')
    }
  }

  // 项目变更时加载文档 / Load documents when project changes
  const handleProjectChange = async (pId: number) => {
    createForm.documentIds = []
    if (!pId) {
      projectDocuments.value = []
      return
    }
    try {
      const res = (await getDocumentList({
        projectId: pId,
        current: 1,
        size: 100
      })) as unknown as BaseResponse<DocumentListResponse>
      if (res.code === 200 && res.data) {
        projectDocuments.value = res.data.records || []
      }
    } catch {
      console.error('加载项目文档失败')
      projectDocuments.value = []
    }
  }

  // 文档选择变更时提取信息 / Extract info when document selection changes
  const handleDocumentChange = (docIds: number[]) => {
    if (!docIds || docIds.length === 0) return

    // 从选中的文档中提取信息 / Extract info from selected documents
    const selectedDocs = projectDocuments.value.filter((doc) => docIds.includes(doc.id))

    // 收集设计要求（摘要和关键词）/ Collect design requirements
    const prompts: string[] = []
    let extractedStyle = ''
    let extractedArea = 0

    selectedDocs.forEach((doc) => {
      // 提取摘要 / Extract summary
      if (doc.summary) {
        prompts.push(doc.summary)
      }
      // 提取关键词 / Extract keywords
      if (doc.keywords && doc.keywords.length > 0) {
        prompts.push(`关键词：${doc.keywords.join('、')}`)
      }
      // 提取功能区域 / Extract functional zones
      if (doc.functionalZones && doc.functionalZones.length > 0) {
        prompts.push(`功能区域：${doc.functionalZones.join('、')}`)
      }
      // 提取设计风格（取第一个有值的）/ Extract style
      if (!extractedStyle && doc.extractedStyle) {
        extractedStyle = doc.extractedStyle
      }
      // 提取面积（取最大值）/ Extract area
      if (doc.extractedArea && doc.extractedArea > extractedArea) {
        extractedArea = doc.extractedArea
      }
    })

    // 自动填充设计要求 / Auto-fill design requirements
    if (prompts.length > 0 && !createForm.prompt) {
      createForm.prompt = prompts.join('\n')
    }

    // 自动填充设计风格 / Auto-fill style
    if (extractedStyle && !createForm.parameters.style) {
      createForm.parameters.style = extractedStyle
    }

    // 自动填充面积（假设为正方形估算宽高）/ Auto-fill dimensions
    if (extractedArea > 0 && !createForm.parameters.width && !createForm.parameters.height) {
      const side = Math.round(Math.sqrt(extractedArea))
      createForm.parameters.width = side
      createForm.parameters.height = side
    }
  }

  // 获取状态类型 / Get status type
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      pending: 'info',
      processing: 'warning',
      completed: 'success',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  // 格式化日期 / Format date
  const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  // 搜索 / Search
  const handleSearch = () => {
    pagination.current = 1
    loadTasks()
  }

  // 重置 / Reset
  const handleReset = () => {
    filterForm.projectId = projectId.value
    filterForm.generationType = ''
    filterForm.status = ''
    pagination.current = 1
    loadTasks()
  }

  // 返回项目 / Back to project
  const handleBack = () => {
    if (projectId.value) {
      router.push(`/designer/designer-assistant/project/ProjectDetail/${projectId.value}`)
    }
  }

  // 新建任务 / Create task
  const handleCreate = () => {
    createForm.projectId = projectId.value
    createForm.documentIds = []
    createForm.generationType = ''
    createForm.prompt = ''
    createForm.parameters = {
      width: undefined,
      height: undefined,
      scale: '1:100',
      style: ''
    }
    if (projectId.value) {
      handleProjectChange(projectId.value)
    }
    createDialogVisible.value = true
  }

  // 提交创建 / Submit create
  const submitCreate = async () => {
    if (!createFormRef.value) return
    await createFormRef.value.validate(async (valid) => {
      if (!valid) return

      creating.value = true
      try {
        const res = (await createCadGeneration({
          projectId: createForm.projectId!,
          documentIds: createForm.documentIds,
          generationType: createForm.generationType,
          prompt: createForm.prompt,
          parameters: createForm.parameters
        })) as unknown as BaseResponse<CadGenerationResponse>

        if (res.code === 200) {
          ElMessage.success('任务创建成功，正在生成中...')
          createDialogVisible.value = false
          loadTasks()
        } else {
          ElMessage.error(res.msg || '创建失败')
        }
      } catch {
        ElMessage.error('创建失败')
      } finally {
        creating.value = false
      }
    })
  }

  // 确认保存 / Confirm save
  const handleConfirm = (row: CadGenerationResponse) => {
    currentTask.value = row
    confirmForm.fileName = ''
    confirmDialogVisible.value = true
  }

  // 提交确认 / Submit confirm
  const submitConfirm = async () => {
    if (!currentTask.value) return

    confirming.value = true
    try {
      const res = (await confirmCadFile({
        id: currentTask.value.id,
        fileName: confirmForm.fileName || undefined
      })) as unknown as BaseResponse<CadGenerationResponse>

      if (res.code === 200) {
        ElMessage.success('CAD文件已保存到项目')
        confirmDialogVisible.value = false
        loadTasks()
      } else {
        ElMessage.error(res.msg || '保存失败')
      }
    } catch {
      ElMessage.error('保存失败')
    } finally {
      confirming.value = false
    }
  }

  // 查看CAD / Preview CAD
  const handlePreview = (row: CadGenerationResponse) => {
    if (row.resultFileId) {
      router.push({
        path: '/designer/designer-assistant/cad-viewer/CadUpload',
        query: { id: row.resultFileId, projectId: row.projectId }
      })
    }
  }

  // 重试 / Retry
  const handleRetry = async (row: CadGenerationResponse) => {
    try {
      await ElMessageBox.confirm('确定要重试此任务吗？', '提示', { type: 'warning' })
      const res = (await retryCadGeneration(
        row.id
      )) as unknown as BaseResponse<CadGenerationResponse>
      if (res.code === 200) {
        ElMessage.success('已重新开始生成')
        loadTasks()
      } else {
        ElMessage.error(res.msg || '重试失败')
      }
    } catch {
      // 用户取消
    }
  }

  // 删除 / Delete
  const handleDelete = async (row: CadGenerationResponse) => {
    try {
      await ElMessageBox.confirm('确定要删除此任务吗？', '提示', { type: 'warning' })
      const res = (await deleteCadGeneration(row.id)) as unknown as BaseResponse<null>
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadTasks()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch {
      // 用户取消
    }
  }

  // 监听路由参数 / Watch route params
  watch(
    () => route.query.projectId,
    (newVal) => {
      if (newVal) {
        projectId.value = Number(newVal)
        filterForm.projectId = projectId.value
        loadTasks()
      }
    }
  )

  // 轮询定时器 / Polling timer
  let pollingTimer: ReturnType<typeof setInterval> | null = null

  // 开始轮询处理中的任务 / Start polling for processing tasks
  const startPolling = () => {
    if (pollingTimer) return

    pollingTimer = setInterval(() => {
      // 检查是否有处理中的任务 / Check if there are processing tasks
      const hasProcessing = taskList.value.some((t) => t.status === 'processing')
      if (hasProcessing) {
        loadTasks()
      } else {
        stopPolling()
      }
    }, 3000) // 每3秒轮询一次 / Poll every 3 seconds
  }

  // 停止轮询 / Stop polling
  const stopPolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  // 监听任务列表变化，自动开始/停止轮询 / Watch task list changes
  watch(
    () => taskList.value,
    (newList) => {
      const hasProcessing = newList.some((t) => t.status === 'processing' || t.status === 'pending')
      if (hasProcessing) {
        startPolling()
      } else {
        stopPolling()
      }
    },
    { deep: true }
  )

  /*** Lifecycle ***/
  onMounted(() => {
    loadTasks()
    loadProjects()
  })

  // 组件卸载时停止轮询 / Stop polling on unmount
  onUnmounted(() => {
    stopPolling()
  })
</script>

<style scoped lang="scss">
  .cad-generation-page {
    padding: 16px;

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      .left {
        display: flex;
        align-items: center;
        gap: 16px;

        h2 {
          margin: 0;
          font-size: 20px;
          font-weight: 500;
        }
      }
    }

    .filter-card {
      margin-bottom: 16px;
    }

    .status-cell {
      display: flex;
      align-items: center;
    }

    .text-secondary {
      color: var(--el-text-color-secondary);
    }

    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      margin-top: 16px;
    }
  }
</style>
