/**
 * 缓存管理工具
 * Cache management utilities
 */

import { CACHE_EXPIRE_TIME } from '@/types/common'
import type { CacheData } from '@/types/store'

/*** 缓存配置选项 - Cache options ***/
export interface CacheOptions {
  expireTime?: number  // 过期时间（毫秒），默认 7 天
}

/*** 缓存键前缀 - Cache key prefix ***/
const CACHE_PREFIX = 'designer_app_cache_'

/*** 缓存管理工具对象 - Cache management utility object ***/
export const cache = {
  /*** 设置缓存 - Set cache ***/
  set<T>(key: string, data: T, options?: CacheOptions): void {
    const expireTime = options?.expireTime || CACHE_EXPIRE_TIME
    const cacheData: CacheData<T> = {
      data,
      timestamp: Date.now(),
      expireTime
    }
    
    try {
      uni.setStorageSync(CACHE_PREFIX + key, JSON.stringify(cacheData))
    } catch (error) {
      console.error('缓存写入失败:', error)
    }
  },

  /*** 获取缓存 - Get cache ***/
  get<T>(key: string): T | null {
    try {
      const cacheStr = uni.getStorageSync(CACHE_PREFIX + key)
      if (!cacheStr) return null
      
      const cacheData: CacheData<T> = JSON.parse(cacheStr)
      
      // 检查是否过期 - Check if expired
      if (this.isExpired(cacheData)) {
        this.remove(key)
        return null
      }
      
      return cacheData.data
    } catch (error) {
      console.error('缓存读取失败:', error)
      return null
    }
  },

  /*** 获取缓存（包含元数据）- Get cache with metadata ***/
  getWithMeta<T>(key: string): CacheData<T> | null {
    try {
      const cacheStr = uni.getStorageSync(CACHE_PREFIX + key)
      if (!cacheStr) return null
      
      return JSON.parse(cacheStr)
    } catch (error) {
      console.error('缓存读取失败:', error)
      return null
    }
  },

  /*** 检查缓存是否过期 - Check if cache is expired ***/
  isExpired<T>(cacheData: CacheData<T>): boolean {
    if (!cacheData || !cacheData.timestamp || !cacheData.expireTime) {
      return true
    }
    return Date.now() - cacheData.timestamp > cacheData.expireTime
  },

  /*** 检查指定键的缓存是否过期 - Check if cache for key is expired ***/
  isKeyExpired(key: string): boolean {
    const cacheData = this.getWithMeta(key)
    if (!cacheData) return true
    return this.isExpired(cacheData)
  },

  /*** 移除缓存 - Remove cache ***/
  remove(key: string): void {
    try {
      uni.removeStorageSync(CACHE_PREFIX + key)
    } catch (error) {
      console.error('缓存删除失败:', error)
    }
  },

  /*** 清除所有应用缓存 - Clear all app cache ***/
  clearAll(): void {
    try {
      const storageInfo = uni.getStorageInfoSync()
      const keys = storageInfo.keys || []
      
      keys.forEach(key => {
        if (key.startsWith(CACHE_PREFIX)) {
          uni.removeStorageSync(key)
        }
      })
    } catch (error) {
      console.error('清除缓存失败:', error)
    }
  },

  /*** 清除过期缓存 - Clear expired cache ***/
  clearExpired(): void {
    try {
      const storageInfo = uni.getStorageInfoSync()
      const keys = storageInfo.keys || []
      
      keys.forEach(key => {
        if (key.startsWith(CACHE_PREFIX)) {
          const cacheStr = uni.getStorageSync(key)
          if (cacheStr) {
            try {
              const cacheData = JSON.parse(cacheStr)
              if (this.isExpired(cacheData)) {
                uni.removeStorageSync(key)
              }
            } catch {
              // 解析失败，删除无效缓存
              uni.removeStorageSync(key)
            }
          }
        }
      })
    } catch (error) {
      console.error('清除过期缓存失败:', error)
    }
  },

  /*** 获取缓存大小 - Get cache size ***/
  getSize(): number {
    try {
      const storageInfo = uni.getStorageInfoSync()
      return storageInfo.currentSize || 0
    } catch (error) {
      console.error('获取缓存大小失败:', error)
      return 0
    }
  },

  /*** 检查缓存是否存在 - Check if cache exists ***/
  has(key: string): boolean {
    try {
      const cacheStr = uni.getStorageSync(CACHE_PREFIX + key)
      return !!cacheStr
    } catch {
      return false
    }
  },

  /*** 更新缓存数据（不改变过期时间）- Update cache data without changing expiration ***/
  update<T>(key: string, updater: (data: T) => T): boolean {
    const cacheData = this.getWithMeta<T>(key)
    if (!cacheData) return false
    
    try {
      cacheData.data = updater(cacheData.data)
      uni.setStorageSync(CACHE_PREFIX + key, JSON.stringify(cacheData))
      return true
    } catch (error) {
      console.error('缓存更新失败:', error)
      return false
    }
  },

  /*** 获取缓存剩余有效时间（毫秒）- Get remaining valid time (ms) ***/
  getRemainingTime(key: string): number {
    const cacheData = this.getWithMeta(key)
    if (!cacheData) return 0
    
    const elapsed = Date.now() - cacheData.timestamp
    const remaining = cacheData.expireTime - elapsed
    return remaining > 0 ? remaining : 0
  }
}

/*** 缓存键常量 - Cache key constants ***/
export const CACHE_KEYS = {
  PROJECT_LIST: 'project_list',
  PROJECT_DETAIL: (id: number) => `project_detail_${id}`,
  MATERIAL_CATEGORIES: 'material_categories',
  MATERIAL_LIST: (category: string) => `material_list_${category}`,
  USER_INFO: 'user_info',
  CHAT_SESSIONS: 'chat_sessions',
  NOTIFICATIONS: 'notifications'
}
