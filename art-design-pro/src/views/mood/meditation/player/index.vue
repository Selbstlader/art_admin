<template>
  <div class="meditation-player">
    <div v-if="content" class="player-container">
      <!-- 封面区域 -->
      <div class="cover-area">
        <img :src="content.cover_image || defaultCover" :alt="content.title" class="cover-image" />
        <div class="cover-overlay">
          <div class="content-title">{{ content.title }}</div>
          <div class="content-meta">
            {{ MeditationTypeConfig[content.category]?.label }} ·
            {{ DifficultyConfig[content.difficulty]?.label }}
          </div>
        </div>
      </div>

      <!-- 播放控制区域 -->
      <div class="control-area">
        <!-- 进度条 -->
        <div class="progress-wrapper">
          <ElSlider
            v-model="currentProgress"
            :max="content.duration"
            :show-tooltip="false"
            @change="handleSeek"
          />
          <div class="time-display">
            <span>{{ formatDuration(currentTime) }}</span>
            <span>{{ formatDuration(content.duration) }}</span>
          </div>
        </div>

        <!-- 控制按钮 -->
        <div class="control-buttons">
          <ElButton circle size="large" @click="seekBackward">
            <span class="control-icon">-15s</span>
          </ElButton>
          <ElButton type="primary" circle size="large" class="play-btn" @click="togglePlay">
            <ElIcon :size="32">
              <VideoPlay v-if="!isPlaying" />
              <VideoPause v-else />
            </ElIcon>
          </ElButton>
          <ElButton circle size="large" @click="seekForward">
            <span class="control-icon">+15s</span>
          </ElButton>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <ElButton :icon="isFavorite ? StarFilled : Star" circle @click="toggleFavorite" />
          <ElButton :icon="Back" circle @click="goBack" />
        </div>
      </div>

      <!-- 描述信息 -->
      <div class="description-area">
        <ElCard shadow="never">
          <template #header>
            <span>内容介绍</span>
          </template>
          <p>{{ content.description || '暂无介绍' }}</p>
        </ElCard>
      </div>

      <!-- 隐藏的音频元素 -->
      <audio
        ref="audioRef"
        :src="content.audio_url"
        @timeupdate="handleTimeUpdate"
        @ended="handleEnded"
        @loadedmetadata="handleLoaded"
      />
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-wrapper">
      <ElSkeleton :rows="5" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { VideoPlay, VideoPause, Star, StarFilled, Back } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import {
    meditationApi,
    MeditationTypeConfig,
    DifficultyConfig,
    type MeditationContent
  } from '@/api/mood'

  defineOptions({ name: 'MeditationPlayer' })

  const route = useRoute()
  const router = useRouter()

  const defaultCover = 'https://via.placeholder.com/400x300?text=Meditation'

  // 内容数据
  const content = ref<MeditationContent | null>(null)
  const isFavorite = ref(false)

  // 播放状态
  const audioRef = ref<HTMLAudioElement | null>(null)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const currentProgress = ref(0)

  // 播放记录
  const playRecord = ref<any>(null)

  // 格式化时长
  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = Math.floor(seconds % 60)
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  // 播放/暂停
  const togglePlay = () => {
    if (!audioRef.value) return

    if (isPlaying.value) {
      audioRef.value.pause()
    } else {
      audioRef.value.play()
    }
    isPlaying.value = !isPlaying.value
  }

  // 快退15秒
  const seekBackward = () => {
    if (!audioRef.value) return
    audioRef.value.currentTime = Math.max(0, audioRef.value.currentTime - 15)
  }

  // 快进15秒
  const seekForward = () => {
    if (!audioRef.value || !content.value) return
    audioRef.value.currentTime = Math.min(content.value.duration, audioRef.value.currentTime + 15)
  }

  // 拖动进度条
  const handleSeek = (value: number) => {
    if (!audioRef.value) return
    audioRef.value.currentTime = value
    currentTime.value = value
  }

  // 时间更新
  const handleTimeUpdate = () => {
    if (!audioRef.value) return
    currentTime.value = audioRef.value.currentTime
    currentProgress.value = audioRef.value.currentTime
  }

  // 播放结束
  const handleEnded = () => {
    isPlaying.value = false
    savePlayRecord(true)
  }

  // 音频加载完成
  const handleLoaded = () => {
    // 恢复上次播放位置
    if (playRecord.value && audioRef.value) {
      const lastPosition = (playRecord.value.progress / 100) * (content.value?.duration || 0)
      audioRef.value.currentTime = lastPosition
      currentTime.value = lastPosition
      currentProgress.value = lastPosition
    }
  }

  // 收藏/取消收藏
  const toggleFavorite = async () => {
    if (!content.value) return

    try {
      if (isFavorite.value) {
        await meditationApi.removeFavorite(content.value.id)
        isFavorite.value = false
        ElMessage.success('已取消收藏')
      } else {
        await meditationApi.addFavorite(content.value.id)
        isFavorite.value = true
        ElMessage.success('已添加收藏')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    }
  }

  // 返回
  const goBack = () => {
    router.back()
  }

  // 保存播放记录
  const savePlayRecord = async (completed = false) => {
    if (!content.value) return

    const progress =
      content.value.duration > 0 ? (currentTime.value / content.value.duration) * 100 : 0

    try {
      await meditationApi.updatePlayRecord({
        content_id: content.value.id,
        duration: Math.floor(currentTime.value),
        progress: Math.min(progress, 100),
        completed
      })
    } catch (error) {
      console.error('保存播放记录失败:', error)
    }
  }

  // 加载内容
  const loadContent = async () => {
    const id = Number(route.params.id)
    if (!id) return

    try {
      const [contentRes, favoriteRes, recordRes] = await Promise.all([
        meditationApi.getDetail(id),
        meditationApi.checkFavorite(id),
        meditationApi.getPlayRecord(id)
      ])

      content.value = contentRes as MeditationContent
      isFavorite.value = (favoriteRes as any)?.is_favorite || false
      playRecord.value = recordRes || null
    } catch (error) {
      console.error('加载内容失败:', error)
      ElMessage.error('加载内容失败')
    }
  }

  // 定时保存播放记录
  let saveTimer: NodeJS.Timeout | null = null

  onMounted(() => {
    loadContent()

    // 每30秒保存一次播放记录
    saveTimer = setInterval(() => {
      if (isPlaying.value) {
        savePlayRecord()
      }
    }, 30000)
  })

  onUnmounted(() => {
    // 离开页面时保存播放记录
    if (currentTime.value > 0) {
      savePlayRecord()
    }

    if (saveTimer) {
      clearInterval(saveTimer)
    }
  })

  // 监听路由变化
  watch(
    () => route.params.id,
    () => {
      loadContent()
    }
  )
</script>

<style scoped lang="scss">
  .meditation-player {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .player-container {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
  }

  .cover-area {
    position: relative;
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 30px;

    .cover-image {
      width: 100%;
      aspect-ratio: 4/3;
      object-fit: cover;
    }

    .cover-overlay {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      padding: 20px;
      background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
      color: #fff;

      .content-title {
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 8px;
      }

      .content-meta {
        font-size: 14px;
        opacity: 0.9;
      }
    }
  }

  .control-area {
    background: rgba(255, 255, 255, 0.95);
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 20px;

    .progress-wrapper {
      margin-bottom: 20px;

      .time-display {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: #909399;
        margin-top: 8px;
      }
    }

    .control-buttons {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 24px;
      margin-bottom: 20px;

      .control-icon {
        font-size: 12px;
        font-weight: 600;
      }

      .play-btn {
        width: 64px;
        height: 64px;
      }
    }

    .action-buttons {
      display: flex;
      justify-content: center;
      gap: 16px;
    }
  }

  .description-area {
    :deep(.el-card) {
      background: rgba(255, 255, 255, 0.95);
      border-radius: 16px;
    }

    p {
      color: #606266;
      line-height: 1.8;
    }
  }

  .loading-wrapper {
    padding: 40px;
    background: #fff;
    border-radius: 16px;
    margin: 20px;
  }
</style>
