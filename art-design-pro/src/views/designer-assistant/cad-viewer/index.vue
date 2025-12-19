<template>
  <div class="cad-viewer-page">
    <CadUpload />
  </div>
</template>

<script setup lang="ts">
/***
 * CAD Viewer Module Entry Page
 * CAD图纸预览模块入口页面
 * Requirements: 8.1, 8.2, 8.3, 8.5
 ***/
import CadUpload from './CadUpload.vue'
</script>

<style scoped lang="scss">
.cad-viewer-page {
  height: 100%;
}
</style>
