<template>
  <div class="dictionary-page art-full-height">
    <ArtSearchBar
      v-show="showSearchBar"
      v-model="searchForm"
      :items="formItems"
      @search="handleTypeSearch"
      @reset="resetTypeSearchParams"
    >
    </ArtSearchBar>

    <ElCard
      class="art-table-card"
      shadow="never"
      :style="{ 'margin-top': showSearchBar ? '12px' : '0' }"
    >
      <ArtTableHeader
        v-model:columns="typeColumnChecks"
        v-model:showSearchBar="showSearchBar"
        :loading="typeLoading"
        @refresh="refreshTypeData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton @click="showTypeDialog('add')" v-ripple>新增字典类型</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 字典类型表格 -->
      <ArtTable
        :loading="typeLoading"
        :data="typeData"
        :columns="typeColumns"
        :pagination="typePagination"
        @pagination:size-change="handleTypeSizeChange"
        @pagination:current-change="handleTypeCurrentChange"
      >
      </ArtTable>
    </ElCard>

    <!-- 字典类型编辑弹窗 -->
    <DictionaryTypeEditDialog
      v-model="typeDialogVisible"
      :dialog-type="typeDialogType"
      :dictionary-type-data="currentTypeData"
      @success="refreshTypeData"
    />

    <!-- 字典数据管理弹窗 -->
    <ElDialog
      v-model="dataDialogVisible"
      :title="`${currentTypeData?.typeName} - 字典数据管理`"
      width="80%"
      :close-on-click-modal="false"
      @close="handleDataDialogClose"
    >
      <div class="dictionary-data-content">
        <DictionarySearch
          v-show="showDataSearchBar"
          v-model="dataSearchForm"
          :dictionary-type-options="dictionaryTypeOptions"
          @search="handleDataSearch"
          @reset="resetDataSearchParams"
        >
        </DictionarySearch>

        <div :style="{ 'margin-top': showDataSearchBar ? '12px' : '0' }">
          <ArtTableHeader
            v-model:columns="dataColumnChecks"
            v-model:showSearchBar="showDataSearchBar"
            :loading="dataLoading"
            @refresh="refreshDataData"
          >
            <template #left>
              <ElSpace wrap>
                <ElButton @click="showDataItemDialog('add')" v-ripple>新增字典数据</ElButton>
              </ElSpace>
            </template>
          </ArtTableHeader>

          <!-- 字典数据表格 -->
          <ArtTable
            :loading="dataLoading"
            :data="dataData"
            :columns="dataColumns"
            :pagination="dataPagination"
            @pagination:size-change="handleDataSizeChange"
            @pagination:current-change="handleDataCurrentChange"
          >
          </ArtTable>
        </div>
      </div>
    </ElDialog>

    <!-- 字典数据编辑弹窗 -->
    <DictionaryEditDialog
      v-model="dataItemDialogVisible"
      :dialog-type="dataItemDialogType"
      :dictionary-data="currentDataData"
      :dictionary-type-options="dictionaryTypeOptions"
      @success="refreshDataData"
    />
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/composables/useTable'
  import {
    fetchGetDictionaryTypeList,
    fetchDeleteDictionaryType,
    fetchGetDictionaryList,
    fetchDeleteDictionary
  } from '@/api/system-manage'
  import DictionarySearch from './modules/dictionary-search.vue'
  // import DictionaryTypeEditDialog from './modules/dictionary-type-edit-dialog.vue'
  import DictionaryEditDialog from './modules/dictionary-edit-dialog.vue'
  import { ElMessage, ElMessageBox, ElTag, ElButton } from 'element-plus'
  import { h } from 'vue'

  defineOptions({ name: 'Dictionary' })

  type DictionaryTypeListItem = Api.SystemManage.DictionaryTypeItem
  type DictionaryListItem = Api.SystemManage.DictionaryItem

  // ========== 字典类型管理 ==========
  const showSearchBar = ref(false)
  const typeDialogVisible = ref(false)
  const typeDialogType = ref<'add' | 'edit'>('add')
  const currentTypeData = ref<DictionaryTypeListItem | undefined>(undefined)

  const searchForm = ref({
    typeName: undefined,
    typeCode: undefined,
    description: undefined,
    enabled: undefined
  })
  const formItems = computed(() => [
    {
      label: '字典类型',
      key: 'typeName',
      type: 'input',
      placeholder: '请输入字典类型名称',
      clearable: true
    },
    {
      label: '字典编码',
      key: 'typeCode',
      type: 'input',
      props: { placeholder: '请输入字典类型编码', maxlength: '11' }
    }
  ])

  const {
    columns: typeColumns,
    columnChecks: typeColumnChecks,
    data: typeData,
    loading: typeLoading,
    pagination: typePagination,
    getData: getTypeData,
    searchParams: typeSearchParams,
    resetSearchParams: resetTypeSearchParams,
    handleSizeChange: handleTypeSizeChange,
    handleCurrentChange: handleTypeCurrentChange,
    refreshData: refreshTypeData
  } = useTable({
    core: {
      apiFn: fetchGetDictionaryTypeList,
      columnsFactory: () => [
        {
          prop: 'id',
          label: 'ID',
          width: 80
        },
        {
          prop: 'typeName',
          label: '字典类型名称',
          minWidth: 150
        },
        {
          prop: 'typeCode',
          label: '字典类型编码',
          width: 150
        },
        {
          prop: 'description',
          label: '描述',
          minWidth: 200,
          showOverflowTooltip: true
        },
        {
          prop: 'enabled',
          label: '状态',
          width: 100,
          formatter: (row) => {
            const statusConfig = row.enabled
              ? { type: 'success', text: '启用' }
              : { type: 'danger', text: '禁用' }
            return h(
              ElTag,
              { type: statusConfig.type as 'success' | 'danger' },
              () => statusConfig.text
            )
          }
        },
        {
          prop: 'remark',
          label: '备注',
          minWidth: 150,
          showOverflowTooltip: true
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
          fixed: 'right',
          formatter: (row) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  type: 'primary',
                  size: 'small',
                  onClick: () => showDataDialog(row)
                },
                () => '字典数据'
              ),
              h(
                ElButton,
                {
                  type: 'warning',
                  size: 'small',
                  onClick: (e: Event) => {
                    e.stopPropagation()
                    showTypeDialog('edit', row)
                  }
                },
                () => '编辑'
              ),
              h(
                ElButton,
                {
                  type: 'danger',
                  size: 'small',
                  onClick: (e: Event) => {
                    e.stopPropagation()
                    handleTypeDelete(row)
                  }
                },
                () => '删除'
              )
            ])
        }
      ]
    }
  })

  // typeSearchParams 的类型注解
  typeSearchParams satisfies Api.SystemManage.DictionaryTypeSearchParams

  // ========== 字典数据管理 ==========
  const showDataSearchBar = ref(false)
  const dataDialogVisible = ref(false)
  const dataItemDialogVisible = ref(false)
  const dataItemDialogType = ref<'add' | 'edit'>('add')
  const currentDataData = ref<DictionaryListItem | undefined>(undefined)
  const dictionaryTypeOptions = ref<Api.SystemManage.DictionaryTypeItem[]>([])

  const dataSearchForm = ref({
    typeCode: undefined,
    label: undefined,
    value: undefined,
    enabled: undefined
  })

  const {
    columns: dataColumns,
    columnChecks: dataColumnChecks,
    data: dataData,
    loading: dataLoading,
    pagination: dataPagination,
    getData: getDataData,
    searchParams: dataSearchParams,
    resetSearchParams: resetDataSearchParams,
    handleSizeChange: handleDataSizeChange,
    handleCurrentChange: handleDataCurrentChange,
    refreshData: refreshDataData
  } = useTable({
    core: {
      apiFn: fetchGetDictionaryList,
      columnsFactory: () => [
        {
          prop: 'id',
          label: 'ID',
          width: 80
        },
        {
          prop: 'typeCode',
          label: '字典类型',
          width: 180,
          formatter: (row) => {
            return h('div', [
              h('div', { class: 'type-code' }, row.typeCode),
              row.dictionaryType && h('div', { class: 'type-name' }, row.dictionaryType.typeName)
            ])
          }
        },
        {
          prop: 'label',
          label: '字典标签',
          minWidth: 150
        },
        {
          prop: 'value',
          label: '字典值',
          width: 150
        },
        {
          prop: 'orderNum',
          label: '排序号',
          width: 100
        },
        {
          prop: 'enabled',
          label: '状态',
          width: 100,
          formatter: (row) => {
            const statusConfig = row.enabled
              ? { type: 'success', text: '启用' }
              : { type: 'danger', text: '禁用' }
            return h(
              ElTag,
              { type: statusConfig.type as 'success' | 'danger' },
              () => statusConfig.text
            )
          }
        },
        {
          prop: 'remark',
          label: '备注',
          minWidth: 150,
          showOverflowTooltip: true
        },
        {
          prop: 'createTime',
          label: '创建时间',
          width: 180
        },
        {
          prop: 'operation',
          label: '操作',
          width: 150,
          fixed: 'right',
          formatter: (row) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  type: 'primary',
                  size: 'small',
                  onClick: () => showDataItemDialog('edit', row)
                },
                () => '编辑'
              ),
              h(
                ElButton,
                {
                  type: 'danger',
                  size: 'small',
                  onClick: () => handleDataDelete(row)
                },
                () => '删除'
              )
            ])
        }
      ]
    }
  })

  // ========== 字典类型相关方法 ==========
  const handleTypeSearch = () => {
    getTypeData({
      ...(searchForm.value as unknown as Record<string, any>),
      current: 1
    })
  }

  const showTypeDialog = (type: 'add' | 'edit', row?: DictionaryTypeListItem) => {
    typeDialogType.value = type
    if (type === 'edit' && row) {
      currentTypeData.value = row
    } else {
      currentTypeData.value = undefined
    }
    typeDialogVisible.value = true
  }

  // 显示字典数据管理弹窗
  const showDataDialog = (row: DictionaryTypeListItem) => {
    currentTypeData.value = row
    dataSearchForm.value.typeCode = row.typeCode as any
    dataDialogVisible.value = true
    // 加载该字典类型的数据
    getDataData({
      ...dataSearchParams,
      typeCode: row.typeCode,
      current: 1
    })
  }

  // 字典数据弹窗关闭处理
  const handleDataDialogClose = () => {
    dataDialogVisible.value = false
    dataSearchForm.value.typeCode = undefined
    dataData.value = []
  }

  const handleTypeDelete = async (row: DictionaryTypeListItem) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除字典类型"${row.typeName}"吗？此操作不可恢复！`,
        '删除确认',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )

      await fetchDeleteDictionaryType(row.id)
      ElMessage.success('删除成功')
      refreshTypeData()
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  // ========== 字典数据相关方法 ==========
  const handleDataSearch = () => {
    getDataData({
      ...(dataSearchParams as unknown as Record<string, any>),
      current: 1
    })
  }

  const showDataItemDialog = (type: 'add' | 'edit', row?: DictionaryListItem) => {
    dataItemDialogType.value = type
    if (type === 'edit' && row) {
      currentDataData.value = row
    } else {
      currentDataData.value = undefined
      // 新增时，如果有选中的字典类型，设置默认值
      if (currentTypeData.value) {
        currentDataData.value = {
          typeCode: currentTypeData.value.typeCode,
          dictionaryType: currentTypeData.value
        } as DictionaryListItem
      }
    }
    dataItemDialogVisible.value = true
  }

  const handleDataDelete = async (row: DictionaryListItem) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除字典数据"${row.label}"吗？此操作不可恢复！`,
        '删除确认',
        {
          type: 'warning',
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }
      )

      await fetchDeleteDictionary(row.id)
      ElMessage.success('删除成功')
      refreshDataData()
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  // 获取字典类型选项（同时加载表格数据）
  const getDictionaryTypeOptions = async () => {
    try {
      // 使用同一个请求获取数据，同时设置给表格和选项
      const response = await fetchGetDictionaryTypeList({ current: 1, size: 100 })
      dictionaryTypeOptions.value = response.records || []

      // 如果表格数据为空，也加载表格数据
      if (typeData.value.length === 0) {
        await getTypeData({ current: 1, size: 20 })
      }
    } catch (error) {
      console.error('获取字典类型失败:', error)
    }
  }

  // 初始化
  onMounted(() => {
    getDictionaryTypeOptions()
  })
</script>

<style lang="scss" scoped>
  .dictionary-page {
    padding-bottom: 15px;
  }

  .dictionary-data-content {
    height: 60vh;
    display: flex;
    flex-direction: column;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }

  .type-code {
    font-weight: 500;
    color: var(--el-text-color-primary);
  }

  .type-name {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 2px;
  }
</style>
