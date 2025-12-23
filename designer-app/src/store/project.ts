/**
 * 项目状态管理
 * Project state management
 */

import { defineStore } from 'pinia'
import type { ProjectState, ProjectListCache } from '@/types/store'
import type { Project, ProjectDetail } from '@/types/api'
import { DEFAULT_PAGE_SIZE } from '@/types/common'
import { projectApi } from '@/api/project'
import { cache, CACHE_KEYS } from '@/utils/cache'

/*** Project store - manages project list and details ***/
export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    list: [],
    currentProject: null,
    total: 0,
    loading: false,
    searchKeyword: '',
    currentPage: 1,
    hasMore: true
  }),

  getters: {
    /*** 过滤后的项目列表 - Filtered project list based on search keyword ***/
    filteredList: (state): Project[] => {
      if (!state.searchKeyword) return state.list
      const keyword = state.searchKeyword.toLowerCase()
      return state.list.filter(p => p.name.toLowerCase().includes(keyword))
    },

    /*** 是否为空列表 - Check if list is empty ***/
    isEmpty: (state): boolean => {
      return state.list.length === 0 && !state.loading
    },

    /*** 是否可以加载更多 - Check if can load more ***/
    canLoadMore: (state): boolean => {
      return state.hasMore && !state.loading
    }
  },

  actions: {
    /*** 获取项目列表 - Fetch project list from API or cache ***/
    async fetchList(refresh: boolean = false): Promise<void> {
      if (refresh) {
        this.currentPage = 1
        this.hasMore = true
      }

      this.loading = true
      
      try {
        // 检查是否有缓存数据 - Check for cached data
        if (!refresh) {
          const cachedData = this.loadFromCache()
          if (cachedData && cachedData.length > 0) {
            this.list = cachedData
            this.loading = false
            return
          }
        }

        // 调用后端接口 - Call backend API
        const res = await projectApi.getList({ 
          current: this.currentPage, 
          size: DEFAULT_PAGE_SIZE 
        })
        
        if (res.code === 200 && res.data) {
          const projects = res.data.records || []
          this.total = res.data.total || 0
          
          if (refresh) {
            this.list = projects
          } else {
            this.list = [...this.list, ...projects]
          }
          
          this.hasMore = projects.length >= DEFAULT_PAGE_SIZE
          
          // 保存到缓存 - Save to cache
          this.saveToCache()
        } else {
          throw new Error(res.msg || res.message || '获取项目列表失败')
        }
      } catch (error) {
        console.error('获取项目列表失败:', error)
        
        // 尝试从缓存加载 - Try to load from cache on error
        const cachedData = this.loadFromCache()
        if (cachedData && cachedData.length > 0) {
          this.list = cachedData
          uni.showToast({ 
            title: '网络异常，已加载缓存数据', 
            icon: 'none',
            duration: 2000
          })
        } else {
          throw error
        }
      } finally {
        this.loading = false
      }
    },

    /*** 加载更多项目 - Load more projects ***/
    async loadMore(): Promise<void> {
      if (!this.hasMore || this.loading) return
      
      this.currentPage++
      
      try {
        const res = await projectApi.getList({ 
          current: this.currentPage, 
          size: DEFAULT_PAGE_SIZE,
          name: this.searchKeyword || undefined
        })
        
        if (res.code === 200 && res.data) {
          const projects = res.data.records || []
          this.list = [...this.list, ...projects]
          this.hasMore = projects.length >= DEFAULT_PAGE_SIZE
          
          // 更新缓存 - Update cache
          this.saveToCache()
        }
      } catch (error) {
        console.error('加载更多项目失败:', error)
        this.currentPage-- // 回退页码 - Rollback page number
        throw error
      }
    },

    /*** 搜索项目 - Search projects by keyword ***/
    async search(keyword: string): Promise<void> {
      this.searchKeyword = keyword
      this.currentPage = 1
      this.hasMore = true
      
      if (!keyword) {
        // 无关键字时从缓存或重新获取 - Load from cache or fetch when no keyword
        const cachedData = this.loadFromCache()
        if (cachedData && cachedData.length > 0) {
          this.list = cachedData
          return
        }
        await this.fetchList(true)
        return
      }
      
      this.loading = true
      
      try {
        const res = await projectApi.search(keyword, this.currentPage, DEFAULT_PAGE_SIZE)
        
        if (res.code === 200 && res.data) {
          this.list = res.data.records || []
          this.total = res.data.total || 0
          this.hasMore = this.list.length >= DEFAULT_PAGE_SIZE
        }
      } catch (error) {
        console.error('搜索项目失败:', error)
        // 本地过滤作为降级方案 - Local filter as fallback
        const cachedData = this.loadFromCache()
        if (cachedData) {
          const lowerKeyword = keyword.toLowerCase()
          this.list = cachedData.filter(p => 
            p.name.toLowerCase().includes(lowerKeyword)
          )
        }
      } finally {
        this.loading = false
      }
    },

    /*** 设置搜索关键字（本地过滤）- Set search keyword for local filtering ***/
    setSearchKeyword(keyword: string): void {
      this.searchKeyword = keyword
    },

    /*** 获取项目详情 - Get project detail ***/
    async getDetail(id: number): Promise<ProjectDetail | null> {
      try {
        // 先检查缓存 - Check cache first
        const cacheKey = CACHE_KEYS.PROJECT_DETAIL(id)
        const cachedDetail = cache.get<ProjectDetail>(cacheKey)
        if (cachedDetail) {
          this.currentProject = cachedDetail
          return cachedDetail
        }

        // 调用后端接口 - Call backend API
        const res = await projectApi.getDetail(id)
        
        if (res.code === 200 && res.data) {
          this.currentProject = res.data
          // 缓存详情数据 - Cache detail data
          cache.set(cacheKey, res.data)
          return res.data
        }
        
        return null
      } catch (error) {
        console.error('获取项目详情失败:', error)
        throw error
      }
    },

    /*** 清除当前项目 - Clear current project ***/
    clearCurrentProject(): void {
      this.currentProject = null
    },

    /*** 保存项目列表到缓存 - Save project list to cache ***/
    saveToCache(): void {
      if (this.list.length > 0) {
        const cacheData: ProjectListCache = {
          data: this.list,
          timestamp: Date.now(),
          expireTime: 7 * 24 * 60 * 60 * 1000, // 7 days
          total: this.total,
          lastPage: this.currentPage
        }
        try {
          uni.setStorageSync('designer_app_cache_' + CACHE_KEYS.PROJECT_LIST, JSON.stringify(cacheData))
        } catch (error) {
          console.error('保存项目缓存失败:', error)
        }
      }
    },

    /*** 从缓存加载项目列表 - Load project list from cache ***/
    loadFromCache(): Project[] | null {
      try {
        const cacheStr = uni.getStorageSync('designer_app_cache_' + CACHE_KEYS.PROJECT_LIST)
        if (!cacheStr) return null
        
        const cacheData: ProjectListCache = JSON.parse(cacheStr)
        
        // 检查是否过期 - Check if expired
        if (Date.now() - cacheData.timestamp > cacheData.expireTime) {
          this.clearCache()
          return null
        }
        
        this.total = cacheData.total
        return cacheData.data
      } catch (error) {
        console.error('加载项目缓存失败:', error)
        return null
      }
    },

    /*** 检查缓存是否过期 - Check if cache is expired ***/
    isCacheExpired(): boolean {
      try {
        const cacheStr = uni.getStorageSync('designer_app_cache_' + CACHE_KEYS.PROJECT_LIST)
        if (!cacheStr) return true
        
        const cacheData: ProjectListCache = JSON.parse(cacheStr)
        return Date.now() - cacheData.timestamp > cacheData.expireTime
      } catch {
        return true
      }
    },

    /*** 清除项目缓存 - Clear project cache ***/
    clearCache(): void {
      try {
        uni.removeStorageSync('designer_app_cache_' + CACHE_KEYS.PROJECT_LIST)
      } catch (error) {
        console.error('清除项目缓存失败:', error)
      }
    },

    /*** 重置状态 - Reset store state ***/
    resetState(): void {
      this.list = []
      this.currentProject = null
      this.total = 0
      this.loading = false
      this.searchKeyword = ''
      this.currentPage = 1
      this.hasMore = true
    }
  }
})
