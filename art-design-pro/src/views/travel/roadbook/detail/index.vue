<template>
  <div class="page-content roadbook-detail">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <ElSkeleton :rows="10" animated />
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <ElResult icon="error" :title="error">
        <template #extra>
          <ElButton type="primary" @click="handleBack">返回列表</ElButton>
        </template>
      </ElResult>
    </div>

    <!-- 路书详情内容 -->
    <template v-else-if="roadbook">
      <!-- 顶部操作栏 -->
      <div class="detail-header">
        <div class="header-left">
          <ElButton :icon="ArrowLeft" @click="handleBack">返回</ElButton>
        </div>
        <div class="header-right">
          <ElButton :icon="Share" @click="handleShare">分享</ElButton>
          <ElButton v-if="isOwner" type="primary" :icon="Edit" @click="handleEdit">编辑</ElButton>
        </div>
      </div>

      <!-- 主体内容 -->
      <div class="detail-body">
        <ElRow :gutter="20">
          <!-- 左侧：路书信息 -->
          <ElCol :lg="14" :md="24">
            <div class="detail-left">
              <!-- 封面和基本信息 -->
              <ElCard class="info-card" shadow="never">
                <!-- 封面图片 -->
                <div v-if="roadbook.coverUrl" class="cover-image">
                  <ElImage :src="roadbook.coverUrl" fit="cover" />
                </div>

                <!-- 标题和作者 -->
                <div class="roadbook-header">
                  <h1 class="roadbook-title">{{ roadbook.title }}</h1>
                  <div class="author-info" v-if="roadbook.author">
                    <ElAvatar :size="32" :src="roadbook.author.avatar">
                      {{ roadbook.author.nickname?.charAt(0) }}
                    </ElAvatar>
                    <span class="author-name">{{ roadbook.author.nickname }}</span>
                    <span class="publish-time">{{ formatDate(roadbook.createdAt) }}</span>
                  </div>
                </div>

                <!-- 标签 -->
                <div v-if="roadbook.tags && roadbook.tags.length > 0" class="tag-list">
                  <ElTag
                    v-for="tag in roadbook.tags"
                    :key="tag.id"
                    :color="tag.color"
                    size="small"
                    effect="light"
                  >
                    {{ tag.name }}
                  </ElTag>
                </div>

                <!-- 统计数据 -->
                <div class="stats-row">
                  <div class="stat-item">
                    <ElIcon><View /></ElIcon>
                    <span>{{ roadbook.viewCount || 0 }}</span>
                  </div>
                  <div class="stat-item">
                    <ElIcon><Star /></ElIcon>
                    <span>{{ roadbook.favoriteCount || 0 }}</span>
                  </div>
                  <div class="stat-item">
                    <ElIcon><ChatDotRound /></ElIcon>
                    <span>{{ roadbook.commentCount || 0 }}</span>
                  </div>
                  <div class="stat-item">
                    <ElIcon><Share /></ElIcon>
                    <span>{{ roadbook.shareCount || 0 }}</span>
                  </div>
                </div>

                <!-- 操作按钮 -->
                <div class="action-buttons">
                  <ElButton
                    :type="isFavorited ? 'warning' : 'default'"
                    :icon="isFavorited ? StarFilled : Star"
                    @click="handleToggleFavorite"
                    :loading="favoriteLoading"
                  >
                    {{ isFavorited ? '已收藏' : '收藏' }}
                  </ElButton>
                  <ElButton
                    :type="isLiked ? 'danger' : 'default'"
                    @click="handleToggleLike"
                    :loading="likeLoading"
                  >
                    <template #icon>
                      <svg v-if="isLiked" class="like-icon filled" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                      </svg>
                      <svg v-else class="like-icon" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M16.5 3c-1.74 0-3.41.81-4.5 2.09C10.91 3.81 9.24 3 7.5 3 4.42 3 2 5.42 2 8.5c0 3.78 3.4 6.86 8.55 11.54L12 21.35l1.45-1.32C18.6 15.36 22 12.28 22 8.5 22 5.42 19.58 3 16.5 3zm-4.4 15.55l-.1.1-.1-.1C7.14 14.24 4 11.39 4 8.5 4 6.5 5.5 5 7.5 5c1.54 0 3.04.99 3.57 2.36h1.87C13.46 5.99 14.96 5 16.5 5c2 0 3.5 1.5 3.5 3.5 0 2.89-3.14 5.74-7.9 10.05z"/>
                      </svg>
                    </template>
                    {{ isLiked ? '已点赞' : '点赞' }} {{ roadbook.likeCount || 0 }}
                  </ElButton>
                </div>

                <!-- 描述 -->
                <div v-if="roadbook.description" class="roadbook-description">
                  <p>{{ roadbook.description }}</p>
                </div>

                <!-- 行程信息 -->
                <div class="trip-info">
                  <div class="info-row">
                    <div class="info-item">
                      <span class="label">行程日期</span>
                      <span class="value">
                        {{ roadbook.startDate ? formatDate(roadbook.startDate) : '未设置' }}
                        <template v-if="roadbook.endDate">
                          ~ {{ formatDate(roadbook.endDate) }}
                        </template>
                      </span>
                    </div>
                    <div class="info-item">
                      <span class="label">行程天数</span>
                      <span class="value">{{ dayCount }} 天</span>
                    </div>
                  </div>
                  <div class="info-row">
                    <div class="info-item">
                      <span class="label">出行方式</span>
                      <span class="value">{{ getTravelModeText(roadbook.travelMode) }}</span>
                    </div>
                    <div class="info-item">
                      <span class="label">预计预算</span>
                      <span class="value">¥{{ roadbook.totalBudget || 0 }}</span>
                    </div>
                  </div>
                  <div class="info-row">
                    <div class="info-item">
                      <span class="label">总距离</span>
                      <span class="value">{{ formattedDistance || '计算中...' }}</span>
                    </div>
                    <div class="info-item">
                      <span class="label">预计时间</span>
                      <span class="value">{{ formattedDuration || '计算中...' }}</span>
                    </div>
                  </div>
                </div>
              </ElCard>

              <!-- 日程时间线 -->
              <ElCard class="timeline-card" shadow="never">
                <template #header>
                  <div class="card-header">
                    <span>行程安排</span>
                    <span class="waypoint-count">共 {{ roadbook.waypoints?.length || 0 }} 个地点</span>
                  </div>
                </template>
                <ScheduleTimeline
                  :waypoints="roadbook.waypoints || []"
                  :days="dayCount"
                  :start-date="roadbook.startDate"
                />
              </ElCard>
            </div>
          </ElCol>

          <!-- 右侧：地图 -->
          <ElCol :lg="10" :md="24">
            <div class="detail-right">
              <!-- 地图卡片 -->
              <ElCard class="map-card" shadow="never">
                <template #header>
                  <div class="card-header">
                    <span>路线地图</span>
                    <ElButton
                      type="primary"
                      link
                      :icon="Position"
                      @click="handleExportNav"
                    >
                      导出导航
                    </ElButton>
                  </div>
                </template>
                <div class="map-container">
                  <AMapContainer
                    ref="mapContainerRef"
                    :center="mapCenter"
                    :zoom="12"
                    @ready="handleMapReady"
                  />
                </div>
              </ElCard>

              <!-- 途经点列表 -->
              <ElCard class="waypoints-card" shadow="never">
                <template #header>
                  <div class="card-header">
                    <span>途经点列表</span>
                  </div>
                </template>
                <div class="waypoint-list">
                  <div
                    v-for="(waypoint, index) in sortedWaypoints"
                    :key="waypoint.id || index"
                    class="waypoint-item"
                    @click="handleWaypointClick(waypoint)"
                  >
                    <div class="waypoint-index" :class="getWaypointClass(index)">
                      {{ index + 1 }}
                    </div>
                    <div class="waypoint-info">
                      <div class="waypoint-name">{{ waypoint.name }}</div>
                      <div class="waypoint-address">{{ waypoint.address }}</div>
                    </div>
                    <div class="waypoint-day">
                      第{{ waypoint.dayIndex || 1 }}天
                    </div>
                  </div>
                </div>
              </ElCard>
            </div>
          </ElCol>
        </ElRow>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Share,
  Edit,
  View,
  Star,
  StarFilled,
  ChatDotRound,
  Position
} from '@element-plus/icons-vue'
import { useDateFormat } from '@vueuse/core'
import AMapContainer from '../../components/AMapContainer.vue'
import ScheduleTimeline from '../../components/ScheduleTimeline.vue'
import { useRoadbookStore } from '../../store/roadbook'
import { useAMap } from '../../composables/useAMap'
import { useRoute as useRouteHook } from '../../composables/useRoute'
import type { Waypoint, TravelMode } from '../../types'

defineOptions({
  name: 'RoadbookDetail'
})

const router = useRouter()
const route = useRoute()
const roadbookStore = useRoadbookStore()

// 路书ID
const roadbookId = computed(() => Number(route.params.id) || 0)

// 状态
const loading = ref(true)
const error = ref<string | null>(null)
const favoriteLoading = ref(false)
const likeLoading = ref(false)

// 路书数据
const roadbook = computed(() => roadbookStore.currentRoadbook)
const isFavorited = computed(() => roadbook.value?.isFavorited ?? false)
const isLiked = computed(() => roadbook.value?.isLiked ?? false)

// 是否是路书所有者（简化判断，实际应该从用户store获取当前用户ID）
const isOwner = computed(() => {
  // TODO: 从用户store获取当前用户ID进行比较
  return false
})

// 计算天数
const dayCount = computed(() => {
  if (!roadbook.value?.startDate || !roadbook.value?.endDate) return 1
  const start = new Date(roadbook.value.startDate)
  const end = new Date(roadbook.value.endDate)
  const diff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
  return Math.max(diff + 1, 1)
})

// 排序后的途经点
const sortedWaypoints = computed(() => {
  if (!roadbook.value?.waypoints) return []
  return [...roadbook.value.waypoints].sort((a, b) => {
    if (a.dayIndex !== b.dayIndex) return (a.dayIndex || 1) - (b.dayIndex || 1)
    return (a.sortOrder || 0) - (b.sortOrder || 0)
  })
})

// 地图相关
const mapContainerRef = ref<InstanceType<typeof AMapContainer> | null>(null)
const mapCenter = ref<[number, number]>([116.397428, 39.90923])
const { setMarkersFromWaypoints, drawPolyline, clearPolylines, panTo } = useAMap()
const { planRouteFromWaypoints, formattedDistance, formattedDuration } = useRouteHook()

// 格式化日期
const formatDate = (dateStr: string): string => {
  if (!dateStr) return ''
  return useDateFormat(new Date(dateStr), 'YYYY-MM-DD').value
}

// 获取出行方式文本
const getTravelModeText = (mode: TravelMode | string): string => {
  const modeMap: Record<string, string> = {
    driving: '自驾',
    walking: '步行',
    cycling: '骑行',
    transit: '公交'
  }
  return modeMap[mode] || '自驾'
}

// 获取途经点样式类
const getWaypointClass = (index: number): string => {
  const total = sortedWaypoints.value.length
  if (index === 0) return 'start'
  if (index === total - 1) return 'end'
  return 'middle'
}

// 地图就绪
const handleMapReady = async () => {
  if (roadbook.value?.waypoints && roadbook.value.waypoints.length > 0) {
    await updateMapDisplay()
  }
}

// 更新地图显示
const updateMapDisplay = async () => {
  const waypoints = roadbook.value?.waypoints
  if (!waypoints || waypoints.length === 0) return

  // 设置标记
  const validWaypoints = waypoints.filter(
    (wp): wp is Waypoint => !!(wp.longitude && wp.latitude)
  )
  setMarkersFromWaypoints(validWaypoints)

  // 设置地图中心
  if (validWaypoints.length > 0) {
    mapCenter.value = [validWaypoints[0].longitude, validWaypoints[0].latitude]
  }

  // 规划路线
  if (validWaypoints.length >= 2) {
    try {
      const result = await planRouteFromWaypoints(
        validWaypoints,
        roadbook.value?.travelMode as TravelMode
      )
      if (result && mapContainerRef.value) {
        clearPolylines()
        const map = mapContainerRef.value.getMap()
        if (map) {
          drawPolyline(result.polyline)
          map.setFitView()
        }
      }
    } catch (err) {
      console.error('Route planning failed:', err)
    }
  }
}

// 点击途经点
const handleWaypointClick = (waypoint: Waypoint) => {
  if (waypoint.longitude && waypoint.latitude) {
    panTo([waypoint.longitude, waypoint.latitude])
  }
}

// 返回
const handleBack = () => {
  router.back()
}

// 编辑
const handleEdit = () => {
  router.push({ name: 'RoadbookEditor', params: { id: roadbookId.value } })
}

// 分享
const handleShare = () => {
  // TODO: 实现分享功能
  ElMessage.info('分享功能开发中')
}

// 导出导航
const handleExportNav = () => {
  // TODO: 实现导出到高德导航功能
  ElMessage.info('导航导出功能开发中')
}

// 切换收藏
const handleToggleFavorite = async () => {
  if (!roadbook.value) return
  
  favoriteLoading.value = true
  try {
    if (isFavorited.value) {
      await roadbookStore.unfavoriteRoadbook(roadbookId.value)
      ElMessage.success('已取消收藏')
    } else {
      await roadbookStore.favoriteRoadbook(roadbookId.value)
      ElMessage.success('收藏成功')
    }
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    favoriteLoading.value = false
  }
}

// 切换点赞
const handleToggleLike = async () => {
  if (!roadbook.value) return
  
  likeLoading.value = true
  try {
    if (isLiked.value) {
      await roadbookStore.unlikeRoadbook(roadbookId.value)
      ElMessage.success('已取消点赞')
    } else {
      await roadbookStore.likeRoadbook(roadbookId.value)
      ElMessage.success('点赞成功')
    }
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    likeLoading.value = false
  }
}

// 加载路书详情
const loadRoadbook = async () => {
  loading.value = true
  error.value = null

  try {
    await roadbookStore.fetchRoadbookDetail(roadbookId.value)
    if (!roadbook.value) {
      error.value = '路书不存在或已被删除'
    }
  } catch (err: any) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

// 监听路书数据变化，更新地图
watch(() => roadbook.value?.waypoints, () => {
  if (mapContainerRef.value?.getMap()) {
    updateMapDisplay()
  }
}, { deep: true })

onMounted(() => {
  loadRoadbook()
})
</script>


<style scoped lang="scss">
.roadbook-detail {
  .loading-container,
  .error-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--art-border-color);

    .header-right {
      display: flex;
      gap: 12px;
    }
  }

  .detail-body {
    .detail-left,
    .detail-right {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }

    .el-card {
      border-radius: calc(var(--custom-radius) / 2 + 2px);

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-weight: 600;

        .waypoint-count {
          font-size: 13px;
          font-weight: normal;
          color: var(--art-text-gray-500);
        }
      }
    }

    // 信息卡片
    .info-card {
      .cover-image {
        margin: -20px -20px 20px -20px;
        height: 300px;
        overflow: hidden;

        .el-image {
          width: 100%;
          height: 100%;
        }
      }

      .roadbook-header {
        margin-bottom: 16px;

        .roadbook-title {
          margin: 0 0 12px 0;
          font-size: 24px;
          font-weight: 700;
          color: var(--art-text-gray-800);
          line-height: 1.4;
        }

        .author-info {
          display: flex;
          align-items: center;
          gap: 8px;

          .author-name {
            font-size: 14px;
            font-weight: 500;
            color: var(--art-text-gray-700);
          }

          .publish-time {
            font-size: 13px;
            color: var(--art-text-gray-500);
          }
        }
      }

      .tag-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-bottom: 16px;
      }

      .stats-row {
        display: flex;
        gap: 24px;
        padding: 12px 0;
        border-top: 1px solid var(--art-border-color);
        border-bottom: 1px solid var(--art-border-color);
        margin-bottom: 16px;

        .stat-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 14px;
          color: var(--art-text-gray-600);

          .el-icon {
            font-size: 16px;
          }
        }
      }

      .action-buttons {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;

        .like-icon {
          width: 16px;
          height: 16px;

          &.filled {
            color: var(--el-color-danger);
          }
        }
      }

      .roadbook-description {
        margin-bottom: 16px;
        padding: 12px 16px;
        background: var(--art-gray-100);
        border-radius: 8px;

        p {
          margin: 0;
          font-size: 14px;
          line-height: 1.8;
          color: var(--art-text-gray-700);
          white-space: pre-wrap;
        }
      }

      .trip-info {
        .info-row {
          display: flex;
          gap: 24px;
          margin-bottom: 12px;

          &:last-child {
            margin-bottom: 0;
          }

          .info-item {
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 4px;

            .label {
              font-size: 12px;
              color: var(--art-text-gray-500);
            }

            .value {
              font-size: 15px;
              font-weight: 500;
              color: var(--art-text-gray-800);
            }
          }
        }
      }
    }

    // 地图卡片
    .map-card {
      position: sticky;
      top: 20px;

      .map-container {
        height: 400px;
        border-radius: 8px;
        overflow: hidden;
      }
    }

    // 途经点列表卡片
    .waypoints-card {
      .waypoint-list {
        max-height: 400px;
        overflow-y: auto;

        .waypoint-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 12px;
          border-radius: 8px;
          cursor: pointer;
          transition: background 0.2s;

          &:hover {
            background: var(--art-gray-100);
          }

          .waypoint-index {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 28px;
            height: 28px;
            border-radius: 50%;
            font-size: 12px;
            font-weight: 600;
            color: #fff;
            flex-shrink: 0;

            &.start {
              background: var(--el-color-success);
            }

            &.end {
              background: var(--el-color-danger);
            }

            &.middle {
              background: var(--el-color-primary);
            }
          }

          .waypoint-info {
            flex: 1;
            min-width: 0;

            .waypoint-name {
              font-size: 14px;
              font-weight: 500;
              color: var(--art-text-gray-800);
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            }

            .waypoint-address {
              margin-top: 2px;
              font-size: 12px;
              color: var(--art-text-gray-500);
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            }
          }

          .waypoint-day {
            font-size: 12px;
            color: var(--art-text-gray-500);
            flex-shrink: 0;
          }
        }
      }
    }
  }
}

// 响应式布局
@media only screen and (max-width: 1200px) {
  .roadbook-detail {
    .detail-body {
      .el-col {
        margin-bottom: 20px;
      }

      .map-card {
        position: static;
      }
    }
  }
}

@media only screen and (max-width: 768px) {
  .roadbook-detail {
    .detail-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;

      .header-right {
        width: 100%;
        justify-content: flex-end;
      }
    }

    .detail-body {
      .info-card {
        .cover-image {
          height: 200px;
        }

        .roadbook-header {
          .roadbook-title {
            font-size: 20px;
          }
        }

        .stats-row {
          flex-wrap: wrap;
          gap: 16px;
        }

        .action-buttons {
          flex-wrap: wrap;
        }

        .trip-info {
          .info-row {
            flex-direction: column;
            gap: 12px;
          }
        }
      }

      .map-card {
        .map-container {
          height: 300px;
        }
      }
    }
  }
}
</style>
