<template>
  <ElDialog
    v-model="visible"
    :title="dialogType === 'add' ? '新增部门' : '编辑部门'"
    width="30%"
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px">
      <ElFormItem label="上级部门" prop="parentId">
        <ElTreeSelect
          v-model="form.parentId"
          :data="departmentTreeData"
          :props="{ label: 'deptName', value: 'deptId', children: 'children' }"
          placeholder="请选择上级部门"
          clearable
          check-strictly
          :render-after-expand="false"
        />
      </ElFormItem>
      <ElFormItem label="部门名称" prop="deptName">
        <ElInput v-model="form.deptName" placeholder="请输入部门名称" />
      </ElFormItem>
      <ElFormItem label="部门编码" prop="deptCode">
        <ElInput v-model="form.deptCode" placeholder="请输入部门编码" />
      </ElFormItem>
      <ElFormItem label="显示排序" prop="orderNum">
        <ElInputNumber v-model="form.orderNum" :min="0" :max="999" placeholder="请输入排序号" />
      </ElFormItem>
      <ElFormItem label="负责人">
        <ElSelect
          v-model="form.leader"
          placeholder="请选择负责人"
          filterable
          clearable
          @change="handleUserChange"
        >
          <ElOption
            v-for="user in userList"
            :key="user.id"
            :label="user.userName"
            :value="user.userName"
          >
            <span>{{ user.userName }}</span>
            <span style="float: right; color: #8492a6; font-size: 13px">{{ user.userPhone }}</span>
          </ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="联系电话">
        <ElInput v-model="form.phone" placeholder="选择负责人后自动带出" readonly />
      </ElFormItem>
      <ElFormItem label="邮箱">
        <ElInput v-model="form.email" placeholder="选择负责人后自动带出" readonly />
      </ElFormItem>
      <ElFormItem label="部门状态">
        <ElRadioGroup v-model="form.status">
          <ElRadio :label="1">正常</ElRadio>
          <ElRadio :label="0">停用</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="handleClose">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">提交</ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import type { FormInstance, FormRules } from 'element-plus'
  import { ElMessage } from 'element-plus'
  import {
    fetchCreateDepartment,
    fetchUpdateDepartment,
    fetchGetDepartmentList,
    fetchGetUserList
  } from '@/api/system-manage'

  type DepartmentListItem = Api.SystemManage.DepartmentListItem
  type UserListItem = Api.SystemManage.UserListItem

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    departmentData?: DepartmentListItem
    departmentOptions?: DepartmentListItem[]
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const formRef = ref<FormInstance>()
  const departmentTreeData = ref<DepartmentListItem[]>([])

  const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  // 表单数据
  const form = ref({
    parentId: undefined as number | undefined,
    deptName: '',
    deptCode: '',
    orderNum: 0,
    leader: '',
    phone: '',
    email: '',
    status: 1
  })

  // 表单校验规则
  const rules: FormRules = {
    deptName: [
      { required: true, message: '请输入部门名称', trigger: 'blur' },
      { min: 2, max: 50, message: '部门名称长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    deptCode: [
      { required: true, message: '请输入部门编码', trigger: 'blur' },
      { min: 2, max: 50, message: '部门编码长度在 2 到 50 个字符', trigger: 'blur' },
      { pattern: /^[A-Za-z0-9_]+$/, message: '部门编码只能包含字母、数字和下划线', trigger: 'blur' }
    ],
    orderNum: [
      { required: true, message: '请输入显示排序', trigger: 'blur' },
      { type: 'number', min: 0, max: 999, message: '排序号范围在 0 到 999', trigger: 'blur' }
    ],
    phone: [],
    email: []
  }

  // 用户列表
  const userList = ref<UserListItem[]>([])

  // 获取用户列表
  const getUserList = async () => {
    try {
      const response = await fetchGetUserList({ current: 1, size: 1000 })
      userList.value = response.records || []
    } catch (error) {
      console.error('获取用户列表失败:', error)
    }
  }

  // 用户选择变化处理
  const handleUserChange = (userName: string) => {
    const selectedUser = userList.value.find((user) => user.userName === userName)
    if (selectedUser) {
      form.value.leader = selectedUser.userName
      form.value.phone = selectedUser.userPhone
      form.value.email = selectedUser.userEmail
    } else {
      form.value.leader = ''
      form.value.phone = ''
      form.value.email = ''
    }
  }

  // 获取部门树形数据
  const getDepartmentTreeData = async () => {
    try {
      const response = await fetchGetDepartmentList()
      // 添加根部门选项
      departmentTreeData.value = [
        {
          deptId: 0,
          deptName: '根部门',
          deptCode: 'ROOT',
          orderNum: 0,
          status: 1,
          createTime: new Date().toISOString(),
          children: response || []
        }
      ]
    } catch (error) {
      console.error('获取部门树数据失败:', error)
      departmentTreeData.value = []
    }
  }

  // 重置表单
  const resetForm = () => {
    form.value = {
      parentId: undefined,
      deptName: '',
      deptCode: '',
      orderNum: 0,
      leader: '',
      phone: '',
      email: '',
      status: 1
    }
  }

  // 设置表单数据（编辑时）
  const setFormData = (data?: DepartmentListItem) => {
    if (data) {
      form.value = {
        parentId: data.parentId || 0,
        deptName: data.deptName || '',
        deptCode: data.deptCode || '',
        orderNum: data.orderNum || 0,
        leader: data.leader || '',
        phone: data.phone || '',
        email: data.email || '',
        status: data.status ?? 1
      }
    } else {
      resetForm()
    }
  }

  // 关闭弹窗
  const handleClose = () => {
    visible.value = false
    resetForm()
  }

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()

      const submitData = {
        ...form.value,
        parentId: form.value.parentId === 0 ? null : form.value.parentId
      }

      if (props.dialogType === 'add') {
        await fetchCreateDepartment(submitData)
        ElMessage.success('新增部门成功')
      } else {
        await fetchUpdateDepartment({
          deptId: props.departmentData!.deptId,
          ...submitData
        })
        ElMessage.success('更新部门成功')
      }

      emit('success')
      handleClose()
    } catch (error: any) {
      console.error('提交失败:', error)
      ElMessage.error(error.message || '提交失败，请重试')
    }
  }

  // 监听弹窗显示状态
  watch(visible, (val) => {
    if (val) {
      getDepartmentTreeData()
      getUserList()
      setFormData(props.departmentData)
    }
  })
</script>
