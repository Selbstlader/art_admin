<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '添加APP用户' : '编辑APP用户'"
    width="30%"
    align-center
  >
    <ElForm ref="formRef" :model="formData" :rules="rules" label-width="90px">
      <ElFormItem label="用户账号" prop="userName">
        <ElInput
          v-model="formData.userName"
          placeholder="请输入用户账号"
          :disabled="dialogType === 'edit'"
        />
      </ElFormItem>
      <ElFormItem label="用户昵称" prop="nickName">
        <ElInput v-model="formData.nickName" placeholder="请输入用户昵称" />
      </ElFormItem>
      <ElFormItem v-if="dialogType === 'add'" label="密码" prop="password">
        <ElInput
          v-model="formData.password"
          type="password"
          placeholder="请输入密码"
          show-password
        />
      </ElFormItem>
      <ElFormItem label="手机号" prop="phone">
        <ElInput v-model="formData.phone" placeholder="请输入手机号" />
      </ElFormItem>
      <ElFormItem label="邮箱" prop="email">
        <ElInput v-model="formData.email" placeholder="请输入邮箱" />
      </ElFormItem>
      <ElFormItem label="用户类型" prop="userType">
        <ElSelect v-model="formData.userType" placeholder="请选择用户类型">
          <ElOption label="管理员" value="1" />
          <ElOption label="普通用户" value="2" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="状态" prop="status">
        <ElSelect v-model="formData.status" placeholder="请选择状态">
          <ElOption label="正常" value="1" />
          <ElOption label="禁用" value="2" />
        </ElSelect>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">提交</ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { fetchCreateAppUser, fetchUpdateAppUser } from '@/api/system-manage'
  import type { FormInstance, FormRules } from 'element-plus'

  interface Props {
    visible: boolean
    type: string
    userData?: Partial<Api.SystemManage.AppUserListItem>
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 对话框显示控制
  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const dialogType = computed(() => props.type)

  // 表单实例
  const formRef = ref<FormInstance>()

  // 表单数据
  const formData = reactive({
    userName: '',
    nickName: '',
    password: '',
    phone: '',
    email: '',
    userType: '2' as Api.SystemManage.AppUserType,
    status: '1'
  })

  // 表单验证规则
  const rules: FormRules = {
    userName: [
      { required: true, message: '请输入用户账号', trigger: 'blur' },
      { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
    ],
    nickName: [
      { required: true, message: '请输入用户昵称', trigger: 'blur' },
      { max: 50, message: '长度不能超过 50 个字符', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
    ],
    phone: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
    ],
    email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
    userType: [{ required: true, message: '请选择用户类型', trigger: 'blur' }],
    status: [{ required: true, message: '请选择状态', trigger: 'blur' }]
  }

  /**
   * 初始化表单数据
   */
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.userData
    const row = props.userData

    Object.assign(formData, {
      userName: isEdit && row ? row.userName || '' : '',
      nickName: isEdit && row ? row.nickName || '' : '',
      password: '',
      phone: isEdit && row ? row.phone || '' : '',
      email: isEdit && row ? row.email || '' : '',
      userType: isEdit && row ? row.userType || '2' : '2',
      status: isEdit && row ? row.status || '1' : '1'
    })
  }

  /**
   * 监听对话框状态变化
   */
  watch(
    () => [props.visible, props.type, props.userData],
    ([visible]) => {
      if (visible) {
        initFormData()
        nextTick(() => {
          formRef.value?.clearValidate()
        })
      }
    },
    { immediate: true }
  )

  /**
   * 提交表单
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          if (dialogType.value === 'add') {
            // 创建APP用户
            const createData: Api.SystemManage.CreateAppUserRequest = {
              userName: formData.userName,
              nickName: formData.nickName,
              password: formData.password,
              phone: formData.phone,
              email: formData.email,
              userType: formData.userType,
              status: formData.status
            }
            await fetchCreateAppUser(createData)
            ElMessage.success('添加APP用户成功')
          } else {
            // 更新APP用户
            const updateData: Api.SystemManage.UpdateAppUserRequest = {
              id: props.userData?.id || 0,
              nickName: formData.nickName,
              phone: formData.phone,
              email: formData.email,
              userType: formData.userType,
              status: formData.status
            }
            await fetchUpdateAppUser(updateData)
            ElMessage.success('更新APP用户成功')
          }

          dialogVisible.value = false
          emit('submit')
        } catch (error) {
          console.error('操作失败:', error)
          ElMessage.error('操作失败，请重试')
        }
      }
    })
  }
</script>
