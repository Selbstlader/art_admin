<template>
  <div class="schedule-timeline">
    <!-- 空状态 -->
    <div v-if="waypoints.length === 0" class="empty-state">
      <ElEmpty description="暂无行程安排" :image-size="60" />
    </div>

    <!-- 日程时间线 -->
    <div v-else class="timeline-container">
      <!-- 日期标签页 -->
      <ElTabs v-model="activeDay" type="card" class="day-tabs">
        <ElTabPane
          v-for="day in days"
          :key="day"
          :label="getDayLabel(day)"
          :name="String(day)"
        >
          <div class="day-content">
            <!-- 日期信息 -->
            <div class="day-header">
              <div class="day-date">
                <span class="date-text">{{ getDateText(day) }}</span>
                <span class="weekday">{{ getWeekday(day) }}</span>
              </div>
              <div class="day-stats">
                <span class="stat-item">
                  <ElIcon><Location /></ElIcon>
                  {{ getDayWaypoints(day).length }} 个地点
                </span>
                <span class="stat-item">
                  <ElIcon><Clock /></ElIcon>
                  {{ getDayDuration(day) }}
                </span>
                <span v-if="getDayBudget(day) > 0" class="stat-item">
                  <ElIcon><Money /></ElIcon>
                  ¥{{ getDayBudget(day) }}
                </span>
              </div>
            </div>

            <!-- 时间线 -->
            <ElTimeline class="waypoint-timeline">
              <ElTimelineItem
                v-for="(waypoint, index) in getDayWaypoints(day)"
                :key="waypoint.id || `wp-${index}`"
                :type="getTimelineType(index, getDayWaypoints(day).length)"
                :hollow="false"
                :size="index === 0 || index === getDayWaypoints(day).length - 1 ? 'large' : 'normal'"
              >
                <template #dot>
                  <div class="timeline-dot" :class="getDotClass(index, getDayWaypoints(day).length)">
                    {{ index + 1 }}
                  </div>
                </template>

                <div class="timeline-content">
                  <div class="waypoint-card">
                    <div class="card-header">
                      <span class="waypoint-name">{{ waypoint.name }}</span>
                      <ElTag v-if="waypoint.waypointType === 1" size="small" type="success">起点</ElTag>
                      <ElTag v-else-if="waypoint.waypointType === 3" size="small" type="danger">终点</ElTag>
                    </div>
                    
                    <div v-if="waypoint.address" class="card-address">
                      <ElIcon><Location /></ElIcon>
                      <span>{{ waypoint.address }}</span>
                    </div>

                    <div class="card-info">
                      <span v-if="waypoint.stayDuration" class="info-item">
                        <ElIcon><Clock /></ElIcon>
                        停留 {{ formatDuration(waypoint.stayDuration) }}
                      </span>
                      <span v-if="waypoint.budget && waypoint.budget > 0" class="info-item">
                        <ElIcon><Money /></ElIcon>
                        ¥{{ waypoint.budget }}
                      </span>
                    </div>

                    <div v-if="waypoint.notes" class="card-notes">
                      {{ waypoint.notes }}
                    </div>
                  </div>
                </div>
              </ElTimelineItem>
            </ElTimeline>

            <!-- 当天无行程 -->
            <div v-if="getDayWaypoints(day).length === 0" class="no-waypoints">
              <ElIcon><Calendar /></ElIcon>
              <span>当天暂无行程安排</span>
            </div>
          </div>
        </ElTabPane>
      </ElTabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Location, Clock, Money, Calendar } from '@element-plus/icons-vue'
import { useDateFormat } from '@vueuse/core'
import type { Waypoint } from '../types'
import { WaypointType } from '../types'

defineOptions({
  name: 'ScheduleTimeline'
})

// Props
interface Props {
  waypoints: Partial<Waypoint>[]
  days: number
  startDate?: string
}

const props = withDefaults(defineProps<Props>(), {
  days: 1,
  startDate: ''
})

// 当前激活的天数
const activeDay = ref('1')

// 获取某天的途经点
const getDayWaypoints = (day: number): Partial<Waypoint>[] => {
  return props.waypoints
    .filter(wp => (wp.dayIndex || 1) === day)
    .sort((a, b) => (a.sortOrder || 0) - (b.sortOrder || 0))
}

// 获取日期标签
const getDayLabel = (day: number): string => {
  return `第${day}天`
}

// 获取日期文本
const getDateText = (day: number): string => {
  if (!props.startDate) return ''
  const date = new Date(props.startDate)
  date.setDate(date.getDate() + day - 1)
  return useDateFormat(date, 'MM月DD日').value
}

// 获取星期几
const getWeekday = (day: number): string => {
  if (!props.startDate) return ''
  const date = new Date(props.startDate)
  date.setDate(date.getDate() + day - 1)
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekdays[date.getDay()]
}

// 获取某天的总时长
const getDayDuration = (day: number): string => {
  const waypoints = getDayWaypoints(day)
  const totalMinutes = waypoints.reduce((sum, wp) => sum + (wp.stayDuration || 0), 0)
  return formatDuration(totalMinutes)
}

// 获取某天的总预算
const getDayBudget = (day: number): number => {
  const waypoints = getDayWaypoints(day)
  return waypoints.reduce((sum, wp) => sum + (wp.budget || 0), 0)
}

// 格式化时长
const formatDuration = (minutes: number): string => {
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
}

// 获取时间线类型
const getTimelineType = (index: number, total: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (index === 0) return 'success'
  if (index === total - 1) return 'danger'
  return 'primary'
}

// 获取圆点样式类
const getDotClass = (index: number, total: number): string => {
  if (index === 0) return 'start'
  if (index === total - 1) return 'end'
  return 'middle'
}
</script>

<style scoped lang="scss">
.schedule-timeline {
  .empty-state {
    padding: 40px 0;
  }

  .timeline-container {
    .day-tabs {
      :deep(.el-tabs__header) {
        margin-bottom: 16px;
      }

      :deep(.el-tabs__item) {
        padding: 0 20px;
      }
    }

    .day-content {
      .day-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        margin-bottom: 16px;
        background: var(--art-gray-100);
        border-radius: 8px;

        .day-date {
          display: flex;
          align-items: baseline;
          gap: 8px;

          .date-text {
            font-size: 16px;
            font-weight: 600;
            color: var(--art-text-gray-800);
          }

          .weekday {
            font-size: 13px;
            color: var(--art-text-gray-500);
          }
        }

        .day-stats {
          display: flex;
          gap: 16px;

          .stat-item {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 13px;
            color: var(--art-text-gray-600);

            .el-icon {
              font-size: 14px;
            }
          }
        }
      }

      .waypoint-timeline {
        padding-left: 8px;

        :deep(.el-timeline-item__wrapper) {
          padding-left: 20px;
        }

        .timeline-dot {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 24px;
          height: 24px;
          border-radius: 50%;
          font-size: 12px;
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

        .timeline-content {
          .waypoint-card {
            padding: 12px 16px;
            background: var(--art-main-bg-color);
            border: 1px solid var(--art-border-color);
            border-radius: 8px;
            transition: all 0.3s;

            &:hover {
              border-color: var(--el-color-primary-light-5);
              box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
            }

            .card-header {
              display: flex;
              align-items: center;
              gap: 8px;
              margin-bottom: 8px;

              .waypoint-name {
                font-size: 15px;
                font-weight: 600;
                color: var(--art-text-gray-800);
              }
            }

            .card-address {
              display: flex;
              align-items: center;
              gap: 4px;
              margin-bottom: 8px;
              font-size: 13px;
              color: var(--art-text-gray-500);

              .el-icon {
                font-size: 14px;
              }
            }

            .card-info {
              display: flex;
              gap: 16px;
              margin-bottom: 8px;

              .info-item {
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

            .card-notes {
              padding: 8px;
              font-size: 13px;
              color: var(--art-text-gray-500);
              background: var(--art-gray-100);
              border-radius: 4px;
            }
          }
        }
      }

      .no-waypoints {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
        padding: 40px 0;
        color: var(--art-text-gray-400);

        .el-icon {
          font-size: 32px;
        }
      }
    }
  }
}

// 响应式
@media only screen and (max-width: 768px) {
  .schedule-timeline {
    .timeline-container {
      .day-content {
        .day-header {
          flex-direction: column;
          gap: 8px;
          align-items: flex-start;

          .day-stats {
            flex-wrap: wrap;
            gap: 12px;
          }
        }
      }
    }
  }
}
</style>
