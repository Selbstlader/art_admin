<template>
  <div class="project-list-page art-full-height">
    <!-- 搜索栏 / Search Bar -->
    <ArtSearchBar
      v-model="searchForm"
      :items="searchItems"
      @search="handleSearch"
      @reset="handleReset"
    />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 / Table Header -->
      <ArtTableHeader :loading="loading" @refresh="getProjectList">
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="handleCreate" v-ripple>新建项目</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 项目列表 / Project List -->
      <ArtTable
        :loading="loading"
        :data="projectList"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handlePageChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  /***
   * Project List Component
   * 项目列表页面组件
   * Requirements: 4.2, 4.3
   ***/
  import { ref, reactive, computed, onMounted, h } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElTag, ElLink } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import dayjs from 'dayjs'
  import {
    getDesignerProjects,
    deleteDesignerProject,
    type ProjectResponse
  } from '@/api/designer-project'
  import type { ColumnOption } from '@/types/component'

  const router = useRouter()

  // 搜索表单 / Search form
  const searchForm = ref({
    name: '',
    status: '',
    style: ''
  })

  // 搜索配置 / Search config
  const searchItems = computed(() => [
    {
      label: '项目名称',
      key: 'name',
      type: 'input',
      placeholder: '请输入项目名称',
      clearable: true
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: {
        placeholder: '请选择状态',
        options: [
          { label: '草稿', value: 'draft' },
          { label: '进行中', value: 'in_progress' },
          { label: '已完成', value: 'completed' },
          { label: '已归档', value: 'archived' }
        ]
      }
    },
    {
      label: '设计风格',
      key: 'style',
      type: 'input',
      placeholder: '请输入设计风格',
      clearable: true
    }
  ])

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

  // 状态配置 / Status config
  const STATUS_CONFIG = {
    draft: { type: 'info' as const, text: '草稿' },
    in_progress: { type: 'warning' as const, text: '进行中' },
    completed: { type: 'success' as const, text: '已完成' },
    archived: { type: 'info' as const, text: '已归档' }
  } as const

  // 格式化金额 / Format money
  const formatMoney = (value: number) => {
    return value ? value.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' }) : '¥0.00'
  }

  // 格式化日期 / Format date
  const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
  }

  // 表格列配置 / Table columns config
  const columns = computed<ColumnOption[]>(() => [
    {
      prop: 'name',
      label: '项目名称',
      minWidth: 200,
      formatter: (row: ProjectResponse) =>
        h(ElLink, { type: 'primary', onClick: () => handleDetail(row.id) }, () => row.name)
    },
    {
      prop: 'area',
      label: '面积(m²)',
      width: 120
    },
    {
      prop: 'budget',
      label: '预算(元)',
      width: 150,
      formatter: (row: ProjectResponse) => formatMoney(row.budget)
    },
    {
      prop: 'style',
      label: '设计风格',
      width: 120
    },
    {
      prop: 'status',
      label: '状态',
      width: 100,
      formatter: (row: ProjectResponse) => {
        const config = STATUS_CONFIG[row.status as keyof typeof STATUS_CONFIG] || {
          type: 'info' as const,
          text: row.status
        }
        return h(ElTag, { type: config.type }, () => config.text)
      }
    },
    {
      prop: 'createdAt',
      label: '创建时间',
      width: 180,
      formatter: (row: ProjectResponse) => formatDate(row.createdAt)
    },
    {
      prop: 'operation',
      label: '操作',
      width: 180,
      fixed: 'right',
      formatter: (row: ProjectResponse) =>
        h('div', { class: 'flex gap-1' }, [
          h(ArtButtonTable, {
            type: 'view',
            onClick: () => handleDetail(row.id)
          }),
          h(ArtButtonTable, {
            type: 'edit',
            onClick: () => handleEdit(row)
          }),
          h(ArtButtonTable, {
            type: 'delete',
            onClick: () => handleDelete(row)
          })
        ])
    }
  ])

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
        name: searchForm.value.name || undefined,
        status: searchForm.value.status || undefined,
        style: searchForm.value.style || undefined
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
    searchForm.value = { name: '', status: '', style: '' }
    pagination.current = 1
    getProjectList()
  }

  // 新建项目 / Create project
  const handleCreate = () => {
    router.push('/designer/designer-assistant/project/ProjectCreate/0')
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
  const handleDelete = async (row: ProjectResponse) => {
    try {
      await ElMessageBox.confirm(`确定要删除项目"${row.name}"吗？`, '删除项目', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      const res = (await deleteDesignerProject(row.id)) as any
      if (res.code === 200) {
        ElMessage.success('删除成功')
        getProjectList()
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除项目失败:', error)
        ElMessage.error('删除失败')
      }
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

  onMounted(() => {
    getProjectList()
  })
</script>

<style scoped lang="scss">
  .project-list-page {
    // 样式由全局组件提供
  }
</style>
