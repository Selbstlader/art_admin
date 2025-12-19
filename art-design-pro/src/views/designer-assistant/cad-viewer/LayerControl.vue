<template>
  <div class="layer-control">
    <div class="layer-header">
      <h4>图层控制</h4>
      <span class="layer-count">共 {{ layers.length }} 个图层</span>
    </div>

    <!-- 批量操作 -->
    <div class="layer-actions">
      <ElButton size="small" link @click="selectAll">
        <ElIcon><Select /></ElIcon>
        全选
      </ElButton>
      <ElButton size="small" link @click="deselectAll">
        <ElIcon><CloseBold /></ElIcon>
        全不选
      </ElButton>
      <ElButton size="small" link @click="invertSelection">
        <ElIcon><Switch /></ElIcon>
        反选
      </ElButton>
    </div>

    <!-- 搜索框 -->
    <ElInput
      v-model="searchKeyword"
      placeholder="搜索图层..."
      size="small"
      clearable
      class="layer-search"
    >
      <template #prefix>
        <ElIcon><Search /></ElIcon>
      </template>
    </ElInput>

    <!-- 图层列表 -->
    <div class="layer-list">
      <div
        v-for="layer in filteredLayers"
        :key="layer.name"
        class="layer-item"
        :class="{ 'is-hidden': !layer.visible, 'is-selected': selectedLayer === layer.name }"
        @click="handleLayerClick(layer)"
      >
        <ElCheckbox v-model="layer.visible" @change="handleVisibilityChange(layer)" @click.stop />

        <div class="layer-info">
          <div class="layer-name" :title="layer.name">
            {{ layer.name }}
          </div>
          <div class="layer-meta">
            <span class="entity-count">{{ layer.entityCount }} 个实体</span>
            <span v-if="layer.frozen" class="frozen-badge">
              <ElIcon><Lock /></ElIcon>
            </span>
          </div>
        </div>

        <div class="layer-color" :style="{ backgroundColor: getColorHex(layer.color) }"></div>
      </div>

      <ElEmpty v-if="filteredLayers.length === 0" description="无匹配图层" :image-size="40" />
    </div>

    <!-- 图层统计 -->
    <div class="layer-stats">
      <div class="stat-item">
        <span class="stat-label">可见图层</span>
        <span class="stat-value">{{ visibleCount }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">隐藏图层</span>
        <span class="stat-value">{{ hiddenCount }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">总实体数</span>
        <span class="stat-value">{{ totalEntityCount }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  /***
   * Layer Control Component
   * 图层控制组件 - 用于控制CAD图层的显示/隐藏
   * Requirements: 8.3, 8.5
   ***/
  import { ref, computed } from 'vue'
  import { Search, Select, CloseBold, Switch, Lock } from '@element-plus/icons-vue'
  import type { CadLayerInfo } from '@/api/designer-cad'

  // Props定义 / Props definition
  interface LayerWithVisibility extends CadLayerInfo {
    visible: boolean
  }

  const props = defineProps<{
    layers: LayerWithVisibility[]
  }>()

  // Emits定义 / Emits definition
  const emit = defineEmits<{
    (e: 'update:layers', layers: LayerWithVisibility[]): void
    (e: 'layerClick', layerName: string): void
    (e: 'visibilityChange', layerName: string, visible: boolean): void
  }>()

  // 搜索关键字 / Search keyword
  const searchKeyword = ref('')

  // 选中的图层 / Selected layer
  const selectedLayer = ref<string | null>(null)

  // 过滤后的图层列表 / Filtered layer list
  const filteredLayers = computed(() => {
    if (!searchKeyword.value) {
      return props.layers
    }
    const keyword = searchKeyword.value.toLowerCase()
    return props.layers.filter((layer) => layer.name.toLowerCase().includes(keyword))
  })

  // 可见图层数量 / Visible layer count
  const visibleCount = computed(() => {
    return props.layers.filter((l) => l.visible).length
  })

  // 隐藏图层数量 / Hidden layer count
  const hiddenCount = computed(() => {
    return props.layers.filter((l) => !l.visible).length
  })

  // 总实体数量 / Total entity count
  const totalEntityCount = computed(() => {
    return props.layers.reduce((sum, layer) => sum + (layer.entityCount || 0), 0)
  })

  // 全选 / Select all
  const selectAll = () => {
    const updatedLayers = props.layers.map((layer) => ({
      ...layer,
      visible: true
    }))
    emit('update:layers', updatedLayers)
  }

  // 全不选 / Deselect all
  const deselectAll = () => {
    const updatedLayers = props.layers.map((layer) => ({
      ...layer,
      visible: false
    }))
    emit('update:layers', updatedLayers)
  }

  // 反选 / Invert selection
  const invertSelection = () => {
    const updatedLayers = props.layers.map((layer) => ({
      ...layer,
      visible: !layer.visible
    }))
    emit('update:layers', updatedLayers)
  }

  // 处理图层点击 / Handle layer click
  const handleLayerClick = (layer: LayerWithVisibility) => {
    selectedLayer.value = layer.name
    emit('layerClick', layer.name)
  }

  // 处理可见性变化 / Handle visibility change
  const handleVisibilityChange = (layer: LayerWithVisibility) => {
    emit('visibilityChange', layer.name, layer.visible)
  }

  // 颜色转换 / Color conversion
  const getColorHex = (colorIndex: number) => {
    const colorMap: Record<number, string> = {
      1: '#ff0000',
      2: '#ffff00',
      3: '#00ff00',
      4: '#00ffff',
      5: '#0000ff',
      6: '#ff00ff',
      7: '#ffffff',
      8: '#808080',
      9: '#c0c0c0'
    }
    return colorMap[colorIndex] || '#ffffff'
  }
</script>

<style scoped lang="scss">
  .layer-control {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--el-fill-color-light);
    border-radius: 4px;
    padding: 16px;

    .layer-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      h4 {
        margin: 0;
        font-size: 14px;
        font-weight: 500;
      }

      .layer-count {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }

    .layer-actions {
      display: flex;
      gap: 8px;
      margin-bottom: 12px;
      flex-wrap: wrap;
    }

    .layer-search {
      margin-bottom: 12px;
    }

    .layer-list {
      flex: 1;
      overflow-y: auto;
      margin: 0 -8px;
      padding: 0 8px;

      .layer-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;

        &:hover {
          background: var(--el-fill-color);
        }

        &.is-hidden {
          opacity: 0.5;
        }

        &.is-selected {
          background: var(--el-color-primary-light-9);
        }

        .layer-info {
          flex: 1;
          min-width: 0;

          .layer-name {
            font-size: 13px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }

          .layer-meta {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-top: 2px;

            .entity-count {
              font-size: 11px;
              color: var(--el-text-color-secondary);
            }

            .frozen-badge {
              font-size: 12px;
              color: var(--el-color-warning);
            }
          }
        }

        .layer-color {
          width: 16px;
          height: 16px;
          border-radius: 2px;
          border: 1px solid rgba(0, 0, 0, 0.1);
          flex-shrink: 0;
        }
      }
    }

    .layer-stats {
      display: flex;
      justify-content: space-between;
      padding-top: 12px;
      border-top: 1px solid var(--el-border-color-lighter);
      margin-top: 12px;

      .stat-item {
        text-align: center;

        .stat-label {
          display: block;
          font-size: 11px;
          color: var(--el-text-color-secondary);
          margin-bottom: 2px;
        }

        .stat-value {
          font-size: 14px;
          font-weight: 500;
          color: var(--el-color-primary);
        }
      }
    }
  }
</style>
