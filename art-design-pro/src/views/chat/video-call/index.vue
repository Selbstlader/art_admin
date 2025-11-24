<template>
  <div class="video-call-page art-full-height">
    <!-- 未开始通话 - 显示房间列表 -->
    <div v-if="!isInCall" class="call-lobby">
      <ElCard class="art-table-card lobby-card" shadow="never">
        <template #header>
          <div class="card-header">
            <h2 class="text-primary">视频通话</h2>
            <p class="subtitle text-gray-500">选择或创建一个房间开始视频通话</p>
          </div>
        </template>

        <!-- 创建房间 -->
        <ElForm :model="createForm" label-width="80px" class="create-room-form">
          <ElFormItem label="房间名称">
            <ElInput
              v-model="createForm.roomName"
              placeholder="输入房间名称"
              @keyup.enter="createRoom"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="createRoom" :loading="creating">
              <ElIcon><Plus /></ElIcon>
              创建房间
            </ElButton>
          </ElFormItem>
        </ElForm>

        <ElDivider>或加入现有房间</ElDivider>

        <!-- 房间列表 -->
        <div v-loading="loadingRooms" class="room-list">
          <div v-for="room in activeRooms" :key="room.id" class="room-item" @click="joinRoom(room)">
            <div class="room-info">
              <h3 class="text-gray-800">
                {{ room.name }}
                <ElTag v-if="room.isActive === false" type="info" size="small" class="status-tag">
                  已关闭
                </ElTag>
              </h3>
              <p class="room-desc text-gray-500">{{ room.description || '暂无描述' }}</p>
            </div>
            <div class="room-meta">
              <ElTag :type="room.onlineCount > 0 ? 'success' : 'info'">
                <ElIcon><User /></ElIcon>
                {{ room.memberCount || 0 }} 人
              </ElTag>
              <ElButton type="primary" size="small">加入</ElButton>
            </div>
          </div>

          <ElEmpty v-if="activeRooms.length === 0 && !loadingRooms" description="暂无房间" />
        </div>
      </ElCard>
    </div>

    <!-- 视频通话中 -->
    <div v-else class="call-container">
      <!-- 视频网格 -->
      <div class="video-grid" :class="`grid-${videoCount}`">
        <!-- 本地视频 -->
        <div class="video-wrapper local-video">
          <video ref="localVideoRef" autoplay muted playsinline></video>
          <div class="video-overlay">
            <div class="user-info">
              <span class="username">{{ currentUser.username }} (你)</span>
            </div>
            <div class="media-status">
              <el-icon v-if="!isVideoEnabled" class="status-icon disabled">
                <VideoCameraFilled />
              </el-icon>
              <el-icon v-if="!isAudioEnabled" class="status-icon disabled">
                <Microphone />
              </el-icon>
            </div>
          </div>
        </div>

        <!-- 远程视频 -->
        <div v-for="peer in remotePeers" :key="peer.userId" class="video-wrapper remote-video">
          <video
            :ref="(el) => setRemoteVideoRef(peer.userId, el as HTMLVideoElement)"
            autoplay
            playsinline
          ></video>
          <div class="video-overlay">
            <div class="user-info">
              <span class="username">{{ peer.username }}</span>
            </div>
            <div class="media-status">
              <el-icon v-if="!peer.videoEnabled" class="status-icon disabled">
                <VideoCameraFilled />
              </el-icon>
              <el-icon v-if="!peer.audioEnabled" class="status-icon disabled">
                <Microphone />
              </el-icon>
            </div>
          </div>
        </div>
      </div>

      <!-- 控制栏 -->
      <div class="control-bar">
        <div class="control-buttons">
          <el-tooltip :content="isAudioEnabled ? '静音' : '取消静音'">
            <el-button
              :type="isAudioEnabled ? 'primary' : 'danger'"
              :icon="Microphone"
              circle
              size="large"
              @click="toggleAudio"
            />
          </el-tooltip>

          <el-tooltip :content="isVideoEnabled ? '关闭摄像头' : '开启摄像头'">
            <el-button
              :type="isVideoEnabled ? 'primary' : 'danger'"
              :icon="VideoCameraFilled"
              circle
              size="large"
              @click="toggleVideo"
            />
          </el-tooltip>

          <el-tooltip content="挂断">
            <el-button
              type="danger"
              :icon="PhoneFilled"
              circle
              size="large"
              @click="handleHangup"
            />
          </el-tooltip>
        </div>

        <div class="room-info-bar">
          <span class="room-name">{{ currentRoom?.name }}</span>
          <span class="participant-count">
            <el-icon><User /></el-icon>
            {{ videoCount }} 人在线
          </span>
          <span v-if="showIdleCountdown && idleCountdownSeconds > 0" class="idle-countdown">
            房间空闲，将在
            {{ Math.ceil(idleCountdownSeconds / 60) }}
            分钟后自动关闭
          </span>
        </div>
      </div>

      <!-- 通话状态提示 -->
      <div v-if="callStatus !== 'connected'" class="call-status-overlay">
        <div class="status-content">
          <el-icon class="is-loading" :size="40"><Loading /></el-icon>
          <p>{{ callStatusText }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import {
    User,
    VideoCameraFilled,
    Microphone,
    PhoneFilled,
    Loading,
    Plus
  } from '@element-plus/icons-vue'
  import { chatRoomApi } from '@/api/chat'
  import { useUserStore } from '@/store/modules/user'
  import { WebRTCManager, CallStatus, type PeerConnection } from '@/utils/webrtc'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()

  // 当前用户信息
  const currentUser = computed(() => ({
    userId: userStore.info?.userId || 0,
    username: userStore.info?.userName || '未知用户'
  }))

  // 房间相关
  const rooms = ref<Api.Chat.ChatRoomItem[]>([])
  const loadingRooms = ref(false)
  const currentRoom = ref<Api.Chat.ChatRoomItem>()
  const creating = ref(false)
  const createForm = ref({
    roomName: ''
  })

  // 通话状态
  const isInCall = ref(false)
  const callStatus = ref<CallStatus>(CallStatus.IDLE)
  const isAudioEnabled = ref(true)
  const isVideoEnabled = ref(true)

  // 房间空闲倒计时（前端提示用，与后端5分钟逻辑保持一致）
  const idleCountdownSeconds = ref(0)
  const showIdleCountdown = ref(false)
  let idleCountdownTimer: number | null = null

  // 视频相关
  const localVideoRef = ref<HTMLVideoElement>()
  const remoteVideoRefs = new Map<number, HTMLVideoElement>()
  const remotePeers = ref<PeerConnection[]>()
  const onlineCount = ref(0) // 真实在线人数（从后端获取）

  // WebRTC 和 WebSocket
  let webrtcManager: WebRTCManager | null = null
  let ws: WebSocket | null = null

  // 计算属性
  const videoCount = computed(() => onlineCount.value)

  // 只展示启用中的房间（isActive !== false），已关闭房间可以通过 tag 标记
  const activeRooms = computed(() => {
    return rooms.value.filter((room) => room.isActive !== false)
  })

  const callStatusText = computed(() => {
    switch (callStatus.value) {
      case CallStatus.CALLING:
        return '正在连接...'
      case CallStatus.RINGING:
        return '等待其他人加入...'
      case CallStatus.CONNECTED:
        return '通话中'
      default:
        return ''
    }
  })

  // 设置远程视频引用
  const setRemoteVideoRef = (userId: number, el: HTMLVideoElement | null) => {
    if (el) {
      remoteVideoRefs.set(userId, el)
      // 如果 peer 已经有 stream，立即设置
      const peer = remotePeers.value?.find((p) => p.userId === userId)
      if (peer?.stream) {
        console.log('Setting existing stream for user:', userId, peer.stream)
        el.srcObject = peer.stream
      }
    } else {
      // 元素被移除时清理引用
      remoteVideoRefs.delete(userId)
    }
  }

  // 获取房间列表
  const fetchRooms = async () => {
    try {
      loadingRooms.value = true
      const res = await chatRoomApi.getRoomList({
        page: 1,
        pageSize: 20,
        type: 'public',
        isActive: true
      })
      rooms.value = res || []
    } catch (error) {
      console.error('Failed to fetch rooms:', error)
      ElMessage.error('获取房间列表失败')
    } finally {
      loadingRooms.value = false
    }
  }

  // 创建房间
  const createRoom = async () => {
    if (!createForm.value.roomName.trim()) {
      ElMessage.warning('请输入房间名称')
      return
    }

    try {
      creating.value = true
      const room = await chatRoomApi.createRoom({
        name: createForm.value.roomName,
        description: '视频通话房间',
        type: 'public',
        maxMembers: 4
      })
      ElMessage.success('房间创建成功')
      createForm.value.roomName = ''
      await joinRoom(room)
    } catch (error) {
      console.error('Failed to create room:', error)
      ElMessage.error('创建房间失败')
    } finally {
      creating.value = false
    }
  }

  // 加入房间
  const joinRoom = async (room: Api.Chat.ChatRoomItem) => {
    try {
      // 先加入聊天室
      await chatRoomApi.joinRoom({ roomId: room.id })

      currentRoom.value = room
      isInCall.value = true

      // 建立 WebSocket 连接
      await connectWebSocket(room.id)

      // 初始化 WebRTC
      await initWebRTC(room.id)

      // 立即获取一次在线人数
      await fetchOnlineCount()
    } catch (error) {
      console.error('Failed to join room:', error)
      ElMessage.error('加入房间失败')
      isInCall.value = false
    }
  }

  // 连接 WebSocket
  const connectWebSocket = (roomId: number): Promise<void> => {
    return new Promise((resolve, reject) => {
      const token = userStore.accessToken || ''
      const wsUrl = `ws://localhost:48080/api/chat/ws?roomId=${roomId}&token=${token}`

      ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        console.log('WebSocket connected')
        resolve()
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        reject(error)
      }

      ws.onclose = () => {
        console.log('WebSocket closed')
      }

      ws.onmessage = (event) => {
        try {
          // 检查消息类型
          if (typeof event.data !== 'string') {
            console.warn('Received non-string WebSocket message:', event.data)
            return
          }

          const data = JSON.parse(event.data)
          console.log('WebSocket message:', data)

          // 监听房间关闭消息
          if (data.type === 'room_closed') {
            console.log('房间已关闭:', data.data)
            ElMessage.warning({
              message: data.data.reason || '房间已关闭',
              duration: 3000
            })
            // 自动退出通话
            handleCallEnded()
            return
          }

          // 监听用户加入/离开消息，实时更新在线人数
          if (data.type === 'user_joined' || data.type === 'user_left') {
            fetchOnlineCount()
          }

          // 转发给 WebRTC Manager
          if (webrtcManager) {
            webrtcManager.handleSignalMessage(data)
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error, 'Raw data:', event.data)
        }
      }
    })
  }

  // 初始化 WebRTC
  const initWebRTC = async (roomId: number) => {
    try {
      if (!ws) {
        throw new Error('WebSocket not connected')
      }

      webrtcManager = new WebRTCManager({
        roomId,
        userId: currentUser.value.userId,
        username: currentUser.value.username,
        webSocket: ws,
        onCallStatusChange: (status) => {
          // 这里只同步状态，避免再次触发 handleCallEnded 导致递归
          callStatus.value = status
        },
        onRemoteStream: (userId, username, stream) => {
          console.log('Received remote stream:', userId, username, stream)
          // 先更新 peer 列表，触发 DOM 渲染
          updateRemotePeers()
          // 等待 DOM 更新完成后设置视频流
          nextTick(() => {
            const videoEl = remoteVideoRefs.get(userId)
            console.log('Setting video stream for user:', userId, videoEl, stream)
            if (videoEl && stream) {
              videoEl.srcObject = stream
              console.log('Video stream set successfully for user:', userId)
            } else {
              console.error('Failed to set video stream:', { userId, videoEl, stream })
            }
          })
        },
        onPeerLeft: (userId) => {
          console.log('Peer left:', userId)
          remoteVideoRefs.delete(userId)
          updateRemotePeers()
        },
        onError: (error) => {
          console.error('WebRTC error:', error)
          ElMessage.error('视频通话错误: ' + error.message)
        }
      })

      // 开始通话
      await webrtcManager.startCall()

      // 设置本地视频
      const localStream = webrtcManager.getLocalStream()
      if (localStream && localVideoRef.value) {
        localVideoRef.value.srcObject = localStream
      }

      callStatus.value = CallStatus.CONNECTED
    } catch (error: any) {
      console.error('Failed to initialize WebRTC:', error)
      ElMessage.error('无法访问摄像头或麦克风: ' + error.message)
      isInCall.value = false
    }
  }

  // 更新远程 Peer 列表
  const updateRemotePeers = () => {
    if (webrtcManager) {
      remotePeers.value = webrtcManager.getAllPeers()
    }
  }

  // 获取房间在线人数
  const fetchOnlineCount = async () => {
    if (!currentRoom.value) return
    try {
      const users = await chatRoomApi.getOnlineUsers(currentRoom.value.id)
      onlineCount.value = users.length

      // 当前房间无人在线时，启动前端提示倒计时；有人在线则取消
      if (onlineCount.value === 0 && isInCall.value) {
        startIdleCountdown()
      } else {
        stopIdleCountdown()
      }
    } catch (error) {
      console.error('获取在线人数失败:', error)
    }
  }

  // 切换音频
  const toggleAudio = () => {
    if (webrtcManager) {
      isAudioEnabled.value = webrtcManager.toggleAudio()
    }
  }

  // 切换视频
  const toggleVideo = () => {
    if (webrtcManager) {
      isVideoEnabled.value = webrtcManager.toggleVideo()
    }
  }

  // 挂断
  const handleHangup = () => {
    // 先让 WebRTCManager 处理信令和内部清理
    if (webrtcManager) {
      webrtcManager.hangup()
      // 置空引用，避免后续回调再次进入 destroy 形成递归
      webrtcManager = null
    }
    // 再做本页面的 UI/状态收尾
    handleCallEnded()
  }

  // 通话结束
  const handleCallEnded = () => {
    isInCall.value = false
    callStatus.value = CallStatus.IDLE

    if (ws) {
      ws.close()
      ws = null
    }

    remotePeers.value = []
    remoteVideoRefs.clear()
    onlineCount.value = 0 // 重置在线人数

    stopIdleCountdown()

    // 刷新房间列表
    fetchRooms()
  }

  // 定时更新远程 Peer
  let updateInterval: number | null = null

  onMounted(async () => {
    // 检查是否从聊天室列表跳转过来，带有 roomId 参数
    const roomIdParam = route.query.roomId as string
    if (roomIdParam) {
      try {
        const roomId = parseInt(roomIdParam)
        // 获取房间详情
        const room = await chatRoomApi.getRoomDetail(roomId)
        // 自动加入该房间
        await joinRoom(room)
      } catch (error) {
        console.error('自动加入房间失败:', error)
        ElMessage.error('加入房间失败')
        // 加入失败，显示房间列表
        fetchRooms()
      }
    } else {
      // 没有 roomId 参数，显示房间列表
      fetchRooms()
    }

    updateInterval = window.setInterval(() => {
      if (isInCall.value) {
        updateRemotePeers()
        fetchOnlineCount() // 定时获取在线人数（兜底方案）
      }
    }, 5000) // 每5秒更新一次作为兜底，主要通过 WebSocket 实时更新
  })

  onUnmounted(() => {
    if (webrtcManager) {
      webrtcManager.destroy()
    }
    if (ws) {
      ws.close()
    }
    if (updateInterval) {
      clearInterval(updateInterval)
    }

    stopIdleCountdown()
  })

  // 房间空闲倒计时相关
  const startIdleCountdown = () => {
    // 与后端逻辑保持一致：5分钟后关闭
    idleCountdownSeconds.value = 5 * 60
    showIdleCountdown.value = true

    if (idleCountdownTimer) {
      window.clearInterval(idleCountdownTimer)
    }

    idleCountdownTimer = window.setInterval(() => {
      if (idleCountdownSeconds.value > 0) {
        idleCountdownSeconds.value -= 1
      } else {
        stopIdleCountdown()
      }
    }, 1000)
  }

  const stopIdleCountdown = () => {
    showIdleCountdown.value = false
    idleCountdownSeconds.value = 0
    if (idleCountdownTimer) {
      window.clearInterval(idleCountdownTimer)
      idleCountdownTimer = null
    }
  }
</script>

<style scoped lang="scss">
  @use '@styles/variables.scss' as *;

  .video-call-page {
    background-color: var(--art-bg-color);
    padding: 20px;

    .call-lobby {
      width: 100%;
      max-width: 1000px;
      margin: 0 auto;

      .lobby-card {
        border: 1px solid var(--art-card-border);
        box-shadow: var(--art-card-shadow);
      }

      .card-header {
        text-align: center;

        h2 {
          margin: 0 0 8px 0;
          font-size: 24px;
          font-weight: 600;
        }

        .subtitle {
          margin: 0;
          font-size: 14px;
        }
      }

      .create-room-form {
        margin-top: 20px;
      }

      .room-list {
        max-height: 400px;
        overflow-y: auto;
      }

      .room-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px;
        margin-bottom: 12px;
        background: var(--art-main-bg-color);
        border: 1px solid var(--art-border-color);
        border-radius: calc(var(--custom-radius) / 2 + 2px);
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          background: var(--art-hoverColor);
          border-color: rgb(var(--art-primary));
          transform: translateY(-1px);
          box-shadow: var(--art-box-shadow-sm);
        }

        .room-info {
          flex: 1;

          h3 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 500;
          }

          .room-desc {
            margin: 0;
            font-size: 14px;
          }
        }

        .room-meta {
          display: flex;
          align-items: center;
          gap: 12px;
        }
      }
    }

    .call-container {
      width: 100%;
      height: 100vh;
      display: flex;
      flex-direction: column;
      background: var(--art-main-bg-color);
    }

    .video-grid {
      flex: 1;
      display: grid;
      gap: 10px;
      padding: 10px;
      overflow: auto;

      &.grid-1 {
        grid-template-columns: 1fr;
      }

      &.grid-2 {
        grid-template-columns: repeat(2, 1fr);
      }

      &.grid-3,
      &.grid-4 {
        grid-template-columns: repeat(2, 1fr);
        grid-template-rows: repeat(2, 1fr);
      }
    }

    .video-wrapper {
      position: relative;
      background: var(--art-gray-800);
      border: 1px solid var(--art-border-color);
      border-radius: calc(var(--custom-radius) / 2 + 2px);
      overflow: hidden;
      min-height: 200px;
      box-shadow: var(--art-box-shadow-sm);

      video {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      &.local-video video {
        transform: scaleX(-1);
      }
    }

    .video-overlay {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      padding: 12px;
      background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
      display: flex;
      justify-content: space-between;
      align-items: center;

      .username {
        color: var(--art-text-gray-100);
        font-size: 14px;
        font-weight: 500;
      }

      .media-status {
        display: flex;
        gap: 8px;

        .status-icon {
          color: var(--el-color-danger);
          font-size: 20px;

          &.disabled {
            opacity: 0.8;
          }
        }
      }
    }

    .control-bar {
      padding: 20px;
      background: var(--art-main-bg-color);
      border-top: 1px solid var(--art-border-color);
      backdrop-filter: blur(10px);
      display: flex;
      flex-direction: column;
      gap: 16px;
      box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.05);

      .control-buttons {
        display: flex;
        justify-content: center;
        gap: 20px;
      }

      .room-info-bar {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 20px;
        color: var(--art-text-gray-700);
        font-size: 14px;

        .room-name {
          font-weight: 500;
          color: var(--art-text-gray-800);
        }

        .participant-count {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }

    .call-status-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(var(--art-gray-800-rgb), 0.8);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 10;

      .status-content {
        text-align: center;
        color: var(--art-text-gray-100);

        p {
          margin-top: 16px;
          font-size: 16px;
        }
      }
    }
  }
</style>
