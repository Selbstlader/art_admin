<template>
  <div class="waypoint-list">
    <!-- 空状态 -->
    <div v-if="modelValue.length === 0" class="empty-state">
      <ElEmpty description="暂无途经点，点击上方按钮添加" :image-size="80" />
    </div>

    <!-- 途经点列表 -->
    <div v-else class="waypoint-items" ref="listRef">
      <TransitionGroup name="list">
        <div
          v-for="(waypoint, index) in modelValue"
          :key="waypoint.id || `temp-${index}`"
          class="waypoint-item"
          :class="{ 
            'is-start': index === 0,
            'is-end': index === modelValue.length - 1,
            'is-dragging': dragIndex === index
          }"
          draggable="true"
          @dragstart="handleDragStart($event, index)"
          @dragover="handleDragOver($event, index)"
          @dragend="handleDragEnd"
        >
          <!-- 序号和连接线 -->
          <div class="waypoint-index">
            <div class="index-circle" :class="getIndexClass(index)">
              {{ index + 1 }}
            </div>
            <div v-if="index < modelValue.length - 1" class="connect-line"></div>
          </div>

          <!-- 途经点信息 -->
          <div class="waypoint-content">
            <div class="waypoint-header">
              <div class="waypoint-name">
                <span class="name">{{ waypoint.name || '未命名地点' }}</span>
                <ElTag v-if="days > 1" size="small" type="info">
                  第{{ waypoint.dayIndex || 1 }}天
                </ElTag>
              </div>
              <div class="waypoint-actions">
                <ElButton
                  type="primary"
                  link
                  :icon="Edit"
                  size="small"
                  @click="handleEdit(waypoint, index)"
                >
                  编辑
                </ElButton>
                <ElPopconfirm
                  title="确定删除该途经点吗？"
                  confirm-button-text="确定"
                  cancel-button-text="取消"
                  @confirm="handleDelete(index)"
                >
                  <template #reference>
                    <ElButton
                      type="danger"
                      link
                      :icon="Delete"
                      size="small"
                    >
                      删除
                    </ElButton>
                  </template>
                </ElPopconfirm>
              </div>
            </div>

            <div class="waypoint-details">
              <div v-if="waypoint.address" class="detail-item">
                <ElIcon><Location /></ElIcon>
                <span>{{ waypoint.address }}</span>
              </div>
              <div v-if="waypoint.stayDuration" class="detail-item">
                <ElIcon><Clock /></ElIcon>
                <span>停留 {{ formatDuration(waypoint.stayDuration) }}</span>
              </div>
              <div v-if="waypoint.budget && waypoint.budget > 0" class="detail-item">
                <ElIcon><Money /></ElIcon>
                <span>预算 ¥{{ waypoint.budget }}</span>
              </div>
            </div>

            <div v-if="waypoint.notes" class="waypoint-notes">
              <span class="notes-label">备注：</span>
              <span class="notes-content">{{ waypoint.notes }}</span>
            </div>
          </div>

          <!-- 拖拽手柄 -->
          <div class="drag-handle">
            <ElIcon><Rank /></ElIcon>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Edit, Delete, Location, Clock, Money, Rank } from '@element-plus/icons-vue'
import type { Waypoint } from '../types'

defineOptions({
  name: 'WaypointList'
})

// Props
interface Props {
  modelValue: Partial<Waypoint>[]
  days?: number
}

const props = withDefaults(defineProps<Props>(), {
  days: 1
})

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: Partial<Waypoint>[]): void
  (e: 'edit', waypoint: Partial<Waypoint>, index: number): void
  (e: 'delete', index: number): void
  (e: 'reorder', waypoints: Partial<Waypoint>[]): void
}>()

// 拖拽相关
const listRef = ref<HTMLElement | null>(null)
const dragIndex = ref<number | null>(null)
const dropIndex = ref<number | null>(null)

// 获取序号样式类
const getIndexClass = (index: number) => {
  if (index === 0) return 'start'
  if (index === props.modelValue.length - 1) return 'end'
  return 'middle'
}

// 格式化停留时间
const formatDuration = (minutes: number): string => {
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
}

// 编辑途经点
const handleEdit = (waypoint: Partial<Waypoint>, index: number) => {
  emit('edit', waypoint, index)
}

// 删除途经点
const handleDelete = (index: number) => {
  emit('delete', index)
}

// 拖拽开始
const handleDragStart = (event: DragEvent, index: number) => {
  dragIndex.value = index
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', String(index))
  }
}

// 拖拽经过
const handleDragOver = (event: DragEvent, index: number) => {
  event.preventDefault()
  if (dragIndex.value === null || dragIndex.value === index) return
  
  dropIndex.value = index
  
  // 重新排序
  const newList = [...props.modelValue]
  const dragItem = newList[dragIndex.value]
  newList.splice(dragIndex.value, 1)
  newList.splice(index, 0, dragItem)
  
  dragIndex.value = index
  emit('update:modelValue', newList)
}

// 拖拽结束
const handleDragEnd = () => {
  if (dragIndex.value !== null) {
    emit('reorder', props.modelValue)
  }
  dragIndex.value = null
  dropIndex.value = null
}
</script>

<style scoped lang="scss">
.waypoint-list {
  .empty-state {
    padding: 40px 0;
  }

  .waypoint-items {
    display: flex;
    flex-direction: column;
  }

  .waypoint-item {
    display: flex;
    gap: 12px;
    padding: 12px;
    margin-bottom: 8px;
    background: var(--art-main-bg-color);
    border: 1px solid var(--art-border-color);
    border-radius: 8px;
    transition: all 0.3s;
    cursor: grab;

    &:hover {
      border-color: var(--el-color-primary-light-5);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    }

    &.is-dragging {
      opacity: 0.5;
      border-color: var(--el-color-primary);
    }

    &:active {
      cursor: grabbing;
    }

    // 序号和连接线
    .waypoint-index {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 32px;
      flex-shrink: 0;

      .index-circle {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        border-radius: 50%;
        font-size: 14px;
        font-weight: 600;
        color: #fff;

        &.start {
          background: var(--el-color-success);
        }

        &.end {
          background: var(--el-color-danger);
        }

        &.middle {
          background: var(--el-color-primary);
        }
      }

      .connect-line {
        flex: 1;
        width: 2px;
        min-height: 20px;
        margin-top: 4px;
        background: var(--art-border-color);
      }
    }

    // 途经点内容
    .waypoint-content {
      flex: 1;
      min-width: 0;

      .waypoint-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 8px;

        .waypoint-name {
          display: flex;
          align-items: center;
          gap: 8px;

          .name {
            font-size: 15px;
            font-weight: 600;
            color: var(--art-text-gray-800);
          }
        }

        .waypoint-actions {
          display: flex;
          gap: 4px;
        }
      }

      .waypoint-details {
        display: flex;
        flex-wrap: wrap;
        gap: 16px;
        margin-bottom: 8px;

        .detail-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          color: var(--art-text-gray-600);

          .el-icon {
            font-size: 14px;
            color: var(--art-text-gray-400);
          }
        }
      }

      .waypoint-notes {
        font-size: 13px;
        color: var(--art-text-gray-500);
        padding: 8px;
        background: var(--art-gray-100);
        border-radius: 4px;

        .notes-label {
          color: var(--art-text-gray-600);
        }
      }
    }

    // 拖拽手柄
    .drag-handle {
      display: flex;
      align-items: center;
      padding: 0 4px;
      color: var(--art-text-gray-400);
      cursor: grab;
      opacity: 0;
      transition: opacity 0.2s;

      &:active {
        cursor: grabbing;
      }
    }

    &:hover {
      .drag-handle {
        opacity: 1;
      }
    }
  }
}

// 列表动画
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

.list-move {
  transition: transform 0.3s ease;
}
</style>
