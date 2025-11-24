/**
 * WebRTC 相关类型定义
 */

// 信令消息类型
export enum SignalType {
  // WebRTC 信令
  OFFER = 'webrtc_offer',
  ANSWER = 'webrtc_answer',
  ICE_CANDIDATE = 'webrtc_ice_candidate',

  // 通话控制
  CALL_REQUEST = 'call_request', // 发起通话请求
  CALL_ACCEPT = 'call_accept', // 接受通话
  CALL_REJECT = 'call_reject', // 拒绝通话
  CALL_CANCEL = 'call_cancel', // 取消通话
  CALL_HANGUP = 'call_hangup', // 挂断通话

  // 用户状态
  USER_JOINED = 'user_joined', // 用户加入
  USER_LEFT = 'user_left', // 用户离开

  // 媒体控制
  TOGGLE_AUDIO = 'toggle_audio', // 切换音频
  TOGGLE_VIDEO = 'toggle_video' // 切换视频
}

// 信令消息基础接口
export interface SignalMessage {
  type: SignalType
  from: number // 发送者用户ID
  to?: number // 接收者用户ID (可选,用于点对点)
  roomId: number
  timestamp: number
}

// WebRTC Offer 消息
export interface OfferMessage extends SignalMessage {
  type: SignalType.OFFER
  sdp: RTCSessionDescriptionInit
}

// WebRTC Answer 消息
export interface AnswerMessage extends SignalMessage {
  type: SignalType.ANSWER
  sdp: RTCSessionDescriptionInit
}

// ICE Candidate 消息
export interface IceCandidateMessage extends SignalMessage {
  type: SignalType.ICE_CANDIDATE
  candidate: RTCIceCandidateInit
}

// 通话请求消息
export interface CallRequestMessage extends SignalMessage {
  type: SignalType.CALL_REQUEST
  username: string
}

// 通话接受消息
export interface CallAcceptMessage extends SignalMessage {
  type: SignalType.CALL_ACCEPT
  username: string
}

// 通话拒绝消息
export interface CallRejectMessage extends SignalMessage {
  type: SignalType.CALL_REJECT
  reason?: string
}

// 媒体状态消息
export interface MediaToggleMessage extends SignalMessage {
  type: SignalType.TOGGLE_AUDIO | SignalType.TOGGLE_VIDEO
  enabled: boolean
}

// Peer 连接信息
export interface PeerConnection {
  userId: number
  username: string
  connection: RTCPeerConnection
  stream?: MediaStream
  audioEnabled: boolean
  videoEnabled: boolean
}

// 通话状态
export enum CallStatus {
  IDLE = 'idle', // 空闲
  CALLING = 'calling', // 呼叫中
  RINGING = 'ringing', // 响铃中
  CONNECTED = 'connected', // 已连接
  ENDED = 'ended' // 已结束
}

// WebRTC 配置
export interface WebRTCConfig {
  iceServers: RTCIceServer[]
  mediaConstraints?: MediaStreamConstraints
}

// 默认 STUN 服务器配置 (Google 公共 STUN)
export const DEFAULT_ICE_SERVERS: RTCIceServer[] = [
  { urls: 'stun:stun.l.google.com:19302' },
  { urls: 'stun:stun1.l.google.com:19302' },
  { urls: 'stun:stun2.l.google.com:19302' }
]

// 默认媒体约束
export const DEFAULT_MEDIA_CONSTRAINTS: MediaStreamConstraints = {
  audio: {
    echoCancellation: true, // 回声消除
    noiseSuppression: true, // 噪音抑制
    autoGainControl: true // 自动增益控制
  },
  video: {
    width: { ideal: 1280 },
    height: { ideal: 720 },
    frameRate: { ideal: 30 }
  }
}
