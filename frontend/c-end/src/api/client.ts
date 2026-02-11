import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

const TOKEN_KEY = 'c_token'
const REFRESH_TOKEN_KEY = 'c_refresh_token'

// 是否正在刷新 Token
let isRefreshing = false
// 等待刷新的请求队列
let refreshSubscribers: ((token: string) => void)[] = []

// 订阅 Token 刷新
function subscribeTokenRefresh(callback: (token: string) => void) {
  refreshSubscribers.push(callback)
}

// 通知所有订阅者
function onTokenRefreshed(token: string) {
  refreshSubscribers.forEach(callback => callback(token))
  refreshSubscribers = []
}

// 创建 Axios 实例
const client: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：添加 Token
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem(TOKEN_KEY)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一错误处理 + Token 刷新
client.interceptors.response.use(
  (response) => {
    const body = response.data as any

    // Backend uses unified wrapper: { msg, data }
    if (body && typeof body === 'object' && 'data' in body && 'msg' in body) {
      return body.data
    }

    return body
  },
  async (error: AxiosError<any>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    // 处理错误响应
    if (error.response) {
      const { status, data } = error.response

      // 开发模式 mock 登录：跳过 401 跳转
      const DEV_MOCK_LOGIN = false // 生产环境
      if (import.meta.env.DEV && DEV_MOCK_LOGIN && status === 401) {
        console.warn('🔧 开发模式：跳过 401 跳转')
        return Promise.reject(error)
      }

      // 401 未授权：尝试刷新 Token
      if (status === 401 && !originalRequest._retry) {
        const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)

        // 如果没有 refresh token 或者是刷新请求本身失败，直接跳转登录
        if (!refreshToken || originalRequest.url?.includes('/auth/refresh')) {
          clearTokensAndRedirect()
          return Promise.reject(error)
        }

        // 如果正在刷新，将请求加入队列
        if (isRefreshing) {
          return new Promise((resolve) => {
            subscribeTokenRefresh((token: string) => {
              originalRequest.headers.Authorization = `Bearer ${token}`
              resolve(client(originalRequest))
            })
          })
        }

        // 开始刷新 Token
        originalRequest._retry = true
        isRefreshing = true

        try {
          const response = await axios.post('/api/v1/c/auth/refresh', {
            refresh_token: refreshToken
          })

          const result = response.data?.data || response.data
          const newToken = result.token
          const newRefreshToken = result.refresh_token

          // 保存新 Token
          localStorage.setItem(TOKEN_KEY, newToken)
          if (newRefreshToken) {
            localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)
          }

          // 通知所有等待的请求
          onTokenRefreshed(newToken)

          // 重试原请求
          originalRequest.headers.Authorization = `Bearer ${newToken}`
          return client(originalRequest)
        } catch (refreshError) {
          // 刷新失败，清除 Token 并跳转登录
          clearTokensAndRedirect()
          return Promise.reject(refreshError)
        } finally {
          isRefreshing = false
        }
      }

      // 其他错误：显示错误消息
      const message = data?.msg || data?.message || data?.error || '请求失败'
      ElMessage.error(message)
    } else if (error.request) {
      // 网络错误
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

// 清除 Token 并跳转登录页
function clearTokensAndRedirect() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  ElMessage.error('登录已过期，请重新登录')

  const redirect = encodeURIComponent(
    `${window.location.pathname}${window.location.search}${window.location.hash}`
  )
  window.location.href = `/login?redirect=${redirect}`
}

export default client
