/***
 * Designer Chat API
 * 设计师AI对话接口
 * Requirements: 5.1, 5.2, 5.4
 ***/
import request from '@/utils/http'

/**
 * 对话请求参数
 * Chat request parameters
 */
export interface ChatRequest {
  query: string
  projectId?: number
  sessionId?: string
}

/**
 * 对话响应
 * Chat response
 */
export interface ChatResponse {
  answer: string
  sessionId: string
  messageId: number
  tokensUsed: number
  responseTime: number
}

/**
 * 对话消息
 * Chat message
 */
export interface ChatMessage {
  id: number
  sessionId: string
  projectId: number
  role: 'user' | 'assistant'
  content: string
  createdAt: string
}

/**
 * 对话会话
 * Chat session
 */
export interface ChatSession {
  id: string
  projectId: number
  userId: number
  title: string
  createdAt: string
  updatedAt: string
}

/**
 * 设计师AI对话API
 * Designer AI chat API
 */
export const designerChatApi = {
  /**
   * 发送对话消息
   * Send chat message
   * Requirements: 5.1 - 基于当前项目信息提供上下文相关的回答
   * Requirements: 5.2 - 引用相关国家标准或行业规范进行回答
   */
  chat: (params: ChatRequest) => {
    return request.post<ChatResponse>({
      url: '/api/designer/chat',
      params
    })
  },

  /**
   * 获取对话历史
   * Get chat history
   * Requirements: 5.4 - 保存对话历史记录
   */
  getHistory: (params: { sessionId?: string; page?: number; pageSize?: number }) => {
    return request.get<{
      data: ChatMessage[]
      total: number
      page: number
      pageSize: number
    }>({
      url: '/api/designer/chat/history',
      params
    })
  },

  /**
   * 获取会话列表
   * Get session list
   */
  getSessions: (projectId?: number) => {
    return request.get<ChatSession[]>({
      url: '/api/designer/chat/sessions',
      params: { projectId }
    })
  },

  /**
   * 删除会话
   * Delete session
   */
  deleteSession: (sessionId: string) => {
    return request.del({
      url: `/api/designer/chat/sessions/${sessionId}`
    })
  },

  /**
   * 导出对话历史
   * Export chat history
   * Requirements: 5.4 - 支持导出为文档格式
   */
  exportHistory: (sessionId: string, format: 'md' | 'json' = 'md') => {
    return request.get<{
      content: string
      format: string
      filename: string
    }>({
      url: `/api/designer/chat/export/${sessionId}`,
      params: { format }
    })
  }
}
