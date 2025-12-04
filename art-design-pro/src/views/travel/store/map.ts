/**
 * 地图状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  Waypoint,
  TravelMode,
  POISearchResult,
  RouteResult,
  MapMarker
} from '../types'

export const useMapStore = defineStore('travelMap', () => {
  // ==================== 状态 ====================

  /** 地图实例 */
  const mapInstance = ref<any>(null)

  /** 地图中心点 */
  const center = ref<[number, number]>([116.397428, 39.90923])

  /** 地图缩放级别 */
  const zoom = ref(12)

  /** 当前出行方式 */
  const travelMode = ref<TravelMode>('driving' as TravelMode)

  /** 地图标记列表 */
  const markers = ref<MapMarker[]>([])

  /** 当前选中的标记 */
  const selectedMarker = ref<MapMarker | null>(null)

  /** POI搜索结果 */
  const searchResults = ref<POISearchResult[]>([])

  /** 搜索关键词 */
  const searchKeyword = ref('')

  /** 搜索加载状态 */
  const searchLoading = ref(false)

  /** 路线规划结果 */
  const routeResult = ref<RouteResult | null>(null)

  /** 路线规划加载状态 */
  const routeLoading = ref(false)

  /** 是否显示路线 */
  const showRoute = ref(true)

  /** 是否显示标记 */
  const showMarkers = ref(true)

  /** 当前编辑的途经点列表 */
  const editingWaypoints = ref<Waypoint[]>([])

  /** 是否处于编辑模式 */
  const isEditing = ref(false)

  /** 地图是否已加载 */
  const mapLoaded = ref(false)

  // ==================== 计算属性 ====================

  /** 标记数量 */
  const markerCount = computed(() => markers.value.length)

  /** 是否有路线 */
  const hasRoute = computed(() => routeResult.value !== null)

  /** 路线总距离（公里） */
  const totalDistance = computed(() => {
    if (!routeResult.value) return 0
    return Math.round(routeResult.value.distance / 100) / 10
  })

  /** 路线总时长（分钟） */
  const totalDuration = computed(() => {
    if (!routeResult.value) return 0
    return Math.round(routeResult.value.duration / 60)
  })

  /** 格式化的路线时长 */
  const formattedDuration = computed(() => {
    const minutes = totalDuration.value
    if (minutes < 60) {
      return `${minutes}分钟`
    }
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
  })

  // ==================== 地图操作 ====================

  /** 设置地图实例 */
  function setMapInstance(instance: any) {
    mapInstance.value = instance
    mapLoaded.value = true
  }

  /** 设置地图中心点 */
  function setCenter(lng: number, lat: number) {
    center.value = [lng, lat]
    if (mapInstance.value) {
      mapInstance.value.setCenter([lng, lat])
    }
  }

  /** 设置地图缩放级别 */
  function setZoom(level: number) {
    zoom.value = level
    if (mapInstance.value) {
      mapInstance.value.setZoom(level)
    }
  }

  /** 设置出行方式 */
  function setTravelMode(mode: TravelMode) {
    travelMode.value = mode
  }

  /** 移动到指定位置 */
  function panTo(lng: number, lat: number) {
    if (mapInstance.value) {
      mapInstance.value.panTo([lng, lat])
    }
    center.value = [lng, lat]
  }

  /** 自适应显示所有标记 */
  function fitBounds() {
    if (!mapInstance.value || markers.value.length === 0) return
    const positions = markers.value.map((m) => m.position)
    mapInstance.value.setFitView(positions)
  }

  // ==================== 标记操作 ====================

  /** 添加标记 */
  function addMarker(marker: MapMarker) {
    markers.value.push(marker)
  }

  /** 批量设置标记 */
  function setMarkers(newMarkers: MapMarker[]) {
    markers.value = newMarkers
  }

  /** 更新标记 */
  function updateMarker(id: string | number, data: Partial<MapMarker>) {
    const index = markers.value.findIndex((m) => m.id === id)
    if (index !== -1) {
      markers.value[index] = { ...markers.value[index], ...data }
    }
  }

  /** 删除标记 */
  function removeMarker(id: string | number) {
    markers.value = markers.value.filter((m) => m.id !== id)
    if (selectedMarker.value?.id === id) {
      selectedMarker.value = null
    }
  }

  /** 清空所有标记 */
  function clearMarkers() {
    markers.value = []
    selectedMarker.value = null
  }

  /** 选中标记 */
  function selectMarker(marker: MapMarker | null) {
    selectedMarker.value = marker
    if (marker) {
      panTo(marker.position[0], marker.position[1])
    }
  }

  // ==================== 搜索操作 ====================

  /** 设置搜索关键词 */
  function setSearchKeyword(keyword: string) {
    searchKeyword.value = keyword
  }

  /** 设置搜索结果 */
  function setSearchResults(results: POISearchResult[]) {
    searchResults.value = results
  }

  /** 清空搜索结果 */
  function clearSearchResults() {
    searchResults.value = []
    searchKeyword.value = ''
  }

  /** 设置搜索加载状态 */
  function setSearchLoading(loading: boolean) {
    searchLoading.value = loading
  }

  // ==================== 路线操作 ====================

  /** 设置路线结果 */
  function setRouteResult(result: RouteResult | null) {
    routeResult.value = result
  }

  /** 清空路线 */
  function clearRoute() {
    routeResult.value = null
  }

  /** 设置路线加载状态 */
  function setRouteLoading(loading: boolean) {
    routeLoading.value = loading
  }

  /** 切换路线显示 */
  function toggleRouteVisibility() {
    showRoute.value = !showRoute.value
  }

  /** 切换标记显示 */
  function toggleMarkersVisibility() {
    showMarkers.value = !showMarkers.value
  }

  // ==================== 编辑模式操作 ====================

  /** 进入编辑模式 */
  function enterEditMode(waypoints: Waypoint[] = []) {
    isEditing.value = true
    editingWaypoints.value = [...waypoints]
  }

  /** 退出编辑模式 */
  function exitEditMode() {
    isEditing.value = false
    editingWaypoints.value = []
  }

  /** 添加编辑中的途经点 */
  function addEditingWaypoint(waypoint: Waypoint) {
    editingWaypoints.value.push(waypoint)
  }

  /** 更新编辑中的途经点 */
  function updateEditingWaypoint(index: number, waypoint: Partial<Waypoint>) {
    if (index >= 0 && index < editingWaypoints.value.length) {
      editingWaypoints.value[index] = { ...editingWaypoints.value[index], ...waypoint }
    }
  }

  /** 删除编辑中的途经点 */
  function removeEditingWaypoint(index: number) {
    if (index >= 0 && index < editingWaypoints.value.length) {
      editingWaypoints.value.splice(index, 1)
    }
  }

  /** 重新排序编辑中的途经点 */
  function reorderEditingWaypoints(fromIndex: number, toIndex: number) {
    const waypoints = [...editingWaypoints.value]
    const [removed] = waypoints.splice(fromIndex, 1)
    waypoints.splice(toIndex, 0, removed)
    // 更新 sortOrder
    waypoints.forEach((wp, index) => {
      wp.sortOrder = index
    })
    editingWaypoints.value = waypoints
  }

  // ==================== 从途经点生成标记 ====================

  /** 从途经点列表生成地图标记 */
  function generateMarkersFromWaypoints(waypoints: Waypoint[]) {
    const newMarkers: MapMarker[] = waypoints.map((wp, index) => ({
      id: wp.id || `temp-${index}`,
      position: [wp.longitude, wp.latitude] as [number, number],
      title: wp.name,
      label: `${index + 1}`,
      data: wp
    }))
    setMarkers(newMarkers)
  }

  // ==================== 重置状态 ====================

  /** 重置所有状态 */
  function $reset() {
    mapInstance.value = null
    center.value = [116.397428, 39.90923]
    zoom.value = 12
    travelMode.value = 'driving' as TravelMode
    markers.value = []
    selectedMarker.value = null
    searchResults.value = []
    searchKeyword.value = ''
    searchLoading.value = false
    routeResult.value = null
    routeLoading.value = false
    showRoute.value = true
    showMarkers.value = true
    editingWaypoints.value = []
    isEditing.value = false
    mapLoaded.value = false
  }

  return {
    // 状态
    mapInstance,
    center,
    zoom,
    travelMode,
    markers,
    selectedMarker,
    searchResults,
    searchKeyword,
    searchLoading,
    routeResult,
    routeLoading,
    showRoute,
    showMarkers,
    editingWaypoints,
    isEditing,
    mapLoaded,
    // 计算属性
    markerCount,
    hasRoute,
    totalDistance,
    totalDuration,
    formattedDuration,
    // 地图操作
    setMapInstance,
    setCenter,
    setZoom,
    setTravelMode,
    panTo,
    fitBounds,
    // 标记操作
    addMarker,
    setMarkers,
    updateMarker,
    removeMarker,
    clearMarkers,
    selectMarker,
    // 搜索操作
    setSearchKeyword,
    setSearchResults,
    clearSearchResults,
    setSearchLoading,
    // 路线操作
    setRouteResult,
    clearRoute,
    setRouteLoading,
    toggleRouteVisibility,
    toggleMarkersVisibility,
    // 编辑模式操作
    enterEditMode,
    exitEditMode,
    addEditingWaypoint,
    updateEditingWaypoint,
    removeEditingWaypoint,
    reorderEditingWaypoints,
    // 工具方法
    generateMarkersFromWaypoints,
    // 重置
    $reset
  }
})
