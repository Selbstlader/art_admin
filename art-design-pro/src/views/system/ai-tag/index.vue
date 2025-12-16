<template>
  <div class="ai-tag-page art-full-height">
    <!-- Search Bar -->
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
            <ElButton type="primary" @click="showDialog('add')" v-ripple>
              <!-- <Icon icon="material-symbols:add" class="mr-1" /> -->
              新增标签
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- AI Tag Table -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>

    <!-- Create/Edit Dialog -->
    <AITagEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :tag-data="currentTagData"
      @success="refreshData"
    />

    <!-- Test Chat Dialog -->
    <AITagTestDialog v-model="testDialogVisible" :tag-data="testTagData" />
  </div>
</template>

<script setup lang="ts">
  /*** AI Tag Management Page ***/
  /*** Requirements: 2.1, 2.2, 2.3 - List, search, and display AI tags ***/

  import { h } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElTag, ElButton, ElSpace } from 'element-plus'
  import { useTable } from '@/composables/useTable'
  import { aiTagApi, type AITag } from '@/api/ai-tag'
  import AITagEditDialog from './modules/ai-tag-edit-dialog.vue'
  import AITagTestDialog from './modules/ai-tag-test-dialog.vue'

  defineOptions({ name: 'AITagManagement' })

  const router = useRouter()

  /*** State ***/
  const showSearchBar = ref(true)
  const dialogVisible = ref(false)
  const dialogType = ref<'add' | 'edit'>('add')
  const currentTagData = ref<AITag | undefined>(undefined)
  const testDialogVisible = ref(false)
  const testTagData = ref<AITag | undefined>(undefined)

  /*** Search Form Configuration ***/
  const searchForm = ref({
    keyword: ''
  })

  const formItems = computed(() => [
    {
      label: '关键词',
      key: 'keyword',
      type: 'input',
      placeholder: '搜索标签名称或描述',
      clearable: true
    }
  ])

  /*** Table Configuration using useTable composable ***/
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
        const response = await aiTagApi.list({
          current: params.current || 1,
          size: params.size || 20,
          keyword: params.keyword || ''
        })
        return {
          records: response.records || [],
          total: response.total || 0
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
          label: '标签名称',
          minWidth: 150
        },
        {
          prop: 'description',
          label: '描述',
          minWidth: 200,
          showOverflowTooltip: true
        },
        {
          prop: 'knowledge_base_name',
          label: '知识库',
          minWidth: 150,
          formatter: (row: AITag) => row.knowledge_base_name || '-'
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row: AITag) => {
            const statusConfig =
              row.status === 1
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
          prop: 'created_at',
          label: '创建时间',
          width: 180
        },
        {
          prop: 'operation',
          label: '操作',
          width: 300,
          fixed: 'right',
          formatter: (row: AITag) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  type: 'primary',
                  size: 'small',
                  disabled: row.status !== 1,
                  onClick: (e: Event) => {
                    e.stopPropagation()
                    handleChat(row)
                  }
                },
                () => '对话'
              ),
              h(
                ElButton,
                {
                  type: 'success',
                  size: 'small',
                  onClick: (e: Event) => {
                    e.stopPropagation()
                    handleTest(row)
                  }
                },
                () => '测试'
              ),
              h(
                ElButton,
                {
                  type: 'warning',
                  size: 'small',
                  onClick: (e: Event) => {
                    e.stopPropagation()
                    showDialog('edit', row)
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
                    handleDelete(row)
                  }
                },
                () => '删除'
              )
            ])
        }
      ]
    }
  })

  /*** Methods ***/
  const handleSearch = () => {
    getData({
      ...searchForm.value,
      current: 1
    })
  }

  const showDialog = (type: 'add' | 'edit', row?: AITag) => {
    dialogType.value = type
    currentTagData.value = type === 'edit' && row ? row : undefined
    dialogVisible.value = true
  }

  const handleDelete = async (row: AITag) => {
    try {
      await ElMessageBox.confirm(`确定要删除标签"${row.name}"吗？删除后将无法恢复。`, '删除确认', {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      })

      await aiTagApi.delete(row.id)
      ElMessage.success('删除成功')
      refreshData()
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
        ElMessage.error(error.message || '删除失败')
      }
    }
  }

  /*** Navigate to chat page ***/
  const handleChat = (row: AITag) => {
    router.push(`/system/ai-tag/chat/${row.id}`)
  }

  /*** Open test dialog ***/
  const handleTest = (row: AITag) => {
    testTagData.value = row
    testDialogVisible.value = true
  }

  /*** Lifecycle ***/
  onMounted(() => {
    getData({ current: 1, size: 20 })
  })
</script>

<style lang="scss" scoped>
  .ai-tag-page {
    padding-bottom: 15px;
  }

  .table-operation {
    display: flex;
    gap: 8px;
  }
</style>
