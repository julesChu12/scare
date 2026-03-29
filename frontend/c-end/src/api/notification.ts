import client from './client'
import type { PageData } from './requests'

export interface Notification {
  id: number
  user_id: number
  title: string
  body: string
  type: string
  is_read: boolean
  created_at: string
}

export const notificationAPI = {
  /**
   * 获取通知列表
   * GET /api/v1/c/notifications
   */
  getList: (params?: { page?: number; page_size?: number }) => {
    return client.get<any, PageData<Notification[]>>('/c/notifications', { params })
  },

  /**
   * 标记通知已读
   * POST /api/v1/c/notifications/:id/read
   */
  markRead: (id: number) => {
    return client.post<any, void>(`/c/notifications/${id}/read`)
  }
}
