import { AppRouteRecord } from '@/types/router'

export const moodRoutes: AppRouteRecord = {
  path: '/mood',
  name: 'Mood',
  component: '/index/index',
  meta: {
    title: '情绪健康',
    icon: '&#xe7b9;',
    keepAlive: true
  },
  children: [
    {
      path: 'dashboard',
      name: 'MoodDashboard',
      component: '/mood/dashboard',
      meta: {
        title: '健康仪表盘',
        keepAlive: true
      }
    },
    {
      path: 'record',
      name: 'MoodRecord',
      component: '/mood/record',
      meta: {
        title: '情绪记录',
        keepAlive: true
      }
    },
    {
      path: 'calendar',
      name: 'MoodCalendar',
      component: '/mood/calendar',
      meta: {
        title: '情绪日历',
        keepAlive: true
      }
    },
    {
      path: 'meditation',
      name: 'Meditation',
      component: '/mood/meditation',
      meta: {
        title: '冥想练习',
        keepAlive: true
      }
    },
    {
      path: 'meditation/player/:id',
      name: 'MeditationPlayer',
      component: '/mood/meditation/player',
      meta: {
        title: '冥想播放',
        keepAlive: false
      }
    },
    {
      path: 'journal',
      name: 'Journal',
      component: '/mood/journal',
      meta: {
        title: '心情日记',
        keepAlive: true
      }
    },
    {
      path: 'goals',
      name: 'Goals',
      component: '/mood/goals',
      meta: {
        title: '目标管理',
        keepAlive: true
      }
    },
    {
      path: 'checkin',
      name: 'Checkin',
      component: '/mood/checkin',
      meta: {
        title: '每日打卡',
        keepAlive: true
      }
    }
  ]
}

export default moodRoutes
