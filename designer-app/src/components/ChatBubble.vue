<template>
  <!-- 
    对话气泡组件 - 玻璃拟态设计
    Chat Bubble Component - Glass morphism design with streaming animation
  -->
  <view :class="['chat-bubble', message.role === 'user' ? 'user-bubble' : 'assistant-bubble']">
    <!-- 头像区域 - Avatar area -->
    <view class="avatar-container">
      <view :class="['avatar', message.role === 'user' ? 'user-avatar' : 'ai-avatar']">
        <!-- 用户头像 - User avatar icon -->
        <svg v-if="message.role === 'user'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
        <!-- AI 头像 - AI avatar icon -->
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 8V4H8"></path>
          <rect width="16" height="12" x="4" y="8" rx="2"></rect>
          <path d="M2 14h2"></path>
          <path d="M20 14h2"></path>
          <path d="M15 13v2"></path>
          <path d="M9 13v2"></path>
        </svg>
      </view>
    </view>

    <!-- 消息内容区域 - Message content area -->
    <view class="bubble-content-wrapper">
      <!-- 消息气泡 - Message bubble with glass morphism -->
      <view :class="['bubble-content', message.role === 'user' ? 'user-content' : 'ai-content']">
        <!-- 消息文本 - Message text -->
        <text class="message-text">{{ message.content }}</text>
        
        <!-- 流式输出光标动画 - Streaming cursor animation -->
        <view v-if="isStreaming && message.role === 'assistant'" class="streaming-cursor">
          <view class="cursor-dot"></view>
        </view>
      </view>

      <!-- 消息时间 - Message timestamp -->
      <view class="message-time">
        <text class="time-text">{{ formattedTime }}</text>
      </view>

      <!-- 玻璃拟态阴影效果 - Glass morphism shadow effect -->
      <view :class="['bubble-shadow', message.role === 'user' ? 'user-shadow' : 'ai-shadow']"></view>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * ChatBubble Component
 * 对话气泡组件 - 玻璃拟态设计，支持流式输出动画
 * Glass morphism design with streaming output animation
 ***/
import { computed } from 'vue'
import type { ChatMessage } from '@/types/api'

/*** Component Props - 组件属性定义 ***/
interface Props {
  message: ChatMessage
  isStreaming?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isStreaming: false
})

/*** 格式化时间 - Format timestamp ***/
const formattedTime = computed((): string => {
  if (!props.message.createdAt) return ''
  
  try {
    const date = new Date(props.message.createdAt)
    const now = new Date()
    const isToday = date.toDateString() === now.toDateString()
    
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    const timeStr = `${hours}:${minutes}`
    
    if (isToday) {
      return timeStr
    } else {
      const month = (date.getMonth() + 1).toString().padStart(2, '0')
      const day = date.getDate().toString().padStart(2, '0')
      return `${month}-${day} ${timeStr}`
    }
  } catch {
    return ''
  }
})
</script>

<style lang="scss" scoped>
/*** 对话气泡容器 - Chat bubble container ***/
.chat-bubble {
  display: flex;
  padding: 16rpx 24rpx;
  margin-bottom: 24rpx;
}

/*** 用户消息 - 右对齐 - User message - right aligned ***/
.user-bubble {
  flex-direction: row-reverse;
}

/*** AI 消息 - 左对齐 - AI message - left aligned ***/
.assistant-bubble {
  flex-direction: row;
}

/*** 头像容器 - Avatar container ***/
.avatar-container {
  flex-shrink: 0;
}

.avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
}

/*** 用户头像样式 - User avatar style ***/
.user-avatar {
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  color: #FFFFFF;
  margin-left: 16rpx;
}

/*** AI 头像样式 - AI avatar style ***/
.ai-avatar {
  background: linear-gradient(135deg, #FFFFFF 0%, #F8FAFC 100%);
  color: #3B82F6;
  border: 2rpx solid rgba(59, 130, 246, 0.2);
  margin-right: 16rpx;
}

/*** 消息内容包装器 - Bubble content wrapper ***/
.bubble-content-wrapper {
  position: relative;
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

/*** 消息气泡内容 - Bubble content with glass morphism ***/
.bubble-content {
  position: relative;
  padding: 24rpx 28rpx;
  border-radius: 24rpx;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

/*** 用户消息样式 - 蓝色渐变 - User message style - blue gradient ***/
.user-content {
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-top-right-radius: 8rpx;
  box-shadow: 
    0 8rpx 24rpx rgba(59, 130, 246, 0.25),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.2);
}

.user-content .message-text {
  color: #FFFFFF;
}

/*** AI 消息样式 - 玻璃拟态白色 - AI message style - glass morphism white ***/
.ai-content {
  background: rgba(255, 255, 255, 0.95);
  border: 2rpx solid rgba(59, 130, 246, 0.1);
  border-top-left-radius: 8rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(0, 0, 0, 0.06),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
}

.ai-content .message-text {
  color: #1E293B;
}

/*** 消息文本 - Message text ***/
.message-text {
  font-size: 28rpx;
  line-height: 1.6;
  word-break: break-word;
  white-space: pre-wrap;
}

/*** 流式输出光标动画 - Streaming cursor animation ***/
.streaming-cursor {
  display: inline-flex;
  align-items: center;
  margin-left: 8rpx;
  vertical-align: middle;
}

.cursor-dot {
  width: 12rpx;
  height: 12rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-radius: 50%;
  animation: cursor-blink 1s ease-in-out infinite;
}

@keyframes cursor-blink {
  0%, 50% {
    opacity: 1;
    transform: scale(1);
  }
  25% {
    opacity: 0.5;
    transform: scale(0.8);
  }
  75% {
    opacity: 0.5;
    transform: scale(0.8);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/*** 消息时间 - Message timestamp ***/
.message-time {
  margin-top: 8rpx;
  padding: 0 8rpx;
}

.user-bubble .message-time {
  text-align: right;
}

.assistant-bubble .message-time {
  text-align: left;
}

.time-text {
  font-size: 20rpx;
  color: #94A3B8;
}

/*** 气泡阴影效果 - Bubble shadow effect ***/
.bubble-shadow {
  position: absolute;
  bottom: -8rpx;
  height: 16rpx;
  border-radius: 50%;
  z-index: -1;
}

.user-shadow {
  right: 16rpx;
  left: 32rpx;
  background: radial-gradient(ellipse, rgba(59, 130, 246, 0.15) 0%, transparent 70%);
}

.ai-shadow {
  left: 16rpx;
  right: 32rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.06) 0%, transparent 70%);
}
</style>
