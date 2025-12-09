<template>
  <div class="process-def-edit-page">
    <ElCard shadow="never" class="header-card">
      <div class="page-header">
        <div class="header-left">
          <ElButton @click="handleBack" v-ripple>
            <el-icon><ArrowLeft /></el-icon>
            返回
          </ElButton>
          <span class="page-title">{{ isEdit ? '编辑流程定义' : '新增流程定义' }}</span>
          <ElTag v-if="isEdit && processData.status" :type="statusConfig[processData.status]?.type || 'info'">
            {{ statusConfig[processData.status]?.text || processData.status }}
          </ElTag>
          <ElTag v-if="isEdit" type="info">v{{ processData.version }}</ElTag>
        </div>
        <div class="header-right">
          <ElButton :loading="saving" @click="handleSave" v-ripple>
            <el-icon><DocumentCopy /></el-icon>
            保存草稿
          </ElButton>
          <ElButton type="success" :loading="publishing" @click="handlePublish" v-ripple>
            <el-icon><Upload /></el-icon>
            发布
          </ElButton>
        </div>
      </div>
    </ElCard>

    <ElCard shadow="never" class="form-card">
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <ElRow :gutter="20">
          <ElCol :span="6">
            <ElFormItem label="流程名称" prop="name">
              <ElInput v-model="formData.name" placeholder="请输入流程名称" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="6">
            <ElFormItem label="流程编码" prop="code">
              <ElInput
                v-model="formData.code"
                placeholder="请输入流程编码"
                :disabled="isEdit"
              />
            </ElFormItem>
          </ElCol>
          <ElCol :span="6">
            <ElFormItem label="分类" prop="category">
              <ElInput v-model="formData.category" placeholder="请输入分类" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="6">
            <ElFormItem label="表单模板" prop="formTemplateId">
              <ElSelect
                v-model="formData.formTemplateId"
                placeholder="请选择表单模板"
                clearable
                filterable
                style="width: 100%"
              >
                <ElOption
                  v-for="item in formTemplateOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </ElSelect>
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElRow :gutter="20">
          <ElCol :span="24">
            <ElFormItem label="描述" prop="description">
              <ElInput
                v-model="formData.description"
                type="textarea"
                :rows="2"
                placeholder="请输入流程描述"
              />
            </ElFormItem>
          </ElCol>
        </ElRow>
      </ElForm>
    </ElCard>

    <ElCard shadow="never" class="designer-card">
      <template #header>
        <div class="card-header">
          <span>流程设计</span>
          <span class="tip">从左侧拖拽节点到画布区域，连接节点并配置属性</span>
        </div>
      </template>
      <ProcessDesigner
        ref="designerRef"
        v-model="formData.graph"
        class="process-designer-container"
        @save="onDesignerSave"
        @publish="onDesignerPublish"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft, DocumentCopy, Upload } from '@element-plus/icons-vue'
import {
  processDefApi,
  formTemplateApi,
  type ProcessGraph,
  type CreateProcessDefParams,
  type UpdateProcessDefParams,
  type FormTemplateResponse
} from '@/api/workflow'
import ProcessDesigner from '@/components/workflow/ProcessDesigner.vue'

defineOptions({ name: 'ProcessDefEdit' })

const route = useRoute()
const router = useRouter()

// 是否编辑模式
const processId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!processId.value)

// 表单引用
const formRef = ref<FormInstance>()
const designerRef = ref<InstanceType<typeof ProcessDesigner>>()

// 状态
const saving = ref(false)
const publishing = ref(false)
const loading = ref(false)

// 流程数据（用于显示状态和版本）
const processData = reactive({
  status: '',
  version: 1
})

// 状态标签配置
const statusConfig: Record<string, { type: 'info' | 'success' | 'danger'; text: string }> = {
  draft: { type: 'info', text: '草稿' },
  published: { type: 'success', text: '已发布' },
  disabled: { type: 'danger', text: '已禁用' }
}

// 表单模板选项
const formTemplateOptions = ref<FormTemplateResponse[]>([])

// 表单数据
const formData = reactive<{
  name: string
  code: string
  description: string
  category: string
  formTemplateId: number | undefined
  graph: ProcessGraph
}>({
  name: '',
  code: '',
  description: '',
  category: '',
  formTemplateId: undefined,
  graph: { nodes: [], edges: [] }
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入流程名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入流程编码', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
      message: '编码必须以字母开头，只能包含字母、数字和下划线',
      trigger: 'blur'
    }
  ]
}

// 加载表单模板列表
const loadFormTemplates = async () => {
  try {
    const res = await formTemplateApi.getList({ pageSize: 1000, status: 'active' })
    formTemplateOptions.value = res.list || []
  } catch (error: any) {
    console.error('加载表单模板失败:', error)
  }
}

// 加载流程详情
const loadProcessDetail = async () => {
  if (!processId.value) return

  loading.value = true
  try {
    const res = await processDefApi.getDetail(Number(processId.value))
    formData.name = res.name
    formData.code = res.code
    formData.description = res.description || ''
    formData.category = res.category || ''
    formData.formTemplateId = res.formTemplateId
    if (res.graph) {
      formData.graph = res.graph
    }
    processData.status = res.status
    processData.version = res.version
  } catch (error: any) {
    console.error('加载流程详情失败:', error)
    ElMessage.error(error.message || '加载流程详情失败')
  } finally {
    loading.value = false
  }
}

// 返回列表
const handleBack = () => {
  router.push('/workflow/process-def')
}

// 验证表单
const validateForm = async (): Promise<boolean> => {
  if (!formRef.value) return false

  try {
    await formRef.value.validate()
    return true
  } catch {
    ElMessage.warning('请完善表单信息')
    return false
  }
}

// 获取流程图数据
const getGraphData = (): ProcessGraph => {
  return designerRef.value?.getGraph() || formData.graph
}

// 保存草稿
const handleSave = async () => {
  if (!(await validateForm())) return

  const graph = getGraphData()

  saving.value = true
  try {
    if (isEdit.value) {
      const params: UpdateProcessDefParams = {
        name: formData.name,
        description: formData.description,
        category: formData.category,
        formTemplateId: formData.formTemplateId,
        graph
      }
      await processDefApi.update(Number(processId.value), params)
      ElMessage.success('保存成功')
    } else {
      const params: CreateProcessDefParams = {
        name: formData.name,
        code: formData.code,
        description: formData.description,
        category: formData.category,
        formTemplateId: formData.formTemplateId,
        graph
      }
      const res = await processDefApi.create(params)
      ElMessage.success('创建成功')
      // 跳转到编辑页面
      router.replace(`/workflow/process-def/edit/${res.id}`)
    }
  } catch (error: any) {
    console.error('保存失败:', error)
    ElMessage.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

// 发布
const handlePublish = async () => {
  if (!(await validateForm())) return

  // 验证流程图
  const validation = designerRef.value?.validate()
  if (validation && !validation.valid) {
    ElMessageBox.alert(
      `<ul style="margin: 0; padding-left: 20px;">${validation.errors.map((e: string) => `<li>${e}</li>`).join('')}</ul>`,
      '流程验证失败',
      {
        dangerouslyUseHTMLString: true,
        type: 'error'
      }
    )
    return
  }

  try {
    await ElMessageBox.confirm(
      '确定要发布此流程吗？发布后将生成新版本，已运行的流程实例不受影响。',
      '发布确认',
      {
        type: 'warning',
        confirmButtonText: '确定发布',
        cancelButtonText: '取消'
      }
    )
  } catch {
    return
  }

  const graph = getGraphData()

  publishing.value = true
  try {
    // 先保存
    if (isEdit.value) {
      const params: UpdateProcessDefParams = {
        name: formData.name,
        description: formData.description,
        category: formData.category,
        formTemplateId: formData.formTemplateId,
        graph
      }
      await processDefApi.update(Number(processId.value), params)
      // 再发布
      await processDefApi.publish(Number(processId.value))
      ElMessage.success('发布成功')
      // 重新加载数据
      await loadProcessDetail()
    } else {
      // 先创建
      const params: CreateProcessDefParams = {
        name: formData.name,
        code: formData.code,
        description: formData.description,
        category: formData.category,
        formTemplateId: formData.formTemplateId,
        graph
      }
      const res = await processDefApi.create(params)
      // 再发布
      await processDefApi.publish(res.id)
      ElMessage.success('发布成功')
      // 跳转到编辑页面
      router.replace(`/workflow/process-def/edit/${res.id}`)
    }
  } catch (error: any) {
    console.error('发布失败:', error)
    ElMessage.error(error.message || '发布失败')
  } finally {
    publishing.value = false
  }
}

// 设计器保存事件
const onDesignerSave = (graph: ProcessGraph) => {
  formData.graph = graph
  handleSave()
}

// 设计器发布事件
const onDesignerPublish = (graph: ProcessGraph) => {
  formData.graph = graph
  handlePublish()
}

// 初始化
onMounted(async () => {
  await loadFormTemplates()
  if (isEdit.value) {
    await loadProcessDetail()
  }
})
</script>

<style lang="scss" scoped>
.process-def-edit-page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  padding-bottom: 15px;

  .header-card {
    flex-shrink: 0;

    :deep(.el-card__body) {
      padding: 12px 20px;
    }
  }

  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .header-left {
      display: flex;
      align-items: center;
      gap: 12px;

      .page-title {
        font-size: 16px;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }
    }

    .header-right {
      display: flex;
      gap: 8px;
    }
  }

  .form-card {
    flex-shrink: 0;

    :deep(.el-card__body) {
      padding: 20px 20px 8px;
    }
  }

  .designer-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 500px;
    overflow: hidden;

    :deep(.el-card__header) {
      padding: 12px 20px;
      border-bottom: 1px solid var(--el-border-color-light);
    }

    :deep(.el-card__body) {
      flex: 1;
      padding: 0;
      overflow: hidden;
    }

    .card-header {
      display: flex;
      align-items: center;
      gap: 12px;

      .tip {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }

    .process-designer-container {
      height: 100%;
    }
  }
}
</style>
