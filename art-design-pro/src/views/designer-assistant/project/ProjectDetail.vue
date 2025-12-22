<template>
  <div class="project-detail-page">
    <ElCard shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <div class="left">
            <ElButton link @click="handleBack">
              <ElIcon><ArrowLeft /></ElIcon>
              返回
            </ElButton>
            <span class="title">{{ project.name || '项目详情' }}</span>
            <ElTag :type="getStatusType(project.status)">{{ getStatusText(project.status) }}</ElTag>
          </div>
          <div class="right">
            <ElButton type="primary" @click="handleEdit">编辑</ElButton>
          </div>
        </div>
      </template>

      <!-- 基本信息 -->
      <ElDescriptions title="基本信息" :column="3" border>
        <ElDescriptionsItem label="项目名称">{{ project.name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="面积">{{ project.area }} m²</ElDescriptionsItem>
        <ElDescriptionsItem label="预算">{{ formatMoney(project.budget) }}</ElDescriptionsItem>
        <ElDescriptionsItem label="设计风格">{{ project.style || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="创建时间">{{
          formatDate(project.createdAt)
        }}</ElDescriptionsItem>
        <ElDescriptionsItem label="更新时间">{{
          formatDate(project.updatedAt)
        }}</ElDescriptionsItem>
        <ElDescriptionsItem label="项目描述" :span="3">
          {{ project.description || '暂无描述' }}
        </ElDescriptionsItem>
      </ElDescriptions>

      <!-- 设计流程模块 / Design Workflow Modules -->
      <div class="module-section">
        <h3>设计流程</h3>
        <div class="workflow-hint">
          <ElAlert type="info" :closable="false" show-icon>
            <template #title>
              推荐流程：文档分析 → AI生成CAD → CAD预览/编辑 → AI生成效果图 → 创建设计版本 →
              施工图标注
            </template>
          </ElAlert>
        </div>
        <ElRow :gutter="16">
          <!-- 第一行：核心设计流程 -->
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-1" @click="goToModule('document')">
              <div class="step-badge">1</div>
              <ElIcon :size="32"><Document /></ElIcon>
              <span>文档分析</span>
              <div class="desc">上传需求文档，AI提取关键信息</div>
              <div class="count" v-if="project.documentCount"
                >{{ project.documentCount }} 个文档</div
              >
            </ElCard>
          </ElCol>
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-2" @click="goToModule('cad-generation')">
              <div class="step-badge">2</div>
              <ElIcon :size="32"><Cpu /></ElIcon>
              <span>AI生成CAD</span>
              <div class="desc">基于文档AI生成CAD图纸</div>
              <div class="count" v-if="cadGenerationCount">{{ cadGenerationCount }} 个任务</div>
            </ElCard>
          </ElCol>
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-3" @click="goToModule('cad-viewer')">
              <div class="step-badge">3</div>
              <ElIcon :size="32"><Files /></ElIcon>
              <span>CAD预览/编辑</span>
              <div class="desc">预览和编辑CAD图纸</div>
              <div class="count" v-if="project.cadFileCount">{{ project.cadFileCount }} 个文件</div>
            </ElCard>
          </ElCol>
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-4" @click="goToModule('render')">
              <div class="step-badge">4</div>
              <ElIcon :size="32"><PictureFilled /></ElIcon>
              <span>AI生成效果图</span>
              <div class="desc">基于CAD生成效果图</div>
              <div class="count" v-if="renderCount">{{ renderCount }} 张效果图</div>
            </ElCard>
          </ElCol>
        </ElRow>
        <ElRow :gutter="16" style="margin-top: 16px">
          <!-- 第二行：版本管理和标注 -->
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-5" @click="goToModule('version')">
              <div class="step-badge">5</div>
              <ElIcon :size="32"><Collection /></ElIcon>
              <span>设计版本管理</span>
              <div class="desc">归档设计图和CAD文件，支持AI智能对比</div>
              <div class="count" v-if="versionCount">{{ versionCount }} 个版本</div>
            </ElCard>
          </ElCol>
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card step-6" @click="goToModule('annotation')">
              <div class="step-badge">6</div>
              <ElIcon :size="32"><EditPen /></ElIcon>
              <span>施工图标注</span>
              <div class="desc">对CAD图纸添加施工标注</div>
              <div class="count" v-if="annotationCount">{{ annotationCount }} 个标注</div>
            </ElCard>
          </ElCol>
          <ElCol :span="6">
            <ElCard shadow="hover" class="module-card" @click="goToModule('design-compare')">
              <ElIcon :size="32"><PictureFilled /></ElIcon>
              <span>设计比对</span>
              <div class="desc">对比设计图与实际效果</div>
              <div class="count" v-if="compareResults.length"
                >{{ compareResults.length }} 个结果</div
              >
            </ElCard>
          </ElCol>
        </ElRow>
      </div>

      <!-- 材料清单与成本估算 / Material List & Cost Estimate -->
      <div class="material-cost-section" v-if="project.id">
        <ElRow :gutter="16">
          <!-- 左侧：材料清单 / Left: Material list -->
          <ElCol :span="16">
            <ElCard shadow="never">
              <template #header>
                <div class="section-header">
                  <ElIcon><Goods /></ElIcon>
                  <span>项目材料清单</span>
                </div>
              </template>
              <ProjectMaterialList
                ref="materialListRef"
                :project-id="project.id"
                @update:material-cost="handleMaterialCostUpdate"
              />
            </ElCard>
          </ElCol>
          <!-- 右侧：成本汇总 / Right: Cost summary -->
          <ElCol :span="8">
            <ProjectCostSummary
              ref="costSummaryRef"
              :project-id="project.id"
              :material-cost="materialCost"
            />
          </ElCol>
        </ElRow>
      </div>

      <!-- 关联文档 -->
      <div class="documents-section">
        <h3>关联文档 ({{ documents.length }})</h3>
        <ElEmpty v-if="documents.length === 0" description="暂无关联文档" />
        <ElTable v-else :data="documents" stripe>
          <ElTableColumn prop="fileName" label="文件名" />
          <ElTableColumn prop="fileType" label="类型" width="100" />
          <ElTableColumn prop="fileSize" label="大小" width="120">
            <template #default="{ row }">
              {{ formatFileSize(row.fileSize) }}
            </template>
          </ElTableColumn>
          <ElTableColumn prop="analysisStatus" label="分析状态" width="120">
            <template #default="{ row }">
              <ElTag :type="getAnalysisStatusType(row.analysisStatus)">
                {{ getAnalysisStatusText(row.analysisStatus) }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="createdAt" label="上传时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </div>

      <!-- CAD文件 -->
      <div class="cad-section" v-if="cadFiles.length > 0">
        <h3>CAD文件 ({{ cadFiles.length }})</h3>
        <ElTable :data="cadFiles" stripe>
          <ElTableColumn prop="fileName" label="文件名" />
          <ElTableColumn prop="fileFormat" label="格式" width="100" />
          <ElTableColumn prop="layerCount" label="图层数" width="100" />
          <ElTableColumn prop="has3d" label="3D" width="80">
            <template #default="{ row }">
              <ElTag :type="row.has3d ? 'success' : 'info'">{{ row.has3d ? '是' : '否' }}</ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="parseStatus" label="解析状态" width="120">
            <template #default="{ row }">
              <ElTag :type="getAnalysisStatusType(row.parseStatus)">
                {{ getAnalysisStatusText(row.parseStatus) }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="createdAt" label="上传时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </div>

      <!-- 智能设计建议 / Design Suggestions -->
      <!-- Requirements: 3.2, 3.3, 3.4 -->
      <div class="suggestions-section" v-if="project.id">
        <h3>
          <ElIcon><MagicStick /></ElIcon>
          智能设计建议
        </h3>
        <SuggestionList :project-id="project.id" />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Project Detail Component
   * 项目详情页面组件
   * Requirements: 4.2, 6.2, 10.1
   ***/
  import { ref, reactive, onMounted } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import {
    ArrowLeft,
    Document,
    PictureFilled,
    Files,
    MagicStick,
    Goods,
    Cpu,
    EditPen,
    Collection
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import {
    getDesignerProject,
    type ProjectDetailResponse,
    type DocumentBrief,
    type CadFileBrief,
    type CostEstimateBrief,
    type CompareResultBrief
  } from '@/api/designer-project'
  import SuggestionList from '../suggestion/SuggestionList.vue'
  import ProjectMaterialList from './components/ProjectMaterialList.vue'
  import ProjectCostSummary from './components/ProjectCostSummary.vue'

  /*** API Response interface / API 响应接口 ***/
  interface ApiResponse<T = unknown> {
    code: number
    msg?: string
    data?: T
  }

  /*** Module route mapping / 模块路由映射 ***/
  const moduleRouteMap: Record<string, string> = {
    document: '/designer/designer-assistant/document/DocumentUpload',
    'cad-generation': '/designer/designer-assistant/cad-generation/CadGenerationList',
    'cad-viewer': '/designer/designer-assistant/cad-viewer/CadUpload',
    render: '/designer/designer-assistant/cad-viewer/CadUpload', // 效果图渲染入口在CAD预览页
    version: '/designer/designer-assistant/version-compare/VersionList',
    annotation: '/designer/designer-assistant/construction-annotation/AnnotationList',
    'design-compare': '/designer/designer-assistant/design-compare/CompareUpload'
  }

  const router = useRouter()
  const route = useRoute()

  const loading = ref(false)
  const project = reactive<Partial<ProjectDetailResponse>>({
    id: 0,
    name: '',
    description: '',
    area: 0,
    budget: 0,
    style: '',
    status: 'draft',
    createdAt: '',
    updatedAt: '',
    documentCount: 0,
    cadFileCount: 0
  })
  const documents = ref<DocumentBrief[]>([])
  const cadFiles = ref<CadFileBrief[]>([])
  const costEstimate = ref<CostEstimateBrief | null>(null)
  const compareResults = ref<CompareResultBrief[]>([])
  const versionCount = ref<number>(0)
  const cadGenerationCount = ref<number>(0)
  const renderCount = ref<number>(0)
  const annotationCount = ref<number>(0)

  // 材料清单和成本汇总组件引用 / Material list and cost summary refs
  const materialListRef = ref<InstanceType<typeof ProjectMaterialList>>()
  const costSummaryRef = ref<InstanceType<typeof ProjectCostSummary>>()
  const materialCost = ref(0)

  // 材料费更新回调 / Material cost update callback
  const handleMaterialCostUpdate = (cost: number) => {
    materialCost.value = cost
  }

  /***
   * Get project detail from API
   * 从API获取项目详情
   ***/
  const getProjectDetail = async () => {
    const id = route.params.id as string
    if (!id) return

    loading.value = true
    try {
      const res = (await getDesignerProject(
        Number(id)
      )) as unknown as ApiResponse<ProjectDetailResponse>
      if (res.code === 200 && res.data) {
        Object.assign(project, res.data)
        documents.value = res.data.documents || []
        cadFiles.value = res.data.cadFiles || []
        costEstimate.value = res.data.costEstimate || null
        compareResults.value = res.data.compareResults || []
        versionCount.value = res.data.versionCount || 0
      } else {
        ElMessage.error(res.msg || '获取项目详情失败')
      }
    } catch (error) {
      console.error('获取项目详情失败:', error)
      ElMessage.error('获取项目详情失败')
    } finally {
      loading.value = false
    }
  }

  // 返回列表 / Back to list
  const handleBack = () => {
    router.push('/designer/designer-assistant/project/ProjectList')
  }

  // 编辑项目 / Edit project
  // 路由格式：/designer/designer-assistant/project/ProjectCreate/:id
  const handleEdit = () => {
    router.push(`/designer/designer-assistant/project/ProjectCreate/${route.params.id}`)
  }

  // 跳转到功能模块 / Go to module
  const goToModule = (module: string) => {
    const path = moduleRouteMap[module]
    if (path) {
      router.push(`${path}?projectId=${route.params.id}`)
    }
  }

  // 格式化金额 / Format money
  const formatMoney = (value: number | undefined) => {
    return value ? value.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' }) : '¥0.00'
  }

  // 格式化日期 / Format date
  const formatDate = (date: string | undefined) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  // 格式化文件大小 / Format file size
  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  // 获取状态类型 / Get status type
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string | undefined): TagType => {
    const map: Record<string, TagType> = {
      draft: 'info',
      in_progress: 'warning',
      completed: 'success',
      archived: 'info'
    }
    return map[status || ''] || 'info'
  }

  // 获取状态文本 / Get status text
  const getStatusText = (status: string | undefined) => {
    const map: Record<string, string> = {
      draft: '草稿',
      in_progress: '进行中',
      completed: '已完成',
      archived: '已归档'
    }
    return map[status || ''] || status || ''
  }

  // 获取分析状态类型 / Get analysis status type
  const getAnalysisStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      pending: 'info',
      processing: 'warning',
      completed: 'success',
      failed: 'danger'
    }
    return map[status] || 'info'
  }

  // 获取分析状态文本 / Get analysis status text
  const getAnalysisStatusText = (status: string) => {
    const map: Record<string, string> = {
      pending: '待处理',
      processing: '处理中',
      completed: '已完成',
      failed: '失败'
    }
    return map[status] || status
  }

  onMounted(() => {
    getProjectDetail()
  })
</script>

<style scoped lang="scss">
  .project-detail-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .left {
        display: flex;
        align-items: center;
        gap: 16px;

        .title {
          font-size: 16px;
          font-weight: 500;
        }
      }
    }

    .module-section {
      margin-top: 24px;

      h3 {
        margin-bottom: 16px;
        font-size: 14px;
        font-weight: 500;
      }

      .workflow-hint {
        margin-bottom: 16px;
      }

      .module-card {
        cursor: pointer;
        text-align: center;
        padding: 24px 16px;
        transition: all 0.3s;
        position: relative;
        min-height: 140px;

        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        .step-badge {
          position: absolute;
          top: 8px;
          left: 8px;
          width: 20px;
          height: 20px;
          border-radius: 50%;
          background: var(--el-color-primary);
          color: #fff;
          font-size: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .el-icon {
          color: var(--el-color-primary);
          margin-bottom: 8px;
        }

        span {
          display: block;
          font-size: 14px;
          font-weight: 500;
        }

        .desc {
          margin-top: 4px;
          font-size: 12px;
          color: var(--el-text-color-secondary);
          line-height: 1.4;
        }

        .count {
          margin-top: 8px;
          font-size: 12px;
          color: var(--el-color-primary);
        }

        // 不同步骤的颜色区分
        &.step-1 .el-icon {
          color: #409eff;
        }
        &.step-2 .el-icon {
          color: #67c23a;
        }
        &.step-3 .el-icon {
          color: #e6a23c;
        }
        &.step-4 .el-icon {
          color: #f56c6c;
        }
        &.step-5 .el-icon {
          color: #909399;
        }
        &.step-6 .el-icon {
          color: #9c27b0;
        }
      }
    }

    .material-cost-section {
      margin-top: 24px;

      .section-header {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 500;

        .el-icon {
          color: var(--el-color-primary);
        }
      }
    }

    .documents-section,
    .cad-section,
    .suggestions-section {
      margin-top: 24px;

      h3 {
        margin-bottom: 16px;
        font-size: 14px;
        font-weight: 500;
        display: flex;
        align-items: center;
        gap: 8px;

        .el-icon {
          color: var(--el-color-primary);
        }
      }
    }
  }
</style>
