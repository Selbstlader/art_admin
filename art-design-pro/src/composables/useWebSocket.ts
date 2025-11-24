import { ref, reactive, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'

interface WebSocketOptions {
  url: string
  reconnectInterval?: number
  maxReconnectAttempts?: number
  heartbeatInterval?: number
  onMessage?: (message: Api.Chat.WebSocketMessage) => void
  onConnect?: () => void
  onDisconnect?: () => void
  onError?: (error: Event) => void
}

export function useWebSocket(options: WebSocketOptions) {
  const {
    url,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5,
    heartbeatInterval = 30000,
    onMessage,
    onConnect,
    onDisconnect,
    onError
  } = options

  const userStore = useUserStore()

  // 状态管理
  const isConnected = ref(false)
  const isConnecting = ref(false)
  const reconnectAttempts = ref(0)

  // WebSocket 实例
  let ws: WebSocket | null = null
  let heartbeatTimer: NodeJS.Timeout | null = null
  let reconnectTimer: NodeJS.Timeout | null = null

  // 连接 WebSocket
  const connect = () => {
    if (isConnecting.value || isConnected.value) {
      return
    }

    isConnecting.value = true

    try {
      // 构建 WebSocket URL，添加认证 token
      const token = userStore.accessToken
      const wsUrl = `${url}${url.includes('?') ? '&' : '?'}token=${encodeURIComponent(token)}`

      ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        console.log('WebSocket 连接成功')
        isConnected.value = true
        isConnecting.value = false
        reconnectAttempts.value = 0

        // 启动心跳
        startHeartbeat()

        onConnect?.()
      }

      ws.onmessage = (event) => {
        try {
          const message: Api.Chat.WebSocketMessage = JSON.parse(event.data)

          // 处理心跳响应
          if (message.type === 'pong') {
            console.log('收到心跳响应')
            return
          }

          onMessage?.(message)
        } catch (error) {
          console.error('解析 WebSocket 消息失败:', error)
        }
      }

      ws.onclose = (event) => {
        console.log('WebSocket 连接关闭:', event.code, event.reason)
        isConnected.value = false
        isConnecting.value = false

        stopHeartbeat()
        onDisconnect?.()

        // 自动重连
        if (reconnectAttempts.value < maxReconnectAttempts) {
          scheduleReconnect()
        } else {
          ElMessage.error('WebSocket 连接失败，请刷新页面重试')
        }
      }

      ws.onerror = (event) => {
        console.error('WebSocket 错误:', event)
        isConnecting.value = false
        onError?.(event)
      }
    } catch (error) {
      console.error('创建 WebSocket 连接失败:', error)
      isConnecting.value = false
    }
  }

  // 断开连接
  const disconnect = () => {
    stopHeartbeat()
    clearReconnectTimer()

    if (ws) {
      ws.close(1000, 'Manual disconnect')
      ws = null
    }

    isConnected.value = false
    isConnecting.value = false
    reconnectAttempts.value = 0
  }

  // 发送消息
  const sendMessage = (message: Api.Chat.WebSocketMessage) => {
    if (!isConnected.value || !ws) {
      console.warn('WebSocket 未连接，无法发送消息')
      return false
    }

    try {
      ws.send(JSON.stringify(message))
      return true
    } catch (error) {
      console.error('发送 WebSocket 消息失败:', error)
      return false
    }
  }

  // 启动心跳
  const startHeartbeat = () => {
    stopHeartbeat()

    heartbeatTimer = setInterval(() => {
      if (isConnected.value) {
        sendMessage({ type: 'ping' })
      }
    }, heartbeatInterval)
  }

  // 停止心跳
  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  // 计划重连
  const scheduleReconnect = () => {
    clearReconnectTimer()

    reconnectAttempts.value++
    console.log(`尝试重连 WebSocket (${reconnectAttempts.value}/${maxReconnectAttempts})`)

    reconnectTimer = setTimeout(() => {
      connect()
    }, reconnectInterval)
  }

  // 清除重连定时器
  const clearReconnectTimer = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
  }

  // 组件卸载时清理
  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected,
    isConnecting,
    reconnectAttempts,
    connect,
    disconnect,
    sendMessage
  }
}

// 聊天室专用的 WebSocket hook
export function useChatWebSocket(roomId: number) {
  const messages = ref<Api.Chat.ChatMessageItem[]>([])
  const onlineUsers = ref<Api.Chat.OnlineUser[]>([])
  const typingUsers = ref<Api.Chat.TypingData[]>([])

  const { isConnected, isConnecting, connect, disconnect, sendMessage } = useWebSocket({
    url: `ws://localhost:48080/api/chat/ws?roomId=${roomId}`,
    onMessage: (message) => {
      handleChatMessage(message)
    },
    onConnect: () => {
      console.log(`已连接到聊天室 ${roomId}`)
    },
    onDisconnect: () => {
      console.log(`已断开聊天室 ${roomId} 连接`)
      // 清空在线用户和正在输入状态
      onlineUsers.value = []
      typingUsers.value = []
    }
  })

  // 处理聊天消息
  const handleChatMessage = (message: Api.Chat.WebSocketMessage) => {
    switch (message.type) {
      case 'join':
        // 用户加入
        if (message.data) {
          ElMessage.success(`${message.data.username} 加入了聊天室`)
          // 更新在线用户列表
          if (!onlineUsers.value.find((u) => u.userId === message.data.userId)) {
            onlineUsers.value.push(message.data)
          }
        }
        break

      case 'leave':
        // 用户离开
        if (message.data) {
          ElMessage.info(`${message.data.username} 离开了聊天室`)
          // 从在线用户列表移除
          onlineUsers.value = onlineUsers.value.filter((u) => u.userId !== message.data.userId)
        }
        break

      case 'typing':
        // 正在输入
        if (message.data) {
          handleTypingStatus(message.data)
        }
        break

      default:
        // 处理未知消息类型，包括error类型
        if ((message as any).type === 'error' && (message as any).message) {
          ElMessage.error((message as any).message)
        } else {
          console.log('收到未知类型消息:', message)
        }
    }
  }

  // 处理正在输入状态
  const handleTypingStatus = (data: Api.Chat.TypingData) => {
    const existingIndex = typingUsers.value.findIndex((u) => u.userId === data.userId)

    if (data.isTyping) {
      if (existingIndex === -1) {
        typingUsers.value.push(data)
      }

      // 3秒后自动移除正在输入状态
      setTimeout(() => {
        const index = typingUsers.value.findIndex((u) => u.userId === data.userId)
        if (index !== -1) {
          typingUsers.value.splice(index, 1)
        }
      }, 3000)
    } else {
      if (existingIndex !== -1) {
        typingUsers.value.splice(existingIndex, 1)
      }
    }
  }

  // 滚动到底部
  const scrollToBottom = () => {
    const messagesContainer = document.querySelector('.chat-messages')
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight
    }
  }

  // 发送聊天消息 - 已移除，改用HTTP API
  // const sendChatMessage = (content: string, replyToId?: number) => {
  //   return sendMessage({
  //     type: 'message',
  //     data: {
  //       roomId,
  //       content,
  //       messageType: 'text',
  //       replyToId
  //     }
  //   })
  // }

  // 发送正在输入状态
  const sendTypingStatus = (isTyping: boolean) => {
    return sendMessage({
      type: 'typing',
      data: {
        roomId,
        isTyping
      }
    })
  }

  return {
    // 状态
    isConnected,
    isConnecting,
    messages,
    onlineUsers,
    typingUsers,

    // 方法
    connect,
    disconnect,
    // sendChatMessage, // 已移除，改用HTTP API
    sendTypingStatus,
    scrollToBottom
  }
}
