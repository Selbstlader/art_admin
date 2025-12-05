<template>
  <div class="page-content roadbook-list">
    <!-- 顶部操作栏 -->
    <div class="list-header">
      <ElRow justify="space-between" :gutter="16" align="middle">
        <!-- 搜索框 -->
        <ElCol :lg="6" :md="8" :sm="12" :xs="16">
          <ElInput
            v-model="searchKeyword"
            :prefix-icon="Search"
            clearable
            placeholder="搜索路书标题或描述"
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          />
        </ElCol>

        <!-- 筛选区域 -->
        <ElCol :lg="12" :md="10" :sm="0" :xs="0">
          <div class="filter-area">
            <!-- 标签筛选 -->
            <div class="filter-item">
              <span class="filter-label">标签：</span>
              <ElSelect
                v-model="selectedTagId"
                placeholder="全部标签"
                clearable
                style="width: 150px"
                @change="handleTagChange"
              >
                <ElOption v-for="tag in tagList" :key="tag.id" :label="tag.name" :value="tag.id">
                  <span class="tag-option">
                    <span class="tag-dot" :style="{ background: tag.color }"></span>
                    {{ tag.name }}
                  </span>
                </ElOption>
              </ElSelect>
            </div>

            <!-- 排序方式 -->
            <div class="filter-item">
              <span class="filter-label">排序：</span>
              <ElSelect v-model="sortBy" style="width: 130px" @change="handleSortChange">
                <ElOption label="最新发布" value="created_at" />
                <ElOption label="最多浏览" value="view_count" />
                <ElOption label="最多收藏" value="favorite_count" />
              </ElSelect>
            </div>
          </div>
        </ElCol>

        <!-- 新建按钮 -->
        <ElCol :lg="6" :md="6" :sm="12" :xs="8" style="display: flex; justify-content: flex-end">
          <ElButton type="primary" :icon="Plus" @click="handleCreate"> 创建路书 </ElButton>
        </ElCol>
      </ElRow>

      <!-- 移动端筛选区域 -->
      <div class="mobile-filter-area">
        <ElRow :gutter="12">
          <ElCol :span="12">
            <ElSelect
              v-model="selectedTagId"
              placeholder="全部标签"
              clearable
              style="width: 100%"
              @change="handleTagChange"
            >
              <ElOption v-for="tag in tagList" :key="tag.id" :label="tag.name" :value="tag.id" />
            </ElSelect>
          </ElCol>
          <ElCol :span="12">
            <ElSelect v-model="sortBy" style="width: 100%" @change="handleSortChange">
              <ElOption label="最新发布" value="created_at" />
              <ElOption label="最多浏览" value="view_count" />
              <ElOption label="最多收藏" value="favorite_count" />
            </ElSelect>
          </ElCol>
        </ElRow>
      </div>

      <!-- 已选筛选条件展示 -->
      <div v-if="hasActiveFilters" class="active-filters">
        <span class="filter-tip">当前筛选：</span>
        <ElTag v-if="searchKeyword" closable size="small" @close="clearSearchKeyword">
          关键词: {{ searchKeyword }}
        </ElTag>
        <ElTag
          v-if="selectedTagId && selectedTagName"
          closable
          size="small"
          :color="selectedTagColor"
          @close="clearTagFilter"
        >
          标签: {{ selectedTagName }}
        </ElTag>
        <ElButton v-if="hasActiveFilters" link type="primary" size="small" @click="clearAllFilters">
          清除全部
        </ElButton>
      </div>
    </div>

    <!-- 路书列表 -->
    <div class="list-content">
      <div class="roadbook-grid">
        <RoadbookCard
          v-for="item in roadbookList"
          :key="item.id"
          :roadbook="item"
          :loading="loading"
          :show-actions="isOwner(item)"
          @click="handleCardClick"
          @edit="handleEdit"
          @delete="handleDelete"
        />
      </div>

      <!-- 空状态 -->
      <div v-if="showEmpty" class="empty-state">
        <ElEmpty :description="emptyText">
          <ElButton type="primary" @click="handleCreate">创建第一个路书</ElButton>
        </ElEmpty>
      </div>

      <!-- 加载中骨架屏 -->
      <div v-if="loading && roadbookList.length === 0" class="roadbook-grid">
        <RoadbookCard
          v-for="i in 8"
          :key="'skeleton-' + i"
          :roadbook="skeletonRoadbook"
          :loading="true"
        />
      </div>
    </div>

    <!-- 分页 -->
    <div class="list-pagination" v-if="total > 0">
      <ElPagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 36, 48]"
        :total="total"
        :pager-count="7"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { Search, Plus } from '@element-plus/icons-vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useRoadbookStore } from '../../store/roadbook'
  import { useUserStore } from '@/store/modules/user'
  import RoadbookCard from '../../components/RoadbookCard.vue'
  import type { Roadbook, Tag, RoadbookListParams } from '../../types'
  import { TravelMode, RoadbookVisibility, RoadbookStatus } from '../../types'

  defineOptions({
    name: 'RoadbookList'
  })

  const router = useRouter()
  const roadbookStore = useRoadbookStore()
  const userStore = useUserStore()

  // 当前用户ID
  const currentUserId = computed(() => userStore.info?.userId)

  // 搜索和筛选状态
  const searchKeyword = ref('')
  const selectedTagId = ref<number | undefined>(undefined)
  const sortBy = ref<'created_at' | 'view_count' | 'favorite_count'>('created_at')

  // 分页状态
  const currentPage = ref(1)
  const pageSize = ref(12)

  // 标签列表
  const tagList = ref<Tag[]>([])

  // 从 store 获取数据
  const roadbookList = computed(() => roadbookStore.roadbookList || [])
  const total = computed(() => roadbookStore.total)
  const loading = computed(() => roadbookStore.loading)

  // 空状态显示
  const showEmpty = computed(() => {
    return !loading.value && roadbookList.value.length === 0
  })

  // 空状态文本
  const emptyText = computed(() => {
    if (searchKeyword.value || selectedTagId.value) {
      return '未找到符合条件的路书'
    }
    return '暂无路书，快来创建第一个吧'
  })

  // 是否有激活的筛选条件
  const hasActiveFilters = computed(() => {
    return !!searchKeyword.value || !!selectedTagId.value
  })

  // 获取选中标签的名称
  const selectedTagName = computed(() => {
    if (!selectedTagId.value) return ''
    const tag = tagList.value.find((t) => t.id === selectedTagId.value)
    return tag?.name || ''
  })

  // 获取选中标签的颜色
  const selectedTagColor = computed(() => {
    if (!selectedTagId.value) return ''
    const tag = tagList.value.find((t) => t.id === selectedTagId.value)
    return tag?.color || ''
  })

  // 骨架屏占位数据
  const skeletonRoadbook: Roadbook = {
    id: 0,
    userId: 0,
    title: '',
    description: '',
    coverUrl: '',
    startDate: '',
    endDate: '',
    visibility: RoadbookVisibility.PUBLIC,
    travelMode: TravelMode.DRIVING,
    totalDistance: 0,
    totalDuration: 0,
    totalBudget: 0,
    viewCount: 0,
    favoriteCount: 0,
    likeCount: 0,
    commentCount: 0,
    shareCount: 0,
    status: RoadbookStatus.PUBLISHED,
    createdAt: '',
    updatedAt: ''
  }

  // 获取路书列表
  const fetchList = async () => {
    const params: RoadbookListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sortBy: sortBy.value,
      sortOrder: 'desc'
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    if (selectedTagId.value) {
      params.tagId = selectedTagId.value
    }

    await roadbookStore.fetchRoadbookList(params)
  }

  // 获取标签列表
  const fetchTags = async () => {
    await roadbookStore.fetchTags()
    tagList.value = roadbookStore.tagList
  }

  // 搜索处理
  const handleSearch = () => {
    currentPage.value = 1
    fetchList()
  }

  // 标签筛选处理
  const handleTagChange = () => {
    currentPage.value = 1
    fetchList()
  }

  // 排序变化处理
  const handleSortChange = () => {
    currentPage.value = 1
    fetchList()
  }

  // 清除搜索关键词
  const clearSearchKeyword = () => {
    searchKeyword.value = ''
    handleSearch()
  }

  // 清除标签筛选
  const clearTagFilter = () => {
    selectedTagId.value = undefined
    handleTagChange()
  }

  // 清除所有筛选条件
  const clearAllFilters = () => {
    searchKeyword.value = ''
    selectedTagId.value = undefined
    currentPage.value = 1
    fetchList()
  }

  // 分页大小变化
  const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchList()
  }

  // 页码变化
  const handlePageChange = (page: number) => {
    currentPage.value = page
    fetchList()
    // 滚动到顶部
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  // 点击卡片
  const handleCardClick = (roadbook: Roadbook) => {
    router.push(`/travel/travel/roadbook/detail/${roadbook.id}`)
  }

  // 编辑路书
  const handleEdit = (roadbook: Roadbook) => {
    router.push(`/travel/travel/roadbook/editor/${roadbook.id}`)
  }

  // 删除路书
  const handleDelete = async (roadbook: Roadbook) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除路书「${roadbook.title}」吗？删除后无法恢复。`,
        '删除确认',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      await roadbookStore.deleteRoadbook(roadbook.id)
      await fetchList()
      ElMessage.success('删除成功')
    } catch (error) {
      // 用户取消删除，不做处理
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }

  // 判断是否是当前用户的路书
  const isOwner = (roadbook: Roadbook): boolean => {
    return !!currentUserId.value && roadbook.userId === currentUserId.value
  }

  // 创建路书
  const handleCreate = () => {
    router.push('/travel/travel/roadbook/editor')
  }

  // 初始化
  onMounted(() => {
    fetchList()
    fetchTags()
  })
</script>

<style scoped lang="scss">
  .roadbook-list {
    .list-header {
      margin-bottom: 20px;

      .filter-area {
        display: flex;
        align-items: center;
        gap: 20px;

        .filter-item {
          display: flex;
          align-items: center;
          gap: 8px;

          .filter-label {
            font-size: 14px;
            color: var(--art-text-gray-600);
            white-space: nowrap;
          }
        }

        .tag-option {
          display: flex;
          align-items: center;
          gap: 8px;

          .tag-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
          }
        }
      }

      .mobile-filter-area {
        display: none;
        margin-top: 12px;
      }

      .active-filters {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;
        margin-top: 12px;
        padding: 10px 12px;
        background: var(--art-gray-100);
        border-radius: 6px;

        .filter-tip {
          font-size: 13px;
          color: var(--art-text-gray-600);
        }

        .el-tag {
          border: none;
        }
      }
    }

    .list-content {
      min-height: 400px;

      .roadbook-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 20px;
      }

      .empty-state {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 400px;
      }
    }

    .list-pagination {
      display: flex;
      justify-content: center;
      margin-top: 32px;
      padding-bottom: 20px;
    }
  }

  // 响应式布局
  @media only screen and (max-width: 1400px) {
    .roadbook-list {
      .list-content {
        .roadbook-grid {
          grid-template-columns: repeat(3, 1fr);
        }
      }
    }
  }

  @media only screen and (max-width: 1024px) {
    .roadbook-list {
      .list-content {
        .roadbook-grid {
          grid-template-columns: repeat(2, 1fr);
        }
      }
    }
  }

  @media only screen and (max-width: 768px) {
    .roadbook-list {
      .list-header {
        .filter-area {
          display: none;
        }

        .mobile-filter-area {
          display: block;
        }
      }

      .list-content {
        .roadbook-grid {
          grid-template-columns: repeat(2, 1fr);
          gap: 12px;
        }
      }

      .list-pagination {
        :deep(.el-pagination) {
          .el-pagination__sizes,
          .el-pagination__jump {
            display: none;
          }
        }
      }
    }
  }

  @media only screen and (max-width: 480px) {
    .roadbook-list {
      .list-content {
        .roadbook-grid {
          grid-template-columns: 1fr;
        }
      }
    }
  }
</style>
