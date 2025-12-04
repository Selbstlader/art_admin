import { AppRouteRecord } from '@/types/router'

export const travelRoutes: AppRouteRecord = {
  path: '/travel',
  name: 'Travel',
  component: '/index/index',
  meta: {
    title: 'menus.travel.title',
    icon: '&#xe6b8;',
    roles: ['R_SUPER', 'R_ADMIN', 'R_USER']
  },
  children: [
    {
      path: 'map',
      name: 'TravelMap',
      component: '/travel/map',
      meta: {
        title: 'menus.travel.map',
        keepAlive: true
      }
    },
    {
      path: 'roadbook',
      name: 'RoadbookList',
      component: '/travel/roadbook/list',
      meta: {
        title: 'menus.travel.roadbookList',
        keepAlive: true,
        authList: [
          { title: '创建', authMark: 'create' },
          { title: '编辑', authMark: 'edit' },
          { title: '删除', authMark: 'delete' }
        ]
      }
    },
    {
      path: 'roadbook/detail/:id',
      name: 'RoadbookDetail',
      component: '/travel/roadbook/detail',
      meta: {
        title: 'menus.travel.roadbookDetail',
        isHide: true,
        keepAlive: true,
        activePath: '/travel/roadbook'
      }
    },
    {
      path: 'roadbook/editor',
      name: 'RoadbookEditor',
      component: '/travel/roadbook/editor',
      meta: {
        title: 'menus.travel.roadbookEditor',
        isHide: true,
        keepAlive: true,
        activePath: '/travel/roadbook'
      }
    },
    {
      path: 'roadbook/editor/:id',
      name: 'RoadbookEditorEdit',
      component: '/travel/roadbook/editor',
      meta: {
        title: 'menus.travel.roadbookEditor',
        isHide: true,
        keepAlive: true,
        activePath: '/travel/roadbook'
      }
    },
    {
      path: 'explore',
      name: 'RoadbookExplore',
      component: '/travel/roadbook/explore',
      meta: {
        title: 'menus.travel.explore',
        keepAlive: true
      }
    },
    {
      path: 'ai',
      name: 'TravelAI',
      component: '/travel/ai',
      meta: {
        title: 'menus.travel.ai',
        keepAlive: true
      }
    },
    {
      path: 'template',
      name: 'TravelTemplate',
      component: '/travel/template',
      meta: {
        title: 'menus.travel.template',
        keepAlive: true
      }
    }
  ]
}

export default travelRoutes
