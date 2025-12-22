<template>
  <div class="version-diff-page">
    <ElCard shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <div class="left">
            <ElButton link @click="handleBack">
              <ElIcon><ArrowLeft /></ElIcon>
              返回
            </ElButton>
            <span class="title">版本对比</span>
          </div>
          <div class="right">
            <ElButton
              type="success"
              @click="handleAIAnalyze"
              :disabled="!versionAId || !versionBId"
              :loading="aiAnalyzing"
            >
              <ElIcon><MagicStick /></ElIcon>
              AI智能分析
            </ElButton>
            <ElButton @click="handleSwap">
              <ElIcon><Switch /></ElIcon>
              交换版本
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 版本选择器 / Version Selector -->
      <div class="version-selector">
        <div class="selector-item">
          <span class="label">版本A:</span>
          <ElSelect v-model="versionAId" placeholder="选择版本A" @change="loadDiff">
            <ElOption
              v-for="v in versions"
              :key="v.id"
              :label="`V${v.versionNumber} - ${v.versionName}`"
              :value="v.id"
              :disabled="v.id === versionBId"
            />
          </ElSelect>
        </div>
        <ElIcon class="vs-icon"><Right /></ElIcon>
        <div class="selector-item">
          <span class="label">版本B:</span>
          <ElSelect v-model="versionBId" placeholder="选择版本B" @change="loadDiff">
            <ElOption
              v-for="v in versions"
              :key="v.id"
              :label="`V${v.versionNumber} - ${v.versionName}`"
              :value="v.id"
              :disabled="v.id === versionAId"
            />
          </ElSelect>
        </div>
      </div>

      <!-- 对比摘要 / Compare Summary -->
      <ElAlert
        v-if="diffData?.summary"
        :title="diffData.summary"
        type="info"
        :closable="false"
        class="summary-alert"
      />

      <!-- 左右分屏对比 / Side by Side Comparison -->
      <div class="diff-container" v-if="diffData">
        <ElRow :gutter="24">
          <!-- 版本A / Version A -->
          <ElCol :span="12">
            <div class="version-panel version-a">
              <div class="panel-header">
                <ElTag type="primary" size="large">
                  V{{ diffData.versionA?.versionNumber }} - {{ diffData.versionA?.versionName }}
                </ElTag>
                <span class="status">
                  <ElTag :type="getStatusType(diffData.versionA?.status || '')">
                    {{ getStatusText(diffData.versionA?.status || '') }}
                  </ElTag>
                </span>
              </div>
              <VersionDetail v-if="diffData.versionA" :version="diffData.versionA" />
            </div>
          </ElCol>

          <!-- 版本B / Version B -->
          <ElCol :span="12">
            <div class="version-panel version-b">
              <div class="panel-header">
                <ElTag type="success" size="large">
                  V{{ diffData.versionB?.versionNumber }} - {{ diffData.versionB?.versionName }}
                </ElTag>
                <span class="status">
                  <ElTag :type="getStatusType(diffData.versionB?.status || '')">
                    {{ getStatusText(diffData.versionB?.status || '') }}
                  </ElTag>
                </span>
              </div>
              <VersionDetail v-if="diffData.versionB" :version="diffData.versionB" />
            </div>
          </ElCol>
        </ElRow>

        <!-- 变化详情 / Change Details -->
        <div class="changes-section">
          <h3>变化详情</h3>

          <!-- 布局变化 / Layout Changes -->
          <ElCollapse v-model="activeCollapse">
            <ElCollapseItem v-if="diffData.layoutChanges?.length" title="布局变化" name="layout">
              <template #title>
                <span>布局变化</span>
                <ElBadge :value="diffData.layoutChanges.length" class="change-badge" />
              </template>
              <ChangeList :changes="diffData.layoutChanges" />
            </ElCollapseItem>

            <ElCollapseItem v-if="diffData.areaChanges?.length" title="面积调整" name="area">
              <template #title>
                <span>面积调整</span>
                <ElBadge :value="diffData.areaChanges.length" class="change-badge" />
              </template>
              <ChangeList :changes="diffData.areaChanges" />
            </ElCollapseItem>

            <ElCollapseItem v-if="diffData.elementChanges?.length" title="元素增减" name="element">
              <template #title>
                <span>元素增减</span>
                <ElBadge :value="diffData.elementChanges.length" class="change-badge" />
              </template>
              <ChangeList :changes="diffData.elementChanges" />
            </ElCollapseItem>

            <ElCollapseItem v-if="diffData.styleChanges?.length" title="风格变化" name="style">
              <template #title>
                <span>风格变化</span>
                <ElBadge :value="diffData.styleChanges.length" class="change-badge" />
              </template>
              <ChangeList :changes="diffData.styleChanges" />
            </ElCollapseItem>

            <ElCollapseItem
              v-if="diffData.materialChanges?.length"
              title="材料变化"
              name="material"
            >
              <template #title>
                <span>材料变化</span>
                <ElBadge :value="diffData.materialChanges.length" class="change-badge" />
              </template>
              <ChangeList :changes="diffData.materialChanges" />
            </ElCollapseItem>
          </ElCollapse>

          <!-- 无变化提示 / No Changes -->
          <ElEmpty v-if="!hasChanges" description="两个版本之间没有检测到差异" />
        </div>
      </div>

      <!-- 空状态 / Empty State -->
      <ElEmpty v-else-if="!loading" description="请选择两个版本进行对比" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Version Diff Component
   * 版本对比视图组件（左右分屏展示差异）
   * Requirements: 7.2, 7.3, 7.4
   ***/
  import { ref, computed, onMounted, watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { ArrowLeft, Switch, Right, MagicStick } from '@element-plus/icons-vue'
  import {
    getDesignVersionList,
    getVersionDiff,
    getVersionCompareDetail,
    startAIVersionAnalysis,
    type DesignVersionResponse,
    type VersionDiffResponse,
    type VersionCompareResponse
  } from '@/api/designer-version'
  import VersionDetail from './VersionDetail.vue'
  import ChangeList from './ChangeList.vue'

  // 定义响应类型 / Define response type
  interface BaseResponse<T = unknown> {
    code: number
    msg: string
    data: T
  }

  const router = useRouter()
  const route = useRoute()

  const loading = ref(false)
  const aiAnalyzing = ref(false)
  const versions = ref<DesignVersionResponse[]>([])
  const versionAId = ref<number | null>(null)
  const versionBId = ref<number | null>(null)
  const diffData = ref<VersionDiffResponse | null>(null)
  const activeCollapse = ref<string[]>(['layout', 'area', 'element', 'style', 'material'])

  // 获取项目ID / Get project ID
  const projectId = ref<number>(0)

  /*** Check if has changes ***/
  const hasChanges = computed(() => {
    if (!diffData.value) return false
    return (
      (diffData.value.layoutChanges?.length || 0) > 0 ||
      (diffData.value.areaChanges?.length || 0) > 0 ||
      (diffData.value.elementChanges?.length || 0) > 0 ||
      (diffData.value.styleChanges?.length || 0) > 0 ||
      (diffData.value.materialChanges?.length || 0) > 0
    )
  })

  /*** Load version list ***/
  const loadVersions = async () => {
    if (!projectId.value) return

    try {
      const res = (await getDesignVersionList({
        projectId: projectId.value,
        current: 1,
        size: 100
      })) as unknown as BaseResponse<{ records: DesignVersionResponse[] }>

      if (res.code === 200 && res.data) {
        versions.value = res.data.records || []
      }
    } catch (error) {
      console.error('获取版本列表失败:', error)
    }
  }

  /*** Load diff data ***/
  const loadDiff = async () => {
    if (!versionAId.value || !versionBId.value) return

    loading.value = true
    try {
      const res = (await getVersionDiff(
        versionAId.value,
        versionBId.value
      )) as unknown as BaseResponse<VersionDiffResponse>

      if (res.code === 200 && res.data) {
        diffData.value = res.data
      } else {
        ElMessage.error(res.msg || '获取版本对比失败')
      }
    } catch (error) {
      console.error('获取版本对比失败:', error)
      ElMessage.error('获取版本对比失败')
    } finally {
      loading.value = false
    }
  }

  /*** Load compare detail by compareId (from notification) ***/
  const loadCompareDetail = async (compareId: number) => {
    loading.value = true
    try {
      const res = (await getVersionCompareDetail(
        compareId
      )) as unknown as BaseResponse<VersionCompareResponse>

      if (res.code === 200 && res.data) {
        const compareData = res.data
        // 设置项目ID和版本ID / Set project ID and version IDs
        projectId.value = compareData.projectId
        versionAId.value = compareData.versionAId
        versionBId.value = compareData.versionBId

        // 加载版本列表 / Load version list
        await loadVersions()

        // 加载对比详情数据 / Load compare detail data
        await loadCompareDetailData(compareData)
      } else {
        ElMessage.error(res.msg || '获取对比详情失败')
      }
    } catch (error) {
      console.error('获取对比详情失败:', error)
      ElMessage.error('获取对比详情失败')
    } finally {
      loading.value = false
    }
  }

  /*** Handle swap versions ***/
  const handleSwap = () => {
    const temp = versionAId.value
    versionAId.value = versionBId.value
    versionBId.value = temp
    loadDiff()
  }

  /*** Handle AI analyze - async analysis ***/
  const handleAIAnalyze = async () => {
    if (!versionAId.value || !versionBId.value || !projectId.value) {
      ElMessage.warning('请先选择两个版本')
      return
    }

    aiAnalyzing.value = true
    try {
      const res = (await startAIVersionAnalysis({
        projectId: projectId.value,
        versionAId: versionAId.value,
        versionBId: versionBId.value
      })) as unknown as BaseResponse<{ taskId: number; message: string }>

      if (res.code === 200 && res.data?.taskId) {
        ElMessage.success('AI分析任务已提交，完成后将通过通知告知您')
        // 开始轮询检查分析状态 / Start polling for analysis status
        pollAnalysisResult(res.data.taskId)
      } else {
        ElMessage.error(res.msg || 'AI分析任务提交失败')
      }
    } catch (error) {
      console.error('AI分析任务提交失败:', error)
      ElMessage.error('AI分析任务提交失败')
    } finally {
      aiAnalyzing.value = false
    }
  }

  /*** Poll for AI analysis result ***/
  const pollAnalysisResult = async (compareId: number) => {
    const maxAttempts = 60 // 最多轮询60次（约2分钟）
    const interval = 2000 // 每2秒轮询一次
    let attempts = 0

    const poll = async () => {
      attempts++
      try {
        const res = (await getVersionCompareDetail(
          compareId
        )) as unknown as BaseResponse<VersionCompareResponse>

        if (res.code === 200 && res.data) {
          if (res.data.compareStatus === 'completed') {
            // 分析完成，更新页面数据 / Analysis completed, update page data
            ElMessage.success('AI分析完成')
            await loadCompareDetailData(res.data)
            return
          } else if (res.data.compareStatus === 'failed') {
            ElMessage.error('AI分析失败')
            return
          }
        }

        // 继续轮询 / Continue polling
        if (attempts < maxAttempts) {
          setTimeout(poll, interval)
        } else {
          ElMessage.info('分析时间较长，完成后将通过通知告知您')
        }
      } catch (error) {
        console.error('轮询分析结果失败:', error)
        if (attempts < maxAttempts) {
          setTimeout(poll, interval)
        }
      }
    }

    // 延迟3秒后开始轮询 / Start polling after 3 seconds delay
    setTimeout(poll, 3000)
  }

  /*** Load compare detail data into diffData ***/
  const loadCompareDetailData = async (compareData: VersionCompareResponse) => {
    // 先加载完整的版本详情 / First load full version details
    if (compareData.versionAId && compareData.versionBId) {
      const fullDiffRes = (await getVersionDiff(
        compareData.versionAId,
        compareData.versionBId
      )) as unknown as BaseResponse<VersionDiffResponse>

      if (fullDiffRes.code === 200 && fullDiffRes.data) {
        // 使用完整的版本详情，但用AI分析的变化数据覆盖 / Use full version details but override with AI analysis changes
        diffData.value = {
          versionA: fullDiffRes.data.versionA,
          versionB: fullDiffRes.data.versionB,
          layoutChanges: compareData.layoutChanges?.length
            ? compareData.layoutChanges
            : fullDiffRes.data.layoutChanges || [],
          areaChanges: compareData.areaChanges?.length
            ? compareData.areaChanges
            : fullDiffRes.data.areaChanges || [],
          elementChanges: compareData.elementChanges?.length
            ? compareData.elementChanges
            : fullDiffRes.data.elementChanges || [],
          styleChanges: compareData.styleChanges?.length
            ? compareData.styleChanges
            : fullDiffRes.data.styleChanges || [],
          materialChanges: compareData.materialChanges?.length
            ? compareData.materialChanges
            : fullDiffRes.data.materialChanges || [],
          summary: compareData.summary || fullDiffRes.data.summary
        }
        return
      }
    }

    // 回退：只使用对比记录数据 / Fallback: use compare record data only
    diffData.value = {
      versionA: compareData.versionA
        ? {
            id: compareData.versionA.id,
            projectId: compareData.projectId,
            versionNumber: compareData.versionA.versionNumber,
            versionName: compareData.versionA.versionName,
            description: '',
            designImages: [],
            cadFileIds: [],
            cadFiles: [],
            layoutInfo: [],
            areaInfo: [],
            styleInfo: [],
            materialInfo: [],
            status: compareData.versionA.status,
            createdBy: 0,
            createdAt: '',
            updatedAt: ''
          }
        : null,
      versionB: compareData.versionB
        ? {
            id: compareData.versionB.id,
            projectId: compareData.projectId,
            versionNumber: compareData.versionB.versionNumber,
            versionName: compareData.versionB.versionName,
            description: '',
            designImages: [],
            cadFileIds: [],
            cadFiles: [],
            layoutInfo: [],
            areaInfo: [],
            styleInfo: [],
            materialInfo: [],
            status: compareData.versionB.status,
            createdBy: 0,
            createdAt: '',
            updatedAt: ''
          }
        : null,
      layoutChanges: compareData.layoutChanges || [],
      areaChanges: compareData.areaChanges || [],
      elementChanges: compareData.elementChanges || [],
      styleChanges: compareData.styleChanges || [],
      materialChanges: compareData.materialChanges || [],
      summary: compareData.summary
    }
  }

  /*** Handle back ***/
  const handleBack = () => {
    router.push({
      path: '/designer-assistant/version-compare',
      query: { projectId: projectId.value }
    })
  }

  /*** Get status type ***/
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      draft: 'info',
      submitted: 'warning',
      approved: 'success',
      rejected: 'danger'
    }
    return map[status] || 'info'
  }

  /*** Get status text ***/
  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      draft: '草稿',
      submitted: '已提交',
      approved: '已批准',
      rejected: '已拒绝'
    }
    return map[status] || status
  }

  // 监听路由参数 / Watch route params
  watch(
    () => route.query,
    (newVal) => {
      // 优先处理compareId（从通知跳转）/ Handle compareId first (from notification)
      if (newVal.compareId) {
        loadCompareDetail(Number(newVal.compareId))
        return
      }

      if (newVal.projectId) {
        projectId.value = Number(newVal.projectId)
        loadVersions()
      }
      if (newVal.versionAId) {
        versionAId.value = Number(newVal.versionAId)
      }
      if (newVal.versionBId) {
        versionBId.value = Number(newVal.versionBId)
      }
      if (versionAId.value && versionBId.value) {
        loadDiff()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    // 优先处理compareId（从通知跳转）/ Handle compareId first (from notification)
    if (route.query.compareId) {
      loadCompareDetail(Number(route.query.compareId))
      return
    }

    if (route.query.projectId) {
      projectId.value = Number(route.query.projectId)
      loadVersions()
    }
    if (route.query.versionAId) {
      versionAId.value = Number(route.query.versionAId)
    }
    if (route.query.versionBId) {
      versionBId.value = Number(route.query.versionBId)
    }
    if (versionAId.value && versionBId.value) {
      loadDiff()
    }
  })
</script>

<style scoped lang="scss">
  .version-diff-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .left {
        display: flex;
        align-items: center;
        gap: 16px;

        .title {
          font-size: 16px;
          font-weight: 500;
        }
      }
    }

    .version-selector {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 24px;
      padding: 16px;
      background: var(--el-fill-color-light);
      border-radius: 8px;
      margin-bottom: 16px;

      .selector-item {
        display: flex;
        align-items: center;
        gap: 8px;

        .label {
          font-weight: 500;
          color: var(--el-text-color-primary);
        }

        .el-select {
          width: 240px;
        }
      }

      .vs-icon {
        font-size: 24px;
        color: var(--el-text-color-secondary);
      }
    }

    .summary-alert {
      margin-bottom: 16px;
    }

    .diff-container {
      .version-panel {
        border: 1px solid var(--el-border-color);
        border-radius: 8px;
        padding: 16px;
        min-height: 400px;

        &.version-a {
          border-top: 3px solid var(--el-color-primary);
        }

        &.version-b {
          border-top: 3px solid var(--el-color-success);
        }

        .panel-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 16px;
          padding-bottom: 12px;
          border-bottom: 1px solid var(--el-border-color-light);
        }
      }

      .changes-section {
        margin-top: 24px;

        h3 {
          margin-bottom: 16px;
          font-size: 16px;
          font-weight: 500;
        }

        .change-badge {
          margin-left: 8px;
        }
      }
    }
  }
</style>
