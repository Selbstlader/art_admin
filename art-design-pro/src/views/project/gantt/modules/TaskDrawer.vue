<template>
  <ElDrawer
    v-model="drawerVisible"
    :title="isEditing ? '编辑任务' : '任务详情'"
    size="500px"
    @close="handleClose"
  >
    <div v-if="taskDetail" class="task-drawer">
      <!-- 查看模式 -->
      <div v-if="!isEditing">
        <ElDescriptions :column="1" border>
          <ElDescriptionsItem label="任务名称">{{ taskDetail.name }}</ElDescriptionsItem>
          <ElDescriptionsItem label="任务类型">
            <ElTag :type="taskDetail.type === 'milestone' ? 'danger' : 'primary'">
              {{ taskDetail.type === 'milestone' ? '里程碑' : '普通任务' }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="优先级">
            <ElTag :type="getPriorityType(taskDetail.priority)">
              {{ getPriorityText(taskDetail.priority) }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="状态">
            <ElTag :type="getStatusType(taskDetail.status)">
              {{ getStatusText(taskDetail.status) }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="进度">
            <ElProgress :percentage="taskDetail.progress" />
          </ElDescriptionsItem>
          <ElDescriptionsItem label="开始日期">
            {{ formatDate(taskDetail.startDate) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="结束日期">
            {{ formatDate(taskDetail.endDate) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="负责人">
            {{ taskDetail.assigneeName || '未分配' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="预估工时">
            {{ taskDetail.estimatedHours || 0 }} 小时
          </ElDescriptionsItem>
          <ElDescriptionsItem label="实际工时">
            {{ taskDetail.actualHours || 0 }} 小时
          </ElDescriptionsItem>
          <ElDescriptionsItem label="是否超期">
            <ElTag :type="taskDetail.isOverdue ? 'danger' : 'success'">
              {{ taskDetail.isOverdue ? '是' : '否' }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="任务描述" v-if="taskDetail.description">
            <div class="description">{{ taskDetail.description }}</div>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="备注" v-if="taskDetail.remark">
            <div class="remark">{{ taskDetail.remark }}</div>
          </ElDescriptionsItem>
        </ElDescriptions>
      </div>

      <!-- 编辑模式 -->
      <div v-else>
        <ElForm ref="formRef" :model="editForm" :rules="rules" label-width="100px">
          <ElFormItem label="任务名称" prop="name">
            <ElInput v-model="editForm.name" placeholder="请输入任务名称" />
          </ElFormItem>
          <ElFormItem label="任务类型" prop="type">
            <ElRadioGroup v-model="editForm.type">
              <ElRadio label="task">普通任务</ElRadio>
              <ElRadio label="milestone">里程碑</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem label="优先级" prop="priority">
            <ElRadioGroup v-model="editForm.priority">
              <ElRadio :label="1">低</ElRadio>
              <ElRadio :label="2">中</ElRadio>
              <ElRadio :label="3">高</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem label="状态" prop="status">
            <ElSelect v-model="editForm.status" style="width: 100%">
              <ElOption :value="1" label="未开始" />
              <ElOption :value="2" label="进行中" />
              <ElOption :value="3" label="已完成" />
              <ElOption :value="4" label="已延期" />
              <ElOption :value="5" label="已取消" />
            </ElSelect>
          </ElFormItem>
          <ElFormItem label="进度" prop="progress">
            <ElSlider v-model="editForm.progress" :min="0" :max="100" :step="5" show-input />
          </ElFormItem>
          <ElFormItem label="开始日期" prop="startDate">
            <ElDatePicker v-model="editForm.startDate" type="date" style="width: 100%" />
          </ElFormItem>
          <ElFormItem label="结束日期" prop="endDate">
            <ElDatePicker v-model="editForm.endDate" type="date" style="width: 100%" />
          </ElFormItem>
          <ElFormItem label="负责人" prop="assigneeId">
            <ElInput v-model.number="editForm.assigneeId" placeholder="请输入负责人ID" />
          </ElFormItem>
          <ElFormItem label="预估工时" prop="estimatedHours">
            <ElInputNumber
              v-model="editForm.estimatedHours"
              :min="0"
              :precision="1"
              style="width: 100%"
            />
          </ElFormItem>
          <ElFormItem label="实际工时" prop="actualHours">
            <ElInputNumber
              v-model="editForm.actualHours"
              :min="0"
              :precision="1"
              style="width: 100%"
            />
          </ElFormItem>
          <ElFormItem label="任务描述" prop="description">
            <ElInput v-model="editForm.description" type="textarea" :rows="4" />
          </ElFormItem>
          <ElFormItem label="备注" prop="remark">
            <ElInput v-model="editForm.remark" placeholder="请输入备注" />
          </ElFormItem>
        </ElForm>
      </div>

      <div class="drawer-footer">
        <ElSpace>
          <ElButton @click="handleClose">关闭</ElButton>
          <ElButton v-if="!isEditing" type="primary" @click="handleEdit">编辑</ElButton>
          <ElButton v-if="isEditing" @click="handleCancelEdit">取消编辑</ElButton>
          <ElButton v-if="isEditing" type="primary" @click="handleSave" :loading="saving"
            >保存</ElButton
          >
          <ElButton type="danger" @click="handleDelete" :loading="deleting">删除任务</ElButton>
        </ElSpace>
      </div>
    </div>
    <ElEmpty v-else description="暂无数据" />
  </ElDrawer>
</template>

<script setup lang="ts">
  import { ref, computed, reactive } from 'vue'
  import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
  import { taskApi } from '@/api/project'

  interface TaskDetail {
    id: number
    name: string
    type: string
    priority: number
    status: number
    progress: number
    startDate: string
    endDate: string
    assigneeId: number
    assigneeName: string
    estimatedHours: number
    actualHours: number
    isOverdue: boolean
    description: string
    remark: string
  }

  interface EditForm {
    name: string
    type: string
    priority: number
    status: number
    progress: number
    startDate: Date
    endDate: Date
    assigneeId: number
    estimatedHours: number
    actualHours: number
    description: string
    remark: string
  }

  const props = defineProps<{
    visible: boolean
    taskId: number | null
  }>()

  const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'deleted'): void
    (e: 'updated'): void
  }>()

  const drawerVisible = computed({
    get: () => props.visible,
    set: (val) => emit('update:visible', val)
  })

  const taskDetail = ref<TaskDetail | null>(null)
  const deleting = ref(false)
  const isEditing = ref(false)
  const saving = ref(false)
  const formRef = ref<FormInstance>()

  const editForm = reactive<EditForm>({
    name: '',
    type: 'task',
    priority: 2,
    status: 1,
    progress: 0,
    startDate: new Date(),
    endDate: new Date(),
    assigneeId: 0,
    estimatedHours: 0,
    actualHours: 0,
    description: '',
    remark: ''
  })

  const rules: FormRules = {
    name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
    startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
    endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
  }

  const loadTaskDetail = async () => {
    if (!props.taskId) return

    try {
      const data = (await taskApi.getTaskDetail(props.taskId)) as any
      taskDetail.value = data
    } catch (error) {
      console.error('加载任务详情失败:', error)
      ElMessage.error('加载任务详情失败')
    }
  }

  const handleClose = () => {
    emit('update:visible', false)
    taskDetail.value = null
    isEditing.value = false
  }

  const handleEdit = () => {
    if (!taskDetail.value) return
    // 将详情数据复制到编辑表单
    Object.assign(editForm, {
      name: taskDetail.value.name,
      type: taskDetail.value.type,
      priority: taskDetail.value.priority,
      status: taskDetail.value.status,
      progress: taskDetail.value.progress,
      startDate: new Date(taskDetail.value.startDate),
      endDate: new Date(taskDetail.value.endDate),
      assigneeId: taskDetail.value.assigneeId,
      estimatedHours: taskDetail.value.estimatedHours,
      actualHours: taskDetail.value.actualHours,
      description: taskDetail.value.description,
      remark: taskDetail.value.remark
    })
    isEditing.value = true
  }

  const handleCancelEdit = () => {
    isEditing.value = false
  }

  const handleSave = async () => {
    if (!formRef.value || !props.taskId) return

    try {
      await formRef.value.validate()
      saving.value = true

      await taskApi.updateTask({
        id: Number(props.taskId),
        name: editForm.name,
        type: editForm.type,
        priority: editForm.priority,
        status: editForm.status,
        progress: editForm.progress,
        startDate: editForm.startDate.toISOString(),
        endDate: editForm.endDate.toISOString(),
        assigneeId: editForm.assigneeId || 0,
        estimatedHours: editForm.estimatedHours || 0,
        actualHours: editForm.actualHours || 0,
        description: editForm.description || '',
        remark: editForm.remark || ''
      })

      ElMessage.success('任务更新成功')
      isEditing.value = false
      await loadTaskDetail()
      emit('updated')
    } catch (error: any) {
      if (error?.message) return
      console.error('更新任务失败:', error)
      ElMessage.error('更新任务失败')
    } finally {
      saving.value = false
    }
  }

  const handleDelete = async () => {
    if (!props.taskId) return

    try {
      await ElMessageBox.confirm('确定要删除该任务吗？删除后无法恢复。', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      deleting.value = true
      await taskApi.deleteTask(props.taskId)
      ElMessage.success('任务已删除')
      handleClose()
      emit('deleted')
    } catch (error: any) {
      if (error === 'cancel') return
      console.error('删除任务失败:', error)
      ElMessage.error('删除任务失败')
    } finally {
      deleting.value = false
    }
  }

  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN')
  }

  const getPriorityText = (priority: number) => {
    const map: Record<number, string> = { 1: '低', 2: '中', 3: '高' }
    return map[priority] || '-'
  }

  const getPriorityType = (priority: number) => {
    const map: Record<number, any> = { 1: 'info', 2: 'warning', 3: 'danger' }
    return map[priority] || 'info'
  }

  const getStatusText = (status: number) => {
    const map: Record<number, string> = {
      1: '未开始',
      2: '进行中',
      3: '已完成',
      4: '已延期',
      5: '已取消'
    }
    return map[status] || '-'
  }

  const getStatusType = (status: number) => {
    const map: Record<number, any> = {
      1: 'info',
      2: 'primary',
      3: 'success',
      4: 'danger',
      5: 'info'
    }
    return map[status] || 'info'
  }

  defineExpose({ loadTaskDetail })
</script>

<style scoped lang="scss">
  .task-drawer {
    .description,
    .remark {
      white-space: pre-wrap;
      word-break: break-all;
    }

    .drawer-footer {
      margin-top: 20px;
      padding-top: 20px;
      border-top: 1px solid var(--el-border-color);
      display: flex;
      justify-content: flex-end;
    }
  }
</style>
