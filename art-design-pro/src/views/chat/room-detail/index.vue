<template>
  <div class="chat-room-detail art-full-height">
    <!-- 聊天室头部 -->
    <div class="chat-header">
      <div class="chat-header-left">
        <ElButton link @click="goBack">
          <ElIcon><ArrowLeft /></ElIcon>
          返回
        </ElButton>
        <ElDivider direction="vertical" />
        <div class="room-info">
          <h3>{{ roomInfo?.name }}</h3>
          <span class="room-desc">{{ roomInfo?.description }}</span>
        </div>
      </div>

      <div class="chat-header-right">
        <ElSpace>
          <ElTag v-if="roomInfo?.type === 'public'" type="success">公开</ElTag>
          <ElTag v-else type="warning">私密</ElTag>

          <ElPopover placement="bottom" width="300" trigger="hover">
            <template #reference>
              <ElButton link>
                <ElIcon><User /></ElIcon>
                在线 {{ onlineUsers.length }}
              </ElButton>
            </template>

            <div class="online-users">
              <div class="online-users-title">在线成员</div>
              <div class="online-users-list">
                <div v-for="user in onlineUsers" :key="user.userId" class="online-user-item">
                  <ElAvatar :size="24" :src="user.avatar">
                    {{ user.username?.charAt(0) }}
                  </ElAvatar>
                  <span>{{ user.username }}</span>
                </div>
              </div>
            </div>
          </ElPopover>

          <ElButton @click="showMembers = true">
            <ElIcon><Setting /></ElIcon>
            管理
          </ElButton>
        </ElSpace>
      </div>
    </div>

    <!-- 聊天内容区域 -->
    <div class="chat-content">
      <!-- 消息列表 -->
      <div class="chat-messages" ref="messagesRef">
        <div v-if="loading" class="loading-more">
          <ElIcon class="is-loading"><Loading /></ElIcon>
          加载中...
        </div>

        <div
          v-for="message in messages"
          :key="message.id"
          class="message-item"
          :class="{ 'is-own': message.userId === currentUserId }"
        >
          <!-- 回复消息 -->
          <div v-if="message.replyTo" class="reply-message">
            <ElIcon><Reply /></ElIcon>
            回复 {{ message.replyTo.username }}: {{ message.replyTo.content }}
          </div>

          <div class="message-content">
            <ElAvatar :size="32" :src="message.avatar" class="message-avatar">
              {{ message.username?.charAt(0) }}
            </ElAvatar>

            <div class="message-body">
              <div class="message-header">
                <span class="message-username">{{ message.username }}</span>
                <span class="message-time">{{ formatTime(message.createdAt) }}</span>
              </div>

              <div class="message-text" :class="{ 'is-recalled': message.isRecalled }">
                <span v-if="message.isRecalled" class="recalled-text">
                  <ElIcon><Delete /></ElIcon>
                  消息已撤回
                </span>
                <span v-else>{{ message.content }}</span>
              </div>

              <!-- 消息操作 -->
              <div v-if="!message.isRecalled" class="message-actions">
                <ElButton link size="small" @click="handleReply(message)">
                  <ElIcon><Reply /></ElIcon>
                </ElButton>
                <ElButton
                  v-if="message.userId === currentUserId && canRecall(message)"
                  link
                  size="small"
                  type="danger"
                  @click="handleRecall(message)"
                >
                  <ElIcon><Delete /></ElIcon>
                </ElButton>
              </div>
            </div>
          </div>
        </div>

        <!-- 正在输入提示 -->
        <div v-if="typingUsers.length > 0" class="typing-indicator">
          <ElIcon class="typing-icon"><Edit /></ElIcon>
          {{ getTypingText() }}
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="chat-input">
        <!-- 回复提示 -->
        <div v-if="replyMessage" class="reply-preview">
          <div class="reply-content">
            <ElIcon><Reply /></ElIcon>
            回复 {{ replyMessage.username }}: {{ replyMessage.content }}
          </div>
          <ElButton link @click="cancelReply">
            <ElIcon><Close /></ElIcon>
          </ElButton>
        </div>

        <div class="input-area">
          <ElInput
            v-model="inputMessage"
            type="textarea"
            :rows="3"
            placeholder="输入消息..."
            :maxlength="1000"
            show-word-limit
            @keydown="handleKeyDown"
            @input="handleInput"
          />

          <div class="input-actions">
            <ElButton type="primary" @click="sendMessage" :disabled="!inputMessage.trim()">
              发送
            </ElButton>
          </div>
        </div>
      </div>
    </div>

    <!-- 成员管理弹窗 -->
    <ElDialog v-model="showMembers" title="聊天室管理" width="600px">
      <ElTabs v-model="activeTab">
        <ElTabPane label="成员列表" name="members">
          <div class="members-list">
            <div v-for="member in members" :key="member.id" class="member-item">
              <div class="member-info">
                <ElAvatar :size="32" :src="member.avatar">
                  {{ member.username?.charAt(0) }}
                </ElAvatar>
                <div class="member-details">
                  <div class="member-name">{{ member.username }}</div>
                  <div class="member-role">
                    <ElTag :type="getRoleType(member.role)" size="small">
                      {{ getRoleText(member.role) }}
                    </ElTag>
                    <ElTag v-if="member.isOnline" type="success" size="small">在线</ElTag>
                  </div>
                </div>
              </div>

              <div class="member-actions">
                <ElButton v-if="member.role !== 'owner'" link type="danger" size="small">
                  移除
                </ElButton>
              </div>
            </div>
          </div>
        </ElTabPane>

        <ElTabPane label="聊天室设置" name="settings">
          <ElForm :model="roomSettings" label-width="100px">
            <ElFormItem label="聊天室名称">
              <ElInput v-model="roomSettings.name" />
            </ElFormItem>
            <ElFormItem label="聊天室描述">
              <ElInput v-model="roomSettings.description" type="textarea" :rows="3" />
            </ElFormItem>
            <ElFormItem label="最大成员数">
              <ElInputNumber v-model="roomSettings.maxMembers" :min="0" />
            </ElFormItem>
          </ElForm>
        </ElTabPane>
      </ElTabs>

      <template #footer>
        <ElButton @click="showMembers = false">关闭</ElButton>
        <ElButton v-if="activeTab === 'settings'" type="primary" @click="saveSettings">
          保存设置
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ArrowLeft, User, Setting, Loading, Delete, Edit, Close } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import { useChatWebSocket } from '@/composables/useWebSocket'
  import { chatRoomApi, chatMessageApi } from '@/api/chat'
  import { useUserStore } from '@/store/modules/user'

  defineOptions({ name: 'ChatRoomDetail' })

  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const roomId = computed(() => Number(route.params.id))
  console.log(userStore.info, 'userStore.info?.userIdƒ')

  const currentUserId = computed(() => userStore.info?.userId?.toString() || 'anonymous')

  // 聊天室信息
  const roomInfo = ref<Api.Chat.ChatRoomItem>()
  const members = ref<Api.Chat.ChatRoomMemberItem[]>([])
  const loading = ref(false)

  // WebSocket 连接
  const {
    isConnected,
    isConnecting,
    messages,
    onlineUsers,
    typingUsers,
    connect,
    disconnect,
    sendChatMessage,
    sendTypingStatus
  } = useChatWebSocket(roomId.value)

  // 输入相关
  const inputMessage = ref('')
  const replyMessage = ref<Api.Chat.ChatMessageItem>()
  let typingTimer: NodeJS.Timeout | null = null

  // 弹窗相关
  const showMembers = ref(false)
  const activeTab = ref('members')
  const roomSettings = reactive({
    name: '',
    description: '',
    maxMembers: 0
  })

  // DOM 引用
  const messagesRef = ref<HTMLElement>()

  // 获取聊天室信息
  const fetchRoomInfo = async () => {
    try {
      const res = await chatRoomApi.getRoomDetail(roomId.value)
      roomInfo.value = res

      // 更新设置表单
      Object.assign(roomSettings, {
        name: res.name,
        description: res.description,
        maxMembers: res.maxMembers
      })
    } catch (error) {
      console.error('获取聊天室信息失败:', error)
      ElMessage.error('获取聊天室信息失败')
    }
  }

  // 获取成员列表
  const fetchMembers = async () => {
    try {
      const res = await chatRoomApi.getRoomMembers(roomId.value)
      members.value = res
    } catch (error) {
      console.error('获取成员列表失败:', error)
    }
  }

  // 获取历史消息
  const fetchMessages = async () => {
    loading.value = true
    try {
      const res = await chatMessageApi.getMessageList({
        roomId: roomId.value,
        page: 1,
        pageSize: 50
      } as any)
      messages.value = res || []

      // 滚动到底部
      nextTick(() => {
        scrollToBottom()
      })
    } catch (error) {
      console.error('获取消息列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 发送消息
  const sendMessage = async () => {
    const content = inputMessage.value.trim()
    if (!content) return

    try {
      // 构建请求参数
      const params: any = {
        roomId: roomId.value,
        content,
        messageType: 'text'
      }

      // 只在有回复消息时才添加 replyToId
      if (replyMessage.value?.id) {
        params.replyToId = replyMessage.value.id
      }

      // 通过 HTTP API 发送消息
      await chatMessageApi.sendMessage(params)

      // 清空输入
      inputMessage.value = ''
      cancelReply()

      // 停止正在输入状态
      sendTypingStatus(false)
    } catch (error) {
      console.error('发送消息失败:', error)
      ElMessage.error('发送消息失败')
    }
  }

  // 处理键盘事件
  const handleKeyDown = (event: Event | KeyboardEvent) => {
    if (event instanceof KeyboardEvent && event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault()
      sendMessage()
    }
  }

  // 处理输入事件
  const handleInput = () => {
    // 发送正在输入状态
    sendTypingStatus(true)

    // 3秒后停止正在输入状态
    if (typingTimer) {
      clearTimeout(typingTimer)
    }

    typingTimer = setTimeout(() => {
      sendTypingStatus(false)
    }, 3000)
  }

  // 回复消息
  const handleReply = (message: Api.Chat.ChatMessageItem) => {
    replyMessage.value = message
  }

  // 取消回复
  const cancelReply = () => {
    replyMessage.value = undefined
  }

  // 撤回消息
  const handleRecall = async (message: Api.Chat.ChatMessageItem) => {
    try {
      await ElMessageBox.confirm('确定要撤回这条消息吗？', '提示', {
        type: 'warning'
      })

      await chatMessageApi.recallMessage({ messageId: message.id })
      ElMessage.success('消息已撤回')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('撤回消息失败:', error)
        ElMessage.error('撤回消息失败')
      }
    }
  }

  // 判断是否可以撤回
  const canRecall = (message: Api.Chat.ChatMessageItem) => {
    const now = dayjs()
    const messageTime = dayjs(message.createdAt)
    return now.diff(messageTime, 'minute') <= 2 // 2分钟内可撤回
  }

  // 滚动到底部
  const scrollToBottom = () => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  }

  // 格式化时间
  const formatTime = (time: string) => {
    return dayjs(time).format('HH:mm')
  }

  // 获取正在输入文本
  const getTypingText = () => {
    const names = typingUsers.value.map((u) => u.username).join('、')
    return `${names} 正在输入...`
  }

  // 获取角色类型
  const getRoleType = (role: Api.Chat.MemberRole) => {
    const types = {
      owner: 'danger',
      admin: 'warning',
      member: 'info'
    }
    return types[role] || 'info'
  }

  // 获取角色文本
  const getRoleText = (role: Api.Chat.MemberRole) => {
    const texts = {
      owner: '群主',
      admin: '管理员',
      member: '成员'
    }
    return texts[role] || '成员'
  }

  // 保存设置
  const saveSettings = async () => {
    try {
      await chatRoomApi.updateRoom({
        id: roomId.value,
        ...roomSettings,
        isActive: true
      })

      ElMessage.success('设置已保存')
      showMembers.value = false
      fetchRoomInfo()
    } catch (error) {
      console.error('保存设置失败:', error)
      ElMessage.error('保存设置失败')
    }
  }

  // 返回
  const goBack = () => {
    router.push('/chat/room')
  }

  // 初始化
  onMounted(async () => {
    await fetchRoomInfo()
    await fetchMembers()
    await fetchMessages()

    // 连接 WebSocket
    connect()
  })

  // 清理
  onUnmounted(() => {
    disconnect()

    if (typingTimer) {
      clearTimeout(typingTimer)
    }
  })
</script>

<style scoped lang="scss">
  .chat-room-detail {
    display: flex;
    flex-direction: column;
    background: #fff;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    border-bottom: 1px solid #e4e7ed;

    .chat-header-left {
      display: flex;
      align-items: center;
      gap: 12px;

      .room-info {
        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
        }

        .room-desc {
          font-size: 12px;
          color: #909399;
        }
      }
    }
  }

  .chat-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .chat-messages {
    flex: 1;
    padding: 16px;
    overflow-y: auto;

    .loading-more {
      text-align: center;
      padding: 16px;
      color: #909399;
      font-size: 14px;
    }

    .message-item {
      margin-bottom: 16px;

      &.is-own {
        .message-content {
          flex-direction: row-reverse;

          .message-body {
            align-items: flex-end;

            .message-header {
              justify-content: flex-end;
            }

            .message-text {
              background: #409eff;
              color: white;
            }
          }
        }
      }

      .reply-message {
        display: flex;
        align-items: center;
        gap: 4px;
        margin-bottom: 8px;
        padding: 8px 12px;
        background: #f5f7fa;
        border-radius: 4px;
        font-size: 12px;
        color: #606266;
      }

      .message-content {
        display: flex;
        gap: 12px;

        .message-avatar {
          flex-shrink: 0;
        }

        .message-body {
          flex: 1;
          display: flex;
          flex-direction: column;

          .message-header {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 4px;

            .message-username {
              font-size: 14px;
              font-weight: 500;
              color: #303133;
            }

            .message-time {
              font-size: 12px;
              color: #909399;
            }
          }

          .message-text {
            padding: 8px 12px;
            background: #f5f7fa;
            border-radius: 8px;
            font-size: 14px;
            line-height: 1.4;
            word-break: break-word;

            &.is-recalled {
              color: #909399;
              font-style: italic;
            }

            .recalled-text {
              display: flex;
              align-items: center;
              gap: 4px;
            }
          }

          .message-actions {
            display: flex;
            gap: 4px;
            margin-top: 4px;
            opacity: 0;
            transition: opacity 0.2s;
          }
        }
      }

      &:hover .message-actions {
        opacity: 1;
      }
    }

    .typing-indicator {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 16px;
      color: #909399;
      font-size: 14px;

      .typing-icon {
        animation: typing 1.5s infinite;
      }
    }
  }

  .chat-input {
    border-top: 1px solid #e4e7ed;

    .reply-preview {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 16px;
      background: #f5f7fa;
      border-bottom: 1px solid #e4e7ed;

      .reply-content {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 14px;
        color: #606266;
      }
    }

    .input-area {
      padding: 16px;

      .input-actions {
        display: flex;
        justify-content: flex-end;
        margin-top: 8px;
      }
    }
  }

  .online-users {
    .online-users-title {
      font-weight: 500;
      margin-bottom: 12px;
    }

    .online-users-list {
      max-height: 200px;
      overflow-y: auto;
    }

    .online-user-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;
      font-size: 14px;
    }
  }

  .members-list {
    max-height: 400px;
    overflow-y: auto;

    .member-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      .member-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .member-details {
          .member-name {
            font-weight: 500;
            margin-bottom: 4px;
          }

          .member-role {
            display: flex;
            gap: 4px;
          }
        }
      }
    }
  }

  @keyframes typing {
    0%,
    60%,
    100% {
      transform: translateY(0);
    }
    30% {
      transform: translateY(-10px);
    }
  }
</style>
