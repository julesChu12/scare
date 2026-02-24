// 服务类型配置（与后端 consts/service_types.go 保持一致）

export interface ServiceType {
  key: string
  name: string
  icon: string
  description?: string
}

// 后端支持的服务类型枚举
export type ServiceTypeKey =
  | 'meal'
  | 'medical'
  | 'care'
  | 'repair'
  | 'cleaning'
  | 'company'
  | 'emergency'
  | 'shopping'
  | 'transport'
  | 'rehab'
  | 'psychology'
  | 'legal_aid'
  | 'other'

// 全部服务类型（与后端 consts/service_types.go 一致）
export const allServiceTypes: ServiceType[] = [
  { key: 'meal', name: '助餐服务', icon: '🍱', description: '营养配餐送上门' },
  { key: 'medical', name: '就医陪护', icon: '🏥', description: '陪同就医看诊' },
  { key: 'care', name: '日常照护', icon: '👴', description: '生活起居照料' },
  { key: 'repair', name: '居家维修', icon: '🔧', description: '家电水电维修' },
  { key: 'cleaning', name: '家政保洁', icon: '🧹', description: '居家清洁服务' },
  { key: 'company', name: '陪伴聊天', icon: '💬', description: '陪伴交流慰藉' },
  { key: 'emergency', name: '紧急救助', icon: '🚨', description: '突发紧急情况' },
  { key: 'shopping', name: '代买代购', icon: '🛒', description: '生活用品代购' },
  { key: 'transport', name: '出行接送', icon: '🚗', description: '出行陪同接送' },
  { key: 'rehab', name: '康复理疗', icon: '💪', description: '康复训练指导' },
  { key: 'psychology', name: '心理慰藉', icon: '🧠', description: '心理咨询辅导' },
  { key: 'legal_aid', name: '法律援助', icon: '⚖️', description: '法律咨询服务' },
  { key: 'other', name: '其他服务', icon: '📋', description: '其他养老服务' }
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
