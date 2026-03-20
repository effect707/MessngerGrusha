import { apiRequest } from './client'
import type { Message, Reaction, Attachment } from '../types'

interface HistoryResponse {
  messages: Message[]
  has_more: boolean
  next_cursor_id?: string
  next_cursor_created_at?: string
}

export const messagesApi = {
  send(chatId: string, type: string, content: string, replyToId?: string) {
    return apiRequest<{ message: Message }>('/api/messages/send', {
      method: 'POST',
      body: JSON.stringify({ chat_id: chatId, type, content, reply_to_id: replyToId }),
    })
  },

  getHistory(chatId: string, limit: number, cursorId?: string, cursorCreatedAt?: string) {
    const params = new URLSearchParams({ chat_id: chatId, limit: String(limit) })
    if (cursorId) params.set('cursor_id', cursorId)
    if (cursorCreatedAt) params.set('cursor_created_at', cursorCreatedAt)
    return apiRequest<HistoryResponse>(`/api/messages/history?${params}`)
  },

  search(chatId: string, query: string, limit: number) {
    const params = new URLSearchParams({ chat_id: chatId, query, limit: String(limit) })
    return apiRequest<{ messages: Message[] }>(`/api/messages/search?${params}`)
  },

  addReaction(messageId: string, emoji: string) {
    return apiRequest<object>(`/api/messages/${messageId}/reactions`, {
      method: 'POST',
      body: JSON.stringify({ emoji }),
    })
  },

  removeReaction(messageId: string, emoji: string) {
    const params = new URLSearchParams({ emoji })
    return apiRequest<object>(`/api/messages/${messageId}/reactions?${params}`, {
      method: 'DELETE',
    })
  },

  getReactions(messageId: string) {
    return apiRequest<{ reactions: Reaction[] }>(`/api/messages/${messageId}/reactions`)
  },

  getAttachments(messageId: string) {
    return apiRequest<{ attachments: Attachment[] }>(`/api/messages/${messageId}/attachments`)
  },
}
