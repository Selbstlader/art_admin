<template>
  <ElDialog
    v-model="visible"
    :title="`测试对话 - ${tagData?.name || ''}`"
    width="800px"
    :close-on-click-modal="false"
    destroy-on-close
    @close="handleClose"
  >
    <div class="test-chat-container">
      <!-- 标签信息 -->
      <div class="tag-info">
        <ElDescriptions :column="2" border size="small">
          <ElDescriptionsItem label="标签名称">{{ tagData?.name }}</ElDescriptionsItem>
          <ElDescriptionsItem label="知识库">
            {{ tagData?.knowledge_base_name || '未配置' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="描述" :span="2">
            {{ tagData?.description || '无' }}
          </ElDescriptionsItem>
        </ElDescriptions>
      </div>

      <!-- 消息列表 -->
      <div class="message-list" ref="messageListRef">
        <div v-if="messages.length === 0" class="empty-tip">
          <ElIcon :size="32"><ChatDotRound /></ElIcon>
          <p>输入问题开始测试对话</p>
        </div>

        <div v-for="(msg, index) in messages" :key="index" class="message-group">
          <!-- 用户消息 -->
          <div class="message user-message">
            <div class="message-avatar">
              <ElIcon><User /></ElIcon>
            </div>
            <div class="message-content">
              <div class="message-text">{{ msg.query }}</div>
            </div>
          </div>

          <!-- AI 回复 -->
          <div class="message ai-message">
            <div class="message-avatar">
              <ElIcon><ChatDotRound /></ElIcon>
            </div>
            <div class="message-content">
              <div class="message-text" v-html="formatMarkdown(msg.answer)"></div>
            </div>
          </div>
        </div>

        <!-- 流式输出中 -->
        <div v-if="streaming && currentAnswer" class="message-group">
          <div class="message ai-message">
            <div class="message-avatar">
              <ElIcon><ChatDotRound /></ElIcon>
            </div>
            <div class="message-content">
              <div class="message-text" v-html="formatMarkdown(currentAnswer)"></div>
              <ElButton type="danger" text size="small" @click="handleStop">
                <ElIcon><CircleClose /></ElIcon>
                停止生成
              </ElButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-area">
        <ElInput
          v-model="query"
          type="textarea"
          :rows="2"
          placeholder="输入测试问题..."
          @keydown.ctrl.enter="handleSend"
        />
        <div class="input-actions">
          <ElButton size="small" @click="handleClear" :disabled="messages.length === 0">
            清空记录
          </ElButton>
          <ElButton
            type="primary"
            :loading="streaming"
            :disabled="!query.trim()"
            @click="handleSend"
          >
            <ElIcon><Promotion /></ElIcon>
            发送 (Ctrl+Enter)
          </ElButton>
        </div>
      </div>
    </div>
  </ElDialog>
</template>

<script setup lang="ts">
  /*** AI Tag Test Dialog Component ***/
  /*** Provides a chat interface for testing AI tag configuration ***/

  import { ref, computed, nextTick, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { User, ChatDotRound, CircleClose, Promotion } from '@element-plus/icons-vue'
  import { difyChatApi } from '@/api/dify'
  import { useUserStore } from '@/store/modules/user'
  import type { AITag } from '@/api/ai-tag'
  import type { DifyChatStreamEvent } from '@/types/dify'
  import { marked } from 'marked'

  interface Message {
    query: string
    answer: string
  }

  interface Props {
    modelValue: boolean
    tagData?: AITag
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
  }

  const props = withDefaults(defineProps<Props>(), {
    tagData: undefined
  })
  const emit = defineEmits<Emits>()

  const userStore = useUserStore()
  const userId = computed(() => userStore.info?.userId?.toString() || 'test_user')

  /*** State ***/
  const query = ref('')
  const messages = ref<Message[]>([])
  const streaming = ref(false)
  const currentAnswer = ref('')
  const currentTaskId = ref('')
  const conversationId = ref<string | undefined>(undefined)
  const messageListRef = ref<HTMLElement>()

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  /*** Send message ***/
  const handleSend = async () => {
    if (!query.value.trim() || !props.tagData) return

    const queryText = query.value.trim()
    query.value = ''
    streaming.value = true
    currentAnswer.value = ''

    try {
      const response = await difyChatApi.chatStream({
        query: queryText,
        user: userId.value,
        conversation_id: conversationId.value,
        tag_id: props.tagData.id
      })

      if (!response.ok) {
        throw new Error('请求失败')
      }

      const reader = response.body?.getReader()
      const decoder = new TextDecoder()

      if (!reader) {
        throw new Error('无法读取响应')
      }

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
                  scrollToBottom()
                }
                if (data.task_id) currentTaskId.value = data.task_id
                if (data.conversation_id) conversationId.value = data.conversation_id
              } else if (data.event === 'message_end') {
                messages.value.push({
                  query: queryText,
                  answer: currentAnswer.value
                })
                currentAnswer.value = ''
                scrollToBottom()
              } else if (data.event === 'error') {
                throw new Error((data as any).message || '对话出错')
              }
            } catch (e: any) {
              if (e.message && !e.message.includes('JSON')) {
                throw e
              }
            }
          }
        }
      }
    } catch (error: any) {
      console.error('测试对话失败:', error)
      ElMessage.error(error.message || '测试对话失败')
      if (currentAnswer.value) {
        messages.value.push({
          query: queryText,
          answer: currentAnswer.value || '[对话中断]'
        })
      }
    } finally {
      streaming.value = false
      currentTaskId.value = ''
      currentAnswer.value = ''
    }
  }

  /*** Stop generation ***/
  const handleStop = async () => {
    if (!currentTaskId.value) return

    try {
      await difyChatApi.stopMessage(currentTaskId.value, userId.value)
    } catch (error) {
      console.error('停止失败:', error)
    } finally {
      streaming.value = false
      if (currentAnswer.value) {
        messages.value.push({
          query: messages.value.length > 0 ? '' : query.value,
          answer: currentAnswer.value
        })
      }
      currentTaskId.value = ''
      currentAnswer.value = ''
    }
  }

  /*** Clear messages ***/
  const handleClear = () => {
    messages.value = []
    conversationId.value = undefined
  }

  /*** Close dialog ***/
  const handleClose = () => {
    visible.value = false
    query.value = ''
    messages.value = []
    conversationId.value = undefined
    currentAnswer.value = ''
    streaming.value = false
  }

  /*** Format markdown ***/
  const formatMarkdown = (text: string) => {
    return marked(text)
  }

  /*** Scroll to bottom ***/
  const scrollToBottom = () => {
    nextTick(() => {
      if (messageListRef.value) {
        messageListRef.value.scrollTop = messageListRef.value.scrollHeight
      }
    })
  }

  /*** Reset on tag change ***/
  watch(
    () => props.tagData?.id,
    () => {
      if (visible.value) {
        handleClear()
      }
    }
  )
</script>

<style lang="scss" scoped>
  .test-chat-container {
    display: flex;
    flex-direction: column;
    height: 500px;
  }

  .tag-info {
    margin-bottom: 16px;
    flex-shrink: 0;
  }

  .message-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    background: var(--el-fill-color-lighter);
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .empty-tip {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--el-text-color-secondary);

    p {
      margin-top: 12px;
    }
  }

  .message-group {
    margin-bottom: 16px;
  }

  .message {
    display: flex;
    gap: 10px;
    margin-bottom: 12px;

    .message-avatar {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 16px;
      flex-shrink: 0;
    }

    .message-content {
      flex: 1;
      min-width: 0;
    }

    .message-text {
      padding: 10px 14px;
      border-radius: 8px;
      line-height: 1.6;
      font-size: 14px;
      word-wrap: break-word;

      :deep(pre) {
        background: #f5f5f5;
        padding: 10px;
        border-radius: 4px;
        overflow-x: auto;
        font-size: 13px;
      }

      :deep(code) {
        background: #f5f5f5;
        padding: 2px 5px;
        border-radius: 3px;
        font-family: 'Courier New', monospace;
        font-size: 13px;
      }

      :deep(p) {
        margin: 0 0 8px;

        &:last-child {
          margin-bottom: 0;
        }
      }
    }

    &.user-message {
      .message-avatar {
        background: var(--el-color-primary);
        color: #fff;
      }

      .message-text {
        background: var(--el-color-primary-light-9);
      }
    }

    &.ai-message {
      .message-avatar {
        background: var(--el-color-success);
        color: #fff;
      }

      .message-text {
        background: #fff;
      }
    }
  }

  .input-area {
    flex-shrink: 0;

    .input-actions {
      margin-top: 10px;
      display: flex;
      justify-content: space-between;
    }
  }
</style>
