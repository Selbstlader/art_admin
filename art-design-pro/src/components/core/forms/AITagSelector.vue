<template>
  <div class="ai-tag-selector">
    <ElSelect
      v-model="selectedTagId"
      :placeholder="placeholder"
      :disabled="disabled"
      :clearable="clearable"
      :loading="loading"
      filterable
      @change="handleChange"
      @visible-change="handleVisibleChange"
    >
      <ElOption
        v-for="tag in tagList"
        :key="tag.id"
        :label="tag.name"
        :value="tag.id"
      >
        <div class="tag-option">
          <div class="tag-option-main">
            <span class="tag-name">{{ tag.name }}</span>
            <ElTag
              v-if="tag.status === 1"
              type="success"
              size="small"
              effect="plain"
            >
              启用
            </ElTag>
            <ElTag
              v-else
              type="info"
              size="small"
              effect="plain"
            >
              禁用
            </ElTag>
          </div>
          <div v-if="tag.description" class="tag-description">
            {{ tag.description }}
          </div>
          <div v-if="tag.knowledge_base_name" class="tag-knowledge">
            <Icon icon="mdi:database-outline" class="kb-icon" />
            {{ tag.knowledge_base_name }}
          </div>
        </div>
      </ElOption>
    </ElSelect>
  </div>
</template>

<script setup lang="ts">
/*** AITagSelector Component ***/
/*** Requirements: 5.1, 5.2 - Select AI tag for chat session ***/

import { ref, computed, watch, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { aiTagApi, type AITag } from '@/api/ai-tag'

defineOptions({
  name: 'AITagSelector'
})

/*** Props Definition ***/
interface Props {
  modelValue?: number | null
  placeholder?: string
  disabled?: boolean
  clearable?: boolean
  onlyActive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  placeholder: '请选择AI标签',
  disabled: false,
  clearable: true,
  onlyActive: true
})

/*** Emits Definition ***/
const emit = defineEmits<{
  (e: 'update:modelValue', value: number | null): void
  (e: 'change', tag: AITag | null): void
  (e: 'tag-config', config: AITagConfig | null): void
}>()

/*** Tag Config Interface - exported for use in chat integration ***/
export interface AITagConfig {
  id: number
  name: string
  knowledge_base_id: string
  knowledge_base_name: string
  system_prompt: string
  chat_api_key: string
}

/*** State ***/
const loading = ref(false)
const tagList = ref<AITag[]>([])
const hasLoaded = ref(false)

/*** Computed ***/
const selectedTagId = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

/*** Methods ***/

/*** Fetch tag list from API ***/
const fetchTags = async () => {
  if (loading.value) return
  
  loading.value = true
  try {
    const response = await aiTagApi.list({
      current: 1,
      size: 100,
      keyword: ''
    })
    
    let tags = response.list || []
    
    /*** Filter only active tags if onlyActive is true ***/
    if (props.onlyActive) {
      tags = tags.filter(tag => tag.status === 1)
    }
    
    tagList.value = tags
    hasLoaded.value = true
  } catch (err) {
    console.error('Failed to fetch AI tags:', err)
    tagList.value = []
  } finally {
    loading.value = false
  }
}

/*** Handle selection change ***/
const handleChange = (tagId: number | null) => {
  const selectedTag = tagId ? tagList.value.find(t => t.id === tagId) : null
  
  /*** Emit change event with full tag object ***/
  emit('change', selectedTag || null)
  
  /*** Emit tag config for chat integration ***/
  if (selectedTag) {
    const config: AITagConfig = {
      id: selectedTag.id,
      name: selectedTag.name,
      knowledge_base_id: selectedTag.knowledge_base_id,
      knowledge_base_name: selectedTag.knowledge_base_name,
      system_prompt: selectedTag.system_prompt,
      chat_api_key: selectedTag.chat_api_key
    }
    emit('tag-config', config)
  } else {
    emit('tag-config', null)
  }
}

/*** Handle dropdown visibility change ***/
const handleVisibleChange = (visible: boolean) => {
  /*** Lazy load tags when dropdown opens for the first time ***/
  if (visible && !hasLoaded.value) {
    fetchTags()
  }
}

/*** Get selected tag object ***/
const getSelectedTag = (): AITag | null => {
  if (!selectedTagId.value) return null
  return tagList.value.find(t => t.id === selectedTagId.value) || null
}

/*** Refresh tag list ***/
const refresh = () => {
  hasLoaded.value = false
  fetchTags()
}

/*** Expose methods for parent component ***/
defineExpose({
  getSelectedTag,
  refresh
})

/*** Watch for external modelValue changes ***/
watch(() => props.modelValue, (newVal) => {
  /*** If a value is set but tags not loaded, fetch them ***/
  if (newVal && !hasLoaded.value) {
    fetchTags()
  }
}, { immediate: true })

/*** Lifecycle ***/
onMounted(() => {
  /*** Pre-load tags if a value is already selected ***/
  if (props.modelValue) {
    fetchTags()
  }
})
</script>

<style scoped lang="scss">
.ai-tag-selector {
  width: 100%;

  :deep(.el-select) {
    width: 100%;
  }
}

.tag-option {
  padding: 4px 0;
  line-height: 1.4;

  .tag-option-main {
    display: flex;
    align-items: center;
    gap: 8px;

    .tag-name {
      font-weight: 500;
      color: var(--el-text-color-primary);
    }

    .el-tag {
      flex-shrink: 0;
    }
  }

  .tag-description {
    margin-top: 4px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 280px;
  }

  .tag-knowledge {
    margin-top: 4px;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
    display: flex;
    align-items: center;
    gap: 4px;

    .kb-icon {
      font-size: 14px;
    }
  }
}
</style>
