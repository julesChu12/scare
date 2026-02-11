import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import type { LoginRequest } from '@/types/api'

/**
 * B端用户信息（本地存储格式）
 */
interface AuthUser {
  id: number
  name: string
  phone: string
  roles: string[]
  station_id: number | null
  status: string
  permissions?: string[]
}

export const useAuthStore = defineStore('auth', () => {
  const B_TOKEN_KEY = 'b_token'
  const B_REFRESH_TOKEN_KEY = 'b_refresh_token'
  const B_USER_KEY = 'b_user'

  // State
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const user = ref<AuthUser | null>(null)

  // Getters
  const isLoggedIn = computed(() => !!token.value)

  /**
   * 检查用户是否拥有指定权限
   * Admin 角色拥有所有权限
   */
  const hasPermission = (permission: string) => {
    // Admin 拥有所有权限
    if (user.value?.roles?.includes('admin')) {
      return true
    }
    if (!user.value?.permissions) return false
    return user.value.permissions.includes(permission)
  }

  /**
   * 检查用户是否拥有指定角色
   */
  const hasRole = (role: string) => {
    if (!user.value?.roles) return false
    return user.value.roles.includes(role)
  }

  /**
   * 获取用户的主要角色
   */
  const primaryRole = computed(() => {
    if (!user.value?.roles?.length) return null
    return user.value.roles[0]
  })

  // Actions

  /**
   * 登录
   */
  async function login(credentials: LoginRequest) {
    const response = await authApi.login(credentials)
    const { token: accessToken, refresh_token, user_id, roles, identities, name, phone, station_id, status } = response.data

    // 保存 token
    setToken(accessToken)
    setRefreshToken(refresh_token)

    // 保存用户信息
    const userData: AuthUser = {
      id: user_id,
      roles: roles || identities || [],
      name,
      phone,
      station_id,
      status,
    }
    setUser(userData)

    // 获取用户权限
    await fetchUserPermissions()

    return userData
  }

  /**
   * 获取当前用户信息和权限
   */
  async function fetchUserPermissions() {
    try {
      const response = await authApi.getCurrentUser()
      const { user: userData, permissions } = response.data

      // 更新用户信息和权限
      const updatedUser: AuthUser = {
        id: userData.id,
        name: userData.name,
        phone: userData.phone,
        roles: userData.roles || userData.identities || [],
        station_id: userData.station_id,
        status: userData.status,
        permissions: permissions || [],
      }
      setUser(updatedUser)

      return updatedUser
    } catch (error) {
      console.error('Failed to fetch user permissions:', error)
      throw error
    }
  }

  /**
   * 登出
   */
  async function logout() {
    try {
      // 调用后端登出接口
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('Logout API error:', error)
    } finally {
      // 清除本地状态
      clearAuth()
    }
  }

  /**
   * 清除认证状态
   */
  function clearAuth() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem(B_TOKEN_KEY)
    localStorage.removeItem(B_REFRESH_TOKEN_KEY)
    localStorage.removeItem(B_USER_KEY)
  }

  /**
   * 设置 Token
   */
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem(B_TOKEN_KEY, newToken)
  }

  /**
   * 设置 Refresh Token
   */
  function setRefreshToken(newRefreshToken: string) {
    refreshToken.value = newRefreshToken
    localStorage.setItem(B_REFRESH_TOKEN_KEY, newRefreshToken)
  }

  /**
   * 设置用户信息
   */
  function setUser(newUser: AuthUser) {
    user.value = newUser
    localStorage.setItem(B_USER_KEY, JSON.stringify(newUser))
  }

  /**
   * 刷新 Token
   */
  async function refreshAccessToken() {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await authApi.refreshToken(refreshToken.value)
      const { token: newToken, refresh_token: newRefreshToken } = response.data

      setToken(newToken)
      setRefreshToken(newRefreshToken)

      return newToken
    } catch (error) {
      // 刷新失败，清除认证状态
      clearAuth()
      throw error
    }
  }

  /**
   * 初始化：从 localStorage 恢复状态
   */
  function init() {
    const savedToken = localStorage.getItem(B_TOKEN_KEY)
    const savedRefreshToken = localStorage.getItem(B_REFRESH_TOKEN_KEY)
    const savedUser = localStorage.getItem(B_USER_KEY)

    if (savedToken) {
      token.value = savedToken
    }

    if (savedRefreshToken) {
      refreshToken.value = savedRefreshToken
    }

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (error) {
        console.error('Failed to parse user data from localStorage:', error)
        clearAuth()
      }
    }
  }

  return {
    // State
    token,
    refreshToken,
    user,
    // Getters
    isLoggedIn,
    primaryRole,
    // Methods
    hasPermission,
    hasRole,
    // Actions
    login,
    logout,
    clearAuth,
    setToken,
    setRefreshToken,
    setUser,
    refreshAccessToken,
    fetchUserPermissions,
    init,
  }
})
