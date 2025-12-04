/**
 * 旅游规划模块组合式函数
 */

// 高德地图操作
export { useAMap, loadAMapScript, getAMap } from './useAMap'
export type { AMapInstance, MarkerInstance, PolylineInstance } from './useAMap'

// 路线规划
export { useRoute } from './useRoute'

// 后续将添加以下 hooks:
// - useNavigation: 导航导出
