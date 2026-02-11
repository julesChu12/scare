import client from './client'
import type { ServiceTypeKey } from '@/config/serviceTypes'

// 请求状态枚举
export type RequestStatus = 'pending' | 'processing' | 'completed' | 'cancelled'

// 请求状态映射
export const REQUEST_STATUS_MAP: Record<RequestStatus, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'warning' },
  processing: { label: '处理中', color: 'primary' },
  completed: { label: '已完成', color: 'success' },
  cancelled: { label: '已取消', color: 'info' }
}

// 接口类型定义
export interface CreateRequestRequest {
  service_type: ServiceTypeKey | string
  description?: string
  address: string
  latitude?: number
  longitude?: number
  contact_name: string
  contact_phone: string
}

export interface CreateRequestResponse {
  id: number
  request_no: string
  service_type: string
  status: RequestStatus
  created_at: string
}

export interface AssignedStaff {
  id: number
  name: string
  phone: string
}

export interface ServiceRequest {
  id: number
  request_no: string
  service_type: string
  status: RequestStatus
  description?: string
  address: string
  contact_name: string
  contact_phone: string
  assigned_staff?: AssignedStaff
  // 扩展字段（详情页可能返回）
  station?: {
    id: number
    name: string
  }
  appointment_time?: string
  reject_reason?: string
  rating?: number
  comment?: string
  created_at: string
  updated_at: string
}

export interface PageData<T> {
  list: T
  total: number
  page: number
  page_size: number
}

export interface RateRequestData {
  rating: number
  comment?: string
}

export interface CancelRequestData {
  reason?: string
}

// API 函数
export const requestsAPI = {
  /**
   * 创建服务请求
   * POST /api/v1/c/requests
   */
  createRequest: (data: CreateRequestRequest) => {
    return client.post<any, CreateRequestResponse>('/c/requests', data)
  },

  /**
   * 获取我的服务请求列表
   * GET /api/v1/c/requests
   */
  getMyRequests: (params?: { page?: number; page_size?: number; status?: string }) => {
    return client.get<any, PageData<ServiceRequest[]>>('/c/requests', { params })
  },

  /**
   * 获取服务请求详情
   * GET /api/v1/c/requests/:id
   */
  getRequestDetail: (id: number) => {
    return client.get<any, ServiceRequest>(`/c/requests/${id}`)
  },

  /**
   * 取消服务请求
   * POST /api/v1/c/requests/:id/cancel
   */
  cancelRequest: (id: number, data?: CancelRequestData) => {
    return client.post<any, void>(`/c/requests/${id}/cancel`, data || {})
  },

  /**
   * 评价服务
   * POST /api/v1/c/requests/:id/rate
   */
  rateRequest: (id: number, data: RateRequestData) => {
    return client.post<any, void>(`/c/requests/${id}/rate`, data)
  }
}
