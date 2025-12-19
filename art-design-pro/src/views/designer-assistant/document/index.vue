<template>
  <div class="document-page">
    <el-row :gutter="20">
      <!-- 左侧：文档列表 -->
      <el-col :span="8">
        <el-card class="list-card">
          <template #header>
            <div class="card-header">
              <span>项目文档</span>
              <el-button type="primary" size="small" @click="showUploadDialog = true">
                上传文档
              </el-button>
            </div>
          </template>

          <div v-loading="loading" class="document-list">
            <div
              v-for="doc in documents"
              :key="doc.id"
              class="document-item"
              :class="{ active: selectedDoc?.id === doc.id }"
              @click="selectDocument(doc)"
            >
              <el-icon class="doc-icon"><Document /></el-icon>
              <div class="doc-info">
                <div class="doc-name">{{ doc.fileName }}</div>
                <div class="doc-meta">
                  <span>{{ formatFileSize(doc.fileSize) }}</span>
                  <el-tag :type="getStatusType(doc.analysisStatus)" size="small">
                    {{ getStatusText(doc.analysisStatus) }}
                  </el-tag>
                </div>
              </div>
              <el-dropdown @command="handleCommand($event, doc)">
                <el-icon class="more-icon"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="analyze" :disabled="doc.analysisStatus === 'processing'">
                      分析文档
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <el-empty v-if="!loading && documents.length === 0" description="暂无文档" />
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              small
              @current-change="loadDocuments"
            />
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：分析结果 -->
      <el-col :span="16">
        <el-row :gutter="20">
          <el-col :span="24">
            <DocumentAnalysis
              :document="selectedDoc"
              @analysis-complete="handleAnalysisComplete"
            />
          </el-col>
          <el-col v-if="analysisResult" :span="24" style="margin-top: 20px">
            <KeywordExtraction
              :analysis-result="analysisResult"
              @apply="handleApplyToProject"
            />
          </el-col>
        </el-row>
      </el-col>
    </el-row>

    <!-- 上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="上传文档" width="600px">
      <DocumentUpload
        :project-id="projectId"
        @upload-success="handleUploadSuccess"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
/*** Document Page - 文档管理页面 ***/
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, MoreFilled } from '@element-plus/icons-vue'
import DocumentUpload from './DocumentUpload.vue'
import DocumentAnalysis from './DocumentAnalysis.vue'
import KeywordExtraction from './KeywordExtraction.vue'
import {
  getDocumentList,
  deleteDocument,
  analyzeDocument,
  type DocumentResponse,
  type DocumentAnalysisResult
} from '@/api/designer-document'

const route = useRoute()
const projectId = ref(Number(route.params.projectId) || 0)

const loading = ref(false)
const documents = ref<DocumentResponse[]>([])
const selectedDoc = ref<DocumentResponse | null>(null)
const analysisResult = ref<DocumentAnalysisResult | null>(null)
const showUploadDialog = ref(false)

const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

/*** Format file size ***/
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/*** Get status type ***/
function getStatusType(status: string) {
  switch (status) {
    case 'completed': return 'success'
    case 'processing': return 'warning'
    case 'failed': return 'danger'
    default: return 'info'
  }
}

/*** Get status text ***/
function getStatusText(status: string) {
  switch (status) {
    case 'completed': return '已分析'
    case 'processing': return '分析中'
    case 'failed': return '失败'
    default: return '待分析'
  }
}

/*** Load documents ***/
async function loadDocuments() {
  if (!projectId.value) return
  
  loading.value = true
  try {
    const res = await getDocumentList({
      projectId: projectId.value,
      current: currentPage.value,
      size: pageSize.value
    })
    documents.value = res.data?.records || []
    total.value = res.data?.total || 0
  } catch (error: any) {
    ElMessage.error(error.message || '加载文档列表失败')
  } finally {
    loading.value = false
  }
}

/*** Select document ***/
function selectDocument(doc: DocumentResponse) {
  selectedDoc.value = doc
  if (doc.analysisStatus === 'completed') {
    analysisResult.value = {
      documentId: doc.id,
      projectName: doc.projectName,
      area: doc.extractedArea,
      budget: doc.extractedBudget,
      style: doc.extractedStyle,
      functionalZones: doc.functionalZones,
      keywords: doc.keywords,
      summary: doc.summary,
      missingFields: [],
      suggestions: []
    }
  } else {
    analysisResult.value = null
  }
}

/*** Handle dropdown command ***/
async function handleCommand(command: string, doc: DocumentResponse) {
  if (command === 'analyze') {
    try {
      const res = await analyzeDocument(doc.id)
      ElMessage.success('分析完成')
      loadDocuments()
      if (selectedDoc.value?.id === doc.id) {
        analysisResult.value = res
      }
    } catch (error: any) {
      ElMessage.error(error.message || '分析失败')
    }
  } else if (command === 'delete') {
    ElMessageBox.confirm('确定要删除该文档吗？', '提示', {
      type: 'warning'
    }).then(async () => {
      await deleteDocument(doc.id)
      ElMessage.success('删除成功')
      if (selectedDoc.value?.id === doc.id) {
        selectedDoc.value = null
        analysisResult.value = null
      }
      loadDocuments()
    }).catch(() => {})
  }
}

/*** Handle upload success ***/
function handleUploadSuccess() {
  showUploadDialog.value = false
  loadDocuments()
}

/*** Handle analysis complete ***/
function handleAnalysisComplete(result: DocumentAnalysisResult) {
  analysisResult.value = result
  loadDocuments()
}

/*** Handle apply to project ***/
function handleApplyToProject(_data: any) {
  ElMessage.success('已应用到项目')
  // TODO: Update project with extracted data
}

onMounted(() => {
  loadDocuments()
})
</script>

<style scoped lang="scss">
.document-page {
  padding: 20px;

  .list-card {
    height: calc(100vh - 140px);
    display: flex;
    flex-direction: column;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    :deep(.el-card__body) {
      flex: 1;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }
  }

  .document-list {
    flex: 1;
    overflow-y: auto;

    .document-item {
      display: flex;
      align-items: center;
      padding: 12px;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.2s;

      &:hover {
        background: var(--el-fill-color-light);
      }

      &.active {
        background: var(--el-color-primary-light-9);
      }

      .doc-icon {
        font-size: 32px;
        color: var(--el-color-primary);
        margin-right: 12px;
      }

      .doc-info {
        flex: 1;
        overflow: hidden;

        .doc-name {
          font-weight: 500;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .doc-meta {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-top: 4px;
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
      }

      .more-icon {
        padding: 4px;
        cursor: pointer;
        color: var(--el-text-color-secondary);

        &:hover {
          color: var(--el-color-primary);
        }
      }
    }
  }

  .pagination {
    padding-top: 12px;
    display: flex;
    justify-content: center;
  }
}
</style>
