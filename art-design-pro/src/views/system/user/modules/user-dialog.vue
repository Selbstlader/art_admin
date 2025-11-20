<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '添加用户' : '编辑用户'"
    width="30%"
    align-center
  >
    <ElForm ref="formRef" :model="formData" :rules="rules" label-width="80px">
      <ElFormItem label="用户名" prop="userName">
        <ElInput
          v-model="formData.userName"
          placeholder="请输入用户名"
          :disabled="dialogType === 'edit'"
        />
      </ElFormItem>
      <ElFormItem label="昵称" prop="nickName">
        <ElInput v-model="formData.nickName" placeholder="请输入昵称" />
      </ElFormItem>
      <ElFormItem v-if="dialogType === 'add'" label="密码" prop="password">
        <ElInput
          v-model="formData.password"
          type="password"
          placeholder="请输入密码"
          show-password
        />
      </ElFormItem>
      <ElFormItem label="邮箱" prop="email">
        <ElInput v-model="formData.email" placeholder="请输入邮箱" />
      </ElFormItem>
      <ElFormItem label="手机号" prop="userPhone">
        <ElInput v-model="formData.userPhone" placeholder="请输入手机号" />
      </ElFormItem>
      <ElFormItem label="性别" prop="userGender">
        <ElSelect v-model="formData.userGender">
          <ElOption label="男" value="male" />
          <ElOption label="女" value="female" />
          <ElOption label="未知" value="unknown" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="状态" prop="status">
        <ElSelect v-model="formData.status">
          <ElOption label="启用" value="1" />
          <ElOption label="禁用" value="0" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="角色" prop="roleIds">
        <ElSelect v-model="formData.roleIds" multiple>
          <ElOption
            v-for="role in roleList"
            :key="role.roleCode"
            :value="role.roleCode"
            :label="role.roleName"
          />
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
  import { ROLE_LIST_DATA } from '@/mock/temp/formData'
  import { fetchCreateUser, fetchUpdateUser } from '@/api/system-manage'
  import type { FormInstance, FormRules } from 'element-plus'

  interface Props {
    visible: boolean
    type: string
    userData?: Partial<Api.SystemManage.UserListItem>
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 角色列表数据
  const roleList = ref(ROLE_LIST_DATA)

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
    email: '',
    userPhone: '',
    userGender: 'male',
    status: '1',
    roleIds: [] as string[] // 前端使用roleCode数组
  })

  // 表单验证规则
  const rules: FormRules = {
    userName: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
    ],
    nickName: [
      { required: true, message: '请输入昵称', trigger: 'blur' },
      { max: 50, message: '长度不能超过 50 个字符', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
    ],
    email: [
      { required: true, message: '请输入邮箱', trigger: 'blur' },
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ],
    userPhone: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
    ],
    userGender: [{ required: true, message: '请选择性别', trigger: 'blur' }],
    status: [{ required: true, message: '请选择状态', trigger: 'blur' }],
    roleIds: [
      { required: true, message: '请选择角色', trigger: 'blur' },
      { type: 'array', min: 1, message: '至少选择一个角色', trigger: 'blur' }
    ]
  }

  /**
   * 初始化表单数据
   * 根据对话框类型（新增/编辑）填充表单
   */
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.userData
    const row = props.userData

    Object.assign(formData, {
      userName: isEdit && row ? row.userName || '' : '',
      nickName: isEdit && row ? row.nickName || '' : '',
      password: '', // 编辑时不显示密码
      email: isEdit && row ? row.userEmail || '' : '',
      userPhone: isEdit && row ? row.userPhone || '' : '',
      userGender: isEdit && row ? row.userGender || 'male' : 'male',
      status: isEdit && row ? row.status || '1' : '1',
      roleIds:
        isEdit && row
          ? Array.isArray(row.userRoles)
            ? row.userRoles.map((role: any) => {
                // 如果是字符串数组，需要转换为角色ID数组
                const roleItem = roleList.value.find(
                  (r) => r.roleName === role || r.roleCode === role
                )
                return roleItem ? roleItem.roleCode : role
              })
            : []
          : []
    })
  }

  /**
   * 监听对话框状态变化
   * 当对话框打开时初始化表单数据并清除验证状态
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
   * 验证通过后调用相应的API
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          if (dialogType.value === 'add') {
            // 创建用户
            const createData: Api.SystemManage.CreateUserRequest = {
              userName: formData.userName,
              nickName: formData.nickName,
              password: formData.password,
              email: formData.email,
              userPhone: formData.userPhone,
              userGender: formData.userGender,
              status: formData.status,
              roleIds: formData.roleIds.map((roleCode: string) => {
                // 将roleCode转换为roleId（这里需要根据实际情况调整）
                const roleItem = roleList.value.find((r) => r.roleCode === roleCode)
                return roleItem ? (roleItem as any).roleId || 1 : 1 // 临时处理，应该从后端获取正确的roleId
              })
            }
            await fetchCreateUser(createData)
            ElMessage.success('添加用户成功')
          } else {
            // 更新用户
            const updateData: Api.SystemManage.UpdateUserRequest = {
              id: props.userData?.id || 0,
              nickName: formData.nickName,
              email: formData.email,
              userPhone: formData.userPhone,
              userGender: formData.userGender,
              status: formData.status,
              roleIds: formData.roleIds.map((roleCode: string) => {
                // 将roleCode转换为roleId
                const roleItem = roleList.value.find((r) => r.roleCode === roleCode)
                return roleItem ? (roleItem as any).roleId || 1 : 1 // 临时处理
              })
            }
            await fetchUpdateUser(updateData)
            ElMessage.success('更新用户成功')
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
