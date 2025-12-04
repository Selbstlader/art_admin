<template>
  <div class="amap-container" ref="containerRef">
    <div class="amap-wrapper" ref="mapRef"></div>
    <!-- 加载状态 -->
    <div v-if="loading" class="amap-loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>地图加载中...</span>
    </div>
    <!-- 错误状态 -->
    <div v-if="error" class="amap-error">
      <el-icon><WarningFilled /></el-icon>
      <span>{{ error }}</span>
      <el-button type="primary" size="small" @click="retry">重试</el-button>
    </div>
    <!-- 地图控件插槽 -->
    <div class="amap-controls">
      <slot name="controls"></slot>
    </div>
    <!-- 信息窗体插槽 -->
    <slot name="infoWindow"></slot>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Loading, WarningFilled } from '@element-plus/icons-vue'
import { loadAMapScript, type AMapInstance } from '../composables/useAMap'

defineOptions({
  name: 'AMapContainer'
})

// Props
interface Props {
  /** 地图中心点 [经度, 纬度] */
  center?: [number, number]
  /** 缩放级别 */
  zoom?: number
  /** 是否显示缩放控件 */
  showZoomControl?: boolean
  /** 是否显示比例尺 */
  showScale?: boolean
  /** 是否显示定位按钮 */
  showGeolocation?: boolean
  /** 地图样式 */
  mapStyle?: string
  /** 是否启用3D视图 */
  viewMode?: '2D' | '3D'
  /** 俯仰角度 */
  pitch?: number
}

const props = withDefaults(defineProps<Props>(), {
  center: () => [116.397428, 39.90923],
  zoom: 12,
  showZoomControl: true,
  showScale: true,
  showGeolocation: false,
  mapStyle: 'amap://styles/normal',
  viewMode: '2D',
  pitch: 0
})

// Emits
const emit = defineEmits<{
  (e: 'ready', map: AMapInstance): void
  (e: 'click', event: any): void
  (e: 'moveend', center: [number, number]): void
  (e: 'zoomend', zoom: number): void
  (e: 'error', error: string): void
}>()

// Refs
const containerRef = ref<HTMLElement | null>(null)
const mapRef = ref<HTMLElement | null>(null)
const mapInstance = ref<AMapInstance | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

// 初始化地图
async function initMap() {
  if (!mapRef.value) return

  loading.value = true
  error.value = null

  try {
    // 加载高德地图脚本
    const AMap = await loadAMapScript()

    // 创建地图实例
    const map = new AMap.Map(mapRef.value, {
      center: props.center,
      zoom: props.zoom,
      mapStyle: props.mapStyle,
      viewMode: props.viewMode,
      pitch: props.pitch,
      resizeEnable: true,
      rotateEnable: true,
      pitchEnable: props.viewMode === '3D'
    })

    mapInstance.value = map

    // 添加控件
    if (props.showZoomControl) {
      AMap.plugin('AMap.ToolBar', () => {
        const toolbar = new AMap.ToolBar({
          position: 'RT'
        })
        map.addControl(toolbar)
      })
    }

    if (props.showScale) {
      AMap.plugin('AMap.Scale', () => {
        const scale = new AMap.Scale()
        map.addControl(scale)
      })
    }

    if (props.showGeolocation) {
      AMap.plugin('AMap.Geolocation', () => {
        const geolocation = new AMap.Geolocation({
          enableHighAccuracy: true,
          timeout: 10000,
          buttonPosition: 'RB'
        })
        map.addControl(geolocation)
      })
    }

    // 绑定事件
    map.on('click', (e: any) => {
      emit('click', {
        lnglat: [e.lnglat.getLng(), e.lnglat.getLat()],
        pixel: [e.pixel.x, e.pixel.y],
        target: e.target
      })
    })

    map.on('moveend', () => {
      const center = map.getCenter()
      emit('moveend', [center.getLng(), center.getLat()])
    })

    map.on('zoomend', () => {
      emit('zoomend', map.getZoom())
    })

    // 地图加载完成
    map.on('complete', () => {
      loading.value = false
      emit('ready', map)
    })
  } catch (err: any) {
    loading.value = false
    error.value = err.message || '地图加载失败'
    emit('error', error.value as string)
  }
}

// 重试加载
function retry() {
  initMap()
}

// 暴露地图实例和方法
defineExpose({
  /** 获取地图实例 */
  getMap: () => mapInstance.value,
  /** 设置中心点 */
  setCenter: (center: [number, number]) => {
    mapInstance.value?.setCenter(center)
  },
  /** 设置缩放级别 */
  setZoom: (zoom: number) => {
    mapInstance.value?.setZoom(zoom)
  },
  /** 平移到指定位置 */
  panTo: (position: [number, number]) => {
    mapInstance.value?.panTo(position)
  },
  /** 自适应显示 */
  setFitView: (overlays?: any[]) => {
    mapInstance.value?.setFitView(overlays)
  },
  /** 销毁地图 */
  destroy: () => {
    mapInstance.value?.destroy()
    mapInstance.value = null
  }
})

// 监听props变化
watch(() => props.center, (newCenter) => {
  if (mapInstance.value && newCenter) {
    mapInstance.value.setCenter(newCenter)
  }
})

watch(() => props.zoom, (newZoom) => {
  if (mapInstance.value && newZoom) {
    mapInstance.value.setZoom(newZoom)
  }
})

// 生命周期
onMounted(() => {
  nextTick(() => {
    initMap()
  })
})

onUnmounted(() => {
  if (mapInstance.value) {
    mapInstance.value.destroy()
    mapInstance.value = null
  }
})
</script>

<style scoped lang="scss">
.amap-container {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 400px;

  .amap-wrapper {
    width: 100%;
    height: 100%;
  }

  .amap-loading,
  .amap-error {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.9);
    z-index: 100;
    gap: 12px;

    .el-icon {
      font-size: 32px;
      color: var(--el-color-primary);
    }

    span {
      font-size: 14px;
      color: var(--el-text-color-secondary);
    }
  }

  .amap-error {
    .el-icon {
      color: var(--el-color-warning);
    }
  }

  .amap-controls {
    position: absolute;
    top: 10px;
    left: 10px;
    z-index: 10;
  }
}

// 暗色模式适配
:deep(.amap-logo),
:deep(.amap-copyright) {
  display: none !important;
}
</style>
