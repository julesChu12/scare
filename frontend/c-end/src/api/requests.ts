import client from './client'
import type { ServiceTypeKey } from '@/config/serviceTypes'

// 请求状态枚举
export type RequestStatus = 'pending' | 'dispatched' | 'claimed' | 'processing' | 'completed' | 'cancelled' | 'rejected'

// 请求状态映射
export const REQUEST_STATUS_MAP: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'warning' },
  dispatched: { label: '已派单', color: 'primary' },
  claimed: { label: '已认领', color: 'primary' },
  processing: { label: '服务中', color: 'primary' },
  completed: { label: '已完成', color: 'success' },
  cancelled: { label: '已取消', color: 'info' },
  rejected: { label: '已拒绝', color: 'danger' }
}

// 创建请求（对齐后端 handler.requestCreate）
export interface CreateRequestRequest {
  service_type: ServiceTypeKey | string
  address?: string
  lat?: number
  lng?: number
  contact_name?: string
  contact_phone?: string
  images?: string[]
}

// 创建响应（后端返回完整 RequestResponse）
export interface CreateRequestResponse {
  id: number
  request_no: string
  service_type: string
  status: string
  created_at: string
}

// 服务请求详情（对齐后端 handler.RequestResponse）
export interface ServiceRequest {
  id: number
  request_no: string
  user_id: number
  service_type: string
  status: string
  description?: string
  contact_name: string
  contact_phone: string
  address: string
  station_id?: number
  station_name?: string
  appointment_time?: string
  urgency?: string
  reject_reason?: string
  images?: string
  rating?: number
  feedback?: string
  created_at: string
  updated_at: string
}

// 列表响应（对齐后端 handler.RequestListResponse）
export interface PageData<T> {
  items: T
  total: number
  page: number
  page_size: number
}

// 评价请求（对齐后端 handler.rateRequest）
export interface RateRequestData {
  rating: number
  feedback?: string
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
