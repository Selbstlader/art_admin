import { AppRouteRecord } from '@/types/router'

export const difyRoutes: AppRouteRecord = {
  path: '/dify',
  name: 'Dify',
  component: '/index/index',
  meta: {
    title: 'AI助手',
    icon: '&#xe7ae;',
    keepAlive: true
  },
  children: [
    {
      path: 'chat',
      name: 'DifyChat',
      component: '/dify/chat',
      meta: {
        title: 'AI对话',
        keepAlive: true
      }
    },
    {
      path: 'knowledge',
      name: 'DifyKnowledge',
      component: '/dify/knowledge',
      meta: {
        title: '知识库管理',
        keepAlive: true
      }
    }
  ]
}

export default difyRoutes
