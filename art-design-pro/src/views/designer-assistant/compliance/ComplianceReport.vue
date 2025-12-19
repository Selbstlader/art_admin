<template>
  <div class="compliance-report-page">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">合规报告详情</span>
          <div class="header-actions">
            <ElButton @click="handleExport"><ElIcon><Download /></ElIcon>导出报告</ElButton>
            <ElButton @click="goBack"><ElIcon><Back /></ElIcon>返回</ElButton>
          </div>
        </div>
      </template>

      <div v-if="reportData" class="report-content">
        <div class="report-header">
          <ElDescriptions :column="3" border>
            <ElDescriptionsItem label="项目名称">{{ reportData.projectName }}</ElDescriptionsItem>
            <ElDescriptionsItem label="检查类型">{{ reportData.checkType === 'full' ? '全面检查' : '部分检查' }}</ElDescriptionsItem>
            <ElDescriptionsItem label="检查时间">{{ reportData.createdAt }}</ElDescriptionsItem>
            <ElDescriptionsItem label="合规分数">
              <ElTag :type="getScoreType(reportData.overallScore)" size="large">{{ reportData.overallScore.toFixed(1) }} 分</ElTag>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="检查状态">
              <ElTag :type="getStatusType(reportData.checkStatus)">{{ getStatusLabel(reportData.checkStatus) }}</ElTag>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="检查规范数">{{ reportData.checkedStandards?.length || 0 }} 项</ElDescriptionsItem>
          </ElDescriptions>
        </div>

        <div class="report-charts">
          <ElRow :gutter="24">
            <ElCol :span="8">
              <ElCard shadow="hover" class="stat-card success">
                <div class="stat-icon"><ElIcon size="48"><CircleCheck /></ElIcon></div>
                <div class="stat-info"><div class="stat-value">{{ reportData.passedItems?.length || 0 }}</div><div class="stat-label">通过项</div></div>
              </ElCard>
            </ElCol>
            <ElCol :span="8">
              <ElCard shadow="hover" class="stat-card danger">
                <div class="stat-icon"><ElIcon size="48"><CircleClose /></ElIcon></div>
                <div class="stat-info"><div class="stat-value">{{ reportData.failedItems?.length || 0 }}</div><div class="stat-label">不合规项</div></div>
              </ElCard>
            </ElCol>
            <ElCol :span="8">
              <ElCard shadow="hover" class="stat-card warning">
                <div class="stat-icon"><ElIcon size="48"><Warning /></ElIcon></div>
                <div class="stat-info"><div class="stat-value">{{ reportData.suggestions?.length || 0 }}</div><div class="stat-label">改进建议</div></div>
              </ElCard>
            </ElCol>
          </ElRow>
        </div>

        <div class="report-section">
          <h3 class="section-title"><ElIcon color="#67c23a"><CircleCheck /></ElIcon>通过项 ({{ reportData.passedItems?.length || 0 }})</h3>
          <ElTable :data="reportData.passedItems" stripe border>
            <ElTableColumn prop="standardCode" label="规范编号" width="150" />
            <ElTableColumn prop="standardName" label="规范名称" min-width="200" />
            <ElTableColumn prop="category" label="类别" width="120">
              <template #default="{ row }"><ElTag size="small" type="success">{{ getCategoryLabel(row.category) }}</ElTag></template>
            </ElTableColumn>
            <ElTableColumn prop="description" label="符合说明" min-width="300" />
          </ElTable>
          <ElEmpty v-if="!reportData.passedItems?.length" description="暂无通过项" />
        </div>

        <div class="report-section">
          <h3 class="section-title"><ElIcon color="#f56c6c"><CircleClose /></ElIcon>不合规项 ({{ reportData.failedItems?.length || 0 }})</h3>
          <div v-for="(item, index) in reportData.failedItems" :key="index" class="failed-card">
            <ElCard shadow="never">
              <template #header>
                <div class="failed-header">
                  <span class="failed-title">[{{ item.standardCode }}] {{ item.standardName }}</span>
                  <ElTag :type="getSeverityType(item.severity)" size="small">{{ getSeverityLabel(item.severity) }}</ElTag>
                </div>
              </template>
              <ElDescriptions :column="1" border size="small">
                <ElDescriptionsItem label="类别">{{ getCategoryLabel(item.category) }}</ElDescriptionsItem>
                <ElDescriptionsItem label="问题位置">{{ item.location }}</ElDescriptionsItem>
                <ElDescriptionsItem label="违规内容">{{ item.violationContent }}</ElDescriptionsItem>
                <ElDescriptionsItem label="原始规范">{{ item.originalRequirement }}</ElDescriptionsItem>
                <ElDescriptionsItem label="修改建议"><span class="suggestion-text">{{ item.suggestion }}</span></ElDescriptionsItem>
              </ElDescriptions>
            </ElCard>
          </div>
          <ElEmpty v-if="!reportData.failedItems?.length" description="暂无不合规项" />
        </div>

        <div class="report-section">
          <h3 class="section-title"><ElIcon color="#e6a23c"><Warning /></ElIcon>改进建议 ({{ reportData.suggestions?.length || 0 }})</h3>
          <ElTimeline>
            <ElTimelineItem v-for="(item, index) in reportData.suggestions" :key="index" :type="getPriorityColor(item.priority)" :hollow="true">
              <ElCard shadow="never" class="suggestion-card">
                <div class="suggestion-content">
                  <div class="suggestion-main">
                    <span class="suggestion-text">{{ item.content }}</span>
                    <ElTag :type="getPriorityType(item.priority)" size="small">{{ getPriorityLabel(item.priority) }}优先级</ElTag>
                  </div>
                  <div class="suggestion-meta">
                    <span><strong>相关类别:</strong> {{ getCategoryLabel(item.category) }}</span>
                    <span><strong>成本影响:</strong> {{ item.costImpact }}</span>
                    <span><strong>参考规范:</strong> {{ item.reference }}</span>
                  </div>
                </div>
              </ElCard>
            </ElTimelineItem>
          </ElTimeline>
          <ElEmpty v-if="!reportData.suggestions?.length" description="暂无改进建议" />
        </div>
      </div>
      <ElEmpty v-else description="暂无报告数据" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
/*** Compliance Report Component - 合规报告详情组件 - Requirements: 11.2, 11.3, 11.4 ***/
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CircleCheck, CircleClose, Warning, Download, Back } from '@element-plus/icons-vue'
import { getComplianceCheckDetail, STANDARD_CATEGORIES, SEVERITY_OPTIONS, PRIORITY_OPTIONS, type ComplianceCheckResult } from '@/api/designer-compliance'

const route = useRoute()
const router = useRouter()
const reportData = ref<ComplianceCheckResult | null>(null)

const loadReport = async () => {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const res = await getComplianceCheckDetail(id)
    if (res.code === 200) reportData.value = res.data
    else ElMessage.error(res.msg || '加载报告失败')
  } catch (error: any) { ElMessage.error(error.message || '加载报告失败') }
}

const handleExport = () => { ElMessage.info('导出功能开发中') }
const goBack = () => { router.back() }
const getCategoryLabel = (category: string): string => STANDARD_CATEGORIES.find((c) => c.value === category)?.label || category
const getSeverityType = (severity: string): 'success' | 'warning' | 'info' | 'danger' => (SEVERITY_OPTIONS.find((o) => o.value === severity)?.type as any) || 'info'
const getSeverityLabel = (severity: string): string => SEVERITY_OPTIONS.find((o) => o.value === severity)?.label || severity
const getPriorityType = (priority: string): 'success' | 'warning' | 'info' | 'danger' => (PRIORITY_OPTIONS.find((o) => o.value === priority)?.type as any) || 'info'
const getPriorityLabel = (priority: string): string => PRIORITY_OPTIONS.find((o) => o.value === priority)?.label || priority
const getPriorityColor = (priority: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' => ({ high: 'danger', medium: 'warning', low: 'info' } as any)[priority] || 'info'
const getScoreType = (score: number): 'success' | 'warning' | 'danger' => score >= 80 ? 'success' : score >= 60 ? 'warning' : 'danger'
const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' => ({ completed: 'success', processing: 'warning', pending: 'info', failed: 'danger' } as any)[status] || 'info'
const getStatusLabel = (status: string): string => ({ completed: '已完成', processing: '检查中', pending: '待检查', failed: '检查失败' } as any)[status] || status

onMounted(() => { loadReport() })
</script>

<style scoped lang="scss">
.compliance-report-page {
  padding: 16px;
  .card-header { display: flex; justify-content: space-between; align-items: center; .title { font-size: 16px; font-weight: 500; } .header-actions { display: flex; gap: 8px; } }
  .report-content {
    .report-header { margin-bottom: 24px; }
    .report-charts { margin-bottom: 32px;
      .stat-card { display: flex; align-items: center; padding: 20px;
        &.success .stat-icon { color: #67c23a; }
        &.danger .stat-icon { color: #f56c6c; }
        &.warning .stat-icon { color: #e6a23c; }
        .stat-icon { margin-right: 16px; }
        .stat-info { .stat-value { font-size: 32px; font-weight: 600; line-height: 1.2; } .stat-label { font-size: 14px; color: var(--el-text-color-secondary); } }
      }
    }
    .report-section { margin-bottom: 32px;
      .section-title { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; font-size: 16px; font-weight: 500; }
      .failed-card { margin-bottom: 16px; .failed-header { display: flex; justify-content: space-between; align-items: center; .failed-title { font-weight: 500; } } .suggestion-text { color: #67c23a; font-weight: 500; } }
      .suggestion-card { .suggestion-content { .suggestion-main { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; .suggestion-text { flex: 1; margin-right: 12px; font-weight: 500; } } .suggestion-meta { display: flex; flex-wrap: wrap; gap: 16px; font-size: 13px; color: var(--el-text-color-secondary); } } }
    }
  }
}
</style>

            .suggestion-meta {
              display: flex;
              flex-wrap: wrap;
              gap: 16px;
              font-size: 13px;
              color: var(--el-text-color-secondary);
            }
          }
        }
      }
    }
  }
</style>
