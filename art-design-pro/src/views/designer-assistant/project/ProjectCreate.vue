<template>
  <div class="project-create-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <ElButton link @click="handleBack">
            <ElIcon>
              <ArrowLeft />
            </ElIcon>
            返回
          </ElButton>
          <span class="title">{{ isEdit ? '编辑项目' : '新建项目' }}</span>
        </div>
      </template>

      <ElForm
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        style="max-width: 600px"
        v-loading="loading"
      >
        <ElFormItem label="项目名称" prop="name">
          <ElInput
            v-model="formData.name"
            placeholder="请输入项目名称"
            maxlength="200"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="项目描述" prop="description">
          <ElInput
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入项目描述"
            maxlength="2000"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="面积(m²)" prop="area">
          <ElInputNumber v-model="formData.area" :min="0" :precision="2" style="width: 100%" />
        </ElFormItem>
        <ElFormItem label="预算(元)" prop="budget">
          <ElInputNumber v-model="formData.budget" :min="0" :precision="2" style="width: 100%" />
        </ElFormItem>
        <ElFormItem label="设计风格" prop="style">
          <ElSelect
            v-model="formData.style"
            placeholder="请选择设计风格"
            style="width: 100%"
            clearable
          >
            <ElOption label="现代简约" value="现代简约" />
            <ElOption label="工业风" value="工业风" />
            <ElOption label="新中式" value="新中式" />
            <ElOption label="北欧风" value="北欧风" />
            <ElOption label="轻奢风" value="轻奢风" />
            <ElOption label="日式" value="日式" />
            <ElOption label="美式" value="美式" />
            <ElOption label="其他" value="其他" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="状态" prop="status" v-if="isEdit">
          <ElSelect v-model="formData.status" placeholder="请选择状态" style="width: 100%">
            <ElOption label="草稿" value="draft" />
            <ElOption label="进行中" value="in_progress" />
            <ElOption label="已完成" value="completed" />
            <ElOption label="已归档" value="archived" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建项目' }}
          </ElButton>
          <ElButton @click="handleReset">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Project Create/Edit Component
   * 创建/编辑项目页面组件
   * Requirements: 4.1
   ***/
  import { ref, reactive, onMounted, computed } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { ArrowLeft } from '@element-plus/icons-vue'
  import {
    createDesignerProject,
    updateDesignerProject,
    getDesignerProject,
    type ProjectDetailResponse,
    type ProjectResponse
  } from '@/api/designer-project'

  /*** API Response interface / API 响应接口 ***/
  interface ApiResponse<T = unknown> {
    code: number
    msg?: string
    data?: T
  }

  const router = useRouter()
  const route = useRoute()
  const formRef = ref<FormInstance>()
  const loading = ref(false)
  const submitting = ref(false)

  // 判断是否为编辑模式 / Check if edit mode
  const isEdit = computed(() => !!route.params.id)
  const projectId = computed(() => Number(route.params.id) || 0)

  // 表单数据 / Form data
  const formData = reactive({
    name: '',
    description: '',
    area: 0,
    budget: 0,
    style: '',
    status: 'draft'
  })

  // 表单验证规则 / Form validation rules
  const formRules: FormRules = {
    name: [
      { required: true, message: '请输入项目名称', trigger: 'blur' },
      { min: 1, max: 200, message: '项目名称长度为1-200个字符', trigger: 'blur' }
    ],
    description: [{ max: 2000, message: '项目描述不能超过2000个字符', trigger: 'blur' }],
    area: [{ required: true, message: '请输入面积', trigger: 'blur' }],
    budget: [{ required: true, message: '请输入预算', trigger: 'blur' }]
  }

  /***
   * Get project detail for edit mode
   * 编辑模式下获取项目详情
   ***/
  const getProjectDetail = async () => {
    if (!isEdit.value) return

    loading.value = true
    try {
      const res = (await getDesignerProject(
        projectId.value
      )) as unknown as ApiResponse<ProjectDetailResponse>
      if (res.code === 200 && res.data) {
        formData.name = res.data.name
        formData.description = res.data.description || ''
        formData.area = res.data.area
        formData.budget = res.data.budget
        formData.style = res.data.style || ''
        formData.status = res.data.status
      } else {
        ElMessage.error(res.msg || '获取项目详情失败')
      }
    } catch (error) {
      console.error('获取项目详情失败:', error)
      ElMessage.error('获取项目详情失败')
    } finally {
      loading.value = false
    }
  }

  // 返回列表 / Back to list
  const handleBack = () => {
    router.push('/designer/designer-assistant/project/ProjectList')
  }

  /***
   * Submit form
   * 提交表单
   ***/
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        submitting.value = true
        try {
          if (isEdit.value) {
            // 更新项目 / Update project
            const res = (await updateDesignerProject(projectId.value, {
              id: projectId.value,
              name: formData.name,
              description: formData.description,
              area: formData.area,
              budget: formData.budget,
              style: formData.style,
              status: formData.status
            })) as unknown as ApiResponse<ProjectResponse>
            if (res.code === 200) {
              ElMessage.success('更新成功')
              router.push('/designer/designer-assistant/project/ProjectList')
            } else {
              ElMessage.error(res.msg || '更新失败')
            }
          } else {
            // 创建项目 / Create project
            const res = (await createDesignerProject({
              name: formData.name,
              description: formData.description,
              area: formData.area,
              budget: formData.budget,
              style: formData.style
            })) as unknown as ApiResponse<ProjectResponse>
            if (res.code === 200) {
              ElMessage.success('创建成功')
              router.push('/designer/designer-assistant/project/ProjectList')
            } else {
              ElMessage.error(res.msg || '创建失败')
            }
          }
        } catch (error) {
          console.error('提交失败:', error)
          ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
        } finally {
          submitting.value = false
        }
      }
    })
  }

  // 重置表单 / Reset form
  const handleReset = () => {
    if (isEdit.value) {
      getProjectDetail()
    } else {
      formRef.value?.resetFields()
    }
  }

  onMounted(() => {
    if (isEdit.value) {
      getProjectDetail()
    }
  })
</script>

<style scoped lang="scss">
  .project-create-page {
    padding: 16px;

    .card-header {
      display: flex;
      align-items: center;
      gap: 16px;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }
  }
</style>
