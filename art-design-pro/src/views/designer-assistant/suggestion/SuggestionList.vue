<template>
  <div class="suggestion-list">
    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="left">
        <ElDropdown split-button type="primary" :loading="generating" @click="handleGenerate">
          <ElIcon><MagicStick /></ElIcon>
          生成建议
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem @click="handleGenerateWithOptions(3)">生成 3 条建议</ElDropdownItem>
              <ElDropdownItem @click="handleGenerateWithOptions(5)">生成 5 条建议</ElDropdownItem>
              <ElDropdownItem @click="handleGenerateWithOptions(10)">生成 10 条建议</ElDropdownItem>
              <ElDropdownItem divided @click="showGenerateDialog = true">
                自定义生成...
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
        <ElSelect v-model="filterCategory" placeholder="类别筛选" clearable style="width: 120px">
          <ElOption
            v-for="cat in categories"
            :key="cat.value"
            :label="cat.label"
            :value="cat.value"
          />
        </ElSelect>
        <ElSelect v-model="filterStatus" placeholder="状态筛选" clearable style="width: 120px">
          <ElOption label="待处理" value="pending" />
          <ElOption label="已采纳" value="adopted" />
          <ElOption label="已忽略" value="ignored" />
        </ElSelect>
      </div>
      <div class="right">
        <ElButton @click="loadSuggestions">刷新</ElButton>
      </div>
    </div>

    <!-- 自定义生成对话框 -->
    <ElDialog v-model="showGenerateDialog" title="自定义生成建议" width="400px">
      <ElForm :model="generateForm" label-width="100px">
        <ElFormItem label="生成数量">
          <ElInputNumber v-model="generateForm.count" :min="1" :max="10" />
        </ElFormItem>
        <ElFormItem label="指定类别">
          <ElSelect
            v-model="generateForm.category"
            placeholder="不限"
            clearable
            style="width: 100%"
          >
            <ElOption
              v-for="cat in categories"
              :key="cat.value"
              :label="cat.label"
              :value="cat.value"
            />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <div class="generate-tips">
        <ElAlert type="info" :closable="false" show-icon>
          <template #title> 建议将基于项目的文档、CAD图纸、材料清单和成本信息综合生成 </template>
        </ElAlert>
      </div>
      <template #footer>
        <ElButton @click="showGenerateDialog = false">取消</ElButton>
        <ElButton type="primary" :loading="generating" @click="handleGenerateCustom">
          生成
        </ElButton>
      </template>
    </ElDialog>

    <!-- 统计信息 -->
    <div class="stats-bar" v-if="stats">
      <ElTag>总计: {{ stats.totalCount }}</ElTag>
      <ElTag type="success">已采纳: {{ stats.adoptedCount }}</ElTag>
      <ElTag type="warning">待处理: {{ stats.pendingCount }}</ElTag>
      <ElTag type="info">已忽略: {{ stats.ignoredCount }}</ElTag>
    </div>

    <!-- 建议列表 -->
    <div class="suggestion-cards" v-loading="loading">
      <ElEmpty
        v-if="suggestions.length === 0 && !loading"
        description="暂无设计建议，点击上方按钮生成"
      />

      <div v-for="suggestion in suggestions" :key="suggestion.id" class="suggestion-card">
        <div class="card-header">
          <div class="left">
            <ElTag :type="getCategoryType(suggestion.category)" size="small">
              {{ getCategoryLabel(suggestion.category) }}
            </ElTag>
            <span class="priority">优先级: {{ suggestion.priority }}</span>
          </div>
          <div class="right">
            <ElTag :type="getStatusType(suggestion.status)" size="small">
              {{ getStatusLabel(suggestion.status) }}
            </ElTag>
          </div>
        </div>

        <div class="card-content">
          <p class="content-text">{{ suggestion.content }}</p>

          <ElCollapse v-if="suggestion.applicableScene || suggestion.costImpact">
            <ElCollapseItem title="查看详情">
              <div class="detail-item" v-if="suggestion.applicableScene">
                <span class="label">适用场景:</span>
                <span class="value">{{ suggestion.applicableScene }}</span>
              </div>
              <div class="detail-item" v-if="suggestion.costImpact">
                <span class="label">成本影响:</span>
                <ElTag size="small" :type="getCostImpactType(suggestion.costImpact)">
                  {{ suggestion.costImpact }}
                </ElTag>
              </div>
              <div class="detail-item" v-if="suggestion.detailInfo">
                <span class="label">详细信息:</span>
                <pre class="detail-json">{{ JSON.stringify(suggestion.detailInfo, null, 2) }}</pre>
              </div>
            </ElCollapseItem>
          </ElCollapse>
        </div>

        <div class="card-footer" v-if="suggestion.status === 'pending'">
          <ElButton type="success" size="small" @click="handleAdopt(suggestion)">
            <ElIcon><Check /></ElIcon>
            采纳
          </ElButton>
          <ElButton type="info" size="small" @click="handleIgnore(suggestion)">
            <ElIcon><Close /></ElIcon>
            忽略
          </ElButton>
        </div>
        <div class="card-footer" v-else>
          <ElButton size="small" @click="handleReset(suggestion)">重置状态</ElButton>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="total > pageSize">
      <ElPagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="loadSuggestions"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  /***
   * Suggestion List Component
   * 设计建议列表组件
   * Requirements: 3.2, 3.3, 3.4
   ***/
  import { ref, watch, onMounted, onUnmounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { MagicStick, Check, Close } from '@element-plus/icons-vue'
  import {
    generateSuggestions,
    getSuggestionList,
    updateSuggestionStatus,
    getSuggestionStats,
    type DesignSuggestion,
    type SuggestionStatsResponse
  } from '@/api/designer-suggestion'

  const props = defineProps<{
    projectId: number
  }>()

  const loading = ref(false)
  const generating = ref(false)
  const suggestions = ref<DesignSuggestion[]>([])
  const stats = ref<SuggestionStatsResponse | null>(null)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const filterCategory = ref('')
  const filterStatus = ref('')

  /*** Generate dialog state ***/
  const showGenerateDialog = ref(false)
  const generateForm = ref({
    count: 5,
    category: ''
  })

  /*** Category options ***/
  const categories = [
    { value: 'layout', label: '空间布局' },
    { value: 'material', label: '材料选择' },
    { value: 'style', label: '设计风格' },
    { value: 'function', label: '功能规划' },
    { value: 'lighting', label: '照明设计' },
    { value: 'hvac', label: '暖通空调' },
    { value: 'general', label: '综合建议' }
  ]

  /*** Load suggestions from API ***/
  const loadSuggestions = async () => {
    if (!props.projectId) return

    loading.value = true
    try {
      const res: any = await getSuggestionList({
        projectId: props.projectId,
        current: currentPage.value,
        size: pageSize.value,
        category: filterCategory.value,
        status: filterStatus.value
      })

      if (res.code === 200 && res.data) {
        suggestions.value = res.data.records || []
        total.value = res.data.total || 0
      }
    } catch (error) {
      console.error('加载建议失败:', error)
      ElMessage.error('加载建议失败')
    } finally {
      loading.value = false
    }
  }

  /*** Load stats ***/
  const loadStats = async () => {
    if (!props.projectId) return

    try {
      const res: any = await getSuggestionStats(props.projectId)
      if (res.code === 200 && res.data) {
        stats.value = res.data
      }
    } catch (error) {
      console.error('加载统计失败:', error)
    }
  }

  /*** Generate new suggestions (default 5) ***/
  const handleGenerate = async () => {
    await doGenerate(5)
  }

  /*** Generate with specific count ***/
  const handleGenerateWithOptions = async (count: number) => {
    await doGenerate(count)
  }

  /*** Generate with custom options from dialog ***/
  const handleGenerateCustom = async () => {
    showGenerateDialog.value = false
    await doGenerate(generateForm.value.count, generateForm.value.category)
  }

  /*** Core generate function (async mode) ***/
  // 修复：异步任务提交后启动轮询机制自动刷新列表
  const doGenerate = async (count: number, category?: string) => {
    generating.value = true
    try {
      const res: any = await generateSuggestions({
        projectId: props.projectId,
        count,
        category
      })

      if (res.code === 200 && res.data) {
        // 异步模式：任务已提交，通过通知提醒用户
        ElMessage.success(res.data.message || '建议生成任务已提交，完成后将通过通知提醒您')
        // 启动轮询，定期刷新列表以显示新生成的建议
        startGenerationPolling()
      } else {
        ElMessage.error(res.msg || '提交生成任务失败')
      }
    } catch (error) {
      console.error('提交生成任务失败:', error)
      ElMessage.error('提交生成任务失败，请稍后重试')
    } finally {
      generating.value = false
    }
  }

  /*** 轮询相关 / Polling related ***/
  let pollingTimer: ReturnType<typeof setInterval> | null = null
  let pollingCount = 0
  const maxPollingCount = 12 // 最多轮询12次（约1分钟）

  const startGenerationPolling = () => {
    stopGenerationPolling()
    pollingCount = 0
    pollingTimer = setInterval(async () => {
      pollingCount++
      await loadSuggestions()
      await loadStats()
      // 达到最大轮询次数后停止
      if (pollingCount >= maxPollingCount) {
        stopGenerationPolling()
      }
    }, 5000) // 每5秒轮询一次
  }

  const stopGenerationPolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  /*** Adopt suggestion ***/
  const handleAdopt = async (suggestion: DesignSuggestion) => {
    try {
      await ElMessageBox.confirm('确定采纳此建议吗？', '确认', { type: 'success' })

      const res: any = await updateSuggestionStatus({
        id: suggestion.id,
        status: 'adopted'
      })

      if (res.code === 200) {
        ElMessage.success('已采纳')
        suggestion.status = 'adopted'
        loadStats()
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('操作失败')
      }
    }
  }

  /*** Ignore suggestion ***/
  const handleIgnore = async (suggestion: DesignSuggestion) => {
    try {
      await ElMessageBox.confirm('确定忽略此建议吗？', '确认', { type: 'warning' })

      const res: any = await updateSuggestionStatus({
        id: suggestion.id,
        status: 'ignored'
      })

      if (res.code === 200) {
        ElMessage.success('已忽略')
        suggestion.status = 'ignored'
        loadStats()
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('操作失败')
      }
    }
  }

  /*** Reset suggestion status ***/
  const handleReset = async (suggestion: DesignSuggestion) => {
    try {
      const res: any = await updateSuggestionStatus({
        id: suggestion.id,
        status: 'pending'
      })

      if (res.code === 200) {
        ElMessage.success('已重置')
        suggestion.status = 'pending'
        loadStats()
      }
    } catch (error) {
      console.error('重置失败:', error)
      ElMessage.error('操作失败')
    }
  }

  /*** Get category tag type ***/
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getCategoryType = (category: string): TagType => {
    const map: Record<string, TagType> = {
      layout: 'primary',
      material: 'success',
      style: 'warning',
      function: 'info',
      lighting: 'primary',
      hvac: 'danger',
      general: 'info'
    }
    return map[category] || 'info'
  }

  /*** Get category label ***/
  const getCategoryLabel = (category: string): string => {
    const cat = categories.find((c) => c.value === category)
    return cat?.label || category
  }

  /*** Get status tag type ***/
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      pending: 'warning',
      adopted: 'success',
      ignored: 'info'
    }
    return map[status] || 'info'
  }

  /*** Get status label ***/
  const getStatusLabel = (status: string): string => {
    const map: Record<string, string> = {
      pending: '待处理',
      adopted: '已采纳',
      ignored: '已忽略'
    }
    return map[status] || status
  }

  /*** Get cost impact tag type ***/
  const getCostImpactType = (impact: string): TagType => {
    if (impact.includes('节省') || impact.includes('降低')) return 'success'
    if (impact.includes('增加') || impact.includes('提高')) return 'warning'
    return 'info'
  }

  /*** Watch filter changes ***/
  watch([filterCategory, filterStatus], () => {
    currentPage.value = 1
    loadSuggestions()
  })

  onMounted(() => {
    loadSuggestions()
    loadStats()
  })

  /*** 组件卸载时清理轮询 / Clean up polling on unmount ***/
  onUnmounted(() => {
    stopGenerationPolling()
  })
</script>

<style scoped lang="scss">
  .suggestion-list {
    .action-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      .left {
        display: flex;
        gap: 12px;
        align-items: center;
      }
    }

    .stats-bar {
      display: flex;
      gap: 12px;
      margin-bottom: 16px;
    }

    .suggestion-cards {
      display: flex;
      flex-direction: column;
      gap: 16px;
      min-height: 200px;
    }

    .suggestion-card {
      border: 1px solid var(--el-border-color-light);
      border-radius: 8px;
      padding: 16px;
      background: var(--el-bg-color);
      transition: box-shadow 0.3s;

      &:hover {
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
      }

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        .left {
          display: flex;
          align-items: center;
          gap: 12px;

          .priority {
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }
      }

      .card-content {
        .content-text {
          margin: 0 0 12px;
          line-height: 1.6;
          color: var(--el-text-color-primary);
        }

        .detail-item {
          margin-bottom: 8px;

          .label {
            font-weight: 500;
            margin-right: 8px;
            color: var(--el-text-color-secondary);
          }

          .value {
            color: var(--el-text-color-primary);
          }
        }

        .detail-json {
          background: var(--el-fill-color-light);
          padding: 8px;
          border-radius: 4px;
          font-size: 12px;
          overflow-x: auto;
          margin-top: 4px;
        }
      }

      .card-footer {
        display: flex;
        gap: 8px;
        margin-top: 12px;
        padding-top: 12px;
        border-top: 1px solid var(--el-border-color-lighter);
      }
    }

    .pagination {
      display: flex;
      justify-content: center;
      margin-top: 24px;
    }

    .generate-tips {
      margin-top: 16px;
    }
  }
</style>
