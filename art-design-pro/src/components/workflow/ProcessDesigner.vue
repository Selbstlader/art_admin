<template>
  <div class="process-designer">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-button-group>
          <el-button size="small" :icon="ArrowLeft" :disabled="!canUndo" @click="undo">
            撤销
          </el-button>
          <el-button size="small" :icon="ArrowRight" :disabled="!canRedo" @click="redo">
            重做
          </el-button>
        </el-button-group>
        <el-divider direction="vertical" />
        <el-button size="small" :icon="ZoomIn" @click="zoomIn">放大</el-button>
        <el-button size="small" :icon="ZoomOut" @click="zoomOut">缩小</el-button>
        <el-button size="small" :icon="FullScreen" @click="fitView">适应</el-button>
      </div>
      <div class="toolbar-right">
        <el-button size="small" :icon="View" @click="handlePreview">预览</el-button>
        <el-button size="small" type="primary" :icon="Check" @click="handleSave">保存</el-button>
        <el-button size="small" type="success" :icon="Upload" @click="handlePublish">发布</el-button>
      </div>
    </div>

    <div class="designer-content">
      <!-- 左侧：节点面板 -->
      <div class="node-panel">
        <div class="panel-header">节点面板</div>
        <div class="node-list">
          <div
            v-for="node in nodePalette"
            :key="node.type"
            class="node-item"
            draggable="true"
            @dragstart="onDragStart($event, node)"
          >
            <div :class="['node-icon', `node-icon-${node.type}`]">
              <el-icon :size="20">
                <component :is="node.icon" />
              </el-icon>
            </div>
            <span class="node-label">{{ node.label }}</span>
          </div>
        </div>
      </div>

      <!-- 中间：画布区域 -->
      <div class="canvas-panel" @drop="onDrop" @dragover.prevent>
        <VueFlow
          ref="vueFlowRef"
          v-model:nodes="nodes"
          v-model:edges="edges"
          :default-viewport="{ zoom: 1 }"
          :min-zoom="0.2"
          :max-zoom="4"
          :snap-to-grid="true"
          :snap-grid="[15, 15]"
          fit-view-on-init
          @node-click="onNodeClick"
          @pane-click="onPaneClick"
          @connect="onConnect"
          @nodes-change="onNodesChange"
          @edges-change="onEdgesChange"
        >
          <!-- 自定义节点 -->
          <template #node-start="nodeProps">
            <StartNode v-bind="nodeProps" :selected="selectedNodeId === nodeProps.id" />
          </template>
          <template #node-end="nodeProps">
            <EndNode v-bind="nodeProps" :selected="selectedNodeId === nodeProps.id" />
          </template>
          <template #node-approval="nodeProps">
            <ApprovalNode v-bind="nodeProps" :selected="selectedNodeId === nodeProps.id" />
          </template>
          <template #node-condition="nodeProps">
            <ConditionNode v-bind="nodeProps" :selected="selectedNodeId === nodeProps.id" />
          </template>

          <!-- 背景 -->
          <Background pattern-color="#aaa" :gap="15" />
          
          <!-- 控制器 -->
          <Controls />
          
          <!-- 小地图 -->
          <MiniMap />
        </VueFlow>
      </div>

      <!-- 右侧：属性面板 -->
      <div class="property-panel">
        <div class="panel-header">属性配置</div>
        <div v-if="selectedNode" class="property-form">
          <NodePropertyPanel
            :node="selectedNode"
            @update="onNodePropertyUpdate"
          />
        </div>
        <div v-else class="no-selection">
          <el-icon :size="48"><Setting /></el-icon>
          <p>请选择一个节点进行配置</p>
        </div>
      </div>
    </div>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="流程预览 (JSON)" width="800px">
      <el-input
        type="textarea"
        :model-value="previewJson"
        :rows="20"
        readonly
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import type { Node, Edge, Connection } from '@vue-flow/core'
import {
  ArrowLeft,
  ArrowRight,
  ZoomIn,
  ZoomOut,
  FullScreen,
  View,
  Check,
  Upload,
  Setting,
  VideoPlay,
  CircleCheck,
  Switch,
  CircleClose
} from '@element-plus/icons-vue'
import StartNode from './nodes/StartNode.vue'
import EndNode from './nodes/EndNode.vue'
import ApprovalNode from './nodes/ApprovalNode.vue'
import ConditionNode from './nodes/ConditionNode.vue'
import NodePropertyPanel from './NodePropertyPanel.vue'
import type { ProcessGraph, ProcessNode, ProcessEdge, NodeProperties } from '@/api/workflow'

defineOptions({ name: 'ProcessDesigner' })

// Props
interface Props {
  modelValue?: ProcessGraph
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({ nodes: [], edges: [] }),
  readonly: false
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: ProcessGraph]
  save: [value: ProcessGraph]
  publish: [value: ProcessGraph]
}>()

// Vue Flow instance
const vueFlowRef = ref()
const { 
  zoomIn: vfZoomIn, 
  zoomOut: vfZoomOut, 
  fitView,
  addNodes,
  addEdges,
  project
} = useVueFlow()

// 节点面板配置
interface NodePaletteItem {
  type: string
  label: string
  icon: any
}

const nodePalette: NodePaletteItem[] = [
  { type: 'start', label: '开始节点', icon: VideoPlay },
  { type: 'approval', label: '审批节点', icon: CircleCheck },
  { type: 'condition', label: '条件节点', icon: Switch },
  { type: 'end', label: '结束节点', icon: CircleClose }
]

// 节点和边数据
const nodes = ref<Node[]>([])
const edges = ref<Edge[]>([])

// 选中的节点
const selectedNodeId = ref<string | null>(null)
const selectedNode = computed(() => {
  if (!selectedNodeId.value) return null
  return nodes.value.find(n => n.id === selectedNodeId.value) || null
})

// 历史记录（用于撤销/重做）
const history = ref<{ nodes: Node[]; edges: Edge[] }[]>([])
const historyIndex = ref(-1)
const canUndo = computed(() => historyIndex.value > 0)
const canRedo = computed(() => historyIndex.value < history.value.length - 1)

// 预览
const previewVisible = ref(false)
const previewJson = computed(() => {
  const graph = toProcessGraph()
  return JSON.stringify(graph, null, 2)
})

// 初始化
onMounted(() => {
  if (props.modelValue?.nodes?.length) {
    loadFromProcessGraph(props.modelValue)
  }
  saveHistory()
})

// 监听外部数据变化
watch(
  () => props.modelValue,
  (val) => {
    if (val?.nodes?.length) {
      loadFromProcessGraph(val)
    }
  },
  { deep: true }
)

// 从 ProcessGraph 加载数据
const loadFromProcessGraph = (graph: ProcessGraph) => {
  nodes.value = graph.nodes.map(n => ({
    id: n.id,
    type: n.type,
    position: n.position,
    data: {
      label: n.name,
      properties: n.properties || {}
    }
  }))
  
  edges.value = graph.edges.map(e => ({
    id: e.id,
    source: e.source,
    target: e.target,
    label: e.condition || '',
    data: {
      condition: e.condition
    }
  }))
}

// 转换为 ProcessGraph
const toProcessGraph = (): ProcessGraph => {
  const processNodes: ProcessNode[] = nodes.value.map(n => ({
    id: n.id,
    type: n.type || 'approval',
    name: n.data?.label || '',
    position: n.position,
    properties: n.data?.properties || {}
  }))
  
  const processEdges: ProcessEdge[] = edges.value.map(e => ({
    id: e.id,
    source: e.source,
    target: e.target,
    condition: e.data?.condition
  }))
  
  return {
    nodes: processNodes,
    edges: processEdges
  }
}

// 生成唯一ID
const generateId = (type: string): string => {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 6)
  return `${type}_${timestamp}_${random}`
}

// 拖拽开始
const onDragStart = (event: DragEvent, node: NodePaletteItem) => {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/vueflow', node.type)
    event.dataTransfer.effectAllowed = 'move'
  }
}

// 拖拽放置
const onDrop = (event: DragEvent) => {
  const type = event.dataTransfer?.getData('application/vueflow')
  if (!type) return

  // 检查开始节点是否已存在
  if (type === 'start' && nodes.value.some(n => n.type === 'start')) {
    ElMessage.warning('流程只能有一个开始节点')
    return
  }

  const { left, top } = (event.target as HTMLElement).getBoundingClientRect()
  const position = project({
    x: event.clientX - left,
    y: event.clientY - top
  })

  const newNode: Node = {
    id: generateId(type),
    type,
    position,
    data: {
      label: getDefaultLabel(type),
      properties: getDefaultProperties(type)
    }
  }

  addNodes([newNode])
  selectedNodeId.value = newNode.id
  saveHistory()
}

// 获取默认标签
const getDefaultLabel = (type: string): string => {
  const labels: Record<string, string> = {
    start: '开始',
    end: '结束',
    approval: '审批节点',
    condition: '条件分支'
  }
  return labels[type] || '节点'
}

// 获取默认属性
const getDefaultProperties = (type: string): NodeProperties => {
  if (type === 'approval') {
    return {
      assigneeRule: { type: 'user', values: [] },
      approvalMode: 'or_sign',
      timeoutHours: 24,
      rejectAction: 'terminate'
    }
  }
  return {}
}

// 节点点击
const onNodeClick = (_event: MouseEvent, node: Node) => {
  selectedNodeId.value = node.id
}

// 画布点击（取消选中）
const onPaneClick = () => {
  selectedNodeId.value = null
}

// 连接节点
const onConnect = (connection: Connection) => {
  const newEdge: Edge = {
    id: generateId('edge'),
    source: connection.source,
    target: connection.target,
    sourceHandle: connection.sourceHandle || undefined,
    targetHandle: connection.targetHandle || undefined
  }
  addEdges([newEdge])
  saveHistory()
}

// 节点变化
const onNodesChange = () => {
  emitUpdate()
}

// 边变化
const onEdgesChange = () => {
  emitUpdate()
}

// 更新节点属性
const onNodePropertyUpdate = (nodeId: string, data: any) => {
  const node = nodes.value.find(n => n.id === nodeId)
  if (node) {
    node.data = { ...node.data, ...data }
    saveHistory()
    emitUpdate()
  }
}

// 发送更新事件
const emitUpdate = () => {
  const graph = toProcessGraph()
  emit('update:modelValue', graph)
}

// 保存历史
const saveHistory = () => {
  // 删除当前位置之后的历史
  history.value = history.value.slice(0, historyIndex.value + 1)
  // 添加新历史
  history.value.push({
    nodes: JSON.parse(JSON.stringify(nodes.value)),
    edges: JSON.parse(JSON.stringify(edges.value))
  })
  historyIndex.value = history.value.length - 1
  
  // 限制历史记录数量
  if (history.value.length > 50) {
    history.value.shift()
    historyIndex.value--
  }
}

// 撤销
const undo = () => {
  if (canUndo.value) {
    historyIndex.value--
    const state = history.value[historyIndex.value]
    nodes.value = JSON.parse(JSON.stringify(state.nodes))
    edges.value = JSON.parse(JSON.stringify(state.edges))
    emitUpdate()
  }
}

// 重做
const redo = () => {
  if (canRedo.value) {
    historyIndex.value++
    const state = history.value[historyIndex.value]
    nodes.value = JSON.parse(JSON.stringify(state.nodes))
    edges.value = JSON.parse(JSON.stringify(state.edges))
    emitUpdate()
  }
}

// 缩放
const zoomIn = () => vfZoomIn()
const zoomOut = () => vfZoomOut()

// 预览
const handlePreview = () => {
  previewVisible.value = true
}

// 验证流程结构
interface ValidationResult {
  valid: boolean
  errors: string[]
}

const validateGraph = (): ValidationResult => {
  const errors: string[] = []
  const graph = toProcessGraph()
  
  // 检查是否有节点
  if (graph.nodes.length === 0) {
    errors.push('流程至少需要一个节点')
    return { valid: false, errors }
  }
  
  // 检查开始节点
  const startNodes = graph.nodes.filter(n => n.type === 'start')
  if (startNodes.length === 0) {
    errors.push('流程必须有一个开始节点')
  } else if (startNodes.length > 1) {
    errors.push('流程只能有一个开始节点')
  }
  
  // 检查结束节点
  const endNodes = graph.nodes.filter(n => n.type === 'end')
  if (endNodes.length === 0) {
    errors.push('流程必须有至少一个结束节点')
  }
  
  // 检查审批节点配置
  const approvalNodes = graph.nodes.filter(n => n.type === 'approval')
  for (const node of approvalNodes) {
    const rule = node.properties?.assigneeRule
    if (!rule || !rule.type) {
      errors.push(`审批节点"${node.name}"未配置审批人规则`)
    } else if ((rule.type === 'user' || rule.type === 'role') && (!rule.values || rule.values.length === 0)) {
      errors.push(`审批节点"${node.name}"未选择审批人`)
    }
  }
  
  // 检查节点连接性
  const nodeIds = new Set(graph.nodes.map(n => n.id))
  const connectedNodes = new Set<string>()
  
  // 从开始节点开始遍历
  if (startNodes.length > 0) {
    const queue = [startNodes[0].id]
    while (queue.length > 0) {
      const currentId = queue.shift()!
      if (connectedNodes.has(currentId)) continue
      connectedNodes.add(currentId)
      
      // 找到所有从当前节点出发的边
      const outEdges = graph.edges.filter(e => e.source === currentId)
      for (const edge of outEdges) {
        if (nodeIds.has(edge.target) && !connectedNodes.has(edge.target)) {
          queue.push(edge.target)
        }
      }
    }
    
    // 检查是否有孤立节点
    const isolatedNodes = graph.nodes.filter(n => !connectedNodes.has(n.id) && n.type !== 'start')
    if (isolatedNodes.length > 0) {
      errors.push(`存在未连接的节点: ${isolatedNodes.map(n => n.name).join(', ')}`)
    }
  }
  
  // 检查是否能到达结束节点
  const reachableEndNodes = endNodes.filter(n => connectedNodes.has(n.id))
  if (endNodes.length > 0 && reachableEndNodes.length === 0) {
    errors.push('流程无法到达结束节点')
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

// 保存
const handleSave = () => {
  const graph = toProcessGraph()
  emit('save', graph)
}

// 发布
const handlePublish = () => {
  // 发布前验证
  const validation = validateGraph()
  if (!validation.valid) {
    ElMessageBox.alert(
      `<ul style="margin: 0; padding-left: 20px;">${validation.errors.map(e => `<li>${e}</li>`).join('')}</ul>`,
      '流程验证失败',
      {
        dangerouslyUseHTMLString: true,
        type: 'error'
      }
    )
    return
  }
  
  const graph = toProcessGraph()
  emit('publish', graph)
}

// 暴露方法
defineExpose({
  getGraph: toProcessGraph,
  setGraph: loadFromProcessGraph,
  validate: validateGraph,
  clear: () => {
    nodes.value = []
    edges.value = []
    selectedNodeId.value = null
    saveHistory()
  }
})
</script>


<style lang="scss" scoped>
/* Import Vue Flow styles */
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';
@import '@vue-flow/minimap/dist/style.css';

.process-designer {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 600px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;

  // 工具栏
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    padding: 0 16px;
    background: var(--el-fill-color-light);
    border-bottom: 1px solid var(--el-border-color-light);

    .toolbar-left,
    .toolbar-right {
      display: flex;
      gap: 8px;
      align-items: center;
    }
  }

  // 设计器内容区
  .designer-content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  // 左侧节点面板
  .node-panel {
    width: 180px;
    border-right: 1px solid var(--el-border-color-light);
    background: var(--el-bg-color);

    .panel-header {
      height: 40px;
      padding: 0 12px;
      font-size: 14px;
      font-weight: 500;
      line-height: 40px;
      color: var(--el-text-color-primary);
      background: var(--el-fill-color-lighter);
      border-bottom: 1px solid var(--el-border-color-light);
    }

    .node-list {
      padding: 12px;
    }

    .node-item {
      display: flex;
      gap: 10px;
      align-items: center;
      padding: 10px 12px;
      margin-bottom: 8px;
      cursor: grab;
      background: var(--el-fill-color-lighter);
      border: 1px solid var(--el-border-color-lighter);
      border-radius: 6px;
      transition: all 0.2s;

      &:hover {
        background: var(--el-color-primary-light-9);
        border-color: var(--el-color-primary-light-5);
      }

      &:active {
        cursor: grabbing;
      }

      .node-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        border-radius: 50%;

        &.node-icon-start {
          color: #fff;
          background: #67c23a;
        }

        &.node-icon-end {
          color: #fff;
          background: #f56c6c;
        }

        &.node-icon-approval {
          color: #fff;
          background: #409eff;
        }

        &.node-icon-condition {
          color: #fff;
          background: #e6a23c;
        }
      }

      .node-label {
        font-size: 13px;
        color: var(--el-text-color-regular);
      }
    }
  }

  // 中间画布区域
  .canvas-panel {
    flex: 1;
    position: relative;
    background: var(--el-fill-color-lighter);

    .vue-flow {
      width: 100%;
      height: 100%;
    }
  }

  // 右侧属性面板
  .property-panel {
    width: 300px;
    border-left: 1px solid var(--el-border-color-light);
    background: var(--el-bg-color);

    .panel-header {
      height: 40px;
      padding: 0 12px;
      font-size: 14px;
      font-weight: 500;
      line-height: 40px;
      color: var(--el-text-color-primary);
      background: var(--el-fill-color-lighter);
      border-bottom: 1px solid var(--el-border-color-light);
    }

    .property-form {
      padding: 12px;
      overflow-y: auto;
      max-height: calc(100% - 40px);
    }

    .no-selection {
      display: flex;
      flex-direction: column;
      gap: 12px;
      align-items: center;
      justify-content: center;
      height: calc(100% - 40px);
      color: var(--el-text-color-secondary);

      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }
}

// Vue Flow 自定义样式
:deep(.vue-flow__node) {
  cursor: pointer;
}

:deep(.vue-flow__edge-path) {
  stroke: var(--el-border-color);
  stroke-width: 2;
}

:deep(.vue-flow__edge.selected .vue-flow__edge-path) {
  stroke: var(--el-color-primary);
}

:deep(.vue-flow__handle) {
  width: 10px;
  height: 10px;
  background: var(--el-color-primary);
  border: 2px solid #fff;
}

:deep(.vue-flow__minimap) {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
}

:deep(.vue-flow__controls) {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  box-shadow: none;
}
</style>
