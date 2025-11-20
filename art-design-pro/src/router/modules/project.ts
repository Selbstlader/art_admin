import { AppRouteRecord } from '@/types/router'

export const projectRoutes: AppRouteRecord = {
  path: '/project',
  name: 'Project',
  component: '/index/index',
  meta: {
    title: '项目管理',
    icon: '&#xe7b9;',
    keepAlive: true
  },
  children: [
    {
      path: 'list',
      name: 'ProjectList',
      component: '/project/index',
      meta: {
        title: '项目列表',
        keepAlive: true
      }
    },
    {
      path: 'template',
      name: 'ProjectTemplate',
      component: '/project/template',
      meta: {
        title: '项目模板',
        keepAlive: true
      }
    },
    {
      path: 'gantt/:id',
      name: 'ProjectGantt',
      component: '/project/gantt',
      meta: {
        title: '甘特图',
        keepAlive: false
      }
    }
  ]
}

export default projectRoutes
