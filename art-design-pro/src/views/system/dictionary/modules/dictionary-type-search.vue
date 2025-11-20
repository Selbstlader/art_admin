<template>
  <ElForm ref="formRef" :model="form" inline>
    <ElFormItem label="字典类型名称" prop="typeName">
      <ElInput
        v-model="form.typeName"
        placeholder="请输入字典类型名称"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
    </ElFormItem>
    <ElFormItem label="字典类型编码" prop="typeCode">
      <ElInput
        v-model="form.typeCode"
        placeholder="请输入字典类型编码"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
    </ElFormItem>
    <ElFormItem label="描述" prop="description">
      <ElInput
        v-model="form.description"
        placeholder="请输入描述"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
    </ElFormItem>
    <ElFormItem label="状态" prop="enabled">
      <ElSelect v-model="form.enabled" placeholder="请选择状态" clearable style="width: 120px">
        <ElOption label="启用" :value="true" />
        <ElOption label="禁用" :value="false" />
      </ElSelect>
    </ElFormItem>
    <ElFormItem>
      <ElSpace>
        <ElButton type="primary" @click="handleSearch">搜索</ElButton>
        <ElButton @click="handleReset">重置</ElButton>
      </ElSpace>
    </ElFormItem>
  </ElForm>
</template>

<script setup lang="ts">
  interface Form {
    typeName?: string
    typeCode?: string
    description?: string
    enabled?: boolean
  }

  interface Props {
    modelValue: Form
  }

  interface Emits {
    (e: 'update:modelValue', value: Form): void
    (e: 'search'): void
    (e: 'reset'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const form = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const handleSearch = () => {
    emit('search')
  }

  const handleReset = () => {
    emit('reset')
  }
</script>
