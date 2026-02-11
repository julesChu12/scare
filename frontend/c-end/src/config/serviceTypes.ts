// 服务类型配置（与后端 API 保持一致）

export interface ServiceType {
  key: string
  name: string
  icon: string
  description?: string
}

// 后端支持的服务类型枚举
export type ServiceTypeKey =
  | 'daily_care'
  | 'medical_care'
  | 'rehabilitation'
  | 'mental_care'
  | 'housekeeping'
  | 'meal_service'
  | 'other'

// 首页展示的主要服务
export const primaryServiceTypes: ServiceType[] = [
  { key: 'meal_service', name: '助餐服务', icon: '🍱', description: '营养配餐送上门' },
  { key: 'medical_care', name: '医疗护理', icon: '🏥', description: '陪同就医看诊' },
  { key: 'daily_care', name: '日常照护', icon: '👴', description: '生活起居照料' },
  { key: 'housekeeping', name: '家政服务', icon: '🧹', description: '居家清洁服务' },
  { key: 'mental_care', name: '心理关怀', icon: '💬', description: '心理陪伴慰藉' }
]

// 全部服务类型（与后端 API 一致）
export const allServiceTypes: ServiceType[] = [
  { key: 'daily_care', name: '日常照护', icon: '👴', description: '生活起居照料' },
  { key: 'medical_care', name: '医疗护理', icon: '🏥', description: '陪同就医看诊' },
  { key: 'rehabilitation', name: '康复训练', icon: '💪', description: '康复训练指导' },
  { key: 'mental_care', name: '心理关怀', icon: '💬', description: '心理陪伴慰藉' },
  { key: 'housekeeping', name: '家政服务', icon: '🧹', description: '居家清洁服务' },
  { key: 'meal_service', name: '助餐服务', icon: '🍱', description: '营养配餐送上门' },
  { key: 'other', name: '其他', icon: '📋', description: '其他养老服务' }
]

// 获取服务类型名称
export function getServiceTypeName(key: string): string {
  const service = allServiceTypes.find(s => s.key === key)
  return service?.name || '未知服务'
}

// 获取服务类型图标
export function getServiceTypeIcon(key: string): string {
  const service = allServiceTypes.find(s => s.key === key)
  return service?.icon || '📋'
}
