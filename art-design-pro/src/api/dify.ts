import request from '@/utils/http'
import { useUserStore } from '@/store/modules/user'

// Dify 知识库相关接口
export const difyDatasetApi = {
  // 获取知识库列表
  getDatasetList: (params: any) => {
    return request.get({ url: '/api/dify/dataset/list', params })
  },
  // 获取知识库详情
  getDatasetDetail: (id: string) => {
    return request.get({ url: `/api/dify/dataset/${id}` })
  },
  // 上传文件到知识库
  uploadFile: (datasetId: string, file: File) => {
    const formData = new FormData()
    formData.append('dataset_id', datasetId)
    formData.append('file', file)
    return request.post({
      url: '/api/dify/dataset/upload',
      params: formData,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  // 删除知识库
  deleteDataset: (id: string) => {
    return request.del({ url: `/api/dify/dataset/${id}` })
  }
}

// Dify AI 对话相关接口
export const difyChatApi = {
  // AI 对话(非流式)
  chat: (params: {
    query: string
    user: string
    conversation_id?: string
    inputs?: Record<string, any>
    files?: any[]
    auto_generate_name?: boolean
    tag_id?: number
  }) => {
    return request.post({ url: '/api/dify/chat', params })
  },

  // AI 对话(流式) - 返回 EventSource
  chatStream: (params: {
    query: string
    user: string
    conversation_id?: string
    inputs?: Record<string, any>
    files?: any[]
    auto_generate_name?: boolean
    tag_id?: number
  }) => {
    // 流式接口需要特殊处理,返回 EventSource 对象
    // 从 pinia store 获取 token
    const userStore = useUserStore()
    const token = userStore.accessToken || ''
    const url = `${import.meta.env.VITE_API_BASE_URL}/api/dify/chat/stream`

    return fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(params)
    })
  },

  // 停止消息生成
  stopMessage: (taskId: string, user: string) => {
    return request.post({
      url: `/api/dify/chat/stop/${taskId}`,
      params: { user }
    })
  },

  // 获取建议问题
  getSuggestedQuestions: (messageId: string, user: string) => {
    return request.get({
      url: `/api/dify/messages/${messageId}/suggested`,
      params: { user }
    })
  }
}

// Dify 会话管理相关接口
export const difyConversationApi = {
  // 获取会话列表
  getConversations: (params: { user: string; last_id?: string; limit?: number }) => {
    return request.get({ url: '/api/dify/conversations', params })
  },

  // 获取会话消息历史
  getMessages: (params: {
    conversation_id: string
    user: string
    first_id?: string
    limit?: number
  }) => {
    return request.get({ url: '/api/dify/messages', params })
  },

  // 删除会话
  deleteConversation: (conversationId: string, user: string) => {
    return request.del({
      url: `/api/dify/conversations/${conversationId}`,
      params: { user }
    })
  },

  // 重命名会话
  renameConversation: (conversationId: string, params: { user: string; name: string }) => {
    return request.post({
      url: `/api/dify/conversations/${conversationId}/name`,
      params
    })
  }
}
