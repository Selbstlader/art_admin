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
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </view>
            <input class="form-input" type="text" placeholder="请输入账号" v-model="formData.username"
              @focus="usernameFocus = true" @blur="usernameFocus = false" @confirm="handleLogin" />
          </view>
        </view>

        <view class="input-group">
          <view class="input-wrapper" :class="{ 'input-focus': passwordFocus }">
            <view class="input-icon-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
            </view>
            <input class="form-input" type="password" placeholder="请输入密码" v-model="formData.password"
              @focus="passwordFocus = true" @blur="passwordFocus = false" @confirm="handleLogin" />
          </view>
        </view>

        <!-- 登录按钮 - Login button with 3D effect -->
        <view class="btn-container">
          <button class="login-btn" :class="{ 'btn-loading': loading, 'btn-disabled': !isFormValid }"
            :disabled="loading || !isFormValid" @click="handleLogin">
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
/*** Login page styles - C4D style premium UI with glass morphism ***/
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40rpx;
  background: linear-gradient(165deg, #E3F0FF 0%, #FFFFFF 40%, #F8FBFF 70%, #EDF5FF 100%);
  position: relative;
  overflow: hidden;
}

/*** Background decoration circles - 背景装饰圆 ***/
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
}

.bg-circle-1 {
  width: 500rpx;
  height: 500rpx;
  top: -150rpx;
  right: -120rpx;
  background: linear-gradient(180deg, #A5D0FF 0%, #60A5FA 50%, #3B82F6 100%);
  opacity: 0.35;
  filter: blur(80rpx);
}

.bg-circle-2 {
  width: 400rpx;
  height: 400rpx;
  bottom: 100rpx;
  left: -150rpx;
  background: linear-gradient(180deg, #BFDBFE 0%, #93C5FD 100%);
  opacity: 0.4;
  filter: blur(70rpx);
}

.bg-circle-3 {
  width: 250rpx;
  height: 250rpx;
  top: 35%;
  right: -80rpx;
  background: linear-gradient(180deg, #DBEAFE 0%, #BFDBFE 100%);
  opacity: 0.5;
  filter: blur(50rpx);
}

/*** Login card with glass morphism - 登录卡片玻璃拟态 ***/
.login-card {
  width: 88%;
  max-width: 680rpx;
  // background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(60rpx);
  -webkit-backdrop-filter: blur(60rpx);
  border-radius: 48rpx;
  padding: 72rpx 56rpx 56rpx;
  box-shadow:
    0 20rpx 60rpx rgba(59, 130, 246, 0.08),
    0 8rpx 24rpx rgba(0, 0, 0, 0.03),
    inset 0 2rpx 0 rgba(255, 255, 255, 0.9);
  border: 1rpx solid rgba(255, 255, 255, 0.7);
  position: relative;
  z-index: 10;
}

/*** Logo area - Logo 区域 ***/
.logo-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 64rpx;
}

.logo-container {
  position: relative;
  margin-bottom: 32rpx;
}

.logo-icon {
  width: 140rpx;
  height: 140rpx;
  background: linear-gradient(155deg, #60A5FA 0%, #3B82F6 40%, #2563EB 100%);
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
  transform: perspective(600rpx) rotateX(8deg) rotateY(-3deg);
  box-shadow:
    0 6rpx 0 #1D4ED8,
    0 12rpx 24rpx rgba(37, 99, 235, 0.35),
    inset 0 2rpx 4rpx rgba(255, 255, 255, 0.3);
}

.logo-inner {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 36rpx;
  background: linear-gradient(155deg, rgba(255, 255, 255, 0.25) 0%, transparent 60%);
}

.logo-text {
  font-size: 48rpx;
  font-weight: 700;
  color: #FFFFFF;
  text-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.15);
  letter-spacing: 2rpx;
}

.logo-shadow {
  position: absolute;
  bottom: -20rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 100rpx;
  height: 20rpx;
  background: radial-gradient(ellipse, rgba(37, 99, 235, 0.25) 0%, transparent 70%);
  border-radius: 50%;
}

.app-name {
  font-size: 48rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 16rpx;
  letter-spacing: 4rpx;
}

.app-desc {
  font-size: 28rpx;
  color: #64748B;
  letter-spacing: 2rpx;
}

/*** Form area - 表单区域 ***/
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
  background: rgba(248, 250, 252, 0.9);
  border-radius: 24rpx;
  padding: 0 32rpx;
  height: 108rpx;
  border: 2rpx solid rgba(226, 232, 240, 0.6);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: inset 0 2rpx 6rpx rgba(0, 0, 0, 0.02);
}

.input-wrapper.input-focus {
  border-color: rgba(59, 130, 246, 0.5);
  background: #FFFFFF;
  box-shadow:
    0 0 0 6rpx rgba(59, 130, 246, 0.08),
    inset 0 2rpx 6rpx rgba(0, 0, 0, 0.01);
}

.input-icon-wrapper {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 24rpx;
  flex-shrink: 0;
  transition: color 0.3s ease;
}

.input-focus .input-icon-wrapper {
  color: #3B82F6;
}

.form-input {
  flex: 1;
  font-size: 32rpx;
  color: #1E293B;
  height: 100%;
  background: transparent;
}

.form-input::placeholder {
  color: #94A3B8;
  font-size: 30rpx;
}

/*** Login button with 3D effect - 登录按钮3D效果 ***/
.btn-container {
  position: relative;
  margin-top: 24rpx;
}

.login-btn {
  width: 100%;
  height: 108rpx;
  background: linear-gradient(155deg, #60A5FA 0%, #3B82F6 40%, #2563EB 100%);
  border: none;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow:
    0 6rpx 0 #1D4ED8,
    0 12rpx 28rpx rgba(37, 99, 235, 0.35);
}

// .login-btn::before {
//   content: '';
//   position: absolute;
//   top: 0;
//   left: 0;
//   right: 0;
//   height: 50%;
//   background: linear-gradient(180deg, rgba(255, 255, 255, 0.2) 0%, transparent 100%);
//   border-radius: 24rpx 24rpx 0 0;
//   pointer-events: none;
// }

.login-btn::after {
  border: none;
}

.login-btn:active {
  transform: translateY(6rpx);
  box-shadow:
    0 0 0 #1D4ED8,
    0 6rpx 16rpx rgba(37, 99, 235, 0.25);
}

// .login-btn.btn-disabled {
//   background: linear-gradient(155deg, #E2E8F0 0%, #CBD5E1 40%, #94A3B8 100%);
//   box-shadow:
//     0 6rpx 0 #64748B,
//     0 12rpx 28rpx rgba(100, 116, 139, 0.2);
// }

.login-btn.btn-loading {
  opacity: 0.9;
}

.btn-shadow {
  position: absolute;
  bottom: -16rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 75%;
  height: 24rpx;
  background: radial-gradient(ellipse, rgba(37, 99, 235, 0.18) 0%, transparent 70%);
  border-radius: 50%;
  z-index: 1;
}

.loading-spinner {
  width: 40rpx;
  height: 40rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
  border-top-color: #FFFFFF;
  border-radius: 50%;
  margin-right: 20rpx;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.btn-text {
  font-size: 34rpx;
  font-weight: 600;
  color: #FFFFFF;
  letter-spacing: 4rpx;
}

/*** Bottom tip - 底部提示 ***/
.bottom-tip {
  margin-top: 56rpx;
  text-align: center;
}

.tip-text {
  font-size: 26rpx;
  color: #94A3B8;
  letter-spacing: 1rpx;
}

/*** Version info - 版本信息 ***/
.version-info {
  position: absolute;
  bottom: 80rpx;
  left: 50%;
  transform: translateX(-50%);
}

.version-text {
  font-size: 24rpx;
  color: #B0BEC5;
  letter-spacing: 1rpx;
}
</style>
