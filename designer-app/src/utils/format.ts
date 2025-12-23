/**
 * 格式化工具函数
 * Format utility functions
 */

/*** 格式化日期 - Format date ***/
export function formatDate(dateStr: string | Date, format: string = 'YYYY-MM-DD'): string {
  if (!dateStr) return '-'
  
  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  
  if (isNaN(date.getTime())) return '-'
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/*** 格式化金额 - Format currency ***/
export function formatCurrency(amount: number, decimals: number = 2): string {
  if (typeof amount !== 'number' || isNaN(amount)) return '¥0.00'
  
  return '¥' + amount.toFixed(decimals).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/*** 格式化面积 - Format area ***/
export function formatArea(area: number): string {
  if (typeof area !== 'number' || isNaN(area)) return '0㎡'
  return `${area}㎡`
}

/*** 格式化文件大小 - Format file size ***/
export function formatFileSize(bytes: number): string {
  if (typeof bytes !== 'number' || isNaN(bytes) || bytes < 0) return '0 B'
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let unitIndex = 0
  let size = bytes
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(unitIndex === 0 ? 0 : 2)} ${units[unitIndex]}`
}

/*** 格式化相对时间 - Format relative time ***/
export function formatRelativeTime(dateStr: string | Date): string {
  if (!dateStr) return '-'
  
  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  
  if (isNaN(date.getTime())) return '-'
  
  const now = Date.now()
  const diff = now - date.getTime()
  
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const week = 7 * day
  const month = 30 * day
  
  if (diff < minute) {
    return '刚刚'
  } else if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`
  } else if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`
  } else if (diff < week) {
    return `${Math.floor(diff / day)}天前`
  } else if (diff < month) {
    return `${Math.floor(diff / week)}周前`
  } else {
    return formatDate(date, 'YYYY-MM-DD')
  }
}

/*** 截断文本 - Truncate text ***/
export function truncateText(text: string, maxLength: number, suffix: string = '...'): string {
  if (!text || text.length <= maxLength) return text || ''
  return text.slice(0, maxLength - suffix.length) + suffix
}
