<template>
  <ArtSearchBar
    ref="searchBarRef"
    v-model="formData"
    :items="formItems"
    :rules="rules"
    @reset="handleReset"
    @search="handleSearch"
  >
  </ArtSearchBar>
</template>

<script setup lang="ts">
  interface Props {
    modelValue: Record<string, any>
  }
  interface Emits {
    (e: 'update:modelValue', value: Record<string, any>): void
    (e: 'search', params: Record<string, any>): void
    (e: 'reset'): void
  }
  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 表单数据双向绑定
  const searchBarRef = ref()
  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  // 校验规则
  const rules = {}

  // 表单配置
  const formItems = computed(() => [
    {
      label: '用户账号',
      key: 'userName',
      type: 'input',
      placeholder: '请输入用户账号',
      clearable: true
    },
    {
      label: '用户昵称',
      key: 'nickName',
      type: 'input',
      placeholder: '请输入用户昵称',
      clearable: true
    },
    {
      label: '手机号',
      key: 'phone',
      type: 'input',
      props: { placeholder: '请输入手机号', maxlength: '11' }
    },
    {
      label: '用户类型',
      key: 'userType',
      type: 'select',
      props: {
        placeholder: '请选择用户类型',
        clearable: true,
        options: [
          { label: '管理员', value: '1' },
          { label: '普通用户', value: '2' }
        ]
      }
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: {
        placeholder: '请选择状态',
        clearable: true,
        options: [
          { label: '正常', value: '1' },
          { label: '禁用', value: '2' }
        ]
      }
    }
  ])

  // 事件
  function handleReset() {
    emit('reset')
  }

  async function handleSearch() {
    await searchBarRef.value.validate()
    emit('search', formData.value)
  }
</script>
