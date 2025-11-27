<template>
  <div class="meditation-page">
    <!-- 继续播放 -->
    <div v-if="continuePlaying.length > 0" class="card art-custom-card continue-card">
      <div class="card-header">
        <div class="title">
          <h4 class="box-title">继续播放</h4>
          <p class="subtitle">未完成的冥想练习</p>
        </div>
      </div>
      <div class="continue-list">
        <div
          v-for="item in continuePlaying"
          :key="item.id"
          class="continue-item"
          @click="goToPlayer(item.content_id)"
        >
          <ElProgress type="circle" :percentage="item.progress" :width="50" :stroke-width="4" />
          <div class="continue-info">
            <div class="continue-title">{{ item.title }}</div>
            <div class="continue-meta">上次播放到 {{ formatDuration(item.duration) }}</div>
          </div>
          <ElButton type="primary" size="small" round>继续</ElButton>
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="card art-custom-card filter-card">
      <div class="filter-bar">
        <div class="filter-group">
          <span class="filter-label">分类:</span>
          <ElRadioGroup v-model="filters.category" @change="loadContent">
            <ElRadioButton label="">全部</ElRadioButton>
            <ElRadioButton
              v-for="(config, type) in MeditationTypeConfig"
              :key="type"
              :label="type"
            >
              {{ config.label }}
            </ElRadioButton>
          </ElRadioGroup>
        </div>
        <div class="filter-group">
          <span class="filter-label">难度:</span>
          <ElRadioGroup v-model="filters.difficulty" @change="loadContent">
            <ElRadioButton label="">全部</ElRadioButton>
            <ElRadioButton
              v-for="(config, level) in DifficultyConfig"
              :key="level"
              :label="level"
            >
              {{ config.label }}
            </ElRadioButton>
          </ElRadioGroup>
        </div>
      </div>
    </div>

    <!-- 内容列表 -->
    <ElRow :gutter="20">
      <ElCol
        v-for="content in contentList"
        :key="content.id"
        :xs="24"
        :sm="12"
        :md="8"
        :lg="6"
      >
        <div class="card art-custom-card content-card" @click="goToPlayer(content.id)">
          <div class="content-cover">
            <img :src="content.cover_image || defaultCover" :alt="content.title" />
            <div class="content-duration">{{ formatDuration(content.duration) }}</div>
            <ElButton
              class="favorite-btn"
              :icon="content.is_favorite ? StarFilled : Star"
              circle
              size="small"
              @click.stop="toggleFavorite(content)"
            />
          </div>
          <div class="content-info">
            <div class="content-title">{{ content.title }}</div>
            <div class="content-meta">
              <ElTag size="small" type="primary">
                {{ MeditationTypeConfig[content.category]?.label }}
              </ElTag>
              <ElTag size="small" :type="getDifficultyType(content.difficulty)">
                {{ DifficultyConfig[content.difficulty]?.label }}
              </ElTag>
            </div>
            <div class="content-stats">
              <span><i class="stat-icon">▶</i> {{ content.play_count }}</span>
              <span><i class="stat-icon">★</i> {{ content.favorite_count }}</span>
            </div>
          </div>
        </div>
      </ElCol>
    </ElRow>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-wrapper">
      <ElSkeleton :rows="3" animated />
    </div>

    <!-- 空状态 -->
    <ElEmpty v-if="!loading && contentList.length === 0" description="暂无冥想内容" :image-size="100" />

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination-wrapper">
      <ElPagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        background
        @current-change="loadContent"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { Star, StarFilled } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import {
    meditationApi,
    MeditationTypeConfig,
    DifficultyConfig,
    type MeditationContent,
    type DifficultyLevel
  } from '@/api/mood'

  defineOptions({ name: 'Meditation' })

  const router = useRouter()

  const defaultCover = 'https://via.placeholder.com/300x200?text=Meditation'

  // 筛选条件
  const filters = reactive({
    category: '',
    difficulty: ''
  })

  // 分页
  const currentPage = ref(1)
  const pageSize = ref(12)
  const total = ref(0)

  // 数据
  const contentList = ref<(MeditationContent & { is_favorite?: boolean })[]>([])
  const continuePlaying = ref<any[]>([])
  const loading = ref(false)

  // 方法
  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  const getDifficultyType = (level: DifficultyLevel): 'success' | 'warning' | 'danger' | 'info' => {
    const types: Record<DifficultyLevel, 'success' | 'warning' | 'danger'> = {
      beginner: 'success',
      intermediate: 'warning',
      advanced: 'danger'
    }
    return types[level] || 'info'
  }

  const goToPlayer = (contentId: number) => {
    router.push(`/mood/mood/meditation/player/${contentId}`)
  }

  const toggleFavorite = async (content: MeditationContent & { is_favorite?: boolean }) => {
    try {
      if (content.is_favorite) {
        await meditationApi.removeFavorite(content.id)
        content.is_favorite = false
        content.favorite_count--
        ElMessage.success('已取消收藏')
      } else {
        await meditationApi.addFavorite(content.id)
        content.is_favorite = true
        content.favorite_count++
        ElMessage.success('已添加收藏')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    }
  }

  const loadContent = async () => {
    loading.value = true
    try {
      const res = await meditationApi.getList({
        page: currentPage.value,
        pageSize: pageSize.value,
        category: filters.category || undefined,
        difficulty: filters.difficulty || undefined
      })
      const data = (res as any)?.data || res
      contentList.value = data?.records || data || []
      total.value = data?.total || contentList.value.length
    } catch (error) {
      console.error('加载内容失败:', error)
    } finally {
      loading.value = false
    }
  }

  const loadContinuePlaying = async () => {
    try {
      const res = await meditationApi.getContinuePlaying()
      continuePlaying.value = (res as any)?.data || res || []
    } catch (error) {
      console.error('加载继续播放失败:', error)
    }
  }

  onMounted(() => {
    loadContent()
    loadContinuePlaying()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/styles/variables.scss' as *;

  .meditation-page {
    --card-spacing: 20px;
  }

  // 统一卡片样式
  .card {
    margin-bottom: var(--card-spacing);
    background: var(--art-main-bg-color);
    border-radius: calc(var(--custom-radius) + 4px) !important;
  }

  // 卡片头部
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 20px 25px 15px;

    .title {
      h4 {
        font-size: 18px;
        font-weight: 500;
        color: var(--art-gray-900);
        margin: 0;
      }

      p {
        margin-top: 5px;
        font-size: 13px;
        color: var(--art-gray-600);
      }
    }
  }

  // 继续播放卡片
  .continue-card {
    padding: 0 25px 25px;

    .continue-list {
      display: flex;
      gap: 16px;
      overflow-x: auto;
      padding-bottom: 8px;

      &::-webkit-scrollbar {
        height: 6px;
      }

      &::-webkit-scrollbar-thumb {
        background: var(--art-gray-300);
        border-radius: 3px;
      }

      .continue-item {
        display: flex;
        align-items: center;
        gap: 14px;
        padding: 14px;
        background: var(--art-gray-100);
        border-radius: 12px;
        cursor: pointer;
        min-width: 300px;
        transition: all 0.3s ease;

        &:hover {
          background: rgb(var(--art-bg-primary));
          transform: translateX(2px);
        }

        .continue-info {
          flex: 1;

          .continue-title {
            font-size: 14px;
            font-weight: 500;
            color: var(--art-gray-800);
            margin-bottom: 4px;
          }

          .continue-meta {
            font-size: 12px;
            color: var(--art-gray-500);
          }
        }
      }
    }
  }

  // 筛选卡片
  .filter-card {
    padding: 20px 25px;

    .filter-bar {
      display: flex;
      flex-direction: column;
      gap: 16px;

      .filter-group {
        display: flex;
        align-items: center;
        gap: 12px;
        flex-wrap: wrap;

        .filter-label {
          font-size: 14px;
          font-weight: 500;
          color: var(--art-gray-700);
          min-width: 50px;
        }
      }
    }
  }

  // 内容卡片
  .content-card {
    cursor: pointer;
    transition: all 0.3s ease;
    overflow: hidden;
    padding: 0;

    &:hover {
      transform: translateY(-4px);
      box-shadow: var(--art-box-shadow);
    }

    .content-cover {
      position: relative;
      aspect-ratio: 16/10;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s ease;
      }

      &:hover img {
        transform: scale(1.05);
      }

      .content-duration {
        position: absolute;
        bottom: 10px;
        right: 10px;
        background: rgba(0, 0, 0, 0.75);
        color: #fff;
        padding: 4px 10px;
        border-radius: 6px;
        font-size: 12px;
        font-weight: 500;
      }

      .favorite-btn {
        position: absolute;
        top: 10px;
        right: 10px;
        background: rgba(255, 255, 255, 0.95);
        border: none;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

        &:hover {
          background: #fff;
          transform: scale(1.1);
        }
      }
    }

    .content-info {
      padding: 14px;

      .content-title {
        font-size: 14px;
        font-weight: 500;
        color: var(--art-gray-800);
        margin-bottom: 10px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .content-meta {
        display: flex;
        gap: 8px;
        margin-bottom: 10px;
      }

      .content-stats {
        display: flex;
        gap: 16px;
        font-size: 12px;
        color: var(--art-gray-500);

        .stat-icon {
          font-style: normal;
          margin-right: 4px;
          font-size: 10px;
        }
      }
    }
  }

  // 加载状态
  .loading-wrapper {
    padding: 40px;
    background: var(--art-main-bg-color);
    border-radius: calc(var(--custom-radius) + 4px);
  }

  // 分页
  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 24px;
    padding: 20px;
    background: var(--art-main-bg-color);
    border-radius: calc(var(--custom-radius) + 4px);
  }

  // 响应式
  @media screen and (max-width: $device-phone) {
    .meditation-page {
      --card-spacing: 15px;
    }

    .continue-card .continue-list .continue-item {
      min-width: 260px;
    }
  }
</style>
