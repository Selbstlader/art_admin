/**
 * 高德地图 Hook
 * 提供地图初始化、标记管理、路线绘制、POI搜索等功能
 */
import { ref, shallowRef, onUnmounted } from 'vue'
import type { MapMarker, POISearchResult, Waypoint } from '../types'

// 高德地图API Key 和安全密钥
const AMAP_KEY = import.meta.env.VITE_AMAP_KEY || '5e55166796e21f88969cf7b11febba65'
const AMAP_SECURITY_CODE =
  import.meta.env.VITE_AMAP_SECURITY_CODE || 'fd63eded3b6ec5e541c171f106fe6416'
const AMAP_VERSION = '2.0'

// 全局AMap对象
let AMapGlobal: any = null
let loadPromise: Promise<any> | null = null
let securityConfigured = false

// 类型定义
export interface AMapInstance {
  setCenter(center: [number, number]): void
  getCenter(): { getLng(): number; getLat(): number }
  setZoom(zoom: number): void
  getZoom(): number
  panTo(position: [number, number]): void
  setFitView(overlays?: any[]): void
  add(overlay: any): void
  remove(overlay: any): void
  clearMap(): void
  on(event: string, handler: Function): void
  off(event: string, handler: Function): void
  destroy(): void
  getContainer(): HTMLElement
  plugin(plugins: string | string[], callback: Function): void
}

export interface MarkerInstance {
  setPosition(position: [number, number]): void
  getPosition(): { getLng(): number; getLat(): number }
  setLabel(label: { content: string; offset?: [number, number] }): void
  setIcon(icon: any): void
  setExtData(data: any): void
  getExtData(): any
  on(event: string, handler: Function): void
  off(event: string, handler: Function): void
}

export interface PolylineInstance {
  setPath(path: [number, number][]): void
  getPath(): any[]
  setOptions(options: any): void
  show(): void
  hide(): void
}

/**
 * 加载高德地图脚本
 */
export function loadAMapScript(): Promise<any> {
  if (AMapGlobal) {
    return Promise.resolve(AMapGlobal)
  }

  if (loadPromise) {
    return loadPromise
  }

  loadPromise = new Promise((resolve, reject) => {
    // 配置安全密钥（必须在加载地图脚本之前）
    if (!securityConfigured) {
      ;(window as any)._AMapSecurityConfig = {
        securityJsCode: AMAP_SECURITY_CODE
      }
      securityConfigured = true
    }

    // 检查是否已加载
    if ((window as any).AMap) {
      AMapGlobal = (window as any).AMap
      resolve(AMapGlobal)
      return
    }

    // 创建script标签
    const script = document.createElement('script')
    script.type = 'text/javascript'
    script.src = `https://webapi.amap.com/maps?v=${AMAP_VERSION}&key=${AMAP_KEY}&plugin=AMap.Scale,AMap.ToolBar,AMap.Geolocation,AMap.PlaceSearch,AMap.AutoComplete,AMap.Geocoder,AMap.Driving,AMap.Walking,AMap.Riding,AMap.Transfer`
    script.async = true

    script.onload = () => {
      AMapGlobal = (window as any).AMap
      if (AMapGlobal) {
        resolve(AMapGlobal)
      } else {
        reject(new Error('高德地图加载失败'))
      }
    }

    script.onerror = () => {
      loadPromise = null
      reject(new Error('高德地图脚本加载失败'))
    }

    document.head.appendChild(script)
  })

  return loadPromise
}

/**
 * 获取AMap对象
 */
export function getAMap(): any {
  return AMapGlobal
}

/**
 * useAMap Hook
 * 提供地图操作的组合式函数
 */
export function useAMap() {
  const map = shallowRef<AMapInstance | null>(null)
  const markers = shallowRef<Map<string | number, MarkerInstance>>(new Map())
  const polylines = shallowRef<PolylineInstance[]>([])
  const infoWindow = shallowRef<any>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // POI搜索相关
  const placeSearch = shallowRef<any>(null)
  const autoComplete = shallowRef<any>(null)
  const searchResults = ref<POISearchResult[]>([])
  const searchLoading = ref(false)

  /**
   * 初始化地图
   */
  async function initMap(
    container: HTMLElement | string,
    options: {
      center?: [number, number]
      zoom?: number
      mapStyle?: string
      viewMode?: '2D' | '3D'
      pitch?: number
    } = {}
  ): Promise<AMapInstance> {
    loading.value = true
    error.value = null

    try {
      const AMap = await loadAMapScript()

      const mapInstance = new AMap.Map(container, {
        center: options.center || [116.397428, 39.90923],
        zoom: options.zoom || 12,
        mapStyle: options.mapStyle || 'amap://styles/normal',
        viewMode: options.viewMode || '2D',
        pitch: options.pitch || 0,
        resizeEnable: true
      })

      map.value = mapInstance
      loading.value = false

      // 初始化信息窗体
      infoWindow.value = new AMap.InfoWindow({
        isCustom: true,
        autoMove: true,
        offset: new AMap.Pixel(0, -30)
      })

      return mapInstance
    } catch (err: any) {
      loading.value = false
      error.value = err.message || '地图初始化失败'
      throw err
    }
  }

  /**
   * 添加标记
   */
  function addMarker(markerData: MapMarker): MarkerInstance | null {
    if (!map.value) return null

    const AMap = getAMap()
    if (!AMap) return null

    const marker = new AMap.Marker({
      position: markerData.position,
      title: markerData.title,
      extData: markerData.data
    })

    // 设置标签
    if (markerData.label) {
      marker.setLabel({
        content: `<div class="amap-marker-label">${markerData.label}</div>`,
        offset: new AMap.Pixel(-10, -10)
      })
    }

    // 设置自定义图标
    if (markerData.icon) {
      marker.setIcon(new AMap.Icon({
        image: markerData.icon,
        size: new AMap.Size(32, 32),
        imageSize: new AMap.Size(32, 32)
      }))
    }

    map.value.add(marker)
    markers.value.set(markerData.id, marker)

    return marker
  }

  /**
   * 批量添加标记
   */
  function addMarkers(markerDataList: MapMarker[]): void {
    markerDataList.forEach(data => addMarker(data))
  }

  /**
   * 更新标记位置
   */
  function updateMarkerPosition(id: string | number, position: [number, number]): void {
    const marker = markers.value.get(id)
    if (marker) {
      marker.setPosition(position)
    }
  }

  /**
   * 删除标记
   */
  function removeMarker(id: string | number): void {
    const marker = markers.value.get(id)
    if (marker && map.value) {
      map.value.remove(marker)
      markers.value.delete(id)
    }
  }

  /**
   * 清空所有标记
   */
  function clearMarkers(): void {
    if (!map.value) return

    markers.value.forEach(marker => {
      map.value!.remove(marker)
    })
    markers.value.clear()
  }

  /**
   * 从途经点生成标记
   */
  function setMarkersFromWaypoints(waypoints: Waypoint[]): void {
    clearMarkers()

    waypoints.forEach((wp, index) => {
      addMarker({
        id: wp.id || `waypoint-${index}`,
        position: [wp.longitude, wp.latitude],
        title: wp.name,
        label: `${index + 1}`,
        data: wp
      })
    })
  }

  /**
   * 绘制路线
   */
  function drawPolyline(
    path: [number, number][],
    options: {
      strokeColor?: string
      strokeWeight?: number
      strokeOpacity?: number
      strokeStyle?: 'solid' | 'dashed'
    } = {}
  ): PolylineInstance | null {
    if (!map.value) return null

    const AMap = getAMap()
    if (!AMap) return null

    const polyline = new AMap.Polyline({
      path,
      strokeColor: options.strokeColor || '#3366FF',
      strokeWeight: options.strokeWeight || 6,
      strokeOpacity: options.strokeOpacity || 0.8,
      strokeStyle: options.strokeStyle || 'solid',
      lineJoin: 'round',
      lineCap: 'round'
    })

    map.value.add(polyline)
    polylines.value.push(polyline)

    return polyline
  }

  /**
   * 清空所有路线
   */
  function clearPolylines(): void {
    if (!map.value) return

    polylines.value.forEach(polyline => {
      map.value!.remove(polyline)
    })
    polylines.value = []
  }

  /**
   * 显示信息窗体
   */
  function showInfoWindow(
    position: [number, number],
    content: string | HTMLElement
  ): void {
    if (!map.value || !infoWindow.value) return

    infoWindow.value.setContent(content)
    infoWindow.value.open(map.value, position)
  }

  /**
   * 关闭信息窗体
   */
  function closeInfoWindow(): void {
    if (infoWindow.value) {
      infoWindow.value.close()
    }
  }

  /**
   * POI搜索
   */
  async function searchPOI(
    keyword: string,
    options: {
      city?: string
      pageSize?: number
      pageIndex?: number
    } = {}
  ): Promise<POISearchResult[]> {
    searchLoading.value = true
    searchResults.value = []

    try {
      const AMap = await loadAMapScript()

      return new Promise((resolve, reject) => {
        if (!placeSearch.value) {
          placeSearch.value = new AMap.PlaceSearch({
            city: options.city || '全国',
            pageSize: options.pageSize || 10,
            pageIndex: options.pageIndex || 1
          })
        }

        placeSearch.value.search(keyword, (status: string, result: any) => {
          searchLoading.value = false

          if (status === 'complete' && result.poiList) {
            const pois: POISearchResult[] = result.poiList.pois.map((poi: any) => ({
              id: poi.id,
              name: poi.name,
              address: poi.address || '',
              location: {
                lng: poi.location.getLng(),
                lat: poi.location.getLat()
              },
              type: poi.type || '',
              distance: poi.distance
            }))
            searchResults.value = pois
            resolve(pois)
          } else {
            searchResults.value = []
            resolve([])
          }
        })
      })
    } catch (err: any) {
      searchLoading.value = false
      throw err
    }
  }

  /**
   * 输入提示搜索
   */
  async function autoCompleteSearch(
    keyword: string,
    city?: string
  ): Promise<POISearchResult[]> {
    try {
      const AMap = await loadAMapScript()

      return new Promise((resolve) => {
        if (!autoComplete.value) {
          autoComplete.value = new AMap.AutoComplete({
            city: city || '全国'
          })
        }

        autoComplete.value.search(keyword, (status: string, result: any) => {
          if (status === 'complete' && result.tips) {
            const tips: POISearchResult[] = result.tips
              .filter((tip: any) => tip.location)
              .map((tip: any) => ({
                id: tip.id,
                name: tip.name,
                address: tip.address || tip.district || '',
                location: {
                  lng: tip.location.getLng(),
                  lat: tip.location.getLat()
                },
                type: tip.typecode || ''
              }))
            resolve(tips)
          } else {
            resolve([])
          }
        })
      })
    } catch {
      return []
    }
  }

  /**
   * 设置地图中心点
   */
  function setCenter(center: [number, number]): void {
    map.value?.setCenter(center)
  }

  /**
   * 设置缩放级别
   */
  function setZoom(zoom: number): void {
    map.value?.setZoom(zoom)
  }

  /**
   * 平移到指定位置
   */
  function panTo(position: [number, number]): void {
    map.value?.panTo(position)
  }

  /**
   * 自适应显示所有覆盖物
   */
  function fitView(): void {
    if (!map.value) return

    const overlays = [...markers.value.values(), ...polylines.value]
    if (overlays.length > 0) {
      map.value.setFitView(overlays)
    }
  }

  /**
   * 清空地图
   */
  function clearMap(): void {
    clearMarkers()
    clearPolylines()
    closeInfoWindow()
  }

  /**
   * 销毁地图
   */
  function destroyMap(): void {
    clearMap()
    if (map.value) {
      map.value.destroy()
      map.value = null
    }
  }

  // 组件卸载时清理
  onUnmounted(() => {
    destroyMap()
  })

  return {
    // 状态
    map,
    markers,
    polylines,
    loading,
    error,
    searchResults,
    searchLoading,
    // 地图操作
    initMap,
    setCenter,
    setZoom,
    panTo,
    fitView,
    clearMap,
    destroyMap,
    // 标记操作
    addMarker,
    addMarkers,
    updateMarkerPosition,
    removeMarker,
    clearMarkers,
    setMarkersFromWaypoints,
    // 路线操作
    drawPolyline,
    clearPolylines,
    // 信息窗体
    showInfoWindow,
    closeInfoWindow,
    // POI搜索
    searchPOI,
    autoCompleteSearch
  }
}

export default useAMap
