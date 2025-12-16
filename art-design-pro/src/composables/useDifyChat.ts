import { ref, computed } from 'vue'
import { difyChatApi, difyConversationApi } from '@/api/dify'
import type {
  DifyConversation,
  DifyMessage,
  DifyChatRequest,
  DifyChatStreamEvent
} from '@/types/dify'
import { useUserStore } from '@/store/modules/user'

export function useDifyChat() {
  const userStore = useUserStore()
  const userId = computed(() => userStore.info?.userId?.toString() || 'anonymous')

  // 状态
  const conversations = ref<DifyConversation[]>([])
  const currentConversation = ref<DifyConversation | null>(null)
  const messages = ref<DifyMessage[]>([])
  const suggestedQuestions = ref<string[]>([])
  const loading = ref(false)
  const streaming = ref(false)
  const currentAnswer = ref('')
  const currentTaskId = ref('')

  // 获取会话列表
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

  // 获取会话消息历史
  const getMessages = async (conversationId: string, limit = 20) => {
    loading.value = true
    try {
      const res = (await difyConversationApi.getMessages({
        conversation_id: conversationId,
        user: userId.value,
        limit
      })) as {
        data: DifyMessage[]
      }
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

  // 发送消息(非流式)
  const sendMessage = async (params: Omit<DifyChatRequest, 'user'>) => {
    loading.value = true
    try {
      const res = (await difyChatApi.chat({
        ...params,
        user: userId.value
      })) as {
        code: number
        data: {
          message_id: string
          conversation_id: string
          answer: string
          created_at: number
        }
      }
      if (res.code === 200) {
        // 添加到消息列表
        const newMessage: DifyMessage = {
          id: res.data.message_id,
          conversation_id: res.data.conversation_id,
          inputs: params.inputs || {},
          query: params.query,
          answer: res.data.answer,
          message_files: [],
          retriever_resources: [],
          created_at: res.data.created_at
        }
        messages.value.push(newMessage)

        // 更新当前会话ID
        if (
          !currentConversation.value ||
          currentConversation.value.id !== res.data.conversation_id
        ) {
          await getConversations()
          currentConversation.value =
            conversations.value.find((c) => c.id === res.data.conversation_id) || null
        }

        return res.data
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 发送消息(流式)
  const sendMessageStream = async (
    params: Omit<DifyChatRequest, 'user'>,
    onChunk?: (chunk: string) => void
  ) => {
    streaming.value = true
    currentAnswer.value = ''

    try {
      const response = await difyChatApi.chatStream({
        ...params,
        user: userId.value
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
                // 消息结束,添加到列表
                const newMessage: DifyMessage = {
                  id: messageId,
                  conversation_id: conversationId,
                  inputs: params.inputs || {},
                  query: params.query,
                  answer: currentAnswer.value,
                  message_files: [],
                  retriever_resources: [],
                  created_at: data.created_at || Date.now() / 1000
                }
                messages.value.push(newMessage)

                // 更新会话列表
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

  // 停止消息生成
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

  // 获取建议问题
  const getSuggestedQuestions = async (messageId: string) => {
    try {
      const res = (await difyChatApi.getSuggestedQuestions(messageId, userId.value)) as {
        code: number
        data: string[]
      }
      if (res.code === 200) {
        suggestedQuestions.value = res.data
        return res.data
      }
    } catch (error: any) {
      // 如果是建议问题功能被禁用，静默处理，不显示错误
      if (error?.message?.includes('Suggested Questions Is Disabled')) {
        console.log('建议问题功能未启用')
      } else {
        console.error('获取建议问题失败:', error)
      }
      return []
    }
  }

  // 删除会话
  const deleteConversation = async (conversationId: string) => {
    try {
      await difyConversationApi.deleteConversation(conversationId, userId.value)
      await getConversations()
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

  // 重命名会话
  const renameConversation = async (conversationId: string, name: string) => {
    try {
      const res = (await difyConversationApi.renameConversation(conversationId, {
        user: userId.value,
        name
      })) as {
        data: DifyConversation
      }
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

  // 选择会话
  const selectConversation = async (conversation: DifyConversation) => {
    currentConversation.value = conversation
    await getMessages(conversation.id)

    // 获取最后一条消息的建议问题
    if (messages.value.length > 0) {
      const lastMessage = messages.value[messages.value.length - 1]
      await getSuggestedQuestions(lastMessage.id)
    }
  }

  // 新建会话
  const newConversation = () => {
    currentConversation.value = null
    messages.value = []
    suggestedQuestions.value = []
    currentAnswer.value = ''
  }

  return {
    // 状态
    conversations,
    currentConversation,
    messages,
    suggestedQuestions,
    loading,
    streaming,
    currentAnswer,

    // 方法
    getConversations,
    getMessages,
    sendMessage,
    sendMessageStream,
    stopMessage,
    getSuggestedQuestions,
    deleteConversation,
    renameConversation,
    selectConversation,
    newConversation
  }
}
