<template>
  <div class="workbench-todo-page art-full-height">
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
            <ElTag type="warning" size="large">
              <el-icon><Bell /></el-icon>
              待办任务: {{ pagination.total || 0 }}
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
  import { Bell } from '@element-plus/icons-vue'

  defineOptions({ name: 'WorkbenchTodo' })

  const router = useRouter()

  // 搜索栏显示状态
  const showSearchBar = ref(true)

  // 搜索表单
  const searchForm = ref({
    processCode: undefined
  })

  // 搜索表单配置
  const formItems = computed(() => [
    {
      label: '流程类型',
      key: 'processCode',
      type: 'input',
      placeholder: '请输入流程编码',
      clearable: true
    }
  ])

  // 状态标签配置
  const statusConfig: Record<
    string,
    { type: 'warning' | 'success' | 'info' | 'danger'; text: string }
  > = {
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
        const res = await workflowQueryApi.getMyTodo({
          page: params.current,
          pageSize: params.size,
          processCode: params.processCode
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
          label: '当前节点',
          width: 120
        },
        {
          prop: 'initiatorName',
          label: '发起人',
          width: 100
        },
        {
          prop: 'status',
          label: '状态',
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
          prop: 'dueAt',
          label: '截止时间',
          width: 180,
          formatter: (row: WfTaskSummaryResponse) => {
            if (!row.dueAt) return '-'
            const isOverdue = new Date(row.dueAt) < new Date()
            return h('span', { style: { color: isOverdue ? '#f56c6c' : 'inherit' } }, row.dueAt)
          }
        },
        {
          prop: 'operation',
          label: '操作',
          width: 120,
          fixed: 'right',
          formatter: (row: WfTaskSummaryResponse) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  type: 'primary',
                  size: 'small',
                  onClick: () => handleApprove(row)
                },
                () => '审批'
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

  // 审批
  const handleApprove = (row: WfTaskSummaryResponse) => {
    router.push(`/workflow/workbench/detail/${row.processInstId}?taskId=${row.id}`)
  }
</script>

<style lang="scss" scoped>
  .workbench-todo-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }
</style>
