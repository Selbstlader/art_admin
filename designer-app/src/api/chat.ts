/**
 * AI 对话相关 API（仅 App 端）
 * AI Chat related API (App only)
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, ChatRequest, ChatResponse, ChatMessage, ChatSession, PageParams, StreamCallbacks } from '@/types/api'
import { AI_REQUEST_TIMEOUT } from '@/types/common'

/*** Chat API - handles AI chat operations (App only) ***/
export const chatApi = {
  /*** 发送消息 - Send message ***/
  sendMessage(data: ChatRequest): Promise<ApiResponse<ChatResponse>> {
    return request({
      url: '/api/designer/chat/send',
      method: 'POST',
      data,
      timeout: AI_REQUEST_TIMEOUT
    })
  },

  /*** 流式发送消息 - Stream message ***/
  async streamMessage(data: ChatRequest, callbacks: StreamCallbacks): Promise<void> {
    // TODO: 实现 SSE 流式请求
    // UniApp 中需要使用 uni.request 配合 enableChunked 或使用 WebSocket
    try {
      const res = await this.sendMessage(data)
      if (res.code === 200 && res.data) {
        callbacks.onChunk(res.data.answer)
        callbacks.onDone({
          sessionId: res.data.sessionId,
          messageId: res.data.messageId
        })
      } else {
        callbacks.onError(res.msg || '请求失败')
      }
    } catch (error: any) {
      callbacks.onError(error.message || 'AI 服务暂时不可用')
    }
  },

  /*** 获取历史消息 - Get history messages ***/
  getHistory(params: PageParams & { sessionId?: string }): Promise<ApiResponse<PageResponse<ChatMessage>>> {
    return request({
      url: '/api/designer/chat/history',
      method: 'GET',
      data: params
    })
  },

  /*** 获取会话列表 - Get sessions ***/
  getSessions(projectId?: number): Promise<ApiResponse<ChatSession[]>> {
    return request({
      url: '/api/designer/chat/sessions',
      method: 'GET',
      data: projectId ? { projectId } : {}
    })
  }
}
