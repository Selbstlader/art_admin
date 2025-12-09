<template>
  <div class="node-property-panel">
    <!-- 基础属性 -->
    <el-divider content-position="left">基础属性</el-divider>
    <el-form label-position="top" size="small">
      <el-form-item label="节点ID">
        <el-input :model-value="node.id" disabled />
      </el-form-item>
      <el-form-item label="节点名称">
        <el-input 
          :model-value="node.data?.label" 
          placeholder="请输入节点名称"
          @input="updateLabel"
        />
      </el-form-item>
    </el-form>

    <!-- 审批节点特有属性 -->
    <template v-if="node.type === 'approval'">
      <el-divider content-position="left">审批人配置</el-divider>
      <el-form label-position="top" size="small">
        <el-form-item label="审批人规则">
          <el-select 
            :model-value="properties.assigneeRule?.type || 'user'"
            placeholder="请选择审批人规则"
            @change="updateAssigneeType"
          >
            <el-option label="指定用户" value="user" />
            <el-option label="指定角色" value="role" />
            <el-option label="部门负责人" value="dept_leader" />
            <el-option label="发起人上级" value="initiator_leader" />
          </el-select>
        </el-form-item>

        <!-- 指定用户时显示用户选择 -->
        <el-form-item 
          v-if="properties.assigneeRule?.type === 'user'" 
          label="选择用户"
        >
          <el-select
            :model-value="properties.assigneeRule?.values || []"
            multiple
            filterable
            placeholder="请选择用户"
            @change="updateAssigneeValues"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>

        <!-- 指定角色时显示角色选择 -->
        <el-form-item 
          v-if="properties.assigneeRule?.type === 'role'" 
          label="选择角色"
        >
          <el-select
            :model-value="properties.assigneeRule?.values || []"
            multiple
            filterable
            placeholder="请选择角色"
            @change="updateAssigneeValues"
          >
            <el-option
              v-for="role in roleOptions"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="审批模式">
          <el-radio-group 
            :model-value="properties.approvalMode || 'or_sign'"
            @change="updateApprovalMode"
          >
            <el-radio value="or_sign">
              <span>或签</span>
              <el-tooltip content="任一审批人通过即可">
                <el-icon class="tip-icon"><QuestionFilled /></el-icon>
              </el-tooltip>
            </el-radio>
            <el-radio value="and_sign">
              <span>会签</span>
              <el-tooltip content="所有审批人都需通过">
                <el-icon class="tip-icon"><QuestionFilled /></el-icon>
              </el-tooltip>
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="超时时间（小时）">
          <el-input-number
            :model-value="properties.timeoutHours || 24"
            :min="1"
            :max="720"
            controls-position="right"
            @change="updateTimeoutHours"
          />
        </el-form-item>

        <el-form-item label="拒绝处理">
          <el-radio-group 
            :model-value="properties.rejectAction || 'terminate'"
            @change="updateRejectAction"
          >
            <el-radio value="terminate">终止流程</el-radio>
            <el-radio value="return_prev">退回上一节点</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <el-divider content-position="left">字段权限</el-divider>
      <el-form label-position="top" size="small">
        <el-form-item>
          <template #label>
            <span>字段权限配置</span>
            <el-tooltip content="配置该节点下表单字段的可见性和可编辑性">
              <el-icon class="tip-icon"><QuestionFilled /></el-icon>
            </el-tooltip>
          </template>
          <div class="field-permissions">
            <div 
              v-for="(permission, fieldKey) in properties.fieldPermissions || {}" 
              :key="fieldKey"
              class="permission-item"
            >
              <span class="field-key">{{ fieldKey }}</span>
              <el-select 
                :model-value="permission"
                size="small"
                @change="(val: string) => updateFieldPermission(fieldKey, val)"
              >
                <el-option label="可见" value="visible" />
                <el-option label="可编辑" value="editable" />
                <el-option label="隐藏" value="hidden" />
              </el-select>
              <el-button 
                type="danger" 
                :icon="Delete" 
                size="small" 
                circle
                @click="removeFieldPermission(fieldKey)"
              />
            </div>
            <el-button type="primary" size="small" @click="addFieldPermission">
              <el-icon><Plus /></el-icon>
              添加字段权限
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </template>

    <!-- 条件节点特有属性 -->
    <template v-if="node.type === 'condition'">
      <el-divider content-position="left">条件配置</el-divider>
      <el-form label-position="top" size="small">
        <el-form-item>
          <template #label>
            <span>条件表达式</span>
            <el-tooltip content="使用表单字段进行条件判断，如：amount > 1000">
              <el-icon class="tip-icon"><QuestionFilled /></el-icon>
            </el-tooltip>
          </template>
          <el-input
            :model-value="properties.condition || ''"
            type="textarea"
            :rows="3"
            placeholder="如：amount > 1000"
            @input="updateCondition"
          />
        </el-form-item>
      </el-form>
    </template>

    <!-- 开始/结束节点 -->
    <template v-if="node.type === 'start' || node.type === 'end'">
      <div class="simple-node-info">
        <el-icon :size="48" :color="node.type === 'start' ? '#67c23a' : '#f56c6c'">
          <component :is="node.type === 'start' ? VideoPlay : CircleClose" />
        </el-icon>
        <p>{{ node.type === 'start' ? '流程开始节点' : '流程结束节点' }}</p>
        <p class="hint">{{ node.type === 'start' ? '流程从此节点开始执行' : '流程在此节点结束' }}</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import type { Node } from '@vue-flow/core'
import { 
  QuestionFilled, 
  Delete, 
  Plus,
  VideoPlay,
  CircleClose
} from '@element-plus/icons-vue'
import type { NodeProperties } from '@/api/workflow'

defineOptions({ name: 'NodePropertyPanel' })

interface Props {
  node: Node
}

const props = defineProps<Props>()

const emit = defineEmits<{
  update: [nodeId: string, data: any]
}>()

// 节点属性
const properties = computed<NodeProperties>(() => {
  return props.node.data?.properties || {}
})

// 用户和角色选项（实际项目中应从API获取）
const userOptions = ref<{ id: number; name: string }[]>([
  { id: 1, name: '张三' },
  { id: 2, name: '李四' },
  { id: 3, name: '王五' },
  { id: 4, name: '赵六' }
])

const roleOptions = ref<{ id: number; name: string }[]>([
  { id: 1, name: '部门经理' },
  { id: 2, name: '财务主管' },
  { id: 3, name: '人事主管' },
  { id: 4, name: '总经理' }
])

// 更新节点标签
const updateLabel = (value: string) => {
  emit('update', props.node.id, {
    label: value,
    properties: properties.value
  })
}

// 更新审批人类型
const updateAssigneeType = (type: string) => {
  const newProperties = {
    ...properties.value,
    assigneeRule: {
      type,
      values: []
    }
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新审批人值
const updateAssigneeValues = (values: number[]) => {
  const newProperties = {
    ...properties.value,
    assigneeRule: {
      ...properties.value.assigneeRule,
      values
    }
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新审批模式
const updateApprovalMode = (mode: string) => {
  const newProperties = {
    ...properties.value,
    approvalMode: mode
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新超时时间
const updateTimeoutHours = (hours: number) => {
  const newProperties = {
    ...properties.value,
    timeoutHours: hours
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新拒绝处理
const updateRejectAction = (action: string) => {
  const newProperties = {
    ...properties.value,
    rejectAction: action
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新条件表达式
const updateCondition = (condition: string) => {
  const newProperties = {
    ...properties.value,
    condition
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 添加字段权限
const addFieldPermission = () => {
  const fieldKey = `field_${Date.now()}`
  const newPermissions = {
    ...properties.value.fieldPermissions,
    [fieldKey]: 'visible'
  }
  const newProperties = {
    ...properties.value,
    fieldPermissions: newPermissions
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 更新字段权限
const updateFieldPermission = (fieldKey: string, permission: string) => {
  const newPermissions = {
    ...properties.value.fieldPermissions,
    [fieldKey]: permission
  }
  const newProperties = {
    ...properties.value,
    fieldPermissions: newPermissions
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}

// 删除字段权限
const removeFieldPermission = (fieldKey: string) => {
  const newPermissions = { ...properties.value.fieldPermissions }
  delete newPermissions[fieldKey]
  const newProperties = {
    ...properties.value,
    fieldPermissions: newPermissions
  }
  emit('update', props.node.id, {
    label: props.node.data?.label,
    properties: newProperties
  })
}
</script>

<style lang="scss" scoped>
.node-property-panel {
  :deep(.el-divider) {
    margin: 16px 0 12px;

    .el-divider__text {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }

  :deep(.el-form-item) {
    margin-bottom: 12px;

    .el-form-item__label {
      display: flex;
      gap: 4px;
      align-items: center;
      font-size: 13px;
      color: var(--el-text-color-regular);
    }
  }

  .tip-icon {
    color: var(--el-text-color-placeholder);
    cursor: help;
  }

  .field-permissions {
    .permission-item {
      display: flex;
      gap: 8px;
      align-items: center;
      margin-bottom: 8px;

      .field-key {
        flex: 1;
        font-size: 12px;
        color: var(--el-text-color-regular);
      }

      .el-select {
        width: 100px;
      }
    }
  }

  .simple-node-info {
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
    padding: 24px 0;
    text-align: center;

    p {
      margin: 0;
      font-size: 14px;
      color: var(--el-text-color-primary);

      &.hint {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}
</style>
