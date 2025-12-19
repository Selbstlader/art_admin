<template>
  <div class="keyword-extraction">
    <el-card class="keyword-card">
      <template #header>
        <div class="card-header">
          <span>关键信息提取</span>
          <el-button v-if="canApply" type="primary" size="small" @click="handleApply">
            应用到项目
          </el-button>
        </div>
      </template>

      <div v-if="analysisResult" class="extraction-content">
        <!-- 项目基本信息 -->
        <div class="info-group">
          <h4>项目基本信息</h4>
          <el-form label-width="80px" size="small">
            <el-form-item label="项目名称">
              <el-input v-model="editableData.projectName" placeholder="未提取到" />
            </el-form-item>
            <el-form-item label="面积">
              <el-input-number v-model="editableData.area" :min="0" :precision="2" placeholder="0" />
              <span class="unit">㎡</span>
            </el-form-item>
            <el-form-item label="预算">
              <el-input-number v-model="editableData.budget" :min="0" :precision="0" placeholder="0" />
              <span class="unit">元</span>
            </el-form-item>
            <el-form-item label="设计风格">
              <el-input v-model="editableData.style" placeholder="未提取到" />
            </el-form-item>
          </el-form>
        </div>

        <!-- 功能分区 -->
        <div class="info-group">
          <h4>功能分区</h4>
          <div class="editable-tags">
            <el-tag
              v-for="(zone, index) in editableData.functionalZones"
              :key="index"
              closable
              @close="removeZone(index)"
              type="info"
            >
              {{ zone }}
            </el-tag>
            <el-input
              v-if="zoneInputVisible"
              ref="zoneInputRef"
              v-model="zoneInputValue"
              class="tag-input"
              size="small"
              @keyup.enter="addZone"
              @blur="addZone"
            />
            <el-button v-else class="add-tag-btn" size="small" @click="showZoneInput">
              + 添加分区
            </el-button>
          </div>
        </div>

        <!-- 关键字 -->
        <div class="info-group">
          <h4>关键字</h4>
          <div class="editable-tags">
            <el-tag
              v-for="(keyword, index) in editableData.keywords"
              :key="index"
              closable
              @close="removeKeyword(index)"
            >
              {{ keyword }}
            </el-tag>
            <el-input
              v-if="keywordInputVisible"
              ref="keywordInputRef"
              v-model="keywordInputValue"
              class="tag-input"
              size="small"
              @keyup.enter="addKeyword"
              @blur="addKeyword"
            />
            <el-button v-else class="add-tag-btn" size="small" @click="showKeywordInput">
              + 添加关键字
            </el-button>
          </div>
        </div>

        <!-- 缺失信息提示 -->
        <div v-if="analysisResult.missingFields?.length" class="info-group warning">
          <h4>缺失信息</h4>
          <ul class="missing-list">
            <li v-for="field in analysisResult.missingFields" :key="field">{{ field }}</li>
          </ul>
        </div>

        <!-- 补充建议 -->
        <div v-if="analysisResult.suggestions?.length" class="info-group">
          <h4>补充建议</h4>
          <ul class="suggestion-list">
            <li v-for="suggestion in analysisResult.suggestions" :key="suggestion">
              {{ suggestion }}
            </li>
          </ul>
        </div>
      </div>

      <el-empty v-else description="暂无分析结果" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
/*** Keyword Extraction Component - 关键字提取展示组件 ***/
import { ref, reactive, computed, nextTick, watch } from 'vue'
import type { DocumentAnalysisResult } from '@/api/designer-document'

const props = defineProps<{
  analysisResult: DocumentAnalysisResult | null
}>()

const emit = defineEmits<{
  (e: 'apply', data: any): void
}>()

/*** Editable data ***/
const editableData = reactive({
  projectName: '',
  area: 0,
  budget: 0,
  style: '',
  functionalZones: [] as string[],
  keywords: [] as string[]
})

/*** Zone input state ***/
const zoneInputVisible = ref(false)
const zoneInputValue = ref('')
const zoneInputRef = ref()

/*** Keyword input state ***/
const keywordInputVisible = ref(false)
const keywordInputValue = ref('')
const keywordInputRef = ref()

/*** Can apply to project ***/
const canApply = computed(() => {
  return editableData.projectName || editableData.area > 0 || editableData.budget > 0
})

/*** Watch analysis result changes ***/
watch(() => props.analysisResult, (newVal) => {
  if (newVal) {
    editableData.projectName = newVal.projectName || ''
    editableData.area = newVal.area || 0
    editableData.budget = newVal.budget || 0
    editableData.style = newVal.style || ''
    editableData.functionalZones = [...(newVal.functionalZones || [])]
    editableData.keywords = [...(newVal.keywords || [])]
  }
}, { immediate: true })

/*** Zone operations ***/
function showZoneInput() {
  zoneInputVisible.value = true
  nextTick(() => zoneInputRef.value?.focus())
}

function addZone() {
  if (zoneInputValue.value.trim()) {
    editableData.functionalZones.push(zoneInputValue.value.trim())
  }
  zoneInputVisible.value = false
  zoneInputValue.value = ''
}

function removeZone(index: number) {
  editableData.functionalZones.splice(index, 1)
}

/*** Keyword operations ***/
function showKeywordInput() {
  keywordInputVisible.value = true
  nextTick(() => keywordInputRef.value?.focus())
}

function addKeyword() {
  if (keywordInputValue.value.trim()) {
    editableData.keywords.push(keywordInputValue.value.trim())
  }
  keywordInputVisible.value = false
  keywordInputValue.value = ''
}

function removeKeyword(index: number) {
  editableData.keywords.splice(index, 1)
}

/*** Apply to project ***/
function handleApply() {
  emit('apply', { ...editableData })
}
</script>

<style scoped lang="scss">
.keyword-extraction {
  .keyword-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }

  .extraction-content {
    .info-group {
      margin-bottom: 24px;

      &:last-child {
        margin-bottom: 0;
      }

      &.warning h4 {
        color: var(--el-color-warning);
      }

      h4 {
        margin: 0 0 12px 0;
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }

      .unit {
        margin-left: 8px;
        color: var(--el-text-color-secondary);
      }

      .editable-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        align-items: center;

        .tag-input {
          width: 100px;
        }

        .add-tag-btn {
          border-style: dashed;
        }
      }

      .missing-list, .suggestion-list {
        margin: 0;
        padding-left: 20px;
        color: var(--el-text-color-regular);

        li {
          margin-bottom: 4px;
          &:last-child {
            margin-bottom: 0;
          }
        }
      }

      .missing-list {
        color: var(--el-color-warning);
      }
    }
  }
}
</style>
