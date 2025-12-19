<template>
  <div class="cad-generation-container">
    <router-view v-if="hasChildRoute" />
    <CadGenerationList v-else />
  </div>
</template>

<script setup lang="ts">
  /***
   * CAD Generation Index
   * AI生成CAD模块入口
   ***/
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import CadGenerationList from './CadGenerationList.vue'

  const route = useRoute()

  const hasChildRoute = computed(() => {
    return route.matched.length > 2
  })
</script>

<style scoped lang="scss">
  .cad-generation-container {
    height: 100%;
  }
</style>
