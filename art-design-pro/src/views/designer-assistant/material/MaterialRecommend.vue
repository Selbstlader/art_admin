<template>
  <div class="material-recommend">
    <!-- 推荐条件配置 / Recommendation config -->
    <el-card class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>智能推荐配置</span>
        </div>
      </template>
      <el-form :model="configForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="关联项目">
              <el-select
                v-model="configForm.projectId"
                placeholder="请选择项目"
                style="width: 100%"
                filterable
              >
                <el-option
                  v-for="project in projectList"
                  :key="project.id"
                  :label="project.name"
                  :value="project.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="空间类型">
              <el-select v-model="configForm.spaceType" placeholder="请选择空间类型" style="width: 100%">
                <el-option label="办公空间" value="办公空间" />
                <el-option label="商业空间" value="商业空间" />
                <el-option label="工业厂房" value="工业厂房" />
                <el-option label="酒店" value="酒店" />
                <el-option label="餐饮" value="餐饮" />
                <el-option label="医疗" value="医疗" />
                <el-option label="教育" value="教育" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="设计风格">
              <el-input v-model="configForm.style" placeholder="如：现代简约、中式等" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="预算(元)">
              <el-input-number
                v-model="configForm.budget"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="项目总预算"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="面积(m²)">
              <el-input-number
                v-model="configForm.area"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="项目面积"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="材料分类">
              <el-select v-model="configForm.category" placeholder="全部分类" clearable style="width: 100%">
                <el-option
                  v-for="cat in categories"
                  :key="cat.category"
                  :label="cat.category"
                  :value="cat.category"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleRecommend">
            <el-icon><MagicStick /></el-icon>
            获取智能推荐
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 推荐结果 / Recommendation results -->
    <el-card v-if="recommendResult" class="result-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>推荐结果 ({{ recommendResult.totalCount }}个)</span>
        </div>
      </template>
      
      <!-- 推荐理由 / Recommendation reason -->
      <el-alert
        v-if="recommendResult.recommendReason"
        :title="recommendResult.recommendReason"
        type="info"
        :closable="false"
        show-icon
        class="reason-alert"
      />

      <!-- 推荐材料列表 / Recommended materials list -->
      <div class="material-grid">
        <el-card
          v-for="material in recommendResult.materials"
          :key="material.id"
          class="material-card"
          shadow="hover"
        >
          <div class="material-image">
            <el-image
              v-if="material.imageUrl"
              :src="material.imageUrl"
              :preview-src-list="[material.imageUrl]"
              fit="cover"
            />
            <div v-else class="no-image">暂无图片</div>
          </div>
          <div class="material-info">
            <h4 class="material-name">{{ material.name }}</h4>
            <div class="material-meta">
              <el-tag size="small">{{ material.category }}</el-tag>
              <el-tag v-if="material.brand" size="small" type="info">{{ material.brand }}</el-tag>
            </div>
            <div class="material-spec">{{ material.specification || '暂无规格' }}</div>
            <div class="material-price">
              <span class="price">¥{{ material.unitPrice.toFixed(2) }}</span>
              <span class="unit">/{{ material.unit }}</span>
            </div>
            <div class="material-actions">
              <el-button type="primary" size="small" @click="handleSelect(material)">
                选择
              </el-button>
              <el-button size="small" @click="handleViewDetail(material)">
                详情
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 空状态 / Empty state -->
      <el-empty v-if="recommendResult.materials.length === 0" description="暂无推荐材料" />
    </el-card>

    <!-- 材料详情弹窗 / Material detail dialog -->
    <el-dialog v-model="detailVisible" title="材料详情" width="600px">
      <el-descriptions v-if="selectedMaterial" :column="2" border>
        <el-descriptions-item label="材料名称">{{ selectedMaterial.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ selectedMaterial.category }}</el-descriptions-item>
        <el-descriptions-item label="品牌">{{ selectedMaterial.brand || '-' }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedMaterial.specification || '-' }}</el-descriptions-item>
        <el-descriptions-item label="单价">¥{{ selectedMaterial.unitPrice.toFixed(2) }}/{{ selectedMaterial.unit }}</el-descriptions-item>
        <el-descriptions-item label="供应商">{{ selectedMaterial.supplier || '-' }}</el-descriptions-item>
        <el-descriptions-item label="适用场景" :span="2">
          <el-tag
            v-for="scene in selectedMaterial.applicableScenes"
            :key="scene"
            size="small"
            style="margin-right: 4px"
          >
            {{ scene }}
          </el-tag>
          <span v-if="!selectedMaterial.applicableScenes?.length">-</span>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ selectedMaterial.description || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { MagicStick } from '@element-plus/icons-vue'
import {
  recommendMaterials,
  getMaterialCategories,
  type MaterialResponse,
  type MaterialRecommendResponse,
  type MaterialCategoryResponse
} from '@/api/designer-material'
import { getDesignerProjects } from '@/api/designer-project'

/*** 配置表单 / Config form ***/
const configForm = reactive({
  projectId: undefined as number | undefined,
  spaceType: '',
  style: '',
  budget: undefined as number | undefined,
  area: undefined as number | undefined,
  category: '',
  limit: 12
})

/*** 数据 / Data ***/
const loading = ref(false)
const projectList = ref<{ id: number; name: string }[]>([])
const categories = ref<MaterialCategoryResponse[]>([])
const recommendResult = ref<MaterialRecommendResponse | null>(null)
const detailVisible = ref(false)
const selectedMaterial = ref<MaterialResponse | null>(null)

/*** Emits ***/
const emit = defineEmits<{
  (e: 'select', material: MaterialResponse): void
}>()

/*** 方法 / Methods ***/

// 加载项目列表 / Load project list
const loadProjectList = async () => {
  try {
    const res = await getDesignerProjects({ current: 1, size: 100 }) as unknown as Http.BaseResponse<{ records: { id: number; name: string }[] }>
    projectList.value = (res.data?.records || []).map((p) => ({ id: p.id, name: p.name }))
  } catch (error) {
    console.error('加载项目列表失败:', error)
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

// 获取推荐 / Get recommendations
const handleRecommend = async () => {
  if (!configForm.projectId) {
    ElMessage.warning('请选择关联项目')
    return
  }

  loading.value = true
  try {
    const res = await recommendMaterials({
      projectId: configForm.projectId,
      spaceType: configForm.spaceType || undefined,
      style: configForm.style || undefined,
      budget: configForm.budget,
      area: configForm.area,
      category: configForm.category || undefined,
      limit: configForm.limit
    }) as unknown as Http.BaseResponse<MaterialRecommendResponse>
    recommendResult.value = res.data || null
    if (res.data?.totalCount === 0) {
      ElMessage.info('暂无符合条件的推荐材料')
    }
  } catch (error) {
    console.error('获取推荐失败:', error)
    ElMessage.error('获取推荐失败')
  } finally {
    loading.value = false
  }
}

// 选择材料 / Select material
const handleSelect = (material: MaterialResponse) => {
  emit('select', material)
  ElMessage.success(`已选择材料: ${material.name}`)
}

// 查看详情 / View detail
const handleViewDetail = (material: MaterialResponse) => {
  selectedMaterial.value = material
  detailVisible.value = true
}

/*** 生命周期 / Lifecycle ***/
onMounted(() => {
  loadProjectList()
  loadCategories()
})
</script>

<style scoped lang="scss">
.material-recommend {
  padding: 16px;

  .config-card {
    margin-bottom: 16px;
  }

  .card-header {
    font-weight: 500;
  }

  .reason-alert {
    margin-bottom: 16px;
  }

  .material-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
  }

  .material-card {
    .material-image {
      height: 160px;
      margin: -20px -20px 12px -20px;
      overflow: hidden;
      border-radius: 4px 4px 0 0;

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
      }
    }

    .material-info {
      .material-name {
        margin: 0 0 8px 0;
        font-size: 16px;
        font-weight: 500;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .material-meta {
        margin-bottom: 8px;
        display: flex;
        gap: 4px;
      }

      .material-spec {
        margin-bottom: 8px;
        color: #909399;
        font-size: 13px;
      }

      .material-price {
        margin-bottom: 12px;

        .price {
          color: #f56c6c;
          font-size: 18px;
          font-weight: 600;
        }

        .unit {
          color: #909399;
          font-size: 13px;
        }
      }

      .material-actions {
        display: flex;
        gap: 8px;
      }
    }
  }
}
</style>
