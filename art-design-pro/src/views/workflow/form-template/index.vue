<template>
  <div class="form-template-page art-full-height">
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
              新增表单模板
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
  import { formTemplateApi, type FormTemplateResponse } from '@/api/workflow'
  import { ElMessage, ElMessageBox, ElTag, ElButton } from 'element-plus'
  import { Plus } from '@element-plus/icons-vue'

  defineOptions({ name: 'FormTemplateList' })

  const router = useRouter()

  // 搜索栏显示状态
  const showSearchBar = ref(true)

  // 搜索表单
  const searchForm = ref({
    name: undefined,
    code: undefined,
    status: undefined
  })

  // 搜索表单配置
  const formItems = computed(() => [
    {
      label: '模板名称',
      key: 'name',
      type: 'input',
      placeholder: '请输入模板名称',
      clearable: true
    },
    {
      label: '模板编码',
      key: 'code',
      type: 'input',
      placeholder: '请输入模板编码',
      clearable: true
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      placeholder: '请选择状态',
      clearable: true,
      options: [
        { label: '启用', value: 'active' },
        { label: '禁用', value: 'disabled' }
      ]
    }
  ])

  // 状态标签配置
  const statusConfig: Record<string, { type: 'success' | 'danger' | 'info'; text: string }> = {
    active: { type: 'success', text: '启用' },
    disabled: { type: 'danger', text: '禁用' }
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
        const res = await formTemplateApi.getList({
          page: params.current,
          pageSize: params.size,
          name: params.name,
          code: params.code,
          status: params.status
        })
        return {
          records: res.records || [],
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
          label: '模板名称',
          minWidth: 150
        },
        {
          prop: 'code',
          label: '模板编码',
          width: 150
        },
        {
          prop: 'description',
          label: '描述',
          minWidth: 200,
          showOverflowTooltip: true
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row: FormTemplateResponse) => {
            const config = statusConfig[row.status] || { type: 'info', text: row.status }
            return h(ElTag, { type: config.type }, () => config.text)
          }
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
          prop: 'updatedAt',
          label: '更新时间',
          width: 180
        },
        {
          prop: 'operation',
          label: '操作',
          width: 180,
          fixed: 'right',
          formatter: (row: FormTemplateResponse) =>
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
              h(
                ElButton,
                {
                  type: 'danger',
                  size: 'small',
                  onClick: () => handleDelete(row)
                },
                () => '删除'
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
    router.push('/workflow/form-template/edit')
  }

  // 编辑
  const handleEdit = (row: FormTemplateResponse) => {
    router.push(`/workflow/form-template/edit/${row.id}`)
  }

  // 删除
  const handleDelete = async (row: FormTemplateResponse) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除表单模板"${row.name}"吗？此操作不可恢复！`,
        '删除确认',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )

      await formTemplateApi.delete(row.id)
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
  .form-template-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }
</style>
