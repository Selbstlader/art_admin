/**
 * 项目状态管理
 * Project state management
 */

import { defineStore } from 'pinia'
import type { ProjectState } from '@/types/store'
import type { Project, ProjectDetail } from '@/types/api'
import { DEFAULT_PAGE_SIZE } from '@/types/common'

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
    /*** 过滤后的项目列表 - Filtered project list ***/
    filteredList: (state) => {
      if (!state.searchKeyword) return state.list
      const keyword = state.searchKeyword.toLowerCase()
      return state.list.filter(p => p.name.toLowerCase().includes(keyword))
    }
  },

  actions: {
    /*** 获取项目列表 - Fetch project list ***/
    async fetchList(refresh: boolean = false): Promise<void> {
      if (refresh) {
        this.currentPage = 1
        this.hasMore = true
      }

      this.loading = true
      try {
        // TODO: 调用后端接口
        // const res = await projectApi.getList({ current: this.currentPage, size: DEFAULT_PAGE_SIZE })
        
        // 模拟数据 - Mock data
        const mockProjects: Project[] = []
        
        if (refresh) {
          this.list = mockProjects
        } else {
          this.list = [...this.list, ...mockProjects]
        }
        
        this.hasMore = mockProjects.length >= DEFAULT_PAGE_SIZE
      } catch (error) {
        console.error('获取项目列表失败:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    /*** 加载更多 - Load more ***/
    async loadMore(): Promise<void> {
      if (!this.hasMore || this.loading) return
      this.currentPage++
      await this.fetchList(false)
    },

    /*** 搜索项目 - Search projects ***/
    setSearchKeyword(keyword: string): void {
      this.searchKeyword = keyword
    },

    /*** 获取项目详情 - Get project detail ***/
    async getDetail(id: number): Promise<ProjectDetail | null> {
      try {
        // TODO: 调用后端接口
        // const res = await projectApi.getDetail(id)
        // this.currentProject = res.data
        
        return this.currentProject
      } catch (error) {
        console.error('获取项目详情失败:', error)
        throw error
      }
    },

    /*** 清除当前项目 - Clear current project ***/
    clearCurrentProject(): void {
      this.currentProject = null
    }
  }
})
