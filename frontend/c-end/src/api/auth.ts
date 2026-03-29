import client from './client'

// 接口类型定义
export interface SendCodeRequest {
  phone: string
}

export interface SendCodeResponse {
  message: string
}

export interface QuickStartRequest {
  phone: string
  code: string
  name: string
  address: string
  latitude?: number
  longitude?: number
  service_type: string
  description?: string
  images?: string[]
  contact_name?: string
  contact_phone?: string
}

export interface QuickStartResponse {
  token: string
  refresh_token: string
  user: {
    id: number
    phone: string
    role: string
  }
  profile: {
    name: string
    address?: string
    user_type?: string
  }
  request: {
    id: number
    request_no: string
    service_type: string
    status: string
    contact_name?: string
    contact_phone?: string
    images?: string
    created_at: string
  }
}

export interface LoginRequest {
  phone: string
  password?: string
  code?: string
  type?: 'password' | 'code'
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user_id: number
  type: string
  customer_type: string
  name: string
  phone: string
  status: string
}

export interface RefreshRequest {
  refresh_token: string
}

export interface RefreshResponse {
  token: string
  refresh_token: string
}

export interface CheckTokenResponse {
  user: {
    id: number
    phone: string
    role: string
  }
  profile?: {
    name: string
    id_number?: string
    address?: string
    user_type?: string
  } | null
}

export interface MeResponse {
  user_id: number
  type: string
  customer_type: string
  name: string
  phone: string
  status: string
}

// API 函数
export const authAPI = {
  // 发送验证码
  sendCode: (data: SendCodeRequest) => {
    return client.post<any, SendCodeResponse>('/c/auth/send-code', data)
  },

  // 快速开通（注册+登录+创建服务请求）
  quickStart: (data: QuickStartRequest) => {
    return client.post<any, QuickStartResponse>('/c/auth/quick-start', data)
  },

  // 登录（支持密码和验证码两种方式）
  login: (data: LoginRequest) => {
    return client.post<any, LoginResponse>('/c/auth/login', data)
  },

  // 刷新 Token
  refresh: (data: RefreshRequest) => {
    return client.post<any, RefreshResponse>('/c/auth/refresh', data)
  },

  // 检查 Token 状态（用于预填充）
  checkToken: () => {
    return client.get<any, CheckTokenResponse>('/c/auth/check')
  },

  // 获取当前用户信息
  getMe: () => {
    return client.get<any, MeResponse>('/c/auth/me')
  },

  // 退出登录
  logout: () => {
    return client.post<any, void>('/c/auth/logout')
  }
}
