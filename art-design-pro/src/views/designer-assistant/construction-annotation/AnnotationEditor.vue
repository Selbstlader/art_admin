<template>
  <div class="annotation-editor">
    <!-- 工具栏 / Toolbar -->
    <div class="editor-toolbar">
      <el-button-group>
        <el-button
          :type="currentTool === 'select' ? 'primary' : 'default'"
          @click="setTool('select')"
        >
          <el-icon><Select /></el-icon>
          选择
        </el-button>
        <el-button
          :type="currentTool === 'dimension' ? 'primary' : 'default'"
          @click="setTool('dimension')"
        >
          <el-icon>
            <Ruler />
          </el-icon>
          尺寸标注
        </el-button>
        <el-button
          :type="currentTool === 'material' ? 'primary' : 'default'"
          @click="setTool('material')"
        >
          <el-icon>
            <Grid />
          </el-icon>
          材料标注
        </el-button>
        <el-button
          :type="currentTool === 'process' ? 'primary' : 'default'"
          @click="setTool('process')"
        >
          <el-icon>
            <Setting />
          </el-icon>
          工艺说明
        </el-button>
      </el-button-group>

      <div class="toolbar-actions">
        <el-button @click="handleAutoAnnotate" :loading="analyzing">
          <el-icon>
            <MagicStick />
          </el-icon>
          自动标注
        </el-button>
        <el-button @click="handleExport">
          <el-icon>
            <Download />
          </el-icon>
          导出
        </el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          <el-icon>
            <Check />
          </el-icon>
          保存
        </el-button>
      </div>
    </div>

    <!-- 主编辑区域 / Main editing area -->
    <div class="editor-main">
      <!-- 画布区域 / Canvas area -->
      <div class="canvas-container" ref="canvasContainer">
        <div
          class="canvas-wrapper"
          :style="canvasStyle"
          @mousedown="handleCanvasMouseDown"
          @mousemove="handleCanvasMouseMove"
          @mouseup="handleCanvasMouseUp"
        >
          <!-- 背景图片 / Background image -->
          <img v-if="imagePath" :src="imagePath" class="background-image" @load="handleImageLoad" />

          <!-- 标注层 / Annotation layer -->
          <svg class="annotation-layer" :viewBox="svgViewBox">
            <!-- 渲染标注项 / Render annotation items -->
            <g v-for="item in annotations" :key="item.id">
              <AnnotationItemRenderer
                :item="item"
                :selected="selectedItemId === item.id"
                @select="handleSelectItem"
                @update="handleUpdateItem"
                @delete="handleDeleteItem"
              />
            </g>

            <!-- 正在绘制的标注 / Currently drawing annotation -->
            <g v-if="drawingItem">
              <AnnotationItemRenderer :item="drawingItem" :drawing="true" />
            </g>
          </svg>
        </div>
      </div>

      <!-- 属性面板 / Properties panel -->
      <div class="properties-panel">
        <el-card v-if="selectedItem" class="property-card">
          <template #header>
            <div class="card-header">
              <span>标注属性</span>
              <el-button type="danger" size="small" @click="handleDeleteSelected">
                <el-icon>
                  <Delete />
                </el-icon>
              </el-button>
            </div>
          </template>

          <el-form label-width="80px" size="small">
            <el-form-item label="类型">
              <el-select v-model="selectedItem.type" @change="handlePropertyChange">
                <el-option label="尺寸标注" value="dimension" />
                <el-option label="材料标注" value="material" />
                <el-option label="工艺说明" value="process" />
              </el-select>
            </el-form-item>

            <el-form-item label="内容">
              <el-input v-model="selectedItem.content" @change="handlePropertyChange" />
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
                :max="48"
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
          </el-form>
        </el-card>

        <!-- 标注列表 / Annotation list -->
        <el-card class="annotation-list-card">
          <template #header>
            <span>标注列表 ({{ annotations.length }})</span>
          </template>

          <el-scrollbar height="300px">
            <div
              v-for="item in annotations"
              :key="item.id"
              class="annotation-list-item"
              :class="{ active: selectedItemId === item.id }"
              @click="handleSelectItem(item.id)"
            >
              <el-icon>
                <Ruler v-if="item.type === 'dimension'" />
                <Grid v-else-if="item.type === 'material'" />
                <Setting v-else />
              </el-icon>
              <span class="item-content">{{ item.content }}</span>
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
  import { ref, computed, onMounted, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    Select,
    Ruler,
    Grid,
    Setting,
    MagicStick,
    Download,
    Check,
    Delete
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

  /*** Component Props ***/
  const props = defineProps<{
    annotationId: number
    projectId: number
    cadFileId: number
  }>()

  /*** Component Emits ***/
  const emit = defineEmits<{
    (e: 'saved'): void
    (e: 'exported', result: any): void
  }>()

  /*** Reactive State ***/
  const currentTool = ref<string>('select')
  const annotations = ref<AnnotationItem[]>([])
  const selectedItemId = ref<string | null>(null)
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
  const imageWidth = ref(800)
  const imageHeight = ref(600)
  const drawingItem = ref<AnnotationItem | null>(null)
  const isDrawing = ref(false)
  const startPoint = ref({ x: 0, y: 0 })

  /*** Computed Properties ***/
  const selectedItem = computed(() => {
    if (!selectedItemId.value) return null
    return annotations.value.find((item) => item.id === selectedItemId.value) || null
  })

  const canvasStyle = computed(() => ({
    width: `${imageWidth.value}px`,
    height: `${imageHeight.value}px`
  }))

  const svgViewBox = computed(() => `0 0 ${imageWidth.value} ${imageHeight.value}`)

  /*** Methods ***/
  // 设置当前工具 / Set current tool
  const setTool = (tool: string) => {
    currentTool.value = tool
    selectedItemId.value = null
  }

  // 加载标注数据 / Load annotation data
  const loadAnnotation = async () => {
    try {
      const res = await getAnnotation(props.annotationId)
      if (res.data) {
        const data = res.data as ConstructionAnnotationResponse
        annotations.value = data.annotations || []
        imagePath.value = data.imagePath
      }
    } catch (error) {
      ElMessage.error('加载标注数据失败')
    }
  }

  // 加载导出格式 / Load export formats
  const loadExportFormats = async () => {
    try {
      const res = await getExportFormats()
      if (res.data) {
        exportFormats.value = res.data as ExportFormat[]
      }
    } catch (error) {
      console.error('加载导出格式失败:', error)
    }
  }

  // 处理图片加载 / Handle image load
  const handleImageLoad = (event: Event) => {
    const img = event.target as HTMLImageElement
    imageWidth.value = img.naturalWidth
    imageHeight.value = img.naturalHeight
  }

  // 处理画布鼠标按下 / Handle canvas mouse down
  const handleCanvasMouseDown = (event: MouseEvent) => {
    if (currentTool.value === 'select') return

    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
    startPoint.value = {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top
    }
    isDrawing.value = true

    // 创建新标注项 / Create new annotation item
    drawingItem.value = {
      id: `temp_${Date.now()}`,
      type: currentTool.value,
      position: { x: startPoint.value.x, y: startPoint.value.y, anchor: 'top-left' },
      content: getDefaultContent(currentTool.value),
      style: getDefaultStyle(currentTool.value),
      properties: {},
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      isGenerated: false
    }
  }

  // 处理画布鼠标移动 / Handle canvas mouse move
  const handleCanvasMouseMove = (event: MouseEvent) => {
    if (!isDrawing.value || !drawingItem.value) return

    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
    const currentX = event.clientX - rect.left
    const currentY = event.clientY - rect.top

    drawingItem.value.position.x = currentX
    drawingItem.value.position.y = currentY
  }

  // 处理画布鼠标抬起 / Handle canvas mouse up
  const handleCanvasMouseUp = async () => {
    if (!isDrawing.value || !drawingItem.value) return

    isDrawing.value = false

    // 添加标注项 / Add annotation item
    try {
      await addAnnotationItem(props.annotationId, drawingItem.value)
      annotations.value.push({ ...drawingItem.value })
      ElMessage.success('标注添加成功')
    } catch (error) {
      ElMessage.error('添加标注失败')
    }

    drawingItem.value = null
    setTool('select')
  }

  // 获取默认内容 / Get default content
  const getDefaultContent = (type: string): string => {
    switch (type) {
      case 'dimension':
        return '尺寸'
      case 'material':
        return '材料'
      case 'process':
        return '工艺说明'
      default:
        return '标注'
    }
  }

  // 获取默认样式 / Get default style
  const getDefaultStyle = (type: string) => {
    switch (type) {
      case 'dimension':
        return {
          fontSize: 12,
          fontColor: '#FF0000',
          lineColor: '#FF0000',
          lineWidth: 1,
          background: '#FFFFFF'
        }
      case 'material':
        return {
          fontSize: 10,
          fontColor: '#0000FF',
          lineColor: '#0000FF',
          lineWidth: 1,
          background: '#FFFFFF'
        }
      case 'process':
        return {
          fontSize: 10,
          fontColor: '#00AA00',
          lineColor: '#00AA00',
          lineWidth: 1,
          background: '#FFFFFF'
        }
      default:
        return {
          fontSize: 10,
          fontColor: '#000000',
          lineColor: '#000000',
          lineWidth: 1,
          background: '#FFFFFF'
        }
    }
  }

  // 选择标注项 / Select annotation item
  const handleSelectItem = (id: string) => {
    selectedItemId.value = id
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
      await deleteAnnotationItem(props.annotationId, id)
      annotations.value = annotations.value.filter((item) => item.id !== id)
      if (selectedItemId.value === id) {
        selectedItemId.value = null
      }
      ElMessage.success('标注删除成功')
    } catch (error) {
      ElMessage.error('删除标注失败')
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
        projectId: props.projectId,
        cadFileId: props.cadFileId,
        imagePath: imagePath.value
      })
      ElMessage.success('自动标注已开始，请稍后刷新查看结果')
      // 延迟重新加载 / Delay reload
      setTimeout(() => loadAnnotation(), 3000)
    } catch (error) {
      ElMessage.error('自动标注失败')
    } finally {
      analyzing.value = false
    }
  }

  // 保存标注 / Save annotations
  const handleSave = async () => {
    saving.value = true
    try {
      await updateAnnotations(props.annotationId, annotations.value)
      ElMessage.success('保存成功')
      emit('saved')
    } catch (error) {
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
      const res = await exportAnnotation(props.annotationId, exportConfig.value)
      ElMessage.success('导出成功')
      exportDialogVisible.value = false
      emit('exported', res.data)
    } catch (error) {
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

  watch(
    () => props.annotationId,
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
  }

  .toolbar-actions {
    display: flex;
    gap: 8px;
  }

  .editor-main {
    display: flex;
    flex: 1;
    overflow: hidden;
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
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .background-image {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: contain;
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

  .properties-panel {
    width: 300px;
    padding: 16px;
    background: #fff;
    border-left: 1px solid #e4e7ed;
    overflow-y: auto;
  }

  .property-card {
    margin-bottom: 16px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .annotation-list-card {
    margin-top: 16px;
  }

  .annotation-list-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    cursor: pointer;
    border-radius: 4px;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: #ecf5ff;
    }

    .item-content {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
</style>
