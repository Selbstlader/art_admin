<template>
  <div class="cost-config-page">
    <ElRow :gutter="16">
      <!-- 左侧配置区域 -->
      <ElCol :span="16">
        <ElCard shadow="never">
          <template #header>
            <div class="card-header">
              <span class="title">成本配置</span>
            </div>
          </template>

          <!-- 项目选择 -->
          <ElForm :model="configForm" label-width="100px">
            <ElFormItem label="关联项目">
              <ElSelect v-model="configForm.projectId" placeholder="请选择项目" style="width: 100%">
                <ElOption
                  v-for="project in projectList"
                  :key="project.id"
                  :label="project.name"
                  :value="project.id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="预算上限">
              <ElInputNumber
                v-model="configForm.budgetLimit"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </ElFormItem>
          </ElForm>

          <ElDivider content-position="left">材料费用</ElDivider>

          <!-- 材料选择表格 -->
          <ElTable :data="selectedMaterials" stripe>
            <ElTableColumn prop="name" label="材料名称" />
            <ElTableColumn prop="unit" label="单位" width="80" />
            <ElTableColumn prop="unitPrice" label="单价" width="120">
              <template #default="{ row }"> ¥{{ row.unitPrice?.toFixed(2) }} </template>
            </ElTableColumn>
            <ElTableColumn label="数量" width="150">
              <template #default="{ row }">
                <ElInputNumber
                  v-model="row.quantity"
                  :min="0"
                  :precision="2"
                  size="small"
                  @change="calculateCost"
                />
              </template>
            </ElTableColumn>
            <ElTableColumn label="小计" width="120">
              <template #default="{ row }">
                ¥{{ (row.unitPrice * row.quantity).toFixed(2) }}
              </template>
            </ElTableColumn>
            <ElTableColumn label="操作" width="80">
              <template #default="{ $index }">
                <ElButton type="danger" link @click="removeMaterial($index)">删除</ElButton>
              </template>
            </ElTableColumn>
          </ElTable>

          <ElButton type="primary" link @click="showMaterialDialog" style="margin-top: 12px">
            <ElIcon><Plus /></ElIcon>
            添加材料
          </ElButton>

          <ElDivider content-position="left">其他费用</ElDivider>

          <ElForm :model="configForm" label-width="100px">
            <ElFormItem label="人工费">
              <ElInputNumber
                v-model="configForm.laborCost"
                :min="0"
                :precision="2"
                style="width: 100%"
                @change="calculateCost"
              />
            </ElFormItem>
            <ElFormItem label="设备费">
              <ElInputNumber
                v-model="configForm.equipmentCost"
                :min="0"
                :precision="2"
                style="width: 100%"
                @change="calculateCost"
              />
            </ElFormItem>
            <ElFormItem label="管理费">
              <ElInputNumber
                v-model="configForm.managementCost"
                :min="0"
                :precision="2"
                style="width: 100%"
                @change="calculateCost"
              />
            </ElFormItem>
          </ElForm>

          <div class="action-area">
            <ElButton type="primary" @click="handleSave" :loading="loading">保存配置</ElButton>
            <ElButton @click="handleExport">导出报告</ElButton>
          </div>
        </ElCard>
      </ElCol>

      <!-- 右侧汇总区域 -->
      <ElCol :span="8">
        <ElCard shadow="never" class="summary-card">
          <template #header>
            <div class="card-header">
              <span class="title">成本汇总</span>
            </div>
          </template>

          <div class="summary-item">
            <span class="label">材料费</span>
            <span class="value">¥{{ costSummary.materialCost.toFixed(2) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">人工费</span>
            <span class="value">¥{{ costSummary.laborCost.toFixed(2) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">设备费</span>
            <span class="value">¥{{ costSummary.equipmentCost.toFixed(2) }}</span>
          </div>
          <div class="summary-item">
            <span class="label">管理费</span>
            <span class="value">¥{{ costSummary.managementCost.toFixed(2) }}</span>
          </div>

          <ElDivider />

          <div class="summary-item total">
            <span class="label">总计</span>
            <span class="value">¥{{ costSummary.totalCost.toFixed(2) }}</span>
          </div>

          <!-- 预算警告 -->
          <ElAlert
            v-if="configForm.budgetLimit > 0 && costSummary.totalCost > configForm.budgetLimit"
            title="预算超出警告"
            type="warning"
            :closable="false"
            show-icon
            style="margin-top: 16px"
          >
            <template #default>
              超出预算 ¥{{ (costSummary.totalCost - configForm.budgetLimit).toFixed(2) }}
            </template>
          </ElAlert>

          <!-- 预算进度 -->
          <div v-if="configForm.budgetLimit > 0" class="budget-progress">
            <span class="label">预算使用率</span>
            <ElProgress
              :percentage="Math.min(100, (costSummary.totalCost / configForm.budgetLimit) * 100)"
              :color="getBudgetColor"
            />
          </div>
        </ElCard>
      </ElCol>
    </ElRow>

    <!-- 材料选择对话框 -->
    <ElDialog v-model="materialDialogVisible" title="选择材料" width="800px">
      <ElTable :data="materialOptions" @selection-change="handleMaterialSelect">
        <ElTableColumn type="selection" width="55" />
        <ElTableColumn prop="name" label="材料名称" />
        <ElTableColumn prop="category" label="分类" width="100" />
        <ElTableColumn prop="unit" label="单位" width="80" />
        <ElTableColumn prop="unitPrice" label="单价" width="120">
          <template #default="{ row }"> ¥{{ row.unitPrice?.toFixed(2) }} </template>
        </ElTableColumn>
      </ElTable>
      <template #footer>
        <ElButton @click="materialDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="confirmMaterialSelect">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  /***
   * Cost Config Component
   * 成本配置页面组件
   * Requirements: 10.1, 10.2, 10.3, 10.4
   ***/
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Plus } from '@element-plus/icons-vue'

  const loading = ref(false)
  const materialDialogVisible = ref(false)

  // 项目列表
  const projectList = ref<any[]>([])

  // 配置表单
  const configForm = reactive({
    projectId: null as number | null,
    budgetLimit: 0,
    laborCost: 0,
    equipmentCost: 0,
    managementCost: 0
  })

  // 已选材料
  const selectedMaterials = ref<any[]>([])

  // 材料选项
  const materialOptions = ref<any[]>([])
  const tempSelectedMaterials = ref<any[]>([])

  // 成本汇总
  const costSummary = computed(() => {
    const materialCost = selectedMaterials.value.reduce(
      (sum, item) => sum + (item.unitPrice || 0) * (item.quantity || 0),
      0
    )
    const laborCost = configForm.laborCost || 0
    const equipmentCost = configForm.equipmentCost || 0
    const managementCost = configForm.managementCost || 0
    const totalCost = materialCost + laborCost + equipmentCost + managementCost

    return {
      materialCost,
      laborCost,
      equipmentCost,
      managementCost,
      totalCost
    }
  })

  // 预算颜色
  const getBudgetColor = computed(() => {
    const percentage = (costSummary.value.totalCost / configForm.budgetLimit) * 100
    if (percentage >= 100) return '#f56c6c'
    if (percentage >= 80) return '#e6a23c'
    return '#67c23a'
  })

  // 计算成本
  const calculateCost = () => {
    // 触发computed重新计算
  }

  // 显示材料选择对话框
  const showMaterialDialog = () => {
    // TODO: 获取材料列表
    materialDialogVisible.value = true
  }

  // 材料选择变化
  const handleMaterialSelect = (selection: any[]) => {
    tempSelectedMaterials.value = selection
  }

  // 确认材料选择
  const confirmMaterialSelect = () => {
    const newMaterials = tempSelectedMaterials.value.map((m) => ({
      ...m,
      quantity: 1
    }))
    selectedMaterials.value.push(...newMaterials)
    materialDialogVisible.value = false
  }

  // 移除材料
  const removeMaterial = (index: number) => {
    selectedMaterials.value.splice(index, 1)
  }

  // 保存配置
  const handleSave = async () => {
    if (!configForm.projectId) {
      ElMessage.warning('请选择关联项目')
      return
    }

    loading.value = true
    try {
      // TODO: 调用后端接口
      // await saveCostEstimate({
      //   ...configForm,
      //   ...costSummary.value,
      //   items: selectedMaterials.value
      // })
      ElMessage.success('保存成功')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  }

  // 导出报告
  const handleExport = () => {
    // TODO: 导出Excel或PDF
    ElMessage.info('导出功能开发中')
  }
</script>

<style scoped lang="scss">
  .cost-config-page {
    padding: 16px;

    .card-header {
      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .action-area {
      margin-top: 24px;
      text-align: center;
    }

    .summary-card {
      position: sticky;
      top: 16px;

      .summary-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 0;
        border-bottom: 1px solid var(--el-border-color-lighter);

        &:last-child {
          border-bottom: none;
        }

        &.total {
          .label,
          .value {
            font-size: 18px;
            font-weight: 600;
          }

          .value {
            color: var(--el-color-primary);
          }
        }

        .label {
          color: var(--el-text-color-secondary);
        }

        .value {
          font-weight: 500;
        }
      }

      .budget-progress {
        margin-top: 16px;

        .label {
          display: block;
          margin-bottom: 8px;
          font-size: 14px;
          color: var(--el-text-color-secondary);
        }
      }
    }
  }
</style>
