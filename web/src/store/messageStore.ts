import { create } from 'zustand'
import type { Message } from '../types'
import { messagesApi } from '../api/messages'

interface Cursor {
  id: string
  createdAt: string
}

interface MessageState {
  messages: Record<string, Message[]>
  cursors: Record<string, Cursor>
  hasMore: Record<string, boolean>
  loaded: Record<string, boolean>
  typing: Record<string, string[]>
  fetchHistory: (chatId: string) => Promise<void>
  addMessage: (message: Message) => void
  setTyping: (chatId: string, userId: string) => void
  clearTyping: (chatId: string, userId: string) => void
}

export const useMessageStore = create<MessageState>((set, get) => ({
  messages: {},
  cursors: {},
  hasMore: {},
  loaded: {},
  typing: {},

  async fetchHistory(chatId) {
    const state = get()

    if (!state.cursors[chatId] && state.loaded[chatId]) return

    const cursor = state.cursors[chatId]
    const res = await messagesApi.getHistory(chatId, 30, cursor?.id, cursor?.createdAt)
    const incoming = res.messages || []
    set((s) => {
      const existing = s.messages[chatId] || []
      const existingIds = new Set(existing.map((m) => m.id))
      const dedupedIncoming = incoming.filter((m: Message) => !existingIds.has(m.id))
      return {
      messages: {
        ...s.messages,
        [chatId]: [...existing, ...dedupedIncoming],
      },
      hasMore: { ...s.hasMore, [chatId]: res.has_more },
      loaded: { ...s.loaded, [chatId]: true },
      cursors: res.next_cursor_id
        ? {
            ...s.cursors,
            [chatId]: { id: res.next_cursor_id, createdAt: res.next_cursor_created_at! },
          }
        : s.cursors,
    }})
  },

  addMessage(message) {
    set((s) => {
      const existing = s.messages[message.chat_id] || []
      if (existing.some((m) => m.id === message.id)) return s
      return {
        messages: {
          ...s.messages,
          [message.chat_id]: [message, ...existing],
        },
      }
    })
  },

  setTyping(chatId, userId) {
    set((s) => {
      const current = s.typing[chatId] || []
      if (current.includes(userId)) return s
      return { typing: { ...s.typing, [chatId]: [...current, userId] } }
    })

    setTimeout(() => get().clearTyping(chatId, userId), 3000)
  },

  clearTyping(chatId, userId) {
    set((s) => ({
      typing: {
        ...s.typing,
        [chatId]: (s.typing[chatId] || []).filter((id) => id !== userId),
      },
    }))
  },
}))
