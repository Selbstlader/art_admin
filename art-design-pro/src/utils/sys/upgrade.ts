import { ElNotification } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { StorageConfig } from '@/utils/storage/storage-config'

/**
 * 版本管理器
 * 负责处理版本比较、升级检测和数据清理
 */
class VersionManager {
  /**
   * 获取存储的版本号
   */
  private getStoredVersion(): string | null {
    return localStorage.getItem(StorageConfig.VERSION_KEY)
  }

  /**
   * 设置版本号到存储
   */
  private setStoredVersion(version: string): void {
    localStorage.setItem(StorageConfig.VERSION_KEY, version)
  }

  /**
   * 检查是否应该跳过升级处理
   */
  private shouldSkipUpgrade(): boolean {
    return StorageConfig.CURRENT_VERSION === StorageConfig.SKIP_UPGRADE_VERSION
  }

  /**
   * 检查是否为首次访问
   */
  private isFirstVisit(storedVersion: string | null): boolean {
    return !storedVersion
  }

  /**
   * 检查版本是否相同
   */
  private isSameVersion(storedVersion: string): boolean {
    return storedVersion === StorageConfig.CURRENT_VERSION
  }

  /**
   * 查找旧的存储结构
   */
  private findLegacyStorage(): { oldSysKey: string | null; oldVersionKeys: string[] } {
    const storageKeys = Object.keys(localStorage)
    const currentVersionPrefix = StorageConfig.generateStorageKey('').slice(0, -1)

    const oldSysKey =
      storageKeys.find(
        (key) =>
          StorageConfig.isVersionedKey(key) && key !== currentVersionPrefix && !key.includes('-')
      ) || null

    const oldVersionKeys = storageKeys.filter(
      (key) =>
        StorageConfig.isVersionedKey(key) &&
        !StorageConfig.isCurrentVersionKey(key) &&
        key.includes('-')
    )

    return { oldSysKey, oldVersionKeys }
  }

  /**
   * 显示升级通知
   */
  private showUpgradeNotification(): void {
    ElNotification({
      title: '系统升级公告',
      message: `系统已升级到 ${StorageConfig.CURRENT_VERSION} 版本`,
      duration: 5000,
      type: 'success'
    })
  }

  /**
   * 清理旧版本数据
   */
  private cleanupLegacyData(oldSysKey: string | null, oldVersionKeys: string[]): void {
    if (oldSysKey) {
      localStorage.removeItem(oldSysKey)
    }
    oldVersionKeys.forEach((key) => {
      localStorage.removeItem(key)
    })
  }

  /**
   * 执行升级流程
   */
  private async executeUpgrade(
    storedVersion: string,
    legacyStorage: ReturnType<typeof this.findLegacyStorage>
  ): Promise<void> {
    try {
      this.showUpgradeNotification()
      this.setStoredVersion(StorageConfig.CURRENT_VERSION)
      this.cleanupLegacyData(legacyStorage.oldSysKey, legacyStorage.oldVersionKeys)
      console.info(`[Upgrade] 升级完成: ${storedVersion} → ${StorageConfig.CURRENT_VERSION}`)
    } catch (error) {
      console.error('[Upgrade] 系统升级处理失败:', error)
    }
  }

  /**
   * 系统升级处理主流程
   */
  async processUpgrade(): Promise<void> {
    if (this.shouldSkipUpgrade()) {
      return
    }

    const storedVersion = this.getStoredVersion()

    if (this.isFirstVisit(storedVersion)) {
      this.setStoredVersion(StorageConfig.CURRENT_VERSION)
      return
    }

    if (this.isSameVersion(storedVersion!)) {
      return
    }

    const legacyStorage = this.findLegacyStorage()
    if (!legacyStorage.oldSysKey && legacyStorage.oldVersionKeys.length === 0) {
      this.setStoredVersion(StorageConfig.CURRENT_VERSION)
      return
    }

    setTimeout(() => {
      this.executeUpgrade(storedVersion!, legacyStorage)
    }, StorageConfig.UPGRADE_DELAY)
  }
}

const versionManager = new VersionManager()

/**
 * 系统升级处理入口函数
 */
export async function systemUpgrade(): Promise<void> {
  await versionManager.processUpgrade()
}
