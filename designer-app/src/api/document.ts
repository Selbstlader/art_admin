/**
 * 文档相关 API
 * Document related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, PageResponse, Document, DocumentSummary, PageParams } from '@/types/api'

/*** Document API - handles document operations ***/
export const documentApi = {
  /*** 获取文档列表 - Get document list ***/
  getList(projectId: number, params: PageParams): Promise<ApiResponse<PageResponse<Document>>> {
    return request({
      url: `/api/designer/projects/${projectId}/documents`,
      method: 'GET',
      data: params
    })
  },

  /*** 获取文档详情 - Get document detail ***/
  getDetail(id: number): Promise<ApiResponse<Document>> {
    return request({
      url: `/api/designer/documents/${id}`,
      method: 'GET'
    })
  },

  /*** 获取文档关键字 - Get document keywords ***/
  getKeywords(id: number): Promise<ApiResponse<string[]>> {
    return request({
      url: `/api/designer/documents/${id}/keywords`,
      method: 'GET'
    })
  },

  /*** 获取文档摘要 - Get document summary ***/
  getSummary(id: number): Promise<ApiResponse<DocumentSummary>> {
    return request({
      url: `/api/designer/documents/${id}/summary`,
      method: 'GET'
    })
  }
}
