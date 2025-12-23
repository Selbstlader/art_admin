/**
 * 设计图相关 API
 * Design image related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, DesignImage, DesignImageListParams } from '@/types/api'

/*** Design API - handles design image operations ***/
export const designApi = {
  /*** 获取设计图列表 - Get design image list ***/
  getList(params: DesignImageListParams): Promise<ApiResponse<PageResponse<DesignImage>>> {
    return request({
      url: `/api/designer/projects/${params.projectId}/designs`,
      method: 'GET',
      data: {
        current: params.current,
        size: params.size
      }
    })
  },

  /*** 获取设计图详情 - Get design image detail ***/
  getDetail(id: number): Promise<ApiResponse<DesignImage>> {
    return request({
      url: `/api/designer/designs/${id}`,
      method: 'GET'
    })
  },

  /*** 获取项目所有设计图 - Get all design images for a project ***/
  getAllByProject(projectId: number): Promise<ApiResponse<DesignImage[]>> {
    return request({
      url: `/api/designer/projects/${projectId}/designs/all`,
      method: 'GET'
    })
  }
}
