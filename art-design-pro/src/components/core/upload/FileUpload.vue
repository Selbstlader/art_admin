<template>
  <div class="file-upload">
    <el-upload
      ref="uploadRef"
      :action="uploadUrl"
      :headers="headers"
      :multiple="multiple"
      :limit="limit"
      :accept="accept"
      :file-list="fileList"
      :list-type="listType"
      :auto-upload="autoUpload"
      :show-file-list="showFileList"
      :before-upload="handleBeforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-exceed="handleExceed"
      :on-remove="handleRemove"
      :on-preview="handlePreview"
      :drag="drag"
      :disabled="disabled"
    >
      <template v-if="drag">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text"> 拖拽文件到此处或 <em>点击上传</em> </div>
      </template>
      <template v-else-if="listType === 'picture-card'">
        <el-icon>
          <Plus />
        </el-icon>
      </template>
      <template v-else>
        <el-button type="primary" :disabled="disabled">
          <el-icon class="mr-1">
            <Upload />
          </el-icon>
          {{ buttonText }}
        </el-button>
      </template>
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
  import { ref, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import { UploadFilled, Plus, Upload } from '@element-plus/icons-vue'
  import type { UploadFile, UploadInstance, UploadRawFile } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'

  interface Props {
    modelValue?: string | string[]
    multiple?: boolean
    limit?: number
    accept?: string
    maxSize?: number // MB
    listType?: 'text' | 'picture' | 'picture-card'
    autoUpload?: boolean
    showFileList?: boolean
    drag?: boolean
    disabled?: boolean
    buttonText?: string
    tip?: string
    category?: string
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
    multiple: false,
    limit: 1,
    accept: '',
    maxSize: 50,
    listType: 'text',
    autoUpload: true,
    showFileList: true,
    drag: false,
    disabled: false,
    buttonText: '选择文件',
    tip: '',
    category: ''
  })

  const emit = defineEmits<{
    'update:modelValue': [value: string | string[]]
    success: [file: any]
    error: [error: any]
    remove: [file: UploadFile]
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
    let url = `${baseUrl}/api/file/upload`
    if (props.category) {
      url += `?category=${props.category}`
    }
    return url
  })

  // 请求头
  const headers = computed(() => ({
    Authorization: `Bearer ${userStore.accessToken}`
  }))

  // 上传前校验
  const handleBeforeUpload = (file: UploadRawFile) => {
    // 检查文件大小
    const isLtSize = file.size / 1024 / 1024 < props.maxSize
    if (!isLtSize) {
      ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
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
      ElMessage.success('上传成功')
    } else {
      ElMessage.error(response.message || '上传失败')
      emit('error', response)
    }
  }

  // 上传失败
  const handleError = (error: any) => {
    ElMessage.error('上传失败')
    emit('error', error)
  }

  // 超出限制
  const handleExceed = () => {
    ElMessage.warning(`最多只能上传 ${props.limit} 个文件`)
  }

  // 移除文件
  const handleRemove = (file: UploadFile) => {
    if (props.multiple) {
      const urls = Array.isArray(props.modelValue)
        ? props.modelValue.filter((url) => url !== file.url)
        : []
      emit('update:modelValue', urls)
    } else {
      emit('update:modelValue', '')
    }
    emit('remove', file)
  }

  // 预览文件
  const handlePreview = (file: UploadFile) => {
    if (file.url && isImage(file.url)) {
      previewList.value = [file.url]
      previewIndex.value = 0
      previewVisible.value = true
    } else if (file.url) {
      window.open(file.url)
    }
  }

  // 判断是否为图片
  const isImage = (url: string) => {
    return /\.(jpg|jpeg|png|gif|webp|bmp|svg)$/i.test(url)
  }

  // 清空文件列表
  const clearFiles = () => {
    uploadRef.value?.clearFiles()
    emit('update:modelValue', props.multiple ? [] : '')
  }

  // 手动上传
  const submit = () => {
    uploadRef.value?.submit()
  }

  defineExpose({
    clearFiles,
    submit
  })
</script>

<style scoped lang="scss">
  .file-upload {
    .el-upload__tip {
      color: var(--el-text-color-secondary);
      font-size: 12px;
      margin-top: 8px;
    }
  }
</style>
