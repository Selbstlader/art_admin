import request from '@/utils/http'

/**
 * 聊天室相关接口
 */
export const chatRoomApi = {
  // 获取聊天室列表
  getRoomList: (params: Api.Chat.ChatRoomSearchParams) => {
    return request.get<Api.Chat.ChatRoomList>({ url: '/api/chat/room/list', params })
  },

  // 获取聊天室详情
  getRoomDetail: (id: number) => {
    return request.get<Api.Chat.ChatRoomItem>({ url: `/api/chat/room/${id}` })
  },

  // 创建聊天室
  createRoom: (params: Api.Chat.CreateChatRoomRequest) => {
    return request.post<Api.Chat.ChatRoomItem>({ url: '/api/chat/room', params })
  },

  // 更新聊天室
  updateRoom: (params: Api.Chat.UpdateChatRoomRequest) => {
    return request.put({ url: '/api/chat/room', params })
  },

  // 删除聊天室
  deleteRoom: (id: number) => {
    return request.del({ url: `/api/chat/room/${id}` })
  },

  // 加入聊天室
  joinRoom: (params: Api.Chat.JoinRoomRequest) => {
    return request.post({ url: '/api/chat/room/join', params })
  },

  // 离开聊天室
  leaveRoom: (id: number) => {
    return request.post({ url: `/api/chat/room/${id}/leave` })
  },

  // 获取聊天室成员列表
  getRoomMembers: (id: number) => {
    return request.get<Api.Chat.ChatRoomMemberItem[]>({ url: `/api/chat/room/${id}/members` })
  },

  // 获取在线用户列表
  getOnlineUsers: (id: number) => {
    return request.get<Api.Chat.OnlineUser[]>({ url: `/api/chat/room/${id}/online-users` })
  },

  // 获取我的聊天室列表
  getMyRooms: () => {
    return request.get<Api.Chat.ChatRoomItem[]>({ url: '/api/chat/my-rooms' })
  }
}

/**
 * 聊天消息相关接口
 */
export const chatMessageApi = {
  // 发送消息
  sendMessage: (params: Api.Chat.SendMessageRequest) => {
    return request.post<Api.Chat.ChatMessageItem>({ url: '/api/chat/message', params })
  },

  // 获取消息列表
  getMessageList: (params: Api.Chat.ChatMessageSearchParams) => {
    return request.get<{
      data: Api.Chat.ChatMessageItem[]
      total: number
      page: number
      pageSize: number
    }>({
      url: '/api/chat/message/list',
      params,
      _fullResponse: true // 返回完整响应，包含分页信息
    })
  },

  // 撤回消息
  recallMessage: (params: Api.Chat.RecallMessageRequest) => {
    return request.post({ url: '/api/chat/message/recall', params })
  }
}
