<template>
  <div class="dify-knowledge-page art-full-height">
    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader
        v-model:showSearchBar="showSearchBar"
        :loading="loading"
        @refresh="refreshData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="handleUpload" v-ripple>
              <ElIcon><Upload /></ElIcon>
              上传文件
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 知识库列表 -->
      <div class="knowledge-list">
        <ElRow :gutter="16">
          <ElCol v-for="dataset in datasets" :key="dataset.id" :xs="24" :sm="12" :md="8" :lg="6">
            <ElCard class="dataset-card" shadow="hover" @click="handleViewDetail(dataset)">
              <div class="dataset-header">
                <ElIcon class="dataset-icon"><FolderOpened /></ElIcon>
                <h3 class="dataset-name">{{ dataset.name }}</h3>
              </div>
              <div class="dataset-info">
                <div class="info-item">
                  <span class="label">文档数:</span>
                  <span class="value">{{ dataset.document_count }}</span>
                </div>
                <div class="info-item">
                  <span class="label">字数:</span>
                  <span class="value">{{ formatNumber(dataset.word_count) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">更新时间:</span>
                  <span class="value">{{ formatTime(dataset.updated_at) }}</span>
                </div>
              </div>
              <div class="dataset-actions">
                <ElButton type="primary" size="small" @click.stop="handleUploadTo(dataset)">
                  上传文件
                </ElButton>
                <ElButton type="danger" size="small" @click.stop="handleDelete(dataset)">
                  删除
                </ElButton>
              </div>
            </ElCard>
          </ElCol>
        </ElRow>

        <ElEmpty v-if="datasets.length === 0 && !loading" description="暂无知识库" />
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="total > 0">
        <ElPagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </ElCard>

    <!-- 上传文件对话框 -->
    <ElDialog
      v-model="uploadVisible"
      title="上传文件到知识库"
      width="500px"
      :close-on-click-modal="false"
    >
      <ElForm label-width="100px">
        <ElFormItem label="选择知识库">
          <ElSelect v-model="selectedDatasetId" placeholder="请选择知识库" style="width: 100%">
            <ElOption
              v-for="dataset in datasets"
              :key="dataset.id"
              :label="dataset.name"
              :value="dataset.id"
            />
          </ElSelect>
        </ElFormItem>

        <ElFormItem label="选择文件">
          <ElUpload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            :accept="acceptFileTypes"
            :before-upload="beforeUpload"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
          >
            <template #trigger>
              <ElButton>
                <ElIcon><Upload /></ElIcon>
                选择文件
              </ElButton>
            </template>
            <template #tip>
              <div class="el-upload__tip">
                支持文档格式: txt, md, pdf, doc, docx, xls, xlsx, csv, html, json, xml<br />
                不支持图片、音频、视频等格式
              </div>
            </template>
          </ElUpload>
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElButton @click="uploadVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="uploading" @click="handleConfirmUpload">
          确定上传
        </ElButton>
      </template>
    </ElDialog>

    <!-- 知识库详情对话框 -->
    <ElDialog v-model="detailVisible" title="知识库详情" width="700px">
      <ElDescriptions v-if="currentDataset" :column="2" border>
        <ElDescriptionsItem label="名称">{{ currentDataset.name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="提供商">{{ currentDataset.provider }}</ElDescriptionsItem>
        <ElDescriptionsItem label="文档数">
          {{ currentDataset.document_count }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="字数">
          {{ formatNumber(currentDataset.word_count) }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="索引技术">
          {{ currentDataset.indexing_technique }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="嵌入模型">
          {{ currentDataset.embedding_model }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="创建时间" :span="2">
          {{ formatDateTime(currentDataset.created_at) }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="更新时间" :span="2">
          {{ formatDateTime(currentDataset.updated_at) }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="描述" :span="2">
          {{ currentDataset.description || '无' }}
        </ElDescriptionsItem>
      </ElDescriptions>

      <template #footer>
        <ElButton @click="detailVisible = false">关闭</ElButton>
        <ElButton type="primary" @click="handleUploadTo(currentDataset!)"> 上传文件 </ElButton>
        <ElButton type="danger" @click="handleDelete(currentDataset!)"> 删除知识库 </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Upload, FolderOpened } from '@element-plus/icons-vue'
  import { difyDatasetApi } from '@/api/dify'
  import type { DifyDataset } from '@/types/dify'
  import dayjs from 'dayjs'

  const showSearchBar = ref(false)
  const loading = ref(false)
  const datasets = ref<DifyDataset[]>([])
  const total = ref(0)
  const pagination = ref({
    page: 1,
    limit: 20
  })

  const uploadVisible = ref(false)
  const selectedDatasetId = ref('')
  const uploadRef = ref()
  const uploading = ref(false)
  const currentFile = ref<File | null>(null)

  const detailVisible = ref(false)
  const currentDataset = ref<DifyDataset | null>(null)

  // 允许的文件格式
  const acceptFileTypes = '.txt,.md,.pdf,.doc,.docx,.xls,.xlsx,.csv,.html,.htm,.json,.xml'

  // 获取知识库列表
  const getDatasets = async () => {
    loading.value = true
    try {
      const res = (await difyDatasetApi.getDatasetList({
        page: pagination.value.page,
        limit: pagination.value.limit
      })) as any
      if (res.data) {
        datasets.value = res.data
        total.value = res.total
      }
    } catch {
      ElMessage.error('获取知识库列表失败')
    } finally {
      loading.value = false
    }
  }

  // 刷新数据
  const refreshData = () => {
    getDatasets()
  }

  // 分页变化
  const handleSizeChange = (size: number) => {
    pagination.value.limit = size
    getDatasets()
  }

  const handleCurrentChange = (page: number) => {
    pagination.value.page = page
    getDatasets()
  }

  // 上传文件
  const handleUpload = () => {
    if (datasets.value.length === 0) {
      ElMessage.warning('请先创建知识库')
      return
    }
    selectedDatasetId.value = ''
    currentFile.value = null
    uploadVisible.value = true
  }

  const handleUploadTo = (dataset: DifyDataset) => {
    selectedDatasetId.value = dataset.id
    currentFile.value = null
    uploadVisible.value = true
  }

  const handleFileChange = (file: any) => {
    currentFile.value = file.raw
  }

  // 上传前验证文件格式
  const beforeUpload = (file: File) => {
    const fileName = file.name.toLowerCase()
    const allowedExtensions = [
      '.txt',
      '.md',
      '.pdf',
      '.doc',
      '.docx',
      '.xls',
      '.xlsx',
      '.csv',
      '.html',
      '.htm',
      '.json',
      '.xml'
    ]

    const isAllowed = allowedExtensions.some((ext) => fileName.endsWith(ext))

    if (!isAllowed) {
      ElMessage.error('不支持该文件格式!请上传文档类型文件(txt, md, pdf, docx, xlsx, csv等)')
      return false
    }

    // 限制文件大小为 50MB
    const maxSize = 50 * 1024 * 1024
    if (file.size > maxSize) {
      ElMessage.error('文件大小不能超过 50MB!')
      return false
    }

    return true
  }

  const handleExceed = () => {
    ElMessage.warning('只能上传一个文件')
  }

  const handleConfirmUpload = async () => {
    if (!selectedDatasetId.value) {
      ElMessage.warning('请选择知识库')
      return
    }
    if (!currentFile.value) {
      ElMessage.warning('请选择文件')
      return
    }

    uploading.value = true
    try {
      const res = await difyDatasetApi.uploadFile(selectedDatasetId.value, currentFile.value)
      if (res.code === 200) {
        ElMessage.success('上传成功')
        uploadVisible.value = false
        uploadRef.value?.clearFiles()
        getDatasets()
      }
    } catch {
      ElMessage.error('上传失败')
    } finally {
      uploading.value = false
    }
  }

  // 查看详情
  const handleViewDetail = async (dataset: DifyDataset) => {
    try {
      const res = await difyDatasetApi.getDatasetDetail(dataset.id)
      if (res.code === 200) {
        currentDataset.value = res.data
        detailVisible.value = true
      }
    } catch {
      ElMessage.error('获取详情失败')
    }
  }

  // 删除知识库
  const handleDelete = async (dataset: DifyDataset) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除知识库"${dataset.name}"吗?此操作不可恢复!`,
        '删除确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )

      const res = (await difyDatasetApi.deleteDataset(dataset.id)) as any
      if (res.code === 200) {
        ElMessage.success('删除成功')
        detailVisible.value = false
        getDatasets()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        const errorMsg = error?.response?.data?.msg || error?.message || '删除失败,请稍后重试'
        ElMessage.error(errorMsg)
      }
    }
  }

  // 格式化数字
  const formatNumber = (num: number) => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    }
    return num.toString()
  }

  // 格式化时间
  const formatTime = (timestamp: number) => {
    return dayjs.unix(timestamp).format('YYYY-MM-DD')
  }

  const formatDateTime = (timestamp: number) => {
    return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')
  }

  onMounted(() => {
    getDatasets()
  })
</script>

<style scoped lang="scss">
  .dify-knowledge-page {
    padding-bottom: 15px;
  }

  .knowledge-list {
    margin-top: 16px;
    min-height: 400px;
  }

  .dataset-card {
    margin-bottom: 16px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-4px);
    }

    .dataset-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;

      .dataset-icon {
        font-size: 32px;
        color: var(--el-color-primary);
      }

      .dataset-name {
        margin: 0;
        font-size: 16px;
        font-weight: 500;
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .dataset-info {
      margin-bottom: 16px;

      .info-item {
        display: flex;
        justify-content: space-between;
        padding: 8px 0;
        border-bottom: 1px solid var(--el-border-color-lighter);

        &:last-child {
          border-bottom: none;
        }

        .label {
          color: var(--el-text-color-secondary);
          font-size: 14px;
        }

        .value {
          font-weight: 500;
          font-size: 14px;
        }
      }
    }

    .dataset-actions {
      display: flex;
      justify-content: center;
    }
  }

  .pagination-wrapper {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
</style>
