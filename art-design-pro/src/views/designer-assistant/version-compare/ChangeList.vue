<template>
  <div class="change-list">
    <ElTable :data="changes" stripe size="small">
      <ElTableColumn prop="field" label="字段" width="150" />
      <ElTableColumn prop="changeType" label="变化类型" width="100">
        <template #default="{ row }">
          <ElTag :type="getChangeTypeTag(row.changeType)" size="small">
            {{ getChangeTypeText(row.changeType) }}
          </ElTag>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="oldValue" label="原值" min-width="150">
        <template #default="{ row }">
          <span :class="{ 'deleted-value': row.changeType === 'removed' }">
            {{ row.oldValue || '-' }}
          </span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="" width="50">
        <template #default>
          <ElIcon><Right /></ElIcon>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="newValue" label="新值" min-width="150">
        <template #default="{ row }">
          <span :class="{ 'added-value': row.changeType === 'added' }">
            {{ row.newValue || '-' }}
          </span>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="description" label="说明" min-width="200" show-overflow-tooltip />
    </ElTable>
  </div>
</template>

<script setup lang="ts">
/***
 * Change List Component
 * 变化列表组件
 * Requirements: 7.2
 ***/
import { Right } from '@element-plus/icons-vue'
import type { ChangeItemResponse } from '@/api/designer-version'

defineProps<{
  changes: ChangeItemResponse[]
}>()

/*** Get change type tag ***/
type TagType = 'success' | 'warning' | 'info' | 'danger' | 'primary'
const getChangeTypeTag = (type: string): TagType => {
  const map: Record<string, TagType> = {
    added: 'success',
    removed: 'danger',
    modified: 'warning'
  }
  return map[type] || 'info'
}

/*** Get change type text ***/
const getChangeTypeText = (type: string) => {
  const map: Record<string, string> = {
    added: '新增',
    removed: '删除',
    modified: '修改'
  }
  return map[type] || type
}
</script>

<style scoped lang="scss">
.change-list {
  .deleted-value {
    color: var(--el-color-danger);
    text-decoration: line-through;
  }

  .added-value {
    color: var(--el-color-success);
    font-weight: 500;
  }
}
</style>
