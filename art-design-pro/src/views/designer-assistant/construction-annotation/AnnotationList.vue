<template>
  <div class="annotation-list-page">
    <!-- 页面标题 / Page title -->
    <div class="page-header">
      <div class="left">
        <ElButton link @click="handleBack" v-if="projectId">
          <ElIcon><ArrowLeft /></ElIcon>
          返回项目
        </ElButton>
        <h2>施工图标注管理</h2>
      </div>
      <ElButton type="primary" @click="handleCreate">
        <ElIcon><Plus /></ElIcon>
        新建标注
      </ElButton>
    </div>

    <!-- 筛选条件 / Filter conditions -->
    <ElCard class="filter-card" shadow="never">
      <ElForm :inline="true" :model="filterForm">
        <ElFormItem label="项目" v-if="!projectId">
          <ElSelect
            v-model="filterForm.projectId"
            placeholder="选择项目"
            clearable
            @change="handleProjectChange"
          >
            <ElOption
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="设计版本">
          <ElSelect v-model="filterForm.versionId" placeholder="选择版本" clearable>
            <ElOption
              v-for="ver in projectVersions"
              :key="ver.id"
              :label="`V${ver.versionNumber} - ${ver.versionName}`"
              :value="ver.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSelect v-model="filterForm.status" placeholder="选择状态" clearable>
            <ElOption label="待分析" value="pending" />
            <ElOption label="分析中" value="processing" />
            <ElOption label="已完成" value="completed" />
            <ElOption label="失败" value="failed" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch">
            <ElIcon><Search /></ElIcon>
            搜索
          </ElButton>
          <ElButton @click="handleReset">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 标注列表 / Annotation list -->
    <ElCard class="list-card" shadow="never">
      <ElTable :data="annotationList" v-loading="loading" stripe>
        <ElTableColumn prop="id" label="ID" width="80" />
        <ElTableColumn label="项目" min-width="120" v-if="!projectId">
          <template #default="{ row }">
            {{ row.project?.name || '-' }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="设计版本" width="120">
          <template #default="{ row }">
            <ElTag v-if="row.version" type="primary"> V{{ row.version.versionNumber }} </ElTag>
            <span v-else class="text-secondary">-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="CAD文件" min-width="150">
          <template #default="{ row }">
            {{ row.cadFile?.fileName || '-' }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.analysisStatus)">
              {{ getStatusText(row.analysisStatus) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="标注数量" width="100">
          <template #default="{ row }">
            {{ row.annotations?.length || 0 }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" size="small" @click="handleEdit(row)">编辑</ElButton>
            <ElButton size="small" @click="handleExport(row)">导出</ElButton>
            <ElButton type="danger" size="small" @click="handleDelete(row)">删除</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <!-- 分页 / Pagination -->
      <div class="pagination-wrapper">
        <ElPagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </ElCard>

    <!-- 新建标注对话框 / Create annotation dialog -->
    <ElDialog v-model="createDialogVisible" title="新建施工图标注" width="550px">
      <ElForm :model="createForm" label-width="100px" :rules="createRules" ref="createFormRef">
        <ElFormItem label="项目" prop="projectId" v-if="!projectId">
          <ElSelect
            v-model="createForm.projectId"
            placeholder="选择项目"
            @change="handleCreateProjectChange"
            style="width: 100%"
          >
            <ElOption
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="设计版本" prop="versionId">
          <ElSelect
            v-model="createForm.versionId"
            placeholder="选择设计版本"
            @change="handleVersionChange"
            style="width: 100%"
          >
            <ElOption
              v-for="ver in createVersions"
              :key="ver.id"
              :label="`V${ver.versionNumber} - ${ver.versionName}`"
              :value="ver.id"
            />
          </ElSelect>
          <div class="form-tip">选择版本后可选择该版本关联的CAD文件</div>
        </ElFormItem>
        <ElFormItem label="CAD文件" prop="cadFileId">
          <ElSelect v-model="createForm.cadFileId" placeholder="选择CAD文件" style="width: 100%">
            <ElOption
              v-for="file in versionCadFiles"
              :key="file.id"
              :label="file.fileName"
              :value="file.id"
            />
          </ElSelect>
          <div class="form-tip" v-if="versionCadFiles.length === 0 && createForm.versionId">
            暂无CAD文件，请先在项目中上传CAD文件
          </div>
        </ElFormItem>
        <ElFormItem label="施工图" prop="imagePath">
          <div class="image-source-section">
            <div v-if="createForm.imagePath" class="uploaded-image">
              <ElImage
                :src="createForm.imagePath"
                fit="contain"
                style="max-width: 300px; max-height: 200px"
              />
              <div class="image-actions">
                <ElButton size="small" @click="createForm.imagePath = ''">清除</ElButton>
              </div>
            </div>
            <div v-else class="upload-section">
              <ElUpload
                class="upload-demo"
                :http-request="customUpload"
                :before-upload="beforeUpload"
                accept=".jpg,.jpeg,.png,.bmp"
                :show-file-list="false"
              >
                <ElButton type="primary">上传施工图</ElButton>
              </ElUpload>
              <div class="el-upload__tip">
                请上传从CAD软件导出的施工图图片（JPG、PNG、BMP格式）
              </div>
            </div>
          </div>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="createDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="confirmCreate" :loading="creating">创建并分析</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * Annotation List Component
   * 施工图标注列表组件（重构版 - 关联设计版本）
   * Requirements: 9.1, 9.2
   ***/
  import { ref, reactive, onMounted, watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
  import { Plus, Search, ArrowLeft } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import {
    type ConstructionAnnotationResponse,
    getAnnotationsList,
    deleteAnnotation,
    analyzeConstructionDrawing
  } from '@/api/designer-construction-annotation'
  import { getDesignerProjects } from '@/api/designer-project'
  import { getDesignVersionList, type DesignVersionResponse } from '@/api/designer-version'
  import { getCadFileList, type CadFileResponse } from '@/api/designer-cad'
  import { uploadFile } from '@/api/file'

  /*** Router ***/
  const router = useRouter()
  const route = useRoute()

  /*** Type definitions for API responses ***/
  interface BaseResponse<T = unknown> {
    code: number
    msg?: string
    data: T
  }

  /*** State ***/
  const loading = ref(false)
  const creating = ref(false)
  const createDialogVisible = ref(false)
  const annotationList = ref<ConstructionAnnotationResponse[]>([])
  const projects = ref<{ id: number; name: string }[]>([])
  const projectVersions = ref<{ id: number; versionNumber: number; versionName: string }[]>([])
  const createVersions = ref<{ id: number; versionNumber: number; versionName: string }[]>([])
  const versionCadFiles = ref<{ id: number; fileName: string }[]>([])
  const createFormRef = ref<FormInstance>()

  // 从路由获取项目ID / Get project ID from route
  const projectId = ref<number | undefined>(
    route.query.projectId ? Number(route.query.projectId) : undefined
  )

  const filterForm = reactive({
    projectId: projectId.value,
    versionId: undefined as number | undefined,
    status: ''
  })

  const pagination = reactive({
    page: 1,
    pageSize: 10,
    total: 0
  })

  const createForm = reactive({
    projectId: projectId.value,
    versionId: undefined as number | undefined,
    cadFileId: undefined as number | undefined,
    imagePath: ''
  })

  const createRules: FormRules = {
    projectId: [{ required: true, message: '请选择项目', trigger: 'change' }],
    versionId: [{ required: true, message: '请选择设计版本', trigger: 'change' }],
    cadFileId: [{ required: true, message: '请选择CAD文件', trigger: 'change' }],
    imagePath: [{ required: true, message: '请上传施工图', trigger: 'change' }]
  }

  /*** Methods ***/
  // 加载标注列表 / Load annotation list
  const loadAnnotations = async () => {
    loading.value = true
    try {
      const res = (await getAnnotationsList({
        projectId: filterForm.projectId,
        versionId: filterForm.versionId,
        status: filterForm.status || undefined,
        page: pagination.page,
        pageSize: pagination.pageSize
      })) as unknown as BaseResponse<{
        records: ConstructionAnnotationResponse[]
        total: number
      }>
      if (res && res.data) {
        annotationList.value = res.data.records || []
        pagination.total = res.data.total || 0
      }
    } catch {
      ElMessage.error('加载标注列表失败')
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
      })) as unknown as BaseResponse<{ records: { id: number; name: string }[] }>
      if (res.code === 200 && res.data) {
        projects.value = res.data.records || []
      }
    } catch {
      console.error('加载项目列表失败')
    }
  }

  // 加载项目版本 / Load project versions
  const loadProjectVersions = async (pId: number) => {
    try {
      const res = (await getDesignVersionList({
        projectId: pId,
        current: 1,
        size: 100
      })) as unknown as BaseResponse<{ records: DesignVersionResponse[] }>
      if (res.code === 200 && res.data) {
        return (
          res.data.records?.map((v) => ({
            id: v.id,
            versionNumber: v.versionNumber,
            versionName: v.versionName
          })) || []
        )
      }
    } catch {
      console.error('加载版本列表失败')
    }
    return []
  }

  // 筛选项目变更 / Filter project change
  const handleProjectChange = async (pId: number) => {
    filterForm.versionId = undefined
    if (pId) {
      projectVersions.value = await loadProjectVersions(pId)
    } else {
      projectVersions.value = []
    }
  }

  // 创建表单项目变更 / Create form project change
  const handleCreateProjectChange = async (pId: number) => {
    createForm.versionId = undefined
    createForm.cadFileId = undefined
    versionCadFiles.value = []
    if (pId) {
      createVersions.value = await loadProjectVersions(pId)
    } else {
      createVersions.value = []
    }
  }

  // 版本变更时加载CAD文件 / Load CAD files when version changes
  const handleVersionChange = async (versionId: number) => {
    createForm.cadFileId = undefined
    createForm.imagePath = ''
    if (!versionId) {
      versionCadFiles.value = []
      return
    }
    // 获取项目的CAD文件列表 / Get project CAD file list
    const pId = createForm.projectId || projectId.value
    if (!pId) {
      versionCadFiles.value = []
      return
    }
    try {
      const res = (await getCadFileList({
        projectId: pId,
        current: 1,
        size: 100,
        parseStatus: 'completed' // 只获取已解析的CAD文件 / Only get parsed CAD files
      })) as unknown as BaseResponse<{ records: CadFileResponse[] }>
      if (res.code === 200 && res.data?.records) {
        versionCadFiles.value = res.data.records.map((f) => ({
          id: f.id,
          fileName: f.fileName
        }))
      } else {
        versionCadFiles.value = []
      }
    } catch {
      console.error('加载CAD文件列表失败')
      versionCadFiles.value = []
    }
  }

  // 获取状态类型 / Get status type
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      completed: 'success',
      processing: 'warning',
      failed: 'danger',
      pending: 'info'
    }
    return map[status] || 'info'
  }

  // 获取状态文本 / Get status text
  const getStatusText = (status: string): string => {
    const map: Record<string, string> = {
      pending: '待分析',
      processing: '分析中',
      completed: '已完成',
      failed: '失败'
    }
    return map[status] || status
  }

  // 格式化日期 / Format date
  const formatDate = (dateStr: string): string => {
    if (!dateStr) return '-'
    return dayjs(dateStr).format('YYYY-MM-DD HH:mm')
  }

  // 搜索 / Search
  const handleSearch = () => {
    pagination.page = 1
    loadAnnotations()
  }

  // 重置 / Reset
  const handleReset = () => {
    filterForm.projectId = projectId.value
    filterForm.versionId = undefined
    filterForm.status = ''
    pagination.page = 1
    loadAnnotations()
  }

  // 分页大小变更 / Page size change
  const handleSizeChange = () => {
    pagination.page = 1
    loadAnnotations()
  }

  // 页码变更 / Page change
  const handlePageChange = () => {
    loadAnnotations()
  }

  // 返回项目 / Back to project
  const handleBack = () => {
    if (projectId.value) {
      router.push(`/designer/project/${projectId.value}`)
    }
  }

  // 新建标注 / Create annotation
  const handleCreate = async () => {
    createForm.projectId = projectId.value
    createForm.versionId = undefined
    createForm.cadFileId = undefined
    createForm.imagePath = ''
    versionCadFiles.value = []
    if (projectId.value) {
      createVersions.value = await loadProjectVersions(projectId.value)
    }
    createDialogVisible.value = true
  }

  // 上传前验证 / Before upload validation
  const beforeUpload = (file: File) => {
    const isImage = ['image/jpeg', 'image/png', 'image/bmp'].includes(file.type)
    const isLt10M = file.size / 1024 / 1024 < 10

    if (!isImage) {
      ElMessage.error('只能上传 JPG/PNG/BMP 格式的图片')
      return false
    }
    if (!isLt10M) {
      ElMessage.error('图片大小不能超过 10MB')
      return false
    }
    return true
  }

  // 自定义上传方法 / Custom upload method
  const customUpload = async (options: { file: File }) => {
    try {
      const res = await uploadFile(options.file, 'image')
      if (res && res.url) {
        createForm.imagePath = res.url
        ElMessage.success('上传成功')
      } else {
        ElMessage.error('上传失败')
      }
    } catch {
      ElMessage.error('上传失败')
    }
  }

  // 确认创建 / Confirm create
  const confirmCreate = async () => {
    if (!createFormRef.value) return

    await createFormRef.value.validate(async (valid) => {
      if (!valid) return

      creating.value = true
      try {
        const res = (await analyzeConstructionDrawing({
          projectId: createForm.projectId!,
          versionId: createForm.versionId,
          cadFileId: createForm.cadFileId!,
          imagePath: createForm.imagePath
        })) as unknown as BaseResponse<ConstructionAnnotationResponse>
        if (res.code === 200) {
          ElMessage.success('创建成功，正在分析中')
          createDialogVisible.value = false
          loadAnnotations()
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

  // 编辑标注 / Edit annotation
  const handleEdit = (row: ConstructionAnnotationResponse) => {
    router.push({
      path: '/designer/designer-assistant/construction-annotation/AnnotationEditor',
      query: {
        annotationId: row.id,
        projectId: row.projectId,
        versionId: row.versionId,
        cadFileId: row.cadFileId
      }
    })
  }

  // 导出标注 / Export annotation
  const handleExport = (row: ConstructionAnnotationResponse) => {
    router.push({
      path: '/designer/designer-assistant/construction-annotation/AnnotationEditor',
      query: {
        annotationId: row.id,
        projectId: row.projectId,
        versionId: row.versionId,
        cadFileId: row.cadFileId,
        action: 'export'
      }
    })
  }

  // 删除标注 / Delete annotation
  const handleDelete = async (row: ConstructionAnnotationResponse) => {
    try {
      await ElMessageBox.confirm('确定要删除这个标注吗？', '提示', {
        type: 'warning'
      })

      const res = (await deleteAnnotation(row.id)) as unknown as BaseResponse<null>
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadAnnotations()
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
    async (newVal) => {
      if (newVal) {
        projectId.value = Number(newVal)
        filterForm.projectId = projectId.value
        projectVersions.value = await loadProjectVersions(projectId.value)
        loadAnnotations()
      }
    }
  )

  /*** Lifecycle Hooks ***/
  onMounted(async () => {
    loadProjects()
    if (projectId.value) {
      projectVersions.value = await loadProjectVersions(projectId.value)
    }
    loadAnnotations()
  })
</script>

<style scoped lang="scss">
  .annotation-list-page {
    padding: 16px;
  }

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

  .list-card {
    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      margin-top: 16px;
    }
  }

  .text-secondary {
    color: var(--el-text-color-secondary);
  }

  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 4px;
  }

  .uploaded-image {
    margin-top: 8px;
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    padding: 8px;

    .image-actions {
      margin-top: 8px;
      text-align: center;
    }
  }

  .image-source-section {
    width: 100%;
  }

  .upload-section {
    .upload-tip {
      margin-bottom: 8px;
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }

    .el-upload__tip {
      margin-top: 8px;
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
</style>
