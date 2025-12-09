<template>
  <div class="apply-form-page">
    <!-- 头部 -->
    <ElCard shadow="never" class="header-card">
      <div class="page-header">
        <div class="header-left">
          <ElButton @click="handleBack" v-ripple>
            <el-icon><ArrowLeft /></el-icon>
            返回
          </ElButton>
          <span class="page-title">发起申请</span>
        </div>
        <div class="header-right">
          <ElButton @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </ElButton>
          <ElButton type="primary" :loading="submitting" @click="handleSubmit" v-ripple>
            <el-icon><Check /></el-icon>
            提交申请
          </ElButton>
        </div>
      </div>
    </ElCard>

    <!-- 流程信息 -->
    <ElCard v-loading="loading" shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><Document /></el-icon>
          <span>流程信息</span>
        </div>
      </template>
      <ElDescriptions :column="3" border v-if="processDetail">
        <ElDescriptionsItem label="流程名称">{{ processDetail.name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="流程编码">{{ processDetail.code }}</ElDescriptionsItem>
        <ElDescriptionsItem label="版本">v{{ processDetail.version }}</ElDescriptionsItem>
        <ElDescriptionsItem label="分类">{{ processDetail.category || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="描述" :span="2">
          {{ processDetail.description || '暂无描述' }}
        </ElDescriptionsItem>
      </ElDescriptions>
    </ElCard>

    <!-- 申请表单 -->
    <ElCard v-loading="loading" shadow="never" class="form-card">
      <template #header>
        <div class="card-header">
          <el-icon><EditPen /></el-icon>
          <span>申请表单</span>
          <span class="required-tip">* 为必填项</span>
        </div>
      </template>

      <ElForm
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        label-position="right"
        class="apply-form"
      >
        <!-- 申请标题（必填） -->
        <ElFormItem label="申请标题" prop="title" required>
          <ElInput
            v-model="formData.title"
            placeholder="请输入申请标题"
            maxlength="100"
            show-word-limit
          />
        </ElFormItem>

        <ElDivider v-if="formSchema && formSchema.fields.length > 0">
          表单字段
        </ElDivider>

        <!-- 动态表单字段 -->
        <template v-if="formSchema && formSchema.fields.length > 0">
          <ElFormItem
            v-for="field in formSchema.fields"
            :key="field.key"
            :label="field.label"
            :prop="`data.${field.key}`"
            :required="field.required"
          >
            <!-- 文本输入 -->
            <ElInput
              v-if="field.type === 'text'"
              v-model="formData.data[field.key]"
              :placeholder="field.placeholder || `请输入${field.label}`"
              :maxlength="field.validation?.maxLength"
              :show-word-limit="!!field.validation?.maxLength"
            />

            <!-- 数字输入 -->
            <ElInputNumber
              v-else-if="field.type === 'number'"
              v-model="formData.data[field.key]"
              :placeholder="field.placeholder || `请输入${field.label}`"
              :min="field.validation?.min"
              :max="field.validation?.max"
              controls-position="right"
              style="width: 100%"
            />

            <!-- 日期选择 -->
            <ElDatePicker
              v-else-if="field.type === 'date'"
              v-model="formData.data[field.key]"
              type="date"
              :placeholder="field.placeholder || `请选择${field.label}`"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />

            <!-- 下拉选择 -->
            <ElSelect
              v-else-if="field.type === 'select'"
              v-model="formData.data[field.key]"
              :placeholder="field.placeholder || `请选择${field.label}`"
              style="width: 100%"
              clearable
            >
              <ElOption
                v-for="option in field.options"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </ElSelect>

            <!-- 多行文本 -->
            <ElInput
              v-else-if="field.type === 'textarea'"
              v-model="formData.data[field.key]"
              type="textarea"
              :rows="4"
              :placeholder="field.placeholder || `请输入${field.label}`"
              :maxlength="field.validation?.maxLength"
              :show-word-limit="!!field.validation?.maxLength"
            />

            <!-- 文件上传 -->
            <ElUpload
              v-else-if="field.type === 'file'"
              v-model:file-list="formData.data[field.key]"
              action="/api/file/upload"
              :headers="uploadHeaders"
              :limit="5"
              :on-exceed="handleExceed"
            >
              <ElButton type="primary">
                <el-icon><Upload /></el-icon>
                点击上传
              </ElButton>
              <template #tip>
                <div class="el-upload__tip">
                  支持上传多个文件，单个文件不超过10MB
                </div>
              </template>
            </ElUpload>

            <!-- 默认文本输入 -->
            <ElInput
              v-else
              v-model="formData.data[field.key]"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
          </ElFormItem>
        </template>

        <!-- 无表单字段提示 -->
        <ElEmpty
          v-else-if="!loading && processDetail"
          description="该流程未配置表单模板"
          :image-size="80"
        />
      </ElForm>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import { ArrowLeft, Check, Refresh, Document, EditPen, Upload } from '@element-plus/icons-vue'
import {
  processDefApi,
  processInstApi,
  formTemplateApi,
  type ProcessDefDetailResponse,
  type FormSchema,
  type FormField
} from '@/api/workflow'
import { useUserStore } from '@/store/modules/user'

defineOptions({ name: 'WorkflowApplyForm' })

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 流程定义ID
const processDefId = computed(() => Number(route.params.id))

// 状态
const loading = ref(false)
const submitting = ref(false)

// 流程详情
const processDetail = ref<ProcessDefDetailResponse | null>(null)

// 表单模板
const formSchema = ref<FormSchema | null>(null)

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<{
  title: string
  data: Record<string, any>
}>({
  title: '',
  data: {}
})

// 上传请求头
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

// 生成表单验证规则
const formRules = computed<FormRules>(() => {
  const rules: FormRules = {
    title: [
      { required: true, message: '请输入申请标题', trigger: 'blur' },
      { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
    ]
  }

  if (formSchema.value?.fields) {
    formSchema.value.fields.forEach((field: FormField) => {
      const fieldRules: any[] = []

      // 必填验证
      if (field.required) {
        fieldRules.push({
          required: true,
          message: `请${field.type === 'select' ? '选择' : '输入'}${field.label}`,
          trigger: field.type === 'select' ? 'change' : 'blur'
        })
      }

      // 文本长度验证
      if (field.validation?.minLength || field.validation?.maxLength) {
        fieldRules.push({
          min: field.validation.minLength,
          max: field.validation.maxLength,
          message: `${field.label}长度应在 ${field.validation.minLength || 0} 到 ${field.validation.maxLength || '∞'} 个字符`,
          trigger: 'blur'
        })
      }

      // 数字范围验证
      if (field.type === 'number' && (field.validation?.min !== undefined || field.validation?.max !== undefined)) {
        fieldRules.push({
          type: 'number',
          min: field.validation?.min,
          max: field.validation?.max,
          message: `${field.label}应在 ${field.validation?.min ?? '-∞'} 到 ${field.validation?.max ?? '∞'} 之间`,
          trigger: 'blur'
        })
      }

      // 正则验证
      if (field.validation?.pattern) {
        fieldRules.push({
          pattern: new RegExp(field.validation.pattern),
          message: `${field.label}格式不正确`,
          trigger: 'blur'
        })
      }

      if (fieldRules.length > 0) {
        rules[`data.${field.key}`] = fieldRules
      }
    })
  }

  return rules
})

// 加载流程详情和表单模板
const loadProcessDetail = async () => {
  if (!processDefId.value) {
    ElMessage.error('流程ID不存在')
    return
  }

  loading.value = true
  try {
    // 加载流程定义详情
    const detail = await processDefApi.getDetail(processDefId.value)
    processDetail.value = detail

    // 检查流程是否已发布
    if (detail.status !== 'published') {
      ElMessage.error('该流程尚未发布，无法发起申请')
      router.push('/workflow/apply')
      return
    }

    // 如果有关联表单模板，加载表单结构
    if (detail.formTemplateId) {
      const templateDetail = await formTemplateApi.getDetail(detail.formTemplateId)
      if (templateDetail.schema) {
        formSchema.value = templateDetail.schema
        // 初始化表单数据默认值
        initFormData(templateDetail.schema)
      }
    }

    // 设置默认标题
    formData.title = `${detail.name} - ${new Date().toLocaleDateString()}`
  } catch (error: any) {
    console.error('加载流程详情失败:', error)
    ElMessage.error(error.message || '加载流程详情失败')
  } finally {
    loading.value = false
  }
}

// 初始化表单数据
const initFormData = (schema: FormSchema) => {
  schema.fields.forEach((field: FormField) => {
    if (field.defaultValue !== undefined) {
      formData.data[field.key] = field.defaultValue
    } else if (field.type === 'number') {
      formData.data[field.key] = undefined
    } else if (field.type === 'file') {
      formData.data[field.key] = []
    } else {
      formData.data[field.key] = ''
    }
  })
}

// 返回
const handleBack = () => {
  router.push('/workflow/apply')
}

// 重置表单
const handleReset = () => {
  formRef.value?.resetFields()
  if (formSchema.value) {
    initFormData(formSchema.value)
  }
  if (processDetail.value) {
    formData.title = `${processDetail.value.name} - ${new Date().toLocaleDateString()}`
  }
}

// 文件上传超出限制
const handleExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning('最多只能上传5个文件')
}

// 提交申请
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    ElMessage.warning('请完善表单信息')
    return
  }

  if (!processDefId.value) {
    ElMessage.error('流程ID不存在')
    return
  }

  submitting.value = true
  try {
    // 处理文件上传数据
    const submitData: Record<string, any> = {}
    if (formSchema.value?.fields) {
      formSchema.value.fields.forEach((field: FormField) => {
        const value = formData.data[field.key]
        if (field.type === 'file' && Array.isArray(value)) {
          // 提取文件URL
          submitData[field.key] = value.map((f: any) => f.response?.url || f.url).filter(Boolean)
        } else {
          submitData[field.key] = value
        }
      })
    }

    await processInstApi.start({
      processDefId: processDefId.value,
      title: formData.title,
      formData: Object.keys(submitData).length > 0 ? submitData : undefined
    })

    ElMessage.success('申请提交成功')
    // 跳转到我发起的列表
    router.push('/workflow/workbench/initiated')
  } catch (error: any) {
    console.error('提交申请失败:', error)
    ElMessage.error(error.message || '提交申请失败')
  } finally {
    submitting.value = false
  }
}

// 初始化
onMounted(() => {
  loadProcessDetail()
})
</script>

<style lang="scss" scoped>
.apply-form-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
  padding-bottom: 32px;

  .header-card {
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
      gap: 16px;

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

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    font-weight: 500;

    .el-icon {
      color: var(--el-color-primary);
    }

    .required-tip {
      margin-left: auto;
      font-size: 12px;
      font-weight: normal;
      color: var(--el-color-danger);
    }
  }

  .info-card {
    :deep(.el-card__header) {
      padding: 12px 20px;
      background: var(--el-fill-color-lighter);
    }
  }

  .form-card {
    :deep(.el-card__header) {
      padding: 12px 20px;
      background: var(--el-fill-color-lighter);
    }

    :deep(.el-card__body) {
      padding: 24px;
    }
  }

  .apply-form {
    max-width: 800px;

    :deep(.el-form-item) {
      margin-bottom: 20px;
    }

    :deep(.el-divider) {
      margin: 24px 0;
    }

    :deep(.el-upload__tip) {
      color: var(--el-text-color-secondary);
    }
  }
}
</style>
