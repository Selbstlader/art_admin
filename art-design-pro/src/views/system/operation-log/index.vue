<template>
  <div class="operation-log-page art-full-height">
    <OperationLogSearch
      v-show="showSearchBar"
      v-model="searchForm"
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
            <ElButton @click="handleBatchDelete" :disabled="selectedRows.length === 0" v-ripple>
              批量删除
            </ElButton>
            <ElButton @click="handleClean" v-ripple>清理日志</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 表格 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
      </ArtTable>
    </ElCard>

    <!-- 日志详情弹窗 -->
    <OperationLogDetailDialog v-model="detailDialogVisible" :log-data="currentLogData" />
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/composables/useTable'
  import {
    fetchGetOperationLogList,
    fetchDeleteOperationLog,
    fetchBatchDeleteOperationLog,
    fetchCleanOperationLog
  } from '@/api/system-manage'
  import OperationLogSearch from './modules/operation-log-search.vue'
  import OperationLogDetailDialog from './modules/operation-log-detail-dialog.vue'
  import { ElTag, ElMessageBox, ElMessage, ElButton } from 'element-plus'
  import { h } from 'vue'

  defineOptions({ name: 'OperationLog' })

  type OperationLogItem = Api.SystemManage.OperationLogItem

  // 搜索表单
  const searchForm = ref({
    module: undefined,
    businessType: undefined,
    operatorName: undefined,
    status: undefined,
    startTime: undefined,
    endTime: undefined
  })

  const showSearchBar = ref(false)
  const detailDialogVisible = ref(false)
  const currentLogData = ref<OperationLogItem | undefined>(undefined)
  const selectedRows = ref<OperationLogItem[]>([])

  // 操作状态配置
  const STATUS_CONFIG = {
    1: { type: 'success' as const, text: '成功' },
    0: { type: 'danger' as const, text: '失败' }
  }

  // 业务类型颜色配置
  const BUSINESS_TYPE_CONFIG: Record<string, any> = {
    新增: { type: 'primary', text: '新增' },
    修改: { type: 'warning', text: '修改' },
    删除: { type: 'danger', text: '删除' },
    查询: { type: 'info', text: '查询' },
    导出: { type: 'success', text: '导出' },
    导入: { type: 'success', text: '导入' }
  }

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
    refreshData
  } = useTable({
    core: {
      apiFn: fetchGetOperationLogList,
      apiParams: {
        current: 1,
        size: 20
      },
      columnsFactory: () => [
        { type: 'selection' },
        { type: 'index', width: 60, label: '序号' },
        {
          prop: 'module',
          label: '操作模块',
          width: 120
        },
        {
          prop: 'businessType',
          label: '业务类型',
          width: 100,
          formatter: (row: OperationLogItem) => {
            const config = BUSINESS_TYPE_CONFIG[row.businessType] || {
              type: 'info',
              text: row.businessType
            }
            return h(ElTag, { type: config.type }, () => config.text)
          }
        },
        {
          prop: 'requestMethod',
          label: '请求方式',
          width: 100
        },
        {
          prop: 'operatorName',
          label: '操作人员',
          width: 120
        },
        {
          prop: 'operatorIp',
          label: '操作IP'
        },
        {
          prop: 'status',
          label: '操作状态',
          formatter: (row: OperationLogItem) => {
            const config = STATUS_CONFIG[row.status as keyof typeof STATUS_CONFIG]
            return h(ElTag, { type: config.type }, () => config.text)
          }
        },
        {
          prop: 'costTime',
          label: '耗时(ms)',
          formatter: (row: OperationLogItem) => {
            const time = row.costTime
            let type: 'success' | 'warning' | 'danger' = 'success'
            if (time > 1000) type = 'danger'
            else if (time > 500) type = 'warning'
            return h(ElTag, { type }, () => `${time}ms`)
          }
        },
        {
          prop: 'operationTime',
          label: '操作时间',
          width: 180,
          sortable: true
        }
        // {
        //   prop: 'operation',
        //   label: '操作',
        //   width: 150,
        //   fixed: 'right',
        //   formatter: (row: OperationLogItem) =>
        //     h('div', { class: 'flex gap-1' }, [
        //       h(
        //         ElButton,
        //         {
        //           type: 'primary',
        //           size: 'small',
        //           onClick: () => showDetail(row)
        //         },
        //         () => '详情'
        //       ),
        //       h(
        //         ElButton,
        //         {
        //           type: 'danger',
        //           size: 'small',
        //           onClick: () => deleteLog(row)
        //         },
        //         () => '删除'
        //       )
        //     ])
        // }
      ]
    }
  })

  /**
   * 搜索处理
   */
  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchParams, params)
    getData()
  }

  /**
   * 显示详情
   */
  const showDetail = (row: OperationLogItem) => {
    currentLogData.value = row
    detailDialogVisible.value = true
  }

  /**
   * 删除日志
   */
  const deleteLog = async (row: OperationLogItem) => {
    try {
      await ElMessageBox.confirm('确定要删除该操作日志吗？', '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await fetchDeleteOperationLog(row.id)
      ElMessage.success('删除成功')
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error('删除失败，请重试')
      }
    }
  }

  /**
   * 批量删除
   */
  const handleBatchDelete = async () => {
    if (selectedRows.value.length === 0) {
      ElMessage.warning('请选择要删除的日志')
      return
    }

    try {
      await ElMessageBox.confirm(
        `确定要删除选中的 ${selectedRows.value.length} 条日志吗？`,
        '批量删除确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )

      const ids = selectedRows.value.map((row) => row.id)
      await fetchBatchDeleteOperationLog(ids)
      ElMessage.success('批量删除成功')
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('批量删除失败:', error)
        ElMessage.error('批量删除失败，请重试')
      }
    }
  }

  /**
   * 清理日志
   */
  const handleClean = async () => {
    try {
      const { value: days } = await ElMessageBox.prompt('请输入要保留的天数', '清理日志', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'number',
        inputValue: '30',
        inputValidator: (value) => {
          const num = parseInt(value)
          if (!value || isNaN(num) || num <= 0) {
            return '请输入大于0的天数'
          }
          return true
        }
      })

      await fetchCleanOperationLog(parseInt(days))
      ElMessage.success('清理成功')
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('清理失败:', error)
        ElMessage.error('清理失败，请重试')
      }
    }
  }

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: OperationLogItem[]) => {
    selectedRows.value = selection
  }
</script>

<style lang="scss" scoped>
  .operation-log-page {
    padding-bottom: 15px;
  }
</style>
