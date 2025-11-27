<template>
  <div class="checkin-page">
    <!-- 今日打卡状态 -->
    <div class="card art-custom-card today-card">
      <div class="card-header">
        <div class="title">
          <h4 class="box-title">今日打卡</h4>
          <p class="subtitle">{{ todayDate }}</p>
        </div>
        <div class="progress-badge">
          <span class="progress-value">{{ todayStatus.total || 0 }}</span>
          <span class="progress-total">/4</span>
        </div>
      </div>
      <div class="checkin-grid">
        <div
          v-for="(config, type) in CheckinTypeConfig"
          :key="type"
          class="checkin-item"
          :class="{ completed: todayStatus[type] }"
          @click="handleCheckin(type)"
        >
          <div class="checkin-icon">{{ config.icon }}</div>
          <div class="checkin-label">{{ config.label }}</div>
          <div class="checkin-status">
            <ElIcon v-if="todayStatus[type]" color="var(--el-color-success)">
              <CircleCheckFilled />
            </ElIcon>
            <ElIcon v-else color="var(--art-gray-400)">
              <CircleClose />
            </ElIcon>
          </div>
        </div>
      </div>
    </div>

    <!-- 连续打卡统计 -->
    <ElRow :gutter="20" class="stat-card-list">
      <ElCol v-for="(stat, index) in streakStats" :key="index" :xs="12" :sm="6">
        <div class="card art-custom-card stat-card">
          <div class="stat-value" :class="stat.colorClass">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-icon-wrapper" :class="stat.bgClass">
            <span class="stat-emoji">{{ stat.icon }}</span>
          </div>
        </div>
      </ElCol>
    </ElRow>

    <!-- 打卡日历 -->
    <div class="card art-custom-card calendar-card">
      <div class="card-header">
        <div class="title">
          <h4 class="box-title">打卡日历</h4>
          <p class="subtitle">查看打卡记录</p>
        </div>
        <div class="month-nav">
          <ElButton :icon="ArrowLeft" circle size="small" @click="prevMonth" />
          <span class="month-text">{{ currentMonth.format('YYYY年MM月') }}</span>
          <ElButton :icon="ArrowRight" circle size="small" @click="nextMonth" />
        </div>
      </div>
      <div class="calendar-wrapper">
        <div class="calendar-header">
          <div v-for="day in weekDays" :key="day" class="calendar-day-name">{{ day }}</div>
        </div>
        <div class="calendar-body">
          <div
            v-for="(day, index) in calendarDays"
            :key="index"
            class="calendar-day"
            :class="{
              'other-month': !day.isCurrentMonth,
              today: day.isToday,
              'has-checkin': day.checkinCount > 0
            }"
          >
            <span class="day-number">{{ day.date }}</span>
            <div v-if="day.checkinCount > 0" class="checkin-dots">
              <span v-for="n in Math.min(day.checkinCount, 4)" :key="n" class="dot"></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 打卡历史 -->
    <div class="card art-custom-card history-card">
      <div class="card-header">
        <div class="title">
          <h4 class="box-title">打卡历史</h4>
          <p class="subtitle">近30天记录</p>
        </div>
      </div>
      <div v-loading="loading" class="history-content">
        <div v-if="historyList.length > 0" class="history-list">
          <div v-for="record in historyList" :key="record.id" class="history-item">
            <div class="history-icon">{{ CheckinTypeConfig[record.checkin_type]?.icon }}</div>
            <div class="history-info">
              <div class="history-type">{{ CheckinTypeConfig[record.checkin_type]?.label }}</div>
              <div v-if="record.note" class="history-note">{{ record.note }}</div>
            </div>
            <div class="history-time">{{ formatTime(record.created_at) }}</div>
          </div>
        </div>
        <ElEmpty v-else description="暂无打卡记录" :image-size="80" />
      </div>
    </div>

    <!-- 打卡弹窗 -->
    <ElDialog
      v-model="checkinDialogVisible"
      title="打卡"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="checkin-dialog-content">
        <div class="checkin-type-display">
          <span class="type-icon">{{ CheckinTypeConfig[checkinType]?.icon }}</span>
          <span class="type-label">{{ CheckinTypeConfig[checkinType]?.label }}打卡</span>
        </div>
        <ElInput
          v-model="checkinNote"
          type="textarea"
          :rows="3"
          placeholder="记录一下今天的感受（可选）"
        />
      </div>
      <template #footer>
        <ElButton @click="checkinDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="submitCheckin" :loading="submitting">完成打卡</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import { ArrowLeft, ArrowRight, CircleCheckFilled, CircleClose } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import {
    checkinApi,
    CheckinTypeConfig,
    type CheckinType,
    type TodayCheckinStatus,
    type CheckinStats,
    type UserCheckin
  } from '@/api/mood'
  import dayjs from 'dayjs'

  defineOptions({ name: 'Checkin' })

  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  const todayDate = dayjs().format('YYYY年MM月DD日')

  // 数据
  const todayStatus = reactive<TodayCheckinStatus>({
    mood: false,
    meditation: false,
    journal: false,
    general: false,
    total: 0
  })
  const stats = reactive<CheckinStats>({
    mood_streak: 0,
    meditation_streak: 0,
    journal_streak: 0,
    total_streak: 0,
    last_30_days: 0,
    type_count: { mood: 0, meditation: 0, journal: 0, general: 0 }
  })
  const historyList = ref<UserCheckin[]>([])
  const calendarData = ref<Record<string, number>>({})
  const loading = ref(false)

  // 连续打卡统计卡片数据
  const streakStats = computed(() => [
    {
      label: '情绪连续',
      value: stats.mood_streak || 0,
      icon: '😊',
      colorClass: 'text-primary',
      bgClass: 'bg-primary-light'
    },
    {
      label: '冥想连续',
      value: stats.meditation_streak || 0,
      icon: '🧘',
      colorClass: 'text-success',
      bgClass: 'bg-success-light'
    },
    {
      label: '日记连续',
      value: stats.journal_streak || 0,
      icon: '📝',
      colorClass: 'text-warning',
      bgClass: 'bg-warning-light'
    },
    {
      label: '近30天',
      value: stats.last_30_days || 0,
      icon: '📅',
      colorClass: 'text-info',
      bgClass: 'bg-info-light'
    }
  ])

  // 日历
  const currentMonth = ref(dayjs())

  const calendarDays = computed(() => {
    const start = currentMonth.value.startOf('month')
    const end = currentMonth.value.endOf('month')
    const startDay = start.day()
    const days: Array<{
      date: number
      isCurrentMonth: boolean
      isToday: boolean
      checkinCount: number
      fullDate: string
    }> = []

    // 上月末尾
    const prevMonth = start.subtract(1, 'month')
    const prevMonthDays = prevMonth.daysInMonth()
    for (let i = startDay - 1; i >= 0; i--) {
      const d = prevMonthDays - i
      const fullDate = prevMonth.date(d).format('YYYY-MM-DD')
      days.push({
        date: d,
        isCurrentMonth: false,
        isToday: false,
        checkinCount: calendarData.value[fullDate] || 0,
        fullDate
      })
    }

    // 当月
    const today = dayjs().format('YYYY-MM-DD')
    for (let i = 1; i <= end.date(); i++) {
      const fullDate = currentMonth.value.date(i).format('YYYY-MM-DD')
      days.push({
        date: i,
        isCurrentMonth: true,
        isToday: fullDate === today,
        checkinCount: calendarData.value[fullDate] || 0,
        fullDate
      })
    }

    // 下月开头
    const remaining = 42 - days.length
    const nextMonth = end.add(1, 'month')
    for (let i = 1; i <= remaining; i++) {
      const fullDate = nextMonth.date(i).format('YYYY-MM-DD')
      days.push({
        date: i,
        isCurrentMonth: false,
        isToday: false,
        checkinCount: calendarData.value[fullDate] || 0,
        fullDate
      })
    }

    return days
  })

  // 打卡弹窗
  const checkinDialogVisible = ref(false)
  const checkinType = ref<CheckinType>('mood')
  const checkinNote = ref('')
  const submitting = ref(false)

  // 方法
  const formatTime = (time: string) => dayjs(time).format('MM-DD HH:mm')

  const prevMonth = () => {
    currentMonth.value = currentMonth.value.subtract(1, 'month')
    loadCalendar()
  }

  const nextMonth = () => {
    currentMonth.value = currentMonth.value.add(1, 'month')
    loadCalendar()
  }

  const handleCheckin = (type: CheckinType) => {
    if (todayStatus[type]) {
      ElMessage.info('今天已经打过卡了')
      return
    }
    checkinType.value = type
    checkinNote.value = ''
    checkinDialogVisible.value = true
  }

  const submitCheckin = async () => {
    submitting.value = true
    try {
      await checkinApi.checkin({
        checkin_type: checkinType.value,
        note: checkinNote.value
      })
      ElMessage.success('打卡成功！')
      checkinDialogVisible.value = false
      loadTodayStatus()
      loadStats()
      loadHistory()
      loadCalendar()
    } catch (error: any) {
      ElMessage.error(error.message || '打卡失败')
    } finally {
      submitting.value = false
    }
  }

  const loadTodayStatus = async () => {
    try {
      const res = await checkinApi.getTodayStatus()
      const data = (res as any)?.data || res
      Object.assign(todayStatus, data)
    } catch (error) {
      console.error('加载今日状态失败:', error)
    }
  }

  const loadStats = async () => {
    try {
      const res = await checkinApi.getStats()
      const data = (res as any)?.data || res
      Object.assign(stats, data)
    } catch (error) {
      console.error('加载统计失败:', error)
    }
  }

  const loadHistory = async () => {
    loading.value = true
    try {
      const res = await checkinApi.getHistory({ days: 30 })
      const data = (res as any)?.data || res
      historyList.value = data || []
    } catch (error) {
      console.error('加载历史失败:', error)
    } finally {
      loading.value = false
    }
  }

  const loadCalendar = async () => {
    try {
      const res = await checkinApi.getCalendar({
        year: currentMonth.value.year(),
        month: currentMonth.value.month() + 1
      })
      const data = (res as any)?.data || res || []
      calendarData.value = {}
      data.forEach((item: any) => {
        calendarData.value[item.date] = item.count || 1
      })
    } catch (error) {
      console.error('加载日历失败:', error)
    }
  }

  onMounted(() => {
    loadTodayStatus()
    loadStats()
    loadHistory()
    loadCalendar()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/styles/variables.scss' as *;

  .checkin-page {
    --card-spacing: 20px;
  }

  // 统一卡片样式
  .card {
    margin-bottom: var(--card-spacing);
    background: var(--art-main-bg-color);
    border-radius: calc(var(--custom-radius) + 4px) !important;
  }

  // 卡片头部
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 20px 25px 15px;

    .title {
      h4 {
        font-size: 18px;
        font-weight: 500;
        color: var(--art-gray-900);
        margin: 0;
      }

      p {
        margin-top: 5px;
        font-size: 13px;
        color: var(--art-gray-600);
      }
    }
  }

  // 今日打卡卡片
  .today-card {
    padding: 0 25px 25px;

    .progress-badge {
      display: flex;
      align-items: baseline;
      padding: 8px 16px;
      background: rgb(var(--art-bg-primary));
      border-radius: 20px;

      .progress-value {
        font-size: 24px;
        font-weight: 600;
        color: rgb(var(--art-primary));
      }

      .progress-total {
        font-size: 14px;
        color: var(--art-gray-500);
      }
    }

    .checkin-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;

      @media (max-width: $device-ipad) {
        grid-template-columns: repeat(2, 1fr);
      }

      .checkin-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 20px 16px;
        background: var(--art-gray-100);
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          background: rgb(var(--art-bg-primary));
          transform: translateY(-2px);
        }

        &.completed {
          background: rgb(var(--art-bg-success));
        }

        .checkin-icon {
          font-size: 32px;
          margin-bottom: 8px;
        }

        .checkin-label {
          font-size: 14px;
          font-weight: 500;
          color: var(--art-gray-700);
          margin-bottom: 8px;
        }

        .checkin-status {
          font-size: 20px;
        }
      }
    }
  }

  // 统计卡片列表
  .stat-card-list {
    background-color: transparent !important;

    .stat-card {
      position: relative;
      box-sizing: border-box;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 120px;
      padding: 16px;
      text-align: center;

      .stat-value {
        font-size: 32px;
        font-weight: 600;
        margin-bottom: 6px;
      }

      .stat-label {
        font-size: 13px;
        color: var(--art-gray-600);
      }

      .stat-icon-wrapper {
        position: absolute;
        top: 12px;
        right: 12px;
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 8px;

        .stat-emoji {
          font-size: 18px;
        }
      }
    }
  }

  // 背景色类
  .bg-primary-light {
    background-color: rgb(var(--art-bg-primary));
  }

  .bg-success-light {
    background-color: rgb(var(--art-bg-success));
  }

  .bg-warning-light {
    background-color: rgb(var(--art-bg-warning));
  }

  .bg-info-light {
    background-color: rgb(var(--art-bg-info));
  }

  // 文字颜色类
  .text-primary {
    color: rgb(var(--art-primary)) !important;
  }

  .text-success {
    color: rgb(var(--art-success)) !important;
  }

  .text-warning {
    color: rgb(var(--art-warning)) !important;
  }

  .text-info {
    color: rgb(var(--art-info)) !important;
  }

  // 日历卡片
  .calendar-card {
    padding: 0 25px 25px;

    .month-nav {
      display: flex;
      align-items: center;
      gap: 12px;

      .month-text {
        font-size: 14px;
        font-weight: 500;
        min-width: 100px;
        text-align: center;
        color: var(--art-gray-800);
      }
    }
  }

  .calendar-wrapper {
    .calendar-header {
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      margin-bottom: 8px;

      .calendar-day-name {
        text-align: center;
        font-size: 12px;
        font-weight: 500;
        color: var(--art-gray-500);
        padding: 8px 0;
      }
    }

    .calendar-body {
      display: grid;
      grid-template-columns: repeat(7, 1fr);
      gap: 4px;

      .calendar-day {
        aspect-ratio: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background: var(--art-gray-100);
        }

        &.other-month {
          .day-number {
            color: var(--art-gray-400);
          }
        }

        &.today {
          background: rgb(var(--art-bg-primary));

          .day-number {
            color: rgb(var(--art-primary));
            font-weight: 600;
          }
        }

        &.has-checkin {
          background: rgb(var(--art-bg-success));
        }

        .day-number {
          font-size: 14px;
          color: var(--art-gray-800);
        }

        .checkin-dots {
          display: flex;
          gap: 2px;
          margin-top: 4px;

          .dot {
            width: 4px;
            height: 4px;
            border-radius: 50%;
            background: rgb(var(--art-success));
          }
        }
      }
    }
  }

  // 历史记录卡片
  .history-card {
    padding: 0 25px 25px;

    .history-content {
      max-height: 400px;
      overflow-y: auto;
    }

    .history-list {
      .history-item {
        display: flex;
        align-items: center;
        gap: 12px;
        height: 60px;
        border-bottom: 1px solid var(--art-border-color);

        &:last-child {
          border-bottom: none;
        }

        .history-icon {
          font-size: 24px;
        }

        .history-info {
          flex: 1;

          .history-type {
            font-size: 14px;
            font-weight: 500;
            color: var(--art-gray-800);
            margin-bottom: 4px;
          }

          .history-note {
            font-size: 12px;
            color: var(--art-gray-500);
          }
        }

        .history-time {
          font-size: 12px;
          color: var(--art-gray-500);
        }
      }
    }
  }

  // 打卡弹窗
  .checkin-dialog-content {
    .checkin-type-display {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 12px;
      padding: 24px;
      background: var(--art-gray-100);
      border-radius: 12px;
      margin-bottom: 20px;

      .type-icon {
        font-size: 40px;
      }

      .type-label {
        font-size: 18px;
        font-weight: 600;
        color: var(--art-gray-800);
      }
    }
  }

  // 响应式
  @media screen and (max-width: $device-phone) {
    .checkin-page {
      --card-spacing: 15px;
    }

    .stat-card-list .stat-card {
      height: 100px;

      .stat-value {
        font-size: 26px;
      }
    }
  }
</style>
