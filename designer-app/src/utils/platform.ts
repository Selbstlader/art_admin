/**
 * 平台检测工具
 * Platform detection utilities
 */

import type { PlatformType } from '@/types/common'

/*** 平台检测工具对象 - Platform detection utility object ***/
export const platform = {
  /*** 获取当前平台类型 - Get current platform type ***/
  getPlatform(): PlatformType {
    // #ifdef APP-PLUS
    return 'app'
    // #endif
    
    // #ifdef MP-WEIXIN
    return 'mp-weixin'
    // #endif
    
    // #ifdef H5
    return 'h5'
    // #endif
    
    // 默认返回 h5 - Default return h5
    // @ts-ignore
    return 'h5'
  },

  /*** 是否为 App 端（iOS/Android）- Is App (iOS/Android) ***/
  isApp(): boolean {
    // #ifdef APP-PLUS
    return true
    // #endif
    
    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 是否为微信小程序 - Is WeChat Mini Program ***/
  isMiniProgram(): boolean {
    // #ifdef MP-WEIXIN
    return true
    // #endif
    
    // #ifndef MP-WEIXIN
    return false
    // #endif
  },

  /*** 是否为 H5 端 - Is H5 ***/
  isH5(): boolean {
    // #ifdef H5
    return true
    // #endif
    
    // #ifndef H5
    return false
    // #endif
  },

  /*** 是否为 iOS 系统 - Is iOS ***/
  isIOS(): boolean {
    // #ifdef APP-PLUS
    const systemInfo = uni.getSystemInfoSync()
    return systemInfo.platform === 'ios'
    // #endif
    
    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 是否为 Android 系统 - Is Android ***/
  isAndroid(): boolean {
    // #ifdef APP-PLUS
    const systemInfo = uni.getSystemInfoSync()
    return systemInfo.platform === 'android'
    // #endif
    
    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 是否可以使用 AI 对话功能 - Can use AI chat feature ***/
  canUseAIChat(): boolean {
    // AI 对话功能仅在 App 端可用，小程序不支持
    // AI chat is only available in App, not supported in Mini Program
    // #ifdef APP-PLUS
    return true
    // #endif
    
    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 是否可以使用推送通知 - Can use push notification ***/
  canUsePush(): boolean {
    // 推送通知仅在 App 端可用
    // Push notification is only available in App
    // #ifdef APP-PLUS
    return true
    // #endif
    
    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 是否可以使用订阅消息（小程序）- Can use subscribe message (Mini Program) ***/
  canUseSubscribeMessage(): boolean {
    // 订阅消息仅在小程序端可用
    // Subscribe message is only available in Mini Program
    // #ifdef MP-WEIXIN
    return true
    // #endif
    
    // #ifndef MP-WEIXIN
    return false
    // #endif
  },

  /*** 获取系统信息 - Get system info ***/
  getSystemInfo(): UniApp.GetSystemInfoResult | null {
    try {
      return uni.getSystemInfoSync()
    } catch (error) {
      console.error('获取系统信息失败:', error)
      return null
    }
  },

  /*** 获取状态栏高度 - Get status bar height ***/
  getStatusBarHeight(): number {
    const systemInfo = this.getSystemInfo()
    return systemInfo?.statusBarHeight || 0
  },

  /*** 获取导航栏高度 - Get navigation bar height ***/
  getNavBarHeight(): number {
    const statusBarHeight = this.getStatusBarHeight()
    // 导航栏高度：状态栏 + 44px（iOS）或 48px（Android）
    // Navigation bar height: status bar + 44px (iOS) or 48px (Android)
    const navHeight = this.isIOS() ? 44 : 48
    return statusBarHeight + navHeight
  },

  /*** 获取安全区域底部高度 - Get safe area bottom height ***/
  getSafeAreaBottom(): number {
    const systemInfo = this.getSystemInfo()
    if (!systemInfo) return 0
    
    const { screenHeight, safeArea } = systemInfo
    if (safeArea) {
      return screenHeight - safeArea.bottom
    }
    return 0
  },

  /*** 获取屏幕宽度 - Get screen width ***/
  getScreenWidth(): number {
    const systemInfo = this.getSystemInfo()
    return systemInfo?.screenWidth || 375
  },

  /*** 获取屏幕高度 - Get screen height ***/
  getScreenHeight(): number {
    const systemInfo = this.getSystemInfo()
    return systemInfo?.screenHeight || 667
  },

  /*** 检查网络状态 - Check network status ***/
  async checkNetwork(): Promise<boolean> {
    return new Promise((resolve) => {
      uni.getNetworkType({
        success: (res) => {
          resolve(res.networkType !== 'none')
        },
        fail: () => {
          resolve(false)
        }
      })
    })
  },

  /*** 监听网络状态变化 - Listen to network status change ***/
  onNetworkStatusChange(callback: (isConnected: boolean, networkType: string) => void): void {
    uni.onNetworkStatusChange((res) => {
      callback(res.isConnected, res.networkType)
    })
  },

  /*** 获取平台特定的 tabBar 配置 - Get platform specific tabBar config ***/
  getTabBarConfig(): { showChat: boolean; showMessage: boolean } {
    // App 端显示 AI 对话和消息入口，小程序端隐藏
    // Show AI chat and message entry in App, hide in Mini Program
    const isAppPlatform = this.isApp()
    return {
      showChat: isAppPlatform,
      showMessage: isAppPlatform
    }
  }
}

/*** 导出平台类型常量 - Export platform type constants ***/
export const PLATFORM = {
  APP: 'app' as PlatformType,
  MP_WEIXIN: 'mp-weixin' as PlatformType,
  H5: 'h5' as PlatformType
}
