/**
 * WebRTC Peer 连接管理器
 * 负责管理多个 P2P 连接
 */

import type {
  PeerConnection,
  WebRTCConfig,
  SignalMessage,
  OfferMessage,
  AnswerMessage,
  IceCandidateMessage
} from './types'
import { SignalType, DEFAULT_ICE_SERVERS } from './types'

export class PeerManager {
  private peers: Map<number, PeerConnection> = new Map()
  private localStream: MediaStream | null = null
  private config: RTCConfiguration
  private onRemoteStream?: (userId: number, stream: MediaStream) => void
  private onPeerDisconnected?: (userId: number) => void
  private sendSignal?: (message: SignalMessage) => void

  constructor(config?: WebRTCConfig) {
    this.config = {
      iceServers: config?.iceServers || DEFAULT_ICE_SERVERS
    }
  }

  /**
   * 设置本地媒体流
   */
  setLocalStream(stream: MediaStream) {
    this.localStream = stream
    // 为已存在的连接添加轨道
    this.peers.forEach((peer) => {
      stream.getTracks().forEach((track) => {
        peer.connection.addTrack(track, stream)
      })
    })
  }

  /**
   * 获取本地媒体流
   */
  getLocalStream(): MediaStream | null {
    return this.localStream
  }

  /**
   * 设置远程流回调
   */
  onRemoteStreamCallback(callback: (userId: number, stream: MediaStream) => void) {
    this.onRemoteStream = callback
  }

  /**
   * 设置 Peer 断开回调
   */
  onPeerDisconnectedCallback(callback: (userId: number) => void) {
    this.onPeerDisconnected = callback
  }

  /**
   * 设置信令发送回调
   */
  setSendSignalCallback(callback: (message: SignalMessage) => void) {
    this.sendSignal = callback
  }

  /**
   * 创建 Peer 连接
   */
  async createPeerConnection(
    userId: number,
    username: string,
    roomId: number,
    isInitiator: boolean
  ): Promise<RTCPeerConnection> {
    // 如果已存在连接,先关闭
    if (this.peers.has(userId)) {
      this.closePeerConnection(userId)
    }

    const connection = new RTCPeerConnection(this.config)

    // 添加本地流
    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => {
        connection.addTrack(track, this.localStream!)
      })
    }

    // ICE 候选事件
    connection.onicecandidate = (event) => {
      if (event.candidate && this.sendSignal) {
        console.log(`Sending ICE candidate to user ${userId}:`, event.candidate.candidate)
        const message: IceCandidateMessage = {
          type: SignalType.ICE_CANDIDATE,
          from: 0, // 会在发送时填充
          to: userId,
          roomId,
          timestamp: Date.now(),
          candidate: event.candidate.toJSON()
        }
        this.sendSignal(message)
      } else if (!event.candidate) {
        console.log(`ICE gathering complete for user ${userId}`)
      }
    }

    // ICE 连接状态变化
    connection.oniceconnectionstatechange = () => {
      console.log(`Peer ${userId} ICE connection state:`, connection.iceConnectionState)
    }

    // 接收远程流
    connection.ontrack = (event) => {
      console.log(
        'Received remote track from user:',
        userId,
        'Track:',
        event.track.kind,
        'Streams:',
        event.streams.length
      )
      const [stream] = event.streams
      if (stream) {
        console.log('Remote stream details:', {
          id: stream.id,
          active: stream.active,
          tracks: stream
            .getTracks()
            .map((t) => ({ kind: t.kind, enabled: t.enabled, readyState: t.readyState }))
        })
        if (this.onRemoteStream) {
          this.onRemoteStream(userId, stream)
        }

        // 更新 peer 的流
        const peer = this.peers.get(userId)
        if (peer) {
          peer.stream = stream
        }
      } else {
        console.error('No stream in track event for user:', userId)
      }
    }

    // 连接状态变化
    connection.onconnectionstatechange = () => {
      console.log(`Peer ${userId} connection state:`, connection.connectionState)
      if (
        connection.connectionState === 'disconnected' ||
        connection.connectionState === 'failed' ||
        connection.connectionState === 'closed'
      ) {
        this.handlePeerDisconnected(userId)
      }
    }

    // ICE 连接状态变化
    connection.oniceconnectionstatechange = () => {
      console.log(`Peer ${userId} ICE state:`, connection.iceConnectionState)
    }

    // 保存 peer 信息
    this.peers.set(userId, {
      userId,
      username,
      connection,
      audioEnabled: true,
      videoEnabled: true
    })

    // 如果是发起方,创建 offer
    if (isInitiator) {
      await this.createOffer(userId, roomId)
    }

    return connection
  }

  /**
   * 创建 Offer
   */
  private async createOffer(userId: number, roomId: number) {
    const peer = this.peers.get(userId)
    if (!peer) {
      console.error(`Cannot create offer: peer ${userId} not found`)
      return
    }

    try {
      console.log(`Creating offer for user ${userId}`)
      const offer = await peer.connection.createOffer()
      console.log(`Setting local description for user ${userId}`, offer)
      await peer.connection.setLocalDescription(offer)
      console.log(`Local description set, ICE gathering state: ${peer.connection.iceGatheringState}`)

      if (this.sendSignal) {
        const message: OfferMessage = {
          type: SignalType.OFFER,
          from: 0, // 会在发送时填充
          to: userId,
          roomId,
          timestamp: Date.now(),
          sdp: offer
        }
        console.log(`Sending offer to user ${userId}`)
        this.sendSignal(message)
      }
    } catch (error) {
      console.error('Error creating offer:', error)
    }
  }

  /**
   * 处理接收到的 Offer
   */
  async handleOffer(
    userId: number,
    username: string,
    roomId: number,
    sdp: RTCSessionDescriptionInit
  ) {
    let peer = this.peers.get(userId)

    // 如果连接不存在,创建新连接
    if (!peer) {
      await this.createPeerConnection(userId, username, roomId, false)
      peer = this.peers.get(userId)
    }

    if (!peer) return

    try {
      await peer.connection.setRemoteDescription(new RTCSessionDescription(sdp))
      const answer = await peer.connection.createAnswer()
      await peer.connection.setLocalDescription(answer)

      if (this.sendSignal) {
        const message: AnswerMessage = {
          type: SignalType.ANSWER,
          from: 0, // 会在发送时填充
          to: userId,
          roomId,
          timestamp: Date.now(),
          sdp: answer
        }
        this.sendSignal(message)
      }
    } catch (error) {
      console.error('Error handling offer:', error)
    }
  }

  /**
   * 处理接收到的 Answer
   */
  async handleAnswer(userId: number, sdp: RTCSessionDescriptionInit) {
    const peer = this.peers.get(userId)
    if (!peer) {
      console.error(`Cannot handle answer: peer ${userId} not found`)
      return
    }

    try {
      console.log(`Handling answer from user ${userId}`, sdp)
      console.log(`Connection state before setting remote description:`, {
        signalingState: peer.connection.signalingState,
        iceGatheringState: peer.connection.iceGatheringState,
        iceConnectionState: peer.connection.iceConnectionState
      })
      await peer.connection.setRemoteDescription(new RTCSessionDescription(sdp))
      console.log(`Remote description set, ICE gathering state: ${peer.connection.iceGatheringState}`)
      console.log(`ICE connection state: ${peer.connection.iceConnectionState}`)
    } catch (error) {
      console.error('Error handling answer:', error)
    }
  }

  /**
   * 处理接收到的 ICE Candidate
   */
  async handleIceCandidate(userId: number, candidate: RTCIceCandidateInit) {
    const peer = this.peers.get(userId)
    if (!peer) return

    try {
      await peer.connection.addIceCandidate(new RTCIceCandidate(candidate))
    } catch (error) {
      console.error('Error adding ICE candidate:', error)
    }
  }

  /**
   * 关闭 Peer 连接
   */
  closePeerConnection(userId: number) {
    const peer = this.peers.get(userId)
    if (peer) {
      peer.connection.close()
      this.peers.delete(userId)
      console.log(`Closed peer connection for user ${userId}`)
    }
  }

  /**
   * 处理 Peer 断开
   */
  private handlePeerDisconnected(userId: number) {
    this.closePeerConnection(userId)
    if (this.onPeerDisconnected) {
      this.onPeerDisconnected(userId)
    }
  }

  /**
   * 切换音频
   */
  toggleAudio(enabled: boolean) {
    if (this.localStream) {
      this.localStream.getAudioTracks().forEach((track) => {
        track.enabled = enabled
      })
    }
  }

  /**
   * 切换视频
   */
  toggleVideo(enabled: boolean) {
    if (this.localStream) {
      this.localStream.getVideoTracks().forEach((track) => {
        track.enabled = enabled
      })
    }
  }

  /**
   * 更新远程 Peer 的媒体状态
   */
  updatePeerMediaStatus(userId: number, type: 'audio' | 'video', enabled: boolean) {
    const peer = this.peers.get(userId)
    if (peer) {
      if (type === 'audio') {
        peer.audioEnabled = enabled
      } else {
        peer.videoEnabled = enabled
      }
    }
  }

  /**
   * 获取所有 Peer
   */
  getAllPeers(): PeerConnection[] {
    return Array.from(this.peers.values())
  }

  /**
   * 获取指定 Peer
   */
  getPeer(userId: number): PeerConnection | undefined {
    return this.peers.get(userId)
  }

  /**
   * 关闭所有连接
   */
  closeAll() {
    this.peers.forEach((peer) => {
      peer.connection.close()
    })
    this.peers.clear()

    // 停止本地流
    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => track.stop())
      this.localStream = null
    }
  }
}
