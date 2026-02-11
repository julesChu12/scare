import client from './client'

export interface Notification {
    id: number
    title: string
    content: string
    type: string
    is_read: boolean
    created_at: string
}

export interface NotificationListResponse {
    list: Notification[]
    total: number
    unread_count: number
}

export const notificationAPI = {
    /**
     * 获取通知列表
     * GET /api/v1/c/notifications
     */
    getList: (params: { page?: number; page_size?: number } = {}) => {
        return client.get<any, NotificationListResponse>('/c/notifications', { params })
    },

    /**
     * 标记通知为已读
     * POST /api/v1/c/notifications/:id/read
     */
    markAsRead: (id: number) => {
        return client.post<any, void>(`/c/notifications/${id}/read`)
    }
}
