<template>
  <div class="template-page art-full-height">
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
            <ElButton @click="handleCreate" v-ripple>新建模板</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 模板列表 -->
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

    <!-- 新建/编辑模板弹窗 -->
    <ElDialog
      v-model="dialogVisible"
      :title="editingTemplate ? '编辑模板' : '新建模板'"
      width="900px"
      :close-on-click-modal="false"
    >
      <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="模板名称" prop="name">
              <ElInput v-model="form.name" placeholder="请输入模板名称" maxlength="200" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="模板分类" prop="category">
              <ElInput v-model="form.category" placeholder="如：软件开发" maxlength="50" />
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="项目周期" prop="duration">
              <ElInputNumber
                v-model="form.duration"
                :min="1"
                :max="365"
                controls-position="right"
                placeholder="天数"
                style="width: 100%"
              />
              <div class="form-tip">预计项目完成天数</div>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="是否公开" prop="isPublic">
              <ElSwitch v-model="form.isPublic" />
              <div class="form-tip">公开模板可被其他用户使用</div>
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElFormItem label="模板描述" prop="description">
          <ElInput
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入模板描述"
            maxlength="1000"
            show-word-limit
          />
        </ElFormItem>

        <ElFormItem label="模板任务">
          <ElAlert
            v-if="editingTemplate"
            title="注意：当前版本不支持编辑模板任务，只能修改模板基本信息"
            type="warning"
            :closable="false"
            style="margin-bottom: 16px"
          />
          <div class="task-list">
            <div class="task-header">
              <span>任务列表</span>
              <ElButton
                type="primary"
                size="small"
                @click="handleAddTask"
                :disabled="!!editingTemplate"
              >
                添加任务
              </ElButton>
            </div>
            <div v-if="form.tasks.length === 0" class="empty-tasks">
              <ElEmpty description="暂无任务，请添加任务" :image-size="80" />
            </div>
            <div v-else class="tasks-container">
              <div
                v-for="(task, index) in form.tasks"
                :key="task.tempId"
                class="task-item"
                :class="{ 'is-milestone': task.isMilestone }"
              >
                <div class="task-content">
                  <div class="task-main">
                    <ElIcon class="drag-handle">
                      <Rank />
                    </ElIcon>
                    <ElInput
                      v-model="task.name"
                      placeholder="任务名称"
                      size="small"
                      style="flex: 1; margin-right: 10px"
                      :disabled="!!editingTemplate"
                    />
                    <ElInputNumber
                      v-model="task.estimatedHours"
                      :min="0"
                      :max="9999"
                      :step="0.5"
                      :precision="1"
                      controls-position="right"
                      placeholder="预估工时"
                      size="small"
                      style="width: 120px; margin-right: 10px"
                    />
                    <ElSelect
                      v-model="task.priority"
                      placeholder="优先级"
                      size="small"
                      style="width: 100px; margin-right: 10px"
                    >
                      <ElOption label="低" :value="1" />
                      <ElOption label="中" :value="2" />
                      <ElOption label="高" :value="3" />
                    </ElSelect>
                    <ElCheckbox v-model="task.isMilestone" size="small"> 里程碑 </ElCheckbox>
                  </div>
                  <div class="task-actions">
                    <ElButton 
                      type="primary" 
                      size="small" 
                      link 
                      @click="handleAddSubTask(index)"
                      :disabled="!!editingTemplate"
                    >
                      添加子任务
                    </ElButton>
                    <ElButton 
                      type="danger" 
                      size="small" 
                      link 
                      @click="handleRemoveTask(index)"
                      :disabled="!!editingTemplate"
                    >
                      删除
                    </ElButton>
                  </div>
                </div>
                <!-- 子任务 -->
                <div v-if="task.children && task.children.length > 0" class="sub-tasks">
                  <div
                    v-for="(subTask, subIndex) in task.children"
                    :key="subTask.tempId"
                    class="sub-task-item"
                  >
                    <ElIcon class="drag-handle">
                      <Rank />
                    </ElIcon>
                    <ElInput
                      v-model="subTask.name"
                      placeholder="子任务名称"
                      size="small"
                      style="flex: 1; margin-right: 10px"
                    />
                    <ElInputNumber
                      v-model="subTask.estimatedHours"
                      :min="0"
                      :max="9999"
                      :step="0.5"
                      :precision="1"
                      controls-position="right"
                      placeholder="预估工时"
                      size="small"
                      style="width: 120px; margin-right: 10px"
                    />
                    <ElSelect
                      v-model="subTask.priority"
                      placeholder="优先级"
                      size="small"
                      style="width: 100px; margin-right: 10px"
                    >
                      <ElOption label="低" :value="1" />
                      <ElOption label="中" :value="2" />
                      <ElOption label="高" :value="3" />
                    </ElSelect>
                    <ElButton
                      type="danger"
                      size="small"
                      link
                      @click="handleRemoveSubTask(index, subIndex)"
                      :disabled="!!editingTemplate"
                    >
                      删除
                    </ElButton>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElSpace>
          <ElButton @click="dialogVisible = false">取消</ElButton>
          <ElButton type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingTemplate ? '更新模板' : '创建模板' }}
          </ElButton>
        </ElSpace>
      </template>
    </ElDialog>

    <!-- 查看模板详情弹窗 -->
    <ElDialog
      v-model="detailDialogVisible"
      title="模板详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="templateDetail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="模板名称">{{ templateDetail.name }}</ElDescriptionsItem>
          <ElDescriptionsItem label="分类">{{ templateDetail.category }}</ElDescriptionsItem>
          <ElDescriptionsItem label="任务数量">
            {{ templateDetail.taskCount || 0 }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="创建人">{{ templateDetail.createdBy }}</ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间" :span="2">
            {{ dayjs(templateDetail.createdAt).format('YYYY-MM-DD HH:mm:ss') }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="描述" :span="2">
            {{ templateDetail.description || '无' }}
          </ElDescriptionsItem>
        </ElDescriptions>

        <ElDivider>任务列表</ElDivider>
        <div v-if="templateDetail.tasks && templateDetail.tasks.length > 0">
          <div v-for="task in templateDetail.tasks" :key="task.id" class="detail-task-item">
            <div class="task-info">
              <ElTag v-if="task.isMilestone" type="warning" size="small">里程碑</ElTag>
              <span class="task-name">{{ task.name }}</span>
              <ElTag type="info" size="small">{{ task.estimatedHours }}h</ElTag>
              <ElTag
                :type="task.priority === 3 ? 'danger' : task.priority === 2 ? 'warning' : 'info'"
                size="small"
              >
                {{ ['', '低', '中', '高'][task.priority] }}
              </ElTag>
            </div>
            <!-- 子任务 -->
            <div v-if="task.children && task.children.length > 0" class="sub-tasks-detail">
              <div v-for="subTask in task.children" :key="subTask.id" class="sub-task-info">
                <span class="sub-task-name">{{ subTask.name }}</span>
                <ElTag type="info" size="small">{{ subTask.estimatedHours }}h</ElTag>
                <ElTag
                  :type="
                    subTask.priority === 3 ? 'danger' : subTask.priority === 2 ? 'warning' : 'info'
                  "
                  size="small"
                >
                  {{ ['', '低', '中', '高'][subTask.priority] }}
                </ElTag>
              </div>
            </div>
          </div>
        </div>
        <ElEmpty v-else description="该模板暂无任务" :image-size="80" />
      </div>

      <template #footer>
        <ElButton @click="detailDialogVisible = false">关闭</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/composables/useTable'
  import { projectTemplateApi } from '@/api/project'
  import { ElMessage, ElMessageBox, ElTag, ElButton, FormInstance, FormRules } from 'element-plus'
  import { Rank } from '@element-plus/icons-vue'
  import { h, reactive } from 'vue'
  import dayjs from 'dayjs'

  defineOptions({ name: 'ProjectTemplate' })

  const showSearchBar = ref(false)

  // 弹窗状态
  const dialogVisible = ref(false)
  const detailDialogVisible = ref(false)
  const formRef = ref<FormInstance>()
  const submitting = ref(false)
  const editingTemplate = ref<any>(null)
  const templateDetail = ref<any>(null)

  // 表单数据
  const form = reactive({
    name: '',
    category: '',
    description: '',
    duration: 30,
    isPublic: false,
    tasks: [] as any[]
  })

  // 表单验证规则
  const rules: FormRules = {
    name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
    category: [{ required: true, message: '请输入模板分类', trigger: 'blur' }],
    duration: [{ required: true, message: '请输入项目周期', trigger: 'blur' }]
  }

  // 搜索表单
  const searchForm = ref({
    name: undefined,
    category: undefined
  })

  const formItems = computed(() => [
    {
      label: '模板名称',
      key: 'name',
      type: 'input',
      placeholder: '请输入模板名称',
      clearable: true
    },
    {
      label: '模板分类',
      key: 'category',
      type: 'input',
      placeholder: '请输入模板分类',
      clearable: true
    }
  ])

  // 任务ID生成器
  let taskIdCounter = 1
  const generateTempId = () => `temp_${Date.now()}_${taskIdCounter++}`

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
      apiFn: projectTemplateApi.getTemplateList,
      paginationKey: {
        current: 'page',
        size: 'pageSize'
      },
      columnsFactory: () => [
        {
          prop: 'name',
          label: '模板名称',
          minWidth: 200
        },
        {
          prop: 'category',
          label: '分类',
          width: 120
        },
        {
          prop: 'description',
          label: '描述',
          minWidth: 250,
          showOverflowTooltip: true
        },
        {
          prop: 'taskCount',
          label: '任务数',
          width: 100,
          formatter: (row: any) => h(ElTag, { size: 'small' }, () => row.taskCount || 0)
        },
        {
          prop: 'createdBy',
          label: '创建人',
          width: 120
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 180,
          formatter: (row: any) => dayjs(row.createdAt).format('YYYY-MM-DD HH:mm')
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
                  onClick: () => handleView(row)
                },
                () => '查看'
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

  // 新建模板
  const handleCreate = () => {
    editingTemplate.value = null
    resetForm()
    dialogVisible.value = true
  }

  // 编辑模板
  const handleEdit = async (row: any) => {
    try {
      editingTemplate.value = row
      const detail: any = await projectTemplateApi.getTemplateDetail(row.id)
      form.name = detail.name
      form.category = detail.category
      form.description = detail.description || ''
      form.duration = detail.duration || 30
      form.isPublic = detail.isPublic || false
      form.tasks = (detail.tasks || []).map((task: any) => ({
        ...task,
        tempId: generateTempId(),
        children: (task.children || []).map((subTask: any) => ({
          ...subTask,
          tempId: generateTempId()
        }))
      }))
      dialogVisible.value = true
    } catch (error) {
      console.error('获取模板详情失败:', error)
      ElMessage.error('获取模板详情失败')
    }
  }

  // 查看模板详情
  const handleView = async (row: any) => {
    try {
      templateDetail.value = await projectTemplateApi.getTemplateDetail(row.id)
      detailDialogVisible.value = true
    } catch (error) {
      console.error('获取模板详情失败:', error)
      ElMessage.error('获取模板详情失败')
    }
  }

  // 删除模板
  const handleDelete = async (row: any) => {
    try {
      await ElMessageBox.confirm(`确定要删除模板"${row.name}"吗？此操作不可恢复！`, '删除确认', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })

      await projectTemplateApi.deleteTemplate(row.id)
      ElMessage.success('删除成功')
      refreshData()
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  // 添加任务
  const handleAddTask = () => {
    form.tasks.push({
      tempId: generateTempId(),
      name: '',
      estimatedHours: 8,
      priority: 2,
      isMilestone: false,
      children: []
    })
  }

  // 添加子任务
  const handleAddSubTask = (taskIndex: number) => {
    if (!form.tasks[taskIndex].children) {
      form.tasks[taskIndex].children = []
    }
    form.tasks[taskIndex].children.push({
      tempId: generateTempId(),
      name: '',
      estimatedHours: 4,
      priority: 2,
      isMilestone: false
    })
  }

  // 删除任务
  const handleRemoveTask = (index: number) => {
    form.tasks.splice(index, 1)
  }

  // 删除子任务
  const handleRemoveSubTask = (taskIndex: number, subIndex: number) => {
    form.tasks[taskIndex].children.splice(subIndex, 1)
  }

  // 重置表单
  const resetForm = () => {
    form.name = ''
    form.category = ''
    form.description = ''
    form.duration = 30
    form.isPublic = false
    form.tasks = []
    formRef.value?.clearValidate()
  }

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()

      // 验证任务
      if (form.tasks.length === 0) {
        ElMessage.warning('请至少添加一个任务')
        return
      }

      for (let i = 0; i < form.tasks.length; i++) {
        const task = form.tasks[i]
        if (!task.name.trim()) {
          ElMessage.warning(`请填写第${i + 1}个任务的名称`)
          return
        }
        if (task.children) {
          for (let j = 0; j < task.children.length; j++) {
            const subTask = task.children[j]
            if (!subTask.name.trim()) {
              ElMessage.warning(`请填写第${i + 1}个任务的第${j + 1}个子任务名称`)
              return
            }
          }
        }
      }

      submitting.value = true

      const params = {
        name: form.name,
        category: form.category,
        description: form.description,
        duration: form.duration,
        isPublic: form.isPublic,
        tasks: form.tasks.map((task, index) => ({
          id: task.id,
          name: task.name,
          estimatedHours: task.estimatedHours,
          priority: task.priority,
          isMilestone: task.isMilestone,
          sortOrder: index + 1,
          children: (task.children || []).map((subTask: any, subIndex: number) => ({
            id: subTask.id,
            name: subTask.name,
            estimatedHours: subTask.estimatedHours,
            priority: subTask.priority,
            isMilestone: false,
            sortOrder: subIndex + 1
          }))
        }))
      }

      if (editingTemplate.value) {
        // 更新模板时只提交基本信息，不包含任务
        const updateParams = {
          id: editingTemplate.value.id,
          name: form.name,
          category: form.category,
          description: form.description,
          duration: form.duration,
          isPublic: form.isPublic,
          status: 1 // 默认状态为启用
        }
        await projectTemplateApi.updateTemplate(updateParams)
        ElMessage.success('模板更新成功')
        ElMessage.warning('注意：任务更新功能需要单独实现，当前只更新了模板基本信息')
      } else {
        await projectTemplateApi.createTemplate(params)
        ElMessage.success('模板创建成功')
      }

      dialogVisible.value = false
      resetForm()
      refreshData()
    } catch (error: any) {
      console.error('操作失败:', error)
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitting.value = false
    }
  }

  // 初始化
  onMounted(() => {
    getData({ page: 1 })
  })
</script>

<style scoped lang="scss">
  .template-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }

  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-placeholder);
    margin-top: 4px;
  }

  .task-list {
    border: 1px solid var(--el-border-color);
    border-radius: 6px;
    overflow: hidden;

    .task-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      background-color: var(--el-bg-color-page);
      border-bottom: 1px solid var(--el-border-color);
      font-weight: 500;
    }

    .empty-tasks {
      padding: 40px 20px;
    }

    .tasks-container {
      padding: 16px;
    }

    .task-item {
      margin-bottom: 16px;
      padding: 16px;
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 6px;
      background-color: var(--el-bg-color);

      &.is-milestone {
        border-color: var(--el-color-warning);
        background-color: var(--el-color-warning-light-9);
      }

      .task-content {
        .task-main {
          display: flex;
          align-items: center;
          margin-bottom: 8px;

          .drag-handle {
            margin-right: 8px;
            cursor: move;
            color: var(--el-text-color-placeholder);
          }
        }

        .task-actions {
          display: flex;
          gap: 8px;
        }
      }

      .sub-tasks {
        margin-top: 12px;
        padding-left: 20px;
        border-left: 2px solid var(--el-border-color);

        .sub-task-item {
          display: flex;
          align-items: center;
          margin-bottom: 8px;
          padding: 8px;
          background-color: var(--el-bg-color-page);
          border-radius: 4px;

          .drag-handle {
            margin-right: 8px;
            cursor: move;
            color: var(--el-text-color-placeholder);
          }
        }
      }
    }
  }

  .detail-task-item {
    margin-bottom: 16px;
    padding: 16px;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 6px;

    .task-info {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;

      .task-name {
        font-weight: 500;
        flex: 1;
      }
    }

    .sub-tasks-detail {
      padding-left: 20px;
      border-left: 2px solid var(--el-border-color);

      .sub-task-info {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;
        padding: 8px;
        background-color: var(--el-bg-color-page);
        border-radius: 4px;

        .sub-task-name {
          flex: 1;
        }
      }
    }
  }
</style>