<template>
  <div class="project-page art-full-height">
    <ArtSearchBar
      v-show="showSearchBar"
      v-model="searchForm"
      :items="formItems"
      @search="handleSearch"
      @reset="handleReset"
    >
    </ArtSearchBar>

    <ElCard
      class="art-table-card"
      shadow="never"
      :style="{ 'margin-top': showSearchBar ? '12px' : '0' }"
    >
      <ArtTableHeader
        v-model:columns="columnChecks"
        v-model:showSearchBar="showSearchBar"
        :loading="loading"
        @refresh="refreshData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton @click="handleCreate" v-ripple>新建项目</ElButton>
            <ElButton @click="handleCreateFromTemplate" v-ripple>从模板创建</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 项目列表 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
      </ArtTable>
    </ElCard>

    <!-- 新建项目弹窗 -->
    <ElDialog
      v-model="createDialogVisible"
      title="新建项目"
      width="700px"
      :close-on-click-modal="false"
    >
      <ElForm ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <ElFormItem label="项目名称" prop="name">
          <ElInput v-model="createForm.name" placeholder="请输入项目名称" maxlength="200" />
        </ElFormItem>

        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="项目分类" prop="category">
              <ElInput v-model="createForm.category" placeholder="如：软件开发" maxlength="50" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="项目经理" prop="managerId">
              <ElSelect
                v-model="createForm.managerId"
                placeholder="请选择项目经理"
                filterable
                clearable
                style="width: 100%"
              >
                <ElOption
                  v-for="user in userList"
                  :key="user.id"
                  :label="user.nickName"
                  :value="user.id"
                />
              </ElSelect>
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="开始日期" prop="startDate">
              <ElDatePicker
                v-model="createForm.startDate"
                type="date"
                placeholder="选择开始日期"
                style="width: 100%"
              />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="结束日期" prop="endDate">
              <ElDatePicker
                v-model="createForm.endDate"
                type="date"
                placeholder="选择结束日期"
                style="width: 100%"
              />
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="项目状态" prop="status">
              <ElSelect v-model="createForm.status" placeholder="选择状态" style="width: 100%">
                <ElOption label="计划中" :value="1" />
                <ElOption label="进行中" :value="2" />
                <ElOption label="已完成" :value="3" />
                <ElOption label="已暂停" :value="4" />
                <ElOption label="已取消" :value="5" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="优先级" prop="priority">
              <ElSelect v-model="createForm.priority" placeholder="选择优先级" style="width: 100%">
                <ElOption label="低" :value="1" />
                <ElOption label="中" :value="2" />
                <ElOption label="高" :value="3" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElFormItem label="预算金额" prop="budget">
          <ElInputNumber
            v-model="createForm.budget"
            :min="0"
            :precision="2"
            :step="1000"
            controls-position="right"
            style="width: 100%"
          />
        </ElFormItem>

        <ElFormItem label="项目描述" prop="description">
          <ElInput
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入项目描述"
            maxlength="1000"
            show-word-limit
          />
        </ElFormItem>

        <ElFormItem label="备注" prop="remark">
          <ElInput v-model="createForm.remark" placeholder="请输入备注" maxlength="500" />
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElSpace>
          <ElButton @click="createDialogVisible = false">取消</ElButton>
          <ElButton type="primary" @click="handleCreateSubmit" :loading="createSubmitting">
            创建项目
          </ElButton>
        </ElSpace>
      </template>
    </ElDialog>

    <!-- 从模板创建弹窗 -->
    <ElDialog
      v-model="templateDialogVisible"
      title="从模板创建项目"
      width="900px"
      :close-on-click-modal="false"
    >
      <ElSteps :active="templateStep" finish-status="success" style="margin-bottom: 30px">
        <ElStep title="选择模板" />
        <ElStep title="填写信息" />
      </ElSteps>

      <!-- 步骤1: 选择模板 -->
      <div v-if="templateStep === 0">
        <ElTable
          v-loading="templateLoading"
          :data="templateList"
          highlight-current-row
          @current-change="handleTemplateSelect"
          max-height="400"
        >
          <ElTableColumn type="index" label="序号" width="60" />
          <ElTableColumn prop="name" label="模板名称" min-width="150" />
          <ElTableColumn prop="category" label="分类" width="100" />
          <ElTableColumn prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <ElTableColumn label="任务数" width="80">
            <template #default="{ row }">
              <ElTag size="small">{{ row.taskCount || 0 }}</ElTag>
            </template>
          </ElTableColumn>
        </ElTable>
      </div>

      <!-- 步骤2: 填写信息 -->
      <div v-if="templateStep === 1">
        <ElForm
          ref="templateFormRef"
          :model="templateForm"
          :rules="templateRules"
          label-width="100px"
        >
          <ElFormItem label="项目名称" prop="name">
            <ElInput v-model="templateForm.name" placeholder="请输入项目名称" />
          </ElFormItem>

          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="项目分类" prop="category">
                <ElInput v-model="templateForm.category" placeholder="如：软件开发" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="项目经理" prop="managerId">
                <ElSelect
                  v-model="templateForm.managerId"
                  placeholder="请选择项目经理"
                  filterable
                  clearable
                  style="width: 100%"
                >
                  <ElOption
                    v-for="user in userList"
                    :key="user.id"
                    :label="user.nickName"
                    :value="user.id"
                  />
                </ElSelect>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="开始日期" prop="startDate">
                <ElDatePicker v-model="templateForm.startDate" type="date" style="width: 100%" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="结束日期" prop="endDate">
                <ElDatePicker v-model="templateForm.endDate" type="date" style="width: 100%" />
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElFormItem label="预算金额" prop="budget">
            <ElInputNumber
              v-model="templateForm.budget"
              :min="0"
              :precision="2"
              controls-position="right"
              style="width: 100%"
            />
          </ElFormItem>

          <ElFormItem label="项目描述" prop="description">
            <ElInput
              v-model="templateForm.description"
              type="textarea"
              :rows="3"
              maxlength="1000"
              show-word-limit
            />
          </ElFormItem>
        </ElForm>
      </div>

      <template #footer>
        <ElSpace>
          <ElButton @click="handleTemplateCancel">取消</ElButton>
          <ElButton v-if="templateStep === 1" @click="templateStep = 0">上一步</ElButton>
          <ElButton
            v-if="templateStep === 0"
            type="primary"
            @click="handleTemplateNext"
            :disabled="!selectedTemplate"
          >
            下一步
          </ElButton>
          <ElButton
            v-if="templateStep === 1"
            type="primary"
            @click="handleTemplateSubmit"
            :loading="templateSubmitting"
          >
            创建项目
          </ElButton>
        </ElSpace>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/composables/useTable'
  import { projectApi, projectTemplateApi } from '@/api/project'
  import { fetchGetUserList } from '@/api/system-manage'
  import {
    ElMessage,
    ElMessageBox,
    ElTag,
    ElButton,
    ElProgress,
    FormInstance,
    FormRules
  } from 'element-plus'
  import { useRouter } from 'vue-router'
  import { h, reactive } from 'vue'
  import dayjs from 'dayjs'

  defineOptions({ name: 'ProjectList' })

  const router = useRouter()

  const showSearchBar = ref(false)
  const userList = ref<any[]>([])

  // 新建项目弹窗
  const createDialogVisible = ref(false)
  const createFormRef = ref<FormInstance>()
  const createSubmitting = ref(false)
  const createForm = reactive({
    name: '',
    category: '',
    status: 1,
    priority: 2,
    startDate: new Date(),
    endDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
    managerId: undefined as number | undefined,
    budget: 0,
    description: '',
    remark: ''
  })

  const createRules: FormRules = {
    name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
    category: [{ required: true, message: '请输入项目分类', trigger: 'blur' }],
    startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
    endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
    managerId: [{ required: true, message: '请选择项目经理', trigger: 'change' }]
  }

  // 从模板创建弹窗
  const templateDialogVisible = ref(false)
  const templateFormRef = ref<FormInstance>()
  const templateStep = ref(0)
  const templateLoading = ref(false)
  const templateSubmitting = ref(false)
  const templateList = ref<any[]>([])
  const selectedTemplate = ref<any>(null)

  const templateForm = reactive({
    name: '',
    category: '',
    startDate: new Date(),
    endDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
    managerId: undefined as number | undefined,
    budget: 0,
    description: '',
    remark: ''
  })

  const templateRules: FormRules = {
    name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
    category: [{ required: true, message: '请输入项目分类', trigger: 'blur' }],
    startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
    endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
    managerId: [{ required: true, message: '请选择项目经理', trigger: 'change' }]
  }

  const searchForm = ref({
    name: undefined,
    code: undefined,
    status: undefined
  })

  const formItems = computed(() => [
    {
      label: '项目名称',
      key: 'name',
      type: 'input',
      placeholder: '请输入项目名称',
      clearable: true
    },
    {
      label: '项目编号',
      key: 'code',
      type: 'input',
      placeholder: '请输入项目编号',
      clearable: true
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      placeholder: '请选择状态',
      clearable: true,
      options: [
        { label: '计划中', value: 1 },
        { label: '进行中', value: 2 },
        { label: '已完成', value: 3 },
        { label: '已暂停', value: 4 },
        { label: '已取消', value: 5 }
      ]
    }
  ])

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    resetSearchParams,
    handleSizeChange,
    handleCurrentChange,
    refreshData
  } = useTable({
    core: {
      apiFn: projectApi.getProjectList,
      paginationKey: {
        current: 'page',
        size: 'pageSize'
      },
      columnsFactory: () => [
        {
          prop: 'code',
          label: '项目编号',
          width: 150
        },
        {
          prop: 'name',
          label: '项目名称',
          minWidth: 200
        },
        {
          prop: 'category',
          label: '分类',
          width: 120
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row: any) => {
            const statusConfig = {
              1: { type: 'info', text: '计划中' },
              2: { type: 'warning', text: '进行中' },
              3: { type: 'success', text: '已完成' },
              4: { type: '', text: '已暂停' },
              5: { type: 'danger', text: '已取消' }
            }[row.status as number] || { type: '', text: '未知' }
            return h(ElTag, { type: statusConfig.type as any }, () => statusConfig.text)
          }
        },
        {
          prop: 'progress',
          label: '进度',
          width: 150,
          formatter: (row: any) => h(ElProgress, { percentage: row.progress })
        },
        {
          prop: 'managerName',
          label: '项目经理',
          width: 120
        },
        {
          prop: 'startDate',
          label: '时间范围',
          width: 220,
          formatter: (row: any) => {
            const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')
            return `${formatDate(row.startDate)} ~ ${formatDate(row.endDate)}`
          }
        },
        {
          prop: 'operation',
          label: '操作',
          width: 250,
          fixed: 'right',
          formatter: (row: any) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  link: true,
                  type: 'primary',
                  size: 'small',
                  onClick: () => handleViewGantt(row)
                },
                () => '甘特图'
              ),
              h(
                ElButton,
                {
                  link: true,
                  type: 'primary',
                  size: 'small',
                  onClick: () => handleEdit(row)
                },
                () => '编辑'
              ),
              h(
                ElButton,
                {
                  link: true,
                  type: 'danger',
                  size: 'small',
                  onClick: () => handleDelete(row)
                },
                () => '删除'
              )
            ])
        }
      ]
    },
    transform: {
      responseAdapter: (response: any) => {
        // 后端返回格式: { code: 200, msg: "success", data: { records: [], current: 1, size: 10, total: 100 } }
        if (response) {
          return {
            records: response.records || [],
            total: response.total || 0,
            current: response.current || 1,
            size: response.size || 10
          }
        }
        return { records: [], total: 0, current: 1, size: 10 }
      }
    }
  })

  // 搜索
  const handleSearch = () => {
    getData({
      ...(searchForm.value as unknown as Record<string, any>),
      page: 1
    })
  }

  // 重置
  const handleReset = () => {
    resetSearchParams()
    getData({ page: 1 })
  }

  // 加载用户列表
  const loadUserList = async () => {
    try {
      const res: any = await fetchGetUserList({
        current: 1,
        size: 1000
      } as any)
      userList.value = res.records || res.list || []
    } catch (error) {
      console.error('加载用户列表失败:', error)
    }
  }

  // 新建项目
  const handleCreate = async () => {
    createDialogVisible.value = true
    await loadUserList()
  }

  // 新建项目提交
  const handleCreateSubmit = async () => {
    if (!createFormRef.value) return

    try {
      await createFormRef.value.validate()
      createSubmitting.value = true

      await projectApi.createProject({
        name: createForm.name,
        category: createForm.category,
        status: createForm.status,
        priority: createForm.priority,
        startDate: createForm.startDate.toISOString(),
        endDate: createForm.endDate.toISOString(),
        managerId: createForm.managerId!,
        budget: createForm.budget,
        description: createForm.description,
        remark: createForm.remark
      })

      ElMessage.success('项目创建成功')
      createDialogVisible.value = false
      createFormRef.value.resetFields()
      refreshData()
    } catch (error: any) {
      console.error('创建项目失败:', error)
      ElMessage.error(error.message || '创建项目失败')
    } finally {
      createSubmitting.value = false
    }
  }

  // 从模板创建
  const handleCreateFromTemplate = async () => {
    templateDialogVisible.value = true
    templateStep.value = 0
    selectedTemplate.value = null

    try {
      templateLoading.value = true
      // 并行加载模板列表和用户列表
      const [templateRes] = await Promise.all([
        projectTemplateApi.getTemplateList({ page: 1, pageSize: 100 }),
        loadUserList()
      ])
      templateList.value = (templateRes as any).records || []
    } catch (error) {
      console.error('加载数据失败:', error)
      ElMessage.error('加载数据失败')
    } finally {
      templateLoading.value = false
    }
  }

  // 选择模板
  const handleTemplateSelect = (row: any) => {
    selectedTemplate.value = row
    if (row) {
      templateForm.category = row.category || ''
      templateForm.description = row.description || ''
    }
  }

  // 模板创建下一步
  const handleTemplateNext = () => {
    if (!selectedTemplate.value) {
      ElMessage.warning('请先选择一个模板')
      return
    }
    templateStep.value = 1
  }

  // 模板创建取消
  const handleTemplateCancel = () => {
    templateDialogVisible.value = false
    templateStep.value = 0
    selectedTemplate.value = null
    templateFormRef.value?.resetFields()
  }

  // 模板创建提交
  const handleTemplateSubmit = async () => {
    if (!templateFormRef.value || !selectedTemplate.value) return

    try {
      await templateFormRef.value.validate()
      templateSubmitting.value = true

      await projectApi.createProjectFromTemplate({
        templateId: selectedTemplate.value.id,
        name: templateForm.name,
        category: templateForm.category,
        startDate: templateForm.startDate.toISOString(),
        endDate: templateForm.endDate.toISOString(),
        managerId: templateForm.managerId!,
        budget: templateForm.budget,
        description: templateForm.description,
        remark: templateForm.remark || ''
      })

      ElMessage.success('项目创建成功')
      templateDialogVisible.value = false
      templateStep.value = 0
      templateFormRef.value.resetFields()
      refreshData()
    } catch (error: any) {
      console.error('创建项目失败:', error)
      ElMessage.error(error.message || '创建项目失败')
    } finally {
      templateSubmitting.value = false
    }
  }

  // 查看甘特图
  const handleViewGantt = (row: any) => {
    router.push(`/project/gantt/${row.id}`)
  }

  // 编辑
  const handleEdit = (row: any) => {
    router.push(`/project/edit/${row.id}`)
  }

  // 删除
  const handleDelete = async (row: any) => {
    try {
      await ElMessageBox.confirm(`确定要删除项目"${row.name}"吗？此操作不可恢复！`, '删除确认', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })

      await projectApi.deleteProject(row.id)
      ElMessage.success('删除成功')
      refreshData()
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  // 初始化
  onMounted(() => {
    getData({ page: 1 })
  })
</script>

<style scoped lang="scss">
  .project-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }
</style>
