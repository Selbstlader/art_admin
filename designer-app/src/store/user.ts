/**
 * 用户状态管理
 * User state management
 */

import { defineStore } from 'pinia'
import type { UserState } from '@/types/store'
import type { UserInfo } from '@/types/api'

/*** User store - manages user authentication and profile ***/
export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: '',
    refreshToken: '',
    tokenExpireTime: 0,
    userInfo: null,
    isLoggedIn: false
  }),

  getters: {
    /*** 获取用户头像 - Get user avatar ***/
    avatar: (state) => state.userInfo?.avatar || '',
    
    /*** 获取用户昵称 - Get user nickname ***/
    nickname: (state) => state.userInfo?.nickname || '',
    
    /*** 获取用户角色 - Get user role ***/
    role: (state) => state.userInfo?.role || '',
    
    /*** 检查 token 是否过期 - Check if token is expired ***/
    isTokenExpired: (state) => {
      if (!state.tokenExpireTime) return true
      return Date.now() > state.tokenExpireTime
    }
  },

  actions: {
    /*** 登录 - Login ***/
    async login(username: string, password: string): Promise<void> {
      // TODO: 调用后端接口
      // const res = await authApi.login(username, password)
      
      // 模拟登录响应 - Mock login response
      const mockResponse = {
        accessToken: 'mock_access_token_' + Date.now(),
        refreshToken: 'mock_refresh_token_' + Date.now(),
        expiresIn: 7200,
        user: {
          id: 1,
          username: username,
          nickname: '设计师',
          avatar: '',
          role: 'designer'
        }
      }

      this.setToken(mockResponse.accessToken, mockResponse.refreshToken, mockResponse.expiresIn)
      this.setUserInfo(mockResponse.user)
    },

    /*** 退出登录 - Logout ***/
    async logout(): Promise<void> {
      // TODO: 调用后端接口
      // await authApi.logout()
      
      this.clearUserData()
    },

    /*** 获取用户信息 - Get user info ***/
    async getUserInfo(): Promise<UserInfo | null> {
      if (!this.token) return null
      
      // TODO: 调用后端接口
      // const res = await authApi.getUserInfo()
      // this.setUserInfo(res.data)
      
      return this.userInfo
    },

    /*** 设置 token - Set token ***/
    setToken(accessToken: string, refreshToken: string, expiresIn: number): void {
      this.token = accessToken
      this.refreshToken = refreshToken
      this.tokenExpireTime = Date.now() + expiresIn * 1000
      this.isLoggedIn = true

      // 持久化存储 - Persist to storage
      uni.setStorageSync('token', accessToken)
      uni.setStorageSync('refreshToken', refreshToken)
      uni.setStorageSync('tokenExpireTime', this.tokenExpireTime)
    },

    /*** 设置用户信息 - Set user info ***/
    setUserInfo(userInfo: UserInfo): void {
      this.userInfo = userInfo
      uni.setStorageSync('userInfo', JSON.stringify(userInfo))
    },

    /*** 清除用户数据 - Clear user data ***/
    clearUserData(): void {
      this.token = ''
      this.refreshToken = ''
      this.tokenExpireTime = 0
      this.userInfo = null
      this.isLoggedIn = false

      // 清除本地存储 - Clear local storage
      uni.removeStorageSync('token')
      uni.removeStorageSync('refreshToken')
      uni.removeStorageSync('tokenExpireTime')
      uni.removeStorageSync('userInfo')
    },

    /*** 从本地存储恢复状态 - Restore state from storage ***/
    restoreFromStorage(): void {
      const token = uni.getStorageSync('token')
      const refreshToken = uni.getStorageSync('refreshToken')
      const tokenExpireTime = uni.getStorageSync('tokenExpireTime')
      const userInfoStr = uni.getStorageSync('userInfo')

      if (token && tokenExpireTime && Date.now() < tokenExpireTime) {
        this.token = token
        this.refreshToken = refreshToken || ''
        this.tokenExpireTime = tokenExpireTime
        this.isLoggedIn = true

        if (userInfoStr) {
          try {
            this.userInfo = JSON.parse(userInfoStr)
          } catch (e) {
            console.error('解析用户信息失败:', e)
          }
        }
      } else {
        // Token 过期，清除数据
        this.clearUserData()
      }
    }
  }
})
