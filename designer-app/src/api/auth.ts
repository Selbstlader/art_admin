/**
 * 认证相关 API
 * Authentication related API
 */

import { request } from '@/utils/request'
import type { ApiResponse, LoginResponse, TokenResponse, UserInfo } from '@/types/api'

/*** Auth API - handles user authentication ***/
export const authApi = {
  /*** 登录 - Login ***/
  login(username: string, password: string): Promise<ApiResponse<LoginResponse>> {
    return request({
      url: '/api/auth/login',
      method: 'POST',
      data: { username, password }
    })
  },

  /*** 退出登录 - Logout ***/
  logout(): Promise<ApiResponse<null>> {
    return request({
      url: '/api/auth/logout',
      method: 'POST'
    })
  },

  /*** 刷新 Token - Refresh token ***/
  refreshToken(refreshToken: string): Promise<ApiResponse<TokenResponse>> {
    return request({
      url: '/api/auth/refresh',
      method: 'POST',
      data: { refreshToken }
    })
  },

  /*** 获取用户信息 - Get user info ***/
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return request({
      url: '/api/auth/user-info',
      method: 'GET'
    })
  }
}
