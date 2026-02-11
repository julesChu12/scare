import { checkPermission } from '@/router/guards/permission.guard'
import { useAuthStore } from '@/store/modules/auth'
import { computed } from 'vue'

/**
 * 权限判断 Hook
 */
export function usePermission() {
  const authStore = useAuthStore()

  /**
   * 检查是否有某个权限
   */
  const hasPermission = async (resource: string, action: string): Promise<boolean> => {
    return await checkPermission(`${resource}:${action}`)
  }

  /**
   * 检查是否有任意一个权限
   */
  const hasAnyPermission = async (
    permissions: Array<{ resource: string; action: string }>
  ): Promise<boolean> => {
    const results: boolean[] = await Promise.all(
      permissions.map((p) => checkPermission(`${p.resource}:${p.action}`))
    )
    return results.some((r) => r)
  }

  /**
   * 检查是否有全部权限
   */
  const hasAllPermissions = async (
    permissions: Array<{ resource: string; action: string }>
  ): Promise<boolean> => {
    const results: boolean[] = await Promise.all(
      permissions.map((p) => checkPermission(`${p.resource}:${p.action}`))
    )
    return results.every((r) => r)
  }

  /**
   * 检查用户角色
   */
  const hasRole = (role: string): boolean => {
    return authStore.user?.roles?.includes(role) ?? false
  }

  /**
   * 检查是否有任意一个角色
   */
  const hasAnyRole = (roles: string[]): boolean => {
    const userRoles = authStore.user?.roles || []
    return roles.some((role) => userRoles.includes(role))
  }

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
    hasAnyRole,
  }
}

/**
 * 任务相关权限 Hook
 */
export function useTaskPermission() {
  const { hasPermission } = usePermission()

  const canClaimTask = computed(async () => {
    return await hasPermission('/api/v1/tasks/:id/claim', 'post')
  })

  const canAssignTask = computed(async () => {
    return await hasPermission('/api/v1/tasks/:id/assign', 'post')
  })

  const canCompleteTask = computed(async () => {
    return await hasPermission('/api/v1/tasks/:id/complete', 'post')
  })

  const canReassignTask = computed(async () => {
    return await hasPermission('/api/v1/tasks/:id/reassign', 'post')
  })

  return {
    canClaimTask,
    canAssignTask,
    canCompleteTask,
    canReassignTask,
  }
}
