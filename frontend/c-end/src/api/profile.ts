import client from './client'
import type { EmergencyContact } from '@/store/userStore'

// 接口类型定义
export interface UpdateProfileRequest {
  name?: string
  id_card?: string
  gender?: string
  birth_date?: string
  address?: string
  latitude?: number
  longitude?: number
  customer_type?: string
  health_status?: string
  disability_level?: string
  medical_history?: string
  special_needs?: string
  emergency_contact?: EmergencyContact
  // 兼容旧字段
  id_number?: string
  user_type?: string
}

export interface ProfileResponse {
  id: number
  user_id: number
  name: string
  id_card?: string
  gender?: string
  birth_date?: string
  address?: string
  latitude?: number
  longitude?: number
  customer_type?: string
  health_status?: string
  disability_level?: string
  medical_history?: string
  special_needs?: string
  emergency_contact?: EmergencyContact
  // 兼容旧字段
  id_number?: string
  user_type?: string
  created_at: string
  updated_at: string
}


// API 函数
export const profileAPI = {
  // 更新个人资料
  updateProfile: (data: UpdateProfileRequest) => {
    return client.put<any, ProfileResponse>('/c/profile', data)
  }
}

