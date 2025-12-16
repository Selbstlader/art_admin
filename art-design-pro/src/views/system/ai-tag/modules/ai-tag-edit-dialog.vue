<template>
  <ElDialog
    v-model="visible"
    :title="dialogTitle"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px" label-position="right">
      <ElFormItem label="标签名称" prop="name">
        <ElInput v-model="form.name" placeholder="请输入标签名称" maxlength="100" show-word-limit />
      </ElFormItem>

      <ElFormItem label="描述" prop="description">
        <ElInput
          v-model="form.description"
          type="textarea"
          placeholder="请输入标签描述"
          maxlength="500"
          show-word-limit
          :rows="2"
        />
      </ElFormItem>

      <ElFormItem label="知识库" prop="knowledge_base_id">
        <ElSelect
          v-model="form.knowledge_base_id"
          placeholder="请选择知识库"
          clearable
          filterable
          :loading="knowledgeBaseLoading"
          style="width: 100%"
          @change="handleKnowledgeBaseChange"
        >
          <ElOption
            v-for="item in knowledgeBaseList"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </ElSelect>
        <div class="form-tip">从 Dify 知识库列表中选择</div>
      </ElFormItem>

      <ElFormItem label="Chat API Key" prop="chat_api_key">
        <ElInput
          v-model="form.chat_api_key"
          :type="showApiKey ? 'text' : 'password'"
          placeholder="请输入 Dify Chat App API Key"
          maxlength="200"
        >
          <template #suffix>
            <ElIcon class="cursor-pointer" @click="showApiKey = !showApiKey">
              <!-- <Icon
                :icon="
                  showApiKey ? 'material-symbols:visibility' : 'material-symbols:visibility-off'
                "
              /> -->
            </ElIcon>
          </template>
        </ElInput>
        <div class="form-tip">Dify Chat App 的 API Key，用于调用对话接口</div>
      </ElFormItem>

      <ElFormItem label="系统提示词" prop="system_prompt">
        <div class="prompt-container">
          <ElInput
            v-model="form.system_prompt"
            type="textarea"
            placeholder="请输入系统提示词，定义 AI 助手的角色和行为"
            maxlength="5000"
            show-word-limit
            :rows="6"
          />
          <div class="ai-generate-section">
            <ElInput
              v-model="aiKeywords"
              placeholder="输入关键词，如：客服助手、技术支持、数据分析..."
              size="small"
              class="keywords-input"
            >
              <template #prepend>AI 关键词</template>
            </ElInput>
            <ElSpace>
              <ElButton
                type="primary"
                size="small"
                :loading="aiGenerating"
                :disabled="!aiKeywords.trim()"
                @click="handleAIGenerate('generate')"
              >
                <ElIcon class="mr-1"><MagicStick /></ElIcon>
                生成提示词
              </ElButton>
              <ElButton
                size="small"
                :loading="aiGenerating"
                :disabled="!aiKeywords.trim() || !form.system_prompt.trim()"
                @click="handleAIGenerate('polish')"
              >
                <ElIcon class="mr-1"><EditPen /></ElIcon>
                润色优化
              </ElButton>
            </ElSpace>
          </div>
        </div>
      </ElFormItem>

      <ElFormItem v-if="dialogType === 'edit'" label="状态" prop="status">
        <ElSwitch v-model="statusEnabled" active-text="启用" inactive-text="禁用" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElSpace>
        <ElButton @click="handleClose">取消</ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit"> 确定 </ElButton>
      </ElSpace>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  /*** AI Tag Edit Dialog Component ***/
  /*** Requirements: 1.1, 1.3, 3.1 - Create and edit AI tags with validation ***/

  import { MagicStick, EditPen } from '@element-plus/icons-vue'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import {
    aiTagApi,
    type AITag,
    type CreateAITagRequest,
    type UpdateAITagRequest
  } from '@/api/ai-tag'
  import { difyDatasetApi } from '@/api/dify'
  import { aiGenerateApi } from '@/api/ai-generate'

  interface Form {
    name: string
    description: string
    knowledge_base_id: string
    knowledge_base_name: string
    system_prompt: string
    chat_api_key: string
    status: number
  }

  interface KnowledgeBase {
    id: string
    name: string
  }

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    tagData?: AITag
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    tagData: undefined
  })
  const emit = defineEmits<Emits>()

  /*** State ***/
  const formRef = ref<FormInstance>()
  const submitLoading = ref(false)
  const showApiKey = ref(false)
  const knowledgeBaseLoading = ref(false)
  const knowledgeBaseList = ref<KnowledgeBase[]>([])
  const aiKeywords = ref('')
  const aiGenerating = ref(false)

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const dialogTitle = computed(() => {
    return props.dialogType === 'add' ? '新增 AI 标签' : '编辑 AI 标签'
  })

  const form = reactive<Form>({
    name: '',
    description: '',
    knowledge_base_id: '',
    knowledge_base_name: '',
    system_prompt: '',
    chat_api_key: '',
    status: 1
  })

  const statusEnabled = computed({
    get: () => form.status === 1,
    set: (value) => {
      form.status = value ? 1 : 0
    }
  })

  /*** Validation Rules ***/
  const rules: FormRules<Form> = {
    name: [
      { required: true, message: '请输入标签名称', trigger: 'blur' },
      { max: 100, message: '标签名称不能超过100个字符', trigger: 'blur' },
      {
        validator: (_rule, value, callback) => {
          if (value && value.trim().length === 0) {
            callback(new Error('标签名称不能为空白字符'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ],
    description: [{ max: 500, message: '描述不能超过500个字符', trigger: 'blur' }],
    system_prompt: [
      { required: true, message: '请输入系统提示词', trigger: 'blur' },
      {
        validator: (_rule, value, callback) => {
          if (value && value.trim().length === 0) {
            callback(new Error('系统提示词不能为空白字符'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ]
  }

  /*** Methods ***/
  const loadKnowledgeBaseList = async () => {
    try {
      knowledgeBaseLoading.value = true
      const response = await difyDatasetApi.getDatasetList({ page: 1, limit: 100 })
      knowledgeBaseList.value = (response as any)?.data || response || []
    } catch (error) {
      console.error('获取知识库列表失败:', error)
      knowledgeBaseList.value = []
    } finally {
      knowledgeBaseLoading.value = false
    }
  }

  const handleKnowledgeBaseChange = (id: string) => {
    const selected = knowledgeBaseList.value.find((item) => item.id === id)
    form.knowledge_base_name = selected?.name || ''
  }

  const initForm = () => {
    if (props.dialogType === 'edit' && props.tagData) {
      Object.assign(form, {
        name: props.tagData.name,
        description: props.tagData.description || '',
        knowledge_base_id: props.tagData.knowledge_base_id || '',
        knowledge_base_name: props.tagData.knowledge_base_name || '',
        system_prompt: props.tagData.system_prompt,
        chat_api_key: props.tagData.chat_api_key || '',
        status: props.tagData.status
      })
    } else {
      Object.assign(form, {
        name: '',
        description: '',
        knowledge_base_id: '',
        knowledge_base_name: '',
        system_prompt: '',
        chat_api_key: '',
        status: 1
      })
    }
    showApiKey.value = false
  }

  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitLoading.value = true

      if (props.dialogType === 'add') {
        const createData: CreateAITagRequest = {
          name: form.name.trim(),
          description: form.description?.trim(),
          knowledge_base_id: form.knowledge_base_id,
          knowledge_base_name: form.knowledge_base_name,
          system_prompt: form.system_prompt.trim(),
          chat_api_key: form.chat_api_key
        }
        await aiTagApi.create(createData)
        ElMessage.success('标签创建成功')
      } else {
        const updateData: UpdateAITagRequest = {
          name: form.name.trim(),
          description: form.description?.trim(),
          knowledge_base_id: form.knowledge_base_id,
          knowledge_base_name: form.knowledge_base_name,
          system_prompt: form.system_prompt.trim(),
          chat_api_key: form.chat_api_key,
          status: form.status
        }
        await aiTagApi.update(props.tagData!.id, updateData)
        ElMessage.success('标签更新成功')
      }

      emit('success')
      handleClose()
    } catch (error: any) {
      console.error('操作失败:', error)
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  }

  const handleClose = () => {
    visible.value = false
    formRef.value?.resetFields()
    aiKeywords.value = ''
  }

  /*** AI Generate Handler ***/
  const handleAIGenerate = async (mode: 'generate' | 'polish') => {
    if (!aiKeywords.value.trim()) {
      ElMessage.warning('请输入关键词')
      return
    }

    if (mode === 'polish' && !form.system_prompt.trim()) {
      ElMessage.warning('润色模式需要先输入现有提示词')
      return
    }

    try {
      aiGenerating.value = true
      const response = await aiGenerateApi.generatePrompt({
        keywords: aiKeywords.value.trim(),
        current_prompt: mode === 'polish' ? form.system_prompt : undefined,
        mode
      })

      if (response && response.prompt) {
        form.system_prompt = response.prompt
        ElMessage.success(mode === 'generate' ? '提示词生成成功' : '提示词润色成功')
      }
    } catch (error: any) {
      console.error('AI生成失败:', error)
      ElMessage.error(error.message || 'AI生成失败，请稍后重试')
    } finally {
      aiGenerating.value = false
    }
  }

  /*** Watchers ***/
  watch(
    () => visible.value,
    (val) => {
      if (val) {
        initForm()
        loadKnowledgeBaseList()
      }
    }
  )
</script>

<style scoped>
  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 4px;
  }

  .cursor-pointer {
    cursor: pointer;
  }

  .prompt-container {
    width: 100%;
  }

  .ai-generate-section {
    margin-top: 12px;
    padding: 12px;
    background: var(--el-fill-color-light);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .keywords-input {
    width: 100%;
  }

  .mr-1 {
    margin-right: 4px;
  }
</style>
