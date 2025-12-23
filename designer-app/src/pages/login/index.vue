<template>
  <view class="login-page">
    <!-- 背景装饰 - Background decoration -->
    <view class="bg-decoration">
      <view class="bg-circle bg-circle-1"></view>
      <view class="bg-circle bg-circle-2"></view>
      <view class="bg-circle bg-circle-3"></view>
    </view>
    
    <!-- 登录卡片 - Login card with glass morphism -->
    <view class="login-card">
      <!-- Logo 区域 - Logo area -->
      <view class="logo-area">
        <view class="logo-container">
          <view class="logo-icon">
            <view class="logo-inner">
              <text class="logo-text">设计</text>
            </view>
          </view>
          <view class="logo-shadow"></view>
        </view>
        <text class="app-name">工装设计助手</text>
        <text class="app-desc">AI赋能设计，让创意更高效</text>
      </view>

      <!-- 表单区域 - Form area -->
      <view class="form-area">
        <view class="input-group">
          <view class="input-wrapper" :class="{ 'input-focus': usernameFocus }">
            <view class="input-icon-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </view>
            <input
              class="form-input"
              type="text"
              placeholder="请输入账号"
              v-model="formData.username"
              @focus="usernameFocus = true"
              @blur="usernameFocus = false"
              @confirm="handleLogin"
            />
          </view>
        </view>

        <view class="input-group">
          <view class="input-wrapper" :class="{ 'input-focus': passwordFocus }">
            <view class="input-icon-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
            </view>
            <input
              class="form-input"
              type="password"
              placeholder="请输入密码"
              v-model="formData.password"
              @focus="passwordFocus = true"
              @blur="passwordFocus = false"
              @confirm="handleLogin"
            />
          </view>
        </view>

        <!-- 登录按钮 - Login button with 3D effect -->
        <view class="btn-container">
          <button
            class="login-btn"
            :class="{ 'btn-loading': loading, 'btn-disabled': !isFormValid }"
            :disabled="loading || !isFormValid"
            @click="handleLogin"
          >
            <view v-if="loading" class="loading-spinner"></view>
            <text class="btn-text">{{ loading ? '登录中...' : '登录' }}</text>
          </button>
          <view class="btn-shadow"></view>
        </view>
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
/*** Login page - C4D style premium UI with glass morphism ***/
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'

const formData = ref({
  username: '',
  password: ''
})

const loading = ref(false)
const usernameFocus = ref(false)
const passwordFocus = ref(false)
const userStore = useUserStore()

const isFormValid = computed(() => {
  return formData.value.username.trim() !== '' && formData.value.password.trim() !== ''
})

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
  padding: 40rpx;
  background: linear-gradient(180deg, #E8F4FD 0%, #FFFFFF 50%, #F0F7FF 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰圆 - Background decoration circles */
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.6;
}

.bg-circle-1 {
  width: 400rpx;
  height: 400rpx;
  top: -100rpx;
  right: -100rpx;
  background: linear-gradient(135deg, #60A5FA 0%, #3B82F6 100%);
  filter: blur(60rpx);
}

.bg-circle-2 {
  width: 300rpx;
  height: 300rpx;
  bottom: 200rpx;
  left: -80rpx;
  background: linear-gradient(135deg, #93C5FD 0%, #60A5FA 100%);
  filter: blur(50rpx);
}

.bg-circle-3 {
  width: 200rpx;
  height: 200rpx;
  top: 40%;
  right: -50rpx;
  background: linear-gradient(135deg, #BFDBFE 0%, #93C5FD 100%);
  filter: blur(40rpx);
}

/* 登录卡片 - Login card with glass morphism */
.login-card {
  width: 100%;
  max-width: 680rpx;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(40rpx);
  -webkit-backdrop-filter: blur(40rpx);
  border-radius: 40rpx;
  padding: 60rpx 48rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.08),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.8);
  border: 1rpx solid rgba(255, 255, 255, 0.6);
  position: relative;
  z-index: 10;
}

/* Logo 区域 - Logo area */
.logo-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 56rpx;
}

.logo-container {
  position: relative;
  margin-bottom: 28rpx;
}

.logo-icon {
  width: 120rpx;
  height: 120rpx;
  background: linear-gradient(145deg, #4F9CF9 0%, #2563EB 50%, #1D4ED8 100%);
  border-radius: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
  transform: perspective(500rpx) rotateX(5deg);
  box-shadow: 
    0 4rpx 0 #1E40AF,
    0 8rpx 16rpx rgba(37, 99, 235, 0.3);
}

.logo-inner {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 28rpx;
  background: linear-gradient(145deg, rgba(255,255,255,0.2) 0%, transparent 50%);
}

.logo-text {
  font-size: 40rpx;
  font-weight: 700;
  color: #FFFFFF;
  text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.2);
}

.logo-shadow {
  position: absolute;
  bottom: -16rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 80rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(37, 99, 235, 0.3) 0%, transparent 70%);
  border-radius: 50%;
}

.app-name {
  font-size: 44rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 12rpx;
  letter-spacing: 2rpx;
}

.app-desc {
  font-size: 26rpx;
  color: #64748B;
  letter-spacing: 1rpx;
}

/* 表单区域 - Form area */
.form-area {
  display: flex;
  flex-direction: column;
  gap: 28rpx;
}

.input-group {
  width: 100%;
}

.input-wrapper {
  display: flex;
  align-items: center;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 20rpx;
  padding: 0 28rpx;
  height: 100rpx;
  border: 2rpx solid rgba(226, 232, 240, 0.8);
  transition: all 0.3s ease;
  box-shadow: inset 0 2rpx 4rpx rgba(0, 0, 0, 0.02);
}

.input-wrapper.input-focus {
  border-color: #3B82F6;
  background: #FFFFFF;
  box-shadow: 
    0 0 0 4rpx rgba(59, 130, 246, 0.1),
    inset 0 2rpx 4rpx rgba(0, 0, 0, 0.02);
}

.input-icon-wrapper {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.input-focus .input-icon-wrapper {
  color: #3B82F6;
}

.form-input {
  flex: 1;
  font-size: 30rpx;
  color: #1E293B;
  height: 100%;
}

.form-input::placeholder {
  color: #94A3B8;
}

/* 登录按钮 - Login button with 3D effect */
.btn-container {
  position: relative;
  margin-top: 16rpx;
}

.login-btn {
  width: 100%;
  height: 100rpx;
  background: linear-gradient(145deg, #4F9CF9 0%, #2563EB 50%, #1D4ED8 100%);
  border: none;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
  transition: all 0.2s ease;
  box-shadow: 
    0 4rpx 0 #1E40AF,
    0 8rpx 20rpx rgba(37, 99, 235, 0.35);
}

.login-btn::after {
  border: none;
}

.login-btn:active {
  transform: translateY(4rpx);
  box-shadow: 
    0 0 0 #1E40AF,
    0 4rpx 12rpx rgba(37, 99, 235, 0.25);
}

.login-btn.btn-disabled {
  background: linear-gradient(145deg, #CBD5E1 0%, #94A3B8 100%);
  box-shadow: 
    0 4rpx 0 #64748B,
    0 8rpx 20rpx rgba(100, 116, 139, 0.2);
}

.login-btn.btn-loading {
  opacity: 0.9;
}

.btn-shadow {
  position: absolute;
  bottom: -12rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 80%;
  height: 20rpx;
  background: radial-gradient(ellipse, rgba(37, 99, 235, 0.2) 0%, transparent 70%);
  border-radius: 50%;
  z-index: 1;
}

.loading-spinner {
  width: 36rpx;
  height: 36rpx;
  border: 3rpx solid rgba(255, 255, 255, 0.3);
  border-top-color: #FFFFFF;
  border-radius: 50%;
  margin-right: 16rpx;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.btn-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #FFFFFF;
  letter-spacing: 2rpx;
}

/* 底部提示 - Bottom tip */
.bottom-tip {
  margin-top: 48rpx;
  text-align: center;
}

.tip-text {
  font-size: 24rpx;
  color: #94A3B8;
}

/* 版本信息 - Version info */
.version-info {
  position: absolute;
  bottom: 60rpx;
  left: 50%;
  transform: translateX(-50%);
}

.version-text {
  font-size: 24rpx;
  color: #94A3B8;
}
</style>
