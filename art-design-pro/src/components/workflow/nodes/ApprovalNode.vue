<template>
  <div :class="['approval-node', { selected }]">
    <Handle type="target" :position="Position.Top" />
    <div class="node-header">
      <el-icon :size="16"><CircleCheck /></el-icon>
      <span>审批节点</span>
    </div>
    <div class="node-body">
      <div class="node-title">{{ data?.label || '审批节点' }}</div>
      <div class="node-info">
        <span class="info-item">
          <el-icon :size="12"><User /></el-icon>
          {{ getAssigneeText() }}
        </span>
        <span class="info-item">
          <el-icon :size="12"><Clock /></el-icon>
          {{ getApprovalModeText() }}
        </span>
      </div>
    </div>
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import { CircleCheck, User, Clock } from '@element-plus/icons-vue'
import type { NodeProperties, AssigneeRule } from '@/api/workflow'

defineOptions({ name: 'ApprovalNode' })

interface Props {
  id: string
  data?: {
    label?: string
    properties?: NodeProperties
  }
  selected?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  data: () => ({}),
  selected: false
})

// 获取审批人文本
const getAssigneeText = (): string => {
  const rule = props.data?.properties?.assigneeRule
  if (!rule) return '未配置'
  
  const typeLabels: Record<string, string> = {
    user: '指定用户',
    role: '指定角色',
    dept_leader: '部门负责人',
    initiator_leader: '发起人上级'
  }
  
  return typeLabels[rule.type] || '未配置'
}

// 获取审批模式文本
const getApprovalModeText = (): string => {
  const mode = props.data?.properties?.approvalMode
  if (mode === 'and_sign') return '会签'
  return '或签'
}
</script>

<style lang="scss" scoped>
.approval-node {
  width: 180px;
  background: #fff;
  border: 2px solid #409eff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
  transition: all 0.2s;

  &:hover {
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.25);
    transform: translateY(-1px);
  }

  &.selected {
    border-color: #409eff;
    box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.3);
  }

  .node-header {
    display: flex;
    gap: 6px;
    align-items: center;
    padding: 8px 12px;
    font-size: 12px;
    font-weight: 500;
    color: #fff;
    background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
    border-radius: 6px 6px 0 0;
  }

  .node-body {
    padding: 10px 12px;

    .node-title {
      margin-bottom: 8px;
      font-size: 14px;
      font-weight: 500;
      color: var(--el-text-color-primary);
    }

    .node-info {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .info-item {
        display: flex;
        gap: 4px;
        align-items: center;
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}

:deep(.vue-flow__handle) {
  width: 10px;
  height: 10px;
  background: #409eff;
  border: 2px solid #fff;
}
</style>
