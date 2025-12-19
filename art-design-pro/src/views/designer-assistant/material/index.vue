<template>
  <div class="material-module">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="材料库管理" name="list">
        <MaterialList />
      </el-tab-pane>
      <el-tab-pane label="智能推荐" name="recommend">
        <MaterialRecommend @select="handleMaterialSelect" />
      </el-tab-pane>
      <el-tab-pane label="材料选择" name="select">
        <MaterialSelect v-model="selectedMaterials" @change="handleSelectionChange" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import MaterialList from './MaterialList.vue'
import MaterialRecommend from './MaterialRecommend.vue'
import MaterialSelect from './MaterialSelect.vue'
import type { MaterialResponse } from '@/api/designer-material'

/*** 数据 / Data ***/
const activeTab = ref('list')
const selectedMaterials = ref<any[]>([])

/*** 方法 / Methods ***/

// 处理推荐材料选择 / Handle recommended material selection
const handleMaterialSelect = (material: MaterialResponse) => {
  // 切换到材料选择标签页 / Switch to material select tab
  activeTab.value = 'select'
  
  // 检查是否已选择 / Check if already selected
  const exists = selectedMaterials.value.some(m => m.id === material.id)
  if (!exists) {
    selectedMaterials.value.push({
      ...material,
      area: 0,
      lossRate: 0.05
    })
    ElMessage.success(`已添加材料: ${material.name}`)
  } else {
    ElMessage.info('该材料已在选择列表中')
  }
}

// 处理选择变化 / Handle selection change
const handleSelectionChange = (materials: any[]) => {
  console.log('Selected materials:', materials)
}
</script>

<style scoped lang="scss">
.material-module {
  height: 100%;
  
  :deep(.el-tabs) {
    height: 100%;
    display: flex;
    flex-direction: column;
    
    .el-tabs__content {
      flex: 1;
      overflow: auto;
      padding: 16px;
    }
  }
}
</style>
