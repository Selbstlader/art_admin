<template>
  <div class="version-list-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <div class="left">
            <ElButton link @click="handleBack">
              <ElIcon><ArrowLeft /></ElIcon>
              返回项目
            </ElButton>
            <span class="title">设计版本管理</span>
          </div>
          <div class="right">
            <ElButton type="primary" @click="handleCreate">
              <ElIcon><Plus /></ElIcon>
              创建新版本
            </ElButton>
            <ElButton @click="handleCompare" :disabled="selectedVersions.length !== 2">
              <ElIcon><Switch /></ElIcon>
              对比选中版本
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 版本列表 / Version List -->
      <ElTable
        v-loading="loading"
        :data="versions"
        stripe
        @selection-change="handleSelectionChange"
      >
        <ElTableColumn type="selection" width="55" />
        <ElTableColumn prop="versionNumber" label="版本号" width="100">
          <template #default="{ row }">
            <ElTag type="primary">V{{ row.versionNumber }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="versionName" label="版本名称" min-width="150" />
        <ElTableColumn label="设计图" width="120">
          <template #default="{ row }">
            <div class="preview-images" v-if="row.designImages?.length">
              <ElImage
                v-for="(img, idx) in row.designImages.slice(0, 2)"
                :key="idx"
                :src="img.url"
                :preview-src-list="row.designImages.map((i: any) => i.url)"
                fit="cover"
                class="preview-thumb"
              />
              <span v-if="row.designImages.length > 2" class="more-count">
                +{{ row.designImages.length - 2 }}
              </span>
            </div>
            <span v-else class="text-secondary">-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="CAD文件" width="100">
          <template #default="{ row }">
            <ElTag v-if="row.cadFileIds?.length" type="info">
              {{ row.cadFileIds.length }} 个
            </ElTag>
            <span v-else class="text-secondary">-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="description" label="版本说明" min-width="200" show-overflow-tooltip />
        <ElTableColumn prop="status" label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="createdAt" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <ElButton link type="primary" @click="handleView(row)">查看</ElButton>
            <ElButton link type="primary" @click="handleEdit(row)">编辑</ElButton>
            <ElButton link type="danger" @click="handleDelete(row)">删除</ElButton>
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
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
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
  </div>
</template>

<script setup lang="ts">
  /***
   * Version List Component
   * 设计版本列表组件（重构版）
   * Requirements: 7.1, 7.2
   ***/
  import { ref, reactive, computed, onMounted, watch } from 'vue'

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
    type FormInstance,
    type FormRules,
    type UploadUserFile
  } from 'element-plus'
  import { ArrowLeft, Switch, Plus } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import {
    getDesignVersionList,
    createDesignVersion,
    updateDesignVersion,
    deleteDesignVersion,
    type DesignVersionResponse
  } from '@/api/designer-version'
  import { getCadFileList, type CadFileResponse } from '@/api/designer-cad'
  import VersionDetail from './VersionDetail.vue'

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

  /*** Upload configuration ***/
  // 上传地址（带认证）
  const uploadUrl = computed(() => {
    const baseUrl = import.meta.env.VITE_API_URL || ''
    return `${baseUrl}/api/file/upload`
  })

  // 上传请求头（携带 token）
  const uploadHeaders = computed(() => ({
    Authorization: `Bearer ${userStore.accessToken}`
  }))

  const pagination = reactive({
    current: 1,
    size: 10,
    total: 0
  })

  // 获取项目ID / Get project ID
  const projectId = ref<number>(0)

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
  const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  /*** Get status type ***/
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      draft: 'info',
      submitted: 'warning',
      approved: 'success',
      rejected: 'danger'
    }
    return map[status] || 'info'
  }

  /*** Get status text ***/
  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      draft: '草稿',
      submitted: '已提交',
      approved: '已批准',
      rejected: '已拒绝'
    }
    return map[status] || status
  }

  // 监听路由参数 / Watch route params
  watch(
    () => route.query.projectId,
    (newVal) => {
      if (newVal) {
        projectId.value = Number(newVal)
        loadVersions()
        loadProjectCadFiles()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (route.query.projectId) {
      projectId.value = Number(route.query.projectId)
      loadVersions()
      loadProjectCadFiles()
    }
  })
</script>

<style scoped lang="scss">
  .version-list-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .left {
        display: flex;
        align-items: center;
        gap: 16px;

        .title {
          font-size: 16px;
          font-weight: 500;
        }
      }

      .right {
        display: flex;
        gap: 12px;
      }
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

    .pagination-wrapper {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
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
  }
</style>
