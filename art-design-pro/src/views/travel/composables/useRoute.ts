/**
 * 路线规划 Hook
 * 提供路线规划API调用、多种出行方式切换功能
 */
import { ref, computed, shallowRef } from 'vue'
import { loadAMapScript, getAMap, type AMapInstance } from './useAMap'
import type { TravelMode, RouteResult, RouteStep, Waypoint } from '../types'

// 路线规划服务实例缓存
interface RouteServices {
  driving: any
  walking: any
  cycling: any
  transit: any
  [key: string]: any
}

/**
 * useRoute Hook
 * 提供路线规划的组合式函数
 */
export function useRoute(mapInstance?: AMapInstance | null) {
  // 状态
  const loading = ref(false)
  const error = ref<string | null>(null)
  const currentMode = ref<TravelMode>('driving' as TravelMode)
  const routeResult = ref<RouteResult | null>(null)
  const routePolyline = shallowRef<any>(null)

  // 路线规划服务实例
  const services = shallowRef<Partial<RouteServices>>({})

  // 计算属性
  const hasRoute = computed(() => routeResult.value !== null)

  const totalDistance = computed(() => {
    if (!routeResult.value) return 0
    return routeResult.value.distance
  })

  const totalDuration = computed(() => {
    if (!routeResult.value) return 0
    return routeResult.value.duration
  })

  const formattedDistance = computed(() => {
    const meters = totalDistance.value
    if (meters < 1000) {
      return `${meters}米`
    }
    return `${(meters / 1000).toFixed(1)}公里`
  })

  const formattedDuration = computed(() => {
    const seconds = totalDuration.value
    const minutes = Math.round(seconds / 60)
    if (minutes < 60) {
      return `${minutes}分钟`
    }
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
  })

  /**
   * 获取或创建路线规划服务
   */
  async function getRouteService(mode: TravelMode): Promise<any> {
    const AMap = await loadAMapScript()

    if (services.value[mode]) {
      return services.value[mode]
    }

    let service: any

    switch (mode) {
      case 'driving':
        service = new AMap.Driving({
          policy: AMap.DrivingPolicy.LEAST_TIME,
          extensions: 'all'
        })
        break
      case 'walking':
        service = new AMap.Walking({
          extensions: 'all'
        })
        break
      case 'cycling':
        service = new AMap.Riding({
          extensions: 'all'
        })
        break
      case 'transit':
        service = new AMap.Transfer({
          city: '北京',
          policy: AMap.TransferPolicy.LEAST_TIME,
          extensions: 'all'
        })
        break
      default:
        service = new AMap.Driving({
          policy: AMap.DrivingPolicy.LEAST_TIME,
          extensions: 'all'
        })
    }

    services.value[mode] = service
    return service
  }

  /**
   * 规划路线
   */
  async function planRoute(
    origin: [number, number],
    destination: [number, number],
    waypoints?: [number, number][],
    mode?: TravelMode
  ): Promise<RouteResult | null> {
    loading.value = true
    error.value = null

    const travelMode = mode || currentMode.value

    try {
      const AMap = await loadAMapScript()
      const service = await getRouteService(travelMode)

      return new Promise((resolve, reject) => {
        const callback = (status: string, result: any) => {
          loading.value = false

          if (status === 'complete') {
            const routeData = parseRouteResult(result, travelMode, AMap)
            routeResult.value = routeData
            resolve(routeData)
          } else {
            const errorMsg = getRouteErrorMessage(result)
            error.value = errorMsg
            routeResult.value = null
            reject(new Error(errorMsg))
          }
        }

        // 根据不同出行方式调用不同的API
        if (travelMode === 'transit') {
          service.search(
            new AMap.LngLat(origin[0], origin[1]),
            new AMap.LngLat(destination[0], destination[1]),
            callback
          )
        } else if (waypoints && waypoints.length > 0) {
          // 有途经点的情况
          const waypointLngLats = waypoints.map(
            (wp) => new AMap.LngLat(wp[0], wp[1])
          )
          service.search(
            new AMap.LngLat(origin[0], origin[1]),
            new AMap.LngLat(destination[0], destination[1]),
            { waypoints: waypointLngLats },
            callback
          )
        } else {
          service.search(
            new AMap.LngLat(origin[0], origin[1]),
            new AMap.LngLat(destination[0], destination[1]),
            callback
          )
        }
      })
    } catch (err: any) {
      loading.value = false
      error.value = err.message || '路线规划失败'
      throw err
    }
  }

  /**
   * 从途经点列表规划路线
   */
  async function planRouteFromWaypoints(
    waypoints: Waypoint[],
    mode?: TravelMode
  ): Promise<RouteResult | null> {
    if (waypoints.length < 2) {
      error.value = '至少需要2个途经点'
      return null
    }

    // 按sortOrder排序
    const sortedWaypoints = [...waypoints].sort((a, b) => a.sortOrder - b.sortOrder)

    const origin: [number, number] = [
      sortedWaypoints[0].longitude,
      sortedWaypoints[0].latitude
    ]
    const destination: [number, number] = [
      sortedWaypoints[sortedWaypoints.length - 1].longitude,
      sortedWaypoints[sortedWaypoints.length - 1].latitude
    ]

    // 中间途经点
    const middleWaypoints: [number, number][] = sortedWaypoints
      .slice(1, -1)
      .map((wp) => [wp.longitude, wp.latitude])

    return planRoute(origin, destination, middleWaypoints, mode)
  }

  /**
   * 解析路线规划结果
   */
  function parseRouteResult(result: any, mode: TravelMode, AMap: any): RouteResult {
    let distance = 0
    let duration = 0
    let polyline: [number, number][] = []
    let steps: RouteStep[] = []

    if (mode === 'transit') {
      // 公交路线解析
      const plan = result.plans?.[0]
      if (plan) {
        distance = plan.distance || 0
        duration = plan.time || 0

        plan.segments?.forEach((segment: any) => {
          if (segment.transit) {
            const path = segment.transit.path || []
            polyline.push(...path.map((p: any) => [p.getLng(), p.getLat()]))
          }
          if (segment.walking) {
            const path = segment.walking.path || []
            polyline.push(...path.map((p: any) => [p.getLng(), p.getLat()]))
          }
        })
      }
    } else {
      // 驾车/步行/骑行路线解析
      const route = result.routes?.[0]
      if (route) {
        distance = route.distance || 0
        duration = route.time || 0

        route.steps?.forEach((step: any) => {
          const stepPath = step.path || []
          const stepPolyline: [number, number][] = stepPath.map((p: any) => [
            p.getLng(),
            p.getLat()
          ])

          polyline.push(...stepPolyline)

          steps.push({
            instruction: step.instruction || '',
            distance: step.distance || 0,
            duration: step.time || 0,
            polyline: stepPolyline
          })
        })
      }
    }

    return {
      distance,
      duration,
      polyline,
      steps
    }
  }

  /**
   * 获取路线规划错误信息
   */
  function getRouteErrorMessage(result: any): string {
    const info = result?.info || ''
    if (info.includes('OVER_DIRECTION_RANGE')) {
      return '起终点距离过长，无法规划路线'
    }
    if (info.includes('NO_ROADS_NEARBY')) {
      return '起点或终点附近没有可用道路'
    }
    if (info.includes('ROUTE_FAIL')) {
      return '路线规划失败，请检查起终点'
    }
    return '路线规划失败'
  }

  /**
   * 在地图上绘制路线
   */
  async function drawRouteOnMap(map: AMapInstance, result?: RouteResult): Promise<void> {
    const AMap = getAMap()
    if (!AMap || !map) return

    // 清除之前的路线
    clearRouteFromMap(map)

    const route = result || routeResult.value
    if (!route || route.polyline.length === 0) return

    // 创建路线折线
    const polyline = new AMap.Polyline({
      path: route.polyline,
      strokeColor: '#3366FF',
      strokeWeight: 6,
      strokeOpacity: 0.8,
      lineJoin: 'round',
      lineCap: 'round',
      showDir: true
    })

    map.add(polyline)
    routePolyline.value = polyline

    // 自适应显示
    map.setFitView([polyline])
  }

  /**
   * 从地图上清除路线
   */
  function clearRouteFromMap(map: AMapInstance): void {
    if (routePolyline.value && map) {
      map.remove(routePolyline.value)
      routePolyline.value = null
    }
  }

  /**
   * 切换出行方式
   */
  function setTravelMode(mode: TravelMode): void {
    currentMode.value = mode
  }

  /**
   * 设置公交城市
   */
  async function setTransitCity(city: string): Promise<void> {
    const AMap = await loadAMapScript()
    if (services.value.transit) {
      services.value.transit.setCity(city)
    } else {
      services.value.transit = new AMap.Transfer({
        city,
        policy: AMap.TransferPolicy.LEAST_TIME,
        extensions: 'all'
      })
    }
  }

  /**
   * 清除路线结果
   */
  function clearRoute(): void {
    routeResult.value = null
    error.value = null
  }

  /**
   * 重置状态
   */
  function reset(): void {
    loading.value = false
    error.value = null
    currentMode.value = 'driving' as TravelMode
    routeResult.value = null
    routePolyline.value = null
  }

  return {
    // 状态
    loading,
    error,
    currentMode,
    routeResult,
    // 计算属性
    hasRoute,
    totalDistance,
    totalDuration,
    formattedDistance,
    formattedDuration,
    // 方法
    planRoute,
    planRouteFromWaypoints,
    drawRouteOnMap,
    clearRouteFromMap,
    setTravelMode,
    setTransitCity,
    clearRoute,
    reset
  }
}

export default useRoute
