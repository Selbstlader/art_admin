<template>
  <g
    class="annotation-item"
    :class="{ selected, drawing, editing }"
    @click.stop="handleClick"
    @dblclick.stop="handleDoubleClick"
  >
    <!-- 引线 / Leader line -->
    <g v-if="hasLeaderLine">
      <!-- 主引线 / Main leader line -->
      <line
        :x1="lineStart.x"
        :y1="lineStart.y"
        :x2="lineEnd.x"
        :y2="lineEnd.y"
        :stroke="item.style.lineColor"
        :stroke-width="item.style.lineWidth"
        stroke-linecap="round"
      />
      <!-- 起点圆点 / Start point dot -->
      <circle
        :cx="lineStart.x"
        :cy="lineStart.y"
        :r="item.style.lineWidth + 2"
        :fill="item.style.lineColor"
      />
      <!-- 终点箭头或圆点 / End point arrow or dot -->
      <circle
        v-if="item.type !== 'dimension'"
        :cx="lineEnd.x"
        :cy="lineEnd.y"
        :r="item.style.lineWidth + 1"
        :fill="item.style.lineColor"
      />
      <!-- 尺寸标注的端点标记 / Dimension end markers -->
      <g v-if="item.type === 'dimension'">
        <line
          :x1="lineStart.x - 6"
          :y1="lineStart.y - 6"
          :x2="lineStart.x + 6"
          :y2="lineStart.y + 6"
          :stroke="item.style.lineColor"
          :stroke-width="item.style.lineWidth"
        />
        <line
          :x1="lineEnd.x - 6"
          :y1="lineEnd.y - 6"
          :x2="lineEnd.x + 6"
          :y2="lineEnd.y + 6"
          :stroke="item.style.lineColor"
          :stroke-width="item.style.lineWidth"
        />
      </g>
    </g>

    <!-- 文字背景 / Text background -->
    <rect
      v-if="item.content && item.style.background !== 'transparent'"
      :x="item.position.x - 4"
      :y="item.position.y - item.style.fontSize - 2"
      :width="textWidth + 8"
      :height="textHeight + 4"
      :fill="item.style.background"
      :stroke="selected ? '#409EFF' : 'none'"
      :stroke-width="selected ? 2 : 0"
      rx="3"
      ry="3"
    />

    <!-- 选中边框（无背景时）/ Selection border (when no background) -->
    <rect
      v-if="selected && item.style.background === 'transparent' && item.content"
      :x="item.position.x - 4"
      :y="item.position.y - item.style.fontSize - 2"
      :width="textWidth + 8"
      :height="textHeight + 4"
      fill="none"
      stroke="#409EFF"
      stroke-width="1"
      stroke-dasharray="4,2"
      rx="3"
      ry="3"
    />

    <!-- 标注文本 / Annotation text -->
    <text
      v-if="item.content"
      :x="item.position.x"
      :y="item.position.y"
      :fill="item.style.fontColor"
      :font-size="item.style.fontSize"
      font-family="Arial, 'Microsoft YaHei', sans-serif"
      dominant-baseline="auto"
    >
      <tspan
        v-for="(line, index) in textLines"
        :key="index"
        :x="item.position.x"
        :dy="index === 0 ? 0 : item.style.fontSize * 1.2"
      >
        {{ line }}
      </tspan>
    </text>

    <!-- 类型图标指示器 / Type indicator -->
    <g v-if="!drawing && item.content" :transform="`translate(${item.position.x - 16}, ${item.position.y - item.style.fontSize / 2 - 4})`">
      <circle r="6" :fill="getTypeColor(item.type)" opacity="0.9" />
    </g>

    <!-- 选中时的控制点 / Control points when selected -->
    <g v-if="selected && !drawing && !editing">
      <!-- 移动控制点 / Move control point -->
      <circle
        :cx="item.position.x + textWidth / 2"
        :cy="item.position.y - item.style.fontSize / 2"
        r="5"
        fill="#409EFF"
        stroke="#fff"
        stroke-width="2"
        class="control-point move"
        @mousedown.stop="startDrag"
      />

      <!-- 删除按钮 / Delete button -->
      <g
        class="delete-button"
        :transform="`translate(${item.position.x + textWidth + 12}, ${item.position.y - item.style.fontSize})`"
        @click.stop="handleDelete"
      >
        <circle r="10" fill="#F56C6C" />
        <text x="0" y="4" fill="#fff" font-size="14" text-anchor="middle" font-weight="bold">×</text>
      </g>

      <!-- 编辑按钮 / Edit button -->
      <g
        class="edit-button"
        :transform="`translate(${item.position.x + textWidth + 34}, ${item.position.y - item.style.fontSize})`"
        @click.stop="handleEdit"
      >
        <circle r="10" fill="#409EFF" />
        <text x="0" y="4" fill="#fff" font-size="10" text-anchor="middle">✎</text>
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
    editing?: boolean
  }>()

  /*** Component Emits ***/
  const emit = defineEmits<{
    (e: 'select', id: string): void
    (e: 'update', item: AnnotationItem): void
    (e: 'delete', id: string): void
    (e: 'startEdit', id: string): void
  }>()

  /*** Reactive State ***/
  const isDragging = ref(false)
  const dragStartPos = ref({ x: 0, y: 0 })
  const itemStartPos = ref({ x: 0, y: 0 })

  /*** Computed Properties ***/
  // 文本行 / Text lines
  const textLines = computed(() => {
    return props.item.content.split('\n')
  })

  // 估算文本宽度 / Estimate text width
  const textWidth = computed(() => {
    const maxLineLength = Math.max(...textLines.value.map((line) => line.length))
    const charWidth = props.item.style.fontSize * 0.6
    return Math.max(maxLineLength * charWidth, 20)
  })

  // 文本高度 / Text height
  const textHeight = computed(() => {
    return textLines.value.length * props.item.style.fontSize * 1.2
  })

  // 是否有引线 / Has leader line
  const hasLeaderLine = computed(() => {
    return props.item.properties?.lineStart && props.item.properties?.lineEnd
  })

  // 引线起点 / Line start point
  const lineStart = computed(() => {
    return props.item.properties?.lineStart || { x: 0, y: 0 }
  })

  // 引线终点 / Line end point
  const lineEnd = computed(() => {
    return props.item.properties?.lineEnd || { x: 0, y: 0 }
  })

  /*** Methods ***/
  // 获取类型颜色 / Get type color
  const getTypeColor = (type: string): string => {
    const colors: Record<string, string> = {
      dimension: '#E53935',
      material: '#1E88E5',
      process: '#43A047',
      line: '#FF9800',
      text: '#666666'
    }
    return colors[type] || '#666666'
  }

  // 处理点击 / Handle click
  const handleClick = () => {
    if (!props.drawing) {
      emit('select', props.item.id)
    }
  }

  // 处理双击 / Handle double click
  const handleDoubleClick = () => {
    if (!props.drawing) {
      emit('startEdit', props.item.id)
    }
  }

  // 处理删除 / Handle delete
  const handleDelete = () => {
    emit('delete', props.item.id)
  }

  // 处理编辑 / Handle edit
  const handleEdit = () => {
    emit('startEdit', props.item.id)
  }

  // 开始拖拽 / Start drag
  const startDrag = (event: MouseEvent) => {
    isDragging.value = true
    
    // 获取SVG元素 / Get SVG element
    const svg = (event.target as Element).closest('svg') as SVGSVGElement
    if (!svg) return

    const rect = svg.getBoundingClientRect()
    const viewBox = svg.viewBox.baseVal
    const scaleX = viewBox.width / rect.width
    const scaleY = viewBox.height / rect.height

    dragStartPos.value = {
      x: (event.clientX - rect.left) * scaleX,
      y: (event.clientY - rect.top) * scaleY
    }
    itemStartPos.value = {
      x: props.item.position.x,
      y: props.item.position.y
    }

    document.addEventListener('mousemove', handleDrag)
    document.addEventListener('mouseup', stopDrag)
  }

  // 处理拖拽 / Handle drag
  const handleDrag = (event: MouseEvent) => {
    if (!isDragging.value) return

    const svgElements = document.querySelectorAll('.annotation-layer')
    if (svgElements.length === 0) return

    const svg = svgElements[0] as SVGSVGElement
    const rect = svg.getBoundingClientRect()
    const viewBox = svg.viewBox.baseVal
    const scaleX = viewBox.width / rect.width
    const scaleY = viewBox.height / rect.height

    const currentX = (event.clientX - rect.left) * scaleX
    const currentY = (event.clientY - rect.top) * scaleY

    const deltaX = currentX - dragStartPos.value.x
    const deltaY = currentY - dragStartPos.value.y

    const newX = Math.max(0, itemStartPos.value.x + deltaX)
    const newY = Math.max(0, itemStartPos.value.y + deltaY)

    // 同时移动引线 / Also move leader line
    let newProperties = { ...props.item.properties }
    if (props.item.properties?.lineStart && props.item.properties?.lineEnd) {
      const lineStartX = props.item.properties.lineStart.x
      const lineStartY = props.item.properties.lineStart.y
      const lineEndX = props.item.properties.lineEnd.x
      const lineEndY = props.item.properties.lineEnd.y
      
      // 计算原始偏移 / Calculate original offset
      const originalOffsetX = props.item.position.x - lineEndX
      const originalOffsetY = props.item.position.y - lineEndY
      
      newProperties = {
        ...newProperties,
        lineStart: {
          x: lineStartX + deltaX,
          y: lineStartY + deltaY
        },
        lineEnd: {
          x: newX - originalOffsetX,
          y: newY - originalOffsetY
        }
      }
    }

    const newItem = {
      ...props.item,
      position: {
        ...props.item.position,
        x: newX,
        y: newY
      },
      properties: newProperties
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
      opacity: 0.95;
    }

    &.selected {
      cursor: move;
    }

    &.drawing {
      opacity: 0.7;
      pointer-events: none;
    }

    &.editing {
      opacity: 0.3;
    }
  }

  .control-point {
    cursor: move;
    transition: r 0.2s;

    &:hover {
      r: 7;
    }
  }

  .delete-button,
  .edit-button {
    cursor: pointer;
    transition: transform 0.2s;

    &:hover {
      transform: scale(1.1);
    }
  }
</style>
