import { AppRouteRecord } from '@/types/router'

export const chatRoutes: AppRouteRecord = {
  path: '/chat',
  name: 'Chat',
  component: '/index/index',
  meta: {
    title: '聊天室',
    icon: '&#xe63a;',
    keepAlive: true
  },
  children: [
    {
      path: 'room',
      name: 'ChatRoom',
      component: '/chat/room',
      meta: {
        title: '聊天室列表',
        keepAlive: true
      }
    },
    {
      path: 'chat/video-call',
      name: 'VideoCall',
      component: '/chat/video-call',
      meta: {
        title: '视频通话',
        keepAlive: false,
        hideInMenu: true
      }
    },
    {
      path: 'room/:id',
      name: 'ChatRoomDetail',
      component: '/chat/room-detail',
      meta: {
        title: '聊天室',
        keepAlive: false,
        hideInMenu: true
      }
    }
  ]
}

export default chatRoutes
