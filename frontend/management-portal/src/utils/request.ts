import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/store/modules/auth'
import type { ApiResponse } from '@/types/api'

/**
 * 创建 Axios 实例
 */
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

/**
 * 请求拦截器
 * 在请求发送之前添加 Token
 */
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()

    // 如果存在 token，添加到请求头
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 * 统一处理响应和错误
 */
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => response,
  (error) => {
    // 处理 HTTP 错误状态码
    if (error.response) {
      const status = error.response.status
      const data = error.response.data

      switch (status) {
        case 401:
          // 未授权，清除登录状态并跳转到登录页
          ElMessage.error(data?.msg || '登录已过期，请重新登录')
          const authStore = useAuthStore()
          authStore.logout()
          router.push('/login')
          break

        case 403:
          // 权限不足
          ElMessage.error(data?.msg || '权限不足，无法访问')
          break

        case 404:
          // 资源不存在
          ElMessage.error(data?.msg || '请求的资源不存在')
          break

        case 409:
          // 冲突（如任务已被认领）
          ElMessage.warning(data?.msg || '操作冲突，请刷新后重试')
          break

        case 500:
          // 服务器内部错误
          ElMessage.error(data?.msg || '服务器错误，请稍后重试')
          break

        default:
          // 其他错误
          ElMessage.error(data?.msg || `请求失败 (${status})`)
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应
      ElMessage.error('网络错误，请检查您的网络连接')
    } else {
      // 请求配置出错
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default request
