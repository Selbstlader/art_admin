<template>
  <!-- 
    我的页面 - C4D风格高级UI
    My profile page - C4D style premium UI with glass morphism
  -->
  <view class="page-container">
    <!-- 顶部区域 - Header area with blue gradient -->
    <view class="header-area">
      <view class="header-bg"></view>
      <view class="header-content">
        <text class="header-title">我的</text>
      </view>
    </view>

    <!-- 用户信息卡片 - User profile card with glass morphism -->
    <view class="profile-card">
      <view class="avatar-wrapper">
        <image 
          v-if="userInfo?.avatar" 
          :src="userInfo.avatar" 
          class="avatar-image"
          mode="aspectFill"
        />
        <view v-else class="avatar-placeholder">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </view>
        <view class="avatar-shadow"></view>
      </view>
      
      <view class="user-info">
        <text class="user-name">{{ userInfo?.nickname || '未登录' }}</text>
        <view class="role-badge">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"></path>
          </svg>
          <text class="role-text">{{ roleText }}</text>
        </view>
      </view>
    </view>

    <!-- 功能菜单区域 - Menu section -->
    <view class="menu-section">
      <text class="section-title">设置</text>
      
      <!-- 关于 - About -->
      <view class="menu-card" @click="showAbout">
        <view class="menu-left">
          <view class="menu-icon-wrapper about-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="M12 16v-4"></path>
              <path d="M12 8h.01"></path>
            </svg>
          </view>
          <view class="menu-info">
            <text class="menu-title">关于</text>
            <text class="menu-desc">版本信息与应用介绍</text>
          </view>
        </view>
        <view class="menu-right">
          <text class="version-text">v{{ appVersion }}</text>
          <view class="menu-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>
      </view>
    </view>

    <!-- 退出登录按钮 - Logout button -->
    <view class="logout-section">
      <view class="logout-btn" @click="handleLogout">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" x2="9" y1="12" y2="12"></line>
        </svg>
        <text class="logout-text">退出登录</text>
      </view>
    </view>

    <!-- 底部版权信息 - Footer copyright -->
    <view class="footer">
      <text class="copyright">© 2024 工装设计师AI助手</text>
    </view>

    <!-- 关于弹窗 - About modal -->
    <view v-if="showAboutModal" class="modal-overlay" @click="closeAbout">
      <view class="modal-content" @click.stop>
        <view class="modal-header">
          <text class="modal-title">关于</text>
          <view class="modal-close" @click="closeAbout">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </view>
        </view>
        <view class="modal-body">
          <view class="app-logo">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"></path>
            </svg>
          </view>
          <text class="app-name">工装设计师AI助手</text>
          <text class="app-version">版本 {{ appVersion }}</text>
          <view class="app-desc">
            <text class="desc-text">专为工装设计师打造的AI辅助工具，</text>
            <text class="desc-text">支持项目管理、文档分析、成本预算等功能。</text>
          </view>
        </view>
        <view class="modal-footer">
          <view class="confirm-btn" @click="closeAbout">
            <text class="confirm-text">确定</text>
          </view>
        </view>
      </view>
    </view>
    
    <!-- 自定义底部导航栏 - Custom TabBar -->
    <CustomTabBar />
  </view>
</template>

<script setup lang="ts">
/*** 
 * My profile page - C4D style premium UI
 * 我的页面 - 展示用户信息、退出登录、关于页面
 * Requirements: 11.1, 11.2, 11.3, 1.5
 ***/
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import CustomTabBar from '@/components/CustomTabBar.vue'

/*** Store instance - 状态管理实例 ***/
const userStore = useUserStore()

/*** App version - 应用版本号 ***/
const appVersion = ref('1.0.0')

/*** About modal state - 关于弹窗状态 ***/
const showAboutModal = ref(false)

/*** User info from store - 从 store 获取用户信息 ***/
const userInfo = computed(() => userStore.userInfo)

/*** Role text mapping - 角色文本映射 ***/
const roleText = computed((): string => {
  const role = userInfo.value?.role
  if (!role) return '未知角色'
  
  const roleMap: Record<string, string> = {
    'admin': '管理员',
    'designer': '设计师',
    'manager': '项目经理',
    'viewer': '访客'
  }
  return roleMap[role] || role
})

/*** Show about modal - 显示关于弹窗 ***/
const showAbout = (): void => {
  showAboutModal.value = true
}

/*** Close about modal - 关闭关于弹窗 ***/
const closeAbout = (): void => {
  showAboutModal.value = false
}

/*** Handle logout with confirmation - 处理退出登录（带确认弹窗） ***/
const handleLogout = (): void => {
  uni.showModal({
    title: '确认退出',
    content: '确定要退出登录吗？',
    confirmText: '退出',
    confirmColor: '#EF4444',
    cancelText: '取消',
    success: async (res) => {
      if (res.confirm) {
        try {
          // 显示加载提示 - Show loading toast
          uni.showLoading({ title: '退出中...', mask: true })
          
          // 调用 store 退出方法 - Call store logout method
          await userStore.logout()
          
          uni.hideLoading()
          
          // 提示退出成功 - Show success toast
          uni.showToast({
            title: '已退出登录',
            icon: 'success',
            duration: 1500
          })
          
          // 延迟跳转登录页 - Delay navigation to login page
          setTimeout(() => {
            uni.reLaunch({ url: '/pages/login/index' })
          }, 1500)
        } catch (error) {
          uni.hideLoading()
          console.error('退出登录失败:', error)
          
          // 即使失败也跳转登录页 - Navigate to login even on failure
          uni.reLaunch({ url: '/pages/login/index' })
        }
      }
    }
  })
}

/*** Initialize on mount - 组件挂载时初始化 ***/
onMounted(() => {
  // 如果没有用户信息，尝试获取 - Fetch user info if not available
  if (!userInfo.value && userStore.isLoggedIn) {
    userStore.getUserInfo()
  }
})
</script>

<style lang="scss" scoped>
/*** 页面容器 - Page container with gradient background ***/
.page-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #E8F4FD 0%, #FFFFFF 30%, #F8FAFC 100%);
  padding-bottom: env(safe-area-inset-bottom);
}

/*** 顶部区域 - Header area with blue gradient ***/
.header-area {
  position: relative;
  padding: 60rpx 32rpx 120rpx;
  overflow: hidden;
}

.header-bg {
  position: absolute;
  top: -100rpx;
  left: -50rpx;
  right: -50rpx;
  height: 400rpx;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 50%, #1D4ED8 100%);
  border-radius: 0 0 60rpx 60rpx;
  transform: rotate(-3deg);
}

.header-content {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #FFFFFF;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
}

/*** 用户信息卡片 - Profile card with glass morphism ***/
.profile-card {
  position: relative;
  z-index: 10;
  margin: -80rpx 32rpx 32rpx;
  padding: 40rpx 32rpx;
  background: #FFFFFF;
  border-radius: 28rpx;
  box-shadow: 
    0 8rpx 32rpx rgba(59, 130, 246, 0.12),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-wrapper {
  position: relative;
  margin-bottom: 20rpx;
}

.avatar-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  border: 6rpx solid #FFFFFF;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.2);
}

.avatar-placeholder {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3B82F6;
  border: 6rpx solid #FFFFFF;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.2);
}

.avatar-shadow {
  position: absolute;
  bottom: -12rpx;
  left: 20rpx;
  right: 20rpx;
  height: 20rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.user-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 12rpx;
}

.role-badge {
  display: flex;
  align-items: center;
  padding: 10rpx 24rpx;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border-radius: 24rpx;
  color: #3B82F6;
}

.role-text {
  font-size: 24rpx;
  font-weight: 500;
  color: #3B82F6;
  margin-left: 8rpx;
}

/*** 菜单区域 - Menu section ***/
.menu-section {
  padding: 0 32rpx;
  margin-bottom: 32rpx;
}

.section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #64748B;
  margin-bottom: 16rpx;
  display: block;
  padding-left: 8rpx;
}

.menu-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

.menu-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.menu-icon-wrapper {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 18rpx;
  margin-right: 20rpx;
}

.about-icon {
  background: linear-gradient(135deg, #E0E7FF 0%, #C7D2FE 100%);
  color: #6366F1;
}

.menu-info {
  display: flex;
  flex-direction: column;
}

.menu-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1E293B;
  margin-bottom: 4rpx;
}

.menu-desc {
  font-size: 24rpx;
  color: #94A3B8;
}

.menu-right {
  display: flex;
  align-items: center;
}

.version-text {
  font-size: 24rpx;
  color: #94A3B8;
  margin-right: 8rpx;
}

.menu-arrow {
  color: #CBD5E1;
}

/*** 退出登录区域 - Logout section ***/
.logout-section {
  padding: 0 32rpx;
  margin-top: 48rpx;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 28rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  color: #EF4444;
  
  &:active {
    transform: scale(0.98);
    background: #FEF2F2;
  }
}

.logout-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #EF4444;
  margin-left: 12rpx;
}

/*** 底部版权 - Footer copyright ***/
.footer {
  padding: 48rpx 32rpx;
  display: flex;
  justify-content: center;
}

.copyright {
  font-size: 22rpx;
  color: #CBD5E1;
}

/*** 弹窗遮罩 - Modal overlay ***/
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 48rpx;
}

.modal-content {
  width: 100%;
  max-width: 600rpx;
  background: #FFFFFF;
  border-radius: 32rpx;
  overflow: hidden;
  box-shadow: 0 24rpx 48rpx rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx;
  border-bottom: 1rpx solid #F1F5F9;
}

.modal-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
}

.modal-close {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #F8FAFC;
  color: #64748B;
  transition: all 0.2s ease;
  
  &:active {
    background: #F1F5F9;
  }
}

.modal-body {
  padding: 48rpx 32rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.app-logo {
  width: 120rpx;
  height: 120rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-radius: 32rpx;
  color: #FFFFFF;
  margin-bottom: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.3);
}

.app-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #1E293B;
  margin-bottom: 8rpx;
}

.app-version {
  font-size: 26rpx;
  color: #64748B;
  margin-bottom: 24rpx;
}

.app-desc {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.desc-text {
  font-size: 26rpx;
  color: #94A3B8;
  line-height: 1.6;
  text-align: center;
}

.modal-footer {
  padding: 24rpx 32rpx 32rpx;
}

.confirm-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  border-radius: 24rpx;
  padding: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.3);
  transition: all 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

.confirm-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #FFFFFF;
}
</style>
