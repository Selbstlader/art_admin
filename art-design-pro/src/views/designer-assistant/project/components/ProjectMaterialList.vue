<template>
  <div class="project-material-list">
    <!-- 操作栏 / Action bar -->
    <div class="action-bar">
      <ElButton type="primary" @click="handleAddMaterial">
        <ElIcon><Plus /></ElIcon>
        添加材料
      </ElButton>
      <ElButton @click="showImportDialog">
        <ElIcon><FolderOpened /></ElIcon>
        从材料库导入
      </ElButton>
      <ElButton @click="handleExport" :disabled="materials.length === 0">
        <ElIcon><Download /></ElIcon>
        导出清单
      </ElButton>
    </div>

    <!-- 材料清单表格 / Material list table -->
    <ElTable :data="materials" v-loading="loading" stripe border>
      <ElTableColumn label="材料名称" min-width="150">
        <template #default="{ row, $index }">
          <ElInput
            v-if="row.editing"
            v-model="row.name"
            placeholder="输入材料名称"
            size="small"
          />
          <span v-else>{{ row.name }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="分类" width="120">
        <template #default="{ row }">
          <ElInput
            v-if="row.editing"
            v-model="row.category"
            placeholder="分类"
            size="small"
          />
          <span v-else>{{ row.category || '-' }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="规格" width="120">
        <template #default="{ row }">
          <ElInput
            v-if="row.editing"
            v-model="row.specification"
            placeholder="规格"
            size="small"
          />
          <span v-else>{{ row.specification || '-' }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="单位" width="80">
        <template #default="{ row }">
          <ElInput
            v-if="row.editing"
            v-model="row.unit"
            placeholder="单位"
            size="small"
          />
          <span v-else>{{ row.unit }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="单价(元)" width="120">
        <template #default="{ row }">
          <ElInputNumber
            v-if="row.editing"
            v-model="row.unitPrice"
            :min="0"
            :precision="2"
            :controls="false"
            size="small"
            style="width: 100%"
          />
          <span v-else class="price">¥{{ row.unitPrice?.toFixed(2) }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="数量" width="100">
        <template #default="{ row }">
          <ElInputNumber
            v-if="row.editing"
            v-model="row.quantity"
            :min="0"
            :precision="2"
            :controls="false"
            size="small"
            style="width: 100%"
          />
          <span v-else>{{ row.quantity }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="小计(元)" width="120">
        <template #default="{ row }">
          <span class="subtotal">¥{{ ((row.unitPrice || 0) * (row.quantity || 0)).toFixed(2) }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="备注" width="150">
        <template #default="{ row }">
          <ElInput
            v-if="row.editing"
            v-model="row.remark"
            placeholder="备注"
            size="small"
          />
          <span v-else>{{ row.remark || '-' }}</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="操作" width="140" fixed="right">
        <template #default="{ row, $index }">
          <template v-if="row.editing">
            <ElButton type="primary" link size="small" @click="handleSaveRow(row, $index)">保存</ElButton>
            <ElButton link size="small" @click="handleCancelEdit(row, $index)">取消</ElButton>
          </template>
          <template v-else>
            <ElButton type="primary" link size="small" @click="handleEditRow(row)">编辑</ElButton>
            <ElButton type="danger" link size="small" @click="handleDeleteRow($index)">删除</ElButton>
          </template>
        </template>
      </ElTableColumn>
    </ElTable>

    <!-- 汇总信息 / Summary -->
    <div class="summary-bar" v-if="materials.length > 0">
      <span>共 <strong>{{ materials.length }}</strong> 项材料</span>
      <span class="total">材料费合计：<strong>¥{{ materialCostTotal.toFixed(2) }}</strong></span>
    </div>

    <!-- 从材料库导入对话框 / Import from library dialog -->
    <ElDialog v-model="importDialogVisible" title="从材料库导入" width="800px" destroy-on-close>
      <div class="import-search">
        <ElInput
          v-model="importKeyword"
          placeholder="搜索材料名称"
          clearable
          style="width: 300px"
          @keyup.enter="searchLibraryMaterials"
        >
          <template #append>
            <ElButton @click="searchLibraryMaterials">搜索</ElButton>
          </template>
        </ElInput>
      </div>
      <ElTable
        :data="libraryMaterials"
        v-loading="importLoading"
        @selection-change="handleImportSelect"
        max-height="400px"
      >
        <ElTableColumn type="selection" width="50" />
        <ElTableColumn prop="name" label="材料名称" />
        <ElTableColumn prop="category" label="分类" width="100" />
        <ElTableColumn prop="specification" label="规格" width="120" />
        <ElTableColumn prop="unit" label="单位" width="80" />
        <ElTableColumn label="单价" width="100">
          <template #default="{ row }">¥{{ row.unitPrice?.toFixed(2) }}</template>
        </ElTableColumn>
      </ElTable>
      <template #footer>
        <ElButton @click="importDialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="confirmImport" :disabled="selectedImportIds.length === 0">
          导入选中 ({{ selectedImportIds.length }})
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
/***
 * Project Material List Component
 * 项目材料清单组件 - 支持输入框编辑
 ***/
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, FolderOpened, Download } from '@element-plus/icons-vue'
import {
  getProjectMaterials,
  saveProjectMaterials,
  exportProjectMaterials,
  type ProjectMaterialItem
} from '@/api/designer-project-material'
import { getMaterialList, type MaterialResponse } from '@/api/designer-material'

interface MaterialRow extends Partial<ProjectMaterialItem> {
  editing?: boolean
  _backup?: Partial<ProjectMaterialItem>
}

const props = defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  (e: 'update:materialCost', cost: number): void
}>()

const loading = ref(false)
const materials = ref<MaterialRow[]>([])

// 材料费合计 / Material cost total
const materialCostTotal = computed(() => {
  return materials.value.reduce((sum, item) => {
    return sum + (item.unitPrice || 0) * (item.quantity || 0)
  }, 0)
})

// 监听材料费变化，通知父组件 / Watch material cost change
watch(materialCostTotal, (val) => {
  emit('update:materialCost', val)
}, { immediate: true })

// 加载项目材料清单 / Load project materials
const loadMaterials = async () => {
  if (!props.projectId) return
  loading.value = true
  try {
    const res: any = await getProjectMaterials(props.projectId)
    if (res.code === 200 && res.data) {
      materials.value = (res.data.items || []).map((item: ProjectMaterialItem) => ({
        ...item,
        editing: false
      }))
    }
  } catch (error) {
    console.error('加载材料清单失败:', error)
  } finally {
    loading.value = false
  }
}

// 添加新材料行 / Add new material row
const handleAddMaterial = () => {
  materials.value.push({
    name: '',
    category: '',
    specification: '',
    unit: '个',
    unitPrice: 0,
    quantity: 1,
    remark: '',
    editing: true
  })
}

// 编辑行 / Edit row
const handleEditRow = (row: MaterialRow) => {
  row._backup = { ...row }
  row.editing = true
}

// 保存行 / Save row
const handleSaveRow = async (row: MaterialRow, index: number) => {
  if (!row.name?.trim()) {
    ElMessage.warning('请输入材料名称')
    return
  }
  row.editing = false
  row._backup = undefined
  row.totalPrice = (row.unitPrice || 0) * (row.quantity || 0)
  
  // 自动保存到后端 / Auto save to backend
  await saveMaterialList()
}

// 取消编辑 / Cancel edit
const handleCancelEdit = (row: MaterialRow, index: number) => {
  if (row._backup) {
    Object.assign(row, row._backup)
    row._backup = undefined
  } else if (!row.id) {
    // 新增行取消则删除 / Remove new row on cancel
    materials.value.splice(index, 1)
  }
  row.editing = false
}

// 删除行 / Delete row
const handleDeleteRow = async (index: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该材料吗？', '提示', { type: 'warning' })
    materials.value.splice(index, 1)
    await saveMaterialList()
    ElMessage.success('删除成功')
  } catch {
    // 取消删除
  }
}

// 保存材料清单 / Save material list
const saveMaterialList = async () => {
  if (!props.projectId) return
  try {
    const items = materials.value
      .filter(m => m.name?.trim())
      .map(m => ({
        materialId: m.materialId,
        name: m.name!,
        category: m.category || '',
        specification: m.specification || '',
        unit: m.unit || '个',
        unitPrice: m.unitPrice || 0,
        quantity: m.quantity || 0,
        totalPrice: (m.unitPrice || 0) * (m.quantity || 0),
        brand: m.brand,
        supplier: m.supplier,
        remark: m.remark
      }))
    
    await saveProjectMaterials({
      projectId: props.projectId,
      items
    })
  } catch (error) {
    console.error('保存材料清单失败:', error)
  }
}

// 导出清单 / Export list
const handleExport = async () => {
  try {
    const res: any = await exportProjectMaterials(props.projectId)
    if (res.code === 200 && res.data?.fileUrl) {
      window.open(res.data.fileUrl, '_blank')
      ElMessage.success('导出成功')
    }
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 从材料库导入相关 / Import from library
const importDialogVisible = ref(false)
const importLoading = ref(false)
const importKeyword = ref('')
const libraryMaterials = ref<MaterialResponse[]>([])
const selectedImportIds = ref<number[]>([])
const selectedImportItems = ref<MaterialResponse[]>([])

const showImportDialog = () => {
  importDialogVisible.value = true
  importKeyword.value = ''
  selectedImportIds.value = []
  searchLibraryMaterials()
}

const searchLibraryMaterials = async () => {
  importLoading.value = true
  try {
    const res: any = await getMaterialList({
      current: 1,
      size: 50,
      keyword: importKeyword.value || undefined,
      status: 'active'
    })
    libraryMaterials.value = res.data?.records || []
  } catch (error) {
    console.error('搜索材料库失败:', error)
  } finally {
    importLoading.value = false
  }
}

const handleImportSelect = (selection: MaterialResponse[]) => {
  selectedImportIds.value = selection.map(item => item.id)
  selectedImportItems.value = selection
}

const confirmImport = () => {
  // 将选中的材料添加到清单 / Add selected materials to list
  selectedImportItems.value.forEach(item => {
    materials.value.push({
      materialId: item.id,
      name: item.name,
      category: item.category,
      specification: item.specification,
      unit: item.unit,
      unitPrice: item.unitPrice,
      quantity: 1,
      brand: item.brand,
      supplier: item.supplier,
      remark: '',
      editing: false
    })
  })
  importDialogVisible.value = false
  ElMessage.success(`已导入 ${selectedImportItems.value.length} 项材料`)
  saveMaterialList()
}

// 监听projectId变化 / Watch projectId change
watch(() => props.projectId, (val) => {
  if (val) loadMaterials()
}, { immediate: true })

onMounted(() => {
  if (props.projectId) loadMaterials()
})

// 暴露方法给父组件 / Expose methods to parent
defineExpose({
  loadMaterials,
  saveMaterialList,
  getMaterialCost: () => materialCostTotal.value
})
</script>

<style scoped lang="scss">
.project-material-list {
  .action-bar {
    margin-bottom: 16px;
    display: flex;
    gap: 12px;
  }

  .price, .subtotal {
    color: #f56c6c;
    font-weight: 500;
  }

  .summary-bar {
    margin-top: 16px;
    padding: 12px 16px;
    background: #f5f7fa;
    border-radius: 4px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .total {
      font-size: 16px;
      strong {
        color: #f56c6c;
      }
    }
  }

  .import-search {
    margin-bottom: 16px;
  }
}
</style>
