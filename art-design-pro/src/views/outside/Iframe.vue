<template>
  <div class="iframe-container">
    <iframe :src="iframeSrc" frameborder="0" class="iframe-content"></iframe>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const iframeSrc = computed(() => {
  const path = route.params.path as string
  if (path) {
    return decodeURIComponent(path)
  }
  return ''
})
</script>

<style scoped lang="scss">
.iframe-container {
  width: 100%;
  height: 100%;
  overflow: hidden;

  .iframe-content {
    width: 100%;
    height: 100%;
    border: none;
  }
}
</style>
