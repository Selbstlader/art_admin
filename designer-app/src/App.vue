<script setup lang="ts">
/*** Root App component - handles app lifecycle events ***/
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'

onLaunch(() => {
  console.log('App Launch')
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style>
/* 全局样式 - Global styles */
page {
  background-color: #F5F7FA;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 28rpx;
  color: #2C3E50;
  box-sizing: border-box;
}

/* 蓝色渐变主题色 - Blue gradient theme */
.gradient-bg {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.gradient-blue {
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
}

/* 玻璃拟态样式 - Glass morphism style */
.glass-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
}

/* 悬浮卡片样式 - Floating card style */
.float-card {
  background: #FFFFFF;
  border-radius: 16rpx;
  box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.float-card:active {
  transform: scale(0.98);
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.08);
}

/* 按钮基础样式 - Button base style */
.btn-primary {
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  color: #FFFFFF;
  border: none;
  border-radius: 12rpx;
  padding: 24rpx 48rpx;
  font-size: 32rpx;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-primary:active {
  opacity: 0.8;
  transform: scale(0.98);
}

/* 文本样式 - Text styles */
.text-primary {
  color: #3B82F6;
}

.text-secondary {
  color: #64748B;
}

.text-danger {
  color: #EF4444;
}

.text-success {
  color: #10B981;
}

/* 间距工具类 - Spacing utilities */
.p-8 { padding: 16rpx; }
.p-16 { padding: 32rpx; }
.p-24 { padding: 48rpx; }

.m-8 { margin: 16rpx; }
.m-16 { margin: 32rpx; }
.m-24 { margin: 48rpx; }

.mb-8 { margin-bottom: 16rpx; }
.mb-16 { margin-bottom: 32rpx; }
.mb-24 { margin-bottom: 48rpx; }

.mt-8 { margin-top: 16rpx; }
.mt-16 { margin-top: 32rpx; }
.mt-24 { margin-top: 48rpx; }

/* Flex 布局工具类 - Flex layout utilities */
.flex { display: flex; }
.flex-col { flex-direction: column; }
.flex-row { flex-direction: row; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.flex-1 { flex: 1; }
.flex-wrap { flex-wrap: wrap; }

/* 安全区域适配 - Safe area adaptation */
.safe-area-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.safe-area-top {
  padding-top: constant(safe-area-inset-top);
  padding-top: env(safe-area-inset-top);
}
</style>
