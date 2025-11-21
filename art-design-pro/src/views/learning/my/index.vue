<template>
  <div class="learning-page art-full-height">
    <ElCard shadow="never">
      <ElPageHeader title="AI 学习系统" content="智能生成个性化学习教材" @back="$router.back()">
        <template #extra>
          <ElButton type="primary" @click="$router.push('/learning/list')">
            <ElIcon><Document /></ElIcon>
            我的教材
          </ElButton>
        </template>
      </ElPageHeader>

      <ElDivider />

      <!-- 步骤条 -->
      <ElSteps :active="currentStep" align-center style="margin-bottom: 32px">
        <ElStep title="选择学科" />
        <ElStep title="填写信息" />
      </ElSteps>

      <!-- 步骤1: 选择学科 -->
      <div v-if="currentStep === 0" class="step-content">
        <div class="subject-grid">
          <ElCard
            v-for="subject in subjects"
            :key="subject.id"
            shadow="hover"
            :class="['subject-card', { active: selectedSubject?.id === subject.id }]"
            @click="selectSubject(subject)"
          >
            <div class="subject-icon">{{ subject.icon }}</div>
            <div class="subject-name">{{ subject.name }}</div>
            <div class="subject-desc">{{ subject.description }}</div>
            <div class="subject-grades">
              <ElTag v-for="grade in subject.grade_levels" :key="grade" size="small">
                {{ grade }}
              </ElTag>
            </div>
          </ElCard>
        </div>
      </div>

      <!-- 步骤2: 填写表单 -->
      <div v-else-if="currentStep === 1" class="step-content">
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-width="100px"
          style="max-width: 600px; margin: 0 auto"
        >
          <ElFormItem label="学科" prop="subject_id">
            <ElSelect v-model="formData.subject_id" placeholder="请选择学科" style="width: 100%">
              <ElOption
                v-for="subject in subjects"
                :key="subject.id"
                :label="`${subject.icon} ${subject.name}`"
                :value="subject.id"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="学段" prop="gradeLevel">
            <ElRadioGroup v-model="formData.gradeLevel" @change="handleGradeLevelChange">
              <ElRadioButton label="小学" />
              <ElRadioButton label="初中" />
              <ElRadioButton label="高中" />
              <ElRadioButton label="大学" />
            </ElRadioGroup>
          </ElFormItem>

          <ElFormItem label="年级" prop="grade">
            <ElSelect
              v-model="formData.grade"
              placeholder="请先选择学段"
              :disabled="!formData.gradeLevel"
              style="width: 100%"
            >
              <ElOption
                v-for="grade in availableGrades"
                :key="grade"
                :label="grade"
                :value="grade"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="学习题材" prop="topic">
            <ElInput
              v-model="formData.topic"
              placeholder="例如：一元二次方程、勾股定理"
              maxlength="50"
              show-word-limit
            />
          </ElFormItem>

          <ElFormItem label="难度等级" prop="difficulty">
            <ElRadioGroup v-model="formData.difficulty">
              <ElRadioButton :value="1">⭐ 基础</ElRadioButton>
              <ElRadioButton :value="2">⭐⭐ 进阶</ElRadioButton>
              <ElRadioButton :value="3">⭐⭐⭐ 高级</ElRadioButton>
            </ElRadioGroup>
          </ElFormItem>

          <ElFormItem>
            <ElSpace>
              <ElButton type="primary" :loading="generating" @click="handleGenerate">
                <ElIcon v-if="!generating"><MagicStick /></ElIcon>
                {{ generating ? '提交中...' : '开始生成教材' }}
              </ElButton>
              <ElButton @click="currentStep = 0">返回</ElButton>
            </ElSpace>
          </ElFormItem>
        </ElForm>

        <ElAlert
          title="异步生成提示"
          description="提交后将在后台异步生成教材，您可以继续创建其他任务，生成进度可在顶部栏查看"
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 20px"
        />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import { ElMessage, type FormInstance } from 'element-plus'
  import { Document, MagicStick } from '@element-plus/icons-vue'
  import { subjectApi, learningMaterialApi, type Subject } from '@/api/learning'

  const currentStep = ref(0)
  const subjects = ref<Subject[]>([])
  const selectedSubject = ref<Subject>()
  const generating = ref(false)
  const formRef = ref<FormInstance>()

  const formData = reactive({
    subject_id: undefined as number | undefined,
    gradeLevel: '',
    grade: '',
    topic: '',
    difficulty: 2
  })

  const rules = {
    subject_id: [{ required: true, message: '请选择学科', trigger: 'change' }],
    gradeLevel: [{ required: true, message: '请选择学段', trigger: 'change' }],
    grade: [{ required: true, message: '请选择年级', trigger: 'change' }],
    topic: [{ required: true, message: '请输入学习题材', trigger: 'blur' }],
    difficulty: [{ required: true, message: '请选择难度等级', trigger: 'change' }]
  }

  // 年级选项映射
  const gradeOptions: Record<string, string[]> = {
    小学: ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级'],
    初中: ['初一', '初二', '初三'],
    高中: ['高一', '高二', '高三'],
    大学: ['大一', '大二', '大三', '大四']
  }

  // 可用的年级列表
  const availableGrades = computed(() => {
    return formData.gradeLevel ? gradeOptions[formData.gradeLevel] || [] : []
  })

  // 学段切换时重置年级
  const handleGradeLevelChange = () => {
    formData.grade = ''
  }

  const loadSubjects = async () => {
    try {
      const res: any = await subjectApi.getSubjectList()
      if (res) {
        subjects.value = res
      }
    } catch (error) {
      console.error('获取学科列表失败:', error)
    }
  }

  const selectSubject = (subject: Subject) => {
    selectedSubject.value = subject
    formData.subject_id = subject.id
    currentStep.value = 1
  }

  const handleGenerate = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (!valid) return
      try {
        generating.value = true
        const res: any = await learningMaterialApi.generateMaterial(formData)
        if (res) {
          ElMessage.success('生成任务已创建，请在顶部栏查看进度')
          // 重置表单
          resetForm()
        } else {
          ElMessage.error(res.message || '创建任务失败')
        }
      } catch (error: any) {
        ElMessage.error(error.message || '创建任务失败，请重试')
      } finally {
        generating.value = false
      }
    })
  }

  const resetForm = () => {
    currentStep.value = 0
    selectedSubject.value = undefined
    formData.subject_id = undefined
    formData.gradeLevel = ''
    formData.grade = ''
    formData.topic = ''
    formData.difficulty = 2
  }

  onMounted(() => {
    loadSubjects()
  })
</script>

<style scoped lang="scss">
  .learning-page {
    .subject-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 16px;
      padding: 20px 0;

      .subject-card {
        cursor: pointer;
        text-align: center;
        transition: all 0.3s;

        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }

        &.active {
          border-color: var(--el-color-primary);
          background: var(--el-color-primary-light-9);
        }

        .subject-icon {
          font-size: 48px;
          margin-bottom: 12px;
        }

        .subject-name {
          font-size: 18px;
          font-weight: 600;
          margin-bottom: 8px;
        }

        .subject-desc {
          font-size: 14px;
          color: var(--el-text-color-secondary);
          margin-bottom: 12px;
          min-height: 40px;
        }

        .subject-grades {
          display: flex;
          justify-content: center;
          gap: 4px;
          flex-wrap: wrap;
        }
      }
    }

    .step-content {
      min-height: 400px;
      padding: 20px;
    }

    .material-detail {
      h2 {
        margin-bottom: 16px;
      }

      :deep(h3) {
        margin: 20px 0 10px;
        color: var(--el-color-primary);
      }

      :deep(p) {
        line-height: 1.8;
        margin: 10px 0;
      }

      :deep(ul) {
        padding-left: 20px;
        li {
          margin: 8px 0;
        }
      }

      :deep(hr) {
        margin: 20px 0;
        border: none;
        border-top: 1px solid var(--el-border-color);
      }
    }
  }
</style>
