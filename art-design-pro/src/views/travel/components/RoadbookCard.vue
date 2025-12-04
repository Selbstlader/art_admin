<template>
  <div class="roadbook-card" @click="handleClick">
    <ElSkeleton animated :loading="loading" style="width: 100%; height: 100%">
      <template #template>
        <div class="card-skeleton">
          <ElSkeletonItem
            variant="image"
            style="width: 100%; aspect-ratio: 16/10; border-radius: 10px"
          />
          <div style="padding: 12px 0">
            <ElSkeletonItem variant="p" style="width: 80%" />
            <ElSkeletonItem variant="p" style="width: 60%; margin-top: 10px" />
            <ElSkeletonItem variant="p" style="width: 40%; margin-top: 10px" />
          </div>
        </div>
      </template>

      <template #default>
        <!-- 封面图片 -->
        <div class="card-cover">
          <ElImage class="cover-image" :src="roadbook.coverUrl" lazy fit="cover">
            <template #error>
              <div class="image-placeholder">
                <ElIcon :size="32"><Picture /></ElIcon>
              </div>
            </template>
          </ElImage>
          <!-- 出行方式标签 -->
          <span class="travel-mode-tag">{{ travelModeText }}</span>
          <!-- 天数标签 -->
          <span v-if="dayCount > 0" class="days-tag">{{ dayCount }}天</span>
        </div>

        <!-- 卡片内容 -->
        <div class="card-content">
          <!-- 标题 -->
          <h3 class="card-title">{{ roadbook.title }}</h3>
          
          <!-- 描述 -->
          <p v-if="roadbook.description" class="card-desc">{{ roadbook.description }}</p>

          <!-- 标签 -->
          <div v-if="roadbook.tags && roadbook.tags.length > 0" class="card-tags">
            <ElTag
              v-for="tag in displayTags"
              :key="tag.id"
              size="small"
              :color="tag.color"
              effect="light"
              round
            >
              {{ tag.name }}
            </ElTag>
            <span v-if="roadbook.tags.length > 3" class="more-tags">
              +{{ roadbook.tags.length - 3 }}
            </span>
          </div>

          <!-- 统计信息 -->
          <div class="card-stats">
            <div class="stat-item">
              <ElIcon><View /></ElIcon>
              <span>{{ formatCount(roadbook.viewCount) }}</span>
            </div>
            <div class="stat-item">
              <ElIcon><Star /></ElIcon>
              <span>{{ formatCount(roadbook.favoriteCount) }}</span>
            </div>
            <div class="stat-item">
              <ElIcon><ChatDotRound /></ElIcon>
              <span>{{ formatCount(roadbook.commentCount) }}</span>
            </div>
          </div>

          <!-- 底部信息 -->
          <div class="card-footer">
            <div class="author-info" v-if="roadbook.author">
              <ElAvatar :size="24" :src="roadbook.author.avatar">
                {{ roadbook.author.nickname?.charAt(0) }}
              </ElAvatar>
              <span class="author-name">{{ roadbook.author.nickname }}</span>
            </div>
            <span class="create-time">{{ formatDate(roadbook.createdAt) }}</span>
          </div>
        </div>
      </template>
    </ElSkeleton>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Picture, View, Star, ChatDotRound } from '@element-plus/icons-vue'
import { useDateFormat } from '@vueuse/core'
import type { Roadbook } from '../types'
import { TravelMode } from '../types'

defineOptions({
  name: 'RoadbookCard'
})

interface Props {
  roadbook: Roadbook
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  click: [roadbook: Roadbook]
}>()

// 出行方式文本映射
const travelModeMap: Record<TravelMode, string> = {
  [TravelMode.DRIVING]: '自驾',
  [TravelMode.WALKING]: '步行',
  [TravelMode.CYCLING]: '骑行',
  [TravelMode.TRANSIT]: '公交'
}

// 计算出行方式文本
const travelModeText = computed(() => {
  return travelModeMap[props.roadbook.travelMode] || '自驾'
})

// 计算天数
const dayCount = computed(() => {
  if (!props.roadbook.startDate || !props.roadbook.endDate) return 0
  const start = new Date(props.roadbook.startDate)
  const end = new Date(props.roadbook.endDate)
  const diff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
  return diff + 1
})

// 显示的标签（最多3个）
const displayTags = computed(() => {
  return props.roadbook.tags?.slice(0, 3) || []
})

// 格式化数量
const formatCount = (count: number): string => {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + 'w'
  }
  if (count >= 1000) {
    return (count / 1000).toFixed(1) + 'k'
  }
  return String(count || 0)
}

// 格式化日期
const formatDate = (dateStr: string): string => {
  if (!dateStr) return ''
  return useDateFormat(dateStr, 'YYYY-MM-DD').value
}

// 点击事件
const handleClick = () => {
  emit('click', props.roadbook)
}
</script>


<style scoped lang="scss">
.roadbook-card {
  box-sizing: border-box;
  cursor: pointer;
  border: 1px solid var(--art-border-color);
  border-radius: calc(var(--custom-radius) / 2 + 2px);
  background: var(--art-main-bg-color);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  }

  .card-skeleton {
    padding: 0;
  }

  .card-cover {
    position: relative;
    aspect-ratio: 16/10;
    overflow: hidden;
    border-radius: calc(var(--custom-radius) / 2 + 2px) calc(var(--custom-radius) / 2 + 2px) 0 0;

    .cover-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
      background: var(--art-gray-200);
    }

    .image-placeholder {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 100%;
      background: var(--art-gray-200);
      color: var(--art-gray-400);
    }

    .travel-mode-tag {
      position: absolute;
      top: 8px;
      left: 8px;
      padding: 4px 8px;
      font-size: 12px;
      color: #fff;
      background: rgba(0, 0, 0, 0.6);
      border-radius: 4px;
    }

    .days-tag {
      position: absolute;
      top: 8px;
      right: 8px;
      padding: 4px 8px;
      font-size: 12px;
      color: #fff;
      background: var(--el-color-primary);
      border-radius: 4px;
    }
  }

  .card-content {
    padding: 12px;

    .card-title {
      margin: 0 0 8px;
      font-size: 16px;
      font-weight: 600;
      color: var(--art-text-gray-800);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .card-desc {
      margin: 0 0 8px;
      font-size: 13px;
      color: var(--art-text-gray-600);
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      line-height: 1.5;
    }

    .card-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
      margin-bottom: 10px;

      .el-tag {
        border: none;
      }

      .more-tags {
        font-size: 12px;
        color: var(--art-text-gray-500);
        line-height: 22px;
      }
    }

    .card-stats {
      display: flex;
      gap: 16px;
      margin-bottom: 10px;

      .stat-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: var(--art-text-gray-500);

        .el-icon {
          font-size: 14px;
        }
      }
    }

    .card-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-top: 10px;
      border-top: 1px solid var(--art-border-color);

      .author-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .author-name {
          font-size: 13px;
          color: var(--art-text-gray-600);
        }
      }

      .create-time {
        font-size: 12px;
        color: var(--art-text-gray-500);
      }
    }
  }
}
</style>
