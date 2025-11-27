<template>
  <div class="mood-dashboard">
    <!-- 今日状态卡片 -->
    <ElRow :gutter="20" class="stat-card-list">
      <ElCol v-for="(stat, index) in statCards" :key="index" :sm="12" :md="6" :lg="6">
        <div class="card art-custom-card stat-card">
          <span class="des subtitle">{{ stat.label }}</span>
          <div class="number box-title">{{ stat.value }}</div>
          <div class="change-box">
            <span class="change-text">{{ stat.subLabel }}</span>
            <span class="change" :class="stat.changeClass">{{ stat.change }}</span>
          </div>
          <div class="stat-icon-wrapper" :class="stat.iconClass">
            <span class="stat-emoji">{{ stat.icon }}</span>
          </div>
        </div>
      </ElCol>
    </ElRow>

    <!-- 快捷入口 -->
    <div class="card art-custom-card quick-entry-card">
      <div class="card-header">
        <div class="title">
          <h4 class="box-title">快捷入口</h4>
          <p class="subtitle">快速访问常用功能</p>
        </div>
      </div>
      <div class="quick-actions">
        <div class="action-item" @click="goTo('/mood/mood/record')">
          <div class="action-icon-wrapper bg-primary-light">
            <span class="action-emoji">📊</span>
          </div>
          <div class="action-info">
            <span class="action-label">记录情绪</span>
            <span class="action-desc">记录当前心情</span>
          </div>
        </div>
        <div class="action-item" @click="goTo('/mood/mood/meditation')">
          <div class="action-icon-wrapper bg-success-light">
            <span class="action-emoji">🧘</span>
          </div>
          <div class="action-info">
            <span class="action-label">开始冥想</span>
            <span class="action-desc">放松身心</span>
          </div>
        </div>
        <div class="action-item" @click="goTo('/mood/mood/journal')">
          <div class="action-icon-wrapper bg-warning-light">
            <span class="action-emoji">📝</span>
          </div>
          <div class="action-info">
            <span class="action-label">写日记</span>
            <span class="action-desc">记录生活点滴</span>
          </div>
        </div>
        <div class="action-item" @click="goTo('/mood/mood/checkin')">
          <div class="action-icon-wrapper bg-info-light">
            <span class="action-emoji">✅</span>
          </div>
          <div class="action-info">
            <span class="action-label">每日打卡</span>
            <span class="action-desc">保持好习惯</span>
          </div>
        </div>
      </div>
    </div>

    <ElRow :gutter="20">
      <!-- 情绪趋势 -->
      <ElCol :xs="24" :lg="12">
        <div class="card art-custom-card trend-card">
          <div class="card-header">
            <div class="title">
              <h4 class="box-title">情绪趋势</h4>
              <p class="subtitle">近7天情绪变化</p>
            </div>
            <ElButton link type="primary" @click="goTo('/mood/calendar')">查看更多</ElButton>
          </div>
          <div class="trend-content" v-loading="trendLoading">
            <div v-if="moodTrends.length > 0" class="trend-list">
              <div v-for="trend in moodTrends" :key="trend.date" class="trend-item">
                <div class="trend-date">{{ formatDate(trend.date) }}</div>
                <div class="trend-mood">
                  <span class="mood-emoji">{{
                    MoodTypeConfig[trend.mood_type]?.icon || '😐'
                  }}</span>
                  <span class="mood-label">{{
                    MoodTypeConfig[trend.mood_type]?.label || '未知'
                  }}</span>
                </div>
                <div class="trend-intensity">
                  <ElProgress
                    :percentage="trend.avg_intensity * 10"
                    :stroke-width="8"
                    :show-text="false"
                    :color="getIntensityColor(trend.avg_intensity)"
                  />
                  <span class="intensity-value">{{ trend.avg_intensity.toFixed(1) }}</span>
                </div>
              </div>
            </div>
            <ElEmpty v-else description="暂无情绪记录" :image-size="80" />
          </div>
        </div>
      </ElCol>

      <!-- 今日打卡状态 -->
      <ElCol :xs="24" :lg="12">
        <div class="card art-custom-card checkin-card">
          <div class="card-header">
            <div class="title">
              <h4 class="box-title">今日打卡</h4>
              <p class="subtitle"
                >完成 <span class="text-success">{{ todayCheckin?.total || 0 }}/4</span></p
              >
            </div>
            <ElButton link type="primary" @click="goTo('/mood/checkin')">去打卡</ElButton>
          </div>
          <div class="checkin-content" v-loading="checkinLoading">
            <div class="checkin-grid">
              <div
                v-for="(item, key) in CheckinTypeConfig"
                :key="key"
                class="checkin-item"
                :class="{ completed: todayCheckin?.[key as CheckinType] }"
              >
                <div class="checkin-icon">{{ item.icon }}</div>
                <div class="checkin-label">{{ item.label }}</div>
                <div class="checkin-status">
                  <ElIcon v-if="todayCheckin?.[key as CheckinType]" color="var(--el-color-success)"
                    ><Check
                  /></ElIcon>
                  <ElIcon v-else color="var(--art-gray-400)"><Close /></ElIcon>
                </div>
              </div>
            </div>
            <div class="streak-info">
              <div class="streak-item">
                <span class="streak-value">{{ checkinStats?.mood_streak || 0 }}</span>
                <span class="streak-label">情绪连续</span>
              </div>
              <div class="streak-item">
                <span class="streak-value">{{ checkinStats?.meditation_streak || 0 }}</span>
                <span class="streak-label">冥想连续</span>
              </div>
              <div class="streak-item">
                <span class="streak-value">{{ checkinStats?.journal_streak || 0 }}</span>
                <span class="streak-label">日记连续</span>
              </div>
            </div>
          </div>
        </div>
      </ElCol>
    </ElRow>

    <ElRow :gutter="20">
      <!-- 目标进度 -->
      <ElCol :xs="24" :lg="12">
        <div class="card art-custom-card goals-card">
          <div class="card-header">
            <div class="title">
              <h4 class="box-title">目标进度</h4>
              <p class="subtitle"
                >已完成 <span class="text-success">{{ completedGoals }}/{{ totalGoals }}</span></p
              >
            </div>
            <ElButton link type="primary" @click="goTo('/mood/goals')">管理目标</ElButton>
          </div>
          <div class="goals-content" v-loading="goalsLoading">
            <div v-if="goals.length > 0" class="goals-list">
              <div v-for="goal in goals" :key="goal.id" class="goal-item">
                <div class="goal-icon">{{ GoalTypeConfig[goal.goal_type]?.icon }}</div>
                <div class="goal-info">
                  <div class="goal-name">{{ GoalTypeConfig[goal.goal_type]?.label }}</div>
                  <div class="goal-progress-bar">
                    <ElProgress
                      :percentage="Math.min((goal.current_value / goal.target_value) * 100, 100)"
                      :stroke-width="6"
                      :show-text="false"
                    />
                  </div>
                </div>
                <div class="goal-stats">
                  <span class="goal-current">{{ goal.current_value }}</span>
                  <span class="goal-divider">/</span>
                  <span class="goal-target">{{ goal.target_value }}</span>
                </div>
              </div>
            </div>
            <ElEmpty v-else description="暂无目标，去设置一个吧" :image-size="80" />
          </div>
        </div>
      </ElCol>

      <!-- 继续冥想 -->
      <ElCol :xs="24" :lg="12">
        <div class="card art-custom-card meditation-card">
          <div class="card-header">
            <div class="title">
              <h4 class="box-title">继续冥想</h4>
              <p class="subtitle">未完成的冥想练习</p>
            </div>
            <ElButton link type="primary" @click="goTo('/mood/meditation')">更多内容</ElButton>
          </div>
          <div class="meditation-content" v-loading="meditationLoading">
            <div v-if="continuePlaying.length > 0" class="continue-list">
              <div
                v-for="item in continuePlaying"
                :key="item.id"
                class="continue-item"
                @click="goToPlayer(item.content_id)"
              >
                <div class="continue-progress">
                  <ElProgress
                    type="circle"
                    :percentage="item.progress"
                    :width="48"
                    :stroke-width="4"
                  />
                </div>
                <div class="continue-info">
                  <div class="continue-title">{{ item.title || '冥想内容' }}</div>
                  <div class="continue-meta">上次播放到 {{ formatDuration(item.duration) }}</div>
                </div>
                <ElButton type="primary" size="small" round>继续</ElButton>
              </div>
            </div>
            <ElEmpty v-else description="暂无未完成的冥想" :image-size="80" />
          </div>
        </div>
      </ElCol>
    </ElRow>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { Check, Close } from '@element-plus/icons-vue'
  import {
    moodRecordApi,
    checkinApi,
    goalApi,
    meditationApi,
    MoodTypeConfig,
    CheckinTypeConfig,
    GoalTypeConfig,
    type MoodTrend,
    type TodayCheckinStatus,
    type CheckinStats,
    type UserGoal,
    type CheckinType,
    type MoodRecord
  } from '@/api/mood'
  import dayjs from 'dayjs'

  defineOptions({ name: 'MoodDashboard' })

  const router = useRouter()

  // 数据状态
  const todayMood = ref<MoodRecord | null>(null)
  const todayCheckin = ref<TodayCheckinStatus | null>(null)
  const checkinStats = ref<CheckinStats | null>(null)
  const moodTrends = ref<MoodTrend[]>([])
  const goals = ref<UserGoal[]>([])
  const continuePlaying = ref<any[]>([])

  // 加载状态
  const trendLoading = ref(false)
  const checkinLoading = ref(false)
  const goalsLoading = ref(false)
  const meditationLoading = ref(false)

  // 计算属性
  const completedGoals = computed(() => {
    return goals.value.filter((g) => g.current_value >= g.target_value).length
  })

  const totalGoals = computed(() => goals.value.length)

  // 统计卡片数据
  const statCards = computed(() => [
    {
      label: '今日情绪',
      value: todayMood.value?.mood_type
        ? MoodTypeConfig[todayMood.value.mood_type]?.label
        : '未记录',
      icon: todayMood.value?.mood_type ? MoodTypeConfig[todayMood.value.mood_type]?.icon : '😊',
      subLabel: '情绪强度',
      change: todayMood.value?.intensity ? `${todayMood.value.intensity}/10` : '-',
      changeClass: 'text-primary',
      iconClass: 'bg-primary-light'
    },
    {
      label: '今日打卡',
      value: `${todayCheckin.value?.total || 0}/4`,
      icon: '✓',
      subLabel: '完成进度',
      change: todayCheckin.value?.total === 4 ? '已完成' : '进行中',
      changeClass: todayCheckin.value?.total === 4 ? 'text-success' : 'text-warning',
      iconClass: 'bg-success-light'
    },
    {
      label: '连续打卡',
      value: `${checkinStats.value?.total_streak || 0}天`,
      icon: '🔥',
      subLabel: '坚持就是胜利',
      change: checkinStats.value?.total_streak ? '+1' : '开始吧',
      changeClass: checkinStats.value?.total_streak ? 'text-success' : 'text-info',
      iconClass: 'bg-warning-light'
    },
    {
      label: '目标完成',
      value: `${completedGoals.value}/${totalGoals.value}`,
      icon: '🎯',
      subLabel: '目标达成率',
      change:
        totalGoals.value > 0
          ? `${Math.round((completedGoals.value / totalGoals.value) * 100)}%`
          : '-',
      changeClass: completedGoals.value > 0 ? 'text-success' : 'text-info',
      iconClass: 'bg-info-light'
    }
  ])

  // 方法
  const goTo = (path: string) => {
    router.push(path)
  }

  const goToPlayer = (contentId: number) => {
    router.push(`/mood/meditation/player/${contentId}`)
  }

  const formatDate = (date: string) => {
    return dayjs(date).format('MM/DD')
  }

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  const getIntensityColor = (intensity: number) => {
    if (intensity <= 3) return '#67C23A'
    if (intensity <= 6) return '#E6A23C'
    return '#F56C6C'
  }

  // 加载数据
  const loadTodayMood = async () => {
    try {
      const today = dayjs().format('YYYY-MM-DD')
      const res = await moodRecordApi.getByDate(today)
      if (res && Array.isArray(res) && res.length > 0) {
        todayMood.value = res[0]
      }
    } catch (error) {
      console.error('加载今日情绪失败:', error)
    }
  }

  const loadMoodTrends = async () => {
    trendLoading.value = true
    try {
      const res = await moodRecordApi.getTrends({ days: 7 })
      moodTrends.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载情绪趋势失败:', error)
    } finally {
      trendLoading.value = false
    }
  }

  const loadCheckinStatus = async () => {
    checkinLoading.value = true
    try {
      const [statusRes, statsRes] = await Promise.all([
        checkinApi.getTodayStatus(),
        checkinApi.getStats()
      ])
      todayCheckin.value = statusRes as TodayCheckinStatus
      checkinStats.value = statsRes as CheckinStats
    } catch (error) {
      console.error('加载打卡状态失败:', error)
    } finally {
      checkinLoading.value = false
    }
  }

  const loadGoals = async () => {
    goalsLoading.value = true
    try {
      const res = await goalApi.getList()
      goals.value = ((res as any)?.data || res || []).filter((g: UserGoal) => g.is_active)
    } catch (error) {
      console.error('加载目标失败:', error)
    } finally {
      goalsLoading.value = false
    }
  }

  const loadContinuePlaying = async () => {
    meditationLoading.value = true
    try {
      const res = await meditationApi.getContinuePlaying()
      continuePlaying.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载继续播放失败:', error)
    } finally {
      meditationLoading.value = false
    }
  }

  onMounted(() => {
    loadTodayMood()
    loadMoodTrends()
    loadCheckinStatus()
    loadGoals()
    loadContinuePlaying()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/styles/variables.scss' as *;

  .mood-dashboard {
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

  // 统计卡片列表
  .stat-card-list {
    background-color: transparent !important;

    .stat-card {
      position: relative;
      box-sizing: border-box;
      display: flex;
      flex-direction: column;
      justify-content: center;
      width: 100%;
      height: 140px;
      padding: 0 20px;
      transition: all 0.3s ease;

      .des {
        display: block;
        height: 14px;
        font-size: 14px;
        line-height: 14px;
        color: var(--art-gray-600);
      }

      .number {
        display: block;
        margin-top: 10px;
        font-size: 28px;
        font-weight: 500;
        color: var(--art-gray-900);
      }

      .change-box {
        display: flex;
        align-items: center;
        margin-top: 10px;

        .change-text {
          font-size: 13px;
          color: var(--art-text-gray-600);
        }

        .change {
          margin-left: 5px;
          font-size: 13px;
          font-weight: bold;
        }
      }

      .stat-icon-wrapper {
        position: absolute;
        top: 0;
        right: 20px;
        bottom: 0;
        width: 52px;
        height: 52px;
        margin: auto;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 12px;

        .stat-emoji {
          font-size: 24px;
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

  // 快捷入口卡片
  .quick-entry-card {
    padding: 0 25px 25px;

    .quick-actions {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
      margin-top: 10px;

      @media (max-width: $device-ipad) {
        grid-template-columns: repeat(2, 1fr);
      }

      .action-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 16px;
        border-radius: 12px;
        background: var(--art-gray-100);
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          background: rgb(var(--art-bg-primary));
          transform: translateY(-2px);
        }

        .action-icon-wrapper {
          width: 48px;
          height: 48px;
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 10px;

          .action-emoji {
            font-size: 24px;
          }
        }

        .action-info {
          display: flex;
          flex-direction: column;

          .action-label {
            font-size: 14px;
            font-weight: 500;
            color: var(--art-gray-800);
          }

          .action-desc {
            font-size: 12px;
            color: var(--art-gray-500);
            margin-top: 2px;
          }
        }
      }
    }
  }

  // 趋势卡片
  .trend-card,
  .checkin-card,
  .goals-card,
  .meditation-card {
    height: 420px;
    padding: 0 25px 25px;
  }

  .trend-content,
  .checkin-content,
  .goals-content,
  .meditation-content {
    height: calc(100% - 80px);
    overflow-y: auto;
  }

  // 趋势列表
  .trend-list {
    .trend-item {
      display: flex;
      align-items: center;
      height: 48px;
      border-bottom: 1px solid var(--art-border-color);

      &:last-child {
        border-bottom: none;
      }

      .trend-date {
        width: 50px;
        font-size: 13px;
        color: var(--art-gray-600);
      }

      .trend-mood {
        width: 90px;
        display: flex;
        align-items: center;
        gap: 6px;

        .mood-emoji {
          font-size: 18px;
        }

        .mood-label {
          font-size: 13px;
          color: var(--art-gray-700);
        }
      }

      .trend-intensity {
        flex: 1;
        display: flex;
        align-items: center;
        gap: 10px;

        .el-progress {
          flex: 1;
        }

        .intensity-value {
          width: 28px;
          font-size: 13px;
          font-weight: 600;
          color: var(--art-gray-800);
        }
      }
    }
  }

  // 打卡内容
  .checkin-content {
    .checkin-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 10px;
      margin-bottom: 16px;

      @media (max-width: $device-ipad) {
        grid-template-columns: repeat(2, 1fr);
      }

      .checkin-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 14px 8px;
        border-radius: 10px;
        background: var(--art-gray-100);
        transition: all 0.3s;

        &.completed {
          background: rgb(var(--art-bg-success));
        }

        .checkin-icon {
          font-size: 24px;
          margin-bottom: 6px;
        }

        .checkin-label {
          font-size: 12px;
          color: var(--art-gray-700);
          margin-bottom: 6px;
        }

        .checkin-status {
          font-size: 16px;
        }
      }
    }

    .streak-info {
      display: flex;
      justify-content: space-around;
      padding: 14px;
      background: var(--art-gray-100);
      border-radius: 10px;

      .streak-item {
        text-align: center;

        .streak-value {
          display: block;
          font-size: 22px;
          font-weight: 600;
          color: rgb(var(--art-primary));
        }

        .streak-label {
          display: block;
          font-size: 12px;
          color: var(--art-gray-500);
          margin-top: 4px;
        }
      }
    }
  }

  // 目标列表
  .goals-list {
    .goal-item {
      display: flex;
      align-items: center;
      gap: 12px;
      height: 60px;
      border-bottom: 1px solid var(--art-border-color);

      &:last-child {
        border-bottom: none;
      }

      .goal-icon {
        font-size: 24px;
      }

      .goal-info {
        flex: 1;

        .goal-name {
          font-size: 14px;
          font-weight: 500;
          color: var(--art-gray-800);
          margin-bottom: 6px;
        }

        .goal-progress-bar {
          .el-progress {
            width: 100%;
          }
        }
      }

      .goal-stats {
        font-size: 13px;
        white-space: nowrap;

        .goal-current {
          font-weight: 600;
          color: rgb(var(--art-primary));
        }

        .goal-divider {
          color: var(--art-gray-400);
          margin: 0 2px;
        }

        .goal-target {
          color: var(--art-gray-600);
        }
      }
    }
  }

  // 继续冥想列表
  .continue-list {
    .continue-item {
      display: flex;
      align-items: center;
      gap: 14px;
      height: 70px;
      border-bottom: 1px solid var(--art-border-color);
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        background: var(--art-gray-100);
        margin: 0 -10px;
        padding: 0 10px;
        border-radius: 8px;
      }

      &:last-child {
        border-bottom: none;
      }

      .continue-info {
        flex: 1;

        .continue-title {
          font-size: 14px;
          font-weight: 500;
          color: var(--art-gray-800);
          margin-bottom: 4px;
        }

        .continue-meta {
          font-size: 12px;
          color: var(--art-gray-500);
        }
      }
    }
  }

  // 响应式
  @media screen and (max-width: $device-phone) {
    .mood-dashboard {
      --card-spacing: 15px;
    }

    .stat-card-list .stat-card {
      height: 120px;

      .number {
        font-size: 24px;
      }
    }
  }
</style>
