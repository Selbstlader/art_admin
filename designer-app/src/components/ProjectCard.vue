<template>
  <!-- 
    项目卡片组件 - C4D风格高级UI
    Project Card Component - C4D style premium UI with 3D floating effect
  -->
  <view class="project-card" @click="handleClick">
    <!-- 卡片左侧装饰条 - Card left accent bar with gradient -->
    <view :class="['card-accent', `accent-${project.status}`]"></view>
    
    <view class="card-content">
      <!-- 卡片头部 - Card header with title and status -->
      <view class="card-header">
        <text class="project-name">{{ project.name }}</text>
        <view :class="['status-badge', `status-${project.status}`]">
          <view class="status-dot"></view>
          <text class="status-text">{{ statusText }}</text>
        </view>
      </view>
      
      <!-- 卡片数据区 - Card data area with icons -->
      <view class="card-data">
        <view class="data-item">
          <view class="data-icon">
            <!-- Area icon -->
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect width="18" height="18" x="3" y="3" rx="2"></rect>
              <path d="M3 9h18"></path>
              <path d="M9 21V9"></path>
            </svg>
          </view>
          <view class="data-content">
            <text class="data-label">面积</text>
            <text class="data-value">{{ project.area }}<text class="data-unit">㎡</text></text>
          </view>
        </view>
        
        <view class="data-divider"></view>
        
        <view class="data-item">
          <view class="data-icon">
            <!-- Time icon -->
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </view>
          <view class="data-content">
            <text class="data-label">更新时间</text>
            <text class="data-value">{{ formattedDate }}</text>
          </view>
        </view>
      </view>
      
      <!-- 卡片底部 - Card footer with budget info -->
      <view class="card-footer">
        <view class="budget-info">
          <text class="budget-label">预算</text>
          <text class="budget-value">¥{{ formattedBudget }}</text>
        </view>
        <view class="arrow-icon">
          <!-- Arrow right icon -->
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m9 18 6-6-6-6"></path>
          </svg>
        </view>
      </view>
    </view>
    
    <!-- 卡片3D阴影效果 - Card 3D shadow effect -->
    <view class="card-shadow"></view>
  </view>
</template>

<script setup lang="ts">
/*** 
 * ProjectCard Component
 * 项目卡片组件 - 悬浮卡片设计，蓝色渐变背景
 * Floating card design with blue gradient background
 ***/
import { computed } from 'vue'
import type { Project } from '@/types/api'
import { ProjectStatusText, ProjectStatus } from '@/types/common'

/*** Component Props - 组件属性定义 ***/
interface Props {
  project: Project
}

const props = defineProps<Props>()

/*** Component Emits - 组件事件定义 ***/
const emit = defineEmits<{
  (e: 'click', id: number): void
}>()

/*** 状态文本 - Status text mapping ***/
const statusText = computed((): string => {
  return ProjectStatusText[props.project.status as ProjectStatus] || props.project.status
})

/*** 格式化日期 - Format date to YYYY-MM-DD ***/
const formattedDate = computed((): string => {
  if (!props.project.updatedAt) return '-'
  const date = new Date(props.project.updatedAt)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
})

/*** 格式化预算 - Format budget with unit ***/
const formattedBudget = computed((): string => {
  const budget = props.project.budget
  if (budget >= 10000) {
    return (budget / 10000).toFixed(1) + '万'
  }
  return budget.toLocaleString()
})

/*** 处理点击事件 - Handle card click event ***/
const handleClick = (): void => {
  emit('click', props.project.id)
}
</script>

<style lang="scss" scoped>
/*** 项目卡片 - Project card with 3D floating effect ***/
.project-card {
  position: relative;
  background: #FFFFFF;
  border-radius: 24rpx;
  overflow: hidden;
  display: flex;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

/*** 卡片左侧装饰条 - Card left accent bar ***/
.card-accent {
  width: 8rpx;
  flex-shrink: 0;
}

.accent-draft {
  background: linear-gradient(180deg, #94A3B8 0%, #64748B 100%);
}

.accent-in_progress {
  background: linear-gradient(180deg, #60A5FA 0%, #3B82F6 100%);
}

.accent-completed {
  background: linear-gradient(180deg, #34D399 0%, #10B981 100%);
}

.accent-archived {
  background: linear-gradient(180deg, #F87171 0%, #EF4444 100%);
}

/*** 卡片内容区 - Card content area ***/
.card-content {
  flex: 1;
  padding: 28rpx 28rpx 24rpx;
}

/*** 卡片3D阴影 - Card 3D shadow ***/
.card-shadow {
  position: absolute;
  bottom: -8rpx;
  left: 24rpx;
  right: 24rpx;
  height: 16rpx;
  background: radial-gradient(ellipse, rgba(0, 0, 0, 0.08) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
}

/*** 卡片头部 - Card header ***/
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20rpx;
}

.project-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #1E293B;
  flex: 1;
  line-height: 1.4;
  margin-right: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/*** 状态徽章 - Status badge ***/
.status-badge {
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
}

.status-draft {
  background: rgba(148, 163, 184, 0.15);
}

.status-in_progress {
  background: rgba(59, 130, 246, 0.12);
}

.status-completed {
  background: rgba(16, 185, 129, 0.12);
}

.status-archived {
  background: rgba(239, 68, 68, 0.12);
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  margin-right: 8rpx;
}

.status-draft .status-dot {
  background: #64748B;
}

.status-in_progress .status-dot {
  background: #3B82F6;
}

.status-completed .status-dot {
  background: #10B981;
}

.status-archived .status-dot {
  background: #EF4444;
}

.status-text {
  font-size: 22rpx;
  font-weight: 500;
}

.status-draft .status-text {
  color: #64748B;
}

.status-in_progress .status-text {
  color: #3B82F6;
}

.status-completed .status-text {
  color: #10B981;
}

.status-archived .status-text {
  color: #EF4444;
}

/*** 卡片数据区 - Card data area ***/
.card-data {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-top: 1rpx solid #F1F5F9;
  border-bottom: 1rpx solid #F1F5F9;
  margin-bottom: 16rpx;
}

.data-item {
  flex: 1;
  display: flex;
  align-items: center;
}

.data-icon {
  width: 40rpx;
  height: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-right: 12rpx;
}

.data-content {
  display: flex;
  flex-direction: column;
}

.data-label {
  font-size: 22rpx;
  color: #94A3B8;
  margin-bottom: 4rpx;
}

.data-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #1E293B;
}

.data-unit {
  font-size: 22rpx;
  font-weight: 400;
  color: #64748B;
}

.data-divider {
  width: 1rpx;
  height: 48rpx;
  background: #E2E8F0;
  margin: 0 24rpx;
}

/*** 卡片底部 - Card footer ***/
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.budget-info {
  display: flex;
  align-items: baseline;
}

.budget-label {
  font-size: 22rpx;
  color: #94A3B8;
  margin-right: 8rpx;
}

.budget-value {
  font-size: 30rpx;
  font-weight: 700;
  color: #3B82F6;
}

.arrow-icon {
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #CBD5E1;
  transition: transform 0.2s ease;
}

.project-card:active .arrow-icon {
  transform: translateX(4rpx);
}
</style>
