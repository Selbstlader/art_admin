/**
 * 动态路由处理
 * 根据接口返回的菜单列表注册动态路由
 */
import type { Router, RouteRecordRaw } from 'vue-router'
import type { AppRouteRecord } from '@/types/router'
import { saveIframeRoutes } from './menuToRouter'
import { h } from 'vue'
import { useMenuStore } from '@/store/modules/menu'
import { RoutesAlias } from '../routesAlias'

/**
 * 动态导入 views 目录下所有 .vue 组件
 */
const modules: Record<string, () => Promise<any>> = import.meta.glob('../../views/**/*.vue')

/**
 * 注册异步路由
 * 将接口返回的菜单列表转换为 Vue Router 路由配置，并添加到传入的 router 实例中
 * @param router Vue Router 实例
 * @param menuList 接口返回的菜单列表
 */
export function registerDynamicRoutes(router: Router, menuList: AppRouteRecord[]): void {
  // 用于局部收集 iframe 类型路由
  const iframeRoutes: AppRouteRecord[] = []
  // 收集路由移除函数
  const removeRouteFns: (() => void)[] = []

  // 检测菜单列表中是否有重复路由
  checkDuplicateRoutes(menuList)

  // 遍历菜单列表，注册路由
  menuList.forEach((route) => {
    // 只有还没注册过的路由才进行注册
    if (route.name && !router.hasRoute(route.name)) {
      const routeConfig = convertRouteComponent(route, iframeRoutes)
      // addRoute 返回移除函数，收集起来
      const removeRouteFn = router.addRoute(routeConfig as RouteRecordRaw)
      removeRouteFns.push(removeRouteFn)
    }
  })

  // 将移除函数存储到 store 中
  const menuStore = useMenuStore()
  menuStore.addRemoveRouteFns(removeRouteFns)

  // 保存 iframe 路由
  saveIframeRoutes(iframeRoutes)
}

/**
 * 路径解析函数：处理父路径和子路径的拼接
 */
function resolvePath(parent: string, child: string): string {
  return [parent.replace(/\/$/, ''), child.replace(/^\//, '')].filter(Boolean).join('/')
}

/**
 * 检测菜单中的重复路由（包括子路由）
 */
function checkDuplicateRoutes(routes: AppRouteRecord[], parentPath = ''): void {
  // 用于检测动态路由中的重复项
  const routeNameMap = new Map<string, string>() // 路由名称 -> 路径
  const componentPathMap = new Map<string, string>() // 组件路径 -> 路由信息

  const checkRoutes = (routes: AppRouteRecord[], parentPath = '') => {
    routes.forEach((route) => {
      // 处理路径拼接
      const currentPath = route.path || ''
      const fullPath = resolvePath(parentPath, currentPath)

      // 名称重复检测
      if (route.name) {
        if (routeNameMap.has(String(route.name))) {
          console.warn(`[路由警告] 名称重复: "${String(route.name)}"`)
        } else {
          routeNameMap.set(String(route.name), fullPath)
        }
      }

      // 组件路径重复检测
      if (route.component) {
        const componentPath = getComponentPathString(route.component)

        if (componentPath && componentPath !== RoutesAlias.Layout) {
          const componentKey = `${parentPath}:${componentPath}`

          if (componentPathMap.has(componentKey)) {
            console.warn(`[路由警告] 路径重复: "${componentPath}"`)
          } else {
            componentPathMap.set(componentKey, fullPath)
          }
        }
      }

      // 递归处理子路由
      if (route.children?.length) {
        checkRoutes(route.children, fullPath)
      }
    })
  }

  checkRoutes(routes, parentPath)
}

/**
 * 获取组件路径的字符串表示
 */
function getComponentPathString(component: any): string {
  if (typeof component === 'string') {
    return component
  }

  return ''
}

/**
 * 根据组件路径动态加载组件
 * @param componentPath 组件路径（不包含 ../../views 前缀和 .vue 后缀）
 * @param routeName 当前路由名称（用于错误提示）
 * @returns 组件加载函数
 */
function loadComponent(componentPath: string, routeName: string): () => Promise<any> {
  // 如果路径为空，直接返回一个空的组件
  if (componentPath === '') {
    return () =>
      Promise.resolve({
        render() {
          return h('div', {})
        }
      })
  }

  // 构建可能的路径
  const fullPath = `../../views${componentPath}.vue`
  const fullPathWithIndex = `../../views${componentPath}/index.vue`

  // 先尝试直接路径，再尝试添加/index的路径
  const module = modules[fullPath] || modules[fullPathWithIndex]

  if (!module) {
    console.error(
      `[路由错误] 未找到组件：${routeName}，尝试过的路径: ${fullPath} 和 ${fullPathWithIndex}`
    )
    return () =>
      Promise.resolve({
        render() {
          return h('div', `组件未找到: ${routeName}`)
        }
      })
  }

  return module
}

/**
 * 转换后的路由配置类型
 */
interface ConvertedRoute extends Omit<RouteRecordRaw, 'children'> {
  id?: number
  children?: ConvertedRoute[]
  component?: RouteRecordRaw['component'] | (() => Promise<any>)
}

/**
 * 转换路由组件配置
 */
function convertRouteComponent(
  route: AppRouteRecord,
  iframeRoutes: AppRouteRecord[],
  depth = 0,
  parentPath = ''
): ConvertedRoute {
  const { component, children, ...routeConfig } = route

  // 基础路由配置
  const converted: ConvertedRoute = {
    ...routeConfig,
    component: undefined
  }

  // 子路由需要将绝对路径转换为相对路径
  // Vue Router 嵌套路由会自动拼接父路由 path，所以子路由应使用相对路径
  if (depth > 0 && route.path && route.path.startsWith('/')) {
    // 如果子路由 path 以父路由 path 开头，则提取相对路径部分
    if (parentPath && route.path.startsWith(parentPath + '/')) {
      converted.path = route.path.slice(parentPath.length + 1)
    } else if (parentPath && route.path.startsWith(parentPath)) {
      converted.path = route.path.slice(parentPath.length) || ''
    } else {
      // 如果不是以父路径开头，取最后一段作为相对路径
      const segments = route.path.split('/').filter(Boolean)
      converted.path = segments[segments.length - 1] || ''
    }
  }

  // 处理动态路由参数
  // 优先使用 meta.dynamicPath（后端配置的完整动态路径）
  // 其次使用 meta.dynamicParams（后端配置的参数名，如 'id' 或 ['id', 'type']）
  // 最后兼容旧的硬编码方式
  if (route.meta?.dynamicPath) {
    // 方式1: 后端直接配置完整的动态路径
    converted.path = route.meta.dynamicPath as string
  } else if (route.meta?.dynamicParams) {
    // 方式2: 后端配置参数名，前端自动拼接
    const params = route.meta.dynamicParams
    const paramStr = Array.isArray(params)
      ? params.map((p: string) => `:${p}`).join('/')
      : `:${params}`
    converted.path = `${converted.path || route.path}/${paramStr}`
  } else {
    // 方式3: 兼容旧的硬编码配置（逐步迁移后可删除）
    // 注意：key 为后端返回的 path，value 为带动态参数的完整路径
    const dynamicRouteMap: Record<string, string> = {
      '/project/gantt': '/project/gantt/:id',
      '/learning/learning/detail': '/learning/learning/detail/:id',
      '/chat/chat/room-detail': '/chat/chat/room-detail/:id',
      '/mood/mood/meditation/player': '/mood/mood/meditation/player/:id',
      '/travel/travel/roadbook/editor': '/travel/travel/roadbook/editor/:id',
      '/travel/travel/roadbook/detail': '/travel/travel/roadbook/detail/:id',
      // 设计师助手模块 - 动态路由配置
      '/designer/designer-assistant/project/ProjectDetail':
        '/designer/designer-assistant/project/ProjectDetail/:id',
      '/designer/designer-assistant/project/ProjectCreate':
        '/designer/designer-assistant/project/ProjectCreate/:id'
      // '/designer/designer-assistant/version-compare/VersionDiff':
      //   '/designer/designer-assistant/version-compare/VersionDiff/:id'
    }
    if (route.path && dynamicRouteMap[route.path]) {
      converted.path = dynamicRouteMap[route.path]
    }
  }

  // 判断是否为一级路由
  const isTopLevel = depth === 0
  // 判断是否有子路由
  const hasChildren = route.children && route.children.length > 0
  // 判断 component 是否为 Layout
  const isLayoutComponent = component === RoutesAlias.Layout
  // 是否为没有子路由的一级菜单（需要 Layout 包裹单个页面）
  const isFirstLevelSinglePage = isTopLevel && !hasChildren && !isLayoutComponent

  if (route.meta?.isIframe) {
    handleIframeRoute(converted, route, iframeRoutes, depth)
  } else if (isTopLevel && (hasChildren || isLayoutComponent)) {
    // 一级路由：有子路由 或 component 是 Layout，使用 Layout 组件
    handleTopLevelWithLayout(converted, route)
  } else if (isFirstLevelSinglePage) {
    // 一级路由：没有子路由且 component 不是 Layout，使用 Layout 包裹单个页面
    handleLayoutRoute(converted, route, component as string)
  } else {
    // 其他情况：普通路由处理
    handleNormalRoute(converted, component as string, String(route.name))
  }

  // 递归处理子路由
  if (children?.length) {
    converted.children = children.map((child) =>
      convertRouteComponent(child, iframeRoutes, depth + 1, route.path || '')
    )
  }

  return converted
}
/**
 * 处理 iframe 类型路由
 */
function handleIframeRoute(
  targetRoute: ConvertedRoute,
  sourceRoute: AppRouteRecord,
  iframeRoutes: AppRouteRecord[],
  depth: number
): void {
  const LAYOUT_VIEW = () => import('@/views/index/index.vue')
  const IFRAME_VIEW = () => import('@/views/outside/Iframe.vue')

  if (depth === 0) {
    // 顶级 iframe：用 Layout 包裹
    targetRoute.component = LAYOUT_VIEW
    targetRoute.path = `/${(sourceRoute.path?.split('/')[1] || '').trim()}`
    targetRoute.name = ''

    targetRoute.children = [
      {
        ...sourceRoute,
        component: IFRAME_VIEW
      } as ConvertedRoute
    ]
  } else {
    // 非顶级（嵌套）iframe：直接使用 Iframe.vue
    targetRoute.component = IFRAME_VIEW
  }

  // 记录 iframe 路由，供 Iframe.vue 查找对应的外链
  iframeRoutes.push(sourceRoute)
}

/**
 * 处理有子路由的一级路由（使用 Layout）
 */
function handleTopLevelWithLayout(converted: ConvertedRoute, route: AppRouteRecord): void {
  converted.component = () => import('@/views/index/index.vue')
  converted.path = `/${(route.path?.split('/')[1] || '').trim()}`
  converted.name = ''
  // 子路由会在递归处理时添加到 converted.children 中
}

/**
 * 处理一级菜单路由（没有子路由，使用 Layout 包裹单个页面）
 */
function handleLayoutRoute(
  converted: ConvertedRoute,
  route: AppRouteRecord,
  component: string | undefined
): void {
  converted.component = () => import('@/views/index/index.vue')
  converted.path = `/${(route.path?.split('/')[1] || '').trim()}`
  converted.name = ''
  route.meta.isFirstLevel = true

  converted.children = [
    {
      ...route,
      component: loadComponent(component as string, String(route.name))
    } as ConvertedRoute
  ]
}

/**
 * 处理普通路由
 */
function handleNormalRoute(
  converted: ConvertedRoute,
  component: string | undefined,
  routeName: string
): void {
  if (component) {
    // 直接使用组件路径加载组件
    converted.component = loadComponent(component as string, routeName)
  }
}
