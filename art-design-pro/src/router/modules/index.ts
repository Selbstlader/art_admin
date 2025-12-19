import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { systemRoutes } from './system'
import { exceptionRoutes } from './exception'
import chatRoutes from './chat'
import { difyRoutes } from './dify'
import { designerRoutes } from './designer'

/**
 * 导出所有模块化路由
 * 核心功能: 仪表盘 + 系统管理 + 聊天室 + Dify AI + 设计师助手
 */
export const routeModules: AppRouteRecord[] = [
  dashboardRoutes,
  chatRoutes,
  difyRoutes,
  designerRoutes,
  systemRoutes,
  exceptionRoutes
]
