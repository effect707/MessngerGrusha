import { apiRequest } from './client'
import type { Chat } from '../types'

export const chatsApi = {
  createDirect(recipientId: string) {
    return apiRequest<{ chat: Chat }>('/api/chats/direct', {
      method: 'POST',
      body: JSON.stringify({ recipient_id: recipientId }),
    })
  },

  createGroup(name: string, memberIds: string[]) {
    return apiRequest<{ chat: Chat }>('/api/chats/group', {
      method: 'POST',
      body: JSON.stringify({ name, member_ids: memberIds }),
    })
  },

  getChat(chatId: string) {
    return apiRequest<{ chat: Chat }>(`/api/chats/${chatId}`)
  },

  getUserChats() {
    return apiRequest<{ chats: Chat[] }>('/api/chats/mine')
  },

  addMember(chatId: string, userId: string) {
    return apiRequest<object>(`/api/chats/${chatId}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  },

  removeMember(chatId: string, userId: string) {
    return apiRequest<object>(`/api/chats/${chatId}/members/${userId}`, {
      method: 'DELETE',
    })
  },
}
