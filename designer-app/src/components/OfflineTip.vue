<template>
  <!--
    离线模式提示组件 - Offline mode tip component
    全局显示网络状态提示 - Global network status indicator
  -->
  <view v-if="showTip" class="offline-tip-container" :class="{ 'is-recovering': networkStore.isRecovering }">
    <view class="offline-tip-content">
      <!-- 离线图标 - Offline icon -->
      <view class="offline-icon">
        <svg v-if="!networkStore.isRecovering" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12.01 21.49 23.64 7c-.45-.34-4.93-4-11.64-4C5.28 3 .81 6.66.36 7l11.63 14.49.01.01.01-.01z"></path>
          <line x1="2" x2="22" y1="2" y2="22"></line>
        </svg>
        <!-- 恢复中图标 - Recovering icon -->
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12.01 21.49 23.64 7c-.45-.34-4.93-4-11.64-4C5.28 3 .81 6.66.36 7l11.63 14.49.01.01.01-.01z"></path>
        </svg>
      </view>
      <!-- 提示文字 - Tip text -->
      <text class="offline-text">{{ tipText }}</text>
      <!-- 关闭按钮 - Close button -->
      <view v-if="closable && !networkStore.isRecovering" class="close-btn" @click="handleClose">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" x2="6" y1="6" y2="18"></line>
          <line x1="6" x2="18" y1="6" y2="18"></line>
        </svg>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
/***
 * Offline tip component - displays network status globally
 * 离线提示组件 - 全局显示网络状态
 ***/
import { computed } from 'vue'
import { useNetworkStore } from '@/store/network'

/*** Props definition - 属性定义 ***/
interface Props {
  /*** 是否可关闭 - Is closable ***/
  closable?: boolean
  /*** 自定义位置（bottom 距离，单位 rpx）- Custom position (bottom distance in rpx) ***/
  bottom?: number
}

const props = withDefaults(defineProps<Props>(), {
  closable: true,
  bottom: 180
})

/*** Store instance - 状态管理实例 ***/
const networkStore = useNetworkStore()

/*** 是否显示提示 - Show tip condition ***/
const showTip = computed(() => {
  return networkStore.showOfflineTip || networkStore.isRecovering
})

/*** 提示文字 - Tip text ***/
const tipText = computed(() => {
  if (networkStore.isRecovering) {
    return '网络已恢复，正在刷新...'
  }
  return '当前为离线模式'
})

/*** 处理关闭 - Handle close ***/
const handleClose = (): void => {
  networkStore.hideOfflineTip()
}
</script>

<style lang="scss" scoped>
/*** 离线提示容器 - Offline tip container ***/
.offline-tip-container {
  position: fixed;
  bottom: v-bind('props.bottom + "rpx"');
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/*** 离线提示内容 - Offline tip content ***/
.offline-tip-content {
  display: flex;
  align-items: center;
  background: rgba(30, 41, 59, 0.92);
  padding: 16rpx 28rpx;
  border-radius: 40rpx;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

/*** 恢复中状态 - Recovering state ***/
.is-recovering .offline-tip-content {
  background: rgba(16, 185, 129, 0.92);
}

/*** 离线图标 - Offline icon ***/
.offline-icon {
  width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FBBF24;
  margin-right: 12rpx;
}

.is-recovering .offline-icon {
  color: #FFFFFF;
}

/*** 提示文字 - Tip text ***/
.offline-text {
  font-size: 24rpx;
  color: #FFFFFF;
  white-space: nowrap;
}

/*** 关闭按钮 - Close button ***/
.close-btn {
  width: 28rpx;
  height: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.7);
  margin-left: 16rpx;
  transition: color 0.2s ease;
  
  &:active {
    color: #FFFFFF;
  }
}
</style>
