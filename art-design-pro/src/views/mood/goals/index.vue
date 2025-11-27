<template>
  <div class="goals-page">
    <!-- 操作栏 -->
    <ElCard shadow="never" class="mb-4">
      <div class="toolbar">
        <div class="filter-area">
          <ElRadioGroup v-model="statusFilter" @change="loadGoals">
            <ElRadioButton label="">全部</ElRadioButton>
            <ElRadioButton label="active">进行中</ElRadioButton>
            <ElRadioButton label="completed">已完成</ElRadioButton>
            <ElRadioButton label="paused">已暂停</ElRadioButton>
          </ElRadioGroup>
        </div>
        <ElButton type="primary" :icon="Plus" @click="openEditor()">新建目标</ElButton>
      </div>
    </ElCard>

    <!-- 目标列表 -->
    <div v-loading="loading" class="goals-list">
      <ElRow :gutter="16">
        <ElCol v-for="goal in goalList" :key="goal.id" :xs="24" :sm="12" :lg="8" class="mb-4">
          <ElCard shadow="hover" class="goal-card">
            <div class="goal-header">
              <div class="goal-icon">{{ getCategoryIcon(goal.category) }}</div>
              <div class="goal-info">
                <div class="goal-title">{{ goal.title }}</div>
                <div class="goal-category">{{ goal.category }}</div>
              </div>
              <ElDropdown trigger="click" @command="handleCommand($event, goal)">
                <ElButton :icon="MoreFilled" link />
                <template #dropdown>
                  <ElDropdownMenu>
                    <ElDropdownItem command="edit">编辑</ElDropdownItem>
                    <ElDropdownItem command="checkin" v-if="goal.status === 'active'"
                      >打卡</ElDropdownItem
                    >
                    <ElDropdownItem command="pause" v-if="goal.status === 'active'"
                      >暂停</ElDropdownItem
                    >
                    <ElDropdownItem command="resume" v-if="goal.status === 'paused'"
                      >恢复</ElDropdownItem
                    >
                    <ElDropdownItem command="complete" v-if="goal.status === 'active'"
                      >完成</ElDropdownItem
                    >
                    <ElDropdownItem command="delete" divided>删除</ElDropdownItem>
                  </ElDropdownMenu>
                </template>
              </ElDropdown>
            </div>

            <div class="goal-progress">
              <div class="progress-info">
                <span>进度</span>
                <span>{{ goal.current_value }} / {{ goal.target_value }} {{ goal.unit }}</span>
              </div>
              <ElProgress :percentage="getProgress(goal)" :stroke-width="8" :show-text="false" />
            </div>

            <div class="goal-stats">
              <div class="stat-item">
                <span class="stat-label">连续打卡</span>
                <span class="stat-value">{{ goal.streak_days || 0 }}天</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">总打卡</span>
                <span class="stat-value">{{ goal.total_checkins || 0 }}次</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">状态</span>
                <ElTag :type="getStatusType(goal.status)" size="small">{{
                  getStatusLabel(goal.status)
                }}</ElTag>
              </div>
            </div>

            <div class="goal-footer">
              <span class="goal-date"
                >{{ formatDate(goal.start_date) }} - {{ formatDate(goal.end_date) }}</span
              >
              <ElButton
                v-if="goal.status === 'active'"
                type="primary"
                size="small"
                @click="openCheckin(goal)"
              >
                今日打卡
              </ElButton>
            </div>
          </ElCard>
        </ElCol>
      </ElRow>

      <ElEmpty v-if="!loading && goalList.length === 0" description="还没有目标，创建一个吧" />
    </div>

    <!-- 新建/编辑目标弹窗 -->
    <ElDialog v-model="editorVisible" :title="editingGoal ? '编辑目标' : '新建目标'" width="600px">
      <ElForm :model="goalForm" :rules="goalRules" ref="goalFormRef" label-width="100px">
        <ElFormItem label="目标名称" prop="title">
          <ElInput v-model="goalForm.title" placeholder="如：每天运动30分钟" />
        </ElFormItem>
        <ElFormItem label="目标分类" prop="category">
          <ElSelect v-model="goalForm.category" placeholder="选择分类">
            <ElOption label="健康" value="健康" />
            <ElOption label="学习" value="学习" />
            <ElOption label="工作" value="工作" />
            <ElOption label="生活" value="生活" />
            <ElOption label="其他" value="其他" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="目标值" prop="target_value">
          <ElInputNumber v-model="goalForm.target_value" :min="1" />
          <ElInput
            v-model="goalForm.unit"
            placeholder="单位"
            style="width: 100px; margin-left: 10px"
          />
        </ElFormItem>
        <ElFormItem label="时间范围">
          <ElDatePicker
            v-model="goalForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </ElFormItem>
        <ElFormItem label="提醒时间">
          <ElTimePicker
            v-model="goalForm.reminder_time"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择提醒时间"
          />
        </ElFormItem>
        <ElFormItem label="目标描述">
          <ElInput
            v-model="goalForm.description"
            type="textarea"
            :rows="3"
            placeholder="描述一下你的目标"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="editorVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveGoal" :loading="saving">保存</ElButton>
      </template>
    </ElDialog>

    <!-- 打卡弹窗 -->
    <ElDialog v-model="checkinVisible" title="目标打卡" width="500px">
      <div v-if="checkinGoal" class="checkin-content">
        <div class="checkin-goal">
          <span class="goal-icon">{{ getCategoryIcon(checkinGoal.category) }}</span>
          <span class="goal-title">{{ checkinGoal.title }}</span>
        </div>
        <ElForm :model="checkinForm" label-width="80px">
          <ElFormItem label="完成值">
            <ElInputNumber v-model="checkinForm.value" :min="0" />
            <span class="unit">{{ checkinGoal.unit }}</span>
          </ElFormItem>
          <ElFormItem label="备注">
            <ElInput
              v-model="checkinForm.note"
              type="textarea"
              :rows="2"
              placeholder="记录一下今天的感受"
            />
          </ElFormItem>
        </ElForm>
      </div>
      <template #footer>
        <ElButton @click="checkinVisible = false">取消</ElButton>
        <ElButton type="primary" @click="submitCheckin" :loading="checkinLoading"
          >完成打卡</ElButton
        >
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { Plus, MoreFilled } from '@element-plus/icons-vue'
  import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
  import { goalApi, checkinApi, type Goal } from '@/api/mood'
  import dayjs from 'dayjs'

  defineOptions({ name: 'Goals' })

  // 筛选
  const statusFilter = ref('')

  // 数据
  const goalList = ref<Goal[]>([])
  const loading = ref(false)

  // 编辑器
  const editorVisible = ref(false)
  const editingGoal = ref<Goal | null>(null)
  const saving = ref(false)
  const goalFormRef = ref<FormInstance>()
  const goalForm = reactive({
    title: '',
    category: '',
    target_value: 1,
    unit: '次',
    dateRange: null as [string, string] | null,
    reminder_time: '',
    description: ''
  })

  const goalRules: FormRules = {
    title: [{ required: true, message: '请输入目标名称', trigger: 'blur' }],
    category: [{ required: true, message: '请选择分类', trigger: 'change' }],
    target_value: [{ required: true, message: '请输入目标值', trigger: 'blur' }]
  }

  // 打卡
  const checkinVisible = ref(false)
  const checkinGoal = ref<Goal | null>(null)
  const checkinLoading = ref(false)
  const checkinForm = reactive({
    value: 1,
    note: ''
  })

  // 方法
  const formatDate = (date: string) => dayjs(date).format('MM/DD')

  const getCategoryIcon = (category: string) => {
    const icons: Record<string, string> = {
      健康: '💪',
      学习: '📚',
      工作: '💼',
      生活: '🏠',
      其他: '🎯'
    }
    return icons[category] || '🎯'
  }

  const getProgress = (goal: Goal) => {
    if (!goal.target_value) return 0
    return Math.min(100, Math.round((goal.current_value / goal.target_value) * 100))
  }

  const getStatusType = (status: string) => {
    const types: Record<string, string> = {
      active: 'success',
      completed: 'primary',
      paused: 'warning',
      abandoned: 'danger'
    }
    return types[status] || 'info'
  }

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      active: '进行中',
      completed: '已完成',
      paused: '已暂停',
      abandoned: '已放弃'
    }
    return labels[status] || status
  }

  const loadGoals = async () => {
    loading.value = true
    try {
      const params: any = {}
      if (statusFilter.value) params.status = statusFilter.value

      const res = await goalApi.getList(params)
      const data = (res as any)?.data || res
      goalList.value = data?.records || data || []
    } catch (error) {
      console.error('加载目标失败:', error)
    } finally {
      loading.value = false
    }
  }

  const openEditor = (goal?: Goal) => {
    editingGoal.value = goal || null
    if (goal) {
      goalForm.title = goal.title
      goalForm.category = goal.category
      goalForm.target_value = goal.target_value
      goalForm.unit = goal.unit
      goalForm.dateRange = [goal.start_date, goal.end_date]
      goalForm.reminder_time = goal.reminder_time || ''
      goalForm.description = goal.description || ''
    } else {
      goalForm.title = ''
      goalForm.category = ''
      goalForm.target_value = 1
      goalForm.unit = '次'
      goalForm.dateRange = null
      goalForm.reminder_time = ''
      goalForm.description = ''
    }
    editorVisible.value = true
  }

  const saveGoal = async () => {
    if (!goalFormRef.value) return
    await goalFormRef.value.validate()

    saving.value = true
    try {
      const data = {
        title: goalForm.title,
        category: goalForm.category,
        target_value: goalForm.target_value,
        unit: goalForm.unit,
        start_date: goalForm.dateRange?.[0] || dayjs().format('YYYY-MM-DD'),
        end_date: goalForm.dateRange?.[1] || dayjs().add(30, 'day').format('YYYY-MM-DD'),
        reminder_time: goalForm.reminder_time,
        description: goalForm.description
      }

      if (editingGoal.value) {
        await goalApi.update({ id: editingGoal.value.id, ...data })
        ElMessage.success('目标已更新')
      } else {
        await goalApi.create(data)
        ElMessage.success('目标已创建')
      }

      editorVisible.value = false
      loadGoals()
    } catch (error: any) {
      ElMessage.error(error.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  const handleCommand = async (command: string, goal: Goal) => {
    switch (command) {
      case 'edit':
        openEditor(goal)
        break
      case 'checkin':
        openCheckin(goal)
        break
      case 'pause':
        await updateGoalStatus(goal, 'paused')
        break
      case 'resume':
        await updateGoalStatus(goal, 'active')
        break
      case 'complete':
        await updateGoalStatus(goal, 'completed')
        break
      case 'delete':
        await deleteGoal(goal)
        break
    }
  }

  const updateGoalStatus = async (goal: Goal, status: Goal['status']) => {
    try {
      await goalApi.update({ id: goal.id, status })
      ElMessage.success('状态已更新')
      loadGoals()
    } catch (error: any) {
      ElMessage.error(error.message || '更新失败')
    }
  }

  const deleteGoal = async (goal: Goal) => {
    try {
      await ElMessageBox.confirm('确定要删除这个目标吗？', '删除确认', { type: 'warning' })
      await goalApi.delete(goal.id)
      ElMessage.success('删除成功')
      loadGoals()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  const openCheckin = (goal: Goal) => {
    checkinGoal.value = goal
    checkinForm.value = 1
    checkinForm.note = ''
    checkinVisible.value = true
  }

  const submitCheckin = async () => {
    if (!checkinGoal.value) return

    checkinLoading.value = true
    try {
      await checkinApi.create({
        goal_id: checkinGoal.value.id,
        value: checkinForm.value,
        note: checkinForm.note
      })
      ElMessage.success('打卡成功！')
      checkinVisible.value = false
      loadGoals()
    } catch (error: any) {
      ElMessage.error(error.message || '打卡失败')
    } finally {
      checkinLoading.value = false
    }
  }

  onMounted(() => {
    loadGoals()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/styles/variables.scss' as *;

  .goals-page {
    --card-spacing: 20px;
  }

  // 统一卡片样式
  :deep(.el-card) {
    margin-bottom: var(--card-spacing);
    background: var(--art-main-bg-color);
    border-radius: calc(var(--custom-radius) + 4px) !important;
    border: none;
    box-shadow: var(--art-root-card-box-shadow);
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: var(--card-spacing);
  }

  .goal-card {
    height: 100%;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);
      box-shadow: var(--art-box-shadow);
    }

    :deep(.el-card__body) {
      padding: 20px;
    }

    .goal-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 14px;

      .goal-title {
        display: flex;
        align-items: center;
        gap: 10px;

        .goal-icon {
          font-size: 28px;
        }

        .goal-name {
          font-size: 16px;
          font-weight: 600;
          color: var(--art-gray-900);
        }
      }
    }

    .goal-meta {
      display: flex;
      gap: 8px;
      margin-bottom: 16px;
    }

    .goal-progress {
      margin-bottom: 16px;

      .progress-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;

        .progress-text {
          font-size: 14px;
          color: var(--art-gray-600);
        }

        .progress-value {
          font-size: 14px;
          font-weight: 600;
          color: rgb(var(--art-primary));
        }
      }
    }

    .goal-stats {
      display: flex;
      justify-content: space-between;
      padding-top: 14px;
      border-top: 1px solid var(--art-border-color);

      .stat-item {
        text-align: center;

        .stat-value {
          font-size: 18px;
          font-weight: 600;
          color: var(--art-gray-900);
        }

        .stat-label {
          font-size: 12px;
          color: var(--art-gray-500);
          margin-top: 2px;
        }
      }
    }

    .goal-actions {
      display: flex;
      gap: 8px;
      margin-top: 16px;
    }
  }

  .goal-form {
    .reminder-row {
      display: flex;
      gap: 16px;
      align-items: center;
    }
  }

  .checkin-form {
    .current-progress {
      text-align: center;
      padding: 24px;
      background: var(--art-gray-100);
      border-radius: 12px;
      margin-bottom: 20px;

      .progress-label {
        font-size: 14px;
        color: var(--art-gray-500);
        margin-bottom: 8px;
      }

      .progress-value {
        font-size: 36px;
        font-weight: 700;
        color: rgb(var(--art-primary));
      }

      .progress-unit {
        font-size: 14px;
        color: var(--art-gray-600);
        margin-left: 4px;
      }
    }
  }

  // 响应式
  @media screen and (max-width: $device-phone) {
    .goals-page {
      --card-spacing: 15px;
    }
  }
</style>
