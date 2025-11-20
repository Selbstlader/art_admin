<template>
  <ElDialog
    v-model="dialogVisible"
    title="操作日志详情"
    width="800px"
    :close-on-click-modal="false"
  >
    <ElDescriptions :column="2" border v-if="logData">
      <ElDescriptionsItem label="日志ID">
        {{ logData.id }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作模块">
        {{ logData.module }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="业务类型">
        <ElTag :type="getBusinessTypeColor(logData.businessType)">
          {{ logData.businessType }}
        </ElTag>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="请求方式">
        <ElTag>{{ logData.requestMethod }}</ElTag>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作人员">
        {{ logData.operatorName }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作IP">
        {{ logData.operatorIp }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作地点" :span="2">
        {{ logData.operatorAddr || '未知' }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="请求URL" :span="2">
        <span class="url-text">{{ logData.requestUrl }}</span>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作状态">
        <ElTag :type="logData.status === 1 ? 'success' : 'danger'">
          {{ logData.status === 1 ? '成功' : '失败' }}
        </ElTag>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="耗时">
        <ElTag :type="getCostTimeColor(logData.costTime)"> {{ logData.costTime }}ms </ElTag>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="操作时间" :span="2">
        {{ logData.operationTime }}
      </ElDescriptionsItem>
      <ElDescriptionsItem label="请求参数" :span="2">
        <div class="code-block">
          <pre>{{ formatJson(logData.requestParam) }}</pre>
        </div>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="响应数据" :span="2">
        <div class="code-block">
          <pre>{{ formatJson(logData.responseData) }}</pre>
        </div>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="错误信息" :span="2" v-if="logData.errorMsg">
        <ElAlert :title="logData.errorMsg" type="error" :closable="false" />
      </ElDescriptionsItem>
      <ElDescriptionsItem label="用户代理" :span="2">
        <div class="user-agent">{{ logData.userAgent }}</div>
      </ElDescriptionsItem>
    </ElDescriptions>

    <template #footer>
      <ElButton @click="dialogVisible = false">关闭</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import {
    ElDialog,
    ElDescriptions,
    ElDescriptionsItem,
    ElTag,
    ElButton,
    ElAlert
  } from 'element-plus'

  type OperationLogItem = Api.SystemManage.OperationLogItem

  defineProps<{
    logData?: OperationLogItem
  }>()

  const dialogVisible = defineModel<boolean>({ required: true })

  /**
   * 获取业务类型颜色
   */
  const getBusinessTypeColor = (
    type: string
  ): 'success' | 'warning' | 'danger' | 'primary' | 'info' => {
    const colorMap: Record<string, 'success' | 'warning' | 'danger' | 'primary' | 'info'> = {
      新增: 'primary',
      修改: 'warning',
      删除: 'danger',
      查询: 'info',
      导出: 'success',
      导入: 'success'
    }
    return colorMap[type] || 'info'
  }

  /**
   * 获取耗时颜色
   */
  const getCostTimeColor = (time: number): 'success' | 'warning' | 'danger' => {
    if (time > 1000) return 'danger'
    if (time > 500) return 'warning'
    return 'success'
  }

  /**
   * 格式化JSON
   */
  const formatJson = (str: string): string => {
    if (!str) return '无'
    try {
      const obj = JSON.parse(str)
      return JSON.stringify(obj, null, 2)
    } catch {
      return str
    }
  }
</script>

<style lang="scss" scoped>
  .url-text {
    word-break: break-all;
    color: var(--el-color-primary);
  }

  .code-block {
    max-height: 300px;
    overflow: auto;
    background-color: #f5f7fa;
    border-radius: 4px;
    padding: 12px;

    pre {
      margin: 0;
      font-family: 'Courier New', monospace;
      font-size: 12px;
      line-height: 1.5;
      white-space: pre-wrap;
      word-wrap: break-word;
    }
  }

  .user-agent {
    word-break: break-all;
    font-size: 12px;
    color: #606266;
  }
</style>
