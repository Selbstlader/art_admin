/**
 * 成本相关 API
 * Cost related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, CostEstimate } from '@/types/api'

/*** Cost API - handles cost estimation operations ***/
export const costApi = {
  /*** 获取项目成本估算 - Get project cost estimate ***/
  getByProjectId(projectId: number): Promise<ApiResponse<CostEstimate>> {
    return request({
      url: `/api/v1/designer/projects/${projectId}/cost`,
      method: 'GET'
    })
  }
}
