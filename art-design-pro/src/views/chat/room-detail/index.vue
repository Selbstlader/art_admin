<template>
  <div class="chat-room-detail art-full-height">
    <!-- 聊天室头部 -->
    <div class="chat-header">
      <div class="chat-header-left">
        <ElButton link @click="goBack">
          <ElIcon>
            <ArrowLeft />
          </ElIcon>
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

          <!-- <ElPopover placement="bottom" width="300" trigger="hover">
            <template #reference>
              <ElButton link>
                <ElIcon>
                  <User />
                </ElIcon>
                在线 {{ onlineMemberUsers.length }}
              </ElButton>
            </template> -->

          <!-- <div class="online-users">
              <div class="online-users-title">在线成员</div>
              {{ onlineMemberUsers }}
              <div class="online-users-list">
                <div v-for="user in onlineMemberUsers" :key="user.userId" class="online-user-item">
                  <ElAvatar :size="24" :src="user.avatar">
                    {{ user.username?.charAt(0) }}
                  </ElAvatar>
                  <span>{{ user.username }}</span>
                </div>
                <div v-if="onlineMemberUsers.length === 0" class="no-online-users">
                  暂无在线成员
                </div>
              </div>
            </div> -->
          <!-- </ElPopover> -->
          <ElButton link>
            <ElIcon>
              <User />
            </ElIcon>
            在线 {{ onlineMemberUsers.length }}
          </ElButton>

          <!-- <ElButton @click="showMembers = true">
            <ElIcon>
              <Setting />
            </ElIcon>
            管理
          </ElButton> -->
        </ElSpace>
      </div>
    </div>

    <!-- 聊天内容区域 -->
    <div class="chat-content">
      <!-- 消息列表 -->
      <div class="chat-messages" ref="messagesRef">
        <div v-if="loading" class="loading-more">
          <ElIcon class="is-loading">
            <Loading />
          </ElIcon>
          加载中...
        </div>

        <!-- 加载更多历史消息 -->
        <div v-if="hasMoreMessages && !loading" class="load-more-container">
          <ElButton link :loading="isLoadingMore" @click="loadMoreMessages" class="load-more-btn">
            {{ isLoadingMore ? '加载中...' : '加载更多历史消息' }}
          </ElButton>
        </div>

        <div
          v-for="message in messages"
          :key="message.id"
          class="message-item"
          :class="{ 'is-own': message.userId === currentUserId }"
        >
          <!-- 回复消息 -->
          <div v-if="message.replyTo" class="reply-message">
            <ElIcon>
              <Reply />
            </ElIcon>
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
                  <ElIcon>
                    <Delete />
                  </ElIcon>
                  消息已撤回
                </span>
                <span v-else>{{ message.content }}</span>
              </div>

              <!-- 消息操作 -->
              <div v-if="!message.isRecalled" class="message-actions">
                <ElButton link size="small" @click="handleReply(message)">
                  <ElIcon>
                    <Reply />
                  </ElIcon>
                </ElButton>
                <ElButton
                  v-if="message.userId === currentUserId && canRecall(message)"
                  link
                  size="small"
                  type="danger"
                  @click="handleRecall(message)"
                >
                  <ElIcon>
                    <Delete />
                  </ElIcon>
                </ElButton>
              </div>
            </div>
          </div>
        </div>

        <!-- 正在输入提示 -->
        <div v-if="typingUsers.length > 0" class="typing-indicator">
          <ElIcon class="typing-icon">
            <Edit />
          </ElIcon>
          {{ getTypingText() }}
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="chat-input">
        <!-- 回复提示 -->
        <div v-if="replyMessage" class="reply-preview">
          <div class="reply-content">
            <ElIcon>
              <Reply />
            </ElIcon>
            回复 {{ replyMessage.username }}: {{ replyMessage.content }}
          </div>
          <ElButton link @click="cancelReply">
            <ElIcon>
              <Close />
            </ElIcon>
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

  const roomId = computed(() => {
    const id = Number(route.params.id)
    return isNaN(id) ? 0 : id
  })
  console.log(userStore.info, 'userStore.info?.userId')

  const currentUserId = computed(() => userStore.info?.userId || 0)

  // 聊天室信息
  const roomInfo = ref<Api.Chat.ChatRoomItem>()
  const members = ref<Api.Chat.ChatRoomMemberItem[]>([])
  const loading = ref(false)

  // 计算在线成员（从成员列表中筛选，不依赖WebSocket）
  const onlineMemberUsers = computed(() => {
    return members.value.map((member) => ({
      userId: member.userId,
      username: member.username,
      avatar: member.avatar
    }))
  })

  // 消息和状态管理（不依赖WebSocket）
  const messages = ref<Api.Chat.ChatMessageItem[]>([])
  const typingUsers = ref<Api.Chat.TypingData[]>([])

  // 输入相关
  const inputMessage = ref('')
  const replyMessage = ref<Api.Chat.ChatMessageItem>()
  let typingTimer: NodeJS.Timeout | null = null

  // 分页相关
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalCount = ref(0)
  const hasMoreMessages = ref(true)
  const isLoadingMore = ref(false)

  // 轮询相关
  const pollingTimer = ref<NodeJS.Timeout | null>(null)
  const lastMessageId = ref<number>(0)
  const isPolling = ref(false)

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

  // 加载消息列表
  const loadMessages = async (loadMore = false) => {
    if (loadMore) {
      isLoadingMore.value = true
    } else {
      loading.value = true
      currentPage.value = 1
      messages.value = []
    }

    try {
      const response = await chatMessageApi.getMessageList({
        roomId: roomId.value,
        page: currentPage.value,
        pageSize: pageSize.value
      } as any)

      if (response && response.data) {
        const newMessages = response.data as Api.Chat.ChatMessageItem[]

        if (loadMore) {
          // 加载更多时，将历史消息添加到前面
          messages.value = [...newMessages.reverse(), ...messages.value]
        } else {
          // 初始加载时，消息按时间正序排列（最新的在底部）
          messages.value = newMessages.reverse()
        }

        // 更新分页信息
        totalCount.value = response.total || 0
        hasMoreMessages.value = messages.value.length < totalCount.value

        // 设置最后一条消息ID用于轮询
        if (newMessages.length > 0) {
          lastMessageId.value = Math.max(...newMessages.map((msg) => msg.id))
        }

        // 滚动到底部（仅初始加载）
        if (!loadMore) {
          nextTick(() => {
            scrollToBottom()
          })
        }
      }
    } catch (error) {
      console.error('获取消息列表失败:', error)
    } finally {
      loading.value = false
      isLoadingMore.value = false
    }
  }

  // 加载更多历史消息
  const loadMoreMessages = async () => {
    if (!hasMoreMessages.value || isLoadingMore.value) return

    // 获取当前第一条消息的ID，作为分页基准
    const firstMessageId = messages.value.length > 0 ? messages.value[0].id : 0

    try {
      isLoadingMore.value = true
      currentPage.value++

      const response = await chatMessageApi.getMessageList({
        roomId: roomId.value,
        page: currentPage.value,
        pageSize: pageSize.value
        // beforeId: firstMessageId
      } as any)

      if (response && response.data) {
        const newMessages = response.data as Api.Chat.ChatMessageItem[]

        if (newMessages.length > 0) {
          // 将历史消息添加到前面
          messages.value = [...newMessages.reverse(), ...messages.value]

          // 更新是否还有更多消息
          totalCount.value = response.total || 0
          hasMoreMessages.value = messages.value.length < totalCount.value
        } else {
          hasMoreMessages.value = false
        }
      }
    } catch (error) {
      console.error('加载更多消息失败:', error)
      currentPage.value-- // 恢复页码
    } finally {
      isLoadingMore.value = false
    }
  }

  // 开始轮询获取新消息
  const startPolling = () => {
    if (isPolling.value) return

    isPolling.value = true
    console.log('开始轮询新消息...')

    // 立即获取一次最新消息
    pollNewMessages()

    // 每3秒轮询一次
    pollingTimer.value = setInterval(() => {
      pollNewMessages()
    }, 3000)
  }

  // 停止轮询
  const stopPolling = () => {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
      pollingTimer.value = null
    }
    isPolling.value = false
    console.log('停止轮询新消息')
  }

  // 轮询获取新消息
  const pollNewMessages = async () => {
    try {
      // 获取当前最后一条消息的ID
      const lastId = messages.value.length > 0 ? messages.value[messages.value.length - 1].id : 0

      const response = await chatMessageApi.getMessageList({
        roomId: roomId.value,
        page: 1,
        pageSize: pageSize.value,
        afterId: lastId
      } as any)

      if (response && response.data && response.data.length > 0) {
        const newMessages = response.data as Api.Chat.ChatMessageItem[]

        // 将新消息添加到列表末尾
        messages.value.push(...newMessages)

        // 滚动到底部
        nextTick(() => {
          scrollToBottom()
        })
      }
    } catch (error) {
      console.error('轮询获取新消息失败:', error)
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
      const response = await chatMessageApi.sendMessage(params)

      // 将发送成功的消息添加到列表中
      if (response) {
        // 发送消息API直接返回消息对象，不是包装在data中
        messages.value.push(response as Api.Chat.ChatMessageItem)

        // 滚动到底部
        nextTick(() => {
          scrollToBottom()
        })
      }

      // 清空输入
      inputMessage.value = ''
      cancelReply()

      // 停止正在输入状态
      // sendTypingStatus(false)
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
    // 暂时移除正在输入状态功能，因为未连接WebSocket
    // sendTypingStatus(true)

    // 3秒后停止正在输入状态
    if (typingTimer) {
      clearTimeout(typingTimer)
    }

    typingTimer = setTimeout(() => {
      // sendTypingStatus(false)
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
  const getRoleType = (
    role: Api.Chat.MemberRole
  ): 'danger' | 'warning' | 'info' | 'success' | 'primary' => {
    const types = {
      owner: 'danger' as const,
      admin: 'warning' as const,
      member: 'info' as const
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
    router.push('/chat/chat/room')
  }

  // 初始化
  onMounted(async () => {
    await fetchRoomInfo()
    await fetchMembers()
    await loadMessages() // 加载初始消息

    // 连接 WebSocket（仅用于在线状态和正在输入功能）
    // connect()

    // 启动消息轮询
    startPolling()

    // 监听页面可见性变化
    document.addEventListener('visibilitychange', handleVisibilityChange)
  })

  // 处理页面可见性变化
  const handleVisibilityChange = () => {
    if (document.hidden) {
      // 页面不可见时停止轮询
      stopPolling()
      console.log('页面不可见，停止轮询')
    } else {
      // 页面可见时重新启动轮询
      if (roomId.value > 0) {
        startPolling()
        console.log('页面可见，重新启动轮询')
      }
    }
  }

  // 清理
  onUnmounted(() => {
    stopPolling()

    // 移除页面可见性监听
    document.removeEventListener('visibilitychange', handleVisibilityChange)

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

    .load-more-container {
      text-align: center;
      padding: 8px 0;
      margin-bottom: 16px;

      .load-more-btn {
        color: #409eff;
        font-size: 14px;

        &:hover {
          color: #66b1ff;
        }
      }
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

      .online-user-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 4px 0;

        span {
          font-size: 14px;
          color: var(--el-text-color-regular);
        }
      }

      .no-online-users {
        padding: 20px 0;
        text-align: center;
        color: var(--el-text-color-placeholder);
        font-size: 14px;
      }
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
