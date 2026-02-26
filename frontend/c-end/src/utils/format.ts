/**
 * 格式化日期时间工具函数
 * 统一输出格式：YYYY-MM-DD HH:mm:ss
 */

/**
 * 格式化日期时间为完整格式
 * @param time ISO 8601 时间字符串或时间戳
 * @returns YYYY-MM-DD HH:mm:ss 格式字符串
 */
export function formatDateTime(time: string | number | Date | null | undefined): string {
  if (!time) return '-'
  
  const date = new Date(time)
  if (isNaN(date.getTime())) return String(time)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

/**
 * 格式化日期（不含时间）
 * @param time ISO 8601 时间字符串或时间戳
 * @returns YYYY-MM-DD 格式字符串
 */
export function formatDate(time: string | number | Date | null | undefined): string {
  if (!time) return '-'
  
  const date = new Date(time)
  if (isNaN(date.getTime())) return String(time)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
}

/**
 * 格式化时间（不含日期）
 * @param time ISO 8601 时间字符串或时间戳
 * @returns HH:mm:ss 格式字符串
 */
export function formatTime(time: string | number | Date | null | undefined): string {
  if (!time) return '-'
  
  const date = new Date(time)
  if (isNaN(date.getTime())) return String(time)
  
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${hours}:${minutes}:${seconds}`
}

/**
 * 格式化为简短日期时间（不含秒）
 * @param time ISO 8601 时间字符串或时间戳
 * @returns YYYY-MM-DD HH:mm 格式字符串
 */
export function formatDateTimeShort(time: string | number | Date | null | undefined): string {
  if (!time) return '-'
  
  const date = new Date(time)
  if (isNaN(date.getTime())) return String(time)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}
