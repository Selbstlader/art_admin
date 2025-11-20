<template>
  <ArtSearchBar
    v-model="formFilters"
    :items="formItems"
    @reset="handleReset"
    @search="handleSearch"
  />
</template>

<script setup lang="ts">
  const emit = defineEmits<{
    search: [params: Record<string, any>]
    reset: []
  }>()

  const model = defineModel<Record<string, any>>({ required: true })

  const formFilters = computed({
    get: () => model.value,
    set: (val) => {
      model.value = val
    }
  })

  const formItems = computed(() => [
    {
      label: '操作模块',
      key: 'module',
      type: 'input',
      props: { clearable: true, placeholder: '请输入操作模块' }
    },
    {
      label: '业务类型',
      key: 'businessType',
      type: 'select',
      props: {
        clearable: true,
        placeholder: '请选择业务类型',
        options: [
          { label: '新增', value: '新增' },
          { label: '修改', value: '修改' },
          { label: '删除', value: '删除' },
          { label: '查询', value: '查询' },
          { label: '导出', value: '导出' },
          { label: '导入', value: '导入' }
        ]
      }
    },
    {
      label: '操作人员',
      key: 'operatorName',
      type: 'input',
      props: { clearable: true, placeholder: '请输入操作人员' }
    },
    {
      label: '操作状态',
      key: 'status',
      type: 'select',
      props: {
        clearable: true,
        placeholder: '请选择操作状态',
        options: [
          { label: '成功', value: 1 },
          { label: '失败', value: 0 }
        ]
      }
    },
    {
      label: '操作时间',
      key: 'daterange',
      type: 'date-picker',
      props: {
        type: 'datetimerange',
        clearable: true,
        startPlaceholder: '开始时间',
        endPlaceholder: '结束时间',
        valueFormat: 'YYYY-MM-DD HH:mm:ss'
      }
    }
  ])

  const handleSearch = () => {
    const { daterange, ...otherParams } = formFilters.value
    const [startTime, endTime] = Array.isArray(daterange) ? daterange : [null, null]

    emit('search', {
      ...otherParams,
      startTime,
      endTime
    })
  }

  const handleReset = () => {
    emit('reset')
  }
</script>
