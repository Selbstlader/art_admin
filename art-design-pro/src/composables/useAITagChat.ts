/*** AI Tag Chat Composable ***/
/*** Handles chat functionality with AI tags, system prompts are hidden from users ***/

import { ref, computed, type Ref } from 'vue'
import { difyChatApi, difyConversationApi } from '@/api/dify'
import type { DifyConversation, DifyMessage, DifyChatStreamEvent } from '@/types/dify'
import { useUserStore } from '@/store/modules/user'
import type { AITag } from '@/api/ai-tag'

export function useAITagChat(tag: Ref<AITag | null>) {
  const userStore = useUserStore()
  const userId = computed(() => userStore.info?.userId?.toString() || 'anonymous')

  /*** State ***/
  const conversations = ref<DifyConversation[]>([])
  const currentConversation = ref<DifyConversation | null>(null)
  const messages = ref<DifyMessage[]>([])
  const loading = ref(false)
  const streaming = ref(false)
  const currentAnswer = ref('')
  const currentTaskId = ref('')

  /*** Get conversations for this tag ***/
  const getConversations = async (limit = 20) => {
    loading.value = true
    try {
      const res = (await difyConversationApi.getConversations({
        user: userId.value,
        limit
      })) as { data: DifyConversation[] }

      if (res.data) {
        conversations.value = res.data
        return res.data
      }
    } catch (error) {
      console.error('获取会话列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  /*** Get messages for a conversation ***/
  const getMessages = async (conversationId: string, limit = 20) => {
    loading.value = true
    try {
      const res = (await difyConversationApi.getMessages({
        conversation_id: conversationId,
        user: userId.value,
        limit
      })) as { data: DifyMessage[] }

      if (res) {
        messages.value = res.data
        return res.data
      }
    } catch (error) {
      console.error('获取消息历史失败:', error)
    } finally {
      loading.value = false
    }
  }

  /*** Send message with tag configuration (streaming) ***/
  const sendMessage = async (query: string, onChunk?: (chunk: string) => void) => {
    if (!tag.value) {
      throw new Error('未选择AI标签')
    }

    streaming.value = true
    currentAnswer.value = ''

    try {
      const response = await difyChatApi.chatStream({
        query,
        user: userId.value,
        conversation_id: currentConversation.value?.id,
        tag_id: tag.value.id
      })

      if (!response.ok) {
        throw new Error('Stream request failed')
      }

      const reader = response.body?.getReader()
      const decoder = new TextDecoder()

      if (!reader) {
        throw new Error('No reader available')
      }

      let conversationId = ''
      let messageId = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            try {
              const data: DifyChatStreamEvent = JSON.parse(line.slice(6))

              if (data.event === 'message') {
                if (data.answer) {
                  currentAnswer.value += data.answer
                  onChunk?.(data.answer)
                }
                if (data.task_id) currentTaskId.value = data.task_id
                if (data.message_id) messageId = data.message_id
                if (data.conversation_id) conversationId = data.conversation_id
              } else if (data.event === 'message_end') {
                const newMessage: DifyMessage = {
                  id: messageId,
                  conversation_id: conversationId,
                  inputs: {},
                  query,
                  answer: currentAnswer.value,
                  message_files: [],
                  retriever_resources: [],
                  created_at: data.created_at || Date.now() / 1000
                }
                messages.value.push(newMessage)

                await getConversations()
                currentConversation.value =
                  conversations.value.find((c) => c.id === conversationId) || null
              }
            } catch (e) {
              console.error('Parse SSE error:', e)
            }
          }
        }
      }

      return { conversation_id: conversationId, message_id: messageId }
    } catch (error) {
      console.error('流式对话失败:', error)
      throw error
    } finally {
      streaming.value = false
      currentTaskId.value = ''
    }
  }

  /*** Stop message generation ***/
  const stopMessage = async () => {
    if (!currentTaskId.value) return

    try {
      await difyChatApi.stopMessage(currentTaskId.value, userId.value)
      streaming.value = false
      currentTaskId.value = ''
    } catch (error) {
      console.error('停止消息失败:', error)
    }
  }

  /*** Delete conversation ***/
  const deleteConversation = async (conversationId: string) => {
    try {
      await difyConversationApi.deleteConversation(conversationId, userId.value)
      conversations.value = conversations.value.filter((c) => c.id !== conversationId)
      if (currentConversation.value?.id === conversationId) {
        currentConversation.value = null
        messages.value = []
      }
      return true
    } catch (error) {
      console.error('删除会话失败:', error)
      return false
    }
  }

  /*** Rename conversation ***/
  const renameConversation = async (conversationId: string, name: string) => {
    try {
      const res = (await difyConversationApi.renameConversation(conversationId, {
        user: userId.value,
        name
      })) as { data: DifyConversation }

      if (res) {
        const index = conversations.value.findIndex((c) => c.id === conversationId)
        if (index !== -1) {
          conversations.value[index] = res.data
        }
        if (currentConversation.value?.id === conversationId) {
          currentConversation.value = res.data
        }
        return true
      }
    } catch (error) {
      console.error('重命名会话失败:', error)
      return false
    }
  }

  /*** Select conversation ***/
  const selectConversation = async (conversation: DifyConversation) => {
    currentConversation.value = conversation
    await getMessages(conversation.id)
  }

  /*** New conversation ***/
  const newConversation = () => {
    currentConversation.value = null
    messages.value = []
    currentAnswer.value = ''
  }

  /*** Reset state ***/
  const reset = () => {
    conversations.value = []
    currentConversation.value = null
    messages.value = []
    loading.value = false
    streaming.value = false
    currentAnswer.value = ''
    currentTaskId.value = ''
  }

  return {
    /*** State ***/
    conversations,
    currentConversation,
    messages,
    loading,
    streaming,
    currentAnswer,

    /*** Methods ***/
    getConversations,
    getMessages,
    sendMessage,
    stopMessage,
    deleteConversation,
    renameConversation,
    selectConversation,
    newConversation,
    reset
  }
}
