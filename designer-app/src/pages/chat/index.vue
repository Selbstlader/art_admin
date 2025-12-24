<template>
  <!-- 
    AI 对话页面 - 仅 App 端可用
    AI Chat Page - App only with glass morphism design
  -->
  <!-- #ifdef APP-PLUS -->
  <view class="chat-page">
    <!-- 页面头部 - Page header with project context -->
    <view class="chat-header" :style="{ paddingTop: statusBarHeight + 'px' }">
      <view class="header-content">
        <!-- 返回按钮 - Back button -->
        <view class="header-left" @click="handleBack">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6"></path>
          </svg>
        </view>
        
        <!-- 标题区域 - Title area -->
        <view class="header-center">
          <text class="header-title">AI 设计助手</text>
          <text v-if="projectContext" class="header-subtitle">{{ projectContextName }}</text>
        </view>
        
        <!-- 新建会话按钮 - New session button -->
        <view class="header-right" @click="handleNewSession">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5v14"></path>
            <path d="M5 12h14"></path>
          </svg>
        </view>
      </view>
    </view>

    <!-- 消息列表区域 - Message list area -->
    <scroll-view 
      class="message-list"
      scroll-y
      :scroll-into-view="scrollToId"
      :scroll-with-animation="true"
      @scrolltoupper="handleLoadMore"
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="handleRefresh"
    >
      <!-- 加载更多提示 - Load more indicator -->
      <view v-if="hasMoreHistory" class="load-more-tip">
        <text class="load-more-text">上拉加载更多历史消息</text>
      </view>

      <!-- 空状态 - Empty state -->
      <view v-if="messages.length === 0 && !isLoading" class="empty-state">
        <view class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 8V4H8"></path>
            <rect width="16" height="12" x="4" y="8" rx="2"></rect>
            <path d="M2 14h2"></path>
            <path d="M20 14h2"></path>
            <path d="M15 13v2"></path>
            <path d="M9 13v2"></path>
          </svg>
        </view>
        <text class="empty-title">开始与 AI 对话</text>
        <text class="empty-desc">我是您的设计助手，可以帮您解答设计问题、提供建议</text>
        
        <!-- 快捷提问 - Quick questions -->
        <view class="quick-questions">
          <view 
            v-for="(question, index) in quickQuestions" 
            :key="index"
            class="quick-question-item"
            @click="handleQuickQuestion(question)"
          >
            <text class="quick-question-text">{{ question }}</text>
          </view>
        </view>
      </view>

      <!-- 消息列表 - Message list -->
      <view v-else class="messages-container">
        <ChatBubble 
          v-for="(msg, index) in messages" 
          :key="msg.id"
          :id="'msg-' + msg.id"
          :message="msg"
          :isStreaming="isStreaming && index === messages.length - 1 && msg.role === 'assistant'"
        />
      </view>

      <!-- 滚动锚点 - Scroll anchor -->
      <view id="scroll-bottom" class="scroll-anchor"></view>
    </scroll-view>

    <!-- 错误提示 - Error toast -->
    <view v-if="errorMessage" class="error-toast">
      <view class="error-content">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <text class="error-text">{{ errorMessage }}</text>
        <view class="retry-btn" @click="handleRetry">
          <text class="retry-text">重试</text>
        </view>
      </view>
    </view>

    <!-- 输入区域 - Input area with glass morphism -->
    <view class="input-area" :style="{ paddingBottom: safeAreaBottom + 'px' }">
      <view class="input-container">
        <!-- 输入框 - Input field -->
        <view class="input-wrapper">
          <textarea
            v-model="inputText"
            class="message-input"
            placeholder="输入您的问题..."
            placeholder-class="input-placeholder"
            :maxlength="2000"
            :auto-height="true"
            :show-confirm-bar="false"
            :adjust-position="true"
            :cursor-spacing="20"
            @confirm="handleSend"
            @focus="handleInputFocus"
            @blur="handleInputBlur"
          />
        </view>
        
        <!-- 发送按钮 - Send button -->
        <view 
          :class="['send-btn', { 'send-btn-active': canSend, 'send-btn-disabled': !canSend || isStreaming }]"
          @click="handleSend"
        >
          <svg v-if="!isStreaming" xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m22 2-7 20-4-9-9-4Z"></path>
            <path d="M22 2 11 13"></path>
          </svg>
          <!-- 发送中状态 - Sending state -->
          <view v-else class="sending-indicator">
            <view class="sending-dot"></view>
          </view>
        </view>
      </view>
      
      <!-- 输入提示 - Input hint -->
      <view class="input-hint">
        <text class="hint-text">{{ inputText.length }}/2000</text>
      </view>
    </view>
  </view>
  <!-- #endif -->

  <!-- #ifndef APP-PLUS -->
  <!-- 非 App 端显示提示 - Show hint for non-App platforms -->
  <view class="not-supported-page">
    <view class="not-supported-content">
      <view class="not-supported-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
      </view>
      <text class="not-supported-title">功能暂不可用</text>
      <text class="not-supported-desc">AI 对话功能仅在 App 端可用，请下载 App 体验完整功能</text>
    </view>
    
    <!-- 自定义底部导航栏 - Custom TabBar (小程序端不显示此页面，但保留组件) -->
    <CustomTabBar />
  </view>
  <!-- #endif -->
</template>

<script setup lang="ts">
/*** 
 * AI Chat Page Component
 * AI 对话页面 - 仅 App 端可用，玻璃拟态设计
 * App only with glass morphism design and streaming support
 ***/
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useChatStore } from '@/store/chat'
import { platform } from '@/utils/platform'
import { ChatBubble } from '@/components'
import CustomTabBar from '@/components/CustomTabBar.vue'

/*** Store 实例 - Store instance ***/
const chatStore = useChatStore()

/*** 响应式状态 - Reactive state ***/
const inputText = ref('')
const isLoading = ref(false)
const isRefreshing = ref(false)
const hasMoreHistory = ref(false)
const currentPage = ref(1)
const scrollToId = ref('')
const errorMessage = ref('')
const isInputFocused = ref(false)
const projectContextName = ref('')

/*** 平台信息 - Platform info ***/
const statusBarHeight = ref(0)
const safeAreaBottom = ref(0)

/*** 快捷提问列表 - Quick questions list ***/
const quickQuestions = [
  '如何选择合适的办公室装修风格？',
  '工装设计中有哪些常见的规范要求？',
  '如何控制装修成本？'
]

/*** 计算属性 - Computed properties ***/
const messages = computed(() => chatStore.currentMessages)
const isStreaming = computed(() => chatStore.streaming)
const projectContext = computed(() => chatStore.currentProjectContext)
const canSend = computed(() => inputText.value.trim().length > 0 && !isStreaming.value)

/*** 
 * 初始化页面 - Initialize page
 * 获取平台信息和项目上下文
 ***/
onMounted(async () => {
  // 获取系统信息 - Get system info
  statusBarHeight.value = platform.getStatusBarHeight()
  safeAreaBottom.value = platform.getSafeAreaBottom()

  // 获取页面参数 - Get page params
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const options = (currentPage as any)?.options || {}
  
  // 如果有项目 ID 参数，设置项目上下文 - Set project context if projectId exists
  if (options.projectId) {
    const projectId = parseInt(options.projectId)
    chatStore.setProjectContext(projectId)
    projectContextName.value = options.projectName || `项目 #${projectId}`
  }

  // 加载历史消息 - Load history messages
  if (chatStore.currentSessionId) {
    await loadHistoryMessages()
  }

  // 滚动到底部 - Scroll to bottom
  scrollToBottom()
})

/*** 
 * 页面卸载时清理 - Cleanup on unmount
 ***/
onUnmounted(() => {
  // 清除错误消息 - Clear error message
  errorMessage.value = ''
})

/*** 
 * 发送消息 - Send message
 ***/
async function handleSend() {
  if (!canSend.value) return

  const query = inputText.value.trim()
  inputText.value = ''
  errorMessage.value = ''

  try {
    await chatStore.sendMessage(query)
    // 滚动到底部 - Scroll to bottom
    await nextTick()
    scrollToBottom()
  } catch (error: any) {
    // 处理超时错误 - Handle timeout error
    if (error?.message?.includes('timeout') || error?.message?.includes('超时')) {
      errorMessage.value = 'AI 服务响应超时，请稍后重试'
    } else {
      errorMessage.value = error?.message || 'AI 服务暂时不可用，请稍后重试'
    }
  }
}

/*** 
 * 加载历史消息 - Load history messages
 ***/
async function loadHistoryMessages() {
  if (isLoading.value) return
  
  isLoading.value = true
  try {
    hasMoreHistory.value = await chatStore.loadHistory(currentPage.value)
    currentPage.value++
  } finally {
    isLoading.value = false
  }
}

/*** 
 * 上拉加载更多 - Load more on scroll to top
 ***/
async function handleLoadMore() {
  if (!hasMoreHistory.value || isLoading.value) return
  await loadHistoryMessages()
}

/*** 
 * 下拉刷新 - Pull to refresh
 ***/
async function handleRefresh() {
  isRefreshing.value = true
  currentPage.value = 1
  try {
    hasMoreHistory.value = await chatStore.loadHistory(1)
  } finally {
    isRefreshing.value = false
  }
}

/*** 
 * 滚动到底部 - Scroll to bottom
 ***/
function scrollToBottom() {
  scrollToId.value = ''
  nextTick(() => {
    scrollToId.value = 'scroll-bottom'
  })
}

/*** 
 * 快捷提问 - Quick question
 ***/
function handleQuickQuestion(question: string) {
  inputText.value = question
  handleSend()
}

/*** 
 * 重试发送 - Retry send
 ***/
async function handleRetry() {
  errorMessage.value = ''
  await chatStore.retryLastMessage()
  scrollToBottom()
}

/*** 
 * 新建会话 - New session
 ***/
function handleNewSession() {
  uni.showModal({
    title: '新建会话',
    content: '确定要开始新的对话吗？当前对话记录将被保留。',
    confirmText: '确定',
    cancelText: '取消',
    success: (res) => {
      if (res.confirm) {
        chatStore.createNewSession()
        currentPage.value = 1
        hasMoreHistory.value = false
      }
    }
  })
}

/*** 
 * 返回上一页 - Go back
 ***/
function handleBack() {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/index/index' })
    }
  })
}

/*** 
 * 输入框聚焦 - Input focus
 ***/
function handleInputFocus() {
  isInputFocused.value = true
  // 延迟滚动到底部 - Delay scroll to bottom
  setTimeout(() => {
    scrollToBottom()
  }, 300)
}

/*** 
 * 输入框失焦 - Input blur
 ***/
function handleInputBlur() {
  isInputFocused.value = false
}
</script>


<style lang="scss" scoped>
/*** 页面容器 - Page container ***/
.chat-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(180deg, #F0F7FF 0%, #F5F7FA 100%);
}

/*** 页面头部 - Page header with glass morphism ***/
.chat-header {
  position: relative;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1rpx solid rgba(59, 130, 246, 0.1);
  box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.04);
  z-index: 100;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}

.header-left,
.header-right {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20rpx;
  background: rgba(59, 130, 246, 0.08);
  color: #3B82F6;
  transition: all 0.2s ease;
}

.header-left:active,
.header-right:active {
  opacity: 0.7;
  transform: scale(0.95);
}

.header-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.header-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
}

.header-subtitle {
  font-size: 22rpx;
  color: #64748B;
  margin-top: 4rpx;
}

/*** 消息列表区域 - Message list area ***/
.message-list {
  flex: 1;
  padding: 24rpx 0;
  overflow-y: auto;
}

.messages-container {
  padding: 0 8rpx;
}

.scroll-anchor {
  height: 1rpx;
}

/*** 加载更多提示 - Load more tip ***/
.load-more-tip {
  text-align: center;
  padding: 24rpx;
}

.load-more-text {
  font-size: 24rpx;
  color: #94A3B8;
}

/*** 空状态 - Empty state ***/
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 48rpx;
}

.empty-icon {
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  border-radius: 40rpx;
  color: #3B82F6;
  margin-bottom: 32rpx;
  box-shadow: 0 8rpx 32rpx rgba(59, 130, 246, 0.15);
}

.empty-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 16rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: #64748B;
  text-align: center;
  line-height: 1.6;
  margin-bottom: 48rpx;
}

/*** 快捷提问 - Quick questions ***/
.quick-questions {
  width: 100%;
}

.quick-question-item {
  background: rgba(255, 255, 255, 0.9);
  border: 2rpx solid rgba(59, 130, 246, 0.15);
  border-radius: 20rpx;
  padding: 24rpx 28rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}

.quick-question-item:active {
  background: rgba(59, 130, 246, 0.08);
  transform: scale(0.98);
}

.quick-question-text {
  font-size: 26rpx;
  color: #3B82F6;
  line-height: 1.5;
}

/*** 错误提示 - Error toast ***/
.error-toast {
  position: fixed;
  left: 24rpx;
  right: 24rpx;
  bottom: 200rpx;
  z-index: 200;
}

.error-content {
  display: flex;
  align-items: center;
  background: rgba(239, 68, 68, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 20rpx;
  padding: 20rpx 24rpx;
  color: #FFFFFF;
  box-shadow: 0 8rpx 32rpx rgba(239, 68, 68, 0.3);
}

.error-text {
  flex: 1;
  font-size: 26rpx;
  margin-left: 12rpx;
}

.retry-btn {
  padding: 8rpx 20rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12rpx;
  margin-left: 16rpx;
}

.retry-text {
  font-size: 24rpx;
  color: #FFFFFF;
}

/*** 输入区域 - Input area with glass morphism ***/
.input-area {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1rpx solid rgba(59, 130, 246, 0.1);
  padding: 16rpx 24rpx;
  box-shadow: 0 -4rpx 24rpx rgba(0, 0, 0, 0.04);
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 16rpx;
}

.input-wrapper {
  flex: 1;
  background: #F8FAFC;
  border: 2rpx solid rgba(59, 130, 246, 0.15);
  border-radius: 24rpx;
  padding: 16rpx 24rpx;
  transition: all 0.2s ease;
}

.input-wrapper:focus-within {
  border-color: #3B82F6;
  box-shadow: 0 0 0 4rpx rgba(59, 130, 246, 0.1);
}

.message-input {
  width: 100%;
  min-height: 40rpx;
  max-height: 200rpx;
  font-size: 28rpx;
  color: #1E293B;
  line-height: 1.5;
}

.input-placeholder {
  color: #94A3B8;
}

/*** 发送按钮 - Send button ***/
.send-btn {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 24rpx;
  transition: all 0.2s ease;
}

.send-btn-active {
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  color: #FFFFFF;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.35);
}

.send-btn-active:active {
  transform: scale(0.95);
  opacity: 0.9;
}

.send-btn-disabled {
  background: #E2E8F0;
  color: #94A3B8;
}

/*** 发送中指示器 - Sending indicator ***/
.sending-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
}

.sending-dot {
  width: 16rpx;
  height: 16rpx;
  background: #FFFFFF;
  border-radius: 50%;
  animation: sending-pulse 1s ease-in-out infinite;
}

@keyframes sending-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(0.6);
    opacity: 0.5;
  }
}

/*** 输入提示 - Input hint ***/
.input-hint {
  display: flex;
  justify-content: flex-end;
  padding: 8rpx 8rpx 0;
}

.hint-text {
  font-size: 20rpx;
  color: #94A3B8;
}

/*** 非 App 端提示页面 - Not supported page ***/
.not-supported-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F5F7FA;
  padding: 48rpx;
}

.not-supported-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.not-supported-icon {
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #FEE2E2 0%, #FECACA 100%);
  border-radius: 40rpx;
  color: #EF4444;
  margin-bottom: 32rpx;
}

.not-supported-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 16rpx;
}

.not-supported-desc {
  font-size: 26rpx;
  color: #64748B;
  line-height: 1.6;
}
</style>
