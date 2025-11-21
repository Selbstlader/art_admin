<template>
  <div class="dify-chat-page art-full-height">
    <div class="chat-layout">
      <!-- 左侧会话列表 -->
      <div class="conversation-sidebar">
        <div class="sidebar-header">
          <h3>对话历史</h3>
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
            <ElDropdown trigger="click" @command="(cmd) => handleConvCommand(cmd, conv)">
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

          <ElEmpty v-if="conversations.length === 0" description="暂无对话" />
        </div>
      </div>

      <!-- 右侧对话区域 -->
      <div class="chat-main">
        <ElCard shadow="never" class="chat-card">
          <!-- 消息列表 -->
          <div class="message-list" ref="messageListRef">
            <!-- 新对话引导 -->
            <div v-if="messages.length === 0 && !currentConversation" class="welcome-guide">
              <div class="guide-icon">
                <ElIcon :size="48"><ChatDotRound /></ElIcon>
              </div>
              <h3>报告分析助手</h3>
              <p class="guide-desc">
                我是您的专业报告分析助手，可以帮您从长篇报告中提取关键洞察、识别潜在风险并提炼核心信息。
              </p>
              <div class="guide-tips">
                <div class="tip-item">
                  <ElIcon><Document /></ElIcon>
                  <span>首次对话时，我会询问您的报告信息</span>
                </div>
                <div class="tip-item">
                  <ElIcon><ChatLineRound /></ElIcon>
                  <span>后续对话可直接提问，无需重复填写</span>
                </div>
              </div>
            </div>

            <ElEmpty v-else-if="messages.length === 0" description="暂无消息记录" />

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
                  <ElButton type="danger" text @click="handleStop">
                    <ElIcon><CircleClose /></ElIcon>
                    停止生成
                  </ElButton>
                </div>
              </div>
            </div>
          </div>

          <!-- 建议问题 -->
          <div v-if="suggestedQuestions.length > 0" class="suggested-questions">
            <div class="suggested-title">您可能想问:</div>
            <div class="question-list">
              <ElTag
                v-for="(question, index) in suggestedQuestions"
                :key="index"
                class="question-tag"
                @click="handleQuestionClick(question)"
              >
                {{ question }}
              </ElTag>
            </div>
          </div>

          <!-- 输入区域 -->
          <div class="input-area">
            <ElInput
              v-model="query"
              type="textarea"
              :rows="3"
              placeholder="请输入您的问题..."
              @keydown.ctrl.enter="handleSend"
            />
            <div class="input-actions">
              <ElButton type="primary" :loading="streaming || loading" @click="handleSend">
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

    <!-- 报告信息对话框 -->
    <ElDialog
      v-model="reportDialogVisible"
      title="填写报告信息"
      width="600px"
      :close-on-click-modal="false"
    >
      <ElForm :model="formData" label-width="100px">
        <ElFormItem label="报告标题" required>
          <ElInput v-model="formData.report_title" placeholder="请输入报告标题" />
        </ElFormItem>
        <ElFormItem label="报告日期">
          <ElDatePicker
            v-model="formData.report_date"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </ElFormItem>
        <ElFormItem label="报告内容" required>
          <ElInput
            v-model="formData.report_content"
            type="textarea"
            :rows="6"
            placeholder="请输入报告内容"
          />
        </ElFormItem>
        <ElFormItem label="分析重点">
          <ElInput
            v-model="formData.analysis_focus"
            type="textarea"
            :rows="3"
            placeholder="请输入您希望重点分析的方向（可选）"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="reportDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleReportConfirm">确定并发送</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, nextTick } from 'vue'
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
    Document,
    ChatLineRound
  } from '@element-plus/icons-vue'
  import { useDifyChat } from '@/composables/useDifyChat'
  import type { DifyConversation } from '@/types/dify'
  import dayjs from 'dayjs'
  import { marked } from 'marked'

  const {
    conversations,
    currentConversation,
    messages,
    suggestedQuestions,
    loading,
    streaming,
    currentAnswer,
    getConversations,
    sendMessageStream,
    stopMessage,
    deleteConversation,
    renameConversation,
    selectConversation,
    newConversation
  } = useDifyChat()

  // 表单数据
  const formData = ref({
    report_title: '',
    report_content: '',
    report_date: '',
    analysis_focus: ''
  })

  const query = ref('')
  const messageListRef = ref<HTMLElement>()
  const renameVisible = ref(false)
  const reportDialogVisible = ref(false)
  const newName = ref('')
  const currentRenameId = ref('')
  const pendingQuery = ref('')

  // 初始化
  onMounted(() => {
    getConversations()
  })

  // 发送消息
  const handleSend = async () => {
    if (!query.value.trim()) {
      ElMessage.warning('请输入问题')
      return
    }

    // 如果是新对话或当前对话没有报告信息，弹出对话框
    if (!currentConversation.value || !currentConversation.value.inputs?.report_title) {
      pendingQuery.value = query.value
      reportDialogVisible.value = true
      return
    }

    // 使用当前对话的报告信息
    await sendMessageWithInputs(query.value, currentConversation.value.inputs)
  }

  // 确认报告信息并发送
  const handleReportConfirm = async () => {
    if (!formData.value.report_title || !formData.value.report_content) {
      ElMessage.warning('请填写报告标题和内容')
      return
    }

    const inputs = {
      report_title: formData.value.report_title,
      report_content: formData.value.report_content,
      report_date: formData.value.report_date || dayjs().format('YYYY-MM-DD'),
      analysis_focus: formData.value.analysis_focus
    }

    reportDialogVisible.value = false
    await sendMessageWithInputs(pendingQuery.value, inputs)
  }

  // 发送消息（带报告信息）
  const sendMessageWithInputs = async (queryText: string, inputs: any) => {
    try {
      await sendMessageStream(
        {
          query: queryText,
          conversation_id: currentConversation.value?.id,
          inputs
        },
        () => {
          scrollToBottom()
        }
      )

      query.value = ''
      pendingQuery.value = ''
      scrollToBottom()
    } catch (error) {
      ElMessage.error('发送失败')
    }
  }

  // 停止生成
  const handleStop = () => {
    stopMessage()
  }

  // 新建对话
  const handleNewChat = () => {
    newConversation()
    query.value = ''
    formData.value = {
      report_title: '',
      report_content: '',
      report_date: '',
      analysis_focus: ''
    }
  }

  // 选择会话
  const handleSelectConversation = (conv: DifyConversation) => {
    selectConversation(conv)
  }

  // 会话操作
  const handleConvCommand = (command: string, conv: DifyConversation) => {
    if (command === 'delete') {
      handleDelete(conv.id)
    } else if (command === 'rename') {
      handleRename(conv)
    }
  }

  // 删除会话
  const handleDelete = async (id: string) => {
    const success = await deleteConversation(id)
    if (success) {
      ElMessage.success('删除成功')
    }
  }

  // 重命名会话
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

  // 点击建议问题
  const handleQuestionClick = (question: string) => {
    query.value = question
    handleSend()
  }

  // 格式化时间
  const formatTime = (timestamp: number) => {
    return dayjs.unix(timestamp).format('MM-DD HH:mm')
  }

  // 格式化 Markdown
  const formatMarkdown = (text: string) => {
    return marked(text)
  }

  // 滚动到底部
  const scrollToBottom = () => {
    nextTick(() => {
      if (messageListRef.value) {
        messageListRef.value.scrollTop = messageListRef.value.scrollHeight
      }
    })
  }
</script>

<style scoped lang="scss">
  .dify-chat-page {
    padding: 16px;
  }

  .chat-layout {
    display: flex;
    height: 100%;
    gap: 16px;
  }

  .conversation-sidebar {
    width: 280px;
    background: #fff;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .sidebar-header {
      padding: 16px;
      border-bottom: 1px solid var(--el-border-color-light);
      display: flex;
      justify-content: space-between;
      align-items: center;

      h3 {
        margin: 0;
        font-size: 16px;
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
      }
    }
  }

  .message-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    margin-bottom: 16px;
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

  .suggested-questions {
    padding: 16px;
    border-top: 1px solid var(--el-border-color-light);
    margin-bottom: 16px;

    .suggested-title {
      font-size: 14px;
      color: var(--el-text-color-secondary);
      margin-bottom: 12px;
    }

    .question-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .question-tag {
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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
      align-items: flex-start;

      .tip-item {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 14px;
        color: var(--el-text-color-regular);
        padding: 12px 20px;
        background: var(--el-fill-color-light);
        border-radius: 8px;
        min-width: 300px;

        .el-icon {
          color: var(--el-color-primary);
          font-size: 18px;
        }
      }
    }
  }
</style>
