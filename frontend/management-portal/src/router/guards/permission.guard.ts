import type { Router } from 'vue-router'
import { useAuthStore } from '@/store/modules/auth'
import { ElMessage } from 'element-plus'

// checkPermission is used by composables/directives to decide whether to show UI actions.
// - not logged in => false
// - has specific permission => true
// - admin => true (admin has all permissions)
// - otherwise => rely on backend 403 for final enforcement
export async function checkPermission(permission: string): Promise<boolean> {
  const authStore = useAuthStore()
  if (!authStore.isLoggedIn) return false
  return authStore.hasPermission(permission)
}

/**
 * 设置权限守卫
 */
export function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // 公开路由（无需认证）
    if (to.meta.public) {
      // 如果已登录且访问登录页，跳转到首页
      if (to.path === '/login' && authStore.isLoggedIn) {
        return next('/')
      }
      return next()
    }

    // 检查登录状态
    if (!authStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      return next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
    }

    // 如果已登录但没有权限数据（undefined），重新获取
    // 注意：空数组 [] 表示已获取但无权限，不应重复获取
    if (!authStore.user?.permissions) {
      try {
        await authStore.fetchUserPermissions()
      } catch (error) {
        console.error('获取用户权限失败:', error)
        // 权限获取失败，清除登录状态
        authStore.clearAuth()
        return next({
          path: '/login',
          query: { redirect: to.fullPath },
        })
      }
    }

    // 权限检查（使用新的 permission_code 字段）
    const requiredPermission = (to.meta.permission_code || to.meta.permission) as string | undefined
    if (requiredPermission) {
      if (!authStore.hasPermission(requiredPermission)) {
        ElMessage.error('权限不足，无法访问此页面')
        return next(from.path || '/')
      }
    }

    // 兼容旧的角色检查（逐步废弃）
    const requiredRoles = to.meta.roles as string[] | undefined
    if (requiredRoles && requiredRoles.length > 0) {
      const userRoles = authStore.user?.roles || []
      const hasRequiredRole = requiredRoles.some((role) => userRoles.includes(role))
      if (!hasRequiredRole) {
        ElMessage.error('权限不足，无法访问此页面')
        return next(from.path || '/')
      }
    }

    next()
  })

  // 全局后置守卫：设置页面标题
  router.afterEach((to) => {
    const title = to.meta.title as string
    document.title = title ? `${title} - ${import.meta.env.VITE_APP_TITLE || 'sCare 管理后台'}` : import.meta.env.VITE_APP_TITLE || 'sCare 管理后台'
  })
}
