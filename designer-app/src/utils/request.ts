/**
 * 网络请求封装
 * Network request wrapper
 */

import type { ApiResponse } from '@/types/api'
import { REQUEST_TIMEOUT, NetworkErrorType } from '@/types/common'

/*** 请求配置选项 - Request options ***/
export interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'OPTIONS' | 'HEAD' | 'TRACE' | 'CONNECT'
  data?: Record<string, any>
  header?: Record<string, string>
  timeout?: number
  showLoading?: boolean
  showError?: boolean
}

/*** API 基础配置 - API base config ***/
// H5 开发模式直接请求后端，其他平台使用完整 URL
// H5 dev mode requests backend directly, other platforms use full URL
function getBaseUrl(): string {
  // @ts-ignore - Vite env types
  const envBaseUrl = typeof import.meta !== 'undefined' && (import.meta as any).env?.VITE_API_BASE_URL
  
  // 开发环境直接使用后端地址，因为 uni.request 不走 Vite 代理
  // Dev environment uses backend URL directly, as uni.request doesn't use Vite proxy
  return envBaseUrl || 'http://localhost:48080'
}

const BASE_URL = getBaseUrl()

/*** 获取存储的 token - Get stored token ***/
function getToken(): string {
  return uni.getStorageSync('token') || ''
}

/*** 检查 token 是否过期 - Check if token is expired ***/
function isTokenExpired(): boolean {
  const expireTime = uni.getStorageSync('tokenExpireTime')
  if (!expireTime) return true
  return Date.now() > expireTime
}

/*** 处理 token 过期 - Handle token expiration ***/
function handleTokenExpired(): void {
  // 清除本地存储 - Clear local storage
  uni.removeStorageSync('token')
  uni.removeStorageSync('refreshToken')
  uni.removeStorageSync('tokenExpireTime')
  uni.removeStorageSync('userInfo')
  
  // 提示用户 - Show toast
  uni.showToast({
    title: '登录已过期，请重新登录',
    icon: 'none',
    duration: 2000
  })
  
  // 跳转登录页 - Navigate to login
  setTimeout(() => {
    uni.reLaunch({ url: '/pages/login/index' })
  }, 1500)
}

/*** 处理网络错误 - Handle network error ***/
function handleNetworkError(errorType: NetworkErrorType, message?: string): void {
  let errorMessage = ''
  
  switch (errorType) {
    case NetworkErrorType.TIMEOUT:
      errorMessage = '请求超时，请检查网络连接'
      break
    case NetworkErrorType.OFFLINE:
      errorMessage = '网络不可用，请检查网络设置'
      break
    case NetworkErrorType.SERVER_ERROR:
      errorMessage = message || '服务器错误，请稍后重试'
      break
    case NetworkErrorType.UNAUTHORIZED:
      handleTokenExpired()
      return
    default:
      errorMessage = message || '请求失败，请重试'
  }
  
  uni.showToast({
    title: errorMessage,
    icon: 'none',
    duration: 2000
  })
}

/*** 请求拦截器 - Request interceptor ***/
function requestInterceptor(options: RequestOptions): RequestOptions {
  const token = getToken()
  
  // 自动携带 token - Auto attach token
  if (token) {
    options.header = {
      ...options.header,
      'Authorization': `Bearer ${token}`
    }
  }
  
  // 设置默认 Content-Type - Set default Content-Type
  options.header = {
    'Content-Type': 'application/json',
    ...options.header
  }
  
  return options
}

/*** 响应拦截器 - Response interceptor ***/
function responseInterceptor<T>(response: UniApp.RequestSuccessCallbackResult): ApiResponse<T> {
  const { statusCode, data } = response
  
  // HTTP 状态码处理 - HTTP status code handling
  if (statusCode === 401) {
    handleTokenExpired()
    throw new Error('登录已过期')
  }
  
  if (statusCode >= 500) {
    handleNetworkError(NetworkErrorType.SERVER_ERROR, '服务器错误')
    throw new Error('服务器错误')
  }
  
  if (statusCode >= 400) {
    const errorData = data as ApiResponse<T>
    throw new Error(errorData?.msg || errorData?.message || '请求失败')
  }
  
  // 业务状态码处理 - Business status code handling
  const responseData = data as ApiResponse<T>
  
  if (responseData.code !== 200 && responseData.code !== 0) {
    // Token 过期 - Token expired
    if (responseData.code === 401 || responseData.code === 4001) {
      handleTokenExpired()
      throw new Error('登录已过期')
    }
    
    throw new Error(responseData.msg || responseData.message || '请求失败')
  }
  
  return responseData
}

/*** 统一请求方法 - Unified request method ***/
export function request<T = any>(options: RequestOptions): Promise<ApiResponse<T>> {
  return new Promise((resolve, reject) => {
    // 检查 token 是否过期（非登录接口）- Check token expiration (non-login APIs)
    if (!options.url.includes('/auth/login') && !options.url.includes('/auth/refresh')) {
      if (isTokenExpired() && getToken()) {
        handleTokenExpired()
        reject(new Error('登录已过期'))
        return
      }
    }
    
    // 应用请求拦截器 - Apply request interceptor
    const interceptedOptions = requestInterceptor(options)
    
    // 显示加载提示 - Show loading
    if (options.showLoading !== false) {
      uni.showLoading({ title: '加载中...', mask: true })
    }
    
    // 发起请求 - Make request
    uni.request({
      url: BASE_URL + interceptedOptions.url,
      method: interceptedOptions.method || 'GET',
      data: interceptedOptions.data,
      header: interceptedOptions.header,
      timeout: interceptedOptions.timeout || REQUEST_TIMEOUT,
      success: (res) => {
        try {
          const result = responseInterceptor<T>(res)
          resolve(result)
        } catch (error: any) {
          if (options.showError !== false) {
            uni.showToast({
              title: error.message || '请求失败',
              icon: 'none',
              duration: 2000
            })
          }
          reject(error)
        }
      },
      fail: (err) => {
        console.error('请求失败:', err)
        
        // 判断错误类型 - Determine error type
        if (err.errMsg?.includes('timeout')) {
          handleNetworkError(NetworkErrorType.TIMEOUT)
          reject(new Error('请求超时'))
        } else if (err.errMsg?.includes('fail')) {
          // 检查网络状态 - Check network status
          uni.getNetworkType({
            success: (networkRes) => {
              if (networkRes.networkType === 'none') {
                handleNetworkError(NetworkErrorType.OFFLINE)
                reject(new Error('网络不可用'))
              } else {
                handleNetworkError(NetworkErrorType.SERVER_ERROR, '连接服务器失败')
                reject(new Error('连接服务器失败'))
              }
            },
            fail: () => {
              handleNetworkError(NetworkErrorType.SERVER_ERROR, '请求失败')
              reject(new Error('请求失败'))
            }
          })
        } else {
          reject(new Error(err.errMsg || '请求失败'))
        }
      },
      complete: () => {
        // 隐藏加载提示 - Hide loading
        if (options.showLoading !== false) {
          uni.hideLoading()
        }
      }
    })
  })
}

/*** 导出请求方法的快捷方式 - Export request method shortcuts ***/
export const http = {
  get<T = any>(url: string, data?: Record<string, any>, options?: Partial<RequestOptions>): Promise<ApiResponse<T>> {
    return request<T>({ url, method: 'GET', data, ...options })
  },
  
  post<T = any>(url: string, data?: Record<string, any>, options?: Partial<RequestOptions>): Promise<ApiResponse<T>> {
    return request<T>({ url, method: 'POST', data, ...options })
  },
  
  put<T = any>(url: string, data?: Record<string, any>, options?: Partial<RequestOptions>): Promise<ApiResponse<T>> {
    return request<T>({ url, method: 'PUT', data, ...options })
  },
  
  delete<T = any>(url: string, data?: Record<string, any>, options?: Partial<RequestOptions>): Promise<ApiResponse<T>> {
    return request<T>({ url, method: 'DELETE', data, ...options })
  },

  /*** OPTIONS 请求 - OPTIONS request ***/
  options<T = any>(url: string, data?: Record<string, any>, options?: Partial<RequestOptions>): Promise<ApiResponse<T>> {
    return request<T>({ url, method: 'OPTIONS', data, ...options })
  }
}
