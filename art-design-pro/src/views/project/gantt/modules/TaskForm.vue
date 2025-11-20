<template>
  <ElDialog
    :model-value="visible"
    :title="title"
    width="600px"
    :close-on-click-modal="false"
    @update:model-value="(val) => emit('update:visible', val)"
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
      <ElFormItem label="任务名称" prop="name">
        <ElInput v-model="form.name" placeholder="请输入任务名称" maxlength="200" />
      </ElFormItem>
      <ElFormItem label="父任务" prop="parentId">
        <ElSelect
          v-model="form.parentId"
          placeholder="请选择父任务（可不选）"
          clearable
          filterable
          style="width: 100%"
        >
          <ElOption
            v-for="task in parentTaskList"
            :key="task.id"
            :label="task.name"
            :value="task.id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="任务类型" prop="type">
        <ElRadioGroup v-model="form.type">
          <ElRadio label="task">普通任务</ElRadio>
          <ElRadio label="milestone">里程碑</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="优先级" prop="priority">
        <ElRadioGroup v-model="form.priority">
          <ElRadio :label="1">低</ElRadio>
          <ElRadio :label="2">中</ElRadio>
          <ElRadio :label="3">高</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="开始日期" prop="startDate">
        <ElDatePicker
          v-model="form.startDate"
          type="date"
          placeholder="选择开始日期"
          style="width: 100%"
        />
      </ElFormItem>
      <ElFormItem label="结束日期" prop="endDate">
        <ElDatePicker
          v-model="form.endDate"
          type="date"
          placeholder="选择结束日期"
          style="width: 100%"
        />
      </ElFormItem>
      <ElFormItem label="负责人" prop="assigneeId">
        <ElInput v-model.number="form.assigneeId" placeholder="请输入负责人ID" />
      </ElFormItem>
      <ElFormItem label="预估工时" prop="estimatedHours">
        <ElInputNumber v-model="form.estimatedHours" :min="0" :precision="1" style="width: 100%" />
      </ElFormItem>
      <ElFormItem label="任务描述" prop="description">
        <ElInput
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入任务描述"
        />
      </ElFormItem>
      <ElFormItem label="备注" prop="remark">
        <ElInput v-model="form.remark" placeholder="请输入备注" maxlength="500" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElSpace>
        <ElButton @click="handleClose">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit" :loading="submitting">确定</ElButton>
      </ElSpace>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ref, reactive, watch } from 'vue'
  import { ElMessage, FormInstance, FormRules } from 'element-plus'

  interface TaskFormData {
    name: string
    parentId?: number
    type: string
    priority: number
    startDate: Date
    endDate: Date
    assigneeId?: number
    estimatedHours: number
    description: string
    remark: string
  }

  interface ParentTask {
    id: number
    name: string
  }

  const props = defineProps<{
    visible: boolean
    title: string
    projectId: number
  }>()

  const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }>()

  const formRef = ref<FormInstance>()
  const submitting = ref(false)
  const parentTaskList = ref<ParentTask[]>([])

  const form = reactive<TaskFormData>({
    name: '',
    parentId: undefined,
    type: 'task',
    priority: 2,
    startDate: new Date(),
    endDate: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000),
    assigneeId: undefined,
    estimatedHours: 0,
    description: '',
    remark: ''
  })

  const rules: FormRules = {
    name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
    priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
    startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
    endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
  }

  const handleClose = () => {
    emit('update:visible', false)
    formRef.value?.resetFields()
  }

  // 加载父任务列表
  const loadParentTasks = async () => {
    try {
      const { taskApi } = await import('@/api/project')
      // 使用甘特图数据接口获取所有任务
      const res: any = await taskApi.getGanttData(props.projectId)
      parentTaskList.value =
        res.tasks?.map((task: any) => ({
          id: task.id,
          name: task.text || task.name
        })) || []
    } catch (error) {
      console.error('加载父任务列表失败:', error)
    }
  }

  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitting.value = true

      const { taskApi } = await import('@/api/project')
      await taskApi.createTask({
        projectId: props.projectId,
        parentId: form.parentId || 0,
        name: form.name,
        description: form.description,
        type: form.type,
        priority: form.priority,
        startDate: form.startDate.toISOString(),
        endDate: form.endDate.toISOString(),
        assigneeId: form.assigneeId || 0,
        estimatedHours: form.estimatedHours,
        sort: 0,
        remark: form.remark
      })

      ElMessage.success('任务创建成功')
      handleClose()
      emit('success')
    } catch (error: any) {
      if (error?.message) {
        return
      }
      console.error('创建任务失败:', error)
      ElMessage.error('创建任务失败')
    } finally {
      submitting.value = false
    }
  }

  watch(
    () => props.visible,
    (val) => {
      if (val) {
        // 重置表单
        Object.assign(form, {
          name: '',
          parentId: undefined,
          type: 'task',
          priority: 2,
          startDate: new Date(),
          endDate: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000),
          assigneeId: undefined,
          estimatedHours: 0,
          description: '',
          remark: ''
        })
        // 加载父任务列表
        loadParentTasks()
      }
    }
  )
</script>
