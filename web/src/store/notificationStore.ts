import { create } from 'zustand'
import type { Notification } from '../types'
import { notificationsApi } from '../api/notifications'

interface NotificationState {
  notifications: Notification[]
  unreadCount: number
  fetchNotifications: (limit?: number) => Promise<void>
  fetchUnreadCount: () => Promise<void>
  markRead: (id: string) => Promise<void>
  markAllRead: () => Promise<void>
}

export const useNotificationStore = create<NotificationState>((set) => ({
  notifications: [],
  unreadCount: 0,

  async fetchNotifications(limit = 20) {
    const res = await notificationsApi.getAll(limit)
    set({ notifications: res.notifications || [] })
  },

  async fetchUnreadCount() {
    const res = await notificationsApi.getUnreadCount()
    set({ unreadCount: Number(res.count) || 0 })
  },

  async markRead(id) {
    await notificationsApi.markRead(id)
    set((s) => ({
      notifications: s.notifications.map((n) => (n.id === id ? { ...n, is_read: true } : n)),
      unreadCount: Math.max(0, s.unreadCount - 1),
    }))
  },

  async markAllRead() {
    await notificationsApi.markAllRead()
    set((s) => ({
      notifications: s.notifications.map((n) => ({ ...n, is_read: true })),
      unreadCount: 0,
    }))
  },
}))
