import { AppRouteRecord } from '@/types/router'

export const learningRoutes: AppRouteRecord = {
  path: '/learning',
  name: 'Learning',
  component: '/index/index',
  meta: {
    title: 'AI学习系统',
    icon: '&#xe7ae;',
    keepAlive: true
  },
  children: [
    {
      path: 'my',
      name: 'LearningMy',
      component: '/learning/my/index',
      meta: {
        title: '生成教材',
        keepAlive: false
      }
    },
    {
      path: 'list',
      name: 'LearningList',
      component: '/learning/list/index',
      meta: {
        title: '我的课程',
        keepAlive: true
      }
    },
    {
      path: 'courses',
      name: 'LearningCourses',
      component: '/learning/courses/index',
      meta: {
        title: '课程列表',
        keepAlive: true,
        hidden: true
      }
    },
    {
      path: 'detail/:id',
      name: 'LearningDetail',
      component: '/learning/detail/index',
      meta: {
        title: '课程详情',
        keepAlive: false,
        hidden: true
      }
    }
  ]
}

export default learningRoutes
