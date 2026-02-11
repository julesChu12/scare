import { computed } from 'vue'
import type { Profile } from '@/store/userStore'

/**
 * 计算资料完善度
 * @param profile 用户资料
 * @returns 完善度百分比 (0-100)
 */
export function useProfileCompleteness(profile: Profile | null) {
  const completeness = computed(() => {
    if (!profile) return 0

    const fields = [
      // 必填字段 (权重更高)
      { key: 'name', weight: 2 },
      { key: 'gender', weight: 1 },
      { key: 'birth_date', weight: 1 },

      // 重要字段
      { key: 'address', weight: 1.5 },
      { key: 'emergency_contact', weight: 2, check: (val: any) => val?.name && val?.phone },

      // 可选字段
      { key: 'id_card', weight: 1 },
      { key: 'health_status', weight: 1 },
      { key: 'medical_history', weight: 0.5 },
      { key: 'special_needs', weight: 0.5 }
    ]

    let totalWeight = 0
    let filledWeight = 0

    fields.forEach(field => {
      totalWeight += field.weight
      const value = (profile as any)[field.key]

      if (field.check) {
        if (field.check(value)) {
          filledWeight += field.weight
        }
      } else {
        if (value && value !== '') {
          filledWeight += field.weight
        }
      }
    })

    return Math.round((filledWeight / totalWeight) * 100)
  })

  const completenessText = computed(() => {
    const percent = completeness.value
    if (percent >= 90) return '资料完善'
    if (percent >= 70) return '资料较完善'
    if (percent >= 50) return '资料待完善'
    return '请完善资料'
  })

  const completenessColor = computed(() => {
    const percent = completeness.value
    if (percent >= 90) return '#67C23A'
    if (percent >= 70) return '#409EFF'
    if (percent >= 50) return '#E6A23C'
    return '#F56C6C'
  })

  return {
    completeness,
    completenessText,
    completenessColor
  }
}
