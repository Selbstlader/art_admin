/**
 * AI 对话状态管理（仅 App 端）
 * AI Chat state management (App only)
 */

import { defineStore } from 'pinia'
import type { ChatState } from '@/types/store'
import type { ChatMessage, ChatSession, StreamCallbacks } from '@/types/api'
import { chatApi } from '@/api/chat'
import { AI_REQUEST_TIMEOUT } from '@/types/common'

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
    streaming: (state) => state.isStreaming,

    /*** 获取项目上下文 - Get project context ***/
    currentProjectContext: (state) => state.projectContext,

    /*** 获取消息数量 - Get message count ***/
    messageCount: (state) => state.messages.length,

    /*** 获取最后一条消息 - Get last message ***/
    lastMessage: (state) => state.messages.length > 0 ? state.messages[state.messages.length - 1] : null
  },

  actions: {
    /*** 
     * 发送消息 - Send message
     * 立即显示用户消息，然后调用 API 获取 AI 回复
     * Immediately show user message, then call API for AI response
     * @param query 用户输入的消息内容
     ***/
    async sendMessage(query: string): Promise<void> {
      if (!query.trim()) return

      // 添加用户消息到列表 - Add user message to list immediately
      const userMessage: ChatMessage = {
        id: Date.now(),
        role: 'user',
        content: query.trim(),
        createdAt: new Date().toISOString()
      }
      this.messages.push(userMessage)

      // 设置流式状态 - Set streaming state
      this.setStreaming(true)

      // 创建 AI 消息占位 - Create AI message placeholder
      const aiMessageId = Date.now() + 1
      const aiMessage: ChatMessage = {
        id: aiMessageId,
        role: 'assistant',
        content: '',
        createdAt: new Date().toISOString()
      }
      this.messages.push(aiMessage)

      try {
        // 调用流式消息 API - Call streaming message API
        const callbacks: StreamCallbacks = {
          onChunk: (content: string) => {
            this.appendStreamContent(content)
          },
          onDone: (result) => {
            this.currentSessionId = result.sessionId
            this.setStreaming(false)
          },
          onError: (error: string) => {
            // 更新 AI 消息为错误提示 - Update AI message with error
            const lastMsg = this.messages[this.messages.length - 1]
            if (lastMsg && lastMsg.role === 'assistant') {
              lastMsg.content = error || 'AI 服务暂时不可用，请稍后重试'
            }
            this.setStreaming(false)
          }
        }

        await chatApi.streamMessage(
          {
            query: query.trim(),
            projectId: this.projectContext || undefined,
            sessionId: this.currentSessionId || undefined
          },
          callbacks
        )
      } catch (error: any) {
        // 处理超时或其他错误 - Handle timeout or other errors
        const lastMsg = this.messages[this.messages.length - 1]
        if (lastMsg && lastMsg.role === 'assistant') {
          lastMsg.content = error?.message || 'AI 服务暂时不可用，请稍后重试'
        }
        this.setStreaming(false)
      }
    },

    /*** 
     * 处理流式消息内容累积 - Handle streaming message content accumulation
     * 每次 chunk 到达后，内容应为之前所有 chunk 的累积拼接
     * @param content 新到达的内容片段
     ***/
    appendStreamContent(content: string): void {
      const lastMessage = this.messages[this.messages.length - 1]
      if (lastMessage && lastMessage.role === 'assistant') {
        // 累积拼接内容 - Accumulate content
        lastMessage.content += content
      }
    },

    /*** 
     * 设置流式状态 - Set streaming state
     * @param streaming 是否正在流式输出
     ***/
    setStreaming(streaming: boolean): void {
      this.isStreaming = streaming
    },

    /*** 
     * 加载历史消息 - Load history messages
     * 支持上拉加载更多历史消息
     * @param page 页码，默认为 1
     * @param size 每页数量，默认为 20
     ***/
    async loadHistory(page: number = 1, size: number = 20): Promise<boolean> {
      if (!this.currentSessionId) return false

      try {
        const res = await chatApi.getHistory({
          sessionId: this.currentSessionId,
          current: page,
          size
        })

        if (res.code === 200 && res.data) {
          const historyMessages = res.data.records || []
          
          if (page === 1) {
            // 首次加载，替换消息列表 - First load, replace message list
            this.messages = historyMessages
          } else {
            // 加载更多，将历史消息插入到列表前面 - Load more, prepend to list
            this.messages = [...historyMessages, ...this.messages]
          }

          // 返回是否还有更多数据 - Return whether there's more data
          return historyMessages.length >= size
        }
        return false
      } catch (error) {
        console.error('加载历史消息失败:', error)
        return false
      }
    },

    /*** 
     * 加载会话列表 - Load session list
     * @param projectId 可选的项目 ID，用于筛选特定项目的会话
     ***/
    async loadSessions(projectId?: number): Promise<void> {
      try {
        const res = await chatApi.getSessions(projectId)
        if (res.code === 200 && res.data) {
          this.sessions = res.data
        }
      } catch (error) {
        console.error('加载会话列表失败:', error)
      }
    },

    /*** 
     * 设置项目上下文 - Set project context
     * 从项目详情页发起对话时，自动携带项目 ID
     * @param projectId 项目 ID，null 表示清除上下文
     ***/
    setProjectContext(projectId: number | null): void {
      this.projectContext = projectId
    },

    /*** 
     * 切换会话 - Switch session
     * @param sessionId 会话 ID
     ***/
    async switchSession(sessionId: string): Promise<void> {
      this.currentSessionId = sessionId
      this.messages = []
      await this.loadHistory(1)
    },

    /*** 
     * 创建新会话 - Create new session
     ***/
    createNewSession(): void {
      this.currentSessionId = ''
      this.messages = []
    },

    /*** 
     * 清除会话 - Clear session
     * 清空当前会话的所有消息和状态
     ***/
    clearSession(): void {
      this.messages = []
      this.currentSessionId = ''
      this.isStreaming = false
      this.projectContext = null
    },

    /*** 
     * 重试发送消息 - Retry sending message
     * 当 AI 响应失败时，可以重试最后一条用户消息
     ***/
    async retryLastMessage(): Promise<void> {
      // 找到最后一条用户消息 - Find last user message
      let lastUserMessageIndex = -1
      for (let i = this.messages.length - 1; i >= 0; i--) {
        if (this.messages[i].role === 'user') {
          lastUserMessageIndex = i
          break
        }
      }

      if (lastUserMessageIndex === -1) return

      const lastUserMessage = this.messages[lastUserMessageIndex]
      
      // 移除最后一条用户消息之后的所有消息（包括失败的 AI 回复）
      // Remove all messages after last user message (including failed AI response)
      this.messages = this.messages.slice(0, lastUserMessageIndex)

      // 重新发送消息 - Resend message
      await this.sendMessage(lastUserMessage.content)
    }
  }
})
