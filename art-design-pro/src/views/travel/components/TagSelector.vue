<template>
  <div class="tag-selector">
    <!-- 已选标签 -->
    <div class="selected-tags">
      <ElTag
        v-for="tag in selectedTags"
        :key="tag.id"
        closable
        :color="tag.color"
        effect="light"
        @close="handleRemoveTag(tag.id)"
      >
        {{ tag.name }}
      </ElTag>
      
      <!-- 添加标签按钮 -->
      <ElPopover
        v-model:visible="popoverVisible"
        placement="bottom-start"
        :width="300"
        trigger="click"
      >
        <template #reference>
          <ElButton class="add-tag-btn" size="small" :icon="Plus">
            添加标签
          </ElButton>
        </template>

        <div class="tag-popover">
          <!-- 搜索/创建输入框 -->
          <ElInput
            v-model="searchKeyword"
            placeholder="搜索或创建标签"
            :prefix-icon="Search"
            clearable
            size="small"
            @input="handleSearch"
          />

          <!-- 系统标签 -->
          <div v-if="systemTags.length > 0" class="tag-section">
            <div class="section-title">系统标签</div>
            <div class="tag-list">
              <ElTag
                v-for="tag in filteredSystemTags"
                :key="tag.id"
                :color="tag.color"
                effect="light"
                class="tag-item"
                :class="{ selected: isSelected(tag.id) }"
                @click="handleToggleTag(tag)"
              >
                {{ tag.name }}
                <ElIcon v-if="isSelected(tag.id)" class="check-icon"><Check /></ElIcon>
              </ElTag>
            </div>
          </div>

          <!-- 自定义标签 -->
          <div v-if="customTags.length > 0" class="tag-section">
            <div class="section-title">自定义标签</div>
            <div class="tag-list">
              <ElTag
                v-for="tag in filteredCustomTags"
                :key="tag.id"
                :color="tag.color"
                effect="light"
                class="tag-item"
                :class="{ selected: isSelected(tag.id) }"
                @click="handleToggleTag(tag)"
              >
                {{ tag.name }}
                <ElIcon v-if="isSelected(tag.id)" class="check-icon"><Check /></ElIcon>
              </ElTag>
            </div>
          </div>

          <!-- 创建新标签 -->
          <div v-if="canCreateTag" class="create-tag">
            <ElButton
              type="primary"
              link
              :icon="Plus"
              @click="handleCreateTag"
              :loading="creating"
            >
              创建标签 "{{ searchKeyword }}"
            </ElButton>
          </div>

          <!-- 空状态 -->
          <div v-if="showEmpty" class="empty-state">
            <span>暂无标签</span>
          </div>
        </div>
      </ElPopover>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Search, Check } from '@element-plus/icons-vue'
import { tagApi } from '../api'
import type { Tag } from '../types'

defineOptions({
  name: 'TagSelector'
})

// Props
interface Props {
  modelValue?: number[]
  maxCount?: number
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [],
  maxCount: 5
})

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: number[]): void
}>()

// 状态
const popoverVisible = ref(false)
const searchKeyword = ref('')
const allTags = ref<Tag[]>([])
const loading = ref(false)
const creating = ref(false)

// 计算属性
const selectedTagIds = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const selectedTags = computed(() => {
  if (!allTags.value || !Array.isArray(allTags.value)) return []
  return allTags.value.filter(tag => selectedTagIds.value.includes(tag.id))
})

const systemTags = computed(() => {
  if (!allTags.value || !Array.isArray(allTags.value)) return []
  return allTags.value.filter(tag => tag.isSystem)
})

const customTags = computed(() => {
  if (!allTags.value || !Array.isArray(allTags.value)) return []
  return allTags.value.filter(tag => !tag.isSystem)
})

const filteredSystemTags = computed(() => {
  if (!searchKeyword.value) return systemTags.value
  const keyword = searchKeyword.value.toLowerCase()
  return systemTags.value.filter(tag => 
    tag.name.toLowerCase().includes(keyword)
  )
})

const filteredCustomTags = computed(() => {
  if (!searchKeyword.value) return customTags.value
  const keyword = searchKeyword.value.toLowerCase()
  return customTags.value.filter(tag => 
    tag.name.toLowerCase().includes(keyword)
  )
})

const canCreateTag = computed(() => {
  if (!searchKeyword.value.trim()) return false
  // 检查是否已存在同名标签
  const keyword = searchKeyword.value.trim().toLowerCase()
  return !allTags.value.some(tag => tag.name.toLowerCase() === keyword)
})

const showEmpty = computed(() => {
  return !loading.value && 
    filteredSystemTags.value.length === 0 && 
    filteredCustomTags.value.length === 0 &&
    !canCreateTag.value
})

// 方法
const isSelected = (tagId: number) => {
  return selectedTagIds.value.includes(tagId)
}

const handleToggleTag = (tag: Tag) => {
  const index = selectedTagIds.value.indexOf(tag.id)
  if (index >= 0) {
    // 移除
    const newIds = [...selectedTagIds.value]
    newIds.splice(index, 1)
    selectedTagIds.value = newIds
  } else {
    // 添加
    if (selectedTagIds.value.length >= props.maxCount) {
      ElMessage.warning(`最多只能选择 ${props.maxCount} 个标签`)
      return
    }
    selectedTagIds.value = [...selectedTagIds.value, tag.id]
  }
}

const handleRemoveTag = (tagId: number) => {
  selectedTagIds.value = selectedTagIds.value.filter(id => id !== tagId)
}

const handleSearch = () => {
  // 搜索逻辑已通过计算属性实现
}

const handleCreateTag = async () => {
  if (!searchKeyword.value.trim()) return

  creating.value = true
  try {
    const newTag = await tagApi.create({
      name: searchKeyword.value.trim(),
      color: getRandomColor()
    }) as unknown as Tag
    
    allTags.value.push(newTag)
    selectedTagIds.value = [...selectedTagIds.value, newTag.id]
    searchKeyword.value = ''
    ElMessage.success('标签创建成功')
  } catch (err: any) {
    ElMessage.error(err.message || '创建标签失败')
  } finally {
    creating.value = false
  }
}

// 生成随机颜色
const getRandomColor = () => {
  const colors = [
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
    '#00CED1', '#FF69B4', '#9370DB', '#20B2AA', '#FF6347'
  ]
  return colors[Math.floor(Math.random() * colors.length)]
}

// 加载标签列表
const fetchTags = async () => {
  loading.value = true
  try {
    const res = await tagApi.getList()
    // 确保返回的是数组
    allTags.value = Array.isArray(res) ? res : []
  } catch (err) {
    console.error('Failed to fetch tags:', err)
    allTags.value = []
  } finally {
    loading.value = false
  }
}

// 监听弹窗打开
watch(popoverVisible, (visible) => {
  if (visible && allTags.value.length === 0) {
    fetchTags()
  }
})

onMounted(() => {
  fetchTags()
})
</script>

<style scoped lang="scss">
.tag-selector {
  width: 100%;

  .selected-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;

    .el-tag {
      border: none;
    }

    .add-tag-btn {
      border-style: dashed;
    }
  }
}

.tag-popover {
  .el-input {
    margin-bottom: 12px;
  }

  .tag-section {
    margin-bottom: 12px;

    .section-title {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-bottom: 8px;
    }

    .tag-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .tag-item {
        cursor: pointer;
        border: none;
        transition: all 0.2s;

        &:hover {
          transform: scale(1.05);
        }

        &.selected {
          box-shadow: 0 0 0 2px var(--el-color-primary);
        }

        .check-icon {
          margin-left: 4px;
          font-size: 12px;
        }
      }
    }
  }

  .create-tag {
    padding-top: 8px;
    border-top: 1px solid var(--el-border-color-lighter);
  }

  .empty-state {
    text-align: center;
    padding: 20px;
    color: var(--el-text-color-secondary);
    font-size: 14px;
  }
}
</style>
