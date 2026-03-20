import { apiRequest } from './client'
import type { User } from '../types'

export const usersApi = {
  getProfile(userId: string) {
    return apiRequest<{ user: User }>(`/api/users/${userId}`)
  },

  getOnlineStatus(userIds: string[]) {
    return apiRequest<{ statuses: Record<string, boolean> }>('/api/users/online-status', {
      method: 'POST',
      body: JSON.stringify({ user_ids: userIds }),
    })
  },
}
