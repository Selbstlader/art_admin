<template>
  <div class="form-template-edit-page">
    <ElCard shadow="never" class="header-card">
      <div class="page-header">
        <div class="header-left">
          <ElButton @click="handleBack" v-ripple>
            <el-icon><ArrowLeft /></el-icon>
            返回
          </ElButton>
          <span class="page-title">{{ isEdit ? '编辑表单模板' : '新增表单模板' }}</span>
        </div>
        <div class="header-right">
          <ElButton type="primary" :loading="saving" @click="handleSave" v-ripple>
            <el-icon><Check /></el-icon>
            保存
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
          <ElCol :span="8">
            <ElFormItem label="模板名称" prop="name">
              <ElInput v-model="formData.name" placeholder="请输入模板名称" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="8">
            <ElFormItem label="模板编码" prop="code">
              <ElInput v-model="formData.code" placeholder="请输入模板编码" :disabled="isEdit" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="8">
            <ElFormItem label="描述" prop="description">
              <ElInput v-model="formData.description" placeholder="请输入描述" />
            </ElFormItem>
          </ElCol>
        </ElRow>
      </ElForm>
    </ElCard>

    <ElCard shadow="never" class="designer-card">
      <template #header>
        <div class="card-header">
          <span>表单设计</span>
          <span class="tip">从左侧拖拽组件到画布区域，配置字段属性</span>
        </div>
      </template>
      <FormDesigner ref="designerRef" v-model="formData.schema" class="form-designer-container" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { ArrowLeft, Check } from '@element-plus/icons-vue'
  import {
    formTemplateApi,
    type FormSchema,
    type CreateFormTemplateParams,
    type UpdateFormTemplateParams
  } from '@/api/workflow'
  import FormDesigner from '@/components/workflow/FormDesigner.vue'

  defineOptions({ name: 'FormTemplateEdit' })

  const route = useRoute()
  const router = useRouter()

  // 是否编辑模式
  const templateId = computed(() => route.params.id as string | undefined)
  const isEdit = computed(() => !!templateId.value)

  // 表单引用
  const formRef = ref<FormInstance>()
  const designerRef = ref<InstanceType<typeof FormDesigner>>()

  // 保存状态
  const saving = ref(false)
  const loading = ref(false)

  // 表单数据
  const formData = reactive<{
    name: string
    code: string
    description: string
    schema: FormSchema
  }>({
    name: '',
    code: '',
    description: '',
    schema: { fields: [] }
  })

  // 表单验证规则
  const formRules: FormRules = {
    name: [
      { required: true, message: '请输入模板名称', trigger: 'blur' },
      { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    code: [
      { required: true, message: '请输入模板编码', trigger: 'blur' },
      {
        pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
        message: '编码必须以字母开头，只能包含字母、数字和下划线',
        trigger: 'blur'
      }
    ]
  }

  // 加载模板详情
  const loadTemplateDetail = async () => {
    if (!templateId.value) return

    loading.value = true
    try {
      const res = await formTemplateApi.getDetail(Number(templateId.value))
      formData.name = res.name
      formData.code = res.code
      formData.description = res.description || ''
      if (res.schema) {
        formData.schema = res.schema
      }
    } catch (error: any) {
      console.error('加载模板详情失败:', error)
      ElMessage.error(error.message || '加载模板详情失败')
    } finally {
      loading.value = false
    }
  }

  // 返回列表
  const handleBack = () => {
    router.push('/workflow/form-template')
  }

  // 保存
  const handleSave = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
    } catch {
      ElMessage.warning('请完善表单信息')
      return
    }

    // 获取设计器中的表单结构
    const schema = designerRef.value?.getSchema() || formData.schema

    if (!schema.fields || schema.fields.length === 0) {
      ElMessage.warning('请至少添加一个表单字段')
      return
    }

    saving.value = true
    try {
      if (isEdit.value) {
        const params: UpdateFormTemplateParams = {
          name: formData.name,
          description: formData.description,
          schema
        }
        await formTemplateApi.update(Number(templateId.value), params)
        ElMessage.success('更新成功')
      } else {
        const params: CreateFormTemplateParams = {
          name: formData.name,
          code: formData.code,
          description: formData.description,
          schema
        }
        await formTemplateApi.create(params)
        ElMessage.success('创建成功')
      }
      router.push('/workflow/form-template')
    } catch (error: any) {
      console.error('保存失败:', error)
      ElMessage.error(error.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  // 初始化
  onMounted(() => {
    if (isEdit.value) {
      loadTemplateDetail()
    }
  })
</script>

<style lang="scss" scoped>
  .form-template-edit-page {
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
        gap: 16px;

        .page-title {
          font-size: 16px;
          font-weight: 500;
          color: var(--el-text-color-primary);
        }
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

      .form-designer-container {
        height: 100%;
      }
    }
  }
</style>
