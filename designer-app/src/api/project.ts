/**
 * 项目相关 API
 * Project related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, Project, ProjectDetail, ProjectListParams } from '@/types/api'

/*** Project API - handles project operations ***/
export const projectApi = {
  /*** 获取项目列表 - Get project list ***/
  getList(params: ProjectListParams): Promise<ApiResponse<PageResponse<Project>>> {
    return request({
      url: '/api/designer/projects',
      method: 'GET',
      data: params
    })
  },

  /*** 获取项目详情 - Get project detail ***/
  getDetail(id: number): Promise<ApiResponse<ProjectDetail>> {
    return request({
      url: `/api/designer/projects/${id}`,
      method: 'GET'
    })
  },

  /*** 搜索项目 - Search projects ***/
  search(keyword: string, page: number, size: number): Promise<ApiResponse<PageResponse<Project>>> {
    return request({
      url: '/api/designer/projects',
      method: 'GET',
      data: { name: keyword, current: page, size }
    })
  }
}
