<template>
  <div class="document-analysis">
    <el-card v-if="document" class="analysis-card">
      <template #header>
        <div class="card-header">
          <span class="title">{{ document.fileName }}</span>
          <el-tag :type="statusType" size="small">{{ statusText }}</el-tag>
        </div>
      </template>

      <!-- 分析中状态 -->
      <div v-if="document.analysisStatus === 'processing'" class="analyzing">
        <el-icon class="loading-icon"><Loading /></el-icon>
        <span>正在分析文档内容...</span>
      </div>

      <!-- 待分析状态 -->
      <div v-else-if="document.analysisStatus === 'pending'" class="pending">
        <el-button type="primary" @click="handleAnalyze" :loading="analyzing">
          开始分析
        </el-button>
        <p class="tip">点击按钮开始AI智能分析文档内容</p>
      </div>

      <!-- 分析完成 -->
      <div v-else-if="document.analysisStatus === 'completed'" class="completed">
        <!-- 提取的关键信息 -->
        <div class="info-section">
          <h4>提取的关键信息</h4>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="项目名称">
              {{ document.projectName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="面积">
              {{ document.extractedArea ? `${document.extractedArea} ㎡` : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="预算">
              {{ document.extractedBudget ? `¥${formatNumber(document.extractedBudget)}` : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="设计风格">
              {{ document.extractedStyle || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 功能分区 -->
        <div v-if="document.functionalZones?.length" class="info-section">
          <h4>功能分区</h4>
          <div class="tags">
            <el-tag v-for="zone in document.functionalZones" :key="zone" type="info" size="small">
              {{ zone }}
            </el-tag>
          </div>
        </div>

        <!-- 关键字 -->
        <div v-if="document.keywords?.length" class="info-section">
          <h4>关键字</h4>
          <div class="tags">
            <el-tag v-for="keyword in document.keywords" :key="keyword" size="small">
              {{ keyword }}
            </el-tag>
          </div>
        </div>

        <!-- 文档摘要 -->
        <div v-if="document.summary" class="info-section">
          <h4>文档摘要</h4>
          <div class="summary">{{ document.summary }}</div>
        </div>
      </div>

      <!-- 分析失败 -->
      <div v-else-if="document.analysisStatus === 'failed'" class="failed">
        <el-icon class="error-icon"><CircleClose /></el-icon>
        <span>分析失败，请重试</span>
        <el-button type="primary" size="small" @click="handleAnalyze" :loading="analyzing">
          重新分析
        </el-button>
      </div>
    </el-card>

    <!-- 空状态 -->
    <el-empty v-else description="请选择要查看的文档" />
  </div>
</template>

<script setup lang="ts">
/*** Document Analysis Component - 文档分析结果展示组件 ***/
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading, CircleClose } from '@element-plus/icons-vue'
import { analyzeDocument } from '@/api/designer-document'
import type { DocumentResponse } from '@/api/designer-document'

const props = defineProps<{
  document: DocumentResponse | null
}>()

const emit = defineEmits<{
  (e: 'analysis-complete', result: any): void
}>()

const analyzing = ref(false)

/*** Status type for tag ***/
const statusType = computed(() => {
  switch (props.document?.analysisStatus) {
    case 'completed': return 'success'
    case 'processing': return 'warning'
    case 'failed': return 'danger'
    default: return 'info'
  }
})

/*** Status text ***/
const statusText = computed(() => {
  switch (props.document?.analysisStatus) {
    case 'completed': return '已分析'
    case 'processing': return '分析中'
    case 'failed': return '分析失败'
    default: return '待分析'
  }
})

/*** Format number with commas ***/
function formatNumber(num: number): string {
  return num.toLocaleString('zh-CN')
}

/*** Handle analyze document ***/
async function handleAnalyze() {
  if (!props.document) return
  
  analyzing.value = true
  try {
    const res = await analyzeDocument(props.document.id)
    ElMessage.success('文档分析完成')
    emit('analysis-complete', res)
  } catch (error: any) {
    ElMessage.error(error.message || '分析失败')
  } finally {
    analyzing.value = false
  }
}
</script>

<style scoped lang="scss">
.document-analysis {
  .analysis-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-weight: 500;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        max-width: 300px;
      }
    }
  }

  .analyzing, .pending, .failed {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px 20px;
    color: var(--el-text-color-secondary);

    .loading-icon {
      font-size: 48px;
      animation: rotate 1.5s linear infinite;
      margin-bottom: 16px;
    }

    .error-icon {
      font-size: 48px;
      color: var(--el-color-danger);
      margin-bottom: 16px;
    }

    .tip {
      margin-top: 12px;
      font-size: 12px;
    }
  }

  .completed {
    .info-section {
      margin-bottom: 20px;

      &:last-child {
        margin-bottom: 0;
      }

      h4 {
        margin: 0 0 12px 0;
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }

      .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }

      .summary {
        padding: 12px;
        background: var(--el-fill-color-lighter);
        border-radius: 4px;
        line-height: 1.6;
        color: var(--el-text-color-regular);
        white-space: pre-wrap;
      }
    }
  }
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
