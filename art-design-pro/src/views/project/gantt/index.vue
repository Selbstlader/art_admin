<template>
  <div class="gantt-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <template #header>
        <div class="gantt-header">
          <span class="gantt-title">项目甘特图 - {{ projectName }}</span>
          <ElSpace wrap>
            <ElButton @click="handleBack" v-ripple>返回</ElButton>
            <ElButton @click="handleRefresh" v-ripple>刷新</ElButton>
            <ElButton type="primary" @click="handleAddTask" v-ripple>添加任务</ElButton>
          </ElSpace>
        </div>
      </template>

      <GanttChart
        v-if="isMounted"
        ref="ganttChartRef"
        :project-id="projectId"
        @task-click="handleTaskClick"
        @data-loaded="handleDataLoaded"
      />
    </ElCard>

    <!-- 任务表单弹窗 -->
    <TaskForm
      v-model:visible="taskFormVisible"
      title="添加任务"
      :project-id="projectId"
      @success="handleTaskCreated"
    />

    <!-- 任务详情抽屉 -->
    <TaskDrawer
      ref="taskDrawerRef"
      v-model:visible="taskDrawerVisible"
      :task-id="selectedTaskId"
      @deleted="handleTaskDeleted"
      @updated="handleTaskUpdated"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onBeforeUnmount } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { projectApi } from '@/api/project'
  import GanttChart from './modules/GanttChart.vue'
  import TaskForm from './modules/TaskForm.vue'
  import TaskDrawer from './modules/TaskDrawer.vue'

  defineOptions({ name: 'ProjectGantt' })

  const route = useRoute()
  const router = useRouter()

  const ganttChartRef = ref()
  const taskDrawerRef = ref()
  const projectId = ref(Number(route.params.id))
  const projectName = ref('')
  const taskFormVisible = ref(false)
  const taskDrawerVisible = ref(false)
  const selectedTaskId = ref<number | null>(null)
  const isMounted = ref(false)

  // 加载项目信息
  const loadProjectInfo = async () => {
    try {
      const data = (await projectApi.getProjectDetail(projectId.value)) as { name: string }
      projectName.value = data.name
    } catch (error) {
      console.error('加载项目信息失败:', error)
      ElMessage.error('加载项目信息失败')
    }
  }

  // 返回
  const handleBack = () => {
    isMounted.value = false
    setTimeout(() => router.back(), 100)
  }

  // 刷新页面数据
  const handleRefresh = () => {
    if (ganttChartRef.value) {
      ganttChartRef.value.refresh()
    }
  }

  // 添加任务
  const handleAddTask = () => {
    taskFormVisible.value = true
  }

  // 任务点击事件
  const handleTaskClick = async (taskId: number) => {
    selectedTaskId.value = taskId
    taskDrawerVisible.value = true
    // 等待抽屉打开后加载详情
    await new Promise((resolve) => setTimeout(resolve, 100))
    taskDrawerRef.value?.loadTaskDetail()
  }

  // 任务操作成功后的回调
  const handleTaskCreated = () => {
    console.log('任务创建成功，刷新甘特图')
    ganttChartRef.value?.refresh()
  }

  const handleTaskDeleted = () => {
    console.log('任务删除成功，刷新甘特图')
    ganttChartRef.value?.refresh()
  }

  const handleTaskUpdated = () => {
    console.log('任务更新成功，刷新甘特图')
    ganttChartRef.value?.refresh()
  }

  // 甘特图数据加载完成
  const handleDataLoaded = () => {}

  // 页面挂载时加载数据
  onMounted(() => {
    isMounted.value = true
    loadProjectInfo()
  })

  // 组件卸载时清理
  onBeforeUnmount(() => {
    isMounted.value = false
  })
</script>

<style scoped lang="scss">
  .gantt-page {
    padding-bottom: 15px;

    :deep(.el-card__header) {
      padding: 16px 20px;
    }

    :deep(.el-card__body) {
      padding: 0;
    }
  }

  .gantt-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .gantt-title {
      font-size: 16px;
      font-weight: 500;
      color: var(--el-text-color-primary);
    }
  }
</style>
