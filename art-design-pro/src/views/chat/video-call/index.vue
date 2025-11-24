<template>
  <div class="video-call-page">
    <!-- 未开始通话 - 显示房间列表 -->
    <div v-if="!isInCall" class="call-lobby">
      <el-card class="lobby-card">
        <template #header>
          <div class="card-header">
            <h2>视频通话</h2>
            <p class="subtitle">选择或创建一个房间开始视频通话</p>
          </div>
        </template>

        <!-- 创建房间 -->
        <el-form :model="createForm" label-width="80px" class="create-room-form">
          <el-form-item label="房间名称">
            <el-input
              v-model="createForm.roomName"
              placeholder="输入房间名称"
              @keyup.enter="createRoom"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="createRoom" :loading="creating"> 创建房间 </el-button>
          </el-form-item>
        </el-form>

        <el-divider>或加入现有房间</el-divider>

        <!-- 房间列表 -->
        <div v-loading="loadingRooms" class="room-list">
          <div v-for="room in rooms" :key="room.id" class="room-item" @click="joinRoom(room)">
            <div class="room-info">
              <h3>{{ room.name }}</h3>
              <p class="room-desc">{{ room.description || '暂无描述' }}</p>
            </div>
            <div class="room-meta">
              <el-tag type="success">
                <el-icon><User /></el-icon>
                {{ room.memberCount || 0 }} 人
              </el-tag>
              <el-button type="primary" size="small">加入</el-button>
            </div>
          </div>

          <el-empty v-if="rooms.length === 0 && !loadingRooms" description="暂无房间" />
        </div>
      </el-card>
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
    Loading
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
        current: 1,
        size: 20,
        type: 'public'
      })
      rooms.value = res.records || []
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
          callStatus.value = status
          if (status === CallStatus.ENDED) {
            handleCallEnded()
          }
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
    if (webrtcManager) {
      webrtcManager.hangup()
    }
    handleCallEnded()
  }

  // 通话结束
  const handleCallEnded = () => {
    isInCall.value = false
    callStatus.value = CallStatus.IDLE

    if (webrtcManager) {
      webrtcManager.destroy()
      webrtcManager = null
    }

    if (ws) {
      ws.close()
      ws = null
    }

    remotePeers.value = []
    remoteVideoRefs.clear()
    onlineCount.value = 0 // 重置在线人数

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
  })
</script>

<style scoped lang="scss">
  .video-call-page {
    height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .call-lobby {
    width: 100%;
    max-width: 800px;
    padding: 20px;

    .lobby-card {
      box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
    }

    .card-header {
      text-align: center;

      h2 {
        margin: 0 0 8px 0;
        color: #303133;
      }

      .subtitle {
        margin: 0;
        color: #909399;
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
      background: #f5f7fa;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        background: #ecf5ff;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }

      .room-info {
        flex: 1;

        h3 {
          margin: 0 0 4px 0;
          font-size: 16px;
          color: #303133;
        }

        .room-desc {
          margin: 0;
          font-size: 14px;
          color: #909399;
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
    background: #1a1a1a;
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
    background: #000;
    border-radius: 12px;
    overflow: hidden;
    min-height: 200px;

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
      color: #fff;
      font-size: 14px;
      font-weight: 500;
    }

    .media-status {
      display: flex;
      gap: 8px;

      .status-icon {
        color: #f56c6c;
        font-size: 20px;

        &.disabled {
          opacity: 0.8;
        }
      }
    }
  }

  .control-bar {
    padding: 20px;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    display: flex;
    flex-direction: column;
    gap: 16px;

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
      color: #fff;
      font-size: 14px;

      .room-name {
        font-weight: 500;
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
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;

    .status-content {
      text-align: center;
      color: #fff;

      p {
        margin-top: 16px;
        font-size: 16px;
      }
    }
  }
</style>
