/**
 * 材料相关 API
 * Material related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, Material, MaterialDetail, MaterialCategory, MaterialListParams } from '@/types/api'

/*** Material API - handles material operations ***/
export const materialApi = {
  /*** 获取材料分类 - Get material categories ***/
  getCategories(): Promise<ApiResponse<MaterialCategory[]>> {
    return request({
      url: '/api/v1/designer/materials/categories',
      method: 'GET'
    })
  },

  /*** 获取材料列表 - Get material list ***/
  getList(params: MaterialListParams): Promise<ApiResponse<PageResponse<Material>>> {
    return request({
      url: '/api/v1/designer/materials',
      method: 'GET',
      data: params
    })
  },

  /*** 获取材料详情 - Get material detail ***/
  getDetail(id: number): Promise<ApiResponse<MaterialDetail>> {
    return request({
      url: `/api/v1/designer/materials/${id}`,
      method: 'GET'
    })
  },

  /*** 搜索材料 - Search materials ***/
  search(keyword: string): Promise<ApiResponse<Material[]>> {
    return request({
      url: '/api/v1/designer/materials/search',
      method: 'GET',
      data: { keyword }
    })
  }
}
