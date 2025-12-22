<template>
  <div class="version-list-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 / Table Header -->
      <ArtTableHeader :loading="loading" @refresh="loadVersions">
        <template #left>
          <ElSpace wrap>
            <ElButton link @click="handleBack">
              <ElIcon><ArrowLeft /></ElIcon>
              返回项目
            </ElButton>
            <span class="page-title">设计版本管理</span>
          </ElSpace>
        </template>
        <template #right>
          <ElButton type="primary" @click="handleCreate" v-ripple>创建新版本</ElButton>
          <ElButton @click="handleCompare" :disabled="selectedVersions.length !== 2" v-ripple>
            对比选中版本
          </ElButton>
        </template>
      </ArtTableHeader>

      <!-- 版本列表 / Version List -->
      <ArtTable
        :loading="loading"
        :data="versions"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>

    <!-- 创建/编辑版本对话框 / Create/Edit Version Dialog -->
    <ElDialog
      v-model="formDialogVisible"
      :title="isEdit ? '编辑版本' : '创建新版本'"
      width="700px"
      destroy-on-close
    >
      <ElForm ref="formRef" :model="versionForm" :rules="formRules" label-width="100px">
        <ElFormItem label="版本名称" prop="versionName">
          <ElInput v-model="versionForm.versionName" placeholder="如：初稿、修改版1、最终版" />
        </ElFormItem>
        <ElFormItem label="版本说明" prop="description">
          <ElInput
            v-model="versionForm.description"
            type="textarea"
            :rows="3"
            placeholder="描述本版本的主要变更内容"
          />
        </ElFormItem>
        <ElFormItem label="设计图" prop="designImages">
          <ElUpload
            v-model:file-list="designImageList"
            :action="uploadUrl"
            :headers="uploadHeaders"
            list-type="picture-card"
            :on-success="handleImageUploadSuccess"
            :on-remove="handleImageRemove"
            accept=".jpg,.jpeg,.png,.webp"
            multiple
          >
            <ElIcon><Plus /></ElIcon>
          </ElUpload>
          <div class="upload-tip">支持上传效果图、平面图等设计图片</div>
        </ElFormItem>
        <ElFormItem label="关联CAD" prop="cadFileIds">
          <ElSelect
            v-model="versionForm.cadFileIds"
            placeholder="选择关联的CAD文件"
            multiple
            style="width: 100%"
            :loading="cadFilesLoading"
          >
            <ElOption
              v-for="cad in projectCadFiles"
              :key="cad.id"
              :label="cad.fileName"
              :value="cad.id"
            >
              <div class="cad-option">
                <span>{{ cad.fileName }}</span>
                <ElTag v-if="cad.parseStatus === 'completed'" type="success" size="small"
                  >已解析</ElTag
                >
                <ElTag v-else-if="cad.parseStatus === 'processing'" type="warning" size="small"
                  >解析中</ElTag
                >
                <ElTag v-else-if="cad.parseStatus === 'failed'" type="danger" size="small"
                  >解析失败</ElTag
                >
              </div>
            </ElOption>
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="版本状态" prop="status">
          <ElSelect v-model="versionForm.status" style="width: 100%">
            <ElOption label="草稿" value="draft" />
            <ElOption label="已提交" value="submitted" />
            <ElOption label="已批准" value="approved" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="formDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="submitForm" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </ElButton>
      </template>
    </ElDialog>

    <!-- 版本详情对话框 / Version Detail Dialog -->
    <ElDialog v-model="detailDialogVisible" title="版本详情" width="800px">
      <VersionDetail v-if="currentVersion" :version="currentVersion" />
    </ElDialog>

    <!-- 历史对比记录 / Compare History -->
    <ElCard shadow="never" class="compare-history-card">
      <template #header>
        <div class="card-header">
          <span class="title">历史对比记录</span>
          <ElButton
            link
            type="primary"
            @click="loadCompareHistory"
            :loading="compareHistoryLoading"
          >
            <ElIcon><Refresh /></ElIcon>
            刷新
          </ElButton>
        </div>
      </template>

      <ElTable v-loading="compareHistoryLoading" :data="compareHistory" stripe>
        <ElTableColumn label="对比版本" min-width="200">
          <template #default="{ row }">
            <div class="compare-versions">
              <ElTag type="primary"
                >V{{ row.versionA?.versionNumber }} {{ row.versionA?.versionName }}</ElTag
              >
              <ElIcon class="vs-icon"><Right /></ElIcon>
              <ElTag type="success"
                >V{{ row.versionB?.versionNumber }} {{ row.versionB?.versionName }}</ElTag
              >
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="summary" label="对比摘要" min-width="250" show-overflow-tooltip />
        <ElTableColumn prop="compareStatus" label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getCompareStatusType(row.compareStatus)">
              {{ getCompareStatusText(row.compareStatus) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="createdAt" label="对比时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton link type="primary" @click="handleViewCompare(row)">查看详情</ElButton>
            <ElButton link type="danger" @click="handleDeleteCompare(row)">删除</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <ElEmpty
        v-if="!compareHistoryLoading && compareHistory.length === 0"
        description="暂无对比记录"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Version List Component
   * 设计版本列表组件（重构版）
   * Requirements: 7.1, 7.2
   ***/
  import { ref, reactive, computed, onMounted, watch, h } from 'vue'

  // 定义响应类型 / Define response type
  interface BaseResponse<T = unknown> {
    code: number
    msg: string
    data: T
  }
  import { useRouter, useRoute } from 'vue-router'
  import { useUserStore } from '@/store/modules/user'
  import {
    ElMessage,
    ElMessageBox,
    ElTag,
    ElImage,
    type FormInstance,
    type FormRules,
    type UploadUserFile
  } from 'element-plus'
  import { ArrowLeft, Plus, Refresh, Right } from '@element-plus/icons-vue'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import dayjs from 'dayjs'
  import {
    getDesignVersionList,
    createDesignVersion,
    updateDesignVersion,
    deleteDesignVersion,
    getVersionCompareList,
    deleteVersionCompare,
    type DesignVersionResponse,
    type VersionCompareResponse
  } from '@/api/designer-version'
  import { getCadFileList, type CadFileResponse } from '@/api/designer-cad'
  import VersionDetail from './VersionDetail.vue'
  import type { ColumnOption } from '@/types/component'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  /*** State ***/
  const loading = ref(false)
  const submitting = ref(false)
  const formDialogVisible = ref(false)
  const detailDialogVisible = ref(false)
  const isEdit = ref(false)
  const versions = ref<DesignVersionResponse[]>([])
  const selectedVersions = ref<DesignVersionResponse[]>([])
  const currentVersion = ref<DesignVersionResponse | null>(null)
  const projectCadFiles = ref<CadFileResponse[]>([])
  const cadFilesLoading = ref(false)
  const designImageList = ref<UploadUserFile[]>([])
  const formRef = ref<FormInstance>()
  const compareHistory = ref<VersionCompareResponse[]>([])
  const compareHistoryLoading = ref(false)

  /*** Upload configuration ***/
  const uploadUrl = computed(() => {
    const baseUrl = import.meta.env.VITE_API_URL || ''
    return `${baseUrl}/api/file/upload`
  })

  const uploadHeaders = computed(() => ({
    Authorization: `Bearer ${userStore.accessToken}`
  }))

  const pagination = reactive({
    current: 1,
    size: 10,
    total: 0
  })

  const projectId = ref<number>(0)

  // 状态配置 / Status config
  const STATUS_CONFIG = {
    draft: { type: 'info' as const, text: '草稿' },
    submitted: { type: 'warning' as const, text: '已提交' },
    approved: { type: 'success' as const, text: '已批准' },
    rejected: { type: 'danger' as const, text: '已拒绝' }
  } as const

  // 格式化日期 / Format date
  const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  // 表格列配置 / Table columns config
  const columns = computed<ColumnOption[]>(() => [
    { type: 'selection', width: 55 },
    {
      prop: 'versionNumber',
      label: '版本号',
      width: 100,
      formatter: (row: DesignVersionResponse) =>
        h(ElTag, { type: 'primary' }, () => `V${row.versionNumber}`)
    },
    { prop: 'versionName', label: '版本名称', minWidth: 150 },
    {
      prop: 'designImages',
      label: '设计图',
      width: 120,
      formatter: (row: DesignVersionResponse) => {
        if (row.designImages?.length) {
          return h('div', { class: 'preview-images' }, [
            ...row.designImages.slice(0, 2).map((img: any, idx: number) =>
              h(ElImage, {
                key: idx,
                src: img.url,
                previewSrcList: row.designImages!.map((i: any) => i.url),
                fit: 'cover',
                class: 'preview-thumb'
              })
            ),
            row.designImages.length > 2
              ? h('span', { class: 'more-count' }, `+${row.designImages.length - 2}`)
              : null
          ])
        }
        return h('span', { class: 'text-secondary' }, '-')
      }
    },
    {
      prop: 'cadFileIds',
      label: 'CAD文件',
      width: 100,
      formatter: (row: DesignVersionResponse) =>
        row.cadFileIds?.length
          ? h(ElTag, { type: 'info' }, () => `${row.cadFileIds!.length} 个`)
          : h('span', { class: 'text-secondary' }, '-')
    },
    { prop: 'description', label: '版本说明', minWidth: 200, showOverflowTooltip: true },
    {
      prop: 'status',
      label: '状态',
      width: 100,
      formatter: (row: DesignVersionResponse) => {
        const config = STATUS_CONFIG[row.status as keyof typeof STATUS_CONFIG] || {
          type: 'info' as const,
          text: row.status
        }
        return h(ElTag, { type: config.type }, () => config.text)
      }
    },
    {
      prop: 'createdAt',
      label: '创建时间',
      width: 170,
      formatter: (row: DesignVersionResponse) => formatDate(row.createdAt)
    },
    {
      prop: 'operation',
      label: '操作',
      width: 180,
      fixed: 'right',
      formatter: (row: DesignVersionResponse) =>
        h('div', { class: 'flex gap-1' }, [
          h(ArtButtonTable, { type: 'view', onClick: () => handleView(row) }),
          h(ArtButtonTable, { type: 'edit', onClick: () => handleEdit(row) }),
          h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
        ])
    }
  ])

  const versionForm = reactive({
    id: 0,
    versionName: '',
    description: '',
    designImages: [] as string[],
    cadFileIds: [] as number[],
    status: 'draft'
  })

  const formRules: FormRules = {
    versionName: [{ required: true, message: '请输入版本名称', trigger: 'blur' }]
  }

  /*** Load version list ***/
  const loadVersions = async () => {
    if (!projectId.value) return

    loading.value = true
    try {
      const res = (await getDesignVersionList({
        projectId: projectId.value,
        current: pagination.current,
        size: pagination.size
      })) as unknown as BaseResponse<{ records: DesignVersionResponse[]; total: number }>

      if (res.code === 200 && res.data) {
        versions.value = res.data.records || []
        pagination.total = res.data.total || 0
      } else {
        ElMessage.error(res.msg || '获取版本列表失败')
      }
    } catch (error) {
      console.error('获取版本列表失败:', error)
      ElMessage.error('获取版本列表失败')
    } finally {
      loading.value = false
    }
  }

  /*** Load project CAD files ***/
  const loadProjectCadFiles = async () => {
    if (!projectId.value) return

    cadFilesLoading.value = true
    try {
      const res = (await getCadFileList({
        projectId: projectId.value,
        current: 1,
        size: 100 // 获取所有CAD文件用于选择
      })) as unknown as BaseResponse<{ records: CadFileResponse[]; total: number }>

      if (res.code === 200 && res.data) {
        projectCadFiles.value = res.data.records || []
      } else {
        console.error('获取CAD文件列表失败:', res.msg)
      }
    } catch (error) {
      console.error('获取CAD文件列表失败:', error)
    } finally {
      cadFilesLoading.value = false
    }
  }

  /*** Load compare history ***/
  const loadCompareHistory = async () => {
    if (!projectId.value) return

    compareHistoryLoading.value = true
    try {
      const res = (await getVersionCompareList({
        projectId: projectId.value,
        current: 1,
        size: 20
      })) as unknown as BaseResponse<{ records: VersionCompareResponse[]; total: number }>

      if (res.code === 200 && res.data) {
        compareHistory.value = res.data.records || []
      } else {
        console.error('获取对比历史失败:', res.msg)
      }
    } catch (error) {
      console.error('获取对比历史失败:', error)
    } finally {
      compareHistoryLoading.value = false
    }
  }

  /*** Handle view compare detail ***/
  const handleViewCompare = (row: VersionCompareResponse) => {
    router.push({
      path: '/designer/designer-assistant/version-compare/VersionDiff',
      query: {
        compareId: row.id
      }
    })
  }

  /*** Handle delete compare record ***/
  const handleDeleteCompare = async (row: VersionCompareResponse) => {
    try {
      await ElMessageBox.confirm('确定要删除这条对比记录吗？此操作不可恢复。', '删除确认', {
        type: 'warning'
      })

      const res = (await deleteVersionCompare(row.id)) as unknown as BaseResponse<null>
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadCompareHistory()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除对比记录失败:', error)
        ElMessage.error('删除失败')
      }
    }
  }

  /*** Get compare status type ***/
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getCompareStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      pending: 'warning',
      processing: 'primary',
      completed: 'success',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  /*** Get compare status text ***/
  const getCompareStatusText = (status: string) => {
    const map: Record<string, string> = {
      pending: '待分析',
      processing: '分析中',
      completed: '已完成',
      failed: '分析失败'
    }
    return map[status] || status
  }

  /*** Handle selection change ***/
  const handleSelectionChange = (selection: DesignVersionResponse[]) => {
    selectedVersions.value = selection
  }

  /*** Handle compare ***/
  const handleCompare = () => {
    if (selectedVersions.value.length !== 2) {
      ElMessage.warning('请选择两个版本进行对比')
      return
    }
    const [versionA, versionB] = selectedVersions.value
    router.push({
      path: '/designer/designer-assistant/version-compare/VersionDiff',
      query: {
        projectId: projectId.value,
        versionAId: versionA.id,
        versionBId: versionB.id
      }
    })
  }

  /*** Handle create ***/
  const handleCreate = () => {
    isEdit.value = false
    versionForm.id = 0
    versionForm.versionName = ''
    versionForm.description = ''
    versionForm.designImages = []
    versionForm.cadFileIds = []
    versionForm.status = 'draft'
    designImageList.value = []
    formDialogVisible.value = true
  }

  /*** Handle edit ***/
  const handleEdit = (row: DesignVersionResponse) => {
    isEdit.value = true
    versionForm.id = row.id
    versionForm.versionName = row.versionName
    versionForm.description = row.description
    versionForm.designImages = row.designImages?.map((img) => img.url) || []
    versionForm.cadFileIds = row.cadFileIds || []
    versionForm.status = row.status
    // 设置图片列表
    designImageList.value =
      row.designImages?.map((img, idx) => ({
        name: img.name || `image-${idx}`,
        url: img.url
      })) || []
    formDialogVisible.value = true
  }

  /*** Handle view version detail ***/
  const handleView = (row: DesignVersionResponse) => {
    currentVersion.value = row
    detailDialogVisible.value = true
  }

  /*** Handle delete version ***/
  const handleDelete = async (row: DesignVersionResponse) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除版本 "${row.versionName}" 吗？此操作不可恢复。`,
        '删除确认',
        { type: 'warning' }
      )

      const res = (await deleteDesignVersion(row.id)) as unknown as BaseResponse<null>
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadVersions()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除版本失败:', error)
        ElMessage.error('删除失败')
      }
    }
  }

  /*** Handle image upload success ***/
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const handleImageUploadSuccess = (response: any, _file: UploadUserFile) => {
    if (response.code === 200 && response.data?.url) {
      versionForm.designImages.push(response.data.url)
    }
  }

  /*** Handle image remove ***/
  const handleImageRemove = (file: UploadUserFile) => {
    const url = file.url || (file.response as any)?.data?.url
    if (url) {
      const idx = versionForm.designImages.indexOf(url)
      if (idx > -1) {
        versionForm.designImages.splice(idx, 1)
      }
    }
  }

  /*** Submit form ***/
  const submitForm = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (!valid) return

      submitting.value = true
      try {
        const data = {
          projectId: projectId.value,
          versionName: versionForm.versionName,
          description: versionForm.description,
          designImages: versionForm.designImages,
          cadFileIds: versionForm.cadFileIds,
          status: versionForm.status
        }

        let res: BaseResponse<DesignVersionResponse>
        if (isEdit.value) {
          res = (await updateDesignVersion({
            id: versionForm.id,
            ...data
          })) as unknown as BaseResponse<DesignVersionResponse>
        } else {
          res = (await createDesignVersion(data)) as unknown as BaseResponse<DesignVersionResponse>
        }

        if (res.code === 200) {
          ElMessage.success(isEdit.value ? '保存成功' : '创建成功')
          formDialogVisible.value = false
          loadVersions()
        } else {
          ElMessage.error(res.msg || (isEdit.value ? '保存失败' : '创建失败'))
        }
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('操作失败')
      } finally {
        submitting.value = false
      }
    })
  }

  /*** Handle back ***/
  const handleBack = () => {
    router.push(`/designer/designer-assistant/project/ProjectDetail/${projectId.value}`)
  }

  /*** Handle page size change ***/
  const handleSizeChange = (size: number) => {
    pagination.size = size
    loadVersions()
  }

  /*** Handle current page change ***/
  const handleCurrentChange = (current: number) => {
    pagination.current = current
    loadVersions()
  }

  /*** Format date ***/
  // 已在上方定义

  // 监听路由参数 / Watch route params
  watch(
    () => route.query.projectId,
    (newVal) => {
      if (newVal) {
        projectId.value = Number(newVal)
        loadVersions()
        loadProjectCadFiles()
        loadCompareHistory()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (route.query.projectId) {
      projectId.value = Number(route.query.projectId)
      loadVersions()
      loadProjectCadFiles()
      loadCompareHistory()
    }
  })
</script>

<style scoped lang="scss">
  .version-list-page {
    // 改为自动高度，允许内容撑开 / Change to auto height, allow content to expand
    height: auto !important;
    min-height: var(--art-full-height);

    .page-title {
      font-size: 16px;
      font-weight: 500;
    }

    .preview-images {
      display: flex;
      align-items: center;
      gap: 4px;

      .preview-thumb {
        width: 40px;
        height: 40px;
        border-radius: 4px;
      }

      .more-count {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }

    .text-secondary {
      color: var(--el-text-color-secondary);
    }

    .upload-tip {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 8px;
    }

    .cad-option {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;
    }

    // 版本列表卡片不使用 flex:1，改为自适应高度 / Version list card uses auto height instead of flex:1
    :deep(.art-table-card) {
      flex: none;
      height: auto;
      min-height: 300px;

      .el-card__body {
        height: auto;
        overflow: visible;
      }

      .art-table {
        height: auto !important;

        .el-table {
          height: auto !important;
        }
      }
    }
  }

  .compare-history-card {
    margin-top: 16px;
    flex: none;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .compare-versions {
      display: flex;
      align-items: center;
      gap: 8px;

      .vs-icon {
        font-size: 16px;
        color: var(--el-text-color-secondary);
      }
    }
  }
</style>
