<template>
  <div class="workbench-initiated-page art-full-height">
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
            <ElTag type="primary" size="large">
              <el-icon><Document /></el-icon>
              我发起的: {{ pagination.total || 0 }}
            </ElTag>
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
import { workflowQueryApi, processInstApi, type ProcessInstSummaryResponse } from '@/api/workflow'
import { ElMessage, ElMessageBox, ElTag, ElButton, ElDropdown, ElDropdownMenu, ElDropdownItem, ElIcon } from 'element-plus'
import { Document, MoreFilled } from '@element-plus/icons-vue'

defineOptions({ name: 'WorkbenchInitiated' })

const router = useRouter()

// 搜索栏显示状态
const showSearchBar = ref(true)

// 搜索表单
const searchForm = ref({
  title: undefined,
  status: undefined
})

// 搜索表单配置
const formItems = computed(() => [
  {
    label: '流程标题',
    key: 'title',
    type: 'input',
    placeholder: '请输入流程标题',
    clearable: true
  },
  {
    label: '状态',
    key: 'status',
    type: 'select',
    placeholder: '请选择状态',
    clearable: true,
    options: [
      { label: '进行中', value: 'running' },
      { label: '已完成', value: 'completed' },
      { label: '已拒绝', value: 'rejected' },
      { label: '已撤回', value: 'withdrawn' }
    ]
  }
])

// 状态标签配置
const statusConfig: Record<string, { type: 'warning' | 'success' | 'info' | 'danger'; text: string }> = {
  running: { type: 'warning', text: '进行中' },
  completed: { type: 'success', text: '已完成' },
  rejected: { type: 'danger', text: '已拒绝' },
  withdrawn: { type: 'info', text: '已撤回' }
}

// 表格配置
const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  getData,
  resetSearchParams,
  handleSizeChange,
  handleCurrentChange,
  refreshData
} = useTable({
  core: {
    apiFn: async (params: any) => {
      const res = await workflowQueryApi.getMyInitiated({
        page: params.current,
        pageSize: params.size,
        title: params.title,
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
        prop: 'title',
        label: '流程标题',
        minWidth: 200,
        showOverflowTooltip: true
      },
      {
        prop: 'processName',
        label: '流程类型',
        width: 150
      },
      {
        prop: 'status',
        label: '状态',
        width: 100,
        formatter: (row: ProcessInstSummaryResponse) => {
          const config = statusConfig[row.status] || { type: 'info', text: row.status }
          return h(ElTag, { type: config.type, size: 'small' }, () => config.text)
        }
      },
      {
        prop: 'startedAt',
        label: '发起时间',
        width: 180
      },
      {
        prop: 'completedAt',
        label: '完成时间',
        width: 180,
        formatter: (row: ProcessInstSummaryResponse) => row.completedAt || '-'
      },
      {
        prop: 'operation',
        label: '操作',
        width: 150,
        fixed: 'right',
        formatter: (row: ProcessInstSummaryResponse) =>
          h('div', { class: 'table-operation' }, [
            h(
              ElButton,
              {
                type: 'primary',
                size: 'small',
                onClick: () => handleView(row)
              },
              () => '查看'
            ),
            row.status === 'running'
              ? h(
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
                        h(ElDropdownItem, { command: 'withdraw' }, () => '撤回')
                      ])
                  }
                )
              : null
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

// 查看详情
const handleView = (row: ProcessInstSummaryResponse) => {
  router.push(`/workflow/workbench/detail/${row.id}`)
}

// 处理下拉菜单命令
const handleCommand = (command: string, row: ProcessInstSummaryResponse) => {
  switch (command) {
    case 'withdraw':
      handleWithdraw(row)
      break
  }
}

// 撤回流程
const handleWithdraw = async (row: ProcessInstSummaryResponse) => {
  try {
    await ElMessageBox.confirm(
      `确定要撤回流程"${row.title}"吗？撤回后流程将终止。`,
      '撤回确认',
      {
        type: 'warning',
        confirmButtonText: '确定撤回',
        cancelButtonText: '取消'
      }
    )

    await processInstApi.withdraw(row.id)
    ElMessage.success('撤回成功')
    refreshData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('撤回失败:', error)
      ElMessage.error(error.message || '撤回失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.workbench-initiated-page {
  padding-bottom: 15px;
}

.table-operation {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
