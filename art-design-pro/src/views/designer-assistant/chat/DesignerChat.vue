<template>
  <div class="designer-chat-page art-full-height">
    <div class="chat-layout">
      <!-- 左侧项目上下文和会话列表 -->
      <div class="context-sidebar">
        <div class="sidebar-header">
          <h3>项目上下文</h3>
        </div>

        <div class="project-select">
          <ElSelect
            v-model="currentProjectId"
            placeholder="选择项目"
            style="width: 100%"
            clearable
            @change="handleProjectChange"
          >
            <ElOption
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </ElSelect>
        </div>

        <div v-if="currentProject" class="project-info">
          <div class="info-item">
            <span class="label">面积</span>
            <span class="value">{{ currentProject.area }} m²</span>
          </div>
          <div class="info-item">
            <span class="label">预算</span>
            <span class="value">¥{{ currentProject.budget?.toLocaleString() }}</span>
          </div>
          <div class="info-item">
            <span class="label">风格</span>
            <span class="value">{{ currentProject.style || '-' }}</span>
          </div>
        </div>

        <ElDivider />

        <!-- 会话列表 -->
        <div class="session-section">
          <div class="section-header">
            <h4>对话历史</h4>
            <ElButton type="primary" size="small" @click="handleNewChat">
              <ElIcon><Plus /></ElIcon>
              新对话
            </ElButton>
          </div>
          <div class="session-list">
            <div
              v-for="session in sessions"
              :key="session.id"
              class="session-item"
              :class="{ active: currentSessionId === session.id }"
              @click="handleSelectSession(session)"
            >
              <div class="session-info">
                <div class="session-title">{{ session.title }}</div>
                <div class="session-time">{{ formatTime(session.updatedAt) }}</div>
              </div>
              <ElDropdown trigger="click" @command="(cmd: string) => handleSessionCommand(cmd, session)">
                <ElIcon class="session-more"><MoreFilled /></ElIcon>
                <template #dropdown>
                  <ElDropdownMenu>
                    <ElDropdownItem command="export">
                      <ElIcon><Download /></ElIcon>
                      导出
                    </ElDropdownItem>
                    <ElDropdownItem command="delete">
                      <ElIcon><Delete /></ElIcon>
                      删除
                    </ElDropdownItem>
                  </ElDropdownMenu>
                </template>
              </ElDropdown>
            </div>
            <ElEmpty v-if="sessions.length === 0" description="暂无对话" :image-size="60" />
          </div>
        </div>

        <ElDivider />

        <div class="quick-actions">
          <h4>快捷操作</h4>
          <ElButton text @click="askQuestion('请分析当前项目的设计要点')">
            <ElIcon><Document /></ElIcon>
            分析设计要点
          </ElButton>
          <ElButton text @click="askQuestion('请推荐适合当前项目的材料')">
            <ElIcon><Goods /></ElIcon>
            材料推荐
          </ElButton>
          <ElButton text @click="askQuestion('请估算当前项目的成本')">
            <ElIcon><Money /></ElIcon>
            成本估算
          </ElButton>
          <ElButton text @click="askQuestion('请检查设计方案的合规性')">
            <ElIcon><Checked /></ElIcon>
            合规检查
          </ElButton>
        </div>
      </div>

      <!-- 右侧对话区域 -->
      <div class="chat-main">
        <ElCard shadow="never" class="chat-card">
          <!-- 消息列表 -->
          <div class="message-list" ref="messageListRef">
            <div v-if="messages.length === 0" class="welcome-guide">
              <div class="guide-icon">
                <ElIcon :size="48"><ChatDotRound /></ElIcon>
              </div>
              <h3>设计师AI助手</h3>
              <p class="guide-desc">
                我是您的专业设计助手，可以帮您分析项目需求、推荐材料、估算成本、检查合规性等。
              </p>
              <div class="guide-tips">
                <div class="tip-item">
                  <ElIcon><InfoFilled /></ElIcon>
                  <span>选择项目后，AI将基于项目信息提供针对性建议</span>
                </div>
                <div class="tip-item">
                  <ElIcon><Document /></ElIcon>
                  <span>涉及设计规范时，将引用相关国家标准</span>
                </div>
              </div>
            </div>

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
                  <div class="message-text markdown-body" v-html="formatMarkdown(msg.answer)"></div>
                </div>
              </div>
            </div>

            <!-- 流式输出 -->
            <div v-if="loading" class="message-group">
              <div class="message user-message">
                <div class="message-avatar">
                  <ElIcon><User /></ElIcon>
                </div>
                <div class="message-content">
                  <div class="message-text">{{ pendingQuery }}</div>
                </div>
              </div>
              <div class="message ai-message">
                <div class="message-avatar">
                  <ElIcon><ChatDotRound /></ElIcon>
                </div>
                <div class="message-content">
                  <div class="message-text">
                    <ElIcon class="loading-icon"><Loading /></ElIcon>
                    正在思考中...
                  </div>
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
              placeholder="请输入您的问题... (Ctrl+Enter 发送)"
              :disabled="loading"
              @keydown.ctrl.enter="handleSend"
            />
            <div class="input-actions">
              <span class="input-tip">支持设计规范咨询、材料推荐、成本估算等</span>
              <ElButton type="primary" :loading="loading" @click="handleSend">
                <ElIcon><Promotion /></ElIcon>
                发送
              </ElButton>
            </div>
          </div>
        </ElCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/***
 * Designer Chat Component
 * 设计师AI对话页面组件
 * Requirements: 5.1, 5.2, 5.3, 5.4
 ***/
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  User,
  ChatDotRound,
  Promotion,
  Document,
  Goods,
  Money,
  Checked,
  Plus,
  MoreFilled,
  Delete,
  Download,
  InfoFilled,
  Loading
} from '@element-plus/icons-vue'
import { marked } from 'marked'
import dayjs from 'dayjs'
import { designerChatApi, type ChatMessage, type ChatSession } from '@/api/designer-chat'
import { getDesignerProjects } from '@/api/designer-project'

const query = ref('')
const loading = ref(false)
const pendingQuery = ref('')
const messageListRef = ref<HTMLElement>()

// 项目相关 / Project related
const currentProjectId = ref<number | null>(null)
const projectList = ref<any[]>([])

// 会话相关 / Session related
const currentSessionId = ref<string>('')
const sessions = ref<ChatSession[]>([])

// 消息列表 / Message list
const messages = ref<{ id: number; query: string; answer: string }[]>([])

// 当前项目 / Current project
const currentProject = computed(() => {
  return projectList.value.find((p) => p.id === currentProjectId.value)
})

// 初始化 / Initialize
onMounted(async () => {
  await loadProjects()
  await loadSessions()
})

// 加载项目列表 / Load project list
const loadProjects = async () => {
  try {
    const res = await getDesignerProjects({ current: 1, size: 100 })
    projectList.value = res.records || []
  } catch (error) {
    console.error('加载项目列表失败:', error)
  }
}

// 加载会话列表 / Load session list
const loadSessions = async () => {
  try {
    const res = await designerChatApi.getSessions(currentProjectId.value || undefined)
    sessions.value = res || []
  } catch (error) {
    console.error('加载会话列表失败:', error)
  }
}

// 项目切换 / Project change
const handleProjectChange = () => {
  currentSessionId.value = ''
  messages.value = []
  loadSessions()
}

// 选择会话 / Select session
const handleSelectSession = async (session: ChatSession) => {
  currentSessionId.value = session.id
  await loadChatHistory(session.id)
}

// 加载对话历史 / Load chat history
const loadChatHistory = async (sessionId: string) => {
  try {
    const res = await designerChatApi.getHistory({ sessionId, page: 1, pageSize: 100 })
    const historyMessages = res.data || []
    
    // 将历史消息转换为显示格式 / Convert history to display format
    const groupedMessages: { id: number; query: string; answer: string }[] = []
    for (let i = 0; i < historyMessages.length; i += 2) {
      const userMsg = historyMessages[i]
      const aiMsg = historyMessages[i + 1]
      if (userMsg && userMsg.role === 'user') {
        groupedMessages.push({
          id: userMsg.id,
          query: userMsg.content,
          answer: aiMsg?.content || ''
        })
      }
    }
    messages.value = groupedMessages.reverse()
    scrollToBottom()
  } catch (error) {
    console.error('加载对话历史失败:', error)
  }
}

// 发送消息 / Send message
const handleSend = async () => {
  if (!query.value.trim()) {
    ElMessage.warning('请输入问题')
    return
  }

  const userQuery = query.value
  query.value = ''
  pendingQuery.value = userQuery
  loading.value = true

  try {
    const res = await designerChatApi.chat({
      query: userQuery,
      projectId: currentProjectId.value || undefined,
      sessionId: currentSessionId.value || undefined
    })

    // 更新会话ID / Update session ID
    if (!currentSessionId.value && res.sessionId) {
      currentSessionId.value = res.sessionId
      await loadSessions()
    }

    // 添加消息 / Add message
    messages.value.push({
      id: res.messageId,
      query: userQuery,
      answer: res.answer
    })

    scrollToBottom()
  } catch (error: any) {
    ElMessage.error(error.message || '发送失败')
    // 恢复输入 / Restore input
    query.value = userQuery
  } finally {
    loading.value = false
    pendingQuery.value = ''
  }
}

// 快捷提问 / Quick question
const askQuestion = (question: string) => {
  query.value = question
  handleSend()
}

// 新建对话 / New chat
const handleNewChat = () => {
  currentSessionId.value = ''
  messages.value = []
  query.value = ''
}

// 会话操作 / Session command
const handleSessionCommand = async (command: string, session: ChatSession) => {
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除这个对话吗？', '提示', {
        type: 'warning'
      })
      await designerChatApi.deleteSession(session.id)
      ElMessage.success('删除成功')
      if (currentSessionId.value === session.id) {
        handleNewChat()
      }
      await loadSessions()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  } else if (command === 'export') {
    try {
      const res = await designerChatApi.exportHistory(session.id, 'md')
      // 创建下载 / Create download
      const blob = new Blob([res.content], { type: 'text/markdown' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = res.filename || 'chat_history.md'
      a.click()
      URL.revokeObjectURL(url)
      ElMessage.success('导出成功')
    } catch (error) {
      ElMessage.error('导出失败')
    }
  }
}

// 格式化Markdown / Format Markdown
const formatMarkdown = (text: string) => {
  return marked(text || '')
}

// 格式化时间 / Format time
const formatTime = (time: string) => {
  return dayjs(time).format('MM-DD HH:mm')
}

// 滚动到底部 / Scroll to bottom
const scrollToBottom = () => {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}
</script>

<style scoped lang="scss">
.designer-chat-page {
  padding: 16px;
}

.chat-layout {
  display: flex;
  height: 100%;
  gap: 16px;
}

.context-sidebar {
  width: 300px;
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .sidebar-header {
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 500;
    }
  }

  .project-select {
    margin-bottom: 16px;
  }

  .project-info {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 8px 0;
      border-bottom: 1px solid var(--el-border-color-lighter);

      .label {
        color: var(--el-text-color-secondary);
        font-size: 13px;
      }

      .value {
        font-weight: 500;
        font-size: 13px;
      }
    }
  }

  .session-section {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      h4 {
        margin: 0;
        font-size: 14px;
        font-weight: 500;
      }
    }

    .session-list {
      flex: 1;
      overflow-y: auto;
    }

    .session-item {
      padding: 10px 12px;
      margin-bottom: 8px;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.3s;
      display: flex;
      justify-content: space-between;
      align-items: center;

      &:hover {
        background: var(--el-fill-color-light);

        .session-more {
          opacity: 1;
        }
      }

      &.active {
        background: var(--el-color-primary-light-9);
        border-left: 3px solid var(--el-color-primary);
      }

      .session-info {
        flex: 1;
        min-width: 0;

        .session-title {
          font-size: 13px;
          font-weight: 500;
          margin-bottom: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .session-time {
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
      }

      .session-more {
        opacity: 0;
        transition: opacity 0.3s;
        cursor: pointer;
      }
    }
  }

  .quick-actions {
    h4 {
      margin-bottom: 12px;
      font-size: 14px;
      font-weight: 500;
    }

    .el-button {
      width: 100%;
      justify-content: flex-start;
      margin-bottom: 8px;
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
      padding: 0;
    }
  }
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
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
    gap: 12px;
    align-items: flex-start;

    .tip-item {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 13px;
      color: var(--el-text-color-regular);
      padding: 10px 16px;
      background: var(--el-fill-color-light);
      border-radius: 6px;

      .el-icon {
        color: var(--el-color-primary);
      }
    }
  }
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

    :deep(ul), :deep(ol) {
      padding-left: 20px;
      margin: 8px 0;
    }

    :deep(p) {
      margin: 8px 0;
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

.loading-icon {
  animation: rotate 1s linear infinite;
  margin-right: 8px;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.input-area {
  padding: 16px;
  border-top: 1px solid var(--el-border-color-light);
  background: var(--el-fill-color-lighter);

  .input-actions {
    margin-top: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .input-tip {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
}
</style>
