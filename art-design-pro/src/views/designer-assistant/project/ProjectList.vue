<template>
  <div class="project-list-page">
    <ElCard shadow="never">
      <!-- 页面头部 -->
      <template #header>
        <div class="card-header">
          <span class="title">项目管理</span>
          <ElButton type="primary" @click="handleCreate">
            <ElIcon><Plus /></ElIcon>
            新建项目
          </ElButton>
        </div>
      </template>

      <!-- 搜索区域 -->
      <div class="search-area">
        <ElForm :inline="true" :model="searchForm">
          <ElFormItem label="项目名称">
            <ElInput v-model="searchForm.name" placeholder="请输入项目名称" clearable />
          </ElFormItem>
          <ElFormItem label="状态">
            <ElSelect v-model="searchForm.status" placeholder="请选择状态" clearable>
              <ElOption label="草稿" value="draft" />
              <ElOption label="进行中" value="in_progress" />
              <ElOption label="已完成" value="completed" />
              <ElOption label="已归档" value="archived" />
            </ElSelect>
          </ElFormItem>
          <ElFormItem label="设计风格">
            <ElInput v-model="searchForm.style" placeholder="请输入设计风格" clearable />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="handleSearch">
              <ElIcon><Search /></ElIcon>
              搜索
            </ElButton>
            <ElButton @click="handleReset">
              <ElIcon><Refresh /></ElIcon>
              重置
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>

      <!-- 项目列表 -->
      <ElTable :data="projectList" v-loading="loading" stripe>
        <ElTableColumn prop="name" label="项目名称" min-width="200">
          <template #default="{ row }">
            <ElLink type="primary" @click="handleDetail(row.id)">{{ row.name }}</ElLink>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="area" label="面积(m²)" width="120" />
        <ElTableColumn prop="budget" label="预算(元)" width="150">
          <template #default="{ row }">
            {{ formatMoney(row.budget) }}
          </template>
        </ElTableColumn>
        <ElTableColumn prop="style" label="设计风格" width="120" />
        <ElTableColumn prop="status" label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" link @click="handleDetail(row.id)">详情</ElButton>
            <ElButton type="primary" link @click="handleEdit(row)">编辑</ElButton>
            <ElPopconfirm title="确定删除该项目吗？" @confirm="handleDelete(row.id)">
              <template #reference>
                <ElButton type="danger" link>删除</ElButton>
              </template>
            </ElPopconfirm>
          </template>
        </ElTableColumn>
      </ElTable>

      <!-- 分页 -->
      <div class="pagination-area">
        <ElPagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Project List Component
   * 项目列表页面组件
   * Requirements: 4.2, 4.3
   ***/
  import { ref, reactive, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { Plus, Search, Refresh } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import {
    getDesignerProjects,
    deleteDesignerProject,
    type ProjectResponse,
    type ProjectListResponse
  } from '@/api/designer-project'

  const router = useRouter()

  // 搜索表单 / Search form
  const searchForm = reactive({
    name: '',
    status: '',
    style: ''
  })

  // 分页 / Pagination
  const pagination = reactive({
    current: 1,
    size: 10,
    total: 0
  })

  // 加载状态 / Loading state
  const loading = ref(false)

  // 项目列表数据 / Project list data
  const projectList = ref<ProjectResponse[]>([])

  /***
   * Get project list from API
   * 从API获取项目列表
   ***/
  const getProjectList = async () => {
    loading.value = true
    try {
      const res = (await getDesignerProjects({
        current: pagination.current,
        size: pagination.size,
        name: searchForm.name || undefined,
        status: searchForm.status || undefined,
        style: searchForm.style || undefined
      })) as any
      if (res.code === 200 && res.data) {
        projectList.value = res.data.records || []
        pagination.total = res.data.total || 0
      } else {
        ElMessage.error(res.msg || '获取项目列表失败')
      }
    } catch (error) {
      console.error('获取项目列表失败:', error)
      ElMessage.error('获取项目列表失败')
    } finally {
      loading.value = false
    }
  }

  // 搜索 / Search
  const handleSearch = () => {
    pagination.current = 1
    getProjectList()
  }

  // 重置 / Reset
  const handleReset = () => {
    searchForm.name = ''
    searchForm.status = ''
    searchForm.style = ''
    pagination.current = 1
    getProjectList()
  }

  // 新建项目 / Create project
  const handleCreate = () => {
    router.push('/designer/designer-assistant/project/ProjectCreate')
  }

  // 查看详情 / View detail
  const handleDetail = (id: number) => {
    router.push(`/designer/designer-assistant/project/ProjectDetail/${id}`)
  }

  // 编辑项目 / Edit project
  const handleEdit = (row: ProjectResponse) => {
    router.push(`/designer/designer-assistant/project/ProjectCreate/${row.id}`)
  }

  /***
   * Delete project
   * 删除项目
   ***/
  const handleDelete = async (id: number) => {
    try {
      const res = (await deleteDesignerProject(id)) as any
      if (res.code === 200) {
        ElMessage.success('删除成功')
        getProjectList()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      console.error('删除项目失败:', error)
      ElMessage.error('删除失败')
    }
  }

  // 分页大小变化 / Page size change
  const handleSizeChange = (size: number) => {
    pagination.size = size
    getProjectList()
  }

  // 页码变化 / Page change
  const handlePageChange = (page: number) => {
    pagination.current = page
    getProjectList()
  }

  // 格式化金额 / Format money
  const formatMoney = (value: number) => {
    return value ? value.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' }) : '¥0.00'
  }

  // 格式化日期 / Format date
  const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  // 获取状态类型 / Get status type
  type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
  const getStatusType = (status: string): TagType => {
    const map: Record<string, TagType> = {
      draft: 'info',
      in_progress: 'warning',
      completed: 'success',
      archived: 'info'
    }
    return map[status] || 'info'
  }

  // 获取状态文本 / Get status text
  const getStatusText = (status: string) => {
    const map: Record<string, string> = {
      draft: '草稿',
      in_progress: '进行中',
      completed: '已完成',
      archived: '已归档'
    }
    return map[status] || status
  }

  onMounted(() => {
    getProjectList()
  })
</script>

<style scoped lang="scss">
  .project-list-page {
    padding: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .search-area {
      margin-bottom: 16px;
    }

    .pagination-area {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }
  }
</style>
