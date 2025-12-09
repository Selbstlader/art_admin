<template>
  <div class="form-designer">
    <!-- 左侧：组件面板 -->
    <div class="component-panel">
      <div class="panel-header">
        <span>组件面板</span>
      </div>
      <div class="component-list">
        <div
          v-for="comp in componentPalette"
          :key="comp.type"
          class="component-item"
          draggable="true"
          @dragstart="handleDragStart($event, comp)"
        >
          <el-icon :size="18">
            <component :is="comp.icon" />
          </el-icon>
          <span>{{ comp.label }}</span>
        </div>
      </div>
    </div>

    <!-- 中间：画布区域 -->
    <div class="canvas-panel">
      <div class="panel-header">
        <span>表单画布</span>
        <div class="canvas-actions">
          <el-button size="small" @click="handlePreview">
            <el-icon><View /></el-icon>
            预览
          </el-button>
          <el-button size="small" type="primary" @click="handleClear">
            <el-icon><Delete /></el-icon>
            清空
          </el-button>
        </div>
      </div>
      <div
        class="canvas-area"
        @dragover.prevent
        @drop="handleDrop"
      >
        <draggable
          v-model="formFields"
          item-key="key"
          group="form-fields"
          :animation="200"
          ghost-class="ghost-field"
          class="field-list"
          @change="handleFieldChange"
        >
          <template #item="{ element, index }">
            <div
              :class="['field-item', { active: selectedFieldKey === element.key }]"
              @click="selectField(element)"
            >
              <div class="field-header">
                <span class="field-type-badge">{{ getFieldTypeLabel(element.type) }}</span>
                <span class="field-label">{{ element.label || '未命名字段' }}</span>
                <el-icon class="delete-btn" @click.stop="removeField(index)">
                  <Close />
                </el-icon>
              </div>
              <div class="field-preview">
                <component
                  :is="getPreviewComponent(element.type)"
                  v-bind="getPreviewProps(element)"
                  disabled
                />
              </div>
            </div>
          </template>
        </draggable>
        <div v-if="formFields.length === 0" class="empty-canvas">
          <el-icon :size="48"><DocumentAdd /></el-icon>
          <p>从左侧拖拽组件到此处</p>
        </div>
      </div>
    </div>

    <!-- 右侧：属性面板 -->
    <div class="property-panel">
      <div class="panel-header">
        <span>属性配置</span>
      </div>
      <div v-if="selectedField" class="property-form">
        <el-form label-position="top" size="small">
          <!-- 基础属性 -->
          <el-divider content-position="left">基础属性</el-divider>
          <el-form-item label="字段标识">
            <el-input v-model="selectedField.key" placeholder="唯一标识" />
          </el-form-item>
          <el-form-item label="字段名称">
            <el-input v-model="selectedField.label" placeholder="显示名称" />
          </el-form-item>
          <el-form-item label="占位提示">
            <el-input v-model="selectedField.placeholder" placeholder="输入提示文字" />
          </el-form-item>
          <el-form-item label="默认值">
            <el-input v-model="selectedField.defaultValue" placeholder="默认值" />
          </el-form-item>
          <el-form-item label="是否必填">
            <el-switch v-model="selectedField.required" />
          </el-form-item>

          <!-- 下拉选项配置（仅select类型显示） -->
          <template v-if="selectedField.type === 'select'">
            <el-divider content-position="left">选项配置</el-divider>
            <el-form-item label="选项列表">
              <div class="options-editor">
                <div
                  v-for="(option, idx) in selectedField.options || []"
                  :key="idx"
                  class="option-row"
                >
                  <el-input
                    v-model="option.label"
                    placeholder="显示文本"
                    size="small"
                  />
                  <el-input
                    v-model="option.value"
                    placeholder="值"
                    size="small"
                  />
                  <el-button
                    type="danger"
                    :icon="Delete"
                    size="small"
                    circle
                    @click="removeOption(idx)"
                  />
                </div>
                <el-button type="primary" size="small" @click="addOption">
                  <el-icon><Plus /></el-icon>
                  添加选项
                </el-button>
              </div>
            </el-form-item>
          </template>

          <!-- 验证规则 -->
          <el-divider content-position="left">验证规则</el-divider>
          <template v-if="selectedField.type === 'text' || selectedField.type === 'textarea'">
            <el-form-item label="最小长度">
              <el-input-number
                v-model="validationRules.minLength"
                :min="0"
                controls-position="right"
                @change="updateValidation"
              />
            </el-form-item>
            <el-form-item label="最大长度">
              <el-input-number
                v-model="validationRules.maxLength"
                :min="0"
                controls-position="right"
                @change="updateValidation"
              />
            </el-form-item>
            <el-form-item label="正则表达式">
              <el-input
                v-model="validationRules.pattern"
                placeholder="如：^[a-zA-Z]+$"
                @change="updateValidation"
              />
            </el-form-item>
          </template>
          <template v-if="selectedField.type === 'number'">
            <el-form-item label="最小值">
              <el-input-number
                v-model="validationRules.min"
                controls-position="right"
                @change="updateValidation"
              />
            </el-form-item>
            <el-form-item label="最大值">
              <el-input-number
                v-model="validationRules.max"
                controls-position="right"
                @change="updateValidation"
              />
            </el-form-item>
          </template>
        </el-form>
      </div>
      <div v-else class="no-selection">
        <el-icon :size="48"><Setting /></el-icon>
        <p>请选择一个字段进行配置</p>
      </div>
    </div>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="表单预览" width="600px">
      <el-form label-position="top">
        <el-form-item
          v-for="field in formFields"
          :key="field.key"
          :label="field.label"
          :required="field.required"
        >
          <component
            :is="getPreviewComponent(field.type)"
            v-bind="getPreviewProps(field)"
          />
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import draggable from 'vuedraggable'
import {
  Edit,
  Calendar,
  List,
  Document,
  Upload,
  View,
  Delete,
  Close,
  Plus,
  Setting,
  DocumentAdd
} from '@element-plus/icons-vue'
import type { FormField, FormSchema, SelectOption, FieldValidation } from '@/api/workflow'

defineOptions({ name: 'FormDesigner' })

// Props
interface Props {
  modelValue?: FormSchema
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({ fields: [] }),
  readonly: false
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: FormSchema]
  change: [value: FormSchema]
}>()

// 组件面板配置
interface ComponentItem {
  type: string
  label: string
  icon: any
}

const componentPalette: ComponentItem[] = [
  { type: 'text', label: '文本输入', icon: Edit },
  { type: 'number', label: '数字输入', icon: Document },
  { type: 'date', label: '日期选择', icon: Calendar },
  { type: 'select', label: '下拉选择', icon: List },
  { type: 'textarea', label: '多行文本', icon: Document },
  { type: 'file', label: '文件上传', icon: Upload }
]

// 字段类型标签映射
const fieldTypeLabels: Record<string, string> = {
  text: '文本',
  number: '数字',
  date: '日期',
  select: '下拉',
  textarea: '多行',
  file: '文件'
}

// 表单字段列表
const formFields = ref<FormField[]>([])

// 选中的字段
const selectedFieldKey = ref<string | null>(null)
const selectedField = computed(() => {
  if (!selectedFieldKey.value) return null
  return formFields.value.find(f => f.key === selectedFieldKey.value) || null
})

// 验证规则（用于属性面板编辑）
const validationRules = ref<FieldValidation>({})

// 预览对话框
const previewVisible = ref(false)

// 初始化
watch(
  () => props.modelValue,
  (val) => {
    if (val?.fields) {
      formFields.value = JSON.parse(JSON.stringify(val.fields))
    }
  },
  { immediate: true, deep: true }
)

// 监听字段变化，同步到父组件
watch(
  formFields,
  (val) => {
    const schema: FormSchema = { fields: val }
    emit('update:modelValue', schema)
    emit('change', schema)
  },
  { deep: true }
)

// 监听选中字段变化，同步验证规则
watch(selectedField, (field) => {
  if (field?.validation) {
    validationRules.value = { ...field.validation }
  } else {
    validationRules.value = {}
  }
})

// 生成唯一key
const generateKey = (type: string): string => {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 6)
  return `${type}_${timestamp}_${random}`
}

// 创建默认字段
const createDefaultField = (type: string): FormField => {
  const field: FormField = {
    key: generateKey(type),
    label: `${fieldTypeLabels[type] || type}字段`,
    type,
    required: false,
    placeholder: ''
  }

  if (type === 'select') {
    field.options = [
      { label: '选项1', value: 'option1' },
      { label: '选项2', value: 'option2' }
    ]
  }

  return field
}

// 拖拽开始
const handleDragStart = (event: DragEvent, comp: ComponentItem) => {
  event.dataTransfer?.setData('fieldType', comp.type)
}

// 拖拽放置
const handleDrop = (event: DragEvent) => {
  const fieldType = event.dataTransfer?.getData('fieldType')
  if (fieldType) {
    const newField = createDefaultField(fieldType)
    formFields.value.push(newField)
    selectField(newField)
  }
}

// 字段排序变化
const handleFieldChange = () => {
  // 排序变化时自动触发watch更新
}

// 选择字段
const selectField = (field: FormField) => {
  selectedFieldKey.value = field.key
}

// 删除字段
const removeField = (index: number) => {
  const removed = formFields.value[index]
  if (removed.key === selectedFieldKey.value) {
    selectedFieldKey.value = null
  }
  formFields.value.splice(index, 1)
}

// 获取字段类型标签
const getFieldTypeLabel = (type: string): string => {
  return fieldTypeLabels[type] || type
}

// 获取预览组件
const getPreviewComponent = (type: string): string => {
  const componentMap: Record<string, string> = {
    text: 'el-input',
    number: 'el-input-number',
    date: 'el-date-picker',
    select: 'el-select',
    textarea: 'el-input',
    file: 'el-upload'
  }
  return componentMap[type] || 'el-input'
}

// 获取预览组件属性
const getPreviewProps = (field: FormField): Record<string, any> => {
  const baseProps: Record<string, any> = {
    placeholder: field.placeholder || `请输入${field.label}`
  }

  switch (field.type) {
    case 'textarea':
      baseProps.type = 'textarea'
      baseProps.rows = 3
      break
    case 'date':
      baseProps.type = 'date'
      baseProps.valueFormat = 'YYYY-MM-DD'
      break
    case 'select':
      // 下拉选择需要特殊处理
      break
    case 'file':
      baseProps.action = '#'
      baseProps.autoUpload = false
      break
  }

  return baseProps
}

// 添加选项
const addOption = () => {
  if (selectedField.value && selectedField.value.type === 'select') {
    if (!selectedField.value.options) {
      selectedField.value.options = []
    }
    selectedField.value.options.push({
      label: `选项${selectedField.value.options.length + 1}`,
      value: `option${selectedField.value.options.length + 1}`
    })
  }
}

// 删除选项
const removeOption = (index: number) => {
  if (selectedField.value?.options) {
    selectedField.value.options.splice(index, 1)
  }
}

// 更新验证规则
const updateValidation = () => {
  if (selectedField.value) {
    const validation: FieldValidation = {}
    if (validationRules.value.minLength !== undefined && validationRules.value.minLength > 0) {
      validation.minLength = validationRules.value.minLength
    }
    if (validationRules.value.maxLength !== undefined && validationRules.value.maxLength > 0) {
      validation.maxLength = validationRules.value.maxLength
    }
    if (validationRules.value.min !== undefined) {
      validation.min = validationRules.value.min
    }
    if (validationRules.value.max !== undefined) {
      validation.max = validationRules.value.max
    }
    if (validationRules.value.pattern) {
      validation.pattern = validationRules.value.pattern
    }
    selectedField.value.validation = Object.keys(validation).length > 0 ? validation : undefined
  }
}

// 预览
const handlePreview = () => {
  previewVisible.value = true
}

// 清空
const handleClear = () => {
  formFields.value = []
  selectedFieldKey.value = null
}

// 暴露方法
defineExpose({
  getSchema: () => ({ fields: formFields.value }),
  setSchema: (schema: FormSchema) => {
    formFields.value = JSON.parse(JSON.stringify(schema.fields || []))
    selectedFieldKey.value = null
  },
  clear: handleClear
})
</script>

<style lang="scss" scoped>
.form-designer {
  display: flex;
  height: 100%;
  min-height: 500px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 40px;
    padding: 0 12px;
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-primary);
    background: var(--el-fill-color-light);
    border-bottom: 1px solid var(--el-border-color-light);
  }

  // 左侧组件面板
  .component-panel {
    width: 200px;
    border-right: 1px solid var(--el-border-color-light);

    .component-list {
      padding: 12px;
    }

    .component-item {
      display: flex;
      gap: 8px;
      align-items: center;
      padding: 10px 12px;
      margin-bottom: 8px;
      font-size: 13px;
      cursor: grab;
      background: var(--el-fill-color-lighter);
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 4px;
      transition: all 0.2s;

      &:hover {
        color: var(--el-color-primary);
        background: var(--el-color-primary-light-9);
        border-color: var(--el-color-primary-light-5);
      }

      &:active {
        cursor: grabbing;
      }
    }
  }

  // 中间画布区域
  .canvas-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 400px;

    .canvas-actions {
      display: flex;
      gap: 8px;
    }

    .canvas-area {
      flex: 1;
      padding: 16px;
      overflow-y: auto;
      background: var(--el-fill-color-lighter);
    }

    .field-list {
      min-height: 100%;
    }

    .field-item {
      padding: 12px;
      margin-bottom: 12px;
      cursor: pointer;
      background: var(--el-bg-color);
      border: 2px solid var(--el-border-color-light);
      border-radius: 6px;
      transition: all 0.2s;

      &:hover {
        border-color: var(--el-color-primary-light-5);
      }

      &.active {
        border-color: var(--el-color-primary);
        box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
      }

      .field-header {
        display: flex;
        gap: 8px;
        align-items: center;
        margin-bottom: 8px;

        .field-type-badge {
          padding: 2px 8px;
          font-size: 12px;
          color: var(--el-color-primary);
          background: var(--el-color-primary-light-9);
          border-radius: 4px;
        }

        .field-label {
          flex: 1;
          font-size: 14px;
          font-weight: 500;
        }

        .delete-btn {
          padding: 4px;
          color: var(--el-text-color-secondary);
          cursor: pointer;
          border-radius: 4px;
          transition: all 0.2s;

          &:hover {
            color: var(--el-color-danger);
            background: var(--el-color-danger-light-9);
          }
        }
      }

      .field-preview {
        pointer-events: none;

        :deep(.el-input),
        :deep(.el-select),
        :deep(.el-date-picker),
        :deep(.el-input-number) {
          width: 100%;
        }
      }
    }

    .ghost-field {
      background: var(--el-color-primary-light-9);
      border: 2px dashed var(--el-color-primary);
      opacity: 0.8;
    }

    .empty-canvas {
      display: flex;
      flex-direction: column;
      gap: 12px;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: var(--el-text-color-secondary);

      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }

  // 右侧属性面板
  .property-panel {
    width: 280px;
    border-left: 1px solid var(--el-border-color-light);

    .property-form {
      padding: 12px;
      overflow-y: auto;
      max-height: calc(100% - 40px);

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
          font-size: 13px;
          color: var(--el-text-color-regular);
        }
      }

      .options-editor {
        .option-row {
          display: flex;
          gap: 8px;
          align-items: center;
          margin-bottom: 8px;

          .el-input {
            flex: 1;
          }
        }
      }
    }

    .no-selection {
      display: flex;
      flex-direction: column;
      gap: 12px;
      align-items: center;
      justify-content: center;
      height: calc(100% - 40px);
      color: var(--el-text-color-secondary);

      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }
}
</style>
