import request from '@/utils/http'

/** 文件信息 */
export interface FileInfo {
  id: number
  fileName: string
  storageName: string
  filePath: string
  fileSize: number
  fileType: string
  fileExt: string
  category: string
  uploadUserId: number
  url: string
  createdAt: string
  updatedAt: string
}

/** 文件列表响应 */
export interface FileListResponse {
  list: FileInfo[]
  total: number
  page: number
  pageSize: number
}

/** 批量上传响应 */
export interface BatchUploadResponse {
  success: FileInfo[]
  errors: string[]
}

/**
 * 上传单个文件
 * @param file 文件
 * @param category 分类(image/document/video/audio/other)
 */
export function uploadFile(file: File, category?: string) {
  const formData = new FormData()
  formData.append('file', file)
  if (category) {
    formData.append('category', category)
  }
  return request.post<FileInfo>({
    url: '/api/file/upload',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 批量上传文件
 * @param files 文件列表
 * @param category 分类
 */
export function uploadMultipleFiles(files: File[], category?: string) {
  const formData = new FormData()
  files.forEach((file) => {
    formData.append('files', file)
  })
  if (category) {
    formData.append('category', category)
  }
  return request.post<BatchUploadResponse>({
    url: '/api/file/upload-multiple',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 获取文件列表
 * @param params 查询参数
 */
export function getFileList(params: { page?: number; pageSize?: number; category?: string }) {
  return request.get<FileListResponse>({
    url: '/api/file/list',
    params
  })
}

/**
 * 获取文件详情
 * @param id 文件ID
 */
export function getFileDetail(id: number) {
  return request.get<FileInfo>({
    url: `/api/file/${id}`
  })
}

/**
 * 删除文件
 * @param id 文件ID
 */
export function deleteFile(id: number) {
  return request.delete({
    url: `/api/file/${id}`
  })
}

/**
 * 批量删除文件
 * @param ids 文件ID列表
 */
export function batchDeleteFiles(ids: number[]) {
  return request.post({
    url: '/api/file/batch-delete',
    data: { ids }
  })
}
