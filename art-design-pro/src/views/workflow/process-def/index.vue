<template>
  <div class="process-def-page art-full-height">
    <ArtSearchBar
      v-show="showSearchBar"
      v-model="searchForm"
      :items="formItems"
      @search="handleSearch"
      @reset="resetSearchParams"
    />

    <ElCard
      class="art-table-card"
      shadow="never"
      :style="{ 'margin-top': showSearchBar ? '12px' : '0' }"
    >
      <ArtTableHeader
        v-model:columns="columnChecks"
        v-model:showSearchBar="showSearchBar"
        :loading="loading"
        @refresh="refreshData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="handleCreate" v-ripple>
              <el-icon><Plus /></el-icon>
              新增流程
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { useRouter } from 'vue-router'
import { useTable } from '@/composables/useTable'
import { processDefApi, type ProcessDefResponse } from '@/api/workflow'
import { ElMessage, ElMessageBox, ElTag, ElButton, ElDropdown, ElDropdownMenu, ElDropdownItem, ElIcon } from 'element-plus'
import { Plus, MoreFilled } from '@element-plus/icons-vue'

defineOptions({ name: 'ProcessDefList' })

const router = useRouter()

// 搜索栏显示状态
const showSearchBar = ref(true)

// 搜索表单
const searchForm = ref({
  name: undefined,
  code: undefined,
  category: undefined,
  status: undefined
})

// 搜索表单配置
const formItems = computed(() => [
  {
    label: '流程名称',
    key: 'name',
    type: 'input',
    placeholder: '请输入流程名称',
    clearable: true
  },
  {
    label: '流程编码',
    key: 'code',
    type: 'input',
    placeholder: '请输入流程编码',
    clearable: true
  },
  {
    label: '分类',
    key: 'category',
    type: 'input',
    placeholder: '请输入分类',
    clearable: true
  },
  {
    label: '状态',
    key: 'status',
    type: 'select',
    placeholder: '请选择状态',
    clearable: true,
    options: [
      { label: '草稿', value: 'draft' },
      { label: '已发布', value: 'published' },
      { label: '已禁用', value: 'disabled' }
    ]
  }
])

// 状态标签配置
const statusConfig: Record<string, { type: 'info' | 'success' | 'danger'; text: string }> = {
  draft: { type: 'info', text: '草稿' },
  published: { type: 'success', text: '已发布' },
  disabled: { type: 'danger', text: '已禁用' }
}

// 表格配置
const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  getData,
  searchParams,
  resetSearchParams,
  handleSizeChange,
  handleCurrentChange,
  refreshData,
  refreshRemove
} = useTable({
  core: {
    apiFn: async (params: any) => {
      const res = await processDefApi.getList({
        page: params.current,
        pageSize: params.size,
        name: params.name,
        code: params.code,
        category: params.category,
        status: params.status
      })
      return {
        records: res.list || [],
        total: res.total || 0,
        current: params.current,
        size: params.size
      }
    },
    columnsFactory: () => [
      {
        prop: 'id',
        label: 'ID',
        width: 80
      },
      {
        prop: 'name',
        label: '流程名称',
        minWidth: 150
      },
      {
        prop: 'code',
        label: '流程编码',
        width: 150
      },
      {
        prop: 'category',
        label: '分类',
        width: 120
      },
      {
        prop: 'version',
        label: '版本',
        width: 80,
        formatter: (row: ProcessDefResponse) => `v${row.version}`
      },
      {
        prop: 'status',
        label: '状态',
        width: 100,
        formatter: (row: ProcessDefResponse) => {
          const config = statusConfig[row.status] || { type: 'info', text: row.status }
          return h(ElTag, { type: config.type }, () => config.text)
        }
      },
      {
        prop: 'description',
        label: '描述',
        minWidth: 200,
        showOverflowTooltip: true
      },
      {
        prop: 'createdByName',
        label: '创建人',
        width: 120
      },
      {
        prop: 'createdAt',
        label: '创建时间',
        width: 180
      },
      {
        prop: 'operation',
        label: '操作',
        width: 200,
        fixed: 'right',
        formatter: (row: ProcessDefResponse) =>
          h('div', { class: 'table-operation' }, [
            h(
              ElButton,
              {
                type: 'primary',
                size: 'small',
                onClick: () => handleEdit(row)
              },
              () => '编辑'
            ),
            row.status === 'draft'
              ? h(
                  ElButton,
                  {
                    type: 'success',
                    size: 'small',
                    onClick: () => handlePublish(row)
                  },
                  () => '发布'
                )
              : null,
            h(
              ElDropdown,
              {
                trigger: 'click',
                onCommand: (command: string) => handleCommand(command, row)
              },
              {
                default: () =>
                  h(
                    ElButton,
                    { size: 'small', style: { marginLeft: '8px' } },
                    () => [h(ElIcon, null, () => h(MoreFilled)), '更多']
                  ),
                dropdown: () =>
                  h(ElDropdownMenu, null, () => [
                    h(ElDropdownItem, { command: 'copy' }, () => '复制'),
                    row.status === 'draft'
                      ? h(ElDropdownItem, { command: 'delete', divided: true }, () => '删除')
                      : null
                  ])
              }
            )
          ])
      }
    ]
  }
})

// 搜索
const handleSearch = () => {
  getData({
    ...searchForm.value,
    current: 1
  })
}

// 新增
const handleCreate = () => {
  router.push('/workflow/process-def/edit')
}

// 编辑
const handleEdit = (row: ProcessDefResponse) => {
  router.push(`/workflow/process-def/edit/${row.id}`)
}

// 发布
const handlePublish = async (row: ProcessDefResponse) => {
  try {
    await ElMessageBox.confirm(
      `确定要发布流程"${row.name}"吗？发布后将生成新版本。`,
      '发布确认',
      {
        type: 'warning',
        confirmButtonText: '确定发布',
        cancelButtonText: '取消'
      }
    )

    await processDefApi.publish(row.id)
    ElMessage.success('发布成功')
    refreshData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('发布失败:', error)
      ElMessage.error(error.message || '发布失败')
    }
  }
}

// 处理下拉菜单命令
const handleCommand = (command: string, row: ProcessDefResponse) => {
  switch (command) {
    case 'copy':
      handleCopy(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

// 复制流程
const handleCopy = async (row: ProcessDefResponse) => {
  try {
    await ElMessageBox.confirm(
      `确定要复制流程"${row.name}"吗？`,
      '复制确认',
      {
        type: 'info',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    // 获取流程详情
    const detail = await processDefApi.getDetail(row.id)
    
    // 创建新流程
    await processDefApi.create({
      name: `${row.name}_副本`,
      code: `${row.code}_copy_${Date.now()}`,
      description: row.description,
      category: row.category,
      formTemplateId: detail.formTemplateId,
      graph: detail.graph || { nodes: [], edges: [] }
    })
    
    ElMessage.success('复制成功')
    refreshData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('复制失败:', error)
      ElMessage.error(error.message || '复制失败')
    }
  }
}

// 删除
const handleDelete = async (row: ProcessDefResponse) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除流程"${row.name}"吗？此操作不可恢复！`,
      '删除确认',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await processDefApi.delete(row.id)
    ElMessage.success('删除成功')
    refreshRemove()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error(error.message || '删除失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.process-def-page {
  padding-bottom: 15px;
}

.table-operation {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
