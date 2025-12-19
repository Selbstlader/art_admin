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
  import { ArrowLeft, Switch, Right } from '@element-plus/icons-vue'
  import {
    getDesignVersionList,
    getVersionDiff,
    type DesignVersionResponse,
    type VersionDiffResponse
  } from '@/api/designer-version'
  import VersionDetail from './VersionDetail.vue'
  import ChangeList from './ChangeList.vue'

  const router = useRouter()
  const route = useRoute()

  const loading = ref(false)
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
      })) as unknown as Http.BaseResponse<{ records: DesignVersionResponse[] }>

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
      )) as unknown as Http.BaseResponse<VersionDiffResponse>

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

  /*** Handle swap versions ***/
  const handleSwap = () => {
    const temp = versionAId.value
    versionAId.value = versionBId.value
    versionBId.value = temp
    loadDiff()
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
