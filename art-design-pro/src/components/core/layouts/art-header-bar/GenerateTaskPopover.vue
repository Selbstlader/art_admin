<template>
  <ElPopover
    placement="bottom-end"
    :width="400"
    trigger="click"
    popper-class="generate-task-popover"
    @show="loadTasks"
  >
    <template #reference>
      <div class="btn-box task-btn">
        <div class="btn">
          <ElIcon><Document /></ElIcon>
          <ElBadge v-if="processingCount > 0" :value="processingCount" class="task-badge" />
        </div>
      </div>
    </template>
    <template #default>
      <div class="task-list-container">
        <div class="task-header">
          <h3>生成任务</h3>
          <ElButton text size="small" @click="loadTasks">
            <ElIcon><Refresh /></ElIcon>
          </ElButton>
        </div>

        <ElDivider style="margin: 12px 0" />

        <div v-if="loading" class="loading-wrapper">
          <ElSkeleton :rows="3" animated />
        </div>

        <ElEmpty v-else-if="tasks.length === 0" description="暂无生成任务" :image-size="80" />

        <div v-else class="task-list">
          <div v-for="task in tasks" :key="task.id" class="task-item">
            <div class="task-info">
              <div class="task-title">{{ task.topic }}</div>
              <div class="task-meta">
                <ElTag size="small">{{ task.grade }}</ElTag>
                <ElTag size="small" :type="getDifficultyType(task.difficulty)">
                  {{ getDifficultyText(task.difficulty) }}
                </ElTag>
              </div>
            </div>
            <div class="task-status">
              <ElTag v-if="task.status === 0" type="info" size="small">待处理</ElTag>
              <ElTag v-else-if="task.status === 1" type="warning" size="small">
                <ElIcon class="is-loading"><Loading /></ElIcon>
                生成中
              </ElTag>
            </div>
          </div>
        </div>

        <div v-if="tasks.length > 0" class="task-footer">
          <ElButton text size="small" @click="$router.push('/learning/learning/my')">
            查看全部
          </ElButton>
        </div>
      </div>
    </template>
  </ElPopover>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { Document, Refresh, Loading } from '@element-plus/icons-vue'
  import { learningMaterialApi, type GenerateTask } from '@/api/learning'

  const loading = ref(false)
  const tasks = ref<GenerateTask[]>([])
  const processingCount = ref(0)

  const loadTasks = async () => {
    try {
      loading.value = true
      const res: any = await learningMaterialApi.getTaskList()
      if (res) {
        tasks.value = res.slice(0, 10) // 只显示最近10个
        // 后端已默认只返回进行中的任务,直接使用数量
        processingCount.value = res.length
      }
    } catch (error) {
      console.error('获取任务列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const getDifficultyType = (difficulty: number) => {
    const types: Record<number, any> = { 1: 'success', 2: 'warning', 3: 'danger' }
    return types[difficulty] || ''
  }

  const getDifficultyText = (difficulty: number) => {
    const texts: Record<number, string> = { 1: '⭐ 基础', 2: '⭐⭐ 进阶', 3: '⭐⭐⭐ 高级' }
    return texts[difficulty] || '未知'
  }
</script>

<style scoped lang="scss">
  .task-btn {
    position: relative;
    cursor: pointer;

    .btn {
      position: relative;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: 6px;
      transition: all 0.3s;
      font-size: 18px;
      &:hover {
        background: var(--el-fill-color-light);
      }
    }

    :deep(.task-badge) {
      position: absolute;
      top: -4px;
      right: -4px;
    }
  }

  .task-list-container {
    max-height: 500px;
    overflow-y: auto;

    .task-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 600;
      }
    }

    .loading-wrapper {
      padding: 20px 0;
    }

    .task-list {
      .task-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 0;
        border-bottom: 1px solid var(--el-border-color-lighter);

        &:last-child {
          border-bottom: none;
        }

        .task-info {
          flex: 1;
          min-width: 0;

          .task-title {
            font-size: 14px;
            font-weight: 500;
            margin-bottom: 6px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .task-meta {
            display: flex;
            gap: 6px;
          }
        }

        .task-status {
          margin-left: 12px;
        }
      }
    }

    .task-footer {
      padding-top: 12px;
      text-align: center;
      border-top: 1px solid var(--el-border-color-lighter);
    }
  }
</style>
