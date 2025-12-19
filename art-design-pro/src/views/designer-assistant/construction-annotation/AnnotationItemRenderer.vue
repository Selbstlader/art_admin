<template>
  <g
    class="annotation-item"
    :class="{ selected, drawing }"
    @click.stop="handleClick"
    @dblclick.stop="handleDoubleClick"
  >
    <!-- 标注背景 / Annotation background -->
    <rect
      :x="item.position.x - 2"
      :y="item.position.y - item.style.fontSize - 2"
      :width="textWidth + 8"
      :height="item.style.fontSize + 8"
      :fill="item.style.background"
      :stroke="selected ? '#409EFF' : item.style.lineColor"
      :stroke-width="selected ? 2 : 1"
      rx="2"
      ry="2"
    />

    <!-- 标注文本 / Annotation text -->
    <text
      :x="item.position.x + 2"
      :y="item.position.y"
      :fill="item.style.fontColor"
      :font-size="item.style.fontSize"
      font-family="Arial, sans-serif"
    >
      {{ item.content }}
    </text>

    <!-- 引线 / Leader line -->
    <line
      v-if="item.type === 'dimension'"
      :x1="item.position.x"
      :y1="item.position.y + 4"
      :x2="item.position.x"
      :y2="item.position.y + 20"
      :stroke="item.style.lineColor"
      :stroke-width="item.style.lineWidth"
    />

    <!-- 类型图标 / Type icon -->
    <circle
      :cx="item.position.x - 10"
      :cy="item.position.y - item.style.fontSize / 2"
      r="6"
      :fill="getTypeColor(item.type)"
      opacity="0.8"
    />

    <!-- 选中时的控制点 / Control points when selected -->
    <g v-if="selected && !drawing">
      <!-- 移动控制点 / Move control point -->
      <circle
        :cx="item.position.x + textWidth / 2"
        :cy="item.position.y - item.style.fontSize / 2"
        r="4"
        fill="#409EFF"
        class="control-point move"
        @mousedown.stop="startDrag"
      />

      <!-- 删除按钮 / Delete button -->
      <g
        class="delete-button"
        :transform="`translate(${item.position.x + textWidth + 10}, ${item.position.y - item.style.fontSize})`"
        @click.stop="handleDelete"
      >
        <circle r="8" fill="#F56C6C" />
        <text x="0" y="4" fill="#fff" font-size="12" text-anchor="middle">×</text>
      </g>
    </g>
  </g>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue'
  import type { AnnotationItem } from '@/api/designer-construction-annotation'

  /*** Component Props ***/
  const props = defineProps<{
    item: AnnotationItem
    selected?: boolean
    drawing?: boolean
  }>()

  /*** Component Emits ***/
  const emit = defineEmits<{
    (e: 'select', id: string): void
    (e: 'update', item: AnnotationItem): void
    (e: 'delete', id: string): void
  }>()

  /*** Reactive State ***/
  const isDragging = ref(false)
  const dragOffset = ref({ x: 0, y: 0 })

  /*** Computed Properties ***/
  // 估算文本宽度 / Estimate text width
  const textWidth = computed(() => {
    const charWidth = props.item.style.fontSize * 0.6
    return props.item.content.length * charWidth
  })

  /*** Methods ***/
  // 获取类型颜色 / Get type color
  const getTypeColor = (type: string): string => {
    switch (type) {
      case 'dimension':
        return '#FF0000'
      case 'material':
        return '#0000FF'
      case 'process':
        return '#00AA00'
      default:
        return '#666666'
    }
  }

  // 处理点击 / Handle click
  const handleClick = () => {
    emit('select', props.item.id)
  }

  // 处理双击 / Handle double click
  const handleDoubleClick = () => {
    // 可以在这里实现内联编辑 / Can implement inline editing here
    console.log('Double click on annotation:', props.item.id)
  }

  // 处理删除 / Handle delete
  const handleDelete = () => {
    emit('delete', props.item.id)
  }

  // 开始拖拽 / Start drag
  const startDrag = (event: MouseEvent) => {
    isDragging.value = true
    dragOffset.value = {
      x: event.clientX - props.item.position.x,
      y: event.clientY - props.item.position.y
    }

    document.addEventListener('mousemove', handleDrag)
    document.addEventListener('mouseup', stopDrag)
  }

  // 处理拖拽 / Handle drag
  const handleDrag = (event: MouseEvent) => {
    if (!isDragging.value) return

    const newItem = {
      ...props.item,
      position: {
        ...props.item.position,
        x: event.clientX - dragOffset.value.x,
        y: event.clientY - dragOffset.value.y
      }
    }

    emit('update', newItem)
  }

  // 停止拖拽 / Stop drag
  const stopDrag = () => {
    isDragging.value = false
    document.removeEventListener('mousemove', handleDrag)
    document.removeEventListener('mouseup', stopDrag)
  }
</script>

<style scoped lang="scss">
  .annotation-item {
    cursor: pointer;
    transition: opacity 0.2s;

    &:hover {
      opacity: 0.9;
    }

    &.selected {
      cursor: move;
    }

    &.drawing {
      opacity: 0.7;
    }
  }

  .control-point {
    cursor: move;

    &:hover {
      r: 6;
    }
  }

  .delete-button {
    cursor: pointer;

    &:hover circle {
      fill: #e6a23c;
    }
  }
</style>
