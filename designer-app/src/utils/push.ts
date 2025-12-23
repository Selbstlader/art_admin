/**
 * 推送通知工具（仅 App 端）
 * Push notification utilities (App only)
 */

import { platform } from './platform'

/*** 声明 UniApp 全局类型 - Declare UniApp global types ***/
declare const uni: any
declare const plus: any

/*** 推送通知管理器 - Push notification manager ***/
export const pushManager = {
  /*** 请求推送权限 - Request push permission ***/
  async requestPermission(): Promise<boolean> {
    // #ifdef APP-PLUS
    try {
      /*** Check if push module is available ***/
      const push = uni.requireNativePlugin?.('push')
      if (!push) {
        console.warn('Push module not available')
        return false
      }

      /*** Request permission based on platform ***/
      if (platform.isIOS()) {
        return await this.requestIOSPermission()
      } else if (platform.isAndroid()) {
        return await this.requestAndroidPermission()
      }
      
      return false
    } catch (error) {
      console.error('请求推送权限失败:', error)
      return false
    }
    // #endif

    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 请求 iOS 推送权限 - Request iOS push permission ***/
  async requestIOSPermission(): Promise<boolean> {
    // #ifdef APP-PLUS
    return new Promise((resolve) => {
      /*** Check current permission status ***/
      plus.push.getClientInfo({
        success: (info: any) => {
          if (info.clientid) {
            console.log('iOS Push ClientID:', info.clientid)
            resolve(true)
          } else {
            /*** Show permission request dialog ***/
            uni.showModal({
              title: '开启通知',
              content: '开启通知权限，及时接收项目更新和审批提醒',
              confirmText: '去开启',
              cancelText: '暂不开启',
              success: (res: any) => {
                if (res.confirm) {
                  /*** Open system settings ***/
                  plus.runtime.openURL('app-settings:')
                }
                resolve(false)
              }
            })
          }
        },
        fail: () => {
          resolve(false)
        }
      })
    })
    // #endif

    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 请求 Android 推送权限 - Request Android push permission ***/
  async requestAndroidPermission(): Promise<boolean> {
    // #ifdef APP-PLUS
    return new Promise((resolve) => {
      /*** Check notification permission on Android ***/
      const main = plus.android.runtimeMainActivity()
      const NotificationManagerCompat = plus.android.importClass('androidx.core.app.NotificationManagerCompat')
      
      if (NotificationManagerCompat) {
        const notificationManager = NotificationManagerCompat.from(main)
        const isEnabled = notificationManager.areNotificationsEnabled()
        
        if (isEnabled) {
          /*** Get push client info ***/
          plus.push.getClientInfo({
            success: (info: any) => {
              console.log('Android Push ClientID:', info.clientid)
              resolve(true)
            },
            fail: () => {
              resolve(false)
            }
          })
        } else {
          /*** Show permission request dialog ***/
          uni.showModal({
            title: '开启通知',
            content: '开启通知权限，及时接收项目更新和审批提醒',
            confirmText: '去开启',
            cancelText: '暂不开启',
            success: (res: any) => {
              if (res.confirm) {
                /*** Open notification settings ***/
                const Intent = plus.android.importClass('android.content.Intent')
                const Settings = plus.android.importClass('android.provider.Settings')
                const Uri = plus.android.importClass('android.net.Uri')
                
                const intent = new Intent(Settings.ACTION_APP_NOTIFICATION_SETTINGS)
                intent.putExtra(Settings.EXTRA_APP_PACKAGE, main.getPackageName())
                main.startActivity(intent)
              }
              resolve(false)
            }
          })
        }
      } else {
        resolve(false)
      }
    })
    // #endif

    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 检查推送权限状态 - Check push permission status ***/
  async checkPermission(): Promise<boolean> {
    // #ifdef APP-PLUS
    return new Promise((resolve) => {
      if (platform.isIOS()) {
        /*** iOS: Check via push client info ***/
        plus.push.getClientInfo({
          success: (info: any) => {
            resolve(!!info.clientid)
          },
          fail: () => {
            resolve(false)
          }
        })
      } else if (platform.isAndroid()) {
        /*** Android: Check notification permission ***/
        try {
          const main = plus.android.runtimeMainActivity()
          const NotificationManagerCompat = plus.android.importClass('androidx.core.app.NotificationManagerCompat')
          
          if (NotificationManagerCompat) {
            const notificationManager = NotificationManagerCompat.from(main)
            resolve(notificationManager.areNotificationsEnabled())
          } else {
            resolve(false)
          }
        } catch (error) {
          console.error('检查 Android 通知权限失败:', error)
          resolve(false)
        }
      } else {
        resolve(false)
      }
    })
    // #endif

    // #ifndef APP-PLUS
    return false
    // #endif
  },

  /*** 获取推送客户端 ID - Get push client ID ***/
  async getClientId(): Promise<string | null> {
    // #ifdef APP-PLUS
    return new Promise((resolve) => {
      plus.push.getClientInfo({
        success: (info: any) => {
          resolve(info.clientid || null)
        },
        fail: () => {
          resolve(null)
        }
      })
    })
    // #endif

    // #ifndef APP-PLUS
    return null
    // #endif
  },

  /*** 监听推送消息 - Listen to push messages ***/
  onPushMessage(callback: (message: PushMessage) => void): void {
    // #ifdef APP-PLUS
    plus.push.addEventListener('click', (msg: any) => {
      /*** Handle push message click ***/
      callback({
        type: 'click',
        title: msg.title || '',
        content: msg.content || '',
        payload: msg.payload
      })
    }, false)

    plus.push.addEventListener('receive', (msg: any) => {
      /*** Handle push message receive ***/
      callback({
        type: 'receive',
        title: msg.title || '',
        content: msg.content || '',
        payload: msg.payload
      })
    }, false)
    // #endif
  },

  /*** 创建本地通知 - Create local notification ***/
  createLocalNotification(options: LocalNotificationOptions): void {
    // #ifdef APP-PLUS
    plus.push.createMessage(options.content, options.payload, {
      title: options.title,
      cover: options.cover ?? true,
      sound: options.sound ?? 'system',
      delay: options.delay ?? 0
    })
    // #endif
  },

  /*** 清除所有通知 - Clear all notifications ***/
  clearAllNotifications(): void {
    // #ifdef APP-PLUS
    plus.push.clear()
    // #endif
  }
}

/*** 微信小程序订阅消息管理器 - WeChat Mini Program subscribe message manager ***/
export const subscribeManager = {
  /*** 请求订阅消息权限 - Request subscribe message permission ***/
  async requestSubscribe(tmplIds: string[]): Promise<SubscribeResult> {
    // #ifdef MP-WEIXIN
    return new Promise((resolve) => {
      uni.requestSubscribeMessage({
        tmplIds,
        success: (res: any) => {
          /*** Check subscription results ***/
          const accepted: string[] = []
          const rejected: string[] = []
          
          tmplIds.forEach((id) => {
            if (res[id] === 'accept') {
              accepted.push(id)
            } else {
              rejected.push(id)
            }
          })
          
          resolve({
            success: accepted.length > 0,
            accepted,
            rejected
          })
        },
        fail: (err: any) => {
          console.error('请求订阅消息失败:', err)
          resolve({
            success: false,
            accepted: [],
            rejected: tmplIds,
            error: err.errMsg
          })
        }
      })
    })
    // #endif

    // #ifndef MP-WEIXIN
    return {
      success: false,
      accepted: [],
      rejected: tmplIds,
      error: 'Not supported on this platform'
    }
    // #endif
  },

  /*** 获取订阅消息设置 - Get subscribe message settings ***/
  async getSettings(): Promise<SubscribeSettings | null> {
    // #ifdef MP-WEIXIN
    return new Promise((resolve) => {
      uni.getSetting({
        withSubscriptions: true,
        success: (res: any) => {
          resolve({
            mainSwitch: res.subscriptionsSetting?.mainSwitch ?? false,
            itemSettings: res.subscriptionsSetting?.itemSettings ?? {}
          })
        },
        fail: () => {
          resolve(null)
        }
      })
    })
    // #endif

    // #ifndef MP-WEIXIN
    return null
    // #endif
  }
}

/*** 推送消息类型 - Push message type ***/
export interface PushMessage {
  type: 'click' | 'receive'
  title: string
  content: string
  payload?: any
}

/*** 本地通知选项 - Local notification options ***/
export interface LocalNotificationOptions {
  title: string
  content: string
  payload?: any
  cover?: boolean
  sound?: string
  delay?: number
}

/*** 订阅结果 - Subscribe result ***/
export interface SubscribeResult {
  success: boolean
  accepted: string[]
  rejected: string[]
  error?: string
}

/*** 订阅设置 - Subscribe settings ***/
export interface SubscribeSettings {
  mainSwitch: boolean
  itemSettings: Record<string, string>
}
