<template>
  <div class="department-page art-full-height">
    <DepartmentSearch
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
            <ElButton @click="showDialog('add')" v-ripple>新增部门</ElButton>
            <ElButton @click="expandAll" v-ripple>展开全部</ElButton>
            <ElButton @click="collapseAll" v-ripple>收起全部</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 部门树形表格 -->
      <ArtTable
        ref="tableRef"
        :loading="loading"
        :data="tableData"
        :columns="columns"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="false"
        row-key="deptId"
        :stripe="false"
      />
    </ElCard>

    <!-- 部门编辑弹窗 -->
    <DepartmentEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :department-data="currentDepartmentData"
      :department-options="departmentOptions"
      @success="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
  import { ElTag, ElMessageBox, ElButton, ElMessage } from 'element-plus'
  import { h } from 'vue'
  import { fetchGetDepartmentList, fetchDeleteDepartment } from '@/api/system-manage'
  import DepartmentSearch from './modules/department-search.vue'
  import DepartmentEditDialog from './modules/department-edit-dialog.vue'

  defineOptions({ name: 'Department' })

  type DepartmentListItem = Api.SystemManage.DepartmentListItem

  // 搜索表单
  const searchForm = ref({
    deptName: undefined,
    deptCode: undefined,
    status: undefined
  })

  const showSearchBar = ref(false)
  const dialogVisible = ref(false)
  const dialogType = ref<'add' | 'edit'>('add')
  const currentDepartmentData = ref<DepartmentListItem | undefined>(undefined)
  const departmentOptions = ref<DepartmentListItem[]>([])
  const tableRef = ref()

  // 获取部门数据（树形结构）
  const getDepartmentData = async (params?: any) => {
    try {
      loading.value = true
      const response = await fetchGetDepartmentList(params || {})
      tableData.value = response || []
      // 构建部门选项（用于上级部门选择）
      departmentOptions.value = buildDepartmentOptions(response || [])
    } catch (error) {
      console.error('获取部门列表失败:', error)
      tableData.value = []
      departmentOptions.value = []
    } finally {
      loading.value = false
    }
  }

  // 构建部门选项（扁平化处理）
  const buildDepartmentOptions = (
    departments: DepartmentListItem[],
    level = 0
  ): DepartmentListItem[] => {
    const options: DepartmentListItem[] = []
    departments.forEach((dept) => {
      options.push({
        ...dept,
        deptName: '　'.repeat(level) + dept.deptName
      })
      if (dept.children && dept.children.length > 0) {
        options.push(...buildDepartmentOptions(dept.children, level + 1))
      }
    })
    return options
  }

  const tableData = ref<DepartmentListItem[]>([])
  const loading = ref(false)

  const columns = [
    {
      prop: 'deptId',
      label: '部门ID',
      width: 100
    },
    {
      prop: 'deptName',
      label: '部门名称',
      minWidth: 200,
      formatter: (row: DepartmentListItem) => {
        return h('div', { class: 'dept-name-cell' }, [
          h('span', { class: 'dept-icon' }, row.children && row.children.length > 0 ? '📁' : '📄'),
          h('span', { class: 'dept-name' }, row.deptName)
        ])
      }
    },
    {
      prop: 'deptCode',
      label: '部门编码',
      width: 150
    },
    {
      prop: 'parentId',
      label: '上级部门',
      width: 150,
      formatter: (row: DepartmentListItem) => {
        return row.parentName || '根部门'
      }
    },
    {
      prop: 'orderNum',
      label: '排序号',
      width: 100
    },
    {
      prop: 'status',
      label: '状态',
      width: 100,
      formatter: (row: DepartmentListItem) => {
        const statusConfig = {
          type: (row.status === 1 ? 'success' : 'danger') as 'success' | 'danger',
          text: row.status === 1 ? '正常' : '停用'
        }
        return h(ElTag, { type: statusConfig.type }, () => statusConfig.text)
      }
    },
    {
      prop: 'createTime',
      label: '创建时间',
      width: 180
    },
    {
      prop: 'operation',
      label: '操作',
      width: 240,
      fixed: 'right' as const,
      formatter: (row: DepartmentListItem) =>
        h('div', { class: 'table-operation' }, [
          h(
            ElButton,
            {
              type: 'primary',
              size: 'small',
              onClick: () => showAddChildDialog(row)
            },
            () => '新增下级'
          ),
          h(
            ElButton,
            {
              type: 'warning',
              size: 'small',
              onClick: () => showDialog('edit', row)
            },
            () => '编辑'
          ),
          h(
            ElButton,
            {
              type: 'danger',
              size: 'small',
              onClick: () => deleteDepartment(row)
            },
            () => '删除'
          )
        ])
    }
  ]

  const columnChecks = ref(
    columns.map((col) => ({ label: col.label, value: col.prop, checked: true }))
  )

  // 搜索处理
  const handleSearch = (params: Record<string, any>) => {
    Object.assign(searchForm.value, params)
    refreshData()
  }

  // 重置搜索参数
  const resetSearchParams = () => {
    searchForm.value = {
      deptName: undefined,
      deptCode: undefined,
      status: undefined
    }
    refreshData()
  }

  // 显示弹窗
  const showDialog = (type: 'add' | 'edit', row?: DepartmentListItem) => {
    dialogType.value = type
    currentDepartmentData.value = row
    dialogVisible.value = true
  }

  // 显示新增下级部门弹窗
  const showAddChildDialog = (row: DepartmentListItem) => {
    dialogType.value = 'add'
    currentDepartmentData.value = {
      parentId: row.deptId,
      parentName: row.deptName
    } as DepartmentListItem
    dialogVisible.value = true
  }

  // 删除部门
  const deleteDepartment = (row: DepartmentListItem) => {
    const hasChildren = row.children && row.children.length > 0
    const confirmMessage = hasChildren
      ? `确定删除部门"${row.deptName}"吗？此部门包含下级部门，删除后下级部门也将被删除！此操作不可恢复！`
      : `确定删除部门"${row.deptName}"吗？此操作不可恢复！`

    ElMessageBox.confirm(confirmMessage, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        try {
          await fetchDeleteDepartment(row.deptId)
          ElMessage.success('删除成功')
          refreshData()
        } catch (error: any) {
          console.error('删除失败:', error)
          ElMessage.error(error.message || '删除失败，请重试')
        }
      })
      .catch(() => {
        ElMessage.info('已取消删除')
      })
  }

  // 刷新数据
  const refreshData = () => {
    getDepartmentData(searchForm.value)
  }

  // 展开全部
  const expandAll = () => {
    if (tableRef.value) {
      tableRef.value.toggleAllSelection()
      // 这里需要调用表格的展开全部方法，具体根据ArtTable组件的API来实现
    }
  }

  // 收起全部
  const collapseAll = () => {
    if (tableRef.value) {
      // 这里需要调用表格的收起全部方法，具体根据ArtTable组件的API来实现
    }
  }

  // 初始化
  onMounted(() => {
    getDepartmentData()
  })
</script>

<style lang="scss" scoped>
  .department-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .dept-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .dept-icon {
    font-size: 14px;
  }

  .dept-name {
    font-weight: 500;
  }
</style>
