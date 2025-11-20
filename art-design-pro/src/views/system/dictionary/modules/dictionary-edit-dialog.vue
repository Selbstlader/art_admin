<template>
  <ElDialog
    v-model="visible"
    :title="dialogTitle"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px" label-position="right">
      <ElFormItem label="字典类型" prop="typeCode">
        <ElInput
          v-model="typeDisplayText"
          placeholder="请选择字典类型"
          readonly
          style="width: 100%"
          :disabled="true"
        />
      </ElFormItem>
      <ElFormItem label="字典标签" prop="label">
        <ElInput
          v-model="form.label"
          placeholder="请输入字典标签"
          maxlength="100"
          show-word-limit
        />
        <div class="form-tip">用于前端显示的文本</div>
      </ElFormItem>
      <ElFormItem label="字典值" prop="value">
        <ElInput v-model="form.value" placeholder="请输入字典值" maxlength="100" show-word-limit />
        <div class="form-tip">用于后端存储和比较的值</div>
      </ElFormItem>
      <ElFormItem label="排序号" prop="orderNum">
        <ElInputNumber
          v-model="form.orderNum"
          :min="0"
          :max="9999"
          placeholder="请输入排序号"
          style="width: 100%"
        />
        <div class="form-tip">数值越小排序越靠前</div>
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
        <ElButton type="primary" :loading="loading" @click="handleSubmit"> 确定 </ElButton>
      </ElSpace>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { fetchCreateDictionary, fetchUpdateDictionary } from '@/api/system-manage'

  interface Form {
    typeCode: string
    label: string
    value: string
    orderNum: number
    enabled: boolean
    remark: string
  }

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    dictionaryData?: Api.SystemManage.DictionaryItem
    dictionaryTypeOptions: Api.SystemManage.DictionaryTypeItem[]
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    dictionaryData: undefined
  })
  const emit = defineEmits<Emits>()

  const formRef = ref<FormInstance>()
  const loading = ref(false)

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const dialogTitle = computed(() => {
    return props.dialogType === 'add' ? '新增字典数据' : '编辑字典数据'
  })

  // 字典类型显示文本
  const typeDisplayText = computed(() => {
    if (form.typeCode) {
      const type = props.dictionaryTypeOptions.find((item) => item.typeCode === form.typeCode)
      return type ? `${type.typeName} (${type.typeCode})` : form.typeCode
    }
    return ''
  })

  const form = reactive<Form>({
    typeCode: '',
    label: '',
    value: '',
    orderNum: 0,
    enabled: true,
    remark: ''
  })

  const rules: FormRules<Form> = {
    label: [
      { required: true, message: '请输入字典标签', trigger: 'blur' },
      { max: 100, message: '字典标签不能超过100个字符', trigger: 'blur' }
    ],
    value: [
      { required: true, message: '请输入字典值', trigger: 'blur' },
      { max: 100, message: '字典值不能超过100个字符', trigger: 'blur' }
    ],
    orderNum: [{ type: 'number', message: '排序号必须为数字', trigger: 'blur' }],
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
    if (props.dialogType === 'edit' && props.dictionaryData) {
      Object.assign(form, {
        id: props.dictionaryData.id,
        typeCode: props.dictionaryData.typeCode,
        label: props.dictionaryData.label,
        value: props.dictionaryData.value,
        orderNum: props.dictionaryData.orderNum,
        enabled: props.dictionaryData.enabled,
        remark: props.dictionaryData.remark
      })
    } else {
      // 新增时，如果有传入的字典类型数据，使用它
      const defaultTypeCode = props.dictionaryData?.typeCode || ''
      Object.assign(form, {
        typeCode: defaultTypeCode,
        label: '',
        value: '',
        orderNum: 0,
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
        await fetchCreateDictionary(form)
        ElMessage.success('字典数据创建成功')
      } else {
        await fetchUpdateDictionary({ ...form, id: props.dictionaryData!.id })
        ElMessage.success('字典数据更新成功')
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
