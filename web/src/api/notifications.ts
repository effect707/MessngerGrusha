import { apiRequest } from './client'
import type { Notification } from '../types'

export const notificationsApi = {
  getAll(limit: number, unreadOnly: boolean = false) {
    const params = new URLSearchParams({ limit: String(limit) })
    if (unreadOnly) params.set('unread_only', 'true')
    return apiRequest<{ notifications: Notification[] }>(`/api/notifications?${params}`)
  },

  markRead(notificationId: string) {
    return apiRequest<object>(`/api/notifications/${notificationId}/read`, { method: 'POST' })
  },

  markAllRead() {
    return apiRequest<object>('/api/notifications/read-all', { method: 'POST' })
  },

  getUnreadCount() {
    return apiRequest<{ count: number }>('/api/notifications/unread-count')
  },
}
