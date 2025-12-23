/**
 * 用户状态管理
 * User state management
 */

import { defineStore } from 'pinia'
import type { UserState } from '@/types/store'
import type { UserInfo } from '@/types/api'
import { authApi } from '@/api/auth'

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
    /*** 登录 - Login with real backend API ***/
    async login(username: string, password: string): Promise<void> {
      try {
        // 调用后端登录接口 - Call backend login API
        const res = await authApi.login(username, password)
        
        if (res.code === 200 && res.data) {
          // 兼容 token 和 accessToken 两种字段名
          // Compatible with both token and accessToken field names
          const accessToken = res.data.token || res.data.accessToken || ''
          const refreshToken = res.data.refreshToken || ''
          // 默认 2 小时过期 - Default 2 hours expiration
          const expiresIn = res.data.expiresIn || 7200
          
          this.setToken(accessToken, refreshToken, expiresIn)
          
          // 如果返回了用户信息则设置，否则后续获取
          // Set user info if returned, otherwise fetch later
          if (res.data.user) {
            this.setUserInfo(res.data.user)
          }
        } else {
          throw new Error(res.msg || res.message || '登录失败')
        }
      } catch (error: any) {
        // 清除可能存在的旧数据 - Clear any existing old data
        this.clearUserData()
        throw error
      }
    },

    /*** 退出登录 - Logout with backend API call ***/
    async logout(): Promise<void> {
      try {
        // 调用后端退出接口 - Call backend logout API
        if (this.token) {
          await authApi.logout().catch(() => {
            // 即使后端退出失败，也清除本地数据 - Clear local data even if backend logout fails
            console.warn('后端退出接口调用失败，但仍清除本地数据')
          })
        }
      } finally {
        // 无论如何都清除本地数据 - Always clear local data
        this.clearUserData()
      }
    },

    /*** 获取用户信息 - Get user info from backend ***/
    async getUserInfo(): Promise<UserInfo | null> {
      if (!this.token) return null
      
      try {
        // 调用后端获取用户信息接口 - Call backend get user info API
        const res = await authApi.getUserInfo()
        
        if (res.code === 200 && res.data) {
          this.setUserInfo(res.data)
          return res.data
        }
        return this.userInfo
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // 返回本地缓存的用户信息 - Return cached user info
        return this.userInfo
      }
    },

    /*** 刷新 Token - Refresh token ***/
    async refreshAccessToken(): Promise<boolean> {
      if (!this.refreshToken) return false
      
      try {
        const res = await authApi.refreshToken(this.refreshToken)
        
        if (res.code === 200 && res.data) {
          // 兼容 token 和 accessToken 两种字段名
          // Compatible with both token and accessToken field names
          const accessToken = res.data.token || res.data.accessToken || ''
          const refreshToken = res.data.refreshToken || ''
          const expiresIn = res.data.expiresIn || 7200
          
          this.setToken(accessToken, refreshToken, expiresIn)
          return true
        }
        return false
      } catch (error) {
        console.error('刷新 Token 失败:', error)
        // Token 刷新失败，清除用户数据 - Clear user data on refresh failure
        this.clearUserData()
        return false
      }
    },

    /*** 设置 token - Set token with persistence ***/
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

    /*** 设置用户信息 - Set user info with persistence ***/
    setUserInfo(userInfo: UserInfo): void {
      this.userInfo = userInfo
      uni.setStorageSync('userInfo', JSON.stringify(userInfo))
    },

    /*** 清除用户数据 - Clear all user data from state and storage ***/
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

    /*** 从本地存储恢复状态 - Restore state from local storage ***/
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
      } else if (token && refreshToken) {
        // Token 过期但有 refreshToken，尝试刷新 - Token expired but has refreshToken, try refresh
        this.refreshAccessToken().then(success => {
          if (!success) {
            this.clearUserData()
          }
        })
      } else {
        // Token 过期且无法刷新，清除数据 - Token expired and cannot refresh, clear data
        this.clearUserData()
      }
    },

    /*** 检查登录状态 - Check login status and handle expiration ***/
    checkLoginStatus(): boolean {
      if (!this.token) return false
      
      if (this.isTokenExpired) {
        // Token 过期，提示用户重新登录 - Token expired, prompt user to re-login
        uni.showToast({
          title: '登录已过期，请重新登录',
          icon: 'none',
          duration: 2000
        })
        
        // 延迟跳转登录页 - Delay navigation to login page
        setTimeout(() => {
          this.clearUserData()
          uni.reLaunch({ url: '/pages/login/index' })
        }, 1500)
        
        return false
      }
      
      return true
    }
  }
})
