/**
 * 网络状态管理
 * Network state management
 */

import { defineStore } from 'pinia'

/*** 网络类型枚举 - Network type enum ***/
export type NetworkType = 'wifi' | '2g' | '3g' | '4g' | '5g' | 'ethernet' | 'unknown' | 'none'

/*** 网络状态接口 - Network state interface ***/
export interface NetworkState {
  /*** 是否在线 - Is online ***/
  isOnline: boolean
  /*** 网络类型 - Network type ***/
  networkType: NetworkType
  /*** 上次在线时间 - Last online time ***/
  lastOnlineTime: number
  /*** 是否显示离线提示 - Show offline tip ***/
  showOfflineTip: boolean
  /*** 是否正在恢复 - Is recovering ***/
  isRecovering: boolean
}

/*** Network store - manages network status globally ***/
export const useNetworkStore = defineStore('network', {
  state: (): NetworkState => ({
    isOnline: true,
    networkType: 'unknown',
    lastOnlineTime: Date.now(),
    showOfflineTip: false,
    isRecovering: false
  }),

  getters: {
    /*** 是否为离线状态 - Is offline ***/
    isOffline: (state): boolean => !state.isOnline,

    /*** 获取网络状态描述 - Get network status description ***/
    networkStatusText: (state): string => {
      if (!state.isOnline) return '网络不可用'
      switch (state.networkType) {
        case 'wifi': return 'WiFi'
        case '4g': return '4G'
        case '5g': return '5G'
        case '3g': return '3G'
        case '2g': return '2G'
        case 'ethernet': return '有线网络'
        default: return '已连接'
      }
    },

    /*** 离线持续时间（毫秒）- Offline duration in ms ***/
    offlineDuration: (state): number => {
      if (state.isOnline) return 0
      return Date.now() - state.lastOnlineTime
    }
  },

  actions: {
    /*** 初始化网络状态监听 - Initialize network status listener ***/
    initNetworkListener(): void {
      /*** Check initial network status ***/
      this.checkNetworkStatus()

      /*** Listen to network status changes ***/
      uni.onNetworkStatusChange((res) => {
        this.handleNetworkChange(res.isConnected, res.networkType as NetworkType)
      })
    },

    /*** 检查当前网络状态 - Check current network status ***/
    async checkNetworkStatus(): Promise<boolean> {
      return new Promise((resolve) => {
        uni.getNetworkType({
          success: (res) => {
            const isConnected = res.networkType !== 'none'
            this.updateNetworkStatus(isConnected, res.networkType as NetworkType)
            resolve(isConnected)
          },
          fail: () => {
            this.updateNetworkStatus(false, 'none')
            resolve(false)
          }
        })
      })
    },

    /*** 更新网络状态 - Update network status ***/
    updateNetworkStatus(isOnline: boolean, networkType: NetworkType): void {
      const wasOffline = !this.isOnline
      
      this.isOnline = isOnline
      this.networkType = networkType

      if (isOnline) {
        this.lastOnlineTime = Date.now()
        /*** Hide offline tip when online ***/
        if (wasOffline) {
          this.showOfflineTip = false
        }
      } else {
        /*** Show offline tip when offline ***/
        this.showOfflineTip = true
      }
    },

    /*** 处理网络状态变化 - Handle network status change ***/
    handleNetworkChange(isConnected: boolean, networkType: NetworkType): void {
      const wasOffline = !this.isOnline
      
      this.updateNetworkStatus(isConnected, networkType)

      /*** Network recovered - trigger data refresh ***/
      if (wasOffline && isConnected) {
        this.onNetworkRecovered()
      }
    },

    /*** 网络恢复时的处理 - Handle network recovery ***/
    async onNetworkRecovered(): Promise<void> {
      this.isRecovering = true
      
      /*** Show recovery toast ***/
      uni.showToast({
        title: '网络已恢复',
        icon: 'none',
        duration: 2000
      })

      /*** Emit custom event for pages to refresh data ***/
      uni.$emit('network:recovered')

      /*** Wait a moment then hide recovering state ***/
      setTimeout(() => {
        this.isRecovering = false
      }, 1000)
    },

    /*** 显示离线提示 - Show offline tip ***/
    showOfflineToast(): void {
      if (!this.isOnline) {
        uni.showToast({
          title: '当前为离线模式',
          icon: 'none',
          duration: 2000
        })
      }
    },

    /*** 设置离线提示显示状态 - Set offline tip visibility ***/
    setOfflineTipVisible(visible: boolean): void {
      this.showOfflineTip = visible
    },

    /*** 隐藏离线提示 - Hide offline tip ***/
    hideOfflineTip(): void {
      this.showOfflineTip = false
    },

    /*** 重置状态 - Reset state ***/
    resetState(): void {
      this.isOnline = true
      this.networkType = 'unknown'
      this.lastOnlineTime = Date.now()
      this.showOfflineTip = false
      this.isRecovering = false
    }
  }
})
