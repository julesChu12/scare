import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

/**
 * 权限指令 v-permission
 *
 * 使用示例：
 * <!-- 单个权限 -->
 * <el-button v-permission="'service:task:claim'">认领任务</el-button>
 *
 * <!-- 多个权限（任一满足） -->
 * <el-button v-permission="['service:task:claim', 'service:task:complete']">处理任务</el-button>
 */
export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    const authStore = useAuthStore()
    const value = binding.value

    if (!value) {
      console.warn('v-permission 指令需要权限码参数')
      return
    }

    const hasPermission = Array.isArray(value)
      ? value.some((code) => authStore.hasPermission(code))
      : authStore.hasPermission(value)

    if (!hasPermission) {
      // 移除元素
      el.parentNode?.removeChild(el)
    }
  },
}
