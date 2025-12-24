<template>
  <!-- 
    项目卡片组件 - C4D风格玻璃拟态UI
    Project Card Component - C4D style glassmorphism premium UI
  -->
  <view class="project-card" @click="handleClick">
    <!-- 卡片玻璃主体 - Card glass body -->
    <view class="card-glass">
      <!-- 玻璃高光层 - Glass highlight layer -->
      <view class="glass-highlight"></view>
      <!-- 玻璃边框光效 - Glass border glow -->
      <view class="glass-border"></view>

      <!-- 卡片内容 - Card content -->
      <view class="card-inner">
        <!-- 头部区域 - Header area -->
        <view class="card-header">
          <view class="title-area">
            <text class="project-name">{{ project.name }}</text>
          </view>
          <view :class="['status-pill', `status-${project.status}`]">
            <view class="pill-glow"></view>
            <view class="pill-dot"></view>
            <text class="pill-text">{{ statusText }}</text>
          </view>
        </view>

        <!-- 数据展示区 - Data display area with 3D effect -->
        <view class="data-section">
          <view class="data-card">
            <view class="data-card-inner">
              <view class="data-icon-box">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="18" height="18" x="3" y="3" rx="2"></rect>
                  <path d="M3 9h18"></path>
                  <path d="M9 21V9"></path>
                </svg>
              </view>
              <view class="data-info">
                <text class="data-label">面积</text>
                <view class="data-value-row">
                  <text class="data-number">{{ project.area }}</text>
                  <text class="data-unit">㎡</text>
                </view>
              </view>
            </view>
          </view>

          <view class="data-card">
            <view class="data-card-inner">
              <view class="data-icon-box">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </view>
              <view class="data-info">
                <text class="data-label">更新时间</text>
                <text class="data-date">{{ formattedDate }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- 底部区域 - Footer area -->
        <view class="card-footer">
          <view class="budget-area">
            <text class="budget-label">预算</text>
            <view class="budget-display">
              <text class="budget-symbol">¥</text>
              <text class="budget-amount">{{ formattedBudget }}</text>
            </view>
          </view>
          <view class="action-btn">
            <view class="btn-glow"></view>
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"></path>
            </svg>
          </view>
        </view>
      </view>
    </view>

    <!-- 3D悬浮阴影 - 3D floating shadow -->
    <view class="card-float-shadow"></view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * ProjectCard Component - C4D Glassmorphism Style
 * 项目卡片组件 - C4D风格玻璃拟态设计
 ***/
import { computed } from 'vue'
import type { Project } from '@/types/api'
import { ProjectStatusText, ProjectStatus } from '@/types/common'

interface Props {
  project: Project
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'click', id: number): void
}>()

/*** 状态文本映射 ***/
const statusText = computed((): string => {
  return ProjectStatusText[props.project.status as ProjectStatus] || props.project.status
})

/*** 格式化日期 ***/
const formattedDate = computed((): string => {
  if (!props.project.updatedAt) return '-'
  const date = new Date(props.project.updatedAt)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
})

/*** 格式化预算 ***/
const formattedBudget = computed((): string => {
  const budget = props.project.budget
  if (budget >= 10000) {
    return (budget / 10000).toFixed(1) + '万'
  }
  return budget.toLocaleString()
})

const handleClick = (): void => {
  emit('click', props.project.id)
}
</script>

<style lang="scss" scoped>
/*** 项目卡片容器 ***/
.project-card {
  position: relative;
  margin-bottom: 8rpx;
  transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);

  &:active {
    transform: scale(0.97) translateY(6rpx);

    .card-float-shadow {
      opacity: 0.3;
      transform: translateY(-8rpx) scale(0.9);
    }

    .action-btn {
      transform: translateX(6rpx);
    }
  }
}

/*** 玻璃卡片主体 ***/
.card-glass {
  position: relative;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.75) 0%, rgba(255, 255, 255, 0.55) 20%, rgba(240, 245, 255, 0.65) 100%);
  -webkit-backdrop-filter: blur(40px);
  border-radius: 32rpx;
  overflow: hidden;
}

/*** 玻璃高光效果 ***/
.glass-highlight {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 50%;
  background: linear-gradient(180deg,
      rgba(255, 255, 255, 0.8) 0%,
      rgba(255, 255, 255, 0.2) 60%,
      transparent 100%);
  pointer-events: none;
}

/*** 玻璃边框光效 ***/
.glass-border {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 32rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 8rpx 32rpx rgba(59, 130, 246, 0.12),
    0 2rpx 8rpx rgba(0, 0, 0, 0.04),
    inset 0 1rpx 0 rgba(255, 255, 255, 1),
    inset 0 -1rpx 0 rgba(0, 0, 0, 0.02);
  pointer-events: none;
}

/*** 卡片内容区 ***/
.card-inner {
  position: relative;
  padding: 32rpx;
  z-index: 1;
}

/*** 头部区域 ***/
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 28rpx;
}

.title-area {
  flex: 1;
  margin-right: 20rpx;
}

.project-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #0F172A;
  line-height: 1.4;
  letter-spacing: 1rpx;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/*** 状态胶囊 ***/
.status-pill {
  position: relative;
  display: flex;
  align-items: center;
  padding: 10rpx 20rpx;
  border-radius: 28rpx;
  flex-shrink: 0;
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.pill-glow {
  position: absolute;
  inset: 0;
  opacity: 0.12;
}

.pill-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  margin-right: 10rpx;
  position: relative;
  z-index: 1;
}

.pill-text {
  font-size: 24rpx;
  font-weight: 600;
  position: relative;
  z-index: 1;
}

/*** 状态颜色变体 ***/
.status-draft {
  background: rgba(113, 113, 122, 0.08);

  .pill-glow {
    background: linear-gradient(135deg, #A1A1AA, #71717A);
  }

  .pill-dot {
    background: linear-gradient(135deg, #A1A1AA, #71717A);
    box-shadow: 0 2rpx 8rpx rgba(113, 113, 122, 0.4);
  }

  .pill-text {
    color: #52525B;
  }
}

.status-in_progress {
  background: rgba(59, 130, 246, 0.08);

  .pill-glow {
    background: linear-gradient(135deg, #60A5FA, #3B82F6);
  }

  .pill-dot {
    background: linear-gradient(135deg, #60A5FA, #3B82F6);
    box-shadow: 0 2rpx 12rpx rgba(59, 130, 246, 0.6);
    animation: pulse-glow 2s ease-in-out infinite;
  }

  .pill-text {
    color: #2563EB;
  }
}

.status-completed {
  background: rgba(34, 197, 94, 0.08);

  .pill-glow {
    background: linear-gradient(135deg, #4ADE80, #22C55E);
  }

  .pill-dot {
    background: linear-gradient(135deg, #4ADE80, #22C55E);
    box-shadow: 0 2rpx 8rpx rgba(34, 197, 94, 0.5);
  }

  .pill-text {
    color: #16A34A;
  }
}

.status-archived {
  background: rgba(249, 115, 22, 0.08);

  .pill-glow {
    background: linear-gradient(135deg, #FB923C, #F97316);
  }

  .pill-dot {
    background: linear-gradient(135deg, #FB923C, #F97316);
    box-shadow: 0 2rpx 8rpx rgba(249, 115, 22, 0.4);
  }

  .pill-text {
    color: #EA580C;
  }
}

@keyframes pulse-glow {

  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.7;
    transform: scale(0.9);
  }
}

/*** 数据展示区 ***/
.data-section {
  display: flex;
  gap: 20rpx;
  margin-bottom: 28rpx;
}

.data-card {
  flex: 1;
  background: linear-gradient(145deg,
      rgba(248, 250, 252, 0.8) 0%,
      rgba(241, 245, 249, 0.6) 100%);
  border-radius: 20rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.8);
  box-shadow: inset 0 2rpx 4rpx rgba(255, 255, 255, 0.9),
    inset 0 -1rpx 2rpx rgba(0, 0, 0, 0.02),
    0 2rpx 8rpx rgba(0, 0, 0, 0.03);
}

.data-card-inner {
  display: flex;
  align-items: center;
  padding: 20rpx;
}

.data-icon-box {
  width: 52rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(145deg, #FFFFFF 0%, #F1F5F9 100%);
  border-radius: 16rpx;
  margin-right: 16rpx;
  color: #3B82F6;
  box-shadow: 0 4rpx 12rpx rgba(59, 130, 246, 0.15),
    inset 0 1rpx 2rpx rgba(255, 255, 255, 1);
}

.data-info {
  display: flex;
  flex-direction: column;
}

.data-label {
  font-size: 22rpx;
  color: #94A3B8;
  font-weight: 500;
  margin-bottom: 6rpx;
}

.data-value-row {
  display: flex;
  align-items: baseline;
}

.data-number {
  font-size: 32rpx;
  font-weight: 800;
  color: #1E293B;
  letter-spacing: 0.5rpx;
}

.data-unit {
  font-size: 20rpx;
  font-weight: 600;
  color: #64748B;
  margin-left: 4rpx;
}

.data-date {
  font-size: 28rpx;
  font-weight: 700;
  color: #1E293B;
}

/*** 底部区域 ***/
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.budget-area {
  display: flex;
  align-items: baseline;
}

.budget-label {
  font-size: 24rpx;
  color: #94A3B8;
  font-weight: 500;
  margin-right: 12rpx;
}

.budget-display {
  display: flex;
  align-items: baseline;
}

.budget-symbol {
  font-size: 28rpx;
  font-weight: 700;
  color: #3B82F6;
  margin-right: 2rpx;
}

.budget-amount {
  font-size: 40rpx;
  font-weight: 800;
  background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/*** 操作按钮 ***/
.action-btn {
  position: relative;
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(145deg, #FFFFFF 0%, #F1F5F9 100%);
  border-radius: 18rpx;
  color: #3B82F6;
  transition: transform 0.3s ease;
  box-shadow: 0 4rpx 16rpx rgba(59, 130, 246, 0.15),
    inset 0 1rpx 2rpx rgba(255, 255, 255, 1);
  overflow: hidden;
}

.btn-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, transparent 50%);
}

/*** 3D悬浮阴影 ***/
.card-float-shadow {
  position: absolute;
  bottom: -16rpx;
  left: 40rpx;
  right: 40rpx;
  height: 32rpx;
  background: radial-gradient(ellipse 50% 100%,
      rgba(59, 130, 246, 0.15) 0%,
      rgba(59, 130, 246, 0.05) 40%,
      transparent 70%);
  border-radius: 50%;
  z-index: -1;
  transition: all 0.35s ease;
}
</style>
