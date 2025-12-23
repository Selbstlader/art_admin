/**
 * AI 对话状态管理（仅 App 端）
 * AI Chat state management (App only)
 */

import { defineStore } from 'pinia'
import type { ChatState } from '@/types/store'
import type { ChatMessage, ChatSession } from '@/types/api'

/*** Chat store - manages AI chat messages and sessions (App only) ***/
export const useChatStore = defineStore('chat', {
  state: (): ChatState => ({
    messages: [],
    currentSessionId: '',
    isStreaming: false,
    sessions: [],
    projectContext: null
  }),

  getters: {
    /*** 获取当前会话消息 - Get current session messages ***/
    currentMessages: (state) => state.messages,
    
    /*** 是否正在流式输出 - Is streaming ***/
    streaming: (state) => state.isStreaming
  },

  actions: {
    /*** 发送消息 - Send message ***/
    async sendMessage(query: string): Promise<void> {
      // 添加用户消息 - Add user message
      const userMessage: ChatMessage = {
        id: Date.now(),
        role: 'user',
        content: query,
        createdAt: new Date().toISOString()
      }
      this.messages.push(userMessage)

      // TODO: 调用后端接口发送消息
      // await chatApi.sendMessage({ query, projectId: this.projectContext, sessionId: this.currentSessionId })
    },

    /*** 处理流式消息 - Handle streaming message ***/
    appendStreamContent(content: string): void {
      const lastMessage = this.messages[this.messages.length - 1]
      if (lastMessage && lastMessage.role === 'assistant') {
        lastMessage.content += content
      } else {
        // 创建新的 AI 消息
        this.messages.push({
          id: Date.now(),
          role: 'assistant',
          content: content,
          createdAt: new Date().toISOString()
        })
      }
    },

    /*** 设置流式状态 - Set streaming state ***/
    setStreaming(streaming: boolean): void {
      this.isStreaming = streaming
    },

    /*** 加载历史消息 - Load history messages ***/
    async loadHistory(): Promise<void> {
      // TODO: 调用后端接口
      // const res = await chatApi.getHistory({ sessionId: this.currentSessionId })
    },

    /*** 设置项目上下文 - Set project context ***/
    setProjectContext(projectId: number | null): void {
      this.projectContext = projectId
    },

    /*** 清除会话 - Clear session ***/
    clearSession(): void {
      this.messages = []
      this.currentSessionId = ''
      this.isStreaming = false
    }
  }
})
