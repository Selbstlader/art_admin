<template>
  <div class="edit-page">
    <!-- 页面头部 / Page header -->
    <div class="page-header">
      <el-button @click="handleBack">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
      <h2>编辑施工图标注</h2>
    </div>

    <!-- 编辑器 / Editor -->
    <div class="editor-wrapper" v-if="annotationId">
      <AnnotationEditor
        :annotation-id="annotationId"
        :project-id="projectId"
        :cad-file-id="cadFileId"
        @saved="handleSaved"
        @exported="handleExported"
      />
    </div>

    <el-empty v-else description="未找到标注信息" />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { ArrowLeft } from '@element-plus/icons-vue'
  import AnnotationEditor from './AnnotationEditor.vue'

  /*** Route & Router ***/
  const route = useRoute()
  const router = useRouter()

  /*** Reactive State ***/
  const annotationId = ref<number>(0)
  const projectId = ref<number>(0)
  const cadFileId = ref<number>(0)

  /*** Methods ***/
  // 返回列表 / Back to list
  const handleBack = () => {
    router.push('/designer-assistant/construction-annotation')
  }

  // 保存成功 / Save success
  const handleSaved = () => {
    ElMessage.success('保存成功')
  }

  // 导出成功 / Export success
  const handleExported = (result: any) => {
    ElMessage.success('导出成功')
    // 可以在这里处理下载逻辑 / Can handle download logic here
    if (result?.downloadUrl) {
      window.open(result.downloadUrl, '_blank')
    }
  }

  /*** Lifecycle Hooks ***/
  onMounted(() => {
    // 从路由参数获取ID / Get ID from route params
    annotationId.value = Number(route.query.id) || 0
    projectId.value = Number(route.query.projectId) || 0
    cadFileId.value = Number(route.query.cadFileId) || 0

    // 如果有导出操作参数，自动打开导出对话框
    // If there's an export action param, auto open export dialog
    if (route.query.action === 'export') {
      // 编辑器组件会处理这个逻辑
      // Editor component will handle this logic
    }
  })
</script>

<style scoped lang="scss">
  .edit-page {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #f5f7fa;
  }

  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 16px;
    background: #fff;
    border-bottom: 1px solid #e4e7ed;

    h2 {
      margin: 0;
      font-size: 18px;
      font-weight: 500;
    }
  }

  .editor-wrapper {
    flex: 1;
    overflow: hidden;
  }
</style>
