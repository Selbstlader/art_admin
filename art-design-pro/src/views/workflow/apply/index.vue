<template>
  <div class="apply-page">
    <!-- 搜索区域 -->
    <ElCard shadow="never" class="search-card">
      <ElForm :inline="true" :model="searchForm" class="search-form">
        <ElFormItem label="流程名称">
          <ElInput
            v-model="searchForm.name"
            placeholder="请输入流程名称"
            clearable
            @keyup.enter="handleSearch"
          />
        </ElFormItem>
        <ElFormItem label="分类">
          <ElSelect
            v-model="searchForm.category"
            placeholder="请选择分类"
            clearable
            style="width: 150px"
          >
            <ElOption
              v-for="cat in categories"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </ElButton>
          <ElButton @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 流程列表 -->
    <ElCard shadow="never" class="list-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="title">可用流程</span>
          <span class="count">共 {{ total }} 个流程</span>
        </div>
      </template>

      <div v-if="processList.length > 0" class="process-grid">
        <div
          v-for="process in processList"
          :key="process.id"
          class="process-card"
          @click="handleApply(process)"
        >
          <div class="process-icon">
            <el-icon :size="32">
              <Document />
            </el-icon>
          </div>
          <div class="process-info">
            <h4 class="process-name">{{ process.name }}</h4>
            <p class="process-desc">{{ process.description || '暂无描述' }}</p>
            <div class="process-meta">
              <ElTag v-if="process.category" size="small" type="info">
                {{ process.category }}
              </ElTag>
              <span class="version">v{{ process.version }}</span>
            </div>
          </div>
          <div class="process-action">
            <ElButton type="primary" size="small" circle>
              <el-icon><Right /></el-icon>
            </ElButton>
          </div>
        </div>
      </div>

      <ElEmpty v-else description="暂无可用流程" />

      <!-- 分页 -->
      <div v-if="total > 0" class="pagination-wrapper">
        <ElPagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="total"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { processDefApi, type ProcessDefResponse } from '@/api/workflow'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Document, Right } from '@element-plus/icons-vue'

defineOptions({ name: 'WorkflowApply' })

const router = useRouter()

// 加载状态
const loading = ref(false)

// 搜索表单
const searchForm = reactive({
  name: '',
  category: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 12
})

// 流程列表
const processList = ref<ProcessDefResponse[]>([])
const total = ref(0)

// 分类列表（从流程列表中提取）
const categories = computed(() => {
  const cats = new Set<string>()
  processList.value.forEach(p => {
    if (p.category) cats.add(p.category)
  })
  return Array.from(cats)
})

// 加载流程列表
const loadProcessList = async () => {
  loading.value = true
  try {
    const res = await processDefApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      name: searchForm.name || undefined,
      category: searchForm.category || undefined,
      status: 'published' // 只显示已发布的流程
    })
    processList.value = res.list || []
    total.value = res.total || 0
  } catch (error: any) {
    console.error('加载流程列表失败:', error)
    ElMessage.error(error.message || '加载流程列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadProcessList()
}

// 重置
const handleReset = () => {
  searchForm.name = ''
  searchForm.category = ''
  pagination.page = 1
  loadProcessList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadProcessList()
}

// 页码变化
const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadProcessList()
}

// 发起申请
const handleApply = (process: ProcessDefResponse) => {
  router.push(`/workflow/apply/form/${process.id}`)
}

// 初始化
onMounted(() => {
  loadProcessList()
})
</script>

<style lang="scss" scoped>
.apply-page {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .search-card {
    :deep(.el-card__body) {
      padding: 16px 20px 0;
    }

    .search-form {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
  }

  .list-card {
    flex: 1;

    .card-header {
      display: flex;
      align-items: center;
      gap: 12px;

      .title {
        font-size: 16px;
        font-weight: 500;
      }

      .count {
        font-size: 13px;
        color: var(--el-text-color-secondary);
      }
    }
  }

  .process-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 16px;
  }

  .process-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: var(--el-fill-color-lighter);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: var(--el-color-primary-light-5);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-2px);

      .process-action {
        opacity: 1;
      }
    }

    .process-icon {
      flex-shrink: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 56px;
      height: 56px;
      background: var(--el-color-primary-light-9);
      border-radius: 12px;
      color: var(--el-color-primary);
    }

    .process-info {
      flex: 1;
      min-width: 0;

      .process-name {
        margin: 0 0 4px;
        font-size: 15px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .process-desc {
        margin: 0 0 8px;
        font-size: 13px;
        color: var(--el-text-color-secondary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .process-meta {
        display: flex;
        align-items: center;
        gap: 8px;

        .version {
          font-size: 12px;
          color: var(--el-text-color-placeholder);
        }
      }
    }

    .process-action {
      flex-shrink: 0;
      opacity: 0;
      transition: opacity 0.3s;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
    padding-top: 16px;
    border-top: 1px solid var(--el-border-color-lighter);
  }
}
</style>
