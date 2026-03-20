import { create } from 'zustand'
import type { Chat } from '../types'
import { chatsApi } from '../api/chats'

interface ChatState {
  chats: Chat[]
  activeChatId: string | null
  fetchChats: () => Promise<void>
  setActiveChat: (chatId: string | null) => void
  createDirect: (recipientId: string) => Promise<Chat>
  createGroup: (name: string, memberIds: string[]) => Promise<Chat>
}

export const useChatStore = create<ChatState>((set) => ({
  chats: [],
  activeChatId: null,

  async fetchChats() {
    const res = await chatsApi.getUserChats()
    set({ chats: res.chats || [] })
  },

  setActiveChat(chatId) {
    set({ activeChatId: chatId })
  },

  async createDirect(recipientId) {
    const res = await chatsApi.createDirect(recipientId)
    set((s) => {
      const exists = s.chats.some((c) => c.id === res.chat.id)
      return {
        chats: exists ? s.chats : [res.chat, ...s.chats],
        activeChatId: res.chat.id,
      }
    })
    return res.chat
  },

  async createGroup(name, memberIds) {
    const res = await chatsApi.createGroup(name, memberIds)
    set((s) => {
      const exists = s.chats.some((c) => c.id === res.chat.id)
      return {
        chats: exists ? s.chats : [res.chat, ...s.chats],
        activeChatId: res.chat.id,
      }
    })
    return res.chat
  },
}))
