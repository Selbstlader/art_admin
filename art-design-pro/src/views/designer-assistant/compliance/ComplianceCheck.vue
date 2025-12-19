<template>
  <div class="compliance-check-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">合规检查</span>
          <ElButton type="primary" link @click="showHistoryDialog">
            <ElIcon><Clock /></ElIcon>
            检查历史
          </ElButton>
        </div>
      </template>

      <!-- 项目选择 -->
      <ElForm :inline="true" :model="checkForm" class="check-form">
        <ElFormItem label="选择项目">
          <ElSelect
            v-model="checkForm.projectId"
            placeholder="请选择项目"
            style="width: 300px"
            filterable
            @change="handleProjectChange"
          >
            <ElOption
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="检查类别">
          <ElCheckboxGroup v-model="checkForm.categories">
            <ElCheckbox v-for="cat in STANDARD_CATEGORIES" :key="cat.value" :label="cat.value">
              {{ cat.label }}
            </ElCheckbox>
          </ElCheckboxGroup>
        </ElFormItem>
        <ElFormItem>
          <ElButton
            type="primary"
            :loading="loading"
            :disabled="!checkForm.projectId"
            @click="handleCheck"
          >
            <ElIcon><Search /></ElIcon>
            开始检查
          </ElButton>
        </ElFormItem>
      </ElForm>

      <!-- 检查摘要 -->
      <div v-if="checkSummary && checkForm.projectId" class="summary-section">
        <ElDivider content-position="left">检查摘要</ElDivider>
        <ElRow :gutter="16">
          <ElCol :span="6">
            <ElStatistic title="历史检查次数" :value="checkSummary.totalChecks" />
          </ElCol>
          <ElCol :span="6">
            <ElStatistic title="最近通过项" :value="checkSummary.passedCount">
              <template #suffix>
                <ElIcon color="#67c23a"><CircleCheck /></ElIcon>
              </template>
            </ElStatistic>
          </ElCol>
          <ElCol :span="6">
            <ElStatistic title="最近不合规项" :value="checkSummary.failedCount">
              <template #suffix>
                <ElIcon color="#f56c6c"><CircleClose /></ElIcon>
              </template>
            </ElStatistic>
          </ElCol>
          <ElCol :span="6">
            <ElStatistic title="最近合规分数" :value="checkSummary.averageScore" suffix="分" />
          </ElCol>
        </ElRow>
        <div v-if="checkSummary.latestCheckAt" class="latest-check-time">
          最近检查时间: {{ checkSummary.latestCheckAt }}
        </div>
      </div>

      <!-- 检查结果 -->
      <div v-if="checkResult" class="result-section">
        <ElDivider content-position="left">
          检查结果
          <ElTag :type="getScoreType(checkResult.overallScore)" style="margin-left: 8px">
            合规分数: {{ checkResult.overallScore.toFixed(1) }}分
          </ElTag>
        </ElDivider>

        <!-- 结果统计 -->
        <ElRow :gutter="16" class="stats-row">
          <ElCol :span="8">
            <ElStatistic title="通过项" :value="checkResult.passedItems?.length || 0">
              <template #suffix>
                <ElIcon color="#67c23a"><CircleCheck /></ElIcon>
              </template>
            </ElStatistic>
          </ElCol>
          <ElCol :span="8">
            <ElStatistic title="不合规项" :value="checkResult.failedItems?.length || 0">
              <template #suffix>
                <ElIcon color="#f56c6c"><CircleClose /></ElIcon>
              </template>
            </ElStatistic>
          </ElCol>
          <ElCol :span="8">
            <ElStatistic title="改进建议" :value="checkResult.suggestions?.length || 0">
              <template #suffix>
                <ElIcon color="#e6a23c"><Warning /></ElIcon>
              </template>
            </ElStatistic>
          </ElCol>
        </ElRow>

        <!-- 详细结果 -->
        <ElTabs v-model="activeTab" class="result-tabs">
          <ElTabPane label="通过项" name="passed">
            <ElTable :data="checkResult.passedItems" stripe>
              <ElTableColumn prop="standardCode" label="规范编号" width="150" />
              <ElTableColumn prop="standardName" label="规范名称" />
              <ElTableColumn prop="category" label="类别" width="120">
                <template #default="{ row }">
                  <ElTag size="small">{{ getCategoryLabel(row.category) }}</ElTag>
                </template>
              </ElTableColumn>
              <ElTableColumn prop="description" label="说明" />
              <ElTableColumn label="状态" width="80">
                <template #default>
                  <ElTag type="success" size="small">通过</ElTag>
                </template>
              </ElTableColumn>
            </ElTable>
            <ElEmpty v-if="!checkResult.passedItems?.length" description="暂无通过项" />
          </ElTabPane>

          <ElTabPane label="不合规项" name="failed">
            <div v-for="(item, index) in checkResult.failedItems" :key="index" class="failed-item">
              <ElAlert
                :title="`[${item.standardCode}] ${item.standardName}`"
                type="error"
                :closable="false"
                show-icon
              >
                <template #default>
                  <div class="failed-detail">
                    <p><strong>类别:</strong> {{ getCategoryLabel(item.category) }}</p>
                    <p><strong>问题位置:</strong> {{ item.location }}</p>
                    <p><strong>违规内容:</strong> {{ item.violationContent }}</p>
                    <p><strong>原始规范:</strong> {{ item.originalRequirement }}</p>
                    <p>
                      <strong>严重程度:</strong>
                      <ElTag :type="getSeverityType(item.severity)" size="small">{{
                        getSeverityLabel(item.severity)
                      }}</ElTag>
                    </p>
                    <p><strong>修改建议:</strong> {{ item.suggestion }}</p>
                  </div>
                </template>
              </ElAlert>
            </div>
            <ElEmpty v-if="!checkResult.failedItems?.length" description="暂无不合规项" />
          </ElTabPane>

          <ElTabPane label="改进建议" name="suggestions">
            <div
              v-for="(item, index) in checkResult.suggestions"
              :key="index"
              class="suggestion-item"
            >
              <ElCard shadow="hover">
                <template #header>
                  <div class="suggestion-header">
                    <span>{{ item.content }}</span>
                    <ElTag size="small" :type="getPriorityType(item.priority)">{{
                      getPriorityLabel(item.priority)
                    }}</ElTag>
                  </div>
                </template>
                <p><strong>相关类别:</strong> {{ getCategoryLabel(item.category) }}</p>
                <p><strong>成本影响:</strong> {{ item.costImpact }}</p>
                <p><strong>参考规范:</strong> {{ item.reference }}</p>
              </ElCard>
            </div>
            <ElEmpty v-if="!checkResult.suggestions?.length" description="暂无改进建议" />
          </ElTabPane>
        </ElTabs>
      </div>

      <ElEmpty v-else-if="!loading" description="请选择项目并开始检查" />
    </ElCard>

    <!-- 检查历史对话框 -->
    <ElDialog v-model="historyDialogVisible" title="检查历史" width="800px">
      <ElTable :data="historyList" stripe v-loading="historyLoading">
        <ElTableColumn prop="id" label="ID" width="80" />
        <ElTableColumn prop="checkType" label="检查类型" width="100">
          <template #default="{ row }">{{
            row.checkType === 'full' ? '全面检查' : '部分检查'
          }}</template>
        </ElTableColumn>
        <ElTableColumn prop="overallScore" label="合规分数" width="100">
          <template #default="{ row }">
            <ElTag :type="getScoreType(row.overallScore)">{{ row.overallScore.toFixed(1) }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="checkStatus" label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.checkStatus)">{{
              getStatusLabel(row.checkStatus)
            }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="createdAt" label="检查时间" width="180" />
        <ElTableColumn label="操作" width="120">
          <template #default="{ row }">
            <ElButton type="primary" link @click="viewHistoryDetail(row)">查看详情</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
      <ElPagination
        v-model:current-page="historyPagination.current"
        v-model:page-size="historyPagination.size"
        :total="historyPagination.total"
        layout="total, prev, pager, next"
        @current-change="loadHistory"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /*** Compliance Check Component - 合规检查页面组件 - Requirements: 11.1, 11.2, 11.3, 11.4 ***/
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { CircleCheck, CircleClose, Warning, Search, Clock } from '@element-plus/icons-vue'
  import { getDesignerProjects } from '@/api/designer-project'
  import {
    performComplianceCheck,
    getComplianceCheckList,
    getComplianceCheckSummary,
    STANDARD_CATEGORIES,
    SEVERITY_OPTIONS,
    PRIORITY_OPTIONS,
    type ComplianceCheckResult,
    type ComplianceCheckSummary
  } from '@/api/designer-compliance'

  const loading = ref(false)
  const activeTab = ref('passed')
  const historyDialogVisible = ref(false)
  const historyLoading = ref(false)
  const projectList = ref<any[]>([])
  const checkForm = reactive({
    projectId: null as number | null,
    categories: ['fire', 'accessibility', 'environmental', 'safety']
  })
  const checkResult = ref<ComplianceCheckResult | null>(null)
  const checkSummary = ref<ComplianceCheckSummary | null>(null)
  const historyList = ref<ComplianceCheckResult[]>([])
  const historyPagination = reactive({ current: 1, size: 10, total: 0 })

  const loadProjects = async () => {
    try {
      const res = await getDesignerProjects({ current: 1, size: 100 })
      if (res.code === 200) projectList.value = res.data?.records || []
    } catch (error) {
      console.error('加载项目列表失败:', error)
    }
  }

  const handleProjectChange = async (projectId: number) => {
    if (projectId) await loadCheckSummary(projectId)
    else checkSummary.value = null
    checkResult.value = null
  }

  const loadCheckSummary = async (projectId: number) => {
    try {
      const res = await getComplianceCheckSummary(projectId)
      if (res.code === 200) checkSummary.value = res.data
    } catch (error) {
      console.error('加载检查摘要失败:', error)
    }
  }

  const handleCheck = async () => {
    if (!checkForm.projectId) {
      ElMessage.warning('请选择项目')
      return
    }
    if (checkForm.categories.length === 0) {
      ElMessage.warning('请选择检查类别')
      return
    }
    loading.value = true
    try {
      const res = await performComplianceCheck({
        projectId: checkForm.projectId,
        categories: checkForm.categories,
        checkType: checkForm.categories.length === STANDARD_CATEGORIES.length ? 'full' : 'partial'
      })
      if (res.code === 200) {
        checkResult.value = res.data
        await loadCheckSummary(checkForm.projectId)
        ElMessage.success('检查完成')
      } else ElMessage.error(res.msg || '检查失败')
    } catch (error: any) {
      ElMessage.error(error.message || '检查失败')
    } finally {
      loading.value = false
    }
  }

  const showHistoryDialog = async () => {
    if (!checkForm.projectId) {
      ElMessage.warning('请先选择项目')
      return
    }
    historyDialogVisible.value = true
    await loadHistory()
  }

  const loadHistory = async () => {
    if (!checkForm.projectId) return
    historyLoading.value = true
    try {
      const res = await getComplianceCheckList({
        projectId: checkForm.projectId,
        current: historyPagination.current,
        size: historyPagination.size
      })
      if (res.code === 200) {
        historyList.value = res.data?.records || []
        historyPagination.total = res.data?.total || 0
      }
    } catch (error) {
      console.error('加载历史记录失败:', error)
    } finally {
      historyLoading.value = false
    }
  }

  const viewHistoryDetail = (row: ComplianceCheckResult) => {
    checkResult.value = row
    historyDialogVisible.value = false
    activeTab.value = 'passed'
  }
  const getCategoryLabel = (category: string): string =>
    STANDARD_CATEGORIES.find((c) => c.value === category)?.label || category
  const getSeverityType = (severity: string): 'success' | 'warning' | 'info' | 'danger' =>
    (SEVERITY_OPTIONS.find((o) => o.value === severity)?.type as any) || 'info'
  const getSeverityLabel = (severity: string): string =>
    SEVERITY_OPTIONS.find((o) => o.value === severity)?.label || severity
  const getPriorityType = (priority: string): 'success' | 'warning' | 'info' | 'danger' =>
    (PRIORITY_OPTIONS.find((o) => o.value === priority)?.type as any) || 'info'
  const getPriorityLabel = (priority: string): string =>
    PRIORITY_OPTIONS.find((o) => o.value === priority)?.label || priority
  const getScoreType = (score: number): 'success' | 'warning' | 'danger' =>
    score >= 80 ? 'success' : score >= 60 ? 'warning' : 'danger'
  const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' =>
    (({ completed: 'success', processing: 'warning', pending: 'info', failed: 'danger' }) as any)[
      status
    ] || 'info'
  const getStatusLabel = (status: string): string =>
    (({ completed: '已完成', processing: '检查中', pending: '待检查', failed: '检查失败' }) as any)[
      status
    ] || status

  onMounted(() => {
    loadProjects()
  })
</script>

<style scoped lang="scss">
  .compliance-check-page {
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
    .check-form {
      margin-bottom: 16px;
    }
    .summary-section {
      margin-bottom: 24px;
      padding: 16px;
      background: var(--el-fill-color-light);
      border-radius: 8px;
      .latest-check-time {
        margin-top: 12px;
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
    .result-section {
      .stats-row {
        margin-bottom: 24px;
      }
      .result-tabs {
        margin-top: 16px;
      }
      .failed-item {
        margin-bottom: 16px;
        .failed-detail {
          p {
            margin: 8px 0;
            line-height: 1.6;
          }
        }
      }
      .suggestion-item {
        margin-bottom: 16px;
        .suggestion-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        p {
          margin: 8px 0;
          color: var(--el-text-color-regular);
        }
      }
    }
  }
</style>
