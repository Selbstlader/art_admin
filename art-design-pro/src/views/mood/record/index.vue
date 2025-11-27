<template>
  <div class="mood-record-page">
    <!-- 记录情绪表单 -->
    <ElCard shadow="never" class="mb-4">
      <template #header>
        <span class="card-title">记录当前情绪</span>
      </template>

      <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
        <!-- 情绪类型选择 -->
        <ElFormItem label="情绪类型" prop="mood_type">
          <div class="mood-type-grid">
            <div
              v-for="(config, type) in MoodTypeConfig"
              :key="type"
              class="mood-type-item"
              :class="{ active: form.mood_type === type }"
              @click="form.mood_type = type as MoodType"
            >
              <span class="mood-icon">{{ config.icon }}</span>
              <span class="mood-label">{{ config.label }}</span>
            </div>
          </div>
        </ElFormItem>

        <!-- 情绪强度 -->
        <ElFormItem label="情绪强度" prop="intensity">
          <div class="intensity-slider">
            <ElSlider
              v-model="form.intensity"
              :min="1"
              :max="10"
              :marks="intensityMarks"
              show-stops
            />
            <div class="intensity-labels">
              <span>轻微</span>
              <span>中等</span>
              <span>强烈</span>
            </div>
          </div>
        </ElFormItem>

        <!-- 触发因素 -->
        <ElFormItem label="触发因素">
          <div class="tag-input-wrapper">
            <ElTag
              v-for="tag in form.triggers"
              :key="tag"
              closable
              @close="removeTag('triggers', tag)"
              class="mr-2 mb-2"
            >
              {{ tag }}
            </ElTag>
            <ElInput
              v-if="triggerInputVisible"
              ref="triggerInputRef"
              v-model="triggerInputValue"
              size="small"
              style="width: 100px"
              @keyup.enter="addTag('triggers')"
              @blur="addTag('triggers')"
            />
            <ElButton v-else size="small" @click="showTagInput('triggers')"> + 添加 </ElButton>
          </div>
          <div class="quick-tags">
            <ElTag
              v-for="tag in quickTriggers"
              :key="tag"
              :type="form.triggers.includes(tag) ? 'primary' : 'info'"
              effect="plain"
              class="mr-2 mb-2 cursor-pointer"
              @click="toggleQuickTag('triggers', tag)"
            >
              {{ tag }}
            </ElTag>
          </div>
        </ElFormItem>

        <!-- 相关活动 -->
        <ElFormItem label="相关活动">
          <div class="tag-input-wrapper">
            <ElTag
              v-for="tag in form.activities"
              :key="tag"
              closable
              type="success"
              @close="removeTag('activities', tag)"
              class="mr-2 mb-2"
            >
              {{ tag }}
            </ElTag>
            <ElInput
              v-if="activityInputVisible"
              ref="activityInputRef"
              v-model="activityInputValue"
              size="small"
              style="width: 100px"
              @keyup.enter="addTag('activities')"
              @blur="addTag('activities')"
            />
            <ElButton v-else size="small" @click="showTagInput('activities')"> + 添加 </ElButton>
          </div>
          <div class="quick-tags">
            <ElTag
              v-for="tag in quickActivities"
              :key="tag"
              :type="form.activities.includes(tag) ? 'success' : 'info'"
              effect="plain"
              class="mr-2 mb-2 cursor-pointer"
              @click="toggleQuickTag('activities', tag)"
            >
              {{ tag }}
            </ElTag>
          </div>
        </ElFormItem>

        <!-- 位置 -->
        <ElFormItem label="当前位置">
          <ElInput v-model="form.location" placeholder="如：办公室、家里" />
        </ElFormItem>

        <!-- 备注 -->
        <ElFormItem label="备注说明">
          <ElInput
            v-model="form.note"
            type="textarea"
            :rows="3"
            placeholder="记录一下此刻的想法..."
            maxlength="500"
            show-word-limit
          />
        </ElFormItem>

        <ElFormItem>
          <ElButton type="primary" @click="handleSubmit" :loading="submitting"> 保存记录 </ElButton>
          <ElButton @click="resetForm">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 今日记录列表 -->
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">今日情绪记录</span>
          <ElButton link type="primary" @click="$router.push('/mood/calendar')">
            查看更多
          </ElButton>
        </div>
      </template>

      <div v-loading="loading">
        <div v-if="todayRecords.length > 0" class="record-list">
          <div v-for="record in todayRecords" :key="record.id" class="record-item">
            <div class="record-mood">
              <span class="mood-emoji">{{ MoodTypeConfig[record.mood_type]?.icon }}</span>
              <div class="mood-info">
                <span class="mood-name">{{ MoodTypeConfig[record.mood_type]?.label }}</span>
                <span class="mood-time">{{ formatTime(record.created_at) }}</span>
              </div>
            </div>
            <div class="record-intensity">
              <ElProgress
                :percentage="record.intensity * 10"
                :stroke-width="6"
                :show-text="false"
                :color="getIntensityColor(record.intensity)"
              />
              <span class="intensity-value">{{ record.intensity }}</span>
            </div>
            <div class="record-tags" v-if="record.triggers?.length || record.activities?.length">
              <ElTag v-for="tag in record.triggers" :key="'t-' + tag" size="small" class="mr-1">
                {{ tag }}
              </ElTag>
              <ElTag
                v-for="tag in record.activities"
                :key="'a-' + tag"
                size="small"
                type="success"
                class="mr-1"
              >
                {{ tag }}
              </ElTag>
            </div>
            <div class="record-note" v-if="record.note">
              {{ record.note }}
            </div>
            <div class="record-actions">
              <ElButton link type="primary" size="small" @click="editRecord(record)">
                编辑
              </ElButton>
              <ElButton link type="danger" size="small" @click="deleteRecord(record)">
                删除
              </ElButton>
            </div>
          </div>
        </div>
        <ElEmpty v-else description="今天还没有记录，快来记录一下吧" />
      </div>
    </ElCard>

    <!-- 编辑弹窗 -->
    <ElDialog v-model="editDialogVisible" title="编辑情绪记录" width="600px">
      <ElForm :model="editForm" label-width="100px">
        <ElFormItem label="情绪类型">
          <div class="mood-type-grid small">
            <div
              v-for="(config, type) in MoodTypeConfig"
              :key="type"
              class="mood-type-item"
              :class="{ active: editForm.mood_type === type }"
              @click="editForm.mood_type = type as MoodType"
            >
              <span class="mood-icon">{{ config.icon }}</span>
              <span class="mood-label">{{ config.label }}</span>
            </div>
          </div>
        </ElFormItem>
        <ElFormItem label="情绪强度">
          <ElSlider v-model="editForm.intensity" :min="1" :max="10" show-stops />
        </ElFormItem>
        <ElFormItem label="备注说明">
          <ElInput v-model="editForm.note" type="textarea" :rows="3" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="editDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleEditSubmit" :loading="editSubmitting">
          保存
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, nextTick } from 'vue'
  import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
  import { moodRecordApi, MoodTypeConfig, type MoodType, type MoodRecord } from '@/api/mood'
  import dayjs from 'dayjs'

  defineOptions({ name: 'MoodRecord' })

  const formRef = ref<FormInstance>()
  const triggerInputRef = ref()
  const activityInputRef = ref()

  // 表单数据
  const form = reactive({
    mood_type: '' as MoodType | '',
    intensity: 5,
    triggers: [] as string[],
    activities: [] as string[],
    location: '',
    note: ''
  })

  const rules: FormRules = {
    mood_type: [{ required: true, message: '请选择情绪类型', trigger: 'change' }],
    intensity: [{ required: true, message: '请设置情绪强度', trigger: 'change' }]
  }

  // 强度标记
  const intensityMarks = {
    1: '1',
    3: '3',
    5: '5',
    7: '7',
    10: '10'
  }

  // 快捷标签
  const quickTriggers = ['工作压力', '人际关系', '健康问题', '经济压力', '家庭事务', '学习考试']
  const quickActivities = ['工作', '运动', '休息', '社交', '娱乐', '学习', '家务', '用餐']

  // 标签输入状态
  const triggerInputVisible = ref(false)
  const triggerInputValue = ref('')
  const activityInputVisible = ref(false)
  const activityInputValue = ref('')

  // 提交状态
  const submitting = ref(false)
  const loading = ref(false)

  // 今日记录
  const todayRecords = ref<MoodRecord[]>([])

  // 编辑相关
  const editDialogVisible = ref(false)
  const editSubmitting = ref(false)
  const editForm = reactive({
    id: 0,
    mood_type: '' as MoodType | '',
    intensity: 5,
    note: ''
  })

  // 方法
  const showTagInput = (type: 'triggers' | 'activities') => {
    if (type === 'triggers') {
      triggerInputVisible.value = true
      nextTick(() => triggerInputRef.value?.focus())
    } else {
      activityInputVisible.value = true
      nextTick(() => activityInputRef.value?.focus())
    }
  }

  const addTag = (type: 'triggers' | 'activities') => {
    const value = type === 'triggers' ? triggerInputValue.value : activityInputValue.value
    if (value && !form[type].includes(value)) {
      form[type].push(value)
    }
    if (type === 'triggers') {
      triggerInputVisible.value = false
      triggerInputValue.value = ''
    } else {
      activityInputVisible.value = false
      activityInputValue.value = ''
    }
  }

  const removeTag = (type: 'triggers' | 'activities', tag: string) => {
    const index = form[type].indexOf(tag)
    if (index > -1) {
      form[type].splice(index, 1)
    }
  }

  const toggleQuickTag = (type: 'triggers' | 'activities', tag: string) => {
    const index = form[type].indexOf(tag)
    if (index > -1) {
      form[type].splice(index, 1)
    } else {
      form[type].push(tag)
    }
  }

  const getIntensityColor = (intensity: number) => {
    if (intensity <= 3) return '#67C23A'
    if (intensity <= 6) return '#E6A23C'
    return '#F56C6C'
  }

  const formatTime = (time: string) => {
    return dayjs(time).format('HH:mm')
  }

  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitting.value = true

      await moodRecordApi.create({
        mood_type: form.mood_type as MoodType,
        intensity: form.intensity,
        triggers: form.triggers,
        activities: form.activities,
        location: form.location,
        note: form.note
      })

      ElMessage.success('情绪记录保存成功')
      resetForm()
      loadTodayRecords()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '保存失败')
      }
    } finally {
      submitting.value = false
    }
  }

  const resetForm = () => {
    form.mood_type = ''
    form.intensity = 5
    form.triggers = []
    form.activities = []
    form.location = ''
    form.note = ''
    formRef.value?.clearValidate()
  }

  const loadTodayRecords = async () => {
    loading.value = true
    try {
      const today = dayjs().format('YYYY-MM-DD')
      const res = await moodRecordApi.getByDate(today)
      todayRecords.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载今日记录失败:', error)
    } finally {
      loading.value = false
    }
  }

  const editRecord = (record: MoodRecord) => {
    editForm.id = record.id
    editForm.mood_type = record.mood_type
    editForm.intensity = record.intensity
    editForm.note = record.note
    editDialogVisible.value = true
  }

  const handleEditSubmit = async () => {
    try {
      editSubmitting.value = true
      await moodRecordApi.update({
        id: editForm.id,
        mood_type: editForm.mood_type as MoodType,
        intensity: editForm.intensity,
        note: editForm.note
      })
      ElMessage.success('更新成功')
      editDialogVisible.value = false
      loadTodayRecords()
    } catch (error: any) {
      ElMessage.error(error.message || '更新失败')
    } finally {
      editSubmitting.value = false
    }
  }

  const deleteRecord = async (record: MoodRecord) => {
    try {
      await ElMessageBox.confirm('确定要删除这条记录吗？', '删除确认', {
        type: 'warning'
      })
      await moodRecordApi.delete(record.id)
      ElMessage.success('删除成功')
      loadTodayRecords()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  onMounted(() => {
    loadTodayRecords()
  })
</script>

<style scoped lang="scss">
  .mood-record-page {
    padding: 16px;
  }

  .mb-4 {
    margin-bottom: 16px;
  }

  .mr-1 {
    margin-right: 4px;
  }

  .mr-2 {
    margin-right: 8px;
  }

  .mb-2 {
    margin-bottom: 8px;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .mood-type-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;

    @media (max-width: 768px) {
      grid-template-columns: repeat(2, 1fr);
    }

    &.small {
      grid-template-columns: repeat(4, 1fr);
      gap: 8px;

      .mood-type-item {
        padding: 8px;

        .mood-icon {
          font-size: 20px;
        }

        .mood-label {
          font-size: 12px;
        }
      }
    }

    .mood-type-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 16px;
      border: 2px solid #ebeef5;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #409eff;
        background: #f0f9ff;
      }

      &.active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .mood-icon {
        font-size: 32px;
        margin-bottom: 8px;
      }

      .mood-label {
        font-size: 14px;
        color: #606266;
      }
    }
  }

  .intensity-slider {
    width: 100%;

    .intensity-labels {
      display: flex;
      justify-content: space-between;
      margin-top: 8px;
      font-size: 12px;
      color: #909399;
    }
  }

  .tag-input-wrapper {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    margin-bottom: 8px;
  }

  .quick-tags {
    display: flex;
    flex-wrap: wrap;

    .cursor-pointer {
      cursor: pointer;
    }
  }

  .record-list {
    .record-item {
      padding: 16px;
      border-radius: 8px;
      background: #fafafa;
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }

      .record-mood {
        display: flex;
        align-items: center;
        margin-bottom: 12px;

        .mood-emoji {
          font-size: 32px;
          margin-right: 12px;
        }

        .mood-info {
          display: flex;
          flex-direction: column;

          .mood-name {
            font-size: 16px;
            font-weight: 500;
            color: #303133;
          }

          .mood-time {
            font-size: 12px;
            color: #909399;
          }
        }
      }

      .record-intensity {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;

        .el-progress {
          flex: 1;
        }

        .intensity-value {
          font-size: 14px;
          font-weight: 600;
          color: #303133;
          width: 20px;
        }
      }

      .record-tags {
        margin-bottom: 8px;
      }

      .record-note {
        font-size: 14px;
        color: #606266;
        margin-bottom: 8px;
        padding: 8px;
        background: #fff;
        border-radius: 4px;
      }

      .record-actions {
        display: flex;
        justify-content: flex-end;
        gap: 8px;
      }
    }
  }
</style>
