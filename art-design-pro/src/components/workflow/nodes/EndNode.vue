<template>
  <div :class="['end-node', { selected }]">
    <Handle type="target" :position="Position.Top" />
    <div class="node-content">
      <el-icon :size="24"><CircleClose /></el-icon>
      <span class="node-label">{{ data?.label || '结束' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import { CircleClose } from '@element-plus/icons-vue'

defineOptions({ name: 'EndNode' })

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
.end-node {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
  transition: all 0.2s;

  &:hover {
    box-shadow: 0 4px 12px rgba(245, 108, 108, 0.4);
    transform: scale(1.02);
  }

  &.selected {
    box-shadow: 0 0 0 3px rgba(245, 108, 108, 0.5);
  }

  .node-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
    align-items: center;
    color: #fff;

    .node-label {
      font-size: 12px;
      font-weight: 500;
    }
  }
}

:deep(.vue-flow__handle) {
  width: 12px;
  height: 12px;
  background: #f56c6c;
  border: 2px solid #fff;
}
</style>
