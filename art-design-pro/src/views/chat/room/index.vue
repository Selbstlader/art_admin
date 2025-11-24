<template>
  <div class="chat-room-page art-full-height">
    <ArtSearchBar
      v-show="showSearchBar"
      v-model="searchForm"
      :items="formItems"
      @search="getData"
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
          <ElButton type="primary" @click="handleCreate" v-ripple>
            <ElIcon><Plus /></ElIcon>
            创建聊天室
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #operation="{ row }">
          <ElButton link type="primary" @click="handleEnter(row)">进入</ElButton>
          <ElButton link type="primary" @click="handleEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="handleDelete(row)">删除</ElButton>
        </template>
      </ArtTable>
    </ElCard>

    <!-- 创建/编辑聊天室弹窗 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <ElForm ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <ElFormItem label="聊天室名称" prop="name">
          <ElInput v-model="form.name" placeholder="请输入聊天室名称" maxlength="100" />
        </ElFormItem>

        <ElFormItem label="聊天室描述" prop="description">
          <ElInput
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入聊天室描述"
            maxlength="500"
          />
        </ElFormItem>

        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="聊天室类型" prop="type">
              <ElSelect v-model="form.type" placeholder="请选择类型" style="width: 100%">
                <ElOption label="公开" value="public" />
                <ElOption label="私密" value="private" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="最大成员数" prop="maxMembers">
              <ElInputNumber
                v-model="form.maxMembers"
                :min="0"
                :max="10000"
                placeholder="0表示无限制"
                style="width: 100%"
              />
            </ElFormItem>
          </ElCol>
        </ElRow>

        <ElFormItem v-if="isEdit" label="状态" prop="isActive">
          <ElSwitch v-model="form.isActive" active-text="启用" inactive-text="禁用" />
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit" :loading="submitting">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, h } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElTag, type FormInstance, type FormRules } from 'element-plus'
  import { Plus } from '@element-plus/icons-vue'
  import { useTable } from '@/composables/useTable'
  import { chatRoomApi } from '@/api/chat'
  import { ElButton } from 'element-plus'
  defineOptions({ name: 'ChatRoomList' })

  const router = useRouter()
  const showSearchBar = ref(false)

  // 搜索表单
  const searchForm = ref({
    keyword: undefined,
    type: undefined,
    isActive: undefined
  })

  const formItems = computed(() => [
    {
      label: '关键词',
      key: 'keyword',
      type: 'input',
      placeholder: '聊天室名称',
      clearable: true
    },
    {
      label: '类型',
      key: 'type',
      type: 'select',
      placeholder: '全部',
      clearable: true,
      options: [
        { label: '公开', value: 'public' },
        { label: '私密', value: 'private' }
      ]
    },
    {
      label: '状态',
      key: 'isActive',
      type: 'select',
      placeholder: '全部',
      clearable: true,
      options: [
        { label: '启用', value: true },
        { label: '禁用', value: false }
      ]
    }
  ])

  // 使用 useTable composable
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
      apiFn: chatRoomApi.getRoomList,
      paginationKey: {
        current: 'page',
        size: 'pageSize'
      },
      columnsFactory: () => [
        {
          type: 'index',
          label: '序号',
          width: 60
        },
        {
          prop: 'name',
          label: '聊天室名称',
          minWidth: 150,
          showOverflowTooltip: true
        },
        {
          prop: 'description',
          label: '描述',
          minWidth: 200,
          showOverflowTooltip: true
        },
        {
          prop: 'type',
          label: '类型',
          width: 80,
          formatter: (row: any) => (row.type === 'public' ? '公开' : '私密')
        },
        {
          prop: 'memberCount',
          label: '成员数',
          width: 100,
          formatter: (row: any) => {
            const max = row.maxMembers > 0 ? `/${row.maxMembers}` : ''
            return `${row.memberCount}${max}`
          }
        },
        {
          prop: 'onlineCount',
          label: '在线人数',
          width: 90
        },
        {
          prop: 'isActive',
          label: '状态',
          width: 80,
          formatter: (row: any) =>
            h(ElTag, { type: row.isActive ? 'success' : 'danger' }, () =>
              row.isActive ? '启用' : '禁用'
            )
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 160
        },
        {
          label: '操作',
          width: 220,
          fixed: 'right',
          formatter: (row: any) =>
            h('div', { class: 'table-operation' }, [
              h(
                ElButton,
                {
                  link: true,
                  type: 'primary',
                  size: 'small',
                  onClick: () => handleEnter(row)
                },
                () => '进入'
              ),
              h(
                ElButton,
                {
                  link: true,
                  type: 'primary',
                  size: 'small',
                  onClick: () => handleEdit(row)
                },
                () => '编辑'
              ),
              h(
                ElButton,
                {
                  link: true,
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

  // 弹窗表单
  const dialogVisible = ref(false)
  const dialogTitle = computed(() => (isEdit.value ? '编辑聊天室' : '创建聊天室'))
  const isEdit = ref(false)
  const submitting = ref(false)
  const formRef = ref<FormInstance>()

  const form = reactive<Api.Chat.CreateChatRoomRequest & { id?: number; isActive?: boolean }>({
    name: '',
    description: '',
    type: 'public',
    maxMembers: 100
  })

  const formRules: FormRules = {
    name: [{ required: true, message: '请输入聊天室名称', trigger: 'blur' }],
    description: [{ required: true, message: '请输入聊天室描述', trigger: 'blur' }],
    type: [{ required: true, message: '请选择聊天室类型', trigger: 'change' }],
    maxMembers: [{ required: true, message: '请输入最大成员数', trigger: 'blur' }]
  }

  // 创建
  const handleCreate = () => {
    isEdit.value = false
    Object.assign(form, {
      name: '',
      description: '',
      type: 'public',
      maxMembers: 100
    })
    dialogVisible.value = true
  }

  // 编辑
  const handleEdit = (row: Api.Chat.ChatRoomItem) => {
    isEdit.value = true
    Object.assign(form, {
      id: row.id,
      name: row.name,
      description: row.description,
      maxMembers: row.maxMembers,
      isActive: row.isActive
    })
    dialogVisible.value = true
  }

  // 提交
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate()

    submitting.value = true
    try {
      if (isEdit.value) {
        await chatRoomApi.updateRoom(form as Api.Chat.UpdateChatRoomRequest)
        ElMessage.success('更新成功')
      } else {
        await chatRoomApi.createRoom(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      getData()
    } catch (error) {
      console.error('操作失败:', error)
    } finally {
      submitting.value = false
    }
  }

  // 删除
  const handleDelete = async (row: Api.Chat.ChatRoomItem) => {
    try {
      await ElMessageBox.confirm('确定要删除该聊天室吗？', '提示', {
        type: 'warning'
      })

      await chatRoomApi.deleteRoom(row.id)
      ElMessage.success('删除成功')
      getData()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除失败:', error)
      }
    }
  }

  // 进入聊天室
  const handleEnter = async (row: Api.Chat.ChatRoomItem) => {
    try {
      // 判断是否为视频通话房间
      const isVideoRoom = row.description === '视频通话房间'

      if (isVideoRoom) {
        // 视频房间直接跳转到视频通话页面，由视频通话页面处理加入逻辑
        router.push({
          path: '/chat/chat/video-call',
          query: { roomId: row.id.toString() }
        })
      } else {
        // 普通聊天室，先加入再跳转
        await chatRoomApi.joinRoom({ roomId: row.id })
        router.push(`/chat/chat/room-detail/${row.id}`)
      }
    } catch (error) {
      console.error('加入聊天室失败:', error)
      ElMessage.error('加入聊天室失败')
    }
  }
</script>

<style scoped lang="scss">
  .chat-room-page {
    display: flex;
    flex-direction: column;
  }
</style>
