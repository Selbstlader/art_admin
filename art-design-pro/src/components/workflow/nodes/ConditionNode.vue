<template>
  <div :class="['condition-node', { selected }]">
    <Handle type="target" :position="Position.Top" />
    <div class="node-diamond">
      <div class="diamond-content">
        <el-icon :size="20"><Switch /></el-icon>
        <span class="node-label">{{ data?.label || '条件' }}</span>
      </div>
    </div>
    <Handle 
      id="yes" 
      type="source" 
      :position="Position.Right" 
      class="handle-yes"
    />
    <Handle 
      id="no" 
      type="source" 
      :position="Position.Bottom" 
      class="handle-no"
    />
    <div class="branch-labels">
      <span class="label-yes">是</span>
      <span class="label-no">否</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import { Switch } from '@element-plus/icons-vue'

defineOptions({ name: 'ConditionNode' })

interface Props {
  id: string
  data?: {
    label?: string
    properties?: Record<string, any>
  }
  selected?: boolean
}

withDefaults(defineProps<Props>(), {
  data: () => ({}),
  selected: false
})
</script>

<style lang="scss" scoped>
.condition-node {
  position: relative;
  width: 100px;
  height: 100px;

  .node-diamond {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
    border: 2px solid #e6a23c;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(230, 162, 60, 0.3);
    transform: translate(-50%, -50%) rotate(45deg);
    transition: all 0.2s;
  }

  &:hover .node-diamond {
    box-shadow: 0 4px 12px rgba(230, 162, 60, 0.4);
  }

  &.selected .node-diamond {
    box-shadow: 0 0 0 3px rgba(230, 162, 60, 0.5);
  }

  .diamond-content {
    position: absolute;
    top: 50%;
    left: 50%;
    display: flex;
    flex-direction: column;
    gap: 2px;
    align-items: center;
    color: #fff;
    transform: translate(-50%, -50%);

    .node-label {
      font-size: 11px;
      font-weight: 500;
      white-space: nowrap;
    }
  }

  .branch-labels {
    position: absolute;
    width: 100%;
    height: 100%;
    pointer-events: none;

    .label-yes {
      position: absolute;
      top: 50%;
      right: -20px;
      font-size: 11px;
      color: #67c23a;
      transform: translateY(-50%);
    }

    .label-no {
      position: absolute;
      bottom: -18px;
      left: 50%;
      font-size: 11px;
      color: #f56c6c;
      transform: translateX(-50%);
    }
  }
}

:deep(.vue-flow__handle) {
  width: 10px;
  height: 10px;
  background: #e6a23c;
  border: 2px solid #fff;

  &.handle-yes {
    background: #67c23a;
  }

  &.handle-no {
    background: #f56c6c;
  }
}
</style>
