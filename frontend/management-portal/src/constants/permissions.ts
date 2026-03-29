/**
 * 权限码常量定义
 * 与后端 permissions 表中的 code 字段保持一致
 */

// ==================== 服务管理 ====================
// 需求管理
export const PERM_REQUEST_CREATE = 'service:request:create'
export const PERM_REQUEST_UPDATE = 'service:request:update'
export const PERM_REQUEST_DELETE = 'service:request:delete'

// 任务管理
export const PERM_TASK_CLAIM = 'service:task:claim'
export const PERM_TASK_COMPLETE = 'service:task:complete'
export const PERM_TASK_ASSIGN = 'service:task:assign'
export const PERM_TASK_TRANSFER = 'service:task:transfer'

// ==================== 系统管理 ====================
// 用户管理
export const PERM_USER_CREATE = 'system:user:create'
export const PERM_USER_UPDATE = 'system:user:update'
export const PERM_USER_DELETE = 'system:user:delete'
export const PERM_USER_ROLE_ASSIGN = 'system:user:role:assign'

// 角色权限管理
export const PERM_ROLE_PERMISSION_UPDATE = 'system:role:permission:update'

// 菜单管理
export const PERM_MENU_CREATE = 'system:menu:create'
export const PERM_MENU_UPDATE = 'system:menu:update'
export const PERM_MENU_DELETE = 'system:menu:delete'

// ==================== 内容管理 ====================
// 新闻管理
export const PERM_NEWS_CREATE = 'content:news:create'
export const PERM_NEWS_UPDATE = 'content:news:update'
export const PERM_NEWS_DELETE = 'content:news:delete'
export const PERM_NEWS_PUBLISH = 'content:news:publish'
export const PERM_NEWS_ARCHIVE = 'content:news:archive'

// 轮播图管理
export const PERM_BANNER_CREATE = 'content:banner:create'
export const PERM_BANNER_UPDATE = 'content:banner:update'
export const PERM_BANNER_DELETE = 'content:banner:delete'

// ==================== 站点管理 ====================
// 站点
export const PERM_STATION_CREATE = 'station:list:create'
export const PERM_STATION_UPDATE = 'station:list:update'
export const PERM_STATION_DELETE = 'station:list:delete'

// 围栏
export const PERM_ZONE_CREATE = 'station:zone:create'
export const PERM_ZONE_UPDATE = 'station:zone:update'
export const PERM_ZONE_DELETE = 'station:zone:delete'

// ==================== 居民管理 ====================
// 老年人档案
export const PERM_ELDERLY_CREATE = 'resident:elderly:create'
export const PERM_ELDERLY_UPDATE = 'resident:elderly:update'
export const PERM_ELDERLY_DELETE = 'resident:elderly:delete'
export const PERM_ELDERLY_CONTACT_UPDATE = 'resident:elderly:contact:update'
export const PERM_ELDERLY_CONTACT_DELETE = 'resident:elderly:contact:delete'
