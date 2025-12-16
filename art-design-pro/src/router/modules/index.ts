import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { systemRoutes } from './system'
import { exceptionRoutes } from './exception'
import chatRoutes from './chat'
import { difyRoutes } from './dify'

/**
 * 导出所有模块化路由
 * 核心功能: 仪表盘 + 系统管理 + 聊天室 + Dify AI
 */
export const routeModules: AppRouteRecord[] = [
  dashboardRoutes,
  chatRoutes,
  difyRoutes,
  systemRoutes,
  exceptionRoutes
]
