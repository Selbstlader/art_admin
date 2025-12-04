<template>
  <div class="image-upload">
    <el-upload
      ref="uploadRef"
      :action="uploadUrl"
      :headers="headers"
      :multiple="multiple"
      :limit="limit"
      accept="image/*"
      :file-list="fileList"
      list-type="picture-card"
      :auto-upload="true"
      :before-upload="handleBeforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-exceed="handleExceed"
      :on-remove="handleRemove"
      :on-preview="handlePreview"
      :disabled="disabled"
    >
      <el-icon><Plus /></el-icon>
      <template #tip>
        <div v-if="tip" class="el-upload__tip">{{ tip }}</div>
      </template>
    </el-upload>

    <!-- 图片预览 -->
    <el-image-viewer
      v-if="previewVisible"
      :url-list="previewList"
      :initial-index="previewIndex"
      @close="previewVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { UploadFile, UploadInstance, UploadRawFile } from 'element-plus'
import { useUserStore } from '@/store/modules/user'

interface Props {
  modelValue?: string | string[]
  multiple?: boolean
  limit?: number
  maxSize?: number // MB
  disabled?: boolean
  tip?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  multiple: false,
  limit: 1,
  maxSize: 10,
  disabled: false,
  tip: '支持 jpg、png、gif 格式，单个文件不超过 10MB'
})

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]]
  'success': [file: any]
}>()

const userStore = useUserStore()
const uploadRef = ref<UploadInstance>()
const fileList = ref<UploadFile[]>([])
const previewVisible = ref(false)
const previewList = ref<string[]>([])
const previewIndex = ref(0)

// 上传地址
const uploadUrl = computed(() => {
  const baseUrl = import.meta.env.VITE_BASE_URL || ''
  return `${baseUrl}/api/file/upload?category=image`
})

// 请求头
const headers = computed(() => ({
  Authorization: `Bearer ${userStore.accessToken}`
}))

// 初始化文件列表
watch(
  () => props.modelValue,
  (val) => {
    if (!val) {
      fileList.value = []
      return
    }
    const urls = Array.isArray(val) ? val : [val]
    fileList.value = urls.map((url, index) => ({
      name: `image-${index}`,
      url
    })) as UploadFile[]
  },
  { immediate: true }
)

// 上传前校验
const handleBeforeUpload = (file: UploadRawFile) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  const isLtSize = file.size / 1024 / 1024 < props.maxSize
  if (!isLtSize) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB`)
    return false
  }
  return true
}

// 上传成功
const handleSuccess = (response: any, file: UploadFile) => {
  if (response.code === 200) {
    const fileInfo = response.data
    file.url = fileInfo.url

    if (props.multiple) {
      const urls = Array.isArray(props.modelValue) ? [...props.modelValue] : []
      urls.push(fileInfo.url)
      emit('update:modelValue', urls)
    } else {
      emit('update:modelValue', fileInfo.url)
    }

    emit('success', fileInfo)
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleError = () => {
  ElMessage.error('上传失败')
}

// 超出限制
const handleExceed = () => {
  ElMessage.warning(`最多只能上传 ${props.limit} 张图片`)
}

// 移除文件
const handleRemove = (file: UploadFile) => {
  if (props.multiple) {
    const urls = Array.isArray(props.modelValue) ? props.modelValue.filter((url) => url !== file.url) : []
    emit('update:modelValue', urls)
  } else {
    emit('update:modelValue', '')
  }
}

// 预览图片
const handlePreview = (file: UploadFile) => {
  if (file.url) {
    previewList.value = [file.url]
    previewIndex.value = 0
    previewVisible.value = true
  }
}

// 清空
const clearFiles = () => {
  uploadRef.value?.clearFiles()
  emit('update:modelValue', props.multiple ? [] : '')
}

defineExpose({ clearFiles })
</script>

<style scoped lang="scss">
.image-upload {
  :deep(.el-upload--picture-card) {
    width: 100px;
    height: 100px;
  }
  :deep(.el-upload-list--picture-card .el-upload-list__item) {
    width: 100px;
    height: 100px;
  }
  .el-upload__tip {
    color: var(--el-text-color-secondary);
    font-size: 12px;
    margin-top: 8px;
  }
}
</style>
