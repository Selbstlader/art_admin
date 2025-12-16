import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { templateRoutes } from './template'
import { widgetsRoutes } from './widgets'
import { examplesRoutes } from './examples'
import { systemRoutes } from './system'
import { resultRoutes } from './result'
import { exceptionRoutes } from './exception'
import { safeguardRoutes } from './safeguard'
import { helpRoutes } from './help'
import chatRoutes from './chat'
import { difyRoutes } from './dify'

/**
 * 导出所有模块化路由
 * 基础功能 + 系统功能 + Dify AI + 聊天室
 */
export const routeModules: AppRouteRecord[] = [
  dashboardRoutes,
  chatRoutes,
  difyRoutes,
  templateRoutes,
  widgetsRoutes,
  examplesRoutes,
  systemRoutes,
  resultRoutes,
  exceptionRoutes,
  safeguardRoutes,
  ...helpRoutes
]
