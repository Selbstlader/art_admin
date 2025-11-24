<template>
  <div v-if="isCallActive" class="video-call-container">
    <!-- 视频网格 -->
    <div class="video-grid" :class="`grid-${videoCount}`">
      <!-- 本地视频 -->
      <div class="video-wrapper local-video">
        <video ref="localVideoRef" autoplay muted playsinline></video>
        <div class="video-overlay">
          <div class="user-info">
            <span class="username">{{ localUsername }} (你)</span>
          </div>
          <div class="media-status">
            <el-icon v-if="!isVideoEnabled" class="status-icon"><VideoCamera /></el-icon>
            <el-icon v-if="!isAudioEnabled" class="status-icon"><Microphone /></el-icon>
          </div>
        </div>
      </div>

      <!-- 远程视频 -->
      <div v-for="peer in remotePeers" :key="peer.userId" class="video-wrapper remote-video">
        <video :ref="(el) => setRemoteVideoRef(peer.userId, el)" autoplay playsinline></video>
        <div class="video-overlay">
          <div class="user-info">
            <span class="username">{{ peer.username }}</span>
          </div>
          <div class="media-status">
            <el-icon v-if="!peer.videoEnabled" class="status-icon"><VideoCamera /></el-icon>
            <el-icon v-if="!peer.audioEnabled" class="status-icon"><Microphone /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 控制栏 -->
    <div class="control-bar">
      <el-button
        :type="isAudioEnabled ? 'primary' : 'danger'"
        :icon="isAudioEnabled ? Microphone : MicrophoneSlash"
        circle
        size="large"
        @click="toggleAudio"
      />
      <el-button
        :type="isVideoEnabled ? 'primary' : 'danger'"
        :icon="isVideoEnabled ? VideoCamera : VideoCameraSlash"
        circle
        size="large"
        @click="toggleVideo"
      />
      <el-button type="danger" :icon="PhoneFilled" circle size="large" @click="handleHangup" />
    </div>

    <!-- 通话状态提示 -->
    <div v-if="callStatus !== 'connected'" class="call-status">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>{{ callStatusText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
  import { VideoCamera, Microphone, PhoneFilled, Loading } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import { WebRTCManager, CallStatus, type PeerConnection } from '@/utils/webrtc'

  interface Props {
    roomId: number
    userId: number
    username: string
    webSocket: WebSocket
    isCallActive: boolean
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'call-ended'): void
    (e: 'error', error: Error): void
  }>()

  // Refs
  const localVideoRef = ref<HTMLVideoElement>()
  const remoteVideoRefs = new Map<number, HTMLVideoElement>()

  // WebRTC Manager
  let webrtcManager: WebRTCManager | null = null

  // 状态
  const callStatus = ref<CallStatus>(CallStatus.IDLE)
  const isAudioEnabled = ref(true)
  const isVideoEnabled = ref(true)
  const remotePeers = ref<PeerConnection[]>([])

  // 计算属性
  const videoCount = computed(() => {
    return 1 + remotePeers.value.length
  })

  const callStatusText = computed(() => {
    switch (callStatus.value) {
      case CallStatus.CALLING:
        return '正在连接...'
      case CallStatus.RINGING:
        return '对方正在响铃...'
      case CallStatus.CONNECTED:
        return '通话中'
      default:
        return ''
    }
  })

  const localUsername = computed(() => props.username)

  // 设置远程视频引用
  const setRemoteVideoRef = (userId: number, el: any) => {
    if (el) {
      remoteVideoRefs.set(userId, el as HTMLVideoElement)
    }
  }

  // 初始化 WebRTC
  const initWebRTC = async () => {
    try {
      webrtcManager = new WebRTCManager({
        roomId: props.roomId,
        userId: props.userId,
        username: props.username,
        webSocket: props.webSocket,
        onCallStatusChange: (status) => {
          callStatus.value = status
          if (status === CallStatus.ENDED) {
            emit('call-ended')
          }
        },
        onRemoteStream: (userId, username, stream) => {
          console.log('Received remote stream:', userId, username)
          // 等待 DOM 更新后设置视频流
          setTimeout(() => {
            const videoEl = remoteVideoRefs.get(userId)
            if (videoEl) {
              videoEl.srcObject = stream
            }
          }, 100)
        },
        onPeerLeft: (userId) => {
          console.log('Peer left:', userId)
          remoteVideoRefs.delete(userId)
          updateRemotePeers()
        },
        onError: (error) => {
          console.error('WebRTC error:', error)
          ElMessage.error('视频通话错误: ' + error.message)
          emit('error', error)
        }
      })

      // 开始通话
      await webrtcManager.startCall()

      // 设置本地视频
      const localStream = webrtcManager.getLocalStream()
      if (localStream && localVideoRef.value) {
        localVideoRef.value.srcObject = localStream
      }

      // 监听 WebSocket 消息
      setupWebSocketListener()
    } catch (error: any) {
      console.error('Failed to initialize WebRTC:', error)
      ElMessage.error('无法访问摄像头或麦克风: ' + error.message)
      emit('error', error)
    }
  }

  // 设置 WebSocket 监听
  const setupWebSocketListener = () => {
    if (!webrtcManager) return

    const originalOnMessage = props.webSocket.onmessage

    props.webSocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        // 处理 WebRTC 信令消息
        const webrtcMessageTypes = [
          'webrtc_offer',
          'webrtc_answer',
          'webrtc_ice_candidate',
          'call_request',
          'call_accept',
          'call_reject',
          'call_cancel',
          'call_hangup',
          'toggle_audio',
          'toggle_video',
          'user_joined',
          'user_left'
        ]

        if (webrtcMessageTypes.includes(data.type)) {
          webrtcManager?.handleSignalMessage(data)
        }

        // 调用原始处理器
        if (originalOnMessage) {
          originalOnMessage.call(props.webSocket, event)
        }
      } catch (error) {
        console.error('Error handling WebSocket message:', error)
      }
    }
  }

  // 更新远程 Peer 列表
  const updateRemotePeers = () => {
    if (webrtcManager) {
      remotePeers.value = webrtcManager.getAllPeers()
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
    emit('call-ended')
  }

  // 监听通话激活状态
  watch(
    () => props.isCallActive,
    (active) => {
      if (active && !webrtcManager) {
        initWebRTC()
      } else if (!active && webrtcManager) {
        webrtcManager.destroy()
        webrtcManager = null
      }
    }
  )

  // 定时更新远程 Peer 列表
  let updateInterval: number | null = null

  onMounted(() => {
    if (props.isCallActive) {
      initWebRTC()
    }

    // 每秒更新一次远程 Peer 列表
    updateInterval = window.setInterval(() => {
      updateRemotePeers()
    }, 1000)
  })

  onUnmounted(() => {
    if (webrtcManager) {
      webrtcManager.destroy()
    }
    if (updateInterval) {
      clearInterval(updateInterval)
    }
  })
</script>

<style scoped lang="scss">
  .video-call-container {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: #1a1a1a;
    z-index: 9999;
    display: flex;
    flex-direction: column;
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

    &.grid-5,
    &.grid-6 {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  .video-wrapper {
    position: relative;
    background: #000;
    border-radius: 8px;
    overflow: hidden;
    min-height: 200px;

    video {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    &.local-video {
      video {
        transform: scaleX(-1); // 镜像翻转本地视频
      }
    }
  }

  .video-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 10px;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .user-info {
    .username {
      color: #fff;
      font-size: 14px;
      font-weight: 500;
    }
  }

  .media-status {
    display: flex;
    gap: 8px;

    .status-icon {
      color: #f56c6c;
      font-size: 18px;
    }
  }

  .control-bar {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 20px;
    padding: 20px;
    background: rgba(0, 0, 0, 0.5);
  }

  .call-status {
    position: absolute;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.8);
    color: #fff;
    padding: 12px 24px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 14px;
  }
</style>
