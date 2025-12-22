/***
 * Designer Chat API
 * 设计师AI对话接口
 * Requirements: 5.1, 5.2, 5.4
 ***/
import request from '@/utils/http'

/*** API响应包装类型 / API response wrapper type ***/
interface ApiResponse<T = unknown> {
  code: number
  msg?: string
  data?: T
}

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
 * 对话历史响应
 * Chat history response
 */
export interface ChatHistoryResponse {
  data: ChatMessage[]
  total: number
  page: number
  pageSize: number
}

/**
 * 导出历史响应
 * Export history response
 */
export interface ExportHistoryResponse {
  content: string
  format: string
  filename: string
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
  chat: (data: ChatRequest): Promise<ApiResponse<ChatResponse>> => {
    return request.post<ApiResponse<ChatResponse>>({
      url: '/api/designer/chat',
      data,
      _fullResponse: true,
      timeout: 100000
    })
  },

  /**
   * 发送流式对话消息
   * Send streaming chat message
   * Requirements: 5.1, 5.2 - 流式输出支持
   */
  chatStream: async (
    data: ChatRequest,
    onChunk: (content: string) => void,
    onDone: (result: { sessionId: string; messageId: number; tokensUsed: number }) => void,
    onError: (error: string) => void
  ): Promise<void> => {
    const { VITE_API_URL } = import.meta.env
    // 从 pinia store 获取 token / Get token from pinia store
    const { useUserStore } = await import('@/store/modules/user')
    const userStore = useUserStore()
    const token = userStore.accessToken || ''

    try {
      const response = await fetch(`${VITE_API_URL}/api/designer/chat/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: token.startsWith('Bearer ') ? token : `Bearer ${token}`
        },
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        const errorData = await response.json()
        onError(errorData.msg || '请求失败')
        return
      }

      const reader = response.body?.getReader()
      if (!reader) {
        onError('不支持流式读取')
        return
      }

      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6).trim()
            if (data === '[DONE]') {
              continue
            }
            // 尝试解析JSON（最终结果或错误）
            if (data.startsWith('{')) {
              try {
                const json = JSON.parse(data)
                if (json.error) {
                  onError(json.error)
                } else if (json.sessionId) {
                  onDone({
                    sessionId: json.sessionId,
                    messageId: json.messageId,
                    tokensUsed: json.tokensUsed
                  })
                }
              } catch {
                // 不是JSON，作为普通内容处理
                onChunk(data)
              }
            } else {
              // 普通文本内容
              onChunk(data)
            }
          }
        }
      }
    } catch (error: any) {
      onError(error.message || '网络错误')
    }
  },

  /**
   * 获取对话历史
   * Get chat history
   * Requirements: 5.4 - 保存对话历史记录
   */
  getHistory: (params: {
    sessionId?: string
    page?: number
    pageSize?: number
  }): Promise<ApiResponse<ChatHistoryResponse>> => {
    return request.get<ApiResponse<ChatHistoryResponse>>({
      url: '/api/designer/chat/history',
      params,
      _fullResponse: true
    })
  },

  /**
   * 获取会话列表
   * Get session list
   */
  getSessions: (projectId?: number): Promise<ApiResponse<ChatSession[]>> => {
    return request.get<ApiResponse<ChatSession[]>>({
      url: '/api/designer/chat/sessions',
      params: { projectId },
      _fullResponse: true
    })
  },

  /**
   * 删除会话
   * Delete session
   */
  deleteSession: (sessionId: string): Promise<ApiResponse<null>> => {
    return request.del<ApiResponse<null>>({
      url: `/api/designer/chat/sessions/${sessionId}`,
      _fullResponse: true
    })
  },

  /**
   * 导出对话历史
   * Export chat history
   * Requirements: 5.4 - 支持导出为文档格式
   */
  exportHistory: (
    sessionId: string,
    format: 'md' | 'json' = 'md'
  ): Promise<ApiResponse<ExportHistoryResponse>> => {
    return request.get<ApiResponse<ExportHistoryResponse>>({
      url: `/api/designer/chat/export/${sessionId}`,
      params: { format },
      _fullResponse: true
    })
  }
}
