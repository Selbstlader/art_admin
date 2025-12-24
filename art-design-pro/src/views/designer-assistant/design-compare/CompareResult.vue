<template>
  <div class="compare-result-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">比对历史</span>
          <ElButton type="primary" size="small" @click="handleNewCompare">
            <ElIcon><Plus /></ElIcon>
            新建比对
          </ElButton>
        </div>
      </template>

      <!-- 筛选区域 -->
      <div class="filter-area">
        <ElSelect v-model="filterStatus" placeholder="分析状态" clearable style="width: 150px">
          <ElOption label="全部" value="" />
          <ElOption label="待分析" value="pending" />
          <ElOption label="分析中" value="processing" />
          <ElOption label="已完成" value="completed" />
          <ElOption label="失败" value="failed" />
        </ElSelect>
        <ElButton :icon="Refresh" @click="loadCompareList">刷新</ElButton>
      </div>

      <!-- 比对列表 -->
      <ElTable :data="compareList" v-loading="loading" stripe>
        <ElTableColumn label="设计图" width="120">
          <template #default="{ row }">
            <ElImage
              v-if="row.designImages?.length"
              :src="row.designImages[0].fileUrl"
              :preview-src-list="row.designImages.map((img: any) => img.fileUrl)"
              fit="cover"
              style="width: 80px; height: 60px; border-radius: 4px"
            />
            <span v-else>-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="id" label="ID" width="80" />
        <ElTableColumn label="匹配度" width="120">
          <template #default="{ row }">
            <ElProgress
              v-if="row.analysisStatus === 'completed'"
              :percentage="Math.round(row.overallScore)"
              :color="getScoreColor(row.overallScore)"
              :stroke-width="8"
            />
            <span v-else>-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.analysisStatus)" size="small">
              {{ getStatusText(row.analysisStatus) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="匹配/偏差/建议" width="150">
          <template #default="{ row }">
            <div v-if="row.analysisStatus === 'completed'" class="stats-cell">
              <ElTag type="success" size="small">{{ row.matchItems?.length || 0 }}</ElTag>
              <ElTag type="warning" size="small">{{ row.deviationItems?.length || 0 }}</ElTag>
              <ElTag type="info" size="small">{{ row.suggestions?.length || 0 }}</ElTag>
            </div>
            <span v-else>-</span>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton
              v-if="row.analysisStatus === 'pending'"
              type="primary"
              size="small"
              link
              @click="handleAnalyze(row)"
            >
              分析
            </ElButton>
            <ElButton type="primary" size="small" link @click="handleView(row)"> 查看 </ElButton>
            <ElButton type="danger" size="small" link @click="handleDelete(row)"> 删除 </ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <!-- 分页 -->
      <div class="pagination-area">
        <ElPagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadCompareList"
          @current-change="loadCompareList"
        />
      </div>
    </ElCard>

    <!-- 详情弹窗 -->
    <ElDialog v-model="detailVisible" title="比对详情" width="800px" destroy-on-close>
      <div v-if="currentDetail" class="detail-content">
        <!-- 设计图预览 -->
        <div class="image-preview">
          <ElImage
            v-if="currentDetail.designImages?.length"
            :src="currentDetail.designImages[0].fileUrl"
            :preview-src-list="currentDetail.designImages.map(img => img.fileUrl)"
            fit="contain"
            style="max-width: 100%; max-height: 300px"
          />
        </div>

        <!-- 匹配度 -->
        <div class="score-display">
          <ElProgress
            type="circle"
            :percentage="Math.round(currentDetail.overallScore)"
            :color="getScoreColor(currentDetail.overallScore)"
            :width="100"
          />
          <span class="score-label">整体匹配度</span>
        </div>

        <!-- 详细结果 -->
        <ElTabs>
          <ElTabPane label="匹配项">
            <div v-for="(item, index) in currentDetail.matchItems" :key="index" class="result-item">
              <p><strong>需求:</strong> {{ item.requirement }}</p>
              <p><strong>匹配:</strong> {{ item.designMatch }}</p>
              <ElProgress :percentage="item.score" :stroke-width="6" />
            </div>
            <ElEmpty v-if="!currentDetail.matchItems?.length" description="暂无匹配项" />
          </ElTabPane>
          <ElTabPane label="偏差项">
            <ElAlert
              v-for="(item, index) in currentDetail.deviationItems"
              :key="index"
              :type="getSeverityType(item.severity)"
              :title="item.content"
              :closable="false"
              style="margin-bottom: 12px"
            >
              <p>位置: {{ item.location }}</p>
              <p>原始需求: {{ item.originalRequirement }}</p>
            </ElAlert>
            <ElEmpty v-if="!currentDetail.deviationItems?.length" description="暂无偏差项" />
          </ElTabPane>
          <ElTabPane label="建议">
            <div
              v-for="(item, index) in currentDetail.suggestions"
              :key="index"
              class="suggestion-item"
            >
              <p>{{ item.content }}</p>
              <div class="tags">
                <ElTag size="small">优先级: {{ item.priority }}</ElTag>
                <ElTag size="small" type="warning">成本: {{ item.costImpact }}</ElTag>
              </div>
            </div>
            <ElEmpty v-if="!currentDetail.suggestions?.length" description="暂无建议" />
          </ElTabPane>
        </ElTabs>
      </div>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * Compare Result Component
   * 设计比对结果列表组件
   * Requirements: 2.2, 2.4
   ***/
  import { ref, computed, onMounted, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Plus, Refresh } from '@element-plus/icons-vue'
  import {
    getDesignCompareList,
    getDesignCompareDetail,
    analyzeDesignCompare,
    deleteDesignCompare,
    type DesignCompareResponse
  } from '@/api/designer-compare'

  const route = useRoute()
  const router = useRouter()

  // 状态
  const loading = ref(false)
  const compareList = ref<DesignCompareResponse[]>([])
  const filterStatus = ref('')
  const detailVisible = ref(false)
  const currentDetail = ref<DesignCompareResponse | null>(null)
  const pagination = ref({
    current: 1,
    size: 10,
    total: 0
  })

  // 计算属性
  const projectId = computed(() => Number(route.query.projectId) || 0)

  // 加载比对列表
  const loadCompareList = async () => {
    if (!projectId.value) return

    loading.value = true
    try {
      const res = await getDesignCompareList({
        projectId: projectId.value,
        current: pagination.value.current,
        size: pagination.value.size,
        analysisStatus: filterStatus.value || undefined
      })
      if (res.code === 200 && res.data) {
        compareList.value = res.data.records || []
        pagination.value.total = res.data.total
      }
    } catch (error) {
      console.error('加载比对列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 新建比对
  const handleNewCompare = () => {
    router.push({
      path: '/designer-assistant/design-compare/upload',
      query: { projectId: projectId.value }
    })
  }

  // 执行分析
  const handleAnalyze = async (row: DesignCompareResponse) => {
    try {
      await analyzeDesignCompare(row.id)
      ElMessage.success('分析完成')
      loadCompareList()
    } catch (error: any) {
      ElMessage.error(error.message || '分析失败')
    }
  }

  // 查看详情
  const handleView = async (row: DesignCompareResponse) => {
    try {
      const res = await getDesignCompareDetail(row.id)
      if (res.code === 200 && res.data) {
        currentDetail.value = res.data
        detailVisible.value = true
      }
    } catch (error) {
      ElMessage.error('获取详情失败')
    }
  }

  // 删除
  const handleDelete = async (row: DesignCompareResponse) => {
    try {
      await ElMessageBox.confirm('确定要删除这条比对记录吗？', '提示', {
        type: 'warning'
      })
      await deleteDesignCompare(row.id)
      ElMessage.success('删除成功')
      loadCompareList()
    } catch (error: any) {
      if (error !== 'cancel') {
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  // 格式化日期
  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleString('zh-CN')
  }

  // 获取分数颜色
  const getScoreColor = (score: number) => {
    if (score >= 80) return '#67c23a'
    if (score >= 60) return '#e6a23c'
    return '#f56c6c'
  }

  // 获取状态类型
  const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' => {
    const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
      completed: 'success',
      processing: 'warning',
      pending: 'info',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  // 获取状态文本
  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      completed: '已完成',
      processing: '分析中',
      pending: '待分析',
      failed: '失败'
    }
    return map[status] || status
  }

  // 获取严重程度类型（用于Alert组件）
  const getSeverityType = (severity: string): 'success' | 'warning' | 'info' | 'error' => {
    const map: Record<string, 'success' | 'warning' | 'info' | 'error'> = {
      low: 'info',
      medium: 'warning',
      high: 'error'
    }
    return map[severity] || 'info'
  }

  // 监听项目ID变化
  watch(projectId, () => {
    loadCompareList()
  })

  // 监听筛选条件变化
  watch(filterStatus, () => {
    pagination.value.current = 1
    loadCompareList()
  })

  // 初始化
  onMounted(() => {
    loadCompareList()
  })
</script>

<style scoped lang="scss">
  .compare-result-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .filter-area {
      display: flex;
      gap: 12px;
      margin-bottom: 16px;
    }

    .stats-cell {
      display: flex;
      gap: 4px;
    }

    .pagination-area {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }

    .detail-content {
      .image-preview {
        text-align: center;
        margin-bottom: 24px;
        padding: 16px;
        background: var(--el-fill-color-light);
        border-radius: 8px;
      }

      .score-display {
        text-align: center;
        margin-bottom: 24px;

        .score-label {
          display: block;
          margin-top: 8px;
          font-size: 14px;
          color: var(--el-text-color-secondary);
        }
      }

      .result-item {
        padding: 12px;
        margin-bottom: 12px;
        background: var(--el-fill-color-light);
        border-radius: 4px;

        p {
          margin: 4px 0;
          font-size: 13px;
        }
      }

      .suggestion-item {
        padding: 12px;
        margin-bottom: 12px;
        background: var(--el-fill-color-light);
        border-radius: 4px;

        p {
          margin-bottom: 8px;
        }

        .tags {
          display: flex;
          gap: 8px;
        }
      }
    }
  }
</style>
