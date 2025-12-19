<template>
  <div class="project-cost-summary">
    <ElCard shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">成本汇总</span>
          <div class="header-actions">
            <ElButton type="primary" link @click="handleReadFromDocument">
              <ElIcon><Document /></ElIcon>
              从文档读取
            </ElButton>
            <ElButton type="primary" link @click="handleExport">
              <ElIcon><Download /></ElIcon>
              导出报告
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 预算设置 / Budget setting -->
      <div class="budget-setting">
        <span class="label">预算上限</span>
        <ElInputNumber
          v-model="costConfig.budgetLimit"
          :min="0"
          :precision="2"
          :controls="false"
          placeholder="输入预算上限"
          style="width: 200px"
          @change="handleConfigChange"
        />
        <span class="unit">元</span>
      </div>

      <ElDivider />

      <!-- 费用明细 / Cost details -->
      <div class="cost-items">
        <div class="cost-item">
          <span class="label">材料费</span>
          <span class="value">¥{{ materialCost.toFixed(2) }}</span>
        </div>
        <div class="cost-item">
          <span class="label">人工费</span>
          <div class="input-wrapper">
            <ElInputNumber
              v-model="costConfig.laborCost"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 150px"
              @change="handleConfigChange"
            />
            <span class="unit">元</span>
          </div>
        </div>
        <div class="cost-item">
          <span class="label">设备费</span>
          <div class="input-wrapper">
            <ElInputNumber
              v-model="costConfig.equipmentCost"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 150px"
              @change="handleConfigChange"
            />
            <span class="unit">元</span>
          </div>
        </div>
        <div class="cost-item">
          <span class="label">管理费</span>
          <div class="input-wrapper">
            <ElInputNumber
              v-model="costConfig.managementCost"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 150px"
              @change="handleConfigChange"
            />
            <span class="unit">元</span>
          </div>
        </div>

        <!-- 自定义费用项 / Custom cost items -->
        <div v-for="(item, index) in customCostItems" :key="index" class="cost-item custom-item">
          <div class="custom-name">
            <ElInput
              v-model="item.name"
              placeholder="费用名称"
              size="small"
              style="width: 80px"
              @change="handleConfigChange"
            />
          </div>
          <div class="input-wrapper">
            <ElInputNumber
              v-model="item.amount"
              :precision="2"
              :controls="false"
              size="small"
              style="width: 150px"
              @change="handleConfigChange"
            />
            <span class="unit">元</span>
            <ElButton type="danger" link size="small" @click="removeCustomItem(index)">
              <ElIcon><Delete /></ElIcon>
            </ElButton>
          </div>
        </div>

        <!-- 添加自定义费用 / Add custom cost -->
        <div class="add-custom">
          <ElButton type="primary" link size="small" @click="addCustomItem">
            <ElIcon><Plus /></ElIcon>
            添加自定义费用（支持负数进行扣减）
          </ElButton>
        </div>
      </div>

      <ElDivider />

      <!-- 自定义费用小计 / Custom cost subtotal -->
      <div v-if="customCostItems.length > 0" class="subtotal-section">
        <div class="subtotal-item">
          <span class="label">自定义费用小计</span>
          <span :class="['value', customCostTotal >= 0 ? 'positive' : 'negative']">
            {{ customCostTotal >= 0 ? '+' : '' }}¥{{ customCostTotal.toFixed(2) }}
          </span>
        </div>
      </div>

      <!-- 总计 / Total -->
      <div class="total-section">
        <div class="total-item">
          <span class="label">总计</span>
          <span class="value">¥{{ totalCost.toFixed(2) }}</span>
        </div>
      </div>

      <!-- 预算警告 / Budget warning -->
      <ElAlert
        v-if="costConfig.budgetLimit > 0 && totalCost > costConfig.budgetLimit"
        title="预算超出警告"
        type="warning"
        :closable="false"
        show-icon
        style="margin-top: 16px"
      >
        <template #default>
          超出预算 ¥{{ (totalCost - costConfig.budgetLimit).toFixed(2) }}
        </template>
      </ElAlert>

      <!-- 预算进度 / Budget progress -->
      <div v-if="costConfig.budgetLimit > 0" class="budget-progress">
        <div class="progress-header">
          <span class="label">预算使用率</span>
          <span class="percentage">{{ budgetPercentage.toFixed(1) }}%</span>
        </div>
        <ElProgress
          :percentage="Math.min(100, budgetPercentage)"
          :color="budgetColor"
          :show-text="false"
        />
      </div>

      <!-- 保存按钮 / Save button -->
      <div class="save-action">
        <ElButton type="primary" @click="onSaveClick" :loading="saving">保存成本配置</ElButton>
      </div>
    </ElCard>

    <!-- 从文档读取弹窗 / Read from document dialog -->
    <ElDialog v-model="documentDialogVisible" title="从项目文档读取费用" width="600px">
      <div class="document-dialog">
        <ElAlert
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #title>
            系统将从项目关联的文档中提取费用信息，提取的数据将覆盖当前输入的值
          </template>
        </ElAlert>

        <div v-if="documentLoading" class="loading-wrapper">
          <ElIcon class="is-loading"><Loading /></ElIcon>
          <span>正在分析文档...</span>
        </div>

        <div v-else-if="extractedCosts" class="extracted-costs">
          <h4>从文档中提取到以下费用信息：</h4>
          <div class="cost-preview">
            <div class="preview-item" v-if="extractedCosts.laborCost !== undefined">
              <span class="label">人工费</span>
              <span class="value">¥{{ extractedCosts.laborCost.toFixed(2) }}</span>
            </div>
            <div class="preview-item" v-if="extractedCosts.equipmentCost !== undefined">
              <span class="label">设备费</span>
              <span class="value">¥{{ extractedCosts.equipmentCost.toFixed(2) }}</span>
            </div>
            <div class="preview-item" v-if="extractedCosts.managementCost !== undefined">
              <span class="label">管理费</span>
              <span class="value">¥{{ extractedCosts.managementCost.toFixed(2) }}</span>
            </div>
            <div class="preview-item" v-if="extractedCosts.budgetLimit !== undefined">
              <span class="label">预算上限</span>
              <span class="value">¥{{ extractedCosts.budgetLimit.toFixed(2) }}</span>
            </div>
          </div>
          <ElEmpty v-if="!hasExtractedData" description="未从文档中提取到费用信息" />
        </div>
      </div>
      <template #footer>
        <ElButton @click="documentDialogVisible = false">取消</ElButton>
        <ElButton
          type="primary"
          :disabled="!hasExtractedData"
          @click="applyExtractedCosts"
        >
          应用到成本配置
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * Project Cost Summary Component
   * 项目成本汇总组件 - 支持从文档读取和自定义费用项
   ***/
  import { ref, reactive, computed, watch, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Download, Document, Plus, Delete, Loading } from '@element-plus/icons-vue'
  import { getProjectCostSummary, saveProjectCost } from '@/api/designer-project-material'
  import { exportCostReport } from '@/api/designer-cost'
  import { getDocumentList } from '@/api/designer-document'

  const props = defineProps<{
    projectId: number
    materialCost: number
  }>()

  const saving = ref(false)

  const costConfig = reactive({
    budgetLimit: 0,
    laborCost: 0,
    equipmentCost: 0,
    managementCost: 0
  })

  // 自定义费用项 / Custom cost items
  interface CustomCostItem {
    name: string
    amount: number
  }
  const customCostItems = ref<CustomCostItem[]>([])

  // 自定义费用小计 / Custom cost subtotal
  const customCostTotal = computed(() => {
    return customCostItems.value.reduce((sum, item) => sum + (item.amount || 0), 0)
  })

  // 总成本 / Total cost
  const totalCost = computed(() => {
    return (
      props.materialCost +
      costConfig.laborCost +
      costConfig.equipmentCost +
      costConfig.managementCost +
      customCostTotal.value
    )
  })

  // 预算使用百分比 / Budget percentage
  const budgetPercentage = computed(() => {
    if (costConfig.budgetLimit <= 0) return 0
    return (totalCost.value / costConfig.budgetLimit) * 100
  })

  // 预算颜色 / Budget color
  const budgetColor = computed(() => {
    if (budgetPercentage.value >= 100) return '#f56c6c'
    if (budgetPercentage.value >= 80) return '#e6a23c'
    return '#67c23a'
  })

  // 添加自定义费用项 / Add custom cost item
  const addCustomItem = () => {
    customCostItems.value.push({ name: '', amount: 0 })
  }

  // 删除自定义费用项 / Remove custom cost item
  const removeCustomItem = (index: number) => {
    customCostItems.value.splice(index, 1)
    handleConfigChange()
  }

  // 加载成本配置 / Load cost config
  const loadCostConfig = async () => {
    if (!props.projectId) return
    try {
      const res: any = await getProjectCostSummary(props.projectId)
      if (res.code === 200 && res.data) {
        costConfig.budgetLimit = res.data.budgetLimit || 0
        costConfig.laborCost = res.data.laborCost || 0
        costConfig.equipmentCost = res.data.equipmentCost || 0
        costConfig.managementCost = res.data.managementCost || 0
        // 加载自定义费用项 / Load custom cost items
        if (res.data.customCosts && Array.isArray(res.data.customCosts)) {
          customCostItems.value = res.data.customCosts
        }
      }
    } catch {
      console.error('加载成本配置失败')
    }
  }

  // 配置变化时自动保存 / Auto save on config change
  let saveTimer: ReturnType<typeof setTimeout> | null = null
  const handleConfigChange = () => {
    if (saveTimer) clearTimeout(saveTimer)
    saveTimer = setTimeout(() => {
      handleSave(true)
    }, 1000)
  }

  // 保存成本配置 / Save cost config
  const handleSave = async (silent = false) => {
    if (!props.projectId) return
    saving.value = true
    try {
      await saveProjectCost({
        projectId: props.projectId,
        budgetLimit: costConfig.budgetLimit,
        laborCost: costConfig.laborCost,
        equipmentCost: costConfig.equipmentCost,
        managementCost: costConfig.managementCost,
        customCosts: customCostItems.value.filter(item => item.name || item.amount !== 0)
      })
      if (!silent) {
        ElMessage.success('保存成功')
      }
    } catch {
      if (!silent) {
        ElMessage.error('保存失败')
      }
    } finally {
      saving.value = false
    }
  }

  // 按钮点击保存 / Button click save
  const onSaveClick = () => {
    handleSave(false)
  }

  // 导出报告 / Export report
  const handleExport = async () => {
    try {
      const res: any = await exportCostReport({
        projectId: props.projectId,
        format: 'excel'
      })
      if (res.code === 200 && res.data?.fileUrl) {
        window.open(res.data.fileUrl, '_blank')
        ElMessage.success('导出成功')
      }
    } catch {
      ElMessage.error('导出失败')
    }
  }

  // 从文档读取相关 / Read from document related
  const documentDialogVisible = ref(false)
  const documentLoading = ref(false)
  const extractedCosts = ref<{
    laborCost?: number
    equipmentCost?: number
    managementCost?: number
    budgetLimit?: number
  } | null>(null)

  const hasExtractedData = computed(() => {
    if (!extractedCosts.value) return false
    return (
      extractedCosts.value.laborCost !== undefined ||
      extractedCosts.value.equipmentCost !== undefined ||
      extractedCosts.value.managementCost !== undefined ||
      extractedCosts.value.budgetLimit !== undefined
    )
  })

  // 从文档读取费用 / Read costs from document
  const handleReadFromDocument = async () => {
    documentDialogVisible.value = true
    documentLoading.value = true
    extractedCosts.value = null

    try {
      // 获取项目关联的文档列表 / Get project documents
      const res: any = await getDocumentList({
        projectId: props.projectId,
        current: 1,
        size: 100,
        analysisStatus: 'completed'
      })

      if (res.code === 200 && res.data?.records?.length > 0) {
        // 从文档中提取费用信息 / Extract cost info from documents
        const costs: {
          laborCost?: number
          equipmentCost?: number
          managementCost?: number
          budgetLimit?: number
        } = {}

        for (const doc of res.data.records) {
          // 从文档摘要或关键词中提取费用 / Extract from summary or keywords
          if (doc.summary) {
            const laborMatch = doc.summary.match(/人工费[：:]\s*([\d,.]+)\s*[元万]?/i)
            const equipmentMatch = doc.summary.match(/设备费[：:]\s*([\d,.]+)\s*[元万]?/i)
            const managementMatch = doc.summary.match(/管理费[：:]\s*([\d,.]+)\s*[元万]?/i)
            const budgetMatch = doc.summary.match(/预算[：:]\s*([\d,.]+)\s*[元万]?/i)

            if (laborMatch && costs.laborCost === undefined) {
              costs.laborCost = parseFloat(laborMatch[1].replace(/,/g, ''))
              if (doc.summary.includes('万')) costs.laborCost *= 10000
            }
            if (equipmentMatch && costs.equipmentCost === undefined) {
              costs.equipmentCost = parseFloat(equipmentMatch[1].replace(/,/g, ''))
              if (doc.summary.includes('万')) costs.equipmentCost *= 10000
            }
            if (managementMatch && costs.managementCost === undefined) {
              costs.managementCost = parseFloat(managementMatch[1].replace(/,/g, ''))
              if (doc.summary.includes('万')) costs.managementCost *= 10000
            }
            if (budgetMatch && costs.budgetLimit === undefined) {
              costs.budgetLimit = parseFloat(budgetMatch[1].replace(/,/g, ''))
              if (doc.summary.includes('万')) costs.budgetLimit *= 10000
            }
          }

          // 从提取的预算字段获取 / Get from extracted budget field
          if (doc.extractedBudget && costs.budgetLimit === undefined) {
            costs.budgetLimit = doc.extractedBudget
          }
        }

        extractedCosts.value = costs
      } else {
        extractedCosts.value = {}
      }
    } catch (error) {
      console.error('读取文档失败:', error)
      ElMessage.error('读取文档失败')
      extractedCosts.value = {}
    } finally {
      documentLoading.value = false
    }
  }

  // 应用提取的费用 / Apply extracted costs
  const applyExtractedCosts = () => {
    if (!extractedCosts.value) return

    if (extractedCosts.value.laborCost !== undefined) {
      costConfig.laborCost = extractedCosts.value.laborCost
    }
    if (extractedCosts.value.equipmentCost !== undefined) {
      costConfig.equipmentCost = extractedCosts.value.equipmentCost
    }
    if (extractedCosts.value.managementCost !== undefined) {
      costConfig.managementCost = extractedCosts.value.managementCost
    }
    if (extractedCosts.value.budgetLimit !== undefined) {
      costConfig.budgetLimit = extractedCosts.value.budgetLimit
    }

    documentDialogVisible.value = false
    handleConfigChange()
    ElMessage.success('已应用文档中的费用数据')
  }

  // 监听projectId变化 / Watch projectId change
  watch(
    () => props.projectId,
    (val) => {
      if (val) loadCostConfig()
    },
    { immediate: true }
  )

  onMounted(() => {
    if (props.projectId) loadCostConfig()
  })

  // 暴露方法 / Expose methods
  defineExpose({
    loadCostConfig,
    handleSave,
    getTotalCost: () => totalCost.value
  })
</script>

<style scoped lang="scss">
  .project-cost-summary {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }

      .header-actions {
        display: flex;
        gap: 8px;
      }
    }

    .budget-setting {
      display: flex;
      align-items: center;
      gap: 12px;

      .label {
        color: var(--el-text-color-secondary);
      }

      .unit {
        color: var(--el-text-color-secondary);
      }
    }

    .cost-items {
      .cost-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 0;
        border-bottom: 1px solid var(--el-border-color-lighter);

        &:last-child {
          border-bottom: none;
        }

        .label {
          color: var(--el-text-color-secondary);
        }

        .value {
          font-weight: 500;
          color: #f56c6c;
        }

        .input-wrapper {
          display: flex;
          align-items: center;
          gap: 8px;

          .unit {
            color: var(--el-text-color-secondary);
            font-size: 12px;
          }
        }

        &.custom-item {
          background: #f5f7fa;
          margin: 4px -12px;
          padding: 12px;
          border-radius: 4px;
          border-bottom: none;

          .custom-name {
            flex-shrink: 0;
          }
        }
      }

      .add-custom {
        padding: 12px 0;
        text-align: center;
      }
    }

    .subtotal-section {
      margin-bottom: 12px;

      .subtotal-item {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .label {
          font-size: 14px;
          color: var(--el-text-color-secondary);
        }

        .value {
          font-size: 14px;
          font-weight: 500;

          &.positive {
            color: #67c23a;
          }

          &.negative {
            color: #f56c6c;
          }
        }
      }
    }

    .total-section {
      .total-item {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .label {
          font-size: 16px;
          font-weight: 600;
        }

        .value {
          font-size: 20px;
          font-weight: 600;
          color: var(--el-color-primary);
        }
      }
    }

    .budget-progress {
      margin-top: 16px;

      .progress-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;

        .label {
          font-size: 14px;
          color: var(--el-text-color-secondary);
        }

        .percentage {
          font-size: 14px;
          font-weight: 500;
        }
      }
    }

    .save-action {
      margin-top: 24px;
      text-align: center;
    }
  }

  /*** 文档读取弹窗样式 / Document dialog styles ***/
  .document-dialog {
    .loading-wrapper {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 40px;
      color: var(--el-text-color-secondary);

      .el-icon {
        font-size: 32px;
        margin-bottom: 12px;
        color: var(--el-color-primary);
      }
    }

    .extracted-costs {
      h4 {
        margin: 0 0 16px;
        font-size: 14px;
        color: var(--el-text-color-primary);
      }

      .cost-preview {
        background: #f5f7fa;
        border-radius: 8px;
        padding: 16px;

        .preview-item {
          display: flex;
          justify-content: space-between;
          padding: 8px 0;
          border-bottom: 1px dashed var(--el-border-color-lighter);

          &:last-child {
            border-bottom: none;
          }

          .label {
            color: var(--el-text-color-secondary);
          }

          .value {
            font-weight: 500;
            color: var(--el-color-primary);
          }
        }
      }
    }
  }
</style>
