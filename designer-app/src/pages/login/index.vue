<template>
  <view class="login-page">
    <!-- 背景渐变 - Background gradient -->
    <view class="bg-gradient"></view>
    
    <!-- 登录卡片 - Login card -->
    <view class="login-card glass-card">
      <!-- Logo 区域 - Logo area -->
      <view class="logo-area">
        <view class="logo-icon">
          <text class="logo-text">设计</text>
        </view>
        <text class="app-name">工装设计助手</text>
        <text class="app-desc">AI赋能设计，让创意更高效</text>
      </view>

      <!-- 表单区域 - Form area -->
      <view class="form-area">
        <view class="input-group">
          <view class="input-wrapper">
            <text class="input-icon">👤</text>
            <input
              class="form-input"
              type="text"
              placeholder="请输入账号"
              v-model="formData.username"
              @confirm="handleLogin"
            />
          </view>
        </view>

        <view class="input-group">
          <view class="input-wrapper">
            <text class="input-icon">🔒</text>
            <input
              class="form-input"
              type="password"
              placeholder="请输入密码"
              v-model="formData.password"
              @confirm="handleLogin"
            />
          </view>
        </view>

        <!-- 登录按钮 - Login button -->
        <button
          class="login-btn"
          :class="{ 'btn-loading': loading }"
          :disabled="loading || !isFormValid"
          @click="handleLogin"
        >
          <text v-if="loading" class="loading-icon">⏳</text>
          <text class="btn-text">{{ loading ? '登录中...' : '登录' }}</text>
        </button>
      </view>

      <!-- 底部提示 - Bottom tip -->
      <view class="bottom-tip">
        <text class="tip-text">使用与 Web 端相同的账号登录</text>
      </view>
    </view>

    <!-- 版本信息 - Version info -->
    <view class="version-info">
      <text class="version-text">v1.0.0</text>
    </view>
  </view>
</template>

<script setup lang="ts">
/*** Login page - user authentication with account and password ***/
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'

// 表单数据 - Form data
const formData = ref({
  username: '',
  password: ''
})

const loading = ref(false)
const userStore = useUserStore()

// 表单验证 - Form validation
const isFormValid = computed(() => {
  return formData.value.username.trim() !== '' && formData.value.password.trim() !== ''
})

// 登录处理 - Handle login
const handleLogin = async () => {
  if (!isFormValid.value || loading.value) return

  loading.value = true
  try {
    await userStore.login(formData.value.username, formData.value.password)
    
    uni.showToast({
      title: '登录成功',
      icon: 'success',
      duration: 1500
    })

    // 跳转到首页 - Navigate to home
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 1500)
  } catch (error: any) {
    uni.showToast({
      title: error.message || '登录失败，请重试',
      icon: 'none',
      duration: 2000
    })
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48rpx;
  position: relative;
  overflow: hidden;
}

.bg-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #3B82F6 100%);
  z-index: -1;
}

.login-card {
  width: 100%;
  max-width: 650rpx;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 32rpx;
  padding: 64rpx 48rpx;
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.15);
}

.logo-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 64rpx;
}

.logo-icon {
  width: 120rpx;
  height: 120rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.3);
}

.logo-text {
  font-size: 40rpx;
  font-weight: bold;
  color: #FFFFFF;
}

.app-name {
  font-size: 40rpx;
  font-weight: 600;
  color: #2C3E50;
  margin-bottom: 12rpx;
}

.app-desc {
  font-size: 26rpx;
  color: #64748B;
}

.form-area {
  display: flex;
  flex-direction: column;
  gap: 32rpx;
}

.input-group {
  width: 100%;
}

.input-wrapper {
  display: flex;
  align-items: center;
  background: #F5F7FA;
  border-radius: 16rpx;
  padding: 24rpx 32rpx;
  border: 2rpx solid transparent;
  transition: all 0.2s ease;
}

.input-wrapper:focus-within {
  border-color: #3B82F6;
  background: #FFFFFF;
  box-shadow: 0 0 0 4rpx rgba(59, 130, 246, 0.1);
}

.input-icon {
  font-size: 36rpx;
  margin-right: 20rpx;
}

.form-input {
  flex: 1;
  font-size: 30rpx;
  color: #2C3E50;
}

.login-btn {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  border: none;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16rpx;
  transition: all 0.2s ease;
}

.login-btn:active {
  opacity: 0.8;
  transform: scale(0.98);
}

.login-btn[disabled] {
  opacity: 0.6;
}

.btn-loading {
  opacity: 0.8;
}

.loading-icon {
  margin-right: 12rpx;
}

.btn-text {
  font-size: 32rpx;
  font-weight: 500;
  color: #FFFFFF;
}

.bottom-tip {
  margin-top: 48rpx;
  text-align: center;
}

.tip-text {
  font-size: 24rpx;
  color: #94A3B8;
}

.version-info {
  position: absolute;
  bottom: 48rpx;
  left: 50%;
  transform: translateX(-50%);
}

.version-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}
</style>
