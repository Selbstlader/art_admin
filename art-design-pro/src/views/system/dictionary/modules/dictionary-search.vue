<template>
  <ElForm ref="formRef" :model="form" inline>
    <ElFormItem label="字典类型编码" prop="typeCode">
      <ElSelect
        v-model="form.typeCode"
        placeholder="请选择字典类型"
        clearable
        filterable
        style="width: 200px"
      >
        <ElOption
          v-for="item in dictionaryTypeOptions"
          :key="item.typeCode"
          :label="`${item.typeName} (${item.typeCode})`"
          :value="item.typeCode"
        />
      </ElSelect>
    </ElFormItem>
    <ElFormItem label="字典标签" prop="label">
      <ElInput
        v-model="form.label"
        placeholder="请输入字典标签"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
    </ElFormItem>
    <ElFormItem label="字典值" prop="value">
      <ElInput
        v-model="form.value"
        placeholder="请输入字典值"
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      />
    </ElFormItem>
    <ElFormItem label="状态" prop="enabled">
      <ElSelect
        v-model="form.enabled"
        placeholder="请选择状态"
        clearable
        style="width: 120px"
      >
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
    typeCode?: string
    label?: string
    value?: string
    enabled?: boolean
  }

  interface Props {
    modelValue: Form
    dictionaryTypeOptions: Api.SystemManage.DictionaryTypeItem[]
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
