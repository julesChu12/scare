import { defineStore } from 'pinia'
import { ref } from 'vue'

// 类型定义
export interface User {
  id: number
  phone: string
  role: string
}

export interface EmergencyContact {
  name: string
  phone: string
  relation: string
}

export interface Profile {
  id?: number
  name: string
  // 基本信息
  id_card?: string
  gender?: string
  birth_date?: string
  address?: string
  latitude?: number
  longitude?: number
  // 客户类型和健康信息
  customer_type?: string // elderly/disabled/pregnant/child/other
  health_status?: string
  disability_level?: string
  medical_history?: string
  special_needs?: string
  // 紧急联系人
  emergency_contact?: EmergencyContact
  // 兼容旧字段
  id_number?: string
  user_type?: string
}

export const useUserStore = defineStore('user', () => {
  // State
  const user = ref<User | null>(null)
  const profile = ref<Profile | null>(null)

  // Actions
  function setUser(userData: User) {
    user.value = userData
  }

  function setProfile(profileData: Profile) {
    profile.value = profileData
  }

  function clearUser() {
    user.value = null
    profile.value = null
  }

  function updateProfile(updates: Partial<Profile>) {
    if (profile.value) {
      profile.value = { ...profile.value, ...updates }
    }
  }

  return {
    user,
    profile,
    setUser,
    setProfile,
    clearUser,
    updateProfile
  }
})
