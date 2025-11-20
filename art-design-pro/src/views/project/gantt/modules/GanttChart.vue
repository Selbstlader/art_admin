<template>
  <div ref="ganttContainer" class="gantt-chart-container"></div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { gantt } from 'dhtmlx-gantt'
  import 'dhtmlx-gantt/codebase/dhtmlxgantt.css'
  import { taskApi } from '@/api/project'

  interface GanttTask {
    id: number
    text: string
    start_date: string
    end_date: string
    duration: number
    progress: number
    parent: number
    type: string
    priority: number
    status: number
    assigneeId: number
    assigneeName: string
    isOverdue: boolean
    open: boolean
  }

  interface GanttDependency {
    id: number
    projectId: number
    predecessorId: number
    predecessorName: string
    successorId: number
    successorName: string
    dependencyType: string
    lagDays: number
    createdAt: string
  }

  interface GanttResponse {
    tasks: GanttTask[]
    dependencies: GanttDependency[]
  }

  const props = defineProps<{
    projectId: number
  }>()

  const emit = defineEmits<{
    (e: 'taskClick', taskId: number): void
    (e: 'dataLoaded'): void
  }>()

  const ganttContainer = ref<HTMLElement>()
  let isUpdating = false
  let isInitialized = false

  // 初始化甘特图
  const initGantt = () => {
    try {
      console.log('开始初始化甘特图...')

      // 如果已经初始化，先清理
      if (gantt.$container) {
        console.log('清理现有甘特图实例')
        gantt.clearAll()
      }

      gantt.i18n.setLocale('cn')

      gantt.config.columns = [
        { name: 'text', label: '任务名称', tree: true, width: 120 },
        { name: 'start_date', label: '开始时间', align: 'center', width: 120 },
        { name: 'duration', label: '工期', align: 'center', width: 60 },
        {
          name: 'progress',
          label: '进度',
          align: 'center',
          width: 80,
          template: (task: any) => Math.round(task.progress * 100) + '%'
        },
        { name: 'assigneeName', label: '负责人', align: 'center', width: 80 }
      ]

      gantt.config.date_format = '%Y-%m-%d %H:%i:%s'
      gantt.config.xml_date = '%Y-%m-%d %H:%i:%s'
      gantt.config.scale_unit = 'day'
      gantt.config.date_grid = '%Y-%m-%d'
      gantt.config.scale_height = 60

      gantt.config.drag_links = true
      gantt.config.drag_progress = true
      gantt.config.drag_resize = true
      gantt.config.drag_move = true

      gantt.config.auto_scheduling = true
      gantt.config.auto_scheduling_strict = true

      // 禁用自带的lightbox弹窗
      gantt.config.readonly = false
      gantt.config.details_on_create = false
      gantt.config.details_on_dblclick = false

      gantt.config.work_time = true
      gantt.config.skip_off_time = true

      gantt.templates.task_class = (start: any, end: any, task: any) => {
        if (task.type === 'milestone') return 'gantt_milestone'
        if (task.isOverdue) return 'gantt_overdue'
        return ''
      }

      gantt.templates.task_text = (start: any, end: any, task: any) => task.text

      gantt.templates.task_row_class = (start: any, end: any, task: any) => {
        if (!task.start_date || !task.end_date) return 'gantt_invalid_task'
        return ''
      }

      gantt.templates.date_grid = (date: any) => gantt.date.date_to_str('%Y-%m-%d')(date)

      // 确保容器存在
      if (ganttContainer.value) {
        console.log('初始化甘特图到容器:', ganttContainer.value)

        gantt.init(ganttContainer.value)
        bindEvents()
        isInitialized = true
        console.log('甘特图初始化完成')
      } else {
        console.error('甘特图容器不存在')
      }
    } catch (error) {
      console.error('初始化甘特图失败:', error)
    }
  }

  // 绑定事件
  const bindEvents = () => {
    // 任务拖拽结束 - 实时更新
    // gantt.attachEvent('onAfterTaskDrag', (id: any) => {
    //   if (isUpdating) return
    //   isUpdating = true

    //   const task = gantt.getTask(id)
    //   quickUpdateTask(task).finally(() => {
    //     setTimeout(() => {
    //       isUpdating = false
    //     }, 300)
    //   })
    // })

    // 任务双击事件
    gantt.attachEvent('onTaskDblClick', (id: any) => {
      emit('taskClick', id)
      return true
    })

    // 依赖创建
    gantt.attachEvent('onAfterLinkAdd', (id: any, link: any) => {
      saveDependency(link)
    })

    // 依赖删除
    gantt.attachEvent('onAfterLinkDelete', (id: any, link: any) => {
      deleteDependency(link.id)
    })

    // 数据加载完成
    gantt.attachEvent('onParse', () => {
      setTimeout(() => gantt.render(), 100)
    })
  }

  // 快速更新任务
  const quickUpdateTask = async (task: any) => {
    try {
      const startDate = new Date(task.start_date).toISOString()
      const endDate = new Date(task.end_date).toISOString()

      await taskApi.quickUpdateTask({
        id: task.id,
        startDate,
        endDate,
        progress: Math.round(task.progress * 100)
      })
      console.log('任务已实时更新')
    } catch (error) {
      console.error('实时更新任务失败:', error)
      ElMessage.error('更新任务失败')
      await loadGanttData()
    }
  }

  // 保存依赖关系
  const saveDependency = async (link: any) => {
    try {
      const typeMap: Record<number, string> = { 0: 'FS', 1: 'SS', 2: 'FF', 3: 'SF' }
      await taskApi.createDependency({
        projectId: props.projectId,
        predecessorId: link.source,
        successorId: link.target,
        dependencyType: typeMap[link.type] || 'FS',
        lagDays: 0
      })
      ElMessage.success('依赖关系已保存')
    } catch {
      ElMessage.error('保存依赖关系失败')
    }
  }

  // 删除依赖关系
  const deleteDependency = async (id: number) => {
    try {
      await taskApi.deleteDependency(id)
      ElMessage.success('依赖关系已删除')
    } catch {
      ElMessage.error('删除依赖关系失败')
    }
  }

  // 加载甘特图数据
  const loadGanttData = async () => {
    try {
      if (!isInitialized) {
        initGantt()
      }

      const res = (await taskApi.getGanttData(props.projectId)) as GanttResponse
      console.log('获取到数据:', res)

      // 处理空数据情况
      if (!res || !res.tasks || res.tasks.length === 0) {
        console.log('没有任务数据，显示空甘特图')
        const emptyData = {
          data: [],
          links: []
        }
        gantt.clearAll()
        gantt.parse(emptyData)
        emit('dataLoaded')
        return
      }

      const ganttData = {
        data: res.tasks
          .map((task: GanttTask) => {
            if (!task) return null

            const startDate = formatDateForGantt(task.start_date)
            const endDate = task.end_date ? formatDateForGantt(task.end_date) : null

            let taskType = gantt.config.types.task
            const duration = Number(task.duration)
            const durationStr = String(task.duration)
            if (task.type === 'milestone' && (duration === 0 || durationStr === '0')) {
              taskType = gantt.config.types.milestone
            }

            return {
              id: task.id,
              text: task.text || '未命名任务',
              start_date: startDate,
              end_date: endDate,
              duration: Number(task.duration) || 1,
              progress: Number(task.progress) || 0,
              parent: task.parent || 0,
              type: taskType,
              assigneeName: task.assigneeName,
              isOverdue: task.isOverdue
            }
          })
          .filter((task) => task !== null),
        links: (res.dependencies || []).map((dep: GanttDependency) => ({
          id: dep.id,
          source: dep.predecessorId,
          target: dep.successorId,
          type: getDependencyType(dep.dependencyType).toString()
        }))
      }

      console.log('处理后的甘特图数据:', ganttData)

      gantt.clearAll()
      gantt.parse(ganttData)

      if (ganttData.data.length > 0) {
        gantt.showTask(ganttData.data[0].id)
      }

      emit('dataLoaded')
    } catch (error) {
      console.error('加载甘特图数据失败:', error)
      ElMessage.error('加载甘特图数据失败')
    }
  }

  // 格式化日期
  const formatDateForGantt = (dateStr: string) => {
    if (!dateStr) return new Date().toISOString().split('T')[0] + ' 00:00:00'

    try {
      let date: Date
      if (dateStr.includes('T')) {
        date = new Date(dateStr)
      } else if (dateStr.includes('-')) {
        date = new Date(dateStr + 'T00:00:00')
      } else {
        date = new Date(dateStr)
      }

      if (isNaN(date.getTime())) {
        console.warn('无效日期格式:', dateStr)
        return new Date().toISOString().split('T')[0] + ' 00:00:00'
      }

      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')

      return `${year}-${month}-${day} 00:00:00`
    } catch (error) {
      console.error('日期格式化错误:', error, dateStr)
      return new Date().toISOString().split('T')[0] + ' 00:00:00'
    }
  }

  // 获取依赖类型
  const getDependencyType = (type: string) => {
    const types: Record<string, number> = { FS: 0, SS: 1, FF: 2, SF: 3 }
    return types[type] || 0
  }

  // 刷新数据
  const refresh = () => {
    loadGanttData()
  }

  // 销毁甘特图组件
  const destroy = () => {
    try {
      if (gantt && gantt.clearAll) {
        gantt.clearAll()
      }
      isInitialized = false
    } catch (error) {
      console.warn('销毁甘特图组件时出错:', error)
    }
  }

  onMounted(() => {
    if (ganttContainer.value && props.projectId) {
      initGantt()
      setTimeout(() => loadGanttData(), 300)
    }
  })

  onBeforeUnmount(() => {
    try {
      if (gantt && gantt.clearAll) {
        gantt.clearAll()
      }
      isInitialized = false
    } catch (error) {
      console.warn('甘特图清理时出错:', error)
    }
  })

  watch(
    () => props.projectId,
    (newId, oldId) => {
      if (newId !== oldId && newId && isInitialized) {
        loadGanttData()
      }
    },
    { immediate: false }
  )

  defineExpose({ refresh, loadGanttData, destroy })
</script>

<style scoped lang="scss">
  .gantt-chart-container {
    width: 100% !important;
    height: calc(100vh - 200px) !important;
    min-height: 600px !important;
    max-height: calc(100vh - 200px) !important;
    overflow: visible;
    position: relative;
    background: #fff;
    border: 1px solid #e4e7ed;
    display: block !important;
  }

  :deep(.gantt_container) {
    width: 100% !important;
    height: 100% !important;
  }

  :deep(.gantt_layout) {
    overflow: visible !important;
  }

  :deep(.gantt_layout_cell) {
    overflow: visible !important;
  }

  // 左侧任务列表区域
  :deep(.gantt_grid) {
    overflow-y: auto !important;
    overflow-x: hidden !important;
  }

  // 右侧时间轴区域 - 关键：启用水平和垂直滚动
  :deep(.gantt_task) {
    overflow: auto !important;
  }

  :deep(.gantt_task_bg) {
    overflow: visible !important;
  }

  :deep(.gantt_data_area) {
    overflow: visible !important;
  }

  :deep(.gantt_milestone) {
    background-color: #f56c6c;
    border-color: #f56c6c;

    .gantt_task_content {
      width: 16px !important;
      height: 16px !important;
      margin: 2px auto !important;
      background-color: #f56c6c !important;
      border: 2px solid #d32f2f !important;
      transform: rotate(45deg) !important;
    }
  }

  :deep(.gantt_overdue) {
    background-color: #e6a23c;
    border-color: #e6a23c;
  }

  :deep(.gantt_invalid_task) {
    .gantt_task_content {
      background-color: #ffcccc !important;
      border: 1px solid #ff0000 !important;
    }
  }

  :deep(.gantt_task_line) {
    min-height: 20px;
    display: flex;
    align-items: center;
  }

  :deep(.gantt_task_content) {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    color: white;
    text-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
  }

  :deep(.gantt_row) {
    min-height: 30px;
  }

  :deep(.gantt_scale_line) {
    border-bottom: 1px solid #e0e0e0;
  }
</style>
