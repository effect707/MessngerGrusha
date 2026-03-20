import { create } from 'zustand'

type SidebarTab = 'chats' | 'channels'

interface UIState {
  sidebarTab: SidebarTab
  searchQuery: string
  activeModal: string | null
  setSidebarTab: (tab: SidebarTab) => void
  setSearchQuery: (q: string) => void
  openModal: (modal: string) => void
  closeModal: () => void
}

export const useUIStore = create<UIState>((set) => ({
  sidebarTab: 'chats',
  searchQuery: '',
  activeModal: null,

  setSidebarTab(tab) { set({ sidebarTab: tab }) },
  setSearchQuery(q) { set({ searchQuery: q }) },
  openModal(modal) { set({ activeModal: modal }) },
  closeModal() { set({ activeModal: null }) },
}))
