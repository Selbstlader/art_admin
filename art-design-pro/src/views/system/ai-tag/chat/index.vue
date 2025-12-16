<template>
  <div class="ai-tag-chat-page art-full-height">
    <div class="chat-layout">
      <!-- 左侧：标签信息 + 会话列表 -->
      <div class="chat-sidebar">
        <!-- 标签信息卡片 -->
        <div class="tag-info-card" v-if="currentTag">
          <div class="tag-header">
            <ElIcon :size="24" class="tag-icon"><Collection /></ElIcon>
            <div class="tag-title">{{ currentTag.name }}</div>
          </div>
          <div class="tag-desc" v-if="currentTag.description">
            {{ currentTag.description }}
          </div>
          <div class="tag-meta">
            <ElTag v-if="currentTag.knowledge_base_name" size="small" type="info">
              {{ currentTag.knowledge_base_name }}
            </ElTag>
          </div>
        </div>

        <!-- 会话列表 -->
        <div class="conversation-section">
          <div class="section-header">
            <span>对话历史</span>
            <ElButton type="primary" size="small" @click="handleNewChat">
              <ElIcon><Plus /></ElIcon>
              新对话
            </ElButton>
          </div>

          <div class="conversation-list">
            <div
              v-for="conv in conversations"
              :key="conv.id"
              class="conversation-item"
              :class="{ active: currentConversation?.id === conv.id }"
              @click="handleSelectConversation(conv)"
            >
              <div class="conv-info">
                <div class="conv-name">{{ conv.name }}</div>
                <div class="conv-time">{{ formatTime(conv.updated_at) }}</div>
              </div>
              <ElDropdown trigger="click" @command="(cmd: string) => handleConvCommand(cmd, conv)">
                <ElIcon class="conv-more"><MoreFilled /></ElIcon>
                <template #dropdown>
                  <ElDropdownMenu>
                    <ElDropdownItem command="rename">
                      <ElIcon><Edit /></ElIcon>
                      重命名
                    </ElDropdownItem>
                    <ElDropdownItem command="delete">
                      <ElIcon><Delete /></ElIcon>
                      删除
                    </ElDropdownItem>
                  </ElDropdownMenu>
                </template>
              </ElDropdown>
            </div>

            <ElEmpty v-if="conversations.length === 0" description="暂无对话" :image-size="60" />
          </div>
        </div>
      </div>

      <!-- 右侧：对话区域 -->
      <div class="chat-main">
        <ElCard shadow="never" class="chat-card">
          <!-- 消息列表 -->
          <div class="message-list" ref="messageListRef">
            <!-- 欢迎引导 -->
            <div v-if="messages.length === 0 && !currentConversation" class="welcome-guide">
              <div class="guide-icon">
                <ElIcon :size="48"><ChatDotRound /></ElIcon>
              </div>
              <h3>{{ currentTag?.name || 'AI 助手' }}</h3>
              <p class="guide-desc">
                {{ currentTag?.description || '开始与 AI 助手对话，获取专业帮助' }}
              </p>
              <div class="guide-tips">
                <div class="tip-item">
                  <ElIcon><Promotion /></ElIcon>
                  <span>直接输入您的问题开始对话</span>
                </div>
              </div>
            </div>

            <ElEmpty v-else-if="messages.length === 0" description="暂无消息记录" />

            <!-- 消息列表 -->
            <div v-for="msg in messages" :key="msg.id" class="message-group">
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

            <!-- 流式输出中的消息 -->
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
              :rows="3"
              placeholder="请输入您的问题..."
              :disabled="!currentTag"
              @keydown.ctrl.enter="handleSend"
            />
            <div class="input-actions">
              <ElButton
                type="primary"
                :loading="streaming || loading"
                :disabled="!currentTag || !query.trim()"
                @click="handleSend"
              >
                <ElIcon><Promotion /></ElIcon>
                发送 (Ctrl+Enter)
              </ElButton>
            </div>
          </div>
        </ElCard>
      </div>
    </div>

    <!-- 重命名对话框 -->
    <ElDialog v-model="renameVisible" title="重命名会话" width="400px">
      <ElInput v-model="newName" placeholder="请输入新名称" />
      <template #footer>
        <ElButton @click="renameVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleRenameConfirm">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">

  /*** AI Tag Chat Page ***/
  /*** Users only see the chat interface, system prompt is applied automatically ***/

  import { ref, onMounted, nextTick, computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import {
    Plus,
    MoreFilled,
    Edit,
    Delete,
    User,
    ChatDotRound,
    CircleClose,
    Promotion,
    Collection
  } from '@element-plus/icons-vue'
  import { aiTagApi, type AITag } from '@/api/ai-tag'
  import { useAITagChat } from '@/composables/useAITagChat'
  import type { DifyConversation } from '@/types/dify'
  import dayjs from 'dayjs'
  import { marked } from 'marked'

  defineOptions({ name: 'AITagChat' })

  const route = useRoute()

  /*** State ***/
  const currentTag = ref<AITag | null>(null)
  const query = ref('')
  const messageListRef = ref<HTMLElement>()
  const renameVisible = ref(false)
  const newName = ref('')
  const currentRenameId = ref('')

  /*** Use AI Tag Chat composable ***/
  const {
    conversations,
    currentConversation,
    messages,
    loading,
    streaming,
    currentAnswer,
    getConversations,
    sendMessage,
    stopMessage,
    deleteConversation,
    renameConversation,
    selectConversation,
    newConversation
  } = useAITagChat(currentTag)

  /*** Load tag data ***/
  const loadTag = async () => {
    const tagId = route.params.id as string
    if (!tagId) {
      ElMessage.error('缺少标签ID')
      return
    }

    try {
      const tag = await aiTagApi.get(Number(tagId))
      currentTag.value = tag as AITag
    } catch (error: any) {
      console.error('加载标签失败:', error)
      ElMessage.error(error.message || '加载标签失败')
    }
  }

  /*** Send message ***/
  const handleSend = async () => {
    if (!query.value.trim()) {
      ElMessage.warning('请输入问题')
      return
    }

    if (!currentTag.value) {
      ElMessage.warning('请先选择AI标签')
      return
    }

    const queryText = query.value
    query.value = ''

    try {
      await sendMessage(queryText, () => {
        scrollToBottom()
      })
      scrollToBottom()
    } catch (error: any) {
      ElMessage.error(error.message || '发送失败')
      query.value = queryText
    }
  }

  /*** Stop generation ***/
  const handleStop = () => {
    stopMessage()
  }

  /*** New chat ***/
  const handleNewChat = () => {
    newConversation()
    query.value = ''
  }

  /*** Select conversation ***/
  const handleSelectConversation = (conv: DifyConversation) => {
    selectConversation(conv)
  }

  /*** Conversation commands ***/
  const handleConvCommand = (command: string, conv: DifyConversation) => {
    if (command === 'delete') {
      handleDelete(conv.id)
    } else if (command === 'rename') {
      handleRename(conv)
    }
  }

  /*** Delete conversation ***/
  const handleDelete = async (id: string) => {
    const success = await deleteConversation(id)
    if (success) {
      ElMessage.success('删除成功')
    }
  }

  /*** Rename conversation ***/
  const handleRename = (conv: DifyConversation) => {
    currentRenameId.value = conv.id
    newName.value = conv.name
    renameVisible.value = true
  }

  const handleRenameConfirm = async () => {
    if (!newName.value.trim()) {
      ElMessage.warning('请输入名称')
      return
    }

    const success = await renameConversation(currentRenameId.value, newName.value)
    if (success) {
      ElMessage.success('重命名成功')
      renameVisible.value = false
    }
  }

  /*** Format time ***/
  const formatTime = (timestamp: number) => {
    return dayjs.unix(timestamp).format('MM-DD HH:mm')
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

  /*** Lifecycle ***/
  onMounted(async () => {
    await loadTag()
    if (currentTag.value) {
      await getConversations()
    }
  })
</script>

<style lang="scss" scoped>
  .ai-tag-chat-page {
    padding: 16px;
  }

  .chat-layout {
    display: flex;
    height: 100%;
    gap: 16px;
  }

  .chat-sidebar {
    width: 300px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .tag-info-card {
    background: #fff;
    border-radius: 8px;
    padding: 16px;

    .tag-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;

      .tag-icon {
        color: var(--el-color-primary);
      }

      .tag-title {
        font-size: 18px;
        font-weight: 600;
      }
    }

    .tag-desc {
      font-size: 14px;
      color: var(--el-text-color-secondary);
      margin-bottom: 12px;
      line-height: 1.5;
    }

    .tag-meta {
      display: flex;
      gap: 8px;
    }
  }

  .conversation-section {
    flex: 1;
    background: #fff;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .section-header {
      padding: 16px;
      border-bottom: 1px solid var(--el-border-color-light);
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 500;
    }
  }

  .conversation-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  .conversation-item {
    padding: 12px;
    margin-bottom: 8px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s;
    display: flex;
    justify-content: space-between;
    align-items: center;

    &:hover {
      background: var(--el-fill-color-light);

      .conv-more {
        opacity: 1;
      }
    }

    &.active {
      background: var(--el-color-primary-light-9);
      border-left: 3px solid var(--el-color-primary);
    }

    .conv-info {
      flex: 1;
      min-width: 0;

      .conv-name {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 4px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .conv-time {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }

    .conv-more {
      opacity: 0;
      transition: opacity 0.3s;
      cursor: pointer;
    }
  }

  .chat-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .chat-card {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;

      :deep(.el-card__body) {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        padding: 0;
      }
    }
  }

  .message-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  .message-group {
    margin-bottom: 24px;
  }

  .message {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;

    .message-avatar {
      width: 36px;
      height: 36px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 18px;
      flex-shrink: 0;
    }

    .message-content {
      flex: 1;
      min-width: 0;
    }

    .message-text {
      padding: 12px 16px;
      border-radius: 8px;
      line-height: 1.6;
      word-wrap: break-word;

      :deep(pre) {
        background: #f5f5f5;
        padding: 12px;
        border-radius: 4px;
        overflow-x: auto;
      }

      :deep(code) {
        background: #f5f5f5;
        padding: 2px 6px;
        border-radius: 3px;
        font-family: 'Courier New', monospace;
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
        background: var(--el-fill-color-light);
      }
    }
  }

  .input-area {
    padding: 16px;
    border-top: 1px solid var(--el-border-color-light);
    background: var(--el-fill-color-lighter);

    .input-actions {
      margin-top: 12px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .welcome-guide {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 40px;
    text-align: center;

    .guide-icon {
      color: var(--el-color-primary);
      margin-bottom: 24px;
    }

    h3 {
      font-size: 24px;
      font-weight: 600;
      margin: 0 0 16px;
      color: var(--el-text-color-primary);
    }

    .guide-desc {
      font-size: 14px;
      color: var(--el-text-color-secondary);
      line-height: 1.6;
      max-width: 500px;
      margin: 0 0 32px;
    }

    .guide-tips {
      display: flex;
      flex-direction: column;
      gap: 16px;

      .tip-item {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 14px;
        color: var(--el-text-color-regular);
        padding: 12px 20px;
        background: var(--el-fill-color-light);
        border-radius: 8px;

        .el-icon {
          color: var(--el-color-primary);
          font-size: 18px;
        }
      }
    }
  }
</style>
