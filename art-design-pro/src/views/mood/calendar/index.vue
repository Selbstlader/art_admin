<template>
  <div class="mood-calendar-page">
    <ElRow :gutter="16">
      <!-- 日历视图 -->
      <ElCol :xs="24" :lg="16" class="mb-4">
        <ElCard shadow="never">
          <template #header>
            <div class="calendar-header">
              <ElButton :icon="ArrowLeft" circle @click="prevMonth" />
              <span class="calendar-title">{{ currentYear }}年{{ currentMonth }}月</span>
              <ElButton :icon="ArrowRight" circle @click="nextMonth" />
            </div>
          </template>

          <div class="calendar-grid" v-loading="loading">
            <!-- 星期标题 -->
            <div class="weekday-header">
              <div v-for="day in weekDays" :key="day" class="weekday-cell">{{ day }}</div>
            </div>

            <!-- 日期格子 -->
            <div class="days-grid">
              <div
                v-for="(day, index) in calendarDays"
                :key="index"
                class="day-cell"
                :class="{
                  'other-month': !day.isCurrentMonth,
                  today: day.isToday,
                  selected: selectedDate === day.date,
                  'has-record': day.hasRecord
                }"
                @click="selectDate(day)"
              >
                <div class="day-number">{{ day.day }}</div>
                <div v-if="day.hasRecord" class="day-mood">
                  <span class="mood-emoji">{{ day.moodIcon }}</span>
                  <span class="mood-intensity">{{ day.avgIntensity?.toFixed(1) }}</span>
                </div>
              </div>
            </div>
          </div>
        </ElCard>
      </ElCol>

      <!-- 选中日期详情 -->
      <ElCol :xs="24" :lg="8" class="mb-4">
        <ElCard shadow="never">
          <template #header>
            <span class="card-title">{{ selectedDate || '请选择日期' }} 情绪记录</span>
          </template>

          <div v-loading="detailLoading">
            <div v-if="selectedRecords.length > 0" class="record-list">
              <div v-for="record in selectedRecords" :key="record.id" class="record-item">
                <div class="record-header">
                  <span class="mood-emoji">{{ MoodTypeConfig[record.mood_type]?.icon }}</span>
                  <span class="mood-label">{{ MoodTypeConfig[record.mood_type]?.label }}</span>
                  <span class="record-time">{{ formatTime(record.created_at) }}</span>
                </div>
                <div class="record-intensity">
                  <span>强度:</span>
                  <ElProgress
                    :percentage="record.intensity * 10"
                    :stroke-width="6"
                    :show-text="false"
                    :color="getIntensityColor(record.intensity)"
                    style="flex: 1; margin: 0 8px"
                  />
                  <span>{{ record.intensity }}</span>
                </div>
                <div v-if="record.note" class="record-note">{{ record.note }}</div>
                <div v-if="record.triggers?.length" class="record-tags">
                  <ElTag v-for="tag in record.triggers" :key="tag" size="small" class="mr-1">
                    {{ tag }}
                  </ElTag>
                </div>
              </div>
            </div>
            <ElEmpty v-else description="该日期暂无记录" />
          </div>
        </ElCard>

        <!-- 月度统计 -->
        <ElCard shadow="never" class="mt-4">
          <template #header>
            <span class="card-title">本月统计</span>
          </template>

          <div v-loading="statsLoading" class="stats-content">
            <div class="stat-item">
              <span class="stat-label">记录天数</span>
              <span class="stat-value">{{ monthStats.recordDays }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总记录数</span>
              <span class="stat-value">{{ monthStats.totalRecords }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">平均强度</span>
              <span class="stat-value">{{ monthStats.avgIntensity?.toFixed(1) || '-' }}</span>
            </div>
            <div class="mood-distribution">
              <div class="dist-title">情绪分布</div>
              <div class="dist-list">
                <div
                  v-for="(count, type) in monthStats.moodDistribution"
                  :key="type"
                  class="dist-item"
                >
                  <span class="dist-icon">{{ MoodTypeConfig[type as MoodType]?.icon }}</span>
                  <span class="dist-label">{{ MoodTypeConfig[type as MoodType]?.label }}</span>
                  <span class="dist-count">{{ count }}</span>
                </div>
              </div>
            </div>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, watch } from 'vue'
  import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
  import { moodRecordApi, MoodTypeConfig, type MoodType, type MoodRecord } from '@/api/mood'
  import dayjs from 'dayjs'

  defineOptions({ name: 'MoodCalendar' })

  const weekDays = ['日', '一', '二', '三', '四', '五', '六']

  // 当前年月
  const currentYear = ref(dayjs().year())
  const currentMonth = ref(dayjs().month() + 1)

  // 选中日期
  const selectedDate = ref(dayjs().format('YYYY-MM-DD'))

  // 加载状态
  const loading = ref(false)
  const detailLoading = ref(false)
  const statsLoading = ref(false)

  // 日历数据
  const calendarData = ref<Record<string, any>>({})
  const selectedRecords = ref<MoodRecord[]>([])

  // 月度统计
  const monthStats = ref({
    recordDays: 0,
    totalRecords: 0,
    avgIntensity: 0,
    moodDistribution: {} as Record<string, number>
  })

  // 计算日历格子
  const calendarDays = computed(() => {
    const year = currentYear.value
    const month = currentMonth.value
    const firstDay = dayjs(`${year}-${month}-01`)
    const lastDay = firstDay.endOf('month')
    const startWeekDay = firstDay.day()
    const daysInMonth = lastDay.date()

    const days: any[] = []

    // 上月填充
    const prevMonth = firstDay.subtract(1, 'month')
    const prevMonthDays = prevMonth.daysInMonth()
    for (let i = startWeekDay - 1; i >= 0; i--) {
      const date = prevMonth.date(prevMonthDays - i).format('YYYY-MM-DD')
      days.push({
        day: prevMonthDays - i,
        date,
        isCurrentMonth: false,
        isToday: false,
        hasRecord: !!calendarData.value[date],
        moodIcon: calendarData.value[date]?.icon,
        avgIntensity: calendarData.value[date]?.avgIntensity
      })
    }

    // 当月
    const today = dayjs().format('YYYY-MM-DD')
    for (let i = 1; i <= daysInMonth; i++) {
      const date = firstDay.date(i).format('YYYY-MM-DD')
      days.push({
        day: i,
        date,
        isCurrentMonth: true,
        isToday: date === today,
        hasRecord: !!calendarData.value[date],
        moodIcon: calendarData.value[date]?.icon,
        avgIntensity: calendarData.value[date]?.avgIntensity
      })
    }

    // 下月填充
    const remaining = 42 - days.length
    const nextMonth = firstDay.add(1, 'month')
    for (let i = 1; i <= remaining; i++) {
      const date = nextMonth.date(i).format('YYYY-MM-DD')
      days.push({
        day: i,
        date,
        isCurrentMonth: false,
        isToday: false,
        hasRecord: !!calendarData.value[date],
        moodIcon: calendarData.value[date]?.icon,
        avgIntensity: calendarData.value[date]?.avgIntensity
      })
    }

    return days
  })

  // 方法
  const prevMonth = () => {
    if (currentMonth.value === 1) {
      currentMonth.value = 12
      currentYear.value--
    } else {
      currentMonth.value--
    }
  }

  const nextMonth = () => {
    if (currentMonth.value === 12) {
      currentMonth.value = 1
      currentYear.value++
    } else {
      currentMonth.value++
    }
  }

  const selectDate = (day: any) => {
    selectedDate.value = day.date
    loadDateRecords(day.date)
  }

  const formatTime = (time: string) => {
    return dayjs(time).format('HH:mm')
  }

  const getIntensityColor = (intensity: number) => {
    if (intensity <= 3) return '#67C23A'
    if (intensity <= 6) return '#E6A23C'
    return '#F56C6C'
  }

  // 加载日历数据
  const loadCalendarData = async () => {
    loading.value = true
    try {
      const res = await moodRecordApi.getCalendar({
        year: currentYear.value,
        month: currentMonth.value
      })
      const data = (res as any)?.data || res || []

      // 转换为日期映射
      const dataMap: Record<string, any> = {}
      data.forEach((item: any) => {
        dataMap[item.date] = {
          icon: MoodTypeConfig[item.mood_type as MoodType]?.icon || '😐',
          avgIntensity: item.avg_intensity
        }
      })
      calendarData.value = dataMap

      // 计算月度统计
      calculateMonthStats(data)
    } catch (error) {
      console.error('加载日历数据失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 加载日期记录
  const loadDateRecords = async (date: string) => {
    detailLoading.value = true
    try {
      const res = await moodRecordApi.getByDate(date)
      selectedRecords.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载日期记录失败:', error)
    } finally {
      detailLoading.value = false
    }
  }

  // 计算月度统计
  const calculateMonthStats = (data: any[]) => {
    const recordDays = new Set(data.map((d: any) => d.date)).size
    const totalRecords = data.reduce((sum: number, d: any) => sum + (d.count || 1), 0)
    const avgIntensity =
      data.length > 0
        ? data.reduce((sum: number, d: any) => sum + d.avg_intensity, 0) / data.length
        : 0

    const moodDistribution: Record<string, number> = {}
    data.forEach((d: any) => {
      moodDistribution[d.mood_type] = (moodDistribution[d.mood_type] || 0) + (d.count || 1)
    })

    monthStats.value = {
      recordDays,
      totalRecords,
      avgIntensity,
      moodDistribution
    }
  }

  // 监听年月变化
  watch([currentYear, currentMonth], () => {
    loadCalendarData()
  })

  onMounted(() => {
    loadCalendarData()
    loadDateRecords(selectedDate.value)
  })
</script>

<style scoped lang="scss">
  .mood-calendar-page {
    padding: 16px;
  }

  .mb-4 {
    margin-bottom: 16px;
  }

  .mt-4 {
    margin-top: 16px;
  }

  .mr-1 {
    margin-right: 4px;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
  }

  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;

    .calendar-title {
      font-size: 18px;
      font-weight: 600;
      min-width: 120px;
      text-align: center;
    }
  }

  .calendar-grid {
    .weekday-header {
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      margin-bottom: 8px;

      .weekday-cell {
        text-align: center;
        padding: 8px;
        font-weight: 500;
        color: #909399;
      }
    }

    .days-grid {
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      gap: 4px;

      .day-cell {
        aspect-ratio: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4px;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          background: #f5f7fa;
        }

        &.other-month {
          opacity: 0.4;
        }

        &.today {
          background: #ecf5ff;
          border: 1px solid #409eff;
        }

        &.selected {
          background: #409eff;
          color: #fff;

          .mood-intensity {
            color: #fff;
          }
        }

        &.has-record {
          background: #f0f9eb;

          &.selected {
            background: #409eff;
          }
        }

        .day-number {
          font-size: 14px;
          font-weight: 500;
        }

        .day-mood {
          display: flex;
          align-items: center;
          gap: 2px;
          margin-top: 2px;

          .mood-emoji {
            font-size: 14px;
          }

          .mood-intensity {
            font-size: 10px;
            color: #909399;
          }
        }
      }
    }
  }

  .record-list {
    .record-item {
      padding: 12px;
      background: #fafafa;
      border-radius: 8px;
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }

      .record-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;

        .mood-emoji {
          font-size: 24px;
        }

        .mood-label {
          font-weight: 500;
          flex: 1;
        }

        .record-time {
          font-size: 12px;
          color: #909399;
        }
      }

      .record-intensity {
        display: flex;
        align-items: center;
        font-size: 12px;
        color: #606266;
        margin-bottom: 8px;
      }

      .record-note {
        font-size: 13px;
        color: #606266;
        padding: 8px;
        background: #fff;
        border-radius: 4px;
        margin-bottom: 8px;
      }

      .record-tags {
        display: flex;
        flex-wrap: wrap;
      }
    }
  }

  .stats-content {
    .stat-item {
      display: flex;
      justify-content: space-between;
      padding: 12px 0;
      border-bottom: 1px solid #ebeef5;

      .stat-label {
        color: #909399;
      }

      .stat-value {
        font-weight: 600;
        color: #303133;
      }
    }

    .mood-distribution {
      margin-top: 16px;

      .dist-title {
        font-weight: 500;
        margin-bottom: 12px;
      }

      .dist-list {
        .dist-item {
          display: flex;
          align-items: center;
          padding: 8px 0;

          .dist-icon {
            font-size: 20px;
            margin-right: 8px;
          }

          .dist-label {
            flex: 1;
            color: #606266;
          }

          .dist-count {
            font-weight: 500;
            color: #303133;
          }
        }
      }
    }
  }
</style>
