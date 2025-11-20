<template>
  <ElDialog
    v-model="visible"
    :title="dialogTitle"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px" label-position="right">
      <ElFormItem label="字典类型名称" prop="typeName">
        <ElInput
          v-model="form.typeName"
          placeholder="请输入字典类型名称"
          maxlength="50"
          show-word-limit
        />
      </ElFormItem>
      <ElFormItem label="字典类型编码" prop="typeCode">
        <ElInput
          v-model="form.typeCode"
          placeholder="请输入字典类型编码"
          maxlength="50"
          show-word-limit
          :disabled="dialogType === 'edit'"
        />
        <div class="form-tip">编码用于系统内部识别，创建后不可修改</div>
      </ElFormItem>
      <ElFormItem label="描述" prop="description">
        <ElInput
          v-model="form.description"
          type="textarea"
          placeholder="请输入描述"
          maxlength="200"
          show-word-limit
          :rows="3"
        />
      </ElFormItem>
      <ElFormItem label="状态" prop="enabled">
        <ElSwitch v-model="form.enabled" active-text="启用" inactive-text="禁用" />
      </ElFormItem>
      <ElFormItem label="备注" prop="remark">
        <ElInput
          v-model="form.remark"
          type="textarea"
          placeholder="请输入备注"
          maxlength="500"
          show-word-limit
          :rows="3"
        />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElSpace>
        <ElButton @click="handleClose">取消</ElButton>
        <ElButton type="primary" :loading="loading" @click="handleSubmit">确定</ElButton>
      </ElSpace>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { fetchCreateDictionaryType, fetchUpdateDictionaryType } from '@/api/system-manage'

  interface Form {
    typeName: string
    typeCode: string
    description: string
    enabled: boolean
    remark: string
  }

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    dictionaryTypeData?: Api.SystemManage.DictionaryTypeItem
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    dictionaryTypeData: undefined
  })
  const emit = defineEmits<Emits>()

  const formRef = ref<FormInstance>()
  const loading = ref(false)

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const dialogTitle = computed(() => {
    return props.dialogType === 'add' ? '新增字典类型' : '编辑字典类型'
  })

  const form = reactive<Form>({
    typeName: '',
    typeCode: '',
    description: '',
    enabled: true,
    remark: ''
  })

  const rules: FormRules<Form> = {
    typeName: [
      { required: true, message: '请输入字典类型名称', trigger: 'blur' },
      { max: 50, message: '字典类型名称不能超过50个字符', trigger: 'blur' }
    ],
    typeCode: [
      { required: true, message: '请输入字典类型编码', trigger: 'blur' },
      { max: 50, message: '字典类型编码不能超过50个字符', trigger: 'blur' },
      {
        pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
        message: '编码必须以字母开头，只能包含字母、数字和下划线',
        trigger: 'blur'
      }
    ],
    description: [{ max: 200, message: '描述不能超过200个字符', trigger: 'blur' }],
    remark: [{ max: 500, message: '备注不能超过500个字符', trigger: 'blur' }]
  }

  // 监听弹窗显示，初始化表单数据
  watch(
    () => visible.value,
    (val) => {
      if (val) {
        initForm()
      }
    }
  )

  const initForm = () => {
    if (props.dialogType === 'edit' && props.dictionaryTypeData) {
      Object.assign(form, {
        id: props.dictionaryTypeData.id,
        typeName: props.dictionaryTypeData.typeName,
        typeCode: props.dictionaryTypeData.typeCode,
        description: props.dictionaryTypeData.description,
        enabled: props.dictionaryTypeData.enabled,
        remark: props.dictionaryTypeData.remark
      })
    } else {
      Object.assign(form, {
        typeName: '',
        typeCode: '',
        description: '',
        enabled: true,
        remark: ''
      })
    }
  }

  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      loading.value = true

      if (props.dialogType === 'add') {
        await fetchCreateDictionaryType(form)
        ElMessage.success('字典类型创建成功')
      } else {
        await fetchUpdateDictionaryType({ ...form, id: props.dictionaryTypeData!.id })
        ElMessage.success('字典类型更新成功')
      }

      emit('success')
      handleClose()
    } catch (error: any) {
      console.error('操作失败:', error)
      ElMessage.error(error.message || '操作失败')
    } finally {
      loading.value = false
    }
  }

  const handleClose = () => {
    visible.value = false
    formRef.value?.resetFields()
  }
</script>

<style scoped>
  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 4px;
  }
</style>
