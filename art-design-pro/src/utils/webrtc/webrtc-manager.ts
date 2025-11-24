/**
 * WebRTC 管理器
 * 整合 PeerManager 和 WebSocket 信令
 */

import { PeerManager } from './peer-manager'
import type {
  SignalMessage,
  OfferMessage,
  AnswerMessage,
  IceCandidateMessage,
  CallRequestMessage,
  CallAcceptMessage,
  MediaToggleMessage,
  WebRTCConfig
} from './types'
import { SignalType, CallStatus, DEFAULT_MEDIA_CONSTRAINTS } from './types'

export interface WebRTCManagerOptions {
  roomId: number
  userId: number
  username: string
  webSocket: WebSocket
  config?: WebRTCConfig
  onCallStatusChange?: (status: CallStatus) => void
  onRemoteStream?: (userId: number, username: string, stream: MediaStream) => void
  onPeerLeft?: (userId: number) => void
  onError?: (error: Error) => void
}

export class WebRTCManager {
  private peerManager: PeerManager
  private roomId: number
  private userId: number
  private username: string
  private webSocket: WebSocket
  private callStatus: CallStatus = CallStatus.IDLE
  private localStream: MediaStream | null = null
  private audioEnabled = true
  private videoEnabled = true

  // 回调函数
  private onCallStatusChange?: (status: CallStatus) => void
  private onRemoteStream?: (userId: number, username: string, stream: MediaStream) => void
  private onPeerLeft?: (userId: number) => void
  private onError?: (error: Error) => void

  // 用户信息映射
  private userMap: Map<number, string> = new Map()

  constructor(options: WebRTCManagerOptions) {
    this.roomId = options.roomId
    this.userId = options.userId
    this.username = options.username
    this.webSocket = options.webSocket
    this.onCallStatusChange = options.onCallStatusChange
    this.onRemoteStream = options.onRemoteStream
    this.onPeerLeft = options.onPeerLeft
    this.onError = options.onError

    // 创建 PeerManager
    this.peerManager = new PeerManager(options.config)

    // 设置回调
    this.peerManager.onRemoteStreamCallback((userId, stream) => {
      const username = this.userMap.get(userId) || `User ${userId}`
      if (this.onRemoteStream) {
        this.onRemoteStream(userId, username, stream)
      }
    })

    this.peerManager.onPeerDisconnectedCallback((userId) => {
      if (this.onPeerLeft) {
        this.onPeerLeft(userId)
      }
    })

    this.peerManager.setSendSignalCallback((message) => {
      this.sendSignal({ ...message, from: this.userId })
    })

    // 监听 WebSocket 消息
    this.setupWebSocketListeners()
  }

  /**
   * 设置 WebSocket 监听器
   */
  private setupWebSocketListeners() {
    // 注意: 这里假设 WebSocket 消息会被外部处理并调用 handleSignalMessage
    // 实际使用时需要在外部调用 handleSignalMessage
  }

  /**
   * 开始通话 - 获取本地媒体流
   */
  async startCall(constraints?: MediaStreamConstraints): Promise<void> {
    try {
      // 获取本地媒体流
      const stream = await navigator.mediaDevices.getUserMedia(
        constraints || DEFAULT_MEDIA_CONSTRAINTS
      )

      this.localStream = stream
      this.peerManager.setLocalStream(stream)

      this.updateCallStatus(CallStatus.CALLING)

      console.log('Local stream obtained:', stream)
    } catch (error) {
      console.error('Error getting local stream:', error)
      if (this.onError) {
        this.onError(error as Error)
      }
      throw error
    }
  }

  /**
   * 连接到指定用户
   */
  async connectToPeer(userId: number, username: string): Promise<void> {
    try {
      this.userMap.set(userId, username)
      await this.peerManager.createPeerConnection(userId, username, this.roomId, true)
      console.log(`Connecting to peer ${userId} (${username})`)
    } catch (error) {
      console.error('Error connecting to peer:', error)
      if (this.onError) {
        this.onError(error as Error)
      }
    }
  }

  /**
   * 处理信令消息
   */
  async handleSignalMessage(message: any): Promise<void> {
    console.log('Received signal:', message.type, 'from:', message.from)

    try {
      // 处理房间用户列表（新用户加入时收到）
      if (message.type === 'room_users') {
        // 等待本地流准备好再处理
        if (!this.localStream) {
          console.log('Local stream not ready, waiting...')
          await new Promise((resolve) => {
            const checkInterval = setInterval(() => {
              if (this.localStream) {
                clearInterval(checkInterval)
                resolve(true)
              }
            }, 100)
            // 最多等待 5 秒
            setTimeout(() => {
              clearInterval(checkInterval)
              resolve(false)
            }, 5000)
          })
        }
        if (this.localStream) {
          await this.handleRoomUsers(message.users)
        } else {
          console.error('Local stream not available after waiting')
        }
        return
      }

      // 忽略自己发送的消息
      if (message.from === this.userId) {
        return
      }

      switch (message.type) {
        case SignalType.OFFER:
          await this.handleOffer(message as OfferMessage)
          break

        case SignalType.ANSWER:
          await this.handleAnswer(message as AnswerMessage)
          break

        case SignalType.ICE_CANDIDATE:
          await this.handleIceCandidate(message as IceCandidateMessage)
          break

        case SignalType.CALL_REQUEST:
          this.handleCallRequest(message as CallRequestMessage)
          break

        case SignalType.CALL_ACCEPT:
          this.handleCallAccept(message as CallAcceptMessage)
          break

        case SignalType.CALL_HANGUP:
          this.handleHangup(message.from)
          break

        case SignalType.TOGGLE_AUDIO:
        case SignalType.TOGGLE_VIDEO:
          this.handleMediaToggle(message as MediaToggleMessage)
          break

        case SignalType.USER_JOINED:
          // 新用户加入，等待本地流准备好再建立连接
          if (this.localStream) {
            await this.handleUserJoined(message)
          } else {
            console.log('Local stream not ready for new user, skipping connection')
          }
          break

        case SignalType.USER_LEFT:
          this.handleUserLeft(message.from)
          break
      }
    } catch (error) {
      console.error('Error handling signal message:', error)
      if (this.onError) {
        this.onError(error as Error)
      }
    }
  }

  /**
   * 处理房间用户列表
   */
  private async handleRoomUsers(users: Array<{ userId: number; username: string }>) {
    console.log('Room users:', users)
    // 向房间内已有的每个用户发起连接
    // 使用 userId 比较来避免信令冲突：只有当自己的 userId 较小时才主动发起连接
    for (const user of users) {
      if (user.userId !== this.userId && this.userId < user.userId) {
        console.log(`Initiating connection to user ${user.userId} (I am ${this.userId})`)
        this.userMap.set(user.userId, user.username)
        await this.connectToPeer(user.userId, user.username)
      } else if (user.userId !== this.userId) {
        console.log(`Waiting for user ${user.userId} to initiate connection (I am ${this.userId})`)
        this.userMap.set(user.userId, user.username)
      }
    }
  }

  /**
   * 处理新用户加入
   */
  private async handleUserJoined(message: any) {
    const userId = message.from || message.data?.userId
    const username = message.data?.username || `User ${userId}`

    if (userId && userId !== this.userId) {
      console.log('New user joined:', userId, username)
      this.userMap.set(userId, username)
      // 使用 userId 比较来避免信令冲突：只有当自己的 userId 较小时才主动发起连接
      if (this.userId < userId) {
        console.log(`Initiating connection to new user ${userId} (I am ${this.userId})`)
        await this.connectToPeer(userId, username)
      } else {
        console.log(`Waiting for new user ${userId} to initiate connection (I am ${this.userId})`)
      }
    }
  }

  /**
   * 处理 Offer
   */
  private async handleOffer(message: OfferMessage) {
    const username = this.userMap.get(message.from) || `User ${message.from}`
    this.userMap.set(message.from, username)

    await this.peerManager.handleOffer(message.from, username, this.roomId, message.sdp)

    if (this.callStatus === CallStatus.IDLE) {
      this.updateCallStatus(CallStatus.CONNECTED)
    }
  }

  /**
   * 处理 Answer
   */
  private async handleAnswer(message: AnswerMessage) {
    await this.peerManager.handleAnswer(message.from, message.sdp)
    this.updateCallStatus(CallStatus.CONNECTED)
  }

  /**
   * 处理 ICE Candidate
   */
  private async handleIceCandidate(message: IceCandidateMessage) {
    await this.peerManager.handleIceCandidate(message.from, message.candidate)
  }

  /**
   * 处理通话请求
   */
  private handleCallRequest(message: CallRequestMessage) {
    this.userMap.set(message.from, message.username)
    this.updateCallStatus(CallStatus.RINGING)
  }

  /**
   * 处理通话接受
   */
  private handleCallAccept(message: CallAcceptMessage) {
    this.userMap.set(message.from, message.username)
    this.updateCallStatus(CallStatus.CONNECTED)
  }

  /**
   * 处理挂断
   */
  private handleHangup(userId: number) {
    this.peerManager.closePeerConnection(userId)
    if (this.onPeerLeft) {
      this.onPeerLeft(userId)
    }
  }

  /**
   * 处理用户离开
   */
  private handleUserLeft(userId: number) {
    this.peerManager.closePeerConnection(userId)
    this.userMap.delete(userId)
    if (this.onPeerLeft) {
      this.onPeerLeft(userId)
    }
  }

  /**
   * 处理媒体切换
   */
  private handleMediaToggle(message: MediaToggleMessage) {
    const type = message.type === SignalType.TOGGLE_AUDIO ? 'audio' : 'video'
    this.peerManager.updatePeerMediaStatus(message.from, type, message.enabled)
  }

  /**
   * 切换音频
   */
  toggleAudio(): boolean {
    this.audioEnabled = !this.audioEnabled
    this.peerManager.toggleAudio(this.audioEnabled)

    // 通知其他用户
    this.sendSignal({
      type: SignalType.TOGGLE_AUDIO,
      from: this.userId,
      roomId: this.roomId,
      timestamp: Date.now(),
      enabled: this.audioEnabled
    } as MediaToggleMessage)

    return this.audioEnabled
  }

  /**
   * 切换视频
   */
  toggleVideo(): boolean {
    this.videoEnabled = !this.videoEnabled
    this.peerManager.toggleVideo(this.videoEnabled)

    // 通知其他用户
    this.sendSignal({
      type: SignalType.TOGGLE_VIDEO,
      from: this.userId,
      roomId: this.roomId,
      timestamp: Date.now(),
      enabled: this.videoEnabled
    } as MediaToggleMessage)

    return this.videoEnabled
  }

  /**
   * 挂断通话
   */
  hangup(): void {
    // 通知其他用户
    this.sendSignal({
      type: SignalType.CALL_HANGUP,
      from: this.userId,
      roomId: this.roomId,
      timestamp: Date.now()
    })

    this.endCall()
  }

  /**
   * 结束通话
   */
  endCall(): void {
    // 关闭所有 Peer 连接
    this.peerManager.closeAll()

    // 停止本地流
    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => track.stop())
      this.localStream = null
    }

    this.updateCallStatus(CallStatus.ENDED)
    this.userMap.clear()
  }

  /**
   * 发送信令消息
   */
  private sendSignal(message: SignalMessage): void {
    if (this.webSocket.readyState === WebSocket.OPEN) {
      this.webSocket.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not open, cannot send signal')
    }
  }

  /**
   * 更新通话状态
   */
  private updateCallStatus(status: CallStatus): void {
    this.callStatus = status
    if (this.onCallStatusChange) {
      this.onCallStatusChange(status)
    }
  }

  /**
   * 获取本地流
   */
  getLocalStream(): MediaStream | null {
    return this.localStream
  }

  /**
   * 获取所有 Peer
   */
  getAllPeers() {
    return this.peerManager.getAllPeers()
  }

  /**
   * 获取通话状态
   */
  getCallStatus(): CallStatus {
    return this.callStatus
  }

  /**
   * 获取音频状态
   */
  isAudioEnabled(): boolean {
    return this.audioEnabled
  }

  /**
   * 获取视频状态
   */
  isVideoEnabled(): boolean {
    return this.videoEnabled
  }

  /**
   * 销毁管理器
   */
  destroy(): void {
    this.endCall()
  }
}
