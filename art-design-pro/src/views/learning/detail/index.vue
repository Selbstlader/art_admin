<template>
  <div class="detail-page art-full-height">
    <ElCard v-if="loading" shadow="never">
      <ElSkeleton :rows="10" animated />
    </ElCard>

    <ElCard v-else-if="material" shadow="never">
      <!-- 页面头部 -->
      <div class="page-header">
        <ElPageHeader @back="$router.back()">
          <template #content>
            <span class="header-title">课程详情</span>
          </template>
          <template #extra>
            <ElSpace>
              <ElButton @click="handleFavorite">
                <ElIcon><Star /></ElIcon>
                {{ material.is_favorite ? '取消收藏' : '收藏' }}
              </ElButton>
              <ElButton type="danger" @click="handleDelete">
                <ElIcon><Delete /></ElIcon>
                删除
              </ElButton>
            </ElSpace>
          </template>
        </ElPageHeader>
      </div>

      <ElDivider />

      <!-- 课程信息 -->
      <div class="material-header">
        <h1>{{ material.title }}</h1>
        <ElSpace wrap style="margin-top: 16px">
          <ElTag>{{ material.grade }}</ElTag>
          <ElTag :type="getDifficultyType(material.difficulty)">
            {{ getDifficultyText(material.difficulty) }}
          </ElTag>
          <ElTag type="info">
            <ElIcon><Clock /></ElIcon>
            预计 {{ material.total_time }} 分钟
          </ElTag>
          <ElTag type="info">
            <ElIcon><View /></ElIcon>
            {{ material.view_count || 0 }} 次浏览
          </ElTag>
          <ElTag type="info">
            <ElIcon><Star /></ElIcon>
            {{ material.favorite_count || 0 }} 次收藏
          </ElTag>
        </ElSpace>

        <ElAlert
          v-if="material.content?.summary"
          :title="material.content.summary"
          type="info"
          :closable="false"
          style="margin-top: 20px"
        />
      </div>

      <ElDivider />

      <!-- 课程内容 -->
      <div class="material-content">
        <div v-html="formatContent(material.content)" />
      </div>
    </ElCard>

    <ElEmpty v-else description="课程不存在" />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Star, Delete, Clock, View } from '@element-plus/icons-vue'
  import { learningMaterialApi, type LearningMaterial } from '@/api/learning'

  const route = useRoute()
  const router = useRouter()
  const loading = ref(false)
  const material = ref<LearningMaterial | null>(null)

  const loadDetail = async () => {
    try {
      loading.value = true
      const id = parseInt(route.params.id as string)
      const res: any = await learningMaterialApi.getMaterialDetail(id)
      material.value = res
    } catch (error) {
      ElMessage.error('加载课程失败')
    } finally {
      loading.value = false
    }
  }

  const handleFavorite = () => {
    ElMessage.info('收藏功能开发中...')
  }

  const handleDelete = async () => {
    try {
      await ElMessageBox.confirm('确定删除这个课程吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      if (material.value) {
        await learningMaterialApi.deleteMaterial(material.value.id)
        ElMessage.success('删除成功')
        router.back()
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }

  const getDifficultyType = (difficulty: number) => {
    const types: Record<number, any> = { 1: 'success', 2: 'warning', 3: 'danger' }
    return types[difficulty] || ''
  }

  const getDifficultyText = (difficulty: number) => {
    const texts: Record<number, string> = { 1: '⭐ 基础', 2: '⭐⭐ 进阶', 3: '⭐⭐⭐ 高级' }
    return texts[difficulty] || '未知'
  }

  const formatContent = (content: any) => {
    if (!content || !content.sections) return ''
    return content.sections
      .map((section: any) => {
        let html = `<h3>${section.title}</h3>`
        if (section.content) {
          html += `<p>${section.content.replace(/\n/g, '<br>')}</p>`
        }
        if (section.key_points) {
          html += '<ul>'
          section.key_points.forEach((point: string) => {
            html += `<li>${point}</li>`
          })
          html += '</ul>'
        }
        if (section.question) {
          html += `<div class="question"><strong>题目：</strong>${section.question}</div>`
        }
        if (section.solution) {
          html += `<div class="solution"><strong>解答：</strong>${section.solution}</div>`
        }
        if (section.answer) {
          html += `<div class="answer"><strong>答案：</strong>${section.answer}</div>`
        }
        return html
      })
      .join('<hr/>')
  }

  onMounted(() => {
    loadDetail()
  })
</script>

<style scoped lang="scss">
  .detail-page {
    height: 100%;
    overflow-y: auto;

    .page-header {
      margin-bottom: 16px;

      .header-title {
        font-size: 20px;
        font-weight: 600;
      }
    }

    .material-header {
      h1 {
        font-size: 28px;
        margin: 0;
        color: var(--el-text-color-primary);
      }
    }

    .material-content {
      max-width: 900px;
      margin: 0 auto;
      padding: 20px 0 40px;

      :deep(h3) {
        margin: 24px 0 12px;
        color: var(--el-color-primary);
        font-size: 20px;
        border-left: 4px solid var(--el-color-primary);
        padding-left: 12px;
      }

      :deep(p) {
        line-height: 1.8;
        margin: 12px 0;
        color: var(--el-text-color-regular);
        font-size: 16px;
      }

      :deep(ul) {
        padding-left: 24px;
        margin: 12px 0;

        li {
          margin: 10px 0;
          line-height: 1.6;
          font-size: 15px;
        }
      }

      :deep(hr) {
        margin: 32px 0;
        border: none;
        border-top: 1px solid var(--el-border-color);
      }

      :deep(.question),
      :deep(.solution),
      :deep(.answer) {
        margin: 16px 0;
        padding: 16px;
        background: var(--el-fill-color-light);
        border-radius: 4px;
        line-height: 1.6;

        strong {
          color: var(--el-color-primary);
          margin-right: 8px;
        }
      }

      :deep(.question) {
        background: var(--el-color-primary-light-9);
        border-left: 4px solid var(--el-color-primary);
      }

      :deep(.solution) {
        background: var(--el-color-success-light-9);
        border-left: 4px solid var(--el-color-success);
      }

      :deep(.answer) {
        background: var(--el-color-warning-light-9);
        border-left: 4px solid var(--el-color-warning);
      }
    }
  }
</style>
