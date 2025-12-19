<!-- APP用户管理 -->
<template>
  <div class="app-user-page art-full-height">
    <!-- 搜索栏 -->
    <AppUserSearch
      v-model="searchForm"
      @search="handleSearch"
      @reset="resetSearchParams"
    ></AppUserSearch>

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 -->
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElSpace wrap>
            <ElButton @click="showDialog('add')" v-ripple>新增APP用户</ElButton>
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

      <!-- APP用户弹窗 -->
      <AppUserDialog
        v-model:visible="dialogVisible"
        :type="dialogType"
        :user-data="currentUserData"
        @submit="handleDialogSubmit"
      />

      <!-- 配额管理弹窗 -->
      <ElDialog v-model="quotaDialogVisible" title="效果图配额管理" width="450px">
        <div v-if="currentQuotaUser" class="quota-dialog-content">
          <p class="user-info">用户：{{ currentQuotaUser.userName }}</p>
          <ElForm :model="quotaForm" label-width="120px" v-loading="quotaLoading">
            <ElFormItem label="每日生成限制">
              <ElInputNumber
                v-model="quotaForm.dailyLimit"
                :min="0"
                :max="100"
                placeholder="每日可生成次数"
              />
              <span class="form-tip">次/天</span>
            </ElFormItem>
          </ElForm>
          <div class="quota-actions">
            <ElButton type="warning" @click="resetTodayUsage">重置今日使用次数</ElButton>
          </div>
        </div>
        <template #footer>
          <ElButton @click="quotaDialogVisible = false">取消</ElButton>
          <ElButton type="primary" @click="saveQuota">保存</ElButton>
        </template>
      </ElDialog>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/composables/useTable'
  import {
    fetchGetAppUserList,
    fetchDeleteAppUser,
    fetchResetAppUserPassword,
    fetchGetUserRenderQuota,
    fetchUpdateUserRenderQuota,
    fetchResetUserRenderQuota,
    type UserRenderQuotaInfo
  } from '@/api/system-manage'
  import AppUserSearch from './modules/app-user-search.vue'
  import AppUserDialog from './modules/app-user-dialog.vue'
  import {
    ElTag,
    ElMessageBox,
    ElMessage,
    ElDialog,
    ElForm,
    ElFormItem,
    ElInputNumber,
    ElButton
  } from 'element-plus'

  defineOptions({ name: 'AppUser' })

  type AppUserListItem = Api.SystemManage.AppUserListItem

  // 弹窗相关
  const dialogType = ref<Form.DialogType>('add')
  const dialogVisible = ref(false)
  const currentUserData = ref<Partial<AppUserListItem>>({})

  // 选中行
  const selectedRows = ref<AppUserListItem[]>([])

  // 配额管理相关 / Quota management related
  const quotaDialogVisible = ref(false)
  const quotaLoading = ref(false)
  const currentQuotaUser = ref<AppUserListItem | null>(null)
  const quotaForm = ref({
    dailyLimit: 3
  })

  // 搜索表单
  const searchForm = ref({
    userName: undefined,
    nickName: undefined,
    phone: undefined,
    userType: undefined,
    status: undefined
  })

  // 用户类型配置
  const USER_TYPE_CONFIG = {
    '1': { type: 'danger' as const, text: '管理员' },
    '2': { type: 'primary' as const, text: '普通用户' }
  } as const

  // 用户状态配置
  const USER_STATUS_CONFIG = {
    '1': { type: 'success' as const, text: '正常' },
    '2': { type: 'info' as const, text: '禁用' }
  } as const

  /**
   * 获取用户类型配置
   */
  const getUserTypeConfig = (userType: string) => {
    return (
      USER_TYPE_CONFIG[userType as keyof typeof USER_TYPE_CONFIG] || {
        type: 'info' as const,
        text: '未知'
      }
    )
  }

  /**
   * 获取用户状态配置
   */
  const getUserStatusConfig = (status: string) => {
    return (
      USER_STATUS_CONFIG[status as keyof typeof USER_STATUS_CONFIG] || {
        type: 'info' as const,
        text: '未知'
      }
    )
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
    // 核心配置
    core: {
      apiFn: fetchGetAppUserList,
      apiParams: {
        current: 1,
        size: 20,
        ...searchForm.value
      },
      columnsFactory: () => [
        { type: 'selection' }, // 勾选列
        { type: 'index', width: 60, label: '序号' }, // 序号
        {
          prop: 'userName',
          label: '用户账号',
          width: 180
        },
        {
          prop: 'nickName',
          label: '用户昵称',
          width: 150
        },
        { prop: 'phone', label: '手机号', width: 140 },
        { prop: 'email', label: '邮箱', width: 200 },
        {
          prop: 'userType',
          label: '用户类型',
          width: 120,
          formatter: (row) => {
            const typeConfig = getUserTypeConfig(row.userType)
            return h(ElTag, { type: typeConfig.type }, () => typeConfig.text)
          }
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) => {
            const statusConfig = getUserStatusConfig(row.status)
            return h(ElTag, { type: statusConfig.type }, () => statusConfig.text)
          }
        },
        {
          prop: 'lastLoginTime',
          label: '最后登录时间',
          width: 180,
          sortable: true
        },
        {
          prop: 'createTime',
          label: '创建时间',
          width: 180,
          sortable: true
        },
        {
          prop: 'operation',
          label: '操作',
          width: 280,
          fixed: 'right', // 固定列
          formatter: (row) =>
            h('div', { class: 'flex gap-1' }, [
              h(ArtButtonTable, {
                type: 'edit',
                onClick: () => showDialog('edit', row)
              }),
              h(ArtButtonTable, {
                type: 'more',
                text: '配额',
                onClick: () => showQuotaDialog(row)
              }),
              h(ArtButtonTable, {
                type: 'more',
                text: '重置密码',
                onClick: () => resetPassword(row)
              }),
              h(ArtButtonTable, {
                type: 'delete',
                onClick: () => deleteUser(row)
              })
            ])
        }
      ]
    }
  })

  /**
   * 搜索处理
   * @param params 参数
   */
  const handleSearch = (params: Record<string, any>) => {
    console.log(params)
    // 搜索参数赋值
    Object.assign(searchParams, params)
    getData()
  }

  /**
   * 显示用户弹窗
   */
  const showDialog = (type: Form.DialogType, row?: AppUserListItem): void => {
    console.log('打开弹窗:', { type, row })
    dialogType.value = type
    currentUserData.value = row || {}
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /**
   * 重置用户密码
   */
  const resetPassword = async (row: AppUserListItem): Promise<void> => {
    console.log('重置密码:', row)
    try {
      const { value: newPassword } = await ElMessageBox.prompt(
        `请输入用户"${row.userName}"的新密码`,
        '重置密码',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          inputType: 'password',
          inputValidator: (value) => {
            if (!value || value.length < 6) {
              return '密码长度不能少于6位'
            }
            return true
          }
        }
      )

      await fetchResetAppUserPassword({
        id: row.id,
        newPassword
      })
      ElMessage.success('密码重置成功')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('重置密码失败:', error)
        ElMessage.error('重置密码失败，请重试')
      }
    }
  }

  /**
   * 删除用户
   */
  const deleteUser = async (row: AppUserListItem): Promise<void> => {
    console.log('删除用户:', row)
    try {
      await ElMessageBox.confirm(`确定要删除用户"${row.userName}"吗？`, '删除用户', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      })

      await fetchDeleteAppUser(row.id)
      ElMessage.success('删除成功')
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除用户失败:', error)
        ElMessage.error('删除失败，请重试')
      }
    }
  }

  /**
   * 处理弹窗提交事件
   */
  const handleDialogSubmit = async () => {
    try {
      dialogVisible.value = false
      currentUserData.value = {}
      // 刷新数据列表
      refreshData()
    } catch (error) {
      console.error('提交失败:', error)
    }
  }

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: AppUserListItem[]): void => {
    selectedRows.value = selection
    console.log('选中行数据:', selectedRows.value)
  }

  /**
   * 显示配额管理弹窗
   */
  const showQuotaDialog = async (row: AppUserListItem): Promise<void> => {
    currentQuotaUser.value = row
    quotaLoading.value = true
    quotaDialogVisible.value = true

    try {
      const res = await fetchGetUserRenderQuota(row.id)
      if (res.data) {
        quotaForm.value.dailyLimit = res.data.dailyLimit || 3
      }
    } catch (error) {
      console.error('获取配额信息失败:', error)
      quotaForm.value.dailyLimit = 3
    } finally {
      quotaLoading.value = false
    }
  }

  /**
   * 保存配额设置
   */
  const saveQuota = async (): Promise<void> => {
    if (!currentQuotaUser.value) return

    try {
      await fetchUpdateUserRenderQuota({
        userId: currentQuotaUser.value.id,
        dailyLimit: quotaForm.value.dailyLimit
      })
      ElMessage.success('配额设置保存成功')
      quotaDialogVisible.value = false
    } catch (error) {
      console.error('保存配额失败:', error)
      ElMessage.error('保存失败，请重试')
    }
  }

  /**
   * 重置今日使用次数
   */
  const resetTodayUsage = async (): Promise<void> => {
    if (!currentQuotaUser.value) return

    try {
      await ElMessageBox.confirm(
        `确定要重置用户"${currentQuotaUser.value.userName}"的今日使用次数吗？`,
        '重置确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )

      await fetchResetUserRenderQuota(currentQuotaUser.value.id)
      ElMessage.success('今日使用次数已重置')
    } catch (error) {
      if (error !== 'cancel') {
        console.error('重置失败:', error)
        ElMessage.error('重置失败，请重试')
      }
    }
  }
</script>

<style lang="scss" scoped>
  .app-user-page {
    :deep(.flex) {
      display: flex;
      gap: 8px;
    }
  }

  .quota-dialog-content {
    .user-info {
      margin-bottom: 16px;
      font-size: 14px;
      color: var(--el-text-color-secondary);
    }

    .form-tip {
      margin-left: 8px;
      color: var(--el-text-color-secondary);
    }

    .quota-actions {
      margin-top: 16px;
      padding-top: 16px;
      border-top: 1px solid var(--el-border-color-light);
    }
  }
</style>
