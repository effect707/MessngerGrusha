import { apiRequest } from './client'
import type { Channel } from '../types'

export const channelsApi = {
  create(slug: string, name: string, description: string, isPrivate: boolean) {
    return apiRequest<{ channel: Channel }>('/api/channels', {
      method: 'POST',
      body: JSON.stringify({ slug, name, description, is_private: isPrivate }),
    })
  },

  getChannel(channelId: string) {
    return apiRequest<{ channel: Channel }>(`/api/channels/${channelId}`)
  },

  update(channelId: string, name: string, description: string, isPrivate: boolean) {
    return apiRequest<{ channel: Channel }>(`/api/channels/${channelId}`, {
      method: 'PUT',
      body: JSON.stringify({ name, description, is_private: isPrivate }),
    })
  },

  delete(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}`, { method: 'DELETE' })
  },

  subscribe(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}/subscribe`, { method: 'POST' })
  },

  unsubscribe(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}/subscribe`, { method: 'DELETE' })
  },

  getPublic(limit: number) {
    return apiRequest<{ channels: Channel[] }>(`/api/channels/public?limit=${limit}`)
  },

  getMine() {
    return apiRequest<{ channels: Channel[] }>('/api/channels/mine')
  },
}
