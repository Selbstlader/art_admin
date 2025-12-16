<template>
  <div class="workbench-detail-page">
    <ElCard v-loading="loading" shadow="never">
      <!-- 流程基本信息 -->
      <template #header>
        <div class="card-header">
          <div class="title-section">
            <h3>{{ detail?.title || '流程详情' }}</h3>
            <ElTag :type="statusConfig[detail?.status || '']?.type || 'info'" size="large">
              {{ statusConfig[detail?.status || '']?.text || detail?.status }}
            </ElTag>
          </div>
          <div class="action-section" v-if="canApprove">
            <ElButton type="success" @click="handleApprove">
              <el-icon><Check /></el-icon>
              通过
            </ElButton>
            <ElButton type="danger" @click="handleReject">
              <el-icon><Close /></el-icon>
              拒绝
            </ElButton>
            <ElButton @click="showDelegateDialog = true">
              <el-icon><Switch /></el-icon>
              委托
            </ElButton>
            <ElButton @click="showTransferDialog = true">
              <el-icon><Right /></el-icon>
              转办
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 流程信息 -->
      <ElDescriptions :column="3" border class="info-section">
        <ElDescriptionsItem label="流程编号">{{ detail?.id }}</ElDescriptionsItem>
        <ElDescriptionsItem label="发起人">{{ detail?.initiatorName }}</ElDescriptionsItem>
        <ElDescriptionsItem label="发起时间">{{ detail?.startedAt }}</ElDescriptionsItem>
        <ElDescriptionsItem label="当前节点">{{ currentNodeName || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="完成时间">{{ detail?.completedAt || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="流程版本">v{{ detail?.processDefVersion }}</ElDescriptionsItem>
      </ElDescriptions>

      <!-- 表单数据 -->
      <div class="form-section" v-if="detail?.formData && Object.keys(detail.formData).length > 0">
        <h4 class="section-title">
          <el-icon><Document /></el-icon>
          表单数据
        </h4>
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem
            v-for="(value, key) in detail.formData"
            :key="key"
            :label="String(key)"
          >
            {{ formatFormValue(value) }}
          </ElDescriptionsItem>
        </ElDescriptions>
      </div>

      <!-- 审批轨迹 -->
      <div class="trail-section">
        <h4 class="section-title">
          <el-icon><List /></el-icon>
          审批轨迹
        </h4>
        <ElTimeline v-if="detail?.approvalTrail && detail.approvalTrail.length > 0">
          <ElTimelineItem
            v-for="record in detail.approvalTrail"
            :key="record.taskId"
            :type="getTimelineType(record.action)"
            :timestamp="record.createdAt"
            placement="top"
          >
            <ElCard shadow="never" class="trail-card">
              <div class="trail-content">
                <div class="trail-header">
                  <span class="operator">{{ record.operatorName }}</span>
                  <ElTag :type="getActionTagType(record.action)" size="small">
                    {{ getActionText(record.action) }}
                  </ElTag>
                </div>
                <div class="trail-node" v-if="record.nodeName"> 节点: {{ record.nodeName }} </div>
                <div class="trail-comment" v-if="record.comment"> 意见: {{ record.comment }} </div>
              </div>
            </ElCard>
          </ElTimelineItem>
        </ElTimeline>
        <ElEmpty v-else description="暂无审批记录" />
      </div>

      <!-- 当前待办任务 -->
      <div class="tasks-section" v-if="detail?.currentTasks && detail.currentTasks.length > 0">
        <h4 class="section-title">
          <el-icon><Bell /></el-icon>
          当前待办
        </h4>
        <ElTable :data="detail.currentTasks" border stripe>
          <ElTableColumn prop="nodeName" label="节点" width="150" />
          <ElTableColumn prop="assigneeName" label="审批人" width="120" />
          <ElTableColumn prop="status" label="状态" width="100">
            <template #default="{ row }">
              <ElTag :type="taskStatusConfig[row.status]?.type || 'info'" size="small">
                {{ taskStatusConfig[row.status]?.text || row.status }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn prop="createdAt" label="创建时间" width="180" />
          <ElTableColumn prop="dueAt" label="截止时间" width="180">
            <template #default="{ row }">
              <span
                v-if="row.dueAt"
                :style="{ color: isOverdue(row.dueAt) ? '#f56c6c' : 'inherit' }"
              >
                {{ row.dueAt }}
              </span>
              <span v-else>-</span>
            </template>
          </ElTableColumn>
        </ElTable>
      </div>
    </ElCard>

    <!-- 审批意见对话框 -->
    <ElDialog v-model="showApproveDialog" :title="approveDialogTitle" width="500px">
      <ElForm :model="approveForm" label-width="80px">
        <ElFormItem label="审批意见">
          <ElInput
            v-model="approveForm.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见（可选）"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showApproveDialog = false">取消</ElButton>
        <ElButton :type="approveAction === 'approve' ? 'success' : 'danger'" @click="submitApprove">
          确定
        </ElButton>
      </template>
    </ElDialog>

    <!-- 委托对话框 -->
    <ElDialog v-model="showDelegateDialog" title="委托任务" width="500px">
      <ElForm :model="delegateForm" label-width="80px">
        <ElFormItem label="委托给" required>
          <ElSelect
            v-model="delegateForm.toUserId"
            filterable
            remote
            :remote-method="searchUsers"
            :loading="userSearchLoading"
            placeholder="请选择委托人"
            style="width: 100%"
          >
            <ElOption
              v-for="user in userOptions"
              :key="user.id"
              :label="user.nickName || user.userName"
              :value="user.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="委托原因">
          <ElInput
            v-model="delegateForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入委托原因（可选）"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showDelegateDialog = false">取消</ElButton>
        <ElButton type="primary" @click="submitDelegate">确定</ElButton>
      </template>
    </ElDialog>

    <!-- 转办对话框 -->
    <ElDialog v-model="showTransferDialog" title="转办任务" width="500px">
      <ElForm :model="transferForm" label-width="80px">
        <ElFormItem label="转办给" required>
          <ElSelect
            v-model="transferForm.toUserId"
            filterable
            remote
            :remote-method="searchUsers"
            :loading="userSearchLoading"
            placeholder="请选择转办人"
            style="width: 100%"
          >
            <ElOption
              v-for="user in userOptions"
              :key="user.id"
              :label="user.nickName || user.userName"
              :value="user.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="转办原因">
          <ElInput
            v-model="transferForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入转办原因（可选）"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showTransferDialog = false">取消</ElButton>
        <ElButton type="primary" @click="submitTransfer">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useRoute, useRouter } from 'vue-router'
  import {
    processInstApi,
    type ProcessInstDetailResponse,
    type ApprovalRecord
  } from '@/api/workflow'
  import { fetchGetUserList } from '@/api/system-manage'
  import { ElMessage } from 'element-plus'
  import { Check, Close, Switch, Right, Document, List, Bell } from '@element-plus/icons-vue'

  defineOptions({ name: 'WorkbenchDetail' })

  const route = useRoute()
  const router = useRouter()

  // 流程实例ID
  const instanceId = computed(() => Number(route.params.id))
  // 任务ID（从待办进入时会带上）
  const taskId = computed(() => (route.query.taskId ? Number(route.query.taskId) : null))

  // 状态
  const loading = ref(false)
  const detail = ref<ProcessInstDetailResponse | null>(null)

  // 状态配置
  const statusConfig: Record<
    string,
    { type: 'warning' | 'success' | 'info' | 'danger'; text: string }
  > = {
    running: { type: 'warning', text: '进行中' },
    completed: { type: 'success', text: '已完成' },
    rejected: { type: 'danger', text: '已拒绝' },
    withdrawn: { type: 'info', text: '已撤回' }
  }

  const taskStatusConfig: Record<
    string,
    { type: 'warning' | 'success' | 'info' | 'danger'; text: string }
  > = {
    pending: { type: 'warning', text: '待处理' },
    approved: { type: 'success', text: '已通过' },
    rejected: { type: 'danger', text: '已拒绝' },
    delegated: { type: 'info', text: '已委托' },
    transferred: { type: 'info', text: '已转办' }
  }

  // 当前节点名称
  const currentNodeName = computed(() => {
    if (!detail.value?.currentTasks?.length) return null
    return detail.value.currentTasks[0]?.nodeName
  })

  // 是否可以审批（有待办任务且taskId匹配）
  const canApprove = computed(() => {
    if (!taskId.value || !detail.value?.currentTasks?.length) return false
    return detail.value.currentTasks.some((t) => t.id === taskId.value && t.status === 'pending')
  })

  // 审批对话框
  const showApproveDialog = ref(false)
  const approveAction = ref<'approve' | 'reject'>('approve')
  const approveDialogTitle = computed(() =>
    approveAction.value === 'approve' ? '审批通过' : '审批拒绝'
  )
  const approveForm = ref({
    comment: ''
  })

  // 委托对话框
  const showDelegateDialog = ref(false)
  const delegateForm = ref({
    toUserId: undefined as number | undefined,
    reason: ''
  })

  // 转办对话框
  const showTransferDialog = ref(false)
  const transferForm = ref({
    toUserId: undefined as number | undefined,
    reason: ''
  })

  // 用户搜索
  const userSearchLoading = ref(false)
  const userOptions = ref<any[]>([])

  // 加载详情
  const loadDetail = async () => {
    if (!instanceId.value) return

    loading.value = true
    try {
      detail.value = await processInstApi.getDetail(instanceId.value)
    } catch (error: any) {
      console.error('加载详情失败:', error)
      ElMessage.error(error.message || '加载详情失败')
    } finally {
      loading.value = false
    }
  }

  // 格式化表单值
  const formatFormValue = (value: any): string => {
    if (value === null || value === undefined) return '-'
    if (typeof value === 'object') return JSON.stringify(value)
    return String(value)
  }

  // 获取时间线类型
  const getTimelineType = (
    action: string
  ): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
    const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
      approve: 'success',
      reject: 'danger',
      delegate: 'warning',
      transfer: 'warning',
      withdraw: 'info',
      start: 'primary'
    }
    return typeMap[action] || 'info'
  }

  // 获取操作标签类型
  const getActionTagType = (action: string): 'success' | 'danger' | 'warning' | 'info' => {
    const typeMap: Record<string, 'success' | 'danger' | 'warning' | 'info'> = {
      approve: 'success',
      reject: 'danger',
      delegate: 'warning',
      transfer: 'warning',
      withdraw: 'info',
      start: 'info'
    }
    return typeMap[action] || 'info'
  }

  // 获取操作文本
  const getActionText = (action: string): string => {
    const textMap: Record<string, string> = {
      approve: '通过',
      reject: '拒绝',
      delegate: '委托',
      transfer: '转办',
      withdraw: '撤回',
      start: '发起'
    }
    return textMap[action] || action
  }

  // 是否超时
  const isOverdue = (dueAt: string): boolean => {
    return new Date(dueAt) < new Date()
  }

  // 搜索用户
  const searchUsers = async (query: string) => {
    if (!query) {
      userOptions.value = []
      return
    }

    userSearchLoading.value = true
    try {
      const res = await fetchGetUserList({
        page: 1,
        pageSize: 20,
        userName: query
      })
      userOptions.value = res.list || []
    } catch (error) {
      console.error('搜索用户失败:', error)
      userOptions.value = []
    } finally {
      userSearchLoading.value = false
    }
  }

  // 审批通过
  const handleApprove = () => {
    approveAction.value = 'approve'
    approveForm.value.comment = ''
    showApproveDialog.value = true
  }

  // 审批拒绝
  const handleReject = () => {
    approveAction.value = 'reject'
    approveForm.value.comment = ''
    showApproveDialog.value = true
  }

  // 提交审批
  const submitApprove = async () => {
    if (!taskId.value) {
      ElMessage.error('任务ID不存在')
      return
    }

    try {
      await processInstApi.completeTask({
        taskId: taskId.value,
        action: approveAction.value,
        comment: approveForm.value.comment
      })
      ElMessage.success(approveAction.value === 'approve' ? '审批通过成功' : '审批拒绝成功')
      showApproveDialog.value = false
      // 返回待办列表
      router.push('/workflow/workbench/todo')
    } catch (error: any) {
      console.error('审批失败:', error)
      ElMessage.error(error.message || '审批失败')
    }
  }

  // 提交委托
  const submitDelegate = async () => {
    if (!taskId.value) {
      ElMessage.error('任务ID不存在')
      return
    }
    if (!delegateForm.value.toUserId) {
      ElMessage.warning('请选择委托人')
      return
    }

    try {
      await processInstApi.delegateTask({
        taskId: taskId.value,
        toUserId: delegateForm.value.toUserId,
        reason: delegateForm.value.reason
      })
      ElMessage.success('委托成功')
      showDelegateDialog.value = false
      // 返回待办列表
      router.push('/workflow/workbench/todo')
    } catch (error: any) {
      console.error('委托失败:', error)
      ElMessage.error(error.message || '委托失败')
    }
  }

  // 提交转办
  const submitTransfer = async () => {
    if (!taskId.value) {
      ElMessage.error('任务ID不存在')
      return
    }
    if (!transferForm.value.toUserId) {
      ElMessage.warning('请选择转办人')
      return
    }

    try {
      await processInstApi.transferTask({
        taskId: taskId.value,
        toUserId: transferForm.value.toUserId,
        reason: transferForm.value.reason
      })
      ElMessage.success('转办成功')
      showTransferDialog.value = false
      // 返回待办列表
      router.push('/workflow/workbench/todo')
    } catch (error: any) {
      console.error('转办失败:', error)
      ElMessage.error(error.message || '转办失败')
    }
  }

  // 初始化
  onMounted(() => {
    loadDetail()
  })
</script>

<style lang="scss" scoped>
  .workbench-detail-page {
    padding: 16px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .title-section {
      display: flex;
      align-items: center;
      gap: 12px;

      h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
      }
    }

    .action-section {
      display: flex;
      gap: 8px;
    }
  }

  .info-section {
    margin-bottom: 24px;
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin: 24px 0 16px;
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-primary);

    .el-icon {
      color: var(--el-color-primary);
    }
  }

  .form-section {
    margin-bottom: 24px;
  }

  .trail-section {
    margin-bottom: 24px;

    .trail-card {
      :deep(.el-card__body) {
        padding: 12px 16px;
      }
    }

    .trail-content {
      .trail-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;

        .operator {
          font-weight: 500;
        }
      }

      .trail-node {
        font-size: 13px;
        color: var(--el-text-color-secondary);
        margin-bottom: 4px;
      }

      .trail-comment {
        font-size: 13px;
        color: var(--el-text-color-regular);
        background: var(--el-fill-color-light);
        padding: 8px 12px;
        border-radius: 4px;
        margin-top: 8px;
      }
    }
  }

  .tasks-section {
    margin-bottom: 24px;
  }
</style>
