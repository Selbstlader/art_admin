<template>
  <div class="annotation-editor">
    <!-- 工具栏 / Toolbar -->
    <div class="editor-toolbar">
      <el-button-group>
        <el-button :type="currentTool === 'select' ? 'primary' : 'default'" @click="setTool('select')">
          <el-icon><Select /></el-icon>
          选择
        </el-button>
        <el-button :type="currentTool === 'dimension' ? 'primary' : 'default'" @click="setTool('dimension')">
          <el-icon><Edit /></el-icon>
          尺寸标注
        </el-button>
        <el-button :type="currentTool === 'material' ? 'primary' : 'default'" @click="setTool('material')">
          <el-icon><Grid /></el-icon>
          材料标注
        </el-button>
        <el-button :type="currentTool === 'process' ? 'primary' : 'default'" @click="setTool('process')">
          <el-icon><Setting /></el-icon>
          工艺说明
        </el-button>
        <el-button :type="currentTool === 'line' ? 'primary' : 'default'" @click="setTool('line')">
          <el-icon><Share /></el-icon>
          引线
        </el-button>
        <el-button :type="currentTool === 'text' ? 'primary' : 'default'" @click="setTool('text')">
          <el-icon><EditPen /></el-icon>
          文字
        </el-button>
      </el-button-group>

      <!-- 缩放控制 / Zoom controls -->
      <div class="zoom-controls">
        <el-button-group>
          <el-button size="small" @click="zoomOut" :disabled="zoomLevel <= 0.25">
            <el-icon><ZoomOut /></el-icon>
          </el-button>
          <el-button size="small" disabled style="min-width: 60px">{{ Math.round(zoomLevel * 100) }}%</el-button>
          <el-button size="small" @click="zoomIn" :disabled="zoomLevel >= 3">
            <el-icon><ZoomIn /></el-icon>
          </el-button>
          <el-button size="small" @click="resetZoom">
            <el-icon><RefreshRight /></el-icon>
          </el-button>
        </el-button-group>
      </div>

      <div class="toolbar-actions">
        <el-button @click="handleAutoAnnotate" :loading="analyzing">
          <el-icon><MagicStick /></el-icon>
          自动标注
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          <el-icon><Check /></el-icon>
          保存
        </el-button>
      </div>
    </div>

    <!-- 主编辑区域 / Main editing area -->
    <div class="editor-main">
      <!-- 画布区域 / Canvas area -->
      <div class="canvas-container" ref="canvasContainer" @wheel="handleWheel">
        <div
          class="canvas-wrapper"
          ref="canvasWrapper"
          :class="{ 'drawing-mode': currentTool !== 'select' }"
          :style="canvasWrapperStyle"
          @mousedown="handleCanvasMouseDown"
          @mousemove="handleCanvasMouseMove"
          @mouseup="handleCanvasMouseUp"
          @mouseleave="handleCanvasMouseLeave"
        >
          <!-- 背景图片 / Background image -->
          <img
            v-if="imagePath"
            :src="imagePath"
            class="background-image"
            ref="backgroundImage"
            @load="handleImageLoad"
            draggable="false"
          />

          <!-- 标注层 / Annotation layer -->
          <svg
            class="annotation-layer"
            :width="displayWidth"
            :height="displayHeight"
            :viewBox="svgViewBox"
          >
            <!-- 渲染已有标注项 / Render existing annotation items -->
            <g v-for="item in annotations" :key="item.id">
              <AnnotationItemRenderer
                :item="item"
                :selected="selectedItemId === item.id"
                :editing="editingItemId === item.id"
                @select="handleSelectItem"
                @update="handleUpdateItem"
                @delete="handleDeleteItem"
                @startEdit="handleStartEdit"
              />
            </g>

            <!-- 正在绘制的标注 / Currently drawing annotation -->
            <g v-if="drawingItem">
              <!-- 引线绘制预览 / Line drawing preview -->
              <line
                v-if="drawingItem.properties?.lineEnd"
                :x1="drawingItem.properties.lineStart?.x || drawingItem.position.x"
                :y1="drawingItem.properties.lineStart?.y || drawingItem.position.y"
                :x2="drawingItem.properties.lineEnd.x"
                :y2="drawingItem.properties.lineEnd.y"
                :stroke="drawingItem.style.lineColor"
                :stroke-width="drawingItem.style.lineWidth"
                stroke-dasharray="5,5"
              />
              <AnnotationItemRenderer :item="drawingItem" :drawing="true" />
            </g>
          </svg>

          <!-- 文字输入框 / Text input overlay -->
          <div
            v-if="showTextInput"
            class="text-input-overlay"
            :style="textInputStyle"
          >
            <textarea
              ref="textInputRef"
              v-model="textInputValue"
              class="text-input"
              :style="textInputTextStyle"
              @blur="handleTextInputBlur"
              @keydown.enter.exact="handleTextInputConfirm"
              @keydown.escape="handleTextInputCancel"
              placeholder="输入标注内容..."
              autofocus
            />
          </div>

          <!-- 绘制提示 / Drawing hint -->
          <div v-if="currentTool !== 'select' && !isDrawing" class="drawing-hint">
            {{ getDrawingHint(currentTool) }}
          </div>
        </div>
      </div>

      <!-- 属性面板 / Properties panel -->
      <div class="properties-panel">
        <el-card v-if="selectedItem" class="property-card">
          <template #header>
            <div class="card-header">
              <span>标注属性</span>
              <el-button type="danger" size="small" @click="handleDeleteSelected">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </template>

          <el-form label-width="80px" size="small">
            <el-form-item label="类型">
              <el-select v-model="selectedItem.type" @change="handlePropertyChange">
                <el-option label="尺寸标注" value="dimension" />
                <el-option label="材料标注" value="material" />
                <el-option label="工艺说明" value="process" />
                <el-option label="引线" value="line" />
                <el-option label="文字" value="text" />
              </el-select>
            </el-form-item>

            <el-form-item label="内容">
              <el-input
                v-model="selectedItem.content"
                type="textarea"
                :rows="2"
                @change="handlePropertyChange"
              />
            </el-form-item>

            <el-form-item label="X坐标">
              <el-input-number v-model="selectedItem.position.x" @change="handlePropertyChange" />
            </el-form-item>

            <el-form-item label="Y坐标">
              <el-input-number v-model="selectedItem.position.y" @change="handlePropertyChange" />
            </el-form-item>

            <el-form-item label="字体大小">
              <el-input-number
                v-model="selectedItem.style.fontSize"
                :min="8"
                :max="72"
                @change="handlePropertyChange"
              />
            </el-form-item>

            <el-form-item label="字体颜色">
              <el-color-picker
                v-model="selectedItem.style.fontColor"
                @change="handlePropertyChange"
              />
            </el-form-item>

            <el-form-item label="线条颜色">
              <el-color-picker
                v-model="selectedItem.style.lineColor"
                @change="handlePropertyChange"
              />
            </el-form-item>

            <el-form-item label="线条粗细">
              <el-input-number
                v-model="selectedItem.style.lineWidth"
                :min="1"
                :max="10"
                @change="handlePropertyChange"
              />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 标注列表 / Annotation list -->
        <el-card class="annotation-list-card">
          <template #header>
            <span>标注列表 ({{ annotations.length }})</span>
          </template>

          <el-scrollbar height="400px">
            <div v-if="annotations.length === 0" class="empty-list">
              暂无标注，请选择工具后在画布上绘制
            </div>
            <div
              v-for="item in annotations"
              :key="item.id"
              class="annotation-list-item"
              :class="{ active: selectedItemId === item.id }"
              @click="handleSelectItem(item.id)"
              @dblclick="handleStartEdit(item.id)"
            >
              <el-icon class="item-icon">
                <Edit v-if="item.type === 'dimension'" />
                <Grid v-else-if="item.type === 'material'" />
                <Setting v-else-if="item.type === 'process'" />
                <Share v-else-if="item.type === 'line'" />
                <EditPen v-else />
              </el-icon>
              <span class="item-content">{{ item.content || '(空)' }}</span>
              <el-tag size="small" :type="item.isGenerated ? 'info' : 'success'">
                {{ item.isGenerated ? '自动' : '手动' }}
              </el-tag>
            </div>
          </el-scrollbar>
        </el-card>
      </div>
    </div>

    <!-- 导出对话框 / Export dialog -->
    <el-dialog v-model="exportDialogVisible" title="导出标注" width="400px">
      <el-form label-width="100px">
        <el-form-item label="导出格式">
          <el-select v-model="exportConfig.format">
            <el-option
              v-for="format in exportFormats"
              :key="format.format"
              :label="format.name"
              :value="format.format"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="图片质量" v-if="exportConfig.format !== 'pdf'">
          <el-slider v-model="exportConfig.quality" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="缩放比例">
          <el-input-number v-model="exportConfig.scale" :min="0.1" :max="10" :step="0.1" />
        </el-form-item>
        <el-form-item label="显示图层">
          <el-switch v-model="exportConfig.showLayers" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="exportDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmExport" :loading="exporting">导出</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, watch, nextTick } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    Select,
    Edit,
    Grid,
    Setting,
    MagicStick,
    Download,
    Check,
    Delete,
    Share,
    EditPen,
    ZoomIn,
    ZoomOut,
    RefreshRight
  } from '@element-plus/icons-vue'
  import AnnotationItemRenderer from './AnnotationItemRenderer.vue'
  import {
    type AnnotationItem,
    type ConstructionAnnotationResponse,
    type AnnotationExportConfig,
    type ExportFormat,
    getAnnotation,
    updateAnnotations,
    addAnnotationItem,
    deleteAnnotationItem,
    exportAnnotation,
    getExportFormats,
    analyzeConstructionDrawing
  } from '@/api/designer-construction-annotation'
  import { useRoute } from 'vue-router'

  const route = useRoute()

  /*** Component Emits ***/
  const emit = defineEmits<{
    (e: 'saved'): void
    (e: 'exported', result: any): void
  }>()

  /*** Computed route params ***/
  const annotationId = computed(() => Number(route.query.annotationId) || 0)
  const projectId = computed(() => Number(route.query.projectId) || 0)
  const cadFileId = computed(() => Number(route.query.cadFileId) || 0)

  /*** Reactive State ***/
  const currentTool = ref<string>('select')
  const annotations = ref<AnnotationItem[]>([])
  const selectedItemId = ref<string | null>(null)
  const editingItemId = ref<string | null>(null)
  const imagePath = ref<string>('')
  const analyzing = ref(false)
  const saving = ref(false)
  const exporting = ref(false)
  const exportDialogVisible = ref(false)
  const exportFormats = ref<ExportFormat[]>([])
  const exportConfig = ref<AnnotationExportConfig>({
    format: 'pdf',
    quality: 90,
    scale: 1,
    showLayers: true
  })

  // 画布相关状态 / Canvas related state
  const canvasContainer = ref<HTMLElement | null>(null)
  const canvasWrapper = ref<HTMLElement | null>(null)
  const backgroundImage = ref<HTMLImageElement | null>(null)
  const imageWidth = ref(800)
  const imageHeight = ref(600)
  const displayWidth = ref(800)
  const displayHeight = ref(600)
  const zoomLevel = ref(1)

  // 绘制状态 / Drawing state
  const drawingItem = ref<AnnotationItem | null>(null)
  const isDrawing = ref(false)
  const startPoint = ref({ x: 0, y: 0 })
  const lastUsedTool = ref<string>('dimension') // 记住上次使用的工具

  // 文字输入状态 / Text input state
  const showTextInput = ref(false)
  const textInputValue = ref('')
  const textInputRef = ref<HTMLTextAreaElement | null>(null)
  const textInputPosition = ref({ x: 0, y: 0 })
  const pendingTextItem = ref<AnnotationItem | null>(null)

  /*** Computed Properties ***/
  const selectedItem = computed(() => {
    if (!selectedItemId.value) return null
    return annotations.value.find((item) => item.id === selectedItemId.value) || null
  })

  const svgViewBox = computed(() => `0 0 ${imageWidth.value} ${imageHeight.value}`)

  // 画布缩放样式 / Canvas zoom style
  const canvasWrapperStyle = computed(() => ({
    transform: `scale(${zoomLevel.value})`,
    transformOrigin: 'top left'
  }))

  const textInputStyle = computed(() => {
    const img = backgroundImage.value
    if (!img) return { left: '0px', top: '0px' }
    
    // 考虑缩放 / Consider zoom
    const scaleX = (img.clientWidth * zoomLevel.value) / imageWidth.value
    const scaleY = (img.clientHeight * zoomLevel.value) / imageHeight.value
    
    return {
      left: `${textInputPosition.value.x * scaleX}px`,
      top: `${textInputPosition.value.y * scaleY}px`
    }
  })

  const textInputTextStyle = computed(() => {
    const item = pendingTextItem.value
    if (!item) return {}
    
    const img = backgroundImage.value
    // 考虑缩放 / Consider zoom
    const scale = img ? (img.clientWidth * zoomLevel.value) / imageWidth.value : 1
    
    return {
      fontSize: `${Math.max(14, item.style.fontSize * scale)}px`,
      color: item.style.fontColor
    }
  })

  /*** Type for API response ***/
  interface BaseResponse<T = unknown> {
    code: number
    msg?: string
    data: T
  }

  /*** Methods ***/
  // 设置当前工具 / Set current tool
  const setTool = (tool: string) => {
    currentTool.value = tool
    if (tool !== 'select') {
      selectedItemId.value = null
      editingItemId.value = null
      lastUsedTool.value = tool // 记住工具
    }
  }

  // 缩放控制 / Zoom controls
  const zoomIn = () => {
    zoomLevel.value = Math.min(3, zoomLevel.value + 0.25)
  }

  const zoomOut = () => {
    zoomLevel.value = Math.max(0.25, zoomLevel.value - 0.25)
  }

  const resetZoom = () => {
    zoomLevel.value = 1
  }

  // 鼠标滚轮缩放 / Mouse wheel zoom
  const handleWheel = (event: WheelEvent) => {
    if (event.ctrlKey || event.metaKey) {
      event.preventDefault()
      if (event.deltaY < 0) {
        zoomIn()
      } else {
        zoomOut()
      }
    }
  }

  // 获取绘制提示 / Get drawing hint
  const getDrawingHint = (tool: string): string => {
    switch (tool) {
      case 'dimension':
        return '点击并拖拽绘制尺寸标注线，松开后输入尺寸'
      case 'material':
        return '点击并拖拽绘制材料标注，松开后输入材料名称'
      case 'process':
        return '点击并拖拽绘制工艺说明，松开后输入说明内容'
      case 'line':
        return '点击起点并拖拽到终点绘制引线'
      case 'text':
        return '点击画布位置添加文字标注'
      default:
        return '选择工具后在画布上绘制'
    }
  }

  // 加载标注数据 / Load annotation data
  const loadAnnotation = async () => {
    if (!annotationId.value) return
    try {
      const res = (await getAnnotation(
        annotationId.value
      )) as unknown as BaseResponse<ConstructionAnnotationResponse>
      if (res.code === 200 && res.data) {
        annotations.value = res.data.annotations || []
        imagePath.value = res.data.imagePath
      }
    } catch {
      ElMessage.error('加载标注数据失败')
    }
  }

  // 加载导出格式 / Load export formats
  const loadExportFormats = async () => {
    try {
      const res = (await getExportFormats()) as unknown as BaseResponse<ExportFormat[]>
      if (res.code === 200 && res.data) {
        exportFormats.value = res.data
      }
    } catch {
      console.error('加载导出格式失败')
    }
  }

  // 处理图片加载 / Handle image load
  const handleImageLoad = (event: Event) => {
    const img = event.target as HTMLImageElement
    imageWidth.value = img.naturalWidth
    imageHeight.value = img.naturalHeight
    displayWidth.value = img.clientWidth
    displayHeight.value = img.clientHeight
  }

  // 获取画布内的鼠标坐标 / Get mouse coordinates within canvas
  const getCanvasCoordinates = (event: MouseEvent): { x: number; y: number } => {
    const wrapper = canvasWrapper.value
    const img = backgroundImage.value
    if (!wrapper || !img) return { x: 0, y: 0 }

    const rect = img.getBoundingClientRect()
    // 考虑缩放因素 / Consider zoom factor
    const scaleX = imageWidth.value / (img.clientWidth * zoomLevel.value)
    const scaleY = imageHeight.value / (img.clientHeight * zoomLevel.value)

    return {
      x: Math.round((event.clientX - rect.left) * scaleX),
      y: Math.round((event.clientY - rect.top) * scaleY)
    }
  }

  // 获取默认样式 / Get default style - 增大字体
  const getDefaultStyle = (type: string) => {
    const styles: Record<string, any> = {
      dimension: {
        fontSize: 24,
        fontColor: '#E53935',
        lineColor: '#E53935',
        lineWidth: 3,
        background: 'rgba(255,255,255,0.95)'
      },
      material: {
        fontSize: 22,
        fontColor: '#1E88E5',
        lineColor: '#1E88E5',
        lineWidth: 3,
        background: 'rgba(255,255,255,0.95)'
      },
      process: {
        fontSize: 22,
        fontColor: '#43A047',
        lineColor: '#43A047',
        lineWidth: 3,
        background: 'rgba(255,255,255,0.95)'
      },
      line: {
        fontSize: 22,
        fontColor: '#FF9800',
        lineColor: '#FF9800',
        lineWidth: 3,
        background: 'rgba(255,255,255,0.95)'
      },
      text: {
        fontSize: 28,
        fontColor: '#333333',
        lineColor: '#333333',
        lineWidth: 2,
        background: 'rgba(255,255,255,0.9)'
      }
    }
    return styles[type] || styles.text
  }

  // 处理画布鼠标按下 / Handle canvas mouse down
  const handleCanvasMouseDown = (event: MouseEvent) => {
    if (event.button !== 0) return
    if (currentTool.value === 'select') return
    if (showTextInput.value) return

    event.preventDefault()
    const coords = getCanvasCoordinates(event)
    startPoint.value = coords
    isDrawing.value = true

    // 创建新标注项 / Create new annotation item
    const newItem: AnnotationItem = {
      id: `annotation_${Date.now()}`,
      type: currentTool.value,
      position: { x: coords.x, y: coords.y, anchor: 'top-left' },
      content: '',
      style: getDefaultStyle(currentTool.value),
      properties: {
        lineStart: { x: coords.x, y: coords.y },
        lineEnd: null
      },
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      isGenerated: false
    }

    drawingItem.value = newItem

    // 纯文字工具直接显示输入框 / Text tool shows input immediately
    if (currentTool.value === 'text') {
      isDrawing.value = false
      showTextInputAt(coords, newItem)
    }
  }

  // 处理画布鼠标移动 / Handle canvas mouse move
  const handleCanvasMouseMove = (event: MouseEvent) => {
    if (!isDrawing.value || !drawingItem.value) return

    const coords = getCanvasCoordinates(event)
    
    // 更新引线终点 / Update line end point
    drawingItem.value = {
      ...drawingItem.value,
      properties: {
        ...drawingItem.value.properties,
        lineEnd: { x: coords.x, y: coords.y }
      }
    }
  }

  // 处理画布鼠标抬起 / Handle canvas mouse up
  const handleCanvasMouseUp = (event: MouseEvent) => {
    if (!isDrawing.value || !drawingItem.value) return

    isDrawing.value = false
    const coords = getCanvasCoordinates(event)
    
    // 计算文字位置（引线终点偏移一点）/ Calculate text position
    const textX = coords.x + 10
    const textY = coords.y - 5

    const finalItem: AnnotationItem = {
      ...drawingItem.value,
      position: { x: textX, y: textY, anchor: 'top-left' },
      properties: {
        lineStart: startPoint.value,
        lineEnd: coords
      }
    }

    // 显示文字输入框 / Show text input
    showTextInputAt({ x: textX, y: textY }, finalItem)
  }

  // 处理鼠标离开画布 / Handle mouse leave canvas
  const handleCanvasMouseLeave = () => {
    if (isDrawing.value && drawingItem.value) {
      // 取消绘制 / Cancel drawing
      isDrawing.value = false
      drawingItem.value = null
    }
  }

  // 显示文字输入框 / Show text input at position
  const showTextInputAt = (coords: { x: number; y: number }, item: AnnotationItem) => {
    textInputPosition.value = coords
    pendingTextItem.value = item
    textInputValue.value = item.content || ''
    showTextInput.value = true
    drawingItem.value = null

    nextTick(() => {
      textInputRef.value?.focus()
    })
  }

  // 处理文字输入确认 / Handle text input confirm
  const handleTextInputConfirm = async (event: KeyboardEvent) => {
    event.preventDefault()
    await saveTextAnnotation()
  }

  // 处理文字输入失焦 / Handle text input blur
  const handleTextInputBlur = async () => {
    await saveTextAnnotation()
  }

  // 处理文字输入取消 / Handle text input cancel
  const handleTextInputCancel = () => {
    showTextInput.value = false
    textInputValue.value = ''
    pendingTextItem.value = null
    // 不切换工具，保持当前工具 / Don't switch tool
  }

  // 保存文字标注 / Save text annotation
  const saveTextAnnotation = async () => {
    if (!pendingTextItem.value) return

    const content = textInputValue.value.trim()
    const toolUsed = pendingTextItem.value.type
    
    if (!content) {
      handleTextInputCancel()
      return
    }

    const newItem: AnnotationItem = {
      ...pendingTextItem.value,
      content
    }

    // 添加到列表 / Add to list
    annotations.value.push(newItem)

    // 保存到后端 / Save to backend
    try {
      if (annotationId.value) {
        await addAnnotationItem(annotationId.value, newItem)
      }
      ElMessage.success('标注添加成功')
    } catch (error) {
      console.error('保存标注失败:', error)
    }

    // 重置状态 / Reset state
    showTextInput.value = false
    textInputValue.value = ''
    pendingTextItem.value = null
    selectedItemId.value = null
    
    // 保持当前工具，继续添加标注 / Keep current tool for continuous annotation
    currentTool.value = toolUsed
  }

  // 开始编辑标注 / Start editing annotation
  const handleStartEdit = (id: string) => {
    const item = annotations.value.find((a) => a.id === id)
    if (!item) return

    editingItemId.value = id
    selectedItemId.value = id
    
    // 显示文字输入框进行编辑 / Show text input for editing
    textInputPosition.value = { x: item.position.x, y: item.position.y }
    pendingTextItem.value = { ...item }
    textInputValue.value = item.content
    showTextInput.value = true

    // 从列表中临时移除 / Temporarily remove from list
    annotations.value = annotations.value.filter((a) => a.id !== id)

    nextTick(() => {
      textInputRef.value?.focus()
      textInputRef.value?.select()
    })
  }

  // 选择标注项 / Select annotation item
  const handleSelectItem = (id: string) => {
    selectedItemId.value = id
    editingItemId.value = null
    currentTool.value = 'select'
  }

  // 更新标注项 / Update annotation item
  const handleUpdateItem = (item: AnnotationItem) => {
    const index = annotations.value.findIndex((a) => a.id === item.id)
    if (index !== -1) {
      annotations.value[index] = { ...item }
    }
  }

  // 删除标注项 / Delete annotation item
  const handleDeleteItem = async (id: string) => {
    try {
      if (annotationId.value) {
        await deleteAnnotationItem(annotationId.value, id)
      }
      annotations.value = annotations.value.filter((item) => item.id !== id)
      if (selectedItemId.value === id) {
        selectedItemId.value = null
      }
      ElMessage.success('标注删除成功')
    } catch {
      // 即使后端失败也删除本地 / Delete local even if backend fails
      annotations.value = annotations.value.filter((item) => item.id !== id)
    }
  }

  // 删除选中的标注 / Delete selected annotation
  const handleDeleteSelected = () => {
    if (selectedItemId.value) {
      handleDeleteItem(selectedItemId.value)
    }
  }

  // 属性变更处理 / Handle property change
  const handlePropertyChange = () => {
    if (selectedItem.value) {
      handleUpdateItem(selectedItem.value)
    }
  }

  // 自动标注 / Auto annotate
  const handleAutoAnnotate = async () => {
    analyzing.value = true
    try {
      await analyzeConstructionDrawing({
        projectId: projectId.value,
        cadFileId: cadFileId.value,
        imagePath: imagePath.value
      })
      ElMessage.success('自动标注已开始，请稍后刷新查看结果')
      setTimeout(() => loadAnnotation(), 3000)
    } catch {
      ElMessage.error('自动标注失败')
    } finally {
      analyzing.value = false
    }
  }

  // 保存标注 / Save annotations
  const handleSave = async () => {
    saving.value = true
    try {
      await updateAnnotations(annotationId.value, annotations.value)
      ElMessage.success('保存成功')
      emit('saved')
    } catch {
      ElMessage.error('保存失败')
    } finally {
      saving.value = false
    }
  }

  // 导出标注 / Export annotation
  const handleExport = () => {
    exportDialogVisible.value = true
  }

  // 确认导出 / Confirm export
  const confirmExport = async () => {
    exporting.value = true
    try {
      const res = (await exportAnnotation(
        annotationId.value,
        exportConfig.value
      )) as unknown as BaseResponse<{ downloadUrl?: string }>
      ElMessage.success('导出成功')
      exportDialogVisible.value = false
      emit('exported', res.data)
    } catch {
      ElMessage.error('导出失败')
    } finally {
      exporting.value = false
    }
  }

  /*** Lifecycle Hooks ***/
  onMounted(() => {
    loadAnnotation()
    loadExportFormats()
  })

  /*** Expose methods for parent component ***/
  defineExpose({
    openExportDialog: () => {
      exportDialogVisible.value = true
    }
  })

  watch(
    () => annotationId.value,
    () => {
      loadAnnotation()
    }
  )
</script>

<style scoped lang="scss">
  .annotation-editor {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #f5f7fa;
  }

  .editor-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #fff;
    border-bottom: 1px solid #e4e7ed;
    flex-shrink: 0;
    flex-wrap: wrap;
    gap: 12px;
  }

  .zoom-controls {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .toolbar-actions {
    display: flex;
    gap: 8px;
  }

  .editor-main {
    display: flex;
    flex: 1;
    overflow: hidden;
    min-height: 0;
  }

  .canvas-container {
    flex: 1;
    overflow: auto;
    padding: 16px;
    background: #e4e7ed;
  }

  .canvas-wrapper {
    position: relative;
    background: #fff;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
    display: inline-block;
    transition: transform 0.1s ease-out;

    &.drawing-mode {
      cursor: crosshair;
    }
  }

  .background-image {
    display: block;
    user-select: none;
    -webkit-user-drag: none;
  }

  .annotation-layer {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }

  .annotation-layer g {
    pointer-events: auto;
  }

  .text-input-overlay {
    position: absolute;
    z-index: 100;
    min-width: 120px;
    max-width: 300px;
  }

  .text-input {
    width: 100%;
    min-width: 120px;
    min-height: 32px;
    padding: 4px 8px;
    border: 2px solid #409eff;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.95);
    outline: none;
    resize: both;
    font-family: inherit;
    line-height: 1.4;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

    &::placeholder {
      color: #c0c4cc;
    }
  }

  .drawing-hint {
    position: absolute;
    bottom: 16px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.75);
    color: #fff;
    padding: 10px 20px;
    border-radius: 6px;
    font-size: 14px;
    pointer-events: none;
    white-space: nowrap;
  }

  .properties-panel {
    width: 320px;
    padding: 16px;
    background: #fff;
    border-left: 1px solid #e4e7ed;
    overflow-y: auto;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .property-card {
    flex-shrink: 0;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .annotation-list-card {
    flex: 1;
    min-height: 200px;
    display: flex;
    flex-direction: column;

    :deep(.el-card__body) {
      flex: 1;
      padding: 0;
      overflow: hidden;
    }
  }

  .empty-list {
    padding: 24px 16px;
    text-align: center;
    color: #909399;
    font-size: 14px;
  }

  .annotation-list-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    cursor: pointer;
    border-bottom: 1px solid #f0f0f0;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: #ecf5ff;
    }

    &:last-child {
      border-bottom: none;
    }

    .item-icon {
      flex-shrink: 0;
      color: #606266;
    }

    .item-content {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      font-size: 14px;
      color: #303133;
    }

    .el-tag {
      flex-shrink: 0;
    }
  }
</style>
