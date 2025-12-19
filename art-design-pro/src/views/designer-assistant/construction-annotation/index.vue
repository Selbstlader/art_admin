<template>
  <div class="construction-annotation-container">
    <router-view v-if="hasChildRoute" />
    <AnnotationList v-else />
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import AnnotationList from './AnnotationList.vue'

  /*** Route ***/
  const route = useRoute()

  /*** Computed ***/
  const hasChildRoute = computed(() => {
    return route.matched.length > 2
  })
</script>

<style scoped lang="scss">
  .construction-annotation-container {
    height: 100%;
  }
</style>
