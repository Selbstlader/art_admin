import { AppRouteRecord } from '@/types/router'

export const xhsRoutes: AppRouteRecord = {
  path: '/xhs',
  name: 'XHS',
  component: '/index/index',
  meta: {
    title: '小红书助手',
    icon: '&#xe726;',
    keepAlive: true
  },
  children: [
    {
      path: 'summary',
      name: 'XHSSummary',
      component: '/xhs/summary',
      meta: {
        title: 'AI笔记总结',
        keepAlive: true
      }
    },
    {
      path: 'history',
      name: 'XHSHistory',
      component: '/xhs/history',
      meta: {
        title: '历史记录',
        keepAlive: true
      }
    },
    {
      path: 'favorites',
      name: 'XHSFavorites',
      component: '/xhs/favorites',
      meta: {
        title: '我的收藏',
        keepAlive: true
      }
    }
  ]
}

export default xhsRoutes
