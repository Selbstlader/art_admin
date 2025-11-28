<template>
  <div class="xhs-history-page">
    <ElCard shadow="never" class="toolbar-card">
      <div class="toolbar">
        <div class="search-area">
          <ElInput
            v-model="searchKeyword"
            placeholder="搜索总结..."
            clearable
            style="width: 250px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <ElIcon><Search /></ElIcon>
            </template>
          </ElInput>
        </div>
        <ElButton type="primary" @click="$router.push('/xhs/summary')">
          <ElIcon class="mr-1"><Plus /></ElIcon>
          新建总结
        </ElButton>
      </div>
    </ElCard>

    <div v-loading="loading" class="history-list">
      <ElRow :gutter="16">
        <ElCol
          v-for="item in historyList"
          :key="item.id"
          :xs="24"
          :sm="12"
          :lg="8"
          :xl="6"
          class="mb-4"
        >
          <ElCard shadow="hover" class="history-card" @click="viewDetail(item)">
            <div class="card-header">
              <span class="card-date">{{ formatDate(item.created_at) }}</span>
              <ElTag :type="getSentimentType(item.sentiment)" size="small">
                {{ SentimentConfig[item.sentiment]?.icon }}
              </ElTag>
            </div>
            <h4 class="card-title">{{ item.summary_title || item.note_title || '未命名总结' }}</h4>
            <p class="card-content">{{ getPreview(item.summary_content) }}</p>
            <div class="card-tags" v-if="item.tags?.length">
              <ElTag
                v-for="tag in item.tags.slice(0, 3)"
                :key="tag"
                size="small"
                effect="plain"
                class="mr-1"
              >
                {{ tag }}
              </ElTag>
            </div>
            <div class="card-footer" @click.stop>
              <ElButton
                link
                :type="item.is_favorite ? 'warning' : 'default'"
                size="small"
                @click="toggleFavorite(item)"
              >
                <ElIcon><component :is="item.is_favorite ? StarFilled : Star" /></ElIcon>
              </ElButton>
              <ElButton link type="danger" size="small" @click="deleteItem(item)">
                <ElIcon><Delete /></ElIcon>
              </ElButton>
            </div>
          </ElCard>
        </ElCol>
      </ElRow>

      <ElEmpty v-if="!loading && historyList.length === 0" description="暂无历史记录">
        <ElButton type="primary" @click="$router.push('/xhs/summary')">去创建总结</ElButton>
      </ElEmpty>

      <div v-if="total > pageSize" class="pagination-wrapper">
        <ElPagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadHistory"
        />
      </div>
    </div>

    <!-- 详情弹窗 -->
    <ElDialog
      v-model="detailVisible"
      :title="currentItem?.summary_title || '总结详情'"
      width="700px"
    >
      <div v-if="currentItem" class="detail-content">
        <div class="detail-meta">
          <span>{{ formatDateTime(currentItem.created_at) }}</span>
          <ElTag :type="getSentimentType(currentItem.sentiment)" size="small">
            {{ SentimentConfig[currentItem.sentiment]?.label }}
          </ElTag>
        </div>
        <div class="detail-section">
          <h4>核心总结</h4>
          <p>{{ currentItem.summary_content }}</p>
        </div>
        <div class="detail-section" v-if="currentItem.key_points?.length">
          <h4>关键信息</h4>
          <ul>
            <li v-for="(point, index) in currentItem.key_points" :key="index">{{ point }}</li>
          </ul>
        </div>
        <div class="detail-section" v-if="currentItem.tags?.length">
          <h4>标签</h4>
          <ElTag v-for="tag in currentItem.tags" :key="tag" class="mr-2">{{ tag }}</ElTag>
        </div>
      </div>
      <template #footer>
        <ElButton @click="detailVisible = false">关闭</ElButton>
        <ElButton type="primary" @click="copyContent">复制总结</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { Search, Plus, Star, StarFilled, Delete } from '@element-plus/icons-vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { xhsApi, SentimentConfig, type XHSSummary, type Sentiment } from '@/api/xhs'
  import dayjs from 'dayjs'

  defineOptions({ name: 'XHSHistory' })

  const searchKeyword = ref('')
  const currentPage = ref(1)
  const pageSize = ref(12)
  const total = ref(0)
  const loading = ref(false)
  const historyList = ref<XHSSummary[]>([])
  const detailVisible = ref(false)
  const currentItem = ref<XHSSummary | null>(null)

  const formatDate = (date: string) => dayjs(date).format('MM-DD HH:mm')
  const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')
  const getPreview = (content: string) =>
    content?.length > 80 ? content.slice(0, 80) + '...' : content

  const getSentimentType = (sentiment: Sentiment) => {
    const map: Record<Sentiment, '' | 'success' | 'warning' | 'info'> = {
      positive: 'success',
      neutral: 'info',
      negative: 'warning'
    }
    return map[sentiment] || 'info'
  }

  const handleSearch = () => {
    currentPage.value = 1
    if (searchKeyword.value) {
      searchHistory()
    } else {
      loadHistory()
    }
  }

  const loadHistory = async () => {
    loading.value = true
    try {
      const res = await xhsApi.getHistory({ page: currentPage.value, page_size: pageSize.value })
      const data = (res as any)?.data || res
      historyList.value = data?.list || []
      total.value = data?.total || 0
    } catch (error) {
      console.error('加载历史失败:', error)
    } finally {
      loading.value = false
    }
  }

  const searchHistory = async () => {
    loading.value = true
    try {
      const res = await xhsApi.search({
        keyword: searchKeyword.value,
        page: currentPage.value,
        page_size: pageSize.value
      })
      const data = (res as any)?.data || res
      historyList.value = data?.list || []
      total.value = data?.total || 0
    } catch (error) {
      console.error('搜索失败:', error)
    } finally {
      loading.value = false
    }
  }

  const viewDetail = (item: XHSSummary) => {
    currentItem.value = item
    detailVisible.value = true
  }

  const toggleFavorite = async (item: XHSSummary) => {
    try {
      await xhsApi.updateFavorite(item.id, { is_favorite: !item.is_favorite })
      item.is_favorite = !item.is_favorite
      ElMessage.success(item.is_favorite ? '已收藏' : '已取消收藏')
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    }
  }

  const deleteItem = async (item: XHSSummary) => {
    try {
      await ElMessageBox.confirm('确定要删除这条总结吗？', '删除确认', { type: 'warning' })
      await xhsApi.deleteSummary(item.id)
      ElMessage.success('删除成功')
      loadHistory()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  const copyContent = () => {
    if (!currentItem.value) return
    const text = `${currentItem.value.summary_title}\n\n${currentItem.value.summary_content}`
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success('已复制到剪贴板')
    })
  }

  onMounted(() => {
    loadHistory()
  })
</script>

<style scoped lang="scss">
  .xhs-history-page {
    padding: 0;
  }

  .toolbar-card {
    margin-bottom: 20px;
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .mr-1 {
    margin-right: 4px;
  }

  .mr-2 {
    margin-right: 8px;
  }

  .mb-4 {
    margin-bottom: 16px;
  }

  .history-card {
    cursor: pointer;
    transition: all 0.3s ease;
    height: 100%;

    &:hover {
      transform: translateY(-4px);
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .card-date {
        font-size: 12px;
        color: var(--art-gray-500);
      }
    }

    .card-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--art-gray-900);
      margin: 0 0 10px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .card-content {
      font-size: 13px;
      color: var(--art-gray-600);
      line-height: 1.6;
      margin-bottom: 12px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .card-tags {
      margin-bottom: 12px;
    }

    .card-footer {
      display: flex;
      justify-content: flex-end;
      gap: 8px;
      padding-top: 12px;
      border-top: 1px solid var(--art-border-color);
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }

  .detail-content {
    .detail-meta {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 20px;
      color: var(--art-gray-500);
    }

    .detail-section {
      margin-bottom: 20px;

      h4 {
        font-size: 14px;
        font-weight: 600;
        color: var(--art-gray-800);
        margin-bottom: 10px;
      }

      p {
        font-size: 14px;
        line-height: 1.8;
        color: var(--art-gray-700);
      }

      ul {
        padding-left: 20px;
        margin: 0;

        li {
          font-size: 14px;
          line-height: 1.8;
          color: var(--art-gray-700);
        }
      }
    }
  }
</style>
