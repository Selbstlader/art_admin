<template>
  <div class="workbench-done-page art-full-height">
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
            <ElTag type="success" size="large">
              <el-icon><Finished /></el-icon>
              已办任务: {{ pagination.total || 0 }}
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
import { workflowQueryApi, type WfTaskSummaryResponse } from '@/api/workflow'
import { ElTag, ElButton } from 'element-plus'
import { Finished } from '@element-plus/icons-vue'

defineOptions({ name: 'WorkbenchDone' })

const router = useRouter()

// 搜索栏显示状态
const showSearchBar = ref(true)

// 搜索表单
const searchForm = ref({
  processCode: undefined,
  status: undefined
})

// 搜索表单配置
const formItems = computed(() => [
  {
    label: '流程类型',
    key: 'processCode',
    type: 'input',
    placeholder: '请输入流程编码',
    clearable: true
  },
  {
    label: '处理结果',
    key: 'status',
    type: 'select',
    placeholder: '请选择处理结果',
    clearable: true,
    options: [
      { label: '已通过', value: 'approved' },
      { label: '已拒绝', value: 'rejected' },
      { label: '已委托', value: 'delegated' },
      { label: '已转办', value: 'transferred' }
    ]
  }
])

// 状态标签配置
const statusConfig: Record<string, { type: 'warning' | 'success' | 'info' | 'danger'; text: string }> = {
  pending: { type: 'warning', text: '待处理' },
  approved: { type: 'success', text: '已通过' },
  rejected: { type: 'danger', text: '已拒绝' },
  delegated: { type: 'info', text: '已委托' },
  transferred: { type: 'info', text: '已转办' }
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
      const res = await workflowQueryApi.getMyDone({
        page: params.current,
        pageSize: params.size,
        processCode: params.processCode,
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
        prop: 'processTitle',
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
        prop: 'nodeName',
        label: '处理节点',
        width: 120
      },
      {
        prop: 'initiatorName',
        label: '发起人',
        width: 100
      },
      {
        prop: 'status',
        label: '处理结果',
        width: 100,
        formatter: (row: WfTaskSummaryResponse) => {
          const config = statusConfig[row.status] || { type: 'info', text: row.status }
          return h(ElTag, { type: config.type, size: 'small' }, () => config.text)
        }
      },
      {
        prop: 'createdAt',
        label: '接收时间',
        width: 180
      },
      {
        prop: 'completedAt',
        label: '处理时间',
        width: 180
      },
      {
        prop: 'operation',
        label: '操作',
        width: 100,
        fixed: 'right',
        formatter: (row: WfTaskSummaryResponse) =>
          h('div', { class: 'table-operation' }, [
            h(
              ElButton,
              {
                type: 'primary',
                size: 'small',
                link: true,
                onClick: () => handleView(row)
              },
              () => '查看'
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

// 查看详情
const handleView = (row: WfTaskSummaryResponse) => {
  router.push(`/workflow/workbench/detail/${row.processInstId}`)
}
</script>

<style lang="scss" scoped>
.workbench-done-page {
  padding-bottom: 15px;
}

.table-operation {
  display: flex;
  gap: 8px;
}
</style>
