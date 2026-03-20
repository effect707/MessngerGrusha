import { create } from 'zustand'
import type { Channel } from '../types'
import { channelsApi } from '../api/channels'

interface ChannelState {
  channels: Channel[]
  publicChannels: Channel[]
  activeChannelId: string | null
  fetchMyChannels: () => Promise<void>
  fetchPublicChannels: () => Promise<void>
  setActiveChannel: (id: string | null) => void
  subscribe: (channelId: string) => Promise<void>
  unsubscribe: (channelId: string) => Promise<void>
}

export const useChannelStore = create<ChannelState>((set) => ({
  channels: [],
  publicChannels: [],
  activeChannelId: null,

  async fetchMyChannels() {
    const res = await channelsApi.getMine()
    set({ channels: res.channels || [] })
  },

  async fetchPublicChannels() {
    const res = await channelsApi.getPublic(50)
    set({ publicChannels: res.channels || [] })
  },

  setActiveChannel(id) {
    set({ activeChannelId: id })
  },

  async subscribe(channelId) {
    await channelsApi.subscribe(channelId)
  },

  async unsubscribe(channelId) {
    await channelsApi.unsubscribe(channelId)
    set((s) => ({ channels: s.channels.filter((c) => c.id !== channelId) }))
  },
}))
