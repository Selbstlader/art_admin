<template>
  <div class="material-select">
    <!-- 已选材料列表 / Selected materials list -->
    <el-card class="selected-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>已选材料 ({{ selectedMaterials.length }})</span>
          <el-button v-if="selectedMaterials.length > 0" type="danger" link @click="handleClearAll">
            清空
          </el-button>
        </div>
      </template>

      <el-table v-if="selectedMaterials.length > 0" :data="selectedMaterials" stripe>
        <el-table-column prop="name" label="材料名称" min-width="150" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column label="单价" width="120">
          <template #default="{ row }">
            ¥{{ row.unitPrice.toFixed(2) }}/{{ row.unit }}
          </template>
        </el-table-column>
        <el-table-column label="面积(m²)" width="140">
          <template #default="{ row }">
            <el-input-number
              v-model="row.area"
              :min="0"
              :precision="2"
              size="small"
              style="width: 100px"
              @change="handleCalculate"
            />
          </template>
        </el-table-column>
        <el-table-column label="损耗率" width="120">
          <template #default="{ row }">
            <el-input-number
              v-model="row.lossRate"
              :min="0"
              :max="1"
              :step="0.01"
              :precision="2"
              size="small"
              style="width: 80px"
              @change="handleCalculate"
            />
          </template>
        </el-table-column>
        <el-table-column label="用量" width="100">
          <template #default="{ row }">
            {{ row.quantity?.toFixed(2) || '-' }} {{ row.unit }}
          </template>
        </el-table-column>
        <el-table-column label="成本" width="120">
          <template #default="{ row }">
            <span class="cost">¥{{ row.totalCost?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row, $index }">
            <el-button type="danger" link size="small" @click="handleRemove($index)">
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-else description="暂未选择材料" />

      <!-- 成本汇总 / Cost summary -->
      <div v-if="selectedMaterials.length > 0" class="cost-summary">
        <div class="summary-item">
          <span class="label">材料总数:</span>
          <span class="value">{{ selectedMaterials.length }} 种</span>
        </div>
        <div class="summary-item">
          <span class="label">总用量:</span>
          <span class="value">{{ totalQuantity.toFixed(2) }}</span>
        </div>
        <div class="summary-item total">
          <span class="label">总成本:</span>
          <span class="value">¥{{ totalCost.toFixed(2) }}</span>
        </div>
      </div>
    </el-card>

    <!-- 材料库选择 / Material library selection -->
    <el-card class="library-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>材料库</span>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索材料"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </template>

      <!-- 分类筛选 / Category filter -->
      <div class="category-filter">
        <el-radio-group v-model="selectedCategory" @change="handleCategoryChange">
          <el-radio-button value="">全部</el-radio-button>
          <el-radio-button
            v-for="cat in categories"
            :key="cat.category"
            :value="cat.category"
          >
            {{ cat.category }}
          </el-radio-button>
        </el-radio-group>
      </div>

      <!-- 材料列表 / Material list -->
      <div v-loading="loading" class="material-list">
        <div
          v-for="material in materialList"
          :key="material.id"
          class="material-item"
          :class="{ selected: isSelected(material.id) }"
          @click="handleSelectMaterial(material)"
        >
          <div class="material-image">
            <el-image
              v-if="material.imageUrl"
              :src="material.imageUrl"
              fit="cover"
            />
            <div v-else class="no-image">暂无</div>
          </div>
          <div class="material-info">
            <div class="name">{{ material.name }}</div>
            <div class="spec">{{ material.specification || material.category }}</div>
            <div class="price">¥{{ material.unitPrice.toFixed(2) }}/{{ material.unit }}</div>
          </div>
          <el-icon v-if="isSelected(material.id)" class="check-icon"><Check /></el-icon>
        </div>
      </div>

      <!-- 分页 / Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[12, 24, 48]"
          layout="total, sizes, prev, pager, next"
          small
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { Search, Check } from '@element-plus/icons-vue'
import {
  getMaterialList,
  getMaterialCategories,
  batchCalculateMaterialUsage,
  type MaterialResponse,
  type MaterialCategoryResponse,
  type MaterialCalculateRequest,
  type MaterialCalculateResponse
} from '@/api/designer-material'

/*** Props ***/
interface Props {
  modelValue?: SelectedMaterial[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => []
})

/*** Emits ***/
const emit = defineEmits<{
  (e: 'update:modelValue', value: SelectedMaterial[]): void
  (e: 'change', value: SelectedMaterial[]): void
}>()

/*** 类型 / Types ***/
interface SelectedMaterial extends MaterialResponse {
  area: number
  lossRate: number
  quantity?: number
  totalCost?: number
}

/*** 数据 / Data ***/
const loading = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref('')
const materialList = ref<MaterialResponse[]>([])
const categories = ref<MaterialCategoryResponse[]>([])
const selectedMaterials = ref<SelectedMaterial[]>([])

const pagination = reactive({
  current: 1,
  size: 12,
  total: 0
})

/*** 计算属性 / Computed ***/
const totalCost = computed(() => {
  return selectedMaterials.value.reduce((sum, m) => sum + (m.totalCost || 0), 0)
})

const totalQuantity = computed(() => {
  return selectedMaterials.value.reduce((sum, m) => sum + (m.quantity || 0), 0)
})

/*** 方法 / Methods ***/

// 加载材料列表 / Load material list
const loadMaterialList = async () => {
  loading.value = true
  try {
    const res = await getMaterialList({
      current: pagination.current,
      size: pagination.size,
      keyword: searchKeyword.value || undefined,
      category: selectedCategory.value || undefined,
      status: 'active'
    }) as unknown as Http.BaseResponse<{ records: MaterialResponse[]; total: number }>
    materialList.value = res.data?.records || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error('加载材料列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载分类 / Load categories
const loadCategories = async () => {
  try {
    const res = await getMaterialCategories() as unknown as Http.BaseResponse<MaterialCategoryResponse[]>
    categories.value = res.data || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 搜索 / Search
const handleSearch = () => {
  pagination.current = 1
  loadMaterialList()
}

// 分类变化 / Category change
const handleCategoryChange = () => {
  pagination.current = 1
  loadMaterialList()
}

// 分页变化 / Pagination change
const handleSizeChange = () => {
  pagination.current = 1
  loadMaterialList()
}

const handleCurrentChange = () => {
  loadMaterialList()
}

// 检查是否已选 / Check if selected
const isSelected = (id: number) => {
  return selectedMaterials.value.some(m => m.id === id)
}

// 选择材料 / Select material
const handleSelectMaterial = (material: MaterialResponse) => {
  if (isSelected(material.id)) {
    // 已选则移除 / Remove if already selected
    const index = selectedMaterials.value.findIndex(m => m.id === material.id)
    if (index > -1) {
      selectedMaterials.value.splice(index, 1)
    }
  } else {
    // 未选则添加 / Add if not selected
    selectedMaterials.value.push({
      ...material,
      area: 0,
      lossRate: 0.05 // 默认5%损耗 / Default 5% loss
    })
  }
  emitChange()
}

// 移除材料 / Remove material
const handleRemove = (index: number) => {
  selectedMaterials.value.splice(index, 1)
  emitChange()
}

// 清空所有 / Clear all
const handleClearAll = () => {
  selectedMaterials.value = []
  emitChange()
}

// 计算用量和成本 / Calculate usage and cost
const handleCalculate = async () => {
  const items: MaterialCalculateRequest[] = selectedMaterials.value
    .filter(m => m.area > 0)
    .map(m => ({
      materialId: m.id,
      area: m.area,
      lossRate: m.lossRate
    }))

  if (items.length === 0) {
    // 清空计算结果 / Clear calculation results
    selectedMaterials.value.forEach(m => {
      m.quantity = undefined
      m.totalCost = undefined
    })
    return
  }

  try {
    const res = await batchCalculateMaterialUsage(items) as unknown as Http.BaseResponse<{ items: MaterialCalculateResponse[] }>
    // 更新计算结果 / Update calculation results
    const resultItems = res.data?.items || []
    resultItems.forEach((result: MaterialCalculateResponse) => {
      const material = selectedMaterials.value.find((m) => m.id === result.materialId)
      if (material) {
        material.quantity = result.quantity
        material.totalCost = result.totalCost
      }
    })
    emitChange()
  } catch (error) {
    console.error('计算失败:', error)
  }
}

// 触发变更事件 / Emit change event
const emitChange = () => {
  emit('update:modelValue', selectedMaterials.value)
  emit('change', selectedMaterials.value)
}

// 监听外部值变化 / Watch external value changes
watch(() => props.modelValue, (newVal) => {
  if (newVal && newVal !== selectedMaterials.value) {
    selectedMaterials.value = [...newVal]
  }
}, { immediate: true })

/*** 生命周期 / Lifecycle ***/
onMounted(() => {
  loadMaterialList()
  loadCategories()
})
</script>

<style scoped lang="scss">
.material-select {
  display: flex;
  gap: 16px;
  height: 100%;

  .selected-card {
    flex: 1;
    min-width: 500px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .cost {
      color: #f56c6c;
      font-weight: 500;
    }

    .cost-summary {
      margin-top: 16px;
      padding: 16px;
      background: #f5f7fa;
      border-radius: 4px;
      display: flex;
      justify-content: flex-end;
      gap: 32px;

      .summary-item {
        .label {
          color: #909399;
          margin-right: 8px;
        }

        .value {
          font-weight: 500;
        }

        &.total .value {
          color: #f56c6c;
          font-size: 18px;
        }
      }
    }
  }

  .library-card {
    width: 400px;
    flex-shrink: 0;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .category-filter {
      margin-bottom: 16px;
      overflow-x: auto;

      :deep(.el-radio-group) {
        flex-wrap: nowrap;
      }
    }

    .material-list {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;
      min-height: 300px;

      .material-item {
        position: relative;
        padding: 12px;
        border: 1px solid #e4e7ed;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          border-color: #409eff;
        }

        &.selected {
          border-color: #67c23a;
          background: #f0f9eb;
        }

        .material-image {
          height: 80px;
          margin-bottom: 8px;
          border-radius: 4px;
          overflow: hidden;

          :deep(.el-image) {
            width: 100%;
            height: 100%;
          }

          .no-image {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f5f7fa;
            color: #909399;
            font-size: 12px;
          }
        }

        .material-info {
          .name {
            font-weight: 500;
            font-size: 13px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .spec {
            color: #909399;
            font-size: 12px;
            margin: 4px 0;
          }

          .price {
            color: #f56c6c;
            font-size: 13px;
            font-weight: 500;
          }
        }

        .check-icon {
          position: absolute;
          top: 8px;
          right: 8px;
          color: #67c23a;
          font-size: 18px;
        }
      }
    }

    .pagination-wrapper {
      margin-top: 16px;
      display: flex;
      justify-content: center;
    }
  }
}
</style>
